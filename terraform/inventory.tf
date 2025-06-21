locals {
  # Consolidate all managed services into a single map for Ansible.
  # This makes it easy to add new services in the future.
  ansible_managed_hosts = merge(
    { for key, service in local.stacks.arr_services.services : service.service_name => {
      ansible_host = service.ip_address
      group        = "arr_services"
      }
    },
    {
      (local.services.flaresolverr.service_name) = {
        ansible_host = local.services.flaresolverr.ip_address
        group        = "utility_services"
      }
    },
    {
      (local.services.deluge-vpn.service_name) = {
        ansible_host = local.services.deluge-vpn.ip_address
        group        = "vpn_services"
      }
    }
  )
}

output "ansible_inventory" {
  description = "Dynamic inventory for Ansible, consumable by a script inventory."
  value = {
    _meta = {
      hostvars = local.ansible_managed_hosts
    }
  }
}