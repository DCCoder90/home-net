package resources

import (
	"context"
	"fmt"
	"strconv"

	npmclient "github.com/DCCoder90/home-net/pulumi/providers/npmproxy/client"
	"github.com/pulumi/pulumi-go-provider/infer"
)

// ProxyHost manages an NPM reverse-proxy entry.
type ProxyHost struct{}

// ProxyHostArgs are the inputs for an NPM proxy host.
type ProxyHostArgs struct {
	DomainName            string `pulumi:"domainName"`
	ForwardScheme         string `pulumi:"forwardScheme,optional"`
	ForwardHost           string `pulumi:"forwardHost"`
	ForwardPort           int    `pulumi:"forwardPort"`
	AccessListID          int    `pulumi:"accessListId,optional"`
	CertificateID         int    `pulumi:"certificateId,optional"`
	CachingEnabled        bool   `pulumi:"cachingEnabled,optional"`
	AllowWebsocketUpgrade bool   `pulumi:"allowWebsocketUpgrade,optional"`
	BlockExploits         bool   `pulumi:"blockExploits,optional"`
	HTTP2Support          bool   `pulumi:"http2Support,optional"`
	SSLForced             bool   `pulumi:"sslForced,optional"`
	HSTSEnabled           bool   `pulumi:"hstsEnabled,optional"`
	HSTSSubdomains        bool   `pulumi:"hstsSubdomains,optional"`
}

// ProxyHostState is the stored state.
type ProxyHostState struct {
	ProxyHostArgs
	ProxyHostID int `pulumi:"proxyHostId"`
}

var _ infer.CustomResource[ProxyHostArgs, ProxyHostState] = (*ProxyHost)(nil)
var _ infer.CustomUpdate[ProxyHostArgs, ProxyHostState] = (*ProxyHost)(nil)
var _ infer.CustomDelete[ProxyHostState] = (*ProxyHost)(nil)

func (ph *ProxyHost) Create(ctx context.Context, name string, input ProxyHostArgs, preview bool) (string, ProxyHostState, error) {
	if preview {
		return name, ProxyHostState{ProxyHostArgs: input}, nil
	}

	cfg := infer.GetConfig[Config](ctx)

	c := npmclient.New(cfg.URL, cfg.Username, cfg.Password)
	if err := c.Authenticate(); err != nil {
		return "", ProxyHostState{}, fmt.Errorf("NPM auth: %w", err)
	}

	scheme := input.ForwardScheme
	if scheme == "" {
		scheme = "http"
	}

	id, err := c.CreateProxyHost(npmclient.ProxyHostInput{
		DomainNames:           []string{input.DomainName},
		ForwardScheme:         scheme,
		ForwardHost:           input.ForwardHost,
		ForwardPort:           input.ForwardPort,
		AccessListID:          input.AccessListID,
		CertificateID:         input.CertificateID,
		CachingEnabled:        input.CachingEnabled,
		AllowWebsocketUpgrade: input.AllowWebsocketUpgrade,
		BlockExploits:         input.BlockExploits,
		HTTP2Support:          input.HTTP2Support,
		SSLForced:             input.SSLForced,
		HSTSEnabled:           input.HSTSEnabled,
		HSTSSubdomains:        input.HSTSSubdomains,
		Enabled:               true,
	})
	if err != nil {
		return "", ProxyHostState{}, fmt.Errorf("creating NPM proxy host: %w", err)
	}

	state := ProxyHostState{ProxyHostArgs: input, ProxyHostID: id}
	return strconv.Itoa(id), state, nil
}

func (ph *ProxyHost) Update(ctx context.Context, id string, olds ProxyHostState, news ProxyHostArgs, preview bool) (ProxyHostState, error) {
	if preview {
		return ProxyHostState{ProxyHostArgs: news, ProxyHostID: olds.ProxyHostID}, nil
	}

	cfg := infer.GetConfig[Config](ctx)

	c := npmclient.New(cfg.URL, cfg.Username, cfg.Password)
	if err := c.Authenticate(); err != nil {
		return ProxyHostState{}, err
	}

	scheme := news.ForwardScheme
	if scheme == "" {
		scheme = "http"
	}

	if err := c.UpdateProxyHost(olds.ProxyHostID, npmclient.ProxyHostInput{
		DomainNames:           []string{news.DomainName},
		ForwardScheme:         scheme,
		ForwardHost:           news.ForwardHost,
		ForwardPort:           news.ForwardPort,
		AccessListID:          news.AccessListID,
		CertificateID:         news.CertificateID,
		CachingEnabled:        news.CachingEnabled,
		AllowWebsocketUpgrade: news.AllowWebsocketUpgrade,
		BlockExploits:         news.BlockExploits,
		HTTP2Support:          news.HTTP2Support,
		SSLForced:             news.SSLForced,
		HSTSEnabled:           news.HSTSEnabled,
		HSTSSubdomains:        news.HSTSSubdomains,
		Enabled:               true,
	}); err != nil {
		return ProxyHostState{}, err
	}

	return ProxyHostState{ProxyHostArgs: news, ProxyHostID: olds.ProxyHostID}, nil
}

func (ph *ProxyHost) Delete(ctx context.Context, id string, props ProxyHostState) error {
	cfg := infer.GetConfig[Config](ctx)

	c := npmclient.New(cfg.URL, cfg.Username, cfg.Password)
	if err := c.Authenticate(); err != nil {
		return err
	}

	return c.DeleteProxyHost(props.ProxyHostID)
}
