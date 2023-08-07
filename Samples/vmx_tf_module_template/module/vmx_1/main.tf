terraform {
  required_providers {
    junos-vmx = {
      source = "juniper/providers/junos-vmx"
      version = "20.41.12"
    }
  }
}

resource "junos-vmx_Logical__SystemsName" "vmx_1" {
    resource_name = "vmx_1"
    name = "logical-system1"
}