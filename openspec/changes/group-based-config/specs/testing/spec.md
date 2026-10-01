# Spec Delta

## ADDED Requirements

### Requirement: A group-based example is maintained and exercised
The examples SHALL include a group-based variant that sits alongside the base-configuration example in its own
directory, covering the same devices, with the bulk of each device's configuration inside a configuration group and
only device-specific values in the base hierarchy. The variant SHALL be buildable and convertible with the documented
commands and SHALL be exercised end to end: build the provider with groups, convert the XML to Terraform
configuration, apply, and confirm a second plan reports no changes.

#### Scenario: Group example builds and converts
- **WHEN** the documented build and convert commands are run for the group-based example
- **THEN** a provider is produced whose schema contains `groups`, and the conversion yields Terraform configuration
  with a `groups` block and an `apply_groups` list for each device

#### Scenario: Group example applies idempotently
- **WHEN** the group-based example is applied to a device holding no prior configuration and then planned again
- **THEN** the first apply succeeds and the second plan reports no changes

#### Scenario: Base example is unaffected
- **WHEN** the existing base-configuration example is built, converted and applied
- **THEN** its behaviour and output are unchanged by the presence of the group-based variant

### Requirement: Default and groups builds are both covered
The test suite SHALL cover provider generation both with and without the groups option, verifying that the default
build omits `groups` and `apply-groups` from the schema and resource, and that the groups build includes them.

#### Scenario: Default build asserted
- **WHEN** the generation tests run without the groups option
- **THEN** they assert that neither `groups` nor `apply-groups` appears in the generated schema or resource schema

#### Scenario: Groups build asserted
- **WHEN** the generation tests run with the groups option
- **THEN** they assert that both appear, with `groups` keyed by `name` and `apply-groups` as an ordered leaf-list
