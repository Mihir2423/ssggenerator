package writer

import (
	"fmt"
	"io"

	"github.com/Mihir2423/ssggenerator/internal/site"
)

type HTMLWriter struct {
	Creator FileCreator
}

func (w HTMLWriter) Write(pages []site.Page) error {
	for _, page := range pages {
		err := w.Creator.MkdirAll(page.OutputPath)
		if err != nil {
			return fmt.Errorf("creating output dir: %w", err)
		}
		f, err := w.Creator.Create(page.OutputPath)
		if err != nil {
			return fmt.Errorf("creating file %s: %w", page.OutputPath, err)
		}
		_, err = io.WriteString(f, page.HTML)
		f.Close()
		if err != nil {
			return fmt.Errorf("writing file %s: %w", page.OutputPath, err)
		}

	}
	return nil
}
