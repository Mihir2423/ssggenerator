package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/Mihir2423/ssggenerator/internal/fs"
	"github.com/Mihir2423/ssggenerator/internal/site"
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

	gen := site.Generator{
		FS: fs.OSReader{},
	}

	pages, err := gen.DiscoverPages(*input, *output)

	if err != nil {
		log.Fatal(err)
	}
	log.Printf("discovered %d markdown files\n", len(pages))
	fmt.Println("Input:", *input)
	fmt.Println("Output:", *output)
}
