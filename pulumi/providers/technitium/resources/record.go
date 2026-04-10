package resources

import (
	"context"
	"fmt"
	"strings"

	"github.com/DCCoder90/home-net/pulumi/providers/technitium/client"
	"github.com/pulumi/pulumi-go-provider/infer"
)

// RecordArgs are the inputs for a Technitium DNS record.
// Tags must match the pulumi.Map keys used in the thin wrapper (technitium.go).
type RecordArgs struct {
	Zone      string `pulumi:"zone"`
	Name      string `pulumi:"name"`
	Type      string `pulumi:"type"`
	TTL       int    `pulumi:"ttl"`
	IPAddress string `pulumi:"ipAddress"`
}

// RecordState is the persisted state for a Technitium DNS record.
type RecordState struct {
	RecordArgs
}

// Record implements CRUD for a Technitium DNS A/AAAA record.
// The composite resource ID is "domain|zone|type|ttl" (e.g. "grafana.dcapi.app|dcapi.app|A|3600").
type Record struct{}

// Ensure Record satisfies the infer interfaces at compile time.
var _ = (infer.CustomCreate[RecordArgs, RecordState])((*Record)(nil))
var _ = (infer.CustomDelete[RecordState])((*Record)(nil))
var _ = (infer.CustomRead[RecordArgs, RecordState])((*Record)(nil))

func newClient(ctx context.Context) (*client.Client, error) {
	cfg := infer.GetConfig[*Config](ctx)
	return client.New(cfg.URL, cfg.Token), nil
}

// recordID builds the stable composite ID for a record (domain|zone|type|ttl).
func recordID(domain, zone, recType string, ttl int) string {
	return fmt.Sprintf("%s|%s|%s|%d", domain, zone, recType, ttl)
}

// parseID splits the composite ID back into its parts.
// Accepts both 4-part "domain|zone|type|ttl" (current) and legacy 3-part "domain|zone|type".
func parseID(id string) (domain, zone, recType string, ttl int, err error) {
	parts := strings.Split(id, "|")
	switch len(parts) {
	case 4:
		var n int
		if _, e := fmt.Sscanf(parts[3], "%d", &n); e != nil {
			return "", "", "", 0, fmt.Errorf("invalid ttl in ID %q: %w", id, e)
		}
		return parts[0], parts[1], parts[2], n, nil
	case 3:
		return parts[0], parts[1], parts[2], 3600, nil
	default:
		return "", "", "", 0, fmt.Errorf("invalid record ID %q: expected domain|zone|type|ttl", id)
	}
}

// Create adds a new DNS record and returns the composite ID.
func (*Record) Create(ctx context.Context, req infer.CreateRequest[RecordArgs]) (infer.CreateResponse[RecordState], error) {
	args := req.Inputs
	ttl := args.TTL
	if ttl <= 0 {
		ttl = 3600
	}
	if req.DryRun {
		return infer.CreateResponse[RecordState]{ID: recordID(args.Name, args.Zone, args.Type, ttl), Output: RecordState{args}}, nil
	}
	c, err := newClient(ctx)
	if err != nil {
		return infer.CreateResponse[RecordState]{}, err
	}
	if err := c.AddRecord(args.Zone, args.Name, args.Type, ttl, args.IPAddress); err != nil {
		return infer.CreateResponse[RecordState]{}, fmt.Errorf("add record: %w", err)
	}
	return infer.CreateResponse[RecordState]{ID: recordID(args.Name, args.Zone, args.Type, ttl), Output: RecordState{args}}, nil
}

// Read fetches the current state of the record (used during refresh and import).
func (*Record) Read(ctx context.Context, req infer.ReadRequest[RecordArgs, RecordState]) (infer.ReadResponse[RecordArgs, RecordState], error) {
	id, args := req.ID, req.Inputs
	domain, zone, recType, _, err := parseID(id)
	if err != nil {
		return infer.ReadResponse[RecordArgs, RecordState]{ID: id, Inputs: args, State: RecordState{args}}, err
	}
	c, err := newClient(ctx)
	if err != nil {
		return infer.ReadResponse[RecordArgs, RecordState]{ID: id, Inputs: args, State: RecordState{args}}, err
	}
	info, err := c.GetRecord(zone, domain, recType)
	if err != nil {
		return infer.ReadResponse[RecordArgs, RecordState]{ID: id, Inputs: args, State: RecordState{args}}, fmt.Errorf("get record: %w", err)
	}
	if info == nil {
		// Record does not exist — Pulumi will treat this as deleted.
		return infer.ReadResponse[RecordArgs, RecordState]{}, nil
	}
	ttl := info.TTL
	if ttl == 0 {
		ttl = args.TTL
		if ttl == 0 {
			ttl = 3600
		}
	}
	newArgs := RecordArgs{
		Zone:      zone,
		Name:      domain,
		Type:      recType,
		TTL:       ttl,
		IPAddress: info.IPAddress,
	}
	return infer.ReadResponse[RecordArgs, RecordState]{ID: recordID(domain, zone, recType, ttl), Inputs: newArgs, State: RecordState{newArgs}}, nil
}

// Delete removes the DNS record.
func (*Record) Delete(ctx context.Context, req infer.DeleteRequest[RecordState]) (infer.DeleteResponse, error) {
	id := req.ID
	state := req.State
	domain, zone, recType, _, err := parseID(id)
	if err != nil {
		return infer.DeleteResponse{}, err
	}
	c, err := newClient(ctx)
	if err != nil {
		return infer.DeleteResponse{}, err
	}
	return infer.DeleteResponse{}, c.DeleteRecord(zone, domain, recType, state.IPAddress)
}
