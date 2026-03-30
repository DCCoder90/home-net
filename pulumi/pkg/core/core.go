package core

import (
	"fmt"

	"github.com/DCCoder90/home-net/pulumi/pkg/servers"
	"github.com/DCCoder90/home-net/pulumi/pkg/types"
	dockerprovider "github.com/pulumi/pulumi-docker/sdk/v4/go/docker"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi/config"
)

// Deploy creates the Technitium DNS and Nginx Proxy Manager containers on Tower.
func Deploy(
	ctx *pulumi.Context,
	srv *servers.ServerContext,
	system *types.SystemConfig,
) error {
	cfg := config.New(ctx, "")

	tz := cfg.Get("tz")
	if tz == "" {
		tz = "America/Chicago"
	}
	hostPathDNS := cfg.Get("hostPathDnsConfig")
	hostPathNPMData := cfg.Get("hostPathNpmData")
	hostPathNPMLetsencrypt := cfg.Get("hostPathNpmLetsencrypt")
	hostPathNPMLogs := cfg.Get("hostPathNpmLogs")
	// Read as plain string — set with: pulumi config set --secret technitiumAdminPassword <value>
	dnsAdminPassword := cfg.Require("technitiumAdminPassword")

	providerOpt := pulumi.Provider(srv.Provider)

	// --- Technitium DNS ---
	dnsImage, err := dockerprovider.NewRemoteImage(ctx, "technitium-image", &dockerprovider.RemoteImageArgs{
		Name:        pulumi.String("technitium/dns-server:latest"),
		KeepLocally: pulumi.Bool(true),
	}, providerOpt)
	if err != nil {
		return err
	}

	_, err = dockerprovider.NewContainer(ctx, "technitium_dns", &dockerprovider.ContainerArgs{
		Name:     pulumi.String("technitium_dns"),
		Image:    dnsImage.RepoDigest,
		Restart:  pulumi.String("unless-stopped"),
		Hostname: pulumi.String("technitium-dns"),
		Envs: pulumi.StringArray{
			pulumi.String(fmt.Sprintf("TZ=%s", tz)),
			pulumi.String(fmt.Sprintf("DNS_SERVER_ADMIN_PASSWORD=%s", dnsAdminPassword)),
		},
		Volumes: dockerprovider.ContainerVolumeArray{
			&dockerprovider.ContainerVolumeArgs{
				HostPath:      pulumi.String(hostPathDNS),
				ContainerPath: pulumi.String("/etc/dns"),
			},
		},
		NetworksAdvanced: dockerprovider.ContainerNetworksAdvancedArray{
			&dockerprovider.ContainerNetworksAdvancedArgs{
				Name:        pulumi.String(srv.Config.Networks.Primary),
				Ipv4Address: pulumi.String(system.DNSIP),
			},
		},
		Labels: dockerprovider.ContainerLabelArray{
			&dockerprovider.ContainerLabelArgs{
				Label: pulumi.String("net.unraid.docker.icon"),
				Value: pulumi.String("https://cdn.jsdelivr.net/gh/homarr-labs/dashboard-icons/png/technitium.png"),
			},
		},
	}, providerOpt, pulumi.IgnoreChanges([]string{"image", "labels", "logOpts"}))
	if err != nil {
		return err
	}

	// --- Nginx Proxy Manager ---
	npmImage, err := dockerprovider.NewRemoteImage(ctx, "npm-image", &dockerprovider.RemoteImageArgs{
		Name:        pulumi.String("jc21/nginx-proxy-manager:latest"),
		KeepLocally: pulumi.Bool(true),
	}, providerOpt)
	if err != nil {
		return err
	}

	_, err = dockerprovider.NewContainer(ctx, "nginx_proxy_manager", &dockerprovider.ContainerArgs{
		Name:    pulumi.String("Nginx-Proxy-Manager-Official"),
		Image:   npmImage.RepoDigest,
		Restart: pulumi.String("unless-stopped"),
		Envs: pulumi.StringArray{
			pulumi.String(fmt.Sprintf("TZ=%s", tz)),
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
				Name:        pulumi.String(srv.Config.Networks.Primary),
				Ipv4Address: pulumi.String(system.ProxyIP),
			},
		},
		Labels: dockerprovider.ContainerLabelArray{
			&dockerprovider.ContainerLabelArgs{
				Label: pulumi.String("net.unraid.docker.icon"),
				Value: pulumi.String("https://cdn.jsdelivr.net/gh/homarr-labs/dashboard-icons/png/nginx-proxy-manager.png"),
			},
		},
	}, providerOpt, pulumi.IgnoreChanges([]string{"image", "labels", "envs", "logOpts", "ports", "volumes"}))

	return err
}
