## ADDED Requirements

### Requirement: Generic provider generation
`jtaf-provider --generic` SHALL write the provider's main package (`main.go` calling `generic.Serve`, `embed_schema.go`) and the schema compact and gzipped.

#### Scenario: Excluding subtrees
- **WHEN** `--exclude PATH` is given (repeatable) to `jtaf-provider`, or to `jtaf-yang2go`, which passes it on
- **THEN** the subtree at PATH, relative to `configuration`, SHALL be left out of the schema; a PATH that does not exist SHALL be an error

#### Scenario: Top-level version
- **WHEN** the schema has a top-level `version` leaf (the Junos release the configuration was committed with)
- **THEN** it SHALL be left out, whether or not the schema is trimmed; nested leaves named `version` SHALL be kept
