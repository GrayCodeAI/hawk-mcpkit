# hawk-mcpkit

Shared MCP server scaffolding for the hawk ecosystem.

`mcpkit` wraps [`mark3labs/mcp-go`](https://github.com/mark3labs/mcp-go) with
the construction, transports, and handler helpers that every hawk-ecosystem
engine (`inspect`, `sight`, ...) would otherwise duplicate. Repos declare their
tools and handlers; mcpkit does the rest.

**Tagline:** Shared MCP server scaffolding for the hawk ecosystem.

## Install

```sh
go get github.com/GrayCodeAI/hawk-mcpkit
```

## Usage

```go
package main

import (
	"context"

	mcpkit "github.com/GrayCodeAI/hawk-mcpkit"
	mcplib "github.com/mark3labs/mcp-go/mcp"
)

func main() {
	s := mcpkit.New("mytool", "0.1.0")

	s.AddTool(
		mcplib.NewTool("greet",
			mcplib.WithDescription("Greets a person by name."),
			mcplib.WithString("name", mcplib.Required(), mcplib.Description("Who to greet")),
		),
		func(ctx context.Context, req mcplib.CallToolRequest) (*mcplib.CallToolResult, error) {
			name := mcpkit.StrArg(req, "name")
			if name == "" {
				return mcplib.NewToolResultError("name is required"), nil
			}
			return mcpkit.JSONResult(map[string]string{"greeting": "hello " + name})
		},
	)

	// stdio transport:
	_ = s.ServeStdio()
	// or streamable HTTP at http://localhost:8080/mcp:
	// _ = s.ServeHTTP(":8080")
}
```

## Quick Reference

| Task | Code |
|------|------|
| Create server | `s := mcpkit.New("name", "0.1.0")` |
| Add a tool | `s.AddTool(tool, handler)` |
| Serve stdio | `s.ServeStdio()` |
| Serve HTTP | `s.ServeHTTP(":8080")` |
| Extract string arg | `mcpkit.StrArg(req, "key")` |
| Return JSON result | `mcpkit.JSONResult(map[string]any{...})` |

## Architecture

```
hawk-mcpkit Server
├── wraps mark3labs/mcp-go MCPServer
├── AddTool() registers tools + handlers
├── ServeStdio() → stdin/stdout transport
├── ServeHTTP(addr) → streamable HTTP at /mcp
├── StrArg() → extract string arguments
└── JSONResult() → marshal values as JSON text results
```

## API Reference

### Server

| Symbol | Purpose |
|--------|---------|
| `New(name, version)` | Create a `*Server` with tool capabilities enabled. Returns `*Server`. |
| `(*Server).AddTool(tool, handler)` | Register a tool and its handler. `handler` is `func(context.Context, mcplib.CallToolRequest) (*mcplib.CallToolResult, error)`. |
| `(*Server).ServeStdio()` | Serve MCP over stdin/stdout. Blocks until stream closes. Returns `error`. |
| `(*Server).ServeHTTP(addr)` | Serve MCP over streamable HTTP at `/mcp`. Blocks until server stops. Returns `error`. |
| `(*Server).MCP()` | Escape hatch to the underlying `*mcpserver.MCPServer`. Use only for capabilities mcpkit does not wrap. |

### Handler Helpers

| Symbol | Purpose |
|--------|---------|
| `StrArg(req, key)` | Extract a string argument from a tool call request. Returns `""` when absent or not a string. |
| `JSONResult(v)` | Marshal `v` as indented JSON and return it as a text tool result. Returns `(*mcplib.CallToolResult, error)`. Error only when marshalling fails. |

## Ecosystem

hawk-mcpkit is a **foundation repo** in the hawk-eco mono-ecosystem:

| Component | Purpose |
|-----------|---------|
| **hawk-mcpkit** | Shared MCP server scaffolding (this repo) |
| **hawk-core-contracts** | Shared cross-repo contracts (types, tools, events, policy, review, verify, sessions) |
| **eyrie** | LLM provider runtime — routing, streaming, retries, caching |
| **yaad** | Graph-based persistent memory for coding agents |
| **tok** | Tokenizer, compression, secrets scanning, rate limiting |
| **sight** | Diff-based code review and static analysis |
| **inspect** | Security audit library (CVE, API security, CI output) |
| **trace** | Session capture and replay CLI |
| **hawk** | AI coding agent (this repo) |

Engines that serve MCP (`sight`, `inspect`) import `hawk-mcpkit`; it never
imports them back.

## Ecosystem Boundaries

Rules that keep this repo at the foundation layer:

- **Zero hawk-eco dependencies.** This repo must never import `hawk`, any
  engine (`eyrie`, `yaad`, `tok`, `trace`, `sight`, `inspect`), any SDK, or
  `hawk-core-contracts`. Its only non-stdlib dependency is upstream
  `mark3labs/mcp-go`. `make boundaries` (also run in CI) enforces this with
  `scripts/check-ecosystem-boundaries.sh`.
- **Implementation-free of product logic.** This repo holds MCP server
  scaffolding only — no hawk orchestration, no engine-specific behavior, no
  provider logic.
- **Consumers, not dependents.** Engines that serve MCP (`sight`, `inspect`)
  import `hawk-mcpkit`; it never imports them back.

If you need a hawk-ecosystem type here, that's a sign it belongs in the
consuming engine instead, not in this repo.

## License

MIT — see [LICENSE](LICENSE).
