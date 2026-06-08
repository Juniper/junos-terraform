---
description: Generate or regenerate Junos Terraform provider artifacts with a unified build flow (build.sh + convert.sh) without running Terraform plan/apply.
---

Generate Terraform provider artifacts.

Use when:
- You want to regenerate provider code/artifacts after schema, XML, or template changes.
- You want a deterministic build flow from the repo scripts.

Accepted input forms (strict):
- /generate-provider generate
- /generate-provider generate force
- /generate-provider generate install vqfx
- /generate-provider generate install srx
- /generate-provider generate install all
- /generate-provider generate install all force

If input is not one of these forms, stop and ask user to choose one accepted form.

Execution rules:

1. Run environment setup only when needed
- First `cd junos-terraform`.
- If already inside an active venv for this repo and `jtaf-yang2go` is available, skip setup commands.
- Otherwise run:
  python3 -m venv venv
  . venv/bin/activate
  pip install -e .

Dependency checks (fail fast):
- Before running generation or install, verify required commands exist: `go`, `python3`.
- After setup/activation, verify `jtaf-yang2go` exists.
- If any command is missing, stop and return `blocking_error` with the missing command name.

2. Resolve working directory deterministically
- Required path: examples/providers
- If current directory is not this path, change into it first.
- If path does not exist, stop and report missing path.

3. Resolve generation mode
- `generate` runs the common repo generation scripts once for all providers.
- `generate install <device>` runs generation, then installs the selected provider binary.
- `generate install all` runs generation, then installs both provider binaries.
- If provider artifacts are already present for requested install targets, skip generation and run install only.
- `generate force` and `generate install all force` always run `./build.sh` and `./convert.sh` even if artifacts exist.

Provider artifact presence checks:
- vqfx: `examples/providers/terraform-provider-junos-vqfx-evpn-vxlan/resource_config_provider.go`
- srx: `examples/providers/terraform-provider-junos-vsrx-evpn-vxlan/resource_config_provider.go`

Provider install directory resolution:
- Resolve provider directories dynamically from `examples/providers`:
  - vqfx dir pattern: `terraform-provider-junos-vqfx*`
  - srx dir pattern: `terraform-provider-junos-vsrx*`
- If no matching directory is found for a requested install target, stop and return `blocking_error`.

4. Run generation-only commands
- Never run terraform apply in this prompt.
- Do not run terraform plan in this prompt.
- Do not run manual preflight checks.
- Show command output directly in terminal/chat.

Mode: generate
If either provider artifact is missing, run:
./build.sh
./convert.sh
Else:
Skip generation (already generated).

Mode: generate force
./build.sh
./convert.sh

Mode: generate install vqfx
If vqfx artifact is missing, run:
./build.sh
./convert.sh
Else:
Skip generation (already generated).
cd $(find . -maxdepth 1 -type d -name 'terraform-provider-junos-vqfx*' | head -n 1)
go install .

Mode: generate install srx
If srx artifact is missing, run:
./build.sh
./convert.sh
Else:
Skip generation (already generated).
cd $(find . -maxdepth 1 -type d -name 'terraform-provider-junos-vsrx*' | head -n 1)
go install .

Mode: generate install all
If either artifact is missing, run:
./build.sh
./convert.sh
Else:
Skip generation (already generated).
vqfx_dir=$(find . -maxdepth 1 -type d -name 'terraform-provider-junos-vqfx*' | head -n 1)
srx_dir=$(find . -maxdepth 1 -type d -name 'terraform-provider-junos-vsrx*' | head -n 1)
cd "$vqfx_dir"
go install .
cd ../"$srx_dir"
go install .

Mode: generate install all force
./build.sh
./convert.sh
vqfx_dir=$(find . -maxdepth 1 -type d -name 'terraform-provider-junos-vqfx*' | head -n 1)
srx_dir=$(find . -maxdepth 1 -type d -name 'terraform-provider-junos-vsrx*' | head -n 1)
cd "$vqfx_dir"
go install .
cd ../"$srx_dir"
go install .

Post-run sanity checks:
- After generation (when run), verify these files exist and are non-empty:
  - `terraform-provider-junos-vqfx-evpn-vxlan/resource_config_provider.go`
  - `terraform-provider-junos-vsrx-evpn-vxlan/resource_config_provider.go`
- After install commands, verify binaries exist under `$(go env GOPATH)/bin` or `$(go env GOBIN)` for both provider names.
- If any check fails, return `blocking_error`.

5. Required output contract (compact by default)
- Return only 1 short summary line and 4 fields:
  mode, exit_code, generation_summary, terminal_output_shown.
- Add warnings only if present.
- Add blocking_error only if present.
- Do not include extra command output unless user asks for full details.

Guardrails:
- Never run terraform apply in this prompt.
- Never run terraform plan in this prompt.
- Do not run terraform init when provider dev_overrides are active.
- Do not run extra commands or checks beyond the requested generation flow.
