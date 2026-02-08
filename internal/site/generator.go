package site

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/Mihir2423/ssggenerator/internal/fs"
)

type Generator struct {
	FS fs.Reader
}

func (g Generator) DiscoverPages(inputDir, outputDir string) ([]Page, error) {
	entries, err := g.FS.ReadDir(inputDir)
	if err != nil {
		return nil, fmt.Errorf("reading input dir: %w", err)
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
			return nil, fmt.Errorf("reading file %s: %w", sourcePath, err)
		}
		outputName := strings.TrimSuffix(entry.Name(), ".md") + ".html"
		outputPath := filepath.Join(outputDir, outputName)

		pages = append(pages, Page{
			SourcePath: sourcePath,
			OutputPath: outputPath,
			Content:    content,
		})

	}
	return pages, nil
}
