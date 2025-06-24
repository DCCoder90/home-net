data "authentik_group" "access_group" {
  count = var.create_access_group ? 0 : 1
  name  = var.access_group_name
}

resource "authentik_group" "group" {
  count      = var.create_access_group ? 1 : 0
  name       = var.access_group_name
  users      = [data.authentik_user.akadmin.id]
  roles      = var.access_group_roles
  attributes = var.access_group_attributes
}

data "authentik_user" "akadmin" {
  username = var.user_to_add_to_access_group
}

resource "authentik_policy_binding" "app_binding" {
  count  = var.access_group_name == null && var.create_access_group == false ? 0 : 1
  target = authentik_application.name.uuid
  group  = authentik_group.group[0].id
  order  = 0
}