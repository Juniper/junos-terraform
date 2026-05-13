# Code Generation Specification

Jinja2 templates that render YANG-derived JSON schema into Go Terraform provider source code. Lives at `junosterraform/templates/`.

## Architecture

```
Filtered JSON schema (from filter_json_using_xml)
    ↓
Jinja2 template engine (Python)
    ↓
For each template:
  - Render with schema data context
  - Write to terraform-provider-junos-{type}/ directory
    ↓
Result: Complete, buildable Go provider
```

---

## Templates

| Template | Generates | Purpose |
|----------|-----------|---------|
| `resource_config_provider.go.j2` | `resource_config_provider.go` | Resource structs, CRUD methods, schema definitions |
| `provider.go.j2` | `provider.go` | Provider struct, plugin framework integration |
| `config.go.j2` | `config.go` | Config struct, NETCONF Client() factory |
| `go.mod.j2` | `go.mod` | Go module definition with correct module path |
| `ansible.j2` | Ansible task files | Ansible role structure with tasks/templates |

---

## Behaviors

### Template Variable Context

- **Given** `data.device_type` is `"vmx-4-topo"`, **When** templates render, **Then** all references use `junos-vmx-4-topo` as provider name
- **Given** `data['root']['children']` contains filtered YANG nodes, **When** `resource_config_provider.go.j2` renders, **Then** one Go struct is generated per top-level YANG container/list

### Resource Code Generation (`resource_config_provider.go.j2`)

- **Given** a YANG container with leaf children, **When** rendered, **Then**:
  - A Terraform schema attribute is generated per leaf
  - String leaves → `schema.StringAttribute`
  - Integer leaves → `schema.Int64Attribute`
  - Boolean leaves → `schema.BoolAttribute`
  - Empty type leaves → `schema.BoolAttribute` (presence semantics)

- **Given** a YANG list with key field, **When** rendered, **Then**:
  - List key becomes a required attribute
  - Key change forces resource replacement

- **Given** a YANG container with nested children, **When** rendered, **Then**:
  - Nested struct type generated
  - Parent attribute uses `schema.ListNestedAttribute` or `schema.SingleNestedAttribute`

### Create Method Generation

- **Given** the resource is being created, **When** rendered Create() runs, **Then**:
  1. Extract all attribute values from Terraform plan
  2. Build XML representation of the configuration
  3. Wrap in `<configuration><groups><name>{resource_name}</name>...</groups></configuration>`
  4. Call `client.SendRPC()` with `<load-configuration>` containing the XML
  5. Call `client.SendRPC()` with `<commit/>`

### Read Method Generation

- **Given** the resource needs state refresh, **When** rendered Read() runs, **Then**:
  1. Call `client.SendRPC()` with `<get-configuration>` for the group
  2. Parse returned XML
  3. Map XML values back to Terraform state attributes
  4. If group not found on device, mark resource as gone

### Update Method Generation

- **Given** the resource has changes, **When** rendered Update() runs, **Then**:
  1. Build current state XML and desired plan XML
  2. Use patch engine: `LeafMapWithSchema` → `ComputeDiff` → `CreateDiffPatchWithSchema`
  3. Send `<edit-config>` with the patch
  4. Commit changes

### Delete Method Generation

- **Given** the resource is being destroyed, **When** rendered Delete() runs, **Then**:
  1. Build delete XML: `<edit-config><configuration><groups><name>{name}</name></groups></configuration></edit-config>` with `operation="delete"`
  2. Commit the deletion

---

## Post-Generation Steps

### Module Name Normalization

- **Given** templates are rendered, **When** `ensure_go_module_name()` runs, **Then** `go.mod` directive set to `module terraform-provider-junos-{type}`

### Import Path Rewriting

- **Given** source files contain `"terraform_provider/"` imports, **When** `rewrite_import_prefixes()` runs, **Then** all occurrences replaced with `"terraform-provider-junos-{type}/"`

### Schema Emission

- **Given** code generation completes, **When** finalization runs, **Then** `trimmed_schema.json` is written to the output directory
- This JSON is loaded at runtime by `ProcessSchema()` in the patch engine for type-aware diffing

---

## Output Structure

**Given** device type `vmx-4-topo`, **When** generation completes, **Then** output directory contains:

```
terraform-provider-junos-vmx-4-topo/
├── main.go                      ← Provider entry point (starts plugin server)
├── provider.go                  ← Provider definition (schema, configure)
├── config.go                    ← NETCONF config struct + Client() factory
├── resource_config_provider.go  ← Generated resource with CRUD
├── go.mod                       ← Module: terraform-provider-junos-vmx-4-topo
├── trimmed_schema.json          ← Schema index for patch engine
├── patch/                       ← Patch engine package (copied from terraform_provider/patch/)
└── netconf/                     ← NETCONF client package (copied from terraform_provider/netconf/)
```

---

## Conventions

- **Given** a generated Go file, **When** inspected, **Then** it follows `gofmt` formatting
- **Given** provider registry, **When** provider built, **Then** registry address is `tf-registry.click/juniper/jtaf-{type}`
- **Given** resource naming, **When** resource created, **Then** type name is `terraform-provider-junos-{type}`
- **Given** YANG namespaces, **When** generating schema, **Then** namespaces are stripped (Terraform users don't see XML namespaces)

---

## Testing

### How to Validate Generated Code

```bash
# Generate a provider
jtaf-yang2go -p examples/yang/18.2/18.2R3/common examples/yang/18.2/18.2R3/junos-qfx/conf/*.yang \
  -x examples/evpn-vxlan-dc/dc1/dc1-spine1.xml -t test-qfx

# Verify it compiles
cd terraform-provider-junos-test-qfx && go build .

# Run generated provider tests
cd terraform-provider-junos-test-qfx && go test ./...

# Verify schema file exists
test -f terraform-provider-junos-test-qfx/trimmed_schema.json && echo "OK"

# Verify import paths are correct (no leftover terraform_provider/ imports)
grep -r '"terraform_provider/' terraform-provider-junos-test-qfx/ && echo "FAIL: leftover imports" || echo "OK"
```

### What to Check After Generation

| Check | Command | Expected |
|-------|---------|----------|
| Compiles | `go build .` | Exit 0, no errors |
| Tests pass | `go test ./...` | All tests pass |
| No leftover imports | `grep -r '"terraform_provider/'` | No matches |
| Schema exists | `test -f trimmed_schema.json` | File exists |
| Module name correct | `head -1 go.mod` | `module terraform-provider-junos-{type}` |
