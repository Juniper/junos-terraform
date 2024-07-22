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
      ipv6_addr   = optional(string),
      proxy_macip_adv = optional(string),
      vga_accept_data = optional(string),
      ipv4_vga = optional(string)
    }))
  }))
}
variable "switch-options" {
  type = object({
    vtep_source_interface = optional(string),
    route_distinguisher = optional(string),
    vrf_target = object({
      rt = optional(string),
      auto = optional(string),
    })
  })
}
variable "evpn" {
  type = object({
    encapsulation = optional(string),
    default_gateway = optional(string),
    extended_vni_list = optional(string)
  })
}
variable "vlans" {
  type = map(object({
    vlan_id = optional(number),
    l3_interface = optional(string),
    vxlan = optional(list(object({
      vni = optional(number)
    })))

  }))
}
variable "routing-options" {
  type = object({
    as_number = optional(number),
    router_id = optional(string)
    forwarding_table = object({
      export_policy = optional(string),
      ecmp_fast_reroute = optional(string),
    })
  })
}
variable "forwarding-options" {
  type = object({
    storm_control_profile_all = optional(string)
  })
}
variable "policy-statement" {
  type = map(object({
    term = map(object({
      from_protocol = optional(string),
      then_accept = optional(string),
      then_reject = optional(string),
      load_balance_per_packet = optional(string)
    }))
  }))
}
variable "bgp" {
  type = object({
    log_updown = optional(string),
    graceful_restart = optional(string),
    groups = list(object({
    group_type = optional(string),
    group_name            = string,
    local_addr            = optional(string),
    local_as              = optional(number),
    cluster               = optional(string),
    af_evpn_signaling     = optional(string),
    af_inet_unicast        = optional(string),
    af_inet_vpn_unicast   = optional(string),
    multipath             = optional(string),
    mtu_discovery         = optional(string),
    multipath_multiple_as = optional(string),
    export_policy = optional(string),
    import_policy = optional(string),
    neighbors = list(object({
      name        = string,
      description = optional(string),
      peer_as     = optional(number),
    })),
    bfd_multiplier       = optional(number),
    bfd_minimum_interval = optional(number),
    vpn_apply_export     = optional(number),
    auth_key             = optional(string)
  }))
})
}
variable "lldp" {
  type = object({
    interface = optional(string),
  })
}
variable "igmp-snooping" {
    type = object({
      vlan = optional(string)
    })
}
