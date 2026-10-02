## Why

The current Jinja2-based code generator produces a per-schema Go source file that scales linearly with the YANG model size. For the full Junos model (33 QFX YANG files, 18.2), the generated `resource_config_provider.go` is 554MB / 13.4M lines with 158K type definitions — far beyond what the Go compiler can handle. This makes a "build once, covers all Junos" universal provider impossible with the current architecture.

A generic, schema-driven Go provider that reads `trimmed_schema.json` at runtime eliminates the code generation bottleneck entirely. The provider binary is compiled once from fixed source code (~800 lines) and works for any YANG model by swapping the embedded JSON.

## What Changes

- **BREAKING**: Remove dependency on Jinja2 template (`resource_config_provider.go.j2`) for provider generation
- Replace 927-line Jinja2 template with a fixed Go module (`terraform_provider/generic/`) that builds Terraform schema, plan→XML, and XML→state at runtime from the schema index
- Embed `trimmed_schema.json` (optionally gzip-compressed) via `go:embed` at compile time
- Reuse the existing `patch/` engine (6.4K lines) for diff, XML↔path-map conversion, and schema indexing — no changes needed
- Reuse existing `netconf/` client and `provider.go` connection setup unchanged
- `jtaf-provider` output changes: instead of rendering Go code, it produces `trimmed_schema.json` + copies the generic provider source

## Capabilities

### New Capabilities
- `dynamic-schema-builder`: Recursively walks `trimmed_schema.json` at startup to build `schema.Schema` with `ListNestedAttribute` for containers/lists and `StringAttribute` for leaves
- `generic-plan-xml-bridge`: Converts Terraform plan `tftypes.Value` tree → Junos XML and device XML → Terraform state `tftypes.Value` tree using the schema index
- `embedded-schema-loading`: Loads and decompresses embedded `trimmed_schema.json` at provider startup into the existing `map[path]NodeInfo` index

### Modified Capabilities
- `code-generation`: The `jtaf-provider` CLI tool no longer renders Go source via Jinja2. Instead it outputs `trimmed_schema.json` and copies the generic provider skeleton. The `-x` flag remains optional for filtering.
- `terraform-provider`: CRUD operations switch from per-schema typed structs to generic tree-walking, while maintaining identical NETCONF behavior (edit-config, get-config, commit).

## Impact

- **Go source**: New package `terraform_provider/generic/` (~800 lines across 4-5 files)
- **Jinja2 template**: `junosterraform/templates/resource_config_provider.go.j2` becomes unused (can be archived)
- **Python CLI**: `jtaf-provider` step 2 (render template) replaced with copy-generic-provider step
- **Binary size**: Provider binary grows by embedded JSON size (268MB raw → ~15-30MB gzip for full QFX 18.2)
- **Dependencies**: No new Go dependencies — uses existing `terraform-plugin-framework` v1.19.0 and `terraform-plugin-go` v0.31.0
- **Existing tests**: Patch engine tests (`patch/*_test.go`), NETCONF tests, provider_test.go remain valid. Generated-provider build tests become generic-provider build tests.
