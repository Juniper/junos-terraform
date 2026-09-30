# Junos Ansible Guide

This guide covers generating Ansible roles, playbooks, and inventory data from YANG models and XML configurations using JTAF. For initial setup and YANG → JSON conversion, see the [main README](README.md).

---

## Generate Ansible Playbook

Create an Ansible role + playbook from a Junos JSON schema and one or more XML configs. The generated playbook runs locally and renders configs (does not connect to devices) by default.

Quick usage:
```
jtaf-ansible -j <junos.json> -x <config1.xml> [-x <config2.xml> ...] -t <device-type>
```

What is created (under ansible-provider-junos-<type>/):
- roles/<type>_role/ (tasks/main.yml, templates/template.j2)
- jtaf-playbook.yml (uses connection: local)
- host_vars/, configs/, trimmed_schema.json.gz

Verify rendering without applying:
```
cd ansible-provider-junos-<type>
ansible-playbook -i "localhost," jtaf-playbook.yml --check --diff
```

---

### Single command to generate ansible role

Generate an Ansible role + playbook in one step from YANG files and XML config(s):

```
jtaf-yang2ansible -p <path-to-common> <path-to-yang-files> -x <xml-config(s)> -t <device-type>
```

Example:
```
jtaf-yang2ansible -p examples/yang/18.2/18.2R3/common examples/yang/18.2/18.2R3/junos-qfx/conf/*.yang -x examples/evpn-vxlan-dc/dc1/*spine*.xml -t vqfx
```

Notes:
- If supplying multiple XML configs they must be for the same device type.
- Output directory: ansible-provider-junos-<type>/ containing roles/<type>_role/ (tasks/templates), jtaf-playbook.yml (connection: local), host_vars/, group_vars/, configs/, trimmed_schema.json.gz.
- Run the generated playbook in check/diff mode to verify rendered configs without applying:
	ansible-playbook -i "localhost," jtaf-playbook.yml --check --diff

---

## Generate YAML for Ansible host_vars (jtaf-xml2yaml)

Convert one or more Junos XML configs into Ansible `host_vars`, `group_vars`, and inventory data.

Important behavior:
- Each run should use one `trimmed_schema.json.gz` (plain JSON is also accepted) and matching XML configs for the generated role.
- `--grouping-hosts-file` is required.
- Inventory groups and optional `:children` sections come from the `grouping.hosts` file, not from a `-t/--type` flag.
- Output can be split across directories so generated role location and provisioning playbook location can differ.
- Repeated runs are merge-safe when they reuse the same `-d/--directory` output directory:
  - `group_vars/all.yaml` stores values shared across all hosts currently tracked in that output directory.
  - `group_vars/<group>/all.yaml` stores per-group deltas for the groups defined in `grouping.hosts`.
  - Host-specific differences are preserved in each `host_vars/<host>.yaml` file.
  - Existing inventory hosts/groups are merged (not clobbered).

Usage:
```
jtaf-xml2yaml -j <trimmed_schema.json.gz> -x <config1.xml> [<config2.xml> ...] -d <output-dir> --grouping-hosts-file <grouping_hosts_file>

```

Example:
```
jtaf-xml2yaml -j ansible-provider-junos-vqfx/trimmed_schema.json.gz \
	-x examples/evpn-vxlan-dc/dc1/dc1-leaf1.xml examples/evpn-vxlan-dc/dc1/dc1-spine1.xml \
  -d ansible_files \
  --grouping-hosts-file examples/ansible/switches_grouping_hosts
```

Output:
- Creates `host_vars/<hostname>.yaml` for every XML file provided (hostname is file base name or `system/host-name` from XML).
- Writes `group_vars/all.yaml` for keys shared across all tracked hosts in the output directory.
- Writes `group_vars/<group>/all.yaml` for each group defined in `grouping.hosts`.
- Writes/updates inventory hosts file with `[all]`, `[group]`, and `[group:children]` sections taken from `grouping.hosts`.

This output feeds into the Ansible role/playbook created by jtaf-ansible/jtaf-yang2ansible.

