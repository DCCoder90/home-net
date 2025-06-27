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
| [docker_network.private_network](https://registry.terraform.io/providers/kreuzwerker/docker/3.6.2/docs/resources/network) | resource |

## Inputs

| Name | Description | Type | Default | Required |
|------|-------------|------|---------|:--------:|
| <a name="input_networks"></a> [networks](#input\_networks) | List of networks to create | <pre>map(object({<br/>    internal = optional(bool, false)<br/>    driver   = optional(string, "bridge")<br/>    options  = optional(map(string), {})<br/>  }))</pre> | `{}` | no |

## Outputs

| Name | Description |
|------|-------------|
| <a name="output_networks"></a> [networks](#output\_networks) | A map of the created Docker networks, including their ID, name, and scope. |
<!-- END_TF_DOCS -->