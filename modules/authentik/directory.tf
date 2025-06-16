resource "authentik_user" "admin_user" {
  username = var.admin_username
  email    = var.admin_email
}

resource "authentik_user" "name" {
  for_each = { for user in var.users : user.username => user }

  username = each.value.username
  email    = each.value.email
  password = each.value.password
}

resource "authentik_group" "group" {
  name         = "tf_admins"
  users        = [authentik_user.admin_user.id]
  is_superuser = true
}