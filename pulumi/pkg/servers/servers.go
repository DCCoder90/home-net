package servers

import (
	"fmt"
	"os"
	"sort"

	"github.com/DCCoder90/home-net/pulumi/pkg/infisical"
	"github.com/DCCoder90/home-net/pulumi/pkg/types"
	dockerprovider "github.com/pulumi/pulumi-docker/sdk/v4/go/docker"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi/config"
	"gopkg.in/yaml.v3"
)

const serversFolder = "/servers"

// ServerContext holds a server's config and its Pulumi Docker provider.
type ServerContext struct {
	Name     string
	Config   types.ServerConfig
	Provider *dockerprovider.Provider
}

// Load reads config/servers.yaml, fetches SSH credentials from Infisical for each
// server, and creates one Docker provider per host.
func Load(
	ctx *pulumi.Context,
	infisicalClient *infisical.Client,
) (map[string]*ServerContext, error) {
	raw, err := os.ReadFile("../config/servers.yaml")
	if err != nil {
		return nil, fmt.Errorf("reading servers.yaml: %w", err)
	}

	var file types.ServersFile
	if err := yaml.Unmarshal(raw, &file); err != nil {
		return nil, fmt.Errorf("parsing servers.yaml: %w", err)
	}

	result := make(map[string]*ServerContext, len(file.Servers))

	for name, cfg := range file.Servers {
		creds, err := fetchSSHCredentials(infisicalClient, name)
		if err != nil {
			return nil, err
		}

		keyPath := fmt.Sprintf("/tmp/pulumi-ssh-%s", name)
		if err := os.WriteFile(keyPath, []byte(creds.privateKey), 0600); err != nil {
			return nil, fmt.Errorf("writing SSH key for %q: %w", name, err)
		}

		provider, err := dockerprovider.NewProvider(ctx, "docker-"+name, &dockerprovider.ProviderArgs{
			Host: pulumi.String(fmt.Sprintf("ssh://%s@%s:%s",
				creds.user, creds.ip, creds.port)),
			SshOpts: pulumi.StringArray{
				pulumi.String("-i"), pulumi.String(keyPath),
				pulumi.String("-o"), pulumi.String("StrictHostKeyChecking=no"),
				pulumi.String("-o"), pulumi.String("UserKnownHostsFile=/dev/null"),
			},
			RegistryAuth: dockerprovider.ProviderRegistryAuthArray{
				&dockerprovider.ProviderRegistryAuthArgs{
					Address:  pulumi.String("ghcr.io"),
					Username: pulumi.String(config.New(ctx, "").Get("ghcrUsername")),
					Password: pulumi.String(config.New(ctx, "").Get("ghcrToken")),
				},
			},
		})
		if err != nil {
			return nil, fmt.Errorf("creating Docker provider for %q: %w", name, err)
		}

		result[name] = &ServerContext{Name: name, Config: cfg, Provider: provider}
	}

	return result, nil
}

// Require returns the ServerContext for hostName, or a clear error if the host
// is not defined in servers.yaml — satisfying the requirement that an undefined
// host causes an explicit failure.
func Require(servers map[string]*ServerContext, hostName string) (*ServerContext, error) {
	if srv, ok := servers[hostName]; ok {
		return srv, nil
	}

	known := make([]string, 0, len(servers))
	for k := range servers {
		known = append(known, k)
	}
	sort.Strings(known)

	return nil, fmt.Errorf(
		"host %q is not defined in servers.yaml — add it before deploying services to it (known hosts: %v)",
		hostName, known,
	)
}

// ServiceHost returns the host for a service, defaulting to "tower".
func ServiceHost(host string) string {
	if host == "" {
		return "tower"
	}
	return host
}

// PrimaryIP walks a service's network list and returns the first IP on the
// server's primary network, then secondary, then any IP.
func PrimaryIP(networks []types.NetworkEntry, srv *ServerContext) string {
	primary := srv.Config.Networks.Primary
	secondary := srv.Config.Networks.Secondary

	for _, n := range networks {
		if n.Name == primary && n.IPAddress != "" {
			return n.IPAddress
		}
	}
	if secondary != "" {
		for _, n := range networks {
			if n.Name == secondary && n.IPAddress != "" {
				return n.IPAddress
			}
		}
	}
	for _, n := range networks {
		if n.IPAddress != "" {
			return n.IPAddress
		}
	}
	return ""
}

type sshCredentials struct {
	ip, user, port, privateKey string
}

func fetchSSHCredentials(client *infisical.Client, serverName string) (*sshCredentials, error) {
	fetch := func(suffix string) (string, error) {
		return client.Get(fmt.Sprintf("server_%s_%s", serverName, suffix), serversFolder)
	}

	ip, err := fetch("ip")
	if err != nil {
		return nil, err
	}
	user, err := fetch("ssh_user")
	if err != nil {
		return nil, err
	}
	port, err := fetch("ssh_port")
	if err != nil {
		return nil, err
	}
	key, err := fetch("ssh_private_key")
	if err != nil {
		return nil, err
	}

	return &sshCredentials{ip: ip, user: user, port: port, privateKey: key}, nil
}
