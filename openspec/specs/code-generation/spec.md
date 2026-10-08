# Code Generation Specification

Turning a YANG-derived JSON schema into a buildable Go Terraform provider directory.

## Purpose

Define what provider generation produces for a given model and device type, so a generated directory builds and carries the schema downstream tools need.

## Architecture

```
JSON schema (trimmed to XML with -x, or the whole model with --generic)
    ↓
jtaf-provider
    ↓
  - Copy terraform_provider's Go sources (minus tests) into
    terraform-provider-junos-{type}/
  - Write main.go and embed_schema.go
  - Write trimmed_schema.json.gz, and compile it to schema.bin.gz
    ↓
Result: Complete, buildable Go provider
```

---

## Generated Files

| File | Origin | Purpose |
|------|--------|---------|
| `generic/`, `patch/`, `netconf/`, `cmd/` | copied from `terraform_provider` | Provider source, identical for every device type |
| `main.go` | written by `jtaf-provider` | Calls `generic.Serve` with the provider type name |
| `embed_schema.go` | written by `jtaf-provider` | `go:embed` of `schema.bin.gz` |
| `schema.bin.gz` | `cmd/compileschema` | The schema as `patch.Schema`'s binary form, read without parsing JSON |
| `trimmed_schema.json.gz` | written by `jtaf-provider` | The same schema as compact JSON, for downstream tools |

The copied sources' own top-level `.go` files are removed; `main.go` and `embed_schema.go` are the generated provider's whole main package.

No Go source is rendered from a template. `ansible.j2` is the Jinja2 template behind `jtaf-ansible` and is unrelated to provider generation.

`go.mod` and `go.sum` are not generated: they are `terraform_provider`'s own, copied with the sources, and `ensure_go_module_name()` sets the module path. They match what the provider module is built and tested with, which Renovate keeps up to date.

---

## Behaviors

### Schema Scope

- **Given** `-x` names XML configuration files, **When** `jtaf-provider` runs, **Then** the embedded schema is trimmed to the paths those files use
- **Given** `--generic` is passed instead, **When** `jtaf-provider` runs, **Then** the whole model is embedded
- **Given** neither is passed, **When** `jtaf-provider` runs, **Then** it exits with an error: the scope is always explicit, because a full model costs far more memory at plan time

### Provider Source

- **Given** any device type and model, **When** generation completes, **Then** the Go sources are the same apart from the module name in `go.mod` and the provider type name in `main.go`
- **Given** `_test.go` files exist in `terraform_provider`, **When** the sources are copied, **Then** they are left out of the generated directory

---

## Post-Generation Steps

### Module Name Normalization

- **Given** the sources are copied, **When** `ensure_go_module_name()` runs, **Then** `go.mod` directive set to `module terraform-provider-junos-{type}`

### Import Path Rewriting

- **Given** source files contain `"terraform_provider/"` imports, **When** `rewrite_import_prefixes()` runs, **Then** all occurrences replaced with `"terraform-provider-junos-{type}/"`

## Requirements

### Requirement: Generated providers carry a compressed schema
Provider generation SHALL write the schema as compact gzip-compressed JSON named `trimmed_schema.json.gz` and SHALL
NOT write a plain `trimmed_schema.json`, whether the schema was trimmed to XML or left untrimmed.

#### Scenario: Schema is available to downstream tools
- **WHEN** provider generation completes
- **THEN** `trimmed_schema.json.gz` contains the schema and `jtaf-xml2tf` can consume it to produce the same output as
  the equivalent plain JSON schema

### Requirement: Generic provider generation
`jtaf-provider --generic` SHALL write the provider's main package (`main.go` calling `generic.Serve`, `embed_schema.go`), and the schema compiled into `patch.Schema`'s binary form and gzipped, which the provider embeds.

#### Scenario: Excluding subtrees
- **WHEN** `--exclude PATH` is given (repeatable) to `jtaf-provider`, or to `jtaf-yang2go`, which passes it on
- **THEN** the subtree at PATH, relative to `configuration`, SHALL be left out of the schema; a PATH that does not exist SHALL be an error

#### Scenario: Top-level version
- **WHEN** the schema has a top-level `version` leaf (the Junos release the configuration was committed with)
- **THEN** it SHALL be left out, whether or not the schema is trimmed; nested leaves named `version` SHALL be kept

### Requirement: One provider generator
Provider generation SHALL produce a single schema-driven Go provider whose source does not depend on the contents of
the YANG model. No Go source SHALL be rendered from a template. The generated output SHALL differ between device types
only in the module name, the provider type name and the embedded schema.

#### Scenario: No templated Go source
- **WHEN** a provider is generated for any device type
- **THEN** the output directory contains no file produced by template expansion of schema contents, and its size does
  not grow with the number of nodes in the model

#### Scenario: Equivalent providers across models
- **WHEN** providers are generated for two different models
- **THEN** their Go sources are identical apart from module name, provider type name and embedded schema

### Requirement: Generated provider carries the group decision
The embedded schema SHALL be the sole record of whether groups are managed. A provider built with groups SHALL expose
them, and a provider built without groups SHALL NOT, with no runtime flag, environment variable or provider-block
setting able to change that.

#### Scenario: Behaviour follows the binary
- **WHEN** a provider built without groups is used with a `.tf` file that declares a `groups` block
- **THEN** Terraform reports an unsupported attribute, regardless of provider configuration

## Output Structure

**Given** device type `vmx-4-topo`, **When** generation completes, **Then** output directory contains:

```
terraform-provider-junos-vmx-4-topo/
├── main.go                      ← Provider entry point (calls generic.Serve)
├── embed_schema.go              ← go:embed of schema.bin.gz
├── schema.bin.gz                ← Compiled schema the provider serves
├── go.mod                       ← Module: terraform-provider-junos-vmx-4-topo
├── trimmed_schema.json.gz       ← Same schema as JSON, for downstream tools (jtaf-xml2tf)
├── generic/                     ← Schema-driven provider (copied from terraform_provider/generic/)
├── patch/                       ← Patch engine package (copied from terraform_provider/patch/)
├── netconf/                     ← NETCONF client package (copied from terraform_provider/netconf/)
└── cmd/                         ← compileschema (copied from terraform_provider/cmd/)
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
test -f terraform-provider-junos-test-qfx/trimmed_schema.json.gz && echo "OK"

# Verify import paths are correct (no leftover terraform_provider/ imports)
grep -r '"terraform_provider/' terraform-provider-junos-test-qfx/ && echo "FAIL: leftover imports" || echo "OK"
```

### What to Check After Generation

| Check | Command | Expected |
|-------|---------|----------|
| Compiles | `go build .` | Exit 0, no errors |
| Tests pass | `go test ./...` | All tests pass |
| No leftover imports | `grep -r '"terraform_provider/'` | No matches |
| Schema exists | `test -f trimmed_schema.json.gz` | File exists |
| Module name correct | `head -1 go.mod` | `module terraform-provider-junos-{type}` |
