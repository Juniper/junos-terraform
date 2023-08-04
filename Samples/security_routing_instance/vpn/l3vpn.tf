# --------------------- #
# Configure L3VPN
# --------------------- #

resource "junos-vsrx_Routing__InstancesInstanceInstance__Type" "l3vpn_instance_type_1" {
    resource_name = "JTAF-L3VPN-1-TYPE"
    name = "L3VPN-1"
    instance__type = "vrf"
}

resource "junos-vsrx_Routing__InstancesInstanceVrf__Table__Label" "l3vpn_tablelabel_1" {
    resource_name = "JTAF-L3VPN-1-TLABEL"
    name = "L3VPN-1"
}

resource "junos-vsrx_Routing__InstancesInstanceRoute__DistinguisherRd__Type" "l3vpn_rd_1" {
    resource_name = "JTAF-L3VPN-1-RD"
    name = "L3VPN-1"
    rd__type = "100:100"
}

resource "junos-vsrx_Routing__InstancesInstanceVrf__TargetCommunity" "l3vpn_rt_1" {
    resource_name = "JTAF-L3VPN-1-RT"
    name = "L3VPN-1"
    community = "target:100:100"
}

resource "junos-vsrx_Routing__InstancesInstanceInterfaceName" "l3vpn_interface_1_1" {
    resource_name = "JTAF-L3VPN-1-INTERFACE-1"
    name = "L3VPN-1"
    name__1 = "ge-0/0/1.0"
}

resource "junos-vsrx_Routing__InstancesInstanceDescription" "l3vpn_description_1" {
    resource_name = "JTAF-L3VPN-1-DESCRIPTION"
    name = "L3VPN-1"
    description = "CUSTOMER-1 VRF"
}