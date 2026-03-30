package resources

import (
	"context"
	"fmt"
	"strconv"

	npmclient "github.com/DCCoder90/home-net/pulumi/providers/npmproxy/client"
	"github.com/pulumi/pulumi-go-provider/infer"
)

// AccessList manages an NPM access list.
type AccessList struct{}

// AccessItem is a single allow/deny rule.
type AccessItem struct {
	Address   string `pulumi:"address"`
	Directive string `pulumi:"directive"`
}

// AccessListArgs are the inputs for an NPM access list.
type AccessListArgs struct {
	Name       string       `pulumi:"name"`
	Items      []AccessItem `pulumi:"items"`
	PassAuth   bool         `pulumi:"passAuth,optional"`
	SatisfyAny bool         `pulumi:"satisfyAny,optional"`
}

// AccessListState is the stored state.
type AccessListState struct {
	AccessListArgs
	AccessListID int `pulumi:"accessListId"`
}

var _ infer.CustomResource[AccessListArgs, AccessListState] = (*AccessList)(nil)
var _ infer.CustomDelete[AccessListState] = (*AccessList)(nil)

func (a *AccessList) Create(ctx context.Context, name string, input AccessListArgs, preview bool) (string, AccessListState, error) {
	if preview {
		return name, AccessListState{AccessListArgs: input}, nil
	}

	cfg := infer.GetConfig[Config](ctx)

	c := npmclient.New(cfg.URL, cfg.Username, cfg.Password)
	if err := c.Authenticate(); err != nil {
		return "", AccessListState{}, fmt.Errorf("NPM auth: %w", err)
	}

	items := make([]npmclient.AccessListItem, len(input.Items))
	for i, item := range input.Items {
		items[i] = npmclient.AccessListItem{Address: item.Address, Directive: item.Directive}
	}

	id, err := c.CreateAccessList(npmclient.CreateAccessListInput{
		Name:       input.Name,
		Items:      items,
		PassAuth:   input.PassAuth,
		SatisfyAny: input.SatisfyAny,
	})
	if err != nil {
		return "", AccessListState{}, fmt.Errorf("creating NPM access list: %w", err)
	}

	state := AccessListState{AccessListArgs: input, AccessListID: id}
	return strconv.Itoa(id), state, nil
}

func (a *AccessList) Delete(ctx context.Context, id string, props AccessListState) error {
	cfg := infer.GetConfig[Config](ctx)

	c := npmclient.New(cfg.URL, cfg.Username, cfg.Password)
	if err := c.Authenticate(); err != nil {
		return err
	}

	return c.DeleteAccessList(props.AccessListID)
}
