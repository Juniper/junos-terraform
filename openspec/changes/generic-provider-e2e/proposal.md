# Proposal

## Why

The existing end-to-end plan flow exercises only the standard generated provider, leaving the `--generic` provider option without equivalent workflow coverage. Running the same Terraform plan through both implementations will expose integration defects and prevent generic-provider regressions from reaching users.

## What Changes

- Extend GitHub Actions to run the existing mock-backed Terraform lifecycle with both standard and generic provider implementations in separately labeled jobs.
- Extend the Ansible mock workflow to run the same role/playbook scenario using both its normal schema and a schema produced by the generic Terraform generator. Ansible has no separate generic role backend.
- Keep provider/schema selection explicit so logs and uploaded artifacts identify which variant failed.
- Fix generic-provider defects that are reproduced by the new end-to-end pass and add regression coverage for those defects.

## Capabilities

### New Capabilities

None.

### Modified Capabilities

- `testing`: Require the Terraform end-to-end workflow to validate both standard and generic provider builds using the same plan scenario.

## Impact

- GitHub workflows for Terraform and Ansible, plus the parameterized Python workflow test.
- Existing provider build/install scripts under `.github/prompts/` and `examples/providers/`.
- Existing Terraform plan inputs under `examples/terraform_files/` and the testing capability documentation/specification.
- Generic provider generation/runtime paths only where failures are demonstrated by the new test.