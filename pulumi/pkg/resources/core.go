package resources

import (
	"fmt"

	"github.com/DCCoder90/home-net/pulumi/pkg/config"
	dockerprovider "github.com/pulumi/pulumi-docker/sdk/v4/go/docker"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
	pulumiconfig "github.com/pulumi/pulumi/sdk/v3/go/pulumi/config"
)

// CoreContainers holds references to the four foundational containers.
type CoreContainers struct {
	TechnitiumDNS   *dockerprovider.Container
	NginxProxy      *dockerprovider.Container
	AuthentikServer *dockerprovider.Container
	AuthentikWorker *dockerprovider.Container
}

// RegisterCoreContainers creates the Technitium DNS, Nginx Proxy Manager, and
// Authentik server/worker containers on the tower Docker provider.
func RegisterCoreContainers(
	ctx *pulumi.Context,
	system *config.SystemConfig,
	secrets map[string]string,
	towerProvider *dockerprovider.Provider,
	importIDs map[string]string,
) (*CoreContainers, error) {
	cfg := pulumiconfig.New(ctx, "")
	tz := cfg.Get("tz")
	if tz == "" {
		tz = "America/Chicago"
	}
	hostPathDNS := cfg.Get("hostPathDnsConfig")
	hostPathNPMData := cfg.Get("hostPathNpmData")
	hostPathNPMLetsencrypt := cfg.Get("hostPathNpmLetsencrypt")
	hostPathNPMLogs := cfg.Get("hostPathNpmLogs")

	provOpt := pulumi.Provider(towerProvider)

	// ── Technitium DNS ──────────────────────────────────────────────────────────
	dnsImage, err := dockerprovider.NewRemoteImage(ctx, "technitium-dns-image", &dockerprovider.RemoteImageArgs{
		Name:        pulumi.String("technitium/dns-server:14.3.0"),
		KeepLocally: pulumi.Bool(true),
	}, provOpt)
	if err != nil {
		return nil, fmt.Errorf("technitium image: %w", err)
	}

	technitiumDNS, err := dockerprovider.NewContainer(ctx, "technitium_dns",
		&dockerprovider.ContainerArgs{
			Name:    pulumi.String("technitium_dns"),
			Image:   dnsImage.RepoDigest,
			Restart: pulumi.String("unless-stopped"),
			Envs: pulumi.StringArray{
				pulumi.String("TZ=" + tz),
				pulumi.String("DNS_SERVER_ADMIN_PASSWORD=" + secrets["technitium_admin_password"]),
			},
			Volumes: dockerprovider.ContainerVolumeArray{
				&dockerprovider.ContainerVolumeArgs{
					HostPath:      pulumi.String(hostPathDNS),
					ContainerPath: pulumi.String("/etc/dns"),
				},
			},
			NetworksAdvanced: dockerprovider.ContainerNetworksAdvancedArray{
				&dockerprovider.ContainerNetworksAdvancedArgs{
					Name:        pulumi.String("br1"),
					Ipv4Address: pulumi.String(system.DNSIP),
				},
			},
			Labels: dockerprovider.ContainerLabelArray{
				&dockerprovider.ContainerLabelArgs{
					Label: pulumi.String("net.unraid.docker.icon"),
					Value: pulumi.String("https://cdn.jsdelivr.net/gh/homarr-labs/dashboard-icons/png/technitium.png"),
				},
			},
		},
		append([]pulumi.ResourceOption{provOpt,
			pulumi.IgnoreChanges([]string{"labels", "logOpts"}),
		}, importOpts("technitium_dns", importIDs)...)...,
	)
	if err != nil {
		return nil, fmt.Errorf("technitium_dns container: %w", err)
	}

	// ── Nginx Proxy Manager ─────────────────────────────────────────────────────
	npmImage, err := dockerprovider.NewRemoteImage(ctx, "npm-image", &dockerprovider.RemoteImageArgs{
		Name:        pulumi.String("jc21/nginx-proxy-manager:2"),
		KeepLocally: pulumi.Bool(true),
	}, provOpt)
	if err != nil {
		return nil, fmt.Errorf("npm image: %w", err)
	}

	nginxProxy, err := dockerprovider.NewContainer(ctx, "nginx_proxy_manager",
		&dockerprovider.ContainerArgs{
			Name:    pulumi.String("Nginx-Proxy-Manager-Official"),
			Image:   npmImage.RepoDigest,
			Restart: pulumi.String("unless-stopped"),
			Envs: pulumi.StringArray{
				pulumi.String("TZ=" + tz),
			},
			Ports: dockerprovider.ContainerPortArray{
				&dockerprovider.ContainerPortArgs{Internal: pulumi.Int(80), External: pulumi.Int(80)},
				&dockerprovider.ContainerPortArgs{Internal: pulumi.Int(81), External: pulumi.Int(81)},
				&dockerprovider.ContainerPortArgs{Internal: pulumi.Int(443), External: pulumi.Int(443)},
			},
			Volumes: dockerprovider.ContainerVolumeArray{
				&dockerprovider.ContainerVolumeArgs{
					HostPath:      pulumi.String(hostPathNPMData),
					ContainerPath: pulumi.String("/data"),
				},
				&dockerprovider.ContainerVolumeArgs{
					HostPath:      pulumi.String(hostPathNPMLetsencrypt),
					ContainerPath: pulumi.String("/etc/letsencrypt"),
				},
				&dockerprovider.ContainerVolumeArgs{
					HostPath:      pulumi.String(hostPathNPMLogs),
					ContainerPath: pulumi.String("/var/log"),
				},
			},
			NetworksAdvanced: dockerprovider.ContainerNetworksAdvancedArray{
				&dockerprovider.ContainerNetworksAdvancedArgs{
					Name:        pulumi.String("br1"),
					Ipv4Address: pulumi.String(system.ProxyIP),
				},
			},
			Labels: dockerprovider.ContainerLabelArray{
				&dockerprovider.ContainerLabelArgs{
					Label: pulumi.String("net.unraid.docker.icon"),
					Value: pulumi.String("https://cdn.jsdelivr.net/gh/homarr-labs/dashboard-icons/png/nginx-proxy-manager.png"),
				},
			},
		},
		append([]pulumi.ResourceOption{provOpt,
			pulumi.IgnoreChanges([]string{"labels", "logOpts"}),
		}, importOpts("nginx_proxy_manager", importIDs)...)...,
	)
	if err != nil {
		return nil, fmt.Errorf("nginx_proxy_manager container: %w", err)
	}

	// ── Authentik ───────────────────────────────────────────────────────────────
	authentikImage, err := dockerprovider.NewRemoteImage(ctx, "authentik-image", &dockerprovider.RemoteImageArgs{
		Name:        pulumi.String(system.Authentik.Image),
		KeepLocally: pulumi.Bool(true),
	}, provOpt)
	if err != nil {
		return nil, fmt.Errorf("authentik image: %w", err)
	}

	authentikEnvs := pulumi.StringArray{
		pulumi.String("TZ=" + tz),
		pulumi.String("AUTHENTIK_POSTGRESQL__HOST=" + system.Authentik.PostgreSQLHost),
		pulumi.String("AUTHENTIK_POSTGRESQL__USER=" + system.Authentik.PostgreSQLUser),
		pulumi.String("AUTHENTIK_POSTGRESQL__NAME=" + system.Authentik.PostgreSQLName),
		pulumi.String("AUTHENTIK_POSTGRESQL__PASSWORD=" + secrets["authentik_postgresql_password"]),
		pulumi.String("AUTHENTIK_SECRET_KEY=" + secrets["authentik_secret_key"]),
		pulumi.String("AUTHENTIK_REDIS__HOST=" + system.Authentik.RedisHost),
	}

	authentikVolumes := dockerprovider.ContainerVolumeArray{
		&dockerprovider.ContainerVolumeArgs{
			HostPath:      pulumi.String(system.Authentik.DataPath + "/media"),
			ContainerPath: pulumi.String("/media"),
		},
		&dockerprovider.ContainerVolumeArgs{
			HostPath:      pulumi.String(system.Authentik.DataPath + "/custom-templates"),
			ContainerPath: pulumi.String("/templates"),
		},
	}

	authentikServer, err := dockerprovider.NewContainer(ctx, "authentik-server",
		&dockerprovider.ContainerArgs{
			Name:    pulumi.String("authentik-server"),
			Image:   authentikImage.RepoDigest,
			Restart: pulumi.String("unless-stopped"),
			Command: pulumi.StringArray{pulumi.String("server")},
			Envs:    authentikEnvs,
			Volumes: authentikVolumes,
			NetworksAdvanced: dockerprovider.ContainerNetworksAdvancedArray{
				&dockerprovider.ContainerNetworksAdvancedArgs{
					Name:        pulumi.String("br1"),
					Ipv4Address: pulumi.String(system.Authentik.IPAddress),
				},
			},
			Labels: dockerprovider.ContainerLabelArray{
				&dockerprovider.ContainerLabelArgs{
					Label: pulumi.String("net.unraid.docker.icon"),
					Value: pulumi.String("https://cdn.jsdelivr.net/gh/homarr-labs/dashboard-icons/png/authentik.png"),
				},
				&dockerprovider.ContainerLabelArgs{
					Label: pulumi.String("net.unraid.docker.webui"),
					Value: pulumi.String(fmt.Sprintf("http://%%IP%%:%d", system.Authentik.Port)),
				},
			},
		},
		append([]pulumi.ResourceOption{provOpt,
			pulumi.IgnoreChanges([]string{"labels", "logOpts"}),
		}, importOpts("authentik-server", importIDs)...)...,
	)
	if err != nil {
		return nil, fmt.Errorf("authentik-server container: %w", err)
	}

	workerVolumes := append(authentikVolumes, &dockerprovider.ContainerVolumeArgs{
		HostPath:      pulumi.String(system.Authentik.DataPath + "/certs"),
		ContainerPath: pulumi.String("/certs"),
	})

	authentikWorker, err := dockerprovider.NewContainer(ctx, "authentik-worker",
		&dockerprovider.ContainerArgs{
			Name:    pulumi.String("authentik-worker"),
			Image:   authentikImage.RepoDigest,
			Restart: pulumi.String("unless-stopped"),
			Command: pulumi.StringArray{pulumi.String("worker")},
			Envs:    authentikEnvs,
			Volumes: workerVolumes,
			// Worker uses the same bridge but without a static IP — it doesn't
			// serve traffic directly; it communicates internally with the server.
			NetworksAdvanced: dockerprovider.ContainerNetworksAdvancedArray{
				&dockerprovider.ContainerNetworksAdvancedArgs{
					Name: pulumi.String("br1"),
				},
			},
		},
		append([]pulumi.ResourceOption{provOpt,
			pulumi.IgnoreChanges([]string{"logOpts"}),
		}, importOpts("authentik-worker", importIDs)...)...,
	)
	if err != nil {
		return nil, fmt.Errorf("authentik-worker container: %w", err)
	}

	return &CoreContainers{
		TechnitiumDNS:   technitiumDNS,
		NginxProxy:      nginxProxy,
		AuthentikServer: authentikServer,
		AuthentikWorker: authentikWorker,
	}, nil
}

