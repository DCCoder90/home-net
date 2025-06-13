module "service_docker" {
  source = "../../modules/docker"
  
  container_name = "Automated_Test"
  container_image = "ghcr.io/flaresolverr/flaresolverr:latest"
  attach_to_br0 = true
  attach_to_br1 = true
  br1_ipv4_addr = "192.168.5.111"
}