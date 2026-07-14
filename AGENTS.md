---
description: hawk-mcpkit — MCP server scaffolding conventions.
globs: "*.go"
alwaysApply: false
---

# hawk-mcpkit Conventions

Shared MCP server scaffolding for the hawk ecosystem.

## Development workflow

When starting any new work (feature, fix, refactor, chore), always create a feature branch from `main` first. Never commit directly to `main`. Use branch naming conventions like `feat/<description>`, `fix/<description>`, or `chore/<description>`. Open a PR, ensure CI is green, then merge.

## Build & Test

```bash
go build ./...                    # Build library
go test ./...                     # Run tests
go vet ./...                      # Static analysis
make boundaries                   # Enforce ecosystem boundary rules
```

## Ecosystem Boundaries

- Zero hawk-eco dependencies (no `hawk`, no engines, no `hawk-core-contracts`)
- Only non-stdlib dependency: `mark3labs/mcp-go`
- Implementation-free of product logic

For full hawk-eco extension guidelines, see [hawk/AGENTS.md](https://github.com/GrayCodeAI/hawk/blob/main/AGENTS.md).
