package mcpkit

import (
	"context"
	"strings"
	"testing"

	mcp "github.com/mark3labs/mcp-go/mcp"
)

func TestNew(t *testing.T) {
	s := New("test-server", "0.0.1")
	if s == nil {
		t.Fatal("New returned nil")
	}
	if s.MCP() == nil {
		t.Fatal("underlying MCP server is nil")
	}
	// WithToolCapabilities enables tool support
	tool := mcp.NewTool("test", mcp.WithDescription("test"))
	s.AddTool(tool, func(_ context.Context, _ mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		return nil, nil
	})
}

func TestAddTool(t *testing.T) {
	s := New("test-server", "0.0.1")
	tool := mcp.NewTool("echo", mcp.WithDescription("echoes input"))
	s.AddTool(tool, func(_ context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		return mcp.NewToolResultText(StrArg(req, "text")), nil
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
			req := mcp.CallToolRequest{}
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
	text, ok := mcp.AsTextContent(result.Content[0])
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

func TestAddPrompt(t *testing.T) {
	s := New("test-server", "0.0.1")
	prompt := mcp.NewPrompt("hello", mcp.WithPromptDescription("Greets someone"))
	s.AddPrompt(prompt, func(_ context.Context, _ mcp.GetPromptRequest) (*mcp.GetPromptResult, error) {
		return nil, nil
	})
}

func TestAddResource(t *testing.T) {
	s := New("test-server", "0.0.1")
	resource := mcp.NewResource("info://server", "Server info")
	s.AddResource(resource, func(_ context.Context, _ mcp.ReadResourceRequest) ([]mcp.ResourceContents, error) {
		return []mcp.ResourceContents{
			&mcp.TextResourceContents{URI: "info://server", Text: "test"},
		}, nil
	})
}

func TestAddResourceTemplate(t *testing.T) {
	s := New("test-server", "0.0.1")
	tmpl := mcp.NewResourceTemplate("info://{path}", "Dynamic resource")
	s.AddResourceTemplate(tmpl, func(_ context.Context, req mcp.ReadResourceRequest) ([]mcp.ResourceContents, error) {
		uri := req.Params.URI
		return []mcp.ResourceContents{
			&mcp.TextResourceContents{URI: uri, Text: "data at " + uri},
		}, nil
	})
}

func TestServer_MCPCapabilities(t *testing.T) {
	s := New("test-server", "0.0.1")
	mcpServer := s.MCP()
	if mcpServer == nil {
		t.Fatal("MCP() returned nil")
	}
	tool := mcp.NewTool("test_tool", mcp.WithDescription("test tool"))
	s.AddTool(tool, func(_ context.Context, _ mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		return mcp.NewToolResultText("ok"), nil
	})
}

func TestStrArg_WithRequest(t *testing.T) {
	tests := []struct {
		name string
		args map[string]any
		key  string
		want string
	}{
		{name: "present string", args: map[string]any{"name": "test"}, key: "name", want: "test"},
		{name: "missing key", args: map[string]any{}, key: "name", want: ""},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			req := mcp.CallToolRequest{}
			req.Params.Arguments = tc.args
			if got := StrArg(req, tc.key); got != tc.want {
				t.Errorf("StrArg(%q) = %q, want %q", tc.key, got, tc.want)
			}
		})
	}
}
