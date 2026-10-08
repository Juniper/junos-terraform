# Junos Ansible Guide

Follow this command-by-command walkthrough to generate roles, create inventory
and variables, render configuration, and deploy to lab devices. Complete
[the initial setup](README.md#setup) first, then start in the repository root.

Generated roles stay separate from the operator-owned playbook and inventory.
Steps 2 and 3 generate QFX and SRX roles; steps 4 through 7 prepare one shared
deployment directory; steps 8 through 10 render, preview, apply, and verify.
Reference material and optional modes follow the walkthrough.

### 1. Install Ansible dependencies on your control node

Install Ansible and its Python dependencies in the same virtual environment as
JTAF. Use an Ansible release compatible with your Python version.

```bash
. venv/bin/activate
python -m pip install --upgrade pip
python -m pip install ansible ncclient junos-eznc jxmlease

# Install Juniper collection used to push config.
ansible-galaxy collection install juniper.device
```

The deployment steps require NETCONF over SSH on every target device and
reachable management addresses. The bundled 18.2 models are for the examples;
choose matching platform and release models for your own devices.

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

`-p` supplies the common model search path and YANG files. `-x` supplies XML
configurations to filter the schema; `-t` names the output directory and role.
`jtaf-yang2ansible` runs `pyang` and `jtaf-ansible` in one command.

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
filter_plugins = ../ansible-provider-junos-vqfx-evpn-vxlan/filter_plugins:../ansible-provider-junos-srx-ansible-role/filter_plugins
interpreter_python = auto_silent
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
  tags: [render]
  vars:
    tmp_dir: ../ansible-provider-junos-vqfx-evpn-vxlan/configs
    jtaf_vars_root: "{{ playbook_dir }}"
  roles:
    - role: vqfx-evpn-vxlan_role
      delegate_to: localhost

- name: Render XML from generated SRX role
  hosts: firewall
  gather_facts: false
  connection: local
  tags: [render]
  vars:
    tmp_dir: ../ansible-provider-junos-srx-ansible-role/configs
    jtaf_vars_root: "{{ playbook_dir }}"
  roles:
    - role: srx-ansible-role_role
      delegate_to: localhost

- name: Preview and apply rendered XML on QFX devices
  hosts: borderleaf:leaf:spine
  gather_facts: false
  connection: local
  tags: [deploy]
  vars:
    netconf_user: "{{ lookup('env', 'NETCONF_USERNAME') }}"
    netconf_pass: "{{ lookup('env', 'NETCONF_PASSWORD') }}"
    tmp_dir: ../ansible-provider-junos-vqfx-evpn-vxlan/configs
  tasks:
    - name: Preview candidate diff without committing
      tags: [preview]
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
      tags: [preview]
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
  tags: [deploy]
  vars:
    netconf_user: "{{ lookup('env', 'NETCONF_USERNAME') }}"
    netconf_pass: "{{ lookup('env', 'NETCONF_PASSWORD') }}"
    tmp_dir: ../ansible-provider-junos-srx-ansible-role/configs
  tasks:
    - name: Preview candidate diff without committing (SRX)
      tags: [preview]
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

    - name: Print diff lines from preview (SRX)
      tags: [preview]
      ansible.builtin.debug:
        var: preview_result_srx.diff_lines

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

    - name: Verify module-reported apply result (SRX)
      ansible.builtin.assert:
        that:
          - not (apply_result_srx.failed | default(false))
          - apply_result_srx.msg is defined
        fail_msg: "Config apply failed on {{ inventory_hostname }}"

    - name: Confirm previously confirmed commit (SRX)
      juniper.device.config:
        host: "{{ ansible_host | default(inventory_hostname) }}"
        port: "{{ ansible_port | default(830) }}"
        user: "{{ netconf_user }}"
        passwd: "{{ netconf_pass }}"
        check: true
        commit: false
        diff: false
```

### 9. Render, inspect, and deploy

Use the active virtual environment's Python for locally executed modules.
Render first without device connections:

```bash
cd ansible-evpn-vxlan-deploy
export ANSIBLE_PYTHON_INTERPRETER="$(command -v python)"
ansible-playbook -i inventory.ini site.yml --syntax-check
ansible-playbook -i inventory.ini site.yml --tags render
```

Inspect XML under each generated role's `configs/` directory. Then provide
NETCONF credentials using your normal secure environment setup and preview:

```bash

# NETCONF credentials used by the playbook
export NETCONF_USERNAME='<junos-netconf-user>'
export NETCONF_PASSWORD='<junos-netconf-password>'

ansible-playbook -i inventory.ini site.yml --tags preview
```

The preview connects to devices but does not commit. Review the diffs, then
explicitly deploy the already rendered configuration:

```bash
ansible-playbook -i inventory.ini site.yml --tags deploy
```

The deploy tasks preview again, commit with a five-minute confirmation timer,
check the module result, and confirm. These examples target every inventory
device in the play's groups; use `--limit <host-or-group>` on each run for a
smaller target set. Never run an unreviewed example against production devices.

### 10. What to check in output

- The preview tasks should show diffs for the generated switch groups (`borderleaf`, `leaf`, and `spine`) and for `firewall`.
- The apply tasks should succeed for both generated roles.
- Inventory and vars should remain merged across repeated `jtaf-xml2yaml` runs.

At this point you have completed render -> preview diff -> push -> plugin-level verification using both generated roles (QFX and SRX).

---

## Reference

### Generate from an existing JSON model

If you already have pyang JSON, use `jtaf-ansible` instead of `jtaf-yang2ansible`:

```bash
jtaf-ansible -j junos.json \
  -x examples/evpn-vxlan-dc/dc1/*spine*.xml -t vqfx
```

`-j` accepts plain or gzipped JSON, or `-` for stdin. Multiple XML files must
match the same device type. The output contains `roles/<type>_role/`,
`filter_plugins/`, `jtaf-playbook.yml`, `configs/`, and `trimmed_schema.json.gz`.
The generated playbook includes device-push tasks: `connection: local` means
modules run on the control node, not that they cannot connect to devices.
Use the operator-owned playbook above to make deployment choices explicit.

### Inventory and variable hierarchy

`jtaf-xml2yaml` requires a grouping definition. `--grouping-hosts-file` supplies
inventory group names and optional `:children` relationships; `--hosts-file`
selects the inventory output. `-j` must reference the matching generated role
schema, `-x` supplies device XML, and `-d` selects the variable output directory.

The role loads global, group, and host values from `jtaf_vars_root` and merges
them in that order. The generated filters handle merge directives. Reusing
the same output directory merges tracked inventory and variable data rather
than discarding other devices. Review generated data after each regeneration.

### How deployment works

Jinja2 renders XML from the merged variables. `juniper.device.config` previews
the candidate diff and pushes XML through NETCONF. This is distinct from the
Terraform provider's schema-aware leaf-level patch engine.
Check module results and device reachability; a successful module result alone
does not prove that the network is operating correctly.

### Override mode

The default `group` mode wraps configuration in a `JTAF_ANSIBLE` group and uses
`load replace`. `override` mode renders bare configuration and requires
`load override`, replacing the entire device configuration.

To generate an override role, add `--mode override` to `jtaf-ansible` or
`jtaf-yang2ansible`. **Do not use the group-mode deployment playbook above
unchanged:** its push tasks use `load: replace`. Change the load operation to
`override` and use a commit-confirmed safety workflow with reachability checks
and confirmation before the timer expires.

Override XML must retain the management interface address, root authentication,
and NETCONF SSH service or you may lose access. All unmanaged configuration is
deleted. Choose a confirmation timer long enough to verify the device; the
generated role's `jtaf_commit_confirm_minutes` default is 2 minutes.
