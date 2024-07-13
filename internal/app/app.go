package app

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/andyfusniak/sitebuild/internal/site"
	log "github.com/sirupsen/logrus"
)

type App struct {
	port   string
	router *http.ServeMux
}

// New creates a new app server.
func New(cfg *site.BuildConfig, port string) (*App, error) {
	app := &App{
		port: port,
	}

	// routing
	app.router = app.routes(cfg)

	return app, nil
}

// Start the app server.
func (a *App) Start(ctx context.Context) error {
	// HTTP Service
	srv := http.Server{
		Addr:    "0.0.0.0:" + a.port,
		Handler: a.router,
	}

	// signals
	idleConnsClosed := make(chan struct{})
	go func() {
		sigint := make(chan os.Signal, 1)

		// interrupt signal sent from terminal
		signal.Notify(sigint, os.Interrupt)
		// sigterm signal sent from kubernetes
		signal.Notify(sigint, syscall.SIGTERM)
		switch sig := <-sigint; sig {
		case syscall.SIGINT:
			log.Info("[main] received signal SIGINT")
		case syscall.SIGTERM:
			log.Info("[main] received signal SIGTERM")
		default:
			log.Info("[main] received unexpected signal", "signal", sig)
		}

		log.Info("[main] gracefully shutting down the server...")

		// We received an interrupt signal, shut down.
		if err := srv.Shutdown(ctx); err != nil {
			// Error from closing listeners, or context timeout:
			log.Infof("[main] HTTP server Shutdown: %+v", err)
		}
		log.Info("[main] HTTP server shutdown complete")
		close(idleConnsClosed)
	}()

	log.Infof("[main] server listening on HTTP port %s", a.port)
	if err := srv.ListenAndServe(); err != http.ErrServerClosed {
		log.Errorf("[main] error: %+v", err)
		return err
	}
	<-idleConnsClosed

	return nil
}
