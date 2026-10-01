# Tasks

## 1. Prepare Both Provider Variants

- [x] 1.1 Update provider generation/build orchestration to create and compile standard and generic variants into separate staging locations before either plan run; verify both expected provider binaries and schemas are present and neither build overwrites the other.
- [x] 1.2 Keep generated Terraform inputs shared between the two variants and verify the existing conversion step produces the same configuration used by both plan passes.
- [x] 1.3 Make `--generic` and `-x` mutually exclusive in `jtaf-yang2go` and `jtaf-provider`; verify both CLIs reject the combination with an argument error and accept generic mode without XML.
- [x] 1.4 Remove XML arguments from generic provider build, Terraform CI, Ansible CI, and pytest generation calls; verify standard XML filtering and downstream XML conversion remain intact.
- [x] 1.5 Update generic-provider documentation and verify no supported example combines `--generic` with XML filtering.

## 2. Run the End-to-End Plan Twice

- [x] 2.1 Extend the existing end-to-end workflow to select the standard staged providers and run the existing Terraform plan first; verify the configuration and plan command are unchanged and no apply or init is run.
- [x] 2.2 Select the generic staged providers and rerun the same plan scenario; verify the workflow keeps per-variant plan output, diagnostics, and exit status, and leaves the standard provider as the default selection.
- [x] 2.3 Verify failure handling reports the provider variant and preserves its diagnostic output, still runs the second plan when the first plan fails, and cleans temporary selection/configuration state on success and failure.
- [x] 2.4 Add a GitHub Actions matrix to run the same Terraform mock lifecycle with standard and generic provider implementations; verify jobs and uploaded logs are labeled by implementation.
- [ ] 2.5 Add an Ansible schema-source matrix that feeds the generic-generated schema into `jtaf-xml2yaml` and runs the same mock playbook lifecycle; verify logs identify the schema source.

## 3. Resolve Reproduced Generic Provider Defects

- [x] 3.1 Run the generic-provider end-to-end pass, fix each in-scope defect it reproduces, and add a focused regression test for every fix; verify each regression test fails before the fix and passes afterward.
- [x] 3.2 Update the end-to-end workflow documentation/output contract to describe both provider results and verify the documented invocation matches the implemented flow.

## 4. Verify Integration

- [ ] 4.1 Run both GitHub Actions matrix workflows and verify standard/generic Terraform and Ansible jobs pass the same scenarios against the NETCONF mock.
- [x] 4.2 Run the focused Python and Go tests for touched components and verify both provider build paths still succeed.