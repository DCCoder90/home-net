package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"
)

// Client talks to the Authentik REST API v3.
type Client struct {
	baseURL    string
	token      string
	httpClient *http.Client
}

// New creates an Authentik API client using a service account token.
func New(baseURL, token string) *Client {
	return &Client{
		baseURL:    baseURL,
		token:      token,
		httpClient: &http.Client{Timeout: 30 * time.Second},
	}
}

// --- Data source lookups (read-only) ---

// Flow represents an Authentik flow.
type Flow struct {
	PK   string `json:"pk"`
	Slug string `json:"slug"`
	Name string `json:"name"`
}

// GetFlowBySlug finds a flow by its slug.
func (c *Client) GetFlowBySlug(slug string) (*Flow, error) {
	params := url.Values{"slug": []string{slug}}
	var result paginatedResult[Flow]
	if err := c.get("/api/v3/flows/instances/?" + params.Encode(), &result); err != nil {
		return nil, err
	}
	if len(result.Results) == 0 {
		return nil, fmt.Errorf("flow with slug %q not found", slug)
	}
	return &result.Results[0], nil
}

// ScopeMapping represents a property mapping for OAuth2 scopes.
type ScopeMapping struct {
	PK   string `json:"pk"`
	Name string `json:"name"`
}

// GetScopeMappingByName finds a scope property mapping by name.
func (c *Client) GetScopeMappingByName(name string) (*ScopeMapping, error) {
	params := url.Values{"name": []string{name}}
	var result paginatedResult[ScopeMapping]
	if err := c.get("/api/v3/propertymappings/provider/scope/?"+params.Encode(), &result); err != nil {
		return nil, err
	}
	if len(result.Results) == 0 {
		return nil, fmt.Errorf("scope mapping %q not found", name)
	}
	return &result.Results[0], nil
}

// CertificateKeyPair represents an Authentik certificate.
type CertificateKeyPair struct {
	PK   string `json:"pk"`
	Name string `json:"name"`
}

// GetCertificateByName finds a certificate key pair by name.
func (c *Client) GetCertificateByName(name string) (*CertificateKeyPair, error) {
	params := url.Values{"name": []string{name}}
	var result paginatedResult[CertificateKeyPair]
	if err := c.get("/api/v3/crypto/certificatekeypairs/?"+params.Encode(), &result); err != nil {
		return nil, err
	}
	if len(result.Results) == 0 {
		return nil, fmt.Errorf("certificate %q not found", name)
	}
	return &result.Results[0], nil
}

// User represents an Authentik user.
type User struct {
	PK       int    `json:"pk"`
	Username string `json:"username"`
}

// GetUserByUsername finds a user by username.
func (c *Client) GetUserByUsername(username string) (*User, error) {
	params := url.Values{"username": []string{username}}
	var result paginatedResult[User]
	if err := c.get("/api/v3/core/users/?"+params.Encode(), &result); err != nil {
		return nil, err
	}
	if len(result.Results) == 0 {
		return nil, fmt.Errorf("user %q not found", username)
	}
	return &result.Results[0], nil
}

// --- Resource CRUD ---

// Application represents an Authentik application.
type Application struct {
	Slug string `json:"slug"`
	Name string `json:"name"`
	UUID string `json:"pk"`
}

// CreateApplicationInput is the payload for creating an application.
type CreateApplicationInput struct {
	Name             string `json:"name"`
	Slug             string `json:"slug"`
	MetaDescription  string `json:"meta_description,omitempty"`
	Group            string `json:"group,omitempty"`
	ProtocolProvider int    `json:"provider,omitempty"`
}

// CreateApplication creates an application and returns it.
func (c *Client) CreateApplication(input CreateApplicationInput) (*Application, error) {
	var out Application
	return &out, c.post("/api/v3/core/applications/", input, &out)
}

// DeleteApplication removes an application by slug.
func (c *Client) DeleteApplication(slug string) error {
	return c.delete("/api/v3/core/applications/" + slug + "/")
}

// ProviderProxy represents an Authentik proxy provider.
type ProviderProxy struct {
	PK   int    `json:"pk"`
	Name string `json:"name"`
}

