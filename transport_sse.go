package mcpkit

// This file adds SSE (Server-Sent Events) transport support to the
// shared MCP server scaffolding. SSE is the classic MCP transport used
// by Claude Desktop and other clients that don't support streamable HTTP.
//
// SSE works by:
//   1. Client GETs /sse to establish an event stream
//   2. Server sends a message endpoint URL via the first SSE event
//   3. Client POSTs JSON-RPC requests to that endpoint
//   4. Server responds over the SSE stream
//
// The SSE transport adds the same auth (bearer token / HTTP token) and
// body-size protection as the streamable HTTP transport.

import (
	"context"
	"fmt"
	"net/http"

	mcpserver "github.com/mark3labs/mcp-go/server"
)

// ServeSSE starts the MCP server on the SSE transport at
// http://<addr>/sse and blocks until the server stops.
//
// Auth precedence: if WithHTTPToken was set, every request is gated at
// the transport boundary (bearer or X-API-Key). Otherwise, if
// RequireBearerToken was set, only tool calls without a matching bearer
// header are rejected. The two modes are mutually exclusive. Every auth
// mode also caps the request body at MaxMCPRequestBodySize — the same
// protection the streamable HTTP transport applies.
//
// ServeStdio is never affected, regardless of these settings.
func (s *Server) ServeSSE(addr string) error {
	sseServer, err := s.buildSSEServer(addr)
	if err != nil {
		return err
	}

	return sseServer.Start(addr)
}

// ServeSSEWithShutdown starts the MCP server on the SSE transport at
// http://<addr>/sse and returns the underlying server so the caller can
// invoke Shutdown(ctx) for graceful teardown. It launches the listener
// in a background goroutine and returns immediately (once the server
// object exists), so the caller owns the lifecycle. See ServeSSE for
// auth and body-cap semantics — they apply identically here.
func (s *Server) ServeSSEWithShutdown(addr string) (*mcpserver.SSEServer, error) {
	sseServer, err := s.buildSSEServer(addr)
	if err != nil {
		return nil, err
	}
	s.serverStartErr = make(chan error, 1)
	go func() { s.serverStartErr <- sseServer.Start(addr) }()
	return sseServer, nil
}

// buildSSEServer constructs the SSE transport, applying the configured
// auth mode. The SSE handler is always served through an http.Server
// wrapped in the same middleware as the streamable HTTP transport:
// httpTokenHandler (auth gate + body cap) on the HTTP-token path, and
// capBodyHandler (body cap only) on the bearer-only and no-auth paths —
// this is what gives the SSE transport the auth and body-size parity
// the package documentation claims.
func (s *Server) buildSSEServer(addr string) (*mcpserver.SSEServer, error) {
	if s.bearerToken != "" && s.httpToken != "" {
		return nil, fmt.Errorf("mcpkit: cannot set both RequireBearerToken and WithHTTPToken; use at most one")
	}

	var opts []mcpserver.SSEOption

	if s.bearerToken != "" && s.httpToken == "" {
		// Bearer-only mode: attach the context func that feeds
		// bearerToolMiddleware (registered in RequireBearerToken).
		// SSEContextFunc and HTTPContextFunc have the same signature;
		// we wrap to satisfy the SSE-specific type.
		bearerCtxFunc := bearerHTTPContextFunc(s.bearerToken)
		opts = append(opts, mcpserver.WithSSEContextFunc(
			func(ctx context.Context, r *http.Request) context.Context {
				return bearerCtxFunc(ctx, r)
			},
		))
	}

	// Route the SSE handler (which implements http.Handler) through our
	// own http.Server so the auth/body-cap middleware applies to every
	// request — and so SSEServer.Start/Shutdown keep managing this exact
	// server. The delegate closes over sseServer, which is assigned right
	// after construction below; it only runs once Start is called, after
	// this function has returned, so there is no race on the variable.
	var sseServer *mcpserver.SSEServer
	delegate := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		sseServer.ServeHTTP(w, r)
	})
	var handler http.Handler = capBodyHandler(delegate)
	if s.httpToken != "" {
		handler = httpTokenHandler(s.httpToken, delegate)
	}
	opts = append(opts, mcpserver.WithHTTPServer(&http.Server{
		Addr:    addr,
		Handler: handler,
	}))

	sseServer = mcpserver.NewSSEServer(s.mcp, opts...)
	return sseServer, nil
}
