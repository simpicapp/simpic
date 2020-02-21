package http

import (
	"github.com/go-chi/chi"
	"net/http"
)

func (s *server) routes() {
	s.router.Post("/login", s.handleAuthenticate())

	s.router.Route("/albums", func(r chi.Router) {
		r.Use(s.authenticatedContext)
		r.Get("/", s.handleGetAlbums())

		r.Group(func(r chi.Router) {
			r.Use(s.requireAnyUser)
			r.Post("/", s.handleAddAlbum())
		})

		r.Group(func(r chi.Router) {
			r.Use(s.albumContext)
			r.Get("/{uuid}", s.handleGetAlbum())
		})

		r.Route("/{uuid}/photos", func(r chi.Router) {
			r.Use(s.albumContext)
			r.Get("/", s.handleGetPhotosForAlbum())
			r.Group(func(r chi.Router) {
				r.Use(s.requireAnyUser)
				r.Post("/", s.handleAddPhotosToAlbum())
			})
		})
	})

	s.router.Group(func(r chi.Router) {
		r.Use(s.authenticatedContext)
		r.Get("/timeline", s.handleTimeline())
	})

	s.router.Group(func(r chi.Router) {
		r.Use(s.authenticatedContext, s.requireAnyUser)
		r.Post("/photo", s.handleStorePhoto())
		r.Get("/users/me", s.handleGetSelf())
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
