// Command pulumi-resource-technitium is a Pulumi resource provider for Technitium DNS.
// It implements a single resource type: technitium:index:Record.
package main

import (
	"context"
	"fmt"
	"os"

	p "github.com/pulumi/pulumi-go-provider"
	"github.com/pulumi/pulumi-go-provider/infer"
	"github.com/pulumi/pulumi/sdk/v3/go/common/tokens"

	"github.com/DCCoder90/home-net/pulumi/providers/technitium/resources"
)

func main() {
	provider := infer.Provider(infer.Options{
		Resources: []infer.InferredResource{
			infer.Resource[*resources.Record, resources.RecordArgs, resources.RecordState](&resources.Record{}),
		},
		Config: infer.Config[*resources.Config](&resources.Config{}),
		// Map the "resources" Go package to the "index" Pulumi module so tokens
		// are technitium:index:Record (not technitium:resources:Record).
		ModuleMap: map[tokens.ModuleName]tokens.ModuleName{
			"resources": "index",
		},
	})

	if err := p.RunProvider(context.Background(), "technitium", "0.0.1", provider); err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}
}
