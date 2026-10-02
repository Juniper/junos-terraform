# Proposal

## Why

Operators keep the bulk of a Junos configuration in configuration groups (`groups` plus `apply-groups`) and leave only
device-specific values in the base hierarchy. JTAF cannot express that today: `groups` and `apply-groups` are stripped
out of the schema and hard-filtered out of the generated Terraform resource, so a group-based network has to be
flattened before it can be managed. Making groups visible by default is not acceptable either — in a full Junos model
the `groups` subtree is a second copy of the entire configuration model, which is why it is excluded now — so the
support has to be opt-in.

A second, blocking problem surfaced while scoping this: there are still two Go provider generators. The default
(non-`--generic`) path renders Go source from Jinja2 templates (`resource_config_provider.go.j2`, `provider.go.j2`,
`config.go.j2`), and that template explicitly drops `groups`/`apply-groups` nodes. Adding groups to two generators is
wasted work, so this change first finishes the migration started by `go-generator-replacement`: one schema-driven Go
provider for both the trimmed and the full-model builds.

## What Changes

- **BREAKING**: Retire the Jinja2 Go code-generation path. `jtaf-provider` without `--generic` produces the same
  schema-driven provider as `--generic`, differing only in that the embedded schema is trimmed to the XML given with
  `-x`. `resource_config_provider.go.j2`, `provider.go.j2` and `config.go.j2` are removed; no Go source is templated.
  `--generic` remains accepted and means "embed the untrimmed model".
- Add `--groups` to `jtaf-provider` (and pass-through on `jtaf-yang2go`). Without it, behaviour is unchanged: `groups`
  and `apply-groups` are excluded from the schema and from the resource. With it, both are kept.
- Groups are modelled as plain Junos hierarchy, not as a provider or resource extension. There is no `group_name`
  attribute: `groups` is the YANG list keyed by `name`, so a group is written as a nested block in the `.tf` file
  exactly as it appears in the device configuration.
- `apply-groups` becomes a managed, order-preserving leaf-list attribute at the top of the configuration, so the
  provider writes group references alongside the group bodies. The legacy implicit apply-groups bookkeeping in
  `netconf/client.go` (`applyGroupsList`, `sendApplyGroupsLocked`, `getGroupXMLStr`) is removed rather than revived.
- The patch engine's leaf map, diff and patch emit paths below `groups[name=X]` and preserve `apply-groups` ordering,
  so a change inside a group produces a minimal `edit-config` instead of a full group reload.
- Add `--groups` to `jtaf-xml2tf` so it keeps the group hierarchy instead of flattening applied groups into the base
  configuration (current, and still default, behaviour).
- Add `examples/evpn-vxlan-dc-groups/`, a sibling of `examples/evpn-vxlan-dc/`, covering the same DC devices with the
  bulk of each configuration inside a group and only device-specific leaves in the base hierarchy, plus the build and
  convert scripts and generated `.tf` files needed to plan and apply it.
- Ansible (`jtaf-ansible`, `jtaf-yang2ansible`, the role templates) is explicitly out of scope and unchanged.

## Capabilities

### New Capabilities
- `config-groups`: opt-in handling of the `groups` subtree and `apply-groups` leaf-list across schema generation, the
  Terraform resource, XML→HCL conversion, and the NETCONF patch it produces.

### Modified Capabilities
- `cli-tools`: `--groups` on `jtaf-provider`, `jtaf-yang2go` and `jtaf-xml2tf`; `--generic` redefined as a schema-scope
  flag now that there is one generator.
- `code-generation`: single schema-driven generator; Jinja2 Go templates removed; trimmed and full builds differ only
  in the embedded schema.
- `terraform-provider`: the resource exposes `groups`/`apply-groups` when the embedded schema carries them, instead of
  filtering those names unconditionally.
- `patch-engine`: paths below `groups[name=X]` and ordered `apply-groups` entries.
- `netconf-mock`: the mock keeps a device's configuration as one tree rather than a blob per group, so a
  configuration holding groups alongside base hierarchy, or more than one group, can be represented and patched.
- `testing`: an end-to-end example whose configuration lives in a group.

## Impact

- **Python CLI**: `junosterraform/jtaf-provider` (one generation path, `--groups`), `junosterraform/jtaf-yang2go`
  (flag pass-through), `junosterraform/jtaf-xml2tf` (`--groups`, keep group hierarchy),
  `junosterraform/jtaf_common.py` (stop stripping `apply-groups` when groups are kept).
- **Templates**: `junosterraform/templates/resource_config_provider.go.j2`, `provider.go.j2`, `config.go.j2` deleted;
  the Ansible template is untouched.
- **Go source**: `terraform_provider/generic/names.go` (schema-driven instead of name-based filtering),
  `generic/convert.go`, `generic/device.go`, `terraform_provider/patch/` (leafmap, diff, order),
  `terraform_provider/netconf/client.go` (legacy group helpers removed).
- **Examples**: new `examples/evpn-vxlan-dc-groups/`; `examples/providers/` scripts gain a groups build/convert;
  `examples/terraform_files/` gains a groups output directory.
- **Tests**: `junosterraform/tests/` (generation, xml2tf flattening vs. preservation), `terraform_provider/` and
  `patch/` Go tests, `netconf_mock/` group handling, and the both-providers e2e script.
- **Users**: anyone relying on the Jinja2-generated provider source must rebuild; the resource type, provider block and
  attribute shape are unchanged, so existing state and `.tf` files keep working. Providers built without `--groups`
  behave exactly as today.
