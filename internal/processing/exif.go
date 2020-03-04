package processing

import (
	uuid "github.com/satori/go.uuid"
	"github.com/simpicapp/goexif/exif"
	"github.com/simpicapp/goexif/mknote"
	"github.com/simpicapp/goexif/tiff"
	"github.com/simpicapp/simpic/internal"
	"io"
	"log"
	"path/filepath"
	"regexp"
	"strings"
	"time"
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

var (
	fileNamePrefixes = regexp.MustCompile("(?i)^(PANO-|Screenshot[-_ ](from )?|IMG[-_])")
	fileNameSuffixes = regexp.MustCompile(`(?i)(\(\d{1,2}\)|[-~]\d|-PANO|-WA\d+|_[a-f0-9]{32})$`)
	fileNameFormats  = []string{
		"20060102 150405",
		"2006 01 02 15 04 05",
		"2006 01 02",
		"20060102",
	}
)

func (p *Processor) calculateTimestamp(photo *internal.Photo) error {
	tag, err := p.db.GetExifTag(photo.Id, "DateTime")
	if err == nil {
		// "2017:04:16 22:05:12"
		format := "2006:01:02 15:04:05"
		t, err := time.Parse(format, strings.Trim(tag.Value, "\""))
		if err == nil {
			photo.TakenTime = &t
			return p.db.UpdatePhoto(photo)
		}
		log.Printf("Unable to parse EXIF DateTime '%s' from photo %s: %v\n", tag.Value, photo.Id, err)
	}

	var fileName string
	fileName = strings.TrimSuffix(photo.FileName, filepath.Ext(photo.FileName))
	fileName = fileNamePrefixes.ReplaceAllString(fileName, "")
	fileName = fileNameSuffixes.ReplaceAllString(fileName, "")
	fileName = strings.ReplaceAll(fileName, "_", " ")
	fileName = strings.ReplaceAll(fileName, "-", " ")
	fileName = strings.ReplaceAll(fileName, ".", " ")

	for _, format := range fileNameFormats {
		t, err := time.Parse(format, fileName)
		if err == nil {
			photo.TakenTime = &t
			return p.db.UpdatePhoto(photo)
		}
	}

	log.Printf("Unable to find time for photo '%s' with filename '%s' (trimmed: '%s')\n", photo.Id, photo.FileName, fileName)
	photo.TakenTime = photo.UploadTime
	return p.db.UpdatePhoto(photo)
}

func init() {
	exif.RegisterParsers(mknote.All...)
}