---

## How the Ansible Role Works

### Overview

The generated Ansible role uses **Jinja2 templates** to render Junos XML configuration from YAML variables. The workflow is:

```
YANG + XML → jtaf-ansible → Role (tasks + template.j2)
                              ↓
XML configs → jtaf-xml2yaml → host_vars/ + group_vars/ + inventory
                              ↓
ansible-playbook → Render XML per host → Preview diff → Push via NETCONF → Verify
```

| Step | What Happens |
|------|-------------|
| **Render** | Jinja2 template produces XML config from host/group vars |
| **Preview** | `juniper.device.config` with `check: true` shows candidate diff |
| **Push** | `load replace` with `commit confirmed` applies config safely |
| **Verify** | Commit confirm ensures rollback on failure |

### Key Differences from Terraform

| Aspect | Terraform Provider | Ansible Role |
|--------|-------------------|--------------|
| **Config format** | HCL (`.tf` files) | YAML (`host_vars/`, `group_vars/`) |
| **Change detection** | NETCONF patch engine (leaf-level diff) | Junos candidate diff (`show | compare`) |
| **Push method** | `edit-config` with `nc:operation` per leaf | `load replace` full XML |
| **Rollback** | Terraform state + fallback to full replace | `commit confirmed` timer |
| **Multi-device** | One provider instance per device | Inventory groups + parallel execution |
| **Variables** | Inline in HCL resource blocks | Hierarchical YAML (all → group → host) |

### Variable Hierarchy

When `jtaf-xml2yaml` processes multiple XML configs, it produces a layered variable structure:

```
group_vars/all.yaml           ← Values shared by ALL hosts
group_vars/<group>/all.yaml   ← Values shared within a group (e.g., all spines)
host_vars/<hostname>.yaml     ← Host-specific overrides
```

At render time, Ansible merges these layers (host > group > all), so the template receives the final merged values for each host.

---

## Override Mode (`load override`)

JTAF supports two operating modes controlled by the `jtaf_mode` variable:

| Mode | Behavior | Use When |
|------|----------|----------|
| `group` (default) | Wraps config in `<groups>JTAF_ANSIBLE</groups>`, pushes with `load replace` | You manage a subset of the device config alongside other config |
| `override` | Renders bare `<configuration>`, pushes with `load override` + `commit confirmed` | You own the entire device config — what you provide IS the box |

To enable override mode, set `jtaf_mode: "override"` in your role defaults, `group_vars`, or `host_vars`. You can also pass `--mode override` when generating the role with `jtaf-ansible` or `jtaf-yang2ansible`.

### Important Notes

**Override mode deletes all unmanaged config.** When using `load override`, the rendered XML replaces the *entire* device configuration. Your rendered XML **must** include:
- Management interface IP address (e.g., `fxp0` or `em0`)
- Root authentication (`system/root-authentication`)
- NETCONF SSH service (`system/services/netconf/ssh`)

If any of these are missing from the rendered template, the device may become unreachable after commit — at which point the `commit confirmed` timer will auto-rollback.

**Large schemas may increase generation time.** If using full-model parse (all YANG files without XML filtering), the trimmed schema and generated Jinja2 template will be larger. This affects only the one-time `jtaf-ansible` / `jtaf-xml2yaml` generation step, not playbook runtime. Template rendering remains fast because Jinja2 `{% if ... is defined %}` conditionals skip undefined paths.

**Set the commit confirmed timer appropriately.** Override mode uses `commit confirmed` as a safety net — if the device becomes unreachable after the override, Junos auto-rolls back when the timer expires. Configure `jtaf_commit_confirm_minutes` (default: 2) based on your expected commit time. For large configurations, consider 5 minutes.

**NETCONF access is required.** The `juniper.device.config` module connects to devices via NETCONF (port 830). Ensure NETCONF SSH is enabled on target devices before running the playbook.

**Backward compatibility.** Group mode remains the default. No existing behavior changes unless you explicitly set `--mode override` or `jtaf_mode: "override"`. Existing roles and playbooks continue to work as before.

