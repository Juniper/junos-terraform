# Tasks: Schema-Driven Generic Provider

## Phase 1: Generic Runtime Library (Core)

### Task 1.1: Schema Loader
**File**: `terraform_provider/generic/schema_loader.go`
**Description**: Extend existing `patch.UnmarshalTrimmedSchemaIndex` to also preserve the schema tree structure (not just the flat path→NodeInfo index). Create a `SchemaIndex` type that holds both the tree and the flat index.
**Acceptance**:
- `LoadSchema(jsonBytes)` returns tree + flat index
- Reuses `patch.SchemaNode` and `patch.NodeInfo` types (import, don't duplicate)
- Unit test with a small representative YANG fragment

### Task 1.2: Terraform Schema Builder
**File**: `terraform_provider/generic/tf_schema_builder.go`
**Description**: Implement `BuildTerraformSchema(idx *SchemaIndex) schema.Schema` that walks the schema tree and produces the Terraform Framework schema dynamically.
**Mapping**:
- leaf → `schema.StringAttribute{Optional: true}`
- leaf-list → `schema.ListAttribute{ElementType: types.StringType, Optional: true}`
- container → `schema.SingleNestedAttribute{Attributes: ..., Optional: true}`
- list → `schema.ListNestedAttribute{NestedObject: ..., Optional: true}`
- Always adds `resource_name` as Required + RequiresReplace
**Acceptance**:
- Given a known schema fragment, produces expected `schema.Schema`
- Handles deeply nested structures (5+ levels)
- Handles empty containers (no children)
- Unit test comparing output against manually-defined expected schema

### Task 1.3: Model → XML Conversion
**File**: `terraform_provider/generic/model_convert.go`
**Description**: Implement `ModelToXMLBytes(ctx, plan, idx)` that converts a Terraform plan (accessed via `tfsdk` attributes) into XML bytes suitable for NETCONF.
**Approach**:
- Walk the schema tree
- For each node, extract the corresponding attribute value from the plan object
- Build XML element tree (can reuse `patch.Node` and serialize)
- Handle null/unknown values (skip them)
**Acceptance**:
- Round-trips: model → XML → model produces same values
- Null attributes are omitted from XML
- Leaf-lists produce repeated XML elements
- Lists produce repeated XML elements with key children

### Task 1.4: XML → Model Conversion
**File**: `terraform_provider/generic/model_convert.go`
**Description**: Implement `XMLBytesToModel(ctx, xmlBytes, idx)` that converts device XML config into Terraform state attribute values.
**Approach**:
- Parse XML with `patch.BuildTree`
- Walk schema tree + XML tree in parallel
- Build `attr.Value` objects (`types.String`, `types.List`, `types.Object`)
- Handle missing XML elements → `types.StringNull()` / `types.ListNull(...)`
**Acceptance**:
- Handles all YANG node types (leaf, leaf-list, container, list)
- Missing elements produce null types (not errors)
- Handles YANG `empty` type (presence-only elements)
- Unit test with representative XML fragments

### Task 1.5: XML Tree Serialization
**File**: `terraform_provider/generic/xml_tree.go`
**Description**: Implement `SerializeTree(root *patch.Node) ([]byte, error)` to convert a generic Node tree back to XML bytes (inverse of `patch.BuildTree`).
**Acceptance**:
- Round-trip: `BuildTree(SerializeTree(BuildTree(xml))) == BuildTree(xml)`
- Produces well-formed XML with declaration header
- Handles empty elements, text content, nested children

## Phase 2: Generic Resource Implementation

### Task 2.1: Generic Config Resource
**File**: `terraform_provider/generic/resource.go`
**Description**: Implement a `configResource` that uses the schema index for all operations. This replaces the generated `configResource` struct and its CRUD methods.
**Methods**:
- `Schema()` → delegates to `BuildTerraformSchema`
- `Create()` → plan → XML → SendDirectTransaction → commit → read back → state
- `Read()` → GetConfigXML → XMLBytesToModel → set state
- `Update()` → plan XML vs state XML → patch engine → apply diff → read back
- `Delete()` → state XML vs empty → patch engine → apply delete patch
**Acceptance**:
- All CRUD paths work with mock NETCONF client
- Update uses patch engine (same behavior as current generated code)
- Delete uses patch engine for targeted removal
- Integrates with existing `ProviderConfig` and NETCONF client interface

### Task 2.2: NETCONF Client Extension
**File**: `terraform_provider/netconf/client.go` (or adapter in generic/)
**Description**: Add a `GetConfigXML() ([]byte, error)` method that returns raw XML bytes from the device, as an alternative to `MarshalConfig(v interface{})`.
**Acceptance**:
- Returns raw `<configuration>...</configuration>` XML from device
- Works with existing NETCONF session management
- Falls back gracefully (can wrap existing `get-config` RPC)

## Phase 3: Integration & Wiring

### Task 3.1: New Jinja2 Template (Minimal)
**File**: `junosterraform/templates/resource_config_provider_generic.go.j2`
**Description**: Create a minimal template that just:
- Declares `package main`
- Embeds `TrimmedSchemaJSON` as a Go raw string literal
- Imports and wires the `generic` package's resource
**Generated output**: ~50 lines regardless of schema size
**Acceptance**:
- Template renders successfully for any schema
- Output compiles with `go build`
- `go vet` passes

### Task 3.2: Update jtaf-provider Script
**File**: `junosterraform/jtaf-provider`
**Description**: Add `--mode` flag (`generic` | `legacy`, default `generic`). In `generic` mode:
- Copy `generic/` package into output directory
- Render minimal template instead of full expansion
- Still embed `trimmed_schema.json`
**Acceptance**:
- `--mode legacy` produces same output as before (no regression)
- `--mode generic` produces small provider that compiles
- Default mode is `generic`

### Task 3.3: Provider.go Template Update
**File**: `junosterraform/templates/provider.go.j2`
**Description**: Update to reference the generic resource factory when in generic mode. The `Resources()` method returns `generic.NewConfigResource(schemaJSON)`.
**Acceptance**:
- Provider correctly initializes generic resource
- Schema JSON is passed from embedded constant to resource constructor

## Phase 4: Testing & Validation

### Task 4.1: Unit Tests for Generic Library
**Files**: `terraform_provider/generic/*_test.go`
**Description**: Comprehensive unit tests for each component:
- Schema loading with various YANG structures
- TF schema building (verify attribute types/nesting)
- Model ↔ XML round-trips
- Edge cases: empty containers, leaf-lists, deeply nested structures
**Acceptance**:
- >90% coverage of generic/ package
- Tests pass with `go test ./generic/...`

### Task 4.2: Integration Test
**File**: `terraform_provider/generic/integration_test.go`
**Description**: End-to-end test using mock NETCONF server:
- Load a representative schema fragment
- Execute Create → Read → Update → Delete cycle
- Verify XML sent to device matches expectations
- Verify Terraform state matches expectations
**Acceptance**:
- Full CRUD cycle passes
- Patch-based update produces minimal NETCONF diff
- Delete removes only managed configuration

### Task 4.3: Regression Comparison
**Description**: Generate a provider using both old (legacy) and new (generic) approaches with the same XML-filtered schema. Run `terraform plan` against both and verify identical behavior.
**Acceptance**:
- Same `.tf` input produces same plan output
- Same NETCONF RPCs sent to device
- State files are semantically equivalent

## Phase 5: Cleanup (Optional, Post-Validation)

### Task 5.1: Remove Legacy Template
**Description**: Once generic mode is validated, remove or deprecate the old `resource_config_provider.go.j2` mega-template and update documentation.
**Acceptance**:
- Old template removed or moved to `legacy/`
- README updated with new workflow
- CI passes with generic mode as default

---

## Dependency Graph

```
1.1 (Schema Loader)
 ├── 1.2 (TF Schema Builder)
 ├── 1.3 (Model → XML)
 │    └── 1.5 (XML Serialization)
 └── 1.4 (XML → Model)

2.2 (NETCONF Extension)
 └── 2.1 (Generic Resource) ← depends on 1.2, 1.3, 1.4, 1.5

3.1 (New Template) ← depends on 2.1
3.2 (jtaf-provider update) ← depends on 3.1
3.3 (provider.go update) ← depends on 2.1

4.1 (Unit Tests) ← parallel with Phase 1-2
4.2 (Integration Test) ← depends on 2.1
4.3 (Regression) ← depends on 3.2
```

## Estimated Size

| Component | Lines of Go |
|-----------|-------------|
| schema_loader.go | ~100 |
| tf_schema_builder.go | ~200 |
| model_convert.go | ~350 |
| xml_tree.go | ~100 |
| resource.go | ~250 |
| Tests | ~500 |
| **Total new code** | **~1500** |
| **Generated output** | **~50 lines + JSON** |
