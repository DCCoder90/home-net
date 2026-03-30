package resources

import (
	"context"
	"fmt"

	authentik "github.com/DCCoder90/home-net/pulumi/providers/authentik/client"
	"github.com/pulumi/pulumi-go-provider/infer"
)

// Group manages an Authentik group.
type Group struct{}

// GroupArgs are the inputs for an Authentik group.
type GroupArgs struct {
	Name       string                 `pulumi:"name"`
	Users      []int                  `pulumi:"users,optional"`
	Attributes map[string]interface{} `pulumi:"attributes,optional"`
}

// GroupState is the stored state.
type GroupState struct {
	GroupArgs
	GroupPK string `pulumi:"groupPk"`
}

var _ infer.CustomResource[GroupArgs, GroupState] = (*Group)(nil)
var _ infer.CustomUpdate[GroupArgs, GroupState] = (*Group)(nil)
var _ infer.CustomDelete[GroupState] = (*Group)(nil)

func (g *Group) Create(ctx context.Context, name string, input GroupArgs, preview bool) (string, GroupState, error) {
	if preview {
		return name, GroupState{GroupArgs: input}, nil
	}

	cfg := infer.GetConfig[Config](ctx)
	c := authentik.New(cfg.Host, cfg.Token)

	grp, err := c.CreateGroup(authentik.CreateGroupInput{
		Name:       input.Name,
		Users:      input.Users,
		Attributes: input.Attributes,
	})
	if err != nil {
		return "", GroupState{}, fmt.Errorf("creating Authentik group: %w", err)
	}

	state := GroupState{GroupArgs: input, GroupPK: grp.PK}
	return grp.PK, state, nil
}

func (g *Group) Update(ctx context.Context, id string, olds GroupState, news GroupArgs, preview bool) (GroupState, error) {
	if preview {
		return GroupState{GroupArgs: news, GroupPK: olds.GroupPK}, nil
	}

	cfg := infer.GetConfig[Config](ctx)
	c := authentik.New(cfg.Host, cfg.Token)

	if err := c.UpdateGroup(olds.GroupPK, authentik.CreateGroupInput{
		Name:       news.Name,
		Users:      news.Users,
		Attributes: news.Attributes,
	}); err != nil {
		return GroupState{}, fmt.Errorf("updating Authentik group: %w", err)
	}

	return GroupState{GroupArgs: news, GroupPK: olds.GroupPK}, nil
}

func (g *Group) Delete(ctx context.Context, id string, props GroupState) error {
	cfg := infer.GetConfig[Config](ctx)
	c := authentik.New(cfg.Host, cfg.Token)
	return c.DeleteGroup(props.GroupPK)
}
