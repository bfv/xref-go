package commands

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/bfv/xref/internal/datafile"
	"github.com/bfv/xref/internal/logging"
	"github.com/bfv/xref/internal/parser"
	"github.com/spf13/cobra"
)

// NewParseCmd returns the parse subcommand.
func NewParseCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "parse",
		Short: "Parse .xref files in a directory",
		RunE: func(cmd *cobra.Command, args []string) error {
			dir, _ := cmd.Flags().GetString("dir")
			output, _ := cmd.Flags().GetString("output")
			srcdir, _ := cmd.Flags().GetString("srcdir")

			logging.Logger.Info().Str("dir", dir).Str("output", output).Msg("parsing xref files")

			t1 := time.Now()
			p := parser.NewParser(nil)
			xrefdata := p.ParseDir(dir, srcdir)
			elapsed := time.Since(t1)

			data, err := json.MarshalIndent(xrefdata, "", "  ")
			if err != nil {
				return fmt.Errorf("cannot marshal xref data: %w", err)
			}

			if err := os.WriteFile(output, data, 0644); err != nil {
				return fmt.Errorf("cannot write output file '%s': %w", output, err)
			}

			fmt.Printf("Parsed %d xref files in %s -> %s\n", len(xrefdata), elapsed, output)
			return nil
		},
	}

	cmd.Flags().StringP("dir", "d", ".", "Directory containing .xref files")
	cmd.Flags().StringP("output", "o", datafile.DefaultDataFile, "Output JSON file path")
	cmd.Flags().StringP("srcdir", "s", "", "Source base directory (stripped from source paths)")

	return cmd
}
