package resources

import (
	"fmt"
	"strings"

	"github.com/DCCoder90/home-net/pulumi/pkg/config"
	dockerprovider "github.com/pulumi/pulumi-docker/sdk/v4/go/docker"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// RegisterContainers creates Docker images and containers for every FlatService.
// Returns a map of service_name → *dockerprovider.Container.
func RegisterContainers(
	ctx *pulumi.Context,
	services []config.FlatService,
	dockerProviders map[string]*dockerprovider.Provider,
	networks map[string]*dockerprovider.Network,
	secrets map[string]string,
	generatedSecrets map[string]pulumi.StringOutput,
	importIDs map[string]string,
) (map[string]*dockerprovider.Container, error) {
	result := make(map[string]*dockerprovider.Container, len(services))

	for _, svc := range services {
		container, err := registerContainer(ctx, svc, dockerProviders, networks, secrets, generatedSecrets, importIDs)
		if err != nil {
			return nil, fmt.Errorf("service %q: %w", svc.ServiceName, err)
		}
		result[svc.ServiceName] = container
	}
	return result, nil
}

func registerContainer(
	ctx *pulumi.Context,
	svc config.FlatService,
	dockerProviders map[string]*dockerprovider.Provider,
	networks map[string]*dockerprovider.Network,
	secrets map[string]string,
	generatedSecrets map[string]pulumi.StringOutput,
	importIDs map[string]string,
) (*dockerprovider.Container, error) {
	// Determine which Docker provider to use (defaults to "tower").
	hostName := svc.Def.Host
	if hostName == "" {
		hostName = "tower"
	}
	provider, ok := dockerProviders[hostName]
	if !ok {
		return nil, fmt.Errorf("no Docker provider for host %q", hostName)
	}
	provOpt := pulumi.Provider(provider)

	// Pull image.
	imageName := svc.Def.ImageName
	img, err := dockerprovider.NewRemoteImage(ctx, svc.ServiceName+"-image", &dockerprovider.RemoteImageArgs{
		Name:        pulumi.String(imageName),
		KeepLocally: pulumi.Bool(true),
	}, provOpt)
	if err != nil {
		return nil, fmt.Errorf("image %q: %w", imageName, err)
	}

	// Build environment variables.
	// 1. Static envs from YAML (already merged from stack + service level in FlatServices).
	staticEnvs := append([]string{}, svc.Def.Env...)
	// 2. Inject secrets (Infisical key → env var value).
	for envVar, secretKey := range svc.Def.Secrets {
		if val, ok := secrets[secretKey]; ok {
			staticEnvs = append(staticEnvs, envVar+"="+val)
		}
	}

	// 3. Resolve ${GENERATED_SECRET} placeholders — requires async Pulumi outputs.
	allEnvs := buildEnvOutput(staticEnvs, svc.GeneratedSecrets, generatedSecrets)

	// Build volume mounts.
	var volumes dockerprovider.ContainerVolumeArray
	for _, mount := range svc.Def.Mounts {
		// Mounts are "host:container[:ro]" strings.
		parts := strings.SplitN(mount, ":", 3)
		if len(parts) < 2 {
			continue
		}
		vol := &dockerprovider.ContainerVolumeArgs{
			HostPath:      pulumi.String(parts[0]),
			ContainerPath: pulumi.String(parts[1]),
		}
		if len(parts) == 3 && parts[2] == "ro" {
			vol.ReadOnly = pulumi.Bool(true)
		}
		volumes = append(volumes, vol)
	}

	// Build network attachments.
	var networksAdvanced dockerprovider.ContainerNetworksAdvancedArray
	if svc.Def.Network != nil && svc.Def.NetworkMode != "host" {
		for _, ne := range svc.Def.Network.Networks {
			entry := &dockerprovider.ContainerNetworksAdvancedArgs{
				Name: pulumi.String(ne.Name),
			}
			if ne.IPAddress != "" {
				entry.Ipv4Address = pulumi.String(ne.IPAddress)
			}
			networksAdvanced = append(networksAdvanced, entry)
		}
	}

	// Build capabilities.
	var caps *dockerprovider.ContainerCapabilitiesArgs
	if svc.Def.Capabilities != nil {
		var adds, drops pulumi.StringArray
		for _, c := range svc.Def.Capabilities.Add {
			adds = append(adds, pulumi.String(c))
		}
		for _, c := range svc.Def.Capabilities.Drop {
			drops = append(drops, pulumi.String(c))
		}
		caps = &dockerprovider.ContainerCapabilitiesArgs{
			Adds:  adds,
			Drops: drops,
		}
	}

	// Build command.
	var command pulumi.StringArray
	for _, cmd := range svc.Def.Commands {
		command = append(command, pulumi.String(cmd))
	}

	// Build device list from the devices config block.
	var devices dockerprovider.ContainerDeviceArray
	if d := svc.Def.Devices; d != nil {
		if d.GPU {
			devices = append(devices, &dockerprovider.ContainerDeviceArgs{
				HostPath:      pulumi.String("/dev/dri"),
				ContainerPath: pulumi.String("/dev/dri"),
			})
		}
		if d.USB {
			devices = append(devices, &dockerprovider.ContainerDeviceArgs{
				HostPath:      pulumi.String("/dev/bus/usb"),
				ContainerPath: pulumi.String("/dev/bus/usb"),
			})
		}
		for _, path := range d.Paths {
			parts := strings.SplitN(path, ":", 2)
			if len(parts) != 2 {
				continue
			}
			devices = append(devices, &dockerprovider.ContainerDeviceArgs{
				HostPath:      pulumi.String(parts[0]),
				ContainerPath: pulumi.String(parts[1]),
			})
		}
	}

	// Build labels.
	var labels dockerprovider.ContainerLabelArray
	if svc.Def.Icon != "" {
		labels = append(labels, &dockerprovider.ContainerLabelArgs{
			Label: pulumi.String("net.unraid.docker.icon"),
			Value: pulumi.String(svc.Def.Icon),
		})
	}
	if svc.Def.Network != nil && svc.Def.Network.ServicePort > 0 {
		ip := serviceIP(svc.Def)
		if ip != "" {
			labels = append(labels, &dockerprovider.ContainerLabelArgs{
				Label: pulumi.String("net.unraid.docker.webui"),
				Value: pulumi.String(fmt.Sprintf("http://%s:%d", ip, svc.Def.Network.ServicePort)),
			})
		}
	}

	args := &dockerprovider.ContainerArgs{
		Name:             pulumi.String(svc.Def.ServiceName),
		Image:            img.RepoDigest,
		Restart:          pulumi.String("unless-stopped"),
		Envs:             allEnvs,
		Volumes:          volumes,
		NetworksAdvanced: networksAdvanced,
		Labels:           labels,
		Command:          command,
		Devices:          devices,
	}

	if svc.Def.NetworkMode != "" {
		args.NetworkMode = pulumi.String(svc.Def.NetworkMode)
		// When using host mode, clear networksAdvanced.
		if svc.Def.NetworkMode == "host" {
			args.NetworksAdvanced = nil
		}
	}
	if caps != nil {
		args.Capabilities = caps
	}


	// Determine which network-related changes to ignore on existing containers.
	ignoreFields := []string{"image", "labels", "logOpts"}

	opts := []pulumi.ResourceOption{
		provOpt,
		pulumi.IgnoreChanges(ignoreFields),
	}
	opts = append(opts, importOpts(svc.ServiceName, importIDs)...)

	// Add network dependencies so containers wait for custom networks.
	for _, ne := range networkEntries(svc.Def.Network) {
		if net, ok := networks[ne.Name]; ok {
			opts = append(opts, pulumi.DependsOn([]pulumi.Resource{net}))
		}
	}

	container, err := dockerprovider.NewContainer(ctx, svc.ServiceName, args, opts...)
	if err != nil {
		return nil, err
	}
	return container, nil
}

// buildEnvOutput returns a StringArrayOutput that:
//  1. Injects every key in neededKeys as KEY=value into the env list.
//  2. Also substitutes any ${KEY} placeholders in staticEnvs (for cases where
//     the env var name differs from the generated secret key name).
//
// If there are no generated secrets, returns a resolved static array.
func buildEnvOutput(staticEnvs []string, neededKeys []string, generatedSecrets map[string]pulumi.StringOutput) pulumi.StringArrayOutput {
	var keys []string
	var outputs []any
	for _, key := range neededKeys {
		if out, ok := generatedSecrets[key]; ok {
			keys = append(keys, key)
			outputs = append(outputs, out)
		}
	}

	if len(keys) == 0 {
		arr := make(pulumi.StringArray, len(staticEnvs))
		for i, e := range staticEnvs {
			arr[i] = pulumi.String(e)
		}
		return arr.ToStringArrayOutput()
	}

	return pulumi.All(outputs...).ApplyT(func(vals []any) ([]string, error) {
		resolved := make(map[string]string, len(keys))
		for i, k := range keys {
			s, ok := vals[i].(string)
			if !ok {
				return nil, fmt.Errorf("generated secret %q resolved to non-string type %T", k, vals[i])
			}
			resolved[k] = s
		}

		// Substitute any ${KEY} placeholders in existing env entries.
		result := make([]string, len(staticEnvs))
		for i, e := range staticEnvs {
			v := e
			for k, val := range resolved {
				v = strings.ReplaceAll(v, "${"+k+"}", val)
			}
			result[i] = v
		}

		// Auto-inject each generated secret as KEY=value.
		for _, k := range keys {
			result = append(result, k+"="+resolved[k])
		}

		return result, nil
	}).(pulumi.StringArrayOutput)
}

// serviceIP returns the primary static IP of a service (first br1 entry, then any).
func serviceIP(svc config.ServiceConfig) string {
	if svc.Network == nil {
		return ""
	}
	for _, ne := range svc.Network.Networks {
		if ne.Name == "br1" && ne.IPAddress != "" {
			return ne.IPAddress
		}
	}
	for _, ne := range svc.Network.Networks {
		if ne.IPAddress != "" {
			return ne.IPAddress
		}
	}
	return ""
}

// networkEntries returns the network attachment list from a (possibly nil) ServiceNetwork.
func networkEntries(sn *config.ServiceNetwork) []config.NetworkEntry {
	if sn == nil {
		return nil
	}
	return sn.Networks
}
