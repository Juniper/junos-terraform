leaf_1 = {
  interfaces = {
    "xe-0/0/0" = {
      description = "to spine1",
      mtu         = 9192,
      unit = {
        0 = {
          ipv4_addr = "10.1.113.2/30",
        },
      }
    },
    "xe-0/0/1" = {
      description = "to spine2",
      mtu         = 9192,
      unit = {
        0 = {
          ipv4_addr = "10.1.123.2/30",
        },
      }
    },
    "lo0" = {
      unit = {
        0 = {
          ipv4_addr = "10.255.1.3/32",
        },
      }
    },
    "irb" = {
      unit = {
        150 = {
          ipv4_addr       = "192.168.150.1/24",
          proxy_macip_adv = "true",
          vga_accept_data = "true",
          ipv4_vga        = "192.168.150.254"
        },
        250 = {
          ipv4_addr       = "192.168.250.1/24",
          proxy_macip_adv = "true",
          vga_accept_data = "true",
          ipv4_vga        = "192.168.250.254"
        }
      }
    }
  }
  switch-options = {
    vtep_source_interface = "lo0.0",
    route_distinguisher   = "10.255.1.3:1",
    vrf_target = {
      rt   = "target:64589:1111",
      auto = "true",
    }
  }
  evpn = {
    encapsulation     = "vxlan",
    default_gateway   = "no-gateway-community",
    extended_vni_list = "all"
  }
  vlans = {
    VNI_15000 = {
      vlan_id      = 150,
      l3_interface = "irb.150",
      vxlan        = [{ vni = 15000 }],
    }
    VNI_25000 = {
      vlan_id      = 250,
      l3_interface = "irb.250"
      vxlan        = [{ vni = 25000 }],
    },
    default = {
      vlan_id = 1,
      vxlan   = [{ vni = 0 }],
    }
  }
  routing-options = {
    router_id = "10.255.1.3",
    as_number = 64589,
    forwarding_table = {
      export_policy     = "PFE-LB",
      ecmp_fast_reroute = "true",
    }
  }
  forwarding-options = {
    storm_control_profile_all = "default"
  }
  policy-statement = {
    IPCLOS_BGP_EXP = {
      term = {
        loopback = {
          from_protocol = "[direct, bgp]",
          then_accept   = "true"
        },
        default = {
          then_reject = "true"
        }
      }
    },
    IPCLOS_BGP_IMP = {
      term = {
        loopback = {
          from_protocol = "[direct, bgp]",
          then_accept   = "true"
        },
        default = {
          then_reject = "true"
        }
      }
    },
    PFE-LB = {
      term = {
        NA = {
          load_balance_per_packet = "true"
        }
      }
    }
  }
  bgp = {
    log_updown       = "true",
    graceful_restart = "true"
    groups = [
      {
        group_name            = "IPCLOS_eBGP",
        mtu_discovery         = "true",
        export_policy         = "IPCLOS_BGP_EXP",
        import_policy         = "IPCLOS_BGP_IMP",
        local_as              = 65513,
        multipath_multiple_as = "true",
        neighbors = [
          { name = "10.1.113.1", description = "EBGP peering to Spine1", peer_as = 65511 },
          { name = "10.1.123.1", description = "EBGP peering to Spine2", peer_as = 65512 },
        ]
      },
      {
        group_name        = "OVERLAY",
        group_type        = "internal",
        af_evpn_signaling = "true",
        neighbors = [
          { name = "10.255.1.1", description = "IBGP/overlay peering to Spine1" },
          { name = "10.255.1.2", description = "IBGP/overlay peering to Spine2" }
        ]
      }
    ],
  }
  lldp = {
    interface = "all"
  }
  igmp-snooping = {
    vlan = "default"
  }
}
leaf_2 = {
  interfaces = {
    "xe-0/0/0" = {
      description = "to spine1",
      mtu         = 9192,
      unit = {
        0 = {
          ipv4_addr = "10.1.114.2/30",
        },
      }
    },
    "xe-0/0/1" = {
      description = "to spine2",
      mtu         = 9192,
      unit = {
        0 = {
          ipv4_addr = "10.1.124.2/30",
        },
      }
    },
    "lo0" = {
      unit = {
        0 = {
          ipv4_addr = "10.255.1.4/32",
        },
      }
    },
    "irb" = {
      unit = {
        150 = {
          ipv4_addr       = "192.168.150.2/24",
          proxy_macip_adv = "true",
          vga_accept_data = "true",
          ipv4_vga        = "192.168.150.254"
        },
        250 = {
          ipv4_addr       = "192.168.250.3/24",
          proxy_macip_adv = "true",
          vga_accept_data = "true",
          ipv4_vga        = "192.168.250.254"
        }
      }
    }
  }
  switch-options = {
    vtep_source_interface = "lo0.0",
    route_distinguisher   = "10.255.1.4:1",
    vrf_target = {
      rt   = "target:64589:1111",
      auto = "true",
    }
  }
  evpn = {
    encapsulation     = "vxlan",
    default_gateway   = "no-gateway-community",
    extended_vni_list = "all"
  }
  vlans = {
    VNI_15000 = {
      vlan_id      = 150,
      l3_interface = "irb.150",
      vxlan        = [{ vni = 15000 }],
    }
    VNI_25000 = {
      vlan_id      = 250,
      l3_interface = "irb.250"
      vxlan        = [{ vni = 25000 }],
    },
    default = {
      vlan_id = 1,
      vxlan   = [{ vni = 0 }],
    }
  }
  routing-options = {
    router_id = "10.255.1.4",
    as_number = 64589,
    forwarding_table = {
      export_policy     = "PFE-LB",
      ecmp_fast_reroute = "true",
    }
  }
  forwarding-options = {
    storm_control_profile_all = "default"
  }
  policy-statement = {
    IPCLOS_BGP_EXP = {
      term = {
        loopback = {
          from_protocol = "[direct, bgp]",
          then_accept   = "true"
        },
        default = {
          then_reject = "true"
        }
      }
    },
    IPCLOS_BGP_IMP = {
      term = {
        loopback = {
          from_protocol = "[direct, bgp]",
          then_accept   = "true"
        },
        default = {
          then_reject = "true"
        }
      }
    },
    PFE-LB = {
      term = {
        NA = {
          load_balance_per_packet = "true"
        }
      }
    }
  }
  bgp = {
    log_updown       = "true",
    graceful_restart = "true"
    groups = [
      {
        group_name            = "IPCLOS_eBGP",
        mtu_discovery         = "true",
        export_policy         = "IPCLOS_BGP_EXP",
        import_policy         = "IPCLOS_BGP_IMP",
        local_as              = 65514,
        multipath_multiple_as = "true",
        neighbors = [
          { name = "10.1.114.1", description = "EBGP peering to Spine1", peer_as = 65511 },
          { name = "10.1.124.1", description = "EBGP peering to Spine2", peer_as = 65512 },
        ]
      },
      {
        group_name        = "OVERLAY",
        group_type        = "internal",
        af_evpn_signaling = "true",
        neighbors = [
          { name = "10.255.1.1", description = "IBGP/overlay peering to Spine1" },
          { name = "10.255.1.2", description = "IBGP/overlay peering to Spine2" }
        ]
      }
    ],
  }
  lldp = {
    interface = "all"
  }
  igmp-snooping = {
    vlan = "default"
  }
}
spine_1 = {
  interfaces = {
    "xe-0/0/0" = {
      description = "to leaf1",
      mtu         = 9192,
      unit = {
        0 = {
          ipv4_addr = "10.1.113.1/30",
        },
      }
    },
    "xe-0/0/1" = {
      description = "to leaf2",
      mtu         = 9192,
      unit = {
        0 = {
          ipv4_addr = "10.1.114.1/30",
        },
      }
    },
    "lo0" = {
      unit = {
        0 = {
          ipv4_addr = "10.255.1.1/32",
        },
      }
    },
  }
  switch-options = {
  }
  evpn = {
  }
  vlans = {
    default = {
      vlan_id = 1,
      vxlan   = [{ vni = 0 }],
    }
  }
  routing-options = {
    router_id = "10.255.1.1",
    as_number = 64589,
    forwarding_table = {
      export_policy     = "PFE-LB",
      ecmp_fast_reroute = "true",
    }
  }
  forwarding-options = {
    storm_control_profile_all = "default"
  }
  policy-statement = {
    IPCLOS_BGP_EXP = {
      term = {
        loopback = {
          from_protocol = "[direct, bgp]",
          then_accept   = "true"
        },
        default = {
          then_reject = "true"
        }
      }
    },
    IPCLOS_BGP_IMP = {
      term = {
        loopback = {
          from_protocol = "[direct, bgp]",
          then_accept   = "true"
        },
        default = {
          then_reject = "true"
        }
      }
    },
    PFE-LB = {
      term = {
        NA = {
          load_balance_per_packet = "true"
        }
      }
    }
  }
  bgp = {
    log_updown       = "true",
    graceful_restart = "true"
    groups = [
      {
        group_name            = "IPCLOS_eBGP",
        mtu_discovery         = "true",
        export_policy         = "IPCLOS_BGP_EXP",
        import_policy         = "IPCLOS_BGP_IMP",
        local_as              = 65511,
        multipath_multiple_as = "true",
        neighbors = [
          { name = "10.1.113.2", description = "EBGP peering to Leaf1", peer_as = 65513 },
          { name = "10.1.114.2", description = "EBGP peering to Leaf2", peer_as = 65514 },
        ]
      },
      {
        group_name        = "OVERLAY",
        group_type        = "internal",
        af_evpn_signaling = "true",
        local-address = "10.255.1.1",
        cluster = "10.255.1.10"
        neighbors = [
          { name = "10.255.1.3", description = "IBGP/overlay peering to Leaf1" },
          { name = "10.255.1.4", description = "IBGP/overlay peering to Leaf2" }
        ]
      },
      {
        group_name        = "OVERLAY_RR_MESH",
        group_type        = "internal",
        af_evpn_signaling = "true",
        local-address = "10.255.1.1",
        neighbors = [
          { name = "10.255.1.2", description = "IBGP/overlay peering to Spine2" },
        ]
      }
    ],
  }
  lldp = {
    interface = "all"
  }
  igmp-snooping = {
    vlan = "default"
  }
}
spine_2 = {
  interfaces = {
    "xe-0/0/0" = {
      description = "to leaf1",
      mtu         = 9192,
      unit = {
        0 = {
          ipv4_addr = "10.1.123.1/30",
        },
      }
    },
    "xe-0/0/1" = {
      description = "to leaf2",
      mtu         = 9192,
      unit = {
        0 = {
          ipv4_addr = "10.1.124.1/30",
        },
      }
    },
    "lo0" = {
      unit = {
        0 = {
          ipv4_addr = "10.255.1.1/32",
        },
      }
    },
  }
  switch-options = {
  }
  evpn = {
  }
  vlans = {
    default = {
      vlan_id = 1,
      vxlan   = [{ vni = 0 }],
    }
  }
  routing-options = {
    router_id = "10.255.1.2",
    as_number = 64589,
    forwarding_table = {
      export_policy     = "PFE-LB",
      ecmp_fast_reroute = "true",
    }
  }
  forwarding-options = {
    storm_control_profile_all = "default"
  }
  policy-statement = {
    IPCLOS_BGP_EXP = {
      term = {
        loopback = {
          from_protocol = "[direct, bgp]",
          then_accept   = "true"
        },
        default = {
          then_reject = "true"
        }
      }
    },
    IPCLOS_BGP_IMP = {
      term = {
        loopback = {
          from_protocol = "[direct, bgp]",
          then_accept   = "true"
        },
        default = {
          then_reject = "true"
        }
      }
    },
    PFE-LB = {
      term = {
        NA = {
          load_balance_per_packet = "true"
        }
      }
    }
  }
  bgp = {
    log_updown       = "true",
    graceful_restart = "true"
    groups = [
      {
        group_name            = "IPCLOS_eBGP",
        mtu_discovery         = "true",
        export_policy         = "IPCLOS_BGP_EXP",
        import_policy         = "IPCLOS_BGP_IMP",
        local_as              = 65511,
        multipath_multiple_as = "true",
        neighbors = [
          { name = "10.1.123.2", description = "EBGP peering to Leaf1", peer_as = 65513 },
          { name = "10.1.124.2", description = "EBGP peering to Leaf2", peer_as = 65514 },
        ]
      },
      {
        group_name        = "OVERLAY",
        group_type        = "internal",
        af_evpn_signaling = "true",
        local-address = "10.255.1.2",
        cluster = "10.255.1.10"
        neighbors = [
          { name = "10.255.1.3", description = "IBGP/overlay peering to Leaf1" },
          { name = "10.255.1.4", description = "IBGP/overlay peering to Leaf2" }
        ]
      },
      {
        group_name        = "OVERLAY_RR_MESH",
        group_type        = "internal",
        af_evpn_signaling = "true",
        local-address = "10.255.1.2",
        neighbors = [
          { name = "10.255.1.1", description = "IBGP/overlay peering to Spine1" },
        ]
      }
    ],
  }
  lldp = {
    interface = "all"
  }
  igmp-snooping = {
    vlan = "default"
  }
}
