package http

import (
	"fmt"
	"github.com/csmith/simpic"
	"github.com/gorilla/mux"
	"net/http"
)

type server struct {
	router    *mux.Router
	db        *simpic.Database
	retriever *simpic.Retriever
	storer    *simpic.Storer
}

func Start(db *simpic.Database, retriever *simpic.Retriever, storer *simpic.Storer, port int) error {
	s := server{
		router:    mux.NewRouter(),
		db:        db,
		retriever: retriever,
		storer:    storer,
	}

	s.routes()

	srv := &http.Server{
		Addr:    fmt.Sprintf("0.0.0.0:%d", port),
		Handler: s.router,
	}

	return srv.ListenAndServe()
}
