package resources

import (
	"context"
	"fmt"
	"strconv"

	"github.com/pulumi/pulumi-go-provider/infer"
)

// CertificateArgs are the inputs for an NPM Let's Encrypt certificate.
// Tags must match the pulumi.Map keys used in the thin wrapper (npmproxy.go).
type CertificateArgs struct {
	DomainNames []string `pulumi:"domainNames"`
	Email       string   `pulumi:"email"`
	CFAPIToken  string   `pulumi:"cfApiToken"`
}

// CertificateState is the persisted state for an NPM certificate.
type CertificateState struct {
	CertificateArgs
}

// Certificate implements Create/Read/Delete for an NPM Let's Encrypt cert.
// The resource ID is the numeric NPM certificate ID as a string.
type Certificate struct{}

var _ = (infer.CustomCreate[CertificateArgs, CertificateState])((*Certificate)(nil))
var _ = (infer.CustomDelete[CertificateState])((*Certificate)(nil))
var _ = (infer.CustomRead[CertificateArgs, CertificateState])((*Certificate)(nil))

func (*Certificate) Create(ctx context.Context, req infer.CreateRequest[CertificateArgs]) (infer.CreateResponse[CertificateState], error) {
	args := req.Inputs
	if req.DryRun {
		return infer.CreateResponse[CertificateState]{ID: "", Output: CertificateState{args}}, nil
	}
	c, err := newNPMClient(ctx)
	if err != nil {
		return infer.CreateResponse[CertificateState]{}, err
	}
	result, err := c.CreateCertificate(args.DomainNames, args.Email, args.CFAPIToken)
	if err != nil {
		return infer.CreateResponse[CertificateState]{}, err
	}
	return infer.CreateResponse[CertificateState]{ID: strconv.Itoa(result.ID), Output: CertificateState{args}}, nil
}

func (*Certificate) Read(ctx context.Context, req infer.ReadRequest[CertificateArgs, CertificateState]) (infer.ReadResponse[CertificateArgs, CertificateState], error) {
	id, args := req.ID, req.Inputs
	numID, err := strconv.Atoi(id)
	if err != nil {
		return infer.ReadResponse[CertificateArgs, CertificateState]{ID: id, Inputs: args, State: CertificateState{args}}, fmt.Errorf("invalid certificate ID %q: %w", id, err)
	}
	c, err := newNPMClient(ctx)
	if err != nil {
		return infer.ReadResponse[CertificateArgs, CertificateState]{ID: id, Inputs: args, State: CertificateState{args}}, err
	}
	cert, err := c.GetCertificate(numID)
	if err != nil {
		// Signal to Pulumi that the resource no longer exists so it will recreate it.
		return infer.ReadResponse[CertificateArgs, CertificateState]{ID: "", Inputs: args, State: CertificateState{args}}, nil
	}
	// Preserve credential inputs that the API does not return.
	newArgs := CertificateArgs{
		DomainNames: cert.DomainNames,
		Email:       args.Email,
		CFAPIToken:  args.CFAPIToken,
	}
	return infer.ReadResponse[CertificateArgs, CertificateState]{ID: id, Inputs: newArgs, State: CertificateState{newArgs}}, nil
}

func (*Certificate) Delete(ctx context.Context, req infer.DeleteRequest[CertificateState]) (infer.DeleteResponse, error) {
	id := req.ID
	numID, err := strconv.Atoi(id)
	if err != nil {
		return infer.DeleteResponse{}, fmt.Errorf("invalid certificate ID %q: %w", id, err)
	}
	c, err := newNPMClient(ctx)
	if err != nil {
		return infer.DeleteResponse{}, err
	}
	return infer.DeleteResponse{}, c.DeleteCertificate(numID)
}
