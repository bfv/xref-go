package commands

import (
	"fmt"
	"time"

	"github.com/bfv/xref/internal/config"
	"github.com/bfv/xref/internal/logging"
	"github.com/bfv/xref/internal/parser"
	"github.com/spf13/cobra"
)

// NewParseCmd returns the parse subcommand.
func NewParseCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "parse",
		Short: "Parse .xref files in a repository",
		RunE: func(cmd *cobra.Command, args []string) error {
			name, _ := cmd.Flags().GetString("name")

			cfg, err := config.NewConfig()
			if err != nil {
				return err
			}

			if name == "" {
				name = cfg.Data.Current
			}
			if name == "" {
				return fmt.Errorf("no repo specified and no current repo set")
			}

			repo, err := cfg.GetRepo(name)
			if err != nil {
				return err
			}

			logging.Logger.Info().Str("repo", name).Str("dir", repo.Dir).Msg("parsing xref files")

			t1 := time.Now()
			p := parser.NewParser(nil)
			xrefdata := p.ParseDir(repo.Dir, repo.SrcDir)
			elapsed := time.Since(t1)

			if err := cfg.WriteRepoData(repo.Name, xrefdata); err != nil {
				return fmt.Errorf("error writing repo data: %w", err)
			}

			fmt.Printf("Parsed %d xref files in %s\n", len(xrefdata), elapsed)
			return nil
		},
	}

	cmd.Flags().StringP("name", "n", "", "Repository name (defaults to current)")

	return cmd
}
