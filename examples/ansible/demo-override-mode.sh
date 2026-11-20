#!/bin/bash
# =============================================================================
# JTAF Ansible Override Mode — End-to-End Demo
# =============================================================================
#
# This script demonstrates the full workflow for JTAF's Ansible override mode:
#   1. Generate roles (override mode)
#   2. Generate host_vars/group_vars
#   3. Dry-run preview (no commit)
#   4. Apply with commit confirmed safety
#   5. Make a change and re-apply
#   6. Verify on device
#
# Prerequisites:
#   - junos-terraform repo checked out on ansible-override-mode branch
#   - pip install -e . (JTAF tools installed)
#   - pip install ansible junos-eznc jxmlease ncclient
#   - ansible-galaxy collection install juniper.device
#   - Network access to Junos devices on port 830
#
# Usage:
#   cd junos-terraform/examples/ansible
#   export NETCONF_USERNAME='{user}'
#   export NETCONF_PASSWORD='{password}'
#   # Then run each phase manually (copy/paste sections)
#
# =============================================================================

set -e

# =============================================================================
# CONFIGURATION — Update these for your environment
# =============================================================================

NETCONF_USERNAME="${NETCONF_USERNAME:-{user}}"
NETCONF_PASSWORD="${NETCONF_PASSWORD:-{password}}"
YANG_COMMON="../yang/18.2/18.2R3/common"
YANG_QFX="../yang/18.2/18.2R3/junos-qfx/conf/*.yang"
DEVICE_TYPE="vqfx-override"
ROLE_NAME="${DEVICE_TYPE}_role"

export NETCONF_USERNAME NETCONF_PASSWORD

echo "============================================="
echo " JTAF Ansible Override Mode Demo"
echo "============================================="
echo ""

# =============================================================================
# PHASE 1: Generate Ansible Role (Override Mode)
# =============================================================================

echo ">>> PHASE 1: Generating Ansible role with --mode override..."
echo ""

# Clean previous output
rm -rf ansible-provider-junos-${DEVICE_TYPE}

# Generate the role
jtaf-yang2ansible \
  -p ${YANG_COMMON} ${YANG_QFX} \
  -x ../evpn-vxlan-dc/dc1/dc1-*leaf* ../evpn-vxlan-dc/dc1/dc1-*spine* \
     ../evpn-vxlan-dc/dc1/dc1-*borderleaf* ../evpn-vxlan-dc/dc2/dc2-*spine* \
  -t ${DEVICE_TYPE} \
  --mode override

echo ""
echo "--- Verifying override mode is set ---"
echo "defaults/main.yml:"
cat ansible-provider-junos-${DEVICE_TYPE}/roles/${ROLE_NAME}/defaults/main.yml
echo ""
echo "Template <groups> count (should be 0):"
grep -c "<groups>" ansible-provider-junos-${DEVICE_TYPE}/roles/${ROLE_NAME}/templates/template.j2 || echo "0"
echo ""
echo ">>> PHASE 1 COMPLETE: Role generated in ansible-provider-junos-${DEVICE_TYPE}/"
echo ""

# =============================================================================
# PHASE 2: Generate host_vars / group_vars
# =============================================================================

echo ">>> PHASE 2: Generating host_vars and group_vars..."
echo ""

rm -rf ansible_override_files

jtaf-xml2yaml \
  -x ../evpn-vxlan-dc/dc1/dc1-*leaf* ../evpn-vxlan-dc/dc1/dc1-*spine* \
     ../evpn-vxlan-dc/dc1/dc1-*borderleaf* ../evpn-vxlan-dc/dc2/dc2-*spine* \
  -j ansible-provider-junos-${DEVICE_TYPE}/trimmed_schema.json.gz \
  -d ansible_override_files \
  --grouping-hosts-file switches_grouping_hosts

echo ""
echo "--- Generated host_vars ---"
ls ansible_override_files/host_vars/
echo ""
echo "--- Generated group_vars ---"
ls ansible_override_files/group_vars/
echo ""
echo ">>> PHASE 2 COMPLETE"
echo ""

# =============================================================================
# PHASE 3: Setup — ansible.cfg, inventory, playbooks
# =============================================================================

echo ">>> PHASE 3: Creating ansible.cfg, inventory, and playbooks..."
echo ""

# ansible.cfg
cat > ansible.cfg << 'EOF'
[defaults]
roles_path = ansible-provider-junos-vqfx-override/roles
filter_plugins = ansible-provider-junos-vqfx-override/filter_plugins
host_key_checking = False
interpreter_python = auto_silent
EOF

