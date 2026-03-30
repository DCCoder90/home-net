package resources

import (
	"context"
	"fmt"
	"strconv"

	npmclient "github.com/DCCoder90/home-net/pulumi/providers/npmproxy/client"
	"github.com/pulumi/pulumi-go-provider/infer"
)

// Certificate manages an NPM Let's Encrypt certificate.
type Certificate struct{}

// CertificateArgs are the inputs for a Let's Encrypt certificate.
type CertificateArgs struct {
	DomainName         string `pulumi:"domainName"`
	LetsEncryptEmail   string `pulumi:"letsEncryptEmail"`
	CloudflareAPIToken string `pulumi:"cloudflareApiToken" provider:"secret"`
	PropagationSeconds int    `pulumi:"propagationSeconds,optional"`
}

// CertificateState is the stored state.
type CertificateState struct {
	CertificateArgs
	CertificateID int `pulumi:"certificateId"`
}

var _ infer.CustomResource[CertificateArgs, CertificateState] = (*Certificate)(nil)
var _ infer.CustomDelete[CertificateState] = (*Certificate)(nil)

func (cert *Certificate) Create(ctx context.Context, name string, input CertificateArgs, preview bool) (string, CertificateState, error) {
	if preview {
		return name, CertificateState{CertificateArgs: input}, nil
	}

	cfg := infer.GetConfig[Config](ctx)

	c := npmclient.New(cfg.URL, cfg.Username, cfg.Password)
	if err := c.Authenticate(); err != nil {
		return "", CertificateState{}, fmt.Errorf("NPM auth: %w", err)
	}

	propagation := input.PropagationSeconds
	if propagation == 0 {
		propagation = 10
	}

	id, err := c.CreateCertificate(npmclient.CertificateInput{
		Provider:    "letsencrypt",
		DomainNames: []string{input.DomainName},
		Meta: map[string]interface{}{
			"letsencrypt_email":        input.LetsEncryptEmail,
			"letsencrypt_agree":        true,
			"dns_challenge":            true,
			"dns_provider":             "cloudflare",
			"dns_provider_credentials": fmt.Sprintf("dns_cloudflare_api_token=%s", input.CloudflareAPIToken),
			"propagation_seconds":      strconv.Itoa(propagation),
		},
	})
	if err != nil {
		return "", CertificateState{}, fmt.Errorf("creating NPM certificate: %w", err)
	}

	state := CertificateState{CertificateArgs: input, CertificateID: id}
	return strconv.Itoa(id), state, nil
}

func (cert *Certificate) Delete(ctx context.Context, id string, props CertificateState) error {
	cfg := infer.GetConfig[Config](ctx)

	c := npmclient.New(cfg.URL, cfg.Username, cfg.Password)
	if err := c.Authenticate(); err != nil {
		return err
	}

	return c.DeleteCertificate(props.CertificateID)
}
