---
name: deploy-multiple-devices
description: "Deploy Terraform config to multiple Junos devices in one operation. Use when user says push vqfx, push srx, deploy borderleaf pair, or run plan/apply for multiple targets."
---

Deploy Terraform configuration to a selected set of devices safely.

## Inputs

- Terraform working directory
- Target mode (all, vqfx, srx, role group, or explicit resource list)
- Optional plan filename

## Steps

1. Validate context
- Confirm working directory contains Terraform files.
- Confirm host and port are configured for each targeted provider alias.

2. Resolve target set
- vqfx: include all terraform-provider-junos-vqfx-evpn-vxlan resources requested by scope.
- srx: include all terraform-provider-junos-vsrx-evpn-vxlan resources requested by scope.
- role group: map requested group (example: borderleaf) to exact resource addresses.
- explicit list: validate each supplied resource address.
- If ambiguous, list candidates and ask user to confirm.

3. Create batch plan
- For all devices, run full terraform plan -out <planfile>.
- For subsets, run one terraform plan with multiple -target flags and -out <planfile>.

4. Apply batch plan
- Run terraform apply <planfile>.

5. Verify post-apply
- Run terraform plan.
- Confirm if final plan is clean or list remaining changes.

6. Report
- Summarize mode (vqfx/srx/group/list), targeted devices, plan filename, apply status, and post-check result.

## Guardrails

- Prefer saved plan files for batch operations.
- Do not use -auto-approve unless the user explicitly requests it.
- Do not run terraform init when provider dev_overrides are active.
- If endpoint lookup or connectivity fails for any target, stop and provide host/port remediation steps.
