package processing

import (
	"fmt"
	"github.com/disintegration/imaging"
	"github.com/rwcarlsen/goexif/exif"
	uuid "github.com/satori/go.uuid"
	"github.com/simpicapp/simpic/internal"
	"github.com/simpicapp/simpic/internal/storage"
	"image"
	"image/jpeg"
	"io"
	"strconv"
)

type saveRawMigration struct{}

func (*saveRawMigration) migrate(c *context, photo *internal.Photo, raw io.Reader) error {
	out, err := c.writer.Write(photo.Id, storage.KindRaw)
	if err != nil {
		return err
	}

	_, err = io.Copy(out, raw)
	if err != nil {
		_ = out.Close()
		return err
	}

	return out.Close()
}

func (*saveRawMigration) rollback(c *context, photo *internal.Photo) error {
	return c.writer.Delete(photo.Id, storage.KindRaw)
}

type saveSampledMigration struct{}

func (s *saveSampledMigration) migrate(c *context, photo *internal.Photo, raw io.Reader) error {
	img, _, err := image.Decode(raw)
	if err != nil {
		return err
	}

	orientation := s.exifOrientation(c, photo)
	if orientation == 5 || orientation == 6 {
		img = imaging.Rotate270(img)
	} else if orientation == 3 || orientation == 4 {
		img = imaging.Rotate180(img)
	} else if orientation == 7 || orientation == 8 {
		img = imaging.Rotate90(img)
	}

	if orientation == 2 || orientation == 4 || orientation == 5 || orientation == 7 {
		img = imaging.FlipH(img)
	}

	_, _, err = s.storeSampled(c.writer, storage.KindThumbnail, c.thumbnailHeight, 80, photo.Id, img)
	if err != nil {
		return err
	}

	width, height, err := s.storeSampled(c.writer, storage.KindScreenJpeg, c.screenHeight, 95, photo.Id, img)
	if err != nil {
		return err
	}

	photo.Width = width
	photo.Height = height
	return c.db.UpdatePhoto(photo)
}

func (*saveSampledMigration) exifOrientation(c *context, photo *internal.Photo) int {
	tag, err := c.db.GetExifTag(photo.Id, string(exif.Orientation))
	if err == nil {
		val, err := strconv.Atoi(tag.Value)
		if err == nil {
			return val
		}
	}
	return 1
}

func (*saveSampledMigration) rollback(c *context, photo *internal.Photo) error {
	_ = c.writer.Delete(photo.Id, storage.KindScreenJpeg)
	_ = c.writer.Delete(photo.Id, storage.KindThumbnail)
	return nil
}

func (*saveSampledMigration) storeSampled(pw PhotoWriter, kind storage.StoreKind, height, quality int, id uuid.UUID, img image.Image) (int, int, error) {
	targetHeight := height
	if img.Bounds().Dy() < targetHeight {
		targetHeight = img.Bounds().Dy()
	}

	thumb := imaging.Resize(img, 0, targetHeight, imaging.Lanczos)
	thumbOut, err := pw.Write(id, kind)
	if err != nil {
		return 0, 0, err
	}

	if err := jpeg.Encode(thumbOut, thumb, &jpeg.Options{Quality: quality}); err != nil {
		_ = thumbOut.Close()
		return 0, 0, fmt.Errorf("unable to create resampled image: %v", err)
	}

	return thumb.Bounds().Dx(), thumb.Bounds().Dy(), thumbOut.Close()
}

func init() {
	migrations[migrationSaveRaw] = &saveRawMigration{}
	migrations[migrationSaveSampled] = &saveSampledMigration{}
}
