package resources

import (
	"context"
	"fmt"
	"strconv"

	authentik "github.com/DCCoder90/home-net/pulumi/providers/authentik/client"
	"github.com/pulumi/pulumi-go-provider/infer"
)

// ProviderProxy manages an Authentik proxy provider.
type ProviderProxy struct{}

// ProviderProxyArgs are the inputs for a proxy provider.
type ProviderProxyArgs struct {
	Name                       string `pulumi:"name"`
	InternalHost               string `pulumi:"internalHost"`
	ExternalHost               string `pulumi:"externalHost"`
	AuthorizationFlow          string `pulumi:"authorizationFlow"`
	InvalidationFlow           string `pulumi:"invalidationFlow"`
	BasicAuthEnabled           bool   `pulumi:"basicAuthEnabled,optional"`
	BasicAuthUsernameAttribute string `pulumi:"basicAuthUsernameAttribute,optional"`
	BasicAuthPasswordAttribute string `pulumi:"basicAuthPasswordAttribute,optional"`
}

// ProviderProxyState is the stored state.
type ProviderProxyState struct {
	ProviderProxyArgs
	ProviderID int `pulumi:"providerId"`
}

var _ infer.CustomResource[ProviderProxyArgs, ProviderProxyState] = (*ProviderProxy)(nil)
var _ infer.CustomDelete[ProviderProxyState] = (*ProviderProxy)(nil)

func (pp *ProviderProxy) Create(ctx context.Context, name string, input ProviderProxyArgs, preview bool) (string, ProviderProxyState, error) {
	if preview {
		return name, ProviderProxyState{ProviderProxyArgs: input}, nil
	}

	cfg := infer.GetConfig[Config](ctx)
	c := authentik.New(cfg.Host, cfg.Token)

	prov, err := c.CreateProviderProxy(authentik.CreateProviderProxyInput{
		Name:                       input.Name,
		InternalHost:               input.InternalHost,
		ExternalHost:               input.ExternalHost,
		AuthorizationFlow:          input.AuthorizationFlow,
		InvalidationFlow:           input.InvalidationFlow,
		BasicAuthEnabled:           input.BasicAuthEnabled,
		BasicAuthUsernameAttribute: input.BasicAuthUsernameAttribute,
		BasicAuthPasswordAttribute: input.BasicAuthPasswordAttribute,
	})
	if err != nil {
		return "", ProviderProxyState{}, fmt.Errorf("creating Authentik proxy provider: %w", err)
	}

	state := ProviderProxyState{ProviderProxyArgs: input, ProviderID: prov.PK}
	return strconv.Itoa(prov.PK), state, nil
}

func (pp *ProviderProxy) Delete(ctx context.Context, id string, props ProviderProxyState) error {
	cfg := infer.GetConfig[Config](ctx)
	c := authentik.New(cfg.Host, cfg.Token)
	return c.DeleteProviderProxy(props.ProviderID)
}
