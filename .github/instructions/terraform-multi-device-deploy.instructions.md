---
applyTo: "examples/terraform_files/**/*.tf"
description: "Use when deploying Terraform config to multiple Junos devices at once, including device-family batches or explicit target sets."
---

# Terraform Multi-Device Deploy Instructions

## Endpoint and Credentials

- Keep provider endpoints configurable through shared settings.
- Keep credentials configurable and treat passwords as sensitive.
- Confirm each targeted alias resolves to a reachable host and port before plan/apply.

## Batch Deployment Patterns

Use one of these patterns depending on intent:

1. Full fleet deployment
- Create one full plan file.
- Apply the same plan file.

2. Provider-type deployment
- Device-family batch (example labels: vqfx, srx).
- Build a family-specific target set from the current environment's resource addresses.
- Build a single plan using multiple -target flags.
- Apply the saved plan file.

3. Role-group deployment
- Example groups: borderleaf pair, all leaves, all firewalls.
- Use one plan file with multiple targets.

## Example Commands (Example Values Only)

The command patterns below are examples. Replace resource addresses and plan file names with values from your environment.

Example family batch (vqfx):

terraform plan \
  -target='terraform-provider-junos-vqfx-evpn-vxlan.dc1-borderleaf1-base-config' \
  -target='terraform-provider-junos-vqfx-evpn-vxlan.dc1-borderleaf2-base-config' \
  -target='terraform-provider-junos-vqfx-evpn-vxlan.dc1-leaf1-base-config' \
  -target='terraform-provider-junos-vqfx-evpn-vxlan.dc1-leaf2-base-config' \
  -target='terraform-provider-junos-vqfx-evpn-vxlan.dc1-leaf3-base-config' \
  -target='terraform-provider-junos-vqfx-evpn-vxlan.dc1-spine1-base-config' \
  -target='terraform-provider-junos-vqfx-evpn-vxlan.dc1-spine2-base-config' \
  -target='terraform-provider-junos-vqfx-evpn-vxlan.dc2-spine1-base-config' \
  -target='terraform-provider-junos-vqfx-evpn-vxlan.dc2-spine2-base-config' \
  -out vqfx.plan
terraform apply vqfx.plan

Example family batch (srx):

terraform plan \
  -target='terraform-provider-junos-vsrx-evpn-vxlan.dc1-firewall1-base-config' \
  -target='terraform-provider-junos-vsrx-evpn-vxlan.dc1-firewall2-base-config' \
  -target='terraform-provider-junos-vsrx-evpn-vxlan.dc2-firewall1-base-config' \
  -target='terraform-provider-junos-vsrx-evpn-vxlan.dc2-firewall2-base-config' \
  -out srx.plan
terraform apply srx.plan

After any targeted batch:

terraform plan

## Guardrails

- Prefer one saved plan file per batch and apply that exact plan.
- Use targeted mode for selective batches; use full plan/apply for routine fleet sync.
- Do not run terraform init when provider dev_overrides are active.
- If a single endpoint fails lookup or connectivity, stop and remediate host/port before retrying the batch.
