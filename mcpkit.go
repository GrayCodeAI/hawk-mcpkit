// Package mcpkit provides the shared scaffolding used by hawk-ecosystem
// libraries (inspect, sight, ...) to expose their functionality as MCP
// servers. It wraps github.com/mark3labs/mcp-go with the ecosystem's
// standard construction, transports, and small handler helpers so that
// individual repos only declare their tools and handlers.
package mcpkit

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	mcp "github.com/mark3labs/mcp-go/mcp"
	mcpserver "github.com/mark3labs/mcp-go/server"
)

// Server wraps an mcp-go MCPServer with the ecosystem's standard
// transports (stdio and streamable HTTP).
type Server struct {
	mcp         *mcpserver.MCPServer
	bearerToken string
}

// New creates a named MCP server with tool, prompt, and resource
// capabilities enabled.
func New(name, version string) *Server {
	return &Server{
		mcp: mcpserver.NewMCPServer(
			name, version,
			mcpserver.WithToolCapabilities(true),
			mcpserver.WithPromptCapabilities(true),
			mcpserver.WithResourceCapabilities(false, true),
		),
	}
}

// AddTool registers a tool and its handler.
func (s *Server) AddTool(tool mcp.Tool, handler mcpserver.ToolHandlerFunc) {
	s.mcp.AddTool(tool, handler)
}

// AddPrompt registers a prompt and its handler.
func (s *Server) AddPrompt(prompt mcp.Prompt, handler mcpserver.PromptHandlerFunc) {
	s.mcp.AddPrompt(prompt, handler)
}

// AddResource registers a resource and its handler.
func (s *Server) AddResource(resource mcp.Resource, handler mcpserver.ResourceHandlerFunc) {
	s.mcp.AddResource(resource, handler)
}

// AddResourceTemplate registers a resource template and its handler.
func (s *Server) AddResourceTemplate(template mcp.ResourceTemplate, handler mcpserver.ResourceTemplateHandlerFunc) {
	s.mcp.AddResourceTemplate(template, handler)
}

// MCP returns the underlying mcp-go server, as an escape hatch for
// capabilities mcpkit does not wrap.
func (s *Server) MCP() *mcpserver.MCPServer {
	return s.mcp
}

// RequireBearerToken configures ServeHTTP and ServeSSE to reject tool
// calls that don't present a matching "Authorization: Bearer <token>"
// header. Pass "" (the default) for no auth requirement.
//
// This only gates tool calls, not resources or prompts: mcp-go's
// resource/prompt middleware are construction-time-only ServerOptions with
// no post-construction equivalent to the tool middleware's Use() method,
// so wiring them here would require restructuring New() itself. Given
// mcpkit's resource capability is already read-only
// (WithResourceCapabilities(false, true)) and tools are the primary
// capability, tool-only gating is the deliberate scope here.
//
// ServeStdio is never affected, regardless of this setting: stdio is a
// locally-spawned child process with no network exposure, not a transport
// this check is meant to protect.
func (s *Server) RequireBearerToken(token string) {
	s.bearerToken = token
}

// ServeStdio serves MCP over stdin/stdout and blocks until the stream
// closes or the context that mcp-go derives internally is done.
func (s *Server) ServeStdio() error {
	return mcpserver.ServeStdio(s.mcp)
}

// ServeHTTP serves MCP over the streamable HTTP transport at
// http://<addr>/mcp and blocks until the server stops. If
// RequireBearerToken was called with a non-empty token, tool calls without
// a matching "Authorization: Bearer <token>" header are rejected.
func (s *Server) ServeHTTP(addr string) error {
	if s.bearerToken == "" {
		return mcpserver.NewStreamableHTTPServer(s.mcp).Start(addr)
	}
	s.mcp.Use(bearerToolMiddleware(s.bearerToken))
	httpServer := mcpserver.NewStreamableHTTPServer(
		s.mcp,
		mcpserver.WithHTTPContextFunc(bearerHTTPContextFunc(s.bearerToken)),
	)
	return httpServer.Start(addr)
}

// ServeSSE serves MCP over the SSE transport at <addr> and blocks until
// the server stops. If RequireBearerToken was called with a non-empty
// token, tool calls without a matching "Authorization: Bearer <token>"
// header are rejected.
func (s *Server) ServeSSE(addr string) error {
	if s.bearerToken == "" {
		return mcpserver.NewSSEServer(s.mcp).Start(addr)
	}
	s.mcp.Use(bearerToolMiddleware(s.bearerToken))
	sseServer := mcpserver.NewSSEServer(
		s.mcp,
		mcpserver.WithSSEContextFunc(bearerSSEContextFunc(s.bearerToken)),
	)
	return sseServer.Start(addr)
}

type bearerAuthorizedKey struct{}

func checkBearer(r *http.Request, token string) bool {
	return r.Header.Get("Authorization") == "Bearer "+token
}

// bearerHTTPContextFunc validates the incoming request's Authorization
// header and stashes the result in context for bearerToolMiddleware to
// check. It never rejects the request itself — mcp-go's HTTPContextFunc
// has no way to do that — enforcement happens in the tool middleware.
func bearerHTTPContextFunc(token string) mcpserver.HTTPContextFunc {
	return func(ctx context.Context, r *http.Request) context.Context {
		return context.WithValue(ctx, bearerAuthorizedKey{}, checkBearer(r, token))
	}
}

func bearerSSEContextFunc(token string) mcpserver.SSEContextFunc {
	return func(ctx context.Context, r *http.Request) context.Context {
		return context.WithValue(ctx, bearerAuthorizedKey{}, checkBearer(r, token))
	}
}

// bearerToolMiddleware rejects a tool call unless bearerHTTPContextFunc /
// bearerSSEContextFunc marked its context as authorized. A plain Go error
// return (rather than a *mcp.CallToolResult) mirrors mcp-go's own
// WithRecovery middleware, which mcp-go turns into a protocol-level error
// response.
func bearerToolMiddleware(token string) mcpserver.ToolHandlerMiddleware {
	return func(next mcpserver.ToolHandlerFunc) mcpserver.ToolHandlerFunc {
		return func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			authorized, _ := ctx.Value(bearerAuthorizedKey{}).(bool)
			if !authorized {
				return nil, fmt.Errorf("unauthorized: missing or invalid bearer token")
			}
			return next(ctx, req)
		}
	}
}

// StrArg extracts a string argument from a tool call request. Returns
// "" when absent or not a string.
func StrArg(req mcp.CallToolRequest, key string) string {
	if v, ok := req.GetArguments()[key].(string); ok {
		return v
	}
	return ""
}

// JSONResult marshals v as indented JSON and returns it as a text tool
// result. Returns a protocol-level error only when marshalling fails.
func JSONResult(v any) (*mcp.CallToolResult, error) {
	b, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		return nil, err
	}
	return mcp.NewToolResultText(string(b)), nil
}
