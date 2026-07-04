# hawk-mcpkit

Shared MCP server scaffolding for the hawk ecosystem.

`mcpkit` wraps [`mark3labs/mcp-go`](https://github.com/mark3labs/mcp-go) with
the construction, transports, and handler helpers that every hawk-ecosystem
library (inspect, sight, ...) would otherwise duplicate. Repos declare their
tools and handlers; mcpkit does the rest.

## Install

```sh
go get github.com/GrayCodeAI/hawk-mcpkit
```

## Usage

```go
package main

import (
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

## API

| Symbol | Purpose |
| --- | --- |
| `New(name, version)` | Create a `*Server` with tool capabilities enabled |
| `(*Server).AddTool(tool, handler)` | Register a tool and its handler |
| `(*Server).ServeStdio()` | Serve MCP over stdin/stdout |
| `(*Server).ServeHTTP(addr)` | Serve MCP over streamable HTTP at `/mcp` |
| `(*Server).MCP()` | Escape hatch to the underlying `*mcpserver.MCPServer` |
| `StrArg(req, key)` | Extract a string tool argument (`""` if absent/mistyped) |
| `JSONResult(v)` | Marshal `v` as indented JSON into a text tool result |

## License

MIT — see [LICENSE](LICENSE).
