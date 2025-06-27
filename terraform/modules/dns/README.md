<!-- BEGIN_TF_DOCS -->
## Requirements

| Name | Version |
|------|---------|
| <a name="requirement_cloudflare"></a> [cloudflare](#requirement\_cloudflare) | ~> 5 |
| <a name="requirement_nginxproxymanager"></a> [nginxproxymanager](#requirement\_nginxproxymanager) | 1.2.0 |
| <a name="requirement_technitium"></a> [technitium](#requirement\_technitium) | 0.2.0 |

## Providers

| Name | Version |
|------|---------|
| <a name="provider_cloudflare"></a> [cloudflare](#provider\_cloudflare) | ~> 5 |
| <a name="provider_nginxproxymanager"></a> [nginxproxymanager](#provider\_nginxproxymanager) | 1.2.0 |
| <a name="provider_technitium"></a> [technitium](#provider\_technitium) | 0.2.0 |

## Resources

| Name | Type |
|------|------|
| [cloudflare_dns_record.ipv4_dns_record](https://registry.terraform.io/providers/cloudflare/cloudflare/latest/docs/resources/dns_record) | resource |
| [cloudflare_dns_record.ipv6_dns_record](https://registry.terraform.io/providers/cloudflare/cloudflare/latest/docs/resources/dns_record) | resource |
| [nginxproxymanager_certificate_letsencrypt.certificate](https://registry.terraform.io/providers/Sander0542/nginxproxymanager/1.2.0/docs/resources/certificate_letsencrypt) | resource |
| [nginxproxymanager_proxy_host.host](https://registry.terraform.io/providers/Sander0542/nginxproxymanager/1.2.0/docs/resources/proxy_host) | resource |
| [technitium_record.ipv4_dns_record](https://registry.terraform.io/providers/kevynb/technitium/0.2.0/docs/resources/record) | resource |
| [technitium_record.ipv6_dns_record](https://registry.terraform.io/providers/kevynb/technitium/0.2.0/docs/resources/record) | resource |
| [cloudflare_zone.zone](https://registry.terraform.io/providers/cloudflare/cloudflare/latest/docs/data-sources/zone) | data source |

## Inputs

| Name | Description | Type | Default | Required |
|------|-------------|------|---------|:--------:|
| <a name="input_admin_email"></a> [admin\_email](#input\_admin\_email) | Email address for the admin user | `string` | n/a | yes |
| <a name="input_domain_name"></a> [domain\_name](#input\_domain\_name) | Domain name for the application | `string` | n/a | yes |
| <a name="input_service_ipv4"></a> [service\_ipv4](#input\_service\_ipv4) | IPv4 address of the service | `string` | n/a | yes |
| <a name="input_zone_name"></a> [zone\_name](#input\_zone\_name) | Zone name for the DNS record | `string` | n/a | yes |
| <a name="input_access_list_id"></a> [access\_list\_id](#input\_access\_list\_id) | ID of access list in NGINX | `string` | `null` | no |
| <a name="input_dns_cloudflare_api_token"></a> [dns\_cloudflare\_api\_token](#input\_dns\_cloudflare\_api\_token) | Cloudflare API token for DNS updates | `string` | `""` | no |
| <a name="input_enable_ipv4"></a> [enable\_ipv4](#input\_enable\_ipv4) | If true, enable IPv4 for the application | `bool` | `true` | no |
| <a name="input_enable_ipv6"></a> [enable\_ipv6](#input\_enable\_ipv6) | If true, enable IPv6 for the application | `bool` | `false` | no |
| <a name="input_external_host_ipv4"></a> [external\_host\_ipv4](#input\_external\_host\_ipv4) | External host for the application | `string` | `""` | no |
| <a name="input_external_host_ipv6"></a> [external\_host\_ipv6](#input\_external\_host\_ipv6) | External host for the application | `string` | `""` | no |
| <a name="input_forward_scheme"></a> [forward\_scheme](#input\_forward\_scheme) | Scheme of forwarded service | `string` | `"http"` | no |
| <a name="input_internal_host_ipv4"></a> [internal\_host\_ipv4](#input\_internal\_host\_ipv4) | Internal host for the application | `string` | `""` | no |
| <a name="input_internal_host_ipv6"></a> [internal\_host\_ipv6](#input\_internal\_host\_ipv6) | Internal host for the application | `string` | `""` | no |
| <a name="input_internal_only"></a> [internal\_only](#input\_internal\_only) | If true, the application is only accessible internally | `bool` | `false` | no |
| <a name="input_proxied_domain"></a> [proxied\_domain](#input\_proxied\_domain) | If true, the DNS record is proxied through Cloudflare | `bool` | `true` | no |
| <a name="input_service_port"></a> [service\_port](#input\_service\_port) | Port on which the service is running | `number` | `80` | no |
| <a name="input_ttl"></a> [ttl](#input\_ttl) | Time to live for the DNS record | `number` | `3600` | no |
<!-- END_TF_DOCS -->