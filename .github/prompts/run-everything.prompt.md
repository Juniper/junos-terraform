---
description: Run complete Junos Terraform workflow end-to-end (setup, generate, install, preview) without apply.
---

Run the full automation flow in one prompt.

Use when:
- You want one command to set up env, regenerate providers, install binaries, and run a preview plan.

Accepted input forms (strict):
- /run-everything run
- /run-everything run force

If input is not one of these forms, stop and ask user to choose one accepted form.

Example Execution rules:

1. Resolve repo and dependencies
- Required repo path: junos-terraform
- If path is missing, stop and return blocking_error.
- Verify commands exist: python3, go, terraform.
- If already in repo venv and jtaf-yang2go is available, skip setup.
- Otherwise run:
  - cd junos-terraform
  - python3 -m venv venv
  - . venv/bin/activate
  - pip install -e .
- Verify jtaf-yang2go exists after setup.

2. Generate and exercise both provider implementations
- Go to the providers path:
  - cd junos-terraform/examples/providers
- Always run `bash ./test-both-providers.sh`, for both accepted modes. It generates and builds standard and generic variants into isolated temporary directories, converts the Terraform scenario once, and runs the same plan first with standard then generic provider dev overrides.
- The helper does not install either provider globally and never runs `terraform init` or `terraform apply`.
- It saves per-variant outputs under `junos-terraform/examples/terraform_files/`:
  - `preview_full_config.txt` for the standard provider
  - `preview_generic_full_config.txt` for the generic provider
  - `preview.plan` and `preview_generic.plan` when the corresponding plans succeed
- If either plan fails, report both exit statuses and retain the per-variant output; do not stop before attempting the second plan.
- Verify both text result files exist and are non-empty, and verify each binary plan file exists when its plan succeeded. If generation fails, return `blocking_error`.

Required output contract (compact):
- One short summary line and these fields:
  - mode
  - exit_code
  - setup_summary
  - generation_summary
  - standard_plan_exit_code
  - generic_plan_exit_code
  - terminal_output_shown
  - standard_plan_file
  - standard_full_output_file
  - generic_plan_file
  - generic_full_output_file
- Add warnings only if present.
- Add blocking_error only if present.

Guardrails:
- Never run terraform apply.
- Never run terraform init when provider dev_overrides are active.
- Do not run extra commands beyond this flow unless user explicitly asks.
