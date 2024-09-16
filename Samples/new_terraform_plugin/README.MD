Terraform has introduced plugin framework to develop custom providers. This example uses the new terraform plugin framework to generate a provider for junos-vsrx. After successful installation of JTAF (using either option 1 or option 2 within generateFiles.sh script)
follow the below procedure to use new terraform plugin for provider generation.

Below files are to be copied into the respective folders within the junos-terraform folder.

    junos-terraform\
        |
        ------ terraform_providers\
                    |
                    ------------ config.go
                    |
                    ------------ provider.go
                    |
                    ------------ resource_InterfacesInterfaceUnitFamilyInetAddressName.go
                    |
                    ------------ resource_junos_destroy_commit.go
                    |
                    ------------ resource_junos_device_commit.go
                    |
                    ------------ main.go
        |
        ------- testbed\
                  |
                  ----------- main.tf
                  |
                  ----------- vsrx_1\
                                |
                                ------- main.tf

After copying the files, follow below procedure to generate the terraform provider for junos-vsrx

1. Run command `go mod tidy -go={GO_VERSION}` from junos-terraform\terraform_providers folder. This will install the dependencies needed for the new terraform plugin framework.
2. Execute `go build` command the junos-terraform\terraform_providers folder.
3. Above step will generate a binary file for the terraform provider, file would be named as `terraform_providers`. Rename it to `terraform-provider-junos-vsrx` using command `mv terraform_providers terraform-provider-junos-vsrx`.
4. Copy the newly generated terraform provider binary file into `~/.terraform.d/plugins/juniper/providers/junos-vsrx/<version>/(linux_amd64|darwin_arm64)` folder.

