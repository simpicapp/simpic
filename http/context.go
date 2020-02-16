package http

import (
	"context"
	"github.com/go-chi/chi"
	"net/http"
)

const ctxPhoto = "photo"

func (s server) photoCtx(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id, err := s.idForKey(chi.URLParam(r, "uuid"))
		if err != nil {
			writeError(w, http.StatusNotFound, "No such photo")
			return
		}

		photo, err := s.db.GetPhoto(id)
		if err != nil {
			writeError(w, http.StatusNotFound, "No such photo")
			return
		}

		ctx := context.WithValue(r.Context(), ctxPhoto, photo)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
