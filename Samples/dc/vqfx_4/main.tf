locals {
  unit_obj = flatten([
    for intf_name, value in var.interfaces :
    [for unit_name, unit_value in value.unit : {
      intf        = intf_name,
      unit_name   = unit_name,
      ipv4_addr   = unit_value.ipv4_addr,
      family_iso  = unit_value.family_iso,
      family_mpls = unit_value.family_mpls,
      iso_addr    = unit_value.iso_addr,
      proxy_macip_adv = unit_value.proxy_macip_adv,
      vga_accept_data = unit_value.vga_accept_data,
      ipv4_vga = unit_value.ipv4_vga,
      }
    ]
  ])
  policy_list = flatten([
    for policy_name, policy_value in var.policy-statement :
    [for term_name, term_value in policy_value.term : {
      policy_name        = policy_name,
      term_name   = term_name,
      from_protocol   =term_value.from_protocol,
      then_accept  = term_value.then_accept,
      then_reject = term_value.then_reject,
      load_balance_per_packet = term_value.load_balance_per_packet,
      }
    ]
  ])
  bgp_neighbor_list = flatten([
    for group in var.bgp.groups :
    [for obj in group.neighbors : {
      group_name  = group.group_name
      neighbor_ip = obj.name
      description = obj.description
      peer_as     = obj.peer_as
      }
    ]
  ])
}
terraform {
  required_providers {
    junos-vqfx = {
      source  = "juniper/providers/junos-vqfx"
      version = "21.31.108"
    }
  }
}
resource "junos-vqfx_InterfacesInterfaceDescription" "intf_desc" {
        resource_name = "test_group"
        for_each = var.interfaces
        name = each.key
        description = each.value.description
}
resource "junos-vqfx_InterfacesInterfaceMtu" "intf_mtu" {
        resource_name = "test_group"
        for_each = var.interfaces
        name = each.key
        mtu = each.value.mtu
}
resource "junos-vqfx_InterfacesInterfaceUnitFamilyInetAddressName" "ipaddr" {
  resource_name = "test_group"
  count         = length(local.unit_obj)
  name          = local.unit_obj[count.index].intf
  name__1       = local.unit_obj[count.index].unit_name
  name__2       = local.unit_obj[count.index].ipv4_addr
}

