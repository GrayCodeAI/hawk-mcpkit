package mcpkit

// This file adds a credential vault to the shared MCP server scaffolding.
// The vault stores API keys, tokens, and other secrets for MCP tools
// that need to call external services.
//
// The vault is in-memory by default. For production use, it can be
// backed by a persistent store (e.g., OS keychain, encrypted file).
//
// This is inspired by composio's credential management pattern.

import (
	"crypto/subtle"
	"fmt"
	"sync"
	"time"
)

// CredentialScope defines the scope at which a credential is valid.
type CredentialScope string

const (
	ScopeGlobal  CredentialScope = "global"
	ScopeProject CredentialScope = "project"
	ScopeSession CredentialScope = "session"
)

// Credential holds a single secret with metadata.
type Credential struct {
	ID        string                 `json:"id"`
	Name      string                 `json:"name"`
	Type      string                 `json:"type"` // "api_key", "oauth_token", "bearer_token", "password"
	Value     string                 `json:"-"`    // never serialized
	Scope     CredentialScope        `json:"scope"`
	ProjectID string                 `json:"project_id,omitempty"`
	SessionID string                 `json:"session_id,omitempty"`
	Tags      []string               `json:"tags,omitempty"`
	Attrs     map[string]interface{} `json:"attrs,omitempty"`
	CreatedAt time.Time              `json:"created_at"`
	UpdatedAt time.Time              `json:"updated_at"`
	ExpiresAt time.Time              `json:"expires_at,omitempty"`
}

// IsExpired reports whether the credential has expired.
func (c *Credential) IsExpired() bool {
	if c.ExpiresAt.IsZero() {
		return false
	}
	return time.Now().After(c.ExpiresAt)
}

// CredentialVault stores and manages credentials for MCP tools.
type CredentialVault struct {
	mu    sync.RWMutex
	items map[string]*Credential
}

// NewCredentialVault creates an empty credential vault.
func NewCredentialVault() *CredentialVault {
	return &CredentialVault{
		items: make(map[string]*Credential),
	}
}

// Store adds or updates a credential in the vault.
func (v *CredentialVault) Store(c *Credential) {
	v.mu.Lock()
	defer v.mu.Unlock()
	if c.CreatedAt.IsZero() {
		c.CreatedAt = time.Now()
	}
	c.UpdatedAt = time.Now()
	v.items[c.ID] = c
}

// Get retrieves a credential by ID. Returns nil if not found or expired.
func (v *CredentialVault) Get(id string) *Credential {
	v.mu.RLock()
	defer v.mu.RUnlock()
	c, ok := v.items[id]
	if !ok {
		return nil
	}
	if c.IsExpired() {
		return nil
	}
	return c
}

// GetByName retrieves a credential by name within a scope.
func (v *CredentialVault) GetByName(name string, scope CredentialScope) *Credential {
	v.mu.RLock()
	defer v.mu.RUnlock()
	for _, c := range v.items {
		if c.Name == name && c.Scope == scope && !c.IsExpired() {
			return c
		}
	}
	return nil
}

// GetByTag retrieves all credentials with a given tag.
func (v *CredentialVault) GetByTag(tag string) []*Credential {
	v.mu.RLock()
	defer v.mu.RUnlock()
	result := make([]*Credential, 0)
	for _, c := range v.items {
		if c.IsExpired() {
			continue
		}
		for _, t := range c.Tags {
			if t == tag {
				result = append(result, c)
				break
			}
		}
	}
	return result
}

// List returns all non-expired credentials.
func (v *CredentialVault) List() []*Credential {
	v.mu.RLock()
	defer v.mu.RUnlock()
	result := make([]*Credential, 0, len(v.items))
	for _, c := range v.items {
		if !c.IsExpired() {
			result = append(result, c)
		}
	}
	return result
}

// Delete removes a credential from the vault.
func (v *CredentialVault) Delete(id string) bool {
	v.mu.Lock()
	defer v.mu.Unlock()
	if _, ok := v.items[id]; !ok {
		return false
	}
	delete(v.items, id)
	return true
}

// Count returns the number of credentials in the vault.
func (v *CredentialVault) Count() int {
	v.mu.RLock()
	defer v.mu.RUnlock()
	return len(v.items)
}

// Verify checks whether a credential's value matches the expected value
// using constant-time comparison.
func (v *CredentialVault) Verify(id, expectedValue string) bool {
	c := v.Get(id)
	if c == nil {
		return false
	}
	return subtle.ConstantTimeCompare([]byte(c.Value), []byte(expectedValue)) == 1
}

// Vault returns the server's credential vault, initializing it if needed.
func (s *Server) Vault() *CredentialVault {
	if s.vault == nil {
		s.vault = NewCredentialVault()
	}
	return s.vault
}

// StoreCredential adds a credential to the server's vault.
func (s *Server) StoreCredential(c *Credential) {
	s.Vault().Store(c)
}

// GetCredential retrieves a credential by ID from the server's vault.
func (s *Server) GetCredential(id string) *Credential {
	return s.Vault().Get(id)
}

// _ ensures fmt is used (for potential future error formatting).
var _ = fmt.Sprintf
