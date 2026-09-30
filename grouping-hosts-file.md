# `--grouping-hosts-file` Knob for `jtaf-xml2yaml`

## Overview

The `--grouping-hosts-file` flag gives operators explicit control over how Ansible inventory groups and `group_vars` directories are structured when converting Junos XML configs into Ansible host/group variable files.

A single grouping file can cover **all device roles and all data centers**. Both the QFX and SRX `jtaf-xml2yaml` runs reference this one file, so there is one consistent inventory definition for the entire deployment.

---

## Why It Was Added

In real deployments, devices are logically grouped by role (e.g., `spine`, `leaf`, `firewall`). With `--grouping-hosts-file`, the operator fully owns the group topology and the tool produces a predictable, stable inventory structure across repeated runs.

---

## The Grouping Hosts File Format

The file uses the same INI syntax as an Ansible inventory file. Two types of sections are supported:

| Section syntax | Purpose |
|---|---|
| `[all]` | Lists every host that should appear in the `[all]` section of the generated inventory |
| `[groupname]` | Lists the hosts that belong to a flat inventory group called `groupname` |

Lines beginning with `#` are treated as comments. Each host or group token is one entry per line.

---

## Recommended Layout — One File for All Devices and Data Centers

Define every device role and every data center in a single `grouping.hosts` file. This example covers a two-DC deployment with QFX switching (borderleaf, leaf, spine) and SRX firewalls:

Create `ansible-evpn-vxlan-deploy/grouping.hosts`:

```ini
[all]
dc1-borderleaf1
dc1-borderleaf2
dc1-leaf1
dc1-leaf2
dc1-leaf3
dc1-spine1
dc1-spine2
dc2-spine1
dc2-spine2
dc1-firewall1
dc1-firewall2
dc2-firewall1
dc2-firewall2

[dc1-borderleaf]
dc1-borderleaf1
dc1-borderleaf2

[dc1-leaf]
dc1-leaf1
dc1-leaf2
dc1-leaf3

[dc1-spine]
dc1-spine1
dc1-spine2

[dc2-leafspine]
dc2-spine1
dc2-spine2

[dc1-firewall]
dc1-firewall1
dc1-firewall2

[dc2-firewall]
dc2-firewall1
dc2-firewall2
```

---

## Running `jtaf-xml2yaml` for Multiple Providers Against the Same File

Both the QFX and SRX roles reference the **same** grouping file. Run one command per generated role, targeting the same `-d` output directory each time.

**QFX role (borderleaf, leaf, spine — both DCs):**

```bash
jtaf-xml2yaml \
  -j ansible-provider-junos-vqfx-evpn-vxlan/trimmed_schema.json.gz \
  -x examples/evpn-vxlan-dc/dc1/*{spine,leaf}*.xml \
     examples/evpn-vxlan-dc/dc2/*spine*.xml \
  -d ansible-evpn-vxlan-deploy \
  --hosts-file ansible-evpn-vxlan-deploy/inventory.ini \
  --grouping-hosts-file ansible-evpn-vxlan-deploy/grouping.hosts
```

**SRX role (firewalls — both DCs):**

```bash
jtaf-xml2yaml \
  -j ansible-provider-junos-srx-ansible-role/trimmed_schema.json.gz \
  -x examples/evpn-vxlan-dc/dc1/dc1-*firewall*.xml \
     examples/evpn-vxlan-dc/dc2/dc2-*firewall*.xml \
  -d ansible-evpn-vxlan-deploy \
  --hosts-file ansible-evpn-vxlan-deploy/inventory.ini \
  --grouping-hosts-file ansible-evpn-vxlan-deploy/grouping.hosts
```

After both runs the playbook project contains:

```
ansible-evpn-vxlan-deploy/
├── inventory.ini
├── group_vars/
│   ├── all.yaml              ← keys shared across ALL 13 devices
│   ├── dc1-borderleaf/
│   │   └── all.yaml          ← keys shared by the two DC1 borderleaf devices
│   ├── dc1-leaf/
│   │   └── all.yaml          ← keys shared by the three DC1 leaf devices
│   ├── dc1-spine/
│   │   └── all.yaml          ← keys shared by the two DC1 spine devices
│   ├── dc2-leafspine/
│   │   └── all.yaml          ← keys shared by the two DC2 spine devices
│   ├── dc1-firewall/
│   │   └── all.yaml          ← keys shared by the two DC1 firewall devices
│   └── dc2-firewall/
│       └── all.yaml          ← keys shared by the two DC2 firewall devices
└── host_vars/
    ├── dc1-borderleaf1.yaml
    ├── dc1-leaf1.yaml
    ├── dc1-spine1.yaml
    ├── dc1-firewall1.yaml
    └── ...                   ← one file per device, host-unique delta only
```

And `inventory.ini` contains all groups and all 13 devices:

```ini
[all]
dc1-borderleaf1
dc1-borderleaf2
dc1-leaf1
dc1-leaf2
dc1-leaf3
dc1-spine1
dc1-spine2
dc2-spine1
dc2-spine2
dc1-firewall1
dc1-firewall2
dc2-firewall1
dc2-firewall2

[dc1-borderleaf]
dc1-borderleaf1
dc1-borderleaf2

[dc1-leaf]
dc1-leaf1
dc1-leaf2
dc1-leaf3

[dc1-spine]
dc1-spine1
dc1-spine2

[dc2-leafspine]
dc2-spine1
dc2-spine2

[dc1-firewall]
dc1-firewall1
dc1-firewall2

[dc2-firewall]
dc2-firewall1
dc2-firewall2
```

`group_vars/dc1-spine/all.yaml` contains config shared by the two DC1 spine devices. DC2 spine devices live in the separate `dc2-leafspine` group, so differences between the two data centers are captured at the group level rather than falling through to `host_vars`.

---

## Variable Deduplication Behaviour

`jtaf-xml2yaml` performs a three-level deduplication pass driven entirely by the group structure declared in the grouping hosts file:

1. **`group_vars/all.yaml`** — config keys whose value is identical across every host currently tracked in the output directory (all roles combined).
2. **`group_vars/<group>/all.yaml`** — keys that are identical within a group but differ from the global shared value.
3. **`host_vars/<hostname>.yaml`** — only the residual, host-unique differences.

If a host's config conflicts with its group siblings, the conflicting keys are placed directly in that host's `host_vars` file rather than in the group, preserving correctness.

The grouping file need not list devices that belong to another role. When the SRX run processes firewall XML files, `jtaf-xml2yaml` only reads the `[all]` and `[firewall]` sections that are relevant to the active hosts — spine/leaf/borderleaf entries in the same file are silently ignored for that run.

---

## Validation and Error Handling

The parser validates the grouping file at startup and exits with a clear message for:
- Missing or unreadable file
- Malformed `[section]` headers
- Duplicate section names
- Hosts in a group that were not provided via `-x` (when `--strict-grouping-known-hosts` is also set)

---

## Summary

| Aspect | Behaviour |
|---|---|
| One file for all roles and DCs | A single `grouping.hosts` covers QFX and SRX devices across both data centers |
| Stable inventory shape | Group topology is version-controlled and reproducible across repeated runs |
| Cross-provider shared vars | `group_vars/all.yaml` accumulates keys common to every tracked host, regardless of which role wrote them |
| Per-group vars | `group_vars/<group>/all.yaml` is written for every `[groupname]` section in the file |
| Safe repeated runs | Re-running either role updates only that role's hosts; other groups are preserved |
