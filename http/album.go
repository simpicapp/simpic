package http

import (
	"fmt"
	"github.com/simpicapp/simpic"
	"log"
	"net/http"
	"strconv"
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

		http.Redirect(w, r, fmt.Sprintf("/album/%s", album.Uuid), http.StatusSeeOther)
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
