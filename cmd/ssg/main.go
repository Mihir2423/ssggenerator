package main

import (
	"flag"
	"fmt"
	"log"
)

func main() {
	input := flag.String("input", "", "Input path")
	output := flag.String("output", "", "Output path")
	config := flag.String("config", "", "Optional JSON configuration (e.g. '{\"key\":\"value\"}')")

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
		return
	}
	if !hasProvidedOuput {
		log.Fatal("You should enter the Output path of html file.")
		return
	}

	fmt.Println("Input:", *input)
	fmt.Println("Output:", *output)
	fmt.Println("Config:", *config)
}
