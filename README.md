# JUNOS Terraform Automation Framework (JTAF)

Terraform is traditionally used for managing virtual infrastructure, but there are organisations out there that use Terraform end-to-end and also want to manage configuration state using the same methods for managing infrastructure. Sure, we can run a provisioner with Terraform, but that wasn't asked for!

Much the same as you can use Terraform to create an AWS EC2 instance, you can manage the configurational state of Junos. In essence, we treat Junos configuration as declarative resources.

So what is JTAF? It's a framework, meaning, it's an opinionated set of tools and steps that allow you to go from YANG models to a custom Junos Terraform provider. With all frameworks, there are some dependencies.

To use JTAF, you'll need machine that can run **Go, Python, Git and Terraform.** This can be Linux, OSX or Windows. Some easy to consume videos are below.

## Quick start

### <u>Setup</u>
Run the following commands to set up the Junos-Terraform Environment and Workflow

```bash
git clone https://github.com/aburston/junos-terraform
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
pyang --plugindir ./pyang_plugin -f jtaf -p <path-to-common> <path-to-yang-files> > junos.json
```
Example: 
```
pyang --plugindir ./pyang_plugin -f jtaf -p ../yang/18.2/18.2R3/common ../yang/18.2/18.2R3/junos-qfx/conf/*.yang > junos.json
```
---

### <u>Add Desired Configuration</u>
Now copy an example OR add a real desired device configuration to home directory. [See below using example for a dc2 qfx spine]
```
cp examples/evpn-vxlan-dc/dc2/dc2-spine1.xml config.xml
```

---

### <u>Generate Resource Provider</u>

Now run the following command to generate a `resource provider`. 

```bash
generate_provider.py -j <json-file> -x <xml-configuration> -t <device-type>
```

Example:
```bash
generate_provider.py -j junos.json -x config.xml -t vqfx
```

All in one example (`-j` accepts `-` for `stdin` for `generate_provider.py`):
```bash
pyang --plugindir ./pyang_plugin -f jtaf -p ../yang/18.2/18.2R3/common ../yang/18.2/18.2R3/junos-qfx/conf/*.yang | generate_provider.py -j - -x config.xml -t vqfx
```

---

### <u>Build the provider and install</u>

cd into the newly created directory starting with `terraform-provider-junos-` then the device-type and then `go install`

example:

```
cd terraform-provider-junos-vqfx
go install
```

---

### <u>Autogenerate Terraform Testing Files</u>

Run this command to create a `.tf` test file to deploy the terraform provider.
```
populate-tf.py config.xml
cd testbed
```

In the `/testbed` folder created by the previous command, create a `.terraform.rc` file with `vi` and add the following contents, replacing any `<elements>` tags with your own information. This is to ensure that the terraform plugin you created and installed to `/go/bin` will be read.
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

---

### <u>Update `provider.go`</u>

Line 73 of `provider.go` needs to match the device type you intend to use (from "registry.terraform.io/hashicorp/junos-vqfx"). 

For example, for vsrx devices, update line 73 in `provider.go` to:
```	
resp.TypeName = "junos-vsrx"
```

---

### <u>Edit Test Files, Plan, and Apply</u>

Once the `.terraform.rc` file is set up, and the `main.tf` test file contains access to the provider, information regarding the desired devices to push the configuration to, and the desired config in `HCL` format, we are now ready to use the provider.

```
terrafrom plan
terraform apply -auto-approve
```
