terraform {
  required_providers {
    junos-vqfx = {
      source = "junos-vqfx"
    }
  }
}

provider "junos-vqfx" {
  host     = ""
  port     = 22
  username = ""
  password = ""
  sshkey   = ""
  alias    = "dc2-spine2"
}

provider "junos-vqfx" {
  host     = ""
  port     = 22
  username = ""
  password = ""
  sshkey   = ""
  alias    = "dc1-1eaf1"
}

provider "junos-vqfx" {
  host     = ""
  port     = 22
  username = ""
  password = ""
  sshkey   = ""
  alias    = "dc1-spine1"
}

resource "junos-vqfx_Apply_Groups" "dc2-spine2" {
  resource_name = "JTAF_dc2-spine2"
  provider = junos-vqfx.dc2-spine2
  system = [
    {
      host_name = "dc2-spine2"
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
                  address = ["0.0.0.0/0"]
                }
              ]
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
                  deployment_scope = ["commercial"]
                }
              ]
            },
            {
              name = "chef"
              license_type = [
                {
                  name = "juniper"
                  deployment_scope = ["commercial"]
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
          description = "*** to wan-pe2 ***"
          unit = [
            {
              name = 0
              family = [
                {
                  inet = [
                    {
                      address = [
                        {
                          name = "10.32.12.1/30"
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
                          name = "10.93.1.1/30"
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
                          name = "10.92.1.1/30"
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
          name = "xe-0/0/4"
          ether_options = [
            {
              ieee_802_3ad = [
                {
                  bundle = "ae1"
                }
              ]
            }
          ]
        },
        {
          name = "xe-0/0/5"
          description = "*** to dc2-spine1 ***"
          unit = [
            {
              name = 0
              family = [
                {
                  inet = [
                    {
                      address = [
                        {
                          name = "10.30.189.2/30"
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
          name = "ae0"
          esi = [
            {
              identifier = "00:00:00:00:00:00:00:01:01:00"
              all_active = ""
            }
          ]
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
          unit = [
            {
              name = 0
              family = [
                {
                  ethernet_switching = [
                    {
                      vlan = [
                        {
                          members = [1002]
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
          name = "ae1"
          esi = [
            {
              identifier = "00:00:00:00:00:00:00:01:02:00"
              all_active = ""
            }
          ]
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
          unit = [
            {
              name = 0
              family = [
                {
                  ethernet_switching = [
                    {
                      vlan = [
                        {
                          members = [1002]
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
                          name = "100.123.24.9/16"
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
              name = 1002
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
                          name = "10.30.100.9/32"
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
                          name = "10.40.101.2/32"
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
              next_hop = ["100.123.0.1"]
            }
          ]
        }
      ]
      router_id = "10.30.100.9"
      forwarding_table = [
        {
          export = ["PFE-LB"]
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
              local_address = "10.30.100.9"
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
                  as_number = 65201
                }
              ]
              multipath = [
                {
                  multiple_as = ""
                }
              ]
              neighbor = [
                {
                  name = "10.30.100.1"
                  description = "DCI EBGP peering to 10.30.100.1"
                  peer_as = 65200
                },
                {
                  name = "10.30.100.2"
                  description = "DCI EBGP peering to 10.30.100.2"
                  peer_as = 65200
                }
              ]
            },
            {
              name = "EVPN_iBGP"
              type = "internal"
              local_address = "10.30.100.9"
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
              cluster = "10.30.100.9"
              local_as = [
                {
                  as_number = 65201
                }
              ]
              multipath = [
                {

                }
              ]
              neighbor = [
                {
                  name = "10.30.100.8"
                }
              ]
            },
            {
              name = "IPCLOS_eBGP"
              type = "external"
              mtu_discovery = ""
              import = ["IPCLOS_BGP_IMP"]
              export = ["IPCLOS_BGP_EXP"]
              vpn_apply_export = ""
              local_as = [
                {
                  as_number = 65521
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
                  name = "10.30.189.1"
                  description = "EBGP peering to 10.30.189.1"
                  peer_as = 65520
                },
                {
                  name = "10.32.12.2"
                  description = "EBGP peering to 10.32.12.2"
                  peer_as = 65401
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
          extended_vni_list = ["all"]
          no_core_isolation = ""
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
                  protocol = ["direct"]
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
                  protocol = ["static"]
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
                  protocol = ["evpn", "ospf"]
                },
                {
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
              name = "fm_v4_host"
              from = [
                {
                  protocol = ["evpn"]
                },
                {
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
                  protocol = ["evpn"]
                },
                                {
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
                  protocol = ["direct", "bgp"]
                },
              ]
              then = [
                {
                  community = [
                    {
                      add = ""
                      community_name = "dc2-spine2"
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
                  protocol = ["bgp", "direct"]
                },
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
                  protocol = ["evpn"]
                },
                {
                  route_filter = [
                    {
                      address = "10.1.2.0/24"
                      orlonger = ""
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
          name = "dc2-spine2"
          members = ["65521:1"]
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
          route_distinguisher = [
            {
              rd_type = "10.40.101.2:10001"
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
                  export = ["to-ospf"]
                  area = [
                    {
                      name = "0.0.0.0"
                      interface = [
                        {
                          name = "xe-0/0/1.1"
                          metric = 100
                        },
                        {
                          name = "xe-0/0/2.1"
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
                      export = ["EVPN_T5_EXPORT"]
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
          rd_type = "10.30.100.9:9999"
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
          name = "vlan_1002"
          vlan_id = 1002
          l3_interface = "irb.1002"
          vxlan = [
            {
              vni = 1002
            }
          ]
        }
      ]
    }
  ]
}

resource "junos-vqfx_Apply_Groups" "dc1-spine1" {
  resource_name = "JTAF_dc1-spine1"
  provider = junos-vqfx.dc1-spine1
  system = [
    {
      host_name = "dc1-spine1"
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
                      address = ["0.0.0.0/0"]
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
                  deployment_scope = ["commercial"]
                }
              ]
            },
            {
              name = "chef"
              license_type = [
                {
                  name = "juniper"
                  deployment_scope = ["commercial"]
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
          description = "*** to dc1-borderleaf1 ***"
          unit = [
            {
              name = 0
              family = [
                {
                  inet = [
                    {
                      address = [
                        {
                          name = "10.30.131.1/30"
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
          description = "*** to dc1-borderleaf2 ***"
          unit = [
            {
              name = 0
              family = [
                {
                  inet = [
                    {
                      address = [
                        {
                          name = "10.30.132.1/30"
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
          description = "*** to dc1-leaf1 ***"
          unit = [
            {
              name = 0
              family = [
                {
                  inet = [
                    {
                      address = [
                        {
                          name = "10.30.135.1/30"
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
          description = "*** to dc1-leaf2 ***"
          unit = [
            {
              name = 0
              family = [
                {
                  inet = [
                    {
                      address = [
                        {
                          name = "10.30.136.1/30"
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
          description = "*** to dc1-leaf3 ***"
          unit = [
            {
              name = 0
              family = [
                {
                  inet = [
                    {
                      address = [
                        {
                          name = "10.30.137.1/30"
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
                          name = "100.123.24.3/16"
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
                          name = "10.30.100.3/32"
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
              next_hop = ["100.123.0.1"]
            }
          ]
        }
      ]
      router_id = "10.30.100.3"
      forwarding_table = [
        {
          export = ["PFE-LB"]
          ecmp_fast_reroute = ""
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
              local_address = "10.30.100.3"
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
              cluster = "10.30.100.3"
              local_as = [
                {
                  as_number = 65200
                }
              ]
              multipath = [
                {

                }
              ]
              allow = "10.30.100.0/24"
              neighbor = [
                {
                  name = "10.30.100.4"
                }
              ]
            },
            {
              name = "IPCLOS_eBGP"
              type = "external"
              mtu_discovery = ""
              import = ["IPCLOS_BGP_IMP"]
              export = ["IPCLOS_BGP_EXP"]
              vpn_apply_export = ""
              local_as = [
                {
                  as_number = 65501
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
                  name = "10.30.135.2"
                  description = "EBGP peering to 10.30.135.2"
                  peer_as = 65503
                },
                {
                  name = "10.30.136.2"
                  description = "EBGP peering to 10.30.136.2"
                  peer_as = 65504
                },
                {
                  name = "10.30.137.2"
                  description = "EBGP peering to 10.30.137.2"
                  peer_as = 65505
                },
                {
                  name = "10.30.131.2"
                  description = "EBGP peering to 10.30.131.2"
                  peer_as = 65506
                },
                {
                  name = "10.30.132.2"
                  description = "EBGP peering to 10.30.132.2"
                  peer_as = 65507
                }
              ]
            }
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
          name = "IPCLOS_BGP_EXP"
          term = [
            {
              name = "loopback"
              from = [
                {
                  protocol = ["direct", "bgp"]
                },
              ]
              then = [
                {
                  community = [
                    {
                      add = ""
                      community_name = "dc1-spine1"
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
                  protocol = ["bgp", "direct"]
                },
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
          name = "dc1-spine1"
          members = ["65501:1"]
        }
      ]
    }
  ]
}

