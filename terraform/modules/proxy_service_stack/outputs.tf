output "container_names" {
  value = [for docker in module.service_docker : docker.name]
}