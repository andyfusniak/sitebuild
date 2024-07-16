package site

import (
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"text/template"
)

type Site struct {
	cwd       string
	destDir   string
	sourceDir string
}

// NewSiteBuilder creates a new site with the destination directory.
// destDir and sourceDir are relative to the current working directory.
func NewSiteBuilder(destDir, sourceDir string) (*Site, error) {
	cwd, err := os.Getwd()
	if err != nil {
		return nil, errors.New("failed to get current working directory")
	}

	return &Site{
		cwd:       cwd,
		destDir:   destDir,
		sourceDir: sourceDir,
	}, nil
}

// GeneratePages generates the pages for the site. It reads the
// sources for each page, parses the templates and writes the output
// to the output directory. The basePath is prepended to the source
// file paths.
func (s *Site) GeneratePages(basePath string, pages map[string]Page) error {
	if err := s.purgeDestDir(); err != nil {
		return err
	}

	if err := s.ensureDestDirExist(); err != nil {
		return err
	}

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
		outfile = filepath.Join(s.destDir, outfile)

		if err := generateSinglePage(outfile, sources); err != nil {
			return err
		}
	}

	if err := s.CopyStaticFiles(); err != nil {
		return err
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

func (s *Site) ensureDestDirExist() error {
	outputDir := filepath.Join(s.cwd, s.destDir)
	if _, err := os.Stat(outputDir); os.IsNotExist(err) {
		if err := os.Mkdir(outputDir, 0700); err != nil {
			return fmt.Errorf("check permissions as failed to mkdir %q", outputDir)
		}
	}
	return nil
}

func (s *Site) purgeDestDir() error {
	destDir := filepath.Join(s.cwd, s.destDir)
	if err := os.RemoveAll(destDir); err != nil {
		return fmt.Errorf("failed to remove %q: %w", destDir, err)
	}
	return nil
}

func (s *Site) CopyStaticFiles() error {
	staticFilesDir := filepath.Join(s.cwd, s.sourceDir, "static")
	outDir := filepath.Join(s.cwd, s.destDir)

	if _, err := os.Stat(staticFilesDir); os.IsNotExist(err) {
		return nil // Source directory does not exist, nothing to copy
	}

	err := filepath.Walk(staticFilesDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err // Handle potential walking error
		}

		// Calculate relative path to preserve directory structure
		relPath, err := filepath.Rel(staticFilesDir, path)
		if err != nil {
			return err
		}
		destPath := filepath.Join(outDir, relPath)

		if info.IsDir() {
			// Create directory if it does not exist
			return os.MkdirAll(destPath, info.Mode())
		} else {
			// Copy file content
			return copyFile(path, destPath)
		}
	})

	return err
}

// copyFile copies a single file from src to dst
func copyFile(src, dst string) error {
	sourceFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer sourceFile.Close()

	// Ensure destination directory exists
	destDir := filepath.Dir(dst)
	if err := os.MkdirAll(destDir, 0755); err != nil {
		return err
	}

	destFile, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer destFile.Close()

	_, err = io.Copy(destFile, sourceFile)
	return err
}