resource "junos-vqfx_Routing__OptionsRouter__Id" "router_id" {
        resource_name = "test_group"
        router__id = var.routing-options.router_id
}
resource "junos-vqfx_Routing__OptionsAutonomous__SystemAs__Number" "as_number" {
        resource_name = "test_group"
        as__number = var.routing-options.as_number
}
resource "junos-vqfx_Routing__OptionsForwarding__TableExport" "ft_export_policy" {
        resource_name = "test_group"
        export = var.routing-options.forwarding_table.export_policy
}
resource "junos-vqfx_VlansVlanVlan__Id" "vlans_vlan_id" {
        resource_name = "test_group"
        for_each = var.vlans
        name = each.key
        vlan__id = each.value.vlan_id
}
resource "junos-vqfx_Forwarding__OptionsStorm__Control__ProfilesAll" "storm_control_profile" {
        resource_name = "test_group"
        name = var.forwarding-options.storm_control_profile_all
}
resource "junos-vqfx_Policy__OptionsPolicy__StatementTermThenAccept" "policy_then_accept" {
        resource_name = "test_group"
        count = length(local.policy_list)
        name = local.policy_list[count.index].policy_name
        name__1 = (local.policy_list[count.index].then_accept == "true"? local.policy_list[count.index].term_name: "")
        accept = ""
}
resource "junos-vqfx_Policy__OptionsPolicy__StatementTermThenReject" "policy_then_reject" {
        resource_name = "test_group"
        count = length(local.policy_list)
        name = local.policy_list[count.index].policy_name
        name__1 = (local.policy_list[count.index].then_reject == "true"? local.policy_list[count.index].term_name: "")
        reject = ""
}
resource "junos-vqfx_Policy__OptionsPolicy__StatementThenLoad__BalancePer__Packet" "load_balance_pp" {
        resource_name = "test_group"
        count = length(local.policy_list)
        name = (local.policy_list[count.index].load_balance_per_packet == "true" ? local.policy_list[count.index].policy_name : "")
        per__packet = ""
}
resource "junos-vqfx_ProtocolsBgpLog__Updown" "bgp_log" {
        resource_name = "test_group"
        log__updown = " "
}
resource "junos-vqfx_ProtocolsBgpGraceful__Restart" "bgp_gr" {
        resource_name = "test_group"
}
resource "junos-vqfx_ProtocolsBgpGroupType" "bgp_group_type" {
  resource_name = "test_group"
  count = length(var.bgp.groups)
  name          = var.bgp.groups[count.index].group_name
  type          = var.bgp.groups[count.index].group_type
}
resource "junos-vqfx_ProtocolsBgpGroupMtu__Discovery" "bgp_mtu" {
        resource_name = "test_group"
        count = length(var.bgp.groups)
        name = (var.bgp.groups[count.index].mtu_discovery == "true" ? var.bgp.groups[count.index].group_name : "")
        mtu__discovery = " "
}
resource "junos-vqfx_ProtocolsBgpGroupMultipathMultiple__As" "bgp_multipath_multias" {
        resource_name = "test_group"
        count = length(var.bgp.groups)
        name = (var.bgp.groups[count.index].multipath_multiple_as == "true" ? var.bgp.groups[count.index].group_name : "")
        multiple__as = " "
}
resource "junos-vqfx_ProtocolsBgpGroupFamilyEvpnSignaling" "bgp_af_evpn" {
        resource_name = "test_group"
        count = length(var.bgp.groups)
        name = (var.bgp.groups[count.index].af_evpn_signaling == "true" ? var.bgp.groups[count.index].group_name : "")
}
resource "junos-vqfx_ProtocolsBgpGroupLocal__AsAs__Number" "bgp_local_as" {
        resource_name = "test_group"
        count = length(var.bgp.groups)
        name = var.bgp.groups[count.index].group_name
        as__number = lookup(var.bgp.groups[count.index], "local_as", "")
}
resource "junos-vqfx_ProtocolsBgpGroupLocal__Address" "bgp_local_address" {
    resource_name = "test_group"
    count = length(var.bgp.groups)
    name = var.bgp.groups[count.index].group_name
    local__address = var.bgp.groups[count.index].local_addr
}
resource "junos-vqfx_ProtocolsBgpGroupImport" "bgp_import" {
        resource_name = "test_group"
        count = length(var.bgp.groups)
        name = var.bgp.groups[count.index].group_name
        import = lookup(var.bgp.groups[count.index], "import_policy", "")
}
resource "junos-vqfx_ProtocolsBgpGroupExport" "bgp_export" {
        resource_name = "test_group"
        count = length(var.bgp.groups)
        name = var.bgp.groups[count.index].group_name
        export = lookup(var.bgp.groups[count.index], "export_policy", "")
}
resource "junos-vqfx_ProtocolsBgpGroupCluster" "bgp_cluster" {
	resource_name = "test_group"
        count = length(var.bgp.groups)
	name = var.bgp.groups[count.index].group_name
	cluster = var.bgp.groups[count.index].cluster
}
resource "junos-vqfx_ProtocolsBgpGroupNeighborName" "bgp_group_neighbors" {
        resource_name = "test_group"
        count = length(local.bgp_neighbor_list)
        name = local.bgp_neighbor_list[count.index].group_name
        name__1 = local.bgp_neighbor_list[count.index].neighbor_ip
}
resource "junos-vqfx_ProtocolsBgpGroupNeighborDescription" "bgp_group_neighbor_desc" {
        resource_name = "test_group"
        count = length(local.bgp_neighbor_list)
        name = local.bgp_neighbor_list[count.index].group_name
        name__1 = local.bgp_neighbor_list[count.index].neighbor_ip
        description = local.bgp_neighbor_list[count.index].description
}
resource "junos-vqfx_ProtocolsBgpGroupNeighborPeer__As" "bgp_group_neighbor_peer_as" {
        resource_name = "test_group"
        count = length(local.bgp_neighbor_list)
        name = local.bgp_neighbor_list[count.index].group_name
        name__1 = local.bgp_neighbor_list[count.index].neighbor_ip
        peer__as = local.bgp_neighbor_list[count.index].peer_as
}
resource "junos-vqfx_ProtocolsLldpInterfaceName" "lldp_intf" {
        resource_name = "test_group"
        name = var.lldp.interface
}
resource "junos-vqfx_ProtocolsIgmp__SnoopingVlanName" "igmp_snooping_vlan" {
        resource_name = "test_group"
        name = var.igmp-snooping.vlan
}
