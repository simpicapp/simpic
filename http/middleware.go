package http

import (
	"context"
	"github.com/go-chi/chi"
	"log"
	"net/http"
	"strings"
)

const ctxPhoto = "photo"
const ctxUser = "user"

func (s *server) photoContext(next http.Handler) http.Handler {
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

func (s *server) authenticatedContext(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		header := r.Header.Get("Authorization")
		if strings.HasPrefix(header, "Bearer ") {
			signedToken := strings.TrimPrefix(header, "Bearer ")
			user, err := s.validateToken(signedToken)
			if err != nil {
				log.Printf("client supplied JWT but it was invalid: %v\n", err)
				writeError(w, http.StatusUnauthorized, "Invalid bearer token supplied")
				return
			}

			ctx := context.WithValue(r.Context(), ctxUser, user)
			next.ServeHTTP(w, r.WithContext(ctx))
		} else {
			next.ServeHTTP(w, r)
		}
	})
}

func (s *server) requireAnyUser(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user := r.Context().Value(ctxUser)
		if user == nil {
			writeError(w, http.StatusUnauthorized, "Authorisation required")
			return
		}

		next.ServeHTTP(w, r)
	})
}
