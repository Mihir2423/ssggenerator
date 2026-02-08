package writer

import (
	"bytes"
	"errors"
	"io"
	"testing"

	"github.com/Mihir2423/ssggenerator/internal/site"
)

type fakeWriteCloser struct {
	bytes.Buffer
}

func (f *fakeWriteCloser) Close() error {
	return nil
}

type fakeFileCreator struct {
	files  map[string]*fakeWriteCloser
	mkErr  error
	crtErr error
}

func (f *fakeFileCreator) Create(path string) (io.WriteCloser, error) {
	if f.crtErr != nil {
		return nil, f.crtErr
	}
	if f.files == nil {
		f.files = make(map[string]*fakeWriteCloser)
	}
	fw := &fakeWriteCloser{}
	f.files[path] = fw
	return fw, nil
}

func (f *fakeFileCreator) MkdirAll(path string) error {
	if f.mkErr != nil {
		return f.mkErr
	}
	return nil
}

var _ FileCreator = &fakeFileCreator{}

func TestHTMLWriter_Write_Success(t *testing.T) {
	creator := &fakeFileCreator{}
	writer := HTMLWriter{Creator: creator}

	pages := []site.Page{
		{
			SourcePath: "/input/index.md",
			OutputPath: "/output/index.html",
			HTML:       "<h1>Home</h1>",
		},
	}

	err := writer.Write(pages)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	if creator.files["/output/index.html"] == nil {
		t.Error("expected file to be created")
	}
}

func TestHTMLWriter_Write_MkdirError(t *testing.T) {
	creator := &fakeFileCreator{
		mkErr: errors.New("permission denied"),
	}
	writer := HTMLWriter{Creator: creator}

	pages := []site.Page{
		{
			OutputPath: "/output/index.html",
			HTML:       "<h1>Home</h1>",
		},
	}

	err := writer.Write(pages)
	if err == nil {
		t.Error("expected error but got nil")
	}
}
