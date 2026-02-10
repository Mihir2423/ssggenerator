package writer

import (
	"bytes"
	"errors"
	"io"
	"testing"

	"github.com/Mihir2423/ssggenerator/internal/cache"
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
	writer := HTMLWriter{
		Creator: creator,
		Cache:   cache.New("/tmp/cache"),
	}

	result := &site.BuildResult{
		ChangedPages: []site.Page{
			{
				SourcePath: "/input/index.md",
				OutputPath: "/output/index.html",
				HTML:       "<h1>Home</h1>",
			},
		},
		UnchangedFiles: []site.UnchangedFile{},
	}

	err := writer.Write(result)
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
	writer := HTMLWriter{
		Creator: creator,
		Cache:   cache.New("/tmp/cache"),
	}

	result := &site.BuildResult{
		ChangedPages: []site.Page{
			{
				OutputPath: "/output/index.html",
				HTML:       "<h1>Home</h1>",
			},
		},
		UnchangedFiles: []site.UnchangedFile{},
	}

	err := writer.Write(result)
	if err == nil {
		t.Error("expected error but got nil")
	}
}

func TestHTMLWriter_Write_WithUnchangedFiles(t *testing.T) {
	creator := &fakeFileCreator{}
	writer := HTMLWriter{
		Creator: creator,
		Cache:   nil,
	}

	result := &site.BuildResult{
		ChangedPages: []site.Page{
			{
				SourcePath: "/input/index.md",
				OutputPath: "/output/index.html",
				HTML:       "<h1>Home</h1>",
			},
		},
		UnchangedFiles: []site.UnchangedFile{
			{
				SourcePath: "/input/about.md",
				OutputPath: "/output/about.html",
			},
		},
	}

	err := writer.Write(result)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	// Changed file should be written
	if creator.files["/output/index.html"] == nil {
		t.Error("expected changed file to be created")
	}
}
