package wrappers

import (
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// AccessList is the Pulumi state for an npmproxy:index:AccessList resource.
type AccessList struct {
	pulumi.CustomResourceState
	AccessListID pulumi.IntOutput `pulumi:"accessListId"`
}

// NewAccessList creates an NPM access list via the npmproxy provider.
func NewAccessList(ctx *pulumi.Context, name string, args pulumi.Map, opts ...pulumi.ResourceOption) (*AccessList, error) {
	r := &AccessList{}
	if err := ctx.RegisterResource("npmproxy:index:AccessList", name, args, r, opts...); err != nil {
		return nil, err
	}
	return r, nil
}

// Certificate is the Pulumi state for an npmproxy:index:Certificate resource.
type Certificate struct {
	pulumi.CustomResourceState
	CertificateID pulumi.IntOutput `pulumi:"certificateId"`
}

// NewCertificate requests a Let's Encrypt certificate via the npmproxy provider.
func NewCertificate(ctx *pulumi.Context, name string, args pulumi.Map, opts ...pulumi.ResourceOption) (*Certificate, error) {
	r := &Certificate{}
	if err := ctx.RegisterResource("npmproxy:index:Certificate", name, args, r, opts...); err != nil {
		return nil, err
	}
	return r, nil
}

// ProxyHost is the Pulumi state for an npmproxy:index:ProxyHost resource.
type ProxyHost struct {
	pulumi.CustomResourceState
	ProxyHostID pulumi.IntOutput `pulumi:"proxyHostId"`
}

// NewProxyHost creates an NPM proxy host entry via the npmproxy provider.
func NewProxyHost(ctx *pulumi.Context, name string, args pulumi.Map, opts ...pulumi.ResourceOption) (*ProxyHost, error) {
	r := &ProxyHost{}
	if err := ctx.RegisterResource("npmproxy:index:ProxyHost", name, args, r, opts...); err != nil {
		return nil, err
	}
	return r, nil
}
