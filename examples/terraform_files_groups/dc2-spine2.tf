resource "terraform-provider-junos-vqfx-evpn-vxlan-groups" "dc2-spine2-base-config" {
  resource_name = "base-config"
  provider = junos-vqfx-evpn-vxlan-groups.dc2_spine2
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
            local.common_g_49453a_groups_spine_interfaces_interface_xe_0_0_3,
            local.common_g_49453a_groups_spine_interfaces_interface_xe_0_0_4,
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
            local.common_g_49453a_groups_spine_interfaces_interface_ae0,
            local.common_g_49453a_groups_spine_interfaces_interface_ae1,
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
            local.common_g_0c49e5_groups_spine_interfaces_interface_em1,
            local.common_g_49453a_groups_spine_interfaces_interface_irb,
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
      snmp = local.common_g_0c49e5_groups_spine_snmp
      forwarding_options = local.common_g_0c49e5_groups_spine_forwarding_options
      routing_options = [
        {
          static = local.common_g_0c49e5_groups_spine_routing_options_static
          router_id = "10.30.100.9"
          forwarding_table = local.common_g_49453a_groups_spine_routing_options_forwarding_table
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
                  multihop = local.common_g_49453a_groups_spine_protocols_bgp_group_WAN_OVERLAY_eBGP_multihop
                  local_address = "10.30.100.9"
                  family = local.common_g_49453a_groups_spine_protocols_bgp_group_WAN_OVERLAY_eBGP_family
                  local_as = local.common_g_49453a_groups_spine_protocols_bgp_group_WAN_OVERLAY_eBGP_local_as
                  multipath = local.common_g_49453a_groups_spine_protocols_bgp_group_WAN_OVERLAY_eBGP_multipath
                  neighbor = local.common_g_49453a_groups_spine_protocols_bgp_group_WAN_OVERLAY_eBGP_neighbor
                },
                {
                  name = "EVPN_iBGP"
                  type = "internal"
                  local_address = "10.30.100.9"
                  family = local.common_g_0c49e5_groups_spine_protocols_bgp_group_EVPN_iBGP_family
                  cluster = "10.30.100.9"
                  local_as = local.common_g_49453a_groups_spine_protocols_bgp_group_EVPN_iBGP_local_as
                  multipath = local.common_g_0c49e5_groups_spine_protocols_bgp_group_EVPN_iBGP_multipath
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
                  import = local.common_g_0c49e5_groups_spine_protocols_bgp_group_IPCLOS_eBGP_import
                  export = local.common_g_0c49e5_groups_spine_protocols_bgp_group_IPCLOS_eBGP_export
                  vpn_apply_export = ""
                  local_as = [
                    {
                      as_number = 65521
                    }
                  ]
                  multipath = local.common_g_0c49e5_groups_spine_protocols_bgp_group_IPCLOS_eBGP_multipath
                  bfd_liveness_detection = local.common_g_0c49e5_groups_spine_protocols_bgp_group_IPCLOS_eBGP_bfd_liveness_detection
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
          evpn = local.common_g_49453a_groups_spine_protocols_evpn
          lldp = local.common_g_0c49e5_groups_spine_protocols_lldp
          igmp_snooping = local.common_g_0c49e5_groups_spine_protocols_igmp_snooping
        }
      ]
      policy_options = [
        {
          policy_statement = [
            local.common_g_49453a_groups_spine_policy_options_policy_statement_EVPN_T5_EXPORT,
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
                          community_name = "dc2-spine2"
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
            local.common_g_0c49e5_groups_spine_policy_options_policy_statement_PFE_LB,
            local.common_g_49453a_groups_spine_policy_options_policy_statement_to_ospf
          ]
          community = [
            {
              name = "dc2-spine2"
              members = [
                "65521:1"
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
              interface = local.common_g_49453a_groups_spine_routing_instances_instance_VRF_10001_interface
              route_distinguisher = [
                {
                  rd_type = "10.40.101.2:10001"
                }
              ]
              vrf_target = local.common_g_49453a_groups_spine_routing_instances_instance_VRF_10001_vrf_target
              vrf_table_label = local.common_g_49453a_groups_spine_routing_instances_instance_VRF_10001_vrf_table_label
              routing_options = local.common_g_49453a_groups_spine_routing_instances_instance_VRF_10001_routing_options
              protocols = local.common_g_49453a_groups_spine_routing_instances_instance_VRF_10001_protocols
            }
          ]
        }
      ]
      switch_options = [
        {
          vtep_source_interface = local.common_g_49453a_groups_spine_switch_options_vtep_source_interface
          route_distinguisher = [
            {
              rd_type = "10.30.100.9:9999"
            }
          ]
          vrf_target = local.common_g_49453a_groups_spine_switch_options_vrf_target
        }
      ]
      vlans = local.common_g_49453a_groups_spine_vlans
    }
  ]
  system = [
    {
      host_name = "dc2-spine2"
    }
  ]
}
