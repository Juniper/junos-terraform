This example is of a Data Center IPCLOS fabric with two leafs and two spine devices respectively. For this example we use junos-vqfx terraform provider for all the devices.
Each device has its own terraform file `main.tf` located in its respective folder which are: `vqfx_1`, `vqfx_2`, `vqfx_3` and `vqfx_4`. 

Objective is to keep the `main.tf` files of each device similar, and pass only the needed variables through the `root` module using `terraform.tfvars`. Variable types are controlled via 
`variables.tf` file located in each of the device folder.

Each device details needs to added in `main.tf` file located in the root directory, in this case under `dc` folder. The terraform apply command needs to executed from `dc` folder.

```
provider "junos-vqfx" {
  host     = "10.52.231.212"
  port     = 830
  username = "regress"
  password = "MaRtInI"
  sshkey   = ""
  alias    = "leaf1"
}
provider "junos-vqfx" {
  host     = "10.52.227.18"
  port     = 830
  username = "regress"
  password = "MaRtInI"
  sshkey   = ""
  alias    = "leaf2"
}
provider "junos-vqfx" {
  host     = "10.52.231.216"
  port     = 830
  username = "regress"
  password = "MaRtInI"
  sshkey   = ""
  alias    = "spine1"
}
provider "junos-vqfx" {
  host     = "10.52.231.214"
  port     = 830
  username = "regress"
  password = "MaRtInI"
  sshkey   = ""
  alias    = "spine2"
}
```
