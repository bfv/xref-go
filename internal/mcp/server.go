package mcp

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/bfv/xref/internal/datafile"
	"github.com/bfv/xref/internal/searcher"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

// Run starts the MCP server with the given configuration.
func Run(input, transport string, port int) error {
	xrefdata, err := datafile.Load(input)
	if err != nil {
		return err
	}

	s := searcher.NewSearcher(xrefdata)
	server := mcp.NewServer(&mcp.Implementation{
		Name:    "xref",
		Version: "1.0.0",
	}, nil)

	registerListTools(server, s)
	registerSearchTools(server, s)
	registerSourceTools(server, s)
	registerDependencyTools(server, s)
	registerMigrationTools(server, s)

	switch transport {
	case "stdio":
		return server.Run(context.Background(), &mcp.StdioTransport{})
	case "http":
		addr := fmt.Sprintf(":%d", port)
		handler := mcp.NewStreamableHTTPHandler(func(req *http.Request) *mcp.Server {
			return server
		}, nil)
		log.Printf("MCP server listening on %s", addr)
		return http.ListenAndServe(addr, handler)
	default:
		return fmt.Errorf("unsupported transport: %s", transport)
	}
}
