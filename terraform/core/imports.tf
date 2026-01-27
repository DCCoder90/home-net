# -------------------------------------------------------------------
# Import blocks for existing core containers
# Remove this file after a successful apply imports them into state.
# -------------------------------------------------------------------

# ---- Technitium DNS ----
import {
  to = module.technitium.docker_image.main
  id = "sha256:84d1bc86034931ce7d20555e456ed6b8f498f0efde396ad68383100aa7415d65"
}

import {
  to = module.technitium.docker_container.container
  id = "ee9ec2823d8e123ce4a87f722c72ec60696a67ba11b7958c118b2fd3930e87c9"
}

# ---- Nginx Proxy Manager ----
import {
  to = module.nginx_proxy_manager.docker_image.main
  id = "sha256:9f5e0949eb63eb73027a58233bc2ae45ea524acdae87636e9b4d95f2c1530304"
}

import {
  to = module.nginx_proxy_manager.docker_container.container
  id = "dc2c940e26c873c4c202b664c32f74e7df9cdb55cd663e52e82ab50469f7d21a"
}

# ---- Authentik PostgreSQL ----
import {
  to = module.authentik_postgres.docker_image.main
  id = "sha256:3f8c1b47f90f0348a6c722bd1a809461ad4fd4286275e6e0750ce67eb3a80ab0"
}

import {
  to = module.authentik_postgres.docker_container.container
  id = "c6f1374cfaae7f8197d9ab6930b101a4972aa8a45e4a26edc3d4f1a32aa5ea61"
}

# ---- Authentik Redis ----
import {
  to = module.authentik_redis.docker_image.main
  id = "sha256:e74faa347ab0b6a3c1b040834a35ea5f20e3aa02460db4bb9d5b4685f3dd3baf"
}

import {
  to = module.authentik_redis.docker_container.container
  id = "221c4708e6c8b2c292a542562b81bce5ca62ea30e55e66945497785c1942587b"
}

# ---- Authentik Server ----
import {
  to = module.authentik_server.docker_image.main
  id = "sha256:e3c993bbf4f6bb9b2bc28f8e81e04ded7fb2e01543ecf0a4777dd7021ce6cc30"
}

import {
  to = module.authentik_server.docker_container.container
  id = "7066425afab43f9121532286fa368e452de9e879b3a08a3d8e2a4e31ec9ed098"
}

# ---- Authentik Worker ----
import {
  to = module.authentik_worker.docker_image.main
  id = "sha256:e3c993bbf4f6bb9b2bc28f8e81e04ded7fb2e01543ecf0a4777dd7021ce6cc30"
}

import {
  to = module.authentik_worker.docker_container.container
  id = "c4de4c260c8f7c41cb4183b5cb2865b279a3766d887fd52d25d8de210efa58db"
}
