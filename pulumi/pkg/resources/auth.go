package resources

import (
	"fmt"
	"strings"

	"github.com/DCCoder90/home-net/pulumi/pkg/config"
	authentik "github.com/OSMIT-GmbH/pulumi-authentik/sdk/v2024/go/authentik"
	dockerprovider "github.com/pulumi/pulumi-docker/sdk/v4/go/docker"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// RegisterAuthResources creates Authentik groups, providers, applications, and policy
// bindings for all services that have auth.enabled: true.
// Called only when an authentik_token is present in secrets.
func RegisterAuthResources(
	ctx *pulumi.Context,
	services []config.FlatService,
	system *config.SystemConfig,
	secrets map[string]string,
	containers map[string]*dockerprovider.Container,
	importIDs map[string]string,
) error {
	authentikURL := fmt.Sprintf("http://%s:%d", system.Authentik.IPAddress, system.Authentik.Port)
	authentikToken := secrets["authentik_token"]

	// Instantiate the Authentik provider.
	authentikProvider, err := authentik.NewProvider(ctx, "authentik-provider", &authentik.ProviderArgs{
		Url:   pulumi.String(authentikURL),
		Token: pulumi.String(authentikToken),
	})
	if err != nil {
		return fmt.Errorf("authentik provider: %w", err)
	}
	authOpt := pulumi.Provider(authentikProvider)

	// Look up the standard authorization and invalidation flows.
	authorizationFlow, err := authentik.LookupFlow(ctx, &authentik.LookupFlowArgs{
		Slug: pulumi.StringRef("default-provider-authorization-implicit-consent"),
	}, authOpt)
	if err != nil {
		return fmt.Errorf("lookup authorization flow: %w", err)
	}

	invalidationFlow, err := authentik.LookupFlow(ctx, &authentik.LookupFlowArgs{
		Slug: pulumi.StringRef("default-provider-invalidation-flow"),
	}, authOpt)
	if err != nil {
		return fmt.Errorf("lookup invalidation flow: %w", err)
	}

	// Deduplicate groups — multiple services can share the same Authentik group.
	// Key is lowercased group name to prevent case-variation duplicates.
	groups := map[string]*authentik.Group{} // strings.ToLower(groupName) → resource

	for _, svc := range services {
		if svc.Def.Auth == nil || !svc.Def.Auth.Enabled {
			continue
		}

		groupName := svc.Def.Auth.Group
		if groupName == "" {
			groupName = svc.ServiceName
		}
		groupKey := strings.ToLower(groupName)

		// Create group if not already created.
		if _, exists := groups[groupKey]; !exists {
			// Resource name and import-ID key both use the lowercased group name so that
			// multiple services sharing the same group use the same import entry.
			resourceName := groupKey + "-group"
			grp, err := authentik.NewGroup(ctx, resourceName, &authentik.GroupArgs{
				Name: pulumi.String(groupName),
			}, append([]pulumi.ResourceOption{authOpt}, importOpts(resourceName, importIDs)...)...)
			if err != nil {
				return fmt.Errorf("authentik group %q: %w", groupName, err)
			}
			groups[groupKey] = grp
		}
		group := groups[groupKey]

		proxyAuth := svc.Def.Auth.Proxy != nil && svc.Def.Auth.Proxy.Enabled
		oauthAuth := svc.Def.Auth.OAuth != nil && svc.Def.Auth.OAuth.Enabled
		ip := serviceIP(svc.Def)
		domain := svc.Def.DNS.DomainName

		// ── Proxy Provider ───────────────────────────────────────────────────────
		if proxyAuth && ip != "" && domain != "" {
			port := 80
			if svc.Def.Network != nil && svc.Def.Network.ServicePort > 0 {
				port = svc.Def.Network.ServicePort
			}
			internalHost := fmt.Sprintf("http://%s:%d", ip, port)
			externalHost := "https://" + domain

			providerName := svc.ServiceName + "-proxy-provider"
			proxyProv, err := authentik.NewProviderProxy(ctx, providerName, &authentik.ProviderProxyArgs{
				Name:              pulumi.String(svc.ServiceName),
				InternalHost:      pulumi.String(internalHost),
				ExternalHost:      pulumi.String(externalHost),
				AuthorizationFlow: pulumi.String(authorizationFlow.Id),
				InvalidationFlow:  pulumi.String(invalidationFlow.Id),
			}, append([]pulumi.ResourceOption{authOpt}, importOpts(providerName, importIDs)...)...)
			if err != nil {
				return fmt.Errorf("proxy provider %q: %w", svc.ServiceName, err)
			}

			appName := svc.ServiceName + "-app"
			app, err := authentik.NewApplication(ctx, appName, &authentik.ApplicationArgs{
				Name: pulumi.String(svc.ServiceName),
				Slug: pulumi.String("tf_" + svc.ServiceName),
				ProtocolProvider: proxyProv.ID().ApplyT(func(id string) *int {
					n := 0
					fmt.Sscanf(id, "%d", &n)
					return &n
				}).(pulumi.IntPtrOutput),
			}, append([]pulumi.ResourceOption{authOpt}, importOpts(appName, importIDs)...)...)
			if err != nil {
				return fmt.Errorf("authentik app %q: %w", svc.ServiceName, err)
			}

			// Use "-proxy-binding" suffix to avoid name collision with OAuth binding.
			proxyBindingName := svc.ServiceName + "-binding"
			if oauthAuth {
				proxyBindingName = svc.ServiceName + "-proxy-binding"
			}
			_, err = authentik.NewPolicyBinding(ctx, proxyBindingName, &authentik.PolicyBindingArgs{
				Target:  app.Uuid,
				Group:   group.ID(),
				Order:   pulumi.Int(0),
				Enabled: pulumi.Bool(true),
			}, append([]pulumi.ResourceOption{authOpt}, importOpts(proxyBindingName, importIDs)...)...)
			if err != nil {
				return fmt.Errorf("policy binding %q: %w", svc.ServiceName, err)
			}
		}

		// ── OAuth2 Provider ──────────────────────────────────────────────────────
		if oauthAuth && domain != "" {
			oauth := svc.Def.Auth.OAuth

			// Build redirect URIs in the format Authentik expects: [{matching_mode, url}].
			var redirectURIs pulumi.StringMapArray
			for _, uri := range oauth.RedirectURIs {
				// Relative paths are expanded to full URLs.
				if !strings.HasPrefix(uri, "http") {
					uri = "https://" + domain + "/" + strings.TrimPrefix(uri, "/")
				}
				redirectURIs = append(redirectURIs, pulumi.StringMap{
					"matching_mode": pulumi.String("strict"),
					"url":           pulumi.String(uri),
				})
			}

			providerName := svc.ServiceName + "-oauth-provider"
			oauthProvArgs := &authentik.ProviderOauth2Args{
				Name:              pulumi.String(svc.ServiceName),
				AuthorizationFlow: pulumi.String(authorizationFlow.Id),
				InvalidationFlow:  pulumi.String(invalidationFlow.Id),
			}
			if len(redirectURIs) > 0 {
				oauthProvArgs.AllowedRedirectUris = redirectURIs
			}
			oauthProv, err := authentik.NewProviderOauth2(ctx, providerName, oauthProvArgs, append([]pulumi.ResourceOption{authOpt}, importOpts(providerName, importIDs)...)...)
			if err != nil {
				return fmt.Errorf("oauth2 provider %q: %w", svc.ServiceName, err)
			}

			appName := svc.ServiceName + "-app"
			// Use "-oauth-app" suffix only if service also has proxy auth (to avoid collision).
			if proxyAuth {
				appName = svc.ServiceName + "-oauth-app"
			}
			appSlug := strings.ToLower(svc.ServiceName)
			app, err := authentik.NewApplication(ctx, appName, &authentik.ApplicationArgs{
				Name: pulumi.String(svc.ServiceName),
				Slug: pulumi.String(appSlug),
				ProtocolProvider: oauthProv.ID().ApplyT(func(id string) *int {
					n := 0
					fmt.Sscanf(id, "%d", &n)
					return &n
				}).(pulumi.IntPtrOutput),
			}, append([]pulumi.ResourceOption{authOpt}, importOpts(appName, importIDs)...)...)
			if err != nil {
				return fmt.Errorf("authentik app %q: %w", svc.ServiceName, err)
			}

			bindingName := svc.ServiceName + "-binding"
			if proxyAuth {
				bindingName = svc.ServiceName + "-oauth-binding"
			}
			_, err = authentik.NewPolicyBinding(ctx, bindingName, &authentik.PolicyBindingArgs{
				Target:  app.Uuid,
				Group:   group.ID(),
				Order:   pulumi.Int(0),
				Enabled: pulumi.Bool(true),
			}, append([]pulumi.ResourceOption{authOpt}, importOpts(bindingName, importIDs)...)...)
			if err != nil {
				return fmt.Errorf("policy binding %q: %w", svc.ServiceName, err)
			}

			// Export OAuth2 client ID and OIDC issuer URL as stack outputs so services
			// like Grafana can read them (e.g. GF_AUTH_GENERIC_OAUTH_CLIENT_ID).
			// Full OIDC config is discoverable at <issuer>/.well-known/openid-configuration.
			ctx.Export(svc.ServiceName+"-oauth-client-id", oauthProv.ClientId)
			ctx.Export(svc.ServiceName+"-oauth-issuer-url", pulumi.String(authentikURL+"/application/o/"+appSlug+"/"))
		}
	}
	return nil
}
