package mcp

import (
	"context"
	"encoding/json"

	"github.com/bfv/xref/internal/searcher"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

type searchSourcesInput struct {
	Pattern string `json:"pattern"`
}

type searchRunReferencesInput struct {
	Program string `json:"program"`
}

type searchClassReferencesInput struct {
	Class    string `json:"class"`
	Detailed *bool  `json:"detailed,omitempty"`
}

func registerAdditionalTools(server *mcp.Server, s *searcher.Searcher) {
	mcp.AddTool(server, &mcp.Tool{
		Name:        "search_sources",
		Description: "Filter sources by prefix pattern (e.g. 'alg/server/' or 'alg/server/*'). Returns matching source file paths.",
	}, func(ctx context.Context, req *mcp.CallToolRequest, input *searchSourcesInput) (*mcp.CallToolResult, any, error) {
		sources := s.SearchSources(input.Pattern)
		data, _ := json.Marshal(sources)
		return &mcp.CallToolResult{
			Content: []mcp.Content{&mcp.TextContent{Text: string(data)}},
		}, nil, nil
	})

	mcp.AddTool(server, &mcp.Tool{
		Name:        "get_summary",
		Description: "Get a single-call overview of the xref dataset: source count, table count, database count, class/interface/procedure/include breakdown",
	}, func(ctx context.Context, req *mcp.CallToolRequest, input *emptyInput) (*mcp.CallToolResult, any, error) {
		summary := s.GetSummary()
		data, _ := json.Marshal(summary)
		return &mcp.CallToolResult{
			Content: []mcp.Content{&mcp.TextContent{Text: string(data)}},
		}, nil, nil
	})

	mcp.AddTool(server, &mcp.Tool{
		Name:        "search_run_references",
		Description: "Find sources that RUN a given program by name",
	}, func(ctx context.Context, req *mcp.CallToolRequest, input *searchRunReferencesInput) (*mcp.CallToolResult, any, error) {
		sources := s.GetRunReferences(input.Program)
		data, _ := json.Marshal(sources)
		return &mcp.CallToolResult{
			Content: []mcp.Content{&mcp.TextContent{Text: string(data)}},
		}, nil, nil
	})

	mcp.AddTool(server, &mcp.Tool{
		Name:        "search_class_references",
		Description: "Find sources that reference a given class by name — includes instantiation, invocation, and inheritance (inherits/implements). Returns source names only by default; set detailed=true for full breakdown per source.",
	}, func(ctx context.Context, req *mcp.CallToolRequest, input *searchClassReferencesInput) (*mcp.CallToolResult, any, error) {
		refs := s.GetClassReferences(input.Class)
		var data []byte
		if input.Detailed != nil && *input.Detailed {
			data, _ = json.Marshal(refs)
		} else {
			sources := make([]string, len(refs))
			for i, ref := range refs {
				sources[i] = ref.Source
			}
			data, _ = json.Marshal(sources)
		}
		return &mcp.CallToolResult{
			Content: []mcp.Content{&mcp.TextContent{Text: string(data)}},
		}, nil, nil
	})

	mcp.AddTool(server, &mcp.Tool{
		Name:        "list_interfaces",
		Description: "List all interface names defined in the xref dataset",
	}, func(ctx context.Context, req *mcp.CallToolRequest, input *emptyInput) (*mcp.CallToolResult, any, error) {
		interfaces := s.GetInterfaceNames()
		data, _ := json.Marshal(interfaces)
		return &mcp.CallToolResult{
			Content: []mcp.Content{&mcp.TextContent{Text: string(data)}},
		}, nil, nil
	})
}
