package site

import (
	"errors"
	"os"
	"testing"

	"github.com/Mihir2423/ssggenerator/internal/buildstate"
	"github.com/Mihir2423/ssggenerator/internal/cache"
	"github.com/Mihir2423/ssggenerator/internal/fs"
)

type fakeDirEntry struct {
	name  string
	isDir bool
}

func (f fakeDirEntry) Name() string               { return f.name }
func (f fakeDirEntry) IsDir() bool                { return f.isDir }
func (f fakeDirEntry) Type() os.FileMode          { return 0 }
func (f fakeDirEntry) Info() (os.FileInfo, error) { return nil, nil }

type fakeReader struct {
	entries []os.DirEntry
	files   map[string][]byte
	err     error
}

func (f fakeReader) ReadFile(path string) ([]byte, error) {
	if f.err != nil {
		return nil, f.err
	}
	return f.files[path], nil
}

func (f fakeReader) ReadDir(path string) ([]os.DirEntry, error) {
	if f.err != nil {
		return nil, f.err
	}
	return f.entries, nil
}

var _ fs.Reader = fakeReader{}

func TestDiscoverAndClassify_Success(t *testing.T) {
	reader := fakeReader{
		entries: []os.DirEntry{
			fakeDirEntry{name: "index.md", isDir: false},
			fakeDirEntry{name: "about.md", isDir: false},
		},
		files: map[string][]byte{
			"/input/index.md": []byte("# Home"),
			"/input/about.md": []byte("## About"),
		},
	}

	state := buildstate.New()
	gen := Generator{
		FS:         reader,
		BuildState: state,
		Cache:      cache.New("/tmp/cache"),
	}

	result, err := gen.DiscoverAndClassify("/input", "/output")

	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if len(result.ChangedPages) != 2 {
		t.Errorf("expected 2 changed pages, got %d", len(result.ChangedPages))
	}
	if len(result.UnchangedFiles) != 0 {
		t.Errorf("expected 0 unchanged files, got %d", len(result.UnchangedFiles))
	}
}

func TestDiscoverAndClassify_ReadDirError(t *testing.T) {
	reader := fakeReader{
		err: errors.New("permission denied"),
	}

	gen := Generator{
		FS:         reader,
		BuildState: buildstate.New(),
		Cache:      cache.New("/tmp/cache"),
	}

	_, err := gen.DiscoverAndClassify("/input", "/output")

	if err == nil {
		t.Error("expected error but got nil")
	}
}

func TestDiscoverAndClassify_IncrementalBuild(t *testing.T) {
	reader := fakeReader{
		entries: []os.DirEntry{
			fakeDirEntry{name: "index.md", isDir: false},
			fakeDirEntry{name: "about.md", isDir: false},
		},
		files: map[string][]byte{
			"/input/index.md": []byte("# Home"),
			"/input/about.md": []byte("## About"),
		},
	}

	state := buildstate.New()
	state.Update("/input/index.md", []byte("# Home"), "/output/index.html")

	gen := Generator{
		FS:         reader,
		BuildState: state,
		Cache:      cache.New("/tmp/cache"),
	}

	result, err := gen.DiscoverAndClassify("/input", "/output")

	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if len(result.ChangedPages) != 1 {
		t.Errorf("expected 1 changed page, got %d", len(result.ChangedPages))
	}
	if len(result.UnchangedFiles) != 1 {
		t.Errorf("expected 1 unchanged file, got %d", len(result.UnchangedFiles))
	}
}
