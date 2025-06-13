terraform {
  cloud { 
    
    organization = "DCCoders" 

    workspaces { 
      name = "home" 
    } 
  } 

  required_providers {
    portainer = {
      source = "portainer/portainer"
      version = "1.4.3"
    }
  }
}

resource "portainer_custom_template" "template" {
  title        = "Arrrrrr"
  description  = "Piratey stuff"
  note = "This is by no means complete.  Need to finish it.  Use Arr env from vaultwarden"
  platform     = 1
  type         = 2
  file_content = file("${path.module}/template.yaml")
}