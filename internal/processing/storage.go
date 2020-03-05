package processing

import (
	uuid "github.com/satori/go.uuid"
	"github.com/simpicapp/simpic/internal"
	"gopkg.in/gographics/imagick.v3/imagick"
	"io"
)

func (p *Processor) saveRaw(id uuid.UUID, raw io.Reader) error {
	out, err := p.store.Write(id, internal.PurposeOriginal, "")
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

	if err := p.saveResampled(mw, largeWidth, largeHeight, 85, photo.Id, "WEBP", internal.PurposeScreen); err != nil {
		return err
	}

	if err := p.saveResampled(mw, largeWidth, largeHeight, 85, photo.Id, "JPEG", internal.PurposeScreen); err != nil {
		return err
	}

	if err := p.saveResampled(mw, smallWidth, smallHeight, 65, photo.Id, "WEBP", internal.PurposePreview); err != nil {
		return err
	}

	if err := p.saveResampled(mw, smallWidth, smallHeight, 65, photo.Id, "JPEG", internal.PurposePreview); err != nil {
		return err
	}

	return p.db.AddFormat(&internal.Format{
		Photo:   photo.Id,
		Purpose: internal.PurposeOriginal,
		Format:  format,
		Width:   width,
		Height:  height,
		Size:    int64(len(bytes)),
	})
}

func (p *Processor) saveResampled(mw *imagick.MagickWand, width, height, quality uint, id uuid.UUID, format string, purpose internal.FormatPurpose) error {
	if err := mw.ResizeImage(width, height, imagick.FILTER_LANCZOS2_SHARP); err != nil {
		return err
	}

	if err := mw.SetImageCompressionQuality(quality); err != nil {
		return err
	}

	if err := mw.SetImageFormat(format); err != nil {
		return err
	}

	if err := p.write(mw, id, purpose, format); err != nil {
		return err
	}

	return p.db.AddFormat(&internal.Format{
		Photo:   id,
		Purpose: purpose,
		Format:  format,
		Width:   width,
		Height:  height,
		Size:    p.store.Size(id, purpose, format),
	})
}

func (p *Processor) write(mw *imagick.MagickWand, id uuid.UUID, purpose internal.FormatPurpose, format string) error {
	file, err := p.store.Write(id, purpose, format)
	if err != nil {
		return err
	}

	defer func() {
		_ = file.Close()
	}()

	return mw.WriteImageFile(file)
}

func (p *Processor) resizedDimensions(inputWidth, inputHeight, targetHeight uint) (width, height uint) {
	if inputHeight <= targetHeight {
		width = inputWidth
		height = inputHeight
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
