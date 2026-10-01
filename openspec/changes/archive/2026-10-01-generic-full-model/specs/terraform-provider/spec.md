## ADDED Requirements

### Requirement: Generic provider serves the plugin protocol directly
The generic provider SHALL implement tfprotov6 with terraform-plugin-go, without terraform-plugin-framework.

#### Scenario: Plan
- **WHEN** Terraform calls `PlanResourceChange`
- **THEN** the planned state SHALL be the proposed state, and a change of `resource_name` SHALL require replacement

#### Scenario: State from an earlier schema
- **WHEN** Terraform calls `UpgradeResourceState` with attributes the schema does not have
- **THEN** they SHALL be ignored and missing attributes SHALL be null

#### Scenario: Update
- **WHEN** Terraform applies a change
- **THEN** the provider SHALL patch the difference between the device's configuration and the plan, commit, and if the device still differs load the plan and commit again

#### Scenario: Unsupported operations
- **WHEN** Terraform calls import, a data source, a function, an ephemeral or list resource, or an action
- **THEN** the provider SHALL return an error diagnostic
