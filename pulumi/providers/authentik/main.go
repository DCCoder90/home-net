package main

import (
	p "github.com/pulumi/pulumi-go-provider"
	"github.com/pulumi/pulumi-go-provider/infer"

	"github.com/DCCoder90/home-net/pulumi/providers/authentik/resources"
)

func main() {
	err := p.RunProvider("authentik", "0.0.1",
		infer.Provider(infer.Options{
			Resources: []infer.InferredResource{
				infer.Resource[*resources.Application, resources.ApplicationArgs, resources.ApplicationState](),
				infer.Resource[*resources.ProviderProxy, resources.ProviderProxyArgs, resources.ProviderProxyState](),
				infer.Resource[*resources.ProviderOAuth2, resources.ProviderOAuth2Args, resources.ProviderOAuth2State](),
				infer.Resource[*resources.Group, resources.GroupArgs, resources.GroupState](),
				infer.Resource[*resources.PolicyBinding, resources.PolicyBindingArgs, resources.PolicyBindingState](),
			},
			Config: infer.Config[resources.Config](),
		}),
	)
	if err != nil {
		panic(err)
	}
}
