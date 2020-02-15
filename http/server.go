package http

import (
	"fmt"
	"github.com/csmith/simpic"
	"github.com/gorilla/mux"
	"gopkg.in/square/go-jose.v2"
	"net/http"
)

type server struct {
	staticDir   string
	router      *mux.Router
	db          *simpic.Database
	retriever   *simpic.Retriever
	storer      *simpic.Storer
	thumbnailer *simpic.Thumbnailer
	usermanager *simpic.UserManager
	signer      jose.Signer
}

func Start(db *simpic.Database, thumbnailer *simpic.Thumbnailer, usermanager *simpic.UserManager, retriever *simpic.Retriever, storer *simpic.Storer, staticDir string, port int) error {
	s := server{
		router:      mux.NewRouter(),
		db:          db,
		retriever:   retriever,
		storer:      storer,
		thumbnailer: thumbnailer,
		staticDir:   staticDir,
		usermanager: usermanager,
	}

	signer, err := s.createSigner()
	if err != nil {
		panic(fmt.Sprintf("Unable to create JWT signer: %v", err))
	}

	s.signer = signer
	s.routes()

	srv := &http.Server{
		Addr:    fmt.Sprintf("0.0.0.0:%d", port),
		Handler: s.router,
	}

	return srv.ListenAndServe()
}
