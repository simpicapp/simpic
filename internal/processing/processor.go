package processing

import (
	"bytes"
	"fmt"
	uuid "github.com/satori/go.uuid"
	"github.com/simpicapp/simpic/internal"
	"github.com/simpicapp/simpic/internal/storage"
	"io"
	"io/ioutil"
	"log"
	"os"
)

type PhotoStore interface {
	Read(id uuid.UUID, kind storage.StoreKind) (io.ReadCloser, error)
	Write(id uuid.UUID, kind storage.StoreKind) (*os.File, error)
	Size(id uuid.UUID, kind storage.StoreKind) int64
	Delete(id uuid.UUID, kind storage.StoreKind) error
}

type Processor struct {
	db              *internal.Database
	store           PhotoStore
	thumbnailHeight uint
	screenHeight    uint
}

func NewProcessor(db *internal.Database, store PhotoStore, thumbnailHeight, screenHeight uint) *Processor {
	m := &Processor{
		db:              db,
		store:           store,
		thumbnailHeight: thumbnailHeight,
		screenHeight:    screenHeight,
	}

	return m
}

func (p *Processor) MigrateAll() {
	photos, err := p.db.GetPhotosByProcessedLevel(len(migrations))
	if err != nil {
		log.Printf("Unable to get photos to be migrated: %v\n", err)
		return
	}

	log.Printf("%d photos need migrating\n", len(photos))

	for i, photo := range photos {
		b, err := p.bytes(photo.Id)
		if err != nil {
			log.Printf("Failed to read photo %s: %v\n", photo.Id, err)
			continue
		}

		if err := p.performActions(&photo, b, p.migrationFrom(photo.Processed)); err != nil {
			log.Printf("Failed to migrate photo %s: %v\n", photo.Id, err)
			continue
		}

		photo.Processed = len(migrations)
		if err := p.db.UpdatePhoto(&photo); err != nil {
			log.Printf("Failed to update photo %s after migrating: %v\n", photo.Id, err)
			continue
		}

		if i%20 == 0 {
			log.Printf("%d photos migrated...\n", i)
		}
	}

	log.Printf("Migration of %d photos completed\n", len(photos))
}

func (p *Processor) Process(photo *internal.Photo, reader io.Reader) error {
	b, err := ioutil.ReadAll(reader)
	if err != nil {
		return err
	}

	if err := p.performActions(photo, b, initialActions); err != nil {
		return err
	}

	photo.Processed = len(migrations)
	return p.db.UpdatePhoto(photo)
}

func (p *Processor) bytes(id uuid.UUID) ([]byte, error) {
	reader, err := p.store.Read(id, storage.KindRaw)
	if err != nil {
		return nil, err
	}

	defer func() {
		_ = reader.Close()
	}()

	return ioutil.ReadAll(reader)
}

func (p *Processor) migrationFrom(oldLevel int) action {
	actions := actionNoop
	for i := oldLevel; i < len(migrations); i++ {
		actions |= migrations[i]
	}
	return actions
}

func (p *Processor) performActions(photo *internal.Photo, b []byte, actions action) error {
	for _, action := range allActions {
		if action&actions == action {
			err := p.performAction(photo, b, action)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func (p *Processor) performAction(photo *internal.Photo, b []byte, a action) error {
	switch a {
	case actionSaveRaw:
		return p.saveRaw(photo.Id, bytes.NewReader(b))
	case actionGenerateSamples:
		return p.generateSamples(photo, b)
	case actionExtractExif:
		return p.extractExif(photo.Id, bytes.NewReader(b))
	default:
		return fmt.Errorf("unknown action %d", a)
	}
}
