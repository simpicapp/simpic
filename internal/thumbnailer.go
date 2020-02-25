package internal

import (
	"github.com/disintegration/imaging"
	uuid "github.com/satori/go.uuid"
	"github.com/simpicapp/simpic/internal/storage"
	"image"
	"io"
)

type PhotoStore interface {
	Write(id uuid.UUID, kind storage.StoreKind) (io.WriteCloser, error)
	Read(id uuid.UUID, kind storage.StoreKind) (io.ReadCloser, error)
	Exists(id uuid.UUID, kind storage.StoreKind) bool
}

type Thumbnailer struct {
	store  PhotoStore
	height int
}

func NewThumbnailer(store PhotoStore, height int) *Thumbnailer {
	return &Thumbnailer{
		store:  store,
		height: height,
	}
}

// Thumbnail provides a reader for a JPEG-encoded thumbnail version of the picture with the given ID.
// If the thumbnail is not cached, it will be generated.
func (t Thumbnailer) Thumbnail(id uuid.UUID) (io.ReadCloser, error) {
	if !t.store.Exists(id, storage.KindThumbnail) {
		if err := t.Generate(id); err != nil {
			return nil, err
		}
	}

	return t.store.Read(id, storage.KindThumbnail)
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
	reader, err := t.store.Read(id, storage.KindPhoto)
	if err != nil {
		return nil, err
	}

	defer func() {
		_ = reader.Close()
	}()

	return imaging.Decode(reader, imaging.AutoOrientation(true))
}

func (t Thumbnailer) save(id uuid.UUID, img image.Image) error {
	writer, err := t.store.Write(id, storage.KindThumbnail)
	if err != nil {
		return err
	}

	defer func() {
		_ = writer.Close()
	}()

	return imaging.Encode(writer, img, imaging.JPEG, imaging.JPEGQuality(80))
}
