# Spec Delta

## ADDED Requirements

### Requirement: Paths below configuration groups
The diff engine SHALL treat configuration below a `groups` entry as ordinary configuration, keyed by the group's
`name`, so that a change inside a group produces the same minimal edit as the equivalent change in the base hierarchy.
Edits SHALL be scoped to the group they belong to and SHALL NOT affect identically named nodes in the base hierarchy or
in another group.

#### Scenario: Minimal edit inside a group
- **WHEN** one leaf inside a group changes
- **THEN** the emitted edit targets that leaf below its `<groups><name>…</name></groups>` element and contains no other
  changed node

#### Scenario: Same path in two groups
- **WHEN** the same configuration path exists in two groups with different values
- **THEN** each is tracked separately and changing one leaves the other untouched

#### Scenario: Group removal
- **WHEN** a managed group is removed from the configuration
- **THEN** the emitted edit deletes that group as a single operation rather than one operation per leaf

### Requirement: Ordered group references
`apply-groups` SHALL be diffed as an order-sensitive leaf-list: a change of order SHALL produce a non-empty diff, and
the emitted edit SHALL result in the declared order on the device.

#### Scenario: Order change produces an edit
- **WHEN** the entries of `apply-groups` are the same but in a different order
- **THEN** the diff is non-empty and applying it leaves the device in the declared order

#### Scenario: Identical order produces no edit
- **WHEN** the entries and their order are unchanged
- **THEN** the diff is empty
