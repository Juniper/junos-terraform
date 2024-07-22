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
      }
    ]
  ])
  bgp_neighbor_list = flatten([
    for group_name, value in var.bgp :
    [for obj in value.neighbors : {
      group_name  = group_name
      neighbor_ip = obj.name
      description = obj.description
      peer-as     = obj.peer-as
      }
    ]
  ])
}
terraform {
  required_providers {
    junos-vmx = {
      source  = "juniper/providers/junos-vmx"
      version = "22.41.101"
    }
  }
}
resource "junos-vmx_InterfacesInterfaceUnitFamilyInetAddressName" "ipaddr" {
  resource_name = "test_group"
  count         = length(local.unit_obj)
  name          = local.unit_obj[count.index].intf
  name__1       = local.unit_obj[count.index].unit_name
  name__2       = local.unit_obj[count.index].ipv4_addr
}
resource "junos-vmx_InterfacesInterfaceUnitFamilyIso" "family_iso" {
  resource_name = "test_group"
  count         = length(local.unit_obj)
  name          = (local.unit_obj[count.index].family_iso ? local.unit_obj[count.index].intf : "")
  name__1       = (local.unit_obj[count.index].family_iso ? local.unit_obj[count.index].unit_name : "")
}
resource "junos-vmx_InterfacesInterfaceUnitFamilyMpls" "family_mpls" {
  resource_name = "test_group"
  count         = length(local.unit_obj)
  name          = (local.unit_obj[count.index].family_mpls ? local.unit_obj[count.index].intf : "")
  name__1       = (local.unit_obj[count.index].family_mpls ? local.unit_obj[count.index].unit_name : "")
}
resource "junos-vmx_InterfacesInterfaceUnitFamilyIsoAddressName" "iso_addr" {
  resource_name = "test_group"
  count         = length(local.unit_obj)
  name          = local.unit_obj[count.index].intf
  name__1       = local.unit_obj[count.index].unit_name
  name__2       = lookup(local.unit_obj[count.index], "iso_addr", "")
}
resource "junos-vmx_Routing__OptionsAutonomous__SystemAs__Number" "as_number" {
  resource_name = "test_group"
  as__number    = var.routing-options.as_number
}
resource "junos-vmx_Routing__OptionsRouter__Id" "router_id" {
  resource_name = "test_group"
  router__id    = var.routing-options.router_id
}
resource "junos-vmx_ProtocolsIsisInterfacePoint__To__Point" "isis_p2p" {
  resource_name    = "test_group"
  count            = length(var.isis.interfaces)
  name             = (var.isis.interfaces[count.index].p2p == "true" ? var.isis.interfaces[count.index].name : "")
  point__to__point = ""
}


resource "junos-vmx_ProtocolsIsisLevelWide__Metrics__Only" "wide_metrics" {
  resource_name       = "test_group"
  count               = length(var.isis.level)
  name                = (var.isis.level[count.index].wide_metrics == "true" ? var.isis.level[count.index].name : "")
  wide__metrics__only = ""
}
resource "junos-vmx_ProtocolsIsisInterfaceLdp__Synchronization" "ldp_sync" {
  resource_name = "test_group"
  count         = length(var.isis.interfaces)
  name          = (var.isis.interfaces[count.index].ldp_sync == "true" ? var.isis.interfaces[count.index].name : "")
}
resource "junos-vmx_ProtocolsIsisInterfacePassive" "intf_passive" {
  resource_name = "test_group"
  count         = length(var.isis.interfaces)
  name          = (var.isis.interfaces[count.index].passive == "true" ? var.isis.interfaces[count.index].name : "")
}
resource "junos-vmx_ProtocolsIsisLevelDisable" "level_disable" {
  resource_name = "test_group"
  count         = length(var.isis.level)
  name          = (var.isis.level[count.index].disable == "true" ? var.isis.level[count.index].name : "")
  disable       = " "
}
resource "junos-vmx_ProtocolsIsisLevelAuthentication__Type" "isis_auth_type" {
  resource_name        = "test_group"
  count                = length(var.isis.level)
  name                 = var.isis.level[count.index].name
  authentication__type = lookup(var.isis.level[count.index], "auth_type", "")
}
resource "junos-vmx_ProtocolsIsisLevelAuthentication__Key" "isis_auth_key" {
  resource_name       = "test_group"
  count               = length(var.isis.level)
  name                = var.isis.level[count.index].name
  authentication__key = lookup(var.isis.level[count.index], "auth_key", "")
}
resource "junos-vmx_ProtocolsLdpInterfaceName" "ldp_intf" {
  resource_name = "test_group"
  count         = length(var.ldp.interfaces)
  name          = var.ldp.interfaces[count.index].name
}
resource "junos-vmx_ProtocolsLdpTrack__Igp__Metric" "track_igp_metric" {
  resource_name      = "test_group"
  track__igp__metric = (var.ldp.track_igp_metric == "true" ? " " : "null")
}
resource "junos-vmx_ProtocolsBgpGroupType" "bgp_group_type" {
  resource_name = "test_group"
  for_each = var.bgp
  name          = each.key
  type          = each.value.group_type
}
resource "junos-vmx_ProtocolsBgpGroupLocal__Address" "bgp_local_address" {
        resource_name = "test_group"
        for_each  = var.bgp
        name = each.key
        local__address = each.value.local_addr
}
resource "junos-vmx_ProtocolsBgpGroupFamilyInetUnicast" "inet_unicast" {
        resource_name = "test_group"
        for_each = var.bgp
        name = (each.value.af_inet_unicast == "true" ? each.key : "")
}
resource "junos-vmx_ProtocolsBgpGroupFamilyInet__VpnUnicast" "inet_vpn_unicast" {
        resource_name = "test_group"
        for_each = var.bgp
        name = (each.value.af_inet_vpn_unicast == "true" ? each.key : "")
}
resource "junos-vmx_ProtocolsBgpGroupFamilyEvpnSignaling" "evpn_signaling" {
        resource_name = "test_group"
        for_each = var.bgp
        name = (each.value.af_evpn_signaling == "true" ? each.key : "")
}
resource "junos-vmx_ProtocolsBgpGroupAuthentication__Key" "bgp_auth_key" {
        resource_name = "test_group"
        for_each = var.bgp
        name = each.key
        authentication__key = lookup(each.value, "auth_key", "")
}
resource "junos-vmx_ProtocolsBgpGroupNeighborName" "bgp_neighbors" {
        resource_name = "test_group"
        count = length(local.bgp_neighbor_list)
        name = local.bgp_neighbor_list[count.index].group_name
        name__1 = local.bgp_neighbor_list[count.index].neighbor_ip
}
