package npmconfig

import (
	"github.com/DCCoder90/home-net/pulumi/pkg/wrappers"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// AccessLists holds references to the two standard NPM access lists.
type AccessLists struct {
	InternalOnly *wrappers.AccessList
	Cloudflare   *wrappers.AccessList
}

// Deploy creates the standard NPM access lists (Internal Only and Cloudflare).
// The returned AccessLists can be referenced when creating proxy hosts.
func Deploy(ctx *pulumi.Context) (*AccessLists, error) {
	internalOnly, err := wrappers.NewAccessList(ctx, "npm-acl-internal", pulumi.Map{
		"name": pulumi.String("Internal Only"),
		"items": pulumi.Array{
			pulumi.Map{"directive": pulumi.String("allow"), "address": pulumi.String("192.168.1.0/24")},
			pulumi.Map{"directive": pulumi.String("allow"), "address": pulumi.String("192.168.2.0/24")},
			pulumi.Map{"directive": pulumi.String("allow"), "address": pulumi.String("192.168.4.0/23")},
		},
		"passAuth":   pulumi.Bool(false),
		"satisfyAny": pulumi.Bool(false),
	})
	if err != nil {
		return nil, err
	}

	cloudflare, err := wrappers.NewAccessList(ctx, "npm-acl-cloudflare", pulumi.Map{
		"name": pulumi.String("Cloudflare"),
		"items": pulumi.Array{
			// Home network subnets
			pulumi.Map{"directive": pulumi.String("allow"), "address": pulumi.String("192.168.1.0/24")},
			pulumi.Map{"directive": pulumi.String("allow"), "address": pulumi.String("192.168.2.0/24")},
			pulumi.Map{"directive": pulumi.String("allow"), "address": pulumi.String("192.168.4.0/23")},
			// Cloudflare subnets
			pulumi.Map{"directive": pulumi.String("allow"), "address": pulumi.String("103.21.244.0/22")},
			pulumi.Map{"directive": pulumi.String("allow"), "address": pulumi.String("103.22.200.0/22")},
			pulumi.Map{"directive": pulumi.String("allow"), "address": pulumi.String("103.31.4.0/22")},
			pulumi.Map{"directive": pulumi.String("allow"), "address": pulumi.String("104.16.0.0/13")},
			pulumi.Map{"directive": pulumi.String("allow"), "address": pulumi.String("104.24.0.0/14")},
			pulumi.Map{"directive": pulumi.String("allow"), "address": pulumi.String("108.162.192.0/18")},
			pulumi.Map{"directive": pulumi.String("allow"), "address": pulumi.String("131.0.72.0/22")},
			pulumi.Map{"directive": pulumi.String("allow"), "address": pulumi.String("141.101.64.0/18")},
			pulumi.Map{"directive": pulumi.String("allow"), "address": pulumi.String("162.158.0.0/15")},
			pulumi.Map{"directive": pulumi.String("allow"), "address": pulumi.String("172.64.0.0/13")},
			pulumi.Map{"directive": pulumi.String("allow"), "address": pulumi.String("173.245.48.0/20")},
			pulumi.Map{"directive": pulumi.String("allow"), "address": pulumi.String("188.114.96.0/20")},
			pulumi.Map{"directive": pulumi.String("allow"), "address": pulumi.String("190.93.240.0/20")},
			pulumi.Map{"directive": pulumi.String("allow"), "address": pulumi.String("197.234.240.0/22")},
			pulumi.Map{"directive": pulumi.String("allow"), "address": pulumi.String("198.41.128.0/17")},
		},
		"passAuth":   pulumi.Bool(false),
		"satisfyAny": pulumi.Bool(false),
	})
	if err != nil {
		return nil, err
	}

	return &AccessLists{InternalOnly: internalOnly, Cloudflare: cloudflare}, nil
}
