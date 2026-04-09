// Command pulumi-resource-npmproxy is a Pulumi resource provider for Nginx Proxy Manager.
// It implements three resource types:
//   - npmproxy:index:AccessList
//   - npmproxy:index:Certificate
//   - npmproxy:index:ProxyHost
package main

import (
	"context"
	"fmt"
	"os"

	p "github.com/pulumi/pulumi-go-provider"
	"github.com/pulumi/pulumi-go-provider/infer"
	"github.com/pulumi/pulumi/sdk/v3/go/common/tokens"

	"github.com/DCCoder90/home-net/pulumi/providers/npmproxy/resources"
)

func main() {
	provider := infer.Provider(infer.Options{
		Resources: []infer.InferredResource{
			infer.Resource[*resources.AccessList, resources.AccessListArgs, resources.AccessListState](&resources.AccessList{}),
			infer.Resource[*resources.Certificate, resources.CertificateArgs, resources.CertificateState](&resources.Certificate{}),
			infer.Resource[*resources.ProxyHost, resources.ProxyHostArgs, resources.ProxyHostState](&resources.ProxyHost{}),
		},
		Config: infer.Config[*resources.Config](&resources.Config{}),
		// Map the "resources" Go package to the "index" Pulumi module so tokens
		// are npmproxy:index:AccessList (not npmproxy:resources:AccessList).
		ModuleMap: map[tokens.ModuleName]tokens.ModuleName{
			"resources": "index",
		},
	})

	if err := p.RunProvider(context.Background(), "npmproxy", "0.0.1", provider); err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}
}
