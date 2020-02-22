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
			log.Printf("Unable to add album '%s': %v\n", album.Id, err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		http.Redirect(w, r, fmt.Sprintf("/albums/%s", album.Id), http.StatusSeeOther)
	}
}

func (s *server) handleDeleteAlbum() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		album := r.Context().Value(ctxAlbum).(*simpic.Album)

		if err := s.db.DeleteAlbum(album); err != nil {
			log.Printf("Unable to delete album '%s': %v\n", album.Id, err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusNoContent)
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

func (s *server) handleAlterPhotosInAlbum() http.HandlerFunc {
	type AlterPhotosData struct {
		AddedPhotos   []uuid.UUID `json:"add_photos"`
		RemovedPhotos []uuid.UUID `json:"remove_photos"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		data := &AlterPhotosData{}
		if !bind(w, r, data) {
			return
		}

		user := r.Context().Value(ctxUser).(*simpic.User)
		album := r.Context().Value(ctxAlbum).(*simpic.Album)

		if len(data.AddedPhotos) > 0 {
			if err := s.addToAlbum(user, album, data.AddedPhotos); err != nil {
				log.Printf("unable to add photos for album %s: %v\n", album.Id, err)
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
		}

		if len(data.RemovedPhotos) > 0 {
			if err := s.db.RemoveAlbumPhotos(album.Id, data.RemovedPhotos); err != nil {
				log.Printf("unable to remove photos for album %s: %v\n", album.Id, err)
				w.WriteHeader(http.StatusInternalServerError)
				return
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

func (s *server) addToAlbum(user *simpic.User, album *simpic.Album, ids []uuid.UUID) error {
	count, err := s.db.GetAlbumOrderMax(album.Id)
	if err != nil {
		return fmt.Errorf("unable to retrieve max order: %v", err)
	}

	var photos []simpic.AlbumEntry
	for i, photo := range ids {
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
		return fmt.Errorf("unable to add photos: %v", err)
	}

	if album.Cover == nil {
		album.Cover = &ids[0]
		err = s.db.UpdateAlbum(album)
		if err != nil {
			log.Printf("unable to update cover for album %s: %v\n", album.Id, err)
		}
	}

	return nil
}
