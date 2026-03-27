package commands

import (
	"fmt"

	"github.com/spf13/cobra"
)

// NewAboutCmd returns the about subcommand.
func NewAboutCmd(version string) *cobra.Command {
	return &cobra.Command{
		Use:   "about",
		Short: "Display information about xref",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("xref - OpenEdge XREF CLI tool")
			fmt.Println("Version:", version)
			fmt.Println("Author:  Bronco Oostermeyer <dev@bfv.io>")
			fmt.Println("License: MIT")
			fmt.Println("Repo:    https://github.com/bfv/xref-go")
		},
	}
}
