package app

import (
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
	"os"
	"path/filepath"

	"github.com/andyfusniak/sitebuild/internal/funcs"
	"github.com/andyfusniak/sitebuild/internal/site"
	log "github.com/sirupsen/logrus"
)

func (a *App) routes(cfg *site.BuildConfig) *http.ServeMux {
	mux := http.NewServeMux()

	for _, p := range cfg.Pages {
		var sources []string
		for _, s := range p.Sources {
			sources = append(sources, filepath.Join(cfg.SourceDir, s))
		}

		// root page is a special case
		if p.URL == "/" {
			staticFilesDir := filepath.Join(cfg.SourceDir, "static")
			mux.Handle(p.URL, customRootHandler(newHandler(sources, p.DataSources), staticFilesDir))
			continue
		}

		mux.Handle(p.URL, newHandler(sources, p.DataSources))
	}

	return mux
}

type vars struct {
	Data map[string]any
}

func newHandler(sources, dataSources []string) http.Handler {
	// get the global templates func from the funcs package
	// and add it to the template.FuncMap
	funcMap := template.FuncMap(funcs.FuncMap())

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		// Ensure there is at least one source file
		if len(sources) == 0 {
			http.Error(w, "no template files provided", http.StatusInternalServerError)
			return
		}

		rootTemplateName := filepath.Base(sources[0])

		// parse the template files using the global funcMap

		// tmpl, err := template.ParseFiles(sources...)

		tmpl := template.New(rootTemplateName).Funcs(funcMap)
		tmpl, err := tmpl.ParseFiles(sources...)
		if err != nil {
			log.Error("error parsing template files", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		var v vars

		v.Data = make(map[string]any)
		for _, ds := range dataSources {
			var d any
			if err := loadJSONFile(filepath.Join("data", ds), &d); err != nil {
				log.Error("error loading JSON data", err)
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			// strip the .json .whatever postfix from the file name
			// leaving only the base name
			name := filepath.Base(ds)
			extension := filepath.Ext(name)
			name = name[:len(name)-len(extension)]
			v.Data[name] = d
		}

		if err := tmpl.ExecuteTemplate(w, rootTemplateName, v); err != nil {
			log.Error("error executing template", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	})
}

func customRootHandler(homePageHandler http.Handler, staticFilesDir string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Check if the request is specifically for the root ("/")
		if r.URL.Path == "/" {
			homePageHandler.ServeHTTP(w, r)
			return
		}

		// Attempt to serve a static file for any other request
		staticFilePath := filepath.Join(staticFilesDir, r.URL.Path)
		if _, err := os.Stat(staticFilePath); err == nil {
			http.ServeFile(w, r, staticFilePath)
			return
		}

		// If the file does not exist, you can decide to serve a 404 page,
		// redirect to the home page, or simply let the request fall through.
		http.NotFound(w, r)
	})
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
