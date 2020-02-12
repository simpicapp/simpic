package storage

import (
	"github.com/satori/go.uuid"
	"io"
	"os"
	"path"
)

// DiskDriver stores files on disk according to their UUID.
// Files are grouped into folders according to the last two bytes of the UUID, with each folder being located under
// the given path.
type DiskDriver struct {
	Path string
}

func (d DiskDriver) Read(id uuid.UUID) (io.ReadCloser, error) {
	return os.Open(d.pathFor(id))
}

func (d DiskDriver) Write(id uuid.UUID) (io.WriteCloser, error) {
	file := d.pathFor(id)

	if err := os.MkdirAll(path.Dir(file), os.ModePerm); err != nil {
		return nil, err
	}

	return os.Create(file)
}

func (d DiskDriver) pathFor(id uuid.UUID) string {
	str := id.String()
	return path.Join(d.Path, str[len(str)-2:], str)
}