// RegisterOutpostContainer deploys the standalone Authentik proxy outpost container.
// It is a no-op unless both authentik_outpost_token is present in secrets and
// outpost_ip_address is set in system.yaml (phase 3 of the Authentik bootstrap).
//
// After phase 2 (RegisterAuthResources), retrieve the outpost token from the
// Authentik UI under Outposts → "Pulumi Proxy Outpost" → View Token, then store
// it in Infisical as authentik_outpost_token and run pulumi up again.
func RegisterOutpostContainer(
	ctx *pulumi.Context,
	system *config.SystemConfig,
	secrets map[string]string,
	towerProvider *dockerprovider.Provider,
	importIDs map[string]string,
) (*dockerprovider.Container, error) {
	token := secrets["authentik_outpost_token"]
	if token == "" || system.Authentik.OutpostIPAddress == "" {
		return nil, nil
	}

	image := system.Authentik.OutpostImage
	if image == "" {
		image = "ghcr.io/goauthentik/proxy:latest"
	}

	authentikHost := "https://" + system.Authentik.DomainName
	if system.Authentik.DomainName == "" {
		authentikHost = fmt.Sprintf("http://%s:%d", system.Authentik.IPAddress, system.Authentik.Port)
	}

	provOpt := pulumi.Provider(towerProvider)

	outpostImage, err := dockerprovider.NewRemoteImage(ctx, "authentik-outpost-image", &dockerprovider.RemoteImageArgs{
		Name:        pulumi.String(image),
		KeepLocally: pulumi.Bool(true),
	}, provOpt)
	if err != nil {
		return nil, fmt.Errorf("authentik outpost image: %w", err)
	}

	return dockerprovider.NewContainer(ctx, "authentik-outpost",
		&dockerprovider.ContainerArgs{
			Name:    pulumi.String("authentik-outpost"),
			Image:   outpostImage.RepoDigest,
			Restart: pulumi.String("unless-stopped"),
			Envs: pulumi.StringArray{
				pulumi.String("AUTHENTIK_HOST=" + authentikHost),
				pulumi.String("AUTHENTIK_TOKEN=" + token),
				pulumi.String("AUTHENTIK_INSECURE=false"),
			},
			NetworksAdvanced: dockerprovider.ContainerNetworksAdvancedArray{
				&dockerprovider.ContainerNetworksAdvancedArgs{
					Name:        pulumi.String("br1"),
					Ipv4Address: pulumi.String(system.Authentik.OutpostIPAddress),
				},
			},
			Labels: dockerprovider.ContainerLabelArray{
				&dockerprovider.ContainerLabelArgs{
					Label: pulumi.String("net.unraid.docker.icon"),
					Value: pulumi.String("https://cdn.jsdelivr.net/gh/homarr-labs/dashboard-icons/png/authentik.png"),
				},
			},
		},
		append([]pulumi.ResourceOption{provOpt,
			pulumi.IgnoreChanges([]string{"labels", "logOpts"}),
		}, importOpts("authentik-outpost", importIDs)...)...,
	)
}
