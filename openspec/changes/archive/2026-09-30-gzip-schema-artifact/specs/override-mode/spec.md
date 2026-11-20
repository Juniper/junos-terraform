# Spec Delta

## MODIFIED Requirements

### Requirement: Unmodeled paths are silently omitted
When XML input contains configuration paths that do not exist in the YANG schema (the trimmed schema given with `-j`, plain or gzipped), those paths SHALL NOT be included in the generated Ansible role template or variables. This matches Terraform provider behavior.

#### Scenario: XML config has paths not in YANG
- **WHEN** input XML contains `<extension-service>` or other elements not modeled in the YANG files
- **THEN** those elements are silently skipped during role/template generation — no error, no warning, no output for those paths

#### Scenario: Schema given gzipped
- **WHEN** an override-mode role is generated with `jtaf-ansible --mode override -j SCHEMA.json.gz -x CONFIG -t vqfx` and its variables with `jtaf-xml2yaml -j ansible-provider-junos-vqfx/trimmed_schema.json.gz -x CONFIG -d OUT`
- **THEN** unmodeled paths are omitted exactly as when the same schema is given as plain JSON
