output "name" {
  description = "The name of the created Docker container."
  value       = docker_container.container.name
}

output "image_name" {
  description = "The name of the image used by this Docker container."
  value       = docker_container.container.image
}

output "created_docker_volumes" {
  description = "Details of the volumes created by the module. Keyed by volume name."
  value = {
    for key, volume_resource in docker_volume.managed_volumes : volume_resource.name => {
      driver      = volume_resource.driver
      driver_opts = volume_resource.driver_opts
      mountpoint  = volume_resource.mountpoint
      labels      = volume_resource.labels
    }
  }
}