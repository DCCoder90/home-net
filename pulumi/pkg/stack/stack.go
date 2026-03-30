package stack

import (
	"github.com/DCCoder90/home-net/pulumi/pkg/npmconfig"
	"github.com/DCCoder90/home-net/pulumi/pkg/servers"
	"github.com/DCCoder90/home-net/pulumi/pkg/service"
	"github.com/DCCoder90/home-net/pulumi/pkg/types"
	dockerprovider "github.com/pulumi/pulumi-docker/sdk/v4/go/docker"
	"github.com/pulumi/pulumi-random/sdk/v4/go/random"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// DeployInput is everything needed to deploy a full stack.
type DeployInput struct {
	Name           string
	Config         *types.StackConfig
	Servers        map[string]*servers.ServerContext
	System         *types.SystemConfig
	AllSecrets     map[string]string
	NPMAccessLists *npmconfig.AccessLists
	PublicIP       string
	CFApiToken     string
	AdminEmail     string
	AdminUsername  string
	Parent         pulumi.Resource
}

// Deploy creates custom Docker networks, then deploys each service in the stack.
func Deploy(ctx *pulumi.Context, in *DeployInput) error {
	// Determine which server hosts this stack (default: tower).
	hostName := servers.ServiceHost(in.Config.Host)
	srv, err := servers.Require(in.Servers, hostName)
	if err != nil {
		return err
	}

	// Generate stable random passwords for each declared generated_secret.
	// Each password is stored in Pulumi state and stable across runs.
	genSecrets := map[string]pulumi.StringOutput{}
	for _, key := range in.Config.GeneratedSecrets {
		if key == "" {
			continue
		}
		pwd, err := random.NewRandomPassword(ctx, in.Name+"-secret-"+key, &random.RandomPasswordArgs{
			Length:  pulumi.Int(32),
			Special: pulumi.Bool(false),
		}, pulumi.Parent(in.Parent))
		if err != nil {
			return err
		}
		genSecrets[key] = pwd.Result
	}

	// Create custom Docker networks defined in the stack.
	for netName, netDef := range in.Config.Networks {
		driver := netDef.Driver
		if driver == "" {
			driver = "bridge"
		}
		if _, err := dockerprovider.NewNetwork(ctx, in.Name+"-net-"+netName, &dockerprovider.NetworkArgs{
			Name:     pulumi.String(netName),
			Driver:   pulumi.String(driver),
			Internal: pulumi.Bool(netDef.Internal),
		}, pulumi.Provider(srv.Provider), pulumi.Parent(in.Parent)); err != nil {
			return err
		}
	}

	// Deploy each service, merging stack-level env and mounts.
	for svcKey, svcConfig := range in.Config.Services {
		cfg := svcConfig // capture loop variable
		// Per-service host override is possible (service.Host takes precedence over stack.Host).
		svcHost := servers.ServiceHost(cfg.Host)
		if cfg.Host == "" {
			svcHost = hostName
		}
		svcSrv, err := servers.Require(in.Servers, svcHost)
		if err != nil {
			return err
		}

		if err := service.Deploy(ctx, &service.DeployInput{
			Name:             in.Name + "-" + svcKey,
			Config:           &cfg,
			Server:           svcSrv,
			System:           in.System,
			AllSecrets:       in.AllSecrets,
			NPMAccessLists:   in.NPMAccessLists,
			PublicIP:         in.PublicIP,
			CFApiToken:       in.CFApiToken,
			AdminEmail:       in.AdminEmail,
			AdminUsername:    in.AdminUsername,
			StackEnvs:        in.Config.Env,
			StackMounts:      in.Config.Mounts,
			GeneratedSecrets: genSecrets,
			Parent:           in.Parent,
		}); err != nil {
			return err
		}
	}

	return nil
}
