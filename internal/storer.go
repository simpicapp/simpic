package internal

import (
	"bytes"
	"fmt"
	"github.com/disintegration/imaging"
	uuid "github.com/satori/go.uuid"
	"github.com/simpicapp/simpic/internal/storage"
	"image"
	"image/jpeg"
	"io"
)

type PhotoWriter interface {
	Write(id uuid.UUID, kind storage.StoreKind) (io.WriteCloser, error)
	DeleteAll(id uuid.UUID) error
}

type Storer struct {
	db              *Database
	writer          PhotoWriter
	thumbnailHeight int
}

func NewStorer(db *Database, driver PhotoWriter, thumbnailHeight int) *Storer {
	return &Storer{
		db:              db,
		writer:          driver,
		thumbnailHeight: thumbnailHeight,
	}
}

func (s *Storer) Store(fileName string, uploader int, stream io.Reader) (*Photo, error) {
	photo := NewPhoto(fileName)
	photo.Uploader = uploader

	var buf bytes.Buffer
	img, format, err := image.Decode(io.TeeReader(stream, &buf))
	if err != nil {
		return nil, err
	}

	photo.Width = img.Bounds().Dx()
	photo.Height = img.Bounds().Dx()
	photo.Type = imageTypeForFormat(format)

	if photo.Type == Unknown {
		return nil, fmt.Errorf("unknown image format: %s", format)
	}

	out, err := s.writer.Write(photo.Id, storage.KindPhoto)
	if err != nil {
		return nil, err
	}

	_, err = io.Copy(out, &buf)
	if err != nil {
		_ = out.Close()
		return nil, err
	}

	if err := out.Close(); err != nil {
		return nil, err
	}

	if err := s.storeThumbnail(photo, img); err != nil {
		_ = s.writer.DeleteAll(photo.Id)
		return nil, err
	}

	if err := s.db.Add(photo); err != nil {
		_ = s.writer.DeleteAll(photo.Id)
		return nil, err
	}

	return photo, nil
}

func (s *Storer) storeThumbnail(photo *Photo, img image.Image) error {
	thumb := imaging.Resize(img, 0, s.thumbnailHeight, imaging.Lanczos)
	thumbOut, err := s.writer.Write(photo.Id, storage.KindThumbnail)
	if err != nil {
		return err
	}

	if err := jpeg.Encode(thumbOut, thumb, &jpeg.Options{Quality: 80}); err != nil {
		_ = thumbOut.Close()
		return fmt.Errorf("unable to create thumbnail: %v", err)
	}

	return thumbOut.Close()
}

func (s *Storer) Delete(id uuid.UUID) error {
	return s.writer.DeleteAll(id)
}

func imageTypeForFormat(format string) PhotoType {
	switch format {
	case "jpeg":
		return Jpeg
	case "png":
		return Png
	case "gif":
		return Gif
	default:
		return Unknown
	}
}
