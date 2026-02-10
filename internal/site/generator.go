package site

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/Mihir2423/ssggenerator/internal/buildstate"
	"github.com/Mihir2423/ssggenerator/internal/cache"
	"github.com/Mihir2423/ssggenerator/internal/fs"
	"github.com/Mihir2423/ssggenerator/internal/markdown"
)

var (
	ErrDiscoveryFailed = errors.New("page discovery failed")
	ErrReadInputDir    = errors.New("failed to read input directory")
	ErrReadFile        = errors.New("failed to read file")
)

type Generator struct {
	FS         fs.Reader
	BuildState *buildstate.BuildState
	Cache      *cache.Manager
}

type BuildResult struct {
	ChangedPages   []Page
	UnchangedFiles []UnchangedFile
}

type UnchangedFile struct {
	SourcePath string
	OutputPath string
}

func (g Generator) DiscoverAndClassify(inputDir, outputDir string) (*BuildResult, error) {
	entries, err := g.FS.ReadDir(inputDir)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrReadInputDir, err)
	}

	result := &BuildResult{
		ChangedPages:   []Page{},
		UnchangedFiles: []UnchangedFile{},
	}

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

			// Check if file has changed
			if g.BuildState != nil && !g.BuildState.HasChanged(sourcePath, content, outputPath) {
				// File unchanged, can copy from cache
				mu.Lock()
				result.UnchangedFiles = append(result.UnchangedFiles, UnchangedFile{
					SourcePath: sourcePath,
					OutputPath: outputPath,
				})
				mu.Unlock()
				return
			}

			// File changed or new, process it
			html := markdown.ToHTML(content)

			mu.Lock()
			result.ChangedPages = append(result.ChangedPages, Page{
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

	return result, nil
}

func (g Generator) UpdateBuildState(result *BuildResult) {
	if g.BuildState == nil {
		return
	}

	for _, page := range result.ChangedPages {
		g.BuildState.Update(page.SourcePath, page.Content, page.OutputPath)
	}
}
