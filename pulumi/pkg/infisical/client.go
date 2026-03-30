package infisical

import (
	"context"
	"fmt"

	infisicalsdk "github.com/infisical/go-sdk"
)

// Client wraps the Infisical SDK with project/environment context baked in.
type Client struct {
	inner     infisicalsdk.InfisicalClientInterface
	projectID string
	env       string
}

// New creates an authenticated Infisical client.
func New(clientID, clientSecret, projectID, environment string) (*Client, error) {
	inner := infisicalsdk.NewInfisicalClient(context.Background(), infisicalsdk.Config{
		SiteUrl: "https://app.infisical.com",
	})
	if _, err := inner.Auth().UniversalAuthLogin(clientID, clientSecret); err != nil {
		return nil, fmt.Errorf("authenticating with Infisical: %w", err)
	}
	return &Client{inner: inner, projectID: projectID, env: environment}, nil
}

// FetchAll fetches all secrets from a folder and returns them as a map[key]value.
func (c *Client) FetchAll(folder string) (map[string]string, error) {
	list, err := c.inner.Secrets().List(infisicalsdk.ListSecretsOptions{
		ProjectID:   c.projectID,
		Environment: c.env,
		SecretPath:  folder,
	})
	if err != nil {
		return nil, fmt.Errorf("listing secrets at %q: %w", folder, err)
	}

	result := make(map[string]string, len(list))
	for _, s := range list {
		result[s.SecretKey] = s.SecretValue
	}
	return result, nil
}

// Get fetches a single secret by key from a folder.
func (c *Client) Get(key, folder string) (string, error) {
	s, err := c.inner.Secrets().Retrieve(infisicalsdk.RetrieveSecretOptions{
		SecretKey:   key,
		ProjectID:   c.projectID,
		Environment: c.env,
		SecretPath:  folder,
	})
	if err != nil {
		return "", fmt.Errorf("fetching secret %q from %q: %w", key, folder, err)
	}
	return s.SecretValue, nil
}
