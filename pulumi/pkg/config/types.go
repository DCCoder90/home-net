package config

// SystemConfig maps to config/system.yaml.
type SystemConfig struct {
	Authentik struct {
		IPAddress      string `yaml:"ip_address"`
		Port           int    `yaml:"port"`
		DomainName     string `yaml:"domain_name"`
		Image          string `yaml:"image"`
		DataPath       string `yaml:"data_path"`
		PostgreSQLHost string `yaml:"postgresql_host"`
		PostgreSQLUser string `yaml:"postgresql_user"`
		PostgreSQLName string `yaml:"postgresql_name"`
		RedisHost      string `yaml:"redis_host"`
		// Dedicated proxy outpost — optional. When OutpostIPAddress is set, Pulumi
		// deploys a standalone ghcr.io/goauthentik/proxy container and NPM routes
		// proxy-auth traffic through it instead of the embedded Authentik outpost.
		OutpostIPAddress string `yaml:"outpost_ip_address"`
		OutpostPort      int    `yaml:"outpost_port"`  // defaults to 9000
		OutpostImage     string `yaml:"outpost_image"` // defaults to ghcr.io/goauthentik/proxy:latest
	} `yaml:"authentik"`
	ProxyIP          string            `yaml:"proxy_ip"`
	DNSIP            string            `yaml:"dns_ip"`
	Zones            map[string]string `yaml:"zones"`
	AdminUsername    string            `yaml:"admin_username"`
	ExistingNetworks []string          `yaml:"existing_networks"`
	Infisical        struct {
		Project     string `yaml:"project"`
		Environment string `yaml:"environment"`
		Folder      string `yaml:"folder"`
	} `yaml:"infisical"`
}

// ServersFile is the top-level structure of config/servers.yaml.
type ServersFile struct {
	Servers map[string]ServerConfig `yaml:"servers"`
}

// ServerConfig describes a single managed host.
type ServerConfig struct {
	Networks NetworkNames `yaml:"networks"`
}

// NetworkNames holds the primary (services) and secondary network interface names.
type NetworkNames struct {
	Primary   string `yaml:"primary"`
	Secondary string `yaml:"secondary,omitempty"`
}

// StackConfig is the top-level structure of a config/stacks/*.yaml entry.
type StackConfig struct {
	Host             string                   `yaml:"host,omitempty"`
	Env              []string                 `yaml:"env,omitempty"`
	Mounts           []string                 `yaml:"mounts,omitempty"`
	Volumes          []string                 `yaml:"volumes,omitempty"`
	Networks         map[string]NetworkDef    `yaml:"networks,omitempty"`
	GeneratedSecrets []string                 `yaml:"generated_secrets,omitempty"`
	Services         map[string]ServiceConfig `yaml:"services"`
}

// NetworkDef describes a Docker network to create for a stack.
type NetworkDef struct {
	Internal bool   `yaml:"internal,omitempty"`
	Driver   string `yaml:"driver,omitempty"`
}

// ConfigFile defines a file to write to the remote host from an Infisical /config secret.
type ConfigFile struct {
	Path        string `yaml:"path"`
	Key         string `yaml:"key"`
	Permissions string `yaml:"permissions,omitempty"` // octal string, e.g. "0600"; defaults to "0600"
}

// ServiceConfig is the full definition of a single service.
type ServiceConfig struct {
	ServiceName  string            `yaml:"service_name"`
	ImageName    string            `yaml:"image_name"`
	Host         string            `yaml:"host,omitempty"`
	Description  string            `yaml:"description,omitempty"`
	Icon         string            `yaml:"icon,omitempty"`
	Commands     []string          `yaml:"commands,omitempty"`
	NetworkMode  string            `yaml:"network_mode,omitempty"`
	Privileged   bool              `yaml:"privileged,omitempty"`
	ShmSize      int               `yaml:"shm_size,omitempty"` // /dev/shm size in MB
	Capabilities *Capabilities     `yaml:"capabilities,omitempty"`
	Devices      *DevicesConfig    `yaml:"devices,omitempty"`
	Mounts       []string          `yaml:"mounts,omitempty"`
	Volumes      []string          `yaml:"volumes,omitempty"`
	Env          []string          `yaml:"env,omitempty"`
	Secrets      map[string]string `yaml:"secrets,omitempty"`
	ConfigFiles  []ConfigFile      `yaml:"configfiles,omitempty"`
	Network      *ServiceNetwork   `yaml:"network,omitempty"`
	DNS          ServiceDNS        `yaml:"dns"`
	Auth         *ServiceAuth      `yaml:"auth,omitempty"`
}

