# Design: Junos Terraform Automation Framework (JTAF)

## Overview

This document describes the technical architecture, component design, and data flow of JTAF — from YANG model ingestion through provider generation to runtime NETCONF operations against Junos devices.

---

## System Architecture

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                        GENERATION PHASE (Python)                            │
│                                                                             │
│  ┌──────────────┐    ┌──────────────┐    ┌──────────────────────────────┐  │
│  │ YANG Models  │───▶│ pyang + JTAF │───▶│ JSON Schema (junos.json)     │  │
│  │ (.yang)      │    │   Plugin     │    │ Full YANG tree as JSON       │  │
│  └──────────────┘    └──────────────┘    └──────────────┬───────────────┘  │
│                                                          │                  │
│  ┌──────────────┐                                        ▼                  │
│  │ XML Configs  │───▶ ┌─────────────────────────────────────────────────┐  │
│  │ (.xml)       │     │ jtaf-provider (filter_json_using_xml)           │  │
│  └──────────────┘───▶ │ 1. Parse XML configs                           │  │
│                        │ 2. Extract XPath-like paths from XML           │  │
│                        │ 3. Filter JSON schema to only matching paths   │  │
│                        │ 4. Render Go code via Jinja2 templates         │  │
│                        │ 5. Write output directory with all files       │  │
│                        └──────────────────────────┬──────────────────────┘  │
│                                                   │                         │
└───────────────────────────────────────────────────┼─────────────────────────┘
                                                    │
                                                    ▼
┌─────────────────────────────────────────────────────────────────────────────┐
│                    OUTPUT: terraform-provider-junos-{type}/                  │
│                                                                             │
│  main.go ─── provider.go ─── config.go ─── resource_config_provider.go     │
│  go.mod ──── trimmed_schema.json ──── patch/ ──── netconf/                  │
│                                                                             │
└───────────────────────────────────────┬─────────────────────────────────────┘
                                        │ go install .
                                        ▼
┌─────────────────────────────────────────────────────────────────────────────┐
│                        RUNTIME PHASE (Go binary)                            │
│                                                                             │
│  ┌────────────┐    ┌────────────────┐    ┌──────────────────────────────┐  │
│  │ Terraform  │───▶│ Plugin         │───▶│ Provider                     │  │
│  │ CLI        │    │ Framework      │    │ Configure → NETCONF Client   │  │
│  └────────────┘    └────────────────┘    └──────────────┬───────────────┘  │
│                                                          │                  │
│                                           ┌──────────────▼───────────────┐  │
│                                           │ Config Resource CRUD         │  │
│                                           │ Create: load-configuration   │  │
│                                           │ Read:   get-configuration    │  │
│                                           │ Update: diff → edit-config   │  │
│                                           │ Delete: edit-config delete   │  │
│                                           └──────────────┬───────────────┘  │
│                                                          │                  │
│                                           ┌──────────────▼───────────────┐  │
│                                           │ NETCONF Client (SSH:830)     │  │
│                                           │ Framed XML over SSH          │  │
│                                           └──────────────┬───────────────┘  │
│                                                          │                  │
└──────────────────────────────────────────────────────────┼──────────────────┘
                                                           │
                                                           ▼
                                                   ┌──────────────┐
                                                   │ Junos Device │
                                                   │ (NETCONF)    │
                                                   └──────────────┘
