package mcpkit

// This file adds tool search functionality to the shared MCP server
// scaffolding. Tool search allows clients to discover tools by name,
// description, or tags without knowing the exact tool name ahead of time.
//
// This is inspired by composio's tool search pattern and the MCP
// tool discovery protocol.

import (
	"strings"
	"sync"
)

// ToolInfo holds metadata about a registered tool for search purposes.
type ToolInfo struct {
	Name        string                 `json:"name"`
	Description string                 `json:"description,omitempty"`
	Tags        []string               `json:"tags,omitempty"`
	Attrs       map[string]interface{} `json:"attrs,omitempty"`
}

// ToolSearchIndex maintains an in-memory index of registered tools for
// fast name/description/tag search.
type ToolSearchIndex struct {
	mu    sync.RWMutex
	tools map[string]*ToolInfo
}

// NewToolSearchIndex creates an empty tool search index.
func NewToolSearchIndex() *ToolSearchIndex {
	return &ToolSearchIndex{
		tools: make(map[string]*ToolInfo),
	}
}

// IndexTool adds or updates a tool in the search index.
func (idx *ToolSearchIndex) IndexTool(info *ToolInfo) {
	idx.mu.Lock()
	defer idx.mu.Unlock()
	idx.tools[info.Name] = info
}

// Search queries the tool index by name, description, or tags.
// Returns tools matching any of the query terms (OR semantics).
// If query is empty, returns all tools.
func (idx *ToolSearchIndex) Search(query string) []*ToolInfo {
	idx.mu.RLock()
	defer idx.mu.RUnlock()

	if query == "" {
		result := make([]*ToolInfo, 0, len(idx.tools))
		for _, t := range idx.tools {
			result = append(result, t)
		}
		return result
	}

	q := strings.ToLower(strings.TrimSpace(query))
	terms := strings.Fields(q)

	result := make([]*ToolInfo, 0)
	for _, t := range idx.tools {
		if matchesSearch(t, terms) {
			result = append(result, t)
		}
	}
	return result
}

// GetByTag retrieves all tools with a given tag.
func (idx *ToolSearchIndex) GetByTag(tag string) []*ToolInfo {
	idx.mu.RLock()
	defer idx.mu.RUnlock()
	result := make([]*ToolInfo, 0)
	for _, t := range idx.tools {
		for _, tg := range t.Tags {
			if tg == tag {
				result = append(result, t)
				break
			}
		}
	}
	return result
}

// GetTool retrieves a tool by name from the index.
func (idx *ToolSearchIndex) GetTool(name string) (*ToolInfo, bool) {
	idx.mu.RLock()
	defer idx.mu.RUnlock()
	t, ok := idx.tools[name]
	return t, ok
}

// ListTools returns all indexed tools.
func (idx *ToolSearchIndex) ListTools() []*ToolInfo {
	idx.mu.RLock()
	defer idx.mu.RUnlock()
	result := make([]*ToolInfo, 0, len(idx.tools))
	for _, t := range idx.tools {
		result = append(result, t)
	}
	return result
}

// Count returns the number of indexed tools.
func (idx *ToolSearchIndex) Count() int {
	idx.mu.RLock()
	defer idx.mu.RUnlock()
	return len(idx.tools)
}

// matchesSearch checks if a tool matches any of the search terms.
func matchesSearch(t *ToolInfo, terms []string) bool {
	if len(terms) == 0 {
		return true
	}

	name := strings.ToLower(t.Name)
	desc := strings.ToLower(t.Description)

	for _, term := range terms {
		if strings.Contains(name, term) {
			return true
		}
		if strings.Contains(desc, term) {
			return true
		}
		for _, tag := range t.Tags {
			if strings.Contains(strings.ToLower(tag), term) {
				return true
			}
		}
	}
	return false
}

// SearchTools searches the server's registered tools by name, description,
// or tags. This is a convenience method that wraps the tool search index.
//
// Usage:
//
//	s := mcpkit.New("my-server", "1.0.0")
//	s.AddTool(tool, handler)
//	s.IndexTool(mcpkit.ToolInfo{Name: "my_tool", Description: "does X"})
//	results := s.SearchTools("X")
func (s *Server) SearchTools(query string) []*ToolInfo {
	if s.toolIndex == nil {
		return nil
	}
	return s.toolIndex.Search(query)
}

// IndexTool adds a tool to the search index.
func (s *Server) IndexTool(info *ToolInfo) {
	if s.toolIndex == nil {
		s.toolIndex = NewToolSearchIndex()
	}
	s.toolIndex.IndexTool(info)
}

// ToolSearch returns the tool search index, initializing it if needed.
func (s *Server) ToolSearch() *ToolSearchIndex {
	if s.toolIndex == nil {
		s.toolIndex = NewToolSearchIndex()
	}
	return s.toolIndex
}
