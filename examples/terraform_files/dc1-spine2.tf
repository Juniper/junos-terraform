resource "terraform-provider-junos-vqfx" "dc1-spine2-base-config" {
  resource_name = "base-config"
  provider = junos-vqfx.dc1_spine2
  system = [
    {
      host_name = "dc1-spine2"
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
                          name = "10.30.141.1/30"
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
                          name = "10.30.142.1/30"
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
                          name = "10.30.145.1/30"
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
                          name = "10.30.146.1/30"
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
                          name = "10.30.147.1/30"
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
                          name = "100.123.24.4/16"
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
                          name = "10.30.100.4/32"
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
      router_id = "10.30.100.4"
      forwarding_table = [
        {
          export = [
            "PFE-LB"
          ]
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
              local_address = "10.30.100.4"
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
              cluster = "10.30.100.4"
              local_as = [
                {
                  as_number = 65200
                }
              ]
              multipath = [
                {

                }
              ]
              allow = [
                "10.30.100.0/24"
              ]
              neighbor = [
                {
                  name = "10.30.100.3"
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
                  as_number = 65502
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
                  name = "10.30.145.2"
                  description = "EBGP peering to 10.30.145.2"
                  peer_as = 65503
                },
                {
                  name = "10.30.146.2"
                  description = "EBGP peering to 10.30.146.2"
                  peer_as = 65504
                },
                {
                  name = "10.30.147.2"
                  description = "EBGP peering to 10.30.147.2"
                  peer_as = 65505
                },
                {
                  name = "10.30.141.2"
                  description = "EBGP peering to 10.30.141.2"
                  peer_as = 65506
                },
                {
                  name = "10.30.142.2"
                  description = "EBGP peering to 10.30.142.2"
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
                      community_name = "dc1-spine2"
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
          name = "dc1-spine2"
          members = [
            "65502:1"
          ]
        }
      ]
    }
  ]
}
