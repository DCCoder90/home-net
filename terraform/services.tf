module "stacks" {
  for_each = local.stacks
  source   = "git@github.com:DCCoder90/home-tf-modules.git//docker-stack?ref=1.0.0"

  stack  = each.value
  system = local.system
}

module "services"{
  for_each = local.services
  source   = "git@github.com:DCCoder90/home-tf-modules.git//docker-service?ref=1.0.0"

  service  = each.value
  system   = local.system
}