# Using OpenSpec with JTAF

OpenSpec is a spec-driven development workflow integrated into this project. It helps you plan, design, and implement changes through structured artifacts and AI-assisted chat commands.

## Prerequisites

- VS Code with GitHub Copilot Chat extension
- Go, Python 3.9+, Git, and Terraform installed
- The [Juniper YANG models repo](https://github.com/Juniper/yang) cloned alongside this repo

> **Important:** You must open the `junos-terraform/` folder as your workspace root in VS Code. The prompt files at `.github/prompts/` are resolved relative to the workspace root — if you open a parent folder instead, the `/opsx-*` commands won't appear in Copilot Chat.

### Environment Setup

```bash
git clone https://github.com/juniper/junos-terraform
git clone https://github.com/juniper/yang
cd junos-terraform
python3 -m venv venv
. venv/bin/activate
pip install -e .
```

If you don't have Terraform installed:

```bash
# macOS
brew tap hashicorp/tap
brew install hashicorp/tap/terraform
```

### Install OpenSpec CLI

```bash
pip install openspec
```

### Verify Your Setup

Run these checks to confirm everything is ready:

```bash
# OpenSpec CLI is installed
openspec --version

# Prompt files are in place (should list 4 .prompt.md files)
ls .github/prompts/

# Python environment is active and JTAF is installed
which jtaf-provider

# Go toolchain is available
go version
```

In VS Code, verify the prompts are loaded in Copilot Chat:
1. Open the `junos-terraform/` folder directly in VS Code (`File > Open Folder...`)
2. Open Copilot Chat (⌘⇧I or Ctrl+Shift+I)
3. Type `/` in the chat input
4. You should see `opsx-explore`, `opsx-propose`, `opsx-apply`, and `opsx-archive` in the prompt list

If the prompts don't appear:
- Make sure `junos-terraform/` is the workspace root (not a parent folder)
- Check that `.github/prompts/` exists with 4 `.prompt.md` files
- Restart VS Code if you just cloned the repo

## Quick Start

1. **Explore** a problem or idea (read-only, no code changes):
   ```
   /opsx-explore how should we handle YANG union types in the diff engine?
   ```

2. **Propose** a change when you're ready to commit:
   ```
   /opsx-propose handle-union-types
   ```
   This generates `openspec/changes/handle-union-types/` with:
   - `proposal.md` — what and why
   - `design.md` — how (architecture, files to modify)
   - `tasks.md` — implementation checklist

3. **Apply** the change (implements tasks one-by-one):
   ```
   /opsx-apply
   ```

4. **Archive** when complete:
   ```
   /opsx-archive
   ```

## Commands Reference

| Command | Purpose | Writes code? |
|---------|---------|:---:|
| `/opsx-explore <topic>` | Investigate, brainstorm, compare approaches | No |
| `/opsx-propose <name>` | Create a change with proposal, design, tasks | No (artifacts only) |
| `/opsx-apply [name]` | Implement tasks from a change | Yes |
| `/opsx-archive [name]` | Archive a completed change | No |

## Directory Structure

```
junos-terraform/
├── .github/
│   └── prompts/             # Copilot Chat slash commands (committed with repo)
│       ├── opsx-apply.prompt.md
│       ├── opsx-archive.prompt.md
│       ├── opsx-explore.prompt.md
│       └── opsx-propose.prompt.md
└── openspec/
    ├── config.yaml          # Project config (schema, conventions)
    ├── changes/             # Active changes
    │   ├── <change-name>/
    │   │   ├── .openspec.yaml
    │   │   ├── proposal.md
    │   │   ├── design.md
    │   │   └── tasks.md
    │   └── archive/         # Completed changes (date-prefixed)
    │       └── 2026-05-06-handle-union-types/
    └── specs/               # Accumulated project specs
        ├── cli-tools/
        ├── code-generation/
        ├── netconf-mock/
        ├── patch-engine/
        ├── pyang-plugin/
        ├── terraform-provider/
        └── testing/
```

## Workflow Examples

### Adding a new YANG type to the patch engine

```
/opsx-explore the patch engine doesn't handle identityref leaves correctly

  → Agent investigates terraform_provider/patch/, finds the gap,
    suggests approaches

/opsx-propose fix-identityref-handling

  → Generates artifacts:
    proposal.md: "Add identityref resolution to LeafMapWithSchema..."
    design.md: "Modify leafmap.go type switch, add namespace lookup..."
    tasks.md:
      - [ ] Add identityref detection in schema type lookup
      - [ ] Implement namespace-qualified value comparison
      - [ ] Add test cases with identity values
      - [ ] Update diff to normalize identityref comparisons

/opsx-apply

  → Implements each task, marks [x] as complete

/opsx-archive
```

### Fixing a bug in provider generation

```
/opsx-propose fix-empty-container-diff

  → You can skip explore and go straight to propose if you
    already know what needs to change

/opsx-apply fix-empty-container-diff

  → Specify the change name if multiple are active
```

## Tips

- **Explore is safe** — it never modifies code, only reads and analyzes
- **Propose captures context** — artifacts persist across chat sessions so you don't lose your plan
- **Apply is resumable** — stop anytime; progress is tracked via checkboxes in tasks.md
- **Use explore for reviews** — great for understanding unfamiliar code before changing it
- **Name changes in kebab-case** — e.g., `add-vrrp-ordering`, `fix-empty-leaf-diff`

## For New Contributors

1. Clone the repo and switch to the `netconf_patch` branch
2. Open the `junos-terraform/` folder directly in VS Code (`File > Open Folder...`)
3. Verify prompts are available: type `/` in Copilot Chat and look for `opsx-*` commands
4. Run `/opsx-explore` to understand any area of the codebase
5. Check `openspec/changes/` for any in-progress work
6. Start your own change with `/opsx-propose`

## Configuration

The project config lives at `openspec/config.yaml`:

```yaml
project:
  name: junos-terraform
  description: >
    Framework that generates custom Junos Terraform providers from YANG models.
  languages: [python, go]
  frameworks: [terraform, pyang, netconf]
  conventions:
    - Go files use gofmt; tests use testing package and testify
    - Python follows PEP 8; tests use pytest
    - NETCONF XML paths use Junos hierarchy notation
    - Commit messages follow conventional commits format

schema: spec-driven
```

Edit this file to update project conventions that guide artifact generation.

## Project Specs

Specs document every component of JTAF. They live at `openspec/specs/<component>/spec.md` and give the AI agent context when you run `/opsx-propose` or `/opsx-explore`.

| Spec | Path | What It Covers |
|------|------|----------------|
| **CLI Tools** | `specs/cli-tools/spec.md` | Python entry points (jtaf-provider, jtaf-yang2go, jtaf-xml2tf, etc.), jtaf_common.py functions |
| **Patch Engine** | `specs/patch-engine/spec.md` | LeafMap, Diff, Patch, Tree, Order — the Go diff/apply engine in `terraform_provider/patch/` |
| **Terraform Provider** | `specs/terraform-provider/spec.md` | Provider struct, Config, NETCONF transport, resource CRUD lifecycle |
| **Code Generation** | `specs/code-generation/spec.md` | Jinja2 templates, template variables, what gets generated |
| **pyang Plugin** | `specs/pyang-plugin/spec.md` | YANG→JSON conversion, FNode classes, type mapping |
| **NETCONF Mock** | `specs/netconf-mock/spec.md` | Mock NETCONF server for testing without live devices |
| **Testing** | `specs/testing/spec.md` | All Go and Python test suites, patterns, test data locations |

### How Specs Work With Commands

Specs provide background context. When you run a command, the agent reads the relevant specs automatically:

```
/opsx-explore why does the diff engine produce spurious deletes for empty containers?

  → Agent reads specs/patch-engine/spec.md for architecture context
  → Investigates terraform_provider/patch/diff.go and leafmap.go
  → Gives an informed answer grounded in how the engine actually works
```

```
/opsx-propose add-snmp-community-support

  → Agent reads specs/cli-tools/spec.md, specs/code-generation/spec.md,
    specs/terraform-provider/spec.md to understand the full pipeline
  → Generates proposal/design/tasks that follow existing patterns
```

You can also point the agent at a specific spec:

```
/opsx-explore review the testing spec — are we missing coverage for the NETCONF client?
```

### Creating a New Spec

When you add a major feature or component, create a spec so future changes have context.

**1. Create the directory and file:**

```bash
mkdir -p openspec/specs/<component-name>
```

**2. Write `spec.md` with these sections:**

```markdown
# Component Name

Brief description of what this component does and where it lives.

## Architecture

How it fits into the overall system. Diagrams welcome.

## Files

| File | Purpose |
|------|---------|
| `path/to/file.go` | What this file does |

## Key Types / Functions

Important interfaces, structs, functions — what a developer needs to know.

## Usage

How to use this component (commands, API calls, examples).

## Tests

Where the tests live and what they cover.
```

**3. Follow these guidelines:**

- **One spec per component** — don't combine unrelated features
- **Name the directory in kebab-case** — e.g., `snmp-support`, `bgp-filter-engine`
- **Focus on what an AI agent needs** — architecture, file locations, key types, corner cases
- **Keep it current** — update the spec when you change the component
- **Don't duplicate code** — reference files and functions by name, don't paste source code

**4. Commit it:**

```bash
git add openspec/specs/<component-name>/spec.md
```

### Example: Adding a Spec for a New Feature

Say you're adding SNMP community management to the provider:

```bash
mkdir -p openspec/specs/snmp-support
```

Then create `openspec/specs/snmp-support/spec.md`:

```markdown
# SNMP Support

Terraform resource for managing Junos SNMP community strings via NETCONF.

## Architecture

Extends the generated provider with an SNMP-specific resource that maps
`configuration/snmp/community` YANG paths to Terraform schema attributes.

## Files

| File | Purpose |
|------|---------|
| `terraform_provider/resource_snmp.go` | SNMP community resource (CRUD) |
| `terraform_provider/resource_snmp_test.go` | Unit tests |
| `examples/snmp_config.xml` | Sample SNMP XML config |

## Key Types

- `snmpCommunityResource` — implements Terraform resource interface
- `snmpCommunityModel` — Terraform state model (name, authorization, clients)

## YANG Path

`configuration/snmp/community[name=<community>]`

## Tests

- `resource_snmp_test.go` — CRUD operations, authorization levels, client restrictions
```

Now when someone runs `/opsx-propose extend-snmp-v3`, the agent has full context about the existing SNMP implementation.
