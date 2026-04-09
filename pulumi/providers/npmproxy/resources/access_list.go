package resources

import (
	"context"
	"fmt"
	"strconv"

	"github.com/DCCoder90/home-net/pulumi/providers/npmproxy/client"
	"github.com/pulumi/pulumi-go-provider/infer"
)

// AccessListClientEntry mirrors the client package type for the Pulumi schema.
type AccessListClientEntry struct {
	Address   string `pulumi:"address"`
	Directive string `pulumi:"directive"`
}

// AccessListArgs are the inputs for an NPM access list.
type AccessListArgs struct {
	Name       string                  `pulumi:"name"`
	PassAuth   bool                    `pulumi:"passAuth"`
	SatisfyAny bool                    `pulumi:"satisfyAll"` // matches wrapper key "satisfyAll"
	Clients    []AccessListClientEntry `pulumi:"clients"`
}

// AccessListState is the persisted state for an NPM access list.
type AccessListState struct {
	AccessListArgs
}

// AccessList implements CRUD for an NPM access list.
// The resource ID is the numeric NPM access list ID as a string (e.g. "3").
type AccessList struct{}

var _ = (infer.CustomCreate[AccessListArgs, AccessListState])((*AccessList)(nil))
var _ = (infer.CustomDelete[AccessListState])((*AccessList)(nil))
var _ = (infer.CustomRead[AccessListArgs, AccessListState])((*AccessList)(nil))
var _ = (infer.CustomUpdate[AccessListArgs, AccessListState])((*AccessList)(nil))

func newNPMClient(ctx context.Context) (*client.Client, error) {
	cfg := infer.GetConfig[*Config](ctx)
	c := client.New(cfg.URL, cfg.Username, cfg.Password)
	if err := c.Authenticate(); err != nil {
		return nil, fmt.Errorf("npm authenticate: %w", err)
	}
	return c, nil
}

func toClientEntries(in []AccessListClientEntry) []client.AccessListClientEntry {
	out := make([]client.AccessListClientEntry, len(in))
	for i, e := range in {
		out[i] = client.AccessListClientEntry{Address: e.Address, Directive: e.Directive}
	}
	return out
}

func fromClientEntries(in []client.AccessListClientEntry) []AccessListClientEntry {
	out := make([]AccessListClientEntry, len(in))
	for i, e := range in {
		out[i] = AccessListClientEntry{Address: e.Address, Directive: e.Directive}
	}
	return out
}

func (*AccessList) Create(ctx context.Context, req infer.CreateRequest[AccessListArgs]) (infer.CreateResponse[AccessListState], error) {
	args := req.Inputs
	if req.DryRun {
		return infer.CreateResponse[AccessListState]{ID: "", Output: AccessListState{args}}, nil
	}
	c, err := newNPMClient(ctx)
	if err != nil {
		return infer.CreateResponse[AccessListState]{}, err
	}
	result, err := c.CreateAccessList(args.Name, args.PassAuth, args.SatisfyAny, toClientEntries(args.Clients))
	if err != nil {
		return infer.CreateResponse[AccessListState]{}, err
	}
	return infer.CreateResponse[AccessListState]{ID: strconv.Itoa(result.ID), Output: AccessListState{args}}, nil
}

func (*AccessList) Read(ctx context.Context, req infer.ReadRequest[AccessListArgs, AccessListState]) (infer.ReadResponse[AccessListArgs, AccessListState], error) {
	id, args := req.ID, req.Inputs
	numID, err := strconv.Atoi(id)
	if err != nil {
		return infer.ReadResponse[AccessListArgs, AccessListState]{ID: id, Inputs: args, State: AccessListState{args}}, fmt.Errorf("invalid access list ID %q: %w", id, err)
	}
	c, err := newNPMClient(ctx)
	if err != nil {
		return infer.ReadResponse[AccessListArgs, AccessListState]{ID: id, Inputs: args, State: AccessListState{args}}, err
	}
	al, err := c.GetAccessList(numID)
	if err != nil {
		// Signal to Pulumi that the resource no longer exists so it will recreate it.
		return infer.ReadResponse[AccessListArgs, AccessListState]{ID: "", Inputs: args, State: AccessListState{args}}, nil
	}
	newArgs := AccessListArgs{
		Name:       al.Name,
		PassAuth:   bool(al.PassAuth),
		SatisfyAny: bool(al.SatisfyAny),
		Clients:    fromClientEntries(al.Clients),
	}
	return infer.ReadResponse[AccessListArgs, AccessListState]{ID: id, Inputs: newArgs, State: AccessListState{newArgs}}, nil
}

func (*AccessList) Update(ctx context.Context, req infer.UpdateRequest[AccessListArgs, AccessListState]) (infer.UpdateResponse[AccessListState], error) {
	id, args := req.ID, req.Inputs
	if req.DryRun {
		return infer.UpdateResponse[AccessListState]{Output: AccessListState{args}}, nil
	}
	numID, err := strconv.Atoi(id)
	if err != nil {
		return infer.UpdateResponse[AccessListState]{Output: AccessListState{args}}, fmt.Errorf("invalid access list ID %q: %w", id, err)
	}
	c, err := newNPMClient(ctx)
	if err != nil {
		return infer.UpdateResponse[AccessListState]{Output: AccessListState{args}}, err
	}
	_, err = c.UpdateAccessList(numID, args.Name, args.PassAuth, args.SatisfyAny, toClientEntries(args.Clients))
	if err != nil {
		return infer.UpdateResponse[AccessListState]{Output: AccessListState{args}}, err
	}
	return infer.UpdateResponse[AccessListState]{Output: AccessListState{args}}, nil
}

func (*AccessList) Delete(ctx context.Context, req infer.DeleteRequest[AccessListState]) (infer.DeleteResponse, error) {
	id := req.ID
	numID, err := strconv.Atoi(id)
	if err != nil {
		return infer.DeleteResponse{}, fmt.Errorf("invalid access list ID %q: %w", id, err)
	}
	c, err := newNPMClient(ctx)
	if err != nil {
		return infer.DeleteResponse{}, err
	}
	return infer.DeleteResponse{}, c.DeleteAccessList(numID)
}
