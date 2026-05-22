---
applyTo: "examples/terraform_files/**/*.tf"
description: "Use when working on Terraform device provider endpoints, targeted one-device plan/apply, and provider alias deployment workflows."
---

# Terraform Device Deploy Instructions

## Endpoint and Credentials

- Prefer variable-driven endpoints and credentials over hardcoded host, port, username, or password values.
- Keep provider endpoint values in tfvars files so environments can switch without editing provider blocks.
- Treat passwords as sensitive variables.

## One-Device Deployment Workflow

Use this sequence for a single device change:

1. Resolve exact target resource address.
2. Run targeted plan with explicit output file.
3. Apply the saved plan file.
4. Run a full plan to verify no pending changes.

Commands:

```bash
terraform plan -target='<resource-address>' -out one-device.plan
terraform apply one-device.plan
terraform plan
```

## Guardrails

- Do not run `terraform init` when provider `dev_overrides` are active in Terraform CLI config.
- If apply fails with DNS/hostname lookup errors, switch endpoint host to a reachable IP and keep alias unchanged.
- For troubleshooting, verify SSH or NETCONF reachability before planning/applying.
