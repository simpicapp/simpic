package processing

import (
	"github.com/rwcarlsen/goexif/exif"
	"github.com/rwcarlsen/goexif/mknote"
	"github.com/rwcarlsen/goexif/tiff"
	uuid "github.com/satori/go.uuid"
	"io"
	"log"
)

func (p *Processor) extractExif(id uuid.UUID, raw io.Reader) error {
	ex, err := exif.Decode(raw)
	if err != nil {
		log.Printf("Unable to read exif data for photo %s: %v\n", id, err)
		return nil
	}

	walker := &exifWalker{
		values: make(map[string]string),
	}
	if err := ex.Walk(walker); err != nil {
		log.Printf("Unable to walk exif tags for photo %s: %v\n", id, err)
	}

	if err := p.db.StoreExifTags(id, walker.values); err != nil {
		return err
	}

	return nil
}

type exifWalker struct {
	values map[string]string
}

func (e *exifWalker) Walk(name exif.FieldName, tag *tiff.Tag) error {
	e.values[string(name)] = tag.String()
	return nil
}

func init() {
	exif.RegisterParsers(mknote.All...)
}
