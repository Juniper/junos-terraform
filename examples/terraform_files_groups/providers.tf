terraform {
  required_providers {
    junos-vqfx-evpn-vxlan-groups = {
      source = "hashicorp/junos-vqfx-evpn-vxlan-groups"
    }
  }
}


provider "junos-vqfx-evpn-vxlan-groups" {
    host     = "dc1-spine1"
    port     = 22
    username = "jcluser"
    password = "Juniper!1"
    alias    = "dc1_spine1"
}


provider "junos-vqfx-evpn-vxlan-groups" {
    host     = "dc1-spine2"
    port     = 22
    username = "jcluser"
    password = "Juniper!1"
    alias    = "dc1_spine2"
}


provider "junos-vqfx-evpn-vxlan-groups" {
    host     = "dc1-borderleaf1"
    port     = 22
    username = "jcluser"
    password = "Juniper!1"
    alias    = "dc1_borderleaf1"
}


provider "junos-vqfx-evpn-vxlan-groups" {
    host     = "dc1-borderleaf2"
    port     = 22
    username = "jcluser"
    password = "Juniper!1"
    alias    = "dc1_borderleaf2"
}


provider "junos-vqfx-evpn-vxlan-groups" {
    host     = "dc1-leaf1"
    port     = 22
    username = "jcluser"
    password = "Juniper!1"
    alias    = "dc1_leaf1"
}


provider "junos-vqfx-evpn-vxlan-groups" {
    host     = "dc1-leaf2"
    port     = 22
    username = "jcluser"
    password = "Juniper!1"
    alias    = "dc1_leaf2"
}


provider "junos-vqfx-evpn-vxlan-groups" {
    host     = "dc1-leaf3"
    port     = 22
    username = "jcluser"
    password = "Juniper!1"
    alias    = "dc1_leaf3"
}


provider "junos-vqfx-evpn-vxlan-groups" {
    host     = "dc2-spine1"
    port     = 22
    username = "jcluser"
    password = "Juniper!1"
    alias    = "dc2_spine1"
}


provider "junos-vqfx-evpn-vxlan-groups" {
    host     = "dc2-spine2"
    port     = 22
    username = "jcluser"
    password = "Juniper!1"
    alias    = "dc2_spine2"
}
