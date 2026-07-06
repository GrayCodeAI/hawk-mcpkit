# AGENTS.md — hawk-mcpkit

This file describes the hawk-mcpkit project for AI agents working in this
codebase. The TUI `/memory` command references this file.

---

## Project Overview

hawk-mcpkit is a shared MCP server scaffolding library for the hawk ecosystem.
It wraps [`mark3labs/mcp-go`](https://github.com/mark3labs/mcp-go) with the
construction, transports, and handler helpers that every hawk-ecosystem engine
(`inspect`, `sight`, ...) would otherwise duplicate.

Consumers declare their tools and handlers; mcpkit does the rest.

**Tagline:** Shared MCP server scaffolding for the hawk ecosystem.

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

## Architecture

```
hawk-mcpkit/
├── mcpkit.go              # Server wrapper, tool registration, transport helpers
├── mcpkit_test.go          # Tests for StrArg, JSONResult, Server construction
├── scripts/
│   └── check-ecosystem-boundaries.sh   # CI guard: zero hawk-eco deps
├── .github/
│   └── workflows/
│       ├── ci.yml                  # CI: format, module hygiene, vet, lint, test, security
│       └── release.yml             # GitHub Release on v* tags
├── Makefile                # Local dev tasks (lint, test, boundaries, ci)
├── lefthook.yml            # Pre-commit hooks (boundary guard, co-author strip)
├── AGENTS.md               # This file
├── README.md               # Public API docs, usage examples
├── CHANGELOG.md            # Keep a Changelog format
├── CODEOWNERS              # Code ownership
├── LICENSE                 # MIT
├── VERSION                 # Source of truth for versioning
└── go.mod / go.sum         # Module files
```

## Key Design Decisions

- **Zero hawk-eco dependencies:** mcpkit imports nothing from `hawk`, any
  engine, any SDK, or `hawk-core-contracts`. Its only non-stdlib dependency
  is upstream `mark3labs/mcp-go`. CI enforces this with
  `make boundaries`.
- **Library only, no binary:** mcpkit is a Go library consumed by engines. It
  does not ship a CLI or standalone binary.
- **Minimal public API:** `New()`, `AddTool()`, `MCP()`, `ServeStdio()`,
  `ServeHTTP()`, `StrArg()`, `JSONResult()`. Everything else is unexported.
- **Transport agnostic:** The `Server` struct wraps the underlying
  `mcpserver.MCPServer`. Consumers choose the transport (stdio or HTTP) via
  `ServeStdio()` or `ServeHTTP()`.
- **No ecosystem product logic:** mcpkit holds MCP server scaffolding only.
  No agent orchestration, no engine-specific behavior, no provider logic.

## Development Guidelines

### Build & Test

```bash
make test          # Run unit tests (no race detector)
make test-race     # Run unit tests with race detector
make lint          # Run golangci-lint
make boundaries    # Enforce zero hawk-eco imports
make ci            # Full CI suite (tidy, fmt, vet, boundaries, lint, test-race, security)
```

### Go Conventions

- Standard Go project layout: tests live alongside source files (`foo.go` → `foo_test.go`)
- Use table-driven tests where practical
- Errors are values — wrap with `fmt.Errorf("context: %w", err)`
- No global mutable state; prefer dependency injection

### Commit Conventions

Use [Conventional Commits](https://www.conventionalcommits.org/):
```
feat: add HTTP transport support to mcpkit
fix: handle missing tool arguments gracefully
refactor: simplify Server struct API
```

## File Organization Notes

| File | Purpose |
|------|---------|
| `mcpkit.go` | Server wrapper, tool registration, transport helpers, StrArg, JSONResult |
| `mcpkit_test.go` | Unit tests for all exported API surface |
| `scripts/check-ecosystem-boundaries.sh` | CI guard against hawk-eco imports |
| `.github/workflows/ci.yml` | CI pipeline (format, module hygiene, vet, lint, test, security, build) |
| `.github/workflows/release.yml` | GitHub Release on `v*` tags |
| `Makefile` | Local dev tasks |
| `lefthook.yml` | Pre-commit hooks (boundary guard, co-author strip) |

## Testing Patterns

- **Table-driven tests** with `t.Run(name, func(t *testing.T){...})` for all multi-case tests
- **`t.Parallel()`** on all tests that don't share mutable state
- **No mocks framework** — use concrete types and test doubles
- All tests are in `mcpkit_test.go` (single test file, library is small)

## How to Add New Transports

1. Add a method on `Server` that accepts transport-specific options
2. The method should create the appropriate server from `mcpserver` and call
   `Listen` or `Start` with the consumer-provided address or streams
3. Add tests in `mcpkit_test.go`
4. Update the README API table

## How to Add New Handler Helpers

1. Add a helper function that wraps common argument extraction or result
   marshalling patterns used across engines
2. The helper should return `(*mcplib.CallToolResult, error)` to be compatible
   with `mcpserver.ToolHandlerFunc`
3. Add tests in `mcpkit_test.go`

## Common Pitfalls

- Do not import `hawk-core-contracts` or any other hawk-eco package in mcpkit
- Do not put engine-specific logic in mcpkit (e.g., agent loop, provider
  routing, permission checking)
- The `Server` struct's unexported `mcp` field must not be accessed directly;
  use `MCP()` for the escape hatch
- `StrArg` returns `""` for missing/non-string arguments; handle this in
  your handler by returning an error
