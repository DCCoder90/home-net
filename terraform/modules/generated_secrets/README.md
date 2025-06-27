<!-- BEGIN_TF_DOCS -->
## Requirements

| Name | Version |
|------|---------|
| <a name="requirement_infisical"></a> [infisical](#requirement\_infisical) | 0.15.19 |

## Providers

| Name | Version |
|------|---------|
| <a name="provider_infisical"></a> [infisical](#provider\_infisical) | 0.15.19 |
| <a name="provider_random"></a> [random](#provider\_random) | n/a |

## Resources

| Name | Type |
|------|------|
| [infisical_secret.stored_secret](https://registry.terraform.io/providers/infisical/infisical/0.15.19/docs/resources/secret) | resource |
| [random_password.generated_secret](https://registry.terraform.io/providers/hashicorp/random/latest/docs/resources/password) | resource |
| [infisical_projects.home-net](https://registry.terraform.io/providers/infisical/infisical/0.15.19/docs/data-sources/projects) | data source |

## Inputs

| Name | Description | Type | Default | Required |
|------|-------------|------|---------|:--------:|
| <a name="input_name"></a> [name](#input\_name) | List of secret names to generate | `list(string)` | n/a | yes |
| <a name="input_project_slug"></a> [project\_slug](#input\_project\_slug) | Slug of project to store secret in | `string` | `"home-net-ln-sy"` | no |

## Outputs

| Name | Description |
|------|-------------|
| <a name="output_generated_secrets_list"></a> [generated\_secrets\_list](#output\_generated\_secrets\_list) | A list of objects, each containing the name and generated value of a secret. |
| <a name="output_generated_secrets_map"></a> [generated\_secrets\_map](#output\_generated\_secrets\_map) | A map of generated secret names to their values. |
<!-- END_TF_DOCS -->