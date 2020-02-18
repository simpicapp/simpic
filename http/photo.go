package http

import (
	"crypto/rand"
	"encoding/hex"
	"flag"
	"fmt"
	uuid "github.com/satori/go.uuid"
	"github.com/simpicapp/simpic"
	"golang.org/x/crypto/nacl/secretbox"
	"io"
	"log"
	"net/http"
)

var (
	pictureKey = flag.String("picture-key", "", "32 byte key to use to encode URLs to pictures, as a hexadecimal string")
	nonce      = [24]byte{0x02, 0xf1, 0xf0, 0x4a, 0xc5, 0x99, 0x34, 0xa0, 0xab, 0x30, 0xac, 0x89, 0x8e, 0x92, 0xbf, 0x73, 0xa4, 0xc6, 0x56, 0xc8, 0xea, 0xfc, 0xda, 0xc9}
)

type photoData struct {
	simpic.Photo

	Key       string `json:"key"`
	Url       string `json:"url"`
	Thumbnail string `json:"thumbnail"`
}

func encryptionKey() (key [32]byte) {
	createRandom := func() {
		log.Printf("Warning: invalid picture key specified. This means that picture URLs will change across restarts.\n")
		data := make([]byte, 32)
		_, _ = rand.Read(data)
		copy(key[:], data)
	}

	if len(*pictureKey) < 64 {
		createRandom()
		return
	}

	data, err := hex.DecodeString(*pictureKey)
	if err != nil {
		log.Printf("Error decoding picture key: %v\n", err)
		createRandom()
	} else {
		copy(key[:], data)
	}
	return
}

func (s *server) keyForId(id uuid.UUID) string {
	var output []byte
	out := secretbox.Seal(output, id.Bytes(), &nonce, &s.key)
	return hex.EncodeToString(out)
}

func (s *server) idForKey(key string) (u uuid.UUID, err error) {
	box, err := hex.DecodeString(key)
	if err != nil {
		return
	}

	var output []byte
	out, _ := secretbox.Open(output, box, &nonce, &s.key)
	u, err = uuid.FromBytes(out)
	return
}

func (s *server) decoratePhotos(photos []simpic.Photo) []photoData {
	res := make([]photoData, 0, len(photos))
	for _, photo := range photos {
		res = append(res, s.decoratePhoto(photo))
	}
	return res
}

func (s *server) decoratePhoto(photo simpic.Photo) photoData {
	key := s.keyForId(photo.Id)
	return photoData{
		Photo:     photo,
		Key:       key,
		Url:       fmt.Sprintf("/data/image/%s", key),
		Thumbnail: fmt.Sprintf("/data/thumb/%s", key),
	}
}

func (s *server) handleGetPhoto() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		photo := r.Context().Value(ctxPhoto).(*simpic.Photo)

		stream, err := s.driver.Read(photo.Id)
		if err != nil {
			log.Printf("unable to retrieve photo '%s': %v\n", photo.Id, err)
			writeError(w, http.StatusInternalServerError, "No photo found")
			return
		}

		defer func() {
			_ = stream.Close()
		}()

		w.Header().Set("Content-Type", mimeTypeFor(photo.Type))
		_, _ = io.Copy(w, stream)
	}
}

func (s *server) handleGetThumbnail() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		photo := r.Context().Value(ctxPhoto).(*simpic.Photo)

		stream, err := s.thumbnailer.Thumbnail(photo.Id)
		if err != nil {
			log.Printf("unable to retrieve thumbnail '%s': %v\n", photo.Id, err)
			w.WriteHeader(http.StatusNotFound)
			return
		}

		defer func() {
			_ = stream.Close()
		}()

		w.Header().Set("Content-Type", "image/jpeg")
		_, _ = io.Copy(w, stream)
	}
}

func (s *server) handleStorePhoto() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := r.ParseMultipartForm(32 << 20); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		file, headers, err := r.FormFile("file")
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		defer func() {
			_ = file.Close()
		}()

		user := r.Context().Value(ctxUser).(*simpic.User)
		photo, writer, err := s.storer.Store(headers.Filename, user.Id)
		if err != nil {
			log.Printf("unable to create photo '%s': %v\n", headers.Filename, err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		if _, err := io.Copy(writer, file); err != nil {
			log.Printf("unable to write photo '%s': %v\n", headers.Filename, err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		if err := writer.Close(); err != nil {
			log.Printf("unable to close photo '%s': %v\n", headers.Filename, err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.Header().Set("Location", fmt.Sprintf("/photo/%s", photo.Id.String()))
		w.WriteHeader(http.StatusSeeOther)
	}
}

func mimeTypeFor(t simpic.Type) string {
	switch t {
	case simpic.Jpeg:
		return "image/jpeg"
	case simpic.Png:
		return "image/png"
	default:
		log.Printf("No known content type for type %d\n", t)
		return "application/octet-stream"
	}
}
