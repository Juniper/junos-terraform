terraform {
  required_providers {
    junos-vmx = {
      source = "juniper/providers/junos-vmx"
      version = "20.41.12"
    }
  }
}

provider "junos-vmx" {
    host = "localhost"
    port = 8300
    username = "root"
    password = "juniper123"
    sshkey = ""
}

module "vmx_1" {
  source = "./vmx_1"

  providers = {junos-vmx = junos-vmx}

  depends_on = [junos-vmx_destroycommit.commit-main]
}

resource "junos-vmx_commit" "commit-main" {
  resource_name = "commit"
  depends_on = [module.vmx_1]
}

resource "junos-vmx_destroycommit" "commit-main" {
  resource_name = "destroycommit"
}