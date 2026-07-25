// Package graph provides graph-based MCP tool definitions for hawk-mcpkit.
package mcpkit

import (
	"encoding/json"
	"fmt"
	"sync"
	"time"

	graphcontracts "github.com/GrayCodeAI/hawk-core-contracts/graph"
)

// MCPNode represents an MCP tool or resource as a graph node.
type MCPNode struct {
	ID          string                 `json:"id"`
	Type        string                 `json:"type"` // "tool", "resource", "prompt"
	Name        string                 `json:"name"`
	Description string                 `json:"description,omitempty"`
	URI         string                 `json:"uri,omitempty"`
	Attrs       map[string]interface{} `json:"attrs,omitempty"`
	CreatedAt   time.Time              `json:"created_at"`
	UpdatedAt   time.Time              `json:"updated_at"`
}

// MCPEdge represents a relationship between MCP tools.
type MCPEdge struct {
	From   string  `json:"from"`
	To     string  `json:"to"`
	Kind   string  `json:"kind"` // "calls", "depends_on", "produces"
	Weight float64 `json:"weight"`
}

// MCPGraph represents a graph of MCP tools and resources.
type MCPGraph struct {
	mu    sync.RWMutex
	ID    string                 `json:"id"`
	Name  string                 `json:"name"`
	Nodes map[string]*MCPNode    `json:"nodes"`
	Edges []MCPEdge              `json:"edges"`
	Attrs map[string]interface{} `json:"attrs,omitempty"`
}

// NewMCPGraph creates a new MCP graph.
func NewMCPGraph(id, name string) *MCPGraph {
	return &MCPGraph{
		ID:    id,
		Name:  name,
		Nodes: make(map[string]*MCPNode),
		Attrs: make(map[string]interface{}),
	}
}

// AddNode adds a node to the graph.
func (g *MCPGraph) AddNode(node *MCPNode) {
	g.mu.Lock()
	defer g.mu.Unlock()
	if node.CreatedAt.IsZero() {
		node.CreatedAt = time.Now()
	}
	node.UpdatedAt = time.Now()
	g.Nodes[node.ID] = node
}

// AddEdge adds an edge to the graph.
func (g *MCPGraph) AddEdge(edge MCPEdge) {
	g.mu.Lock()
	defer g.mu.Unlock()
	g.Edges = append(g.Edges, edge)
}

// GetNode retrieves a node by ID.
func (g *MCPGraph) GetNode(id string) (*MCPNode, bool) {
	g.mu.RLock()
	defer g.mu.RUnlock()
	node, ok := g.Nodes[id]
	return node, ok
}

// GetNodes returns all nodes.
func (g *MCPGraph) GetNodes() []*MCPNode {
	g.mu.RLock()
	defer g.mu.RUnlock()
	result := make([]*MCPNode, 0, len(g.Nodes))
	for _, node := range g.Nodes {
		result = append(result, node)
	}
	return result
}

// GetEdges returns all edges.
func (g *MCPGraph) GetEdges() []MCPEdge {
	g.mu.RLock()
	defer g.mu.RUnlock()
	return g.Edges
}

// FindByType finds all nodes of a specific type.
func (g *MCPGraph) FindByType(nodeType string) []*MCPNode {
	g.mu.RLock()
	defer g.mu.RUnlock()
	result := []*MCPNode{}
	for _, node := range g.Nodes {
		if node.Type == nodeType {
			result = append(result, node)
		}
	}
	return result
}

// ToJSON serializes the MCP graph to JSON.
func (g *MCPGraph) ToJSON() ([]byte, error) {
	g.mu.RLock()
	defer g.mu.RUnlock()
	return json.Marshal(g)
}

// ToGraphSpec converts the MCP graph to a portable graph spec.
func (g *MCPGraph) ToGraphSpec() *graphcontracts.GraphSpec {
	g.mu.RLock()
	defer g.mu.RUnlock()

	nodes := make([]graphcontracts.NodeSpec, 0, len(g.Nodes))
	for id, node := range g.Nodes {
		config := map[string]string{
			"type":        node.Type,
			"name":        node.Name,
			"description": node.Description,
		}
		if node.URI != "" {
			config["uri"] = node.URI
		}
		for k, v := range node.Attrs {
			config[k] = fmt.Sprintf("%v", v)
		}

		nodes = append(nodes, graphcontracts.NodeSpec{
			ID:     id,
			Type:   graphcontracts.NodeTypeTool,
			Name:   node.Name,
			Config: config,
		})
	}

	edges := make([]graphcontracts.EdgeSpec, 0, len(g.Edges))
	for _, edge := range g.Edges {
		edges = append(edges, graphcontracts.EdgeSpec{
			From:   edge.From,
			To:     edge.To,
			Weight: edge.Weight,
		})
	}

	return &graphcontracts.GraphSpec{
		ID:    g.ID,
		Name:  g.Name,
		Nodes: nodes,
		Edges: edges,
	}
}
