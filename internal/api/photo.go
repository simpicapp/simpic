package api

import (
	"fmt"
	uuid "github.com/satori/go.uuid"
	"github.com/simpicapp/simpic/internal"
	"github.com/simpicapp/simpic/internal/storage"
	"io"
	"log"
	"net/http"
)

func (s *server) handleGetPhotoInfo() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		photo := r.Context().Value(ctxPhoto).(*internal.Photo)
		writeJSON(w, http.StatusOK, photo)
	}
}

func (s *server) handleDeletePhotos() http.HandlerFunc {
	type DeletePhotosData struct {
		Ids []uuid.UUID `json:"photos"`
	}
	return func(w http.ResponseWriter, r *http.Request) {
		data := &DeletePhotosData{}
		if !bind(w, r, data) {
			return
		}

		if err := s.db.DeletePhotos(data.Ids); err != nil {
			log.Printf("Failed to delete photo batch: %v\n", err)
			writeError(w, http.StatusInternalServerError, "Unexpected error")
			return
		}

		w.WriteHeader(http.StatusNoContent)
	}
}

func (s *server) handleGetPhoto() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		photo := r.Context().Value(ctxPhoto).(*internal.Photo)

		stream, err := s.photoReader.Read(photo.Id, storage.KindPhoto)
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
		photo := r.Context().Value(ctxPhoto).(*internal.Photo)

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
			log.Printf("Failed to parse multipart form: %v\n", err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		file, headers, err := r.FormFile("file")
		if err != nil {
			log.Printf("No file present with key 'file'\n")
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		defer func() {
			_ = file.Close()
		}()

		user := r.Context().Value(ctxUser).(*internal.User)
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

		go func() {
			if err := s.thumbnailer.Generate(photo.Id); err != nil {
				log.Printf("Failed to generate thumbnail for uploaded image %s: %v\n", photo.Id, err)
			}
		}()

		http.Redirect(w, r, fmt.Sprintf("/photos/%s", photo.Id.String()), http.StatusSeeOther)
	}
}

func mimeTypeFor(t internal.Type) string {
	switch t {
	case internal.Jpeg:
		return "image/jpeg"
	case internal.Png:
		return "image/png"
	default:
		log.Printf("No known content type for type %d\n", t)
		return "application/octet-stream"
	}
}
