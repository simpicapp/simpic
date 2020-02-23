package http

import (
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"net/http"
)

func (s *server) routes() http.Handler {
	r := createRouter()
	r.Use(s.authenticatedContext)
	r.Post("/login", s.handleAuthenticate())
	r.Get("/timeline", s.handleTimeline())

	r.Route("/albums", func(r chi.Router) {
		r.Get("/", s.handleGetAlbums())

		r.Group(func(r chi.Router) {
			r.Use(s.requireAnyUser)
			r.Post("/", s.handleAddAlbum())
		})

		r.Group(func(r chi.Router) {
			r.Use(s.albumContext)
			r.Get("/{uuid}", s.handleGetAlbum())

			r.Group(func(r chi.Router) {
				r.Use(s.requireAnyUser)
				r.Delete("/{uuid}", s.handleDeleteAlbum())
			})
		})

		r.Route("/{uuid}/photos", func(r chi.Router) {
			r.Use(s.albumContext)
			r.Get("/", s.handleGetPhotosForAlbum())
			r.Group(func(r chi.Router) {
				r.Use(s.requireAnyUser)
				r.Post("/", s.handleAlterPhotosInAlbum())
			})
		})
	})

	r.Group(func(r chi.Router) {
		r.Use(s.requireAnyUser)
		r.Post("/photo", s.handleStorePhoto())
		r.Get("/users/me", s.handleGetSelf())
	})

	r.Route("/data", func(r chi.Router) {
		r.Group(func(r chi.Router) {
			r.Use(s.photoContext)
			r.Get("/image/{uuid}", s.handleGetPhoto())
			r.Get("/thumb/{uuid}", s.handleGetThumbnail())
		})
	})

	r.Mount("/", http.FileServer(http.Dir(*frontendDir)))
	return r
}

func createRouter() *chi.Mux {
	r := chi.NewRouter()
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	return r
}
