## ADDED Requirements

### Requirement: Build the protocol schema from schema nodes
The provider SHALL build a tfprotov6 schema and its value type from the flattened schema nodes: a leaf as an optional string, a leaf-list as an optional list of strings, a container or list as an optional list of nested objects, plus a required `resource_name`.

#### Scenario: Attribute names
- **WHEN** a YANG name contains `-`, `.` or capitals
- **THEN** the attribute name SHALL replace `-` and `.` with `_` and be lowercase
