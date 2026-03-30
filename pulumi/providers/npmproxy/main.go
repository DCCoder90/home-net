package main

import (
	p "github.com/pulumi/pulumi-go-provider"
	"github.com/pulumi/pulumi-go-provider/infer"

	"github.com/DCCoder90/home-net/pulumi/providers/npmproxy/resources"
)

func main() {
	err := p.RunProvider("npmproxy", "0.0.1",
		infer.Provider(infer.Options{
			Resources: []infer.InferredResource{
				infer.Resource[*resources.AccessList, resources.AccessListArgs, resources.AccessListState](),
				infer.Resource[*resources.Certificate, resources.CertificateArgs, resources.CertificateState](),
				infer.Resource[*resources.ProxyHost, resources.ProxyHostArgs, resources.ProxyHostState](),
			},
			Config: infer.Config[resources.Config](),
		}),
	)
	if err != nil {
		panic(err)
	}
}
