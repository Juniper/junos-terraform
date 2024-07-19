
This is an example of 3 nodes (vMX) with ISIS/LDP/BGP protocols configuration. All three devices are interconnected to each other in a ring topology.

All devices have identical files (`main.tf` and `variables.tf`) within the respective folders. Only difference is in the value of the variables being passed from the root module. The `terraform.tfvars` file contains key value pair based attributes for all 3 VMX devices. 

3 devices use the same `terraform-provider-junos-vmx`, although each device is known distinctly with `alias` statement in their respective provider definition (as highighted below) within `main.tf` file.

```
provider "junos-vmx" {
  host     = "aa.aa.aa.aa"
  port     = 830
  username = "username"
  password = "password"
  sshkey   = ""
  alias    = "R1"
}
provider "junos-vmx" {
  host     = "bb.bb.bb.bb"
  port     = 830
  username = "username"
  password = "password"
  sshkey   = ""
  alias    = "R2"
}
provider "junos-vmx" {
  host     = "cc.cc.cc.cc"
  port     = 830
  username = "username"
  password = "password"
  sshkey   = ""
  alias    = "R3"
}
```
