package http

import (
	"net/http"
)

func (s *server) handleTimeline() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		paginate(w, r, func(offset, count int) (i interface{}, err error) {
			return s.db.GetPhotosByTime(offset, count)
		})
	}
}
