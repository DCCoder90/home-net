resource "random_password" "generated" {
  for_each = toset(coalesce(var.stack.generated_secrets, []))

  length  = 48
  special = true
}


module "custom_network" {
  count    = length(local.creatable_networks) > 0 ? 1 : 0
  source   = "../../modules/docker-network"
  networks = local.creatable_networks
}

module "service" {
  for_each = var.stack.services
  source   = "../../modules/service"

  service = each.value
  system = var.system
}
