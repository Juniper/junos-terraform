**INTRODUCTION**

This is an example which exhibits a way to templatize the terraform files written for Junos OS devices. Terraform allows to declare input variables in a file `variables.tf` and value for each variable can be passed in separate file called as `terraform.tfvars`.

**Modules**

Terraform has modules, which are basically containers to store multiple `.tf` files. In junos-terraform, each device is defined as a module e.g. vsrx_1 hence a folder is created to package all `.tf` files associated with that device.

1. Root Module - Each Terraform configuration has atleast one module, which is known as "root module". Basically root module is path from where `terraform plan/apply` command is run.
2. Child Module - A module that can be called by other modules is called as child modules, in junos-terraform case vsrx_1 is a child module. It is called from the root module within `main.tf`.

Input variables can be declared within the root module in file named `variables.tf` and can be passed over to the child module. Within the child module similar variables must be defined within `variables.tf` file as well. Values from root module can be passed to child module, when the module is being called. In this case `bgp` variable declared in root module is being passed to vsrx_1 child module within `main.tf` file like below:

```
module "vsrx_1" {
  source    = "./vsrx_1"
  bgp       = var.bgp
  providers = {junos-vsrx = junos-vsrx}
  depends_on = [junos-vsrx_destroycommit.commit-main]
}
```
In the statement `bgp = var.bgp` left side variable `bgp` is of child module, where as right side variable `bgp` is of root module. And to access the value of `bgp` variable of root module `var.bgp` is used. Within the vsrx_1 child module, variable `bgp` is defined with its arguments and type constraints but variable `bgp` in the root module is not defined with constraints etc. This is done so that modules can be reused easily.

Value of `bgp` variable is defined using `terraform.tfvars` file within the root module. `main.tf` file within the vsrx_1 module can refer to the value of `bgp` variable using `var.bgp`. And each resource statement can be defined once and argument to each of its attribute can be substituted by variable `bgp` using `for_each` expression of Terraform. Like in example below the resource sets the BGP neighbor type for each of BGP group.

```
resource "junos-vsrx_ProtocolsBgpGroupType" "vsrx_r1" {
  resource_name = "my_group"
  for_each = var.bgp
  name          = each.key
  type          = each.value.neighbor_type
}
```

Since nested for loops are not possible within Terraform, hence a local variable `bgp_neighbor_list` is defined within `vsrx_1/main.tf` this creates a new object which is specifically for BGP neighbors and attributes of each neighbor. In this example, its limited to name, description and peer-as only but can be expanded as well. 
