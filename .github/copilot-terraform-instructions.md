# JTAF - Junos Terraform Automation Framework

## Getting Started
- Read `junos-terraform/README.md` first for the full JTAF workflow (YANG -> JSON -> Terraform provider).

## General Rules
- Keep responses environment-neutral by default.
- Do not assume specific device names, aliases, hostnames, usernames, or passwords unless the user provides them.
- If concrete values are shown for illustration, mark them explicitly as examples.

## Execution Workflow (When User Asks To Run JTAF Steps)
- Ask for the target device or alias if user intent does not include one.
- Build SSH command using resolved values: `ssh -p <PubPort> <Username>@<PubAddr>`.
- If password prompt appears and user asked to log in, send the resolved password value.
- After login, run Quick Start setup.
- For provider build requests, run from `junos-terraform/examples/providers` using `./build.sh`.
- If user asks to wait, keep monitoring long-running commands and provide status updates.
- Report completion with key output lines once the command returns.

## Inventory Source
- In JCL environments, inventory can be in `JCL-Sandbox-Resources.csv` at workspace root.
- Resolve endpoint details from structured columns (for example: `PubAddr`, `PubPort`, `Username`, `Password`) instead of free-form URL text.
- Do not hardcode credentials in generated files.

## SSH Workflow

### Connection Resolution
- Read inventory file first (when available).
- Match target by alias or name.
- Prefer row where `Service=SSH` unless user requests another protocol.

### SSH Command Format
```bash
ssh -p <PubPort> <Username>@<PubAddr>
```

- Do not add extra SSH flags unless user asks.
- If first attempt times out, check reachability with `nc -vz <PubAddr> <PubPort>` and report result.
- Report exactly which resolved values were used.

## Example Block (Explicit Example)
The values below are examples only. Do not assume they apply to user environments.

- Alias: `NITA`
- Service: `SSH`
- Public IP: `66.129.234.205`
- Public port: `40106`
- Username: `jcluser`
- Password: `Juniper!1`

## Project Structure
- `junos-terraform/` - main JTAF tooling (Python scripts, templates, provider code)
- `yang/` - Junos YANG models by version
- `JCL-Sandbox-Resources.csv` - optional JCL inventory file

## Terraform Provider Build Workflow
Use these steps on the selected build VM.

### Step 1: Setup
- Run commands from `junos-terraform/README.md` Quick Start section.

### Step 2: Build Providers
```bash
cd junos-terraform/examples/providers
./build.sh
```

### Step 3: Convert Generated Configs
```bash
./convert.sh
```

### Step 4: Install Generated Provider
```bash
cd terraform-provider-junos-<device-type>
go install .
```

Example:
```bash
cd terraform-provider-junos-vqfx
go install .
```

### Step 5: Rebuild After Updates
```bash
cd junos-terraform
git pull
pip install .
cd examples/providers
./build.sh
./convert.sh
```

### Step 6: Configure Terraform CLI
Create `~/.terraformrc` and add:

```hcl
provider_installation {
  dev_overrides {
    "registry.terraform.io/hashicorp/junos-<device-type>" = "<path-to-go/bin>"
  }
  direct {}
}
```

Explicit example only:
```hcl
provider_installation {
  dev_overrides {
    "registry.terraform.io/hashicorp/junos-vqfx" = "/home/<username>/go/bin"
    "registry.terraform.io/hashicorp/junos-vsrx" = "/home/<username>/go/bin"
  }
  direct {}
}
```

### Step 7: Host Names In Test Files
Example provider block (sample values):

```hcl
provider "junos-vqfx" {
  host     = "dc1-leaf1"
  port     = 22
  username = ""
  password = ""
  alias    = "dc1_leaf1"
}
```

- Use direct IPs, or
- Map hostnames in `/etc/hosts`.
