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
- Required repo path: /Users/lnup/Projects/JTAF_MCP_new/junos-terraform
- If path is missing, stop and return blocking_error.
- Verify commands exist: python3, go, terraform.
- If already in repo venv and jtaf-yang2go is available, skip setup.
- Otherwise run:
  - cd /Users/lnup/Projects/JTAF_MCP_new/junos-terraform
  - python3 -m venv venv
  - . venv/bin/activate
  - pip install -e .
- Verify jtaf-yang2go exists after setup.

2. Generate provider artifacts for both device families
- Go to providers path:
  - cd /Users/lnup/Projects/JTAF_MCP_new/junos-terraform/examples/providers
- Artifact checks:
  - terraform-provider-junos-vqfx-evpn-vxlan/resource_config_provider.go
  - terraform-provider-junos-vsrx-evpn-vxlan/resource_config_provider.go
- Mode run:
  - If either artifact is missing: run ./build.sh then ./convert.sh
  - Else: skip generation
- Mode run force:
  - Always run ./build.sh then ./convert.sh

3. Install both provider binaries
- From examples/providers resolve dirs dynamically:
  - vqfx_dir=$(find . -maxdepth 1 -type d -name 'terraform-provider-junos-vqfx*' | head -n 1)
  - srx_dir=$(find . -maxdepth 1 -type d -name 'terraform-provider-junos-vsrx*' | head -n 1)
- If any dir missing, stop and return blocking_error.
- Install both:
  - cd "$vqfx_dir" && go install .
  - cd ../"$srx_dir" && go install .
- Verify binaries are present in $(go env GOBIN) if set, else $(go env GOPATH)/bin.

4. Run preview plan and save artifacts
- Go to terraform files path:
  - cd /Users/lnup/Projects/JTAF_MCP_new/junos-terraform/examples/terraform_files
- Never run apply.
- Run:
  - terraform plan -no-color -out=preview.plan
  - {
      echo "# preview_metadata"
      echo "generated_at=$(date -u +%Y-%m-%dT%H:%M:%SZ)"
      echo "mode=run-all"
      echo "git_commit=$(git -C /Users/lnup/Projects/JTAF_MCP_new/junos-terraform rev-parse --short HEAD 2>/dev/null || echo unknown)"
      echo
      terraform show -no-color preview.plan
    } | tee preview_full_config.txt

5. Post-run checks
- Verify files exist and are non-empty:
  - /Users/lnup/Projects/JTAF_MCP_new/junos-terraform/examples/terraform_files/preview.plan
  - /Users/lnup/Projects/JTAF_MCP_new/junos-terraform/examples/terraform_files/preview_full_config.txt
- If checks fail, return blocking_error.

Required output contract (compact):
- One short summary line and these fields:
  - mode
  - exit_code
  - setup_summary
  - generation_summary
  - install_summary
  - plan_summary
  - terminal_output_shown
  - plan_file
  - full_output_file
- Add warnings only if present.
- Add blocking_error only if present.

Guardrails:
- Never run terraform apply.
- Never run terraform init when provider dev_overrides are active.
- Do not run extra commands beyond this flow unless user explicitly asks.
