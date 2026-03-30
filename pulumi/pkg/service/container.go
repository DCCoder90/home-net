package service

import (
	"strings"

	"github.com/DCCoder90/home-net/pulumi/pkg/servers"
	"github.com/DCCoder90/home-net/pulumi/pkg/types"
	dockerprovider "github.com/pulumi/pulumi-docker/sdk/v4/go/docker"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// DeployContainer creates a Docker image pull + container for the given service.
// envVars may include both static env vars and async OAuth env vars.
func DeployContainer(
	ctx *pulumi.Context,
	svc *types.ServiceConfig,
	srv *servers.ServerContext,
	envVars pulumi.StringArrayInput,
	parent pulumi.Resource,
) (*dockerprovider.Container, error) {
	providerOpt := pulumi.Provider(srv.Provider)
	parentOpt := pulumi.Parent(parent)

	img, err := dockerprovider.NewRemoteImage(ctx, svc.ServiceName+"-image", &dockerprovider.RemoteImageArgs{
		Name:        pulumi.String(svc.ImageName),
		KeepLocally: pulumi.Bool(true),
	}, providerOpt, parentOpt)
	if err != nil {
		return nil, err
	}

	args := &dockerprovider.ContainerArgs{
		Name:    pulumi.String(svc.ServiceName),
		Image:   img.RepoDigest,
		Restart: pulumi.String("unless-stopped"),
		Envs:    envVars,
	}

	// Mounts (bind mounts from "host:container[:ro]" strings)
	if len(svc.Mounts) > 0 {
		mounts := make(dockerprovider.ContainerMountArray, 0, len(svc.Mounts))
		for _, m := range svc.Mounts {
			source, target, readOnly := parseMountString(m)
			mounts = append(mounts, &dockerprovider.ContainerMountArgs{
				Source:   pulumi.String(source),
				Target:   pulumi.String(target),
				Type:     pulumi.String("bind"),
				ReadOnly: pulumi.Bool(readOnly),
			})
		}
		args.Mounts = mounts
	}

	// Named volumes from "volume-name:container-path" strings
	if len(svc.Volumes) > 0 {
		vols := make(dockerprovider.ContainerVolumeArray, 0, len(svc.Volumes))
		for _, v := range svc.Volumes {
			parts := strings.SplitN(v, ":", 2)
			if len(parts) != 2 {
				continue
			}
			vols = append(vols, &dockerprovider.ContainerVolumeArgs{
				VolumeName:    pulumi.String(parts[0]),
				ContainerPath: pulumi.String(parts[1]),
			})
		}
		args.Volumes = vols
	}

	// Network mode and advanced network attachments
	if svc.NetworkMode != "" {
		args.NetworkMode = pulumi.String(svc.NetworkMode)
	} else if svc.Network != nil && len(svc.Network.Networks) > 0 {
		nets := make(dockerprovider.ContainerNetworksAdvancedArray, 0, len(svc.Network.Networks))
		for _, n := range svc.Network.Networks {
			netArgs := &dockerprovider.ContainerNetworksAdvancedArgs{
				Name: pulumi.String(n.Name),
			}
			if n.IPAddress != "" {
				netArgs.Ipv4Address = pulumi.String(n.IPAddress)
			}
			nets = append(nets, netArgs)
		}
		args.NetworksAdvanced = nets
	}

	// Docker capabilities
	if svc.Capabilities != nil {
		caps := &dockerprovider.ContainerCapabilitiesArgs{}
		if len(svc.Capabilities.Add) > 0 {
			adds := make(pulumi.StringArray, len(svc.Capabilities.Add))
			for i, a := range svc.Capabilities.Add {
				adds[i] = pulumi.String(a)
			}
			caps.Adds = adds
		}
		if len(svc.Capabilities.Drop) > 0 {
			drops := make(pulumi.StringArray, len(svc.Capabilities.Drop))
			for i, d := range svc.Capabilities.Drop {
				drops[i] = pulumi.String(d)
			}
			caps.Drops = drops
		}
		args.Capabilities = caps
	}

	// Command override
	if len(svc.Commands) > 0 {
		cmds := make(pulumi.StringArray, len(svc.Commands))
		for i, c := range svc.Commands {
			cmds[i] = pulumi.String(c)
		}
		args.Command = cmds
	}

	// GPU support
	if svc.EnableGPU {
		args.Runtime = pulumi.String("nvidia")
	}

	// Icon / webui labels
	labels := dockerprovider.ContainerLabelArray{}
	if svc.Icon != "" {
		labels = append(labels, &dockerprovider.ContainerLabelArgs{
			Label: pulumi.String("net.unraid.docker.icon"),
			Value: pulumi.String(svc.Icon),
		})
	}
	if len(labels) > 0 {
		args.Labels = labels
	}

	return dockerprovider.NewContainer(ctx, svc.ServiceName, args,
		providerOpt, parentOpt,
		pulumi.IgnoreChanges([]string{"logOpts"}),
	)
}

// parseMountString splits "host:container[:ro]" into its components.
func parseMountString(s string) (source, target string, readOnly bool) {
	parts := strings.SplitN(s, ":", 3)
	if len(parts) < 2 {
		return s, s, false
	}
	source = parts[0]
	target = parts[1]
	readOnly = len(parts) == 3 && parts[2] == "ro"
	return
}
