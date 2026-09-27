## ADDED Requirements

### Requirement: Schema lookups through a compact table
The patch engine SHALL look schema nodes up in a `patch.Schema`: one record per node, children consecutive, names stored once, compiled from the pyang JSON.

#### Scenario: Lookup by schema path
- **WHEN** the engine looks up a path relative to `<configuration>`, without list predicates
- **THEN** it SHALL get the node's kind, list key, presence and ordered-by-user, or not found

#### Scenario: Same-named siblings
- **WHEN** two siblings have the same name (the same node in two cases of a choice)
- **THEN** they SHALL be one node, with the last one's kind and key, either one's flags, and the children of both
