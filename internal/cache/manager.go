package cache

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
)

type Manager struct {
	cacheDir string
}

func New(cacheDir string) *Manager {
	return &Manager{cacheDir: cacheDir}
}

func (m *Manager) Init() error {
	return os.MkdirAll(m.cacheDir, 0755)
}

func (m *Manager) Save(sourcePath, outputPath string, content []byte) error {
	cachePath := m.getCachePath(outputPath)

	dir := filepath.Dir(cachePath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create cache directory: %w", err)
	}

	return os.WriteFile(cachePath, content, 0644)
}

func (m *Manager) CopyToOutput(outputPath string) error {
	cachePath := m.getCachePath(outputPath)

	src, err := os.Open(cachePath)
	if err != nil {
		return fmt.Errorf("failed to open cached file: %w", err)
	}
	defer src.Close()

	dir := filepath.Dir(outputPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create output directory: %w", err)
	}

	dst, err := os.Create(outputPath)
	if err != nil {
		return fmt.Errorf("failed to create output file: %w", err)
	}
	defer dst.Close()

	if _, err := io.Copy(dst, src); err != nil {
		return fmt.Errorf("failed to copy cached file: %w", err)
	}

	return nil
}

func (m *Manager) HasCached(outputPath string) bool {
	cachePath := m.getCachePath(outputPath)
	_, err := os.Stat(cachePath)
	return err == nil
}

func (m *Manager) getCachePath(outputPath string) string {
	filename := filepath.Base(outputPath)
	return filepath.Join(m.cacheDir, filename)
}
