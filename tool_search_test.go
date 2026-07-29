package mcpkit

import (
	"testing"
	"time"
)

// --- Tool Search Tests ---

func TestToolSearchIndexEmpty(t *testing.T) {
	t.Parallel()
	idx := NewToolSearchIndex()
	if idx.Count() != 0 {
		t.Errorf("expected 0 tools, got %d", idx.Count())
	}
	if len(idx.ListTools()) != 0 {
		t.Error("expected empty tool list")
	}
}

func TestToolSearchIndexPutAndGet(t *testing.T) {
	t.Parallel()
	idx := NewToolSearchIndex()

	info := &ToolInfo{
		Name:        "test_tool",
		Description: "A test tool for testing",
		Tags:        []string{"test", "utility"},
	}
	idx.IndexTool(info)

	if idx.Count() != 1 {
		t.Errorf("expected 1 tool, got %d", idx.Count())
	}

	retrieved, ok := idx.GetTool("test_tool")
	if !ok {
		t.Fatal("expected to find tool")
	}
	if retrieved.Name != "test_tool" {
		t.Errorf("expected name 'test_tool', got %q", retrieved.Name)
	}
	if retrieved.Description != "A test tool for testing" {
		t.Errorf("expected description, got %q", retrieved.Description)
	}
}

func TestToolSearchIndexSearch(t *testing.T) {
	t.Parallel()
	idx := NewToolSearchIndex()

	idx.IndexTool(&ToolInfo{
		Name:        "search_files",
		Description: "Search for files in the codebase",
		Tags:        []string{"search", "files"},
	})
	idx.IndexTool(&ToolInfo{
		Name:        "list_tools",
		Description: "List available MCP tools",
		Tags:        []string{"list", "tools"},
	})
	idx.IndexTool(&ToolInfo{
		Name:        "edit_file",
		Description: "Edit a file's contents",
		Tags:        []string{"edit", "files"},
	})

	// Search by name
	results := idx.Search("search")
	if len(results) != 1 {
		t.Errorf("expected 1 result for 'search', got %d", len(results))
	}
	if results[0].Name != "search_files" {
		t.Errorf("expected 'search_files', got %q", results[0].Name)
	}

	// Search by tag
	results = idx.Search("files")
	if len(results) != 2 {
		t.Errorf("expected 2 results for 'files', got %d", len(results))
	}

	// Search by description
	results = idx.Search("available")
	if len(results) != 1 {
		t.Errorf("expected 1 result for 'available', got %d", len(results))
	}

	// Empty query returns all
	results = idx.Search("")
	if len(results) != 3 {
		t.Errorf("expected 3 results for empty query, got %d", len(results))
	}

	// No match
	results = idx.Search("nonexistent")
	if len(results) != 0 {
		t.Errorf("expected 0 results for 'nonexistent', got %d", len(results))
	}
}

func TestToolSearchIndexGetByTag(t *testing.T) {
	t.Parallel()
	idx := NewToolSearchIndex()

	idx.IndexTool(&ToolInfo{
		Name: "tool1",
		Tags: []string{"alpha", "beta"},
	})
	idx.IndexTool(&ToolInfo{
		Name: "tool2",
		Tags: []string{"beta", "gamma"},
	})
	idx.IndexTool(&ToolInfo{
		Name: "tool3",
		Tags: []string{"gamma"},
	})

	results := idx.GetByTag("beta")
	if len(results) != 2 {
		t.Errorf("expected 2 tools with tag 'beta', got %d", len(results))
	}

	results = idx.GetByTag("alpha")
	if len(results) != 1 {
		t.Errorf("expected 1 tool with tag 'alpha', got %d", len(results))
	}
}

func TestServerSearchTools(t *testing.T) {
	t.Parallel()
	s := New("test", "0.0.1")

	// No index yet
	results := s.SearchTools("test")
	if results != nil {
		t.Error("expected nil results before index creation")
	}

	// Index a tool
	s.IndexTool(&ToolInfo{
		Name:        "my_tool",
		Description: "Does something useful",
		Tags:        []string{"utility"},
	})

	results = s.SearchTools("useful")
	if len(results) != 1 {
		t.Errorf("expected 1 result, got %d", len(results))
	}
	if results[0].Name != "my_tool" {
		t.Errorf("expected 'my_tool', got %q", results[0].Name)
	}
}

