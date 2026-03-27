package commands

import (
	"fmt"
	"os"

	"github.com/bfv/xref/internal/config"
	"github.com/bfv/xref/internal/logging"
	"github.com/spf13/cobra"
)

// NewRemoveCmd returns the remove subcommand.
func NewRemoveCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "remove",
		Short: "Remove a repository from the configuration",
		RunE: func(cmd *cobra.Command, args []string) error {
			name, _ := cmd.Flags().GetString("name")

			if name == "" && len(args) > 0 {
				name = args[0]
			}
			if name == "" {
				return fmt.Errorf("specify a repo name with --name or as argument")
			}

			cfg, err := config.NewConfig()
			if err != nil {
				return err
			}

			// Remove data file if it exists
			dataPath := cfg.RepoDataPath(name)
			if _, err := os.Stat(dataPath); err == nil {
				if err := os.Remove(dataPath); err != nil {
					logging.Logger.Warn().Err(err).Str("file", dataPath).Msg("could not remove data file")
				}
			}

			if err := cfg.RemoveRepo(name); err != nil {
				return err
			}

			if err := cfg.Save(); err != nil {
				return err
			}

			fmt.Printf("Repo '%s' removed.\n", name)
			return nil
		},
	}

	cmd.Flags().StringP("name", "n", "", "Repository name to remove")

	return cmd
}
