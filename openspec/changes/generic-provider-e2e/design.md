# Design

## Context

See proposal.md for motivation. The current end-to-end command builds and installs standard providers before running one plan. The standard and `--generic` generation scripts write to the same provider module directories, and the generic build script also installs binaries, so a naive second build replaces the first implementation before both can be exercised. The plan workflow uses dev overrides and does not run `terraform init` or `terraform apply`.

## Goals / Non-Goals

**Goals:**
- Keep separately selectable standard and generic provider artifacts for the duration of one workflow run.
- Reuse the same Terraform files and plan command for both passes, with clear per-variant results.
- Leave the standard provider selected as the normal/default provider after the workflow.
- Make failures actionable and fix generic-provider defects demonstrated by the end-to-end pass.

**Non-Goals:**
- Applying configuration to devices or changing the Terraform scenario to make the two providers produce identical plan text.
- Replacing the existing provider-generation implementation or adding a new test framework.
- Fixing unrelated failures discovered while running the broader repository test suites.

## Decisions

### Stage both builds before running either plan

Generate the standard and generic provider modules into separate staging locations, then compile each variant to its own provider-binary directory. The current generators share output paths, so the workflow must preserve each module before invoking the other mode, or direct generation to distinct workspaces. The Terraform configuration is generated once from the existing scenario and is not changed between runs.

**Alternative considered:** Generate, test, and then regenerate for the second mode. Rejected because it does not satisfy the requirement that both versions be prepared up front and makes a later generic build failure occur after the standard test has already run.

### Select each implementation per Terraform process

Use a temporary Terraform CLI configuration with dev overrides pointing at the selected staged binaries. Invoke the existing plan scenario first for the standard provider, then for the generic provider. Use separate plan/output artifacts and label each invocation so the result and diagnostics identify their provider variant. Keep selection process-local instead of replacing binaries in the user's global Go bin directory or editing persistent Terraform CLI configuration.

**Alternative considered:** Install each variant in turn under the same global provider binary name. Rejected because it mutates user state, risks leaving the generic binary selected after failure, and couples test selection to global install paths.

### Preserve both outcomes and fail after both passes

Capture each plan's exit status and output independently. If both builds succeed, run both plans even when the standard plan fails, then return an overall failure if either plan fails. A build/setup failure prevents plan execution because a comparable pass cannot be made for both variants. Do not run `terraform apply`.

### Fix only test-reproduced generic-provider defects

Use the generic pass to reproduce defects and add the narrowest regression check at the owning layer. Keep unrelated failures out of this change; report them rather than expanding scope.

## Risks / Trade-offs

- [The providers may produce different plan details despite both being functionally valid] → Treat each plan's success and diagnostics as the parity gate; do not require byte-identical plan output.
- [Staging can consume additional disk space and increase workflow time] → Reuse existing generated inputs and remove temporary artifacts after preserving the reported plan outputs.
- [Terraform may pick up a user's global CLI configuration unexpectedly] → Set the temporary CLI config explicitly for each process and retain the existing no-init behavior.
- [A provider-specific failure could obscure the other implementation's result] → Record both plan outcomes independently and aggregate status after both invocations.

## Migration Plan

No user migration is required. The workflow remains no-apply; temporary staging and CLI configuration are cleaned up on success or failure, and the standard provider remains the default installed/selected variant.