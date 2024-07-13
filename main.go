package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"os"
	"path/filepath"
)

const (
	siteBuildFile = "sitebuild.json"
	outputDir     = "dist"
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}
}

func run() error {
	cfg, err := readSiteBuildConfig(siteBuildFile)
	if err != nil {
		return err
	}
	fmt.Printf("%#v\n", cfg)

	generatePages(cfg.BasePath, cfg.Pages)

	return nil
}

func generatePages(basePath string, pages map[string]page) error {
	for outfile, p := range pages {
		fmt.Printf("Generating page %s with URL %s\n", outfile, p.URL)
		for _, s := range p.Sources {
			fmt.Printf("  - source: %s\n", s)
		}

		// add the base path to the outFile and each source
		// files in the slice
		var sources []string
		for _, s := range p.Sources {
			sources = append(sources, filepath.Join(basePath, s))
		}
		outfile = filepath.Join(outputDir, outfile)

		if err := generateSinglePage(outfile, sources); err != nil {
			return err
		}
	}

	return nil
}

func generateSinglePage(outfile string, sources []string) error {
	fmt.Printf("Generating page %s with sources %v\n", outfile, sources)

	// parse html templates
	tmpl, err := template.ParseFiles(sources...)
	if err != nil {
		log.Fatal(err)
	}

	// create output file
	f, err := os.Create(outfile)
	if err != nil {
		return err
	}
	defer f.Close()

	// execute template
	if err := tmpl.Execute(f, nil); err != nil {
		return fmt.Errorf("error executing template: %w", err)
	}

	return nil
}

type page struct {
	URL     string   `json:"url"`
	Sources []string `json:"sources"`
}

type siteBuildConfig struct {
	BasePath string          `json:"basepath"`
	Pages    map[string]page `json:"pages"`
}

func readSiteBuildConfig(fileName string) (*siteBuildConfig, error) {
	f, err := os.Open(fileName)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	var cfg siteBuildConfig
	decoder := json.NewDecoder(f)
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&cfg); err != nil {
		return nil, fmt.Errorf("error decoding file %s: %w", fileName, err)
	}

	return &cfg, nil
}
