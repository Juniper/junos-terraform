
provider "junos-qfx-evpn-vxlan" {
    host     = "dc1-borderleaf2"
    port     = 22
    username = "jcluser"
    password = "Juniper!1"
    alias    = "dc1_borderleaf2"
}

resource "junos-qfx-evpn-vxlan-base-config" "dc1-borderleaf2" {
  resource_name = "dc1-borderleaf2"
  provider = junos-qfx-evpn-vxlan.dc1_borderleaf2
  system = [
    {
      host_name = "dc1-borderleaf2"
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
                          name = "10.30.132.2/30"
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
                          name = "10.30.142.2/30"
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
          vlan_tagging = ""
          unit = [
            {
              name = 1
              vlan_id = 1
              family = [
                {
                  inet = [
                    {
                      address = [
                        {
                          name = "10.97.1.1/30"
                        }
                      ]
                    }
                  ]
                }
              ]
            },
            {
              name = 2
              vlan_id = 2
              family = [
                {
                  inet = [
                    {
                      address = [
                        {
                          name = "10.97.2.1/30"
                        }
                      ]
                    }
                  ]
                }
              ]
            },
            {
              name = 3
              vlan_id = 3
              family = [
                {
                  inet = [
                    {
                      address = [
                        {
                          name = "10.97.3.1/30"
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
          name = "xe-0/0/3"
          vlan_tagging = ""
          unit = [
            {
              name = 1
              vlan_id = 1
              family = [
                {
                  inet = [
                    {
                      address = [
                        {
                          name = "10.96.1.1/30"
                        }
                      ]
                    }
                  ]
                }
              ]
            },
            {
              name = 2
              vlan_id = 2
              family = [
                {
                  inet = [
                    {
                      address = [
                        {
                          name = "10.96.2.1/30"
                        }
                      ]
                    }
                  ]
                }
              ]
            },
            {
              name = 3
              vlan_id = 3
              family = [
                {
                  inet = [
                    {
                      address = [
                        {
                          name = "10.96.3.1/30"
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
          name = "xe-0/0/4"
          description = "*** to wan-pe1 ***"
          unit = [
            {
              name = 0
              family = [
                {
                  inet = [
                    {
                      address = [
                        {
                          name = "10.32.10.1/30"
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
                          name = "100.123.24.2/16"
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
                          name = "10.30.100.2/32"
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
                          name = "10.40.100.2/32"
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
                          name = "10.40.100.7/32"
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
                          name = "10.40.100.12/32"
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
      router_id = "10.30.100.2"
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
              name = "WAN_OVERLAY_eBGP"
              type = "external"
              multihop = [
                {
                  no_nexthop_change = ""
                }
              ]
              local_address = "10.30.100.2"
              family = [
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
              local_as = [
                {
                  as_number = 65200
                }
              ]
              multipath = [
                {
                  multiple_as = ""
                }
              ]
              neighbor = [
                {
                  name = "10.30.100.8"
                  description = "DCI EBGP peering to 10.30.100.8"
                  peer_as = 65201
                },
                {
                  name = "10.30.100.9"
                  description = "DCI EBGP peering to 10.30.100.9"
                  peer_as = 65201
                }
              ]
            },
            {
              name = "EVPN_iBGP"
              type = "internal"
              local_address = "10.30.100.2"
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
              cluster = "10.30.100.2"
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
                  as_number = 65507
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
                  name = "10.30.132.1"
                  description = "EBGP peering to 10.30.132.1"
                  peer_as = 65501
                },
                {
                  name = "10.30.142.1"
                  description = "EBGP peering to 10.30.142.1"
                  peer_as = 65502
                },
                {
                  name = "10.32.10.2"
                  description = "EBGP peering to 10.32.10.2"
                  peer_as = 65400
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
              name = "fm_v4_default"
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
                      community_name = "dc1-borderleaf2"
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
        },
        {
          name = "to-ospf"
          term = [
            {
              name = 10
              from = [
                {
                  protocol = [
                    "evpn"
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
              name = 100
              then = [
                {
                  reject = ""
                }
              ]
            }
          ]
        }
      ]
      community = [
        {
          name = "dc1-borderleaf2"
          members = [
            "65507:1"
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
              name = "xe-0/0/2.1"
            },
            {
              name = "xe-0/0/3.1"
            },
            {
              name = "lo0.10001"
            }
          ]
          route_distinguisher = [
            {
              rd_type = "10.40.100.2:10001"
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
              ospf = [
                {
                  export = [
                    "to-ospf"
                  ]
                  area = [
                    {
                      name = "0.0.0.0"
                      interface = [
                        {
                          name = "xe-0/0/2.1"
                          metric = 100
                        },
                        {
                          name = "xe-0/0/3.1"
                          metric = 200
                        }
                      ]
                    }
                  ]
                }
              ]
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
              name = "xe-0/0/2.2"
            },
            {
              name = "xe-0/0/3.2"
            },
            {
              name = "lo0.10002"
            }
          ]
          route_distinguisher = [
            {
              rd_type = "10.40.100.7:10002"
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
              ospf = [
                {
                  export = [
                    "to-ospf"
                  ]
                  area = [
                    {
                      name = "0.0.0.0"
                      interface = [
                        {
                          name = "xe-0/0/2.2"
                          metric = 100
                        },
                        {
                          name = "xe-0/0/3.2"
                          metric = 200
                        }
                      ]
                    }
                  ]
                }
              ]
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
              name = "xe-0/0/2.3"
            },
            {
              name = "xe-0/0/3.3"
            },
            {
              name = "lo0.10003"
            }
          ]
          route_distinguisher = [
            {
              rd_type = "10.40.100.12:10003"
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
              ospf = [
                {
                  export = [
                    "to-ospf"
                  ]
                  area = [
                    {
                      name = "0.0.0.0"
                      interface = [
                        {
                          name = "xe-0/0/2.3"
                          metric = 100
                        },
                        {
                          name = "xe-0/0/3.3"
                          metric = 200
                        }
                      ]
                    }
                  ]
                }
              ]
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
          rd_type = "10.30.100.2:9999"
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
}
