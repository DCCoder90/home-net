package resources

// Config holds provider-level configuration passed via Pulumi stack config.
type Config struct {
	Host  string `pulumi:"host"`
	Token string `pulumi:"token" provider:"secret"`
}
