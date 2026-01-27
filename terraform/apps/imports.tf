# -------------------------------------------------------------------
# Import blocks for existing app containers
# Remove this file after a successful apply imports them into state.
# -------------------------------------------------------------------

# ---- Vaultwarden ----
import {
  to = module.services["vaultwarden"].module.service_container.docker_image.main
  id = "sha256:032fced6c02ed83d6a5da88a2c629de425788ce1c6cc45775766d7cfbeb9f054"
}

import {
  to = module.services["vaultwarden"].module.service_container.docker_container.container
  id = "1544936f387df2886202cd24d9923939674a1dad8d7ba2b3d46e62d07723218e"
}
