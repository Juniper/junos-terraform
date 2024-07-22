terraform {
  required_providers {
    junos-vqfx = {
      source  = "juniper/providers/junos-vqfx"
      version = "21.31.108"
    }
  }
}

provider "junos-vqfx" {
  host     = "aa.aa.aa.aa"
  port     = 830
  username = "username"
  password = "password"
  sshkey   = ""
  alias    = "leaf1"
}
provider "junos-vqfx" {
  host     = "bb.bb.bb.bb"
  port     = 830
  username = "username"
  password = "password"
  sshkey   = ""
  alias    = "leaf2"
}
provider "junos-vqfx" {
  host     = "cc.cc.cc.cc"
  port     = 830
  username = "username"
  password = "password"
  sshkey   = ""
  alias    = "spine1"
}
provider "junos-vqfx" {
  host     = "dd.dd.dd.dd"
  port     = 830
  username = "username"
  password = "password"
  sshkey   = ""
  alias    = "spine2"
}

module "vqfx_1" {
  source             = "./vqfx_1"
  interfaces         = var.leaf_1.interfaces
  switch-options     = var.leaf_1.switch-options
  evpn               = var.leaf_1.evpn
  vlans              = var.leaf_1.vlans
  routing-options    = var.leaf_1.routing-options
  forwarding-options = var.leaf_1.forwarding-options
  policy-statement   = var.leaf_1.policy-statement
  bgp                = var.leaf_1.bgp
  igmp-snooping      = var.leaf_1.igmp-snooping
  lldp               = var.leaf_1.lldp
  providers          = { junos-vqfx = junos-vqfx.leaf1 }
  depends_on         = [junos-vqfx_destroycommit.commit-main_leaf1]
}

resource "junos-vqfx_commit" "commit-main_leaf1" {
  provider      = junos-vqfx.leaf1
  resource_name = "commit"
  depends_on    = [module.vqfx_1]
}

resource "junos-vqfx_destroycommit" "commit-main_leaf1" {
  provider      = junos-vqfx.leaf1
  resource_name = "destroycommit"
}

module "vqfx_2" {
  source             = "./vqfx_2"
  interfaces         = var.leaf_2.interfaces
  switch-options     = var.leaf_2.switch-options
  evpn               = var.leaf_2.evpn
  vlans              = var.leaf_2.vlans
  routing-options    = var.leaf_2.routing-options
  forwarding-options = var.leaf_2.forwarding-options
  policy-statement   = var.leaf_2.policy-statement
  bgp                = var.leaf_2.bgp
  igmp-snooping      = var.leaf_2.igmp-snooping
  lldp               = var.leaf_2.lldp
  providers          = { junos-vqfx = junos-vqfx.leaf2 }
  depends_on         = [junos-vqfx_destroycommit.commit-main_leaf2]
}

resource "junos-vqfx_commit" "commit-main_leaf2" {
  provider      = junos-vqfx.leaf2
  resource_name = "commit"
  depends_on    = [module.vqfx_2]
}

resource "junos-vqfx_destroycommit" "commit-main_leaf2" {
  provider      = junos-vqfx.leaf2
  resource_name = "destroycommit"
}

module "vqfx_3" {
  source             = "./vqfx_3"
  interfaces         = var.spine_1.interfaces
  switch-options     = var.spine_1.switch-options
  evpn               = var.spine_1.evpn
  vlans              = var.spine_1.vlans
  routing-options    = var.spine_1.routing-options
  forwarding-options = var.spine_1.forwarding-options
  policy-statement   = var.spine_1.policy-statement
  bgp                = var.spine_1.bgp
  igmp-snooping      = var.spine_1.igmp-snooping
  lldp               = var.spine_1.lldp
  providers          = { junos-vqfx = junos-vqfx.spine1 }
  depends_on         = [junos-vqfx_destroycommit.commit-main_spine1]
}

resource "junos-vqfx_commit" "commit-main_spine1" {
  provider      = junos-vqfx.spine1
  resource_name = "commit"
  depends_on    = [module.vqfx_3]
}

resource "junos-vqfx_destroycommit" "commit-main_spine1" {
  provider      = junos-vqfx.spine1
  resource_name = "destroycommit"
}

module "vqfx_4" {
  source             = "./vqfx_4"
  interfaces         = var.spine_2.interfaces
  switch-options     = var.spine_2.switch-options
  evpn               = var.spine_2.evpn
  vlans              = var.spine_2.vlans
  routing-options    = var.spine_2.routing-options
  forwarding-options = var.spine_2.forwarding-options
  policy-statement   = var.spine_2.policy-statement
  bgp                = var.spine_2.bgp
  igmp-snooping      = var.spine_2.igmp-snooping
  lldp               = var.spine_2.lldp
  providers          = { junos-vqfx = junos-vqfx.spine2 }
  depends_on         = [junos-vqfx_destroycommit.commit-main_spine2]
}

resource "junos-vqfx_commit" "commit-main_spine2" {
  provider      = junos-vqfx.spine2
  resource_name = "commit"
  depends_on    = [module.vqfx_4]
}

resource "junos-vqfx_destroycommit" "commit-main_spine2" {
  provider      = junos-vqfx.spine2
  resource_name = "destroycommit"
}
