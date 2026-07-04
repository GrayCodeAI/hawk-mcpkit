// Package mcpkit provides the shared scaffolding used by hawk-ecosystem
// libraries (inspect, sight, ...) to expose their functionality as MCP
// servers. It wraps github.com/mark3labs/mcp-go with the ecosystem's
// standard construction, transports, and small handler helpers so that
// individual repos only declare their tools and handlers.
package mcpkit

import (
	"context"
	"encoding/json"
	"os"

	mcplib "github.com/mark3labs/mcp-go/mcp"
	mcpserver "github.com/mark3labs/mcp-go/server"
)

// Server wraps an mcp-go MCPServer with the ecosystem's standard
// transports (stdio and streamable HTTP).
type Server struct {
	mcp *mcpserver.MCPServer
}

// New creates a named MCP server with tool capabilities enabled.
func New(name, version string) *Server {
	return &Server{
		mcp: mcpserver.NewMCPServer(name, version, mcpserver.WithToolCapabilities(true)),
	}
}

// AddTool registers a tool and its handler.
func (s *Server) AddTool(tool mcplib.Tool, handler mcpserver.ToolHandlerFunc) {
	s.mcp.AddTool(tool, handler)
}

// MCP returns the underlying mcp-go server, as an escape hatch for
// capabilities mcpkit does not wrap.
func (s *Server) MCP() *mcpserver.MCPServer {
	return s.mcp
}

// ServeStdio serves MCP over stdin/stdout and blocks until the stream is
// closed or the context that mcp-go derives internally is done.
func (s *Server) ServeStdio() error {
	return mcpserver.NewStdioServer(s.mcp).Listen(context.Background(), os.Stdin, os.Stdout)
}

// ServeHTTP serves MCP over the streamable HTTP transport at
// http://<addr>/mcp and blocks until the server stops.
func (s *Server) ServeHTTP(addr string) error {
	return mcpserver.NewStreamableHTTPServer(s.mcp).Start(addr)
}

// StrArg extracts a string argument from a tool call request. It returns
// "" when the argument is absent or is not a string.
func StrArg(req mcplib.CallToolRequest, key string) string {
	if v, ok := req.GetArguments()[key].(string); ok {
		return v
	}
	return ""
}

// JSONResult marshals v as indented JSON and returns it as a text tool
// result. It returns a protocol-level error only when marshalling fails.
func JSONResult(v any) (*mcplib.CallToolResult, error) {
	b, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		return nil, err
	}
	return mcplib.NewToolResultText(string(b)), nil
}
