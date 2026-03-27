package commands

import (
	"fmt"
	"os"

	"github.com/bfv/xref/internal/config"
	"github.com/spf13/cobra"
)

// NewValidateCmd returns the validate subcommand.
func NewValidateCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "validate",
		Short: "Validate a repository configuration",
		RunE: func(cmd *cobra.Command, args []string) error {
			repoName, _ := cmd.Flags().GetString("name")

			cfg, err := config.NewConfig()
			if err != nil {
				return err
			}

			if repoName == "" {
				repoName = cfg.Data.Current
			}
			if repoName == "" {
				return fmt.Errorf("no repo specified and no current repo set")
			}

			repo, err := cfg.GetRepo(repoName)
			if err != nil {
				return err
			}

			errors := 0

			// Check xref directory exists
			info, err := os.Stat(repo.Dir)
			if err != nil || !info.IsDir() {
				fmt.Printf("ERROR: xref directory does not exist: %s\n", repo.Dir)
				errors++
			} else {
				fmt.Printf("OK:    xref directory exists: %s\n", repo.Dir)
			}

			// Check srcdir if set
			if repo.SrcDir != "" {
				info, err := os.Stat(repo.SrcDir)
				if err != nil || !info.IsDir() {
					fmt.Printf("WARN:  source directory does not exist: %s\n", repo.SrcDir)
				} else {
					fmt.Printf("OK:    source directory exists: %s\n", repo.SrcDir)
				}
			}

			// Check data file
			dataPath := cfg.RepoDataPath(repoName)
			if _, err := os.Stat(dataPath); err != nil {
				fmt.Printf("INFO:  no parsed data found (run 'xref parse' first)\n")
			} else {
				fmt.Printf("OK:    parsed data file exists: %s\n", dataPath)
			}

			if errors > 0 {
				return fmt.Errorf("validation found %d error(s)", errors)
			}
			return nil
		},
	}

	cmd.Flags().StringP("name", "n", "", "Repository name (defaults to current)")

	return cmd
}
