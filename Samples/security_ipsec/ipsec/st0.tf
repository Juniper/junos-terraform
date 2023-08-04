# --------------------- #
# Configure st0 interface
# --------------------- #
# Configure st0 Interface and place in untrust zone

resource "junos-vsrx_InterfacesInterfaceDescription" "InterfacesInterfaceDescription_st0" {
    	resource_name = "JTAF-IPSEC-ST-Interfaces1"
    	name = "st0"
    	description = "customer interface 1"
}

resource "junos-vsrx_InterfacesInterfaceUnitFamilyInetAddressName" "InterfacesInterfaceUnitFamilyInetAddressName_st0" {
    	resource_name = "JTAF-IPSEC-ST-Interfaces2"
    	name = "st0"
    	name__1 = "0" # Unit number
    	name__2 = "172.10.1.1/30"
}

resource "junos-vsrx_SecurityZonesSecurity__ZoneInterfacesHost__Inbound__TrafficProtocolsName" "seczonesprotocol_st0" {
	resource_name = "JTAF-IPSEC-ST-Interfaces3"
	name = "untrust"
	name__1 = "st0.0"
	name__2 = "all"  # protocol all 
}

resource "junos-vsrx_SecurityZonesSecurity__ZoneInterfacesHost__Inbound__TrafficSystem__ServicesName" "seczonesservices_st0" {
	resource_name = "JTAF-IPSEC-ST-Interfaces4"
    	name = "untrust"
    	name__1 = "st0.0"
    	name__2 = "all" # services all
}