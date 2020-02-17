package http

import (
	"crypto/ecdsa"
	"fmt"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/simpicapp/simpic"
	"github.com/simpicapp/simpic/storage"
	"gopkg.in/square/go-jose.v2"
	"net/http"
)

type server struct {
	staticDir   string
	router      *chi.Mux
	db          *simpic.Database
	storer      *simpic.Storer
	thumbnailer *simpic.Thumbnailer
	usermanager *simpic.UserManager
	driver      storage.Driver
	signer      jose.Signer
	jwtKey      *ecdsa.PrivateKey
	key         [32]byte
}

func Start(db *simpic.Database, thumbnailer *simpic.Thumbnailer, usermanager *simpic.UserManager, driver storage.Driver, storer *simpic.Storer, staticDir string, port int) error {
	s := server{
		router:      createRouter(),
		db:          db,
		driver:      driver,
		storer:      storer,
		thumbnailer: thumbnailer,
		staticDir:   staticDir,
		usermanager: usermanager,
	}

	if err := s.createSigner(); err != nil {
		panic(fmt.Sprintf("Unable to create JWT signer: %v", err))
	}

	s.key = encryptionKey()
	s.routes()

	srv := &http.Server{
		Addr:    fmt.Sprintf("0.0.0.0:%d", port),
		Handler: s.router,
	}

	return srv.ListenAndServe()
}

func createRouter() *chi.Mux {
	r := chi.NewRouter()
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	return r
}
