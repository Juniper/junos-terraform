vmx_1 = {
  interfaces = {
    "ge-0/0/0" = {
      unit = {
        0 = {
          ipv4_addr   = "10.0.1.0/31",
          family_iso  = "true",
          family_mpls = "true",
        },
      }
    },
    "ge-0/0/1" = {
      unit = {
        0 = {
          ipv4_addr   = "10.0.2.0/31",
          family_iso  = "true",
          family_mpls = "true",
        },
      }
    },
    "lo0" = {
      unit = {
        0 = {
          ipv4_addr   = "1.1.1.1/32",
          family_iso  = "true",
          family_mpls = "false",
          iso_addr    = "49.0001.0010.0100.1001.00",
        },
      },
    }
  }
  routing-options = {
    router_id = "1.1.1.1",
    as_number = 64512,
  }
  isis = {
    interfaces = [
      {
        name     = "ge-0/0/0.0",
        ldp_sync = "true",
        p2p      = "true",
      },
      {
        name     = "ge-0/0/1.0",
        ldp_sync = "true",
        p2p      = "true",
      },
      {
        name    = "lo0.0",
        passive = "true",
      }
    ],
    level = [
      {
        name    = "1",
        disable = "true"
      },
      {
        name         = "2",
        wide_metrics = "true",
        auth_type    = "md5",
        auth_key     = "testing123",
      }
    ]
  }
  ldp = {
    interfaces = [
      { name = "ge-0/0/0.0" },
      { name = "ge-0/0/1.0" }
    ],
    track_igp_metric = "true",
    deaggrgate       = "true"
  }
  bgp = {
    IBGP = {
      group_type          = "internal",
      local_addr          = "1.1.1.1",
      af_inet_unicast      = "true",
      af_evpn_signaling   = "true",
      af_inet_vpn_unicast = "true",
      auth_key            = "testing123",
      neighbors = [
        { name = "2.2.2.2" },
        { name = "3.3.3.3" },
      ]
    }
  }
}
vmx_2 = {
  interfaces = {
    "ge-0/0/0" = {
      unit = {
        0 = {
          ipv4_addr   = "10.0.1.1/31",
          family_iso  = "true",
          family_mpls = "true",
        },
      }
    },
    "ge-0/0/1" = {
      unit = {
        0 = {
          ipv4_addr   = "10.0.3.0/31",
          family_iso  = "true",
          family_mpls = "true",
        },
      }
    },
    "lo0" = {
      unit = {
        0 = {
          ipv4_addr   = "2.2.2.2/32",
          family_iso  = "true",
          family_mpls = "false",
          iso_addr    = "49.0002.0020.0200.2002.00",
        },
      },
    }
  }
  routing-options = {
    router_id = "2.2.2.2",
    as_number = 64512,
  }
  isis = {
    interfaces = [
      {
        name     = "ge-0/0/0.0",
        ldp_sync = "true",
        p2p      = "true",
      },
      {
        name     = "ge-0/0/1.0",
        ldp_sync = "true",
        p2p      = "true",
      },
      {
        name    = "lo0.0",
        passive = "true",
      }
    ],
    level = [
      {
        name    = "1",
        disable = "true"
      },
      {
        name         = "2",
        wide_metrics = "true",
        auth_type    = "md5",
        auth_key     = "testing123",
      }
    ]
  }
  ldp = {
    interfaces = [
      { name = "ge-0/0/0.0" },
      { name = "ge-0/0/1.0" }
    ],
    track_igp_metric = "true",
    deaggrgate       = "true"
  }
  bgp = {
    IBGP = {
      group_type          = "internal",
      local_addr          = "2.2.2.2",
      af_inet_unicast      = "true",
      af_evpn_signaling   = "true",
      af_inet_vpn_unicast = "true",
      auth_key            = "testing123",
      neighbors = [
        { name = "1.1.1.1" },
        { name = "3.3.3.3" },
      ]
    }
  }
}
vmx_3 = {
  interfaces = {
    "ge-0/0/0" = {
      unit = {
        0 = {
          ipv4_addr   = "10.0.2.1/31",
          family_iso  = "true",
          family_mpls = "true",
        },
      }
    },
    "ge-0/0/1" = {
      unit = {
        0 = {
          ipv4_addr   = "10.0.3.1/31",
          family_iso  = "true",
          family_mpls = "true",
        },
      }
    },
    "lo0" = {
      unit = {
        0 = {
          ipv4_addr   = "3.3.3.3/32",
          family_iso  = "true",
          family_mpls = "false",
          iso_addr    = "49.0003.0030.0300.3003.00",
        },
      },
    }
  }
  routing-options = {
    router_id = "3.3.3.3",
    as_number = 64512,
  }
  isis = {
    interfaces = [
      {
        name     = "ge-0/0/0.0",
        ldp_sync = "true",
        p2p      = "true",
      },
      {
        name     = "ge-0/0/1.0",
        ldp_sync = "true",
        p2p      = "true",
      },
      {
        name    = "lo0.0",
        passive = "true",
      }
    ],
    level = [
      {
        name    = "1",
        disable = "true"
      },
      {
        name         = "2",
        wide_metrics = "true",
        auth_type    = "md5",
        auth_key     = "testing123",
      }
    ]
  }
  ldp = {
    interfaces = [
      { name = "ge-0/0/0.0" },
      { name = "ge-0/0/1.0" }
    ],
    track_igp_metric = "true",
    deaggrgate       = "true"
  }
  bgp = {
    IBGP = {
      group_type          = "internal",
      local_addr          = "3.3.3.3",
      af_inet_unicast      = "true",
      af_evpn_signaling   = "true",
      af_inet_vpn_unicast = "true",
      auth_key            = "testing123",
      neighbors = [
        { name = "1.1.1.1" },
        { name = "2.2.2.2" },
      ]
    }
  }
}