// ServiceNetwork describes how a service attaches to Docker networks.
type ServiceNetwork struct {
	Internal    bool           `yaml:"internal,omitempty"`
	ServicePort int            `yaml:"service_port,omitempty"`
	Networks    []NetworkEntry `yaml:"networks,omitempty"`
}

// NetworkEntry is a single network attachment with an optional static IP.
type NetworkEntry struct {
	Name      string `yaml:"name"`
	IPAddress string `yaml:"ip_address,omitempty"`
}

// ServiceDNS controls DNS and reverse-proxy registration.
type ServiceDNS struct {
	Enabled    bool   `yaml:"enabled"`
	DomainName string `yaml:"domain_name,omitempty"`
	// ForwardHost and ForwardPort override the NPM proxy destination for services
	// that cannot supply a network.networks IP — e.g. host-networked containers
	// on a separate subnet. When set, proxy.go uses these instead of serviceIP()
	// and network.service_port.
	ForwardHost string `yaml:"forward_host,omitempty"`
	ForwardPort int    `yaml:"forward_port,omitempty"`
	// Internal, when explicitly set, overrides the network.internal default.
	// nil = not set (dns.go falls back to network.internal).
	// true = internal-only (no Cloudflare record).
	// false = explicitly public (overrides network.internal: true, e.g. kavita).
	Internal *bool `yaml:"internal,omitempty"`
}

// ServiceAuth configures Authentik SSO for a service.
type ServiceAuth struct {
	Enabled          bool       `yaml:"enabled"`
	Group            string     `yaml:"group,omitempty"`
	AdditionalGroups []string   `yaml:"additional_groups,omitempty"`
	Proxy            *ProxyAuth `yaml:"proxy,omitempty"`
	OAuth            *OAuthAuth `yaml:"oauth,omitempty"`
}

// ProxyAuth configures Authentik forward-auth proxy.
type ProxyAuth struct {
	Enabled        bool   `yaml:"enabled"`
	UserSecret     string `yaml:"user_secret,omitempty"`
	PassSecret     string `yaml:"pass_secret,omitempty"`
	AuthSecretName string `yaml:"auth_secret_name,omitempty"` // Infisical key for X-Proxy-Secret header value
}

// OAuthAuth configures Authentik OAuth2/OIDC.
type OAuthAuth struct {
	Enabled      bool              `yaml:"enabled"`
	Keys         map[string]string `yaml:"keys,omitempty"`
	Scopes       []string          `yaml:"scopes,omitempty"`
	RedirectURIs []string          `yaml:"redirect_uris,omitempty"`
}

// Capabilities maps to Docker's cap_add / cap_drop.
type Capabilities struct {
	Add  []string `yaml:"add,omitempty"`
	Drop []string `yaml:"drop,omitempty"`
}

// DevicesConfig controls device passthrough into the container.
type DevicesConfig struct {
	GPU   bool     `yaml:"gpu,omitempty"`
	USB   bool     `yaml:"usb,omitempty"`
	Paths []string `yaml:"paths,omitempty"`
}

// FlatService is a pre-merged service definition ready for resource registration.
// Stack-level env and mounts are already merged into Def.
type FlatService struct {
	// ServiceName is the service_name from YAML — used as the Pulumi resource name.
	ServiceName      string
	Def              ServiceConfig
	StackName        string           // empty for standalone services
	StackNetworks    map[string]NetworkDef
	GeneratedSecrets []string
}
