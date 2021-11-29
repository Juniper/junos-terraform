terraform {
  required_providers {
    junos-vsrx = {
      source = "juniper/providers/junos-vsrx"
      version = "20.32.0101"
    }
  }
}

// To test Terraform's ability to remove a single resource,
// comment out one resource below and use `terraform taint` on the commit
// resource so that an apply also applies the commit.
resource "junos-vsrx_InterfacesInterfaceDescription" "vsrx_1" {
    resource_name = "vsrx_1"
    name = "ge-0/0/0"
    description = "LEFT-UNTRUST"
}

resource "junos-vsrx_InterfacesInterfaceUnitFamilyInetAddressName" "vsrx_2" {
    resource_name = "vsrx_2"
    name = "ge-0/0/0"
    name__1 = "0"
    name__2 = "10.0.0.0/24"
}


resource "junos-vsrx_Policy__OptionsPolicy__StatementTermFromRoute__FilterAddress" "vsrx_3" {
    resource_name = "vsrx_3"
    name = "DEF-IMPORT-FWTRUST-TABLE"
    name__1 = "t1"
    address = "10.0.0.0/17"
    orlonger = " "
}

resource "junos-vsrx_Policy__OptionsPolicy__StatementTermThenAccept" "vsrx_4" {
    resource_name = "vsrx_4"
    name = "DEF-IMPORT-FWTRUST-TABLE"
    name__1 = "t1"
    accept = ""
}

resource "junos-vsrx_Policy__OptionsPolicy__StatementTermThenReject" "vsrx_5" {
    resource_name = "vsrx_5"
    name = "DEF-IMPORT-FWTRUST-TABLE"
    name__1 = "t2"
    reject = ""
}

resource "junos-vsrx_SecurityAddress__BookAddressIp__Prefix" "vsrx_6" {
    resource_name = "vsrx_6"
    name = "global"
    name__1 = "eastus2-aks1-az1"
    ip__prefix = "10.0.1.0/24"
}

resource "junos-vsrx_SecurityAddress__BookAddressIp__Prefix" "vsrx_7" {
    resource_name = "vsrx_7"
    name = "global"
    name__1 = "eastus2-aks1-az2"
    ip__prefix = "10.0.3.0/24"
}

resource "junos-vsrx_SecurityAddress__BookAddress__SetAddressName" "vsrx_8" {
    resource_name = "vsrx_8"
    name = "global"
    name__1 = "set1"
    name__2 = "eastus2-aks1-az1"
}

resource "junos-vsrx_SecurityAddress__BookAddress__SetAddressName" "vsrx_9" {
    resource_name = "vsrx_9"
    name = "global"
    name__1 = "set1"
    name__2 = "eastus2-aks1-az2"
}

resource "junos-vsrx_SecurityNatSourcePoolAddressToIpaddr" "vsrx_10" {
    resource_name = "vsrx_10"
    name = "ut-pool"
    name__1 = "10.1.1.1/32"
    ipaddr = "10.1.1.1/32"
}

resource "junos-vsrx_ApplicationsApplicationProtocol" "vsrx_11" {
    resource_name = "vsrx_11"
    name = "tcp-993"
    protocol = "tcp"
}

resource "junos-vsrx_ApplicationsApplicationDestination__Port" "vsrx_12" {
    resource_name = "vsrx_12"
    name = "tcp-993"
    destination__port = "993"
}
