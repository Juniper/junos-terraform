# Hierarchical Groups with Merge Directives

This document describes the implementation of hierarchical YAML variable merging that eliminates JTAF-specific containers while supporting device-type groups and optional inline merge directives.

## Overview

This approach provides:

1. **✅ No JTAF containers** - Variables are flat (no `jtaf_shared` or `jtaf_override` wrappers)
2. **✅ Device-type hierarchy** - Variables organized by device type (QFX, SRX, MX, etc.)
3. **✅ Replace as default** - Scalar values replace by default; lists can be customized per-key
4. **✅ Automatic device detection** - Device type extracted from XML config 
5. **✅ Merge directives** - Optional `_merge_directive` meta-instructions for per-key control
6. **✅ Clean separation** - Global defaults, device-type overrides, host-specific deltas

## Quick Start

### 1. Generate Hierarchical Variables from XML

```bash
cd /path/to/ansible-provider-junos-qfx

# Generate YAML with automatic device type detection
jtaf-xml2yaml \
  -j trimmed_schema.json \
  -x dc1-spine1.xml dc1-spine2.xml dc1-leaf1.xml dc1-leaf2.xml \
  -d .
```

### 2. Generated Directory Structure

```
ansible-provider-junos-qfx/
├── group_vars/
│   ├── all.yml                    # Global defaults (all hosts)
│   └── device_qfx/
│       └── all.yml               # QFX-specific defaults
├── host_vars/
│   ├── dc1-spine1.yaml           # Spine1 device config (delta only)
│   ├── dc1-spine2.yaml           # Spine2 device config (delta only)
│   ├── dc1-leaf1.yaml            # Leaf1 device config (delta only)
│   └── dc1-leaf2.yaml            # Leaf2 device config (delta only)
├── hosts                          # Inventory with device type groups
└── jtaf-playbook.yml
```

### 3. Run Playbook

```bash
ansible-playbook -i hosts jtaf-playbook.yml
```

Merge order in playbook (precedence from lowest to highest):
```
group_vars/all.yml
        ↓
group_vars/device_qfx/all.yml
        ↓
host_vars/dc1-spine1.yaml
        ↓
jtaf_effective (final merged result)
```

## Merge Hierarchy

The merge order (precedence from lowest to highest) is:

```
1. group_vars/all.yml              (global defaults - all hosts)
2. group_vars/device_<type>/all.yml (e.g., device_qfx, device_srx)
3. host_vars/<hostname>.yaml       (host-specific overrides - device delta)
```

**Merge strategy:** Each level recursively combines with the previous level using Ansible's `combine()` filter with `recursive=True` and `list_merge='replace'` as the default. Merge directives can override per-key behavior.

**Result:** `jtaf_effective` variable containing the fully merged configuration.

## File Examples

### Example: group_vars/all.yml (Global Defaults)

Common configuration shared by all hosts (all device types):

```yaml
# Global defaults - apply to ALL hosts regardless of type
chassis:
  aggregated_devices:
    ethernet:
      device_count: 24

system:
  services:
    ssh:
      protocol_version: v2
      enabled: true
      port: 22
    netconf:
      enabled: true

syslog:
  host:
    - name: 10.0.0.1
      port: 514
```

### Example: group_vars/device_qfx/all.yml (Device Type Specific)

QFX-specific defaults that override global values:

```yaml
# QFX (data center) specific overrides
chassis:
  redundancy:
    enabled: true
    mode: single-re

system:
  services:
    netconf:
      enabled: true
      ssh_port: 830
  time_zone: UTC

# QFX-specific interfaces
interfaces:
  - name: et-0/0/16
    mtu: 9192
    speed: 400g
```

### Example: host_vars/dc1-spine1.yaml (Host Specific)

Host-specific configuration (device delta - only the differences):

```yaml
# Only host-specific overrides from the device-type defaults
system:
  host_name: dc1-spine1

routing:
  bgp:
    local_as: 65001
    group:
      - name: spines
        type: external

interfaces:
  - name: et-0/0/0
    description: "Link to dc1-leaf1"
    unit: "0"
    family_inet:
      address: 10.0.1.0/31

syslog:
  host:
    - name: 10.0.0.1
      _merge_directive: keep_parent  # Use global syslog config, don't override
```

## Merge Directive Reference

The `_merge_directive` meta-instruction allows fine-grained control over how specific values merge at the Ansible task level.

