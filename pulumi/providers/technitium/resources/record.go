package resources

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/DCCoder90/home-net/pulumi/providers/technitium/client"
	"github.com/pulumi/pulumi-go-provider/infer"
)

// Record manages a Technitium DNS A or AAAA record.
type Record struct{}

// RecordArgs are the inputs for a DNS record.
type RecordArgs struct {
	Domain    string `pulumi:"domain"`
	Zone      string `pulumi:"zone"`
	Type      string `pulumi:"type"`
	IPAddress string `pulumi:"ipAddress"`
	TTL       int    `pulumi:"ttl,optional"`
}

// RecordState is the stored state after creation.
type RecordState struct {
	RecordArgs
}

var _ infer.CustomResource[RecordArgs, RecordState] = (*Record)(nil)
var _ infer.CustomDelete[RecordState] = (*Record)(nil)
var _ infer.CustomRead[RecordArgs, RecordState] = (*Record)(nil)

func (r *Record) Create(ctx context.Context, name string, input RecordArgs, preview bool) (string, RecordState, error) {
	if preview {
		return recordID(input), RecordState{RecordArgs: input}, nil
	}

	cfg := infer.GetConfig[Config](ctx)

	ttl := input.TTL
	if ttl == 0 {
		ttl = 3600
	}

	c := client.New(cfg.Host, cfg.APIToken)
	if err := c.AddRecord(input.Domain, input.Zone, input.Type, input.IPAddress, ttl); err != nil {
		return "", RecordState{}, fmt.Errorf("creating DNS record: %w", err)
	}

	return recordID(input), RecordState{RecordArgs: input}, nil
}

func (r *Record) Read(ctx context.Context, id string, inputs RecordArgs, state RecordState) (string, RecordArgs, RecordState, error) {
	cfg := infer.GetConfig[Config](ctx)

	c := client.New(cfg.Host, cfg.APIToken)
	exists, ip, err := c.RecordExists(state.Domain, state.Zone, state.Type)
	if err != nil {
		return id, inputs, state, err
	}
	if !exists {
		return "", inputs, RecordState{}, nil
	}

	state.IPAddress = ip
	return id, inputs, state, nil
}

func (r *Record) Delete(ctx context.Context, id string, props RecordState) error {
	cfg := infer.GetConfig[Config](ctx)

	c := client.New(cfg.Host, cfg.APIToken)
	return c.DeleteRecord(props.Domain, props.Zone, props.Type, props.IPAddress)
}

// recordID builds a stable, unique ID from the record's domain and type.
func recordID(a RecordArgs) string {
	return strings.Join([]string{a.Domain, a.Zone, a.Type, strconv.Itoa(a.TTL)}, "|")
}
