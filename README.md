# JUNOS Terraform Automation Framework (JTAF)

Terraform is traditionally used for managing virtual infrastructure, but there are organisations out there that use Terraform end-to-end and also want to manage configuration state using the same methods for managing infrastructure. Sure, we can run a provisioner with Terraform, but that wasn't asked for!

Much the same as you can use Terraform to create an AWS EC2 instance, you can manage the configurational state of Junos. In essence, we treat Junos configuration as declarative resources.

So what is JTAF? It's a framework, meaning, it's an opinionated set of tools and steps that allow you to go from YANG models to a custom Junos Terraform provider. With all frameworks, there are some dependencies.

To use JTAF, you'll need machine that can run **Go, Python, Git and Terraform.** This can be Linux, OSX or Windows. Some easy to consume videos are below.

Introduction: https://youtu.be/eH24eCZc7pE

Installation: https://youtu.be/aTF7_Uscd9Q

Generate: https://youtu.be/UgsFU7UplRE

Execution: https://youtu.be/Lfkc38wzhNg

Interface Configuration: https://youtu.be/iCnnkDodUgQ

BGP Configuration: https://youtu.be/nQVNCNCJZRc

# Section Breakdown
* [Junos-Terraform Demo (Build from scratch)](#demo)
* [Junos-Terraform Developer Guide (Build from pre-existing junos config)](#guide)
* [Using the Provider](#provider)
* [Testing with Terraform](#testing)
* [Question & Answers/ Common Problems](#questions)

<a id="demo"></a>
# Junos-Terraform Demo

Following this demo, developers will learn to be able to manually create custom junos-terrraform providers with the ability to configure any junos device given explicitly defined xpaths and yang files.

In this demo specifically, we will create a simple provider for a vSRX of version 19.4R1 that has the capability to do two things:
* Add a description to an interface
* Add an `inet` address to a sub-interface
  
Don't worry about commits just yet, we'll discuss that later. This demo is for the purpose of creating terrafom providers from scratch with custom device capabilities related to the use of Juniper's yang file resources.

In this document, if you see `$JTAF_PROJECT` you can replace it with the path of the JTAF project on your system or create an environment variable. On the author's system, this happens to be below:

```bash
cd Documents/GoDev/src/github.com/Juniper
git clone https://github.com/Juniper/junos-terraform.git
export JTAF_PROJECT=/Users/dgee/Documents/GoDev/src/github.com/Juniper/junos-terraform

# Let's check
echo $JTAF_PROJECT
/Users/dgee/Documents/GoDev/src/github.com/Juniper/junos-terraform
```

Having this variable set, helps to navigate the system without lots of tedious typing.

## Install Python, Go & Terraform

Go version tested with: `go1.14.2 darwin/amd64`
Python version tested with: `3.7`
Terraform version: `v0.12.26`

Other versions beyond these will work, but this is what was tested for the writing of this document.

## Copy the YANG Files 
For our example, our provider will only be able to create an interface description and place an inet address on a sub-interface. We only need a handful of YANG models for this.

 > For developers wanting to add more capabilities to the provider, they will need to also add the neccesary `yang_files` required for those capabilites outlined by the `xpath` inputs which are listed in the `xpath_inputs.xml` file created later on 
 > * say you want to add firewalls or policy-option options in addition to interfaces; you will need to add the yang files along with the xpaths associated with the custom capability (See `/Samples/vsrx_tf_module_template` directory for examples of an xpath file, `xpath_example.xml`, and the corresponding `yang_files` directory)

Let's put the YANG files from [Juniper's YANG GitHub repository](https://github.com/Juniper/yang.git) in to a memorable location. Let's use the home directory `/junos-terraform`. Don't worry, once the necesary files are copied over, the `/yang` folder can be removed.

```bash
/junos-terraform $ [from this directory]

git clone https://github.com/Juniper/yang.git
# this can take some time depending on your internet connection

mkdir yang_files

# Note, the common-types and conf-root YANG models are dependencies
cp yang/19.4/19.4R1/common/junos-common-types@2019-01-01.yang ./yang_files
cp yang/19.4/19.4R1/junos-es/conf/junos-es-conf-root@2019-01-01.yang ./yang_files
# Insert the models required here
cp yang/19.4/19.4R1/junos-es/conf/junos-es-conf-interfaces@2019-01-01.yang ./yang_files
```

If you wanted to remove the YANG directory, you can do it like this:

```bash
cd /junos-terraform
rm -rf yang
```

# Run First Shell Script:  `generateFiles.sh`

**Prior to this step, ensure python and go is installed.**

>This file can be compiled by running `chmod +x generateFiles.sh` from the `/junos-terraform` directory followed by 
`./generateFiles.sh` to run the script. 

`Select "Build a provider from scratch" option: [1] out of the options if following along with the demo`

## Below describes what the script does:


### 1. Generates a `config.toml` File

*If you've never seen a TOML file before, don't worry! It's just a structured file containing configuration that can be parsed by a program, in this case the two main compiled programs that form JTAF. TOML stands for Tom's Obvious Minimal Langage.*

Creates a config file in the home directory. Don't worry about the xPath or fileType keys. They'll be explained shortly.
You can find this file `config.toml` in the home directory (/junos-terraform)

```bash
yangDir = "$(pwd)/yang_files"
providerDir = "$(pwd)/terraform_providers"
xpathPath = "$(pwd)/xpath_inputs.xml"
providerName = "vsrx"
fileType = "both"
```

You can also replace the fileType field to `text` or `xml`. The text files are for us humans.

### 2. Generates the YIN and XPath Files based on YANG files

The next step, depending on the size of YANG model/s, **may take some time**. Prepare some popcorn!
This step will activate a python vitual enviornment (make sure python is downloaded) and install `pyang` so it can be used
to generate the `yin` files.

```bash
cd $JTAF_PROJECT/cmd/processYang
go build
./processYang -config /path_to_junos-terraform/junos-terraform/config.toml
# OUTPUT - WARNING >> This can take some time. Lack of activity does not mean broken!

                    ___ _____ ___  ______
                   |_  |_   _/ _ \ |  ___|
                     | | | |/ /_\ \| |_
                     | | | ||  _  ||  _|
                 /\__/ / | || | | || |
                 \____/  \_/\_| |_/\_|
                           0.1.5
-------------------------------------------------------------------------
- Creating Yin files from Yang file directory: /path_to_junos-terraform/junos-terraform/yang_files -
-------------------------------------------------------------------------
Yin file for junos-common-types@2019-01-01 is generated
Yin file for junos-es-conf-interfaces@2019-01-01 is generated
Yin file for junos-es-conf-root@2019-01-01 is generated
--------------------------------------------
- Creating _xpath files from the Yin files -
--------------------------------------------
Creating Xpath file: junos-common-types@2019-01-01_xpath.txt
Creating Xpath file: junos-es-conf-interfaces@2019-01-01_xpath.txt
Creating Xpath file: junos-es-conf-root@2019-01-01_xpath.txt
```

At this point, `venv` is `deactivated` and the first script has terminated.

# Create an XPath Input XML File

Great, at this point now we have text file and YIN versions of the YANG files. We need those for the next step.

Let's create a file, which provides a list of inputs to the part of JTAF which writes the `.go` code automagically.
This input identifies the content of the provider that JTAF will create.

> For developers wanting to add more capabilities to the provider, this is where the additional xpath inputs need to be added. Assuming that the neccesary `yang_files` required for those capabilites outlined by the `xpath` inputs are added, the inputs can be incorporated into this file following the format below.
 > * Again, examples of this implementation cam be found in the `/Samples/vsrx_tf_module_template` directory which include examples of an xpath file, `xpath_example.xml`, and the corresponding `yang_files` directory)

Create a file `/junos-terraform/xpath_inputs.xml` and populate it with the content below.

## **Make sure that the file is named `xpath_inputs.xml`** 
* `The name must match the name defined in config.toml`

```bash
<file-list>
        <xpath name="/interfaces/interface/description"/>
        <xpath name="/interfaces/interface/unit/family/inet/address/name"/>
</file-list>
```

A simple explanation of the above XPaths:

* The first xpath entry is for the interface description and references a YANG leaf.
* The second xpath entry identifies the inet YANG leaf in the YANG model.

You can view these expressions as a simple way to identify the fields inside the YANG model we're interested in.
JTAF generated providers has a requirement of the smallest data set possible for each resources. That means, in a single resource you would place a description, and in another, you will place the inet address. Terraform is essentially a dependency aware declarative resource manager and so, we have to model resources in a way that's compatible with Terraform and Junos.


# Run Second Shell Script:  `buildProvider.sh`

> ### **Prior to this step**
> * IF you want to test the output of the configuration from the `terraform test` files, set the ENV variable `MOCK_FILE` to the path of an empty `xml` file you create where the system can display the expected config defined in the .tf files. 
>   * example: create a `jtaf_output.xml` file in the `/junos-terraform` directory and run `export MOCK_FILE=/path_to_junos-terraform/jtaf_output.xml`)
> * **Warning:** `MOCK_FILE` should be `unset` unless wanting to enter Mock mode, otherwise system will look for path setup in the terraform `.tf` test files created later on. `Mock Mode allows developers to test terraform commands and commits to a local file prior to device communication.` If `MOCK_FILE` is set, terraform commands will output to the file declared by the variable and not the device declared in the `main.tf` file.

### Now let's run the script: 
>This file can be compiled by running `chmod +x buildProvider.sh` from the home directory followed by 
`./buildProvider.sh` to run the script. 

Below describes what the script does:

## 1. Builds the Provider Resources

First, we need JTAF to create some `.go` code from the YANG models and XML data we provided.

```bash
cd $JTAF_PROJECT/cmd/processProviders
go build
./processProviders -config /path_to_config/config.toml
# This next step is rapid

                    ___ _____ ___  ______
                   |_  |_   _/ _ \ |  ___|
                     | | | |/ /_\ \| |_
                     | | | ||  _  ||  _|
                 /\__/ / | || | | || |
                 \____/  \_/\_| |_/\_|
                           0.1.5
------------------------------------------------------------
- Autogenerating Terraform Provider code from _xpath files -
------------------------------------------------------------
Terraform API resource_InterfacesInterfaceDescription created
Terraform API resource_InterfacesInterfaceUnitFamilyInetAddressName created
# --------------------------------------------------------------------------------
Number of Xpaths processed:  2
Number of potential issues:  0

---------------------------------------------
- Copying the rest of the required Go files -
---------------------------------------------
Copied file: config.go to /path_to_junos-terraform/junos-terraform/terraform_providers
Copied file: main.go to /path_to_junos-terraform/junos-terraform/terraform_providers
Copied file: resource_junos_destroy_commit.go to /path_to_junos-terraform/junos-terraform/terraform_providers
Copied file: resource_junos_device_commit.go to /path_to_junos-terraform/junos-terraform/terraform_providers
-------------------
- Creating Go Mod -
-------------------
```

The output of this step is written to the `/junos-terraform/terraform_provider` directory. Let's build the provider!


## 2. Creates updated xpath input file (if needed)

If there are any found issues during the building of the provider resources, a new file called `updated_xpath_inputs.xml` will be created with a trimmed version of the provided xpath inputs file. This file can be used to replace the data in the `xpath_inputs.xml`.


## 3. Builds the Provider

```bash
cd /terraform_providers
go build
```
This provider without any Go cross-compilation directives, will work on the system it's been generated with. If you happen to be on an OSX machine, then the provider will work for Terraform on OSX and the same is true for Linux, if you use JTAF on Linux, then natively the generated provider will operate on Linux. However, you can cross-compile the provider so that it will operate on another operating system and even CPU architecture.

```bash
# Validate the file kind
file terraform-provider-junos-device
# terraform-provider-junos-device: Mach-O 64-bit executable x86_64
```

If you want this provider to work with Linux, then you can cross-compile using the `GOOS` input. See below.

```bash
GOOS=linux go build -o terraform-provider-junos-device

file terraform-provider-junos-device

# terraform-provider-junos-device: ELF 64-bit LSB executable, x86-64, version 1 (SYSV), statically linked, Go BuildID=SWvAslM7UiUlMNJJOG8f/MV8jDWinx0vKkuo7Zmec/-2fk9ZDz88J7folCoc0q/ftWLT5N4tiPWQ8DlXY2J, not stripped
```

The binary file `terraform-provider-junos-vsrx` is actually our fresh new and shiny Terraform Provider. If you got this far, congratulations. You just created a Terraform provider for Junos and you are ready to use it. Jump to [Using the Provider](#provider)


<a id="guide"></a>
# Junos-Terraform Developer Guide

This section is aimed at users who are a little more comfortable with the way junos-terraform works and want to develop, test, configure Juniper devices with a pre-existing configuration.

> **The only requirement for this section is to upload a configuration in `xml` format to the `/junos-terraform/user_config_files` folder.**
> * When adding this `xml` file to the folder, remove any starting and ending `<configuration>` and `<cli>` tags.
> * For reference, see `test.xml` in the `/Samples/user_config_files` folder which contains configuration for `vqfx` spine for Junos version `23.1`
> * After configuration is uploaded, the rest is taken care of (provider build and test file templates created)

### **Prior to start (for testing)**
* IF you want to test the output of the configuration from the `terraform test` files, set the ENV variable `MOCK_FILE` to the path of an empty `xml` file you create where the system can display the expected config defined in the .tf files. 
   * example: create a `jtaf_output.xml` file in the `/junos-terraform` directory and run `export MOCK_FILE=/path_to_junos-terraform/jtaf_output.xml`)

* **Warning:** `MOCK_FILE` should be `unset` unless wanting to enter Mock mode, otherwise system will look for path setup. `Mock Mode allows developers to test terraform commands and commits to a local file prior to device communication.` If `MOCK_FILE` is set, terraform commands will output to the file declared by the variable and not the device declared in the `main.tf` file.

# Run Shell Script: `generateFiles.sh`

**Prior to this step, ensure python and go is installed.**

> This file can be compiled by running `chmod +x generateFiles.sh` from the home directory followed by `./generateFiles.sh` to run the script. 

`Select "Provide a configuration" option: [2] out of the options if wanting to create a resource provider based on an existing configuration`

Below describes what the script does:

### 1. Generates a `config.toml` File

* Creates a config file in the home directory. Don't worry about the xPath or fileType keys. They'll be explained shortly.
* You can find this file `config.toml` in the home directory `/junos-terraform`

### 2. Copies YANG Files 

* This part of the script calls for the cloning of the Juniper `yang` directory in order to copy necessary yang files for the given `device` and `version`. **This will only happen if the folder `yang_files` does not exist already.**

### 3. Generates the YIN and XPath Files based on YANG files

* This part of the script enables the call to create YIN and Xpath files. For more details, look at the demo for this section.

### 4. Creates an XPath Input XML File

* This part of the script uses a go file called `createXpathInputs.go` to parse the configuration loaded by the user and creates an `xpath_inputs.xml` file. This file contains a skeleton of all the configured xpaths contained in the pre-loaded configuration. These xpaths will be used to generate the provider resources.

### 5. Builds the Provider Resources

* Creates updated xpath input file (if needed)
  * If there are any found issues during the building of the provider resources, a new file called `updated_xpath_inputs.xml` will be created with a trimmed version of the provided xpath inputs file. This file can be used to replace the data in the `xpath_inputs.xml`.
* We need JTAF to create some `.go` code from the YANG models and XML data we provided which is written to the `/terraform_provider` directory.

### 6. Build the Provider

* Similar to the demo, the script run the creation of the binary file `terraform-provider-junos-[device-name]` which is actually our fresh new and shiny Terraform Provider. If you got this far, congratulations. You just created a Terraform provider for Junos.

### 7. Using the resouce

* For the next section, refer to the `/TFtemaplates` directory created by the script providing a basic `.tf` template for the `main` and `test` files


<a id="provider"></a>
# Using the new Provider

To test the provider we need to do two more things, one, put the provider where Terraform can find it and two, create a simple set of `.tf` files as inputs to Terraform!

__Terraform Search Locations__

More recent versions of Terraform will look for providers (plugins) in the Hashicorp hosted registry and in a number of local locations. 

Firstly, there is a default location, which is located under `~/.terraform.d/`. You need to create a hierarchy for each provider, it's version and the CPU architecture type the provider was compiled for. Here is what the author's looks like:

```bash
tree /Users/dgee/.terraform.d
├── checkpoint_cache
├── checkpoint_signature
└── plugins
    └── juniper
        └── providers
            └── junos-vsrx
                └── 19.41.101
                    ├── darwin_amd64
                    │   └── terraform-provider-junos-vsrx
                    └── linux_amd64
                        └── terraform-provider-junos-vsrx
```

Notice that the R is missing in the version `19.41.101`. This is an official Juniper method of naming providers. It basically says: 19.4R1.01 of Junos, with version 01 of the provider. Due to semantic versioning, you'll notice that the leading zero has been stripped, leaving us with the last two digits for the provider number and any other digits in front being the Junos version patch number.

If you're building providers locally, it's worth considering how to version control them. Each provider generated can have a different set of capabilities, even for the same software release, so it's important that you keep a track through a simple version control system. MD5 hashing of the binary is also recommended, so as a worst case, you can identify the binaries by their computed hash.

You can use the method above, enabling Terraform to find the provider locally. Here's how:

`mkdir -p ~/.terraform.d/plugins/juniper/providers/junos-vsrx/19.41.101/darwin_amd64`

It's probably a good idea to replace the `juniper` part with your own organisation's name to prevent any confusion.

**The second way** is to create a `.terraformrc` file in your home directory, where you tell Terraform where to look for the same file structure. That can be in a project directory if you so wish. Here is an example of that file.

```bash
provider_installation {
  filesystem_mirror {
    path = "/path_to_junos-terraform/junos-terraform/plugins"
    include = ["*/*/*"]
  }
}
```

Make sure that the same file tree exists from `plugins` as before.

The other option of course, is to publish to your provider/s to the Hashicorp registry and not have them stored locally.

<a id="testing"></a>
## __Testing With Terraform__

Ok, now we've got the Terraform provider in place, we can actually test Terraform! For this section, you will replace the `provider` section with the information for the device which is being configured. 

> If testing the provider from scatch, skip this message. If building a provider from a pre-loaded configuration, the following steps have been more or less done for you. Look for the `junos-terraform/TFtemplates` which will have prcompiled test files to use for testing. The `testbed` and required files and folders have also already been made for you. The only requirment is to manually fill in the resource information.

You are free to choose a directory in which to test this. I'm going to stick with the home `/junos-terrafrom` directory.

```bash
cd /junos-terraform
mkdir testbed && cd testbed
```

We need to create a number of files.

```bash
mkdir vsrx_1
touch main.tf
touch vsrx_1/main.tf
```

Content of `main.tf`
```bash
terraform {
  required_providers {
    junos-vsrx = {
      source = "juniper/providers/junos-vsrx"
      version = "19.41.101"
    }
  }
}

provider "junos-vsrx" {
    host = "localhost"
    port = 8300
    username = "root"
    password = "juniper123"
    sshkey = ""
}

module "vsrx_1" {
  source = "./vsrx_1"

  providers = {junos-vsrx = junos-vsrx}

  depends_on = [junos-vsrx_destroycommit.commit-main]
}

resource "junos-vsrx_commit" "commit-main" {
  resource_name = "commit"
  depends_on = [module.vsrx_1]
}

resource "junos-vsrx_destroycommit" "commit-main" {
  resource_name = "destroycommit"
}
```

Content of `vsrx_1/main.tf`

```bash
terraform {
  required_providers {
    junos-vsrx = {
      source = "juniper/providers/junos-vsrx"
      version = "19.41.101"
    }
  }
}

resource "junos-vsrx_InterfacesInterfaceDescription" "vsrx_1" {
    resource_name = "vsrx_1"
    name = "ge-0/0/0"
    description = "Test description"
}

resource "junos-vsrx_InterfacesInterfaceUnitFamilyInetAddressName" "vsrx_2" {
    resource_name = "vsrx_2"
    name = "ge-0/0/0"
    name__1 = "0"
    name__2 = "10.0.0.1/24"
}
```

__Let's Initialise Terraform__

We're getting so close! Let's initialize Terraform. From the `testbed` folder, run: 

```bash
testbed $ terrafrom init

Initializing modules...

Initializing the backend...

Initializing provider plugins...
- Finding juniper/providers/junos-vsrx versions matching "19.41.101"...
- Installing juniper/providers/junos-vsrx v19.41.101...
- Installed juniper/providers/junos-vsrx v19.41.101 (unauthenticated)

Terraform has created a lock file .terraform.lock.hcl to record the provider
selections it made above. Include this file in your version control repository
so that Terraform can guarantee to make the same selections by default when
you run "terraform init" in the future.

Terraform has been successfully initialized!

You may now begin working with Terraform. Try running "terraform plan" to see
any changes that are required for your infrastructure. All Terraform commands
should now work.

If you ever set or change modules or backend configuration for Terraform,
rerun this command to reinitialize your working directory. If you forget, other
commands will detect it and remind you to do so if necessary.
```

The next step is to actually run the plan and apply steps.

```bash
testbed $ terraform plan

Terraform used the selected providers to generate the following execution plan. Resource actions are indicated with the following symbols:
  + create

Terraform will perform the following actions:

  # junos-vsrx_commit.commit-main will be created
  + resource "junos-vsrx_commit" "commit-main" {
      + id            = (known after apply)
      + resource_name = "commit"
    }

  # junos-vsrx_destroycommit.commit-main will be created
  + resource "junos-vsrx_destroycommit" "commit-main" {
      + id            = (known after apply)
      + resource_name = "destroycommit"
    }

  # module.vsrx_1.junos-vsrx_InterfacesInterfaceDescription.vsrx_1 will be created
  + resource "junos-vsrx_InterfacesInterfaceDescription" "vsrx_1" {
      + description   = "Test description"
      + id            = (known after apply)
      + name          = "ge-0/0/0"
      + resource_name = "vsrx_1"
    }

  # module.vsrx_1.junos-vsrx_InterfacesInterfaceUnitFamilyInetAddressName.vsrx_2 will be created
  + resource "junos-vsrx_InterfacesInterfaceUnitFamilyInetAddressName" "vsrx_2" {
      + id            = (known after apply)
      + name          = "ge-0/0/0"
      + name__1       = "0"
      + name__2       = "10.0.0.1/24"
      + resource_name = "vsrx_2"
    }

Plan: 4 to add, 0 to change, 0 to destroy.

─────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────

Note: You didn't use the -out option to save this plan, so Terraform can't guarantee to take exactly these actions if you run "terraform apply" now.
```

Our plan is simple! Because we do not have any local Terraform state, the plan has been generated quickly and it's straight forward to read. The `commit` and `destroycommit` resources will be covered after this step. We can also tell Terraform to auto-approve the apply without any further manual input.

```bash
testbed $ terraform apply -auto-approve

Terraform used the selected providers to generate the following execution plan. Resource actions are indicated with the following symbols:
  + create

Terraform will perform the following actions:

  # junos-vsrx_commit.commit-main will be created
  + resource "junos-vsrx_commit" "commit-main" {
      + id            = (known after apply)
      + resource_name = "commit"
    }

  # junos-vsrx_destroycommit.commit-main will be created
  + resource "junos-vsrx_destroycommit" "commit-main" {
      + id            = (known after apply)
      + resource_name = "destroycommit"
    }

  # module.vsrx_1.junos-vsrx_InterfacesInterfaceDescription.vsrx_1 will be created
  + resource "junos-vsrx_InterfacesInterfaceDescription" "vsrx_1" {
      + description   = "Test description"
      + id            = (known after apply)
      + name          = "ge-0/0/0"
      + resource_name = "vsrx_1"
    }

  # module.vsrx_1.junos-vsrx_InterfacesInterfaceUnitFamilyInetAddressName.vsrx_2 will be created
  + resource "junos-vsrx_InterfacesInterfaceUnitFamilyInetAddressName" "vsrx_2" {
      + id            = (known after apply)
      + name          = "ge-0/0/0"
      + name__1       = "0"
      + name__2       = "10.0.0.1/24"
      + resource_name = "vsrx_2"
    }

Plan: 4 to add, 0 to change, 0 to destroy.
junos-vsrx_destroycommit.commit-main: Creating...
junos-vsrx_destroycommit.commit-main: Creation complete after 0s [id=localhost_destroycommit]
module.vsrx_1.junos-vsrx_InterfacesInterfaceDescription.vsrx_1: Creating...
module.vsrx_1.junos-vsrx_InterfacesInterfaceUnitFamilyInetAddressName.vsrx_2: Creating...
module.vsrx_1.junos-vsrx_InterfacesInterfaceUnitFamilyInetAddressName.vsrx_2: Creation complete after 5s [id=localhost_vsrx_2]
module.vsrx_1.junos-vsrx_InterfacesInterfaceDescription.vsrx_1: Still creating... [10s elapsed]
module.vsrx_1.junos-vsrx_InterfacesInterfaceDescription.vsrx_1: Creation complete after 10s [id=localhost_vsrx_1]
junos-vsrx_commit.commit-main: Creating...
junos-vsrx_commit.commit-main: Creation complete after 6s [id=localhost_commit]

Apply complete! Resources: 4 added, 0 changed, 0 destroyed.
```

Let's just check the input on the vSRX instance to make sure we have the correct configs!

Terraform deals with configuration as a set of remote resources, which are stored in Junos `configuration groups`. These groups are then inherited by Junos at commit time. Here is the inherited configuration.

```bash
# show interfaces ge-0/0/0 | display inheritance no-comments
description "Test description";
unit 0 {
    family inet {
        address 10.0.0.1/24;
    }
}
```

Here are the groups:

```bash
# show groups
vsrx_2 {
    interfaces {
        ge-0/0/0 {
            unit 0 {
                family inet {
                    address 10.0.0.1/24;
                }
            }
        }
    }
}
vsrx_1 {
    interfaces {
        ge-0/0/0 {
            description "Test description";
        }
    }
}
```

## Terraform & Junos Commits

Junos and Terraform are not natural buddies. Their life-cycles are somewhat orthogonal in nature.

We handle the commit based transactional nature of Junos by using a simple commit pattern, in which Terraform creates phony commit resources. The resources are only tracked locally as commits are nothing more than Junos remote procedure calls. 

Thankfully, Terraform has a way of creating logical groupings, called `modules`. Terraform modules represent re-usable chunks of logic, in which we can give a name. Because we can name these groups, it means we can also place dependencies up on them!

You might have noticed in the top level `main.tf` file, there was some `depends_on` Terraform HCL keys.

It's by using these concrete dependencies, we are able to trigger the resources `commit` and `destroycommit` to be created.

The dependency order is thus:

1. The `commit` resource depends on the module
2. The module contains the actual desired state (which may have further ordered structure)
3. The module depends on the `destroycommit` resource


This ordering means the `destroycommit` is created first. No action is taken when this resource is created other than local state is stored on the system. The contents of the module are executed next, which consists of NETCONF sessions being made against the target system and stored in configuration groups, which are applied. Lastly, the commit resource is created, which actually runs a commit on Junos via a NETCONF RPC.

__Making Changes__

When you make either a change on Junos, or a change to the local resorce `.tf` files, then run a Terraform plan, Terraform doesn't understand a commit must be run. Therefore you must a `terraform taint $terraform_address.commit` command against the commit resource, which tells Terraform to re-run the commit after making other changes in the module.

Let's try it out in context of this demonstration. In the `vsrx_1/main.tf`, change the inet address to `.2` instead of `.1`.
Run the `terraform taint` command and re-run the plan and apply sequences.

```bash
terraform taint junos-vsrx_commit.commit-main

junos-vsrx_destroycommit.commit-main: Refreshing state... [id=localhost_destroycommit]
module.vsrx_1.junos-vsrx_InterfacesInterfaceUnitFamilyInetAddressName.vsrx_2: Refreshing state... [id=localhost_vsrx_2]
module.vsrx_1.junos-vsrx_InterfacesInterfaceDescription.vsrx_1: Refreshing state... [id=localhost_vsrx_1]
junos-vsrx_commit.commit-main: Refreshing state... [id=localhost_commit]

Terraform used the selected providers to generate the following execution plan. Resource actions are indicated with the following symbols:
  ~ update in-place
-/+ destroy and then create replacement

Terraform will perform the following actions:

  # junos-vsrx_commit.commit-main is tainted, so must be replaced
-/+ resource "junos-vsrx_commit" "commit-main" {
      ~ id            = "localhost_commit" -> (known after apply)
        # (1 unchanged attribute hidden)
    }

  # module.vsrx_1.junos-vsrx_InterfacesInterfaceUnitFamilyInetAddressName.vsrx_2 will be updated in-place
  ~ resource "junos-vsrx_InterfacesInterfaceUnitFamilyInetAddressName" "vsrx_2" {
        id            = "localhost_vsrx_2"
        name          = "ge-0/0/0"
      ~ name__2       = "10.0.0.1/24" -> "10.0.0.2/24"
        # (2 unchanged attributes hidden)
    }

Plan: 1 to add, 1 to change, 1 to destroy.

─────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────

Note: You didn't use the -out option to save this plan, so Terraform can't guarantee to take exactly these actions if you run "terraform apply" now.

# Now run the apply
 terraform apply -auto-approve
junos-vsrx_destroycommit.commit-main: Refreshing state... [id=localhost_destroycommit]
module.vsrx_1.junos-vsrx_InterfacesInterfaceDescription.vsrx_1: Refreshing state... [id=localhost_vsrx_1]
module.vsrx_1.junos-vsrx_InterfacesInterfaceUnitFamilyInetAddressName.vsrx_2: Refreshing state... [id=localhost_vsrx_2]
junos-vsrx_commit.commit-main: Refreshing state... [id=localhost_commit]

Terraform used the selected providers to generate the following execution plan. Resource actions are indicated with the following symbols:
  ~ update in-place
-/+ destroy and then create replacement

Terraform will perform the following actions:

  # junos-vsrx_commit.commit-main is tainted, so must be replaced
-/+ resource "junos-vsrx_commit" "commit-main" {
      ~ id            = "localhost_commit" -> (known after apply)
        # (1 unchanged attribute hidden)
    }

  # module.vsrx_1.junos-vsrx_InterfacesInterfaceUnitFamilyInetAddressName.vsrx_2 will be updated in-place
  ~ resource "junos-vsrx_InterfacesInterfaceUnitFamilyInetAddressName" "vsrx_2" {
        id            = "localhost_vsrx_2"
        name          = "ge-0/0/0"
      ~ name__2       = "10.0.0.1/24" -> "10.0.0.2/24"
        # (2 unchanged attributes hidden)
    }

Plan: 1 to add, 1 to change, 1 to destroy.
junos-vsrx_commit.commit-main: Destroying... [id=localhost_commit]
junos-vsrx_commit.commit-main: Destruction complete after 0s
module.vsrx_1.junos-vsrx_InterfacesInterfaceUnitFamilyInetAddressName.vsrx_2: Modifying... [id=localhost_vsrx_2]
module.vsrx_1.junos-vsrx_InterfacesInterfaceUnitFamilyInetAddressName.vsrx_2: Modifications complete after 5s [id=localhost_vsrx_2]
junos-vsrx_commit.commit-main: Creating...
junos-vsrx_commit.commit-main: Creation complete after 6s [id=localhost_commit]

Apply complete! Resources: 1 added, 1 changed, 1 destroyed.

# Re-run the plan to confirm!
terraform plan
junos-vsrx_destroycommit.commit-main: Refreshing state... [id=localhost_destroycommit]
module.vsrx_1.junos-vsrx_InterfacesInterfaceUnitFamilyInetAddressName.vsrx_2: Refreshing state... [id=localhost_vsrx_2]
module.vsrx_1.junos-vsrx_InterfacesInterfaceDescription.vsrx_1: Refreshing state... [id=localhost_vsrx_1]
junos-vsrx_commit.commit-main: Refreshing state... [id=localhost_commit]

No changes. Your infrastructure matches the configuration.

Terraform has compared your real infrastructure against your configuration and found no differences, so no changes are needed.
```

You can also check the Junos config to make sure the reflected change has happened and that the commit has been executed.

```bash
show interfaces ge-0/0/0 | display inheritance no-comments
description "Test description";
unit 0 {
    family inet {
        address 10.0.0.2/24;
    }
}
```

There we have it.

Now for the fun part! Let's clean up and destroy the Terraform state.

```bash
terraform destroy -auto-approve
junos-vsrx_destroycommit.commit-main: Refreshing state... [id=localhost_destroycommit]
module.vsrx_1.junos-vsrx_InterfacesInterfaceDescription.vsrx_1: Refreshing state... [id=localhost_vsrx_1]
module.vsrx_1.junos-vsrx_InterfacesInterfaceUnitFamilyInetAddressName.vsrx_2: Refreshing state... [id=localhost_vsrx_2]
junos-vsrx_commit.commit-main: Refreshing state... [id=localhost_commit]

Terraform used the selected providers to generate the following execution plan. Resource actions are indicated with the following symbols:
  - destroy

Terraform will perform the following actions:

  # junos-vsrx_commit.commit-main will be destroyed
  - resource "junos-vsrx_commit" "commit-main" {
      - id            = "localhost_commit" -> null
      - resource_name = "commit" -> null
    }

  # junos-vsrx_destroycommit.commit-main will be destroyed
  - resource "junos-vsrx_destroycommit" "commit-main" {
      - id            = "localhost_destroycommit" -> null
      - resource_name = "destroycommit" -> null
    }

  # module.vsrx_1.junos-vsrx_InterfacesInterfaceDescription.vsrx_1 will be destroyed
  - resource "junos-vsrx_InterfacesInterfaceDescription" "vsrx_1" {
      - description   = "Test description" -> null
      - id            = "localhost_vsrx_1" -> null
      - name          = "ge-0/0/0" -> null
      - resource_name = "vsrx_1" -> null
    }

  # module.vsrx_1.junos-vsrx_InterfacesInterfaceUnitFamilyInetAddressName.vsrx_2 will be destroyed
  - resource "junos-vsrx_InterfacesInterfaceUnitFamilyInetAddressName" "vsrx_2" {
      - id            = "localhost_vsrx_2" -> null
      - name          = "ge-0/0/0" -> null
      - name__1       = "0" -> null
      - name__2       = "10.0.0.2/24" -> null
      - resource_name = "vsrx_2" -> null
    }

Plan: 0 to add, 0 to change, 4 to destroy.
junos-vsrx_commit.commit-main: Destroying... [id=localhost_commit]
junos-vsrx_commit.commit-main: Destruction complete after 0s
module.vsrx_1.junos-vsrx_InterfacesInterfaceUnitFamilyInetAddressName.vsrx_2: Destroying... [id=localhost_vsrx_2]
module.vsrx_1.junos-vsrx_InterfacesInterfaceDescription.vsrx_1: Destroying... [id=localhost_vsrx_1]
module.vsrx_1.junos-vsrx_InterfacesInterfaceDescription.vsrx_1: Destruction complete after 3s
module.vsrx_1.junos-vsrx_InterfacesInterfaceUnitFamilyInetAddressName.vsrx_2: Destruction complete after 5s
junos-vsrx_destroycommit.commit-main: Destroying... [id=localhost_destroycommit]
junos-vsrx_destroycommit.commit-main: Destruction complete after 4s

Destroy complete! Resources: 4 destroyed.

# You can check the vSRX manually to make sure the config has gone

show interfaces ge-0/0/0 | display inheritance no-comments
```

There is some inverse logic here. The `destroycommit` is 'created' on the delete cycle within the provider, whereas the `commit` is 'created' on the create cycle. Terraform providers do nothing more than CRUD (create/read/update/delete) on data structures, which in turn have methods on them. The Terraform apply cycle runs `create` and `update`. Terraform destroy in turns runs the `delete` method.

You can read more on this commit 'bookend' pattern [here](https://dave.dev/blog/2021/11/).


<a id="questions"></a>
# Q&A

__1. When I downloaded, terraform, the commands are not recognized by my device__

Ensure the that the terraform download is found in the `/usr/local/bin` directory

* Go to where terraform download is located
 > `cd /Downloads`

* Use the sudo `command` to move it to the correct location.
> `sudo mv terraform /usr/local/bin/`

__2. Where do I find the names of the resources I have created?__

Check in the `provider.go` file that is dynamically generated. You will find a data structure with this signature: `ResourcesMap: map[string]*schema.Resource`. Your resources are named in a map. Here is an example:

```bash
// Output of a provider.go example
func Provider() *schema.Provider {
	return &schema.Provider{

		Schema: map[string]*schema.Schema{
			"host": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},

			"port": &schema.Schema{
				Type:     schema.TypeInt,
				Required: true,
			},

			"username": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},

			"password": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"sshkey": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
		},

        // These are your resources
		ResourcesMap: map[string]*schema.Resource{
			"junos-vsrx_InterfacesInterfaceDescription": junosInterfacesInterfaceDescription(),
			"junos-vsrx_InterfacesInterfaceUnitFamilyInetAddressName": junosInterfacesInterfaceUnitFamilyInetAddressName(),
			"junos-vsrx_commit": junosCommit(),
	        "junos-vsrx_destroycommit": junosDestroyCommit(),
			},
		ConfigureFunc: providerConfigure,
	}
}
```

__3. Running `.generte.sh` script sometimes fails due to go runtime error or pyang error.__
* Ensure that pyang and go are installed. If so, re-run the script and it should continue.
* If errors persist, try examining to ensure the right version and yang file types are being used. Deleting the autogenerated folders may be required. 

__4. Running `terraform init` is giving me errors regarding the location of the provider.__

* Some systems have issues recognizing the `.terraformrc` file so if using this method, delete this file and follow below.
* Ensure that you place the provider in the `~/.terraform.d/plugins/juniper/providers/junos-vsrx/19.41.101/darwin_amd64/` folder. 
  * If using `darwin_arm64`, make sure to rename the folder to match the device's core. 
* Once this is double-checked, check the `main.tf` to ensure that the terrafrom `required_providers` matches the naming of the path and try again. 


__5. Running `terraform apply` is giving me errors realted to Plugin not Responding. How do I fix this?__
> This only applies when NOT using the `$MOCK_FILE` testing env variable 
* If terrafrom cannot connect to a Juniper Device during the `apply` stage, the message `Error: Plugin did not respond` will occur.
  * To fix this issue, ensure that the `provider "junos-deviceName"` section in the `main.tf` file is correctly filled out and matches a running Junos device which can recieves data. If the host and port does not connect to a running device, the `terraform apply` will not work.
* If using the `$MOCK_FILE` env variable --> the information in  the `provider "junos-deviceName"` section in the `main.tf` is not relevant to the output in the log file defined by the varible.

## Close

This covers the use of JTAF with a full example. Please post any issues or bug reports using GitHub issues at the top of this page. Enjoy!


# CONTRIBUTORS
Juniper Networks is actively contributing to and maintaining this repo.
 
*Contributors:*
* [Rahul Kumar](https://github.com/rahkumar651991)
* [David Gee](https://github.com/davedotdev)
