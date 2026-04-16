# xref

A CLI tool for parsing and querying OpenEdge `.xref` files. It reads compiler-generated cross-reference files, extracts structured data (tables, fields, classes, includes, procedures, etc.), and writes the result as JSON. The JSON output can then be queried for table references, field usage, database dependencies, and more.

Designed to work well in CI/CD pipelines: no configuration files, no state — just files in, JSON out.

## Installation

Download a binary from the [releases](https://github.com/bfv/xref/releases) page, or use Docker:

```sh
docker run --rm -v $(pwd):/src devbfvio/xref parse
```

## Quick start

```sh
# Parse .xref files in the current directory, writes xref.json
xref parse

# Parse from a specific directory with a custom output file
xref parse --dir ./build/xref --output result.json --srcdir ./src

# Query the parsed data
xref search --table Customer
xref search --field CustNum --table Customer --updates
xref list --tables
xref list --databases
xref show
xref show --source src/Customer.cls
xref matrix
```

## Commands

### parse

Parse `.xref` files in a directory and write structured JSON output.

```
xref parse [flags]
```

| Flag | Short | Default | Description |
|------|-------|---------|-------------|
| `--dir` | `-d` | `.` | Directory containing `.xref` files (searched recursively) |
| `--output` | `-o` | `xref.json` | Output JSON file path |
| `--srcdir` | `-s` | | Source base directory (stripped from source paths in the output) |

### search

Search for table, field, or database references in the parsed data.

```
xref search [flags]
```

| Flag | Short | Default | Description |
|------|-------|---------|-------------|
| `--input` | `-i` | `xref.json` | Input JSON data file |
| `--table` | `-t` | | Table name to search for |
| `--field` | `-f` | | Field name to search for |
| `--database` | `-d` | | Database name to search for |
| `--creates` | | `false` | Filter on creates |
| `--updates` | | `false` | Filter on updates |
| `--deletes` | | `false` | Filter on deletes |

Combine `--table` and `--field` to search for a field within a specific table. Use `--creates`, `--updates`, `--deletes` to filter by CRUD operation.

### list

List all databases or tables found in the parsed data.

```
xref list [flags]
```

| Flag | Short | Default | Description |
|------|-------|---------|-------------|
| `--input` | `-i` | `xref.json` | Input JSON data file |
| `--databases` | | | List database names |
| `--tables` | | | List table names |

### show

Show parsed xref data for a specific source file, or list all source files.

```
xref show [flags]
```

| Flag | Short | Default | Description |
|------|-------|---------|-------------|
| `--input` | `-i` | `xref.json` | Input JSON data file |
| `--source` | `-s` | | Source file to show details for (lists all sources if omitted) |

### matrix

Show a source/table CRUD matrix.

```
xref matrix [flags]
```

| Flag | Short | Default | Description |
|------|-------|---------|-------------|
| `--input` | `-i` | `xref.json` | Input JSON data file |
| `--database` | `-d` | | Filter by database name |
| `--tables` | `-t` | | Comma-separated list of table names to include |
| `--tablesfile` | `-f` | | File with table names (one or more per line, comma-separated) |
| `--noreads` | `-n` | `false` | Only show tables/sources that have creates, updates or deletes |

The matrix shows each source file as a row and each table as a column. Cells contain `C` (create), `U` (update), `D` (delete), `R` (read-only), or `-` (no reference). Both rows and columns are sorted alphabetically (case-insensitive).

Use `--tables` or `--tablesfile` to limit the matrix to a specific set of tables. Only sources that reference at least one of the specified tables are included. Both options can be combined; all table names are merged.

The tables file supports one table per line, or multiple tables per line separated by commas:

```
Customer
Order,OrderLine
Item
```

Use `--noreads` to exclude tables and sources that only have read references, keeping only those with create, update or delete operations.

## MCP server

The `xref mcp` command starts a [Model Context Protocol](https://modelcontextprotocol.io/) server, exposing the parsed xref data as tools for AI agents.

```
xref mcp [flags]
```

| Flag | Short | Default | Description |
|------|-------|---------|-------------|
| `--input` | `-i` | `xref.json` | Input JSON data file |
| `--transport` | `-t` | `stdio` | Transport: `stdio` or `http` |
| `--port` | `-p` | `8080` | HTTP port (only used with `--transport http`) |

### Available tools

| Tool | Description |
|------|-------------|
| `list_sources` | List all source files |
| `list_tables` | List all database.table references |
| `list_databases` | List all database names |
| `list_interfaces` | List all interfaces |
| `search_table_references` | Find sources referencing a table (with CRUD filters) |
| `search_field_references` | Find sources referencing a field |
| `search_database_references` | Find sources referencing a database |
| `search_include_references` | Find sources that include a given file |
| `search_implementations` | Find classes implementing an interface |
| `search_sources` | Filter sources by glob/prefix pattern |
| `search_run_references` | Find sources that RUN a given program |
| `search_class_references` | Find sources that instantiate or invoke a class |
| `get_source_details` | Full details for a source file |
| `get_dependencies` | All dependencies of a source (tables, includes, runs, classes) |
| `get_class_hierarchy` | Resolve full inheritance chain for a class |
| `get_reverse_dependencies` | Find sources that depend on a given source |
| `get_migration_scope` | Transitive closure of related sources for migration analysis |
| `get_crud_matrix` | Source × table CRUD matrix for a set of sources |
| `get_summary` | Dataset overview (source/table/database counts, type breakdown) |

### VS Code configuration

Add the following to your VS Code `settings.json` (or `.vscode/mcp.json`) to use xref as an MCP server:

```json
{
  "mcp": {
    "servers": {
      "xref": {
        "command": "xref",
        "args": ["mcp", "--input", "${workspaceFolder}/xref.json"]
      }
    }
  }
}
```

For HTTP transport:

```sh
xref mcp --transport http --port 8080 --input xref.json
```

## CI/CD example

```yaml
steps:
  - name: Parse xref
    run: xref parse --dir ./build/xref --srcdir ./src --output xref.json

  - name: Check for deletes on Customer
    run: xref search --table Customer --deletes
```

## Global flags

| Flag | Default | Description |
|------|---------|-------------|
| `--log-level` | `info` | Log level (`trace`, `debug`, `info`, `warn`, `error`) |

## License

[MIT](LICENSE)
