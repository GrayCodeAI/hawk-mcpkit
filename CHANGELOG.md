# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Fixed

- SSE transport now enforces `MaxMCPRequestBodySize` on the bearer-only and no-auth paths, and `ServeSSEWithShutdown` now applies the `WithHTTPToken` gate (previously only `ServeSSE` did) — restoring the auth and body-size parity with the streamable HTTP transport that the package documentation claims. Added a regression test that POSTs an oversized body to a live SSE session.

### Removed

- Dead `MCPGraph`, `MCPNode`, `MCPEdge`, and `NewMCPGraph` graph types (`graph.go`): zero references in mcpkit or any ecosystem consumer (yaad, sight, inspect, hawk). `AddGraphResource`/`GraphMIMEType` are unaffected.
- Unused `fmt` import held alive by a `var _ = fmt.Sprintf` placeholder in `vault.go`.

### Changed

- `vault.go`'s package documentation now states the store's actual behavior (in-memory only, no persistent backing, secrets not zeroed on `Delete`) instead of implying a pluggable persistent store exists.
- `ToolSearchIndex.Search` (and therefore `Server.SearchTools`) now returns results sorted by tool name for deterministic output; added `TestToolSearchIndexSearch_SortedOrder`.
- `IndexTool` and `SearchTools` doc comments now warn loudly that the index is not kept in sync with `AddTool` — tools must be indexed explicitly to be searchable.

## [0.1.5] - 2026-07-24

### Fixed

- Request body cap (`MaxMCPRequestBodySize`) is now applied to **all** HTTP surfaces — previously it was only enforced on the `WithHTTPToken` path, leaving the bearer-only and no-auth paths unprotected against resource exhaustion.

### Added

- `WithHTTPToken` is now documented in the README Security section and API Reference.
- Tests: `constantTimeCompareStrings`, `BuildHTTPServer` mutual-exclusivity, HTTP-token wrong-token rejection, `MaxMCPRequestBodySize` value, and an intentional-skip note for `ServeStdio`.
- Removed redundant `TestStrArg_WithRequest` and `TestServer_MCPCapabilities`.

## [0.1.0] - 2026-07-04

### Added

- `Server` wrapper over `mark3labs/mcp-go` with tool capabilities enabled by default.
- `ServeStdio` and `ServeHTTP` (streamable HTTP) transports.
- `StrArg` and `JSONResult` handler helpers.
- `MCP()` escape hatch to the underlying mcp-go server.
