# JUNOS Terraform Automation Framework (JTAF)

Terraform is traditionally used for managing virtual infrastructure, but there are organisations out there that use Terraform end-to-end and also want to manage configuration state using the same methods for managing infrastructure. Sure, we can run a provisioner with Terraform, but that wasn't asked for!

Much the same as you can use Terraform to create an AWS EC2 instance, you can manage the configurational state of Junos. In essence, we treat Junos configuration as declarative resources.

So what is JTAF? It's a framework, meaning, it's an opinionated set of tools and steps that allow you to go from YANG models to a custom Junos Terraform provider. With all frameworks, there are some dependencies.

To use JTAF, you'll need machine that can run **Go, Python, Git and Terraform.** This can be Linux, OSX or Windows. Some easy to consume videos are below.

### How it works

JTAF no longer generates Go source for your device. There is **one** provider, written once, that
reads a YANG schema into memory at startup and builds its Terraform schema from it. `jtaf-provider`
compiles your model into a compact table, embeds it in the binary, and the provider reads it in a few
milliseconds. That means the build does not get slower as the model gets bigger, and supporting a new
Junos release is a matter of regenerating the schema rather than regenerating and recompiling
thousands of lines of Go.

The one decision you have to make is **how much of the model to embed** — see
[Schema scope](README-terraform.md#schema-scope). You must choose explicitly; there is no default.

## Quick start

### <u>Setup</u>
Run the following commands to set up the Junos-Terraform Environment and Workflow

```bash
git clone https://github.com/juniper/junos-terraform
git clone https://github.com/juniper/yang
cd junos-terraform
python3 -m venv venv
. venv/bin/activate
pip install -e .
```

If you do not already have Terraform installed (in general), for macOS, run the following:
```bash
brew tap hashicorp/tap
brew install hashicorp/tap/terraform
```

For more information, refer to the Terraform website: https://developer.hashicorp.com/terraform/install.

---
### <u>Yang File(s) to JSON Conversion</u>

Find the device's Junos Version that is running, and locate the corresponding yang and common folders. Run the below `pyang` command to generate a `.json` file containing `.yang` information for that version. [See below example for Junos version 18.2]
```
pyang --plugindir $(jtaf-pyang-plugindir) -f jtaf -p <path-to-common> <path-to-yang-files> > junos.json
```
Example: 
```
pyang --plugindir $(jtaf-pyang-plugindir) -f jtaf -p examples/yang/18.2/18.2R3/common examples/yang/18.2/18.2R3/junos-qfx/conf/*.yang > junos.json
```

NOTE: This repository includes YANG examples for 18.2 under `examples/yang/18.2`.

---

## Choose Your Workflow

JTAF supports two output modes from the same YANG models and XML configurations:

| Workflow | Output | Guide |
|----------|--------|-------|
| **Terraform Provider** | Custom Go Terraform provider with NETCONF patch engine | [Junos Terraform Guide](README-terraform.md) |
| **Ansible Role** | Ansible role + playbook with Jinja2 templates | [Junos Ansible Guide](README-ansible.md) |

Both workflows start from the same YANG → JSON conversion above, then diverge into their respective toolchains.

### Choose how much of the model to embed

The Terraform workflow needs exactly one of these. Passing neither, or both, is an error.

| | What you get | Choose it when |
|---|---|---|
| `-x <xml>` | A schema trimmed to the configuration in your XML | You know which stanzas you manage |
| `--generic` | The whole Junos model | You want to configure anything without rebuilding |

Trimming keeps the provider smallest; the full model means you never regenerate just to reach a new
stanza. Either way the provider source and the build time are the same. Costs and timings are in
[Schema scope](README-terraform.md#schema-scope).

Worked examples for the Terraform workflow:

| Example | Shows |
|---------|-------|
| [Untrimmed model](examples/DEMO-GENERIC-PROVIDER.md) | `--generic`, embedding a whole Junos model |
| [Configuration groups](examples/DEMO-GROUPS.md) | `--groups`, managing configuration that lives in Junos groups |
