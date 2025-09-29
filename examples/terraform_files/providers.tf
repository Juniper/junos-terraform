terraform {
  required_providers {
    junos-vqfx = {
      source = "hashicorp/junos-vqfx"
    }
		junos-vsrx = {
      source = "hashicorp/junos-vsrx"
    }
  }
}

provider "junos-vqfx" {
    host     = "dc1-spine1"
    port     = 22
    username = "jcluser"
    password = "Juniper!1"
    alias    = "dc1_spine1"
}

provider "junos-vqfx" {
    host     = "dc1-spine2"
    port     = 22
    username = "jcluser"
    password = "Juniper!1"
    alias    = "dc1_spine2"
}

provider "junos-vqfx" {
    host     = "dc1-borderleaf1"
    port     = 22
    username = "jcluser"
    password = "Juniper!1"
    alias    = "dc1_borderleaf1"
}

provider "junos-vqfx" {
    host     = "dc1-borderleaf2"
    port     = 22
    username = "jcluser"
    password = "Juniper!1"
    alias    = "dc1_borderleaf2"
}

provider "junos-vqfx" {
    host     = "dc1-leaf1"
    port     = 22
    username = "jcluser"
    password = "Juniper!1"
    alias    = "dc1_leaf1"
}

provider "junos-vqfx" {
    host     = "dc1-leaf2"
    port     = 22
    username = "jcluser"
    password = "Juniper!1"
    alias    = "dc1_leaf2"
}

provider "junos-vqfx" {
    host     = "dc1-leaf3"
    port     = 22
    username = "jcluser"
    password = "Juniper!1"
    alias    = "dc1_leaf3"
}

provider "junos-vqfx" {
    host     = "dc2-spine1"
    port     = 22
    username = "jcluser"
    password = "Juniper!1"
    alias    = "dc2_spine1"
}

provider "junos-vqfx" {
    host     = "dc2-spine2"
    port     = 22
    username = "jcluser"
    password = "Juniper!1"
    alias    = "dc2_spine2"
}

provider "junos-vsrx" {
    host     = "dc1-firewall1"
    port     = 22
    username = "jcluser"
    password = "Juniper!1"
    alias    = "dc1_firewall1"
}

provider "junos-vsrx" {
    host     = "dc1-firewall2"
    port     = 22
    username = "jcluser"
    password = "Juniper!1"
    alias    = "dc1_firewall2"
}

provider "junos-vsrx" {
    host     = "dc2-firewall1"
    port     = 22
    username = "jcluser"
    password = "Juniper!1"
    alias    = "dc2_firewall1"
}

provider "junos-vsrx" {
    host     = "dc2-firewall2"
    port     = 22
    username = "jcluser"
    password = "Juniper!1"
    alias    = "dc2_firewall2"
}
