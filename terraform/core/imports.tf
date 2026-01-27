# -------------------------------------------------------------------
# Import blocks for existing core containers
# Remove this file after a successful apply imports them into state.
# -------------------------------------------------------------------

# ---- Technitium DNS ----
import {
  to = module.technitium.docker_container.container
  id = "ee9ec2823d8e123ce4a87f722c72ec60696a67ba11b7958c118b2fd3930e87c9"
}

# ---- Nginx Proxy Manager ----
import {
  to = module.nginx_proxy_manager.docker_container.container
  id = "dc2c940e26c873c4c202b664c32f74e7df9cdb55cd663e52e82ab50469f7d21a"
}

# ---- Authentik PostgreSQL ----
import {
  to = module.authentik_postgres.docker_container.container
  id = "c6f1374cfaae7f8197d9ab6930b101a4972aa8a45e4a26edc3d4f1a32aa5ea61"
}

# ---- Authentik Redis ----
import {
  to = module.authentik_redis.docker_container.container
  id = "221c4708e6c8b2c292a542562b81bce5ca62ea30e55e66945497785c1942587b"
}

# ---- Authentik Server ----
import {
  to = module.authentik_server.docker_container.container
  id = "7066425afab43f9121532286fa368e452de9e879b3a08a3d8e2a4e31ec9ed098"
}

# ---- Authentik Worker ----
import {
  to = module.authentik_worker.docker_container.container
  id = "c4de4c260c8f7c41cb4183b5cb2865b279a3766d887fd52d25d8de210efa58db"
}
