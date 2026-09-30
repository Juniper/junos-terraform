# Proposal

## Why

The existing end-to-end plan flow exercises only the standard generated provider, leaving the `--generic` provider option without equivalent workflow coverage. Running the same Terraform plan through both implementations will expose integration defects and prevent generic-provider regressions from reaching users.

## What Changes

- Extend the end-to-end workflow to generate both standard and generic provider builds before running either test pass, preserving each build so their shared output paths do not overwrite one another.
- Run the existing Terraform plan first with the standard provider, then switch to the generic provider and rerun the same plan and configuration without applying changes.
- Keep provider selection and restoration explicit so the workflow can report which implementation failed and leave the standard provider as the default afterward.
- Fix generic-provider defects that are reproduced by the new end-to-end pass and add regression coverage for those defects.

## Capabilities

### New Capabilities

None.

### Modified Capabilities

- `testing`: Require the Terraform end-to-end workflow to validate both standard and generic provider builds using the same plan scenario.

## Impact

- End-to-end workflow and provider build/install scripts under `.github/prompts/` and `examples/providers/`.
- Existing Terraform plan inputs under `examples/terraform_files/` and the testing capability documentation/specification.
- Generic provider generation/runtime paths only where failures are demonstrated by the new test.