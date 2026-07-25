package mcpkit

import (
	"context"
	"errors"
	"strings"
	"testing"

	mcp "github.com/mark3labs/mcp-go/mcp"
)

func TestGraphResourceHandler(t *testing.T) {
	handler := graphResourceHandler(func(context.Context) (any, error) {
		return map[string]any{
			"schema_version": "graph/v1",
			"nodes":          []map[string]string{{"id": "node-1"}},
		}, nil
	})
	request := mcp.ReadResourceRequest{}
	request.Params.URI = "hawk://graph/current"

	contents, err := handler(context.Background(), request)
	if err != nil {
		t.Fatal(err)
	}
	if len(contents) != 1 {
		t.Fatalf("got %d resource contents, want 1", len(contents))
	}

	content, ok := contents[0].(*mcp.TextResourceContents)
	if !ok {
		t.Fatalf("got %T, want *mcp.TextResourceContents", contents[0])
	}
	if content.URI != request.Params.URI {
		t.Errorf("URI = %q, want %q", content.URI, request.Params.URI)
	}
	if content.MIMEType != GraphMIMEType {
		t.Errorf("MIME type = %q, want %q", content.MIMEType, GraphMIMEType)
	}
	if !strings.Contains(content.Text, `"schema_version": "graph/v1"`) {
		t.Errorf("unexpected graph payload: %s", content.Text)
	}
}

func TestGraphResourceHandlerProviderError(t *testing.T) {
	want := errors.New("graph unavailable")
	handler := graphResourceHandler(func(context.Context) (any, error) {
		return nil, want
	})

	_, err := handler(context.Background(), mcp.ReadResourceRequest{})
	if !errors.Is(err, want) {
		t.Fatalf("error = %v, want wrapped %v", err, want)
	}
}

func TestGraphResourceHandlerMarshalError(t *testing.T) {
	handler := graphResourceHandler(func(context.Context) (any, error) {
		return make(chan int), nil
	})

	if _, err := handler(context.Background(), mcp.ReadResourceRequest{}); err == nil {
		t.Fatal("expected marshal error")
	}
}

func TestAddGraphResource(t *testing.T) {
	server := New("test-server", "0.0.1")
	server.AddGraphResource("hawk://graph/current", "Current graph", func(context.Context) (any, error) {
		return map[string]string{"schema_version": "graph/v1"}, nil
	})
}
