package commands

import (
	"fmt"

	"github.com/spf13/cobra"
)

// NewVersionCmd returns the version subcommand.
func NewVersionCmd(version string) *cobra.Command {
	return &cobra.Command{
		Use:   "version",
		Short: "Print the version number",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("xref version", version)
		},
	}
}
