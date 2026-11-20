# Demo: Configuration Groups

This walkthrough manages a network whose configuration lives in Junos
configuration groups, using `examples/evpn-vxlan-dc-groups/`.

It sits beside `examples/evpn-vxlan-dc/`, the base-configuration example, and
covers the same devices. The difference is where each device's configuration
sits: everything except the host name moves into one group named after the
device's role, and the base hierarchy carries the `apply-groups` reference that
activates it.

```xml
<configuration>
  <apply-groups>spine</apply-groups>
  <groups>
    <name>spine</name>
    <system>...</system>
    <interfaces>...</interfaces>
    <protocols>...</protocols>
  </groups>
  <system>
    <host-name>dc1-spine1</host-name>
  </system>
</configuration>
```

The example is derived from the base one by
`examples/evpn-vxlan-dc-groups/derive-from-base.py`, so the two stay device for
device in step. The output is committed; the script only needs running when the
base example changes.

---

## Prerequisites

The same as the other examples: `pyang`, the JTAF CLI tools, Go and Terraform on
PATH. See [README-terraform.md](../README-terraform.md).

---

## 1. Build a provider that manages groups

```bash
cd examples/providers
bash build-groups.sh
```

This runs `jtaf-yang2go --groups`, which keeps the `groups` subtree and the
`apply-groups` leaf-list in the embedded schema instead of leaving them out, and
trims both to the paths the example XML uses.

Trimming is what makes this affordable. The group subtree repeats the whole
configuration model, so in an untrimmed model it is about half the nodes; here
it is trimmed to the configuration the example actually has:

| Build | Schema nodes |
|-------|--------------|
| Base example (`build.sh`) | 212 |
| Group example (`build-groups.sh`) | 216 |

Building with `--generic --groups` instead would advertise the whole model twice
over, which needs several GB at plan time and warns when generated.

---

## 2. Convert the configuration

```bash
bash convert-groups.sh
```

`jtaf-xml2tf --groups` keeps the hierarchy as the devices hold it. Without the
flag it would flatten what each device inherits through `apply-groups` into the
base hierarchy, which is what a provider built without `--groups` expects.

The result is written to `examples/terraform_files_groups/`, one file per
device:

```hcl
resource "terraform-provider-junos-vqfx-evpn-vxlan-groups" "dc1-spine1-base-config" {
  resource_name = "base-config"
  provider      = junos-vqfx-evpn-vxlan-groups.dc1_spine1
  apply_groups  = local.common_g_0c49e5_apply_groups
  groups = [
    {
      name       = "spine"
      system     = local.common_g_0c49e5_groups_spine_system
      interfaces = [...]
    }
  ]
}
```

Configuration that is identical across devices — which a group-based network has
a lot of — is extracted into `common.tf` as locals and referenced, so the group
bodies are not repeated per device. Pass `--no-extract-common` to
`jtaf-xml2tf` to emit everything inline instead.

A group is an entry of the `groups` list keyed by its `name`, holding the same
attributes it would have in the base hierarchy. `apply_groups` is an ordered
list: Junos applies groups in the order they are listed, so changing the order
is a change.

---

## 3. Plan and apply

```bash
cd ../terraform_files_groups

# Point Terraform at the provider you just built; see examples/example-terraformrc
export TF_CLI_CONFIG_FILE=~/.terraformrc

terraform validate
terraform plan
terraform apply -auto-approve
```

A second plan reports no changes.

---

## What to expect from the provider

- **A change inside a group** is sent as an edit below that group's `<groups>`
  element, leaving the rest of the group, the other groups and the base
  hierarchy alone.
- **Removing a group** is a single delete of the group, not one delete per leaf.
- **Group references are never invented.** A group declared without an
  `apply_groups` entry is written to the device with no reference to it; the
  provider does not add, reorder or retain references the configuration does not
  declare.
- **Reads do not resolve inheritance.** The provider asks the device for its
  committed configuration, not the inherited view. A device that resolves
  `apply-groups` reports a group's values in the base hierarchy too, and
  configuration that was never declared there would otherwise be read back as
  drift on every plan.
