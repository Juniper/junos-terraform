
This is an example of 3 nodes (vMX) with ISIS/LDP/BGP protocols configuration. All three devices are interconnected to each other in a ring topology.

All devices have identical files (`main.tf` and `variables.tf`) within the respective folders. Only difference is in the value of the variables being passed from the root module. The `terraform.tfvars` file contains key value pair based attributes for all 3 VMX devices. 

