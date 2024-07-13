package app

import (
	"html/template"
	"net/http"
	"path/filepath"

	"github.com/andyfusniak/sitebuild/internal/site"
	log "github.com/sirupsen/logrus"
)

func (a *App) routes(cfg *site.BuildConfig) *http.ServeMux {
	mux := http.NewServeMux()

	for _, p := range cfg.Pages {
		var sources []string
		for _, s := range p.Sources {
			sources = append(sources, filepath.Join(cfg.BasePath, s))
		}
		mux.Handle(p.URL, newHandler(sources))
	}
	return mux
}

func newHandler(sources []string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tmpl, err := template.ParseFiles(sources...)
		if err != nil {
			log.Fatal(err)
		}

		if err := tmpl.Execute(w, nil); err != nil {
			log.Error("error executing template", err)
		}
	})
}
