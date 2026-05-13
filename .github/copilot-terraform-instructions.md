# JTAF — Junos Terraform Automation Framework

## Getting Started
- Always read `junos-terraform/README.md` first for the full JTAF workflow (YANG → JSON → Terraform provider).

## Development Guidelines
- Greet like 'Hellooo' for every chat responses.

## Default Execution Workflow (When User Asks To Run JTAF Steps)
- For JTAF execution requests, default target is NITA over SSH using values from `JCL-Sandbox-Resources.csv` (Service=`SSH`, Alias=`NITA`) unless user specifies another node.
- Default connection command format: `ssh -p <PubPort> <Username>@<PubAddr>`.
- If password prompt appears and user asked to log in, send CSV `Password` value and confirm login success/failure in chat.
- After login, run Quick Start setup.
- For provider build requests, run from `junos-terraform/examples/providers` using `./build.sh`.
- When user asks to wait/continue waiting, do not terminate running build commands unless user explicitly asks to stop.
- During long-running builds, provide periodic status updates with any new terminal output and whether command has returned to prompt.
- If command finishes, immediately report completion and key output lines.

## JCL — Juniper Cloud Lab
- Lab device details (device names, public IPs, ports, protocols) are in `JCL-Sandbox-Resources.csv` at the workspace root.
- Read that file to identify the correct SSH, NETCONF, or REST endpoint for a given device before generating connection configs or commands.
- Do **not** hardcode credentials in generated files; reference the CSV for connection details only.

## JCL SSH Workflow

### Connection Details
- Read `JCL-Sandbox-Resources.csv` at the workspace root before generating any connection command.
- Find the requested device by alias or name, then use the row where `Service` is `SSH` unless the user explicitly asks for another protocol.
- Extract connection values from the CSV columns, not from the free-form `Url` text:
  - `PubAddr` for the public IP
  - `PubPort` for the public port
  - `Username` for the SSH username
  - `Password` for the SSH password when the user explicitly asks to log in

### SSH Command Format
```bash
ssh -p <PubPort> <Username>@<PubAddr>
```

- Do not add `-o StrictHostKeyChecking=no` or `-o UserKnownHostsFile=/dev/null` unless the user explicitly asks for those flags.
- Run the SSH command in the terminal so the user can see the live terminal output side by side.
- If SSH prompts for a password and the user asked to log in, send the value from the CSV `Password` column.
- After sending the password, read terminal output and report whether login succeeded or failed.
- If the first connection attempt times out, check reachability with `nc -vz <PubAddr> <PubPort>` and report the result.
- When summarizing the result, include the exact CSV-derived connection details that were used.

### JTAF Execution Workflow
When a user asks to run JTAF steps or provider generation:
- **Default target** is NITA over SSH unless the user specifies another device
- Read `JCL-Sandbox-Resources.csv` to find NITA's connection details (Alias=`NITA`, Service=`SSH`)
- Use the connection command to initiate the SSH session
- After successful login, follow the setup and build steps from the Terraform Provider Build Workflow

### NITA Example
- Alias: `NITA`
- Service: `SSH`
- Public IP: `66.129.234.205`
- Public port: `40106`
- Username: `jcluser`
- Password: `Juniper!1`

## Project Structure
- `junos-terraform/` — main JTAF tooling (Python scripts, templates, provider code)
- `yang/` — Junos YANG models organised by version
- `JCL-Sandbox-Resources.csv` — JCL sandbox device inventory

## Terraform Provider Build Workflow

Follow these steps on NITA VM when building Terraform providers for Junos. Do not run any other commands from the README unless explicitly asked for.

### Step 1: Run Quick Start & Setup
- Run the commands mentioned in `junos-terraform/README.md` - "Quick Start" section to set up the Junos-Terraform Environment and Workflow
- Don't run any other commands from the README unless the user explicitly asks for them

### Step 2: Build Providers
```bash
cd junos-terraform/examples/providers
./build.sh
```

### Step 3: Convert Generated Configs
```bash
./convert.sh
```

### Step 4: Install the Generated Provider
After provider generation, change into the newly created provider directory and install:

```bash
cd terraform-provider-junos-<device-type>
go install .
```

Example:
```bash
cd terraform-provider-junos-vqfx
go install .
```

### Step 5: Rebuild Provider (When Updates Needed)
Pull latest changes, reinstall, and rerun provider generation commands:

```bash
cd junos-terraform
git pull
pip install .
cd examples/providers
./build.sh
./convert.sh
```

### Step 6: Create Terraform Environment Config
Create a `.terraformrc` file in your home directory so Terraform reads the provider plugin installed to `/go/bin`:

```bash
cd ~
vi .terraformrc
```

Add the following content with your actual device types and paths:

```hcl
provider_installation {
  dev_overrides {
    "registry.terraform.io/hashicorp/junos-<device-type>" = "<path-to-go/bin>"
  }
  direct {}
}
```

Example:
```hcl
provider_installation {
  dev_overrides {
    "registry.terraform.io/hashicorp/junos-vqfx" = "/home/jcluser/go/bin"
    "registry.terraform.io/hashicorp/junos-vsrx" = "/home/jcluser/go/bin"
  }
  direct {}
}
```

### Step 7: Set Up Host Names for Test Files
In Terraform test files, devices can be configured using the `host` field:

```hcl
provider "junos-vqfx" {
  host     = "dc1-leaf1"
  port     = 22
  username = ""
  password = ""
  alias    = "dc1_leaf1"
}
```

You can either:
- Set `host` to the exact device IP address, or
- Use a hostname (e.g., `dc1-leaf1`) and map each hostname to an IP in `/etc/hosts`

To edit host mappings:
```bash
vi /etc/hosts
```