# Single-device inventory for testing
cat > test-inventory.ini << 'EOF'
[dc1-spine]
dc1-spine1 ansible_host=100.123.24.3 ansible_port=830
EOF

# Full inventory (for later expansion)
cat > full-inventory.ini << 'EOF'
[dc1-borderleaf]
dc1-borderleaf1 ansible_host=100.123.24.1 ansible_port=830
dc1-borderleaf2 ansible_host=100.123.24.2 ansible_port=830

[dc1-spine]
dc1-spine1 ansible_host=100.123.24.3 ansible_port=830
dc1-spine2 ansible_host=100.123.24.4 ansible_port=830

[dc1-leaf]
dc1-leaf1 ansible_host=100.123.24.5 ansible_port=830
dc1-leaf2 ansible_host=100.123.24.6 ansible_port=830
dc1-leaf3 ansible_host=100.123.24.7 ansible_port=830

[dc2-spine]
dc2-spine1 ansible_host=100.123.24.8 ansible_port=830
dc2-spine2 ansible_host=100.123.24.9 ansible_port=830
EOF

# Dry-run playbook (preview only, no commit)
cat > test-override-dryrun.yml << 'EOF'
---
- name: Render config (override mode)
  hosts: all
  connection: local
  gather_facts: false
  vars:
    tmp_dir: ansible-provider-junos-vqfx-override/configs
    jtaf_vars_root: ansible_override_files
    jtaf_mode: "override"
  roles:
    - role: vqfx-override_role
      delegate_to: localhost

- name: Preview override diff (NO commit)
  hosts: all
  connection: local
  gather_facts: false
  vars:
    netconf_user: "{{ lookup('env', 'NETCONF_USERNAME') }}"
    netconf_pass: "{{ lookup('env', 'NETCONF_PASSWORD') }}"
    tmp_dir: ansible-provider-junos-vqfx-override/configs
  tasks:
    - name: Load override (candidate only, NO commit)
      juniper.device.config:
        host: "{{ ansible_host }}"
        port: "{{ ansible_port | default(830) }}"
        user: "{{ netconf_user }}"
        passwd: "{{ netconf_pass }}"
        load: override
        src: "{{ tmp_dir }}/{{ inventory_hostname }}.xml"
        check: false
        commit: false
        diff: true
      register: preview

    - name: Show what would change
      debug:
        var: preview.diff_lines
      when: preview.diff_lines | default([]) | length > 0

    - name: No changes needed
      debug:
        msg: "No diff — device config matches rendered template exactly."
      when: preview.diff_lines | default([]) | length == 0
EOF

# Apply playbook (commit confirmed)
cat > test-override-apply.yml << 'EOF'
---
- name: Render config (override mode)
  hosts: all
  connection: local
  gather_facts: false
  vars:
    tmp_dir: ansible-provider-junos-vqfx-override/configs
    jtaf_vars_root: ansible_override_files
    jtaf_mode: "override"
  roles:
    - role: vqfx-override_role
      delegate_to: localhost

- name: Apply override with commit confirmed
  hosts: all
  connection: local
  gather_facts: false
  vars:
    netconf_user: "{{ lookup('env', 'NETCONF_USERNAME') }}"
    netconf_pass: "{{ lookup('env', 'NETCONF_PASSWORD') }}"
    tmp_dir: ansible-provider-junos-vqfx-override/configs
    jtaf_commit_confirm_minutes: 5
  tasks:
    - name: Load override + commit confirmed
      juniper.device.config:
        host: "{{ ansible_host }}"
        port: "{{ ansible_port | default(830) }}"
        user: "{{ netconf_user }}"
        passwd: "{{ netconf_pass }}"
        load: override
        src: "{{ tmp_dir }}/{{ inventory_hostname }}.xml"
        confirmed: "{{ jtaf_commit_confirm_minutes }}"
        diff: true
        comment: "JTAF override mode"
      register: override_result

    - name: Show diff applied
      debug:
        var: override_result.diff_lines
      when: override_result.diff_lines | default([]) | length > 0

    - name: Verify device still reachable (port 830)
      wait_for:
        host: "{{ ansible_host }}"
        port: "{{ ansible_port | default(830) }}"
        timeout: 280

    - name: Confirm commit (clear rollback timer)
      juniper.device.config:
        host: "{{ ansible_host }}"
        port: "{{ ansible_port | default(830) }}"
        user: "{{ netconf_user }}"
        passwd: "{{ netconf_pass }}"
        check: true
        diff: false
        commit: false

    - name: Success
      debug:
        msg: "Override applied and confirmed on {{ inventory_hostname }}"
EOF

