variable "interfaces" {
  type = map(object({
    description = optional(string),
    mtu         = optional(number),
    unit = map(object({
      ipv4_addr   = optional(string),
      family_iso  = optional(bool),
      iso_addr    = optional(string),
      vlan_id     = optional(number),
      family_mpls = optional(bool),
      ipv6_addr   = optional(string)
    }))
  }))
}
variable "routing-options" {
  type = object({
    as_number = optional(number),
    router_id = optional(string)
  })
}
variable "isis" {
  type = object({
    interfaces = list(object({
      name     = optional(string),
      ldp_sync = optional(string),
      p2p      = optional(string),
      passive  = optional(string),
    })),
    level = list(object({
      name         = optional(string),
      wide_metrics = optional(string),
      auth_type    = optional(string),
      auth_key     = optional(string),
      disable      = optional(string)
    }))
  })
}
variable "ldp" {
  type = object({
    interfaces = list(object({
      name = optional(string),
    })),
    track_igp_metric = optional(string),
    deaggrgate       = optional(string)
  })
}
variable "bgp" {
  type = map(object({
    group_type            = string,
    local_addr            = optional(string),
    local-as              = optional(string),
    cluster               = optional(string),
    af_evpn_signaling     = optional(string),
    af_inet_unicast        = optional(string),
    af_inet_vpn_unicast   = optional(string),
    multipath             = optional(string),
    mtu_discovery         = optional(string),
    multipath_multiple_as = optional(string),
    neighbors = list(object({
      name        = string,
      description = optional(string),
      peer-as     = optional(string),
    })),
    bfd_multiplier       = optional(number),
    bfd_minimum_interval = optional(number),
    vpn_apply_export     = optional(number),
    auth_key             = optional(string)
  }))
}
