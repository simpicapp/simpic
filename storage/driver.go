package storage

import (
	uuid "github.com/satori/go.uuid"
	"io"
)

type Driver interface {
	Read(id uuid.UUID) (io.ReadCloser, error)
	Write(id uuid.UUID) (io.WriteCloser, error)
}
