locals {
  # Create a map of service IP addresses, trying different network names.
  service_ip_addresses = {
    # This logic finds the first available IP address for a service based on a preferred network order (br1, then br0, then any).
    # It safely handles cases where a service has no network configuration, no networks defined, or no static IP on any network.
    # It concatenates lists of IPs from each check and adds a final `[null]` to ensure the list is never empty, then takes the first element.
    for service_key, service_config in var.stack.services : service_key => concat(
      [for n in try(service_config.network.networks, []) : n.ip_address if n.name == "br1" && n.ip_address != null],
      [for n in try(service_config.network.networks, []) : n.ip_address if n.name == "br0" && n.ip_address != null],
      [for n in try(service_config.network.networks, []) : n.ip_address if n.ip_address != null],
      [null]
    )[0]
  }

  # Filter out networks that already exist.
  creatable_networks = {
    for name, config in coalesce(var.stack.networks, {}) : name => config
    if !contains(var.system.existing_networks, name)
  }

  # Create a map of Nginx Proxy Manager access lists by name.
  npm_access_lists_by_name = { for al in data.nginxproxymanager_access_lists.access_lists.access_lists : al.name => al.id }
}

data "nginxproxymanager_access_lists" "access_lists" {}
