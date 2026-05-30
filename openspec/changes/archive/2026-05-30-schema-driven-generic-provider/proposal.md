# Proposal: Schema-Driven Generic Provider

## Summary

Replace the current Jinja2 code-generation architecture (which emits a unique Go struct + methods for every YANG node path) with a **schema-driven generic runtime** that loads `trimmed_schema.json` at startup and handles all Terraform ↔ XML ↔ NETCONF conversion generically. This reduces the generated Go output from potentially hundreds of thousands of lines to ~200 lines of wiring code plus the existing JSON schema.

## Problem

When `jtaf-provider` runs against the full Junos YANG model **without** an `-x` XML config filter, the Jinja2 template (`resource_config_provider.go.j2`) generates an enormous Go file because it expands:

1. **One unique XML struct** per YANG container/list node (e.g., `xml_Interfaces_Interface_Unit_Family_Inet_Address`)
2. **One unique Terraform model struct** per node (with `AttrTypes()` + `Attributes()` methods)
3. **Unrolled nested for-loops** in `modelToConfig()` — per-path conversion code
4. **Unrolled nested for-loops** in `configToModel()` — per-path reverse conversion code

This yields ~7 code blocks per YANG node × thousands of nodes = 100K-1M+ lines of Go. The file may not even compile due to Go's practical limits on file/package size.

The underlying issue: the Jinja2 template treats every YANG path as structurally unique, when in reality:
- The CRUD operations are identical regardless of path
- The XML ↔ model conversion is purely mechanical tree walking
- The Terraform schema is a direct reflection of the YANG tree structure
- The `patch/` engine already works generically on XML trees

## Proposed Solution

Implement a **generic runtime library** in Go that:
1. Reads the YANG schema tree from `trimmed_schema.json` (already generated and embedded)
2. Builds the Terraform `schema.Schema` dynamically from the schema tree
3. Converts between Terraform state (nested `map[string]any`) and XML using generic tree walkers
4. Delegates CRUD to the existing `patch/` engine and NETCONF client

The generated provider becomes just:
- `provider.go` — wiring (already mostly generic)
- `config.go` — NETCONF client setup (already generic)
- `trimmed_schema.json` — the schema data
- A small registration file pointing to the schema

## Goals

- **100-1000x reduction** in generated Go code size
- Full Junos YANG support without XML filter requirement
- Preserve existing NETCONF patch/diff behavior
- Maintain Terraform plan/apply/destroy semantics
- Keep the `patch/` engine unchanged

## Non-Goals

- Changing the pyang plugin or JSON schema format
- Modifying the NETCONF protocol layer
- Supporting multiple resources per provider (keep single `configResource`)
- Backwards compatibility with previously generated providers (this is a new generation approach)

## Scope

- New Go package: `terraform_provider/generic/` (runtime library)
- Modified Jinja2 template: minimal wiring instead of full expansion
- Modified `jtaf-provider` script: uses new template
- Updated tests: adapt to generic runtime

## Risks

| Risk | Mitigation |
|------|-----------|
| Terraform Framework may not fully support dynamic nested schemas | Prototype `schema.Schema` construction from JSON early; fallback to `schema.DynamicAttribute` if needed |
| Performance regression from dynamic dispatch vs. static types | Benchmark; config management workloads are low-throughput (seconds per operation acceptable) |
| `encoding/xml` edge cases with generic struct | The `patch/` engine's `BuildTree` already solves this with `etree`; can use same approach |
| Loss of Go compile-time type safety | Compensate with comprehensive unit tests on schema→TF and XML→model paths |