| Directive | Behavior | Use Case |
|-----------|----------|----------|
| `replace` | Use override value (default) | Most scalars and dicts |
| `keep_parent` | Ignore override, use parent | Don't override inherited values |
| `append` | Add override items to parent list | Accumulate list items |
| `prepend` | Add override items before parent | Priority items first |
| `extend` | Concatenate lists (type-safe) | Merge list contents |
| `merge_recursive` | Deep merge of dicts | Partial dict updates |

### Using _merge_directive

Any key in a YAML dict can have a `_merge_directive` sibling to control merge behavior:

```yaml
# group_vars/all.yml
interfaces:
  - name: lo0
    address: 127.0.0.1/32

# host_vars/dc1-spine1.yaml
interfaces:
  _merge_directive: append          # Append to parent list, don't replace
  - name: et-0/0/0
    address: 10.0.1.0/31
```

Result in `jtaf_effective`:
```yaml
interfaces:
  - name: lo0
    address: 127.0.0.1/32
  - name: et-0/0/0
    address: 10.0.1.0/31
```

Without the directive (default `replace`):
```yaml
# Would only have:
interfaces:
  - name: et-0/0/0
    address: 10.0.1.0/31
```

## Device Type Detection

Device type is automatically detected from the XML configuration and used to:
1. Create `group_vars/device_<type>/all.yml` group
2. Assign host to appropriate inventory group
3. Enable device-type-specific defaults in the playbook

### Supported Device Types

Detected from `system.product_name` or `chassis.product_name`:

- `qfx` - QFX series (data center switches)
- `srx` - SRX series (security appliances)  
- `mx` - MX series (core routers)
- `ptx` - PTX series (packet transport)
- `acx` - ACX series (access routers)
- `ex` - EX series (ethernet switches)

**Examples:**
- `Juniper Networks QFX5100-96S` → device type: `qfx`
- `SRX1500` → device type: `srx`
- `MX480` → device type: `mx`

**Implementation:** `junosterraform/jtaf-xml2yaml` → `detect_device_type()` function

### Customizing Device Detection

Edit `detect_device_type()` in `junosterraform/jtaf-xml2yaml` to add custom logic:

```python
def detect_device_type(xml_dict: dict[str, Any]) -> Optional[str]:
    """Detect device type from XML config."""
    if "system" in xml_dict and isinstance(xml_dict["system"], dict):
        product_name = xml_dict["system"].get("product_name", "").lower()
        if "custom_device" in product_name:
            return "custom"
    return None
```

## Inventory Generation

Inventory file (`hosts`) is auto-generated with device-type groups:

```ini
[all]
dc1-spine1
dc1-spine2
dc1-leaf1
dc1-leaf2

[device_qfx]
dc1-spine1
dc1-spine2
dc1-leaf1
dc1-leaf2
```

Use in playbooks:

```yaml
- hosts: all
  tasks:
    - name: Apply to all hosts
      ...

- hosts: device_qfx
  tasks:
    - name: QFX-specific configuration
      ...
```

## Ansible Playbook Integration

### In roles/qfx_role/tasks/main.yml

The generated Ansible role automatically:

1. **Detects device type** from host group membership
2. **Merges variables** in hierarchy: global → device-type → host
3. **Applies merge directives** via custom Ansible filter
4. **Creates jtaf_effective** for use in templates

```yaml
---
# Hierarchical variable merging with automatic device type detection
# Merge order: global (all) -> device_type -> host

# Detect device type group membership from inventory
- name: Detect host device type group
  set_fact:
    host_device_type: >-
      {{ group_names
         | select('match', '^device_')
         | map('regex_replace', '^device_', '')
         | first | default('unknown') }}

# Merge variables: device-type group vars first, then host-specific overrides
- name: Merge variables from hierarchy
  set_fact:
    jtaf_effective: >-
      {{
        (vars.get(host_device_type, {}) | default({}))
        | combine(
            (hostvars[inventory_hostname] | default({})),
            recursive=True,
            list_merge='replace'
          )
      }}
  tags: [merge_vars]

# Handle _merge_directive meta-instructions for custom merge behavior
- name: Apply custom merge directives
  set_fact:
    jtaf_effective: "{{ jtaf_effective | jtaf_apply_merge_directives }}"
  when: jtaf_effective is defined
  tags: [merge_directives]

- name: Applying template for qfx_role
  template: src=template.j2 dest={{ tmp_dir }}/{{ inventory_hostname }}.xml
```

### Using jtaf_effective in Templates

Access merged configuration in `roles/qfx_role/templates/template.j2`:

