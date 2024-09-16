## Introduction ##
Terraform has introduced plugin framework to develop custom providers. This example uses the new terraform plugin framework to generate a provider for junos-vsrx. This example demonstrates configuring IP address to an interface. 

## Summary ##
After successful installation of JTAF (using either option 1 or option 2 within generateFiles.sh script) follow below procedure to use new terraform plugin for provider generation.

Below files are to be copied into the respective folders within the junos-terraform folder.


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
       |                                                                                             
       |                                                                                             
       |                                                                                             
       |                                                                                             
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
                                +--------                                                            
## Steps ##
After copying the above files, follow below procedure to generate the terraform provider for junos-vsrx

1. Run command `go mod tidy -go={GO_VERSION}` from junos-terraform\terraform_providers folder. This will install the dependencies needed for the new terraform plugin framework.
2. Execute `go build` command the junos-terraform\terraform_providers folder.
3. Above step will generate a binary file for the terraform provider, file would be named as `terraform_providers`. Rename it to `terraform-provider-junos-vsrx` using command `mv terraform_providers terraform-provider-junos-vsrx`.
4. Copy the newly generated terraform provider binary file into `~/.terraform.d/plugins/juniper/providers/junos-vsrx/<version>/(linux_amd64|darwin_arm64)` folder.

Now that new terraform provider has been build and copied into the respective location, IP address configuration can be pushed to Junos device using terraform. Before pushing the configuration, any old state needs to be removed with the help below commands that needs to be executed within testbed folder.

        `rm -rf .terraform`
        `rm .terraform.lock.hcl`

With a clean slate below commands can now be executed from `junos-terraform\testbed` folder.

        `terraform init`
        `terraform plan`
        `terraform apply -auto-approve`
