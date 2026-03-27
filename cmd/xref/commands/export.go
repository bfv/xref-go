package commands

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/bfv/xref/internal/datafile"
	"github.com/bfv/xref/internal/logging"
	"github.com/spf13/cobra"
)

// NewExportCmd returns the export subcommand.
func NewExportCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "export",
		Short: "Export parsed xref data to a JSON file",
		RunE: func(cmd *cobra.Command, args []string) error {
			input, _ := cmd.Flags().GetString("input")
			output, _ := cmd.Flags().GetString("output")

			xrefdata, err := datafile.Load(input)
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

	cmd.Flags().StringP("input", "i", datafile.DefaultDataFile, "Input JSON data file")
	cmd.Flags().StringP("output", "o", "", "Output file path (stdout if omitted)")

	return cmd
}
