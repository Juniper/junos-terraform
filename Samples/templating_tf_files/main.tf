terraform {
  required_providers {
    junos-vsrx = {
      source  = "juniper/providers/junos-vsrx"
      version = "19.41.101"
    }
  }
}

provider "junos-vsrx" {
  host     = "XX.XX.XX.XX"
  port     = 830
  username = "username"
  password = "password"
  sshkey   = ""
}

module "vsrx_1" {
  source    = "./vsrx_1"
  bgp       = var.bgp
  providers = { junos-vsrx = junos-vsrx }

  depends_on = [junos-vsrx_destroycommit.commit-main]
}

resource "junos-vsrx_commit" "commit-main" {
  resource_name = "commit"
  depends_on    = [module.vsrx_1]
}

resource "junos-vsrx_destroycommit" "commit-main" {
  resource_name = "destroycommit"
}
