locals {

  # ── Group 0c49e5: dc1-spine1, dc1-spine2, dc2-spine1, dc2-spine2 ────────────────────────────────────────────────
  common_g_0c49e5_apply_groups = [
    "spine"
  ]
  common_g_0c49e5_groups_spine_chassis = [
    {
      aggregated_devices = [
        {
          ethernet = [
            {
              device_count = 24
            }
          ]
        }
      ]
    }
  ]
  common_g_0c49e5_groups_spine_forwarding_options = [
    {
      storm_control_profiles = [
        {
          all = [
            {

            }
          ]
          name = "default"
        }
      ]
    }
  ]
  common_g_0c49e5_groups_spine_interfaces_interface_em1 = {
    name = "em1"
    unit = [
      {
        description = "*** to pfe ***"
        family = [
          {
            inet = [
              {
                address = [
                  {
                    name = "169.254.0.2/24"
                  }
                ]
              }
            ]
          }
        ]
        name = 0
      }
    ]
  }
  common_g_0c49e5_groups_spine_policy_options_policy_statement_IPCLOS_BGP_EXP_term_default = {
    name = "default"
    then = [
      {
        reject = ""
      }
    ]
  }
  common_g_0c49e5_groups_spine_policy_options_policy_statement_IPCLOS_BGP_EXP_term_loopback_from = [
    {
      protocol = [
        "direct",
        "bgp"
      ]
    }
  ]
  common_g_0c49e5_groups_spine_policy_options_policy_statement_IPCLOS_BGP_EXP_term_loopback_then_accept = ""
  common_g_0c49e5_groups_spine_policy_options_policy_statement_IPCLOS_BGP_EXP_term_loopback_then_community_add = ""
  common_g_0c49e5_groups_spine_policy_options_policy_statement_IPCLOS_BGP_IMP = {
    name = "IPCLOS_BGP_IMP"
    term = [
      {
        from = [
          {
            protocol = [
              "bgp",
              "direct"
            ]
          }
        ]
        name = "loopback"
        then = [
          {
            accept = ""
          }
        ]
      },
      {
        name = "default"
        then = [
          {
            reject = ""
          }
        ]
      }
    ]
  }
  common_g_0c49e5_groups_spine_policy_options_policy_statement_PFE_LB = {
    name = "PFE-LB"
    then = [
      {
        load_balance = [
          {
            per_packet = ""
          }
        ]
      }
    ]
  }
  common_g_0c49e5_groups_spine_protocols_bgp_group_EVPN_iBGP_family = [
    {
      evpn = [
        {
          signaling = [
            {

            }
          ]
        }
      ]
    }
  ]
  common_g_0c49e5_groups_spine_protocols_bgp_group_EVPN_iBGP_multipath = [
    {

    }
  ]
  common_g_0c49e5_groups_spine_protocols_bgp_group_IPCLOS_eBGP_bfd_liveness_detection = [
    {
      minimum_interval = 1000
      multiplier = 3
    }
  ]
  common_g_0c49e5_groups_spine_protocols_bgp_group_IPCLOS_eBGP_export = [
    "IPCLOS_BGP_EXP"
  ]
  common_g_0c49e5_groups_spine_protocols_bgp_group_IPCLOS_eBGP_import = [
    "IPCLOS_BGP_IMP"
  ]
  common_g_0c49e5_groups_spine_protocols_bgp_group_IPCLOS_eBGP_multipath = [
    {
      multiple_as = ""
    }
  ]
  common_g_0c49e5_groups_spine_protocols_igmp_snooping = [
    {
      vlan = [
        {
          name = "default"
        }
      ]
    }
  ]
  common_g_0c49e5_groups_spine_protocols_lldp = [
    {
      interface = [
        {
          name = "all"
        }
      ]
    }
  ]
  common_g_0c49e5_groups_spine_routing_options_static = [
    {
      route = [
        {
          name = "0.0.0.0/0"
          next_hop = [
            "100.123.0.1"
          ]
        }
      ]
    }
  ]
  common_g_0c49e5_groups_spine_snmp = [
    {
      community = [
        {
          authorization = "read-only"
          name = "public"
        }
      ]
      contact = "aburston@juniper.net"
      location = "JCL Labs"
    }
  ]
  common_g_0c49e5_groups_spine_system = [
    {
      extensions = [
        {
          providers = [
            {
              license_type = [
                {
                  deployment_scope = [
                    "commercial"
                  ]
                  name = "juniper"
                }
              ]
              name = "juniper"
            },
            {
              license_type = [
                {
                  deployment_scope = [
                    "commercial"
                  ]
                  name = "juniper"
                }
              ]
              name = "chef"
            }
          ]
        }
      ]
      login = [
        {
          message = "***********************************************************************\nThis system is restricted to __________, authorized users for legitimate\nbusiness purposes only. All activity on the system will be logged and\nis subject to monitoring. Unauthorized access, use or modification\nof computers, data therein or data in transit to or from the computers\nis a violation of state and federal laws. Unauthorized activity will\nbe reported to the law enforcement for investigation and possible\nprosecution. __________ reserves the right to investigate, refer for\nprosecution and pursue monetary damages in civil actions in the event\nof unauthorized access.\n***********************************************************************\n"
          user = [
            {
              authentication = [
                {
                  encrypted_password = "$1$a31gJmWG$h9ohikT1ajySf/tVH.gmv1"
                }
              ]
              class = "super-user"
              name = "jcluser"
              uid = 2000
            }
          ]
        }
      ]
      root_authentication = [
        {
          encrypted_password = "$1$DbZ1Q3pj$s48cZytjsmSJRUJAf4LdM."
        }
      ]
      services = [
        {
          extension_service = [
            {
              notification = [
                {
                  allow_clients = [
                    {
                      address = [
                        "0.0.0.0/0"
                      ]
                    }
                  ]
                }
              ]
              request_response = [
                {
                  grpc = [
                    {
                      max_connections = 30
                    }
                  ]
                }
              ]
            }
          ]
          netconf = [
            {
              ssh = [
                {

                }
              ]
            }
          ]
          rest = [
            {
              enable_explorer = ""
              http = [
                {
                  port = 3000
                }
              ]
            }
          ]
          ssh = [
            {
              root_login = "allow"
            }
          ]
        }
      ]
      syslog = [
        {
          file = [
            {
              contents = [
                {
                  name = "any"
                  notice = ""
                },
                {
                  info = ""
                  name = "authorization"
                }
              ]
              name = "messages"
            },
            {
              contents = [
                {
                  any = ""
                  name = "interactive-commands"
                }
              ]
              name = "interactive-commands"
            }
          ]
          user = [
            {
              contents = [
                {
                  emergency = ""
                  name = "any"
                }
              ]
              name = "*"
            }
          ]
        }
      ]
    }
  ]

  # ── Group 49453a: dc2-spine1, dc2-spine2 ────────────────────────────────────────────────
  common_g_49453a_groups_spine_interfaces_interface_ae0 = {
    aggregated_ether_options = [
      {
        lacp = [
          {
            active = ""
            periodic = "fast"
            system_id = "00:00:00:01:01:00"
          }
        ]
      }
    ]
    esi = [
      {
        all_active = ""
        identifier = "00:00:00:00:00:00:00:01:01:00"
      }
    ]
    name = "ae0"
    unit = [
      {
        family = [
          {
            ethernet_switching = [
              {
                vlan = [
                  {
                    members = [
                      1002
                    ]
                  }
                ]
              }
            ]
          }
        ]
        name = 0
      }
    ]
  }
  common_g_49453a_groups_spine_interfaces_interface_ae1 = {
    aggregated_ether_options = [
      {
        lacp = [
          {
            active = ""
            periodic = "fast"
            system_id = "00:00:00:01:02:00"
          }
        ]
      }
    ]
    esi = [
      {
        all_active = ""
        identifier = "00:00:00:00:00:00:00:01:02:00"
      }
    ]
    name = "ae1"
    unit = [
      {
        family = [
          {
            ethernet_switching = [
              {
                vlan = [
                  {
                    members = [
                      1002
                    ]
                  }
                ]
              }
            ]
          }
        ]
        name = 0
      }
    ]
  }
  common_g_49453a_groups_spine_interfaces_interface_irb = {
    name = "irb"
    unit = [
      {
        family = [
          {
            inet = [
              {
                address = [
                  {
                    name = "10.1.2.1/24"
                  }
                ]
              }
            ]
          }
        ]
        mac = "02:0a:01:02:01:18"
        name = 1002
      }
    ]
  }
  common_g_49453a_groups_spine_interfaces_interface_xe_0_0_3 = {
    ether_options = [
      {
        ieee_802_3ad = [
          {
            bundle = "ae0"
          }
        ]
      }
    ]
    name = "xe-0/0/3"
  }
  common_g_49453a_groups_spine_interfaces_interface_xe_0_0_4 = {
    ether_options = [
      {
        ieee_802_3ad = [
          {
            bundle = "ae1"
          }
        ]
      }
    ]
    name = "xe-0/0/4"
  }
  common_g_49453a_groups_spine_policy_options_policy_statement_EVPN_T5_EXPORT = {
    name = "EVPN_T5_EXPORT"
    term = [
      {
        from = [
          {
            protocol = [
              "direct"
            ]
          }
        ]
        name = "fm_direct"
        then = [
          {
            accept = ""
          }
        ]
      },
      {
        from = [
          {
            protocol = [
              "static"
            ]
          }
        ]
        name = "fm_static"
        then = [
          {
            accept = ""
          }
        ]
      },
      {
        from = [
          {
            protocol = [
              "evpn",
              "ospf"
            ]
            route_filter = [
              {
                address = "0.0.0.0/0"
                exact = ""
              }
            ]
          }
        ]
        name = "fm_v4_default"
        then = [
          {
            accept = ""
          }
        ]
      },
      {
        from = [
          {
            protocol = [
              "evpn"
            ]
            route_filter = [
              {
                address = "0.0.0.0/0"
                prefix_length_range = "/32-/32"
              }
            ]
          }
        ]
        name = "fm_v4_host"
        then = [
          {
            accept = ""
          }
        ]
      },
      {
        from = [
          {
            protocol = [
              "evpn"
            ]
            route_filter = [
              {
                address = "0::0/0"
                prefix_length_range = "/128-/128"
              }
            ]
          }
        ]
        name = "fm_v6_host"
        then = [
          {
            accept = ""
          }
        ]
      }
    ]
  }
  common_g_49453a_groups_spine_policy_options_policy_statement_to_ospf = {
    name = "to-ospf"
    term = [
      {
        from = [
          {
            protocol = [
              "evpn"
            ]
            route_filter = [
              {
                address = "10.1.2.0/24"
                orlonger = ""
              }
            ]
          }
        ]
        name = 10
        then = [
          {
            accept = ""
          }
        ]
      },
      {
        name = 100
        then = [
          {
            reject = ""
          }
        ]
      }
    ]
  }
  common_g_49453a_groups_spine_protocols_bgp_group_EVPN_iBGP_local_as = [
    {
      as_number = 65201
    }
  ]
  common_g_49453a_groups_spine_protocols_bgp_group_WAN_OVERLAY_eBGP_family = [
    {
      evpn = [
        {
          signaling = [
            {
              delay_route_advertisements = [
                {
                  minimum_delay = [
                    {
                      routing_uptime = 480
                    }
                  ]
                }
              ]
            }
          ]
        }
      ]
    }
  ]
  common_g_49453a_groups_spine_protocols_bgp_group_WAN_OVERLAY_eBGP_local_as = [
    {
      as_number = 65201
    }
  ]
  common_g_49453a_groups_spine_protocols_bgp_group_WAN_OVERLAY_eBGP_multihop = [
    {
      no_nexthop_change = ""
    }
  ]
  common_g_49453a_groups_spine_protocols_bgp_group_WAN_OVERLAY_eBGP_multipath = [
    {
      multiple_as = ""
    }
  ]
  common_g_49453a_groups_spine_protocols_bgp_group_WAN_OVERLAY_eBGP_neighbor = [
    {
      description = "DCI EBGP peering to 10.30.100.1"
      name = "10.30.100.1"
      peer_as = 65200
    },
    {
      description = "DCI EBGP peering to 10.30.100.2"
      name = "10.30.100.2"
      peer_as = 65200
    }
  ]
  common_g_49453a_groups_spine_protocols_evpn = [
    {
      default_gateway = "do-not-advertise"
      encapsulation = "vxlan"
      extended_vni_list = [
        "all"
      ]
      multicast_mode = "ingress-replication"
      no_core_isolation = ""
    }
  ]
  common_g_49453a_groups_spine_routing_instances_instance_VRF_10001_interface = [
    {
      name = "xe-0/0/1.1"
    },
    {
      name = "xe-0/0/2.1"
    },
    {
      name = "irb.1002"
    },
    {
      name = "lo0.10001"
    }
  ]
  common_g_49453a_groups_spine_routing_instances_instance_VRF_10001_protocols = [
    {
      evpn = [
        {
          ip_prefix_routes = [
            {
              advertise = "direct-nexthop"
              encapsulation = "vxlan"
              export = [
                "EVPN_T5_EXPORT"
              ]
              vni = 10001
            }
          ]
        }
      ]
      ospf = [
        {
          area = [
            {
              interface = [
                {
                  metric = 100
                  name = "xe-0/0/1.1"
                },
                {
                  metric = 200
                  name = "xe-0/0/2.1"
                }
              ]
              name = "0.0.0.0"
            }
          ]
          export = [
            "to-ospf"
          ]
        }
      ]
    }
  ]
  common_g_49453a_groups_spine_routing_instances_instance_VRF_10001_routing_options = [
    {
      auto_export = [
        {

        }
      ]
    }
  ]
  common_g_49453a_groups_spine_routing_instances_instance_VRF_10001_vrf_table_label = [
    {

    }
  ]
  common_g_49453a_groups_spine_routing_instances_instance_VRF_10001_vrf_target = [
    {
      community = "target:1:10001"
    }
  ]
  common_g_49453a_groups_spine_routing_options_forwarding_table = [
    {
      chained_composite_next_hop = [
        {
          ingress = [
            {
              evpn = ""
            }
          ]
        }
      ]
      ecmp_fast_reroute = ""
      export = [
        "PFE-LB"
      ]
    }
  ]
  common_g_49453a_groups_spine_switch_options_vrf_target = [
    {
      auto = [
        {

        }
      ]
      community = "target:9999:9999"
    }
  ]
  common_g_49453a_groups_spine_switch_options_vtep_source_interface = [
    {
      interface_name = "lo0.0"
    }
  ]
  common_g_49453a_groups_spine_vlans = [
    {
      vlan = [
        {
          l3_interface = "irb.1002"
          name = "vlan_1002"
          vlan_id = 1002
          vxlan = [
            {
              vni = 1002
            }
          ]
        }
      ]
    }
  ]

  # ── Group 84da76: dc1-leaf1, dc1-leaf2, dc1-leaf3 ────────────────────────────────────────────────
  common_g_84da76_apply_groups = [
    "leaf"
  ]
  common_g_84da76_groups_leaf_chassis = [
    {
      aggregated_devices = [
        {
          ethernet = [
            {
              device_count = 24
            }
          ]
        }
      ]
    }
  ]
  common_g_84da76_groups_leaf_forwarding_options = [
    {
      storm_control_profiles = [
        {
          all = [
            {

            }
          ]
          name = "default"
        }
      ]
    }
  ]
  common_g_84da76_groups_leaf_interfaces_interface_em1 = {
    name = "em1"
    unit = [
      {
        description = "*** to pfe ***"
        family = [
          {
            inet = [
              {
                address = [
                  {
                    name = "169.254.0.2/24"
                  }
                ]
              }
            ]
          }
        ]
        name = 0
      }
    ]
  }
  common_g_84da76_groups_leaf_interfaces_interface_irb = {
    name = "irb"
    unit = [
      {
        family = [
          {
            inet = [
              {
                address = [
                  {
                    name = "10.1.1.1/24"
                  }
                ]
              }
            ]
          }
        ]
        mac = "02:0a:01:01:01:18"
        name = 1001
      },
      {
        family = [
          {
            inet = [
              {
                address = [
                  {
                    name = "10.2.1.1/24"
                  }
                ]
              }
            ]
          }
        ]
        mac = "02:0a:02:01:01:18"
        name = 2001
      },
      {
        family = [
          {
            inet = [
              {
                address = [
                  {
                    name = "10.3.1.1/24"
                  }
                ]
              }
            ]
          }
        ]
        mac = "02:0a:03:01:01:18"
        name = 3001
      }
    ]
  }
  common_g_84da76_groups_leaf_interfaces_interface_xe_0_0_2 = {
    ether_options = [
      {
        ieee_802_3ad = [
          {
            bundle = "ae0"
          }
        ]
      }
    ]
    name = "xe-0/0/2"
  }
  common_g_84da76_groups_leaf_policy_options_policy_statement_EVPN_T5_EXPORT = {
    name = "EVPN_T5_EXPORT"
    term = [
      {
        from = [
          {
            protocol = [
              "direct"
            ]
          }
        ]
        name = "fm_direct"
        then = [
          {
            accept = ""
          }
        ]
      },
      {
        from = [
          {
            protocol = [
              "static"
            ]
          }
        ]
        name = "fm_static"
        then = [
          {
            accept = ""
          }
        ]
      },
      {
        from = [
          {
            protocol = [
              "evpn"
            ]
            route_filter = [
              {
                address = "0.0.0.0/0"
                prefix_length_range = "/32-/32"
              }
            ]
          }
        ]
        name = "fm_v4_host"
        then = [
          {
            accept = ""
          }
        ]
      },
      {
        from = [
          {
            protocol = [
              "evpn"
            ]
            route_filter = [
              {
                address = "0::0/0"
                prefix_length_range = "/128-/128"
              }
            ]
          }
        ]
        name = "fm_v6_host"
        then = [
          {
            accept = ""
          }
        ]
      }
    ]
  }
  common_g_84da76_groups_leaf_policy_options_policy_statement_IPCLOS_BGP_EXP_term_default = {
    name = "default"
    then = [
      {
        reject = ""
      }
    ]
  }
  common_g_84da76_groups_leaf_policy_options_policy_statement_IPCLOS_BGP_EXP_term_loopback_from = [
    {
      protocol = [
        "direct",
        "bgp"
      ]
    }
  ]
  common_g_84da76_groups_leaf_policy_options_policy_statement_IPCLOS_BGP_EXP_term_loopback_then_accept = ""
  common_g_84da76_groups_leaf_policy_options_policy_statement_IPCLOS_BGP_EXP_term_loopback_then_community_add = ""
  common_g_84da76_groups_leaf_policy_options_policy_statement_IPCLOS_BGP_IMP = {
    name = "IPCLOS_BGP_IMP"
    term = [
      {
        from = [
          {
            protocol = [
              "bgp",
              "direct"
            ]
          }
        ]
        name = "loopback"
        then = [
          {
            accept = ""
          }
        ]
      },
      {
        name = "default"
        then = [
          {
            reject = ""
          }
        ]
      }
    ]
  }
  common_g_84da76_groups_leaf_policy_options_policy_statement_PFE_LB = {
    name = "PFE-LB"
    then = [
      {
        load_balance = [
          {
            per_packet = ""
          }
        ]
      }
    ]
  }
  common_g_84da76_groups_leaf_protocols_bgp_group_EVPN_iBGP_family = [
    {
      evpn = [
        {
          signaling = [
            {

            }
          ]
        }
      ]
    }
  ]
  common_g_84da76_groups_leaf_protocols_bgp_group_EVPN_iBGP_local_as = [
    {
      as_number = 65200
    }
  ]
  common_g_84da76_groups_leaf_protocols_bgp_group_EVPN_iBGP_multipath = [
    {

    }
  ]
  common_g_84da76_groups_leaf_protocols_bgp_group_EVPN_iBGP_neighbor = [
    {
      name = "10.30.100.3"
    },
    {
      name = "10.30.100.4"
    }
  ]
  common_g_84da76_groups_leaf_protocols_bgp_group_IPCLOS_eBGP_bfd_liveness_detection = [
    {
      minimum_interval = 1000
      multiplier = 3
    }
  ]
  common_g_84da76_groups_leaf_protocols_bgp_group_IPCLOS_eBGP_export = [
    "IPCLOS_BGP_EXP"
  ]
  common_g_84da76_groups_leaf_protocols_bgp_group_IPCLOS_eBGP_import = [
    "IPCLOS_BGP_IMP"
  ]
  common_g_84da76_groups_leaf_protocols_bgp_group_IPCLOS_eBGP_multipath = [
    {
      multiple_as = ""
    }
  ]
  common_g_84da76_groups_leaf_protocols_evpn = [
    {
      default_gateway = "do-not-advertise"
      encapsulation = "vxlan"
      extended_vni_list = [
        "all"
      ]
      multicast_mode = "ingress-replication"
    }
  ]
  common_g_84da76_groups_leaf_protocols_igmp_snooping = [
    {
      vlan = [
        {
          name = "default"
        }
      ]
    }
  ]
  common_g_84da76_groups_leaf_protocols_lldp = [
    {
      interface = [
        {
          name = "all"
        }
      ]
    }
  ]
  common_g_84da76_groups_leaf_routing_instances_instance_VRF_10001_interface = [
    {
      name = "irb.1001"
    },
    {
      name = "lo0.10001"
    }
  ]
  common_g_84da76_groups_leaf_routing_instances_instance_VRF_10001_protocols = [
    {
      evpn = [
        {
          ip_prefix_routes = [
            {
              advertise = "direct-nexthop"
              encapsulation = "vxlan"
              export = [
                "EVPN_T5_EXPORT"
              ]
              vni = 10001
            }
          ]
        }
      ]
    }
  ]
  common_g_84da76_groups_leaf_routing_instances_instance_VRF_10001_routing_options = [
    {
      auto_export = [
        {

        }
      ]
    }
  ]
  common_g_84da76_groups_leaf_routing_instances_instance_VRF_10001_vrf_table_label = [
    {

    }
  ]
  common_g_84da76_groups_leaf_routing_instances_instance_VRF_10001_vrf_target = [
    {
      community = "target:1:10001"
    }
  ]
  common_g_84da76_groups_leaf_routing_instances_instance_VRF_10002_interface = [
    {
      name = "irb.2001"
    },
    {
      name = "lo0.10002"
    }
  ]
  common_g_84da76_groups_leaf_routing_instances_instance_VRF_10002_protocols = [
    {
      evpn = [
        {
          ip_prefix_routes = [
            {
              advertise = "direct-nexthop"
              encapsulation = "vxlan"
              export = [
                "EVPN_T5_EXPORT"
              ]
              vni = 10002
            }
          ]
        }
      ]
    }
  ]
  common_g_84da76_groups_leaf_routing_instances_instance_VRF_10002_routing_options = [
    {
      auto_export = [
        {

        }
      ]
    }
  ]
  common_g_84da76_groups_leaf_routing_instances_instance_VRF_10002_vrf_table_label = [
    {

    }
  ]
  common_g_84da76_groups_leaf_routing_instances_instance_VRF_10002_vrf_target = [
    {
      community = "target:1:10002"
    }
  ]
  common_g_84da76_groups_leaf_routing_instances_instance_VRF_10003_interface = [
    {
      name = "irb.3001"
    },
    {
      name = "lo0.10003"
    }
  ]
  common_g_84da76_groups_leaf_routing_instances_instance_VRF_10003_protocols = [
    {
      evpn = [
        {
          ip_prefix_routes = [
            {
              advertise = "direct-nexthop"
              encapsulation = "vxlan"
              export = [
                "EVPN_T5_EXPORT"
              ]
              vni = 10003
            }
          ]
        }
      ]
    }
  ]
  common_g_84da76_groups_leaf_routing_instances_instance_VRF_10003_routing_options = [
    {
      auto_export = [
        {

        }
      ]
    }
  ]
  common_g_84da76_groups_leaf_routing_instances_instance_VRF_10003_vrf_table_label = [
    {

    }
  ]
  common_g_84da76_groups_leaf_routing_instances_instance_VRF_10003_vrf_target = [
    {
      community = "target:1:10003"
    }
  ]
  common_g_84da76_groups_leaf_routing_options_forwarding_table = [
    {
      chained_composite_next_hop = [
        {
          ingress = [
            {
              evpn = ""
            }
          ]
        }
      ]
      ecmp_fast_reroute = ""
      export = [
        "PFE-LB"
      ]
    }
  ]
  common_g_84da76_groups_leaf_routing_options_static = [
    {
      route = [
        {
          name = "0.0.0.0/0"
          next_hop = [
            "100.123.0.1"
          ]
        }
      ]
    }
  ]
  common_g_84da76_groups_leaf_snmp = [
    {
      community = [
        {
          authorization = "read-only"
          name = "public"
        }
      ]
      contact = "aburston@juniper.net"
      location = "JCL Labs"
    }
  ]
  common_g_84da76_groups_leaf_switch_options_vrf_target = [
    {
      auto = [
        {

        }
      ]
      community = "target:9999:9999"
    }
  ]
  common_g_84da76_groups_leaf_switch_options_vtep_source_interface = [
    {
      interface_name = "lo0.0"
    }
  ]
  common_g_84da76_groups_leaf_system = [
    {
      extensions = [
        {
          providers = [
            {
              license_type = [
                {
                  deployment_scope = [
                    "commercial"
                  ]
                  name = "juniper"
                }
              ]
              name = "juniper"
            },
            {
              license_type = [
                {
                  deployment_scope = [
                    "commercial"
                  ]
                  name = "juniper"
                }
              ]
              name = "chef"
            }
          ]
        }
      ]
      login = [
        {
          message = "***********************************************************************\nThis system is restricted to __________, authorized users for legitimate\nbusiness purposes only. All activity on the system will be logged and\nis subject to monitoring. Unauthorized access, use or modification\nof computers, data therein or data in transit to or from the computers\nis a violation of state and federal laws. Unauthorized activity will\nbe reported to the law enforcement for investigation and possible\nprosecution. __________ reserves the right to investigate, refer for\nprosecution and pursue monetary damages in civil actions in the event\nof unauthorized access.\n***********************************************************************\n"
          user = [
            {
              authentication = [
                {
                  encrypted_password = "$1$a31gJmWG$h9ohikT1ajySf/tVH.gmv1"
                }
              ]
              class = "super-user"
              name = "jcluser"
              uid = 2000
            }
          ]
        }
      ]
      root_authentication = [
        {
          encrypted_password = "$1$DbZ1Q3pj$s48cZytjsmSJRUJAf4LdM."
        }
      ]
      services = [
        {
          extension_service = [
            {
              notification = [
                {
                  allow_clients = [
                    {
                      address = [
                        "0.0.0.0/0"
                      ]
                    }
                  ]
                }
              ]
              request_response = [
                {
                  grpc = [
                    {
                      max_connections = 30
                    }
                  ]
                }
              ]
            }
          ]
          netconf = [
            {
              ssh = [
                {

                }
              ]
            }
          ]
          rest = [
            {
              enable_explorer = ""
              http = [
                {
                  port = 3000
                }
              ]
            }
          ]
          ssh = [
            {
              root_login = "allow"
            }
          ]
        }
      ]
      syslog = [
        {
          file = [
            {
              contents = [
                {
                  name = "any"
                  notice = ""
                },
                {
                  info = ""
                  name = "authorization"
                }
              ]
              name = "messages"
            },
            {
              contents = [
                {
                  any = ""
                  name = "interactive-commands"
                }
              ]
              name = "interactive-commands"
            }
          ]
          user = [
            {
              contents = [
                {
                  emergency = ""
                  name = "any"
                }
              ]
              name = "*"
            }
          ]
        }
      ]
    }
  ]
  common_g_84da76_groups_leaf_vlans = [
    {
      vlan = [
        {
          l3_interface = "irb.1001"
          name = "vlan_1001"
          vlan_id = 1001
          vxlan = [
            {
              vni = 1001
            }
          ]
        },
        {
          l3_interface = "irb.2001"
          name = "vlan_2001"
          vlan_id = 2001
          vxlan = [
            {
              vni = 2001
            }
          ]
        },
        {
          l3_interface = "irb.3001"
          name = "vlan_3001"
          vlan_id = 3001
          vxlan = [
            {
              vni = 3001
            }
          ]
        }
      ]
    }
  ]

  # ── Group 9a472e: dc1-spine1, dc1-spine2 ────────────────────────────────────────────────
  common_g_9a472e_groups_spine_protocols_bgp_group_EVPN_iBGP_allow = [
    "10.30.100.0/24"
  ]
  common_g_9a472e_groups_spine_protocols_bgp_group_EVPN_iBGP_local_as = [
    {
      as_number = 65200
    }
  ]
  common_g_9a472e_groups_spine_routing_options_forwarding_table = [
    {
      ecmp_fast_reroute = ""
      export = [
        "PFE-LB"
      ]
    }
  ]

  # ── Group a654ee: dc1-borderleaf1, dc1-borderleaf2 ────────────────────────────────────────────────
  common_g_a654ee_apply_groups = [
    "borderleaf"
  ]
  common_g_a654ee_groups_borderleaf_chassis = [
    {
      aggregated_devices = [
        {
          ethernet = [
            {
              device_count = 24
            }
          ]
        }
      ]
    }
  ]
  common_g_a654ee_groups_borderleaf_forwarding_options = [
    {
      storm_control_profiles = [
        {
          all = [
            {

            }
          ]
          name = "default"
        }
      ]
    }
  ]
  common_g_a654ee_groups_borderleaf_interfaces_interface_em1 = {
    name = "em1"
    unit = [
      {
        description = "*** to pfe ***"
        family = [
          {
            inet = [
              {
                address = [
                  {
                    name = "169.254.0.2/24"
                  }
                ]
              }
            ]
          }
        ]
        name = 0
      }
    ]
  }
  common_g_a654ee_groups_borderleaf_policy_options_policy_statement_EVPN_T5_EXPORT = {
    name = "EVPN_T5_EXPORT"
    term = [
      {
        from = [
          {
            protocol = [
              "direct"
            ]
          }
        ]
        name = "fm_direct"
        then = [
          {
            accept = ""
          }
        ]
      },
      {
        from = [
          {
            protocol = [
              "static"
            ]
          }
        ]
        name = "fm_static"
        then = [
          {
            accept = ""
          }
        ]
      },
      {
        from = [
          {
            protocol = [
              "evpn",
              "ospf"
            ]
            route_filter = [
              {
                address = "0.0.0.0/0"
                exact = ""
              }
            ]
          }
        ]
        name = "fm_v4_default"
        then = [
          {
            accept = ""
          }
        ]
      },
      {
        from = [
          {
            protocol = [
              "evpn"
            ]
            route_filter = [
              {
                address = "0::0/0"
                prefix_length_range = "/128-/128"
              }
            ]
          }
        ]
        name = "fm_v6_host"
        then = [
          {
            accept = ""
          }
        ]
      }
    ]
  }
  common_g_a654ee_groups_borderleaf_policy_options_policy_statement_IPCLOS_BGP_EXP_term_default = {
    name = "default"
    then = [
      {
        reject = ""
      }
    ]
  }
  common_g_a654ee_groups_borderleaf_policy_options_policy_statement_IPCLOS_BGP_EXP_term_loopback_from = [
    {
      protocol = [
        "direct",
        "bgp"
      ]
    }
  ]
  common_g_a654ee_groups_borderleaf_policy_options_policy_statement_IPCLOS_BGP_EXP_term_loopback_then_accept = ""
  common_g_a654ee_groups_borderleaf_policy_options_policy_statement_IPCLOS_BGP_EXP_term_loopback_then_community_add = ""
  common_g_a654ee_groups_borderleaf_policy_options_policy_statement_IPCLOS_BGP_IMP = {
    name = "IPCLOS_BGP_IMP"
    term = [
      {
        from = [
          {
            protocol = [
              "bgp",
              "direct"
            ]
          }
        ]
        name = "loopback"
        then = [
          {
            accept = ""
          }
        ]
      },
      {
        name = "default"
        then = [
          {
            reject = ""
          }
        ]
      }
    ]
  }
  common_g_a654ee_groups_borderleaf_policy_options_policy_statement_PFE_LB = {
    name = "PFE-LB"
    then = [
      {
        load_balance = [
          {
            per_packet = ""
          }
        ]
      }
    ]
  }
  common_g_a654ee_groups_borderleaf_policy_options_policy_statement_to_ospf = {
    name = "to-ospf"
    term = [
      {
        from = [
          {
            protocol = [
              "evpn"
            ]
          }
        ]
        name = 10
        then = [
          {
            accept = ""
          }
        ]
      },
      {
        name = 100
        then = [
          {
            reject = ""
          }
        ]
      }
    ]
  }
  common_g_a654ee_groups_borderleaf_protocols_bgp_group_EVPN_iBGP_family = [
    {
      evpn = [
        {
          signaling = [
            {

            }
          ]
        }
      ]
    }
  ]
  common_g_a654ee_groups_borderleaf_protocols_bgp_group_EVPN_iBGP_local_as = [
    {
      as_number = 65200
    }
  ]
  common_g_a654ee_groups_borderleaf_protocols_bgp_group_EVPN_iBGP_multipath = [
    {

    }
  ]
  common_g_a654ee_groups_borderleaf_protocols_bgp_group_EVPN_iBGP_neighbor = [
    {
      name = "10.30.100.3"
    },
    {
      name = "10.30.100.4"
    }
  ]
  common_g_a654ee_groups_borderleaf_protocols_bgp_group_IPCLOS_eBGP_bfd_liveness_detection = [
    {
      minimum_interval = 1000
      multiplier = 3
    }
  ]
  common_g_a654ee_groups_borderleaf_protocols_bgp_group_IPCLOS_eBGP_export = [
    "IPCLOS_BGP_EXP"
  ]
  common_g_a654ee_groups_borderleaf_protocols_bgp_group_IPCLOS_eBGP_import = [
    "IPCLOS_BGP_IMP"
  ]
  common_g_a654ee_groups_borderleaf_protocols_bgp_group_IPCLOS_eBGP_multipath = [
    {
      multiple_as = ""
    }
  ]
  common_g_a654ee_groups_borderleaf_protocols_bgp_group_WAN_OVERLAY_eBGP_family = [
    {
      evpn = [
        {
          signaling = [
            {
              delay_route_advertisements = [
                {
                  minimum_delay = [
                    {
                      routing_uptime = 480
                    }
                  ]
                }
              ]
            }
          ]
        }
      ]
    }
  ]
  common_g_a654ee_groups_borderleaf_protocols_bgp_group_WAN_OVERLAY_eBGP_local_as = [
    {
      as_number = 65200
    }
  ]
  common_g_a654ee_groups_borderleaf_protocols_bgp_group_WAN_OVERLAY_eBGP_multihop = [
    {
      no_nexthop_change = ""
    }
  ]
  common_g_a654ee_groups_borderleaf_protocols_bgp_group_WAN_OVERLAY_eBGP_multipath = [
    {
      multiple_as = ""
    }
  ]
  common_g_a654ee_groups_borderleaf_protocols_bgp_group_WAN_OVERLAY_eBGP_neighbor = [
    {
      description = "DCI EBGP peering to 10.30.100.8"
      name = "10.30.100.8"
      peer_as = 65201
    },
    {
      description = "DCI EBGP peering to 10.30.100.9"
      name = "10.30.100.9"
      peer_as = 65201
    }
  ]
  common_g_a654ee_groups_borderleaf_protocols_evpn = [
    {
      default_gateway = "do-not-advertise"
      encapsulation = "vxlan"
      multicast_mode = "ingress-replication"
    }
  ]
  common_g_a654ee_groups_borderleaf_protocols_igmp_snooping = [
    {
      vlan = [
        {
          name = "default"
        }
      ]
    }
  ]
  common_g_a654ee_groups_borderleaf_protocols_lldp = [
    {
      interface = [
        {
          name = "all"
        }
      ]
    }
  ]
  common_g_a654ee_groups_borderleaf_routing_instances_instance_VRF_10001_interface = [
    {
      name = "xe-0/0/2.1"
    },
    {
      name = "xe-0/0/3.1"
    },
    {
      name = "lo0.10001"
    }
  ]
  common_g_a654ee_groups_borderleaf_routing_instances_instance_VRF_10001_protocols = [
    {
      evpn = [
        {
          ip_prefix_routes = [
            {
              advertise = "direct-nexthop"
              encapsulation = "vxlan"
              export = [
                "EVPN_T5_EXPORT"
              ]
              vni = 10001
            }
          ]
        }
      ]
      ospf = [
        {
          area = [
            {
              interface = [
                {
                  metric = 100
                  name = "xe-0/0/2.1"
                },
                {
                  metric = 200
                  name = "xe-0/0/3.1"
                }
              ]
              name = "0.0.0.0"
            }
          ]
          export = [
            "to-ospf"
          ]
        }
      ]
    }
  ]
  common_g_a654ee_groups_borderleaf_routing_instances_instance_VRF_10001_routing_options = [
    {
      auto_export = [
        {

        }
      ]
    }
  ]
  common_g_a654ee_groups_borderleaf_routing_instances_instance_VRF_10001_vrf_table_label = [
    {

    }
  ]
  common_g_a654ee_groups_borderleaf_routing_instances_instance_VRF_10001_vrf_target = [
    {
      community = "target:1:10001"
    }
  ]
  common_g_a654ee_groups_borderleaf_routing_instances_instance_VRF_10002_interface = [
    {
      name = "xe-0/0/2.2"
    },
    {
      name = "xe-0/0/3.2"
    },
    {
      name = "lo0.10002"
    }
  ]
  common_g_a654ee_groups_borderleaf_routing_instances_instance_VRF_10002_protocols = [
    {
      evpn = [
        {
          ip_prefix_routes = [
            {
              advertise = "direct-nexthop"
              encapsulation = "vxlan"
              export = [
                "EVPN_T5_EXPORT"
              ]
              vni = 10002
            }
          ]
        }
      ]
      ospf = [
        {
          area = [
            {
              interface = [
                {
                  metric = 100
                  name = "xe-0/0/2.2"
                },
                {
                  metric = 200
                  name = "xe-0/0/3.2"
                }
              ]
              name = "0.0.0.0"
            }
          ]
          export = [
            "to-ospf"
          ]
        }
      ]
    }
  ]
  common_g_a654ee_groups_borderleaf_routing_instances_instance_VRF_10002_routing_options = [
    {
      auto_export = [
        {

        }
      ]
    }
  ]
  common_g_a654ee_groups_borderleaf_routing_instances_instance_VRF_10002_vrf_table_label = [
    {

    }
  ]
  common_g_a654ee_groups_borderleaf_routing_instances_instance_VRF_10002_vrf_target = [
    {
      community = "target:1:10002"
    }
  ]
  common_g_a654ee_groups_borderleaf_routing_instances_instance_VRF_10003_interface = [
    {
      name = "xe-0/0/2.3"
    },
    {
      name = "xe-0/0/3.3"
    },
    {
      name = "lo0.10003"
    }
  ]
  common_g_a654ee_groups_borderleaf_routing_instances_instance_VRF_10003_protocols = [
    {
      evpn = [
        {
          ip_prefix_routes = [
            {
              advertise = "direct-nexthop"
              encapsulation = "vxlan"
              export = [
                "EVPN_T5_EXPORT"
              ]
              vni = 10003
            }
          ]
        }
      ]
      ospf = [
        {
          area = [
            {
              interface = [
                {
                  metric = 100
                  name = "xe-0/0/2.3"
                },
                {
                  metric = 200
                  name = "xe-0/0/3.3"
                }
              ]
              name = "0.0.0.0"
            }
          ]
          export = [
            "to-ospf"
          ]
        }
      ]
    }
  ]
  common_g_a654ee_groups_borderleaf_routing_instances_instance_VRF_10003_routing_options = [
    {
      auto_export = [
        {

        }
      ]
    }
  ]
  common_g_a654ee_groups_borderleaf_routing_instances_instance_VRF_10003_vrf_table_label = [
    {

    }
  ]
  common_g_a654ee_groups_borderleaf_routing_instances_instance_VRF_10003_vrf_target = [
    {
      community = "target:1:10003"
    }
  ]
  common_g_a654ee_groups_borderleaf_routing_options_forwarding_table = [
    {
      chained_composite_next_hop = [
        {
          ingress = [
            {
              evpn = ""
            }
          ]
        }
      ]
      ecmp_fast_reroute = ""
      export = [
        "PFE-LB"
      ]
    }
  ]
  common_g_a654ee_groups_borderleaf_routing_options_static = [
    {
      route = [
        {
          name = "0.0.0.0/0"
          next_hop = [
            "100.123.0.1"
          ]
        }
      ]
    }
  ]
  common_g_a654ee_groups_borderleaf_snmp = [
    {
      community = [
        {
          authorization = "read-only"
          name = "public"
        }
      ]
      contact = "aburston@juniper.net"
      location = "JCL Labs"
    }
  ]
  common_g_a654ee_groups_borderleaf_switch_options_vrf_target = [
    {
      auto = [
        {

        }
      ]
      community = "target:9999:9999"
    }
  ]
  common_g_a654ee_groups_borderleaf_switch_options_vtep_source_interface = [
    {
      interface_name = "lo0.0"
    }
  ]
  common_g_a654ee_groups_borderleaf_system = [
    {
      extensions = [
        {
          providers = [
            {
              license_type = [
                {
                  deployment_scope = [
                    "commercial"
                  ]
                  name = "juniper"
                }
              ]
              name = "juniper"
            },
            {
              license_type = [
                {
                  deployment_scope = [
                    "commercial"
                  ]
                  name = "juniper"
                }
              ]
              name = "chef"
            }
          ]
        }
      ]
      login = [
        {
          message = "***********************************************************************\nThis system is restricted to __________, authorized users for legitimate\nbusiness purposes only. All activity on the system will be logged and\nis subject to monitoring. Unauthorized access, use or modification\nof computers, data therein or data in transit to or from the computers\nis a violation of state and federal laws. Unauthorized activity will\nbe reported to the law enforcement for investigation and possible\nprosecution. __________ reserves the right to investigate, refer for\nprosecution and pursue monetary damages in civil actions in the event\nof unauthorized access.\n***********************************************************************\n"
          user = [
            {
              authentication = [
                {
                  encrypted_password = "$1$a31gJmWG$h9ohikT1ajySf/tVH.gmv1"
                }
              ]
              class = "super-user"
              name = "jcluser"
              uid = 2000
            }
          ]
        }
      ]
      root_authentication = [
        {
          encrypted_password = "$1$DbZ1Q3pj$s48cZytjsmSJRUJAf4LdM."
        }
      ]
      services = [
        {
          extension_service = [
            {
              notification = [
                {
                  allow_clients = [
                    {
                      address = [
                        "0.0.0.0/0"
                      ]
                    }
                  ]
                }
              ]
              request_response = [
                {
                  grpc = [
                    {
                      max_connections = 30
                    }
                  ]
                }
              ]
            }
          ]
          netconf = [
            {
              ssh = [
                {

                }
              ]
            }
          ]
          rest = [
            {
              enable_explorer = ""
              http = [
                {
                  port = 3000
                }
              ]
            }
          ]
          ssh = [
            {
              root_login = "allow"
            }
          ]
        }
      ]
      syslog = [
        {
          file = [
            {
              contents = [
                {
                  name = "any"
                  notice = ""
                },
                {
                  info = ""
                  name = "authorization"
                }
              ]
              name = "messages"
            },
            {
              contents = [
                {
                  any = ""
                  name = "interactive-commands"
                }
              ]
              name = "interactive-commands"
            }
          ]
          user = [
            {
              contents = [
                {
                  emergency = ""
                  name = "any"
                }
              ]
              name = "*"
            }
          ]
        }
      ]
    }
  ]

  # ── Group fb64d2: dc1-leaf1, dc1-leaf2 ────────────────────────────────────────────────
  common_g_fb64d2_groups_leaf_interfaces_interface_ae0 = {
    aggregated_ether_options = [
      {
        lacp = [
          {
            active = ""
            periodic = "fast"
            system_id = "00:00:00:00:01:00"
          }
        ]
      }
    ]
    esi = [
      {
        all_active = ""
        identifier = "00:00:00:00:00:00:00:00:01:00"
      }
    ]
    name = "ae0"
    unit = [
      {
        family = [
          {
            ethernet_switching = [
              {
                vlan = [
                  {
                    members = [
                      1001
                    ]
                  }
                ]
              }
            ]
          }
        ]
        name = 0
      }
    ]
  }
  common_g_fb64d2_groups_leaf_interfaces_interface_ae1 = {
    aggregated_ether_options = [
      {
        lacp = [
          {
            active = ""
            periodic = "fast"
            system_id = "00:00:00:00:02:00"
          }
        ]
      }
    ]
    esi = [
      {
        all_active = ""
        identifier = "00:00:00:00:00:00:00:00:02:00"
      }
    ]
    name = "ae1"
    unit = [
      {
        family = [
          {
            ethernet_switching = [
              {
                vlan = [
                  {
                    members = [
                      2001
                    ]
                  }
                ]
              }
            ]
          }
        ]
        name = 0
      }
    ]
  }
  common_g_fb64d2_groups_leaf_interfaces_interface_xe_0_0_3 = {
    ether_options = [
      {
        ieee_802_3ad = [
          {
            bundle = "ae1"
          }
        ]
      }
    ]
    name = "xe-0/0/3"
  }
}
