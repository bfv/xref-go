# Migration Strategy — Service Extraction from OpenEdge Fat Clients

## Goal

Extract services from an OpenEdge fat client codebase. These applications typically have hundreds of screens, each containing both business and UI logic. The same table (e.g., `Customer`) may be updated from multiple screens, each with slightly different validation logic. Collecting information on what and where these changes happen is essential for creating services.

## Required components

### 1. Xref MCP server (read-only) — *this repo*
What code touches what. Provides:
- Per-source: which tables it touches, CRUD flags, field-level access
- Dependency graph: includes, RUN, class hierarchy, invokes, instantiates
- CRUD matrix: for a set of sources, build a table → {source, C/R/U/D} matrix
- Migration scope: BFS/DFS graph traversal to find transitive set of related sources

### 2. Schema MCP server (read-only)
What the data looks like. Provides:
- Field types and sizes — a service API needs to know `CreditLimit` is a `DECIMAL(15,2)`
- Indexes — which fields form unique/primary keys, essential for service method signatures
- Mandatory/format constraints — `NOT NULL`, validation expressions, defaults
- Sequences — key generation strategy
- Table triggers — OpenEdge schema triggers (CREATE, WRITE, DELETE, ASSIGN) contain business logic invisible to xref; the service must absorb these

**Note:** OpenEdge does not have foreign keys or formal relationships between tables. Relationships are purely convention (e.g., `Order.CustNum` matches `Customer.CustNum` but the database doesn't enforce this). Relationships must be inferred by combining:
- Matching field names and types across tables (schema MCP)
- Index structures that suggest parent-child patterns (schema MCP)
- Tables that co-occur in the same source (xref MCP)
- Actual `FIND`, `FOR EACH`, `CAN-FIND` statements in source code (source reading)

### 3. Knowledge MCP server (read-write)
What we've learned and decided. Persists accumulated analysis results:

| Category | Examples |
|---|---|
| Inferred relationships | `Order.CustNum → Customer.CustNum` (confidence: high, evidence: 12 sources, index match) |
| Service candidates | "CustomerService" owns `Customer`, `CustFinancial`; 8 sources assigned |
| Business rules extracted | `CreditLimit` validation: 4 variants found, canonical version proposed |
| Table ownership | `Customer` → CustomerService (primary), OrderService (read-only consumer) |
| Conflict flags | `SharedTable: SalesRep` — claimed by both CustomerService and SalesService |
| Migration status | `fmCustMaint.w` — analyzed, rules extracted; `fmOrderEntry.w` — pending |

Stateless is out of the question — the analysis is iterative and cumulative. This store enables:
- Incremental work across sessions (analyze 10 screens today, continue tomorrow)
- Parallel analysis (two people work on different table domains)
- The agent asking "what do I already know about Order?" before re-analyzing

Backed by something simple (SQLite or structured JSON), but must be queryable.

### 4. Source file access (read-only)
VS Code's built-in file tools — no separate MCP needed. The agent reads `.p`, `.cls`, `.w`, `.i` files directly.

### 5. Orchestrating agent
An `.agent.md` or custom agent mode that drives the analysis workflow, tying together all MCP servers and source reading.

## Workflow for service identification (per table)

1. **Xref MCP:** `get_crud_matrix` → "Customer is updated by these 12 sources"
2. **Xref MCP:** `search_field_references("CreditLimit", "Customer", updates=true)` → narrow to the 4 that actually change credit limits
3. **Schema MCP:** "what are Customer's constraints, indexes, triggers?"
4. **Source reading:** read each of the 4 sources
5. **AI analysis:** extract the business logic around CreditLimit updates
6. **AI comparison:** compare the 4 implementations → identify canonical rules vs. variants
7. **AI proposal:** `CustomerService.UpdateCreditLimit(custNum, newLimit, reason)` with consolidated validation
8. **Knowledge MCP:** store the findings, service definition, and decision

## Key insight: shared tables

The hardest problem in service extraction is shared tables. When multiple would-be services all CRUD the same table, you need to decide on data ownership.

Typical finding in these migrations: 80% of the screens doing `UPDATE Customer` have copy-pasted or slightly drifted variants of the same validation. The service should unify that into one canonical implementation. The 20% that differ represent genuine business rule variations (e.g., different credit limit rules for domestic vs. export customers) that become parameters or separate service methods.

The field-level granularity in xref is the differentiator — knowing that screen A updates `CreditLimit` while screen B only reads it means B doesn't need the write service, just a query endpoint.

## Architecture summary

```
xref MCP (read-only)       → what code touches what
schema MCP (read-only)     → what the data looks like
knowledge MCP (read-write) → what we've learned and decided
source files (read-only)   → VS Code built-in file access
orchestrating agent        → drives the analysis workflow
```
