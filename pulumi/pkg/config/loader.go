package config

import (
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

// LoadSystem reads config/system.yaml.
func LoadSystem(configRoot string) (*SystemConfig, error) {
	data, err := os.ReadFile(filepath.Join(configRoot, "system.yaml"))
	if err != nil {
		return nil, fmt.Errorf("reading system.yaml: %w", err)
	}
	var cfg SystemConfig
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("parsing system.yaml: %w", err)
	}
	return &cfg, nil
}

// LoadServers reads config/servers.yaml.
func LoadServers(configRoot string) (*ServersFile, error) {
	data, err := os.ReadFile(filepath.Join(configRoot, "servers.yaml"))
	if err != nil {
		return nil, fmt.Errorf("reading servers.yaml: %w", err)
	}
	var sf ServersFile
	if err := yaml.Unmarshal(data, &sf); err != nil {
		return nil, fmt.Errorf("parsing servers.yaml: %w", err)
	}
	return &sf, nil
}

// LoadStacks reads all *.yaml files in config/stacks/ (skipping *.yaml.bak).
// Returns a map of stack name → StackConfig.
func LoadStacks(configRoot string) (map[string]StackConfig, error) {
	return loadYAMLDir[StackConfig](filepath.Join(configRoot, "stacks"))
}

// LoadServices reads all *.yaml files in config/services/.
// Returns a map of service name → ServiceConfig.
func LoadServices(configRoot string) (map[string]ServiceConfig, error) {
	return loadYAMLDir[ServiceConfig](filepath.Join(configRoot, "services"))
}

// loadYAMLDir reads every *.yaml file in dir and merges all top-level keys into one map.
func loadYAMLDir[T any](dir string) (map[string]T, error) {
	entries, err := filepath.Glob(filepath.Join(dir, "*.yaml"))
	if err != nil {
		return nil, err
	}
	result := map[string]T{}
	for _, f := range entries {
		data, err := os.ReadFile(f)
		if err != nil {
			return nil, fmt.Errorf("reading %s: %w", f, err)
		}
		var parsed map[string]T
		if err := yaml.Unmarshal(data, &parsed); err != nil {
			return nil, fmt.Errorf("parsing %s: %w", f, err)
		}
		for k, v := range parsed {
			result[k] = v
		}
	}
	return result, nil
}

// FlatServices merges stack-level env and mounts into each service definition and
// returns a flat list of all services. Stack services carry their stack's network
// definitions and generated secret names.
//
// The FlatService.ServiceName field is set to service.ServiceName (the YAML
// service_name value), which becomes the Pulumi logical resource name and
// matches the keys in config/tf-imports.yaml.
func FlatServices(stacks map[string]StackConfig, services map[string]ServiceConfig) []FlatService {
	var flat []FlatService

	// Stack services — merge stack-level env/mounts into each service.
	for stackName, stack := range stacks {
		for _, svcDef := range stack.Services {
			merged := svcDef
			// Stack env prepended; service-level env appended (takes precedence via later values).
			merged.Env = mergeEnv(stack.Env, svcDef.Env)
			merged.Mounts = append(append([]string{}, stack.Mounts...), svcDef.Mounts...)

			flat = append(flat, FlatService{
				ServiceName:      svcDef.ServiceName,
				Def:              merged,
				StackName:        stackName,
				StackNetworks:    stack.Networks,
				GeneratedSecrets: stack.GeneratedSecrets,
			})
		}
	}

	// Standalone services — no stack-level merging needed.
	for _, svcDef := range services {
		flat = append(flat, FlatService{
			ServiceName: svcDef.ServiceName,
			Def:         svcDef,
		})
	}

	return flat
}

// mergeEnv merges stack-level and service-level env slices.
// Service-level values for the same KEY override stack-level values.
func mergeEnv(stackEnv, svcEnv []string) []string {
	// Build map of service-level keys for deduplication.
	svcKeys := make(map[string]bool, len(svcEnv))
	for _, e := range svcEnv {
		if k := envKey(e); k != "" {
			svcKeys[k] = true
		}
	}
	var result []string
	// Append stack env entries that aren't overridden by the service.
	for _, e := range stackEnv {
		if !svcKeys[envKey(e)] {
			result = append(result, e)
		}
	}
	result = append(result, svcEnv...)
	return result
}

// envKey returns the KEY part of a "KEY=VALUE" env string.
func envKey(e string) string {
	for i, c := range e {
		if c == '=' {
			return e[:i]
		}
	}
	return e
}

// CollectGeneratedSecretNames returns the deduplicated list of all generated_secrets
// names across all stacks (used to pre-populate Pulumi config secrets).
func CollectGeneratedSecretNames(stacks map[string]StackConfig) []string {
	seen := map[string]bool{}
	var names []string
	for _, stack := range stacks {
		for _, name := range stack.GeneratedSecrets {
			if name != "" && !seen[name] {
				seen[name] = true
				names = append(names, name)
			}
		}
	}
	return names
}

// LoadImportIDs reads a YAML file of the form "resourceName: importID" and returns
// the map. Empty-string values are silently excluded.
func LoadImportIDs(path string) (map[string]string, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("reading import IDs file %s: %w", path, err)
	}
	raw := map[string]interface{}{}
	if err := yaml.Unmarshal(data, &raw); err != nil {
		return nil, fmt.Errorf("parsing import IDs file %s: %w", path, err)
	}
	ids := map[string]string{}
	for k, v := range raw {
		if v == nil {
			continue
		}
		switch val := v.(type) {
		case string:
			if val != "" {
				ids[k] = val
			}
		case int:
			ids[k] = fmt.Sprintf("%d", val)
		case uint64:
			ids[k] = fmt.Sprintf("%d", val)
		case float64:
			ids[k] = fmt.Sprintf("%d", int64(val))
		}
	}
	return ids, nil
}
