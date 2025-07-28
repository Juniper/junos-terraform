# JUNOS Terraform Automation Framework (JTAF)

Terraform is traditionally used for managing virtual infrastructure, but there are organisations out there that use Terraform end-to-end and also want to manage configuration state using the same methods for managing infrastructure. Sure, we can run a provisioner with Terraform, but that wasn't asked for!

Much the same as you can use Terraform to create an AWS EC2 instance, you can manage the configurational state of Junos. In essence, we treat Junos configuration as declarative resources.

So what is JTAF? It's a framework, meaning, it's an opinionated set of tools and steps that allow you to go from YANG models to a custom Junos Terraform provider. With all frameworks, there are some dependencies.

To use JTAF, you'll need machine that can run **Go, Python, Git and Terraform.** This can be Linux, OSX or Windows. Some easy to consume videos are below.

## Quick start

### <u>Setup</u>
Run the following commands to set up the Junos-Terraform Environment and Workflow

```bash
git clone https://github.com/juniper/junos-terraform
git clone https://github.com/juniper/yang
python3 -m venv venv
. venv/bin/activate
pip install ./junos-terraform
cd junos-terraform
```
---
### <u>Yang File(s) to JSON Conversion</u>

Find the device's Junos Version that is running, and locate the corresponding yang and common folders. Run the below `pyang` command to generate a `.json` file containing `.yang` information for that version. [See below example for Junos version 18.2]
```
pyang --plugindir $(jtaf-pyang-plugindir) -f jtaf -p <path-to-common> <path-to-yang-files> > junos.json
```
Example: 
```
pyang --plugindir $(jtaf-pyang-plugindir) -f jtaf -p ../yang/18.2/18.2R3/common ../yang/18.2/18.2R3/junos-qfx/conf/*.yang > junos.json
```
---

### <u>Generate Resource Provider</u>

Now run the following command to generate a `resource provider`. 

```bash
jtaf-provider -j <json-file> -x <xml-configuration> -t <device-type>
```

Example:
```bash
jtaf-provider -j junos.json -x examples/evpn-vxlan-dc/dc2/dc2-spine1.xml -t vqfx
```
All in one example (`-j` accepts `-` for `stdin` for `jtaf-provider`):
```bash
pyang --plugindir $(jtaf-pyang-plugindir) -f jtaf -p ../yang/18.2/18.2R3/common ../yang/18.2/18.2R3/junos-qfx/conf/*.yang | jtaf-provider -j - -x examples/evpn-vxlan-dc/dc2/dc2-spine1.xml -t vqfx
```

---

### <u>Single command to generate resource provider</u>

Use `jtaf-yang2go` command to generate a resource provider in a single step by supplying all YANG files with the `-p` option, the device XML configuration with `-x`, and the device type with `-t`.

```bash
jtaf-yang2go -p <path-to-common> <path-to-yang-files> -x <xml-configuration> -t <device-type>
```

Example:

```bash
jtaf-yang2go -p ../yang/18.2/18.2R3/common ../yang/18.2/18.2R3/junos-qfx/conf/*.yang -x examples/evpn-vxlan-dc/dc2/dc2-spine1.xml -t vqfx
```

---

### <u>Build the provider and install</u>

cd into the newly created directory starting with `terraform-provider-junos-` then the device-type and then `go install`

Example:

```
cd terraform-provider-junos-vqfx
go install
```

---

### <u>Autogenerate Terraform Testing Files</u>

Run this command to generate a `.tf` test file to deploy the Terraform provider.
Output will be returned in the terminal.
```
jtaf-xml2tf -x <device-xml-config> -t <device-type> -n <deivce-host-name>
```

Example: 
```
jtaf-xml2tf -x config.xml -t vqfx -n dc2-spine1
```

Using the output from the terminal, which represents a template for the HCL .tf file, we can create our testing environment and fill in the template with the necessary device information.

### <u>Setting up Testing environment</u>

Create a testing folder which can be used to write .tf files and apply terraform configuration.   
Example
	```
	mkdir testbed
	```

In the `/testbed` folder created:  
* Create a `main.tf` file with the content of terminal output from the `jtaf-xml2tf` command.  
	* Fill in any missing information

Next, create a `.terraformrc` file in your home directory  with `vi` and add the following contents, replacing any `<elements>` tags with your own information. This is to ensure that the terraform plugin you created and installed to `/go/bin` will be read.

**.terraformrc example**
```
provider_installation {
	dev_overrides {
		"registry.terraform.io/hashicorp/junos-<device-type>" = "<path-to-go/bin>"
	}
}
```

Example:
```
provider_installation {
	dev_overrides {
		"registry.terraform.io/hashicorp/junos-vqfx" = "/Users/patelv/go/bin"
	}
}
```

You should know have a file structure which looks similar to:

```
/junos-terraform/<testing-folder-name>/
/junos-terraform/<testing-folder-name>/main.tf         <-- contents of jtaf-xml2tf command

/Users/<username>/.terraformrc     <-- link to provider created in /usr/go/bin/ [see details above]
```

---

### <u>Edit Test Files, Plan, and Apply</u>

Once the `.terraform.rc` file is set up, and the `main.tf` test file contains access to the provider, information regarding the desired devices to push the configuration to, and the desired config in `HCL` format, we are now ready to use the provider.

```
terrafrom plan
terraform apply -auto-approve
```
