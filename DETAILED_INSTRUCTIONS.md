# JUNOS Terraform Automation Framework (JTAF) Detailed Build Instructions

JTAF is a framework, meaning, it's an opinionated set of tools and steps that allow you to go from YANG models to a Junos Terraform provider. With all frameworks, there are some dependencies.

The tool you'll need to use JTAF is a bash shell. If you're on OSX or a Linux user, you're set for success out of the box. On Windows you can use WSL to install a Linux flavour like Ubuntu and use its bash shell.

In this document, if you see `$JTAF_PROJECT` you can replace it with the path of the JTAF project on your system or create an environment variable. On the author's sytem, this happens to be below:

```bash
export JTAF_PROJECT=/Users/dgee/Documents/GoDev/src/github.com/Juniper/junos-terraform

# Let's check
echo $JTAF_PROJECT
/Users/dgee/Documents/GoDev/src/github.com/Juniper/junos-terraform
```

### 1. Install Python, Go & Terraform

Go version tested with: `go1.14.2 darwin/amd64`
Python version tested with: `3.7`
Terraform version: `v0.12.26`

Other versions beyond these will work, but this is what was tested for the writing of this document.

### 2.  Create and Activate a Python `venv`

A virtual environment does nothing more than provide a separate Python environment, that's safe to install whatever you need, without affecting the global Python install on a given system. 

*For this step, ensure that you're in the `processYang` directory of JTAF.*

```bash
python -m venv ./venv
source venv/bin/activate
```

The prompt will change at this point indicating you are running with a Python virtual environment.

### 3.	Install Pyang and Check the Version

```bash
pip install pyang
pyang -v
```

On the system under test, the version installed by `pip` in to the `venv` was: `2.5.0`.

