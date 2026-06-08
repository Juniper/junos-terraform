---
description: Preview Terraform changes for one or more Junos devices by mode (preview vqfx, preview srx, preview all, or preview explicit targets) without applying.
---

Preview Terraform configuration changes by device scope.

Use when:
- You want a safe plan-only preview for VQFX, SRX, all devices, or explicit targets.

Accepted input forms (strict):
- /preview-devices preview vqfx
- /preview-devices preview qfx
- /preview-devices preview srx
- /preview-devices preview vsrx
- /preview-devices preview all
- /preview-devices preview targets <resource-address> [<resource-address> ...]
- /preview-devices preview <mode> full

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
- Use one reusable plan file for all preview modes: preview.plan
- Treat vsrx as an alias of srx (same target set).
- Treat qfx as an alias of vqfx (same target set).

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

4. Run plan-only command
- Never run apply in this prompt.
- Do not run manual preflight connectivity checks (for example, nc).
- Let terraform plan be the reachability gate. If any target is unreachable, plan fails and you must stop.
- After a successful plan, always render readable output with terraform show -no-color <plan_file>.

Mode: preview vqfx (or preview qfx)
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
  -out preview.plan

Mode: preview srx (or preview vsrx)
terraform plan \
  -var-file=<resolved-var-file> \
  -target='terraform-provider-junos-vsrx-evpn-vxlan.dc1-firewall1-base-config' \
  -target='terraform-provider-junos-vsrx-evpn-vxlan.dc1-firewall2-base-config' \
  -target='terraform-provider-junos-vsrx-evpn-vxlan.dc2-firewall1-base-config' \
  -target='terraform-provider-junos-vsrx-evpn-vxlan.dc2-firewall2-base-config' \
  -out preview.plan

Mode: preview all
terraform plan -var-file=<resolved-var-file> -out preview.plan

Mode: preview targets <resource...>
terraform plan -var-file=<resolved-var-file> -target='<resource-1>' -target='<resource-2>' ... -out preview.plan

After successful plan (all modes)
terraform show -no-color preview.plan > preview.txt

Readable output file (overwritten each run)
- preview.txt

Optional full readable output:
- If input includes full, render the saved plan with:
  terraform show -no-color <plan_file> > <plan_file>.txt
- Return the output file path and line count.

5. Required output contract (compact by default)
- Return only 1 short summary line and 4 fields:
  mode, plan_file, exit_code, plan_summary.
- Add warnings only if present.
- Add blocking_error only if present.
- Do not include extra command output unless user asks for full details.

Guardrails:
- Never run terraform apply in this prompt. If user requests apply, stop and redirect to /deploy-devices.
- Use saved plan files for every preview.
- Do not run terraform init when provider dev_overrides are active.
- Do not run extra commands or checks beyond the requested preview flow.
