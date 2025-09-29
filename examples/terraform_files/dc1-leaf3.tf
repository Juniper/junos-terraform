resource "terraform-provider-junos-vqfx-evpn-vxlan" "dc1-leaf3-base-config" {
  resource_name = "base-config"
  provider = junos-vqfx-evpn-vxlan.dc1_leaf3
  system = [
    {
      host_name = "dc1-leaf3"
      root_authentication = [
        {
          encrypted_password = "$1$DbZ1Q3pj$s48cZytjsmSJRUJAf4LdM."
        }
      ]
      login = [
        {
          message = "***********************************************************************\nThis system is restricted to __________, authorized users for legitimate\nbusiness purposes only. All activity on the system will be logged and\nis subject to monitoring. Unauthorized access, use or modification\nof computers, data therein or data in transit to or from the computers\nis a violation of state and federal laws. Unauthorized activity will\nbe reported to the law enforcement for investigation and possible\nprosecution. __________ reserves the right to investigate, refer for\nprosecution and pursue monetary damages in civil actions in the event\nof unauthorized access.\n***********************************************************************\n"
          user = [
            {
              name = "jcluser"
              uid = 2000
              class = "super-user"
              authentication = [
                {
                  encrypted_password = "$1$a31gJmWG$h9ohikT1ajySf/tVH.gmv1"
                }
              ]
            }
          ]
        }
      ]
      services = [
        {
          ssh = [
            {
              root_login = "allow"
            }
          ]
          extension_service = [
            {
              request_response = [
                {
                  grpc = [
                    {
                      undocumented = [
                        {
                          clear_text = [
                            {
                              address = "0.0.0.0"
                              port = 32767
                            }
                          ]
                        },
                        {
                          skip_authentication = ""
                        }
                      ]
                      max_connections = 30
                    }
                  ]
                }
              ]
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
              http = [
                {
                  port = 3000
                }
              ]
              enable_explorer = ""
            }
          ]
        }
      ]
      syslog = [
        {
          user = [
            {
              name = "*"
              contents = [
                {
                  name = "any"
                  emergency = ""
                }
              ]
            }
          ]
          file = [
            {
              name = "messages"
              contents = [
                {
                  name = "any"
                  notice = ""
                },
                {
                  name = "authorization"
                  info = ""
                }
              ]
            },
            {
              name = "interactive-commands"
              contents = [
                {
                  name = "interactive-commands"
                  any = ""
                }
              ]
            }
          ]
        }
      ]
      extensions = [
        {
          providers = [
            {
              name = "juniper"
              license_type = [
                {
                  name = "juniper"
                  deployment_scope = [
                    "commercial"
                  ]
                }
              ]
            },
            {
              name = "chef"
              license_type = [
                {
                  name = "juniper"
                  deployment_scope = [
                    "commercial"
                  ]
                }
              ]
            }
          ]
        }
      ]
    }
  ]
  chassis = [
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
  interfaces = [
    {
      interface = [
        {
          name = "xe-0/0/0"
          description = "*** to dc1-spine1 ***"
          unit = [
            {
              name = 0
              family = [
                {
                  inet = [
                    {
                      address = [
                        {
                          name = "10.30.137.2/30"
                        }
                      ]
                    }
                  ]
                }
              ]
            }
          ]
        },
        {
          name = "xe-0/0/1"
          description = "*** to dc1-spine2 ***"
          unit = [
            {
              name = 0
              family = [
                {
                  inet = [
                    {
                      address = [
                        {
                          name = "10.30.147.2/30"
                        }
                      ]
                    }
                  ]
                }
              ]
            }
          ]
        },
        {
          name = "xe-0/0/2"
          ether_options = [
            {
              ieee_802_3ad = [
                {
                  bundle = "ae0"
                }
              ]
            }
          ]
        },
        {
          name = "ae0"
          esi = [
            {
              identifier = "00:00:00:00:00:00:00:00:03:00"
              all_active = ""
            }
          ]
          aggregated_ether_options = [
            {
              lacp = [
                {
                  active = ""
                  periodic = "fast"
                  system_id = "00:00:00:00:03:00"
                }
              ]
            }
          ]
          unit = [
            {
              name = 0
              family = [
                {
                  ethernet_switching = [
                    {
                      vlan = [
                        {
                          members = [
                            3001
                          ]
                        }
                      ]
                    }
                  ]
                }
              ]
            }
          ]
        },
        {
          name = "em0"
          unit = [
            {
              name = 0
              description = "*** management ***"
              family = [
                {
                  inet = [
                    {
                      address = [
                        {
                          name = "100.123.24.7/16"
                        }
                      ]
                    }
                  ]
                }
              ]
            }
          ]
        },
        {
          name = "em1"
          unit = [
            {
              name = 0
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
            }
          ]
        },
        {
          name = "irb"
          unit = [
            {
              name = 1001
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
            },
            {
              name = 2001
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
            },
            {
              name = 3001
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
            }
          ]
        },
        {
          name = "lo0"
          unit = [
            {
              name = 0
              description = "*** loopback ***"
              family = [
                {
                  inet = [
                    {
                      address = [
                        {
                          name = "10.30.100.7/32"
                        }
                      ]
                    }
                  ]
                }
              ]
            },
            {
              name = 10001
              description = "Loopback for VXLAN control packets for VRF_10001"
              family = [
                {
                  inet = [
                    {
                      address = [
                        {
                          name = "10.40.100.5/32"
                        }
                      ]
                    }
                  ]
                }
              ]
            },
            {
              name = 10002
              description = "Loopback for VXLAN control packets for VRF_10002"
              family = [
                {
                  inet = [
                    {
                      address = [
                        {
                          name = "10.40.100.10/32"
                        }
                      ]
                    }
                  ]
                }
              ]
            },
            {
              name = 10003
              description = "Loopback for VXLAN control packets for VRF_10003"
              family = [
                {
                  inet = [
                    {
                      address = [
                        {
                          name = "10.40.100.15/32"
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
    }
  ]
  snmp = [
    {
      location = "JCL Labs"
      contact = "aburston@juniper.net"
      community = [
        {
          name = "public"
          authorization = "read-only"
        }
      ]
    }
  ]
  forwarding_options = [
    {
      storm_control_profiles = [
        {
          name = "default"
          all = [
            {

            }
          ]
        }
      ]
    }
  ]
  routing_options = [
    {
      static = [
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
      router_id = "10.30.100.7"
      forwarding_table = [
        {
          export = [
            "PFE-LB"
          ]
          ecmp_fast_reroute = ""
          chained_composite_next_hop = [
            {
              ingress = [
                {
                  evpn = ""
                }
              ]
            }
          ]
        }
      ]
    }
  ]
  protocols = [
    {
      bgp = [
        {
          group = [
            {
              name = "EVPN_iBGP"
              type = "internal"
              local_address = "10.30.100.7"
              family = [
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
              cluster = "10.30.100.7"
              local_as = [
                {
                  as_number = 65200
                }
              ]
              multipath = [
                {

                }
              ]
              neighbor = [
                {
                  name = "10.30.100.3"
                },
                {
                  name = "10.30.100.4"
                }
              ]
            },
            {
              name = "IPCLOS_eBGP"
              type = "external"
              mtu_discovery = ""
              import = [
                "IPCLOS_BGP_IMP"
              ]
              export = [
                "IPCLOS_BGP_EXP"
              ]
              vpn_apply_export = ""
              local_as = [
                {
                  as_number = 65505
                }
              ]
              multipath = [
                {
                  multiple_as = ""
                }
              ]
              bfd_liveness_detection = [
                {
                  minimum_interval = 1000
                  multiplier = 3
                }
              ]
              neighbor = [
                {
                  name = "10.30.137.1"
                  description = "EBGP peering to 10.30.137.1"
                  peer_as = 65501
                },
                {
                  name = "10.30.147.1"
                  description = "EBGP peering to 10.30.147.1"
                  peer_as = 65502
                }
              ]
            }
          ]
        }
      ]
      evpn = [
        {
          encapsulation = "vxlan"
          multicast_mode = "ingress-replication"
          default_gateway = "do-not-advertise"
          extended_vni_list = [
            "all"
          ]
        }
      ]
      lldp = [
        {
          interface = [
            {
              name = "all"
            }
          ]
        }
      ]
      igmp_snooping = [
        {
          vlan = [
            {
              name = "default"
            }
          ]
        }
      ]
    }
  ]
  policy_options = [
    {
      policy_statement = [
        {
          name = "EVPN_T5_EXPORT"
          term = [
            {
              name = "fm_direct"
              from = [
                {
                  protocol = [
                    "direct"
                  ]
                }
              ]
              then = [
                {
                  accept = ""
                }
              ]
            },
            {
              name = "fm_static"
              from = [
                {
                  protocol = [
                    "static"
                  ]
                }
              ]
              then = [
                {
                  accept = ""
                }
              ]
            },
            {
              name = "fm_v4_host"
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
              then = [
                {
                  accept = ""
                }
              ]
            },
            {
              name = "fm_v6_host"
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
              then = [
                {
                  accept = ""
                }
              ]
            }
          ]
        },
        {
          name = "IPCLOS_BGP_EXP"
          term = [
            {
              name = "loopback"
              from = [
                {
                  protocol = [
                    "direct",
                    "bgp"
                  ]
                }
              ]
              then = [
                {
                  community = [
                    {
                      add = ""
                      community_name = "dc1-leaf3"
                    }
                  ]
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
        },
        {
          name = "IPCLOS_BGP_IMP"
          term = [
            {
              name = "loopback"
              from = [
                {
                  protocol = [
                    "bgp",
                    "direct"
                  ]
                }
              ]
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
        },
        {
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
      ]
      community = [
        {
          name = "dc1-leaf3"
          members = [
            "65505:1"
          ]
        }
      ]
    }
  ]
  routing_instances = [
    {
      instance = [
        {
          name = "VRF_10001"
          instance_type = "vrf"
          interface = [
            {
              name = "irb.1001"
            },
            {
              name = "lo0.10001"
            }
          ]
          route_distinguisher = [
            {
              rd_type = "10.40.100.5:10001"
            }
          ]
          vrf_target = [
            {
              community = "target:1:10001"
            }
          ]
          vrf_table_label = [
            {

            }
          ]
          routing_options = [
            {
              auto_export = [
                {

                }
              ]
            }
          ]
          protocols = [
            {
              evpn = [
                {
                  ip_prefix_routes = [
                    {
                      advertise = "direct-nexthop"
                      encapsulation = "vxlan"
                      vni = 10001
                      export = [
                        "EVPN_T5_EXPORT"
                      ]
                    }
                  ]
                }
              ]
            }
          ]
        },
        {
          name = "VRF_10002"
          instance_type = "vrf"
          interface = [
            {
              name = "irb.2001"
            },
            {
              name = "lo0.10002"
            }
          ]
          route_distinguisher = [
            {
              rd_type = "10.40.100.10:10002"
            }
          ]
          vrf_target = [
            {
              community = "target:1:10002"
            }
          ]
          vrf_table_label = [
            {

            }
          ]
          routing_options = [
            {
              auto_export = [
                {

                }
              ]
            }
          ]
          protocols = [
            {
              evpn = [
                {
                  ip_prefix_routes = [
                    {
                      advertise = "direct-nexthop"
                      encapsulation = "vxlan"
                      vni = 10002
                      export = [
                        "EVPN_T5_EXPORT"
                      ]
                    }
                  ]
                }
              ]
            }
          ]
        },
        {
          name = "VRF_10003"
          instance_type = "vrf"
          interface = [
            {
              name = "irb.3001"
            },
            {
              name = "lo0.10003"
            }
          ]
          route_distinguisher = [
            {
              rd_type = "10.40.100.15:10003"
            }
          ]
          vrf_target = [
            {
              community = "target:1:10003"
            }
          ]
          vrf_table_label = [
            {

            }
          ]
          routing_options = [
            {
              auto_export = [
                {

                }
              ]
            }
          ]
          protocols = [
            {
              evpn = [
                {
                  ip_prefix_routes = [
                    {
                      advertise = "direct-nexthop"
                      encapsulation = "vxlan"
                      vni = 10003
                      export = [
                        "EVPN_T5_EXPORT"
                      ]
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
  switch_options = [
    {
      vtep_source_interface = [
        {
          interface_name = "lo0.0"
        }
      ]
      route_distinguisher = [
        {
          rd_type = "10.30.100.7:9999"
        }
      ]
      vrf_target = [
        {
          community = "target:9999:9999"
          auto = [
            {

            }
          ]
        }
      ]
    }
  ]
  vlans = [
    {
      vlan = [
        {
          name = "vlan_1001"
          vlan_id = 1001
          l3_interface = "irb.1001"
          vxlan = [
            {
              vni = 1001
            }
          ]
        },
        {
          name = "vlan_2001"
          vlan_id = 2001
          l3_interface = "irb.2001"
          vxlan = [
            {
              vni = 2001
            }
          ]
        },
        {
          name = "vlan_3001"
          vlan_id = 3001
          l3_interface = "irb.3001"
          vxlan = [
            {
              vni = 3001
            }
          ]
        }
      ]
    }
  ]
}
