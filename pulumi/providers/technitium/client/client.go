// Package client provides a lightweight HTTP client for the Technitium DNS REST API.
// Authentication uses a token query parameter on every request.
package client

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"
)

// Client holds the base URL and API token for a Technitium DNS server.
type Client struct {
	BaseURL    string
	Token      string
	httpClient *http.Client
}

// New creates a Client. baseURL should be like "http://192.168.4.53:5380".
func New(baseURL, token string) *Client {
	return &Client{
		BaseURL: baseURL,
		Token:   token,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// RecordInfo holds a single DNS record returned by the API.
// This type is populated manually in GetRecord — json tags are for reference only.
type RecordInfo struct {
	Domain    string
	Zone      string
	Type      string
	TTL       int
	IPAddress string
}

// apiResponse is the envelope returned by most Technitium endpoints.
type apiResponse struct {
	Status   string          `json:"status"`
	Response json.RawMessage `json:"response"`
	ErrorMessage string      `json:"errorMessage"`
}

// get performs a GET request to path with the given query params (token is added automatically).
func (c *Client) get(path string, params url.Values) ([]byte, error) {
	if params == nil {
		params = url.Values{}
	}
	params.Set("token", c.Token)
	u := fmt.Sprintf("%s%s?%s", c.BaseURL, path, params.Encode())
	resp, err := c.httpClient.Get(u)
	if err != nil {
		return nil, fmt.Errorf("GET %s: %w", path, err)
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("reading response: %w", err)
	}
	var ar apiResponse
	if err := json.Unmarshal(body, &ar); err != nil {
		return nil, fmt.Errorf("parsing response: %w", err)
	}
	if ar.Status != "ok" {
		return nil, fmt.Errorf("API error: %s", ar.ErrorMessage)
	}
	return ar.Response, nil
}

// AddRecord creates an A (or other type) record in the given zone.
func (c *Client) AddRecord(zone, domain, recType string, ttl int, ipAddress string) error {
	params := url.Values{
		"zone":      {zone},
		"domain":    {domain},
		"type":      {recType},
		"ttl":       {fmt.Sprintf("%d", ttl)},
		"ipAddress": {ipAddress},
		"overwrite": {"false"},
	}
	_, err := c.get("/api/zones/records/add", params)
	return err
}

// DeleteRecord removes the matching record from the zone.
func (c *Client) DeleteRecord(zone, domain, recType, ipAddress string) error {
	params := url.Values{
		"zone":      {zone},
		"domain":    {domain},
		"type":      {recType},
		"ipAddress": {ipAddress},
	}
	_, err := c.get("/api/zones/records/delete", params)
	return err
}

// GetRecord fetches the first matching record for domain+type in the zone.
// Returns nil if no matching record is found.
func (c *Client) GetRecord(zone, domain, recType string) (*RecordInfo, error) {
	params := url.Values{
		"zone":   {zone},
		"domain": {domain},
	}
	raw, err := c.get("/api/zones/records/get", params)
	if err != nil {
		return nil, err
	}

	var response struct {
		Records []struct {
			Name string `json:"name"`
			Type string `json:"type"`
			TTL  int    `json:"ttl"`
			RData struct {
				IPAddress string `json:"ipAddress"`
			} `json:"rData"`
		} `json:"records"`
	}
	if err := json.Unmarshal(raw, &response); err != nil {
		return nil, fmt.Errorf("parsing records: %w", err)
	}
	for _, r := range response.Records {
		if r.Type == recType {
			return &RecordInfo{
				Domain:    r.Name,
				Zone:      zone,
				Type:      r.Type,
				TTL:       r.TTL,
				IPAddress: r.RData.IPAddress,
			}, nil
		}
	}
	return nil, nil
}
