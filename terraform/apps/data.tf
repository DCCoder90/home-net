data "terraform_remote_state" "core" {
  backend = "remote"

  config = {
    organization = "DCCoders"
    workspaces = {
      name = "home-core"
    }
  }
}
