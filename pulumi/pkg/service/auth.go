package service

import (
	"fmt"

	"github.com/DCCoder90/home-net/pulumi/pkg/types"
	"github.com/DCCoder90/home-net/pulumi/pkg/wrappers"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// AuthResult contains the outcome of deploying Authentik auth for a service.
type AuthResult struct {
	// OAuthEnvVars is a Pulumi output containing OAuth env-var strings ("KEY=VALUE").
	// It resolves to an empty slice when oauth auth is not enabled.
	OAuthEnvVars pulumi.StringArrayOutput
}

// AuthInput holds everything needed to wire up Authentik auth for a service.
type AuthInput struct {
	Svc        *types.ServiceConfig
	System     *types.SystemConfig
	ServiceIP  string // primary IP of the service container
	AllSecrets map[string]string
	AdminUser  string // Authentik admin username to add to the access group
	Parent     pulumi.Resource
}

// DeployAuth creates Authentik resources (group, policy binding, provider, application)
// for a service. It returns an AuthResult whose OAuthEnvVars can be injected into the
// service container when oauth auth is enabled.
func DeployAuth(ctx *pulumi.Context, in *AuthInput) (*AuthResult, error) {
	svc := in.Svc
	emptyArray := pulumi.StringArray{}
	noopResult := &AuthResult{
		OAuthEnvVars: emptyArray.ToStringArrayOutput(),
	}

	if svc.Auth == nil || !svc.Auth.Enabled {
		return noopResult, nil
	}

	parentOpt := pulumi.Parent(in.Parent)
	groupName := fmt.Sprintf("tf_%s", svc.ServiceName)

	// Look up the admin user ID so we can add them to the access group.
	// We use a one-time Authentik API call through the provider's read mechanism
	// by creating the group with an empty user list first, then updating it.
	// In practice, the provider resource handles group membership directly.

	// --- Proxy auth ---
	if svc.Auth.Proxy != nil && svc.Auth.Proxy.Enabled {
		return deployProxyAuth(ctx, in, groupName, parentOpt)
	}

	// --- OAuth2 auth ---
	if svc.Auth.OAuth != nil && svc.Auth.OAuth.Enabled {
		return deployOAuth(ctx, in, groupName, parentOpt)
	}

	return noopResult, nil
}

func deployProxyAuth(
	ctx *pulumi.Context,
	in *AuthInput,
	groupName string,
	parentOpt pulumi.ResourceOption,
) (*AuthResult, error) {
	svc := in.Svc

	// Resolve basic-auth credentials from Infisical.
	userKey := fmt.Sprintf("%s_username", svc.ServiceName)
	passKey := fmt.Sprintf("%s_password", svc.ServiceName)
	if svc.Auth.Proxy.UserSecret != "" {
		userKey = svc.Auth.Proxy.UserSecret
	}
	if svc.Auth.Proxy.PassSecret != "" {
		passKey = svc.Auth.Proxy.PassSecret
	}

	groupAttrs := map[string]interface{}{
		fmt.Sprintf("%s_username", svc.ServiceName): GetSecret(userKey, in.AllSecrets),
		fmt.Sprintf("%s_password", svc.ServiceName): GetSecret(passKey, in.AllSecrets),
	}

	// Convert map to pulumi.Map for the wrapper.
	attrsMap := pulumi.Map{}
	for k, v := range groupAttrs {
		attrsMap[k] = pulumi.String(v.(string))
	}

	group, err := wrappers.NewAuthentikGroup(ctx, svc.ServiceName+"-group", pulumi.Map{
		"name":       pulumi.String(groupName),
		"attributes": attrsMap,
	}, parentOpt)
	if err != nil {
		return nil, err
	}

	internalHost := in.ServiceIP
	if internalHost == "" {
		internalHost = svc.ServiceName
	}

	// Standard Authentik flow slugs (created by default in every Authentik instance).
	authzFlow := "default-provider-authorization-implicit-consent"
	invalidationFlow := "default-invalidation-flow"

	provider, err := wrappers.NewAuthentikProviderProxy(ctx, svc.ServiceName+"-proxy-provider", pulumi.Map{
		"name":                       pulumi.String(svc.ServiceName),
		"internalHost":               pulumi.Sprintf("http://%s:%d", internalHost, svc.Network.ServicePort),
		"externalHost":               pulumi.Sprintf("http://%s", svc.DNS.DomainName),
		"authorizationFlow":          pulumi.String(authzFlow),
		"invalidationFlow":           pulumi.String(invalidationFlow),
		"basicAuthEnabled":           pulumi.Bool(true),
		"basicAuthUsernameAttribute": pulumi.Sprintf("%s_username", svc.ServiceName),
		"basicAuthPasswordAttribute": pulumi.Sprintf("%s_password", svc.ServiceName),
	}, parentOpt)
	if err != nil {
		return nil, err
	}

	app, err := wrappers.NewAuthentikApplication(ctx, svc.ServiceName+"-app", pulumi.Map{
		"name":        pulumi.String(svc.ServiceName),
		"slug":        pulumi.Sprintf("tf_%s", svc.ServiceName),
		"group":       pulumi.String(svc.Auth.Group),
		"providerId":  provider.ProviderID,
	}, parentOpt)
	if err != nil {
		return nil, err
	}

	if _, err := wrappers.NewAuthentikPolicyBinding(ctx, svc.ServiceName+"-binding", pulumi.Map{
		"target": app.UUID,
		"group":  group.GroupPK,
		"order":  pulumi.Int(0),
	}, parentOpt); err != nil {
		return nil, err
	}

	return &AuthResult{OAuthEnvVars: (pulumi.StringArray{}).ToStringArrayOutput()}, nil
}

func deployOAuth(
	ctx *pulumi.Context,
	in *AuthInput,
	groupName string,
	parentOpt pulumi.ResourceOption,
) (*AuthResult, error) {
	svc := in.Svc

	group, err := wrappers.NewAuthentikGroup(ctx, svc.ServiceName+"-group", pulumi.Map{
		"name": pulumi.String(groupName),
	}, parentOpt)
	if err != nil {
		return nil, err
	}

	authzFlow := "default-provider-authorization-implicit-consent"
	invalidationFlow := "default-invalidation-flow"

	// Build redirect URIs list.
	redirectURIs := pulumi.Array{
		pulumi.Map{
			"matchingMode": pulumi.String("strict"),
			"url":          pulumi.Sprintf("https://%s", svc.DNS.DomainName),
		},
	}
	for _, uri := range svc.Auth.OAuth.RedirectURIs {
		redirectURIs = append(redirectURIs, pulumi.Map{
			"matchingMode": pulumi.String("strict"),
			"url":          pulumi.String(uri),
		})
	}

	provider, err := wrappers.NewAuthentikProviderOAuth2(ctx, svc.ServiceName+"-oauth-provider", pulumi.Map{
		"name":                pulumi.String(svc.ServiceName),
		"clientId":            pulumi.String(svc.ServiceName),
		"authorizationFlow":   pulumi.String(authzFlow),
		"invalidationFlow":    pulumi.String(invalidationFlow),
		"allowedRedirectUris": redirectURIs,
	}, parentOpt)
	if err != nil {
		return nil, err
	}

	app, err := wrappers.NewAuthentikApplication(ctx, svc.ServiceName+"-app", pulumi.Map{
		"name":       pulumi.String(svc.ServiceName),
		"slug":       pulumi.Sprintf("%s", svc.ServiceName),
		"group":      pulumi.String(svc.Auth.Group),
		"providerId": provider.ProviderID,
	}, parentOpt)
	if err != nil {
		return nil, err
	}

	if _, err := wrappers.NewAuthentikPolicyBinding(ctx, svc.ServiceName+"-binding", pulumi.Map{
		"target": app.UUID,
		"group":  group.GroupPK,
		"order":  pulumi.Int(0),
	}, parentOpt); err != nil {
		return nil, err
	}

	// Build OAuth env vars from provider outputs.
	oauthKeys := svc.Auth.OAuth.Keys
	oauthEnvVars := pulumi.All(
		provider.ClientSecret,
		provider.IssuerURL,
	).ApplyT(func(args []interface{}) ([]string, error) {
		clientSecret := args[0].(string)
		issuerURL := args[1].(string)

		outputValues := map[string]string{
			"client_id":       svc.ServiceName,
			"client_secret":   clientSecret,
			"issuer_url":      issuerURL,
			"well_known_url":  issuerURL + "/.well-known/openid-configuration",
			"authorize_url":   issuerURL + "/authorize/",
			"token_url":       issuerURL + "/token/",
			"user_info_url":   issuerURL + "/userinfo/",
			"jwks_url":        issuerURL + "/jwks/",
			"logout_url":      issuerURL + "/end-session/",
			"sign_out_url":    issuerURL + "/end-session/",
		}

		envs := make([]string, 0, len(oauthKeys))
		for envVar, outputKey := range oauthKeys {
			if val, ok := outputValues[outputKey]; ok {
				envs = append(envs, fmt.Sprintf("%s=%s", envVar, val))
			}
		}
		return envs, nil
	}).(pulumi.StringArrayOutput)

	return &AuthResult{OAuthEnvVars: oauthEnvVars}, nil
}
