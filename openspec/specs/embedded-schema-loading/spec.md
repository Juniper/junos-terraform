# embedded-schema-loading Specification

## Purpose
Carries a Junos model inside the provider binary and makes it available at start-up, including models that cover the whole configuration hierarchy.

## Requirements

### Requirement: Load a full model
The provider SHALL load a schema generated from a full Junos model, compiled or as JSON, raw or gzipped.

#### Scenario: Decimal64 range bounds
- **WHEN** a range bound is a numeric string (`"9223372036.854775807"`)
- **THEN** it SHALL be parsed as a number

#### Scenario: Choice and case nodes
- **WHEN** the schema has YANG `choice` or `case` nodes
- **THEN** their children SHALL be indexed and exposed directly under the choice's parent

#### Scenario: Nodes named configuration
- **WHEN** a node below the root is named `configuration` (`system archival configuration`)
- **THEN** it and its children SHALL be indexed at their own path, not merged into its parent

#### Scenario: Compiled schema
- **WHEN** the embedded schema is `patch.Schema`'s binary form
- **THEN** the provider SHALL read it without parsing JSON, and serve the same schema as the JSON it was compiled from

#### Scenario: Schema fails to load
- **WHEN** the schema cannot be parsed or read
- **THEN** `GetProviderSchema` SHALL return an error diagnostic naming the cause
