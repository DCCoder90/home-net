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
| [authentik_group.group](https://registry.terraform.io/providers/goauthentik/authentik/2025.4.0/docs/resources/group) | resource |
| [authentik_user.admin_user](https://registry.terraform.io/providers/goauthentik/authentik/2025.4.0/docs/resources/user) | resource |
| [authentik_user.name](https://registry.terraform.io/providers/goauthentik/authentik/2025.4.0/docs/resources/user) | resource |

## Inputs

| Name | Description | Type | Default | Required |
|------|-------------|------|---------|:--------:|
| <a name="input_admin_email"></a> [admin\_email](#input\_admin\_email) | Admin email address | `any` | n/a | yes |
| <a name="input_admin_username"></a> [admin\_username](#input\_admin\_username) | Admin username | `any` | n/a | yes |
| <a name="input_users"></a> [users](#input\_users) | List of users to create in authentik | <pre>list(object({<br/>    username = string<br/>    email    = string<br/>    password = string<br/>  }))</pre> | `[]` | no |
<!-- END_TF_DOCS -->