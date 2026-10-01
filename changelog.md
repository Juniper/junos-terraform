
# Change Log
All notable changes to this project will be documented in this file.
 
The format is based on [Keep a Changelog](http://keepachangelog.com/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Changed
- **BREAKING:** `jtaf-provider` no longer renders Go source from Jinja2 templates. Both invocations now generate the
  same schema-driven provider and differ only in the model it embeds: trimmed to the XML given with `-x`, or untrimmed
  with `--generic`. `junosterraform/templates/resource_config_provider.go.j2`, `provider.go.j2` and `config.go.j2` are
  removed. The resource type, provider block and attribute names are unchanged, so existing `.tf` files and Terraform
  state keep working; rebuild the provider binary to pick this up.
- **BREAKING:** `jtaf-provider` and `jtaf-ansible` now write the trimmed schema as compact, gzipped `trimmed_schema.json.gz`; the plain `trimmed_schema.json` is no longer written. Pass the `.gz` file to `jtaf-xml2tf -j` / `jtaf-xml2yaml -j`.
- `jtaf-provider`, `jtaf-ansible`, `jtaf-xml2tf` and `jtaf-xml2yaml` read the schema given with `-j` (file or `-` for stdin) as plain or gzipped JSON, detected by content, through a shared loader (`jtaf_common.load_schema_json`); directories generated before this change still load.

## [1.2.0] - 2026-06-17

### Added
- NETCONF patch engine (`terraform_provider/patch/`) — computes minimal leaf-level diffs and emits targeted `edit-config` operations instead of full configuration replacement
- Schema-aware leaf map flattening using YANG schema (`trimmed_schema.json`) for correct list key detection, leaf type handling, and ordered leaf-list support
- XML element order stabilization (`AlignXMLOrderToReference`) — eliminates phantom diffs caused by non-deterministic Junos `get-config` element ordering
- UTF-8 double-encoding repair (`NormalizeLeafMapUTF8`) — fixes false diffs caused by multi-byte characters (e.g., em-dash) being double-encoded through Go XML marshal/unmarshal cycles
- Container delete coalescing — replaces N individual leaf deletes with a single container-level `nc:operation="delete"` when all children are being removed
- Ordered leaf-list position tracking — detects reordering of `ordered-by user` leaf-lists (e.g., VRRP virtual-address) using positional keys `[pos=N]`
- Fallback safety net — if patch verification detects residual differences, emits a Terraform warning and falls back to full replace
- `SendUpdate` and `SendDirectTransaction` NETCONF client methods for `edit-config` and group-less `load-configuration` RPCs
- 57+ new Go unit tests covering all YANG node kinds, CRUD operation permutations, and 10 validated corner cases

### Changed
- Update operation now uses minimal diff/patch (`edit-config` with per-element `nc:operation`) instead of full group replacement (`load-configuration`)
- Delete operation now uses targeted per-leaf/container deletes instead of group deletion
- Configuration model changed from apply-groups wrapping to direct base configuration
- Read operation now stabilizes XML element ordering to prevent spurious diffs

### Fixed
- Phantom diffs from XML element order variation across `get-config` reads
- UTF-8 encoding drift causing false `terraform plan` changes on strings with special characters
- Silent fallback to full replace now emits Terraform warning diagnostics
- Ordered leaf-list reorder detection (previously invisible due to set semantics)
- Container deletion generating excessive individual leaf delete operations

## [1.1.0] - 2025-10-03
 
### Added
- jtaf-yang2go command to combine yang files to JSON conversion and provider creation ([#6](https://github.com/Juniper/junos-terraform/issues/6))

### Changed
- jtaf-provider, jtaf-xml2tf, and jtaf-yang2go to accept multiple xml configurations of the same device type ([#72](https://github.com/Juniper/junos-terraform/issues/72))
- jtaf-xml2tf to support a base configuration, groups, and apply groups ([#65](https://github.com/Juniper/junos-terraform/issues/65))
 
### Fixed
- Dependency on private go-netconf repository ([#61](https://github.com/Juniper/junos-terraform/issues/61))
- Unexpected output rpc-reply messages ([#71](https://github.com/Juniper/junos-terraform/issues/71))
- Leaf-list error in jtaf-xml2tf ([#65](https://github.com/Juniper/junos-terraform/issues/65))
 
## [1.0.0] - 2025-07-01

### Added
- Many updates to make JTAF production ready ([Release 1.0.0](https://github.com/Juniper/junos-terraform/releases/tag/1.0.0))

## [0.1.1] - 2025-06-26

### Added
- Many updates and examples ([Release 0.1.1](https://github.com/Juniper/junos-terraform/releases/tag/0.1.1))

## [0.1] - 2021-04-14

### Added
- First release of API to generate Junos modules for Terraform ([Release 0.1](https://github.com/Juniper/junos-terraform/releases/tag/0.1))
 
