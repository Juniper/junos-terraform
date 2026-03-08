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

If you do not already have Terraform installed (in general), for macOS, run the following:
```bash
brew tap hashicorp/tap
brew install hashicorp/tap/terraform
```

For more information, refer to the Terraform website: https://developer.hashicorp.com/terraform/install.

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

NOTE: For Junos version >23.2 (i.e. starting from 23.4 onwards), the file path in the `yang` directory is slightly different as shown in the example below.
```
pyang --plugindir $(jtaf-pyang-plugindir) -f jtaf -p ../yang/23.4/23.4R1/native/conf-and-rpcs/common/models ../yang/23.4/23.4R1/native/conf-and-rpcs/junos/conf/models/*.yang > junos.json
```
 
---

### <u>Generate Resource Provider</u>

Now run the following command to generate a `resource provider`. 

```bash
jtaf-provider -j <json-file> -x <xml-configuration(s)> -t <device-type>
```

Example:
```bash
jtaf-provider -j junos.json -x examples/evpn-vxlan-dc/dc1/*{spine,leaf}*.xml examples/evpn-vxlan-dc/dc2/*spine*.xml -t vqfx
```
NOTE: If using multiple xml configurations (like the example above), ensure that the configurations are for the same device type

All in one example (`-j` accepts `-` for `stdin` for `jtaf-provider`):
```bash
pyang --plugindir $(jtaf-pyang-plugindir) -f jtaf -p ../yang/18.2/18.2R3/common ../yang/18.2/18.2R3/junos-qfx/conf/*.yang | jtaf-provider -j - -x examples/evpn-vxlan-dc/dc1/*{spine,leaf}*.xml examples/evpn-vxlan-dc/dc2/*spine*.xml  -t vqfx
```

---

### <u>Single command to generate resource provider</u>

Use `jtaf-yang2go` command to generate a resource provider in a single step by supplying all YANG files with the `-p` option, the device XML configuration with `-x`, and the device type with `-t`.

```bash
jtaf-yang2go -p <path-to-common> <path-to-yang-files> -x <xml-configuration(s)> -t <device-type>
```

Example:

```bash
jtaf-yang2go -p ../yang/18.2/18.2R3/common ../yang/18.2/18.2R3/junos-qfx/conf/*.yang -x examples/evpn-vxlan-dc/dc1/*{spine,leaf}*.xml examples/evpn-vxlan-dc/dc2/*spine*.xml -t vqfx
```
NOTE: If using multiple xml configurations (like the example above), ensure that the configurations are for the same device type

