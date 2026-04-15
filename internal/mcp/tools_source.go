package mcp

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/bfv/xref/internal/models"
	"github.com/bfv/xref/internal/searcher"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

type sourceInput struct {
	Source string `json:"source"`
}

func registerSourceTools(server *mcp.Server, s *searcher.Searcher) {
	mcp.AddTool(server, &mcp.Tool{
		Name:        "get_source_details",
		Description: "Get full details of a source file: class, tables, fields, includes, procedures, runs, invokes, temp-tables, annotations",
	}, func(ctx context.Context, req *mcp.CallToolRequest, input *sourceInput) (*mcp.CallToolResult, any, error) {
		xf := s.GetSourceByName(input.Source)
		if xf == nil {
			return &mcp.CallToolResult{
				Content: []mcp.Content{&mcp.TextContent{Text: fmt.Sprintf("source %q not found", input.Source)}},
				IsError: true,
			}, nil, nil
		}
		data, _ := json.Marshal(xf)
		return &mcp.CallToolResult{
			Content: []mcp.Content{&mcp.TextContent{Text: string(data)}},
		}, nil, nil
	})
}

// sourcesFromXrefFiles extracts source file paths from a slice of XrefFile.
func sourcesFromXrefFiles(xreffiles []*models.XrefFile) []string {
	sources := make([]string, 0, len(xreffiles))
	for _, xf := range xreffiles {
		sources = append(sources, xf.SourceFile)
	}
	return sources
}

// sourceMatchesCaseInsensitive checks if a source file path matches case-insensitively.
func sourceMatchesCaseInsensitive(a, b string) bool {
	return strings.EqualFold(a, b)
}
