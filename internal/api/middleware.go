package api

import (
	"context"
	"github.com/go-chi/chi"
	uuid "github.com/satori/go.uuid"
	"github.com/simpicapp/simpic/internal"
	"net/http"
)

const (
	ctxAlbum   = "album"
	ctxPhoto   = "photo"
	ctxSession = "session"
	ctxUser    = "user"
)

func (s *server) albumContext(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id, err := uuid.FromString(chi.URLParam(r, "uuid"))
		if err != nil {
			writeError(w, http.StatusNotFound, "No such album")
			return
		}

		album, err := s.db.GetAlbum(id)
		if err != nil {
			writeError(w, http.StatusNotFound, "No such album")
			return
		}

		ctx := context.WithValue(r.Context(), ctxAlbum, album)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (s *server) photoContext(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id, err := uuid.FromString(chi.URLParam(r, "uuid"))
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
		cookie, err := r.Cookie(cookieName)
		if err != nil {
			s.clearAuthCookie(w)
			next.ServeHTTP(w, r)
			return
		}

		session, err := s.db.GetSession(cookie.Value)
		if err != nil {
			s.clearAuthCookie(w)
			next.ServeHTTP(w, r)
			return
		}

		ctx := context.WithValue(context.WithValue(r.Context(), ctxUser, &session.User), ctxSession, &session.Session)
		next.ServeHTTP(w, r.WithContext(ctx))
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

func (s *server) cacheContext(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Cache-Control", "max-age=31536000")
		next.ServeHTTP(w, r)
	})
}

func (s *server) provideVersion(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Server", internal.GetVersionString())
		next.ServeHTTP(w, r)
	})
}
