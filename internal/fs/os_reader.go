package fs

import "os"

type OSReader struct{}

func (OSReader) ReadFile(path string) ([]byte, error) {
	return os.ReadFile(path)
}

func (OSReader) ReadDir(path string) ([]os.DirEntry, error) {
	return os.ReadDir(path)
}
