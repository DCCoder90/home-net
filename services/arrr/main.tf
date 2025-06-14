module "service_docker" {
  source = "../../modules/docker"
  
  container_name = "AAutomated_Test"
  container_image = "ghcr.io/flaresolverr/flaresolverr:v3.3.21"
  attach_to_br0 = false
  attach_to_br1 = true
  br1_ipv4_addr = "192.168.5.13"
}