```

---

## Component Design

### 1. pyang Plugin (`jtaf_pyang_plugin/jtaf_json.py`)

**Purpose:** Convert YANG AST into a JSON tree consumable by the code generator.

**Design Decisions:**
- Registered as pyang output format `"jtaf"` so it integrates into the standard pyang pipeline
- Walks YANG AST depth-first, building `FNode` objects for each statement
- Resolves types at extraction time (leafref → target type, typedef → base type)
- Flattens `choice`/`case` — Terraform users shouldn't see YANG's branching constructs
- Expands `grouping`/`uses` — JSON output is fully denormalized
- Preserves `ordered-by user` flag — critical for leaf-lists where order matters (VRRP virtual-addresses)
- Outputs to stdout for piping to `jtaf-provider`

**Type Mapping:**
```
YANG Type          → JSON Representation     → Go/Terraform Type
─────────────────────────────────────────────────────────────────
string             → "string"                → schema.StringAttribute
uint8..uint64      → "integer"               → schema.Int64Attribute
int8..int64        → "integer"               → schema.Int64Attribute
boolean            → "boolean"               → schema.BoolAttribute
empty              → "empty"                 → schema.BoolAttribute (presence)
identityref        → "identityref"           → schema.StringAttribute
leafref            → resolved target type    → depends on target
union              → "union"                 → schema.StringAttribute
enumeration        → "enumeration"           → schema.StringAttribute
inet:*             → "string"                → schema.StringAttribute
```

---

### 2. CLI Tools (`junosterraform/`)

**Architecture:** Click-based Python CLI entry points defined in `setup.py`.

#### jtaf-provider

**Core algorithm (`jtaf_common.py`):**

```
1. load_json(schema)        → Parse full YANG JSON tree
2. parse_xml(configs)       → Extract XML element paths from device configs
3. filter_json_using_xml()  → Intersect: keep only JSON nodes matching XML paths
4. render_templates()       → Jinja2 templates → Go source files
5. post_process()           → Rewrite imports, ensure go.mod module name
6. write_output()           → Write to terraform-provider-junos-{type}/
```

**Filtering Design:**
- XML configs define the "scope" of the provider — which parts of the YANG tree to include
- The XML is not deployed directly; it's used as a pattern for what Terraform should manage
- Multiple XML configs are merged (union of paths) so one provider covers multiple device roles

#### jtaf-yang2go

**Design:** Single-command wrapper that internally:
1. Invokes `pyang` with the JTAF plugin (YANG → JSON)
2. Passes output directly to `jtaf-provider` logic (JSON + XML → Go code)

Eliminates the need for intermediate JSON files or shell piping.

#### jtaf-xml2tf

**Purpose:** Generate Terraform `.tf` test files from XML configurations.

**Design:**
- Reads `trimmed_schema.json` (generated by jtaf-provider) for type information
- Parses each XML config, maps elements to Terraform HCL attributes
- Generates one `.tf` file per XML config + a shared `providers.tf`
- Fills in resource definitions with HCL syntax matching the generated provider schema

#### jtaf-xml2yaml

**Purpose:** Convert XML configs into Ansible `host_vars` and `group_vars`.

**Design:**
- Uses `trimmed_schema.json` for type-aware XML→YAML conversion
- Requires `--grouping-hosts-file` for inventory group definitions
- Merge-safe repeated runs: shared values float to `group_vars/all.yaml`, per-group to `group_vars/<group>/all.yaml`, host-specific stay in `host_vars/`
- Produces inventory file with `[all]`, `[group]`, and `[group:children]` sections

#### jtaf-ansible / jtaf-yang2ansible

**Purpose:** Generate complete Ansible role + playbook from YANG + XML.

**Output structure:**
```
ansible-provider-junos-{type}/
├── roles/{type}_role/
│   ├── tasks/main.yml
│   └── templates/template.j2
├── jtaf-playbook.yml        (connection: local, renders config)
├── host_vars/
├── group_vars/
├── configs/
└── trimmed_schema.json
```

---

### 3. Code Generation Templates (`junosterraform/templates/`)

**Template Engine:** Jinja2 (Python)

**Key Context Variables:**
| Variable | Source | Purpose |
|----------|--------|---------|
| `data.device_type` | `-t` flag | Provider name suffix |
| `data['root']['children']` | Filtered JSON | YANG nodes to generate code for |
| `data.schema_json` | Filtered tree | Inlined as `trimmed_schema.json` |

**Templates and Their Outputs:**

| Template | Output | Design Intent |
|----------|--------|---------------|
| `resource_config_provider.go.j2` | `resource_config_provider.go` | One resource type with CRUD, schema attributes generated per YANG leaf |
| `provider.go.j2` | `provider.go` | Provider struct with schema (host, port, username, password, sshkey) |
| `config.go.j2` | `config.go` | Config struct holding connection params, `Client()` factory method |
| `go.mod.j2` | `go.mod` | Module declaration with correct path and dependencies |

**Post-Generation Steps:**
1. `ensure_go_module_name()` — Set go.mod module to `terraform-provider-junos-{type}`
2. `rewrite_import_prefixes()` — Replace `terraform_provider/` imports with the correct module path
3. Copy `patch/` and `netconf/` packages into output directory
4. Write `trimmed_schema.json` for runtime schema lookups

---

### 4. Terraform Provider Runtime (`terraform_provider/`)

**Framework:** HashiCorp Terraform Plugin Framework v1.18.0

**Provider Initialization Flow:**
```
main.go
  → providerserver.Serve()
    → Provider.Metadata() → returns type name
    → Provider.Schema() → returns config schema (host, port, user, pass, sshkey)
    → Provider.Configure() → establishes NETCONF client
    → Provider.Resources() → returns [configResource]
