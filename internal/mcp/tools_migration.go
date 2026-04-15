package mcp

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/bfv/xref/internal/searcher"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

type migrationScopeInput struct {
	Source string `json:"source"`
}

type crudMatrixInput struct {
	Sources []string `json:"sources,omitempty"`
}

func registerMigrationTools(server *mcp.Server, s *searcher.Searcher) {
	mcp.AddTool(server, &mcp.Tool{
		Name:        "get_migration_scope",
		Description: "Starting from a source, follow shared tables, class hierarchy, and include chains to find the transitive set of related sources that would need to be migrated together",
	}, func(ctx context.Context, req *mcp.CallToolRequest, input *migrationScopeInput) (*mcp.CallToolResult, any, error) {
		scope := s.GetMigrationScope(input.Source)
		if scope == nil {
			return &mcp.CallToolResult{
				Content: []mcp.Content{&mcp.TextContent{Text: fmt.Sprintf("source %q not found", input.Source)}},
				IsError: true,
			}, nil, nil
		}
		data, _ := json.Marshal(scope)
		return &mcp.CallToolResult{
			Content: []mcp.Content{&mcp.TextContent{Text: string(data)}},
		}, nil, nil
	})

	mcp.AddTool(server, &mcp.Tool{
		Name:        "get_crud_matrix",
		Description: "Build a table-to-source CRUD matrix showing which sources create, read, update, or delete which tables. Pass specific sources to filter, or omit for all sources.",
	}, func(ctx context.Context, req *mcp.CallToolRequest, input *crudMatrixInput) (*mcp.CallToolResult, any, error) {
		var sources []string
		if len(input.Sources) > 0 {
			sources = input.Sources
		}
		matrix := s.GetCrudMatrix(sources)
		data, _ := json.Marshal(matrix)
		return &mcp.CallToolResult{
			Content: []mcp.Content{&mcp.TextContent{Text: string(data)}},
		}, nil, nil
	})
}
