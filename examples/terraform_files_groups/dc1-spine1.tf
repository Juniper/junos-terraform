resource "terraform-provider-junos-vqfx-evpn-vxlan-groups" "dc1-spine1-base-config" {
  resource_name = "base-config"
  provider = junos-vqfx-evpn-vxlan-groups.dc1_spine1
  apply_groups = local.common_g_0c49e5_apply_groups
  groups = [
    {
      name = "spine"
      system = local.common_g_0c49e5_groups_spine_system
      chassis = local.common_g_0c49e5_groups_spine_chassis
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
            local.common_g_0c49e5_groups_spine_interfaces_interface_em1,
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
      snmp = local.common_g_0c49e5_groups_spine_snmp
      forwarding_options = local.common_g_0c49e5_groups_spine_forwarding_options
      routing_options = [
        {
          static = local.common_g_0c49e5_groups_spine_routing_options_static
          router_id = "10.30.100.3"
          forwarding_table = local.common_g_9a472e_groups_spine_routing_options_forwarding_table
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
                  family = local.common_g_0c49e5_groups_spine_protocols_bgp_group_EVPN_iBGP_family
                  cluster = "10.30.100.3"
                  local_as = local.common_g_9a472e_groups_spine_protocols_bgp_group_EVPN_iBGP_local_as
                  multipath = local.common_g_0c49e5_groups_spine_protocols_bgp_group_EVPN_iBGP_multipath
                  allow = local.common_g_9a472e_groups_spine_protocols_bgp_group_EVPN_iBGP_allow
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
                  import = local.common_g_0c49e5_groups_spine_protocols_bgp_group_IPCLOS_eBGP_import
                  export = local.common_g_0c49e5_groups_spine_protocols_bgp_group_IPCLOS_eBGP_export
                  vpn_apply_export = ""
                  local_as = [
                    {
                      as_number = 65501
                    }
                  ]
                  multipath = local.common_g_0c49e5_groups_spine_protocols_bgp_group_IPCLOS_eBGP_multipath
                  bfd_liveness_detection = local.common_g_0c49e5_groups_spine_protocols_bgp_group_IPCLOS_eBGP_bfd_liveness_detection
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
          lldp = local.common_g_0c49e5_groups_spine_protocols_lldp
          igmp_snooping = local.common_g_0c49e5_groups_spine_protocols_igmp_snooping
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
                  from = local.common_g_0c49e5_groups_spine_policy_options_policy_statement_IPCLOS_BGP_EXP_term_loopback_from
                  then = [
                    {
                      community = [
                        {
                          add = local.common_g_0c49e5_groups_spine_policy_options_policy_statement_IPCLOS_BGP_EXP_term_loopback_then_community_add
                          community_name = "dc1-spine1"
                        }
                      ]
                      accept = local.common_g_0c49e5_groups_spine_policy_options_policy_statement_IPCLOS_BGP_EXP_term_loopback_then_accept
                    }
                  ]
                },
                local.common_g_0c49e5_groups_spine_policy_options_policy_statement_IPCLOS_BGP_EXP_term_default
              ]
            },
            local.common_g_0c49e5_groups_spine_policy_options_policy_statement_IPCLOS_BGP_IMP,
            local.common_g_0c49e5_groups_spine_policy_options_policy_statement_PFE_LB
          ]
          community = [
            {
              name = "dc1-spine1"
              members = [
                "65501:1"
              ]
            }
          ]
        }
      ]
    }
  ]
  system = [
    {
      host_name = "dc1-spine1"
    }
  ]
}
