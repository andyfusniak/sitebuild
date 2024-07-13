package app

import (
	"net/http"
)

func (a *App) routes() *http.ServeMux {
	mux := http.NewServeMux()
	// mux.Use(a.handler.JSONHeader)

	// auth
	mux.HandleFunc("/test", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Hello, World"))
	})

	return mux
}
