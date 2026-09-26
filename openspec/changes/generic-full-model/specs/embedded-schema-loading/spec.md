## ADDED Requirements

### Requirement: Load a full model
The provider SHALL load a schema generated from a full Junos model, raw or gzipped.

#### Scenario: Decimal64 range bounds
- **WHEN** a range bound is a numeric string (`"9223372036.854775807"`)
- **THEN** it SHALL be parsed as a number

#### Scenario: Choice and case nodes
- **WHEN** the schema has YANG `choice` or `case` nodes
- **THEN** their children SHALL be indexed and exposed directly under the choice's parent

#### Scenario: Schema fails to load
- **WHEN** the schema cannot be parsed, or two sibling names map to the same attribute
- **THEN** `GetProviderSchema` SHALL return an error diagnostic naming the cause
