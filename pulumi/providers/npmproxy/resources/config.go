package resources

// Config holds provider-level configuration passed via Pulumi stack config.
type Config struct {
	URL      string `pulumi:"url"`
	Username string `pulumi:"username"`
	Password string `pulumi:"password" provider:"secret"`
}
