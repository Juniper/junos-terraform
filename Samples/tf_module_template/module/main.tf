terraform {
  required_providers {
    junos-vsrx = {
      source = "juniper/providers/junos-vsrx"
      version = "20.32.0101"
    }
  }
}

provider "junos-vsrx" {
    host = "localhost"
    port = 8300
    username = "root"
    password = "juniper123"
    sshkey = ""
}

module "qfx_1" {
  source = "./qfx_1"

  providers = {junos-vsrx = junos-vsrx}

  depends_on = [junos-vsrx_destroycommit.commit-main]
}


resource "junos-vsrx_commit" "commit-main" {
  resource_name = "commit"
  depends_on = [module.qfx_1]
}

resource "junos-vsrx_destroycommit" "commit-main" {
  resource_name = "destroycommit"
}
