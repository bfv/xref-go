package mcp

import (
	"context"
	"encoding/json"

	"github.com/bfv/xref/internal/searcher"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

type searchTableInput struct {
	Table   string `json:"table"`
	Creates *bool  `json:"creates,omitempty"`
	Updates *bool  `json:"updates,omitempty"`
	Deletes *bool  `json:"deletes,omitempty"`
}

type searchFieldInput struct {
	Field   string `json:"field"`
	Table   string `json:"table,omitempty"`
	Updates *bool  `json:"updates,omitempty"`
}

type searchDatabaseInput struct {
	Database string `json:"database"`
}

type searchIncludeInput struct {
	Include string `json:"include"`
}

type searchImplementationsInput struct {
	Interface string `json:"interface"`
}

func registerSearchTools(server *mcp.Server, s *searcher.Searcher) {
	mcp.AddTool(server, &mcp.Tool{
		Name:        "search_table_references",
		Description: "Find sources that reference a given table, with optional CRUD filters",
	}, func(ctx context.Context, req *mcp.CallToolRequest, input *searchTableInput) (*mcp.CallToolResult, any, error) {
		refs := s.GetTableReferences(input.Table, input.Creates, input.Updates, input.Deletes)
		data, _ := json.Marshal(sourcesFromXrefFiles(refs))
		return &mcp.CallToolResult{
			Content: []mcp.Content{&mcp.TextContent{Text: string(data)}},
		}, nil, nil
	})

	mcp.AddTool(server, &mcp.Tool{
		Name:        "search_field_references",
		Description: "Find sources that reference a given field, optionally within a table",
	}, func(ctx context.Context, req *mcp.CallToolRequest, input *searchFieldInput) (*mcp.CallToolResult, any, error) {
		var tablename *string
		if input.Table != "" {
			tablename = &input.Table
		}
		refs := s.GetFieldReferences(input.Field, tablename, input.Updates)
		data, _ := json.Marshal(sourcesFromXrefFiles(refs))
		return &mcp.CallToolResult{
			Content: []mcp.Content{&mcp.TextContent{Text: string(data)}},
		}, nil, nil
	})

	mcp.AddTool(server, &mcp.Tool{
		Name:        "search_database_references",
		Description: "Find sources that reference a given database",
	}, func(ctx context.Context, req *mcp.CallToolRequest, input *searchDatabaseInput) (*mcp.CallToolResult, any, error) {
		refs := s.GetDatabaseReferences(input.Database)
		data, _ := json.Marshal(sourcesFromXrefFiles(refs))
		return &mcp.CallToolResult{
			Content: []mcp.Content{&mcp.TextContent{Text: string(data)}},
		}, nil, nil
	})

	mcp.AddTool(server, &mcp.Tool{
		Name:        "search_include_references",
		Description: "Find sources that include a given file",
	}, func(ctx context.Context, req *mcp.CallToolRequest, input *searchIncludeInput) (*mcp.CallToolResult, any, error) {
		refs := s.GetIncludeReferences(input.Include)
		data, _ := json.Marshal(sourcesFromXrefFiles(refs))
		return &mcp.CallToolResult{
			Content: []mcp.Content{&mcp.TextContent{Text: string(data)}},
		}, nil, nil
	})

	mcp.AddTool(server, &mcp.Tool{
		Name:        "search_implementations",
		Description: "Find sources whose class implements a given interface",
	}, func(ctx context.Context, req *mcp.CallToolRequest, input *searchImplementationsInput) (*mcp.CallToolResult, any, error) {
		refs := s.GetImplementations(input.Interface)
		data, _ := json.Marshal(sourcesFromXrefFiles(refs))
		return &mcp.CallToolResult{
			Content: []mcp.Content{&mcp.TextContent{Text: string(data)}},
		}, nil, nil
	})
}
