package processing

import (
	"github.com/rwcarlsen/goexif/exif"
	"github.com/rwcarlsen/goexif/mknote"
	"github.com/rwcarlsen/goexif/tiff"
	uuid "github.com/satori/go.uuid"
	"github.com/simpicapp/simpic/internal"
	"io"
	"log"
)

type exifMigration struct{}

func (*exifMigration) migrate(c *context, photo *internal.Photo, raw io.Reader) error {
	ex, err := exif.Decode(raw)
	if err != nil {
		log.Printf("Unable to read exif data for photo %s: %v\n", photo.Id, err)
		return nil
	}

	walker := &exifWalker{
		db:      c.db,
		photoId: photo.Id,
		values:  make(map[string]string),
	}
	if err := ex.Walk(walker); err != nil {
		log.Printf("Unable to walk exif tags for photo %s: %v\n", photo.Id, err)
	}

	if err := c.db.StoreExifTags(photo.Id, walker.values); err != nil {
		return err
	}

	return nil
}

func (*exifMigration) rollback(_ *context, _ *internal.Photo) error {
	// No need to do anything, foreign key references on the database will tidy up for us
	return nil
}

type exifWalker struct {
	db      *internal.Database
	photoId uuid.UUID
	values  map[string]string
}

func (e *exifWalker) Walk(name exif.FieldName, tag *tiff.Tag) error {
	e.values[string(name)] = tag.String()
	return nil
}

func init() {
	exif.RegisterParsers(mknote.All...)
	migrations[migrationReadExif] = &exifMigration{}
}
