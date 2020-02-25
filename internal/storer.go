package internal

import (
	"github.com/simpicapp/simpic/internal/storage"
	"io"
)

type Storer struct {
	db     *Database
	driver storage.Driver
}

func NewStorer(db *Database, driver storage.Driver) *Storer {
	return &Storer{
		db:     db,
		driver: driver,
	}
}

func (s *Storer) Store(fileName string, uploader int) (*Photo, io.WriteCloser, error) {
	photo := NewPhoto(fileName)
	photo.Uploader = uploader

	err := s.db.Add(photo)
	if err != nil {
		return nil, nil, err
	}

	writer, err := s.driver.Write(photo.Id)
	if err != nil {
		_ = s.db.DeletePhoto(photo)
		return nil, nil, err
	}

	return photo, writer, err
}
