terraform {
  required_providers {
    junos-vsrx = {
      source = "juniper/providers/junos-vsrx"
      version = "21.31.108"
    }
  }
}


resource "junos-vsrx_Interfaces" "vsrx_2" {
  resource_name = "example_resource"
  name = "ge-0/0/0"  
  description = "Main Ethernet"
  mtu = 9192
  vlan_tagging = true
  
}