### 4.	Copy the YANG Files 
Let's put the YANG files from [Juniper's YANG GitHub repository](https://github.com/Juniper/yang.git) in to a memorable location. Let's use `/var/tmp/yang`.

```bash
cd /var/tmp/
git clone https://github.com/Juniper/yang.git
```

*Note - this may take some time*

Next, to match one of the existing samples, let's copy a specific YANG model into working directory.

*I've repeated the directory change in the bash steps below (just in case!)*

```bash
cd /var/tmp
mkdir jtafwd && cd jtafwd
mkdir terraform_providers
cd ../
mv yang/19.4/19.4R1/junos-qfx/conf/junos-qfx-conf-protocols@2019-01-01.yang ./jtafwd
```

If you wanted to remove the YANG directory, you can do it like this:

```bash
cd /var/tmp
rm -rf yang
```

### 5.  Create a `config.toml` File

*If you've never seen a TOML file before, don't worry! It's just a structured file containing configuration that can be parsed by a program, in this case the two main compiled programs that form JTAF. TOML stands for Tom's Obvious Minimal Langage.*

Create a config file somewhere memorable. I'll use `/var/tmp/jtafwd/config.toml` because why not.

Using your favourite text editor, create a file here: `/var/tmp/jtafwd/config.toml` and put the content below into the file. Don't worry about the xPath or fileType keys. They'll be explained shortly.
The providerName in the toml file will be used to represent the name of the provider in the automagically generated provider files within the "providerDir"

```bash
yangDir = "/var/tmp/jtafwd"
providerDir = "/var/tmp/jtafwd/terraform_providers"
xpathPath = "/var/tmp/jtafwd/xpath_test.xml"
fileType = "both"
providerName = "vsrx"
```

You can also replace the fileType field to "text" or "xml". The text files are for us humans.

### 6.	Generate the YIN and XPath Files

The next step, depending on the size of YANG model/s, may take some time. Prepare some popcorn! Ensure you have copied all the dependent YANG modules into the specified yangDir, else you may run into errors when you run the processYang program. 

```bash
cd cmd/processYang
go build
./processYang -config /var/tmp/jtafwd/config.toml
# OUTPUT - WARNING >> This can take some time. Lack of activity does not mean broken!
Yin file for junos-qfx-conf-protocols@2019-01-01 is generated
Creating Xpath file: junos-qfx-conf-protocols@2019-01-01_xpath.txt
Creating Xpath file: junos-qfx-conf-protocols@2019-01-01_xpath.xml
```

Note: Copy the below YANG files from the common directory as the other junos YANG modules such as junos, junos-es, junos-nfx etc are dependent on them. Any dependent YANG files used should be from the same version. Mixing versions may cause issues.
```
-rw-r--r-- 1 root root 3398 Dec 4 06:26 junos-common-types@2019-01-01.yang
-rw-r--r-- 1 root root 2346 Dec 4 06:26 junos-common-odl-extensions@2019-01-01.yang
-rw-r--r-- 1 root root 1806 Dec 4 06:26 junos-common-ddl-extensions@2019-01-01.yang
```

### 7.	Create an XML XPath File
This file acts as an input to JTAF. This input identifies the content of the provider that JTAF will create. Some `xpath_test.xml` files are scattered are in the `Samples` directory. 

Create a file `/var/tmp/jtafwd/xpath_test.xml` and populate it with the content below.

```bash
<file-list>
    <xpath name="/protocols/bgp/group/traceoptions/file/filename">
    </xpath>
</file-list>
```

The content above does this:

*   The path `/protocols/bgp/group/traceoptions/file/filename` tells JTAF to use the YANG element filename.


### 8.	Build the Provider

First, we need JTAF to create some `.go` code from the YANG models and XML data we provided.

```bash
cd cmd/processProviders
go build
./processProviders -config /var/tmp/jtafwd/config.toml
```

The output of this step is written to the `/var/tmp/jtafwd/terraform_provider` directory. You'll see a `.go` source file.
notice that the generated `.go` files are shown as "resource_<name>", however the provider name is appended to the resources based on declared name from the config.toml file.

Example of generated providers
```
-rw-r--r-- 1 root root  3814 Dec  6 06:22 resource_ApplicationsApplicationDestination__Port.go
-rw-r--r-- 1 root root  3555 Dec  6 06:22 resource_ApplicationsApplicationProtocol.go
-rw-r--r-- 1 root root  3669 Dec  6 06:22 resource_ApplicationsApplicationSource__Port.go
-rw-r--r-- 1 root root  4096 Dec  6 06:22 resource_FirewallFilterTermFromProtocol.go
```

View the provider.go file and notice the declared provider name is appended to the resources generated. The .tf file must call the resources based on the resource map definition within provider.go file.

```
"junos-vsrx_ApplicationsApplicationProtocol": junosApplicationsApplicationProtocol(),
"junos-vsrx_ApplicationsApplicationSource__Port": junosApplicationsApplicationSource__Port(),
"junos-vsrx_ApplicationsApplicationDestination__Port": junosApplicationsApplicationDestination__Port(),
"junos-vsrx_InterfacesInterfaceDescription": junosInterfacesInterfaceDescription(),
...
```


Next, copy this `.go` file to the `terraform_providers` directory within the JTAF project.

```bash
cd $JTAF_PROJECT/terraform_providers
cp /var/tmp/jtafwd/terraform_providers/*.go ./
```

The last step is to actually build the provider!

```bash
go mod init terraform-provider-junos-device
go mod tidy
go build -o terraform-provider-junos-device
```

This provider without any Go cross-compilation directives, will work on the system it's been generated with. If you happen to be on an OSX machine, then the provider will work for Terraform on OSX and the same is true for Linux, if you use JTAF on Linux, then natively the generated provider will operate on Linux. However, you can cross-compile the provider so that it will operate on another operating system and even CPU architecture.

```bash
$ file terraform-provider-junos-device

terraform-provider-junos-device: Mach-O 64-bit executable x86_64
```

If you did want this provider to work with Linux, then you can cross-compile using the `GOOS` input. See below.

```bash
$ GOOS=linux go build -o terraform-provider-junos-device

$ file terraform-provider-junos-device
terraform-provider-junos-device: ELF 64-bit LSB executable, x86-64, version 1 (SYSV), statically linked, Go BuildID=SWvAslM7UiUlMNJJOG8f/MV8jDWinx0vKkuo7Zmec/-2fk9ZDz88J7folCoc0q/ftWLT5N4tiPWQ8DlXY2J, not stripped
```

The file `terraform-provider-junos-device` is actually our fresh new and shiny Terraform Provider. If you got this far, congratulations. You just created a Terraform provider for Junos.

### 9.  Deactivate `venv`

```bash
deactivate
```

# Using the new Provider

You are free to choose a directory in which to test this. I'm going to stick with the `/var/tmp/jtafwd` root directory.

```bash
cd 
mkdir /var/tmp/jtafwd/testtf
mv $JTAF_PROJECT/terraform_providers/terraform-provider-junos-device /var/tmp/jtafwd/testtf
cd /var/tmp/jtafwd/testtf
```

Create a `test.tf` file with the contents below in the current working directory (which we set above). There is enough information in this Terraform config file to connect to the Junos device and have Terraform process the resources and HCL statements.

```bash
provider "junos-device" {
    host = "10.x.x.x"
    port = 22
    username = "user"
    password = "user123"
    sshkey = ""
}

resource "junos-device_commit" "commit" {
    resource_name = "commit"
    depends_on = [
        junos-device_ProtocolsBgpGroupTraceoptionsFileFilename.demo
    ]
}

resource "junos-device_ProtocolsBgpGroupTraceoptionsFileFilename" "demo" {
    resource_name = "bgp_trace_file"
    name = "demo"
    filename = "temp.log"
}
```

In this directory, we have the Terraform provider and a Terraform config file which will set a trace options file name for the BGP protocol.

__Let's Initialise Terraform__

```bash
$ terraform init

Initializing the backend...

Initializing provider plugins...

Terraform has been successfully initialized!

You may now begin working with Terraform. Try running "terraform plan" to see
any changes that are required for your infrastructure. All Terraform commands
should now work.

If you ever set or change modules or backend configuration for Terraform,
rerun this command to reinitialize your working directory. If you forget, other
commands will detect it and remind you to do so if necessary.
```

Great! Terraform initialised and now you can run other commands like  `terraform plan` and `terraform apply` providing your device is accessible via NETCONF and the credentials work.







