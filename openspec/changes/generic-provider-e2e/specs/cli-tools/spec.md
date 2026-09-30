# Spec Delta

## ADDED Requirements

### Requirement: Generic provider generation excludes XML filtering
The `jtaf-yang2go` and `jtaf-provider` commands SHALL reject invocations that combine `--generic` with `-x` or `--xml-config`. Generic provider generation SHALL consume the unfiltered YANG schema; XML files remain valid inputs to downstream configuration conversion tools.

#### Scenario: End-to-end generic generation rejects XML filtering
- **WHEN** `jtaf-yang2go` is invoked with both `--generic` and `-x`
- **THEN** it exits with an argument error before invoking pyang or generating a provider

#### Scenario: Direct generic provider generation rejects XML filtering
- **WHEN** `jtaf-provider` is invoked with both `--generic` and `-x` or `--xml-config`
- **THEN** it exits with an argument error before loading or generating a provider

#### Scenario: Generic schema is converted using XML downstream
- **WHEN** a full generic provider schema is passed to `jtaf-xml2tf` or `jtaf-xml2yaml` with XML configuration inputs
- **THEN** the downstream tool uses those XML files for configuration conversion without applying provider-generation XML filtering