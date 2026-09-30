# Spec Delta

## ADDED Requirements

### Requirement: Terraform end-to-end plans cover both provider implementations
The Terraform end-to-end workflow SHALL generate the standard and generic provider variants for every configured device family before either test pass begins, and SHALL preserve both variants for independent selection.

#### Scenario: Both variants are prepared before testing
- **WHEN** the Terraform end-to-end workflow starts
- **THEN** it generates and builds both provider variants before running a Terraform plan, without one variant overwriting the only copy of the other

### Requirement: Both provider variants run the same end-to-end plan
The Terraform end-to-end workflow SHALL run the existing plan scenario first with the standard provider and then with the generic provider, using the same Terraform configuration and inputs. It SHALL not apply changes, SHALL retain separately identifiable results for each run, and SHALL report failures with the provider variant that produced them.

#### Scenario: Standard and generic providers pass the plan scenario
- **WHEN** both provider variants are built successfully
- **THEN** the workflow runs the same plan scenario with the standard provider followed by the generic provider and records each result separately

#### Scenario: A provider variant fails its plan
- **WHEN** either plan invocation exits unsuccessfully
- **THEN** the workflow identifies the failing provider variant, retains that invocation's diagnostic output, and reports an unsuccessful overall result
