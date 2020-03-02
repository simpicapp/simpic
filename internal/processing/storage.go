package processing

import (
	uuid "github.com/satori/go.uuid"
	"github.com/simpicapp/simpic/internal"
	"github.com/simpicapp/simpic/internal/storage"
	"gopkg.in/gographics/imagick.v3/imagick"
	"io"
)

func (p *Processor) saveRaw(id uuid.UUID, raw io.Reader) error {
	out, err := p.store.Write(id, storage.KindRaw)
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

func (p *Processor) generateSamples(photo *internal.Photo, bytes []byte) error {
	mw := imagick.NewMagickWand()
	defer mw.Destroy()

	if err := mw.ReadImageBlob(bytes); err != nil {
		return err
	}

	if err := mw.AutoOrientImage(); err != nil {
		return err
	}

	format := mw.GetImageFormat()
	width := mw.GetImageWidth()
	height := mw.GetImageHeight()

	largeWidth, largeHeight := p.resizedDimensions(width, height, p.screenHeight)
	smallWidth, smallHeight := p.resizedDimensions(width, height, p.thumbnailHeight)

	if err := p.saveResampled(mw, largeWidth, largeHeight, 95, photo.Id, storage.KindScreenJpeg); err != nil {
		return err
	}

	if err := p.saveResampled(mw, smallWidth, smallHeight, 80, photo.Id, storage.KindThumbnailJpeg); err != nil {
		return err
	}

	photo.Width = largeWidth
	photo.Height = largeHeight
	photo.Format = format
	return p.db.UpdatePhoto(photo)
}

func (p *Processor) saveResampled(mw *imagick.MagickWand, width, height, quality uint, id uuid.UUID, kind storage.StoreKind) error {
	if err := mw.ResizeImage(width, height, imagick.FILTER_LANCZOS2_SHARP); err != nil {
		return err
	}

	if err := mw.SetImageCompressionQuality(quality); err != nil {
		return err
	}

	if err := mw.SetImageFormat("JPEG"); err != nil {
		return err
	}

	file, err := p.store.Write(id, kind)
	if err != nil {
		return err
	}

	defer func() {
		_ = file.Close()
	}()

	if err := mw.WriteImageFile(file); err != nil {
		return err
	}

	return nil
}

func (p *Processor) resizedDimensions(inputWidth, inputHeight, targetHeight uint) (width, height uint) {
	if inputHeight <= targetHeight {
		width = inputWidth
		height = targetHeight
	} else {
		ratio := float32(inputHeight) / float32(targetHeight)
		height = uint(float32(inputHeight) / ratio)
		width = uint(float32(inputWidth) / ratio)
	}

	if width > 5*height {
		ratio := float32(width) / float32(height*5)
		height = uint(float32(height) / ratio)
		width = uint(float32(width) / ratio)
	}

	return
}

func init() {
	imagick.Initialize()
}
