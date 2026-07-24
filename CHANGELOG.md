# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

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
