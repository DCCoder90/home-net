package resources

import (
	"context"
	"fmt"
	"strconv"

	authentik "github.com/DCCoder90/home-net/pulumi/providers/authentik/client"
	"github.com/pulumi/pulumi-go-provider/infer"
)

// ProviderOAuth2 manages an Authentik OAuth2/OIDC provider.
type ProviderOAuth2 struct{}

// RedirectURI is a single redirect URI with its matching mode.
type RedirectURI struct {
	MatchingMode string `pulumi:"matchingMode"`
	URL          string `pulumi:"url"`
}

// ProviderOAuth2Args are the inputs for an OAuth2 provider.
type ProviderOAuth2Args struct {
	Name                string        `pulumi:"name"`
	ClientID            string        `pulumi:"clientId"`
	AuthorizationFlow   string        `pulumi:"authorizationFlow"`
	InvalidationFlow    string        `pulumi:"invalidationFlow"`
	SigningKey           string        `pulumi:"signingKey,optional"`
	AllowedRedirectURIs []RedirectURI `pulumi:"allowedRedirectUris"`
	PropertyMappings    []string      `pulumi:"propertyMappings,optional"`
}

// ProviderOAuth2State is the stored state.
type ProviderOAuth2State struct {
	ProviderOAuth2Args
	ProviderID   int    `pulumi:"providerId"`
	ClientSecret string `pulumi:"clientSecret"`
	IssuerURL    string `pulumi:"issuerUrl"`
}

var _ infer.CustomResource[ProviderOAuth2Args, ProviderOAuth2State] = (*ProviderOAuth2)(nil)
var _ infer.CustomDelete[ProviderOAuth2State] = (*ProviderOAuth2)(nil)

func (po *ProviderOAuth2) Create(ctx context.Context, name string, input ProviderOAuth2Args, preview bool) (string, ProviderOAuth2State, error) {
	if preview {
		return name, ProviderOAuth2State{ProviderOAuth2Args: input}, nil
	}

	cfg := infer.GetConfig[Config](ctx)
	c := authentik.New(cfg.Host, cfg.Token)

	uris := make([]authentik.AllowedRedirectURI, len(input.AllowedRedirectURIs))
	for i, u := range input.AllowedRedirectURIs {
		uris[i] = authentik.AllowedRedirectURI{
			MatchingMode: u.MatchingMode,
			URL:          u.URL,
		}
	}

	prov, err := c.CreateProviderOAuth2(authentik.CreateProviderOAuth2Input{
		Name:                input.Name,
		ClientID:            input.ClientID,
		AuthorizationFlow:   input.AuthorizationFlow,
		InvalidationFlow:    input.InvalidationFlow,
		SigningKey:          input.SigningKey,
		AllowedRedirectURIs: uris,
		PropertyMappings:    input.PropertyMappings,
	})
	if err != nil {
		return "", ProviderOAuth2State{}, fmt.Errorf("creating Authentik OAuth2 provider: %w", err)
	}

	oidcCfg, err := c.GetProviderOAuth2Config(prov.PK)
	if err != nil {
		return "", ProviderOAuth2State{}, fmt.Errorf("fetching OAuth2 provider config: %w", err)
	}

	state := ProviderOAuth2State{
		ProviderOAuth2Args: input,
		ProviderID:         prov.PK,
		ClientSecret:       prov.ClientSecret,
		IssuerURL:          oidcCfg.Issuer,
	}
	return strconv.Itoa(prov.PK), state, nil
}

func (po *ProviderOAuth2) Delete(ctx context.Context, id string, props ProviderOAuth2State) error {
	cfg := infer.GetConfig[Config](ctx)
	c := authentik.New(cfg.Host, cfg.Token)
	return c.DeleteProviderOAuth2(props.ProviderID)
}
