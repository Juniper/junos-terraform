# ----------------------- #
# Required providers
# ----------------------- #

terraform {
	required_providers {
		junos-vsrx = {
			source = "juniper/providers/junos-vsrx"
			version = "1.1"
		}
	}
}

# --------------------- #
# Configure interfaces
# --------------------- #
# Configure Left Interface and place in untrust zone

resource "junos-vsrx_InterfacesInterfaceDescription" "InterfacesInterfaceDescription_1" {
    	resource_name = "JTAF-vSRX-Interfaces"
    	name = "ge-0/0/0"
    	description = "LEFT-UNTRUST"
}

resource "junos-vsrx_InterfacesInterfaceUnitFamilyInetAddressName" "InterfacesInterfaceUnitFamilyInetAddressName_1" {
    	resource_name = "JTAF-vSRX-Interfaces"
    	name = "ge-0/0/0"
    	name__1 = "0" # Unit number
    	name__2 = "10.1.1.1/30"
}

resource "junos-vsrx_SecurityZonesSecurity__ZoneInterfacesHost__Inbound__TrafficProtocolsName" "seczonesprotocol_1" {
	resource_name = "JTAF-vSRX-Interfaces"
	name = "untrust"
	name__1 = "ge-0/0/0.0"
	name__2 = "all"  # protocol all 
}

resource "junos-vsrx_SecurityZonesSecurity__ZoneInterfacesHost__Inbound__TrafficSystem__ServicesName" "seczonesservices_1" {
	resource_name = "JTAF-vSRX-Interfaces"
    	name = "untrust"
    	name__1 = "ge-0/0/0.0"
    	name__2 = "all" # services all
}

# Configure Right Interface and place the interface in trust zone
resource "junos-vsrx_InterfacesInterfaceDescription" "InterfacesInterfaceDescription_2" {
    	resource_name = "JTAF-vSRX-Interfaces"
    	name = "ge-0/0/1"
    	description = "RIGHT-TRUST"
}

resource "junos-vsrx_InterfacesInterfaceUnitFamilyInetAddressName" "InterfacesInterfaceUnitFamilyInetAddressName_2" {
    	resource_name = "JTAF-vSRX-Interfaces"
    	name = "ge-0/0/1"
    	name__1 = "0" # Unit number
    	name__2 = "10.2.1.1/30"
}

resource "junos-vsrx_SecurityZonesSecurity__ZoneInterfacesHost__Inbound__TrafficProtocolsName" "seczonesprotocol_2" {
	resource_name = "JTAF-vSRX-Interfaces"
	name = "untrust"
	name__1 = "ge-0/0/1.0"
	name__2 = "all"  # protocol all 
}

resource "junos-vsrx_SecurityZonesSecurity__ZoneInterfacesHost__Inbound__TrafficSystem__ServicesName" "seczonesservices_2" {
	resource_name = "JTAF-vSRX-Interfaces"
    	name = "untrust"
    	name__1 = "ge-0/0/1.0"
    	name__2 = "all" # services all
}
