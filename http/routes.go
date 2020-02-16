package http

import (
	"github.com/go-chi/chi"
	"net/http"
)

func (s *server) routes() {
	s.router.Post("/login", s.handleAuthenticate())

	s.router.Group(func(r chi.Router) {
		r.Use(s.photoCtx)
		r.Get("/thumbnail/{uuid}", s.handleGetThumbnail())
		r.Get("/photo/{uuid}", s.handleGetPhoto())
	})

	s.router.Post("/photo", s.handleStorePhoto())
	s.router.Get("/timeline", s.handleTimeline())

	s.router.Mount("/", http.FileServer(http.Dir(s.staticDir)))
}
