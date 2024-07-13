package cli

import (
	"fmt"
	"io"
	"os"
)

type AppKey string

type App struct {
	version   string
	gitCommit string
	stdout    io.Writer
	stderr    io.Writer
}

type Option func(*App)

// NewApp creates a new CLI application.
func NewApp(options ...Option) *App {
	a := &App{
		stdout: os.Stdout,
		stderr: os.Stderr,
	}

	for _, o := range options {
		o(a)
	}
	return a
}

// WithVersion option to set the cli version.
func WithVersion(s string) Option {
	return func(a *App) {
		a.version = s
	}
}

// WithGitCommit option to set the git commit hash.
func WithGitCommit(s string) Option {
	return func(a *App) {
		a.gitCommit = s
	}
}

// WithStdOut option to set default output stream.
func WithStdOut(w io.Writer) Option {
	return func(a *App) {
		a.stdout = w
	}
}

// WithStdErr option to set default error stream.
func WithStdErr(w io.Writer) Option {
	return func(a *App) {
		a.stderr = w
	}
}

// Version returns the cli application version.
func (a *App) Version() string {
	return fmt.Sprintf("%s (git commit: %s)", a.version, a.gitCommit)
}
