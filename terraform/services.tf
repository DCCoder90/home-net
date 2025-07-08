module "stacks" {
  for_each = local.stacks
  source   = "./modules/docker-stack"

  stack  = each.value
  system = local.system
}

module "services"{
  for_each = local.services
  source   = "./modules/docker-service"

  service  = each.value
  system   = local.system
}