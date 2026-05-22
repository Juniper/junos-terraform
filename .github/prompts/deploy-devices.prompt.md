---
description: Deploy Terraform config to one or more Junos devices by mode (push vqfx, push srx, all devices, or explicit targets).
---

Deploy Terraform configuration by device scope.

For plan-only previews, use `/preview-devices`.

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

2. Resolve variable file deterministically
- If terraform.tfvars exists, use it.
- Else stop and report missing var file.

3. Resolve target set and output plan name
- push vqfx -> vqfx.plan
- push srx -> srx.plan
- push all -> all-devices.plan
- push targets -> custom-targets.plan
- Treat qfx as an alias of vqfx (same target set and plan file).
- Treat vsrx as an alias of srx (same target set and plan file).

Target sets:

VQFX targets:
- terraform-provider-junos-vqfx-evpn-vxlan.dc1-borderleaf1-base-config
- terraform-provider-junos-vqfx-evpn-vxlan.dc1-borderleaf2-base-config
- terraform-provider-junos-vqfx-evpn-vxlan.dc1-leaf1-base-config
- terraform-provider-junos-vqfx-evpn-vxlan.dc1-leaf2-base-config
- terraform-provider-junos-vqfx-evpn-vxlan.dc1-leaf3-base-config
- terraform-provider-junos-vqfx-evpn-vxlan.dc1-spine1-base-config
- terraform-provider-junos-vqfx-evpn-vxlan.dc1-spine2-base-config
- terraform-provider-junos-vqfx-evpn-vxlan.dc2-spine1-base-config
- terraform-provider-junos-vqfx-evpn-vxlan.dc2-spine2-base-config

SRX targets:
- terraform-provider-junos-vsrx-evpn-vxlan.dc1-firewall1-base-config
- terraform-provider-junos-vsrx-evpn-vxlan.dc1-firewall2-base-config
- terraform-provider-junos-vsrx-evpn-vxlan.dc2-firewall1-base-config
- terraform-provider-junos-vsrx-evpn-vxlan.dc2-firewall2-base-config

4. Plan and apply sequence
- Do not run manual preflight connectivity checks (for example, nc).
- Run terraform plan first and treat it as the reachability gate.
- Run terraform apply only if plan exits with code 0 and a plan file is saved.

Mode: push vqfx (or push qfx)
terraform plan \
  -var-file=<resolved-var-file> \
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

Mode: push srx (or push vsrx)
terraform plan \
  -var-file=<resolved-var-file> \
  -target='terraform-provider-junos-vsrx-evpn-vxlan.dc1-firewall1-base-config' \
  -target='terraform-provider-junos-vsrx-evpn-vxlan.dc1-firewall2-base-config' \
  -target='terraform-provider-junos-vsrx-evpn-vxlan.dc2-firewall1-base-config' \
  -target='terraform-provider-junos-vsrx-evpn-vxlan.dc2-firewall2-base-config' \
  -out srx.plan
terraform apply srx.plan

Mode: push all
terraform plan -var-file=<resolved-var-file> -out all-devices.plan
terraform apply all-devices.plan

Mode: push targets <resource...>
terraform plan -var-file=<resolved-var-file> -target='<resource-1>' -target='<resource-2>' ... -out custom-targets.plan
terraform apply custom-targets.plan

5. Post-apply verification
terraform plan -var-file=<resolved-var-file>

6. Required output contract (compact by default)
- Return only 1 short summary line and 6 fields:
  mode, plan_file, plan_exit_code, apply_exit_code, verify_exit_code, plan_summary.
- Add warnings only if present.
- Add blocking_error only if present.
- Do not include extra command output unless user asks for full details.

Guardrails:
- Use saved plan files for every deployment.
- Do not use -auto-approve unless explicitly requested.
- Do not run terraform init when provider dev_overrides are active.
