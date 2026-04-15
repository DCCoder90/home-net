package main

import (
	"fmt"
	"os"

	"github.com/DCCoder90/home-net/pulumi/pkg/config"
	"github.com/DCCoder90/home-net/pulumi/pkg/npmproxy"
	"github.com/DCCoder90/home-net/pulumi/pkg/resources"
	"github.com/DCCoder90/home-net/pulumi/pkg/secrets"
	"github.com/DCCoder90/home-net/pulumi/pkg/technitium"
	dockerprovider "github.com/pulumi/pulumi-docker/sdk/v4/go/docker"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
	pulumiconfig "github.com/pulumi/pulumi/sdk/v3/go/pulumi/config"
)

func main() {
	pulumi.Run(func(ctx *pulumi.Context) error {
		cfg := pulumiconfig.New(ctx, "")

		// ── 1. Load all YAML config ────────────────────────────────────────────────
		configRoot := "../config"
		system, err := config.LoadSystem(configRoot)
		if err != nil {
			return err
		}
		serversFile, err := config.LoadServers(configRoot)
		if err != nil {
			return err
		}
		stacks, err := config.LoadStacks(configRoot)
		if err != nil {
			return err
		}
		services, err := config.LoadServices(configRoot)
		if err != nil {
			return err
		}
		flat := config.FlatServices(stacks, services)

		// ── 2. Initialise Infisical client ─────────────────────────────────────────
		clientID := cfg.Require("infisicalClientId")
		clientSecret := cfg.Require("infisicalClientSecret")

		infisicalClient, err := secrets.New(clientID, clientSecret,
			system.Infisical.Project, system.Infisical.Environment)
		if err != nil {
			return fmt.Errorf("infisical: %w", err)
		}

		// ── 3. Fetch secrets & server access ──────────────────────────────────────
		allSecrets, err := infisicalClient.FetchAll(system.Infisical.Folder)
		if err != nil {
			return fmt.Errorf("fetching secrets: %w", err)
		}
		serverAccess, err := infisicalClient.FetchServerAccess()
		if err != nil {
			return fmt.Errorf("fetching server access: %w", err)
		}

		// ── 3b. Fetch /config secrets (only if any service defines configfiles) ──
		configSecrets := map[string]string{}
		for _, svc := range flat {
			if len(svc.Def.ConfigFiles) > 0 {
				configSecrets, err = infisicalClient.FetchAll("/config")
				if err != nil {
					return fmt.Errorf("fetching config secrets: %w", err)
				}
				break
			}
		}

		// ── 4. Load import IDs (Option B import mechanism) ─────────────────────────
		var importIDs map[string]string
		if importFile := os.Getenv("PULUMI_IMPORT_IDS_FILE"); importFile != "" {
			importIDs, err = config.LoadImportIDs(importFile)
			if err != nil {
				return fmt.Errorf("loading import IDs: %w", err)
			}
		}

		// ── 5. Ensure generated secrets exist in Infisical /generated folder ─────────
		genSecretNames := config.CollectGeneratedSecretNames(stacks)
		plainGenSecrets, err := infisicalClient.EnsureGenerated(genSecretNames)
		if err != nil {
			return fmt.Errorf("ensuring generated secrets: %w", err)
		}
		generatedSecrets := make(map[string]pulumi.StringOutput, len(plainGenSecrets))
		for name, value := range plainGenSecrets {
			generatedSecrets[name] = pulumi.ToSecret(pulumi.String(value)).ApplyT(func(v any) (string, error) {
				return v.(string), nil
			}).(pulumi.StringOutput)
		}

		// ── 6. Build per-server Docker providers via SSH ───────────────────────────
		dockerProviders := map[string]*dockerprovider.Provider{}
		for serverName, serverCfg := range serversFile.Servers {
			access, ok := serverAccess[serverName]
			if !ok {
				return fmt.Errorf("no server_access entry for server %q", serverName)
			}
			_ = serverCfg // server-level network names are used per-container in YAML

			prov, err := dockerprovider.NewProvider(ctx, "docker-"+serverName, &dockerprovider.ProviderArgs{
				Host: pulumi.String(fmt.Sprintf("ssh://%s@%s:%d", access.User, access.Host, access.Port)),
				SshOpts: pulumi.StringArray{
					pulumi.String("-i"),
					pulumi.String(access.PrivKey),
					pulumi.String("-o StrictHostKeyChecking=no"),
					pulumi.String("-o BatchMode=yes"),
				},
				RegistryAuth: dockerprovider.ProviderRegistryAuthArray{
					&dockerprovider.ProviderRegistryAuthArgs{
						Address:  pulumi.String("ghcr.io"),
						Username: pulumi.String(allSecrets["ghcr_username"]),
						Password: pulumi.String(allSecrets["ghcr_token"]),
					},
				},
			})
			if err != nil {
				return fmt.Errorf("docker provider for %q: %w", serverName, err)
			}
			dockerProviders[serverName] = prov
		}

		towerProvider, ok := dockerProviders["tower"]
		if !ok {
			return fmt.Errorf("required server 'tower' not found in servers.yaml")
		}

		// ── 7. Register core containers ────────────────────────────────────────────
		_, err = resources.RegisterCoreContainers(ctx, system, allSecrets, towerProvider, importIDs)
		if err != nil {
			return fmt.Errorf("core containers: %w", err)
		}

		// ── 8. Register stack-internal Docker networks ─────────────────────────────
		networks, err := resources.RegisterNetworks(ctx, stacks, system.ExistingNetworks, towerProvider, importIDs)
		if err != nil {
			return fmt.Errorf("networks: %w", err)
		}

		// ── 8b. Write config files to remote hosts ────────────────────────────────
		configFileDeps, err := resources.RegisterConfigFiles(ctx, flat, configSecrets, serverAccess)
		if err != nil {
			return fmt.Errorf("config files: %w", err)
		}

		// ── 9. Register all service containers ────────────────────────────────────
		containers, err := resources.RegisterContainers(ctx, flat, dockerProviders, networks,
			allSecrets, generatedSecrets, importIDs, configFileDeps)
		if err != nil {
			return fmt.Errorf("containers: %w", err)
		}

		// ── 9b. Prepend synthetic Authentik service so DNS/proxy pipelines handle it ─
		if system.Authentik.DomainName != "" {
			publicDNS := false
			flat = append([]config.FlatService{{
				ServiceName: "authentik-server",
				Def: config.ServiceConfig{
					ServiceName: "authentik-server",
					DNS: config.ServiceDNS{
						Enabled:    true,
						DomainName: system.Authentik.DomainName,
						Internal:   &publicDNS, // explicit false → public CF record
					},
					Network: &config.ServiceNetwork{
						Internal:    false,
						ServicePort: system.Authentik.Port,
						Networks: []config.NetworkEntry{
							{Name: "br1", IPAddress: system.Authentik.IPAddress},
						},
					},
				},
			}}, flat...)
		}

		// ── 10. Register NPM resources (access lists, certs, proxy hosts) ──────────
		npmURL := cfg.Require("npmUrl")
		npmUser := cfg.Require("npmUsername")
		npmPass := cfg.Require("npmPassword")

		npmProv, err := npmproxy.NewProvider(ctx, "npm-provider", npmURL, npmUser, npmPass)
		if err != nil {
			return fmt.Errorf("npm provider: %w", err)
		}

		err = resources.RegisterProxyResources(ctx, flat, system, allSecrets,
			containers, npmProv, importIDs)
		if err != nil {
			return fmt.Errorf("proxy resources: %w", err)
		}

		// ── 11. Register Technitium DNS records ────────────────────────────────────
		techURL := cfg.Require("technitiumUrl")
		techToken := cfg.Require("technitiumToken")

		techProv, err := technitium.NewProvider(ctx, "technitium-provider", techURL, techToken)
		if err != nil {
			return fmt.Errorf("technitium provider: %w", err)
		}

		publicFacingIP := cfg.Get("publicFacingIp")
		err = resources.RegisterDNSRecords(ctx, flat, system, allSecrets, containers, techProv, publicFacingIP, importIDs)
		if err != nil {
			return fmt.Errorf("dns records: %w", err)
		}

		// ── 12. Register Authentik resources (phase-2 gate) ────────────────────────
		if token := allSecrets["authentik_token"]; token != "" {
			err = resources.RegisterAuthResources(ctx, flat, system, allSecrets, containers, importIDs)
			if err != nil {
				return fmt.Errorf("auth resources: %w", err)
			}
		}

		// ── 13. Deploy Authentik proxy outpost container (phase-3 gate) ────────────
		if _, err = resources.RegisterOutpostContainer(ctx, system, allSecrets, towerProvider, importIDs); err != nil {
			return fmt.Errorf("authentik outpost container: %w", err)
		}

		return nil
	})
}
