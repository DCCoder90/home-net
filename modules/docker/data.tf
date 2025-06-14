//Data for br1 and br0 are manually configured in Unraid
//eth0 - br0 - 192.168.1.0/24
//eth1 - br1 - 192.168.4.0/23
data "docker_network" "main_host" {
  provider = docker
  name     = "br1"
}

data "docker_network" "secondary_host" {
  provider = docker
  name     = "br0"
}