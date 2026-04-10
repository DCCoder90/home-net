package npmproxy

// PublicIPPlaceholder is replaced at registration time with the actual public IP
// read from Pulumi config (publicFacingIp).
// This constant is referenced by dns.go to avoid a circular import.
const PublicIPPlaceholder = "0.0.0.0"
