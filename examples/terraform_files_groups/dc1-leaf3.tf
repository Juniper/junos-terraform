resource "terraform-provider-junos-vqfx-evpn-vxlan-groups" "dc1-leaf3-base-config" {
  resource_name = "base-config"
  provider = junos-vqfx-evpn-vxlan-groups.dc1_leaf3
  apply_groups = local.common_g_84da76_apply_groups
  groups = [
    {
      name = "leaf"
      system = local.common_g_84da76_groups_leaf_system
      chassis = local.common_g_84da76_groups_leaf_chassis
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
            local.common_g_84da76_groups_leaf_interfaces_interface_xe_0_0_2,
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
            local.common_g_84da76_groups_leaf_interfaces_interface_em1,
            local.common_g_84da76_groups_leaf_interfaces_interface_irb,
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
      snmp = local.common_g_84da76_groups_leaf_snmp
      forwarding_options = local.common_g_84da76_groups_leaf_forwarding_options
      routing_options = [
        {
          static = local.common_g_84da76_groups_leaf_routing_options_static
          router_id = "10.30.100.7"
          forwarding_table = local.common_g_84da76_groups_leaf_routing_options_forwarding_table
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
                  family = local.common_g_84da76_groups_leaf_protocols_bgp_group_EVPN_iBGP_family
                  cluster = "10.30.100.7"
                  local_as = local.common_g_84da76_groups_leaf_protocols_bgp_group_EVPN_iBGP_local_as
                  multipath = local.common_g_84da76_groups_leaf_protocols_bgp_group_EVPN_iBGP_multipath
                  neighbor = local.common_g_84da76_groups_leaf_protocols_bgp_group_EVPN_iBGP_neighbor
                },
                {
                  name = "IPCLOS_eBGP"
                  type = "external"
                  mtu_discovery = ""
                  import = local.common_g_84da76_groups_leaf_protocols_bgp_group_IPCLOS_eBGP_import
                  export = local.common_g_84da76_groups_leaf_protocols_bgp_group_IPCLOS_eBGP_export
                  vpn_apply_export = ""
                  local_as = [
                    {
                      as_number = 65505
                    }
                  ]
                  multipath = local.common_g_84da76_groups_leaf_protocols_bgp_group_IPCLOS_eBGP_multipath
                  bfd_liveness_detection = local.common_g_84da76_groups_leaf_protocols_bgp_group_IPCLOS_eBGP_bfd_liveness_detection
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
          evpn = local.common_g_84da76_groups_leaf_protocols_evpn
          lldp = local.common_g_84da76_groups_leaf_protocols_lldp
          igmp_snooping = local.common_g_84da76_groups_leaf_protocols_igmp_snooping
        }
      ]
      policy_options = [
        {
          policy_statement = [
            local.common_g_84da76_groups_leaf_policy_options_policy_statement_EVPN_T5_EXPORT,
            {
              name = "IPCLOS_BGP_EXP"
              term = [
                {
                  name = "loopback"
                  from = local.common_g_84da76_groups_leaf_policy_options_policy_statement_IPCLOS_BGP_EXP_term_loopback_from
                  then = [
                    {
                      community = [
                        {
                          add = local.common_g_84da76_groups_leaf_policy_options_policy_statement_IPCLOS_BGP_EXP_term_loopback_then_community_add
                          community_name = "dc1-leaf3"
                        }
                      ]
                      accept = local.common_g_84da76_groups_leaf_policy_options_policy_statement_IPCLOS_BGP_EXP_term_loopback_then_accept
                    }
                  ]
                },
                local.common_g_84da76_groups_leaf_policy_options_policy_statement_IPCLOS_BGP_EXP_term_default
              ]
            },
            local.common_g_84da76_groups_leaf_policy_options_policy_statement_IPCLOS_BGP_IMP,
            local.common_g_84da76_groups_leaf_policy_options_policy_statement_PFE_LB
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
              interface = local.common_g_84da76_groups_leaf_routing_instances_instance_VRF_10001_interface
              route_distinguisher = [
                {
                  rd_type = "10.40.100.5:10001"
                }
              ]
              vrf_target = local.common_g_84da76_groups_leaf_routing_instances_instance_VRF_10001_vrf_target
              vrf_table_label = local.common_g_84da76_groups_leaf_routing_instances_instance_VRF_10001_vrf_table_label
              routing_options = local.common_g_84da76_groups_leaf_routing_instances_instance_VRF_10001_routing_options
              protocols = local.common_g_84da76_groups_leaf_routing_instances_instance_VRF_10001_protocols
            },
            {
              name = "VRF_10002"
              instance_type = "vrf"
              interface = local.common_g_84da76_groups_leaf_routing_instances_instance_VRF_10002_interface
              route_distinguisher = [
                {
                  rd_type = "10.40.100.10:10002"
                }
              ]
              vrf_target = local.common_g_84da76_groups_leaf_routing_instances_instance_VRF_10002_vrf_target
              vrf_table_label = local.common_g_84da76_groups_leaf_routing_instances_instance_VRF_10002_vrf_table_label
              routing_options = local.common_g_84da76_groups_leaf_routing_instances_instance_VRF_10002_routing_options
              protocols = local.common_g_84da76_groups_leaf_routing_instances_instance_VRF_10002_protocols
            },
            {
              name = "VRF_10003"
              instance_type = "vrf"
              interface = local.common_g_84da76_groups_leaf_routing_instances_instance_VRF_10003_interface
              route_distinguisher = [
                {
                  rd_type = "10.40.100.15:10003"
                }
              ]
              vrf_target = local.common_g_84da76_groups_leaf_routing_instances_instance_VRF_10003_vrf_target
              vrf_table_label = local.common_g_84da76_groups_leaf_routing_instances_instance_VRF_10003_vrf_table_label
              routing_options = local.common_g_84da76_groups_leaf_routing_instances_instance_VRF_10003_routing_options
              protocols = local.common_g_84da76_groups_leaf_routing_instances_instance_VRF_10003_protocols
            }
          ]
        }
      ]
      switch_options = [
        {
          vtep_source_interface = local.common_g_84da76_groups_leaf_switch_options_vtep_source_interface
          route_distinguisher = [
            {
              rd_type = "10.30.100.7:9999"
            }
          ]
          vrf_target = local.common_g_84da76_groups_leaf_switch_options_vrf_target
        }
      ]
      vlans = local.common_g_84da76_groups_leaf_vlans
    }
  ]
  system = [
    {
      host_name = "dc1-leaf3"
    }
  ]
}
