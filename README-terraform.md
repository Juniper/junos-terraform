# Junos Terraform Provider Guide

This guide covers generating, building, testing, and using a custom Junos Terraform provider with JTAF. For initial setup and YANG → JSON conversion, see the [main README](README.md).

---

## Generate Resource Provider

Run the following command to generate a `resource provider`.

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
pyang --plugindir $(jtaf-pyang-plugindir) -f jtaf -p examples/yang/18.2/18.2R3/common examples/yang/18.2/18.2R3/junos-qfx/conf/*.yang | jtaf-provider -j - -x examples/evpn-vxlan-dc/dc1/*{spine,leaf}*.xml examples/evpn-vxlan-dc/dc2/*spine*.xml  -t vqfx
```

---

### Single command to generate resource provider

Use `jtaf-yang2go` command to generate a resource provider in a single step by supplying all YANG files with the `-p` option, the device XML configuration with `-x`, and the device type with `-t`.

```bash
jtaf-yang2go -p <path-to-common> <path-to-yang-files> -x <xml-configuration(s)> -t <device-type>
```

Example:

```bash
jtaf-yang2go -p examples/yang/18.2/18.2R3/common examples/yang/18.2/18.2R3/junos-qfx/conf/*.yang -x examples/evpn-vxlan-dc/dc1/*{spine,leaf}*.xml examples/evpn-vxlan-dc/dc2/*spine*.xml -t vqfx
```
NOTE: If using multiple xml configurations (like the example above), ensure that the configurations are for the same device type

NOTE: The examples in this README use the YANG files shipped in this repository under `examples/yang/18.2`.

### Schema scope

The provider source is the same whichever way it is generated; only the model it embeds differs.

| Option | Embedded model |
|--------|----------------|
| `-x <xml>` | Trimmed to the paths the XML configuration uses |
| `--generic` | The whole model, untrimmed (mutually exclusive with `-x`) |
| `--exclude PATH` | Leaves the subtree at PATH, relative to `configuration`, out. Repeatable |

`--exclude` takes any configuration path, at any depth, naming the nodes a device would show:
`logical-systems`, `system/services/web-management`, `routing-instances/instance/protocols`,
`vlans/vlan/vlan-id`. YANG `choice` and `case` nodes group nodes in the model but are not
configuration, so a path reaches through them and cannot name one. A path that does not exist is an
error rather than a silent no-op, so a typo does not leave the subtree in place.

### Configuration groups

Junos configuration groups are left out by default: in a full model the `groups` subtree repeats the
whole configuration hierarchy and is about half of its nodes. Pass `--groups` to keep it, along with
the `apply-groups` leaf-list, so the provider manages groups as ordinary configuration:

```bash
jtaf-yang2go -p <path-to-common> <path-to-yang-files> -x <xml-configuration(s)> -t <device-type> --groups
```

A group is then written as an entry of the `groups` list, keyed by its `name`, holding the same
attributes it would have in the base hierarchy, and `apply_groups` is an ordered list of group names.

`jtaf-xml2tf` takes the same flag, and the two need to agree. By default it flattens: the
configuration a device inherits through `apply-groups` is merged into the base hierarchy and the
groups themselves are dropped, which is what a provider built without `--groups` expects. With
`--groups` it converts the hierarchy as it stands:

```bash
jtaf-xml2tf -j <trimmed_schema.json.gz> -x <xml-configuration(s)> -t <device-type> -d <output-dir> --groups
```

```hcl
resource "terraform-provider-junos-<device-type>" "dev1-base-config" {
  resource_name = "base-config"
  apply_groups  = ["base"]
  groups = [
    {
      name   = "base"
      system = [{ host_name = "from-group" }]
    }
  ]
}
```

NOTE: `--groups` with `--generic` makes the provider advertise the whole model twice over, which needs
several GB of memory at plan time and warns when generated. Prefer trimming with `-x`, or use
`--exclude` on paths inside `groups`.

---

## Build the Provider and Install

cd into the newly created directory starting with `terraform-provider-junos-` then the device-type and then `go install`

Example:

```
cd terraform-provider-junos-vqfx
go install .
```

---

## Autogenerate Terraform Testing Files

### Overview

Run a command to generate a `.tf` test file to deploy the Terraform provider.

**NOTE:** Output is written to a directory (`-d`) as `providers.tf` plus one `.tf` file per XML input.

**Flag Options:**
 * -j 
	* **Required:** `trimmed_schema.json.gz` output file from jtaf-provider (stored in terraform provider folder /terraform-provider-junos-"device-type"); plain JSON is also accepted
 * -x
	* **Required:** File(s) of xml config to create terraform files for
 * -t
	* **Required:** Junos device type
 * -d
	* **Required:** Output directory where providers.tf and per-device Terraform files are written
 * -u
	* **Optional:** Device username
 * -p
	* **Optional:** Device password

---

### Creating Terraform Testing Files

To create multiple Terraform (.tf) files from multiple config files, where each .tf file will represent one xml file, use the following command (output returned to specified directory name):

```
jtaf-xml2tf -j <path-to-trimmed-schema> -x <path-to-config-files(s)> -t <device-type> -d <testing-folder-name>
```

Example: 

* **trimmed_schema** - `trimmed_schema.json.gz`, stored in the terraform provider folder created from running the jtaf-provider module command (usually in terraform-provider-junos-'device-type')
* **xml_files** - directory containing xml file(s) (ensure xml file(s) are for the same device type)

```
jtaf-xml2tf -j terraform-provider-junos-vqfx/trimmed_schema.json.gz -x examples/evpn-vxlan-dc/dc1/*{spine,leaf}*.xml examples/evpn-vxlan-dc/dc2/*spine*.xml -t vqfx -d testbed
```
* If the user wants to provide the device(s) **username** and **password**, those additional flags can be added as well
```
jtaf-xml2tf -j terraform-provider-junos-vqfx/trimmed_schema.json.gz -x examples/evpn-vxlan-dc/dc1/*{spine,leaf}*.xml examples/evpn-vxlan-dc/dc2/*spine*.xml -t vqfx -d testbed -u root -p password
```

Using the output which is outputted to the specified directory from the command, which represents a template for the HCL .tf file for each input XML file, we can now create our testing environment and fill in the template with any remaining necessary device or config information.

---

### Setting up Testing Environment

Now that we ran the `jtaf-xml2tf` command and have our testing folder setup:
* The command writes files directly under your test folder in the `/junos-terraform` directory.

#### Creating the Environment

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

You should now have a file structure which looks similar to: 
* (if you created one terraform test file)

```
/junos-terraform/<testing-folder-name>/
/junos-terraform/<testing-folder-name>/providers.tf
/junos-terraform/<testing-folder-name>/<hostname>.tf

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

You can either specify the exact IP address in the host field OR use a hostname (like in the example above) and provide the IP address for every hostname in the system file `/etc/hosts` using `vi`.

*NOTE:* If the `/etc/hosts` file is a **READ-ONLY** file, then try using `sudo su` then re-run `vi /etc/hosts`. Exit after editing and return back to user control. 

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

### Edit Test Files, Plan, and Apply

Once the `.terraformrc` file is set up, and the generated test file(s) contain access to the provider, information regarding the desired devices to push the configuration to, and the desired config in `HCL` format, we are now ready to use the provider.

```
terraform plan
terraform apply -auto-approve
```

---

## How the Provider Works

### Overview

The generated Terraform provider communicates with Junos devices over **NETCONF** (port 830). It manages configuration as a single Terraform resource — the entire config block you defined in your `.tf` file is treated as a declarative resource with full CRUD lifecycle:

| Operation | What Happens |
|-----------|-------------|
| **Create** | Full configuration is pushed via `load-configuration` merge |
| **Read** | Device state is fetched via `get-config` and compared to Terraform state |
| **Update** | Only changed leaves are sent via minimal `edit-config` (patch engine) |
| **Delete** | Targeted deletes are sent for each managed leaf/container |

### NETCONF Patch Engine

Starting with JTAF 2.0.0, generated providers use a **NETCONF patch engine** for Update and Delete operations. Instead of replacing the entire configuration on every change, the provider computes a minimal leaf-level diff and sends only the changed elements via NETCONF `edit-config`.

#### Pipeline

The patch engine lives at `terraform_provider/patch/` and implements a five-stage pipeline:

```
Device XML (get-config) ──┐
                          ├→ BuildTree() → LeafMapWithSchema() → map[path]value
Plan XML (Terraform) ────┘                        ↑
                                            YANG Schema
                                       (trimmed_schema.json.gz)
                                                  ↓
                                         ComputeDiff()
                                                  ↓
                                   Create | Replace | Delete operations
                                                  ↓
                            CreateDiffPatchWithSchema() → NETCONF XML
                                                  ↓
                            AlignXMLOrderToReference() → ordered XML
                                                  ↓
                                  edit-config RPC → Junos device
```

**Stage 1 — BuildTree:** Parses raw XML into an in-memory tree structure.

**Stage 2 — LeafMapWithSchema:** Flattens the XML tree into a `map[string]string` where each key is an XPath-like path (e.g., `configuration/interfaces/interface[name=ge-0/0/0]/description`) and each value is the leaf's text content. Uses the YANG schema to correctly identify list keys, handle empty-type leaves, and track ordered leaf-list positions.

**Stage 3 — ComputeDiff:** Compares the device state map against the Terraform plan map and produces a set of Create, Replace, and Delete operations — only for leaves that actually differ.

**Stage 4 — CreateDiffPatchWithSchema:** Converts the diff operations into NETCONF `<configuration>` XML with per-element `nc:operation` attributes (`create`, `replace`, or `delete`).

**Stage 5 — AlignXMLOrderToReference:** Reorders XML siblings to match a reference document, preventing spurious diffs caused by non-deterministic element ordering in Junos `get-config` responses.

#### Before vs. After

| Aspect | Before (1.1.0) | Now (2.0.0) |
|--------|----------------|-------------|
| **Update strategy** | Delete entire apply-group + re-push all leaves | Compute diff → send only changed leaves |
| **Delete strategy** | `delete groups/<name>` | Targeted per-leaf/container `nc:operation="delete"` |
| **NETCONF RPC** | `load-configuration` (full XML body) | `edit-config` with per-element operations |
| **Blast radius** | Entire resource group | Only changed leaves |
| **Phantom diffs** | Common (element order, encoding) | Eliminated |

#### Schema-Aware Features

The patch engine uses the trimmed schema (`trimmed_schema.json.gz`, generated alongside the provider) to make intelligent decisions:

- **List key detection:** Uses YANG `key` statement instead of hardcoded `name`/`id`/`type` guessing
- **Leaf-list semantics:** Distinguishes `ordered-by user` (position matters) from `ordered-by system` (set semantics)
- **Container presence:** Empty containers are skipped (no spurious "container exists" diffs)
- **Type classification:** Correctly handles `empty`, `identityref`, `union`, `enumeration`, `leafref` types
- **Compound keys:** Supports multi-field keys (e.g., `choice-ident choice-value`)

#### Operation Ordering

The patch engine ensures NETCONF operations are applied in a safe order:

1. **Deletes first** (deepest path first — delete leaf before parent)
2. **Replaces second** (in-place value changes)
3. **Creates last** (shallowest path first — create parent before leaf)

This prevents Junos candidate validation failures from referencing nodes that don't exist yet or deleting nodes that still have children.

#### Fallback Safety

If the patch produces residual differences after verification (e.g., due to Junos auto-generated config not in the YANG schema), the provider:

1. Emits a Terraform **warning** listing the unresolved leaves
2. Falls back to full `load-configuration` replace (same as 1.1.0 behavior)
3. Re-verifies the final state

The provider is never worse than the previous release — it's strictly better or equivalent.

---

## Running Tests

### Patch Engine Tests

```bash
# All patch engine tests
cd terraform_provider && go test ./patch/ -v

# Corner case tests only
cd terraform_provider && go test ./patch/ -run TestCC -v

# With coverage
cd terraform_provider && go test ./patch/ -coverprofile=coverage.out && go tool cover -html=coverage.out
```

### Provider Tests

```bash
# All provider tests
cd terraform_provider && go test . -v

# With race detection
cd terraform_provider && go test ./... -race -v
```