---

## Beginner Guide: Apply a generated role from a separate playbook directory

This walkthrough is for first-time Ansible users and reflects the recommended split layout:
- JTAF-generated role in one directory.
- Operator playbook/inventory/vars in another directory.

### 1. Install Ansible dependencies on your control node

Note that Ansible must be installed in the virtual environment (venv). Additionally, the following system packages and Python modules are required:

```bash
python -m pip install --upgrade pip
sudo dnf install ansible -y
sudo dnf install python3-pip -y
/usr/bin/python3 -m pip install ncclient junos-eznc jxmlease

# Install Juniper collection used to push config.
ansible-galaxy collection install juniper.device
```

### 2. Generate the first role (QFX EVPN-VXLAN) from YANG + XML

```bash
# Generate role + templates from YANG + XML
jtaf-yang2ansible \
	-p examples/yang/18.2/18.2R3/common \
	examples/yang/18.2/18.2R3/junos-qfx/conf/*.yang \
	-x \
	examples/evpn-vxlan-dc/dc1/dc1-borderleaf1.xml \
	examples/evpn-vxlan-dc/dc1/dc1-borderleaf2.xml \
	examples/evpn-vxlan-dc/dc1/dc1-leaf1.xml \
	examples/evpn-vxlan-dc/dc1/dc1-leaf2.xml \
	examples/evpn-vxlan-dc/dc1/dc1-leaf3.xml \
	examples/evpn-vxlan-dc/dc1/dc1-spine1.xml \
	examples/evpn-vxlan-dc/dc1/dc1-spine2.xml \
	examples/evpn-vxlan-dc/dc2/dc2-spine1.xml \
	examples/evpn-vxlan-dc/dc2/dc2-spine2.xml \
	-t vqfx-evpn-vxlan
```

### 3. Generate a second role (SRX firewalls) from YANG + XML

```bash
jtaf-yang2ansible \
  -p examples/yang/18.2/18.2R3/common \
  examples/yang/18.2/18.2R3/junos-es/conf/*.yang \
  -x examples/evpn-vxlan-dc/dc1/dc1-*firewall*.xml examples/evpn-vxlan-dc/dc2/dc2-*firewall*.xml \
  -t srx-ansible-role
```

### 4. Create a separate provisioning playbook project

Create a separate directory for your operator playbook:

```bash
mkdir -p ansible-evpn-vxlan-deploy
```

Create `ansible-evpn-vxlan-deploy/ansible.cfg`:

```ini
[defaults]
roles_path = ../ansible-provider-junos-vqfx-evpn-vxlan/roles:../ansible-provider-junos-srx-ansible-role/roles
host_key_checking = False
```

For first-time Ansible users: `roles_path` tells Ansible where custom roles live. In this workflow, both generated roles are referenced, while your operator playbook stays in `ansible-evpn-vxlan-deploy/`.

### 5. Create `grouping.hosts` files for the inventory hierarchy

`jtaf-xml2yaml` now requires a grouping definition. The section names in these files become your generated inventory groups and `group_vars/<group>/all.yaml` directories.

Create `ansible-evpn-vxlan-deploy/qfx.grouping.hosts`:

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

[borderleaf]
dc1-borderleaf1
dc1-borderleaf2

[leaf]
dc1-leaf1
dc1-leaf2
dc1-leaf3

[spine]
dc1-spine1
dc1-spine2
dc2-spine1
dc2-spine2
```

Create `ansible-evpn-vxlan-deploy/firewall.grouping.hosts`:

```ini
[all]
dc1-firewall1
dc1-firewall2
dc2-firewall1
dc2-firewall2

