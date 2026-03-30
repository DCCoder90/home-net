package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/DCCoder90/home-net/pulumi/pkg/core"
	"github.com/DCCoder90/home-net/pulumi/pkg/infisical"
	"github.com/DCCoder90/home-net/pulumi/pkg/npmconfig"
	"github.com/DCCoder90/home-net/pulumi/pkg/servers"
	"github.com/DCCoder90/home-net/pulumi/pkg/service"
	"github.com/DCCoder90/home-net/pulumi/pkg/stack"
	"github.com/DCCoder90/home-net/pulumi/pkg/types"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi/config"
	"gopkg.in/yaml.v3"
)

func main() {
	pulumi.Run(func(ctx *pulumi.Context) error {
		cfg := config.New(ctx, "")

		// --- Load system.yaml ---
		systemRaw, err := os.ReadFile("../config/system.yaml")
		if err != nil {
			return fmt.Errorf("reading system.yaml: %w", err)
		}
		var system types.SystemConfig
		if err := yaml.Unmarshal(systemRaw, &system); err != nil {
			return fmt.Errorf("parsing system.yaml: %w", err)
		}

		// Read config values as plain strings.
		// Values stored with `pulumi config set --secret` are still decrypted here;
		// using Require() (vs RequireSecret) means they won't be tagged as outputs,
		// which is fine for values used only in API calls.
		infisicalClientID := cfg.Require("infisicalClientId")
		infisicalClientSecret := cfg.Require("infisicalClientSecret")
		infisicalProject := cfg.Require("infisicalProject")
		infisicalEnv := cfg.Get("infisicalEnvironment")
		if infisicalEnv == "" {
			infisicalEnv = system.Infisical.Environment
		}

		publicIP := cfg.Get("publicFacingIp")
		cfApiToken := cfg.Require("cloudflareApiToken")
		adminEmail := cfg.Require("networkAdminEmail")
		adminUsername := cfg.Require("networkAdminUsername")

		// --- Infisical client ---
		infisicalClient, err := infisical.New(
			infisicalClientID,
			infisicalClientSecret,
			infisicalProject,
			infisicalEnv,
		)
		if err != nil {
			return fmt.Errorf("creating Infisical client: %w", err)
		}

		// Pre-fetch all service-level secrets once.
		allSecrets, err := infisicalClient.FetchAll(system.Infisical.Folder)
		if err != nil {
			return fmt.Errorf("fetching secrets from Infisical: %w", err)
		}

		// --- Server contexts (one Docker provider per host) ---
		serverMap, err := servers.Load(ctx, infisicalClient)
		if err != nil {
			return fmt.Errorf("loading servers: %w", err)
		}

		// --- Core infrastructure (NPM + Technitium on tower) ---
		towerSrv, err := servers.Require(serverMap, "tower")
		if err != nil {
			return err
		}
		if err := core.Deploy(ctx, towerSrv, &system); err != nil {
			return fmt.Errorf("deploying core: %w", err)
		}

		// --- NPM access lists ---
		accessLists, err := npmconfig.Deploy(ctx)
		if err != nil {
			return fmt.Errorf("deploying NPM access lists: %w", err)
		}

		// --- Load and deploy stacks ---
		stacks, err := loadStacks("../config/stacks")
		if err != nil {
			return fmt.Errorf("loading stacks: %w", err)
		}
		for stackName, stackCfg := range stacks {
			cfg := stackCfg // capture
			if err := stack.Deploy(ctx, &stack.DeployInput{
				Name:           stackName,
				Config:         &cfg,
				Servers:        serverMap,
				System:         &system,
				AllSecrets:     allSecrets,
				NPMAccessLists: accessLists,
				PublicIP:       publicIP,
				CFApiToken:     cfApiToken,
				AdminEmail:     adminEmail,
				AdminUsername:  adminUsername,
			}); err != nil {
				return fmt.Errorf("deploying stack %q: %w", stackName, err)
			}
		}

		// --- Load and deploy standalone services ---
		services, err := loadServices("../config/services")
		if err != nil {
			return fmt.Errorf("loading services: %w", err)
		}
		for svcName, svcCfg := range services {
			cfg := svcCfg // capture
			hostName := servers.ServiceHost(cfg.Host)
			srv, err := servers.Require(serverMap, hostName)
			if err != nil {
				return fmt.Errorf("service %q: %w", svcName, err)
			}
			if err := service.Deploy(ctx, &service.DeployInput{
				Name:           svcName,
				Config:         &cfg,
				Server:         srv,
				System:         &system,
				AllSecrets:     allSecrets,
				NPMAccessLists: accessLists,
				PublicIP:       publicIP,
				CFApiToken:     cfApiToken,
				AdminEmail:     adminEmail,
				AdminUsername:  adminUsername,
			}); err != nil {
				return fmt.Errorf("deploying service %q: %w", svcName, err)
			}
		}

		return nil
	})
}

// loadStacks reads all *.yaml files in dir and merges them into one map[stackName]StackConfig.
func loadStacks(dir string) (map[string]types.StackConfig, error) {
	result := map[string]types.StackConfig{}
	entries, err := filepath.Glob(filepath.Join(dir, "*.yaml"))
	if err != nil {
		return nil, err
	}
	for _, f := range entries {
		raw, err := os.ReadFile(f)
		if err != nil {
			return nil, fmt.Errorf("reading %s: %w", f, err)
		}
		var parsed map[string]types.StackConfig
		if err := yaml.Unmarshal(raw, &parsed); err != nil {
			return nil, fmt.Errorf("parsing %s: %w", f, err)
		}
		for k, v := range parsed {
			result[k] = v
		}
	}
	return result, nil
}

// loadServices reads all *.yaml files in dir and merges them into one map[svcName]ServiceConfig.
func loadServices(dir string) (map[string]types.ServiceConfig, error) {
	result := map[string]types.ServiceConfig{}
	entries, err := filepath.Glob(filepath.Join(dir, "*.yaml"))
	if err != nil {
		return nil, err
	}
	for _, f := range entries {
		raw, err := os.ReadFile(f)
		if err != nil {
			return nil, fmt.Errorf("reading %s: %w", f, err)
		}
		var parsed map[string]types.ServiceConfig
		if err := yaml.Unmarshal(raw, &parsed); err != nil {
			return nil, fmt.Errorf("parsing %s: %w", f, err)
		}
		for k, v := range parsed {
			result[k] = v
		}
	}
	return result, nil
}
