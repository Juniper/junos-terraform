# Spec Delta

## ADDED Requirements

### Requirement: A device's configuration is one tree
The mock SHALL keep each device's committed and candidate configuration as a single configuration tree, not as a
separate payload per configuration group. A loaded payload SHALL be stored as it was sent, including any `groups` and
`apply-groups` it carries, and SHALL NOT be wrapped in, or keyed by, a group that the payload does not contain. A read
SHALL return the committed tree in the same shape.

#### Scenario: Groups alongside base hierarchy
- **WHEN** a configuration holding both a `groups` entry and base-hierarchy configuration is loaded and committed
- **THEN** a read returns both, with the group's configuration below its own `groups` element and the base
  configuration at the top

#### Scenario: More than one group
- **WHEN** a configuration holding two groups is loaded and committed
- **THEN** a read returns both groups, each identified by its own `name`

#### Scenario: Configuration without groups
- **WHEN** a configuration with no `groups` element is loaded and committed
- **THEN** a read returns it unchanged, with no `groups` element introduced

### Requirement: Edits apply inside the tree
An `edit-config` patch SHALL be applied at the path it names, including paths below a `groups` entry, leaving every
other part of the configuration as it was. A delete naming a group SHALL remove that group and nothing else.

#### Scenario: Edit below a group
- **WHEN** a patch replaces one leaf below a `groups` entry
- **THEN** that leaf changes, and the rest of the group, the other groups and the base hierarchy are unchanged

#### Scenario: Delete a group
- **WHEN** a patch deletes a `groups` entry
- **THEN** that group is gone from the committed configuration and any other group remains

#### Scenario: Group references keep their order
- **WHEN** a configuration carrying several `apply-groups` entries is committed
- **THEN** a read returns them in the order they were sent
