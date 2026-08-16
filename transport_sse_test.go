package mcpkit

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"strings"
	"testing"
	"time"
)

// TestServeSSE_BearerOnly_BodyCapRegression mirrors the streamable HTTP
// transport's MaxMCPRequestBodySize protection for the SSE transport's
// bearer-only path: an established SSE session must reject a request body
// larger than MaxMCPRequestBodySize instead of accepting it. The oversized
// payload is valid JSON-RPC, so an uncapped server would parse it and answer
// 202 Accepted; with the cap the body read fails inside mcp-go and the
// request is answered 400 Bad Request.
func TestServeSSE_BearerOnly_BodyCapRegression(t *testing.T) {
	s := New("test-server", "0.0.1")
	s.RequireBearerToken("secret-123")

	addr := "127.0.0.1:18861"
	srv, err := s.ServeSSEWithShutdown(addr)
	if err != nil {
		t.Fatalf("ServeSSEWithShutdown: %v", err)
	}
	defer func() { _ = srv.Shutdown(context.Background()) }()

	conn := &http.Client{Timeout: 10 * time.Second}

	// Establish an SSE session; the first event on the stream carries the
	// message endpoint (e.g. "/message?sessionId=<uuid>"). Keep the stream
	// open for the whole test — the session is unregistered as soon as the
	// stream disconnects, which would invalidate the endpoint.
	var sseBody io.ReadCloser
	defer func() {
		if sseBody != nil {
			_ = sseBody.Close()
		}
	}()
	var endpoint string
	for i := 0; i < 40 && endpoint == ""; i++ {
		resp, err := conn.Get("http://" + addr + "/sse")
		if err != nil {
			time.Sleep(25 * time.Millisecond)
			continue
		}
		sseBody = resp.Body
		sc := bufio.NewScanner(resp.Body)
		for sc.Scan() {
			if v, ok := strings.CutPrefix(sc.Text(), "data: "); ok {
				endpoint = strings.TrimRight(v, "\r")
				break
			}
		}
	}
	if endpoint == "" {
		t.Fatal("SSE endpoint event never received")
	}

	smallBody, _ := json.Marshal(map[string]any{
		"jsonrpc": "2.0", "id": 1, "method": "initialize",
		"params": map[string]any{
			"protocolVersion": "2025-03-26",
			"capabilities":    map[string]any{},
			"clientInfo":      map[string]any{"name": "test", "version": "0.0.1"},
		},
	})

	// Control: a small, well-formed message on the live session is accepted
	// (202). This proves any 400 below comes from the body cap, not from a
	// dead session or broken routing.
	if code := postSSEMessage(t, conn, "http://"+addr+endpoint, smallBody); code != http.StatusAccepted {
		t.Fatalf("small body: got %d, want %d (session must be valid)", code, http.StatusAccepted)
	}

	// Regression: a body over MaxMCPRequestBodySize must be rejected.
	bigBody, _ := json.Marshal(map[string]any{
		"jsonrpc": "2.0", "id": 2, "method": "initialize",
		"params": map[string]any{
			"protocolVersion": "2025-03-26",
			"capabilities":    map[string]any{},
			"clientInfo":      map[string]any{"name": "test", "version": "0.0.1"},
			"padding":         strings.Repeat("a", MaxMCPRequestBodySize),
		},
	})
	if len(bigBody) <= MaxMCPRequestBodySize {
		t.Fatalf("test bug: oversized body is only %d bytes, want > %d", len(bigBody), MaxMCPRequestBodySize)
	}
	if code := postSSEMessage(t, conn, "http://"+addr+endpoint, bigBody); code != http.StatusBadRequest {
		t.Fatalf("oversized body: got %d, want %d (MaxMCPRequestBodySize not enforced on the SSE bearer-only path)", code, http.StatusBadRequest)
	}
}

// postSSEMessage POSTs a JSON-RPC message with the bearer header and returns
// the HTTP status code.
func postSSEMessage(t *testing.T, conn *http.Client, url string, body []byte) int {
	t.Helper()
	req, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(body))
	if err != nil {
		t.Fatalf("build request: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer secret-123")
	resp, err := conn.Do(req)
	if err != nil {
		t.Fatalf("POST %s: %v", url, err)
	}
	defer func() { _ = resp.Body.Close() }()
	return resp.StatusCode
}
