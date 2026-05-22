variable "junos_username" {
  description = "Username used by Junos Terraform providers"
  type        = string
  default     = "jcluser"
}

variable "junos_password" {
  description = "Password used by Junos Terraform providers"
  type        = string
  sensitive   = true
  default     = "Juniper!1"
}

variable "junos_vqfx_endpoints" {
  description = "Per-device endpoint map for junos-vqfx-evpn-vxlan provider aliases"
  type = map(object({
    host = string
    port = number
  }))
  default = {
    dc1_spine1      = { host = "dc1-spine1",      port = 22 }
    dc1_spine2      = { host = "dc1-spine2",      port = 22 }
    dc1_borderleaf1 = { host = "dc1-borderleaf1", port = 22 }
    dc1_borderleaf2 = { host = "dc1-borderleaf2", port = 22 }
    dc1_leaf1       = { host = "dc1-leaf1",       port = 22 }
    dc1_leaf2       = { host = "dc1-leaf2",       port = 22 }
    dc1_leaf3       = { host = "dc1-leaf3",       port = 22 }
    dc2_spine1      = { host = "dc2-spine1",      port = 22 }
    dc2_spine2      = { host = "dc2-spine2",      port = 22 }
  }
}

variable "junos_vsrx_endpoints" {
  description = "Per-device endpoint map for junos-vsrx-evpn-vxlan provider aliases"
  type = map(object({
    host = string
    port = number
  }))
  default = {
    dc1_firewall1 = { host = "dc1-firewall1", port = 22 }
    dc1_firewall2 = { host = "dc1-firewall2", port = 22 }
    dc2_firewall1 = { host = "dc2-firewall1", port = 22 }
    dc2_firewall2 = { host = "dc2-firewall2", port = 22 }
  }
}
