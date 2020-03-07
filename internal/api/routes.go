package api

import (
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"net/http"
)

func (s *server) routes() http.Handler {
	r := createRouter()
	r.Use(s.authenticatedContext, s.provideVersion)

	r.Route("/api", s.apiRoutes)

	r.Route("/data", func(r chi.Router) {
		r.Group(func(r chi.Router) {
			r.Use(s.photoContext, s.cacheContext)
			r.Get("/{uuid}/{purpose}.{format}", s.handleGetData())
		})
	})

	r.Route("/", s.frontendRoutes)
	r.NotFound(http.FileServer(http.Dir(*frontendDir)).ServeHTTP)
	return r
}

func (s *server) apiRoutes(r chi.Router) {
	r.Post("/login", s.handleAuthenticate())
	r.Get("/timeline", s.handleTimeline())

	r.Route("/albums", s.albumRoutes)
	r.Route("/photos", s.photoRoutes)

	r.Group(func(r chi.Router) {
		r.Use(s.requireAnyUser)
		r.Get("/users/me", s.handleGetSelf())
		r.Get("/logout", s.handleLogout())
	})
}

func (s *server) photoRoutes(r chi.Router) {
	r.Group(func(r chi.Router) {
		r.Use(s.requireAnyUser)
		r.Post("/", s.handleStorePhoto())
		r.Post("/delete", s.handleDeletePhotos())
		r.Post("/update", s.handleUpdatePhotos())
	})

	r.Group(func(r chi.Router) {
		r.Use(s.photoContext)
		r.Get("/{uuid}", s.handleGetPhotoInfo())
	})
}

func (s *server) albumRoutes(r chi.Router) {
	r.Get("/", s.handleGetAlbums())

	r.Group(func(r chi.Router) {
		r.Use(s.requireAnyUser)
		r.Post("/", s.handleAddAlbum())
	})

	r.Group(func(r chi.Router) {
		r.Use(s.albumContext)
		r.Get("/{uuid}", s.handleGetAlbumInfo())

		r.Group(func(r chi.Router) {
			r.Use(s.requireAnyUser)
			r.Delete("/{uuid}", s.handleDeleteAlbum())
			r.Post("/{uuid}", s.handleUpdateAlbum())
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
}

func (s *server) frontendRoutes(r chi.Router) {
	r.Get("/timeline/", s.handleFrontendPath("Timeline"))
	r.Get("/timeline/photo/{photoId}/", s.handleFrontendPath("Timeline"))
	r.Get("/albums/", s.handleFrontendPath("Albums"))
	r.Get("/albums/{albumId}/", s.handleFrontendPath("Albums"))
	r.Get("/albums/{albumId}/photo/{photoId}/", s.handleFrontendPath("Albums"))
}

func createRouter() *chi.Mux {
	r := chi.NewRouter()
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	return r
}
