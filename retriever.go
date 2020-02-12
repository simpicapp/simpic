package simpic

import (
	"github.com/csmith/simpic/storage"
	uuid "github.com/satori/go.uuid"
	"io"
)

type Retriever struct {
	db     *Database
	driver storage.Driver
}

func NewRetriever(db *Database, driver storage.Driver) *Retriever {
	return &Retriever{
		db:     db,
		driver: driver,
	}
}

func (r *Retriever) Get(id uuid.UUID) (*Photo, io.ReadCloser, error) {
	photo, err := r.db.GetPhoto(id)
	if err != nil {
		return nil, nil, err
	}

	stream, err := r.driver.Read(id)
	return photo, stream, err
}
