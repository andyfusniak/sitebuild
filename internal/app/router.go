package app

import (
	"html/template"
	"net/http"
	"os"
	"path/filepath"

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
			mux.Handle(p.URL, customRootHandler(newHandler(sources), staticFilesDir))
			continue
		}

		mux.Handle(p.URL, newHandler(sources))
	}

	return mux
}

func newHandler(sources []string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tmpl, err := template.ParseFiles(sources...)
		if err != nil {
			log.Error("error parsing template files", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		if err := tmpl.Execute(w, nil); err != nil {
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
