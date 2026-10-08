# JUNOS Terraform Automation Framework (JTAF)

Terraform is traditionally used for managing virtual infrastructure, but there are organisations out there that use Terraform end-to-end and also want to manage configuration state using the same methods for managing infrastructure. Sure, we can run a provisioner with Terraform, but that wasn't asked for!

Much the same as you can use Terraform to create an AWS EC2 instance, you can manage the configurational state of Junos. In essence, we treat Junos configuration as declarative resources.

So what is JTAF? It's a framework, meaning, it's an opinionated set of tools and steps that allow you to go from YANG models to a custom Junos Terraform provider. With all frameworks, there are some dependencies.

You need **Python 3.9+ and Git**, plus **Go and Terraform** for the Terraform
workflow or **Ansible** for the Ansible workflow. The commands below use a Bash
shell on Linux or macOS; on Windows, use WSL or adapt the environment activation
and paths for your shell.

## Setup

Install the required tools, then clone JTAF and install its Python commands:

```bash
git clone https://github.com/juniper/junos-terraform
cd junos-terraform
python3 -m venv venv
. venv/bin/activate
python -m pip install -e .
```

If you do not already have Terraform installed (in general), for macOS, run the following:
```bash
brew tap hashicorp/tap
brew install hashicorp/tap/terraform
```

For more information, refer to the Terraform website: https://developer.hashicorp.com/terraform/install.

The walkthroughs use bundled 18.2 YANG examples under `examples/yang/18.2`.
For your own devices, obtain models matching the platform and Junos release
from [Juniper's YANG repository](https://github.com/Juniper/yang).

---

## Choose Your Workflow

JTAF supports two output modes from the same YANG models and XML configurations:

| Workflow | Output | Guide |
|----------|--------|-------|
| **Terraform Provider** | Schema-driven Go Terraform provider with NETCONF patch engine | [Junos Terraform Guide](README-terraform.md) |
| **Ansible Role** | Ansible role + playbook with Jinja2 templates | [Junos Ansible Guide](README-ansible.md) |

Follow one guide from beginning to end. Each explains its model conversion,
generation commands, options, and deployment steps without requiring example
wrapper scripts. Start in the repository root with the virtual environment active.

Worked examples for the Terraform workflow:

| Example | Shows |
|---------|-------|
| [Untrimmed model](examples/DEMO-GENERIC-PROVIDER.md) | `--generic`, embedding a whole Junos model |
| [Configuration groups](examples/DEMO-GROUPS.md) | `--groups`, managing configuration that lives in Junos groups |

## Reference: How Terraform Generation Works

JTAF uses one schema-driven Go provider implementation. `jtaf-provider` compiles
your YANG model into a compact table embedded in the binary; it does not generate
thousands of Go resources for each device. A new Junos release requires rebuilding
with the matching schema.

Choose explicitly between a schema trimmed to XML with `-x` and a full model
with `--generic`. The [schema-scope decision](README-terraform.md#decide-the-schema-scope)
and performance tradeoffs are explained in the Terraform walkthrough.
