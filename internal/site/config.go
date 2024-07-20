package site

import (
	"encoding/json"
	"fmt"
	"os"
)

type BuildConfig struct {
	SourceDir string          `json:"sourceDir"`
	Pages     map[string]Page `json:"pages"`
}

type Page struct {
	URL         string   `json:"url"`
	Sources     []string `json:"sources"`
	DataSources []string `json:"dataSources"`
}

// NewSiteBuildConfigFromFile reads a site build configuration from a file.
func NewSiteBuildConfigFromFile(fileName string) (*BuildConfig, error) {
	f, err := os.Open(fileName)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	var cfg BuildConfig
	decoder := json.NewDecoder(f)
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&cfg); err != nil {
		return nil, fmt.Errorf("error decoding file %s: %w", fileName, err)
	}

	return &cfg, nil
}
