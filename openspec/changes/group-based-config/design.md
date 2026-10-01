# Design

## Context

See proposal.md — Why. Three facts about the current code shape the approach:

- `jtaf-provider` has two generation paths. The default one renders `resource_config_provider.go.j2`, `provider.go.j2`
  and `config.go.j2` into Go source; `--generic` instead copies `terraform_provider/` as-is, writes `embed_schema.go`
  and compiles the schema with `terraform_provider/cmd/compileschema`. Both already write `trimmed_schema.json.gz`.
- Groups are suppressed in three independent places: `jtaf_common.filter_json_using_xml()` strips `apply-groups` from
  the XML before computing paths, the Jinja2 template rejects `groups`/`apply-groups` nodes, and
  `terraform_provider/generic/names.go` skips those two names when listing resource attributes.
- `terraform_provider/netconf/client.go` still carries the pre-1.2 group model (`groupStrXML`, `getGroupXMLStr`,
  `applyGroupsList`, `sendApplyGroupsLocked`), where the provider invented an `apply-groups` list from the resources it
  wrote. Nothing in the current provider calls it.

The `groups` subtree in a YANG model is a near-copy of the whole configuration model, which is why including it is a
size and performance question, not just a correctness one.

## Goals / Non-Goals

**Goals:**

- One Go provider implementation, selected by nothing; the embedded schema is the only variable.
- Groups handled as data: no code path in the provider knows the word "groups".
- Default builds unchanged in output and in cost.

**Non-Goals:**

- Ansible. `jtaf-ansible`, `jtaf-yang2ansible` and the Ansible templates keep their current group handling.
- Junos group inheritance semantics. The provider manages the literal `groups` and `apply-groups` configuration; it
  does not resolve, validate or display what a group contributes to the inherited configuration.
- `groups` matched by wildcards (`<name>ge-*</name>` style interface wildcards inside a group) are carried through as
  literal key values; no wildcard expansion is attempted.
- Terraform import of existing groups.

## Decisions

### Retire the Jinja2 Go path rather than teach it groups

The schema-driven provider already covers everything the generated one does, and the generated one cannot hold a full
model at all. Keeping both means every feature is written twice and the groups subtree would push the generated source
past what the Go compiler accepts. So the default path becomes: filter the schema with `-x` as it does today, then run
the same generation as `--generic`. `--generic` degrades to "do not filter".

*Alternative considered:* keep the Jinja2 path and gate groups to `--generic`. Rejected — it leaves the two
implementations diverging on a visible behaviour, and the trimmed build is the one most users run.

*Consequence:* `resource_config_provider.go.j2`, `provider.go.j2` and `config.go.j2` are deleted, along with
`ensure_go_module_name`/`rewrite_import_prefixes` usage that only the template path needed. The generated provider's
resource type name, provider block and attribute names do not change, so existing `.tf` files and state files load
unchanged; users only have to rebuild the provider binary.

### Groups are a schema decision, not a provider decision

`--groups` controls what goes into the schema; the provider then has no special case. Concretely the default build
behaves as if `--exclude groups` were passed and additionally drops the top-level `apply-groups` leaf-list, which is
how `version` is already handled by `drop_version()`. `names.go` loses its name-based skip, so attribute listing
becomes purely structural.

*Alternative considered:* a runtime flag or provider-block setting. Rejected — it would make two providers with the
same binary and schema behave differently, and the spec requires behaviour to follow from the artifact.

*Consequence:* `--groups` together with `--exclude groups` is contradictory and is rejected at generation time.

### Trimming is what makes groups affordable

In the `-x` path the group body is trimmed to the paths the example XML actually uses, so a group-based build is about
the same size as the equivalent base-hierarchy build. `--generic --groups` embeds the group subtree in full and roughly
doubles the model; it is supported but documented as the expensive combination, and `--exclude` remains the lever for
trimming it.

For trimming to see inside groups, `filter_json_using_xml()` must stop removing `apply-groups` and must let `<groups>`
children contribute paths. The XPaths it derives from `<groups><interfaces>…` already line up with the schema path
`configuration/groups/interfaces/…`, so no path rewriting is needed — only the removal of the current suppression.

### `apply-groups` is an ordered leaf-list, nothing more

Order matters to Junos, and the patch engine already tracks `ordered-by user` leaf-lists (`patch/order.go`, `[pos=N]`
keys in the leaf map). `apply-groups` is therefore handled by marking it ordered in the schema and letting the existing
machinery run. The legacy bookkeeping in `netconf/client.go` is deleted rather than reused: it synthesised references
the user never asked for, which the spec now forbids.

*Risk of the deletion:* those helpers are exported. They are unused inside the repo, and the package is an internal
support package of the generated provider, so removal is treated as part of the same breaking change as the generator
retirement.

### Read back without inheritance

A device with `apply-groups` reports inherited values in the base hierarchy when configuration is fetched with
inheritance resolved. If the provider read that, every group-based configuration would show permanent drift: values it
never declared would appear in the base hierarchy. The read path therefore must continue to request the raw committed
configuration (no `inherit` attribute on `get-configuration`), and this is asserted by a test rather than left implicit.

### The example is a transformation of the existing one

`examples/evpn-vxlan-dc-groups/` mirrors `examples/evpn-vxlan-dc/` device for device. Per device, configuration that is
not device-identity (interfaces, protocols, policy, routing instances) moves inside a single group named after the
device role, and the base hierarchy keeps only `system host-name`, management addressing, and the `apply-groups`
reference. Keeping the device set identical means the two examples can be compared directly, and the group build can
reuse the same YANG inputs and the same build/convert script shape under `examples/providers/`.

*Alternative considered:* a small purpose-built example with two interfaces. Rejected — it would not exercise trimming
of a realistic group body, which is the part most likely to break.

## Risks / Trade-offs

- **Retiring the Jinja2 path breaks anyone pinned to generated source** → The resource and provider surface is
  unchanged, so only a rebuild is needed; the change is called out as BREAKING in the proposal and the changelog, and
  the old generator stays reachable through git history.
- **`--generic --groups` roughly doubles an already large model**, raising provider start-up memory and plan time →
  Documented as the expensive combination; `--exclude` trims the group subtree; the shipped example uses the trimmed
  path.
- **Inherited configuration read back as drift** → Addressed by the no-inheritance read decision above, with a test
  that applies a group-based configuration and asserts an empty second plan.
- **Group deletion is destructive on the device**: deleting a group removes configuration from every device that
  applies it → The diff scopes deletes to the declared group and emits a single container delete; the e2e example
  covers removing a group and asserts the base hierarchy is untouched.
- **Wildcard group keys** (`<name>ge-*</name>`) may not survive the key-escaping used in leaf-map paths → Covered by a
  patch-engine test with a wildcard key; if escaping proves lossy the key encoding is fixed in the patch engine rather
  than special-cased for groups.
- **The mock NETCONF server models groups as whole-blob state** (`running_groups` keyed by group name), which does not
  match a patch-based edit inside a group → The mock needs to apply edits to a single configuration tree and let
  `groups` be part of it; otherwise the e2e tests would pass against behaviour no device has.

## Migration Plan

1. Land the generator unification first, with `--groups` absent. Generated output for existing examples must be
   functionally identical; the both-providers e2e script reduces to a single build and is updated accordingly.
2. Add `--groups` and the schema-side changes; default builds stay byte-identical.
3. Add the patch-engine and read-path work, then the example and its e2e run.

Rollback is per step: each step leaves the default build working, so reverting the last commit restores a usable
pipeline. Users roll back by rebuilding from the previous tag; no state migration is involved in either direction.
