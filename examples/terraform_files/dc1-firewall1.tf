resource "terraform-provider-junos-vsrx" "dc1-firewall1-base-config" {
  resource_name = "base-config"
  provider = junos-vsrx.dc1_firewall1
  system = [
    {
      login = [
        {
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
          message = "***********************************************************************\nThis system is restricted to __________, authorized users for legitimate\nbusiness purposes only. All activity on the system will be logged and\nis subject to monitoring. Unauthorized access, use or modification\nof computers, data therein or data in transit to or from the computers\nis a violation of state and federal laws. Unauthorized activity will\nbe reported to the law enforcement for investigation and possible\nprosecution. __________ reserves the right to investigate, refer for\nprosecution and pursue monetary damages in civil actions in the event\nof unauthorized access.\n***********************************************************************\n"
        }
      ]
      root_authentication = [
        {
          encrypted_password = "$1$DbZ1Q3pj$s48cZytjsmSJRUJAf4LdM."
        }
      ]
      host_name = "dc1-firewall1"
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
                  any = ""
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
      license = [
        {
          autoupdate = [
            {
              url = [
                {
                  name = "https://ae1.juniper.net/junos/key_retrieval"
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
  security = [
    {
      log = [
        {
          mode = "stream"
          report = [
            {

            }
          ]
        }
      ]
      screen = [
        {
          ids_option = [
            {
              name = "untrust-screen"
              icmp = [
                {
                  ping_death = ""
                }
              ]
              ip = [
                {
                  source_route_option = ""
                  tear_drop = ""
                }
              ]
              tcp = [
                {
                  syn_flood = [
                    {
                      alarm_threshold = 1024
                      attack_threshold = 200
                      source_threshold = 1024
                      destination_threshold = 2048
                      undocumented = [
                        {
                          queue_size = 2000
                        }
                      ]
                      timeout = 20
                    }
                  ]
                  land = ""
                }
              ]
            }
          ]
        }
      ]
      policies = [
        {
          policy = [
            {
              from_zone_name = "trust"
              to_zone_name = "trust"
              policy = [
                {
                  name = "default-permit"
                  match = [
                    {
                      source_address = [
                        "any"
                      ]
                      destination_address = [
                        "any"
                      ]
                      application = [
                        "any"
                      ]
                    }
                  ]
                  then = [
                    {
                      permit = [
                        {

                        }
                      ]
                    }
                  ]
                }
              ]
            },
            {
              from_zone_name = "trust"
              to_zone_name = "untrust"
              policy = [
                {
                  name = "default-permit"
                  match = [
                    {
                      source_address = [
                        "any"
                      ]
                      destination_address = [
                        "any"
                      ]
                      application = [
                        "any"
                      ]
                    }
                  ]
                  then = [
                    {
                      permit = [
                        {

                        }
                      ]
                    }
                  ]
                }
              ]
            },
            {
              from_zone_name = "VRF_10001"
              to_zone_name = "VRF_10002"
              policy = [
                {
                  name = "default-permit"
                  match = [
                    {
                      source_address = [
                        "any"
                      ]
                      destination_address = [
                        "any"
                      ]
                      application = [
                        "any"
                      ]
                    }
                  ]
                  then = [
                    {
                      permit = [
                        {

                        }
                      ]
                    }
                  ]
                }
              ]
            },
            {
              from_zone_name = "VRF_10002"
              to_zone_name = "VRF_10001"
              policy = [
                {
                  name = "default-permit"
                  match = [
                    {
                      source_address = [
                        "any"
                      ]
                      destination_address = [
                        "any"
                      ]
                      application = [
                        "any"
                      ]
                    }
                  ]
                  then = [
                    {
                      permit = [
                        {

                        }
                      ]
                    }
                  ]
                }
              ]
            },
            {
              from_zone_name = "VRF_10002"
              to_zone_name = "VRF_10003"
              policy = [
                {
                  name = "default-permit"
                  match = [
                    {
                      source_address = [
                        "any"
                      ]
                      destination_address = [
                        "any"
                      ]
                      application = [
                        "any"
                      ]
                    }
                  ]
                  then = [
                    {
                      permit = [
                        {

                        }
                      ]
                    }
                  ]
                }
              ]
            },
            {
              from_zone_name = "VRF_10003"
              to_zone_name = "VRF_10002"
              policy = [
                {
                  name = "default-permit"
                  match = [
                    {
                      source_address = [
                        "any"
                      ]
                      destination_address = [
                        "any"
                      ]
                      application = [
                        "any"
                      ]
                    }
                  ]
                  then = [
                    {
                      permit = [
                        {

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
      zones = [
        {
          security_zone = [
            {
              name = "trust"
              tcp_rst = ""
              host_inbound_traffic = [
                {
                  system_services = [
                    {
                      name = "all"
                    }
                  ]
                  protocols = [
                    {
                      name = "all"
                    }
                  ]
                }
              ]
            },
            {
              name = "untrust"
              screen = "untrust-screen"
            },
            {
              name = "VRF_10001"
              host_inbound_traffic = [
                {
                  system_services = [
                    {
                      name = "all"
                    }
                  ]
                  protocols = [
                    {
                      name = "all"
                    }
                  ]
                }
              ]
              interfaces = [
                {
                  name = "ge-0/0/0.1"
                },
                {
                  name = "ge-0/0/1.1"
                }
              ]
            },
            {
              name = "VRF_10002"
              host_inbound_traffic = [
                {
                  system_services = [
                    {
                      name = "all"
                    }
                  ]
                  protocols = [
                    {
                      name = "all"
                    }
                  ]
                }
              ]
              interfaces = [
                {
                  name = "ge-0/0/0.2"
                },
                {
                  name = "ge-0/0/1.2"
                }
              ]
            },
            {
              name = "VRF_10003"
              host_inbound_traffic = [
                {
                  system_services = [
                    {
                      name = "all"
                    }
                  ]
                  protocols = [
                    {
                      name = "all"
                    }
                  ]
                }
              ]
              interfaces = [
                {
                  name = "ge-0/0/0.3"
                },
                {
                  name = "ge-0/0/1.3"
                }
              ]
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
          name = "ge-0/0/0"
          vlan_tagging = ""
          unit = [
            {
              name = 1
              description = "*** to dc1-borderleaf1 vlan 1 ***"
              vlan_id = 1
              family = [
                {
                  inet = [
                    {
                      address = [
                        {
                          name = "10.99.1.2/30"
                        }
                      ]
                    }
                  ]
                }
              ]
            },
            {
              name = 2
              description = "*** to dc1-borderleaf1 vlan 2 ***"
              vlan_id = 2
              family = [
                {
                  inet = [
                    {
                      address = [
                        {
                          name = "10.99.2.2/30"
                        }
                      ]
                    }
                  ]
                }
              ]
            },
            {
              name = 3
              description = "*** to dc1-borderleaf1 vlan 3 ***"
              vlan_id = 3
              family = [
                {
                  inet = [
                    {
                      address = [
                        {
                          name = "10.99.3.2/30"
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
          name = "ge-0/0/1"
          vlan_tagging = ""
          unit = [
            {
              name = 1
              description = "*** to dc1-borderleaf2 vlan 1 ***"
              vlan_id = 1
              family = [
                {
                  inet = [
                    {
                      address = [
                        {
                          name = "10.97.1.2/30"
                        }
                      ]
                    }
                  ]
                }
              ]
            },
            {
              name = 2
              description = "*** to dc1-borderleaf2 vlan 2 ***"
              vlan_id = 2
              family = [
                {
                  inet = [
                    {
                      address = [
                        {
                          name = "10.97.2.2/30"
                        }
                      ]
                    }
                  ]
                }
              ]
            },
            {
              name = 3
              description = "*** to dc1-borderleaf2 vlan 3 ***"
              vlan_id = 3
              family = [
                {
                  inet = [
                    {
                      address = [
                        {
                          name = "10.97.3.2/30"
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
          name = "fxp0"
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
                          name = "100.123.26.1/16"
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
    }
  ]
  protocols = [
    {
      lldp = [
        {
          interface = [
            {
              name = "all"
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
          name = "ospf-default"
          term = [
            {
              name = 1
              from = [
                {
                  route_filter = [
                    {
                      address = "0.0.0.0/0"
                      exact = ""
                      accept = ""
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
  routing_instances = [
    {
      instance = [
        {
          name = "fabric"
          instance_type = "virtual-router"
          interface = [
            {
              name = "ge-0/0/0.1"
            },
            {
              name = "ge-0/0/0.2"
            },
            {
              name = "ge-0/0/0.3"
            },
            {
              name = "ge-0/0/1.1"
            },
            {
              name = "ge-0/0/1.2"
            },
            {
              name = "ge-0/0/1.3"
            }
          ]
          routing_options = [
            {
              static = [
                {
                  route = [
                    {
                      name = "0.0.0.0/0"
                      discard = ""
                    }
                  ]
                }
              ]
            }
          ]
          protocols = [
            {
              ospf = [
                {
                  export = [
                    "ospf-default"
                  ]
                  area = [
                    {
                      name = "0.0.0.0"
                      interface = [
                        {
                          name = "ge-0/0/0.1"
                        },
                        {
                          name = "ge-0/0/1.1"
                        },
                        {
                          name = "ge-0/0/0.2"
                        },
                        {
                          name = "ge-0/0/1.2"
                        },
                        {
                          name = "ge-0/0/0.3"
                        },
                        {
                          name = "ge-0/0/1.3"
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
}
