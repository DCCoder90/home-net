resource "authentik_user" "dccoder" {
  username = "DCCoder"
  email    = var.dccoder_email
}

resource "authentik_user" "name" {
  for_each = { for user in var.users : user.username => user }

  username = each.value.username
  email    = each.value.email
  password = each.value.password
}

resource "authentik_group" "group" {
  name         = "tf_admins"
  users        = [authentik_user.dccoder.id]
  is_superuser = true
}