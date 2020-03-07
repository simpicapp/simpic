package api

import (
	"flag"
	"fmt"
	"github.com/go-chi/chi"
	uuid "github.com/satori/go.uuid"
	"github.com/simpicapp/simpic/internal"
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"
)

var (
	defaultVisibility = flag.String("default-visibility", "public", "Default visibility for newly uploaded photos: public, unlisted or private")
)

func (s *server) handleGetPhotoInfo() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		photo := r.Context().Value(ctxPhoto).(*internal.Photo)
		formats, err := s.db.GetFormats(photo.Id)
		if err != nil {
			log.Printf("Unable to get formats for photo %s: %v\n", photo.Id, err)
			writeError(w, http.StatusInternalServerError, "Unexpected error")
			return
		}

		writeJSON(w, http.StatusOK, internal.PhotoWithFormats{
			Photo:   photo,
			Formats: formats,
		})
	}
}

func (s *server) handleUpdatePhotos() http.HandlerFunc {
	type UpdatePhotosData struct {
		Ids        []uuid.UUID         `json:"photos"`
		Visibility internal.Visibility `json:"visibility" validate:"required,min=0,max=2"`
	}
	return func(w http.ResponseWriter, r *http.Request) {
		data := &UpdatePhotosData{}
		if !bind(w, r, data) {
			return
		}

		if err := s.db.SetPhotosVisibility(data.Ids, data.Visibility); err != nil {
			log.Printf("Unable to update photos: %v\n", err)
			writeError(w, http.StatusInternalServerError, "Unexpected error")
			return
		}

		w.WriteHeader(http.StatusNoContent)
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
			log.Printf("Unable to delete photos: %v\n", err)
			writeError(w, http.StatusInternalServerError, "Unexpected error")
			return
		}

		for _, id := range data.Ids {
			if err := s.store.DeleteAll(id); err != nil {
				log.Printf("Unable to remove files for photo %s: %v\n", id, err)
			}
		}

		if err := s.db.RefreshMissingCoverImages(); err != nil {
			log.Printf("Unable to refresh cover images: %v\n", err)
		}

		w.WriteHeader(http.StatusNoContent)
	}
}

func (s *server) handleGetData() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		photo := r.Context().Value(ctxPhoto).(*internal.Photo)
		purposeStr := chi.URLParam(r, "purpose")
		purpose, err := strconv.Atoi(purposeStr)
		if err != nil {
			writeError(w, http.StatusBadRequest, "Unparsable purpose")
			return
		}

		format := chi.URLParam(r, "format")

		stream, err := s.store.Read(photo.Id, internal.FormatPurpose(purpose), format)
		if err != nil {
			log.Printf("unable to retrieve photo '%s': %v\n", photo.Id, err)
			writeError(w, http.StatusInternalServerError, "No photo found")
			return
		}

		defer func() {
			_ = stream.Close()
		}()

		if _, dl := r.URL.Query()["download"]; dl {
			fileName := photo.FileName
			if !strings.HasSuffix(strings.ToUpper(fileName), fmt.Sprintf(".%s", format)) {
				fileName = fmt.Sprintf("%s.%s", fileName, format)
			}
			w.Header().Set("Content-Disposition", fmt.Sprintf(`attachment; filename="%s"`, fileName))
		}

		w.Header().Set("Content-Type", mimeTypeFor(format))
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
		photo := internal.NewPhoto(headers.Filename)
		photo.Uploader = user.Id
		photo.Visibility = getDefaultVisibility()

		if err := s.db.AddPhoto(photo); err != nil {
			log.Printf("Unable to save photo '%s': %v\n", headers.Filename, err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		if err := s.processor.Process(photo, file); err != nil {
			_ = s.db.DeletePhoto(photo)
			_ = s.store.DeleteAll(photo.Id)
			log.Printf("unable to create photo '%s': %v\n", headers.Filename, err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		http.Redirect(w, r, fmt.Sprintf("/photos/%s", photo.Id.String()), http.StatusSeeOther)
	}
}

var mimeTypes = map[string]string{
	"ARW":  "image/x-sony-arw",
	"CR2":  "image/x-canon-cr2",
	"CRW":  "image/x-canon-crw",
	"DCR":  "image/x-kodak-dcr",
	"DNG":  "image/x-adobe-dng",
	"GIF":  "image/gif",
	"JPEG": "image/jpeg",
	"MRW":  "image/x-minolta-mrw",
	"NEF":  "image/x-nikon-nef",
	"ORF":  "image/x-olympus-orf",
	"PNG":  "image/png",
	"RAF":  "image/x-fuji-raf",
	"TIFF": "image/tiff",
	"WEBP": "image/webp",
}

func mimeTypeFor(t string) string {
	mt, ok := mimeTypes[strings.ToUpper(t)]
	if ok {
		return mt
	} else {
		log.Printf("No known content type for type %s\n", t)
		return "application/octet-stream"
	}
}

func getDefaultVisibility() internal.Visibility {
	switch *defaultVisibility {
	case "public":
		return internal.VisPublic
	case "unlisted":
		return internal.VisUnlisted
	case "private":
		return internal.VisPrivate
	default:
		log.Printf("Warning: unknown default visibility %s. Setting photos to private.\n", *defaultVisibility)
		return internal.VisPrivate
	}
}
