# Demo: Generic Schema-Driven Terraform Provider

This walkthrough demonstrates the new `--generic` provider workflow using the
EVPN-VXLAN example configs already in this repository. Every step uses the
existing XML files in `examples/evpn-vxlan-dc/` and the YANG models in
`examples/yang/18.2/`.

---

## Prerequisites

```bash
# From the repository root
cd junos-terraform

# Activate the Python virtual environment
source ../venv/bin/activate      # adjust path to your venv

# Install / update JTAF CLI tools
pip install -e .

# Verify everything is available
which pyang jtaf-provider jtaf-yang2go terraform go
pip install asyncssh   # needed for mock NETCONF server
```

---

## Part 1 — Build the Generic Provider

### 1.1 Generate and compile (filtered schema, same XMLs as build.sh)

```bash
cd examples/providers

bash build-generic.sh
```

This runs `jtaf-yang2go --generic` with the same YANG files and XML configs
that `build.sh` uses. It produces two providers:

| Provider | YANG Platform | XML Configs Used | Binary Size |
|----------|---------------|------------------|-------------|
| `terraform-provider-junos-vqfx-evpn-vxlan` | QFX (junos-qfx) | dc1 spines + leaves | ~26 MB |
| `terraform-provider-junos-vsrx-evpn-vxlan` | SRX (junos-es) | dc1/dc2 firewalls | ~26 MB |

Both are compiled and installed to `$GOPATH/bin/`.

### 1.2 Verify the binary works

```bash
terraform-provider-junos-vqfx-evpn-vxlan --help
```

Expected output:
```
Usage of terraform-provider-junos-vqfx-evpn-vxlan:
  -debug
        set to true to run the provider with support for debuggers like delve
```

---

## Part 2 — Generate Terraform Test Files

### 2.1 Run the existing convert script

```bash
cd examples/providers

bash convert.sh
```

This runs `jtaf-xml2tf` to generate `.tf` files from the XML device configs.
Output goes to `examples/terraform_files/`:

```
examples/terraform_files/
├── providers.tf              ← provider blocks with host/port/credentials
├── dc1-spine1.tf             ← config from dc1-spine1.xml
├── dc1-spine2.tf             ← config from dc1-spine2.xml
├── dc1-leaf1.tf              ← config from dc1-leaf1.xml
├── dc1-leaf2.tf
├── dc1-leaf3.tf
├── dc1-borderleaf1.tf
├── dc1-borderleaf2.tf
├── dc2-spine1.tf
├── dc2-spine2.tf
├── dc1-firewall1.tf          ← SRX configs
├── dc1-firewall2.tf
├── dc2-firewall1.tf
└── dc2-firewall2.tf
```

These `.tf` files are identical to what the old Jinja2 provider expects —
the generic provider is fully backward-compatible.

---

## Part 3 — Set Up Terraform

### 3.1 Create `~/.terraformrc`

```bash
cat > ~/.terraformrc << EOF
provider_installation {
  dev_overrides {
    "hashicorp/junos-vqfx-evpn-vxlan" = "$(go env GOPATH)/bin"
    "hashicorp/junos-vsrx-evpn-vxlan" = "$(go env GOPATH)/bin"
  }
  direct {}
}
EOF
```

This tells Terraform to use your locally built providers instead of
the HashiCorp registry. Same pattern as `examples/example-terraformrc`.

### 3.2 Validate the configuration

```bash
cd examples/terraform_files

terraform validate
```

Expected:
```
Success! The configuration is valid.
```

This confirms the generic provider's dynamic schema (built from JSON at
startup) accepts every attribute in the existing `.tf` files — `system`,
`interfaces`, `protocols`, `policy-options`, `routing-options`, etc.

---

## Part 4 — Test Against Mock NETCONF Server

### 4.1 Start the mock server

Open a **second terminal**:

```bash
cd junos-terraform

python netconf_mock/netconf_mock_server.py \
  --host 127.0.0.1 \
  --username jcluser \
  --password 'Juniper!1' \
  --disable-auth \
  --device dc1-spine1:8301 \
  --device dc1-spine2:8302 \
  --log-level INFO
```

You should see:
```
INFO:netconf-mock:Device dc1-spine1 listening on 127.0.0.1:8301
INFO:netconf-mock:Device dc1-spine2 listening on 127.0.0.1:8302
```

### 4.2 Create a test workspace pointing at the mock

Back in the **first terminal**:

```bash
mkdir -p examples/demo-generic-test
cd examples/demo-generic-test

# Copy a single device's .tf file
cp ../terraform_files/dc1-spine1.tf .

# Create a provider block pointing at the mock server
cat > providers.tf << 'EOF'
terraform {
  required_providers {
    junos-vqfx-evpn-vxlan = {
      source = "hashicorp/junos-vqfx-evpn-vxlan"
    }
  }
}

provider "junos-vqfx-evpn-vxlan" {
  host     = "127.0.0.1"
  port     = 8301
  username = "jcluser"
  password = "Juniper!1"
  alias    = "dc1_spine1"
}
EOF
```

