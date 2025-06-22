output "networks" {
  description = "A map of the created Docker networks, including their ID, name, and scope."
  value = {
    for name, network in docker_network.private_network : name => {
      id    = network.id
      scope = network.scope
      name  = network.name
    }
  }
}