output "internal_access_list_id" {
  description = "ID of the internal access list"
  value       = nginxproxymanager_access_list.internal_access.id

}

output "cloudflare_access_list_id" {
  description = "ID of the Cloudflare access list"
  value       = nginxproxymanager_access_list.cloudflare_access.id

}