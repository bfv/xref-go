
## Implementation phases

### Phase 1 — Minimal viable MCP server
1. Add Go MCP SDK dependency (`github.com/modelcontextprotocol/go-sdk`)
2. Create `cmd/xref/commands/mcp.go` — Cobra command with `--input`, `--transport`, `--port` flags
3. Create `internal/mcp/server.go` — server bootstrap, loads xref.json, initializes Searcher, registers tools, starts stdio or HTTP transport
4. Implement Layer 1 tools (`list_sources`, `list_tables`, `list_databases`)
5. Implement Layer 2 tools (all `search_*` tools)
6. Implement Layer 3 tool (`get_source_details`)
7. Test with VS Code MCP client (stdio) and manual HTTP

### Phase 2 — Dependency analysis
8. Add `GetDependencies()` to Searcher — aggregates tables, includes, runs, instantiates, invokes for a source
9. Add `GetClassHierarchy()` to Searcher — resolves full inheritance chain across all sources
10. Add `GetReverseDependencies()` to Searcher — finds sources that reference a given source via includes, RUN, inheritance
11. Implement Layer 4 tool handlers
12. Tests for new Searcher methods

### Phase 3 — Migration tools
13. Add `GetMigrationScope()` to Searcher — BFS/DFS graph traversal: starting from a source, follow shared tables, class hierarchy, and include chains to find the transitive set of related sources
14. Add `GetCrudMatrix()` to Searcher — for a set of sources, build a table → {source, C/R/U/D} matrix
15. Implement Layer 5 tool handlers
16. Tests for migration methods

### Phase 4 — Polish
17. Add MCP server info (name, version) from build vars
18. Update README with MCP usage and VS Code config example
19. Verify both transports work end-to-end
20. Add to Makefile / goreleaser

## Key decisions
- **Single binary** — `xref mcp` is a subcommand, not a separate binary
- **Read-only** — the MCP server never modifies xref.json
- **In-memory** — xref.json is loaded once at startup; no persistence needed
- **Tool responses as JSON** — all tool handlers return structured JSON, not plain text, so AI agents can parse reliably
- **`go-sdk`** — uses the official `github.com/modelcontextprotocol/go-sdk` (maintained with Google); type-safe struct-based tool handlers, supports both stdio and streamable HTTP