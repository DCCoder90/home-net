package secrets

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	infisical "github.com/infisical/go-sdk"
)

// Client wraps the Infisical SDK with project/environment context.
type Client struct {
	inner     infisical.InfisicalClientInterface
	projectID string
	env       string
}

// SSHConfig holds the credentials needed to open an SSH connection to a server.
type SSHConfig struct {
	Host    string
	User    string
	Port    int
	PrivKey string // raw PEM content of the private key
}

// New creates an authenticated Infisical client using Universal Auth.
func New(clientID, clientSecret, projectID, environment string) (*Client, error) {
	inner := infisical.NewInfisicalClient(context.Background(), infisical.Config{
		SiteUrl: "https://app.infisical.com",
	})
	if _, err := inner.Auth().UniversalAuthLogin(clientID, clientSecret); err != nil {
		return nil, fmt.Errorf("authenticating with Infisical: %w", err)
	}
	return &Client{inner: inner, projectID: projectID, env: environment}, nil
}

// FetchAll fetches all secrets from the given folder and returns them as map[key]value.
func (c *Client) FetchAll(folder string) (map[string]string, error) {
	list, err := c.inner.Secrets().List(infisical.ListSecretsOptions{
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

// FetchServerAccess fetches SSH credentials from the /server_access folder and
// writes each private key to a temp file. Returns a map of server name → SSHConfig.
// The caller is responsible for cleaning up the temp key files after provider construction.
func (c *Client) FetchServerAccess() (map[string]SSHConfig, error) {
	list, err := c.inner.Secrets().List(infisical.ListSecretsOptions{
		ProjectID:   c.projectID,
		Environment: c.env,
		SecretPath:  "/server_access",
	})
	if err != nil {
		return nil, fmt.Errorf("listing server_access secrets: %w", err)
	}

	raw := make(map[string]string, len(list))
	for _, s := range list {
		raw[s.SecretKey] = s.SecretValue
	}

	// Identify server names from keys like "server_<name>_ip"
	servers := map[string]bool{}
	for k := range raw {
		// keys are: server_tower_ip, server_tower_ssh_user, server_tower_private_key, etc.
		// We only care about the ones with known fields.
		if len(k) > 7 && k[:7] == "server_" {
			// extract name: server_<name>_ip → name
			rest := k[7:]
			for _, suffix := range []string{"_ip", "_password", "_private_key", "_public_key", "_ssh_port", "_ssh_user"} {
				if len(rest) > len(suffix) && rest[len(rest)-len(suffix):] == suffix {
					name := rest[:len(rest)-len(suffix)]
					servers[name] = true
				}
			}
		}
	}

	result := map[string]SSHConfig{}
	for name := range servers {
		ip := raw["server_"+name+"_ip"]
		user := raw["server_"+name+"_ssh_user"]
		privKey := raw["server_"+name+"_private_key"]
		port := 22
		if portStr := raw["server_"+name+"_ssh_port"]; portStr != "" {
			if n, err := strconv.Atoi(portStr); err == nil {
				port = n
			}
		}
		// Write private key to a temp file (600 perms) for use by the SSH Docker provider.
		keyPath, err := writeTempKey(name, privKey)
		if err != nil {
			return nil, fmt.Errorf("writing SSH key for server %q: %w", name, err)
		}
		result[name] = SSHConfig{
			Host:    ip,
			User:    user,
			Port:    port,
			PrivKey: keyPath,
		}
	}
	return result, nil
}

// EnsureGenerated ensures each name in names exists as a secret under the /generated
// folder in Infisical. Missing secrets are created with a cryptographically random
// 32-character hex value. Returns the final map[name]value for all names.
func (c *Client) EnsureGenerated(names []string) (map[string]string, error) {
	existing, err := c.FetchAll("/generated")
	if err != nil {
		// Folder may not exist yet; treat as empty and let Create establish it.
		existing = map[string]string{}
	}

	result := make(map[string]string, len(names))
	for _, name := range names {
		if val, ok := existing[name]; ok {
			result[name] = val
			continue
		}
		b := make([]byte, 16)
		if _, err := rand.Read(b); err != nil {
			return nil, fmt.Errorf("generating random value for secret %q: %w", name, err)
		}
		value := hex.EncodeToString(b)
		if _, err := c.inner.Secrets().Create(infisical.CreateSecretOptions{
			ProjectID:   c.projectID,
			Environment: c.env,
			SecretPath:  "/generated",
			SecretKey:   name,
			SecretValue: value,
		}); err != nil {
			return nil, fmt.Errorf("creating generated secret %q in Infisical: %w", name, err)
		}
		result[name] = value
	}
	return result, nil
}

// writeTempKey writes a PEM private key to a temp file with 0600 permissions.
// Strips Windows CRLF line endings before writing so OpenSSH can parse the key.
// Returns the file path.
func writeTempKey(serverName, privKeyPEM string) (string, error) {
	path := filepath.Join(os.TempDir(), "pulumi-ssh-key-"+serverName)
	cleaned := strings.ReplaceAll(privKeyPEM, "\r\n", "\n")
	cleaned = strings.ReplaceAll(cleaned, "\r", "\n")
	if !strings.HasSuffix(cleaned, "\n") {
		cleaned += "\n"
	}
	if err := os.WriteFile(path, []byte(cleaned), 0600); err != nil {
		return "", err
	}
	return path, nil
}
