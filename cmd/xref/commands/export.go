package commands

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/bfv/xref/internal/config"
	"github.com/bfv/xref/internal/logging"
	"github.com/spf13/cobra"
)

// NewExportCmd returns the export subcommand.
func NewExportCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "export",
		Short: "Export parsed xref data to a JSON file",
		RunE: func(cmd *cobra.Command, args []string) error {
			repoName, _ := cmd.Flags().GetString("name")
			output, _ := cmd.Flags().GetString("output")

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

			data, err := json.MarshalIndent(xrefdata, "", "  ")
			if err != nil {
				return fmt.Errorf("cannot marshal data: %w", err)
			}

			if output == "" {
				fmt.Println(string(data))
				return nil
			}

			if err := os.WriteFile(output, data, 0644); err != nil {
				return fmt.Errorf("cannot write file: %w", err)
			}

			logging.Logger.Info().Str("file", output).Int("sources", len(xrefdata)).Msg("exported")
			fmt.Printf("Exported %d sources to %s\n", len(xrefdata), output)
			return nil
		},
	}

	cmd.Flags().StringP("name", "n", "", "Repository name (defaults to current)")
	cmd.Flags().StringP("output", "o", "", "Output file path (stdout if omitted)")

	return cmd
}
