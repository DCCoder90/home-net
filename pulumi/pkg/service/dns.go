package service

import (
	"fmt"
	"strings"

	"github.com/DCCoder90/home-net/pulumi/pkg/npmconfig"
	"github.com/DCCoder90/home-net/pulumi/pkg/servers"
	"github.com/DCCoder90/home-net/pulumi/pkg/types"
	"github.com/DCCoder90/home-net/pulumi/pkg/wrappers"
	cloudflare "github.com/pulumi/pulumi-cloudflare/sdk/v4/go/cloudflare"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// DNSInput holds everything needed to wire up DNS and reverse-proxy for a service.
type DNSInput struct {
	Svc              *types.ServiceConfig
	Srv              *servers.ServerContext
	System           *types.SystemConfig
	NPMAccessLists   *npmconfig.AccessLists
	PublicIP         string
	CFApiToken       string
	AdminEmail       string
	ProxyAuthEnabled bool // if true, NPM should forward to Authentik
	Parent           pulumi.Resource
}

// DeployDNS creates Technitium + Cloudflare DNS records and the NPM certificate + proxy host.
func DeployDNS(ctx *pulumi.Context, in *DNSInput) error {
	svc := in.Svc
	if !svc.DNS.Enabled {
		return nil
	}

	domain := svc.DNS.DomainName
	serviceIP := servers.PrimaryIP(svc.Network.Networks, in.Srv)
	parentOpt := pulumi.Parent(in.Parent)

	// --- Technitium internal A record ---
	if _, err := wrappers.NewDNSRecord(ctx, svc.ServiceName+"-dns", pulumi.Map{
		"domain":    pulumi.String(domain),
		"zone":      pulumi.String(extractZone(domain)),
		"type":      pulumi.String("A"),
		"ipAddress": pulumi.String(serviceIP),
		"ttl":       pulumi.Int(3600),
	}, parentOpt); err != nil {
		return fmt.Errorf("creating Technitium record for %s: %w", svc.ServiceName, err)
	}

	// --- Cloudflare external A record (skip for internal-only services) ---
	if !svc.DNS.Internal && in.PublicIP != "" {
		zoneName := extractZone(domain)
		zone, err := cloudflare.LookupZone(ctx, &cloudflare.LookupZoneArgs{
			Name: &zoneName,
		})
		if err != nil {
			return fmt.Errorf("looking up Cloudflare zone %q: %w", zoneName, err)
		}

		if _, err := cloudflare.NewRecord(ctx, svc.ServiceName+"-cf-record", &cloudflare.RecordArgs{
			ZoneId:  pulumi.String(zone.ZoneId),
			Name:    pulumi.String(domain),
			Type:    pulumi.String("A"),
			Value:   pulumi.String(in.PublicIP),
			Ttl:     pulumi.Int(1), // 1 = auto when proxied=true
			Proxied: pulumi.Bool(true),
		}, parentOpt); err != nil {
			return fmt.Errorf("creating Cloudflare record for %s: %w", svc.ServiceName, err)
		}
	}

	// --- NPM Let's Encrypt certificate ---
	cert, err := wrappers.NewCertificate(ctx, svc.ServiceName+"-cert", pulumi.Map{
		"domainName":         pulumi.String(domain),
		"letsEncryptEmail":   pulumi.String(in.AdminEmail),
		"cloudflareApiToken": pulumi.String(in.CFApiToken),
		"propagationSeconds": pulumi.Int(30),
	}, parentOpt)
	if err != nil {
		return fmt.Errorf("creating NPM certificate for %s: %w", svc.ServiceName, err)
	}

	// Determine forward host/port based on whether proxy auth is enabled.
	var forwardHost pulumi.StringInput
	var forwardPort pulumi.IntInput
	if in.ProxyAuthEnabled {
		forwardHost = pulumi.String(in.System.Authentik.IPAddress)
		forwardPort = pulumi.Int(in.System.Authentik.Port)
	} else {
		host := serviceIP
		if host == "" {
			host = svc.ServiceName // Docker network name resolution fallback
		}
		forwardHost = pulumi.String(host)
		forwardPort = pulumi.Int(svc.Network.ServicePort)
	}

	// Pick access list ID based on internal flag.
	var accessListID pulumi.IntInput
	if svc.Network.Internal {
		accessListID = in.NPMAccessLists.InternalOnly.AccessListID
	} else {
		accessListID = in.NPMAccessLists.Cloudflare.AccessListID
	}

	// --- NPM proxy host ---
	if _, err = wrappers.NewProxyHost(ctx, svc.ServiceName+"-proxy", pulumi.Map{
		"domainName":            pulumi.String(domain),
		"forwardScheme":         pulumi.String("http"),
		"forwardHost":           forwardHost,
		"forwardPort":           forwardPort,
		"accessListId":          accessListID,
		"certificateId":         cert.CertificateID,
		"cachingEnabled":        pulumi.Bool(true),
		"allowWebsocketUpgrade": pulumi.Bool(true),
		"blockExploits":         pulumi.Bool(true),
		"http2Support":          pulumi.Bool(true),
		"sslForced":             pulumi.Bool(false),
		"hstsEnabled":           pulumi.Bool(false),
		"hstsSubdomains":        pulumi.Bool(false),
	}, parentOpt); err != nil {
		return fmt.Errorf("creating NPM proxy host for %s: %w", svc.ServiceName, err)
	}

	return nil
}

// extractZone returns the registrable zone (last two DNS labels) from a FQDN.
// e.g. "chat.mallett.family" → "mallett.family"
func extractZone(domain string) string {
	parts := strings.Split(domain, ".")
	if len(parts) < 2 {
		return domain
	}
	return strings.Join(parts[len(parts)-2:], ".")
}
