package storage

import (
	"fmt"
	"github.com/satori/go.uuid"
	"io"
	"os"
	"path"
)

// DiskStore stores files on disk according to their UUID.
// Files are grouped into folders according to the last two bytes of the UUID, with each folder being located under
// the given path.
type DiskStore struct {
	Path string
}

func (d DiskStore) Read(id uuid.UUID, kind StoreKind) (io.ReadCloser, error) {
	return os.Open(d.pathFor(id, kind))
}

func (d DiskStore) Write(id uuid.UUID, kind StoreKind) (*os.File, error) {
	file := d.pathFor(id, kind)

	if err := os.MkdirAll(path.Dir(file), os.ModePerm); err != nil {
		return nil, err
	}

	return os.Create(file)
}

func (d DiskStore) Exists(id uuid.UUID, kind StoreKind) bool {
	_, err := os.Stat(d.pathFor(id, kind))
	return err == nil
}

func (d DiskStore) Size(id uuid.UUID, kind StoreKind) int64 {
	stat, err := os.Stat(d.pathFor(id, kind))
	if err != nil {
		return 0
	}
	return stat.Size()
}

func (d DiskStore) Delete(id uuid.UUID, kind StoreKind) error {
	return os.Remove(d.pathFor(id, kind))
}

func (d DiskStore) DeleteAll(id uuid.UUID) error {
	for _, kind := range StoreKinds {
		if d.Exists(id, kind) {
			err := d.Delete(id, kind)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func (d DiskStore) pathFor(id uuid.UUID, kind StoreKind) string {
	var suffix string

	switch kind {
	case KindRaw:
		suffix = ".photo"
	case KindThumbnailJpeg:
		suffix = ".thumb"
	case KindScreenJpeg:
		suffix = ".screen.jpg"
	default:
		suffix = fmt.Sprintf(".%d", kind)
	}

	str := id.String()
	return path.Join(d.Path, str[len(str)-2:], str+suffix)
}
