package resources

import "github.com/pulumi/pulumi/sdk/v3/go/pulumi"

// importOpts returns a pulumi.Import resource option if the given resource name
// has a non-empty entry in the importIDs map. Used when PULUMI_IMPORT_IDS_FILE is set.
func importOpts(resourceName string, importIDs map[string]string) []pulumi.ResourceOption {
	if importIDs == nil {
		return nil
	}
	id, ok := importIDs[resourceName]
	if !ok || id == "" {
		return nil
	}
	return []pulumi.ResourceOption{pulumi.Import(pulumi.ID(id))}
}
