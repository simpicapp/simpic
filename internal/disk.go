package internal

import (
	"fmt"
	"github.com/satori/go.uuid"
	"io"
	"os"
	"path"
	"strings"
)

// DiskStore stores files on disk according to their UUID.
// Files are grouped into folders according to the last two bytes of the UUID, with each folder being located under
// the given path.
type DiskStore struct {
	Path string
}

func (d DiskStore) Read(id uuid.UUID, purpose FormatPurpose, format string) (io.ReadCloser, error) {
	return os.Open(d.pathFor(id, purpose, format))
}

func (d DiskStore) Write(id uuid.UUID, purpose FormatPurpose, format string) (*os.File, error) {
	file := d.pathFor(id, purpose, format)

	if err := os.MkdirAll(path.Dir(file), os.ModePerm); err != nil {
		return nil, err
	}

	return os.Create(file)
}

func (d DiskStore) Exists(id uuid.UUID, purpose FormatPurpose, format string) bool {
	_, err := os.Stat(d.pathFor(id, purpose, format))
	return err == nil
}

func (d DiskStore) Size(id uuid.UUID, purpose FormatPurpose, format string) int64 {
	stat, err := os.Stat(d.pathFor(id, purpose, format))
	if err != nil {
		return 0
	}
	return stat.Size()
}

func (d DiskStore) Delete(id uuid.UUID, purpose FormatPurpose, format string) error {
	return os.Remove(d.pathFor(id, purpose, format))
}

func (d DiskStore) DeleteAll(id uuid.UUID) error {
	for _, purpose := range []FormatPurpose{PurposePreview, PurposeScreen, PurposeOriginal} {
		for _, format := range []string{"JPEG", "WEBP"} {
			if d.Exists(id, purpose, format) {
				err := d.Delete(id, purpose, format)
				if err != nil {
					return err
				}
			}
		}
	}

	return nil
}

func (d DiskStore) pathFor(id uuid.UUID, purpose FormatPurpose, format string) string {
	var jpegExt = ""
	if strings.ToUpper(format) == "JPEG" {
		jpegExt = ".jpg"
	}

	var webpExt = ""
	if strings.ToUpper(format) == "WEBP" {
		webpExt = ".webp"
	}

	// For legacy reasons different purposes have different suffixes (e.g. it's a jpeg thumbnail is just ".thumb"
	// while its WEBP equivalent is ".thumb.webp"). It would be nice to normalise this at some point in the future.
	var purposeSuffixes = map[FormatPurpose]string{
		PurposePreview:  fmt.Sprintf("thumb%s", webpExt),
		PurposeScreen:   fmt.Sprintf("screen%s%s", jpegExt, webpExt),
		PurposeOriginal: "photo",
	}

	str := id.String()
	return path.Join(d.Path, str[len(str)-2:], fmt.Sprintf("%s.%s", str, purposeSuffixes[purpose]))
}
