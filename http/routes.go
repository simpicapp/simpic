package http

import (
	"github.com/go-chi/chi"
	"net/http"
)

func (s *server) routes() {
	s.router.Post("/login", s.handleAuthenticate())

	s.router.Group(func(r chi.Router) {
		r.Use(s.authenticatedContext)

		s.router.Group(func(r chi.Router) {
			r.Use(s.requireAnyUser)
			r.Post("/photo", s.handleStorePhoto())
		})

		r.Get("/timeline", s.handleTimeline())
	})

	s.router.Group(func(r chi.Router) {
		r.Use(s.photoContext)
		r.Get("/thumbnail/{uuid}", s.handleGetThumbnail())
		r.Get("/photo/{uuid}", s.handleGetPhoto())
	})

	s.router.Mount("/", http.FileServer(http.Dir(s.staticDir)))
}
