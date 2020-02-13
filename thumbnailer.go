package simpic

import (
	"github.com/csmith/simpic/storage"
	"github.com/disintegration/imaging"
	uuid "github.com/satori/go.uuid"
	"image"
	"io"
)

type Thumbnailer struct {
	driver storage.Driver
	cache  storage.Driver
	height int
}

func NewThumbnailer(driver, cache storage.Driver, height int) *Thumbnailer {
	return &Thumbnailer{
		driver: driver,
		cache:  cache,
		height: height,
	}
}

// Thumbnail provides a reader for a JPEG-encoded thumbnail version of the picture with the given ID.
// If the thumbnail is not cached, it will be generated.
func (t Thumbnailer) Thumbnail(id uuid.UUID) (io.ReadCloser, error) {
	if !t.cache.Exists(id) {
		if err := t.Generate(id); err != nil {
			return nil, err
		}
	}

	return t.cache.Read(id)
}

// Generate creates a new JPEG-encoded thumbnail for the picture with the given ID.
func (t Thumbnailer) Generate(id uuid.UUID) error {
	img, err := t.load(id)
	if err != nil {
		return err
	}

	out := imaging.Resize(img, 0, t.height, imaging.Lanczos)
	return t.save(id, out)
}

func (t Thumbnailer) load(id uuid.UUID) (image.Image, error) {
	reader, err := t.driver.Read(id)
	if err != nil {
		return nil, err
	}

	defer func() {
		_ = reader.Close()
	}()

	return imaging.Decode(reader, imaging.AutoOrientation(true))
}

func (t Thumbnailer) save(id uuid.UUID, img image.Image) error {
	writer, err := t.cache.Write(id)
	if err != nil {
		return err
	}

	defer func() {
		_ = writer.Close()
	}()

	return imaging.Encode(writer, img, imaging.JPEG, imaging.JPEGQuality(80))
}
