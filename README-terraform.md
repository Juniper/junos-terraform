# Junos Terraform Provider Guide

This guide walks through building a custom Junos Terraform provider and using it, one command at a
time. For installing the tools themselves, see the [main README](README.md).

Start in the repository root with its Python virtual environment active. The concrete
commands below use the bundled QFX 18.2 models; substitute models matching your
platform and Junos release for a real deployment. Angle-bracket arguments are
placeholders to replace, not literal shell syntax.

The path is six steps:

| Step | Command | Produces |
|------|---------|----------|
| 1 | `pyang` | `junos.json` — the YANG models as JSON |
| 2 | `jtaf-provider` | A provider source directory |
| 3 | `go install .` | The provider binary |
| 4 | `jtaf-xml2tf` | `.tf` files describing your devices |
| 5 | edit `~/.terraformrc` | Terraform can find the binary |
| 6 | `terraform plan` / `apply` | Configuration on the device |

Steps 1 and 2 can be combined; see [Shortcut](#shortcut-run-steps-1-and-2-together).

---

## Step 1 — Convert the YANG models to JSON

Find the Junos version your device runs and locate the matching `yang` and `common` folders, then:

```bash
pyang --plugindir $(jtaf-pyang-plugindir) -f jtaf -p <path-to-common> <path-to-yang-files> > junos.json
```

Example, for 18.2:

```bash
pyang --plugindir $(jtaf-pyang-plugindir) -f jtaf -p examples/yang/18.2/18.2R3/common examples/yang/18.2/18.2R3/junos-qfx/conf/*.yang > junos.json
```

This repository ships YANG examples under `examples/yang/18.2`.

---

## Step 2 — Generate the provider

Every provider is built from the same Go source. What differs between them is only the YANG model
embedded in the binary — so before you run anything, you have one decision to make.

### Decide the schema scope

**You must pass exactly one of `-x` or `--generic`.** Passing neither is an error, and so is passing
both.

| | `-x <xml>` | `--generic` |
|---|---|---|
| Embeds | Only the paths your XML configuration uses | The whole model |
| You can configure | Only stanzas present in that XML | Any stanza the device supports |
| To manage a new stanza | Re-generate and re-build the provider | Just write it in the `.tf` |
| Memory at plan time | Lowest | Higher |

Start with `-x` if you know which stanzas you manage. Move to `--generic` when re-generating the
provider every time you want a new stanza becomes the annoyance. Measured timings and sizes for both
are in [What it costs](#what-it-costs).

### Run jtaf-provider

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

| Flag | Required | Meaning |
|------|----------|---------|
| `-j` | yes | The JSON from step 1. `-` reads it from stdin |
| `-x` | one of | XML configuration(s) to trim the schema to |
| `--generic` | one of | Embed the whole model instead |
| `-t` | yes | Device type; names the output directory |
| `--exclude PATH` | no | Leave a subtree out. Repeatable |
| `--groups` | no | Keep the `groups` subtree. See below |

If you pass several XML files, they must all be for the same device type.

Because `-j` accepts `-`, steps 1 and 2 can be piped together:

```bash
pyang --plugindir $(jtaf-pyang-plugindir) -f jtaf -p examples/yang/18.2/18.2R3/common examples/yang/18.2/18.2R3/junos-qfx/conf/*.yang | jtaf-provider -j - -x examples/evpn-vxlan-dc/dc1/*{spine,leaf}*.xml -t vqfx
```

Generating a provider needs **Go on PATH**, because the schema is compiled during generation.

### Narrow a full model with --exclude

`--exclude` takes any configuration path, at any depth, naming the nodes a device would show:
`logical-systems`, `system/services/web-management`, `routing-instances/instance/protocols`,
`vlans/vlan/vlan-id`. YANG `choice` and `case` nodes group nodes in the model but are not
configuration, so a path reaches through them and cannot name one. A path that does not exist is an
error rather than a silent no-op, so a typo does not leave the subtree in place.

Use it to cut a full model down without going back to `-x`:

```bash
jtaf-provider -j junos.json --generic -t srx \
  --exclude groups --exclude logical-systems --exclude tenants --exclude dynamic-profiles
```

### Keep configuration groups with --groups

Junos configuration groups are left out by default: in a full model the `groups` subtree repeats the
whole configuration hierarchy and is about half of its nodes. Pass `--groups` to keep it, along with
the `apply-groups` leaf-list, so the provider manages groups as ordinary configuration:

```bash
jtaf-provider -j junos.json -x <xml-configuration(s)> -t <device-type> --groups
```

A group is then written as an entry of the `groups` list, keyed by its `name`, holding the same
attributes it would have in the base hierarchy, and `apply_groups` is an ordered list of group names:

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

`jtaf-xml2tf` in step 4 takes the same flag, and the two must agree.

NOTE: `--groups` with `--generic` makes the provider advertise the whole model twice over, which needs
several GB of memory at plan time and warns when generated. Prefer trimming with `-x`, or use
`--exclude` on paths inside `groups`.

### Shortcut: run steps 1 and 2 together

`jtaf-yang2go` runs `pyang` and `jtaf-provider` for you. It takes the YANG files with `-p` and passes
`-x`, `--generic`, `--exclude` and `--groups` straight through, so the same scope rule applies:

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

---

## Step 3 — Build and install the provider

`cd` into the directory just created — `terraform-provider-junos-` followed by your device type — and
install it:

```bash
cd terraform-provider-junos-vqfx
go install .
cd ..
```

That puts the binary in `$(go env GOBIN)` if set, otherwise `$(go env GOPATH)/bin`.
Step 5 points Terraform at that directory. `cd ..` returns you to the repository root
for the remaining generation commands.

---

## Step 4 - Generate Terraform configuration

Convert your device XML into Terraform files using the schema from step 2:

```bash
jtaf-xml2tf -j terraform-provider-junos-vqfx/trimmed_schema.json.gz \
  -x examples/evpn-vxlan-dc/dc1/*{spine,leaf}*.xml examples/evpn-vxlan-dc/dc2/*spine*.xml \
  -t vqfx -d testbed
```

| Flag | Meaning |
|------|---------|
| `-j` | Schema from the provider you built; plain or gzipped JSON is accepted |
| `-x` | XML configurations for devices of the same type |
| `-t` | The same device-type suffix used in step 2 |
| `-d` | Output directory for `providers.tf` and per-device `.tf` files |
| `-u`, `-p` | Optional username and password written into provider blocks |
| `--groups` | Preserve groups; use only with a provider generated with `--groups` |
| `--no-extract-common` | Keep values inline instead of extracting shared values into `common.tf` |

By default, inherited configuration groups are flattened into the base hierarchy.
If you kept groups in step 2, add `--groups` here too. The generated files describe
the desired configuration; they are not just test fixtures. Review them before deployment.
Avoid passing real passwords on the command line, where shell history can retain them.

---

## Step 5 - Tell Terraform where to find the provider

Find the directory containing the installed binary:

```bash
go env GOBIN GOPATH
```

Use `GOBIN` when nonempty, otherwise the `bin` directory under `GOPATH`.
Add the following block to `~/.terraformrc`, preserving any existing settings and
replacing the path with your actual binary directory:

```hcl
provider_installation {
  dev_overrides {
    "registry.terraform.io/hashicorp/junos-vqfx" = "/absolute/path/to/go/bin"
  }
  direct {}
}
```

Use the provider type from step 2 in the registry address. On Windows, use
`terraform.rc` in `%APPDATA%` instead of `~/.terraformrc`.

---

## Step 6 - Review, plan, and apply

Enable NETCONF over SSH on the devices and ensure their management addresses are
reachable. In `testbed/providers.tf`, replace host names or addresses, ports, and
credentials with your device settings. Use port 830 for the NETCONF SSH service,
or the port your device actually exposes. Each resource must reference the correct
provider alias. Review the per-device configuration and any shared `common.tf` values.

From the repository root:

```bash
cd testbed
terraform validate
terraform plan
terraform apply
```

With this development override, Terraform uses the locally installed provider;
you do not need `terraform init` just to install it. If you add modules or other
providers, initialize those dependencies separately. Terraform warns that a
development override is active; this is expected.

`terraform apply` displays the plan and asks for confirmation. Check the devices
and changes before approving. These example configurations are intended for a lab,
not an unreviewed deployment to production devices.

The command walkthrough ends here. The following sections explain generated files,
upgrade considerations, performance, and provider internals.

---

## Reference

### What it costs

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

### Why it is fast at runtime

`jtaf-provider` compiles the schema into a fixed-size record table (`schema.bin.gz`) and the provider
embeds that. At startup it is read straight into memory with no JSON parsing:

| | Embedded as JSON | Embedded compiled |
|---|---|---|
| Size on disk | 267.7 MB (7.1 MB gzipped) | 11.4 MB (1.7 MB gzipped) |
| Load time | 2.53 s | **4.7 ms** |

That is ~538x faster to load. Terraform starts the provider several times for a single plan, so this
is paid repeatedly. After loading, the full model occupies about **14 MB** of live heap.

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

Run each block from the repository root, not from the Terraform workspace.

### Patch Engine Tests

```bash
cd terraform_provider
# All patch engine tests
go test ./patch/ -v

# Corner case tests only
go test ./patch/ -run TestCC -v

# With coverage
go test ./patch/ -coverprofile=coverage.out
go tool cover -html=coverage.out
cd ..
```

### Provider Tests

```bash
cd terraform_provider
# All provider tests
go test . -v

# With race detection
go test ./... -race -v
cd ..
```
