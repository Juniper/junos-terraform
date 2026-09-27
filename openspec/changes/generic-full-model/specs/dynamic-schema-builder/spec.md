## ADDED Requirements

### Requirement: Build the protocol schema from schema nodes
The provider SHALL build a tfprotov6 schema and its value type from the flattened schema nodes: a leaf as an optional string, a leaf-list as an optional list of strings, a container or list as an optional list of nested objects, plus a required `resource_name`.

#### Scenario: Attribute names
- **WHEN** a YANG name contains `-`, `.` or capitals
- **THEN** the attribute name SHALL replace `-` and `.` with `_` and be lowercase

#### Scenario: Invalid or colliding attribute names
- **WHEN** `jtaf-provider --generic` generates a provider whose schema has a name that is not a valid attribute name, or two sibling nodes (with choices flattened) with the same attribute name
- **THEN** it SHALL fail, naming the node's path