```

**Config Resource CRUD:**

| Operation | NETCONF RPC | Behavior |
|-----------|-------------|----------|
| **Create** | `<load-configuration>` + `<commit/>` | Build XML from Terraform plan, load as named group, commit |
| **Read** | `<get-configuration>` | Fetch group config from device, parse back to state |
| **Update** | `<edit-config>` + `<commit/>` | Compute diff via patch engine, send minimal changes |
| **Delete** | `<edit-config operation="delete">` + `<commit/>` | Delete the named group and commit |

**Error Recovery:**
- If `<commit/>` fails → send `<discard-changes/>` → return error to Terraform
- If `<get-configuration>` returns empty for the group → mark resource as gone (triggers recreate)

---

### 5. NETCONF Client (`terraform_provider/netconf/`)

**Design:** Low-level SSH-over-NETCONF client using `nemith.io/netconf` + `golang.org/x/crypto/ssh`.

**Connection Flow:**
```
1. SSH dial to host:port
2. Authenticate (password or key-based)
3. Request "netconf" subsystem
4. Exchange <hello> capabilities
5. Session ready for RPC
```

**RPC Protocol:**
- Messages framed with `]]>]]>` end-of-message delimiter (NETCONF 1.0)
- Each RPC has unique `message-id` for correlating replies
- Client serializes requests as XML, parses XML responses
- Supports `<rpc-error>` detection and propagation

---

### 6. Patch Engine (`terraform_provider/patch/`)

**Purpose:** Compute and apply minimal NETCONF operations to move from current state to desired state, rather than replacing entire configuration groups.

**Three-Stage Pipeline:**

```
                    ┌─────────────────┐
Current XML ───────▶│ LeafMapWithSchema│───▶ map[string]string (current)
                    └─────────────────┘
                    ┌─────────────────┐
Desired XML ───────▶│ LeafMapWithSchema│───▶ map[string]string (desired)
                    └─────────────────┘
                              │
                              ▼
                    ┌─────────────────┐
                    │   ComputeDiff   │───▶ []DiffEntry (adds, deletes, modifies)
                    └─────────────────┘
                              │
                              ▼
                    ┌──────────────────────────┐
                    │ CreateDiffPatchWithSchema │───▶ XML <edit-config> payload
                    └──────────────────────────┘
                              │
                              ▼
                    ┌──────────────────────────┐
                    │ AlignXMLOrderToReference  │───▶ Reordered XML (Junos-compatible)
                    └──────────────────────────┘
```

#### LeafMapWithSchema

**Design:** Recursively walk XML tree, using `trimmed_schema.json` to determine:
- Which elements are **list keys** (identity path segments, not leaves)
- Which elements are **leaf-lists** (indexed by `[value=X]` in path)
- Which are **containers** vs. **leaves**

**Path Format:** `configuration/interfaces/interface[name=ge-0/0/0]/unit[name=0]/family/inet/address[name=10.0.0.1/30]`

#### ComputeDiff

**Design:** Set comparison of flattened paths:
- Paths in desired but not current → `add`
- Paths in current but not desired → `delete`
- Paths in both with different values → `modify`

#### CreateDiffPatchWithSchema (Two-Pass)

**Design:**
- **Pass 1 (deletes):** Generate `<element nc:operation="delete">` entries
- **Pass 2 (adds/modifies):** Generate `<element nc:operation="create|merge">` entries
- Uses schema to reconstruct hierarchical XML from flat paths
- Inserts list key elements for context even when they're not changing

#### AlignXMLOrderToReference

**Design:** Junos is order-sensitive for certain elements. This function reorders the generated patch XML to match the order found in the reference (current) config, preventing spurious ordering diffs.

#### ProcessSchema

**Design:** Loads `trimmed_schema.json` at provider startup into a `map[string]*NodeInfo` index keyed by XPath-like path. Used by LeafMap to make schema-informed decisions about node types.

---

### 7. Configuration Groups Model

**Design Decision:** JTAF uses Junos **configuration groups** as the deployment mechanism.

**How it works:**
```xml
<!-- What JTAF creates on the device -->
<configuration>
  <groups>
    <name>evpn-underlay</name>
    <interfaces>
      <interface>
        <name>ge-0/0/0</name>
        <unit><name>0</name>...</unit>
      </interface>
    </interfaces>
  </groups>
  <apply-groups>evpn-underlay</apply-groups>
