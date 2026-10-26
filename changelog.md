
# Change Log
All notable changes to this project will be documented in this file.
 
The format is based on [Keep a Changelog](http://keepachangelog.com/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Added
- `jtaf-provider --yang-compliant` keeps a choice as the model declares it, a `choice-ident` and a `choice-value`,
  instead of flattening it to an element named after the case. A device with `system services netconf yang-compliant`
  returns the modelled form, so a `--generic` provider that flattened it could not read the choice back and the device
  rejected the whole-tree write with "syntax error, expecting <choice-ident>". The schema fingerprint records the
  choice, so a mismatch with the Terraform files is warned about.
- `jtaf-provider --groups`, passed on by `jtaf-yang2go`, keeps the `groups` subtree and the `apply-groups` leaf-list so
  the provider manages Junos configuration groups as ordinary hierarchy: a group is an entry of the `groups` list keyed
  by its `name`, and `apply_groups` is an ordered list of group names. With `--generic` it warns, because advertising
  the whole model twice over needs several GB at plan time.
- `jtaf-xml2tf --groups` converts a group-based configuration as it stands instead of flattening what a device inherits
  through `apply-groups` into the base hierarchy. Use it with a provider generated with `--groups`.
- `examples/evpn-vxlan-dc-groups/`, a group-based example beside the base-configuration one, covering the same devices
  with each device's configuration in a role-named group. `examples/providers/build-groups.sh` and `convert-groups.sh`
  build and convert it; see [examples/DEMO-GROUPS.md](examples/DEMO-GROUPS.md). A CI job applies it against the NETCONF
  mock and checks idempotency, drift inside a group, and group removal.
- The provider and `jtaf-xml2tf` record the schema they were built from (its scope and the YANG modules' revision
  dates), and the provider warns when a `.tf` file was generated from a different one. Generating files from a trimmed
  schema and applying them with a full-model provider passes `terraform validate` and then fails at the device, because
  the two describe the same configuration differently. Files or providers generated before this carry no fingerprint and
  stay silent.

### Fixed
- A load or commit the device rejected was read as success. Junos nests the `rpc-error` inside
  `<load-configuration-results>` or `<commit-results>`, and `parseReply` matched only direct children of `<rpc-reply>`,
  so the caller carried on and left its staged configuration in the candidate — which is shared rather than
  per-session, so any later commit from any source would have picked it up. Errors are now collected at any depth, and
  the candidate is discarded after a failed load, edit or commit.
- `create` loaded the plan as a merge, which cannot express a removal, so configuration the device already had that the
  plan did not want survived — and because these devices ship a `unit 0` on every port and dhcp on `em0`, the commit was
  rejected outright. `create` now reconciles against the device exactly as `update` does. Both also returned the device
  read-back as the new state, which differs from the plan once Junos normalises it, so Terraform failed the apply with
  "Provider produced inconsistent result after apply"; they now return the planned value, and what the device holds
  beyond it is drift for the next refresh to report.
- A YANG `choice` is now described the same way whichever schema scope is used. A device writes a choice as an element
  naming the case (`<add/>`, `<exact/>`), which trimming already produced as a side effect of matching cases against the
  XML, while `--generic` kept the modelled `choice-ident`/`choice-value` pair. Terraform files generated against one
  scope lost the choice against the other and the device rejected the edit-config with "syntax error, expecting
  <choice-ident>".
- `jtaf-yang2go` never examined `jtaf-provider`'s exit status, so a fatal error was printed and the script still exited
  0; under `set -e`, `build.sh` carried on over a provider that was never finished and reported success. It also now
  names each XML path the YANG model does not cover, and summarises them, instead of repeating the whole set collected
  so far — XML taken from a different Junos release than the model lost configuration silently.
- `jtaf-xml2tf` now names attributes as the provider does. `SanitizeName` lower-cases a node name and `normalize_tag`
  did not, so a Junos name with capitals became a different attribute in the `.tf` than the provider advertises:
  `do-not-translate-AAAA-query-to-A-query`, and `AH_header`/`ESP_header` in firewall filters.
- `jtaf-xml2tf` escapes values that HCL reads as syntax. Values were written straight into a quoted string, so a real
  newline from `get-configuration` produced a file Terraform refused to parse, and `${` or `%{` started an
  interpolation that was evaluated rather than sent to the device. Escaping is used rather than a heredoc, which would
  append a trailing newline that Junos strips again, leaving every plan reporting a change.
- The example configurations are stored as a device returns them. The login banner held the two-character escape rather
  than real newlines, which no device produces and which cannot be read back, so every plan reported a change.
- `jtaf-provider --exclude` now reaches configuration nodes held inside YANG `choice` and `case` nodes, such as
  `vlans/vlan/vlan-id`. Those group nodes in the model but are not configuration, and the provider already flattens
  them away, so a path that a device would show was rejected as not found.

### Changed
- **BREAKING:** `jtaf-provider` and `jtaf-yang2go` require exactly one of `-x` or `--generic`. Omitting both used to
  embed the whole model silently, which made `--generic` a no-op and gave a full model to anyone who forgot the flag;
  a full model costs far more memory at plan time, so it is now opted into. Every script in the repository already
  passes one of the two.
- **BREAKING:** The `groups` subtree and the `apply-groups` leaf-list are now left out of a generated provider unless
  `--groups` is given. In a full Junos model `groups` repeats the whole configuration hierarchy and is about half of
  its nodes, and the provider never exposed it, so this removes a cost that bought nothing.
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
 
