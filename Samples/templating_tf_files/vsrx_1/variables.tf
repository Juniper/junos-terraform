variable "bgp" {
  type = map(object({
    neighbor_type  = string,
    local-address  = optional(string),
    local-as = optional(string),
    cluster = optional(string),
    af-evpn-signaling = optional(number),
    multipath      = optional(number),
    mtu-discovery = optional(number),
    multipath-multiple-as = optional(number),
    neighbors = list(object({
      name        = string,
      description = optional(string),
      peer-as     = optional(string),
    })),
    bfd-multiplier = optional(number),
    bfd-minimum-interval = optional(number),
    vpn-apply-export = optional(number)
  }))
}
