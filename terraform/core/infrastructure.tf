import {
  to = module.nginx_conf.nginxproxymanager_access_list.internal_access
  id = "3"
}

import {
  to = module.nginx_conf.nginxproxymanager_access_list.cloudflare_access
  id = "2"
}

module "nginx_conf" {
  source = "git@github.com:DCCoder90/home-tf-modules.git//nginx_config?ref=1.0.0"
}
