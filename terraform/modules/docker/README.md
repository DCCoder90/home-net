<!-- BEGIN_TF_DOCS -->
## Requirements

| Name | Version |
|------|---------|
| <a name="requirement_docker"></a> [docker](#requirement\_docker) | 3.6.2 |

## Providers

| Name | Version |
|------|---------|
| <a name="provider_docker"></a> [docker](#provider\_docker) | 3.6.2 |

## Resources

| Name | Type |
|------|------|
| [docker_container.container](https://registry.terraform.io/providers/kreuzwerker/docker/3.6.2/docs/resources/container) | resource |
| [docker_image.main](https://registry.terraform.io/providers/kreuzwerker/docker/3.6.2/docs/resources/image) | resource |
| [docker_volume.managed_volumes](https://registry.terraform.io/providers/kreuzwerker/docker/3.6.2/docs/resources/volume) | resource |
| [docker_network.main_host](https://registry.terraform.io/providers/kreuzwerker/docker/3.6.2/docs/data-sources/network) | data source |
| [docker_network.secondary_host](https://registry.terraform.io/providers/kreuzwerker/docker/3.6.2/docs/data-sources/network) | data source |

## Inputs

| Name | Description | Type | Default | Required |
|------|-------------|------|---------|:--------:|
| <a name="input_container_image"></a> [container\_image](#input\_container\_image) | Name and tag of the image to use (ex. ubuntu:latest) | `string` | n/a | yes |
| <a name="input_container_name"></a> [container\_name](#input\_container\_name) | Name for the container to be created | `string` | n/a | yes |
| <a name="input_commands"></a> [commands](#input\_commands) | n/a | `list(string)` | `null` | no |
| <a name="input_container_capabilities"></a> [container\_capabilities](#input\_container\_capabilities) | Linux capabilities to add or drop for the container. Example: { add = ["NET\_ADMIN"], drop = ["SYS\_ADMIN"] } | <pre>object({<br/>    add  = optional(list(string))<br/>    drop = optional(list(string))<br/>  })</pre> | `{}` | no |
| <a name="input_container_dns_servers"></a> [container\_dns\_servers](#input\_container\_dns\_servers) | List of DNS servers to use for the container. | `list(string)` | <pre>[<br/>  "8.8.8.8",<br/>  "1.1.1.1"<br/>]</pre> | no |
| <a name="input_container_network_mode"></a> [container\_network\_mode](#input\_container\_network\_mode) | Base network mode for the container. Defaults to 'null'. Common values: 'none', 'bridge', 'host', 'container:<name\|id>'.<br/>Behavior for modes:<br/>- 'bridge': (Default if not overridden by br1/br0) Connects to the default Docker bridge.<br/>- 'none': If neither br1 nor br0 are attached, creates the container without a network interface. If br1 or br0 are attached, one of them becomes the primary.<br/>- 'host': Uses the host's network stack. DNS, port mappings, and `attach_to_br0`/`attach_to_br1` are ignored by the module.<br/>DNS and port mappings are applied unless the effective network mode is 'host'. | `string` | `null` | no |
| <a name="input_container_ports"></a> [container\_ports](#input\_container\_ports) | List of port mappings for the container. Each object can have 'internal' (required, number), 'external' (optional, number), 'ip' (optional, string), and 'protocol' (optional, string, defaults to 'tcp'). | <pre>list(object({<br/>    internal = number<br/>    external = optional(number)<br/>    ip       = optional(string)<br/>    protocol = optional(string)<br/>  }))</pre> | `[]` | no |
| <a name="input_container_privileged_mode"></a> [container\_privileged\_mode](#input\_container\_privileged\_mode) | Run the container in privileged mode. Grants the container almost all the capabilities of the host machine. | `bool` | `false` | no |
| <a name="input_container_restart"></a> [container\_restart](#input\_container\_restart) | The restart policy for the container. Must be one of 'no', 'on-failure', 'always', 'unless-stopped'. | `string` | `"unless-stopped"` | no |
| <a name="input_container_user"></a> [container\_user](#input\_container\_user) | User used for run the first process. Format is user or user:group which user and group can be passed literraly or by name. | `string` | `null` | no |
| <a name="input_container_volumes"></a> [container\_volumes](#input\_container\_volumes) | List of volume mounts. Each object defines 'container\_path' and EITHER 'host\_path' (for bind mounts) OR 'volume\_name' (for Docker-managed volumes). The 'manage\_volume\_lifecycle' flag determines if named volumes are created by this module. | <pre>list(object({<br/>    container_path = string<br/>    read_only      = optional(bool, false)<br/><br/>    # For bind mounts<br/>    host_path = optional(string)<br/><br/>    # For Docker-managed named volumes<br/>    volume_name             = optional(string)<br/>    manage_volume_lifecycle = optional(bool, true) # If true (default), module creates/manages the named volume. If false, assumes volume_name exists.<br/>    # Options for creating the named volume (used only if volume_name is set AND manage_volume_lifecycle is true)<br/>    driver      = optional(string)<br/>    driver_opts = optional(map(string))<br/>    labels      = optional(map(string))<br/>  }))</pre> | `[]` | no |
| <a name="input_enable_gpu"></a> [enable\_gpu](#input\_enable\_gpu) | If true, use nvidia runtime to add GPU support to the container. | `bool` | `false` | no |
| <a name="input_environment_vars"></a> [environment\_vars](#input\_environment\_vars) | Environment variables to set in the form of KEY=VALUE | `set(string)` | `null` | no |
| <a name="input_icon"></a> [icon](#input\_icon) | n/a | `string` | `null` | no |
| <a name="input_labels"></a> [labels](#input\_labels) | n/a | <pre>set(object({<br/>    label = string<br/>    value = string<br/>  }))</pre> | `null` | no |
| <a name="input_mounts"></a> [mounts](#input\_mounts) | Specification for mounts to be added to containers created as part of the service. | `list(string)` | `[]` | no |
| <a name="input_networks"></a> [networks](#input\_networks) | List of networks to attach to | <pre>list(object({<br/>    name       = string<br/>    ip_address = optional(string, "") # Optional static IP address for the network<br/>  }))</pre> | `[]` | no |
| <a name="input_web_ui"></a> [web\_ui](#input\_web\_ui) | n/a | `string` | `null` | no |

## Outputs

| Name | Description |
|------|-------------|
| <a name="output_created_docker_volumes"></a> [created\_docker\_volumes](#output\_created\_docker\_volumes) | Details of the volumes created by the module. Keyed by volume name. |
| <a name="output_image_name"></a> [image\_name](#output\_image\_name) | The name of the image used by this Docker container. |
| <a name="output_name"></a> [name](#output\_name) | The name of the created Docker container. |
<!-- END_TF_DOCS -->