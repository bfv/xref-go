package commands

import (
	"fmt"

	"github.com/bfv/xref/internal/config"
	"github.com/spf13/cobra"
)

// NewReposCmd returns the repos subcommand.
func NewReposCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "repos",
		Short: "List configured repositories",
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg, err := config.NewConfig()
			if err != nil {
				return err
			}

			if len(cfg.Data.Repos) == 0 {
				fmt.Println("No repositories configured.")
				return nil
			}

			for _, r := range cfg.Data.Repos {
				marker := "  "
				if r.Name == cfg.Data.Current {
					marker = "* "
				}
				fmt.Printf("%s%s\n", marker, r.Name)
				fmt.Printf("    dir:    %s\n", r.Dir)
				if r.SrcDir != "" {
					fmt.Printf("    srcdir: %s\n", r.SrcDir)
				}
			}

			return nil
		},
	}
}
