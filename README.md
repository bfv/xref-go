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

The matrix shows each source file as a row and each table as a column. Cells contain `C` (create), `U` (update), `D` (delete), `R` (read-only), or `-` (no reference).

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
