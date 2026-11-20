## Why

The JTAF Ansible workflow currently wraps all generated configuration inside a Junos `<groups><name>JTAF_ANSIBLE</name>...</groups>` block and pushes it with `load replace`. This limits Ansible to partial device management — the operator can only control config within the JTAF group while other config coexists untouched. Many operators want full ownership of the device configuration through Ansible, where what they provide IS the entire box config. Additionally, there is no `commit confirmed` safety mechanism and no way for operators to continuously edit YAML vars and have changes flow back into the Ansible role and playbook automatically.

## What Changes

- **Two operating modes via a flag**: Introduce a `jtaf_mode` variable (set in role vars, group_vars, or host_vars) that switches between `group` mode (current behavior, default) and `override` mode (new: bare `<configuration>` with `load override` + `commit confirmed` safety).
- **New override template**: A variant of `ansible.j2` that renders `<configuration>` without the `<groups>` wrapper, producing a complete device configuration.
- **Commit confirmed safety pattern**: Override mode uses `commit confirmed <N>` → verify device reachability → `commit` (confirm). If the device becomes unreachable after override (e.g., management IP was removed), the timer expires and Junos auto-rolls back.
- **Continuous YAML-to-role updates**: When an operator edits `host_vars` or `group_vars` YAML files, re-running the playbook with override mode automatically picks up the changes — Junos candidate config computes the diff internally. No patch engine needed; Junos IS the diff engine.
- **Generated playbook variants**: `jtaf-ansible` and `jtaf-yang2ansible` gain a `--mode group|override` flag. Override mode generates a `site.yml` with the commit confirmed pattern. Group mode generates the current `load replace` playbook.
- **Config statements not in YANG**: Paths in the input XML that don't exist in the YANG schema are silently omitted from the generated role/template (same behavior as the Terraform provider). Users can manually add unmodeled config via raw XML passthrough if needed.

## Capabilities

### New Capabilities
- `override-mode`: The `load override` + `commit confirmed` operating mode — template without groups wrapper, playbook with commit confirmed safety pattern, mode flag in role variables.
- `mode-flag`: The mechanism to switch between `group` and `override` modes — variable in role vars, CLI flag on `jtaf-ansible`/`jtaf-yang2ansible`, conditional template rendering.

### Modified Capabilities
- None. Group mode (current default) remains unchanged and backward-compatible.

## Impact

- **Templates**: `ansible.j2` needs a conditional or a sibling `ansible-override.j2` template (4 lines differ — the `<groups>` wrapper).
- **CLI tools**: `jtaf-ansible`, `jtaf-yang2ansible` gain `--mode` flag.
- **Generated playbooks**: `jtaf-playbook.yml` and example `site.yml` conditionally include override+confirm tasks.
- **Generated roles**: `tasks/main.yml` unchanged — rendering logic is the same regardless of mode.
- **Schema/performance**: For override mode to be safe, the trimmed schema should cover all paths the user intends to manage. A schema index (in-memory lookup) may be needed if the full schema is too large for repeated comparisons.
- **No breaking changes**: All existing behavior preserved. Group mode remains the default.
- **Dependencies**: No new Python or Go dependencies. Uses existing `juniper.device.config` module which already supports `load: override` and `confirmed:`.
