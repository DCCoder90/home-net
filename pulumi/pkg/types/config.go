package types

// ServersFile is the top-level structure of config/servers.yaml.
type ServersFile struct {
	Servers map[string]ServerConfig `yaml:"servers"`
}

// ServerConfig describes a single managed host.
type ServerConfig struct {
	Networks NetworkNames `yaml:"networks"`
}

// NetworkNames holds the primary (services) and optional secondary network interface names.
type NetworkNames struct {
	Primary   string `yaml:"primary"`
	Secondary string `yaml:"secondary,omitempty"`
}

// SystemConfig maps to config/system.yaml.
type SystemConfig struct {
	Authentik struct {
		IPAddress string `yaml:"ip_address"`
		Port      int    `yaml:"port"`
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

// ServiceConfig is the full definition of a single service from config/services/*.yaml.
type ServiceConfig struct {
	ServiceName  string            `yaml:"service_name"`
	ImageName    string            `yaml:"image_name"`
	Host         string            `yaml:"host,omitempty"` // defaults to "tower"
	Description  string            `yaml:"description,omitempty"`
	Icon         string            `yaml:"icon,omitempty"`
	Commands     []string          `yaml:"commands,omitempty"`
	NetworkMode  string            `yaml:"network_mode,omitempty"`
	EnableGPU    bool              `yaml:"enable_gpu,omitempty"`
	Capabilities *Capabilities     `yaml:"capabilities,omitempty"`
	Mounts       []string          `yaml:"mounts,omitempty"`
	Volumes      []string          `yaml:"volumes,omitempty"`
	Env          []string          `yaml:"env,omitempty"`
	Secrets      map[string]string `yaml:"secrets,omitempty"`
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

// ServiceDNS controls external DNS and reverse-proxy registration.
type ServiceDNS struct {
	Enabled    bool   `yaml:"enabled"`
	DomainName string `yaml:"domain_name,omitempty"`
	Internal   bool   `yaml:"internal,omitempty"`
}

// ServiceAuth configures Authentik SSO for a service.
type ServiceAuth struct {
	Enabled bool       `yaml:"enabled"`
	Group   string     `yaml:"group,omitempty"`
	Proxy   *ProxyAuth `yaml:"proxy,omitempty"`
	OAuth   *OAuthAuth `yaml:"oauth,omitempty"`
}

// ProxyAuth configures Authentik forward-auth proxy with optional basic-auth injection.
type ProxyAuth struct {
	Enabled    bool   `yaml:"enabled"`
	UserSecret string `yaml:"user_secret,omitempty"` // Infisical key for username
	PassSecret string `yaml:"pass_secret,omitempty"` // Infisical key for password
}

// OAuthAuth configures Authentik OAuth2/OIDC.
type OAuthAuth struct {
	Enabled      bool              `yaml:"enabled"`
	Keys         map[string]string `yaml:"keys,omitempty"`          // env-var → oidc-output-key
	Scopes       []string          `yaml:"scopes,omitempty"`
	RedirectURIs []string          `yaml:"redirect_uris,omitempty"`
}

// Capabilities maps to Docker's cap_add / cap_drop.
type Capabilities struct {
	Add  []string `yaml:"add,omitempty"`
	Drop []string `yaml:"drop,omitempty"`
}

// StackConfig is the top-level structure of a config/stacks/*.yaml file.
type StackConfig struct {
	Host             string                   `yaml:"host,omitempty"`
	Env              []string                 `yaml:"env,omitempty"`
	Mounts           []string                 `yaml:"mounts,omitempty"`
	Networks         map[string]NetworkDef    `yaml:"networks,omitempty"`
	GeneratedSecrets []string                 `yaml:"generated_secrets,omitempty"`
	Services         map[string]ServiceConfig `yaml:"services"`
}

// NetworkDef describes a custom Docker network to create for a stack.
type NetworkDef struct {
	Internal bool   `yaml:"internal,omitempty"`
	Driver   string `yaml:"driver,omitempty"`
}
