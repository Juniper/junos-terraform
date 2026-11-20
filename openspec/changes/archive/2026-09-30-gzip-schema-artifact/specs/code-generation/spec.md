# Spec Delta

## ADDED Requirements

### Requirement: Generated providers carry a gzipped schema
A generated Terraform provider directory SHALL contain the trimmed schema as `trimmed_schema.json.gz`: the pyang JSON (filtered by `-x` if given, with `--exclude` and the top-level `version` applied), serialised compactly and gzip-compressed. No plain `trimmed_schema.json` SHALL be written. This holds for both generation paths, Jinja2 and `--generic`.

#### Scenario: Jinja2 provider generation
- **WHEN** `jtaf-provider -j SCHEMA -x CONFIG -t vqfx` runs without `--generic`
- **THEN** `terraform-provider-junos-vqfx/trimmed_schema.json.gz` exists, decompresses to the filtered schema as a JSON object with a `root` node, and `terraform-provider-junos-vqfx/trimmed_schema.json` does not exist

#### Scenario: Generic provider generation
- **WHEN** `jtaf-provider --generic -j SCHEMA -t srx` runs
- **THEN** `terraform-provider-junos-srx/trimmed_schema.json.gz` exists alongside `schema.bin.gz`, and no plain `trimmed_schema.json` exists

#### Scenario: Downstream tool reads the generated schema
- **WHEN** `jtaf-xml2tf -j terraform-provider-junos-vqfx/trimmed_schema.json.gz -x CONFIG -t vqfx -d OUT` is run on the output of either generation path
- **THEN** the `.tf` files are generated exactly as they would be from the equivalent plain JSON schema
