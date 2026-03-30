package main

import (
	p "github.com/pulumi/pulumi-go-provider"
	"github.com/pulumi/pulumi-go-provider/infer"

	"github.com/DCCoder90/home-net/pulumi/providers/technitium/resources"
)

func main() {
	err := p.RunProvider("technitium", "0.0.1",
		infer.Provider(infer.Options{
			Resources: []infer.InferredResource{
				infer.Resource[*resources.Record, resources.RecordArgs, resources.RecordState](),
			},
			Config: infer.Config[resources.Config](),
		}),
	)
	if err != nil {
		panic(err)
	}
}
