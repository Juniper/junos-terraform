# JTAF — Junos Terraform Automation Framework

## Project Overview

JTAF is a framework that generates custom Junos Terraform providers from YANG models. It converts YANG schemas → JSON → Go resource providers, enabling declarative Junos configuration management via Terraform.

**Branch:** `netconf_patch` — adds NETCONF-based patch/diff engine for granular config operations.

## Architecture

```
YANG models → pyang plugin (jtaf_json.py) → JSON schema
JSON schema + XML configs → jtaf-provider → Go Terraform provider
Go provider → NETCONF → Junos device
```

### Key Components

| Directory | Purpose |
|-----------|---------|
| `junosterraform/` | Python CLI tools (jtaf-provider, jtaf-xml2tf, jtaf-yang2go) |
| `terraform_provider/` | Go provider source with patch/diff engine |
| `terraform_provider/patch/` | NETCONF patch engine (leafmap, diff, apply) |
| `jtaf_pyang_plugin/` | pyang plugin to extract YANG → JSON |
| `netconf_mock/` | Mock NETCONF server for testing |
| `examples/` | YANG files, XML configs, Terraform examples |
| `tests/` | Integration/E2E tests |

### Languages & Tools

- **Python 3.9+** — CLI tools, pyang plugin, mock server
- **Go 1.21+** — Terraform provider, patch engine
- **Terraform** — Infrastructure-as-code consumer
- **NETCONF/XML** — Device communication protocol
- **YANG** — Schema modeling language for network config

## Development Workflow

```bash
# Setup
python3 -m venv venv && . venv/bin/activate && pip install -e .

# Generate provider
pyang --plugindir $(jtaf-pyang-plugindir) -f jtaf -p <common> <yang-files> | \
  jtaf-provider -j - -x <xml-configs> -t <device-type>

# Build provider
cd terraform-provider-junos-<type> && go install .

# Run tests
cd terraform_provider && go test ./...
pytest junosterraform/unit_tests/
pytest netconf_mock/tests/
```

## Coding Conventions

- Go files use standard `gofmt` formatting
- Python follows PEP 8; uses `pytest` for testing
- XML paths use Junos hierarchy notation: `configuration/interfaces/interface[name=ge-0/0/0]`
- YANG types to handle: `empty`, `identityref`, `leafref`, `union`, `ordered-by user` leaf-lists
- NETCONF operations: `merge`, `replace`, `create`, `delete`

## Current Work (netconf_patch branch)

The patch engine (`terraform_provider/patch/`) implements:
- `LeafMapWithSchema` — flattens Junos XML config into path→value map respecting YANG schema
- `Diff` — computes minimal NETCONF edit-config operations between desired and actual state
- `Apply` — executes patch via NETCONF `edit-config` RPC

### Known Corner Cases
- Ordered leaf-lists (VRRP virtual-address) — order matters for Junos
- YANG `empty` type leaves — presence semantics (disable, vlan-tagging)
- Container presence vs. content — empty containers should not generate spurious diffs
- Multi-value keys — interfaces keyed by `[name=X]`

## Testing Strategy

- **Go unit tests:** `terraform_provider/patch/*_test.go`, `terraform_provider/*_test.go`
- **Python unit tests:** `junosterraform/unit_tests/`, `jtaf_pyang_plugin/tests/`
- **Mock NETCONF tests:** `netconf_mock/tests/`
- **E2E:** Terraform plan/apply against mock or live devices

## OpenSpec Integration

This project uses [OpenSpec](https://github.com/Fission-AI/OpenSpec) for spec-driven development. Prompt files live at `.github/prompts/` in this repo. Use slash commands to plan and track changes:

- `/opsx-propose <change>` — Create a change with proposal, specs, design, tasks
- `/opsx-explore <topic>` — Investigate before committing to a change
- `/opsx-apply` — Implement tasks from the active change
- `/opsx-archive` — Archive a completed change

Specs live in `openspec/specs/`, active changes in `openspec/changes/`.
See `docs/openspec-guide.md` for full setup and usage instructions.
