terraform {
  required_providers {
    junos-vmx = {
      source  = "juniper/providers/junos-vmx"
      version = "22.41.101"
    }
  }
}

provider "junos-vmx" {
  host     = "aa.aa.aa.aa"
  port     = 830
  username = "username"
  password = "password"
  sshkey   = ""
  alias    = "R1"
}
provider "junos-vmx" {
  host     = "bb.bb.bb.bb"
  port     = 830
  username = "username"
  password = "password"
  sshkey   = ""
  alias    = "R2"
}
provider "junos-vmx" {
  host     = "cc.cc.cc.cc"
  port     = 830
  username = "username"
  password = "password"
  sshkey   = ""
  alias    = "R3"
}

module "vmx_1" {
  source          = "./vmx_1"
  interfaces      = var.vmx_1.interfaces
  routing-options = var.vmx_1.routing-options
  isis            = var.vmx_1.isis
  ldp             = var.vmx_1.ldp
  bgp             = var.vmx_1.bgp
  providers       = { junos-vmx = junos-vmx.R1 }
  depends_on      = [junos-vmx_destroycommit.commit-main]
}

resource "junos-vmx_commit" "commit-main" {
  provider      = junos-vmx.R1
  resource_name = "commit"
  depends_on    = [module.vmx_1]
}

resource "junos-vmx_destroycommit" "commit-main" {
  provider      = junos-vmx.R1
  resource_name = "destroycommit"
}

module "vmx_2" {
  source     = "./vmx_2"
  interfaces = var.vmx_2.interfaces
  routing-options = var.vmx_2.routing-options
  isis            = var.vmx_2.isis
  ldp             = var.vmx_2.ldp
  bgp             = var.vmx_2.bgp
  providers  = { junos-vmx = junos-vmx.R2 }
  depends_on = [junos-vmx_destroycommit.commit-main_vmx2]
}

resource "junos-vmx_commit" "commit-main_vmx2" {
  provider      = junos-vmx.R2
  resource_name = "commit"
  depends_on    = [module.vmx_2]
}

resource "junos-vmx_destroycommit" "commit-main_vmx2" {
  provider      = junos-vmx.R2
  resource_name = "destroycommit"
}

module "vmx_3" {
  source     = "./vmx_3"
  interfaces = var.vmx_3.interfaces
  routing-options = var.vmx_3.routing-options
  isis            = var.vmx_3.isis
  ldp             = var.vmx_3.ldp
  bgp             = var.vmx_3.bgp
  providers  = { junos-vmx = junos-vmx.R3 }
  depends_on = [junos-vmx_destroycommit.commit-main_vmx3]
}

resource "junos-vmx_commit" "commit-main_vmx3" {
  provider      = junos-vmx.R3
  resource_name = "commit"
  depends_on    = [module.vmx_3]
}

resource "junos-vmx_destroycommit" "commit-main_vmx3" {
  provider      = junos-vmx.R3
  resource_name = "destroycommit"
}
