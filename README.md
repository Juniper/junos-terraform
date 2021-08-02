## JUNOS Terraform Automation Framework (JTAF)

Notes contained herein are for the custom development of a Terraform automation framework for Junos.
This is a work under development.

For more information on the JTAF project, you can watch the videos below.

Introduction: https://youtu.be/eH24eCZc7pE

Installation: https://youtu.be/aTF7_Uscd9Q

Generate: https://youtu.be/UgsFU7UplRE

Execution: https://youtu.be/Lfkc38wzhNg

Interface Configuration: https://youtu.be/iCnnkDodUgQ

BGP Configuration: https://youtu.be/nQVNCNCJZRc

There is a detailed README which takes you through building a complete example Provider from scratch. That can be found here:
https://github.com/Juniper/junos-terraform/blob/master/DETAILED_INSTRUCTIONS.md

## Pre-Requisite Tools 

The following tools needs to be installed
* [Terraform](https://learn.hashicorp.com/tutorials/terraform/install-cli)  
* golang >= 1.12
* pyang

## Steps to Generate the Terraform APIs

### Step 1 : Copy YANG Files

Copy the compatible yang files to local setup.
Yang files for junos devices can be downloaded from the following github repo -
``https://github.com/Juniper/yang.git`` 

### Step 2 : Generate YIN and XPath Files

Generate the yin file and xpath for the yangs. Execute the following processYang binary -

`` cd cmd/processYang `` 

`` go build ``

`` ./processYang -config /var/tmp/config.toml ``

[Note] Sample config file is provided in jtaf/Samples/config.toml

### Step 3 : Update the Input API Details

Update the xpath.xml file with the xpaths and group details. Refer xpaths generated in 
Step: 2 for the same.

### Step 4 : Generate the Providers

Generate the terraform api by executing the following command -

`` cd cmd/processProviders `` 

`` go build ``

`` ./processProviders -config /var/tmp/config.toml ``

The terraform api will be generated in the repository provided in config file. 

[Note] It is recommended to use terraform_providers present in jtaf directory for generating providers. 
Otherwise copy the files present in terraform_providers to the directory where modules are generated

### Step 5 : Create Terraform Binary 

Create the terraform binary by executing below command at the terraform modules generated repository. The binary name must be same as below.

``go build -o terraform-provider-junos-device``


## CONTRIBUTORS
Juniper Networks is actively contributing to and maintaining this repo.
 
*Contributors:*

* [Rahul Kumar](https://github.com/rahkumar651991)
* [David Gee](https://github.com/davedotdev)