echo "Created: ansible.cfg, test-inventory.ini, full-inventory.ini"
echo "Created: test-override-dryrun.yml, test-override-apply.yml"
echo ""
echo ">>> PHASE 3 COMPLETE"
echo ""

# =============================================================================
# PHASE 4: Dry-Run (Preview Only — No Commit)
# =============================================================================

echo ">>> PHASE 4: Running dry-run against dc1-spine1 (preview only, no commit)..."
echo ""
echo "Command: ansible-playbook -i test-inventory.ini test-override-dryrun.yml"
echo ""

ansible-playbook -i test-inventory.ini test-override-dryrun.yml

echo ""
echo ">>> PHASE 4 COMPLETE: Review the diff above. No changes were committed."
echo ""

# =============================================================================
# PHASE 5: Apply Override with Commit Confirmed
# =============================================================================

echo ">>> PHASE 5: Applying override with commit confirmed (5 min safety timer)..."
echo ""
echo "Command: ansible-playbook -i test-inventory.ini test-override-apply.yml"
echo ""
echo "What happens:"
echo "  1. load override → replaces entire device config with rendered XML"
echo "  2. commit confirmed 5 → Junos auto-rolls back in 5 min if unreachable"
echo "  3. wait_for → verify device is still reachable on port 830"
echo "  4. commit → confirms the change, clears the rollback timer"
echo ""

ansible-playbook -i test-inventory.ini test-override-apply.yml

echo ""
echo ">>> PHASE 5 COMPLETE: Override applied and confirmed!"
echo ""

# =============================================================================
# PHASE 6: Make a Change and Re-Apply (Day 2 Operation)
# =============================================================================

echo ">>> PHASE 6: Making a Day-2 change..."
echo ""
echo "Changing dc1-spine1 SNMP contact from 'aburston@juniper.net' to 'jtaf-demo@juniper.net'"
echo ""

# Edit the host_vars to make a change
if grep -q "aburston@juniper.net" ansible_override_files/host_vars/dc1-spine1.yaml 2>/dev/null; then
  sed -i 's/aburston@juniper.net/jtaf-demo@juniper.net/' ansible_override_files/host_vars/dc1-spine1.yaml
  echo "Changed in host_vars/dc1-spine1.yaml"
elif grep -q "aburston@juniper.net" ansible_override_files/group_vars/all.yaml 2>/dev/null; then
  # If it's in group_vars, we need to add a host-specific override
  echo "  snmp:" >> ansible_override_files/host_vars/dc1-spine1.yaml
  echo "    contact: jtaf-demo@juniper.net" >> ansible_override_files/host_vars/dc1-spine1.yaml
  echo "Added override in host_vars/dc1-spine1.yaml"
else
  echo "SNMP contact not found — adding to host_vars"
  echo "snmp:" >> ansible_override_files/host_vars/dc1-spine1.yaml
  echo "  contact: jtaf-demo@juniper.net" >> ansible_override_files/host_vars/dc1-spine1.yaml
fi

echo ""
echo "--- Dry-run to see just the SNMP change ---"
ansible-playbook -i test-inventory.ini test-override-dryrun.yml

echo ""
echo "--- Applying the change ---"
ansible-playbook -i test-inventory.ini test-override-apply.yml

echo ""
echo ">>> PHASE 6 COMPLETE: Day-2 change applied!"
echo ""

# =============================================================================
# PHASE 7: Verify on Device
# =============================================================================

echo ">>> PHASE 7: Verify on device"
echo ""
echo "Run these commands on the device to verify:"
echo ""
echo "  ssh {user}@100.123.24.3"
echo "  show system commit"
echo "  show configuration | compare rollback 1"
echo "  show snmp contact"
echo "  show configuration snmp"
echo ""
echo "============================================="
echo " DEMO COMPLETE"
echo "============================================="
echo ""
echo "Summary of what was demonstrated:"
echo "  1. Generated Ansible role with --mode override (no <groups> wrapper)"
echo "  2. Generated host_vars/group_vars from XML configs"
echo "  3. Previewed config diff without committing (dry-run)"
echo "  4. Applied full config with load override + commit confirmed"
echo "  5. Made a Day-2 YAML change → re-ran playbook → only diff was applied"
echo "  6. Junos calculated the diff internally (no patch engine needed)"
echo "  7. commit confirmed provided auto-rollback safety"
echo ""
echo "Key points:"
echo "  - Override mode: what you provide IS the entire device config"
echo "  - Junos candidate config is the diff engine"
echo "  - commit confirmed = auto-rollback if device becomes unreachable"
echo "  - Edit YAML → re-run playbook → Junos figures out what changed"
echo ""
