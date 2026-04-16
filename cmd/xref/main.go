package main

import (
	"fmt"
	"os"

	"github.com/bfv/xref/cmd/xref/commands"
	"github.com/bfv/xref/internal/logging"
	"github.com/spf13/cobra"
)

var version = "dev"

func main() {
	rootCmd := &cobra.Command{
		Use:   "xref",
		Short: "OpenEdge XREF CLI tool",
		Long:  "CLI tool for parsing and searching OpenEdge .xref files.",
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			logLevel, _ := cmd.Flags().GetString("log-level")
			logging.Setup(logLevel)
		},
	}

	rootCmd.PersistentFlags().String("log-level", "info", "Log level (trace, debug, info, warn, error)")

	rootCmd.AddCommand(commands.NewVersionCmd(version))
	rootCmd.AddCommand(commands.NewAboutCmd(version))
	rootCmd.AddCommand(commands.NewParseCmd())
	rootCmd.AddCommand(commands.NewSearchCmd())
	rootCmd.AddCommand(commands.NewListCmd())
	rootCmd.AddCommand(commands.NewShowCmd())
	rootCmd.AddCommand(commands.NewMatrixCmd())
	rootCmd.AddCommand(commands.NewMcpCmd(version))

	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
