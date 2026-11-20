## ADDED Requirements

### Requirement: Mode flag on CLI tools
The `jtaf-ansible` and `jtaf-yang2ansible` CLI tools SHALL accept a `--mode` flag with values `group` (default) or `override`. This flag controls which template variant is used and what playbook pattern is generated.

#### Scenario: jtaf-ansible with --mode override
- **WHEN** the operator runs `jtaf-ansible -j schema.json -x config.xml -t vqfx --mode override`
- **THEN** the generated role uses the override template (no `<groups>` wrapper) and the generated playbook includes `load: override` with `commit confirmed` tasks

#### Scenario: jtaf-ansible without --mode flag
- **WHEN** the operator runs `jtaf-ansible -j schema.json -x config.xml -t vqfx` (no `--mode` flag)
- **THEN** the generated role uses the group template (current behavior) and the generated playbook uses `load: replace`

#### Scenario: jtaf-yang2ansible passes mode through
- **WHEN** the operator runs `jtaf-yang2ansible -p <yang> ... -x <xml> -t vqfx --mode override`
- **THEN** the `--mode` flag is passed through to the underlying `jtaf-ansible` invocation

### Requirement: Mode variable in role vars
The generated role SHALL include a `jtaf_mode` variable that can be set in the role defaults, group_vars, or host_vars. This variable SHALL control template selection at render time.

#### Scenario: jtaf_mode set in group_vars
- **WHEN** `jtaf_mode: "override"` is set in `group_vars/all.yaml`
- **THEN** all hosts in that inventory render the override template variant

#### Scenario: jtaf_mode defaults to group
- **WHEN** `jtaf_mode` is not set anywhere in the variable hierarchy
- **THEN** the role defaults to `jtaf_mode: "group"` and uses the groups-wrapped template

### Requirement: Mode determines template selection
The role's `tasks/main.yml` SHALL select the appropriate template based on the `jtaf_mode` variable.

#### Scenario: Override mode selects override template
- **WHEN** `jtaf_mode` is `"override"`
- **THEN** the task uses `template-override.j2` (or the conditional path within a single template) to render bare `<configuration>` XML

#### Scenario: Group mode selects group template
- **WHEN** `jtaf_mode` is `"group"`
- **THEN** the task uses `template.j2` to render `<configuration><groups>...</groups><apply-groups>...</apply-groups></configuration>` XML

### Requirement: Mode determines playbook push strategy
The generated or operator-authored playbook SHALL use different push strategies based on the mode.

#### Scenario: Override mode playbook tasks
- **WHEN** the playbook detects `jtaf_mode == "override"` for a host
- **THEN** the push tasks use `load: override`, `confirmed: {{ jtaf_commit_confirm_minutes }}`, followed by reachability verification and commit confirmation

#### Scenario: Group mode playbook tasks
- **WHEN** the playbook detects `jtaf_mode == "group"` for a host
- **THEN** the push tasks use `load: replace`, `commit: true` (immediate, no confirmed timer)
