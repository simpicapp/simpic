package http

import (
	"fmt"
	"github.com/csmith/simpic"
	"github.com/gorilla/mux"
	"net/http"
)

type server struct {
	router      *mux.Router
	db          *simpic.Database
	retriever   *simpic.Retriever
	storer      *simpic.Storer
	thumbnailer *simpic.Thumbnailer
}

func Start(db *simpic.Database, thumbnailer *simpic.Thumbnailer, retriever *simpic.Retriever, storer *simpic.Storer, port int) error {
	s := server{
		router:      mux.NewRouter(),
		db:          db,
		retriever:   retriever,
		storer:      storer,
		thumbnailer: thumbnailer,
	}

	s.routes()

	srv := &http.Server{
		Addr:    fmt.Sprintf("0.0.0.0:%d", port),
		Handler: s.router,
	}

	return srv.ListenAndServe()
}
