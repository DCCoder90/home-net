package resources

import (
	"context"
	"fmt"
	"strconv"

	"github.com/DCCoder90/home-net/pulumi/providers/npmproxy/client"
	"github.com/pulumi/pulumi-go-provider/infer"
)

// ProxyHostArgs are the inputs for an NPM proxy host.
// Tags must match the pulumi.Map keys used in the thin wrapper (npmproxy.go).
type ProxyHostArgs struct {
	DomainNames           []string `pulumi:"domainNames"`
	ForwardHost           string   `pulumi:"forwardHost"`
	ForwardPort           int      `pulumi:"forwardPort"`
	ForwardScheme         string   `pulumi:"forwardScheme"`
	CertificateID         string   `pulumi:"certificateId"`
	AccessListID          string   `pulumi:"accessListId"`
	SSLForced             bool     `pulumi:"sslForced"`
	HTTP2Support          bool     `pulumi:"http2Support"`
	BlockExploits         bool     `pulumi:"blockExploits"`
	AllowWebsocketUpgrade bool     `pulumi:"allowWebsocketUpgrade"`
	AdvancedConfig        string   `pulumi:"advancedConfig,optional"`
}

// ProxyHostState is the persisted state for an NPM proxy host.
type ProxyHostState struct {
	ProxyHostArgs
}

// ProxyHost implements CRUD for an NPM proxy host entry.
// The resource ID is the numeric NPM proxy host ID as a string.
type ProxyHost struct{}

var _ = (infer.CustomCreate[ProxyHostArgs, ProxyHostState])((*ProxyHost)(nil))
var _ = (infer.CustomDelete[ProxyHostState])((*ProxyHost)(nil))
var _ = (infer.CustomRead[ProxyHostArgs, ProxyHostState])((*ProxyHost)(nil))
var _ = (infer.CustomUpdate[ProxyHostArgs, ProxyHostState])((*ProxyHost)(nil))

func (*ProxyHost) Create(ctx context.Context, req infer.CreateRequest[ProxyHostArgs]) (infer.CreateResponse[ProxyHostState], error) {
	args := req.Inputs
	if req.DryRun {
		return infer.CreateResponse[ProxyHostState]{ID: "", Output: ProxyHostState{args}}, nil
	}
	c, err := newNPMClient(ctx)
	if err != nil {
		return infer.CreateResponse[ProxyHostState]{}, err
	}
	certID, err := strconv.Atoi(args.CertificateID)
	if err != nil {
		return infer.CreateResponse[ProxyHostState]{}, fmt.Errorf("invalid certificateId %q: %w", args.CertificateID, err)
	}
	ph := client.ProxyHost{
		DomainNames:           args.DomainNames,
		ForwardHost:           args.ForwardHost,
		ForwardPort:           args.ForwardPort,
		ForwardScheme:         args.ForwardScheme,
		CertificateID:         certID,
		SSLForced:             args.SSLForced,
		HTTP2Support:          args.HTTP2Support,
		BlockExploits:         args.BlockExploits,
		AllowWebsocketUpgrade: args.AllowWebsocketUpgrade,
		AdvancedConfig:        args.AdvancedConfig,
	}
	result, err := c.CreateProxyHost(ph)
	if err != nil {
		return infer.CreateResponse[ProxyHostState]{}, err
	}
	return infer.CreateResponse[ProxyHostState]{ID: strconv.Itoa(result.ID), Output: ProxyHostState{args}}, nil
}

func (*ProxyHost) Read(ctx context.Context, req infer.ReadRequest[ProxyHostArgs, ProxyHostState]) (infer.ReadResponse[ProxyHostArgs, ProxyHostState], error) {
	id, args := req.ID, req.Inputs
	numID, err := strconv.Atoi(id)
	if err != nil {
		return infer.ReadResponse[ProxyHostArgs, ProxyHostState]{ID: id, Inputs: args, State: ProxyHostState{args}}, fmt.Errorf("invalid proxy host ID %q: %w", id, err)
	}
	c, err := newNPMClient(ctx)
	if err != nil {
		return infer.ReadResponse[ProxyHostArgs, ProxyHostState]{ID: id, Inputs: args, State: ProxyHostState{args}}, err
	}
	ph, err := c.GetProxyHost(numID)
	if err != nil {
		// Signal to Pulumi that the resource no longer exists so it will recreate it.
		return infer.ReadResponse[ProxyHostArgs, ProxyHostState]{ID: "", Inputs: args, State: ProxyHostState{args}}, nil
	}
	newArgs := ProxyHostArgs{
		DomainNames:           ph.DomainNames,
		ForwardHost:           ph.ForwardHost,
		ForwardPort:           ph.ForwardPort,
		ForwardScheme:         ph.ForwardScheme,
		CertificateID:         strconv.Itoa(ph.CertificateID),
		SSLForced:             ph.SSLForced,
		HTTP2Support:          ph.HTTP2Support,
		BlockExploits:         ph.BlockExploits,
		AllowWebsocketUpgrade: ph.AllowWebsocketUpgrade,
		AdvancedConfig:        ph.AdvancedConfig,
	}
	return infer.ReadResponse[ProxyHostArgs, ProxyHostState]{ID: id, Inputs: newArgs, State: ProxyHostState{newArgs}}, nil
}

func (*ProxyHost) Update(ctx context.Context, req infer.UpdateRequest[ProxyHostArgs, ProxyHostState]) (infer.UpdateResponse[ProxyHostState], error) {
	id, args := req.ID, req.Inputs
	if req.DryRun {
		return infer.UpdateResponse[ProxyHostState]{Output: ProxyHostState{args}}, nil
	}
	numID, err := strconv.Atoi(id)
	if err != nil {
		return infer.UpdateResponse[ProxyHostState]{Output: ProxyHostState{args}}, fmt.Errorf("invalid proxy host ID %q: %w", id, err)
	}
	c, err := newNPMClient(ctx)
	if err != nil {
		return infer.UpdateResponse[ProxyHostState]{Output: ProxyHostState{args}}, err
	}
	certID, err := strconv.Atoi(args.CertificateID)
	if err != nil {
		return infer.UpdateResponse[ProxyHostState]{Output: ProxyHostState{args}}, fmt.Errorf("invalid certificateId: %w", err)
	}
	ph := client.ProxyHost{
		DomainNames:           args.DomainNames,
		ForwardHost:           args.ForwardHost,
		ForwardPort:           args.ForwardPort,
		ForwardScheme:         args.ForwardScheme,
		CertificateID:         certID,
		SSLForced:             args.SSLForced,
		HTTP2Support:          args.HTTP2Support,
		BlockExploits:         args.BlockExploits,
		AllowWebsocketUpgrade: args.AllowWebsocketUpgrade,
		AdvancedConfig:        args.AdvancedConfig,
	}
	_, err = c.UpdateProxyHost(numID, ph)
	if err != nil {
		return infer.UpdateResponse[ProxyHostState]{Output: ProxyHostState{args}}, err
	}
	return infer.UpdateResponse[ProxyHostState]{Output: ProxyHostState{args}}, nil
}

func (*ProxyHost) Delete(ctx context.Context, req infer.DeleteRequest[ProxyHostState]) (infer.DeleteResponse, error) {
	id := req.ID
	numID, err := strconv.Atoi(id)
	if err != nil {
		return infer.DeleteResponse{}, fmt.Errorf("invalid proxy host ID %q: %w", id, err)
	}
	c, err := newNPMClient(ctx)
	if err != nil {
		return infer.DeleteResponse{}, err
	}
	return infer.DeleteResponse{}, c.DeleteProxyHost(numID)
}
