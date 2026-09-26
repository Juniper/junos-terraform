## ADDED Requirements

### Requirement: Convert between values and configuration XML with list keys
`ValueToConfig` and `ConfigToValue` SHALL convert between the resource's value and a `<configuration>` tree by walking the schema, keeping every list entry.

#### Scenario: Several list entries
- **WHEN** the configuration has several entries of a list
- **THEN** the value SHALL have each entry, in document order

#### Scenario: Keys first
- **WHEN** a list entry is written as XML
- **THEN** its key leaves SHALL come first, in key order

#### Scenario: Empty and missing elements
- **WHEN** an element is empty
- **THEN** its leaf value SHALL be `""`; a missing element SHALL be null

#### Scenario: Unmodeled elements
- **WHEN** the device returns elements not in the schema
- **THEN** they SHALL be ignored
