package resources

import (
	"context"
	"fmt"

	authentik "github.com/DCCoder90/home-net/pulumi/providers/authentik/client"
	"github.com/pulumi/pulumi-go-provider/infer"
)

// PolicyBinding manages an Authentik policy binding (group → application).
type PolicyBinding struct{}

// PolicyBindingArgs are the inputs for a policy binding.
type PolicyBindingArgs struct {
	Target string `pulumi:"target"` // application UUID
	Group  string `pulumi:"group"`  // group PK
	Order  int    `pulumi:"order,optional"`
}

// PolicyBindingState is the stored state.
type PolicyBindingState struct {
	PolicyBindingArgs
	BindingPK string `pulumi:"bindingPk"`
}

var _ infer.CustomResource[PolicyBindingArgs, PolicyBindingState] = (*PolicyBinding)(nil)
var _ infer.CustomDelete[PolicyBindingState] = (*PolicyBinding)(nil)

func (pb *PolicyBinding) Create(ctx context.Context, name string, input PolicyBindingArgs, preview bool) (string, PolicyBindingState, error) {
	if preview {
		return name, PolicyBindingState{PolicyBindingArgs: input}, nil
	}

	cfg := infer.GetConfig[Config](ctx)
	c := authentik.New(cfg.Host, cfg.Token)

	binding, err := c.CreatePolicyBinding(authentik.CreatePolicyBindingInput{
		Target: input.Target,
		Group:  input.Group,
		Order:  input.Order,
	})
	if err != nil {
		return "", PolicyBindingState{}, fmt.Errorf("creating Authentik policy binding: %w", err)
	}

	state := PolicyBindingState{PolicyBindingArgs: input, BindingPK: binding.PK}
	return binding.PK, state, nil
}

func (pb *PolicyBinding) Delete(ctx context.Context, id string, props PolicyBindingState) error {
	cfg := infer.GetConfig[Config](ctx)
	c := authentik.New(cfg.Host, cfg.Token)
	return c.DeletePolicyBinding(props.BindingPK)
}
