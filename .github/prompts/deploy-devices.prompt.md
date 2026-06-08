---
description: Deploy Terraform config to one or more Junos devices by mode (push vqfx, push srx, all devices, or explicit targets).
---

Deploy Terraform configuration by device scope.

For plan-only previews, use `/preview-devices`.

Mode keyword note:
- `vqfx`, `qfx`, `srx`, and `vsrx` are command-mode keywords, not fixed inventory names.

Accepted input forms (strict):
- /deploy-devices push vqfx
- /deploy-devices push qfx
- /deploy-devices push srx
- /deploy-devices push vsrx
- /deploy-devices push all
- /deploy-devices push targets <resource-address> [<resource-address> ...]

If input is not one of these forms, stop and ask user to choose one accepted form.

Execution rules:

1. Resolve working directory deterministically
- Required path: examples/terraform_files
- If current directory is not this path, change into it first.
- If path does not exist, stop and report missing path.

2. Resolve deployment mode
- Treat qfx as an alias of vqfx.
- Treat vsrx as an alias of srx.
- For `push vqfx`, `push srx`, and `push all`, use full `terraform plan` and `terraform apply` (no `-target`).
- Use `-target` only for `push targets <resource...>` or explicit emergency recovery.

Example target addresses (only for `push targets` or emergency targeted mode):
- The values below are examples only. Replace them with resource addresses from your environment.
- terraform-provider-junos-vqfx-evpn-vxlan.dc1-borderleaf1-base-config
- terraform-provider-junos-vqfx-evpn-vxlan.dc1-borderleaf2-base-config
- terraform-provider-junos-vqfx-evpn-vxlan.dc1-leaf1-base-config
- terraform-provider-junos-vqfx-evpn-vxlan.dc1-leaf2-base-config
- terraform-provider-junos-vqfx-evpn-vxlan.dc1-leaf3-base-config
- terraform-provider-junos-vqfx-evpn-vxlan.dc1-spine1-base-config
- terraform-provider-junos-vqfx-evpn-vxlan.dc1-spine2-base-config
- terraform-provider-junos-vqfx-evpn-vxlan.dc2-spine1-base-config
- terraform-provider-junos-vqfx-evpn-vxlan.dc2-spine2-base-config
- terraform-provider-junos-vsrx-evpn-vxlan.dc1-firewall1-base-config
- terraform-provider-junos-vsrx-evpn-vxlan.dc1-firewall2-base-config
- terraform-provider-junos-vsrx-evpn-vxlan.dc2-firewall1-base-config
- terraform-provider-junos-vsrx-evpn-vxlan.dc2-firewall2-base-config

3. Plan and apply sequence
- Do not run manual preflight connectivity checks (for example, nc).
- Run terraform plan first and treat it as the reachability gate.
- Run terraform apply only if plan exits with code 0.

Mode: push vqfx (or push qfx)
terraform plan
terraform apply

Mode: push srx (or push vsrx)
terraform plan
terraform apply

Mode: push all
terraform plan
terraform apply

Mode: push targets <resource...>
terraform plan -target='<resource-1>' -target='<resource-2>' ...
terraform apply -target='<resource-1>' -target='<resource-2>' ...

Targeted mode warning:
- `-target` is not routine mode. Use only for explicit target deploys or recovery/troubleshooting.

4. Post-apply verification
terraform plan

5. Required output contract (compact by default)
- Return only 1 short summary line and 6 fields:
  mode, plan_exit_code, apply_exit_code, verify_exit_code, plan_summary, terminal_output_shown.
- Add warnings only if present.
- Add blocking_error only if present.
- Do not include extra command output unless user asks for full details.

Guardrails:
- Show plan/apply output in terminal/chat directly; do not create plan or text artifact files.
- Do not use -auto-approve unless explicitly requested.
- Do not run terraform init when provider dev_overrides are active.