func TestServerToolSearch(t *testing.T) {
	t.Parallel()
	s := New("test", "0.0.1")
	idx := s.ToolSearch()
	if idx == nil {
		t.Fatal("expected non-nil ToolSearch index")
	}
	if idx.Count() != 0 {
		t.Errorf("expected 0 tools, got %d", idx.Count())
	}
}

// --- Credential Vault Tests ---

func TestCredentialVaultEmpty(t *testing.T) {
	t.Parallel()
	v := NewCredentialVault()
	if v.Count() != 0 {
		t.Errorf("expected 0 credentials, got %d", v.Count())
	}
	if v.Get("nonexistent") != nil {
		t.Error("expected nil for nonexistent credential")
	}
}

func TestCredentialVaultStoreAndGet(t *testing.T) {
	t.Parallel()
	v := NewCredentialVault()

	c := &Credential{
		ID:    "cred-1",
		Name:  "api-key-1",
		Type:  "api_key",
		Value: "secret123",
		Scope: ScopeGlobal,
	}
	v.Store(c)

	if v.Count() != 1 {
		t.Errorf("expected 1 credential, got %d", v.Count())
	}

	retrieved := v.Get("cred-1")
	if retrieved == nil {
		t.Fatal("expected to find credential")
	}
	if retrieved.Name != "api-key-1" {
		t.Errorf("expected name 'api-key-1', got %q", retrieved.Name)
	}
	if retrieved.Value != "secret123" {
		t.Errorf("expected value 'secret123', got %q", retrieved.Value)
	}
}

func TestCredentialVaultGetByName(t *testing.T) {
	t.Parallel()
	v := NewCredentialVault()

	v.Store(&Credential{
		ID:    "cred-1",
		Name:  "github-token",
		Type:  "oauth_token",
		Value: "ghp_xxx",
		Scope: ScopeGlobal,
	})
	v.Store(&Credential{
		ID:    "cred-2",
		Name:  "github-token",
		Type:  "oauth_token",
		Value: "ghp_yyy",
		Scope: ScopeProject,
	})

	// Get by name and scope
	retrieved := v.GetByName("github-token", ScopeGlobal)
	if retrieved == nil {
		t.Fatal("expected to find global github-token")
	}
	if retrieved.Value != "ghp_xxx" {
		t.Errorf("expected value 'ghp_xxx', got %q", retrieved.Value)
	}

	// Project-scoped credential
	retrieved = v.GetByName("github-token", ScopeProject)
	if retrieved == nil {
		t.Fatal("expected to find project github-token")
	}
	if retrieved.Value != "ghp_yyy" {
		t.Errorf("expected value 'ghp_yyy', got %q", retrieved.Value)
	}

	// Nonexistent
	retrieved = v.GetByName("nonexistent", ScopeGlobal)
	if retrieved != nil {
		t.Error("expected nil for nonexistent credential")
	}
}

func TestCredentialVaultGetByTag(t *testing.T) {
	t.Parallel()
	v := NewCredentialVault()

	v.Store(&Credential{
		ID:    "cred-1",
		Name:  "key1",
		Value: "val1",
		Tags:  []string{"production", "critical"},
	})
	v.Store(&Credential{
		ID:    "cred-2",
		Name:  "key2",
		Value: "val2",
		Tags:  []string{"staging", "critical"},
	})
	v.Store(&Credential{
		ID:    "cred-3",
		Name:  "key3",
		Value: "val3",
		Tags:  []string{"dev"},
	})

	results := v.GetByTag("critical")
	if len(results) != 2 {
		t.Errorf("expected 2 credentials with tag 'critical', got %d", len(results))
	}

	results = v.GetByTag("production")
	if len(results) != 1 {
		t.Errorf("expected 1 credential with tag 'production', got %d", len(results))
	}
}