NOTE: For Junos version >23.2, the file path for the folder containing the yang files for each device is slightly different. Refer to section [Yang File(s) to JSON Conversion](./README.md#yang-files-to-json-conversion) for more information and examples.

---

### <u>Build the provider and install</u>

cd into the newly created directory starting with `terraform-provider-junos-` then the device-type and then `go install`

Example:

```
cd terraform-provider-junos-vqfx
go install
```


## <u>Autogenerate Terraform Testing Files</u>
---

### <u>Overview</u>

Run a command to generate a `.tf` test file to deploy the Terraform provider.

**NOTE:** Output will be returned to the terminal **OR** created in a directory depending on your passed flags.

**<u>Flag Options:</u>**
 * -j 
	* **Required:** trimmed_json output file from jtaf-provider (stored in terraform provider folder /terraform-provider-junos-"device-type")
 * -x
	* **Required:** File(s) of xml config to create terraform files for
 * -t
	* **Required:** Junos device type
 * -d
	* **Optional:** Flag to create multiple Terraform files under specified directory name, one for each xml config
 * -u
	* **Optional:** Device username
 * -p
	* **Optional:** Device password

---

### <u>Creating a single Terraform Testing File</u>

To create a single Terraform (.tf) file from a config file(s) use the following command (output returned to terminal):
```
jtaf-xml2tf -j <path-to-trimmed-schema> -x <path-to-config-files(s)> -t <device-type>
```

Example: 

* **trimmed_schema** - stored in terraform provider folder created from running the jtaf-provider module command (usually in terraform-provider-junos-'device-type')
* **xml_files** - directory containing xml file(s) (ensure xml file(s) are for the same device type)

```
jtaf-xml2tf -j terraform-provider-junos-vqfx/trimmed_schema.json -x examples/evpn-vxlan-dc/dc1/*{spine,leaf}*.xml examples/evpn-vxlan-dc/dc2/*spine*.xml -t vqfx
```
* If the user wants to provide the device **username** and **password**, those additional flags can be added as well
```
jtaf-xml2tf -j terraform-provider-junos-vqfx/trimmed_schema.json -x examples/evpn-vxlan-dc/dc1/*{spine,leaf}*.xml examples/evpn-vxlan-dc/dc2/*spine*.xml -t vqfx -u root -p password
```

Using the output from the terminal, which represents a template for the HCL .tf file, we can create our testing folder, copy the output into a terraform file, and fill in the template with the necessary device information.

#### <u>Create testing folder</u>

Create a testing folder which can be used to write .tf files and apply terraform configuration.   

Example
	```
	mkdir testbed
	```

In the `/testbed` folder created:  
* Create a `main.tf` file with the content of terminal output from the `jtaf-xml2tf` command.  
	* Fill in any missing information

Jump to [Setting up the Test Environment](./README.md#setting-up-testing-environment)

---

### <u>Creating multiple Terraform Testing Files</u>

To create multiple Terraform (.tf) files from multiple config files, where each .tf file will represent one xml file, use the following command (output returned to specified directory name):

```
jtaf-xml2tf -j <path-to-trimmed-schema> -x <path-to-config-files(s)> -t <device-type> -d <testing-folder-name>
```

Example: 

* **trimmed_schema** - stored in terraform provider folder created from running the jtaf-provider module command (usually in terraform-provider-junos-'device-type')
* **xml_files** - directory containing xml file(s) (ensure xml file(s) are for the same device type)

```
jtaf-xml2tf -j terraform-provider-junos-vqfx/trimmed_schema.json -x examples/evpn-vxlan-dc/dc1/*{spine,leaf}*.xml examples/evpn-vxlan-dc/dc2/*spine*.xml -t vqfx -d testbed
```
* If the user wants to provide the device(s) **username** and **password**, those additional flags can be added as well
```
jtaf-xml2tf -j terraform-provider-junos-vqfx/trimmed_schema.json -x examples/evpn-vxlan-dc/dc1/*{spine,leaf}*.xml examples/evpn-vxlan-dc/dc2/*spine*.xml -t vqfx -d testbed -u root -p password
```

Using the output which is outputted to the specifed directory from the command, which represents a template for the HCL .tf file for each input XML file, we can now create our testing environment and fill in the template with any remaining necessary device or config information.

---

### <u>Setting up Testing environment</u>

Now that we ran the `jtaf-xml2tf` command and have our testing folder setup:
* Note: if you created a single terraform file, you should have copied that output to a `.tf` file in a test folder in the `/junos-terrafom` directory:
	* ex: `junos-terraform/testbed/main.tf` <-- stores output from command

#### Creating the Enviornment

Next, create a `.terraformrc` file in your home directory, `(cd ~)`, with `vi` and add the following contents, replacing any `<elements>` tags with your own information. This is to ensure that the terraform plugin you created and installed to `/go/bin` will be read.

**.terraformrc example**
```
provider_installation {
	dev_overrides {
		"registry.terraform.io/hashicorp/junos-<device-type>" = "<path-to-go/bin>"
	}
	direct {}
}
```

Example:
```
provider_installation {
	dev_overrides {
		"registry.terraform.io/hashicorp/junos-vqfx" = "/Users/patelv/go/bin"
	}
	direct {}
}
```

You should know have a file structure which looks similar to: 
* (if you created one terraform test file)

```
/junos-terraform/<testing-folder-name>/
/junos-terraform/<testing-folder-name>/main.tf     <-- contents of jtaf-xml2tf command

/Users/<username>/.terraformrc     <-- link to provider created in /usr/go/bin/ [see details above]
```

OR:
* (if you used the -d flag during the `jtaf-xml2tf` command and created a directory of multiple terraform test files)

```
/junos-terraform/<testing-folder-name>/	 <-- contents of jtaf-xml2tf command
/junos-terraform/<testing-folder-name>/dc1-borderleaf1.tf
/junos-terraform/<testing-folder-name>/dc1-borderleaf2.tf
/junos-terraform/<testing-folder-name>/dc1-leaf1.tf
/junos-terraform/<testing-folder-name>/dc1-leaf2.tf  
/junos-terraform/<testing-folder-name>/dc1-leaf3.tf 
/junos-terraform/<testing-folder-name>/dc1-spine1.tf
/junos-terraform/<testing-folder-name>/dc1-spine2.tf 
/junos-terraform/<testing-folder-name>/dc2-spine1.tf
/junos-terraform/<testing-folder-name>/dc2-spine2.tf 

/Users/<username>/.terraformrc     <-- link to provider created in /usr/go/bin/ [see details above]
```

#### Setting Up Host Names

In the test file(s), devices being configured are specified using the `host` field as shown below:
```
provider "junos-vqfx" {
    host     = "dc1-leaf1"
    port     = 22
    username = ""
    password = ""
    alias    = "dc1_leaf1"
}
```

You can either specify the exact IP address in the host field OR use a hostname (like in the example above) and provide the IP address for every hostname in the system file `/etc/hosts`.

Example:
```
127.0.0.1       localhost
<IP address>    dc1-leaf1
<IP address> 	dc1-leaf2
<IP address> 	dc1-leaf3
<IP address> 	dc2-spine1
<IP address> 	dc2-spine2
<IP address> 	dc1-spine1
<IP address> 	dc1-borderleaf2
<IP address> 	dc1-borderleaf1
<IP address> 	dc1-firewall1
<IP address> 	dc1-firewall2
<IP address> 	dc2-firewall1
<IP address> 	dc1-spine2
<IP address>	dc2-firewall2
```

---

### <u>Edit Test Files, Plan, and Apply</u>

Once the `.terraform.rc` file is set up, and the `main.tf` OR group of test file(s) contains access to the provider, information regarding the desired devices to push the configuration to, and the desired config in `HCL` format, we are now ready to use the provider.

```
terraform plan
terraform apply -auto-approve
```

---
### <u>Generate Ansible Playbook</u>

Create an Ansible role + playbook from a Junos JSON schema and one or more XML configs. The generated playbook runs locally and renders configs (does not connect to devices) by default.

Quick usage:
```
jtaf-ansible -j <junos.json> -x <config1.xml> [-x <config2.xml> ...] -t <device-type>
```

What is created (under ansible-provider-junos-<type>/):
- roles/<type>_role/ (tasks/main.yml, templates/template.j2)
- jtaf-playbook.yml (uses connection: local)
- host_vars/, configs/, trimmed_schema.json

Verify rendering without applying:
```
cd ansible-provider-junos-<type>
ansible-playbook -i hosts jtaf-playbook.yml --check --diff
```

---

### <u>Single command to generate ansible role</u>

Generate an Ansible role + playbook in one step from YANG files and XML config(s):

```
jtaf-yang2ansible -p <path-to-common> <path-to-yang-files> -x <xml-config(s)> -t <device-type>
```

Example:
```
jtaf-yang2ansible -p ../yang/18.2/18.2R3/common ../yang/18.2/18.2R3/junos-qfx/conf/*.yang -x examples/evpn-vxlan-dc/dc1/*spine*.xml -t qfx
```

Notes:
- If supplying multiple XML configs they must be for the same device type.
- Output directory: ansible-provider-junos-<type>/ containing roles/<type>_role/ (tasks/templates), jtaf-playbook.yml (connection: local), host_vars/, configs/, trimmed_schema.json.
- Run the generated playbook in check/diff mode to verify rendered configs without applying:
  ansible-playbook -i hosts jtaf-playbook.yml --check --diff

---

### <u>Generate YAML for Ansible host_vars (jtaf-xml2yaml)</u>

Convert one or more Junos XML configs into Ansible host_vars YAML and a simple hosts file.

Usage:
```
jtaf-xml2yaml -j <trimmed_schema.json> -x <config1.xml> [<config2.xml> ...] -d <output-dir>
```

Example:
```
jtaf-xml2yaml -j ansible-provider-junos-qfx/trimmed_schema.json \
  -x examples/qfx/device1.xml examples/qfx/device2.xml \
  -d ansible-provider-junos-qfx
```

Output:
- Creates host_vars/<hostname>.yaml for every XML file provided (hostname is file base name or system/host-name from XML).
- Writes a simple hosts file at <output-dir>/hosts listing all hostnames.

This is useful to feed generated host_vars into the Ansible role/playbook created by jtaf-ansible/jtaf-yang2ansible.


---

### Running Unit Tests

Unit tests are created and stored under /junosterraform/unit_tests/

Run from home /junos-terraform directory:
```
pytest -v
```

Current unit tests:
	- test_worflow.py
		- Ensures that the current workflow will emit the correct output .tf files from a specified set of yang_files (ex: 18.2) based the files located in /examples directory which is considered a source of truth.