```jinja2
{# Access fully merged variables #}
<system>
    <host-name>{{ jtaf_effective.system.host_name }}</host-name>
    
    {% if jtaf_effective.system.services.ssh.enabled %}
    <services>
        <ssh>
            <enabled/>
            <protocol-version>{{ jtaf_effective.system.services.ssh.protocol_version }}</protocol-version>
        </ssh>
    </services>
    {% endif %}
</system>

{% for interface in jtaf_effective.interfaces | default([]) %}
<interfaces>
    <interface>
        <name>{{ interface.name }}</name>
        <mtu>{{ interface.mtu | default(1500) }}</mtu>
    </interface>
</interfaces>
{% endfor %}
```

## Filter Plugins

The `jtaf_apply_merge_directives` filter processes `_merge_directive` meta-instructions in the merged variables. It's automatically generated in `filter_plugins/jtaf_filters.py` when the Ansible role is created.

**Supported directives:**
- `replace` - Replace parent value (default)
- `keep_parent` - Keep parent, ignore override
- `append` - Append override to parent list
- `prepend` - Prepend override to parent list
- `extend` - Concatenate lists
- `merge_recursive` - Deep merge dicts

**Directive location in YAML:**

```yaml
# ✅ Correct - directive is a sibling key of the value being merged
config:
  _merge_directive: append
  - item1
  - item2

# ❌ Wrong - directive on parent is not recognized
_merge_directive: append
config:
  - item1
  - item2
```

## Command Line Options

```bash
jtaf-xml2yaml --help

  -j, --json-schema JSON_FILE
        Path to trimmed_schema.json (required)

  -x, --xml-config FILE [FILE ...]
        One or more XML config files to merge and process (required)

  -d, --directory DIR
        Output directory for host_vars, group_vars, and hosts (required)

  --auto-detect-hierarchy
        Automatically detect device types from XML and create device-type groups
        (enabled by default)
```

### Examples

```bash
# Initialize ansible role from YANG files
jtaf-yang2ansible -p ~/yang/common ~/yang/junos-qfx/conf/*.yang \
  -x config1.xml config2.xml -t qfx

# Generate host_vars and group_vars from XML configs
jtaf-xml2yaml -j trimmed_schema.json -x *.xml -d .

# Generate with explicit device detection
jtaf-xml2yaml -j trimmed_schema.json -x dc1-*.xml -d dc1_vars --auto-detect-hierarchy
```

## Workflow Example: Multi-Datacenter EVPN/VXLAN

See [HIERARCHICAL_GROUPS_EXAMPLE.md](./examples/HIERARCHICAL_GROUPS_EXAMPLE.md) for a complete working example using the EVPN/VXLAN sample data.

Quick reference:

```bash
# 1. Generate Ansible role from YANG files
cd examples
jtaf-yang2ansible \
  -p examples/yang/18.2/18.2R3/common examples/yang/18.2/18.2R3/junos-qfx/conf/*.yang \
  -x examples/evpn-vxlan-dc/dc1/*.xml \
  -t vqfx

# 2. Generate host_vars/group_vars
cd ansible-provider-junos-vqfx
jtaf-xml2yaml -j trimmed_schema.json -x ../examples/evpn-vxlan-dc/dc1/*.xml -d .

# 3. Run the playbook
ansible-playbook -i hosts jtaf-playbook.yml --check --diff
```

## Troubleshooting

### Merge not working as expected

1. **Check merge directive syntax** - Ensure `_merge_directive` is a sibling key
2. **Verify group membership** - Ensure hosts are in correct device groups via `ansible-inventory -i hosts --list`
3. **Debug merged vars** - Add debug task: `debug: var=jtaf_effective`

### Missing configuration in result

1. **Check hierarchy** - Values from higher precedence (host_vars) override lower precedence (group_vars)
2. **Check for typos** - YAML key names are case-sensitive
3. **Verify variables loaded** - Check that host_vars/<hostname>.yaml and group_vars/<type>/all.yml exist

### Custom filter not found

Ensure `filter_plugins/jtaf_filters.py` exists in the Ansible role directory. It's automatically created by `jtaf-ansible` generator.

## See Also

- [Ansible Variable Precedence](https://docs.ansible.com/ansible/latest/user_guide/playbooks_variables.html#variable-precedence-where-should-i-put-a-variable)
- [Ansible Group Variables](https://docs.ansible.com/ansible/latest/user_guide/intro_inventory.html#group-variable-files)
- [Ansible combine Filter](https://docs.ansible.com/ansible/latest/collections/ansible/builtin/combine_filter.html)
- [JTAF Project README](./README.md)
