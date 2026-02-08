package writer

import (
	"errors"
	"fmt"
	"io"

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
	for _, page := range pages {
		err := w.Creator.MkdirAll(page.OutputPath)
		if err != nil {
			return fmt.Errorf("%w: %w", ErrCreateOutputDir, err)
		}
		f, err := w.Creator.Create(page.OutputPath)
		if err != nil {
			return fmt.Errorf("%w %s: %w", ErrCreateFile, page.OutputPath, err)
		}
		_, err = io.WriteString(f, page.HTML)
		f.Close()
		if err != nil {
			return fmt.Errorf("%w %s: %w", ErrWriteFile, page.OutputPath, err)
		}

	}
	return nil
}
