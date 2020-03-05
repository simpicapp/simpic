package api

import (
	"context"
	"flag"
	"fmt"
	uuid "github.com/satori/go.uuid"
	"github.com/simpicapp/simpic/internal"
	"github.com/simpicapp/simpic/internal/processing"
	"io"
	"net/http"
)

var (
	port        = flag.Int("port", 8080, "the port to listen on")
	frontendDir = flag.String("frontend", "dist", "the path to serve frontend files from")
)

type PhotoStore interface {
	Read(id uuid.UUID, purpose internal.FormatPurpose, format string) (io.ReadCloser, error)
	DeleteAll(id uuid.UUID) error
}

type Server interface {
	Start() error
	Stop(ctx context.Context) error
}

type server struct {
	db          *internal.Database
	processor   *processing.Processor
	usermanager *internal.UserManager
	store       PhotoStore
	srv         *http.Server
}

func NewServer(db *internal.Database, usermanager *internal.UserManager, store PhotoStore, processor *processing.Processor) Server {
	s := server{
		db:          db,
		store:       store,
		processor:   processor,
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
