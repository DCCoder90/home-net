package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

// Client talks to the Nginx Proxy Manager REST API.
type Client struct {
	baseURL    string
	username   string
	password   string
	token      string
	httpClient *http.Client
}

// New creates an NPM API client. Call Authenticate() before use.
func New(baseURL, username, password string) *Client {
	return &Client{
		baseURL:    baseURL,
		username:   username,
		password:   password,
		httpClient: &http.Client{Timeout: 60 * time.Second},
	}
}

// Authenticate fetches a JWT from NPM and stores it for subsequent calls.
func (c *Client) Authenticate() error {
	payload := map[string]string{
		"identity": c.username,
		"secret":   c.password,
	}
	body, _ := json.Marshal(payload)

	resp, err := c.httpClient.Post(
		fmt.Sprintf("%s/api/tokens", c.baseURL),
		"application/json",
		bytes.NewReader(body),
	)
	if err != nil {
		return fmt.Errorf("NPM auth: %w", err)
	}
	defer resp.Body.Close()

	var result struct {
		Token string `json:"token"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return fmt.Errorf("NPM auth decode: %w", err)
	}
	if result.Token == "" {
		return fmt.Errorf("NPM auth: empty token returned")
	}
	c.token = result.Token
	return nil
}

// AccessList represents a single NPM access list.
type AccessList struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

// ListAccessLists returns all configured NPM access lists.
func (c *Client) ListAccessLists() ([]AccessList, error) {
	var lists []AccessList
	if err := c.get("/api/nginx/access-lists", &lists); err != nil {
		return nil, err
	}
	return lists, nil
}

// AccessListItem represents a single allow/deny rule inside an access list.
type AccessListItem struct {
	Address   string `json:"address"`
	Directive string `json:"directive"`
}

// CreateAccessListInput is the request body for POST /api/nginx/access-lists.
type CreateAccessListInput struct {
	Name        string           `json:"name"`
	Items       []AccessListItem `json:"items"`
	PassAuth    bool             `json:"pass_auth"`
	SatisfyAny  bool             `json:"satisfy_any"`
}

// CreateAccessList creates a new NPM access list and returns its ID.
func (c *Client) CreateAccessList(input CreateAccessListInput) (int, error) {
	var result struct {
		ID int `json:"id"`
	}
	if err := c.post("/api/nginx/access-lists", input, &result); err != nil {
		return 0, err
	}
	return result.ID, nil
}

// DeleteAccessList removes an access list by ID.
func (c *Client) DeleteAccessList(id int) error {
	return c.delete(fmt.Sprintf("/api/nginx/access-lists/%d", id))
}

// CertificateInput is the request body for creating a Let's Encrypt cert.
type CertificateInput struct {
	Provider    string                 `json:"provider"`
	DomainNames []string               `json:"domain_names"`
	Meta        map[string]interface{} `json:"meta"`
}

// CreateCertificate requests a new Let's Encrypt certificate and returns its ID.
func (c *Client) CreateCertificate(input CertificateInput) (int, error) {
	var result struct {
		ID int `json:"id"`
	}
	if err := c.post("/api/nginx/certificates", input, &result); err != nil {
		return 0, err
	}
	return result.ID, nil
}

// DeleteCertificate removes a certificate by ID.
func (c *Client) DeleteCertificate(id int) error {
	return c.delete(fmt.Sprintf("/api/nginx/certificates/%d", id))
}

// ProxyHostInput is the request body for creating an NPM proxy host.
type ProxyHostInput struct {
	DomainNames           []string `json:"domain_names"`
	ForwardScheme         string   `json:"forward_scheme"`
	ForwardHost           string   `json:"forward_host"`
	ForwardPort           int      `json:"forward_port"`
	CachingEnabled        bool     `json:"caching_enabled"`
	AllowWebsocketUpgrade bool     `json:"allow_websocket_upgrade"`
	BlockExploits         bool     `json:"block_exploits"`
	AccessListID          int      `json:"access_list_id"`
	CertificateID         int      `json:"certificate_id"`
	SSLForced             bool     `json:"ssl_forced"`
	HSTSEnabled           bool     `json:"hsts_enabled"`
	HSTSSubdomains        bool     `json:"hsts_subdomains"`
	HTTP2Support          bool     `json:"http2_support"`
	Enabled               bool     `json:"enabled"`
}

// CreateProxyHost creates a new NPM proxy host and returns its ID.
func (c *Client) CreateProxyHost(input ProxyHostInput) (int, error) {
	var result struct {
		ID int `json:"id"`
	}
	if err := c.post("/api/nginx/proxy-hosts", input, &result); err != nil {
		return 0, err
	}
	return result.ID, nil
}

// UpdateProxyHost replaces a proxy host by ID.
func (c *Client) UpdateProxyHost(id int, input ProxyHostInput) error {
	return c.put(fmt.Sprintf("/api/nginx/proxy-hosts/%d", id), input, nil)
}

// DeleteProxyHost removes a proxy host by ID.
func (c *Client) DeleteProxyHost(id int) error {
	return c.delete(fmt.Sprintf("/api/nginx/proxy-hosts/%d", id))
}

// --- HTTP helpers ---

func (c *Client) get(path string, out interface{}) error {
	req, _ := http.NewRequest(http.MethodGet, c.baseURL+path, nil)
	req.Header.Set("Authorization", "Bearer "+c.token)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("GET %s: %w", path, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("GET %s: HTTP %d: %s", path, resp.StatusCode, string(body))
	}
	return json.NewDecoder(resp.Body).Decode(out)
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
