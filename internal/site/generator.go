package site

import (
	"errors"
	"fmt"
	"path/filepath"
	"strings"

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
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		if !strings.HasSuffix(entry.Name(), ".md") {
			continue
		}
		sourcePath := filepath.Join(inputDir, entry.Name())
		content, err := g.FS.ReadFile(sourcePath)
		if err != nil {
			return nil, fmt.Errorf("%w %s: %w", ErrReadFile, sourcePath, err)
		}
		outputName := strings.TrimSuffix(entry.Name(), ".md") + ".html"
		outputPath := filepath.Join(outputDir, outputName)
		html := markdown.ToHTML(content)

		pages = append(pages, Page{
			SourcePath: sourcePath,
			OutputPath: outputPath,
			Content:    content,
			HTML:       html,
		})

	}
	return pages, nil
}
