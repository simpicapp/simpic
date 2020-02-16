package http

import (
	"fmt"
	"github.com/csmith/simpic"
	"io"
	"log"
	"net/http"
)

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

		photo, writer, err := s.storer.Store(headers.Filename)
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
