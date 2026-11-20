## Context

The JTAF Ansible workflow generates Jinja2 templates and roles from YANG models + XML configs. Currently, all generated templates wrap configuration inside a `<groups><name>JTAF_ANSIBLE</name>...</groups>` block and push with `load replace`. This limits operators to partial device management.

On the Terraform side, the `netconf_patch` branch already moved away from groups — the provider pushes config directly into the base configuration and uses a Go patch engine for minimal `edit-config` operations. The Ansible side needs a similar "direct config" capability, but with a different mechanism: Junos's `load override` + candidate config IS the diff engine (no patch engine needed), and `commit confirmed` IS the rollback safety net.

Key constraint: the `juniper.device.config` Ansible module already supports `load: override` and `confirmed: <minutes>` natively. No new Ansible modules are needed.

## Goals / Non-Goals

**Goals:**
- Add an `override` operating mode that pushes bare `<configuration>` with `load override` and `commit confirmed` safety
- Keep `group` mode as the default for backward compatibility
- Provide a mode-switching mechanism (CLI flag + role variable) so operators choose per-deployment
- Ensure continuous YAML edits by operators flow into the playbook without regeneration — Junos handles the diff

**Non-Goals:**
- Building a patch engine for Ansible (Junos candidate config handles diffing)
- Changing the Terraform provider or its patch engine
- Supporting per-host mixed modes in a single playbook run (all hosts in a play use the same mode)
- Implementing a YAML-level diff tool for Ansible (not needed — `load override` makes Junos the diff engine)

## Decisions

### Decision 1: Single template with conditional vs. two separate templates

**Choice:** Single template (`ansible.j2`) with a Jinja2 conditional at the top/bottom.

**Rationale:** The override and group templates are identical except for the 4-line `<groups>` wrapper. Maintaining two nearly identical templates creates a maintenance burden. A simple `{% if jtaf_mode == 'override' %}` conditional keeps everything in one file.

**Alternative considered:** Two separate templates (`ansible.j2` and `ansible-override.j2`). Rejected because the diff is only 4 lines, and having two files means every future template change needs to be applied twice.

**Implementation:**
```jinja2
{# at the top, after rendering starts #}
{% if jtaf_mode | default('group') == 'override' %}
<configuration>
{% else %}
<configuration>
    <groups>
        <name>{{ jtaf_group_name | default('JTAF_ANSIBLE') }}</name>
{% endif %}

{# ... all the recursive rendering (identical) ... #}

{% if jtaf_mode | default('group') == 'override' %}
</configuration>
{% else %}
    </groups>
    <apply-groups>{{ jtaf_group_name | default('JTAF_ANSIBLE') }}</apply-groups>
</configuration>
{% endif %}
```

### Decision 2: Mode flag mechanism

**Choice:** Three-layer mechanism:
1. CLI flag `--mode group|override` on `jtaf-ansible` / `jtaf-yang2ansible` — controls what gets generated
2. Role defaults `jtaf_mode: "group"` in generated `defaults/main.yml` — runtime default
3. Operator can override via `group_vars` or `host_vars` YAML — runtime override

**Rationale:** CLI flag controls generation-time behavior (which playbook pattern to scaffold). Role variable controls runtime behavior (which template path to render). This separation means a single generated role can be switched between modes without regeneration.

### Decision 3: Playbook structure for override mode

**Choice:** Generate a `site.yml` with conditional task blocks based on `jtaf_mode`.

**Implementation pattern:**
```yaml
- name: Push config (override mode)
  juniper.device.config:
    load: override
    src: "{{ tmp_dir }}/{{ inventory_hostname }}.xml"
    confirmed: "{{ jtaf_commit_confirm_minutes | default(2) }}"
    diff: true
  when: jtaf_mode | default('group') == 'override'

- name: Push config (group mode)
  juniper.device.config:
    load: replace
    src: "{{ tmp_dir }}/{{ inventory_hostname }}.xml"
  when: jtaf_mode | default('group') == 'group'

- name: Verify reachability (override only)
  wait_for:
    host: "{{ ansible_host }}"
    port: "{{ ansible_port | default(830) }}"
    timeout: "{{ (jtaf_commit_confirm_minutes | default(2)) * 60 - 10 }}"
  when: jtaf_mode | default('group') == 'override'

- name: Confirm commit (override only)
  juniper.device.config:
    check: true
    diff: false
    commit: false
  when: jtaf_mode | default('group') == 'override'
```

### Decision 4: Schema performance for override mode

**Choice:** Use the existing `trimmed_schema.json` as-is. If performance becomes an issue with very large schemas (full-model parse), add a one-time index build step that creates an in-memory `dict[path] → node_info` lookup (same pattern as the Go patch engine's `ProcessSchema`).

**Rationale:** The trimmed schema is already used by `jtaf-xml2yaml` for variable extraction. For override mode, the rendering pipeline is the same — the schema is read once at generation time, not at every playbook run. Runtime performance is determined by Jinja2 template rendering speed (reading YAML vars, producing XML), which is independent of schema size.

The concern about "timing issues with large trimmed schemas" applies only during `jtaf-xml2yaml` extraction and `jtaf-ansible` template generation — both are one-time operations, not per-playbook-run. At playbook runtime, only the YAML vars and Jinja2 template are involved; the trimmed schema is not consulted.

## Risks / Trade-offs

**[Override mode deletes unmanaged config]** → Mitigation: This is by design. The operator accepts that what they provide IS the entire config. Documentation MUST clearly state that override mode requires the rendered XML to include management IP, root authentication, and NETCONF SSH service. A pre-flight check in the playbook can validate these critical paths exist in the rendered XML before pushing.

**[Large schema produces very large template]** → Mitigation: Full-model parse produces templates with all possible YANG paths, but Jinja2 conditionals (`{% if ... is defined %}`) skip undefined paths at render time. The template file itself may be large, but rendering is fast because most paths are skipped. If the template file size becomes problematic, investigate template compilation/caching.

**[Commit confirmed timer too short]** → Mitigation: Default to 2 minutes (configurable via `jtaf_commit_confirm_minutes`). Document that operators should set this based on their commit time expectations. For large configs, 5 minutes may be needed.


**[Backward compatibility]** → Mitigation: Group mode is the default everywhere. No existing behavior changes unless the operator explicitly sets `--mode override` or `jtaf_mode: "override"`.
