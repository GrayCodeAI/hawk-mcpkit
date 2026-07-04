package mcpkit

import (
	"context"
	"strings"
	"testing"

	mcplib "github.com/mark3labs/mcp-go/mcp"
)

func TestNew(t *testing.T) {
	s := New("test-server", "0.0.1")
	if s == nil {
		t.Fatal("New returned nil")
	}
	if s.MCP() == nil {
		t.Fatal("underlying MCP server is nil")
	}
}

func TestAddTool(t *testing.T) {
	s := New("test-server", "0.0.1")
	tool := mcplib.NewTool("echo", mcplib.WithDescription("echoes input"))
	s.AddTool(tool, func(ctx context.Context, req mcplib.CallToolRequest) (*mcplib.CallToolResult, error) {
		return mcplib.NewToolResultText(StrArg(req, "text")), nil
	})
}

func TestStrArg(t *testing.T) {
	tests := []struct {
		name string
		args map[string]any
		key  string
		want string
	}{
		{name: "present string", args: map[string]any{"url": "http://example.com"}, key: "url", want: "http://example.com"},
		{name: "missing key", args: map[string]any{}, key: "url", want: ""},
		{name: "nil arguments", args: nil, key: "url", want: ""},
		{name: "wrong type", args: map[string]any{"url": 42}, key: "url", want: ""},
		{name: "empty string", args: map[string]any{"url": ""}, key: "url", want: ""},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			req := mcplib.CallToolRequest{}
			req.Params.Arguments = tc.args
			if got := StrArg(req, tc.key); got != tc.want {
				t.Errorf("StrArg(%q) = %q, want %q", tc.key, got, tc.want)
			}
		})
	}
}

func TestJSONResult(t *testing.T) {
	result, err := JSONResult(map[string]any{"status": "ok", "count": 3})
	if err != nil {
		t.Fatal(err)
	}
	if result == nil || result.IsError {
		t.Fatal("expected successful result")
	}
	text, ok := mcplib.AsTextContent(result.Content[0])
	if !ok {
		t.Fatal("expected text content")
	}
	if !strings.Contains(text.Text, `"status": "ok"`) {
		t.Errorf("unexpected payload: %s", text.Text)
	}
}

func TestJSONResult_MarshalError(t *testing.T) {
	if _, err := JSONResult(make(chan int)); err == nil {
		t.Fatal("expected marshal error for unsupported type")
	}
}
