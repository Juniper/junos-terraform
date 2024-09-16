## Introduction ##
Terraform has introduced a new plugin framework to develop custom providers. This example uses new terraform plugin framework to generate a provider for junos-vsrx. This example demonstrates configuring IP address on an interface. 

## Summary ##
After successful installation of JTAF (using either option 1 or option 2 provided by generateFiles.sh script) follow below procedure to use new terraform plugin for provider generation.

Below files are to be copied into the respective folders within the junos-terraform folder.

```
junos-terraform\                                                                                     
       |                                                                                             
       |                                                                                             
       |                                                                                             
       +----------- terraform_providers\                                                             
       |                   |                                                                         
       |                   |                                                                         
       |                   +--------------- main.go                                                  
       |                   |                                                                         
       |                   |                                                                         
       |                   |                                                                         
       |                   +--------------- config.go                                                
       |                   |                                                                         
       |                   |                                                                         
       |                   |                                                                         
       |                   +---------------- provider.go                                             
       |                   |                                                                         
       |                   |                                                                         
       |                   |                                                                         
       |                   +---------------- resource_InterfacesInterfaceUnitFamilyInetAddressName.go
       |                   |                                                                         
       |                   |                                                                         
       |                   |                                                                         
       |                   +---------------- resource_junos_destroy_commit.go                        
       |                   |                                                                         
       |                   |                                                                         
       |                   |                                                                         
       |                   +---------------- resource_junos_device_commit.go                         
       |                   |                                                                         
       |                   |                                                                         
       |                   |                                                                         
       |                   +---------------- resource.go                                                        
       |                                                                                             
       +----------testbed\                                                                           
                    |                                                                                
                    |                                                                                
                    +--------- main.tf                                                               
                    |                                                                                
                    |                                                                                
                    |                                                                                
                    +---------vsrx_1\                                                                
                                |                                                                    
                                |                                                                    
                                +--------main.tf                                                     
```
## Steps ##
After copying the above files, follow below procedure to generate the terraform provider for junos-vsrx

1. Run command `go mod tidy -go={GO_VERSION}` from junos-terraform\terraform_providers folder. This will install the dependencies needed for the new terraform plugin framework.

2. Execute `go build` command from junos-terraform\terraform_providers folder.

3. Above steps will generate a binary file for the terraform provider, file would be named as `terraform_providers`. Rename it to `terraform-provider-junos-vsrx` using command `mv terraform_providers terraform-provider-junos-vsrx`.

4. Copy the newly generated terraform provider binary file into `~/.terraform.d/plugins/juniper/providers/junos-vsrx/<version>/(linux_amd64|darwin_arm64)` folder.

### Clean up old state ###
Now that new terraform provider has been build and copied into the respective location, IP address configuration can be pushed to Junos device using terraform. Before pushing the configuration, any old state needs to be removed with the help below commands. These need to be executed within `testbed` folder.

        rm -rf .terraform
        rm .terraform.lock.hcl

### Push Configuration ###
With a clean slate now, below commands can be executed from `junos-terraform\testbed` folder.

        terraform init
        terraform plan
        terraform apply -auto-approve
