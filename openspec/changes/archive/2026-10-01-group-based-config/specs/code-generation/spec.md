# Spec Delta

## ADDED Requirements

### Requirement: One provider generator
Provider generation SHALL produce a single schema-driven Go provider whose source does not depend on the contents of
the YANG model. No Go source SHALL be rendered from a template. The generated output SHALL differ between device types
only in the module name, the provider type name and the embedded schema.

#### Scenario: No templated Go source
- **WHEN** a provider is generated for any device type
- **THEN** the output directory contains no file produced by template expansion of schema contents, and its size does
  not grow with the number of nodes in the model

#### Scenario: Equivalent providers across models
- **WHEN** providers are generated for two different models
- **THEN** their Go sources are identical apart from module name, provider type name and embedded schema

### Requirement: Generated provider carries the group decision
The embedded schema SHALL be the sole record of whether groups are managed. A provider built with groups SHALL expose
them, and a provider built without groups SHALL NOT, with no runtime flag, environment variable or provider-block
setting able to change that.

#### Scenario: Behaviour follows the binary
- **WHEN** a provider built without groups is used with a `.tf` file that declares a `groups` block
- **THEN** Terraform reports an unsupported attribute, regardless of provider configuration

## MODIFIED Requirements

### Requirement: Generated providers carry a compressed schema
Provider generation SHALL write the schema as compact gzip-compressed JSON named `trimmed_schema.json.gz` and SHALL
NOT write a plain `trimmed_schema.json`, whether the schema was trimmed to XML or left untrimmed.

#### Scenario: Schema is available to downstream tools
- **WHEN** provider generation completes
- **THEN** `trimmed_schema.json.gz` contains the schema and `jtaf-xml2tf` can consume it to produce the same output as
  the equivalent plain JSON schema
