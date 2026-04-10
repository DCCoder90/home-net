package resources

import (
	"strings"

	"github.com/DCCoder90/home-net/pulumi/pkg/config"
	"github.com/DCCoder90/home-net/pulumi/pkg/npmproxy"
	"github.com/DCCoder90/home-net/pulumi/pkg/technitium"
	cloudflare "github.com/pulumi/pulumi-cloudflare/sdk/v5/go/cloudflare"
	dockerprovider "github.com/pulumi/pulumi-docker/sdk/v4/go/docker"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// RegisterDNSRecords creates Technitium internal DNS records and, for
// externally-facing services, Cloudflare public DNS records.
func RegisterDNSRecords(
	ctx *pulumi.Context,
	services []config.FlatService,
	system *config.SystemConfig,
	secrets map[string]string,
	containers map[string]*dockerprovider.Container,
	techProvider *technitium.Provider,
	publicFacingIP string,
	importIDs map[string]string,
) error {
	cfApiToken := secrets["cloudflare_api_token"]

	// Create the Cloudflare provider once — used for all external DNS records.
	cfProvider, err := cloudflare.NewProvider(ctx, "cloudflare-provider", &cloudflare.ProviderArgs{
		ApiToken: pulumi.ToSecret(pulumi.String(cfApiToken)).(pulumi.StringOutput),
	})
	if err != nil {
		return err
	}

	// Cache Cloudflare zone IDs to avoid redundant API calls per zone.
	zoneIDCache := map[string]string{}

	for _, svc := range services {
		if !svc.Def.DNS.Enabled {
			continue
		}
		domain := svc.Def.DNS.DomainName
		if domain == "" {
			continue
		}
		zone := zoneSuffix(domain, system.Zones)
		resourceName := svc.ServiceName + "-dns"

		// Container dependency — Technitium must already be running.
		var depOpts []pulumi.ResourceOption
		if containers[svc.ServiceName] != nil {
			depOpts = append(depOpts, pulumi.DependsOn([]pulumi.Resource{containers[svc.ServiceName]}))
		}
		techOpt := pulumi.Provider(techProvider)
		opts := append(append([]pulumi.ResourceOption{techOpt}, depOpts...), importOpts(resourceName, importIDs)...)

		// Technitium always points to NPM (the reverse proxy), which then forwards
		// to the service's actual IP:port via the proxy host entry.
		_, err := technitium.NewRecord(ctx, resourceName, zone, domain, "A", system.ProxyIP, 3600, opts...)
		if err != nil {
			return err
		}

		// isInternal: explicit dns.internal overrides network.internal.
		// nil = not set → fall back to network.internal.
		// *false = explicitly public (e.g. kavita: network.internal=true but dns.internal=false).
		// *true = explicitly internal.
		networkInternal := svc.Def.Network != nil && svc.Def.Network.Internal
		var isInternal bool
		if svc.Def.DNS.Internal != nil {
			isInternal = *svc.Def.DNS.Internal
		} else {
			isInternal = networkInternal
		}

		// Cloudflare record for public-facing services only.
		if !isInternal {
			cfResourceName := svc.ServiceName + "-cf-record"
			zoneID, cached := zoneIDCache[zone]
			if !cached {
				zoneID, err = cloudflareZoneID(ctx, zone, cfProvider)
				if err != nil {
					return err
				}
				zoneIDCache[zone] = zoneID
			}
			ip := publicFacingIP
			if ip == "" {
				ip = npmproxy.PublicIPPlaceholder
			}
			cfOpts := append(depOpts, importOpts(cfResourceName, importIDs)...)
			cfOpts = append(cfOpts, pulumi.Provider(cfProvider))
			_, err = cloudflare.NewRecord(ctx, cfResourceName, &cloudflare.RecordArgs{
				ZoneId:  pulumi.String(zoneID),
				Name:    pulumi.String(domain),
				Type:    pulumi.String("A"),
				Value:   pulumi.String(ip),
				Ttl:     pulumi.Int(1), // 1 = auto (proxied)
				Proxied: pulumi.Bool(true),
			}, cfOpts...)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

// zoneSuffix finds the matching zone name (e.g. "dcapi.app") for a domain.
func zoneSuffix(domain string, zones map[string]string) string {
	for _, zone := range zones {
		if strings.HasSuffix(domain, "."+zone) || domain == zone {
			return zone
		}
	}
	// Fall back to last two labels.
	parts := strings.Split(domain, ".")
	if len(parts) >= 2 {
		return strings.Join(parts[len(parts)-2:], ".")
	}
	return domain
}

// cloudflareZoneID looks up the Cloudflare zone ID for zoneName using the given provider.
func cloudflareZoneID(ctx *pulumi.Context, zoneName string, cfProvider *cloudflare.Provider) (string, error) {
	zone, err := cloudflare.LookupZone(ctx, &cloudflare.LookupZoneArgs{
		Name: pulumi.StringRef(zoneName),
	}, pulumi.Provider(cfProvider))
	if err != nil {
		return "", err
	}
	return zone.Id, nil
}
