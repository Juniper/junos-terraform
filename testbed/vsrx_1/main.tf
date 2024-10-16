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
  name          = "ge-0/0/0"
  description   = "Main Ethernet"
  mtu           = 9192
  vlan_tagging = true
  units = [
    {
      name        = "0"
      description = "unit_desc"
      vlan_id = 10
      family = {
          inet = [
            {
              address = "192.168.20.1/24"
            }
          ]
          inet6 = [
            {
              address = "2400:9800:7000:1009:1000::/127"
            }
          ]
        }
    },
    {
      name        = "1"
      description = "unit_desc"
      vlan_id = 20
      family = {
          inet = [
            {
              address = "192.169.20.1/24"
            }
          ]
          inet6 = [
            {
              address = "2400:9800:7000:100a:1000::/127"
            }
          ]
      }
    },
  ]
}