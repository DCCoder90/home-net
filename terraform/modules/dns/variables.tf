variable "internal_only" {
  description = "If true, the application is only accessible internally"
  type        = bool
  default     = false
}

variable "access_list_id" {
  type        = string
  description = "ID of access list in NGINX"
  default     = null
}

variable "service_port" {
  description = "Port on which the service is running"
  type        = number
  default     = 80
}

variable "service_ipv4" {
  description = "IPv4 address of the service"
  type        = string
}

variable "forward_scheme" {
  description = "Scheme of forwarded service"
  type        = string
  default     = "http"
}

variable "dns_cloudflare_api_token" {
  description = "Cloudflare API token for DNS updates"
  type        = string
  default     = ""
  sensitive   = true
}

variable "admin_email" {
  description = "Email address for the admin user"
  type        = string
  sensitive   = true
}

variable "external_host_ipv4" {
  description = "External host for the application"
  type        = string
  default     = ""
}

variable "external_host_ipv6" {
  description = "External host for the application"
  type        = string
  default     = ""
}

variable "internal_host_ipv4" {
  description = "Internal host for the application"
  type        = string
  default     = ""
}

variable "internal_host_ipv6" {
  description = "Internal host for the application"
  type        = string
  default     = ""
}

variable "enable_ipv4" {
  description = "If true, enable IPv4 for the application"
  type        = bool
  default     = true
}

variable "enable_ipv6" {
  description = "If true, enable IPv6 for the application"
  type        = bool
  default     = false
}

variable "zone_name" {
  description = "Zone name for the DNS record"
  type        = string
}

variable "domain_name" {
  description = "Domain name for the application"
  type        = string
}

variable "proxied_domain" {
  description = "If true, the DNS record is proxied through Cloudflare"
  type        = bool
  default     = true
}

variable "ttl" {
  description = "Time to live for the DNS record"
  type        = number
  default     = 3600
}