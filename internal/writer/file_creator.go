package writer

import "io"

type FileCreator interface {
	Create(path string) (io.WriteCloser, error)
	MkdirAll(path string) error
}
