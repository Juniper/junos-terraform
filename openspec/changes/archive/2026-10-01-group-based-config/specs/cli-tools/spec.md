# Spec Delta

## ADDED Requirements

### Requirement: Groups option on the Terraform tools
`jtaf-provider` SHALL accept a `--groups` flag that keeps the `groups` subtree and the `apply-groups` leaf-list in the
generated schema. `jtaf-yang2go` SHALL accept the same flag and pass it to `jtaf-provider`. `jtaf-xml2tf` SHALL accept
a `--groups` flag that preserves the group hierarchy instead of flattening applied groups into the base configuration.
Omitting the flag SHALL preserve today's behaviour in every tool. The Ansible tools SHALL NOT accept the flag.

#### Scenario: Flag reaches the generator
- **WHEN** `jtaf-yang2go --groups` runs
- **THEN** it invokes `jtaf-provider` with `--groups`, and the resulting `trimmed_schema.json.gz` contains the `groups`
  subtree and the `apply-groups` leaf-list

#### Scenario: Flag is absent
- **WHEN** any of the three tools runs without `--groups`
- **THEN** its output is byte-identical to the output of the same invocation before this change

#### Scenario: Conversion preserves groups
- **WHEN** `jtaf-xml2tf --groups` converts an XML configuration containing `<groups>` and `<apply-groups>`
- **THEN** the `.tf` output contains a `groups` block and an `apply_groups` list rather than the flattened result

### Requirement: Schema scope flag
`jtaf-provider --generic` SHALL select the scope of the embedded schema rather than a provider implementation: with
`--generic` the untrimmed model is embedded, and without it the model is trimmed to the XML given with `-x`. The two
invocations SHALL produce the same provider source and SHALL differ only in the embedded schema.

#### Scenario: Same source, different schema
- **WHEN** a provider is generated with `--generic` and another is generated with `-x`
- **THEN** their Go source files are identical and only the embedded schema differs

#### Scenario: Trimming still requires XML
- **WHEN** `--generic` is combined with `-x`
- **THEN** the tool exits with a non-zero status and an error stating the two are mutually exclusive

## MODIFIED Requirements

### Requirement: Provider and role generators write compressed schemas
`jtaf-provider` and `jtaf-ansible` SHALL write their trimmed schema as compact gzip-compressed JSON named
`trimmed_schema.json.gz` and SHALL NOT write `trimmed_schema.json`. This applies whether the schema is trimmed to XML
or embedded untrimmed.

#### Scenario: Generated schema handoff
- **WHEN** either generator completes
- **THEN** its output directory contains `trimmed_schema.json.gz`, and downstream tools can consume that file

#### Scenario: Schema help text
- **WHEN** a user checks the `-j` help text for any of the four tools
- **THEN** it names `trimmed_schema.json.gz` and states that plain JSON is also accepted
