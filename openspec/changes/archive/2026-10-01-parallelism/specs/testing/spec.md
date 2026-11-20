# Spec Delta

## ADDED Requirements

### Requirement: Terraform end-to-end applies use bounded device parallelism
The GitHub Terraform end-to-end workflow SHALL apply independent device resources with an explicit parallelism greater than one and SHALL use the same bounded value for the standard and generic provider implementations and for every apply phase in the lifecycle.

#### Scenario: Initial and changed configurations apply concurrently
- **WHEN** either provider implementation runs the initial apply or applies a saved changed plan
- **THEN** Terraform executes the apply with the workflow's bounded parallelism value

#### Scenario: Drift reconciliation uses the same bound
- **WHEN** either provider implementation reconciles simulated out-of-band drift
- **THEN** Terraform executes the reconciliation apply with the same bounded parallelism value used for the initial and changed-plan applies

#### Scenario: Negative apply remains bounded and observable
- **WHEN** either provider implementation runs the bad-credential negative apply
- **THEN** Terraform executes it with the same bounded parallelism value, reports a non-zero exit code, and retains the existing failure diagnostics

### Requirement: Parallel Terraform execution preserves lifecycle validation
The GitHub Terraform end-to-end workflow SHALL preserve its existing state, no-op, changed-plan, drift, reconciliation, and authentication-failure assertions when apply operations run concurrently.

#### Scenario: Parallel lifecycle completes successfully
- **WHEN** the standard and generic matrix jobs run against independent per-device NETCONF mock listeners
- **THEN** every existing lifecycle assertion passes and each job remains separately identifiable by provider implementation
