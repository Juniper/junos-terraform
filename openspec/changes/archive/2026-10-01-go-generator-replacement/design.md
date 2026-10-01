## Context

The current provider generation pipeline uses a 927-line Jinja2 template to produce a Go source file with per-node typed structs, schema definitions, and CRUD functions. This works for XML-filtered providers (~5K lines for a small EVPN deployment) but fails catastrophically for the full YANG model (554MB, 13.4M lines, uncompilable).

The existing `patch/` package already implements a runtime schema interpreter: `ProcessSchema` builds `map[path]NodeInfo` from `trimmed_schema.json`, `LeafMapWithSchema` converts XML ↔ flat path→value maps, and `ComputeDiff`/`CreateDiffPatch` produce NETCONF edit-config patches. The Update/Delete CRUD methods already use this generic path — only Create/Read and schema registration still depend on per-schema generated code.

The Terraform Plugin Framework v1.19.0 supports building `schema.Schema` programmatically from `map[string]schema.Attribute` at startup — no typed structs or compile-time definitions needed.

## Goals / Non-Goals

**Goals:**
- Single fixed Go program that works for any Junos YANG version/platform
- "Build once, update on new release" — only regenerate when YANG models change
- Compile time: seconds (fixed ~800 lines) instead of impossible (13M lines)
- Preserve existing patch engine, NETCONF client, and provider connection logic
- Maintain identical Terraform UX: nested `.tf` attributes, proper `terraform plan` diffs

**Non-Goals:**
- Removing pyang or the JSON schema pipeline (still needed to produce `trimmed_schema.json`)
- Supporting multiple Terraform resources per provider (keep single `config` resource)
- Schema validation beyond what Terraform Framework provides (no custom validators in v1)
- Backward compatibility with existing generated providers (this is a new architecture)

## Decisions

### 1. Schema representation: `ListNestedAttribute` tree (not `DynamicAttribute`)

**Choice:** Build fully-typed nested `schema.Schema` from JSON at startup.

**Alternatives considered:**
- `DynamicAttribute` with raw XML/JSON string — trivial to implement but terrible UX (no plan diffs on fields, no IDE autocomplete)
- `ObjectAttribute` with `AttributeTypes` map — doesn't support nesting beyond one level
- `SingleNestedAttribute` — doesn't model YANG lists (which can have multiple instances)

**Rationale:** `ListNestedAttribute` with recursive `NestedAttributeObject{Attributes: ...}` gives users the same experience as the current generated provider — every YANG leaf appears as a named field in `.tf` files.

### 2. Data bridge: Generic `tftypes.Value` tree walking (not Go structs)

**Choice:** Use `terraform-plugin-go`'s `tftypes.Value` for plan/state manipulation instead of typed Go structs.

**Alternatives considered:**
- Generate Go structs at compile time (current approach) — doesn't scale
- Use `map[string]interface{}` with reflection — fragile, poor type safety at TF layer

**Rationale:** The framework's `req.Plan.GetAttribute()` / `resp.State.SetAttribute()` can work with `tftypes.Value` directly. Combined with the schema index, we walk the plan tree generically: leaf → string value → XML element, container/list → recurse into children.

### 3. JSON embedding: `go:embed` with gzip compression

**Choice:** Embed gzipped `trimmed_schema.json` in the binary, decompress once at startup.

**Alternatives considered:**
- Raw `go:embed` — works but 268MB binary for full QFX model
- External file loaded at runtime — adds deployment complexity (path configuration)
- Pre-compiled protobuf index — faster startup but adds build complexity

**Rationale:** JSON compresses ~90% (repetitive structure). ~15-30MB binary is acceptable. One-time decompression at startup (~2-3s) is fine since provider processes are long-lived.

### 4. CRUD implementation: Reuse patch engine for all operations

**Choice:** All CRUD operations use the same flow: plan→XML via tree walk, then existing patch engine for diff/apply.

```
Create: plan→XML → load-configuration (full replace)
Read:   get-config → XML→state via tree walk
Update: plan→XML vs state→XML → ComputeDiff → CreateDiffPatch → edit-config
Delete: state→XML vs empty → ComputeDiff → CreateDiffPatch (delete ops)
```

**Rationale:** The Update and Delete paths already work this way in the current provider. Create and Read just need the generic XML↔TF bridge instead of generated struct conversion.

### 5. Architecture: New `generic/` package alongside existing code

**Choice:** Create `terraform_provider/generic/` package with the new components. Keep existing code as-is during development.

```
terraform_provider/
├── generic/
│   ├── schema_builder.go      (~150 lines)
│   ├── plan_xml_bridge.go     (~200 lines)
│   ├── xml_state_bridge.go    (~200 lines)
│   ├── resource.go            (~150 lines)
│   └── embed.go               (~100 lines)
├── patch/                     (unchanged)
├── netconf/                   (unchanged)
├── provider.go                (minor: register generic resource)
└── main.go                    (unchanged)
```

## Risks / Trade-offs

**[Large schema startup time]** → Mitigation: Parsing 268MB JSON takes ~2-3s. Provider processes are long-lived (one per `terraform apply`), so this is acceptable. Can optimize later with binary serialization if needed.

**[Terraform schema depth limits]** → Mitigation: TF Framework has no documented nesting depth limit. The full Junos model has ~15 levels max. Need to verify with a spike that deeply nested `ListNestedAttribute` works.

**[`tftypes.Value` API complexity]** → Mitigation: The `terraform-plugin-go` API for building/reading `tftypes.Value` trees is low-level but well-documented. The pattern is mechanical (recursion on node type) and testable in isolation.

**[Binary distribution size]** → Mitigation: Even gzipped, full model is ~15-30MB. For targeted deployments, users can still use `-x` filtering to produce smaller schemas. Terraform Registry has no hard binary size limits.

**[State migration from existing providers]** → Mitigation: Non-goal for v1. Users upgrading from generated providers would need to `terraform state rm` and re-import. Document this clearly.
