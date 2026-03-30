package wrappers

import (
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// AuthentikApplication is the Pulumi state for an authentik:index:Application resource.
type AuthentikApplication struct {
	pulumi.CustomResourceState
	UUID pulumi.StringOutput `pulumi:"uuid"`
}

// NewAuthentikApplication creates an Authentik application.
func NewAuthentikApplication(ctx *pulumi.Context, name string, args pulumi.Map, opts ...pulumi.ResourceOption) (*AuthentikApplication, error) {
	r := &AuthentikApplication{}
	if err := ctx.RegisterResource("authentik:index:Application", name, args, r, opts...); err != nil {
		return nil, err
	}
	return r, nil
}

// AuthentikProviderProxy is the Pulumi state for an authentik:index:ProviderProxy resource.
type AuthentikProviderProxy struct {
	pulumi.CustomResourceState
	ProviderID pulumi.IntOutput `pulumi:"providerId"`
}

// NewAuthentikProviderProxy creates an Authentik proxy provider.
func NewAuthentikProviderProxy(ctx *pulumi.Context, name string, args pulumi.Map, opts ...pulumi.ResourceOption) (*AuthentikProviderProxy, error) {
	r := &AuthentikProviderProxy{}
	if err := ctx.RegisterResource("authentik:index:ProviderProxy", name, args, r, opts...); err != nil {
		return nil, err
	}
	return r, nil
}

// AuthentikProviderOAuth2 is the Pulumi state for an authentik:index:ProviderOAuth2 resource.
type AuthentikProviderOAuth2 struct {
	pulumi.CustomResourceState
	ProviderID   pulumi.IntOutput    `pulumi:"providerId"`
	ClientSecret pulumi.StringOutput `pulumi:"clientSecret"`
	IssuerURL    pulumi.StringOutput `pulumi:"issuerUrl"`
}

// NewAuthentikProviderOAuth2 creates an Authentik OAuth2 provider.
func NewAuthentikProviderOAuth2(ctx *pulumi.Context, name string, args pulumi.Map, opts ...pulumi.ResourceOption) (*AuthentikProviderOAuth2, error) {
	r := &AuthentikProviderOAuth2{}
	if err := ctx.RegisterResource("authentik:index:ProviderOAuth2", name, args, r, opts...); err != nil {
		return nil, err
	}
	return r, nil
}

// AuthentikGroup is the Pulumi state for an authentik:index:Group resource.
type AuthentikGroup struct {
	pulumi.CustomResourceState
	GroupPK pulumi.StringOutput `pulumi:"groupPk"`
}

// NewAuthentikGroup creates an Authentik group.
func NewAuthentikGroup(ctx *pulumi.Context, name string, args pulumi.Map, opts ...pulumi.ResourceOption) (*AuthentikGroup, error) {
	r := &AuthentikGroup{}
	if err := ctx.RegisterResource("authentik:index:Group", name, args, r, opts...); err != nil {
		return nil, err
	}
	return r, nil
}

// AuthentikPolicyBinding is the Pulumi state for an authentik:index:PolicyBinding resource.
type AuthentikPolicyBinding struct {
	pulumi.CustomResourceState
	BindingPK pulumi.StringOutput `pulumi:"bindingPk"`
}

// NewAuthentikPolicyBinding creates an Authentik policy binding (group → application).
func NewAuthentikPolicyBinding(ctx *pulumi.Context, name string, args pulumi.Map, opts ...pulumi.ResourceOption) (*AuthentikPolicyBinding, error) {
	r := &AuthentikPolicyBinding{}
	if err := ctx.RegisterResource("authentik:index:PolicyBinding", name, args, r, opts...); err != nil {
		return nil, err
	}
	return r, nil
}
