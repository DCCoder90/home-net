locals {
  zone_name = "mallett.family"
  
  username = "default_user"
  password = "SecretPassword1!"

  proxy_ip = "192.168.4.2"
  authentik_ip = "192.168.4.55"
  authentik_port = 9000

  services = {
    1 = {
      service_name = "prowlarr"
      image_name = "linuxserver/prowlarr:1.37.0"
      username = "prowlarr_user"
      password = "ProwlarrPassword1!"
      domain_name = "prowlarr.mallett.family"
      service_port = 9696
      ip_address = "192.168.5.13"
    }
  }
}