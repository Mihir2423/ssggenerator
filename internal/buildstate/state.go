package buildstate

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

type FileState struct {
	Hash       string    `json:"hash"`
	ModTime    time.Time `json:"mod_time"`
	OutputPath string    `json:"output_path"`
}

type BuildState struct {
	Files map[string]FileState `json:"files"`
}

func New() *BuildState {
	return &BuildState{
		Files: make(map[string]FileState),
	}
}

func (bs *BuildState) Load(path string) error {
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return fmt.Errorf("failed to read build state: %w", err)
	}

	return json.Unmarshal(data, bs)
}

func (bs *BuildState) Save(path string) error {
	data, err := json.MarshalIndent(bs, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal build state: %w", err)
	}

	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create build state directory: %w", err)
	}

	return os.WriteFile(path, data, 0644)
}

func (bs *BuildState) HasChanged(sourcePath string, content []byte, outputPath string) bool {
	currentHash := computeHash(content)

	if state, exists := bs.Files[sourcePath]; exists {
		return state.Hash != currentHash || state.OutputPath != outputPath
	}

	return true
}

func (bs *BuildState) Update(sourcePath string, content []byte, outputPath string) {
	bs.Files[sourcePath] = FileState{
		Hash:       computeHash(content),
		ModTime:    time.Now(),
		OutputPath: outputPath,
	}
}

func computeHash(content []byte) string {
	hash := sha256.Sum256(content)
	return hex.EncodeToString(hash[:])
}
