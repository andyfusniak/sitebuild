package main

import (
	"context"
	"fmt"
	"os"

	"github.com/andyfusniak/sitebuild/internal/cli"
	"github.com/spf13/cobra"
)

var (
	version   string
	gitCommit string
)

const (
	siteBuildFile = "sitebuild.json"
	outputDir     = "dist"
	defaultPort   = "7000"
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}
}

func run() error {
	// cli application
	cliApp := cli.NewApp(
		cli.WithVersion(version),
		cli.WithGitCommit(gitCommit),
	)

	root := cobra.Command{
		Use:     "sitebuild",
		Short:   "sitebuild static site generator",
		Version: cliApp.Version(),
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			ctx := cmd.Context()
			v := ctx.Value(cli.AppKey("app"))
			_ = v.(*cli.App)
		},
	}
	root.AddCommand(cli.NewCmdBuild(outputDir, siteBuildFile))
	root.AddCommand(cli.NewCmdServer(version, gitCommit, siteBuildFile, defaultPort))

	ctx := context.WithValue(context.Background(), cli.AppKey("app"), cliApp)
	if err := root.ExecuteContext(ctx); err != nil {
		fmt.Fprintf(os.Stderr, "%+v\n", err)
		os.Exit(1)
	}

	return nil

}
