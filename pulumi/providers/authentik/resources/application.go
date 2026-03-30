package resources

import (
	"context"
	"fmt"

	authentik "github.com/DCCoder90/home-net/pulumi/providers/authentik/client"
	"github.com/pulumi/pulumi-go-provider/infer"
)

// Application manages an Authentik application.
type Application struct{}

// ApplicationArgs are the inputs for an Authentik application.
type ApplicationArgs struct {
	Name            string `pulumi:"name"`
	Slug            string `pulumi:"slug"`
	MetaDescription string `pulumi:"metaDescription,optional"`
	Group           string `pulumi:"group,optional"`
	ProviderID      int    `pulumi:"providerId,optional"`
}

// ApplicationState is the stored state.
type ApplicationState struct {
	ApplicationArgs
	UUID string `pulumi:"uuid"`
}

var _ infer.CustomResource[ApplicationArgs, ApplicationState] = (*Application)(nil)
var _ infer.CustomDelete[ApplicationState] = (*Application)(nil)

func (a *Application) Create(ctx context.Context, name string, input ApplicationArgs, preview bool) (string, ApplicationState, error) {
	if preview {
		return input.Slug, ApplicationState{ApplicationArgs: input}, nil
	}

	cfg := infer.GetConfig[Config](ctx)
	c := authentik.New(cfg.Host, cfg.Token)

	app, err := c.CreateApplication(authentik.CreateApplicationInput{
		Name:             input.Name,
		Slug:             input.Slug,
		MetaDescription:  input.MetaDescription,
		Group:            input.Group,
		ProtocolProvider: input.ProviderID,
	})
	if err != nil {
		return "", ApplicationState{}, fmt.Errorf("creating Authentik application: %w", err)
	}

	state := ApplicationState{ApplicationArgs: input, UUID: app.UUID}
	return app.UUID, state, nil
}

func (a *Application) Delete(ctx context.Context, id string, props ApplicationState) error {
	cfg := infer.GetConfig[Config](ctx)
	c := authentik.New(cfg.Host, cfg.Token)
	return c.DeleteApplication(props.Slug)
}
