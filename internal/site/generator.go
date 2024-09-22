package site

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"text/template"

	"github.com/andyfusniak/sitebuild/internal/funcs"
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

		if err := generateSinglePage(outfile, sources, p.DataSources); err != nil {
			return err
		}
	}

	if err := s.CopyStaticFiles(); err != nil {
		return err
	}

	return nil
}

type Vars struct {
	Data map[string]any
}

func generateSinglePage(outfile string, sources, dataSources []string) error {
	if len(sources) == 0 {
		return fmt.Errorf("no template files provided for %s", outfile)
	}

	fmt.Printf("Generating page %s with sources %v and dataSources %v\n",
		outfile, sources, dataSources)

	rootTemplateName := filepath.Base(sources[0])
	funcMap := template.FuncMap(funcs.FuncMap())

	// parse html templates
	tmpl := template.New(rootTemplateName).Funcs(funcMap)
	tmpl, err := tmpl.ParseFiles(sources...)
	if err != nil {
		log.Fatal(err)
	}

	// create output file
	f, err := os.Create(outfile)
	if err != nil {
		return err
	}
	defer f.Close()

	var v Vars

	v.Data = make(map[string]any)
	for _, ds := range dataSources {
		var d any
		if err := loadJSONFile(filepath.Join("data", ds), &d); err != nil {
			return err
		}

		// strip the .json .whatever postfix from the file name
		// leaving only the base name
		name := filepath.Base(ds)
		extension := filepath.Ext(name)
		name = name[:len(name)-len(extension)]
		v.Data[name] = d
	}

	// execute template
	if err := tmpl.ExecuteTemplate(f, rootTemplateName, v); err != nil {
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

func loadJSONFile(fileName string, v interface{}) error {
	f, err := os.Open(fileName)
	if err != nil {
		return err
	}
	defer f.Close()

	decoder := json.NewDecoder(f)
	if err := decoder.Decode(v); err != nil {
		return fmt.Errorf("error decoding file %s: %w", fileName, err)
	}

	return nil
}
