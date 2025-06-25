resource "authentik_group" "access_groups" {
  for_each = var.create_access_group ? toset([var.access_group_name]) : toset([])

  name = each.key
  # is_superuser = false # Default, can be made a variable if needed
  # attributes   = {}  # Default, can be made a variable if needed

  users = [
    data.authentik_user.users_to_bind.id
    //for user_config in var.user_to_add_to_access_group :
    //data.authentik_user.users_to_bind[user_config.username].id if contains(user_config.groups, each.key)
  ]
}


data "authentik_user" "users_to_bind" {
  //for_each = { for user_config in var.user_to_add_to_access_group : user_config.username => user_config }
  username = var.user_to_add_to_access_group
}

data "authentik_group" "preexisting_groups_to_bind" {
  for_each = !var.create_access_group && length(var.access_group_name) > 0 ? toset([var.access_group_name]) : toset([])
  name     = each.key
}

resource "authentik_policy_binding" "app_binding" {
  count  = length(var.access_group_name) > 0 ? 1 : 0
  target = authentik_application.name.uuid
  group = var.create_access_group ? (
    length(var.access_group_name) > 0 && length(authentik_group.access_groups) > 0 ? values(authentik_group.access_groups)[0].id : null
    ) : (
    length(var.access_group_name) > 0 && length(data.authentik_group.preexisting_groups_to_bind) > 0 ? values(data.authentik_group.preexisting_groups_to_bind)[0].id : null
  )
  order = 0
}