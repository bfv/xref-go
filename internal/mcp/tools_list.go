package mcp

import (
	"context"
	"encoding/json"

	"github.com/bfv/xref/internal/searcher"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

type emptyInput struct{}

func registerListTools(server *mcp.Server, s *searcher.Searcher) {
	mcp.AddTool(server, &mcp.Tool{
		Name:        "list_sources",
		Description: "List all source files in the xref data",
	}, func(ctx context.Context, req *mcp.CallToolRequest, input *emptyInput) (*mcp.CallToolResult, any, error) {
		sources := s.GetSourceNames()
		data, _ := json.Marshal(sources)
		return &mcp.CallToolResult{
			Content: []mcp.Content{&mcp.TextContent{Text: string(data)}},
		}, nil, nil
	})

	mcp.AddTool(server, &mcp.Tool{
		Name:        "list_tables",
		Description: "List all unique database.table references in the xref data",
	}, func(ctx context.Context, req *mcp.CallToolRequest, input *emptyInput) (*mcp.CallToolResult, any, error) {
		tables := s.GetTableNames(nil)
		data, _ := json.Marshal(tables)
		return &mcp.CallToolResult{
			Content: []mcp.Content{&mcp.TextContent{Text: string(data)}},
		}, nil, nil
	})

	mcp.AddTool(server, &mcp.Tool{
		Name:        "list_databases",
		Description: "List all unique database names in the xref data",
	}, func(ctx context.Context, req *mcp.CallToolRequest, input *emptyInput) (*mcp.CallToolResult, any, error) {
		databases := s.GetDatabaseNames(nil)
		data, _ := json.Marshal(databases)
		return &mcp.CallToolResult{
			Content: []mcp.Content{&mcp.TextContent{Text: string(data)}},
		}, nil, nil
	})
}
