package processing

import (
	"bytes"
	"fmt"
	uuid "github.com/satori/go.uuid"
	"github.com/simpicapp/simpic/internal"
	"io"
	"io/ioutil"
	"log"
	"os"
)

type PhotoStore interface {
	Read(id uuid.UUID, purpose internal.FormatPurpose, format string) (io.ReadCloser, error)
	Write(id uuid.UUID, purpose internal.FormatPurpose, format string) (*os.File, error)
	Size(id uuid.UUID, purpose internal.FormatPurpose, format string) int64
	Delete(id uuid.UUID, purpose internal.FormatPurpose, format string) error
}

type Processor struct {
	db              *internal.Database
	store           PhotoStore
	thumbnailHeight uint
	screenHeight    uint
}

type byteProvider func() ([]byte, error)

type cachingBytesProvider struct {
	cache    []byte
	provider byteProvider
}

func (cp *cachingBytesProvider) bytes() ([]byte, error) {
	if len(cp.cache) > 0 {
		return cp.cache, nil
	}

	b, err := cp.provider()
	if err != nil {
		return nil, err
	}

	cp.cache = b
	return b, nil
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
	if len(photos) == 0 {
		return
	}

	for i, photo := range photos {
		cache := &cachingBytesProvider{
			provider: func() (i []byte, err error) {
				return p.bytes(photo.Id)
			},
		}

		if err := p.performActions(&photo, cache.bytes, p.migrationFrom(photo.Processed)); err != nil {
			log.Printf("Failed to migrate photo %s: %v\n", photo.Id, err)
			continue
		}

		photo.Processed = len(migrations)
		if err := p.db.UpdatePhoto(&photo); err != nil {
			log.Printf("Failed to update photo %s after migrating: %v\n", photo.Id, err)
			continue
		}

		if i%10 == 0 {
			log.Printf("%d photos migrated...\n", i)
		}
	}

	log.Printf("Migration of %d photos completed\n", len(photos))
}

func (p *Processor) Process(photo *internal.Photo, reader io.Reader) error {
	cache := &cachingBytesProvider{
		provider: func() (i []byte, err error) {
			return ioutil.ReadAll(reader)
		},
	}

	if err := p.performActions(photo, cache.bytes, initialActions); err != nil {
		return err
	}

	photo.Processed = len(migrations)
	return p.db.UpdatePhoto(photo)
}

func (p *Processor) bytes(id uuid.UUID) ([]byte, error) {
	reader, err := p.store.Read(id, internal.PurposeOriginal, "")
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

func (p *Processor) performActions(photo *internal.Photo, bp byteProvider, actions action) error {
	for _, action := range allActions {
		if action&actions == action {
			err := p.performAction(photo, bp, action)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func (p *Processor) performAction(photo *internal.Photo, bp byteProvider, a action) (err error) {
	var b []byte

	switch a {
	case actionSaveRaw:
		if b, err = bp(); err == nil {
			err = p.saveRaw(photo.Id, bytes.NewReader(b))
		}
	case actionGenerateSamples:
		if b, err = bp(); err == nil {
			err = p.generateSamples(photo, b)
		}
	case actionExtractExif:
		if b, err = bp(); err == nil {
			err = p.extractExif(photo.Id, bytes.NewReader(b))
		}
	case actionCalculateTimestamp:
		return p.calculateTimestamp(photo)
	default:
		err = fmt.Errorf("unknown action %d", a)
	}

	return
}
