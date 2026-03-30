package resources

// Config holds provider-level configuration passed via Pulumi stack config.
type Config struct {
	Host     string `pulumi:"host"`
	APIToken string `pulumi:"apiToken" provider:"secret"`
}
