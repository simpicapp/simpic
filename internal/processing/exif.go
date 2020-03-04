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
	"strconv"
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

const (
	// minValueForUnixTimestamp is the minimum value we'll treat as a unix timestamp, to avoid interpreting random
	// numbers as dates close to 1970.
	minValueForUnixTimestamp = 561_859_200

	// minJavaStyleUnixTimestamp is the minimum value we'll treat as a Java-style unix timestamp (i.e., including
	// milliseconds).
	minJavaStyleUnixTimestamp = 10_000_000_000
)

var (
	fileNamePrefixes = regexp.MustCompile("(?i)^([a-z]+[-_ ])+")
	fileNameSuffixes = regexp.MustCompile(`(?i)(\(\d{1,2}\)|\s*[-~]\s*\d|[-_][a-z]+\d*|-WA\d+|_[a-f0-9]{32}|[-_ ]\d{3,5})+$`)
	fileNameFormats  = []string{
		"20060102 150405",
		"2006 01 02 15 04 05",
		"2006 01 02 150405",
		"2006 01 02",
		"20060102",
	}
)

func unixTimestampFromFilename(fileName string) *time.Time {
	i, err := strconv.ParseInt(fileName, 10, 64)
	if err == nil && i > minValueForUnixTimestamp {
		if i > minJavaStyleUnixTimestamp {
			// Java-style with milliseconds
			t := time.Unix(i/1000, (i%1000)*1_000_000).UTC()
			return &t
		} else {
			t := time.Unix(i, 0).UTC()
			return &t
		}
	}
	return nil
}

func timeFromFilename(originalName string) *time.Time {
	var fileName string
	fileName = strings.TrimSuffix(originalName, filepath.Ext(originalName))
	fileName = fileNamePrefixes.ReplaceAllString(fileName, "")
	fileName = fileNameSuffixes.ReplaceAllString(fileName, "")
	fileName = strings.ReplaceAll(fileName, "_", " ")
	fileName = strings.ReplaceAll(fileName, "-", " ")
	fileName = strings.ReplaceAll(fileName, ".", " ")

	for _, format := range fileNameFormats {
		t, err := time.Parse(format, fileName)
		if err == nil {
			return &t
		}
	}

	return unixTimestampFromFilename(fileName)
}

func (p *Processor) calculateTimestamp(photo *internal.Photo) error {
	tag, err := p.db.GetExifTag(photo.Id, "DateTime")
	if err == nil && len(tag.Value) > 2 {
		format := "2006:01:02 15:04:05"
		t, err := time.Parse(format, strings.Trim(tag.Value, "\""))
		if err == nil {
			photo.TakenTime = &t
			return p.db.UpdatePhoto(photo)
		}
		log.Printf("Unable to parse EXIF DateTime '%s' from photo %s: %v\n", tag.Value, photo.Id, err)
	}

	t := timeFromFilename(photo.FileName)
	if t == nil {
		log.Printf("Unable to find time for photo '%s' with filename '%s'\n", photo.Id, photo.FileName)
		photo.TakenTime = photo.UploadTime
	} else {
		photo.TakenTime = t
	}

	return p.db.UpdatePhoto(photo)
}

func init() {
	exif.RegisterParsers(mknote.All...)
}