func TestCredentialVaultDelete(t *testing.T) {
	t.Parallel()
	v := NewCredentialVault()

	v.Store(&Credential{
		ID:    "cred-1",
		Name:  "key1",
		Value: "val1",
	})

	if !v.Delete("cred-1") {
		t.Error("expected Delete to return true")
	}
	if v.Count() != 0 {
		t.Errorf("expected 0 credentials after delete, got %d", v.Count())
	}

	// Delete nonexistent
	if v.Delete("nonexistent") {
		t.Error("expected Delete to return false for nonexistent")
	}
}

func TestCredentialVaultVerify(t *testing.T) {
	t.Parallel()
	v := NewCredentialVault()

	v.Store(&Credential{
		ID:    "cred-1",
		Name:  "key1",
		Value: "secret123",
	})

	if !v.Verify("cred-1", "secret123") {
		t.Error("expected Verify to return true for correct value")
	}
	if v.Verify("cred-1", "wrong") {
		t.Error("expected Verify to return false for wrong value")
	}
	if v.Verify("nonexistent", "secret123") {
		t.Error("expected Verify to return false for nonexistent credential")
	}
}

func TestCredentialVaultList(t *testing.T) {
	t.Parallel()
	v := NewCredentialVault()

	v.Store(&Credential{ID: "cred-1", Name: "key1", Value: "val1"})
	v.Store(&Credential{ID: "cred-2", Name: "key2", Value: "val2"})

	list := v.List()
	if len(list) != 2 {
		t.Errorf("expected 2 credentials, got %d", len(list))
	}
}

func TestCredentialIsExpired(t *testing.T) {
	t.Parallel()

	// No expiry
	c := &Credential{ID: "1", Value: "val"}
	if c.IsExpired() {
		t.Error("expected non-expired credential to not be expired")
	}

	// Expired
	c.ExpiresAt = time.Now().Add(-1 * time.Hour)
	if !c.IsExpired() {
		t.Error("expected expired credential to be expired")
	}
}

func TestServerVault(t *testing.T) {
	t.Parallel()
	s := New("test", "0.0.1")
	v := s.Vault()
	if v == nil {
		t.Fatal("expected non-nil vault")
	}
	if v.Count() != 0 {
		t.Errorf("expected 0 credentials, got %d", v.Count())
	}
}

func TestServerStoreAndGetCredential(t *testing.T) {
	t.Parallel()
	s := New("test", "0.0.1")

	s.StoreCredential(&Credential{
		ID:    "cred-1",
		Name:  "api-key",
		Value: "secret",
	})

	c := s.GetCredential("cred-1")
	if c == nil {
		t.Fatal("expected to find credential")
	}
	if c.Value != "secret" {
		t.Errorf("expected value 'secret', got %q", c.Value)
	}
}

// --- SSE Transport Tests ---

func TestServerServeSSENil(t *testing.T) {
	t.Parallel()
	s := New("test", "0.0.1")

	// buildSSEServer should succeed with no auth
	sseServer, err := s.buildSSEServer(":0")
	if err != nil {
		t.Fatalf("buildSSEServer failed: %v", err)
	}
	if sseServer == nil {
		t.Fatal("expected non-nil SSE server")
	}
}

func TestServerServeSSEConflictingAuth(t *testing.T) {
	t.Parallel()
	s := New("test", "0.0.1")
	s.RequireBearerToken("token1")
	s.WithHTTPToken("token2")

	_, err := s.buildSSEServer(":0")
	if err == nil {
		t.Error("expected error for conflicting auth modes")
	}
}

func TestServerServeSSEBearerOnly(t *testing.T) {
	t.Parallel()
	s := New("test", "0.0.1")
	s.RequireBearerToken("mytoken")

	sseServer, err := s.buildSSEServer(":0")
	if err != nil {
		t.Fatalf("buildSSEServer failed: %v", err)
	}
	if sseServer == nil {
		t.Fatal("expected non-nil SSE server")
	}
}

func TestServerServeSSEHTTPToken(t *testing.T) {
	t.Parallel()
	s := New("test", "0.0.1")
	s.WithHTTPToken("mytoken")

	sseServer, err := s.buildSSEServer(":0")
	if err != nil {
		t.Fatalf("buildSSEServer failed: %v", err)
	}
	if sseServer == nil {
		t.Fatal("expected non-nil SSE server")
	}
}
