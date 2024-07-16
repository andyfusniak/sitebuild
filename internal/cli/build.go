package cli

import (
	"github.com/andyfusniak/sitebuild/internal/site"
	"github.com/spf13/cobra"
)

// NewCmdBuild creates a new build command. This command builds the site.
func NewCmdBuild(destDir, siteBuildFile string) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "build",
		Short: "builds the site",
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg, err := site.NewSiteBuildConfigFromFile(siteBuildFile)
			if err != nil {
				return err
			}

			site, err := site.NewSiteBuilder(destDir, cfg.SourceDir)
			if err != nil {
				return err
			}

			if err := site.GeneratePages(cfg.SourceDir, cfg.Pages); err != nil {
				return err
			}

			return nil
		},
	}

	return cmd
}
