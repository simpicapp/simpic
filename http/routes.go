package http

import (
	"github.com/go-chi/chi"
	"net/http"
)

func (s *server) routes() {
	s.router.Post("/login", s.handleAuthenticate())

	s.router.Group(func(r chi.Router) {
		r.Use(s.authenticatedContext)
		r.Get("/timeline", s.handleTimeline())
	})

	s.router.Group(func(r chi.Router) {
		r.Use(s.authenticatedContext, s.requireAnyUser)
		r.Post("/photo", s.handleStorePhoto())
	})

	s.router.Route("/data", func(r chi.Router) {
		r.Group(func(r chi.Router) {
			r.Use(s.photoContext)
			r.Get("/image/{uuid}", s.handleGetPhoto())
			r.Get("/thumb/{uuid}", s.handleGetThumbnail())
		})
	})

	s.router.Mount("/", http.FileServer(http.Dir(s.staticDir)))
}
