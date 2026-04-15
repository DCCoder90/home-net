package resources

import (
	"encoding/base64"
	"fmt"
	"path"
	"strings"

	"github.com/DCCoder90/home-net/pulumi/pkg/config"
	"github.com/DCCoder90/home-net/pulumi/pkg/secrets"
	"github.com/pulumi/pulumi-command/sdk/go/command/remote"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// RegisterConfigFiles writes config files to remote hosts for all services that define
// configfiles. Each file's content is fetched from the Infisical /config folder and
// written to the remote host via SSH using the pulumi-command provider.
//
// Files are only re-written when their content changes (tracked via Pulumi triggers).
// Returns a map of service_name → remote.Command resources for use as container
// dependencies, ensuring config files exist before containers start.
func RegisterConfigFiles(
	ctx *pulumi.Context,
	services []config.FlatService,
	configSecrets map[string]string,
	serverAccess map[string]secrets.SSHConfig,
) (map[string][]pulumi.Resource, error) {
	deps := make(map[string][]pulumi.Resource)

	// Build SSH connections once per host rather than per service.
	conns := map[string]remote.ConnectionArgs{}

	for _, svc := range services {
		if len(svc.Def.ConfigFiles) == 0 {
			continue
		}

		hostName := svc.Def.Host
		if hostName == "" {
			hostName = "tower"
		}

		if _, seen := conns[hostName]; !seen {
			access, ok := serverAccess[hostName]
			if !ok {
				return nil, fmt.Errorf("service %q: no server_access entry for host %q", svc.ServiceName, hostName)
			}
			privKey := pulumi.ToSecret(pulumi.String(access.PrivKeyContent)).ApplyT(
				func(v any) (string, error) { return v.(string), nil },
			).(pulumi.StringOutput)
			conns[hostName] = remote.ConnectionArgs{
				Host:       pulumi.String(access.Host),
				User:       pulumi.String(access.User),
				Port:       pulumi.Float64(float64(access.Port)),
				PrivateKey: privKey,
			}
		}
		conn := conns[hostName]

		for i, cf := range svc.Def.ConfigFiles {
			content, ok := configSecrets[cf.Key]
			if !ok {
				return nil, fmt.Errorf("service %q configfile %q: key %q not found in Infisical /config", svc.ServiceName, cf.Path, cf.Key)
			}

			perms := cf.Permissions
			if perms == "" {
				perms = "0600"
			}

			resourceName := fmt.Sprintf("%s-configfile-%d", svc.ServiceName, i)
			createCmd, deleteCmd := buildFileCommands(cf.Path, content, perms)

			// Wrap content as a Pulumi secret so it is encrypted in state and
			// redacted from logs. Used as the trigger so the file is re-written
			// only when the content changes.
			trigger := pulumi.ToSecret(pulumi.String(content)).ApplyT(
				func(v any) (string, error) { return v.(string), nil },
			).(pulumi.StringOutput)

			cmd, err := remote.NewCommand(ctx, resourceName, &remote.CommandArgs{
				Connection: conn,
				Create:     pulumi.String(createCmd),
				Delete:     pulumi.String(deleteCmd),
				Triggers:   pulumi.Array{trigger},
			})
			if err != nil {
				return nil, fmt.Errorf("service %q configfile %q: %w", svc.ServiceName, cf.Path, err)
			}

			deps[svc.ServiceName] = append(deps[svc.ServiceName], cmd)
		}
	}

	return deps, nil
}

// buildFileCommands returns the create and delete shell commands for a remote config file.
// Content is base64-encoded to safely handle arbitrary characters (newlines, quotes, etc.).
func buildFileCommands(filePath, content, permissions string) (create, delete string) {
	dir := path.Dir(filePath)
	encoded := base64.StdEncoding.EncodeToString([]byte(content))

	// Escape the path for use in single-quoted shell arguments.
	escapedPath := strings.ReplaceAll(filePath, "'", "'\\''")
	escapedDir := strings.ReplaceAll(dir, "'", "'\\''")

	create = fmt.Sprintf(
		"mkdir -p '%s' && printf '%%s' '%s' | base64 -d > '%s' && chmod %s '%s'",
		escapedDir, encoded, escapedPath, permissions, escapedPath,
	)
	delete = fmt.Sprintf("rm -f '%s'", escapedPath)
	return
}
