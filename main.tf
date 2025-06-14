terraform {
  cloud { 
    
    organization = "DCCoders" 

    workspaces { 
      name = "home" 
    } 
  }  

  required_providers {
    nginxproxymanager = {
      source  = "Sander0542/nginxproxymanager"
      version = "1.1.1"
    }

    portainer = {
      source = "portainer/portainer"
      version = "1.4.2"
    }

    docker = {
      source  = "kreuzwerker/docker"
      version = "3.6.2"
    }
  }
}

module "nginx_conf" {
  source = "./modules/nginx_config"
}
