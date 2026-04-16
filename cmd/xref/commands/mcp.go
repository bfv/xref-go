package commands

import (
	"fmt"

	"github.com/bfv/xref/internal/datafile"
	mcpserver "github.com/bfv/xref/internal/mcp"
	"github.com/spf13/cobra"
)

// NewMcpCmd returns the mcp subcommand.
func NewMcpCmd(version string) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "mcp",
		Short: "Start the MCP server",
		Long:  "Start a Model Context Protocol server exposing xref data as tools.",
		RunE: func(cmd *cobra.Command, args []string) error {
			input, _ := cmd.Flags().GetString("input")
			transport, _ := cmd.Flags().GetString("transport")
			port, _ := cmd.Flags().GetInt("port")

			if transport != "stdio" && transport != "http" {
				return fmt.Errorf("invalid transport %q: must be stdio or http", transport)
			}

			return mcpserver.Run(input, transport, version, port)
		},
	}

	cmd.Flags().StringP("input", "i", datafile.DefaultDataFile, "Input JSON data file")
	cmd.Flags().StringP("transport", "t", "stdio", "Transport: stdio or http")
	cmd.Flags().IntP("port", "p", 8080, "HTTP port (only used with --transport http)")

	return cmd
}
