package wrappers

import (
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// DNSRecord is the Pulumi state for a technitium:index:Record resource.
type DNSRecord struct {
	pulumi.CustomResourceState
}

// NewDNSRecord creates a Technitium DNS record via the technitium provider.
func NewDNSRecord(ctx *pulumi.Context, name string, args pulumi.Map, opts ...pulumi.ResourceOption) (*DNSRecord, error) {
	r := &DNSRecord{}
	if err := ctx.RegisterResource("technitium:index:Record", name, args, r, opts...); err != nil {
		return nil, err
	}
	return r, nil
}
