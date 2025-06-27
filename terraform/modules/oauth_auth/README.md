<!-- BEGIN_TF_DOCS -->
## Requirements

| Name | Version |
|------|---------|
| <a name="requirement_authentik"></a> [authentik](#requirement\_authentik) | 2025.4.0 |

## Providers

| Name | Version |
|------|---------|
| <a name="provider_authentik"></a> [authentik](#provider\_authentik) | 2025.4.0 |
| <a name="provider_random"></a> [random](#provider\_random) | n/a |

## Resources

| Name | Type |
|------|------|
| [authentik_application.name](https://registry.terraform.io/providers/goauthentik/authentik/2025.4.0/docs/resources/application) | resource |
| [authentik_group.access_groups](https://registry.terraform.io/providers/goauthentik/authentik/2025.4.0/docs/resources/group) | resource |
| [authentik_policy_binding.app_binding](https://registry.terraform.io/providers/goauthentik/authentik/2025.4.0/docs/resources/policy_binding) | resource |
| [authentik_provider_oauth2.name](https://registry.terraform.io/providers/goauthentik/authentik/2025.4.0/docs/resources/provider_oauth2) | resource |
| [random_string.example](https://registry.terraform.io/providers/hashicorp/random/latest/docs/resources/string) | resource |
| [authentik_certificate_key_pair.generated](https://registry.terraform.io/providers/goauthentik/authentik/2025.4.0/docs/data-sources/certificate_key_pair) | data source |
| [authentik_flow.default-authorization-flow](https://registry.terraform.io/providers/goauthentik/authentik/2025.4.0/docs/data-sources/flow) | data source |
| [authentik_flow.default-invalidation-flow](https://registry.terraform.io/providers/goauthentik/authentik/2025.4.0/docs/data-sources/flow) | data source |
| [authentik_flow.default-provider-authorization-implicit-consent](https://registry.terraform.io/providers/goauthentik/authentik/2025.4.0/docs/data-sources/flow) | data source |
| [authentik_group.preexisting_groups_to_bind](https://registry.terraform.io/providers/goauthentik/authentik/2025.4.0/docs/data-sources/group) | data source |
| [authentik_property_mapping_provider_scope.scope-email](https://registry.terraform.io/providers/goauthentik/authentik/2025.4.0/docs/data-sources/property_mapping_provider_scope) | data source |
| [authentik_property_mapping_provider_scope.scope-openid](https://registry.terraform.io/providers/goauthentik/authentik/2025.4.0/docs/data-sources/property_mapping_provider_scope) | data source |
| [authentik_property_mapping_provider_scope.scope-profile](https://registry.terraform.io/providers/goauthentik/authentik/2025.4.0/docs/data-sources/property_mapping_provider_scope) | data source |
| [authentik_provider_oauth2_config.oauth_config](https://registry.terraform.io/providers/goauthentik/authentik/2025.4.0/docs/data-sources/provider_oauth2_config) | data source |
| [authentik_user.users_to_bind](https://registry.terraform.io/providers/goauthentik/authentik/2025.4.0/docs/data-sources/user) | data source |

## Inputs

| Name | Description | Type | Default | Required |
|------|-------------|------|---------|:--------:|
| <a name="input_name"></a> [name](#input\_name) | Name of the application | `string` | n/a | yes |
| <a name="input_access_group_name"></a> [access\_group\_name](#input\_access\_group\_name) | Group to assign access to the application | `string` | `""` | no |
| <a name="input_access_group_roles"></a> [access\_group\_roles](#input\_access\_group\_roles) | Roles to add to access group (if create\_access\_group = true) | `list(string)` | `[]` | no |
| <a name="input_allowed_redirect_uris"></a> [allowed\_redirect\_uris](#input\_allowed\_redirect\_uris) | List of redirect URIs for the OAuth2 provider | `list(map(string))` | <pre>[<br/>  {<br/>    "matching_mode": "strict",<br/>    "url": "http://localhost"<br/>  }<br/>]</pre> | no |
| <a name="input_create_access_group"></a> [create\_access\_group](#input\_create\_access\_group) | If true, create new access group(s), otherwise use pre-created access group(s) | `bool` | `false` | no |
| <a name="input_description"></a> [description](#input\_description) | Description of the application | `string` | `""` | no |
| <a name="input_group"></a> [group](#input\_group) | Group to assign the application to | `string` | `null` | no |
| <a name="input_user_to_add_to_access_group"></a> [user\_to\_add\_to\_access\_group](#input\_user\_to\_add\_to\_access\_group) | User to add to access group | `string` | `""` | no |

## Outputs

| Name | Description |
|------|-------------|
| <a name="output_application_name"></a> [application\_name](#output\_application\_name) | Name of application |
| <a name="output_application_slug"></a> [application\_slug](#output\_application\_slug) | Slug of application |
| <a name="output_authorize_url"></a> [authorize\_url](#output\_authorize\_url) | n/a |
| <a name="output_client_id"></a> [client\_id](#output\_client\_id) | Shared client ID for the OAuth2 provider |
| <a name="output_client_secret"></a> [client\_secret](#output\_client\_secret) | Shared client secret for the OAuth2 provider |
| <a name="output_issuer_url"></a> [issuer\_url](#output\_issuer\_url) | n/a |
| <a name="output_jwks_url"></a> [jwks\_url](#output\_jwks\_url) | n/a |
| <a name="output_logout_url"></a> [logout\_url](#output\_logout\_url) | n/a |
| <a name="output_provider_info_url"></a> [provider\_info\_url](#output\_provider\_info\_url) | n/a |
| <a name="output_token_url"></a> [token\_url](#output\_token\_url) | n/a |
| <a name="output_user_info_url"></a> [user\_info\_url](#output\_user\_info\_url) | n/a |
<!-- END_TF_DOCS -->