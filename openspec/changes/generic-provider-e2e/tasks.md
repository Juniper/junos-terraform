# Tasks

## 1. Prepare Both Provider Variants

- [x] 1.1 Update provider generation/build orchestration to create and compile standard and generic variants into separate staging locations before either plan run; verify both expected provider binaries and schemas are present and neither build overwrites the other.
- [x] 1.2 Keep generated Terraform inputs shared between the two variants and verify the existing conversion step produces the same configuration used by both plan passes.

## 2. Run the End-to-End Plan Twice

- [x] 2.1 Extend the existing end-to-end workflow to select the standard staged providers and run the existing Terraform plan first; verify the configuration and plan command are unchanged and no apply or init is run.
- [x] 2.2 Select the generic staged providers and rerun the same plan scenario; verify the workflow keeps per-variant plan output, diagnostics, and exit status, and leaves the standard provider as the default selection.
- [x] 2.3 Verify failure handling reports the provider variant and preserves its diagnostic output, still runs the second plan when the first plan fails, and cleans temporary selection/configuration state on success and failure.

## 3. Resolve Reproduced Generic Provider Defects

- [ ] 3.1 Run the generic-provider end-to-end pass, fix each in-scope defect it reproduces, and add a focused regression test for every fix; verify each regression test fails before the fix and passes afterward.
- [x] 3.2 Update the end-to-end workflow documentation/output contract to describe both provider results and verify the documented invocation matches the implemented flow.

## 4. Verify Integration

- [ ] 4.1 Run the end-to-end workflow against the configured test target and verify both provider plans complete against the identical Terraform scenario without applying changes.
- [x] 4.2 Run the focused Python and Go tests for touched components and verify both provider build paths still succeed.