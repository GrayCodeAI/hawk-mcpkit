# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [0.1.0] - 2026-07-04

### Added

- `Server` wrapper over `mark3labs/mcp-go` with tool capabilities enabled by default.
- `ServeStdio` and `ServeHTTP` (streamable HTTP) transports.
- `StrArg` and `JSONResult` handler helpers.
- `MCP()` escape hatch to the underlying mcp-go server.
