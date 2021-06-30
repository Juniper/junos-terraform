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

## PRE-REQUISITE Tools 

The following tools needs to be installed
* [Terraform](https://learn.hashicorp.com/tutorials/terraform/install-cli)  
* golang
* pyang

## Steps to Generate the Terraform Api's

### Step 1 : Copy Yang Files

Copy the compatible yang files to local setup.
Yang files for junos devices can be downloaded from the following github repo -
``https://github.com/Juniper/yang.git`` 

### Step 2 : Generate yin and xpath files

Generate the yin file and xpath for the yangs. Execute the following processYang binary -

`` cd cmd/processyang `` 

`` go build ``

`` ./processYang -config /var/tmp/config.toml ``

[Note] Sample config file is provided in jtaf/Samples/config.toml

### Step 3 : Update the Input Api Details

Update the xpath.xml file with the xpaths and group details. Refer xpaths generated in 
Step: 2 for the same.

### Step 4 : Generate the providers

Generate the terraform api by executing the following command -

`` cd cmd/processyang `` 

`` go build ``

`` ./processProviders -config /var/tmp/config.toml ``

The terraform api will be generated in the repository provided in config file. 

[Note] It is recommended to use terraform_providers present in jtaf directory for generating providers. 
Otherwise copy the files present in terraform_providers to the directory where modules are generated

### Step 5 : Create terraform binary 

Create the terraform binary by executing below command at the terraform modules generated repository

``go build -o terraform-provider-junos-qfx``


## CONTRIBUTORS
Juniper Networks is actively contributing to and maintaining this repo.
 
*Contributors:*

* [Rahul Kumar](https://github.com/rahkumar651991)
* [David Gee](https://github.com/davedotdev)
