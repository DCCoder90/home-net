package client

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"
)

// Client talks to the Technitium DNS Server HTTP API.
type Client struct {
	baseURL    string
	token      string
	httpClient *http.Client
}

// New creates an authenticated Technitium API client.
// It logs in and stores the session token.
func New(host, apiToken string) *Client {
	return &Client{
		baseURL:    host,
		token:      apiToken,
		httpClient: &http.Client{Timeout: 30 * time.Second},
	}
}

// AddRecord creates an A or AAAA DNS record.
func (c *Client) AddRecord(domain, zone, recordType, ipAddress string, ttl int) error {
	params := url.Values{}
	params.Set("token", c.token)
	params.Set("domain", domain)
	params.Set("zone", zone)
	params.Set("type", recordType)
	params.Set("ipAddress", ipAddress)
	params.Set("ttl", fmt.Sprintf("%d", ttl))

	resp, err := c.httpClient.Get(fmt.Sprintf("%s/api/zones/records/add?%s", c.baseURL, params.Encode()))
	if err != nil {
		return fmt.Errorf("adding DNS record: %w", err)
	}
	defer resp.Body.Close()

	return c.checkResponse(resp, "add record")
}

// DeleteRecord removes a DNS record by domain, type, and IP address.
func (c *Client) DeleteRecord(domain, zone, recordType, ipAddress string) error {
	params := url.Values{}
	params.Set("token", c.token)
	params.Set("domain", domain)
	params.Set("zone", zone)
	params.Set("type", recordType)
	params.Set("ipAddress", ipAddress)

	resp, err := c.httpClient.Get(fmt.Sprintf("%s/api/zones/records/delete?%s", c.baseURL, params.Encode()))
	if err != nil {
		return fmt.Errorf("deleting DNS record: %w", err)
	}
	defer resp.Body.Close()

	return c.checkResponse(resp, "delete record")
}

// RecordExists checks whether a record exists by looking it up.
func (c *Client) RecordExists(domain, zone, recordType string) (bool, string, error) {
	params := url.Values{}
	params.Set("token", c.token)
	params.Set("domain", domain)
	params.Set("zone", zone)
	params.Set("type", recordType)

	resp, err := c.httpClient.Get(fmt.Sprintf("%s/api/zones/records?%s", c.baseURL, params.Encode()))
	if err != nil {
		return false, "", fmt.Errorf("querying DNS record: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return false, "", err
	}

	var result struct {
		Status   string `json:"status"`
		Response struct {
			Records []struct {
				Name  string `json:"name"`
				Type  string `json:"type"`
				RData struct {
					IPAddress string `json:"ipAddress"`
				} `json:"rData"`
			} `json:"records"`
		} `json:"response"`
	}
	if err := json.Unmarshal(body, &result); err != nil {
		return false, "", fmt.Errorf("parsing record response: %w", err)
	}

	if result.Status != "ok" {
		return false, "", nil
	}

	for _, r := range result.Response.Records {
		return true, r.RData.IPAddress, nil
	}
	return false, "", nil
}

func (c *Client) checkResponse(resp *http.Response, action string) error {
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("%s: reading response: %w", action, err)
	}

	var result struct {
		Status  string `json:"status"`
		ErrorMessage string `json:"errorMessage"`
	}
	if err := json.Unmarshal(body, &result); err != nil {
		return fmt.Errorf("%s: parsing response: %w", action, err)
	}

	if result.Status != "ok" {
		return fmt.Errorf("%s failed: %s", action, result.ErrorMessage)
	}
	return nil
}
