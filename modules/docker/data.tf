data "docker_network" "main_host" {
  name = "br1"
}

data "docker_network" "secondary_host" {
  name = "br0"
}