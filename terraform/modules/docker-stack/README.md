<!-- BEGIN_TF_DOCS -->
## Requirements

| Name | Version |
|------|---------|
| <a name="requirement_infisical"></a> [infisical](#requirement\_infisical) | 0.15.19 |
| <a name="requirement_nginxproxymanager"></a> [nginxproxymanager](#requirement\_nginxproxymanager) | 1.2.0 |

## Providers

| Name | Version |
|------|---------|
| <a name="provider_infisical"></a> [infisical](#provider\_infisical) | 0.15.19 |
| <a name="provider_nginxproxymanager"></a> [nginxproxymanager](#provider\_nginxproxymanager) | 1.2.0 |
| <a name="provider_random"></a> [random](#provider\_random) | n/a |

## Modules

| Name | Source | Version |
|------|--------|---------|
| <a name="module_custom_network"></a> [custom\_network](#module\_custom\_network) | ../../modules/docker-network | n/a |
| <a name="module_oauth_authentication"></a> [oauth\_authentication](#module\_oauth\_authentication) | ../oauth_auth | n/a |
| <a name="module_proxy_authentication"></a> [proxy\_authentication](#module\_proxy\_authentication) | ../proxy_auth | n/a |
| <a name="module_service_container"></a> [service\_container](#module\_service\_container) | ../../modules/docker | n/a |
| <a name="module_service_dns"></a> [service\_dns](#module\_service\_dns) | ../dns | n/a |

## Resources

| Name | Type |
|------|------|
| [random_password.service_password](https://registry.terraform.io/providers/hashicorp/random/latest/docs/resources/password) | resource |
| [infisical_projects.home-net](https://registry.terraform.io/providers/infisical/infisical/0.15.19/docs/data-sources/projects) | data source |
| [infisical_secrets.generated_secrets](https://registry.terraform.io/providers/infisical/infisical/0.15.19/docs/data-sources/secrets) | data source |
| [nginxproxymanager_access_lists.access_lists](https://registry.terraform.io/providers/Sander0542/nginxproxymanager/1.2.0/docs/data-sources/access_lists) | data source |

## Inputs

| Name | Description | Type | Default | Required |
|------|-------------|------|---------|:--------:|
| <a name="input_stack"></a> [stack](#input\_stack) | Map of services to be deployed with their configurations | <pre>object({<br/>    env               = optional(list(string))<br/>    mounts            = optional(list(string))<br/>    volumes           = optional(list(string))<br/>    generated_secrets = optional(list(string))<br/>    zone_name         = optional(string, null)<br/>    networks = optional(map(object({<br/>      internal = optional(bool, false)<br/>      driver   = optional(string, "bridge")<br/>      options  = optional(map(string), {})<br/>    })))<br/>    services = map(object({<br/>      service_name = string<br/>      image_name   = string<br/>      icon         = optional(string, "https://static-00.iconduck.com/assets.00/container-optimized-os-icon-512x512-9od9q084.png")<br/>      network_mode = optional(string)<br/>      description  = optional(string, "")<br/>      commands     = optional(list(string))<br/>      mounts       = optional(list(string))<br/>      enable_gpu   = optional(bool, false)<br/>      env          = optional(list(string))<br/>      volumes      = optional(list(string))<br/>      capabilities = optional(object({<br/>        add  = optional(list(string))<br/>        drop = optional(list(string))<br/>      }), {})<br/>      dns = object({<br/>        enabled     = optional(bool, false),<br/>        internal    = optional(bool, true),<br/>        domain_name = optional(string)<br/>      })<br/>      auth = optional(object({<br/>        enabled = optional(bool, false),<br/>        proxy   = optional(bool, false),<br/>        group   = optional(string, "Uncategorized"),<br/>        oauth = optional(object({<br/>          enabled       = optional(bool, false),<br/>          keys          = optional(map(string), {}),<br/>          scopes        = optional(list(string)),<br/>          redirect_uris = optional(list(string))<br/>        }), {})<br/>      }))<br/>      network = optional(object({<br/>        internal     = optional(bool, false)<br/>        service_port = optional(number)<br/>        networks = optional(list(object({<br/>          name       = string<br/>          ip_address = optional(string)<br/>        })), [])<br/>      }))<br/>    }))<br/>  })</pre> | n/a | yes |
| <a name="input_system"></a> [system](#input\_system) | A dynamic object containing system-wide configuration variables. Its attributes are accessed via lookup() for flexibility. | `any` | n/a | yes |

## Outputs

| Name | Description |
|------|-------------|
| <a name="output_generated_secrets"></a> [generated\_secrets](#output\_generated\_secrets) | A map of the dynamically generated secrets for the stack. |
<!-- END_TF_DOCS -->