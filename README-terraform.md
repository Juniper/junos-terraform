# Junos Terraform Provider Guide

This guide covers generating, building, testing, and using a custom Junos Terraform provider with JTAF. For initial setup and YANG → JSON conversion, see the [main README](README.md).

---

## Generate Resource Provider

Every provider is built from the same Go source. What differs between them is only the YANG model
embedded in the binary, and **you must say which model you want**: `-x` to trim it to your
configuration, or `--generic` to embed the whole thing. Passing neither is an error, and passing both
is an error.

```bash
# Trimmed to the configuration in the XML you give it
jtaf-provider -j <json-file> -x <xml-configuration(s)> -t <device-type>

# The whole model, no XML needed
jtaf-provider -j <json-file> --generic -t <device-type>
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

Use `jtaf-yang2go` command to generate a resource provider in a single step by supplying all YANG files with the `-p` option, the device XML configuration with `-x`, and the device type with `-t`. As with `jtaf-provider`, exactly one of `-x` or `--generic` is required.

```bash
# Trimmed
jtaf-yang2go -p <path-to-common> <path-to-yang-files> -x <xml-configuration(s)> -t <device-type>

# Whole model
jtaf-yang2go -p <path-to-common> <path-to-yang-files> --generic -t <device-type>
```

Example:

```bash
jtaf-yang2go -p examples/yang/18.2/18.2R3/common examples/yang/18.2/18.2R3/junos-qfx/conf/*.yang -x examples/evpn-vxlan-dc/dc1/*{spine,leaf}*.xml examples/evpn-vxlan-dc/dc2/*spine*.xml -t vqfx
```
NOTE: If using multiple xml configurations (like the example above), ensure that the configurations are for the same device type

NOTE: The examples in this README use the YANG files shipped in this repository under `examples/yang/18.2`.

### Schema scope

The provider source is the same whichever way it is generated; only the model it embeds differs.
Exactly one of `-x` or `--generic` is required.

| Option | Embedded model |
|--------|----------------|
| `-x <xml>` | Trimmed to the paths the XML configuration uses |
| `--generic` | The whole model, untrimmed (mutually exclusive with `-x`) |
| `--exclude PATH` | Leaves the subtree at PATH, relative to `configuration`, out. Repeatable |

#### Which one should I use?

The trade is **what you can write in your `.tf`** against **how much the provider costs to run**.

| | `-x` (trimmed) | `--generic` (full model) |
|---|---|---|
| You can configure | Only stanzas present in the XML you trimmed against | Any stanza the device supports |
| Adding a new stanza later | Re-generate and re-build the provider | Just write it in the `.tf` |
| Memory at plan time | Lowest | Higher |
| Best for | A stable, known configuration set | Exploring, or many devices with differing config |

Start with `-x` if you know which stanzas you manage. Move to `--generic` when re-generating the
provider every time you want a new stanza becomes the annoyance.

#### What it costs

Measured on an Apple M4 Pro against the QFX 18.2 model (`junos.json`, 268 MB of pyang JSON,
593,402 nodes). Your numbers will differ with the model and machine, but the ratios hold.

| | `-x` (one device XML) | `--generic` |
|---|---|---|
| `jtaf-provider` run | ~6 s | ~10 s |
| `go build .` | ~2 s | ~2 s |
| `schema.bin.gz` embedded | 1.8 KB | 883 KB |
| Provider binary | 22 MB | 23 MB |

The build does not get slower with the full model, because no Go source is generated from it — the
schema is data the binary reads at startup, not code the compiler has to chew through.

#### Why it is fast at runtime

`jtaf-provider` compiles the schema into a fixed-size record table (`schema.bin.gz`) and the provider
embeds that. At startup it is read straight into memory with no JSON parsing:

| | Embedded as JSON | Embedded compiled |
|---|---|---|
| Size on disk | 267.7 MB (7.1 MB gzipped) | 11.4 MB (1.7 MB gzipped) |
| Load time | 2.53 s | **4.7 ms** |

That is ~538x faster to load. Terraform starts the provider several times for a single plan, so this
is paid repeatedly. After loading, the full model occupies about **14 MB** of live heap.

`--exclude` takes any configuration path, at any depth, naming the nodes a device would show:
`logical-systems`, `system/services/web-management`, `routing-instances/instance/protocols`,
`vlans/vlan/vlan-id`. YANG `choice` and `case` nodes group nodes in the model but are not
configuration, so a path reaches through them and cannot name one. A path that does not exist is an
error rather than a silent no-op, so a typo does not leave the subtree in place.

Use it to trim a full model down without going back to `-x`:

```bash
jtaf-yang2go --generic -p <common> <yang-files> -t srx \
  --exclude groups --exclude logical-systems --exclude tenants --exclude dynamic-profiles
```

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

### What the generated directory contains

```
terraform-provider-junos-<type>/
├── main.go                  ← entry point; calls generic.Serve
├── embed_schema.go          ← go:embed of schema.bin.gz
├── schema.bin.gz            ← the compiled schema the provider serves
├── trimmed_schema.json.gz   ← the same schema as JSON, for jtaf-xml2tf
├── go.mod / go.sum          ← copied; only the module name is rewritten
├── generic/                 ← schema-driven provider
├── patch/                   ← NETCONF patch engine
├── netconf/                 ← NETCONF client
└── cmd/                     ← compileschema
```

Only `main.go`, `embed_schema.go` and the two schema files are written per device type. Everything
else is copied unchanged, so two providers for different devices differ only in their module name,
provider type name and embedded schema.

Generating a provider needs **Go on PATH**, because the schema is compiled with `cmd/compileschema`.

### Upgrading from JTAF 2.x

Go source is no longer rendered from Jinja2 templates. `resource_config_provider.go.j2`,
`provider.go.j2` and `config.go.j2` are gone, and generated directories no longer contain
`resource_config_provider.go`, `provider.go` or `config.go`. The resource type, provider block and
attribute names are unchanged, so **existing `.tf` files and Terraform state keep working** — rebuild
the provider binary to pick this up.

The plain `trimmed_schema.json` is no longer written; pass the `.gz` file to `jtaf-xml2tf -j` and
`jtaf-xml2yaml -j`. Both accept plain or gzipped JSON, detected by content, so directories generated
before this change still load.

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
| **Create** | Reconciles the device against the plan: reads the running configuration, computes the difference, and sends it as one `edit-config`. A device is not a blank slate, and a plain merge cannot remove configuration it already has. |
| **Read** | Device state is fetched via `get-config` and compared to Terraform state |
| **Update** | Same reconcile as Create: only changed leaves are sent via minimal `edit-config` (patch engine) |
| **Delete** | Targeted deletes are sent for each managed leaf/container |

Create and Update return the **planned** value as the new state, not what the device reads back.
Junos normalises what it is given and keeps configuration the plan never mentioned, so returning the
read-back would differ from the plan and Terraform would fail the apply with "Provider produced
inconsistent result after apply". Anything the device holds beyond the plan is drift, which Read
reports on the next refresh.

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

After the patch is committed the provider re-reads the device. If the configuration still differs
from the plan — Junos auto-generates configuration that is not in the YANG schema, for instance — it
merges the whole planned configuration with `load-configuration` and commits again, so an incomplete
patch cannot leave the device half-configured.

A load or commit the device rejects is treated as an error, and the candidate is discarded. The
candidate is shared rather than per-session, so configuration left staged there would otherwise be
picked up by the next commit from any source.

#### Schema mismatch warning

The provider and your `.tf` files are produced by two separate commands, and only the `-j` path given
to `jtaf-xml2tf` ties them together. Generating files from a trimmed schema and applying them with a
full-model provider passes `terraform validate` and then fails at the device, because the two
schemas describe the same configuration differently.

Each provider records the schema it was built from — its scope plus the YANG modules' revision dates
— and warns when a `.tf` file declares a different one, naming both. Files or providers generated
before this carry no fingerprint and stay silent. Changing which paths are excluded does not
invalidate existing files.

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
