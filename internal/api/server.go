package api

import (
	"context"
	"flag"
	"fmt"
	"github.com/simpicapp/simpic/internal"
	"github.com/simpicapp/simpic/internal/storage"
	"net/http"
)

var (
	port        = flag.Int("port", 8080, "the port to listen on")
	frontendDir = flag.String("frontend", "dist", "the path to serve frontend files from")
)

type Server interface {
	Start() error
	Stop(ctx context.Context) error
}

type server struct {
	db          *internal.Database
	storer      *internal.Storer
	thumbnailer *internal.Thumbnailer
	usermanager *internal.UserManager
	driver      storage.Driver
	srv         *http.Server
}

func NewServer(db *internal.Database, thumbnailer *internal.Thumbnailer, usermanager *internal.UserManager, driver storage.Driver, storer *internal.Storer) Server {
	s := server{
		db:          db,
		driver:      driver,
		storer:      storer,
		thumbnailer: thumbnailer,
		usermanager: usermanager,
	}

	s.srv = &http.Server{
		Addr:    fmt.Sprintf("0.0.0.0:%d", *port),
		Handler: s.routes(),
	}

	return &s
}

func (s *server) Start() error {
	err := s.srv.ListenAndServe()
	if err == http.ErrServerClosed {
		return nil
	} else {
		return err
	}
}

func (s *server) Stop(ctx context.Context) error {
	return s.srv.Shutdown(ctx)
}