</configuration>
```

**Why groups:**
1. Clean separation — JTAF-managed config doesn't mix with manual/other-tool config
2. Atomic operations — entire group can be loaded, replaced, or deleted as a unit
3. Terraform resource mapping — one resource = one group = one declarative config blob
4. `apply-groups` activates the group without modifying the base config
5. Multiple Terraform resources can coexist on one device (different groups)

---

### 8. Mock NETCONF Server (`netconf_mock/`)

**Purpose:** Enable full CI/CD testing without physical Junos devices.

**Design:**
- AsyncSSH-based server with configurable per-device TCP ports
- Implements Junos candidate/running config lifecycle
- DeviceState per device holds `running_groups`, `candidate_groups`, `deleted_candidate_groups`
- Supports all RPCs the provider uses: load-configuration, edit-config (patch + delete), commit, discard-changes, get-configuration
- Recursive merge logic for `edit-config` patch mode mirrors how real Junos handles NETCONF operations

**Multi-Device Architecture:**
```
                ┌─ Port 8301 ─── DeviceState("dc1-leaf1")
Mock Server ────┼─ Port 8302 ─── DeviceState("dc1-leaf2")
                └─ Port 8303 ─── DeviceState("dc1-spine1")
```

Each device is fully independent — separate state, separate SSH listener.

---

### 9. Testing Architecture

**Design Principle:** Tests at every layer, no device dependency for CI.

```
┌─────────────────────────────────────────────────────────────┐
│ Unit Tests (fast, no I/O)                                   │
│  • Go: patch/*, provider structs, file resource helpers     │
│  • Python: pyang plugin type resolution, CLI arg parsing    │
├─────────────────────────────────────────────────────────────┤
│ Integration Tests (mock server)                             │
│  • Go: provider CRUD against mock client                    │
│  • Python: mock server state machine, RPC dispatch          │
├─────────────────────────────────────────────────────────────┤
│ E2E Tests (terraform CLI + mock server)                     │
│  • terraform init/plan/apply/destroy against mock           │
│  • Validates full provider lifecycle                        │
├─────────────────────────────────────────────────────────────┤
│ Live Device Tests (optional, manual)                        │
│  • Connectivity check                                       │
│  • Commit/discard cycle                                     │
│  • Full config matrix                                       │
└─────────────────────────────────────────────────────────────┘
```

**Go Test Pattern:** Table-driven with `t.Run()`:
```go
tests := []struct {
    name    string
    current map[string]string
    desired map[string]string
    want    []DiffEntry
}{...}
```

**Python Test Pattern:** pytest with fixtures for DeviceState setup.

---

## Data Flow: End-to-End Example

### Generating a Provider for EVPN-VXLAN DC Fabric

```
1. User has: QFX switches running Junos 18.2, 9 devices (spines + leaves + border-leaves)
2. User has: XML configs exported from each device

Step 1: YANG → JSON
   pyang --plugindir $(jtaf-pyang-plugindir) -f jtaf \
     -p yang/18.2/18.2R3/common \
     yang/18.2/18.2R3/junos-qfx/conf/*.yang > junos.json
   
   Output: ~50MB JSON file with entire Junos QFX 18.2 config schema

Step 2: JSON + XML → Go Provider
   jtaf-provider -j junos.json \
     -x evpn-vxlan-dc/dc1/*.xml evpn-vxlan-dc/dc2/*.xml \
     -t vqfx
   
   What happens:
   a. Parse all 9 XML configs
   b. Extract unique config paths (interfaces, protocols, routing-options, etc.)
   c. Filter 50MB JSON to only the paths present in the XML (~2MB)
   d. Render Jinja2 templates with filtered schema
   e. Copy patch/ and netconf/ packages
   f. Write trimmed_schema.json
   g. Output: terraform-provider-junos-vqfx/

Step 3: Compile
   cd terraform-provider-junos-vqfx && go install .
   → Binary in $GOPATH/bin/terraform-provider-junos-vqfx

Step 4: Configure Terraform
   # ~/.terraformrc
   provider_installation {
     dev_overrides {
       "registry.terraform.io/hashicorp/junos-vqfx" = "/Users/user/go/bin"
     }
     direct {}
   }

Step 5: Write Terraform Config
   # dc1-spine1.tf
   provider "junos-vqfx" {
     host     = "10.0.0.1"
     port     = 830
     username = "admin"
     password = "secret"
     alias    = "dc1_spine1"
   }
   
   resource "junos-vqfx_ConfigProvider" "spine1-underlay" {
     provider      = junos-vqfx.dc1_spine1
     resource_name = "evpn-underlay"
     # ... attributes from YANG schema
   }

Step 6: Deploy
   terraform plan    → Shows what will change
   terraform apply   → Sends NETCONF RPCs to device
```

### Runtime Update Flow (Patch Engine)

```
1. User modifies HCL attribute (e.g., changes interface description)
2. terraform plan detects difference between state and plan
3. Provider Update() called with old state and new plan

4. Build leaf maps:
   current = LeafMapWithSchema(current_xml, schema_idx)
   desired = LeafMapWithSchema(desired_xml, schema_idx)
   
   current: {"configuration/interfaces/interface[name=ge-0/0/0]/description": "old"}
   desired: {"configuration/interfaces/interface[name=ge-0/0/0]/description": "new"}

5. Compute diff:
   diff = ComputeDiff(current, desired)
   → [{Op: Modify, Path: "...description", Old: "old", New: "new"}]

6. Generate patch XML:
   <edit-config>
     <config>
       <configuration>
         <groups>
           <name>evpn-underlay</name>
           <interfaces>
             <interface>
               <name>ge-0/0/0</name>
               <description nc:operation="merge">new</description>
             </interface>
           </interfaces>
         </groups>
       </configuration>
     </config>
   </edit-config>

7. Send via NETCONF → commit → update Terraform state
```

---

## Key Design Decisions

| Decision | Rationale |
|----------|-----------|
| Generate Go code rather than interpret at runtime | Compiled binary = fast, no Python dependency at runtime, standard Terraform plugin |
| Use YANG as single source of truth | Vendor-published, versioned, machine-readable, covers all platforms |
| XML configs as scope filter | Users already have device configs; reuse them to define what Terraform manages |
| Configuration groups as deployment model | Clean isolation, atomic ops, compatible with manual config |
| Patch engine for updates | Avoid replacing entire groups on minor changes; efficient NETCONF operations |
| trimmed_schema.json in output | Runtime schema awareness without re-parsing YANG |
| Mock server for testing | CI/CD without physical devices, reproducible test state |
| Multi-tool output (Terraform + Ansible) | Same YANG + XML analysis serves both automation camps |
| Plugin Framework (not SDKv2) | Modern HashiCorp API, better type system, protocol v6 |

---

## Security Considerations

- SSH credentials (password/key) stored in Terraform state — users should encrypt state files
- NETCONF over SSH (port 830) — all device communication is encrypted
- Provider binary is compiled Go — no interpreted code at runtime
- Mock server accepts `--disable-auth` flag for testing only — never use in production
- No secrets hardcoded in generated code — credentials provided via Terraform config or env vars

---

## Scalability

| Dimension | Design Approach |
|-----------|-----------------|
| Many devices | One provider binary serves all devices of same type/version; Terraform manages per-device state via resource instances |
| Many config hierarchies | XML filter defines scope; provider includes only relevant schema paths |
| Large YANG trees | JSON filter step reduces 50MB+ full schema to ~2MB trimmed schema |
| Multiple platforms | Generate separate providers per platform/version combination |
| CI/CD pipelines | Mock server enables parallel testing without device contention |

---

## Dependencies

### Generation-Time (Python)

| Package | Version | Purpose |
|---------|---------|---------|
| pyang | ≥2.0 | YANG parser and plugin host |
| jinja2 | ≥3.0 | Template rendering |
| click | ≥8.0 | CLI framework |
| lxml | ≥4.0 | XML parsing |
| pyyaml | ≥5.0 | YAML output for Ansible |

### Runtime (Go)

| Package | Version | Purpose |
|---------|---------|---------|
| terraform-plugin-framework | v1.18.0 | Terraform provider interface |
| nemith.io/netconf | v0.0.4 | NETCONF protocol implementation |
| golang.org/x/crypto | latest | SSH client |

### Testing

| Package | Purpose |
|---------|---------|
| asyncssh (Python) | Mock NETCONF server |
| pytest (Python) | Python test runner |
| go test (Go) | Go test runner |
| google/go-cmp (Go) | Deep equality assertions |
