// Package npmproxy provides thin Pulumi resource wrappers for the NPM Proxy Manager
// custom provider (pulumi-resource-npmproxy). No generated SDK is required.
package npmproxy

import (
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// Provider is the Pulumi provider resource for Nginx Proxy Manager.
type Provider struct {
	pulumi.ProviderResourceState
}

// NewProvider creates an NPM provider instance.
func NewProvider(ctx *pulumi.Context, name, url, username, password string, opts ...pulumi.ResourceOption) (*Provider, error) {
	var p Provider
	err := ctx.RegisterResource("pulumi:providers:npmproxy", name, pulumi.Map{
		"url":      pulumi.String(url),
		"username": pulumi.String(username),
		"password": pulumi.ToSecret(pulumi.String(password)),
	}, &p, opts...)
	return &p, err
}

// AccessList represents an NPM access list resource.
// The import ID is the numeric NPM access list ID as a string.
type AccessList struct {
	pulumi.CustomResourceState
}

// AccessListClient is a single client entry in an NPM access list.
type AccessListClient struct {
	// Address is a CIDR or IP range (e.g., "192.168.1.0/24").
	Address   string
	Directive string // "allow" or "deny"
}

// NewAccessList creates an NPM access list.
func NewAccessList(ctx *pulumi.Context, name string, listName string, passAuth bool, clients []AccessListClient, opts ...pulumi.ResourceOption) (*AccessList, error) {
	clientMaps := pulumi.Array{}
	for _, c := range clients {
		clientMaps = append(clientMaps, pulumi.Map{
			"address":   pulumi.String(c.Address),
			"directive": pulumi.String(c.Directive),
		})
	}
	var al AccessList
	err := ctx.RegisterResource("npmproxy:index:AccessList", name, pulumi.Map{
		"name":       pulumi.String(listName),
		"passAuth":   pulumi.Bool(passAuth),
		"satisfyAll": pulumi.Bool(false),
		"clients":    clientMaps,
	}, &al, opts...)
	return &al, err
}

// Certificate represents an NPM Let's Encrypt certificate.
// The import ID is the numeric NPM certificate ID as a string.
type Certificate struct {
	pulumi.CustomResourceState
}

// NewCertificate requests a Let's Encrypt certificate via DNS challenge.
func NewCertificate(ctx *pulumi.Context, name string, domainNames []string, email, cloudflareToken string, opts ...pulumi.ResourceOption) (*Certificate, error) {
	domains := make(pulumi.StringArray, len(domainNames))
	for i, d := range domainNames {
		domains[i] = pulumi.String(d)
	}
	var cert Certificate
	err := ctx.RegisterResource("npmproxy:index:Certificate", name, pulumi.Map{
		"domainNames": domains,
		"email":       pulumi.String(email),
		"cfApiToken":  pulumi.ToSecret(pulumi.String(cloudflareToken)),
	}, &cert, opts...)
	return &cert, err
}

// ProxyHost represents an NPM proxy host entry.
// The import ID is the numeric NPM proxy host ID as a string.
type ProxyHost struct {
	pulumi.CustomResourceState
}

// ProxyHostArgs defines inputs for NewProxyHost.
type ProxyHostArgs struct {
	DomainNames           []string
	ForwardHost           string
	ForwardPort           int
	ForwardScheme         string // "http" or "https"
	CertificateID         pulumi.IDOutput
	SSLForced             bool
	HTTP2Support          bool
	BlockExploits         bool
	AllowWebsocketUpgrade bool
	AdvancedConfig        string // raw nginx config injected into the server block
}

// NewProxyHost creates an NPM proxy host with TLS termination.
func NewProxyHost(ctx *pulumi.Context, name string, args ProxyHostArgs, opts ...pulumi.ResourceOption) (*ProxyHost, error) {
	domains := make(pulumi.StringArray, len(args.DomainNames))
	for i, d := range args.DomainNames {
		domains[i] = pulumi.String(d)
	}
	scheme := args.ForwardScheme
	if scheme == "" {
		scheme = "http"
	}
	var ph ProxyHost
	inputs := pulumi.Map{
		"domainNames":           domains,
		"forwardHost":           pulumi.String(args.ForwardHost),
		"forwardPort":           pulumi.Int(args.ForwardPort),
		"forwardScheme":         pulumi.String(scheme),
		"certificateId":         args.CertificateID,
		"accessListId":          pulumi.String("0"), // 0 = no ACL (public access)
		"sslForced":             pulumi.Bool(args.SSLForced),
		"http2Support":          pulumi.Bool(args.HTTP2Support),
		"blockExploits":         pulumi.Bool(args.BlockExploits),
		"allowWebsocketUpgrade": pulumi.Bool(args.AllowWebsocketUpgrade),
	}
	if args.AdvancedConfig != "" {
		inputs["advancedConfig"] = pulumi.String(args.AdvancedConfig)
	}
	err := ctx.RegisterResource("npmproxy:index:ProxyHost", name, inputs, &ph, opts...)
	return &ph, err
}
