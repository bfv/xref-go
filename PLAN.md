# Port xrefparser + xrefcli (Node/TypeScript) to Go

Consolidate two Node.js repos (`bfv/xrefparser` and `bfv/xrefcli`) into a single Go CLI tool at `github.com/bfv/xref`, following the conventions in REQUIREMENTS.md.

## Source inventory

### xrefparser (library)
- **model.ts** — data types: XrefLine, Field, Table, Class, Interface, Method, Constructor, Parameter, Procedure, Run, TempTable, etc.
- **parser.ts** — `Parser` + `ParserConfig`: reads `.xref` files from a directory, parses each line, builds `XrefFile[]`
- **searcher.ts** — `Searcher`: queries parsed data (table/field/database/include/implementation references)
- **xreffile.ts** — `XrefFile`: per-source aggregation of tables, fields, classes, includes, annotations, etc.
- **util.ts** — `replaceAll` helper

### xrefcli (CLI)
- **12 commands**: about, export, init, list, matrix, parse, remove, repos, search, show, switch, validate
- **config.ts** — multi-repo config stored in `~/.xrefcli/xrefconfig.json`
- **argsparser.ts** — CLI argument parsing (replaced by Cobra)
- **help.ts** — help text (replaced by Cobra)
- **executable.ts** — command interface (`execute` + `validate`)

## Target Go project layout

```
cmd/xref/
  main.go              # Cobra root command, version flag
  commands/
    about.go
    export.go
    init.go
    list.go
    matrix.go
    parse.go
    remove.go
    repos.go
    search.go
    show.go
    switch.go
    validate.go
internal/
  parser/
    parser.go           # Parser + ParserConfig (port of parser.ts)
    parser_test.go
    xrefline.go          # XrefLine parsing (port of XrefLine.parse)
  searcher/
    searcher.go          # Searcher (port of searcher.ts)
    searcher_test.go
  models/
    models.go            # All data structs: XrefFile, Table, Field, Class, etc.
  config/
    config.go            # Repo/config management (Viper-based, port of config.ts)
  logging/
    logging.go           # zerolog setup
```

## Implementation phases

### Phase 1 — Core library (`internal/`)
1. **models.go** — Port all types from `model.ts` + `xreffile.ts` as Go structs
2. **parser.go + xrefline.go** — Port `Parser`, `ParserConfig`, `XrefLine`, and all `process*` methods
3. **parser_test.go** — Port test cases from `parser.spec.ts`; copy testcase `.xref` files from xrefparser repo
4. **searcher.go** — Port `Searcher` with all query methods
5. **searcher_test.go** — Port tests from `searcher.spec.ts`

### Phase 2 — Config & logging
6. **config.go** — Multi-repo config using Viper (`~/.xrefcli/xrefconfig.json` compat or new path)
7. **logging.go** — zerolog initialization, log-level flag support

### Phase 3 — CLI commands (`cmd/xref/`)
8. **main.go** — Cobra root command with `--version`, global flags (log-level, config path)
9. **parse.go** — `xref parse` command (calls parser, writes repo data)
10. **search.go** — `xref search` (table/field/database search flags)
11. Remaining commands: init, list, show, export, matrix, repos, switch, remove, validate, about

### Phase 4 — Build & release
12. Verify `Makefile` and `.goreleaser.yaml` work with the new code
13. Ensure `go test ./...` passes
14. CI pipeline (`.github/` workflows)

## Key design decisions
- **No 1:1 class port** — use idiomatic Go (structs + methods, no class hierarchy)
- **Cobra subcommands** replace the hand-rolled argsparser + help
- **Viper** for config file handling (JSON config compat with the existing Node format where practical)
- **zerolog** for structured logging
- **`internal/`** package for library code (parser, searcher, models) to keep it unexported
- Test data (`.xref` files) copied into `testdata/` dirs alongside test files