### 4.3 Validate

```bash
terraform validate
```

Expected: `Success! The configuration is valid.`

### 4.4 Plan

```bash
terraform plan
```

Terraform will:
1. Load the generic provider binary
2. Build the schema dynamically from embedded `trimmed_schema.json`
3. Connect to the mock on `127.0.0.1:8301` via NETCONF/SSH
4. Read the current device config (empty on mock)
5. Show a plan with all the `dc1-spine1` config to be created

### 4.5 Apply

```bash
terraform apply -auto-approve
```

This sends the full configuration to the mock via `load-configuration`
and commits it.

### 4.6 Verify state

```bash
terraform show
```

The state should contain the full dc1-spine1 configuration — system,
interfaces, protocols, routing-options, etc.

### 4.7 Test an update

Edit `dc1-spine1.tf` to change something:

```bash
# Change the hostname
sed -i '' 's/host_name = "dc1-spine1"/host_name = "dc1-spine1-updated"/' dc1-spine1.tf

terraform plan    # shows the diff
terraform apply -auto-approve
```

The update uses the **patch engine**: it computes a minimal diff between
the plan and device state, then sends only the changed leaves via
`edit-config` — not a full config replace.

### 4.8 Test destroy

```bash
terraform destroy -auto-approve
```

This diffs the current state against an empty config and sends delete
operations for every managed leaf.

### 4.9 Clean up

```bash
cd ..
rm -rf demo-generic-test
```

Stop the mock server in the second terminal with `Ctrl+C`.

---

## Part 5 — Build Full-Model Provider (No XML Filtering)

This is the key capability enabled by the generic architecture.

### 5.1 Generate the full-model provider

```bash
cd examples/providers

bash build-generic-full.sh
```

This runs `jtaf-yang2go --generic` with all 33 QFX YANG files and **no
`-x` flag** — no XML filtering. pyang takes ~2 minutes to process the
full model.

Expected output:
```
Schema:  268M
Binary:  294M
```

### 5.2 What this means

| Aspect | Filtered Provider | Full-Model Provider |
|--------|-------------------|---------------------|
| Schema coverage | Only nodes present in your XML configs | Every QFX YANG node |
| Binary size | ~26 MB | ~294 MB |
| When to rebuild | When you change which XML configs to cover | Only when Juniper releases a new YANG version |
| `.tf` attributes | Only attributes from your XMLs | Every possible Junos QFX attribute |

The full-model provider is a **build once, use forever** binary. You only
need to regenerate when a new Junos release ships new YANG models.

### 5.3 Compare with old approach

The old Jinja2 approach cannot build a full-model provider:

| Metric | Old (Jinja2) | New (Generic) |
|--------|-------------|---------------|
| Generated Go source | 554 MB / 13.4M lines | 564 lines (fixed) |
| Go compilation | **Fails** (file too large) | **Passes** in seconds |
| Provider binary | N/A | 294 MB |

---

## Part 6 — Run the Automated Test Suite

For a quick pass/fail check of everything above:

```bash
cd test-generic-provider

bash run_tests.sh
```

This runs 11 automated tests covering:

| # | Test |
|---|------|
| 1 | Go unit tests (105+ tests across 4 packages) |
| 2 | Generate filtered provider with `--generic -x` |
| 3 | `go build .` the filtered provider |
| 4 | Verify all generated files present, no test files leaked |
| 5 | Validate `trimmed_schema.json` structure |
| 6 | Run the provider binary `--help` |
| 7 | `terraform init` with dev_overrides |
| 8 | `terraform validate` — confirms schema loads |
| 9 | Generate full model provider (no -x, 33 YANG files) |
| 10 | `go build .` the full 268MB model |
| 11 | 85 Python tests (`pytest junosterraform/tests/`) |

Expected: `ALL TESTS PASSED`

---

## Summary of Commands

```bash
# Build
cd examples/providers && bash build-generic.sh

# Generate .tf files
bash convert.sh

# Set up Terraform
cat > ~/.terraformrc << EOF
provider_installation {
  dev_overrides {
    "hashicorp/junos-vqfx-evpn-vxlan" = "$(go env GOPATH)/bin"
    "hashicorp/junos-vsrx-evpn-vxlan" = "$(go env GOPATH)/bin"
  }
  direct {}
}
EOF

# Test
cd examples/terraform_files
terraform validate
terraform plan
terraform apply -auto-approve
terraform destroy -auto-approve

# Full model (covers all Junos)
cd examples/providers && bash build-generic-full.sh

# Automated tests
cd test-generic-provider && bash run_tests.sh
```
