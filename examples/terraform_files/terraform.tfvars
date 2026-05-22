# Copy this file to terraform.tfvars (or another *.tfvars) and adjust values per environment.

junos_username = "jcluser"
junos_password = "Juniper!1"

junos_vqfx_endpoints = {
  dc1_spine1      = { host = "66.129.234.204",  port = 45046 }
  dc1_spine2      = { host = "66.129.234.204",  port = 45054 }
  dc1_borderleaf1 = { host = "66.129.234.204",  port = 45003 }
  dc1_borderleaf2 = { host = "66.129.234.204",  port = 45028 }
  dc1_leaf1       = { host = "66.129.234.204",  port = 45069 }
  dc1_leaf2       = { host = "66.129.234.204",  port = 45080 }
  dc1_leaf3       = { host = "66.129.234.204",  port = 45084 }
  dc2_spine1      = { host = "66.129.234.204",  port = 45088 }
  dc2_spine2      = { host = "66.129.234.204",  port = 45092 }
}

junos_vsrx_endpoints = {
  dc1_firewall1 = { host = "66.129.234.204", port = 45007 }
  dc1_firewall2 = { host = "66.129.234.204", port = 45032 }
  dc2_firewall1 = { host = "66.129.234.204", port = 45050 }
  dc2_firewall2 = { host = "66.129.234.204", port = 45058 }
}