// CreateProviderProxyInput is the payload for creating a proxy provider.
type CreateProviderProxyInput struct {
	Name                       string `json:"name"`
	InternalHost               string `json:"internal_host"`
	ExternalHost               string `json:"external_host"`
	BasicAuthEnabled           bool   `json:"basic_auth_enabled"`
	BasicAuthUsernameAttribute string `json:"basic_auth_username_attribute,omitempty"`
	BasicAuthPasswordAttribute string `json:"basic_auth_password_attribute,omitempty"`
	AuthorizationFlow          string `json:"authorization_flow"`
	InvalidationFlow           string `json:"invalidation_flow"`
}

// CreateProviderProxy creates a proxy provider and returns its PK.
func (c *Client) CreateProviderProxy(input CreateProviderProxyInput) (*ProviderProxy, error) {
	var out ProviderProxy
	return &out, c.post("/api/v3/providers/proxy/", input, &out)
}

// DeleteProviderProxy removes a proxy provider by PK.
func (c *Client) DeleteProviderProxy(pk int) error {
	return c.delete(fmt.Sprintf("/api/v3/providers/proxy/%d/", pk))
}

// ProviderOAuth2 represents an Authentik OAuth2 provider.
type ProviderOAuth2 struct {
	PK                 int    `json:"pk"`
	Name               string `json:"name"`
	ClientID           string `json:"client_id"`
	ClientSecret       string `json:"client_secret"`
	AuthorizationFlow  string `json:"authorization_flow"`
}

// OAuth2ProviderConfig holds the OIDC discovery URLs for an OAuth2 provider.
type OAuth2ProviderConfig struct {
	Issuer        string `json:"issuer"`
	AuthorizationEndpoint string `json:"authorization_endpoint"`
	TokenEndpoint string `json:"token_endpoint"`
	UserinfoEndpoint string `json:"userinfo_endpoint"`
	JwksURI       string `json:"jwks_uri"`
	EndSessionEndpoint string `json:"end_session_endpoint"`
}

// AllowedRedirectURI is a single redirect URI entry.
type AllowedRedirectURI struct {
	MatchingMode string `json:"matching_mode"`
	URL          string `json:"url"`
}

// CreateProviderOAuth2Input is the payload for creating an OAuth2 provider.
type CreateProviderOAuth2Input struct {
	Name                string               `json:"name"`
	ClientID            string               `json:"client_id"`
	AuthorizationFlow   string               `json:"authorization_flow"`
	InvalidationFlow    string               `json:"invalidation_flow"`
	SigningKey           string               `json:"signing_key,omitempty"`
	AllowedRedirectURIs []AllowedRedirectURI `json:"redirect_uris"`
	PropertyMappings    []string             `json:"property_mappings,omitempty"`
}

// CreateProviderOAuth2 creates an OAuth2 provider.
func (c *Client) CreateProviderOAuth2(input CreateProviderOAuth2Input) (*ProviderOAuth2, error) {
	var out ProviderOAuth2
	return &out, c.post("/api/v3/providers/oauth2/", input, &out)
}

// GetProviderOAuth2Config returns the OIDC discovery config for a provider.
func (c *Client) GetProviderOAuth2Config(pk int) (*OAuth2ProviderConfig, error) {
	// Fetch the OIDC discovery document via well-known URL pattern
	var out struct {
		Issuer string `json:"issuer_url"`
	}
	if err := c.get(fmt.Sprintf("/api/v3/providers/oauth2/%d/", pk), &out); err != nil {
		return nil, err
	}

	cfg := &OAuth2ProviderConfig{
		Issuer:                out.Issuer,
		AuthorizationEndpoint: out.Issuer + "/authorize/",
		TokenEndpoint:         out.Issuer + "/token/",
		UserinfoEndpoint:      out.Issuer + "/userinfo/",
		JwksURI:               out.Issuer + "/jwks/",
		EndSessionEndpoint:    out.Issuer + "/end-session/",
	}
	return cfg, nil
}

// DeleteProviderOAuth2 removes an OAuth2 provider by PK.
func (c *Client) DeleteProviderOAuth2(pk int) error {
	return c.delete(fmt.Sprintf("/api/v3/providers/oauth2/%d/", pk))
}

