package writer

import (
	"io"
	"os"
	"path/filepath"
)

type OSCreator struct{}

func (OSCreator) Create(path string) (io.WriteCloser, error) {
	return os.Create(path)
}

func (OSCreator) MkdirAll(path string) error {
	return os.MkdirAll(filepath.Dir(path), 0755)
}
