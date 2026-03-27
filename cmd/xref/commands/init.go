package commands

import (
	"fmt"

	"github.com/bfv/xref/internal/config"
	"github.com/bfv/xref/internal/logging"
	"github.com/spf13/cobra"
)

// NewInitCmd returns the init subcommand.
func NewInitCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "init",
		Short: "Initialize a new xref repository",
		RunE: func(cmd *cobra.Command, args []string) error {
			name, _ := cmd.Flags().GetString("name")
			dir, _ := cmd.Flags().GetString("dir")
			srcdir, _ := cmd.Flags().GetString("srcdir")

			if name == "" || dir == "" {
				return fmt.Errorf("--name and --dir are required")
			}

			cfg, err := config.NewConfig()
			if err != nil {
				return err
			}

			if err := cfg.AddRepo(name, dir, srcdir); err != nil {
				return err
			}

			cfg.Data.Current = name

			if err := cfg.Save(); err != nil {
				return err
			}

			logging.Logger.Info().Str("repo", name).Msg("repo initialized")
			fmt.Printf("Repo '%s' initialized and set as current.\n", name)
			return nil
		},
	}

	cmd.Flags().StringP("name", "n", "", "Repository name")
	cmd.Flags().StringP("dir", "d", "", "Directory containing .xref files")
	cmd.Flags().StringP("srcdir", "s", "", "Source base directory (stripped from source paths)")

	return cmd
}