// Group represents an Authentik group.
type Group struct {
	PK   string `json:"pk"`
	Name string `json:"name"`
}

// CreateGroupInput is the payload for creating a group.
type CreateGroupInput struct {
	Name       string                 `json:"name"`
	Users      []int                  `json:"users,omitempty"`
	Attributes map[string]interface{} `json:"attributes,omitempty"`
}

// CreateGroup creates a group.
func (c *Client) CreateGroup(input CreateGroupInput) (*Group, error) {
	var out Group
	return &out, c.post("/api/v3/core/groups/", input, &out)
}

// UpdateGroup replaces a group's attributes by PK.
func (c *Client) UpdateGroup(pk string, input CreateGroupInput) error {
	return c.put(fmt.Sprintf("/api/v3/core/groups/%s/", pk), input, nil)
}

// DeleteGroup removes a group by PK.
func (c *Client) DeleteGroup(pk string) error {
	return c.delete(fmt.Sprintf("/api/v3/core/groups/%s/", pk))
}

// PolicyBinding represents an Authentik policy binding.
type PolicyBinding struct {
	PK    string `json:"pk"`
	Group string `json:"group"`
}

// CreatePolicyBindingInput is the payload for creating a policy binding.
type CreatePolicyBindingInput struct {
	Target string `json:"target"` // application UUID
	Group  string `json:"group"`  // group PK
	Order  int    `json:"order"`
}

// CreatePolicyBinding creates a policy binding.
func (c *Client) CreatePolicyBinding(input CreatePolicyBindingInput) (*PolicyBinding, error) {
	var out PolicyBinding
	return &out, c.post("/api/v3/policies/bindings/", input, &out)
}

// DeletePolicyBinding removes a policy binding by PK.
func (c *Client) DeletePolicyBinding(pk string) error {
	return c.delete(fmt.Sprintf("/api/v3/policies/bindings/%s/", pk))
}

// --- HTTP helpers ---

type paginatedResult[T any] struct {
	Results []T `json:"results"`
}

func (c *Client) get(path string, out interface{}) error {
	req, _ := http.NewRequest(http.MethodGet, c.baseURL+path, nil)
	req.Header.Set("Authorization", "Bearer "+c.token)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("GET %s: %w", path, err)
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	if resp.StatusCode >= 400 {
		return fmt.Errorf("GET %s: HTTP %d: %s", path, resp.StatusCode, string(body))
	}
	return json.Unmarshal(body, out)
}

func (c *Client) post(path string, in, out interface{}) error {
	body, _ := json.Marshal(in)
	req, _ := http.NewRequest(http.MethodPost, c.baseURL+path, bytes.NewReader(body))
	req.Header.Set("Authorization", "Bearer "+c.token)
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("POST %s: %w", path, err)
	}
	defer resp.Body.Close()

	respBody, _ := io.ReadAll(resp.Body)
	if resp.StatusCode >= 400 {
		return fmt.Errorf("POST %s: HTTP %d: %s", path, resp.StatusCode, string(respBody))
	}
	if out != nil {
		return json.Unmarshal(respBody, out)
	}
	return nil
}

func (c *Client) put(path string, in, out interface{}) error {
	body, _ := json.Marshal(in)
	req, _ := http.NewRequest(http.MethodPut, c.baseURL+path, bytes.NewReader(body))
	req.Header.Set("Authorization", "Bearer "+c.token)
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("PUT %s: %w", path, err)
	}
	defer resp.Body.Close()

	respBody, _ := io.ReadAll(resp.Body)
	if resp.StatusCode >= 400 {
		return fmt.Errorf("PUT %s: HTTP %d: %s", path, resp.StatusCode, string(respBody))
	}
	if out != nil {
		return json.Unmarshal(respBody, out)
	}
	return nil
}

func (c *Client) delete(path string) error {
	req, _ := http.NewRequest(http.MethodDelete, c.baseURL+path, nil)
	req.Header.Set("Authorization", "Bearer "+c.token)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("DELETE %s: %w", path, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("DELETE %s: HTTP %d: %s", path, resp.StatusCode, string(body))
	}
	return nil
}
