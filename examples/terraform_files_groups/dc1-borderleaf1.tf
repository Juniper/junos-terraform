resource "terraform-provider-junos-vqfx-evpn-vxlan-groups" "dc1-borderleaf1-base-config" {
  resource_name = "base-config"
  provider = junos-vqfx-evpn-vxlan-groups.dc1_borderleaf1
  apply_groups = local.common_g_a654ee_apply_groups
  groups = [
    {
      name = "borderleaf"
      system = local.common_g_a654ee_groups_borderleaf_system
      chassis = local.common_g_a654ee_groups_borderleaf_chassis
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
                              name = "10.30.131.2/30"
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
                              name = "10.30.141.2/30"
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
                              name = "10.99.1.1/30"
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
                              name = "10.99.2.1/30"
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
                              name = "10.99.3.1/30"
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
                              name = "10.98.1.1/30"
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
                              name = "10.98.2.1/30"
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
                              name = "10.98.3.1/30"
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
                              name = "10.32.6.1/30"
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
                              name = "100.123.24.1/16"
                            }
                          ]
                        }
                      ]
                    }
                  ]
                }
              ]
            },
            local.common_g_a654ee_groups_borderleaf_interfaces_interface_em1,
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
                              name = "10.30.100.1/32"
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
                              name = "10.40.100.1/32"
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
                              name = "10.40.100.6/32"
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
                              name = "10.40.100.11/32"
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
      snmp = local.common_g_a654ee_groups_borderleaf_snmp
      forwarding_options = local.common_g_a654ee_groups_borderleaf_forwarding_options
      routing_options = [
        {
          static = local.common_g_a654ee_groups_borderleaf_routing_options_static
          router_id = "10.30.100.1"
          forwarding_table = local.common_g_a654ee_groups_borderleaf_routing_options_forwarding_table
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
                  multihop = local.common_g_a654ee_groups_borderleaf_protocols_bgp_group_WAN_OVERLAY_eBGP_multihop
                  local_address = "10.30.100.1"
                  family = local.common_g_a654ee_groups_borderleaf_protocols_bgp_group_WAN_OVERLAY_eBGP_family
                  local_as = local.common_g_a654ee_groups_borderleaf_protocols_bgp_group_WAN_OVERLAY_eBGP_local_as
                  multipath = local.common_g_a654ee_groups_borderleaf_protocols_bgp_group_WAN_OVERLAY_eBGP_multipath
                  neighbor = local.common_g_a654ee_groups_borderleaf_protocols_bgp_group_WAN_OVERLAY_eBGP_neighbor
                },
                {
                  name = "EVPN_iBGP"
                  type = "internal"
                  local_address = "10.30.100.1"
                  family = local.common_g_a654ee_groups_borderleaf_protocols_bgp_group_EVPN_iBGP_family
                  cluster = "10.30.100.1"
                  local_as = local.common_g_a654ee_groups_borderleaf_protocols_bgp_group_EVPN_iBGP_local_as
                  multipath = local.common_g_a654ee_groups_borderleaf_protocols_bgp_group_EVPN_iBGP_multipath
                  neighbor = local.common_g_a654ee_groups_borderleaf_protocols_bgp_group_EVPN_iBGP_neighbor
                },
                {
                  name = "IPCLOS_eBGP"
                  type = "external"
                  mtu_discovery = ""
                  import = local.common_g_a654ee_groups_borderleaf_protocols_bgp_group_IPCLOS_eBGP_import
                  export = local.common_g_a654ee_groups_borderleaf_protocols_bgp_group_IPCLOS_eBGP_export
                  vpn_apply_export = ""
                  local_as = [
                    {
                      as_number = 65506
                    }
                  ]
                  multipath = local.common_g_a654ee_groups_borderleaf_protocols_bgp_group_IPCLOS_eBGP_multipath
                  bfd_liveness_detection = local.common_g_a654ee_groups_borderleaf_protocols_bgp_group_IPCLOS_eBGP_bfd_liveness_detection
                  neighbor = [
                    {
                      name = "10.30.131.1"
                      description = "EBGP peering to 10.30.131.1"
                      peer_as = 65501
                    },
                    {
                      name = "10.30.141.1"
                      description = "EBGP peering to 10.30.141.1"
                      peer_as = 65502
                    },
                    {
                      name = "10.32.6.2"
                      description = "EBGP peering to 10.32.6.2"
                      peer_as = 65400
                    }
                  ]
                }
              ]
            }
          ]
          evpn = local.common_g_a654ee_groups_borderleaf_protocols_evpn
          lldp = local.common_g_a654ee_groups_borderleaf_protocols_lldp
          igmp_snooping = local.common_g_a654ee_groups_borderleaf_protocols_igmp_snooping
        }
      ]
      policy_options = [
        {
          policy_statement = [
            local.common_g_a654ee_groups_borderleaf_policy_options_policy_statement_EVPN_T5_EXPORT,
            {
              name = "IPCLOS_BGP_EXP"
              term = [
                {
                  name = "loopback"
                  from = local.common_g_a654ee_groups_borderleaf_policy_options_policy_statement_IPCLOS_BGP_EXP_term_loopback_from
                  then = [
                    {
                      community = [
                        {
                          add = local.common_g_a654ee_groups_borderleaf_policy_options_policy_statement_IPCLOS_BGP_EXP_term_loopback_then_community_add
                          community_name = "dc1-borderleaf1"
                        }
                      ]
                      accept = local.common_g_a654ee_groups_borderleaf_policy_options_policy_statement_IPCLOS_BGP_EXP_term_loopback_then_accept
                    }
                  ]
                },
                local.common_g_a654ee_groups_borderleaf_policy_options_policy_statement_IPCLOS_BGP_EXP_term_default
              ]
            },
            local.common_g_a654ee_groups_borderleaf_policy_options_policy_statement_IPCLOS_BGP_IMP,
            local.common_g_a654ee_groups_borderleaf_policy_options_policy_statement_PFE_LB,
            local.common_g_a654ee_groups_borderleaf_policy_options_policy_statement_to_ospf
          ]
          community = [
            {
              name = "dc1-borderleaf1"
              members = [
                "65506:1"
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
              interface = local.common_g_a654ee_groups_borderleaf_routing_instances_instance_VRF_10001_interface
              route_distinguisher = [
                {
                  rd_type = "10.40.100.1:10001"
                }
              ]
              vrf_target = local.common_g_a654ee_groups_borderleaf_routing_instances_instance_VRF_10001_vrf_target
              vrf_table_label = local.common_g_a654ee_groups_borderleaf_routing_instances_instance_VRF_10001_vrf_table_label
              routing_options = local.common_g_a654ee_groups_borderleaf_routing_instances_instance_VRF_10001_routing_options
              protocols = local.common_g_a654ee_groups_borderleaf_routing_instances_instance_VRF_10001_protocols
            },
            {
              name = "VRF_10002"
              instance_type = "vrf"
              interface = local.common_g_a654ee_groups_borderleaf_routing_instances_instance_VRF_10002_interface
              route_distinguisher = [
                {
                  rd_type = "10.40.100.6:10002"
                }
              ]
              vrf_target = local.common_g_a654ee_groups_borderleaf_routing_instances_instance_VRF_10002_vrf_target
              vrf_table_label = local.common_g_a654ee_groups_borderleaf_routing_instances_instance_VRF_10002_vrf_table_label
              routing_options = local.common_g_a654ee_groups_borderleaf_routing_instances_instance_VRF_10002_routing_options
              protocols = local.common_g_a654ee_groups_borderleaf_routing_instances_instance_VRF_10002_protocols
            },
            {
              name = "VRF_10003"
              instance_type = "vrf"
              interface = local.common_g_a654ee_groups_borderleaf_routing_instances_instance_VRF_10003_interface
              route_distinguisher = [
                {
                  rd_type = "10.40.100.11:10003"
                }
              ]
              vrf_target = local.common_g_a654ee_groups_borderleaf_routing_instances_instance_VRF_10003_vrf_target
              vrf_table_label = local.common_g_a654ee_groups_borderleaf_routing_instances_instance_VRF_10003_vrf_table_label
              routing_options = local.common_g_a654ee_groups_borderleaf_routing_instances_instance_VRF_10003_routing_options
              protocols = local.common_g_a654ee_groups_borderleaf_routing_instances_instance_VRF_10003_protocols
            }
          ]
        }
      ]
      switch_options = [
        {
          vtep_source_interface = local.common_g_a654ee_groups_borderleaf_switch_options_vtep_source_interface
          route_distinguisher = [
            {
              rd_type = "10.30.100.1:9999"
            }
          ]
          vrf_target = local.common_g_a654ee_groups_borderleaf_switch_options_vrf_target
        }
      ]
    }
  ]
  system = [
    {
      host_name = "dc1-borderleaf1"
    }
  ]
}
