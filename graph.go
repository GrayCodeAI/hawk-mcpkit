package mcpkit

import (
	"context"
	"encoding/json"
	"fmt"

	mcp "github.com/mark3labs/mcp-go/mcp"
	mcpserver "github.com/mark3labs/mcp-go/server"
)

// GraphMIMEType identifies a Hawk ecosystem graph JSON document exposed as
// an MCP resource. The producing repository remains responsible for the
// document's schema, authorization, redaction, and size bounds.
const GraphMIMEType = "application/vnd.hawk.graph+json"

// GraphResourceProvider supplies the current graph document for a read-only
// MCP resource. It deliberately returns any rather than importing an
// ecosystem graph contract: mcpkit is transport scaffolding and has no
// hawk-eco dependencies.
type GraphResourceProvider func(context.Context) (any, error)

// AddGraphResource registers a read-only JSON graph resource. Use
// WithHTTPToken before serving over HTTP when the graph contains data that
// must not be publicly discoverable; RequireBearerToken gates tools only.
func (s *Server) AddGraphResource(uri, name string, provider GraphResourceProvider) {
	resource := mcp.NewResource(
		uri,
		name,
		mcp.WithMIMEType(GraphMIMEType),
		mcp.WithResourceDescription("Read-only Hawk graph document"),
	)
	s.AddResource(resource, graphResourceHandler(provider))
}

func graphResourceHandler(provider GraphResourceProvider) mcpserver.ResourceHandlerFunc {
	return func(ctx context.Context, req mcp.ReadResourceRequest) ([]mcp.ResourceContents, error) {
		value, err := provider(ctx)
		if err != nil {
			return nil, fmt.Errorf("provide graph resource: %w", err)
		}

		payload, err := json.MarshalIndent(value, "", "  ")
		if err != nil {
			return nil, fmt.Errorf("marshal graph resource: %w", err)
		}

		return []mcp.ResourceContents{
			&mcp.TextResourceContents{
				URI:      req.Params.URI,
				MIMEType: GraphMIMEType,
				Text:     string(payload),
			},
		}, nil
	}
}
