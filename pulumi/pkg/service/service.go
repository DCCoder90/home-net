package service

import (
	"strings"

	"github.com/DCCoder90/home-net/pulumi/pkg/npmconfig"
	"github.com/DCCoder90/home-net/pulumi/pkg/servers"
	"github.com/DCCoder90/home-net/pulumi/pkg/types"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// DeployInput is everything needed to deploy a single service end-to-end.
type DeployInput struct {
	Name           string
	Config         *types.ServiceConfig
	Server         *servers.ServerContext
	System         *types.SystemConfig
	AllSecrets     map[string]string
	NPMAccessLists *npmconfig.AccessLists
	PublicIP       string
	CFApiToken     string
	AdminEmail     string
	AdminUsername  string
	// StackEnvs / StackMounts are merged from the parent stack (if any).
	StackEnvs   []string
	StackMounts []string
	// GeneratedSecrets maps secret-key names to their async Pulumi string outputs.
	// Values are substituted into env vars containing ${KEY} placeholders.
	GeneratedSecrets map[string]pulumi.StringOutput
	Parent           pulumi.Resource
}

// Deploy orchestrates auth → container → DNS for a single service.
func Deploy(ctx *pulumi.Context, in *DeployInput) error {
	svc := in.Config

	// Merge stack-level env/mounts into the service config copy.
	mergedEnvs := append(append([]string{}, in.StackEnvs...), svc.Env...)
	mergedMounts := append(append([]string{}, in.StackMounts...), svc.Mounts...)
	svc.Env = mergedEnvs
	svc.Mounts = mergedMounts

	// Resolve secret env vars synchronously from pre-fetched Infisical map.
	secretEnvs := ResolveSecrets(svc.Secrets, in.AllSecrets)

	// Deploy Authentik auth resources (if enabled).
	var networkEntries []types.NetworkEntry
	if svc.Network != nil {
		networkEntries = svc.Network.Networks
	}
	serviceIP := servers.PrimaryIP(networkEntries, in.Server)
	authResult, err := DeployAuth(ctx, &AuthInput{
		Svc:        svc,
		System:     in.System,
		ServiceIP:  serviceIP,
		AllSecrets: in.AllSecrets,
		AdminUser:  in.AdminUsername,
		Parent:     in.Parent,
	})
	if err != nil {
		return err
	}

	proxyAuthEnabled := svc.Auth != nil && svc.Auth.Enabled &&
		svc.Auth.Proxy != nil && svc.Auth.Proxy.Enabled

	// Combine static env vars with async OAuth env vars, substituting ${KEY} placeholders
	// from stack-level generated secrets.
	staticEnv := append(append([]string{}, svc.Env...), secretEnvs...)

	// Collect generated secret keys/outputs that are actually referenced in env vars.
	var genKeys []string
	var genOutputs []interface{}
	for key, out := range in.GeneratedSecrets {
		placeholder := "${" + key + "}"
		for _, e := range staticEnv {
			if strings.Contains(e, placeholder) {
				genKeys = append(genKeys, key)
				genOutputs = append(genOutputs, out)
				break
			}
		}
	}

	var allEnvs pulumi.StringArrayOutput
	if len(genOutputs) > 0 {
		// Wait for OAuthEnvVars + all needed generated secrets, then substitute.
		inputs := append([]interface{}{authResult.OAuthEnvVars}, genOutputs...)
		allEnvs = pulumi.All(inputs...).ApplyT(func(args []interface{}) ([]string, error) {
			oauthEnvs := args[0].([]string)
			resolved := make(map[string]string, len(genKeys))
			for i, k := range genKeys {
				resolved[k] = args[i+1].(string)
			}
			combined := append(append([]string{}, staticEnv...), oauthEnvs...)
			result := make([]string, len(combined))
			for i, env := range combined {
				v := env
				for k, val := range resolved {
					v = strings.ReplaceAll(v, "${"+k+"}", val)
				}
				result[i] = v
			}
			return result, nil
		}).(pulumi.StringArrayOutput)
	} else {
		allEnvs = authResult.OAuthEnvVars.ApplyT(func(oauthEnvs []string) []string {
			return append(staticEnv, oauthEnvs...)
		}).(pulumi.StringArrayOutput)
	}

	// Deploy Docker container.
	if _, err := DeployContainer(ctx, svc, in.Server, allEnvs, in.Parent); err != nil {
		return err
	}

	// Deploy DNS + reverse-proxy (if enabled and network is configured).
	if svc.DNS.Enabled && svc.Network != nil && in.NPMAccessLists != nil {
		if err := DeployDNS(ctx, &DNSInput{
			Svc:              svc,
			Srv:              in.Server,
			System:           in.System,
			NPMAccessLists:   in.NPMAccessLists,
			PublicIP:         in.PublicIP,
			CFApiToken:       in.CFApiToken,
			AdminEmail:       in.AdminEmail,
			ProxyAuthEnabled: proxyAuthEnabled,
			Parent:           in.Parent,
		}); err != nil {
			return err
		}
	}

	return nil
}
