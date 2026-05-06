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
