# VMX - Two Spine Two Leaf

This is an example that demonstrates how to use JTAF on a two spine two leaf VMX setup. This topology was built using Apstra, which is an automation that builds data center fabrics from the ground up. Currently, this topology is two spine two leaf only, but it can be easily scaled horizantally and vertically. 
Here are the details of this example:
* Spine1: 10.56.16.246
* Spine2: 10.56.12.9
* Leaf1: 10.56.17.17
* Leaf2: 10.56.16.194

Credentials:
* Username: regress
* Password: MaRtInI


# Setup
* Install Go on your machine by running the following commands. Commands may vary by machine.
```bash
    wget https://go.dev/dl/go1.23.0.linux-amd64.tar.gz
    sudo rm -rf /usr/local/go
    sudo tar -C /usr/local -xzf go1.23.0.linux-amd64.tar.gz
    echo "export PATH=\$PATH:/usr/local/go/bin" >> ~/.bashrc
    source ~/.bashrc
```
* Install Terraform on your machine by running the following commands. Commands may vary by machine.
```bash
    wget https://releases.hashicorp.com/terraform/1.5.0/terraform_1.5.0_linux_amd64.zip
    unzip terraform_1.5.0_linux_amd64.zip
    sudo mv terraform /usr/local/bin/
```
* Clone the JTAF repository and move into the JTAF directory. 
```bash
    git clone https://github.com/Juniper/junos-terraform.git
    cd junos-terraform
```
* Remove the current generateFiles.sh and replace it with Samples/two_spine_two_leaf/generateFiles.sh
```bash
    rm generateFiles.sh
    mv Samples/two_spine_two_leaf/generateFiles.sh .
```
This updates the two lines below. You could also use vim to update these lines. 
```bash
# Line 198
# common_path="yang/$junos_version/$junos_version_combined/common/junos-common-types@2023-01-01.yang"
common_path="yang/22.3/22.3R1/common/junos-common-types@2022-01-01.yang"
# Line 299 
# go mod tidy -go=1.16 && go mod tidy -go=1.17
go mod tidy -go=1.23
```
* Move the VMX configuration file into to a directory titled user_config_files.
```bash
    mkdir user_config_files
    mv Samples/two_spine_two_leaf/config_files/spine1-config.xml user_config_files
```
**OR** if you feel adventurous, you can retrieve the configuration from the vmx device itself.
```bash
show configuration | display xml | no-more
```
**Note:** you must remove the rpc-reply and configuration level of the file. 

# Run Option 2
Junos-Terraform Developer Guide (Build from pre-existing junos config): https://github.com/Juniper/junos-terraform?tab=readme-ov-file#guide

As a brief summary, run generateFiles.sh.
```bash
./generateFiles.sh
Do you want to:
1. Build a provider from scratch
2. Provide a configuration
Enter your choice (1/2): 2
You chose to provide a configuration.

Enter the configuration file name: spine1-config.xml
Enter a valid device option (vsrx, vmx, vqfx, vptx): vmx
Enter the Junos version: 22.3
```
This takes a LONG time, so don't worry if it is taking a while! That just means it is working.

# Using the Provider 
Yay! The provider has now been created and can be found in terraform_providers/terraform-provider-junos-vmx. 

Enable terraform to find the provider locally. 
```bash
mkdir -p ~/.terraform.d/plugins/juniper/providers/junos-vmx/22.3/linux_amd64
mv terraform_providers/terraform-provider-junos-vmx ~/.terraform.d/plugins/juniper/providers/junos-vmx/22.3/darwin_arm64
```
# Testing 
Testing with Terraform: https://github.com/Juniper/junos-terraform?tab=readme-ov-file#testing

For this portion, you can find the main.tf and vmx_1/main.tf files located in Samples/static_vxlan/testbed.

An example change is changing a device name in vmx_1/main.tf. You can see this change reflected in the configuration of the device. 
```bash
resource "junos-vmx_InterfacesInterfaceName" "vmx_1_new_name" {
	resource_name = "vmx_1_new_name"
	name = "/interfaces/interface/vmx_1_new_name"
}
```