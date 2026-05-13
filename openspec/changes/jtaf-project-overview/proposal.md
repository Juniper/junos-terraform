# Proposal: Junos Terraform Automation Framework (JTAF)

## Summary

JTAF is a framework that bridges the gap between Juniper network device configuration and modern Infrastructure-as-Code (IaC) tooling. It generates custom Terraform providers (and Ansible roles) directly from YANG data models, enabling network engineers to manage Junos device configuration state using the same declarative workflows they use for cloud infrastructure.

---

## Why This Project Exists

### The Problem

Organizations adopting Terraform end-to-end want a unified approach to managing **all** infrastructure — from cloud VMs to physical network switches. However:

1. **No native Junos Terraform provider covers all configurations** — Juniper devices run hundreds of software versions across many platforms (QFX, MX, SRX, PTX, ACX, EX). A single static provider cannot cover the permutations of platform × version × config scope.

2. **Network config is not cloud infrastructure** — Traditional Terraform providers manage the lifecycle of discrete resources (create a VM, destroy a VM). Network configuration is different: it's a continuous state document managed through NETCONF, not individual resource instantiation.

3. **Manual provider development doesn't scale** — Writing Go code by hand for every Junos configuration hierarchy is infeasible. The YANG models for a single Junos version contain thousands of nodes.

4. **Config drift is a real threat** — Without declarative state management, network configs drift from intent over time. Terraform's plan/apply/destroy lifecycle naturally solves this for infrastructure; networks need the same treatment.

5. **Existing tools (Ansible, scripts) lack declarative state tracking** — They push config but don't maintain a state file or detect drift. Terraform's state model fills this gap.

### The Insight

Junos already has a machine-readable, versioned schema for its entire configuration: **YANG models**. These models are published per platform per version. By treating YANG as the single source of truth and generating Terraform providers programmatically, JTAF can:

- Cover **any** Junos version (14.2 through 25.4+)
- Support **any** platform (QFX, MX, VMX, SRX, PTX, ACX, EX)
- Scope providers to **exactly** the configuration needed (right-sizing)
- Leverage Terraform's state engine for drift detection

---

## What JTAF Does

### Core Capability

JTAF takes three inputs and produces a complete, buildable, installable Terraform provider:

| Input | Description |
|-------|-------------|
| **YANG models** | Platform-specific, version-specific schema from Juniper's published YANG repo |
| **XML configuration(s)** | Example device configs defining the scope of what the provider manages |
| **Device type label** | A name for the generated provider (e.g., `vqfx`, `vmx`) |

**Output:** A ready-to-compile Go Terraform provider at `terraform-provider-junos-{type}/`

### What Users Can Do With It

1. **Generate a provider** tailored to their exact network topology and software version
2. **Manage Junos config as Terraform resources** — plan, apply, destroy configuration groups
3. **Detect drift** — `terraform plan` shows differences between desired state and actual device state
4. **Use standard Terraform workflows** — state files, workspaces, modules, CI/CD pipelines
5. **Generate Terraform test files** from XML configs for immediate deployment testing
6. **Generate Ansible roles** from the same YANG + XML inputs for teams preferring Ansible
7. **Apply granular updates** — the patch engine computes minimal NETCONF `edit-config` operations rather than replacing entire configuration groups

### Secondary Capabilities

| Tool | Purpose |
|------|---------|
| `jtaf-xml2tf` | Generate `.tf` test files from XML configs |
| `jtaf-ansible` | Generate Ansible role + playbook from JSON schema + XML |
| `jtaf-yang2ansible` | Single-command YANG → Ansible role |
| `jtaf-xml2yaml` | Convert XML configs to Ansible `host_vars`/`group_vars` with inventory |

---

## Intent and Vision

### Treat Network Config as Declarative Resources

Just as `aws_instance` in Terraform represents a cloud VM, JTAF generates resources that represent Junos configuration groups. A user declares their desired network state in HCL, and Terraform ensures the device matches that state.

```hcl
resource "junos-vqfx_ConfigProvider" "dc1-spine1" {
  resource_name = "evpn-underlay"
  # ... configuration attributes from YANG schema
}
```

### Right-Sized Providers (Goldilocks Principle)

JTAF doesn't try to generate a universal Junos provider with every possible configuration knob. Instead, it follows a "right-sizing" philosophy:

1. **Platform-specific** — Use the correct YANG models for your hardware
2. **Version-specific** — Use the models matching your running Junos version
3. **Scope-specific** — XML configs define exactly which configuration hierarchies are included
4. **Single-resource approach** — Each Terraform resource manages the smallest meaningful unit of config

### Configuration Groups as the Deployment Model

JTAF providers use Junos configuration **groups** as the deployment mechanism:

