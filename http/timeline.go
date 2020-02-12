package http

import (
	"encoding/json"
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
			offset, err = strconv.Atoi(param[0])
		}

		photos, err := s.db.GetPhotosByTime(offset, 100)
		if err != nil {
			log.Printf("unable to retrieve timeline: %v\n", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		_ = json.NewEncoder(w).Encode(photos)
	}
}
