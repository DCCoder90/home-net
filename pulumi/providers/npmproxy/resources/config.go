package resources

// Config holds the provider-level configuration for the NPM Proxy Manager provider.
// Field tags must match the property names used in ctx.RegisterResource for the provider.
type Config struct {
	// URL is the base URL of the NPM instance (e.g. "http://192.168.4.2:81").
	URL string `pulumi:"url"`
	// Username is the NPM admin username.
	Username string `pulumi:"username"`
	// Password is the NPM admin password.
	Password string `pulumi:"password"`
}