- Each Terraform resource creates a named group (e.g., `evpn-underlay`, `bgp-overlay`)
- Groups are applied to the device via `apply-groups`
- This provides clean separation between JTAF-managed config and manually-managed config
- Groups can be independently created, updated, and destroyed without affecting other config

### Multi-Tool Output from Single Source

The same YANG + XML inputs produce both Terraform providers and Ansible roles, allowing organizations to choose their preferred automation tool without re-doing the schema analysis work.

---

## Scope

### In Scope

- **Generation pipeline:** YANG → JSON → Go provider source code
- **Provider runtime:** NETCONF client, CRUD operations, state management
- **Patch engine:** Granular diff/apply for efficient config updates (netconf_patch branch)
- **Testing tools:** `.tf` file generation, mock NETCONF server
- **Ansible output:** Role/playbook/vars generation from same inputs
- **YANG coverage:** All published Juniper YANG versions (14.2–25.4)
- **Platform support:** QFX, MX, VMX, SRX, PTX, ACX, EX (any device with YANG models)

### Out of Scope

- **Runtime device discovery** — Users must know their device type and version
- **Provider registry hosting** — Users build and install providers locally (or via CI)
- **Operational commands** — JTAF manages configuration state, not operational/show commands
- **Multi-vendor** — JTAF is Junos-specific (YANG + NETCONF + Junos groups model)
- **GUI/web interface** — CLI-first tooling designed for automation pipelines

---

## Approach

### Framework, Not Product

JTAF is explicitly a **framework** — an opinionated set of tools and steps. It requires:
- Go (to compile the generated provider)
- Python (to run the generation tools and pyang plugin)
- Terraform (to use the generated provider)
- Git (to manage YANG models and generated code)

### Generation-Time vs. Runtime

JTAF separates concerns into two phases:

| Phase | Language | What Happens |
|-------|----------|-------------|
| **Generation** | Python | YANG models parsed, JSON schema extracted, Go code rendered via Jinja2 templates |
| **Runtime** | Go | Compiled provider binary communicates with devices via NETCONF |

This means the Python tools are development-time dependencies only — production deployments just need the compiled Go binary and Terraform.

### NETCONF as the Transport

All device communication uses NETCONF (RFC 6241) over SSH:
- `<load-configuration>` for full config replacement (Create)
- `<edit-config>` with `nc:operation` for granular patches (Update)
- `<get-configuration>` for state reads (Read)
- `<commit>` / `<discard-changes>` for transaction control
- Configuration group delete for resource destruction (Delete)

### Patch Engine for Efficient Updates

The `netconf_patch` branch adds a sophisticated diff/patch engine that:
1. Flattens XML config trees into path→value leaf maps (respecting YANG schema types)
2. Computes minimal diffs between current state and desired state
3. Generates targeted `<edit-config>` operations instead of full config replacement
4. Handles corner cases: empty leaves, ordered leaf-lists, identityrefs, multi-key lists

---

## Target Users

| User Persona | How They Use JTAF |
|--------------|-------------------|
| **Network Engineer** | Generates provider for their DC fabric, writes `.tf` files for each device role |
| **NetDevOps Team** | Integrates JTAF generation into CI/CD, runs `terraform apply` in pipelines |
| **Platform Team** | Maintains versioned providers per Junos release, distributes to dev teams |
| **Ansible User** | Uses `jtaf-yang2ansible` to get roles + playbooks without Terraform |

---

## Success Criteria

1. A user can go from YANG models to a working `terraform apply` against a Junos device in under 30 minutes
2. Generated providers compile cleanly and pass `go vet` / `go test`
3. `terraform plan` correctly detects configuration drift on managed groups
4. Updates apply minimal changes (not full-group replacement) via the patch engine
5. The same source artifacts produce both Terraform and Ansible outputs
6. Providers work across Junos versions 14.2 through 25.4+ without code changes

---

## References

- **Repository:** https://github.com/Juniper/junos-terraform
- **YANG Models:** https://github.com/Juniper/yang
- **Wiki (FAQ, Rules, Build Steps):** https://github.com/Juniper/junos-terraform/wiki
- **Video Introduction:** https://youtu.be/eH24eCZc7pE
- **Video Installation:** https://youtu.be/aTF7_Uscd9Q
- **Video Generation:** https://youtu.be/UgsFU7UplRE
- **Video Execution:** https://youtu.be/Lfkc38wzhNg
- **Video Interface Config:** https://youtu.be/iCnnkDodUgQ
- **Video BGP Config:** https://youtu.be/nQVNCNCJZRc
- **Terraform Provider Docs:** https://www.terraform.io/docs/language/providers/index.html
- **Terraform Registry:** https://registry.terraform.io/
- **NETCONF RFC 6241:** https://datatracker.ietf.org/doc/html/rfc6241
- **YANG RFC 7950:** https://datatracker.ietf.org/doc/html/rfc7950
- **License:** Apache-2.0
