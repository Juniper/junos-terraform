# VMX - Static VXLAN with IPSEC

This is an example that demonstrates how to use JTAF on the Static VXLAN Sandbox available in JCL. 

# Setup
* Create a JCL sandbox of the 'VMX - Static VXLAN with IPSEC' blueprint. Select 22.3 for the device version. 
* Install Go on the Helper VM by running the following commands:
```bash
    wget https://go.dev/dl/go1.23.0.linux-amd64.tar.gz
    sudo rm -rf /usr/local/go
    sudo tar -C /usr/local -xzf go1.23.0.linux-amd64.tar.gz
    echo "export PATH=\$PATH:/usr/local/go/bin" >> ~/.bashrc
    source ~/.bashrc
```
* Install Terraform on the Helper VM by running the following commands:  
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
* Remove the current generateFiles.sh and replace it with Samples/static_vxlan/generateFiles.sh
```bash
    rm generateFiles.sh
    mv Samples/static_vxlan/generateFiles.sh junos-terraform
```
This updates the two lines below. You could also use vim to update these lines. 
```bash
Line 198
# common_path="yang/$junos_version/$junos_version_combined/common/junos-common-types@2023-01-01.yang"
common_path="yang/$junos_version/$junos_version_combined/common/junos-common-types@2022-01-01.yang"
Line 299 
# go mod tidy -go=1.16 && go mod tidy -go=1.17
go mod tidy -go=1.23
```
* Move the vMX-A1 configuration file into to a directory titled user_config_files.
```bash
    mkdir user_config_files
    mv Samples/static_vxlan/vmx_config.xml user_config_files
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

Enter the configuration file name: vmx_config.xml
Enter a valid device option (vsrx, vmx, vqfx, vptx): vmx
Enter the Junos version: 22.3
```
This takes a LONG time, so don't worry if it is taking a while! That just means it is working.

# Using the Provider 
Yay! The provider has now been created and can be found in terraform_providers/terraform-provider-junos-vmx. 

Enable terraform to find the provider locally.
```bash
mkdir -p ~/.terraform.d/plugins/juniper/providers/junos-vmx/22.3/linux_amd64
mv terraform_providers/terraform-provider-junos-vmx ~/.terraform.d/plugins/juniper/providers/junos-vmx/22.3/linux_amd64
```
# Testing 
Testing with Terraform: https://github.com/Juniper/junos-terraform?tab=readme-ov-file#testing

For this portion, you can find the main.tf and vmx_1/main.tf files located in Samples/static_vxlan/testbed.