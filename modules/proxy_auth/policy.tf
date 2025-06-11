data "authentik_group" "access_group" {
  count = var.access_group == null ? 0 : 1
  name  = var.access_group_name
}

# Create a policy for the specified group
resource "authentik_policy_group" "access_policy" {
  count = var.access_group == null ? 0 : 1
  name  = "${var.name}-group-policy"
  group = data.authentik_group.access_group[0].id
}

resource "authentik_policy_binding" "app_binding" {
  count  = var.access_group == null ? 0 : 1
  target = authentik_application.name.uuid
  policy = authentik_policy_group.access_policy[0].id
  order  = 0
}