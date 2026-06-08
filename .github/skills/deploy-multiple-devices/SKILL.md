---
name: deploy-multiple-devices
description: "Deploy Terraform config to multiple Junos devices in one operation. Example intents: push a device family, deploy a role group, or run plan/apply for explicit targets."
---

Deploy Terraform configuration to a selected set of devices safely.

## Inputs

- Terraform working directory
- Target mode (all devices, device-family keyword, role group, or explicit resource list)
- Optional plan filename

## Steps

1. Validate context
- Confirm working directory contains Terraform files.
- Confirm host and port are configured for each targeted provider alias.

2. Resolve target set
- device-family keyword: include all resources requested by the selected family scope.
- role group: map requested group to exact resource addresses.
- explicit list: validate each supplied resource address.
- If ambiguous, list candidates and ask user to confirm.

Example family mapping (example only):
- `vqfx` -> `terraform-provider-junos-vqfx-evpn-vxlan.*`
- `srx` -> `terraform-provider-junos-vsrx-evpn-vxlan.*`

3. Create batch plan
- For all devices, run full terraform plan -out <planfile>.
- For subsets, run one terraform plan with multiple -target flags and -out <planfile>.

4. Apply batch plan
- Run terraform apply <planfile>.

5. Verify post-apply
- Run terraform plan.
- Confirm if final plan is clean or list remaining changes.

6. Report
- Summarize mode (family/group/list/all), targeted devices, plan filename, apply status, and post-check result.

## Guardrails

- Prefer saved plan files for batch operations.
- Do not use -auto-approve unless the user explicitly requests it.
- Do not run terraform init when provider dev_overrides are active.
- If endpoint lookup or connectivity fails for any target, stop and provide host/port remediation steps.
