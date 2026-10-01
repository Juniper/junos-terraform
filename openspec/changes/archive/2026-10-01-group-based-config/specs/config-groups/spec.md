# Spec Delta

## Purpose

Opt-in management of Junos configuration groups: the `groups` subtree and the `apply-groups` references that activate
it, carried through schema generation, the Terraform resource, XML-to-HCL conversion, and the NETCONF edit that is sent
to the device.

## ADDED Requirements

### Requirement: Group support is opt-in
Configuration groups SHALL be excluded unless explicitly requested. When groups are not requested, the generated
schema SHALL NOT contain the `groups` subtree or the `apply-groups` leaf-list, and the Terraform resource SHALL NOT
expose them. When groups are requested, both SHALL be present and SHALL be treated as ordinary configuration nodes.
Requesting groups SHALL be recorded in the generated artifacts so that the provider's behaviour follows from its
embedded schema rather than from a runtime setting.

#### Scenario: Default build hides groups
- **WHEN** a provider is generated without the groups option
- **THEN** the embedded schema and the resource schema contain neither `groups` nor `apply_groups`, and a `.tf` file
  that sets either attribute fails to plan with an unsupported-attribute error

#### Scenario: Groups build exposes groups
- **WHEN** a provider is generated with the groups option
- **THEN** the resource schema contains `groups` as a keyed nested block and `apply_groups` as a list of strings

#### Scenario: Groups option and subtree exclusion conflict
- **WHEN** the groups option is combined with an exclusion of the `groups` subtree
- **THEN** generation SHALL fail with a non-zero status and an error naming the conflict

### Requirement: Groups are modelled as Junos hierarchy
A group SHALL be expressed as an entry of the `groups` list keyed by its `name` leaf, with the group's configuration
nested inside it exactly as the Junos data model defines it. There SHALL be no provider-level, resource-level or
attribute-level setting that names a group; the group name SHALL come only from the hierarchy. Configuration below a
group SHALL use the same attribute names, nesting and types as the equivalent configuration in the base hierarchy.

#### Scenario: Group body mirrors base hierarchy
- **WHEN** the same interface configuration is written in the base hierarchy and inside a group entry
- **THEN** both use identical attribute names and nesting, differing only by the enclosing `groups` block

#### Scenario: Group name is the list key
- **WHEN** two entries of `groups` are declared
- **THEN** each is identified by its `name` value, and the provider sends each group body under its own
  `<groups><name>…</name></groups>` element

### Requirement: apply-groups is managed and ordered
`apply-groups` SHALL be managed as an order-preserving list of group names at the top of the configuration. The
provider SHALL send the list in the order declared, SHALL treat a reordering as a change, and SHALL remove references
that are no longer declared. The provider SHALL NOT add, reorder or retain `apply-groups` entries that the
configuration does not declare.

#### Scenario: Group reference is written with the group
- **WHEN** a configuration declares a group and lists that group in `apply_groups`
- **THEN** the committed device configuration contains both the group body and the `apply-groups` reference

#### Scenario: Reordering is a change
- **WHEN** only the order of `apply_groups` entries changes
- **THEN** the plan is non-empty and the applied configuration reflects the new order

#### Scenario: Removed reference is deleted
- **WHEN** an entry is removed from `apply_groups`
- **THEN** the provider deletes that reference from the device and leaves the remaining references in order

### Requirement: Round-trip of group-based configuration
Converting a Junos XML configuration that uses groups into Terraform configuration SHALL be able to preserve the group
hierarchy. Preservation SHALL be opt-in; by default the existing behaviour of flattening applied groups into the base
hierarchy SHALL be kept. A preserved conversion SHALL produce configuration that, when applied to a device holding the
original configuration, yields an empty plan.

#### Scenario: Preserved conversion keeps the group
- **WHEN** an XML configuration whose bulk sits inside a group is converted with group preservation requested
- **THEN** the generated Terraform configuration contains a `groups` block with that configuration and an
  `apply_groups` entry naming it, and nothing from the group is copied into the base hierarchy

#### Scenario: Default conversion still flattens
- **WHEN** the same XML is converted without requesting group preservation
- **THEN** the output is unchanged from the current flattened form

#### Scenario: Round-trip is a no-op
- **WHEN** the preserved conversion is applied to a device already holding the original configuration
- **THEN** the plan reports no changes
