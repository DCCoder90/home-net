package resources

import (
	"github.com/DCCoder90/home-net/pulumi/pkg/config"
	dockerprovider "github.com/pulumi/pulumi-docker/sdk/v4/go/docker"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// RegisterNetworks creates Docker networks declared in stack configs.
// Networks listed in system.yaml's existing_networks are skipped (they pre-exist).
// Returns a map of network name → *dockerprovider.Network.
func RegisterNetworks(
	ctx *pulumi.Context,
	stacks map[string]config.StackConfig,
	existingNetworks []string,
	towerProvider *dockerprovider.Provider,
	importIDs map[string]string,
) (map[string]*dockerprovider.Network, error) {
	existing := make(map[string]bool, len(existingNetworks))
	for _, n := range existingNetworks {
		existing[n] = true
	}

	provOpt := pulumi.Provider(towerProvider)
	networks := map[string]*dockerprovider.Network{}

	for stackName, stack := range stacks {
		for netName, netDef := range stack.Networks {
			if existing[netName] {
				continue
			}
			if _, alreadyCreated := networks[netName]; alreadyCreated {
				continue
			}
			driver := netDef.Driver
			if driver == "" {
				driver = "bridge"
			}
			resourceName := stackName + "-net-" + netName
			net, err := dockerprovider.NewNetwork(ctx, resourceName, &dockerprovider.NetworkArgs{
				Name:       pulumi.String(netName),
				Driver:     pulumi.String(driver),
				Internal:   pulumi.Bool(netDef.Internal),
				Attachable: pulumi.Bool(true),
			}, append([]pulumi.ResourceOption{provOpt}, importOpts(netName, importIDs)...)...)
			if err != nil {
				return nil, err
			}
			networks[netName] = net
		}
	}

	return networks, nil
}
