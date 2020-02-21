package http

import (
	"fmt"
	uuid "github.com/satori/go.uuid"
	"github.com/simpicapp/simpic"
	"log"
	"net/http"
	"strconv"
	"time"
)

func (s *server) handleAddAlbum() http.HandlerFunc {
	type AlbumData struct {
		Name string `json:"name" validate:"required,max=128"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		data := &AlbumData{}
		if !bind(w, r, data) {
			return
		}

		user := r.Context().Value(ctxUser).(*simpic.User)
		album := simpic.NewAlbum(data.Name, user.Id)

		if err := s.db.AddAlbum(album); err != nil {
			log.Printf("Unable to add album '%s': %v\n", album.Name, err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		http.Redirect(w, r, fmt.Sprintf("/albums/%s", album.Id), http.StatusSeeOther)
	}
}

func (s *server) handleGetAlbums() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var (
			offset int
			err    error
		)
		param, ok := r.URL.Query()["offset"]
		if ok && len(param) > 0 {
			offset, _ = strconv.Atoi(param[0])
		}

		albums, err := s.db.GetAlbums(offset, 100)
		if err != nil {
			log.Printf("unable to retrieve albums: %v\n", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		writeJSON(w, http.StatusOK, albums)
	}
}

func (s *server) handleGetPhotosForAlbum() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		album := r.Context().Value(ctxAlbum).(*simpic.Album)
		var offset int
		param, ok := r.URL.Query()["offset"]
		if ok && len(param) > 0 {
			offset, _ = strconv.Atoi(param[0])
		}

		photos, err := s.db.GetAlbumPhotos(album.Id, offset, 100)
		if err != nil {
			log.Printf("unable to retrieve photos for album %s: %v\n", album.Id, err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		writeJSON(w, http.StatusOK, photos)
	}
}

func (s *server) handleAddPhotosToAlbum() http.HandlerFunc {
	type AddPhotosData struct {
		AddedPhotos []uuid.UUID `json:"add_photos"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		data := &AddPhotosData{}
		if !bind(w, r, data) {
			return
		}

		user := r.Context().Value(ctxUser).(*simpic.User)
		album := r.Context().Value(ctxAlbum).(*simpic.Album)
		count, err := s.db.GetAlbumOrderMax(album.Id)
		if err != nil {
			log.Printf("unable to retrieve max order for album %s: %v\n", album.Id, err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		var photos []simpic.AlbumEntry
		for i, photo := range data.AddedPhotos {
			photos = append(photos, simpic.AlbumEntry{
				Photo:   photo,
				Album:   album.Id,
				Creator: user.Id,
				Order:   count + 1 + i,
				Added:   time.Now(),
			})
		}

		err = s.db.AddAlbumPhotos(photos)
		if err != nil {
			log.Printf("unable to add photos to album %s: %v\n", album.Id, err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		if album.Cover == nil {
			album.Cover = &data.AddedPhotos[0]
			err = s.db.UpdateAlbum(album)
			if err != nil {
				log.Printf("unable to update cover for album %s: %v\n", album.Id, err)
			}
		}

		w.WriteHeader(http.StatusNoContent)
	}
}

func (s *server) handleGetAlbum() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		album := r.Context().Value(ctxAlbum).(*simpic.Album)
		writeJSON(w, http.StatusOK, album)
	}
}
