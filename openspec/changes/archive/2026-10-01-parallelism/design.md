# Design

## Context

The Terraform integration workflow creates nine resources, each bound to a distinct provider alias and NETCONF mock listener. Terraform plans already use their normal graph concurrency, but all four apply commands override Terraform with `-parallelism=1` even though no resource shares a device endpoint.

Local reproduction of the same nine-resource generic saved-plan apply measured 45.67 seconds at parallelism one, 19.49 seconds at parallelism three, and 13.81 seconds at parallelism nine. Three captures most of the practical gain without scheduling every provider process and NETCONF session simultaneously.

## Goals / Non-Goals

**Goals:**

- Remove unnecessary serialization from all Terraform apply phases in the standard/generic integration matrix.
- Keep one explicit, auditable concurrency bound across the entire workflow.
- Preserve device isolation, lifecycle assertions, and useful negative-test diagnostics.

**Non-Goals:**

- Change provider runtime behavior, Terraform resource semantics, or generated configuration.
- Change plan concurrency or the matrix job concurrency.
- Establish a wall-clock CI performance requirement.
- Maximize concurrency based on the number of devices.

## Decisions

### Use a workflow-level apply parallelism of three

Define one workflow environment value for Terraform apply parallelism and pass it explicitly to every `terraform apply` invocation. Both matrix implementations therefore exercise identical scheduling.

Three is selected because the measured generic create improved by 2.34 times over serial execution, while avoiding the higher simultaneous process, memory, SSH, and device load of one operation per device. A named workflow value keeps future tuning localized.

Alternatives considered:

- Keep `-parallelism=1`: deterministic but retains the measured bottleneck.
- Use Terraform's default of ten: fastest in the local mock for nine devices, but needlessly aggressive for shared GitHub runners and less representative of cautious network automation.
- Configure `TF_CLI_ARGS_apply`: concise, but hides the effective argument from individual commands and can make diagnostics harder to interpret.

### Apply the bound to all apply phases

Initial create, changed-plan apply, drift reconciliation, and bad-credential negative apply use the same value. This prevents the performance behavior and coverage from drifting between scenarios.

Plans remain unchanged because they do not currently force serial execution. The workflow's matrix `max-parallel` remains independent: it controls jobs, whereas this change controls resource operations inside one Terraform process.

### Preserve independent-device safety assumptions

Parallel execution is valid because each generated resource uses a distinct provider alias and each alias targets a distinct mock listener. No two operations mutate the same candidate datastore. Existing post-apply state and mock-history assertions remain the correctness guard.

## Risks / Trade-offs

- [Concurrent execution makes log ordering nondeterministic] -> Keep implementation-specific job names and existing uploaded diagnostics; assertions must not depend on operation order.
- [Three simultaneous generic provider operations increase peak memory] -> Use a conservative bound of three rather than the measured maximum of nine.
- [The negative test may report a different failing device first] -> Assert failure by exit code and preserve logs rather than matching device order.
- [Future generated resources could share a device] -> Retain the invariant that this E2E fixture maps one resource to one provider alias/listener; revisit the bound if that fixture changes.

## Migration Plan

1. Add the workflow-level parallelism value.
2. Replace each serial apply flag with the shared value.
3. Validate workflow syntax and command coverage.
4. Run standard and generic mock-backed Terraform matrix jobs.

Rollback is a one-line value change back to one if runner stability or device-isolation assumptions fail.
