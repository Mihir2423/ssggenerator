package writer

import (
	"errors"
	"fmt"
	"io"
	"sync"

	"github.com/Mihir2423/ssggenerator/internal/cache"
	"github.com/Mihir2423/ssggenerator/internal/site"
)

var (
	ErrCreateOutputDir = errors.New("failed to create output directory")
	ErrCreateFile      = errors.New("failed to create file")
	ErrWriteFile       = errors.New("failed to write file")
	ErrCacheCopy       = errors.New("failed to copy from cache")
)

type HTMLWriter struct {
	Creator FileCreator
	Cache   *cache.Manager
}

func (w HTMLWriter) Write(result *site.BuildResult) error {
	var wg sync.WaitGroup
	errChan := make(chan error, len(result.ChangedPages)+len(result.UnchangedFiles))

	// Process changed pages (write new HTML)
	for _, page := range result.ChangedPages {
		wg.Add(1)
		go func(p site.Page) {
			defer wg.Done()

			err := w.Creator.MkdirAll(p.OutputPath)
			if err != nil {
				errChan <- fmt.Errorf("%w: %w", ErrCreateOutputDir, err)
				return
			}
			f, err := w.Creator.Create(p.OutputPath)
			if err != nil {
				errChan <- fmt.Errorf("%w %s: %w", ErrCreateFile, p.OutputPath, err)
				return
			}
			_, err = io.WriteString(f, p.HTML)
			f.Close()
			if err != nil {
				errChan <- fmt.Errorf("%w %s: %w", ErrWriteFile, p.OutputPath, err)
				return
			}

			// Save to cache for future incremental builds
			if w.Cache != nil {
				if err := w.Cache.Save(p.SourcePath, p.OutputPath, []byte(p.HTML)); err != nil {
					// Non-fatal: cache failure shouldn't break the build
					fmt.Printf("Warning: failed to cache %s: %v\n", p.SourcePath, err)
				}
			}
		}(page)
	}

	// Copy unchanged files from cache
	for _, file := range result.UnchangedFiles {
		wg.Add(1)
		go func(f site.UnchangedFile) {
			defer wg.Done()

			if w.Cache != nil {
				if err := w.Cache.CopyToOutput(f.OutputPath); err != nil {
					errChan <- fmt.Errorf("%w %s: %w", ErrCacheCopy, f.OutputPath, err)
					return
				}
			}
		}(file)
	}

	wg.Wait()
	close(errChan)

	for err := range errChan {
		if err != nil {
			return err
		}
	}

	return nil
}
