package main

import (
	"flag"
	"fmt"
	"log"
	"path/filepath"

	"github.com/Mihir2423/ssggenerator/internal/buildstate"
	"github.com/Mihir2423/ssggenerator/internal/cache"
	"github.com/Mihir2423/ssggenerator/internal/fs"
	"github.com/Mihir2423/ssggenerator/internal/site"
	"github.com/Mihir2423/ssggenerator/internal/writer"
)

func main() {
	input := flag.String("input", "", "Input path")
	output := flag.String("output", "", "Output path")

	flag.Parse()

	hasProvidedInput := false
	hasProvidedOuput := false

	flag.Visit(func(f *flag.Flag) {
		switch f.Name {
		case "input":
			hasProvidedInput = true
		case "output":
			hasProvidedOuput = true
		}
	})

	if !hasProvidedInput {
		log.Fatal("You should enter the Input path of markdown file.")
	}
	if !hasProvidedOuput {
		log.Fatal("You should enter the Output path of html file.")
	}

	buildStatePath := filepath.Join(*output, ".ssg", "build-state.json")
	cacheDir := filepath.Join(*output, ".ssg", "cache")

	state := buildstate.New()
	if err := state.Load(buildStatePath); err != nil {
		log.Printf("Warning: could not load build state: %v", err)
	}

	cacheManager := cache.New(cacheDir)
	if err := cacheManager.Init(); err != nil {
		log.Fatalf("Failed to initialize cache: %v", err)
	}

	gen := site.Generator{
		FS:         fs.OSReader{},
		BuildState: state,
		Cache:      cacheManager,
	}

	result, err := gen.DiscoverAndClassify(*input, *output)
	if err != nil {
		log.Fatalf("Error discovering pages: %v", err)
	}

	htmlWriter := writer.HTMLWriter{
		Creator: writer.OSCreator{},
		Cache:   cacheManager,
	}

	err = htmlWriter.Write(result)
	if err != nil {
		log.Fatalf("Error writing HTML files: %v", err)
	}

	gen.UpdateBuildState(result)

	if err := state.Save(buildStatePath); err != nil {
		log.Printf("Warning: failed to save build state: %v", err)
	}

	totalFiles := len(result.ChangedPages) + len(result.UnchangedFiles)
	log.Printf("Processed %d markdown files (%d changed, %d cached)\n",
		totalFiles, len(result.ChangedPages), len(result.UnchangedFiles))
	fmt.Println("Input:", *input)
	fmt.Println("Output:", *output)
}
