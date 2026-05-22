terraform {
  required_providers {
    junos-vqfx-evpn-vxlan = {
      source = "hashicorp/junos-vqfx-evpn-vxlan"
    }
	junos-vsrx-evpn-vxlan = {
      source = "hashicorp/junos-vsrx-evpn-vxlan"
    }
  }
}


provider "junos-vqfx-evpn-vxlan" {
    host     = var.junos_vqfx_endpoints["dc1_spine1"].host
    port     = var.junos_vqfx_endpoints["dc1_spine1"].port
    username = var.junos_username
    password = var.junos_password
    alias    = "dc1_spine1"
}

provider "junos-vqfx-evpn-vxlan" {
    host     = var.junos_vqfx_endpoints["dc1_spine2"].host
    port     = var.junos_vqfx_endpoints["dc1_spine2"].port
    username = var.junos_username
    password = var.junos_password
    alias    = "dc1_spine2"
}

provider "junos-vqfx-evpn-vxlan" {
    host     = var.junos_vqfx_endpoints["dc1_borderleaf1"].host
    port     = var.junos_vqfx_endpoints["dc1_borderleaf1"].port
    username = var.junos_username
    password = var.junos_password
    alias    = "dc1_borderleaf1"
}

provider "junos-vqfx-evpn-vxlan" {
    host     = var.junos_vqfx_endpoints["dc1_borderleaf2"].host
    port     = var.junos_vqfx_endpoints["dc1_borderleaf2"].port
    username = var.junos_username
    password = var.junos_password
    alias    = "dc1_borderleaf2"
}

provider "junos-vqfx-evpn-vxlan" {
    host     = var.junos_vqfx_endpoints["dc1_leaf1"].host
    port     = var.junos_vqfx_endpoints["dc1_leaf1"].port
    username = var.junos_username
    password = var.junos_password
    alias    = "dc1_leaf1"
}

provider "junos-vqfx-evpn-vxlan" {
    host     = var.junos_vqfx_endpoints["dc1_leaf2"].host
    port     = var.junos_vqfx_endpoints["dc1_leaf2"].port
    username = var.junos_username
    password = var.junos_password
    alias    = "dc1_leaf2"
}

provider "junos-vqfx-evpn-vxlan" {
    host     = var.junos_vqfx_endpoints["dc1_leaf3"].host
    port     = var.junos_vqfx_endpoints["dc1_leaf3"].port
    username = var.junos_username
    password = var.junos_password
    alias    = "dc1_leaf3"
}

provider "junos-vqfx-evpn-vxlan" {
    host     = var.junos_vqfx_endpoints["dc2_spine1"].host
    port     = var.junos_vqfx_endpoints["dc2_spine1"].port
    username = var.junos_username
    password = var.junos_password
    alias    = "dc2_spine1"
}

provider "junos-vqfx-evpn-vxlan" {
    host     = var.junos_vqfx_endpoints["dc2_spine2"].host
    port     = var.junos_vqfx_endpoints["dc2_spine2"].port
    username = var.junos_username
    password = var.junos_password
    alias    = "dc2_spine2"
}

provider "junos-vsrx-evpn-vxlan" {
    host     = var.junos_vsrx_endpoints["dc1_firewall1"].host
    port     = var.junos_vsrx_endpoints["dc1_firewall1"].port
    username = var.junos_username
    password = var.junos_password
    alias    = "dc1_firewall1"
}

provider "junos-vsrx-evpn-vxlan" {
    host     = var.junos_vsrx_endpoints["dc1_firewall2"].host
    port     = var.junos_vsrx_endpoints["dc1_firewall2"].port
    username = var.junos_username
    password = var.junos_password
    alias    = "dc1_firewall2"
}

provider "junos-vsrx-evpn-vxlan" {
    host     = var.junos_vsrx_endpoints["dc2_firewall1"].host
    port     = var.junos_vsrx_endpoints["dc2_firewall1"].port
    username = var.junos_username
    password = var.junos_password
    alias    = "dc2_firewall1"
}

provider "junos-vsrx-evpn-vxlan" {
    host     = var.junos_vsrx_endpoints["dc2_firewall2"].host
    port     = var.junos_vsrx_endpoints["dc2_firewall2"].port
    username = var.junos_username
    password = var.junos_password
    alias    = "dc2_firewall2"
}
