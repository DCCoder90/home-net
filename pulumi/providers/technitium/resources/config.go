package resources

// Config holds the provider-level configuration for the Technitium provider.
// Field tags must match the property names used in ctx.RegisterResource for the provider.
type Config struct {
	// URL is the base URL of the Technitium DNS server (e.g. "http://192.168.4.53:5380").
	URL string `pulumi:"url"`
	// Token is the Technitium API token.
	Token string `pulumi:"token"`
}
