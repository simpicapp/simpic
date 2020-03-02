package processing

import (
	"fmt"
	"github.com/simpicapp/simpic/internal"
	"io"
)

type photoTypeMigration struct{}

func (*photoTypeMigration) migrate(_ *context, photo *internal.Photo, raw io.Reader) error {
	photo.Type = internal.SniffPhotoType(raw)

	if photo.Type == internal.TypeUnknown {
		return fmt.Errorf("unknown image format: %d", photo.Type)
	}

	return nil
}

func init() {
	migrations[migrationUpdateType] = &photoTypeMigration{}
}
