package service

import "fmt"

// ResolveSecrets maps a service's secrets definition (envVar → infisicalKey)
// to actual "KEY=VALUE" env var strings using the pre-fetched secrets map.
func ResolveSecrets(secretsMap map[string]string, allSecrets map[string]string) []string {
	if len(secretsMap) == 0 {
		return nil
	}
	envs := make([]string, 0, len(secretsMap))
	for envVar, infisicalKey := range secretsMap {
		if val, ok := allSecrets[infisicalKey]; ok {
			envs = append(envs, fmt.Sprintf("%s=%s", envVar, val))
		}
	}
	return envs
}

// GetSecret returns a single secret value, falling back to an empty string if not found.
func GetSecret(key string, allSecrets map[string]string) string {
	return allSecrets[key]
}
