<!-- BEGIN_TF_DOCS -->
## Requirements

| Name | Version |
|------|---------|
| <a name="requirement_authentik"></a> [authentik](#requirement\_authentik) | 2025.4.0 |

## Providers

| Name | Version |
|------|---------|
| <a name="provider_authentik"></a> [authentik](#provider\_authentik) | 2025.4.0 |

## Resources

| Name | Type |
|------|------|
| [authentik_application.name](https://registry.terraform.io/providers/goauthentik/authentik/2025.4.0/docs/resources/application) | resource |
| [authentik_group.group](https://registry.terraform.io/providers/goauthentik/authentik/2025.4.0/docs/resources/group) | resource |
| [authentik_policy_binding.app_binding](https://registry.terraform.io/providers/goauthentik/authentik/2025.4.0/docs/resources/policy_binding) | resource |
| [authentik_provider_proxy.name](https://registry.terraform.io/providers/goauthentik/authentik/2025.4.0/docs/resources/provider_proxy) | resource |
| [authentik_flow.default-authorization-flow](https://registry.terraform.io/providers/goauthentik/authentik/2025.4.0/docs/data-sources/flow) | data source |
| [authentik_flow.default-invalidation-flow](https://registry.terraform.io/providers/goauthentik/authentik/2025.4.0/docs/data-sources/flow) | data source |
| [authentik_group.access_group](https://registry.terraform.io/providers/goauthentik/authentik/2025.4.0/docs/data-sources/group) | data source |
| [authentik_user.akadmin](https://registry.terraform.io/providers/goauthentik/authentik/2025.4.0/docs/data-sources/user) | data source |

## Inputs

| Name | Description | Type | Default | Required |
|------|-------------|------|---------|:--------:|
| <a name="input_external_host"></a> [external\_host](#input\_external\_host) | The domain name used to access the application, (includes HTTPS) | `string` | n/a | yes |
| <a name="input_internal_host"></a> [internal\_host](#input\_internal\_host) | The IP and port used to connect the the application | `string` | n/a | yes |
| <a name="input_name"></a> [name](#input\_name) | Name of the application | `string` | n/a | yes |
| <a name="input_access_group_attributes"></a> [access\_group\_attributes](#input\_access\_group\_attributes) | JSON format expected. Use jsonencode() to pass objects. | `string` | `"{}"` | no |
| <a name="input_access_group_name"></a> [access\_group\_name](#input\_access\_group\_name) | Group to assign access to the application | `string` | `""` | no |
| <a name="input_access_group_roles"></a> [access\_group\_roles](#input\_access\_group\_roles) | Roles to add to access group (if create\_access\_group = true) | `list(string)` | `[]` | no |
| <a name="input_create_access_group"></a> [create\_access\_group](#input\_create\_access\_group) | If true, create new access group, otherwise use pre-created access group | `bool` | `false` | no |
| <a name="input_description"></a> [description](#input\_description) | Description of the application | `string` | `""` | no |
| <a name="input_group"></a> [group](#input\_group) | Group to assign the application to | `string` | `""` | no |
| <a name="input_password_attribute"></a> [password\_attribute](#input\_password\_attribute) | The attribute to use for the password in basic auth | `string` | `"password"` | no |
| <a name="input_user_to_add_to_access_group"></a> [user\_to\_add\_to\_access\_group](#input\_user\_to\_add\_to\_access\_group) | User to add to access group | `string` | `""` | no |
| <a name="input_username_attribute"></a> [username\_attribute](#input\_username\_attribute) | The attribute to use for the username in basic auth | `string` | `"username"` | no |

## Outputs

| Name | Description |
|------|-------------|
| <a name="output_client_id"></a> [client\_id](#output\_client\_id) | Auth client ID |
<!-- END_TF_DOCS -->