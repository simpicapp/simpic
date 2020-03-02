package processing

import (
	"bytes"
	"fmt"
	uuid "github.com/satori/go.uuid"
	"github.com/simpicapp/simpic/internal"
	"github.com/simpicapp/simpic/internal/storage"
	"io"
	"io/ioutil"
)

type PhotoWriter interface {
	Write(id uuid.UUID, kind storage.StoreKind) (io.WriteCloser, error)
	Delete(id uuid.UUID, kind storage.StoreKind) error
}

type context struct {
	db              *internal.Database
	writer          PhotoWriter
	thumbnailHeight int
	screenHeight    int
}

type migration interface {
	migrate(c *context, photo *internal.Photo, raw io.Reader) error
	rollback(c *context, photo *internal.Photo) error
}

type Processor struct {
	context *context
}

var migrations = make(map[int]migration)

const (
	_ = iota
	migrationUpdateType
	migrationSaveRaw
	migrationReadExif
	migrationSaveSampled
)

func NewProcessor(db *internal.Database, writer PhotoWriter, thumbnailHeight, screenHeight int) *Processor {
	m := &Processor{
		context: &context{
			db:              db,
			writer:          writer,
			thumbnailHeight: thumbnailHeight,
			screenHeight:    screenHeight,
		},
	}

	return m
}

func (m *Processor) Migrate(photo *internal.Photo, raw io.Reader) error {
	if photo.Processed >= len(migrations) {
		return nil
	}

	b, err := ioutil.ReadAll(raw)
	if err != nil {
		return err
	}

	r := bytes.NewReader(b)
	for photo.Processed < len(migrations) {
		if _, err := r.Seek(0, 0); err != nil {
			return fmt.Errorf("unable to seek buffer: %v", err)
		}

		if err := migrations[photo.Processed+1].migrate(m.context, photo, r); err != nil {
			m.RollBack(photo)
			return fmt.Errorf("migration %d failed: %v", photo.Processed+1, err)
		}

		photo.Processed++
	}

	return m.context.db.UpdatePhoto(photo)
}

func (m *Processor) RollBack(photo *internal.Photo) {
	for photo.Processed > 0 {
		if err := migrations[photo.Processed].rollback(m.context, photo); err != nil {
			fmt.Printf("Rollback of migration %d failed for photo %s: %v\n", photo.Processed, photo.Id, err)
		}
		photo.Processed--
	}
}
