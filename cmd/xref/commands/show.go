package commands

import (
	"encoding/json"
	"fmt"

	"github.com/bfv/xref/internal/datafile"
	"github.com/spf13/cobra"
)

// NewShowCmd returns the show subcommand.
func NewShowCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "show",
		Short: "Show parsed xref data for a source file",
		RunE: func(cmd *cobra.Command, args []string) error {
			input, _ := cmd.Flags().GetString("input")
			source, _ := cmd.Flags().GetString("source")

			xrefdata, err := datafile.Load(input)
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

	cmd.Flags().StringP("input", "i", datafile.DefaultDataFile, "Input JSON data file")
	cmd.Flags().StringP("source", "s", "", "Source file to show details for")

	return cmd
}
