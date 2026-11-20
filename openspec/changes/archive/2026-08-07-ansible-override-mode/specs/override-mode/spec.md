## ADDED Requirements

### Requirement: Override template renders bare configuration XML
The Ansible Jinja2 template in override mode SHALL render `<configuration>...</configuration>` without any `<groups>` or `<apply-groups>` wrapper. The rendered XML SHALL contain all configuration subtrees covered by the trimmed schema, directly under `<configuration>`.

#### Scenario: Override template produces bare XML
- **WHEN** `jtaf_mode` is set to `override` and the role renders the template
- **THEN** the output XML starts with `<configuration>` and contains no `<groups>` or `<apply-groups>` elements

#### Scenario: Group template still wraps in groups
- **WHEN** `jtaf_mode` is set to `group` (or unset, defaulting to `group`) and the role renders the template
- **THEN** the output XML wraps all config in `<groups><name>JTAF_ANSIBLE</name>...</groups>` with `<apply-groups>JTAF_ANSIBLE</apply-groups>`, identical to current behavior

### Requirement: Override mode uses load override NETCONF operation
The generated playbook in override mode SHALL use `load: override` when pushing configuration to the device. This replaces the entire device configuration with the rendered XML.

#### Scenario: Load override replaces full device config
- **WHEN** the playbook runs in override mode against a device
- **THEN** the `juniper.device.config` task uses `load: override` (not `load: replace` or `load: merge`)

### Requirement: Override mode uses commit confirmed safety pattern
The generated playbook in override mode SHALL use `commit confirmed <N>` where N is a configurable timeout in minutes. After committing, the playbook SHALL verify device reachability, and upon successful verification SHALL send a second commit to clear the confirmed timer.

#### Scenario: Successful override with commit confirmed
- **WHEN** the playbook applies an override and the device remains reachable
- **THEN** the playbook sends `commit confirmed <N>`, waits, verifies NETCONF reachability on port 830, then sends a confirming `commit` to clear the timer

#### Scenario: Override causes device unreachability
- **WHEN** the playbook applies an override and the device becomes unreachable (e.g., management IP was removed from the rendered config)
- **THEN** the `commit confirmed` timer expires on the device and Junos auto-rolls back to the previous configuration

#### Scenario: Commit confirm timeout is configurable
- **WHEN** the operator sets `jtaf_commit_confirm_minutes` variable (default: 2)
- **THEN** the playbook uses that value as the `confirmed:` parameter in the `juniper.device.config` task

### Requirement: Unmodeled paths are silently omitted
When XML input contains configuration paths that do not exist in the YANG schema (trimmed_schema.json), those paths SHALL NOT be included in the generated Ansible role template or variables. This matches Terraform provider behavior.

#### Scenario: XML config has paths not in YANG
- **WHEN** input XML contains `<extension-service>` or other elements not modeled in the YANG files
- **THEN** those elements are silently skipped during role/template generation — no error, no warning, no output for those paths

### Requirement: Continuous updates via YAML edits
When an operator edits `host_vars` or `group_vars` YAML files and re-runs the playbook in override mode, the changes SHALL be automatically reflected. No intermediate tools or regeneration steps are needed — the Jinja2 template reads YAML vars at render time, and Junos candidate config computes the diff internally.

#### Scenario: Operator changes a YAML variable and re-runs
- **WHEN** an operator edits `host_vars/dc1-leaf1.yaml` to change an interface description and re-runs `ansible-playbook site.yml`
- **THEN** the template renders the updated config, `load override` loads it into the candidate, Junos computes the diff, and only the changed line is committed
