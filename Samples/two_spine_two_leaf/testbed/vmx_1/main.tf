
terraform {
	required_providers {
		junos-vmx = {
			source = "juniper/providers/junos-vmx"
			version = "22.3"
		}
	}
}

resource "junos-vmx_InterfacesInterfaceName" "vmx_1" {
	resource_name = "vmx_1"
	name = "/interfaces/interface/vmx_1"
}

resource "junos-vmx_InterfacesInterfaceUnit" "vmx_2" {
	resource_name = "vmx_2"
	name = "/interfaces/interface/name"
}

resource "junos-vmx_InterfacesInterfaceUnitFamilyInetAddress" "vmx_3" {
	resource_name = "vmx_3"
	name = "/interfaces/interface/name"
	name__1 = "/interfaces/interface/unit/name"
}

resource "junos-vmx_SystemBackup__Router" "vmx_4" {
	resource_name = "vmx_4"
}

resource "junos-vmx_SystemBackup__RouterAddress" "vmx_5" {
	resource_name = "vmx_5"
	address = "/system/backup-router/address"
}

resource "junos-vmx_InterfacesInterfaceUnitName" "vmx_6" {
	resource_name = "vmx_6"
	name = "/interfaces/interface/name"
	name__1 = "/interfaces/interface/unit/name"
}

resource "junos-vmx_InterfacesInterfaceUnitFamily" "vmx_7" {
	resource_name = "vmx_7"
	name = "/interfaces/interface/name"
	name__1 = "/interfaces/interface/unit/name"
}

resource "junos-vmx_InterfacesInterfaceUnitFamilyInet" "vmx_8" {
	resource_name = "vmx_8"
	name = "/interfaces/interface/name"
	name__1 = "/interfaces/interface/unit/name"
}

resource "junos-vmx_InterfacesInterfaceUnitFamilyInetAddressName" "vmx_9" {
	resource_name = "vmx_9"
	name = "/interfaces/interface/name"
	name__1 = "/interfaces/interface/unit/name"
	name__2 = "/interfaces/interface/unit/family/inet/address/name"
}

resource "junos-vmx_SystemHost__Name" "vmx_10" {
	resource_name = "vmx_10"
	host__name = "/host-name"
}

resource "junos-vmx_InterfacesInterface" "vmx_11" {
	resource_name = "vmx_11"
}

