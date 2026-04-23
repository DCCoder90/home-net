// Package client provides a lightweight HTTP client for the Nginx Proxy Manager REST API.
// Authentication is done by exchanging username/password for a Bearer JWT token.
package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

// Client holds connection details and the cached auth token for NPM.
type Client struct {
	BaseURL    string
	Username   string
	Password   string
	token      string
	httpClient *http.Client
}

// New returns a new Client. Call Authenticate before making API calls.
func New(baseURL, username, password string) *Client {
	return &Client{
		BaseURL:  baseURL,
		Username: username,
		Password: password,
		httpClient: &http.Client{
			Timeout: 3 * time.Minute,
		},
	}
}

// Authenticate exchanges credentials for a Bearer token stored on the client.
func (c *Client) Authenticate() error {
	payload := map[string]string{
		"identity": c.Username,
		"secret":   c.Password,
	}
	body, err := json.Marshal(payload)
	if err != nil {
		return err
	}
	resp, err := c.httpClient.Post(c.BaseURL+"/api/tokens", "application/json", bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("authenticate: %w", err)
	}
	defer resp.Body.Close()
	var result struct {
		Token string `json:"token"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return fmt.Errorf("decode auth response: %w", err)
	}
	if result.Token == "" {
		return fmt.Errorf("authenticate: empty token returned (check credentials)")
	}
	c.token = result.Token
	return nil
}

// request performs an authenticated HTTP request and decodes the JSON response into out.
func (c *Client) request(method, path string, payload interface{}, out interface{}) error {
	var bodyReader io.Reader
	if payload != nil {
		b, err := json.Marshal(payload)
		if err != nil {
			return err
		}
		bodyReader = bytes.NewReader(b)
	}
	req, err := http.NewRequest(method, c.BaseURL+path, bodyReader)
	if err != nil {
		return err
	}
	req.Header.Set("Authorization", "Bearer "+c.token)
	if payload != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("%s %s: %w", method, path, err)
	}
	defer resp.Body.Close()
	if resp.StatusCode >= 400 {
		b, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("%s %s: HTTP %d: %s", method, path, resp.StatusCode, string(b))
	}
	if out != nil {
		return json.NewDecoder(resp.Body).Decode(out)
	}
	return nil
}

// ── Access Lists ─────────────────────────────────────────────────────────────

// AccessListClientEntry is a single IP rule entry in an NPM access list.
// These map to the "clients" field in the NPM API (IP allow/deny rules).
// The separate "items" field is for HTTP Basic Auth credentials and is not used here.
type AccessListClientEntry struct {
	Address   string `json:"address"`
	Directive string `json:"directive"` // "allow" or "deny"
}

// flexBool unmarshals both JSON booleans (true/false) and integers (0/1).
// NPM's POST response returns booleans; GET responses return integers.
type flexBool bool

func (f *flexBool) UnmarshalJSON(data []byte) error {
	switch string(data) {
	case "true", "1":
		*f = true
	case "false", "0":
		*f = false
	default:
		return fmt.Errorf("cannot unmarshal %s into bool", string(data))
	}
	return nil
}

// AccessList is the NPM API representation of an access list.
// satisfy_any and pass_auth are returned as booleans in POST responses and
// integers in GET responses; flexBool handles both.
// IP rules are under "clients"; "items" is for HTTP Basic Auth users (unused here).
type AccessList struct {
	ID         int                     `json:"id"`
	Name       string                  `json:"name"`
	SatisfyAny flexBool                `json:"satisfy_any"`
	PassAuth   flexBool                `json:"pass_auth"`
	Clients    []AccessListClientEntry `json:"clients"`
}

func (c *Client) CreateAccessList(name string, passAuth, satisfyAny bool, clients []AccessListClientEntry) (*AccessList, error) {
	payload := map[string]interface{}{
		"name":        name,
		"satisfy_any": satisfyAny,
		"pass_auth":   passAuth,
		"items":       []struct{}{}, // NPM unconditionally calls .map() on items — must not be absent
		"clients":     clients,
	}
	var result AccessList
	if err := c.request("POST", "/api/nginx/access-lists", payload, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

func (c *Client) GetAccessList(id int) (*AccessList, error) {
	// NPM may not support GET /access-lists/:id — list all and filter.
	var results []AccessList
	if err := c.request("GET", "/api/nginx/access-lists", nil, &results); err != nil {
		return nil, err
	}
	for _, al := range results {
		if al.ID == id {
			return &al, nil
		}
	}
	return nil, fmt.Errorf("access list with ID %d not found", id)
}

func (c *Client) UpdateAccessList(id int, name string, passAuth, satisfyAny bool, clients []AccessListClientEntry) (*AccessList, error) {
	payload := map[string]interface{}{
		"name":        name,
		"satisfy_any": satisfyAny,
		"pass_auth":   passAuth,
		"clients":     clients,
	}
	var result AccessList
	if err := c.request("PUT", fmt.Sprintf("/api/nginx/access-lists/%d", id), payload, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

func (c *Client) DeleteAccessList(id int) error {
	err := c.request("DELETE", fmt.Sprintf("/api/nginx/access-lists/%d", id), nil, nil)
	if isNotFound(err) {
		return nil
	}
	return err
}

// ── Certificates ─────────────────────────────────────────────────────────────

// Certificate is the NPM API representation of a Let's Encrypt certificate.
type Certificate struct {
	ID          int      `json:"id"`
	DomainNames []string `json:"domain_names"`
	Provider    string   `json:"provider"`
}

func (c *Client) CreateCertificate(domainNames []string, email, cfAPIToken string) (*Certificate, error) {
	payload := map[string]interface{}{
		"domain_names": domainNames,
		"meta": map[string]interface{}{
			"dns_challenge":            true,
			"dns_provider":             "cloudflare",
			"dns_provider_credentials": "dns_cloudflare_api_token = " + cfAPIToken,
		},
		"provider": "letsencrypt",
	}
	const maxAttempts = 10
	for attempt := 1; attempt <= maxAttempts; attempt++ {
		var result Certificate
		err := c.request("POST", "/api/nginx/certificates", payload, &result)
		if err == nil {
			return &result, nil
		}
		if strings.Contains(err.Error(), "Another instance of Certbot is already running") {
			time.Sleep(40 * time.Second)
			continue
		}
		return nil, err
	}
	return nil, fmt.Errorf("certificate creation failed after %d attempts: certbot lock never released", maxAttempts)
}

func (c *Client) GetCertificate(id int) (*Certificate, error) {
	// NPM does not have a GET /certificates/:id endpoint — list all and filter.
	var results []Certificate
	if err := c.request("GET", "/api/nginx/certificates", nil, &results); err != nil {
		return nil, err
	}
	for _, cert := range results {
		if cert.ID == id {
			return &cert, nil
		}
	}
	return nil, fmt.Errorf("certificate with ID %d not found", id)
}

func (c *Client) DeleteCertificate(id int) error {
	err := c.request("DELETE", fmt.Sprintf("/api/nginx/certificates/%d", id), nil, nil)
	if isNotFound(err) {
		return nil
	}
	return err
}

// ── Proxy Hosts ───────────────────────────────────────────────────────────────

// ProxyHost is the NPM API representation of a proxy host.
type ProxyHost struct {
	ID                    int      `json:"id,omitempty"`
	DomainNames           []string `json:"domain_names"`
	ForwardHost           string   `json:"forward_host"`
	ForwardPort           int      `json:"forward_port"`
	ForwardScheme         string   `json:"forward_scheme"`
	CertificateID         int      `json:"certificate_id,omitempty"`
	AccessListID          int      `json:"access_list_id"`
	SSLForced             bool     `json:"ssl_forced"`
	HTTP2Support          bool     `json:"http2_support"`
	BlockExploits         bool     `json:"block_exploits"`
	AllowWebsocketUpgrade bool     `json:"allow_websocket_upgrade"`
	AdvancedConfig        string   `json:"advanced_config,omitempty"`
}

func (c *Client) CreateProxyHost(ph ProxyHost) (*ProxyHost, error) {
	var result ProxyHost
	if err := c.request("POST", "/api/nginx/proxy-hosts", ph, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

func (c *Client) GetProxyHost(id int) (*ProxyHost, error) {
	var result ProxyHost
	if err := c.request("GET", fmt.Sprintf("/api/nginx/proxy-hosts/%d", id), nil, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

func (c *Client) UpdateProxyHost(id int, ph ProxyHost) (*ProxyHost, error) {
	var result ProxyHost
	if err := c.request("PUT", fmt.Sprintf("/api/nginx/proxy-hosts/%d", id), ph, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

func (c *Client) DeleteProxyHost(id int) error {
	err := c.request("DELETE", fmt.Sprintf("/api/nginx/proxy-hosts/%d", id), nil, nil)
	if isNotFound(err) {
		return nil
	}
	return err
}

// isNotFound returns true when err is an HTTP 404 response from the NPM API,
// meaning the resource was already deleted and the operation can be treated as a no-op.
func isNotFound(err error) bool {
	return err != nil && strings.Contains(err.Error(), "HTTP 404")
}
