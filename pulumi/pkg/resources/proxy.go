package resources

import (
	"github.com/DCCoder90/home-net/pulumi/pkg/config"
	"github.com/DCCoder90/home-net/pulumi/pkg/npmproxy"
	dockerprovider "github.com/pulumi/pulumi-docker/sdk/v4/go/docker"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// RegisterProxyResources creates an NPM certificate and proxy host for every
// DNS-enabled service. Access lists are not used.
func RegisterProxyResources(
	ctx *pulumi.Context,
	services []config.FlatService,
	system *config.SystemConfig,
	secrets map[string]string,
	containers map[string]*dockerprovider.Container,
	npmProv *npmproxy.Provider,
	importIDs map[string]string,
) error {
	npmOpt := pulumi.Provider(npmProv)
	adminEmail := secrets["network_admin_email"]
	cfApiToken := secrets["cloudflare_api_token"]

	for _, svc := range services {
		if !svc.Def.DNS.Enabled || svc.Def.DNS.DomainName == "" {
			continue
		}

		domain := svc.Def.DNS.DomainName
		ip := serviceIP(svc.Def)
		if ip == "" {
			continue
		}
		port := 80
		if svc.Def.Network != nil && svc.Def.Network.ServicePort > 0 {
			port = svc.Def.Network.ServicePort
		}

		// When proxy auth is enabled, NPM forwards to the Authentik outpost, which
		// handles authentication and proxies authenticated requests to the service.
		// If a dedicated outpost IP is configured, use that; otherwise fall back to
		// the Authentik server (embedded outpost).
		forwardHost := ip
		forwardPort := port
		if svc.Def.Auth != nil && svc.Def.Auth.Proxy != nil && svc.Def.Auth.Proxy.Enabled {
			forwardHost = system.Authentik.IPAddress
			forwardPort = system.Authentik.Port
			if system.Authentik.OutpostIPAddress != "" {
				forwardHost = system.Authentik.OutpostIPAddress
				if system.Authentik.OutpostPort > 0 {
					forwardPort = system.Authentik.OutpostPort
				}
			}
		}

		var depOpts []pulumi.ResourceOption
		if c, ok := containers[svc.ServiceName]; ok {
			depOpts = append(depOpts, pulumi.DependsOn([]pulumi.Resource{c}))
		}

		// Certificate
		certName := svc.ServiceName + "-cert"
		cert, err := npmproxy.NewCertificate(ctx, certName, []string{domain}, adminEmail, cfApiToken,
			append(append([]pulumi.ResourceOption{npmOpt}, depOpts...), importOpts(certName, importIDs)...)...,
		)
		if err != nil {
			return err
		}

		// Proxy host
		proxyName := svc.ServiceName + "-proxy"
		_, err = npmproxy.NewProxyHost(ctx, proxyName, npmproxy.ProxyHostArgs{
			DomainNames:           []string{domain},
			ForwardHost:           forwardHost,
			ForwardPort:           forwardPort,
			ForwardScheme:         "http",
			CertificateID:         cert.ID(),
			SSLForced:             true,
			HTTP2Support:          true,
			BlockExploits:         true,
			AllowWebsocketUpgrade: true,
		}, append(append([]pulumi.ResourceOption{npmOpt}, depOpts...), importOpts(proxyName, importIDs)...)...)
		if err != nil {
			return err
		}
	}

	return nil
}
