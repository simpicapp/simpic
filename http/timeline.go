package http

import (
	"encoding/json"
	"log"
	"net/http"
)

func (s *server) handleTimeline() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		photos, err := s.db.GetPhotosByTime(0, 100)
		if err != nil {
			log.Printf("unable to retrieve timeline: %v\n", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		_ = json.NewEncoder(w).Encode(photos)
	}
}
