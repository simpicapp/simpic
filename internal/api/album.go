package api

import (
	"fmt"
	uuid "github.com/satori/go.uuid"
	"github.com/simpicapp/simpic/internal"
	"log"
	"net/http"
	"strings"
	"time"
)

func (s *server) handleAddAlbum() http.HandlerFunc {
	type AlbumData struct {
		Name       string              `json:"name" validate:"required,min=1,max=128"`
		Visibility internal.Visibility `json:"visibility" validate:"required,min=0,max=2"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		data := &AlbumData{}
		if !bind(w, r, data) {
			return
		}

		user := r.Context().Value(ctxUser).(*internal.User)
		album := internal.NewAlbum(data.Name, user.Id, data.Visibility)

		if err := s.db.AddAlbum(album); err != nil {
			if strings.Contains(err.Error(), "albums_album_name_unique") {
				writeError(w, http.StatusUnprocessableEntity, "An album with that name already exists")
			} else {
				log.Printf("Unable to add album '%s': %v\n", album.Id, err)
				w.WriteHeader(http.StatusInternalServerError)
			}
			return
		}

		http.Redirect(w, r, fmt.Sprintf("/albums/%s", album.Id), http.StatusSeeOther)
	}
}

func (s *server) handleDeleteAlbum() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		album := r.Context().Value(ctxAlbum).(*internal.Album)

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
		paginate(w, r, func(offset, count int) (i interface{}, err error) {
			return s.db.GetAlbums(visForBrowsing(r), offset, count)
		})
	}
}

func (s *server) handleGetPhotosForAlbum() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		album := r.Context().Value(ctxAlbum).(*internal.Album)
		paginate(w, r, func(offset, count int) (i interface{}, err error) {
			return s.db.GetAlbumPhotos(album.Id, visForAccess(r), offset, count)
		})
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

		user := r.Context().Value(ctxUser).(*internal.User)
		album := r.Context().Value(ctxAlbum).(*internal.Album)

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

		if err := s.db.RefreshCoverImage(album.Id); err != nil {
			log.Printf("Unable to update cover image for album %s: %v\n", album.Id, err)
		}

		w.WriteHeader(http.StatusNoContent)
	}
}

func (s *server) handleGetAlbumInfo() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		album := r.Context().Value(ctxAlbum).(*internal.Album)
		writeJSON(w, http.StatusOK, album)
	}
}

func (s *server) addToAlbum(user *internal.User, album *internal.Album, ids []uuid.UUID) error {
	count, err := s.db.GetAlbumOrderMax(album.Id)
	if err != nil {
		return fmt.Errorf("unable to retrieve max order: %v", err)
	}

	var photos []internal.AlbumEntry
	for i, photo := range ids {
		photos = append(photos, internal.AlbumEntry{
			Photo:   photo,
			Album:   album.Id,
			Creator: user.Id,
			Order:   count + 1 + i,
			Added:   time.Now(),
		})
	}

	return s.db.AddAlbumPhotos(photos)
}
