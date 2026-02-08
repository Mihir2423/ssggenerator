package writer

import (
	"errors"
	"fmt"
	"io"
	"sync"

	"github.com/Mihir2423/ssggenerator/internal/site"
)

var (
	ErrCreateOutputDir = errors.New("failed to create output directory")
	ErrCreateFile      = errors.New("failed to create file")
	ErrWriteFile       = errors.New("failed to write file")
)

type HTMLWriter struct {
	Creator FileCreator
}

func (w HTMLWriter) Write(pages []site.Page) error {
	var wg sync.WaitGroup
	errChan := make(chan error, len(pages))

	for _, page := range pages {
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
		}(page)
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
