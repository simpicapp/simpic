package http

import (
	"log"
	"net/http"
	"strconv"
)

func (s *server) handleTimeline() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var (
			offset int
			err    error
		)
		param, ok := r.URL.Query()["offset"]
		if ok && len(param) > 0 {
			offset, _ = strconv.Atoi(param[0])
		}

		photos, err := s.db.GetPhotosByTime(offset, 100)
		if err != nil {
			log.Printf("unable to retrieve timeline: %v\n", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		writeJSON(w, http.StatusOK, s.decoratePhotos(photos))
	}
}
