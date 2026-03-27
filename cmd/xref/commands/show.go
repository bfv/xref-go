package commands

import (
	"encoding/json"
	"fmt"

	"github.com/bfv/xref/internal/config"
	"github.com/spf13/cobra"
)

// NewShowCmd returns the show subcommand.
func NewShowCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "show",
		Short: "Show parsed xref data for a source file",
		RunE: func(cmd *cobra.Command, args []string) error {
			repoName, _ := cmd.Flags().GetString("name")
			source, _ := cmd.Flags().GetString("source")

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

			xrefdata, err := cfg.ReadRepoData(repoName)
			if err != nil {
				return err
			}

			if source == "" {
				fmt.Printf("Sources (%d):\n", len(xrefdata))
				for _, xf := range xrefdata {
					fmt.Println(" ", xf.SourceFile)
				}
				return nil
			}

			for _, xf := range xrefdata {
				if xf.SourceFile == source {
					data, err := json.MarshalIndent(xf, "", "  ")
					if err != nil {
						return err
					}
					fmt.Println(string(data))
					return nil
				}
			}

			return fmt.Errorf("source '%s' not found", source)
		},
	}

	cmd.Flags().StringP("name", "n", "", "Repository name (defaults to current)")
	cmd.Flags().StringP("source", "s", "", "Source file to show details for")

	return cmd
}
