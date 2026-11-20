## 1. Template Conditional (Override vs Group)

- [x] 1.1 Add `jtaf_mode` conditional to `junosterraform/templates/ansible.j2` — wrap the `<groups>` open/close and `<apply-groups>` lines in `{% if jtaf_mode | default('group') != 'override' %}` blocks so override mode renders bare `<configuration>`
- [x] 1.2 Add unit test: render the template with `jtaf_mode='override'` and verify output has no `<groups>` or `<apply-groups>` elements
- [x] 1.3 Add unit test: render the template with `jtaf_mode='group'` (or unset) and verify output matches current behavior (groups wrapper present)

## 2. CLI Flag on jtaf-ansible and jtaf-yang2ansible

- [x] 2.1 Add `--mode` argument to `junosterraform/jtaf-ansible` argparse (choices: `group`, `override`; default: `group`)
- [x] 2.2 Pass `jtaf_mode` value into the Jinja2 template rendering context in `jtaf-ansible`
- [x] 2.3 Add `--mode` argument to `junosterraform/jtaf-yang2ansible` argparse and pass it through to the `jtaf-ansible` subprocess call
- [x] 2.4 Generate `defaults/main.yml` in the role directory with `jtaf_mode: "<value>"` so the role has a default

## 3. Generated Playbook — Override Mode Tasks

- [x] 3.1 Update `jtaf-ansible` playbook generation to include conditional push tasks: `load: override` + `confirmed:` when `jtaf_mode == 'override'`, `load: replace` when `jtaf_mode == 'group'`
- [x] 3.2 Add reachability verification task (`wait_for` on port 830) for override mode, gated by `when: jtaf_mode == 'override'`
- [x] 3.3 Add commit confirmation task (second `juniper.device.config` call with `check: true, commit: false`) for override mode
- [x] 3.4 Add `jtaf_commit_confirm_minutes` variable (default: 2) to generated role defaults

## 4. Documentation

- [x] 5.1 Update `README-ansible.md` — add Override Mode section explaining the two modes, the flag, and the commit confirmed pattern
- [x] 5.3 Add override mode example to `examples/ansible/` with a sample `site.yml` showing the override + commit confirmed flow
- [x] 5.4 Update `README.md` main README

## 5. Testing

- [x] 5.1 Add Python unit tests for `jtaf-ansible` with `--mode override`: verify generated template, playbook, and role defaults
- [x] 5.2 Add Python unit tests for `jtaf-ansible` with `--mode group`: verify backward compatibility (identical to current output)
