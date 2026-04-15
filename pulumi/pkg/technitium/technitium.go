// Package technitium provides thin Pulumi resource wrappers for the Technitium DNS
// custom provider (pulumi-resource-technitium). These types use ctx.RegisterResource
// so no generated SDK is required — the provider binary handles all CRUD logic.
package technitium

import (
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// Provider is the Pulumi provider resource for Technitium DNS.
type Provider struct {
	pulumi.ProviderResourceState
}

// NewProvider creates a Technitium DNS provider instance pointing at the given URL.
// token is the Technitium API token.
func NewProvider(ctx *pulumi.Context, name, url, token string, opts ...pulumi.ResourceOption) (*Provider, error) {
	var p Provider
	err := ctx.RegisterResource("pulumi:providers:technitium", name, pulumi.Map{
		"url":   pulumi.String(url),
		"token": pulumi.String(token),
	}, &p, opts...)
	return &p, err
}

// Record represents a single Technitium DNS record.
// The Pulumi resource ID has the composite format "name|zone|type|ttl".
type Record struct {
	pulumi.CustomResourceState
}

// NewRecord creates a Technitium DNS A record. domainName is the fully-qualified
// domain name (e.g. "grafana.dcapi.app"), zone is the zone (e.g. "dcapi.app").
func NewRecord(ctx *pulumi.Context, resourceName, zone, domainName, recType, ipAddress string, ttl int, opts ...pulumi.ResourceOption) (*Record, error) {
	var r Record
	err := ctx.RegisterResource("technitium:index:Record", resourceName, pulumi.Map{
		"zone":      pulumi.String(zone),
		"name":      pulumi.String(domainName),
		"type":      pulumi.String(recType),
		"ttl":       pulumi.Int(ttl),
		"ipAddress": pulumi.String(ipAddress),
	}, &r, opts...)
	return &r, err
}
