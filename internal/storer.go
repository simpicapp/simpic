package internal

import (
	uuid "github.com/satori/go.uuid"
	"github.com/simpicapp/simpic/internal/storage"
	"io"
)

type PhotoWriter interface {
	Write(id uuid.UUID, kind storage.StoreKind) (io.WriteCloser, error)
}

type Storer struct {
	db     *Database
	writer PhotoWriter
}

func NewStorer(db *Database, driver PhotoWriter) *Storer {
	return &Storer{
		db:     db,
		writer: driver,
	}
}

func (s *Storer) Store(fileName string, uploader int) (*Photo, io.WriteCloser, error) {
	photo := NewPhoto(fileName)
	photo.Uploader = uploader

	err := s.db.Add(photo)
	if err != nil {
		return nil, nil, err
	}

	writer, err := s.writer.Write(photo.Id, storage.KindPhoto)
	if err != nil {
		_ = s.db.DeletePhoto(photo)
		return nil, nil, err
	}

	return photo, writer, err
}
