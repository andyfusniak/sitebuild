package cli

import (
	"github.com/andyfusniak/sitebuild/internal/routegen"
	"github.com/andyfusniak/sitebuild/internal/site"
	"github.com/spf13/cobra"
)

// NewCmdGenRoutes creates a new genroutes command to generate the
// routes for the site.
func NewCmdGenRoutes(destDir, siteBuildFile string) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "genroutes",
		Short: "generate the routes for the site",
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg, err := site.NewSiteBuildConfigFromFile(siteBuildFile)
			if err != nil {
				return err
			}

			routegen := routegen.NewRouteGenerator(cfg.Pages)
			if err := routegen.GenerateRoutes(); err != nil {
				return err
			}

			return nil
		},
	}
	return cmd
}
