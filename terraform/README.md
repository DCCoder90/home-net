<!-- BEGIN_TF_DOCS -->
## Requirements

| Name | Version |
|------|---------|
| <a name="requirement_authentik"></a> [authentik](#requirement\_authentik) | 2025.4.0 |
| <a name="requirement_cloudflare"></a> [cloudflare](#requirement\_cloudflare) | ~> 5 |
| <a name="requirement_docker"></a> [docker](#requirement\_docker) | 3.6.2 |
| <a name="requirement_infisical"></a> [infisical](#requirement\_infisical) | 0.15.19 |
| <a name="requirement_nginxproxymanager"></a> [nginxproxymanager](#requirement\_nginxproxymanager) | 1.2.0 |
| <a name="requirement_technitium"></a> [technitium](#requirement\_technitium) | 0.2.0 |

## Modules

| Name | Source | Version |
|------|--------|---------|
| <a name="module_delugevpn_service"></a> [delugevpn\_service](#module\_delugevpn\_service) | ./services/deluge-vpn | n/a |
| <a name="module_flaresolverr_service"></a> [flaresolverr\_service](#module\_flaresolverr\_service) | ./modules/docker | n/a |
| <a name="module_nginx_conf"></a> [nginx\_conf](#module\_nginx\_conf) | ./modules/nginx_config | n/a |
| <a name="module_secrets"></a> [secrets](#module\_secrets) | ./modules/generated_secrets | n/a |
| <a name="module_stacks"></a> [stacks](#module\_stacks) | ./modules/docker-stack | n/a |

## Inputs

| Name | Description | Type | Default | Required |
|------|-------------|------|---------|:--------:|
| <a name="input_cloudflare_api_token"></a> [cloudflare\_api\_token](#input\_cloudflare\_api\_token) | Cloudflare API token | `string` | n/a | yes |
| <a name="input_infiscal_client_id"></a> [infiscal\_client\_id](#input\_infiscal\_client\_id) | The password for the VPN connection. | `string` | n/a | yes |
| <a name="input_infiscal_client_secret"></a> [infiscal\_client\_secret](#input\_infiscal\_client\_secret) | The password for the VPN connection. | `string` | n/a | yes |
| <a name="input_network_admin_email"></a> [network\_admin\_email](#input\_network\_admin\_email) | Network admin email address | `string` | n/a | yes |
| <a name="input_network_admin_username"></a> [network\_admin\_username](#input\_network\_admin\_username) | Admin username | `string` | n/a | yes |
| <a name="input_nginx_proxy_address"></a> [nginx\_proxy\_address](#input\_nginx\_proxy\_address) | n/a | `string` | n/a | yes |
| <a name="input_nginx_proxy_pass"></a> [nginx\_proxy\_pass](#input\_nginx\_proxy\_pass) | n/a | `string` | n/a | yes |
| <a name="input_nginx_proxy_user"></a> [nginx\_proxy\_user](#input\_nginx\_proxy\_user) | n/a | `string` | n/a | yes |
| <a name="input_public_facing_ip"></a> [public\_facing\_ip](#input\_public\_facing\_ip) | My public facing IP address | `string` | n/a | yes |
| <a name="input_technitium_api_token"></a> [technitium\_api\_token](#input\_technitium\_api\_token) | API token for technitium | `string` | n/a | yes |
| <a name="input_technitium_host"></a> [technitium\_host](#input\_technitium\_host) | Host for technitium | `string` | n/a | yes |
| <a name="input_vpn_pass"></a> [vpn\_pass](#input\_vpn\_pass) | The password for the VPN connection. | `string` | n/a | yes |
| <a name="input_vpn_user"></a> [vpn\_user](#input\_vpn\_user) | The username for the VPN connection. | `string` | n/a | yes |
<!-- END_TF_DOCS -->