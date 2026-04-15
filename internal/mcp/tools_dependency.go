package mcp

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/bfv/xref/internal/searcher"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

type dependencyInput struct {
	Source string `json:"source"`
}

type classHierarchyInput struct {
	Name string `json:"name"`
}

type reverseDependencyInput struct {
	Source string `json:"source"`
}

func registerDependencyTools(server *mcp.Server, s *searcher.Searcher) {
	mcp.AddTool(server, &mcp.Tool{
		Name:        "get_dependencies",
		Description: "Get all dependencies of a source: tables, includes, runs, instantiates, invokes, class/interface info",
	}, func(ctx context.Context, req *mcp.CallToolRequest, input *dependencyInput) (*mcp.CallToolResult, any, error) {
		deps := s.GetDependencies(input.Source)
		if deps == nil {
			return &mcp.CallToolResult{
				Content: []mcp.Content{&mcp.TextContent{Text: fmt.Sprintf("source %q not found", input.Source)}},
				IsError: true,
			}, nil, nil
		}
		data, _ := json.Marshal(deps)
		return &mcp.CallToolResult{
			Content: []mcp.Content{&mcp.TextContent{Text: string(data)}},
		}, nil, nil
	})

	mcp.AddTool(server, &mcp.Tool{
		Name:        "get_class_hierarchy",
		Description: "Resolve the full inheritance chain for a class or interface, walking up through all ancestors",
	}, func(ctx context.Context, req *mcp.CallToolRequest, input *classHierarchyInput) (*mcp.CallToolResult, any, error) {
		hierarchy := s.GetClassHierarchy(input.Name)
		data, _ := json.Marshal(hierarchy)
		return &mcp.CallToolResult{
			Content: []mcp.Content{&mcp.TextContent{Text: string(data)}},
		}, nil, nil
	})

	mcp.AddTool(server, &mcp.Tool{
		Name:        "get_reverse_dependencies",
		Description: "Find sources that reference a given source via includes, RUN, inheritance, invokes, or instantiation",
	}, func(ctx context.Context, req *mcp.CallToolRequest, input *reverseDependencyInput) (*mcp.CallToolResult, any, error) {
		rd := s.GetReverseDependencies(input.Source)
		if rd == nil {
			return &mcp.CallToolResult{
				Content: []mcp.Content{&mcp.TextContent{Text: fmt.Sprintf("source %q not found", input.Source)}},
				IsError: true,
			}, nil, nil
		}
		data, _ := json.Marshal(rd)
		return &mcp.CallToolResult{
			Content: []mcp.Content{&mcp.TextContent{Text: string(data)}},
		}, nil, nil
	})
}
