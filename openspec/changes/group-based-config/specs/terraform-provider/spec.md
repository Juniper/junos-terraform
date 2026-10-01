# Spec Delta

## ADDED Requirements

### Requirement: Resource attributes follow the embedded schema
The resource SHALL expose exactly the configuration nodes present in its embedded schema, with no node name excluded
by the provider itself. Where `groups` and `apply-groups` are present they SHALL be exposed like any other list and
leaf-list; where they are absent the resource SHALL have no corresponding attributes.

#### Scenario: Groups present in schema
- **WHEN** the embedded schema contains `groups` and `apply-groups`
- **THEN** the resource exposes `groups` and `apply_groups`, and values written to them are sent to the device

#### Scenario: Groups absent from schema
- **WHEN** the embedded schema does not contain `groups`
- **THEN** the resource has no `groups` attribute and the provider reports no error about the missing node

### Requirement: Group state is read back from the device
Reading device state SHALL return configuration below `groups` and the `apply-groups` references with the same
fidelity as base-hierarchy configuration, so that out-of-band changes inside a group are detected as drift and
reconciled. Groups and references that the configuration does not manage SHALL be left untouched.

#### Scenario: Drift inside a group
- **WHEN** a leaf inside a managed group is changed on the device out of band
- **THEN** the next plan reports that leaf as changed and the next apply restores it

#### Scenario: Unmanaged groups are preserved
- **WHEN** the device holds a group that the configuration does not declare
- **THEN** applying the configuration leaves that group and any reference to it unchanged

### Requirement: Group references are not synthesised
The provider SHALL NOT maintain `apply-groups` outside the configuration: it SHALL NOT append a reference for a group
it writes, SHALL NOT sort references, and SHALL NOT carry references between resources.

#### Scenario: Group without a reference
- **WHEN** a configuration declares a group but does not list it in `apply_groups`
- **THEN** the committed device configuration contains the group body and no reference to it
