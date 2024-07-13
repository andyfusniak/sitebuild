package cli

import (
	"context"
	"os"
	"runtime"

	"github.com/andyfusniak/sitebuild/internal/app"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

// NewCmdServer creates a new server command. This command starts the web service.
func NewCmdServer(version, gitcommit, defaultPort string) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "server",
		Aliases: []string{"serve"},
		Short:   "start the web service",
		RunE: func(cmd *cobra.Command, args []string) error {
			// set up logging
			defer func() {
				log.Infof("[main] goodbye from sitebuild server version %s (%s)", version, gitcommit)
			}()
			initLogging("info")
			log.Infof("[main] hello from sitebuild version %s (%s) %s for %s %s",
				version, gitcommit, runtime.Version(), runtime.GOOS, runtime.GOARCH)

			// HTTP application server
			app, err := app.New(defaultPort)
			if err != nil {
				return err
			}
			if err := app.Start(context.Background()); err != nil {
				return err
			}

			return nil
		},
	}

	return cmd
}

func initLogging(logLevel string) {
	// Output logs with colour
	log.SetFormatter(&log.TextFormatter{
		ForceColors: true,
	})

	// Output to stdout instead of the default stderr
	log.SetOutput(os.Stdout)

	// Log debug level severity or above.
	logrusLevel := logLevelToLogrusLevel(logLevel)
	log.SetLevel(logrusLevel)
}

func logLevelToLogrusLevel(v string) log.Level {
	switch v {
	case "panic":
		return log.PanicLevel
	case "fatal":
		return log.FatalLevel
	case "error":
		return log.ErrorLevel
	case "warn":
		return log.WarnLevel
	case "info":
		return log.InfoLevel
	case "debug":
		return log.DebugLevel
	case "trace":
		return log.TraceLevel
	default:
		return log.DebugLevel
	}
}
