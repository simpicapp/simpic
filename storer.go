package simpic

import (
	"github.com/simpicapp/simpic/storage"
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

func (s *Storer) Store(fileName string) (*Photo, io.WriteCloser, error) {
	photo := NewPhoto(fileName)

	err := s.db.StorePhoto(photo)
	if err != nil {
		return nil, nil, err
	}

	writer, err := s.driver.Write(photo.Id)
	return photo, writer, err
}
