# Proposal

## Why

The Terraform end-to-end workflow forces independent device resources to apply serially with `-parallelism=1`, making the generic lifecycle substantially slower than necessary. Local reproduction of the GitHub scenario reduced a nine-device generic apply from 45.67 seconds to 19.49 seconds at parallelism three while preserving successful state transitions.

## What Changes

- Run every Terraform apply phase in the mock-backed standard/generic end-to-end workflow with an explicit parallelism of three.
- Define the parallelism once for the workflow so initial apply, changed-plan apply, drift reconciliation, and the negative credential apply cannot diverge.
- Preserve the existing lifecycle assertions, per-device NETCONF isolation, diagnostics, and implementation-specific job labeling.
- Document the measured baseline and the reason for selecting three rather than unrestricted concurrency.

## Capabilities

### New Capabilities

None.

### Modified Capabilities

- `testing`: Require the Terraform standard/generic end-to-end lifecycle to apply independent device resources with bounded parallel execution while preserving deterministic validation and failure diagnostics.

## Impact

- Affects `.github/workflows/go-terraform-provider.yml` and its workflow validation coverage.
- Changes only CI execution scheduling; provider APIs, generated Terraform configuration, resource semantics, and device behavior are unchanged.
- Increases simultaneous NETCONF sessions within each integration job from one to at most three.
