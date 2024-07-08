locals {
  bgp_neighbor_list = flatten([
    for group_name, value in var.bgp :
    [for obj in value.neighbors : {
      group_name  = group_name
      neighbor_ip = obj.name
      description = obj.description
      peer-as = obj.peer-as
      }
    ]
  ])
}
terraform {
  required_providers {
    junos-vsrx = {
      source = "juniper/providers/junos-vsrx"
      version = "19.41.101"
    }
  }
}
resource "junos-vsrx_ProtocolsBgpGroupType" "vsrx_r1" {
  resource_name = "my_group"
  for_each = var.bgp
  name          = each.key
  type          = each.value.neighbor_type
}
resource "junos-vsrx_ProtocolsBgpGroupLocal__Address" "vsrx_2" {
  resource_name  = "my_group"
  for_each = var.bgp
  name           = each.key
  local__address = lookup(each.value, "local-address", "")
}
resource "junos-vsrx_ProtocolsBgpGroupLocal__AsAs__Number" "vsrx_3" {
  resource_name = "my_group"
  for_each = var.bgp
  name          = each.key
  as__number    = lookup(each.value, "local-as", "")
}
resource "junos-vsrx_ProtocolsBgpGroupCluster" "vsrx_4" {
  resource_name = "my_group"
  for_each = var.bgp
  name          = each.key
  cluster       = lookup(each.value, "cluster", "")
}
resource "junos-vsrx_ProtocolsBgpGroupMtu__Discovery" "vsrx_5" {
  resource_name = "my_group"
  for_each = var.bgp
  name          = (each.value.mtu-discovery == 1 ? each.key : "" )
}
resource "junos-vsrx_ProtocolsBgpGroupMultipath" "vsrx_6" {
  resource_name = "my_group"
  for_each = var.bgp
  name          = (each.value.multipath == 1 ? each.key : "")
}
resource "junos-vsrx_ProtocolsBgpGroupMultipathMultiple__As" "vsrx_7" {
  resource_name = "my_group"
  for_each = var.bgp
  name          = (each.value.multipath-multiple-as == 1 ? each.key : "")
}
resource "junos-vsrx_ProtocolsBgpGroupFamilyEvpnSignaling" "vsrx_8" {
  resource_name = "my_group"
  for_each = var.bgp
  name          = (each.value.af-evpn-signaling == 1 ? each.key : "")
}
resource "junos-vsrx_ProtocolsBgpGroupBfd__Liveness__DetectionMultiplier" "vsrx_9" {
  resource_name = "my_group"
  for_each = var.bgp
  name          = each.key
  multiplier    = lookup(each.value, "bfd-multiplier", "")
}
resource "junos-vsrx_ProtocolsBgpGroupBfd__Liveness__DetectionMinimum__Interval" "vsrx_10" {
  resource_name     = "my_group"
  for_each = var.bgp
  name              = each.key
  minimum__interval = lookup(each.value, "bfd-minimum-interval", "")
}
resource "junos-vsrx_ProtocolsBgpGroupVpn__Apply__Export" "vsrx_11" {
  resource_name = "my_group"
  for_each = var.bgp
  name          = (each.value.vpn-apply-export == 1 ? each.key : "")
}
resource "junos-vsrx_ProtocolsBgpGroupNeighborName" "vsrx_12" {
  resource_name = "my_group"
  count = length(local.bgp_neighbor_list)
  name          = local.bgp_neighbor_list[count.index].group_name
  name__1       = local.bgp_neighbor_list[count.index].neighbor_ip
}
resource "junos-vsrx_ProtocolsBgpGroupNeighborDescription" "vsrx_13" {
  resource_name = "my_group"
  count = length(local.bgp_neighbor_list)
  name          = local.bgp_neighbor_list[count.index].group_name
  name__1       = local.bgp_neighbor_list[count.index].neighbor_ip
  description   = lookup(local.bgp_neighbor_list[count.index], "description", "")
}
resource "junos-vsrx_ProtocolsBgpGroupNeighborPeer__As" "vsrx_14" {
  resource_name = "my_group"
  count = length(local.bgp_neighbor_list)
  name          = local.bgp_neighbor_list[count.index].group_name
  name__1       = local.bgp_neighbor_list[count.index].neighbor_ip
  peer__as      = lookup(local.bgp_neighbor_list[count.index], "peer-as", "")
}
