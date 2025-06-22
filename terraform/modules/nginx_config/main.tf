resource "nginxproxymanager_access_list" "internal_access" {
  name = "Internal Only"

  access = [
    {
      directive = "allow"
      address   = "192.168.1.0/24"
    },
    {
      directive = "allow"
      address   = "192.168.2.0/24"
    },
    {
      directive = "allow"
      address   = "192.168.4.0/23"
    },
  ]

  pass_auth   = false
  satisfy_any = false
}

resource "nginxproxymanager_access_list" "cloudflare_access" {
  name = "Cloudflare"

  access = [
    //Home network subnets
    {
      directive = "allow"
      address   = "192.168.1.0/24"
    },
    {
      directive = "allow"
      address   = "192.168.2.0/24"
    },
    {
      directive = "allow"
      address   = "192.168.4.0/23"
    },
    //Cloudflare subnets
    {
      directive = "allow"
      address   = "103.21.244.0/22"
    },
    {
      directive = "allow"
      address   = "103.22.200.0/22"
    },
    {
      directive = "allow"
      address   = "103.31.4.0/22"
    },
    {
      directive = "allow"
      address   = "104.16.0.0/13"
    },
    {
      directive = "allow"
      address   = "104.24.0.0/14"
    },
    {
      directive = "allow"
      address   = "108.162.192.0/18"
    },
    {
      directive = "allow"
      address   = "131.0.72.0/22"
    },
    {
      directive = "allow"
      address   = "141.101.64.0/18"
    },
    {
      directive = "allow"
      address   = "162.158.0.0/15"
    },
    {
      directive = "allow"
      address   = "172.64.0.0/13"
    },
    {
      directive = "allow"
      address   = "173.245.48.0/20"
    },
    {
      directive = "allow"
      address   = "188.114.96.0/20"
    },
    {
      directive = "allow"
      address   = "190.93.240.0/20"
    },
    {
      directive = "allow"
      address   = "197.234.240.0/22"
    },
    {
      directive = "allow"
      address   = "198.41.128.0/17"
    },
    {
      directive = "allow"
      address   = "2400:cb00::/32"
    },
    {
      directive = "allow"
      address   = "2606:4700::/32"
    },
    {
      directive = "allow"
      address   = "2803:f800::/32"
    },
    {
      directive = "allow"
      address   = "2405:b500::/32"
    },
    {
      directive = "allow"
      address   = "2405:8100::/32"
    },
    {
      directive = "allow"
      address   = "2a06:98c0::/29"
    },
    {
      directive = "allow"
      address   = "2c0f:f248::/32"
    },
  ]

  pass_auth   = false
  satisfy_any = false
}