[firewall]
dc1-firewall1
dc1-firewall2
dc2-firewall1
dc2-firewall2
```

### 6. Generate inventory + vars for the first role into the playbook project

Use the same `-d` directory for every `jtaf-xml2yaml` run that should share one inventory, `group_vars`, `host_vars`, and payload cache.

```bash
jtaf-xml2yaml \
	-x examples/evpn-vxlan-dc/dc1/*{spine,leaf}*.xml examples/evpn-vxlan-dc/dc2/*spine*.xml \
	-j ansible-provider-junos-vqfx-evpn-vxlan/trimmed_schema.json.gz \
  -d ansible-evpn-vxlan-deploy \
  --hosts-file ansible-evpn-vxlan-deploy/inventory.ini \
  --grouping-hosts-file ansible-evpn-vxlan-deploy/qfx.grouping.hosts
```

### 7. Generate inventory + vars for the second role into the same playbook project

```bash
jtaf-xml2yaml \
  -x examples/evpn-vxlan-dc/dc1/dc1-*firewall*.xml examples/evpn-vxlan-dc/dc2/dc2-*firewall*.xml \
  -j ansible-provider-junos-srx-ansible-role/trimmed_schema.json.gz \
  -d ansible-evpn-vxlan-deploy \
  --hosts-file ansible-evpn-vxlan-deploy/inventory.ini \
  --grouping-hosts-file ansible-evpn-vxlan-deploy/firewall.grouping.hosts
```

After both runs, your playbook project should contain at least:
- `inventory.ini`
- `group_vars/all.yaml`
- `group_vars/borderleaf/all.yaml`
- `group_vars/leaf/all.yaml`
- `group_vars/spine/all.yaml`
- `group_vars/firewall/all.yaml`
- `host_vars/<hostname>.yaml`

Update `ansible-evpn-vxlan-deploy/inventory.ini` with reachable management addresses while keeping the generated group names:

```ini
[borderleaf]
dc1-borderleaf1 ansible_host=192.0.2.101 ansible_port=830
dc1-borderleaf2 ansible_host=192.0.2.102 ansible_port=830

[leaf]
dc1-leaf1 ansible_host=192.0.2.11 ansible_port=830
dc1-leaf2 ansible_host=192.0.2.12 ansible_port=830
dc1-leaf3 ansible_host=192.0.2.13 ansible_port=830

[spine]
dc1-spine1 ansible_host=192.0.2.21 ansible_port=830
dc1-spine2 ansible_host=192.0.2.22 ansible_port=830
dc2-spine1 ansible_host=192.0.2.31 ansible_port=830
dc2-spine2 ansible_host=192.0.2.32 ansible_port=830

[firewall]
dc1-firewall1 ansible_host=192.0.2.201 ansible_port=830
dc1-firewall2 ansible_host=192.0.2.202 ansible_port=830
dc2-firewall1 ansible_host=192.0.2.203 ansible_port=830
dc2-firewall2 ansible_host=192.0.2.204 ansible_port=830
```

Notes on repeated runs:
- Reuse the same `-d` directory whenever you want one merged inventory and var tree.
- `group_vars/all.yaml` contains values shared across every tracked host in that output directory.
- `group_vars/<group>/all.yaml` contains per-group deltas for groups declared in the relevant `grouping.hosts` file.
- Host-specific differences remain in `host_vars/<hostname>.yaml`.

### 8. Create a playbook that renders, previews diff, pushes, and verifies

Create `ansible-evpn-vxlan-deploy/site.yml`:

```yaml
---
- name: Render XML from generated QFX role
  hosts: borderleaf:leaf:spine
  gather_facts: false
  connection: local
  vars:
    tmp_dir: ../ansible-provider-junos-vqfx-evpn-vxlan/configs
  roles:
    - role: vqfx-evpn-vxlan_role
      delegate_to: localhost

- name: Render XML from generated SRX role
  hosts: firewall
  gather_facts: false
  connection: local
  vars:
    tmp_dir: ../ansible-provider-junos-srx-ansible-role/configs
  roles:
    - role: srx-ansible-role_role
      delegate_to: localhost

- name: Preview and apply rendered XML on QFX devices
  hosts: borderleaf:leaf:spine
  gather_facts: false
  connection: local
  vars:
    netconf_user: "{{ lookup('env', 'NETCONF_USERNAME') }}"
    netconf_pass: "{{ lookup('env', 'NETCONF_PASSWORD') }}"
    tmp_dir: ../ansible-provider-junos-vqfx-evpn-vxlan/configs
  tasks:
    - name: Preview candidate diff without committing
      juniper.device.config:
        host: "{{ ansible_host | default(inventory_hostname) }}"
        port: "{{ ansible_port | default(830) }}"
        user: "{{ netconf_user }}"
        passwd: "{{ netconf_pass }}"
        load: replace
        src: "{{ tmp_dir }}/{{ inventory_hostname }}.xml"
        check: true
        commit: false
        diff: true
      register: preview_result

    - name: Print diff lines from preview
      ansible.builtin.debug:
        var: preview_result.diff_lines

    - name: Load and commit with commit-confirm safeguard
      juniper.device.config:
        host: "{{ ansible_host | default(inventory_hostname) }}"
        port: "{{ ansible_port | default(830) }}"
        user: "{{ netconf_user }}"
        passwd: "{{ netconf_pass }}"
        load: replace
        src: "{{ tmp_dir }}/{{ inventory_hostname }}.xml"
        confirmed: 5
        check_commit_wait: 5
        comment: "Apply EVPN-VXLAN config generated by JTAF"
      register: apply_result

    - name: Verify module-reported apply result
      ansible.builtin.assert:
        that:
          - not (apply_result.failed | default(false))
          - apply_result.msg is defined
        fail_msg: "Config apply failed on {{ inventory_hostname }}"

    - name: Confirm previously confirmed commit
      juniper.device.config:
        host: "{{ ansible_host | default(inventory_hostname) }}"
        port: "{{ ansible_port | default(830) }}"
        user: "{{ netconf_user }}"
        passwd: "{{ netconf_pass }}"
        check: true
        commit: false
        diff: false
      register: confirm_result

    - name: Print apply and confirm summaries
      ansible.builtin.debug:
        msg:
          - "apply={{ apply_result.msg | default('no message') }}"
          - "confirm={{ confirm_result.msg | default('no message') }}"

- name: Preview and apply rendered XML on SRX devices
  hosts: firewall
  gather_facts: false
  connection: local
  vars:
    netconf_user: "{{ lookup('env', 'NETCONF_USERNAME') }}"
    netconf_pass: "{{ lookup('env', 'NETCONF_PASSWORD') }}"
    tmp_dir: ../ansible-provider-junos-srx-ansible-role/configs
  tasks:
    - name: Preview candidate diff without committing (SRX)
      juniper.device.config:
        host: "{{ ansible_host | default(inventory_hostname) }}"
        port: "{{ ansible_port | default(830) }}"
        user: "{{ netconf_user }}"
        passwd: "{{ netconf_pass }}"
        load: replace
        src: "{{ tmp_dir }}/{{ inventory_hostname }}.xml"
        check: true
        commit: false
        diff: true
      register: preview_result_srx

    - name: Load and commit with commit-confirm safeguard (SRX)
      juniper.device.config:
        host: "{{ ansible_host | default(inventory_hostname) }}"
        port: "{{ ansible_port | default(830) }}"
        user: "{{ netconf_user }}"
        passwd: "{{ netconf_pass }}"
        load: replace
        src: "{{ tmp_dir }}/{{ inventory_hostname }}.xml"
        confirmed: 5
        check_commit_wait: 5
        comment: "Apply SRX config generated by JTAF"
      register: apply_result_srx
```

### 9. Run the playbook

```bash
cd ansible-evpn-vxlan-deploy

# NETCONF credentials used by the playbook
export NETCONF_USERNAME='<junos-netconf-user>'
export NETCONF_PASSWORD='<junos-netconf-password>'

ansible-playbook -i inventory.ini site.yml
```

### 10. What to check in output

- The preview tasks should show diffs for the generated switch groups (`borderleaf`, `leaf`, and `spine`) and for `firewall`.
- The apply tasks should succeed for both generated roles.
- Inventory and vars should remain merged across repeated `jtaf-xml2yaml` runs.

At this point you have completed render -> preview diff -> push -> plugin-level verification using both generated roles (QFX and SRX).
