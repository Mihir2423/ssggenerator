package site

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/Mihir2423/ssggenerator/internal/fs"
	"github.com/Mihir2423/ssggenerator/internal/markdown"
)

var (
	ErrDiscoveryFailed = errors.New("page discovery failed")
	ErrReadInputDir    = errors.New("failed to read input directory")
	ErrReadFile        = errors.New("failed to read file")
)

type Generator struct {
	FS fs.Reader
}

func (g Generator) DiscoverPages(inputDir, outputDir string) ([]Page, error) {
	entries, err := g.FS.ReadDir(inputDir)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrReadInputDir, err)
	}

	var pages []Page
	var mu sync.Mutex
	var wg sync.WaitGroup
	errChan := make(chan error, len(entries))

	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		if !strings.HasSuffix(entry.Name(), ".md") {
			continue
		}

		wg.Add(1)
		go func(entry os.DirEntry) {
			defer wg.Done()

			sourcePath := filepath.Join(inputDir, entry.Name())
			content, err := g.FS.ReadFile(sourcePath)
			if err != nil {
				errChan <- fmt.Errorf("%w %s: %w", ErrReadFile, sourcePath, err)
				return
			}
			outputName := strings.TrimSuffix(entry.Name(), ".md") + ".html"
			outputPath := filepath.Join(outputDir, outputName)
			html := markdown.ToHTML(content)

			mu.Lock()
			pages = append(pages, Page{
				SourcePath: sourcePath,
				OutputPath: outputPath,
				Content:    content,
				HTML:       html,
			})
			mu.Unlock()
		}(entry)
	}

	wg.Wait()
	close(errChan)

	for e := range errChan {
		if e != nil {
			return nil, e
		}
	}

	return pages, nil
}
