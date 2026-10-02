# Tasks

## 1. Bounded Terraform Apply Parallelism

- [x] 1.1 Define a single Terraform apply parallelism value of three in `.github/workflows/go-terraform-provider.yml`, use it for the initial, changed-plan, drift-reconciliation, and bad-credential apply commands in both matrix variants, and verify no `terraform apply` command retains `-parallelism=1`.
- [x] 1.2 Add focused workflow regression coverage that verifies the shared value is three and every Terraform apply phase references it, then run the targeted Python test successfully.

## 2. Integration Verification

- [x] 2.1 Validate the workflow YAML and run the standard mock-backed Terraform matrix job, verifying all existing lifecycle and negative-auth assertions pass under parallel execution.
- [x] 2.2 Run the generic mock-backed Terraform matrix job, verifying all existing lifecycle and negative-auth assertions pass and record its apply timing for comparison with the 45.67-second serial baseline.
- [x] 2.3 Run `openspec validate parallelism --strict` and verify the completed implementation satisfies the testing spec delta.
