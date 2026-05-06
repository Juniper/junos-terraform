# Full-Stack Terraform Provider Test Report

**Date:** April 30, 2026  
**Status:** ALL TESTS PASSING ✅ — Full Terraform lifecycle verified on live device

---

## Test Environment

| Component | Details |
|-----------|---------|
| **Device** | vpaa — 10.54.21.174:830 (MX960, Junos 26.2I-20260201) |
| **Provider** | `junos-vmx-4-topo` (generated from device config) |
| **YANG Models** | 25.4R1 (`yang/25.4/25.4R1/native/conf-and-rpcs/`) |
| **Input XML** | `examples/patch_test/vpaa_config.xml` (live device config, groups removed) |
| **Terraform** | v1.9.4 (darwin_arm64) |
| **Go** | 1.25.6 |
| **Branch** | `netconf_patch` |

---

## Provider Generation

The provider was generated from the **actual device running config** to ensure the Terraform schema matches what exists on the device.

### Steps
```bash
# 1. Retrieve config from device
python3 -c "from ncclient import manager; ..." > examples/patch_test/vpaa_config.xml

# 2. Remove <groups> and <apply-groups> (not in YANG schema)
# Removed lines 4-267 from vpaa_config.xml

# 3. Generate Go structs from YANG
jtaf-yang2go -config examples/patch_test/vpaa_config.xml \
  -yang yang/25.4/25.4R1/native/conf-and-rpcs/ \
  -providerName vmx-4-topo

# 4. Generate Terraform resource file
jtaf-xml2tf -config examples/patch_test/vpaa_config.xml \
  -providerName vmx-4-topo

# 5. Build provider
jtaf-provider -j trimmed_schema.json
cd examples/providers/terraform-provider-junos-vmx-4-topo
go build -o terraform-provider-junos-vmx-4-topo
```

### Build Result
- Binary: `terraform-provider-junos-vmx-4-topo` (29 MB)
- Compile: **Clean** — no errors or warnings

---

## Unit Test Results

### Patch Engine Tests (`terraform_provider/patch/`)
```
$ cd terraform_provider && go test ./patch/ -v
--- PASS: 44 tests
PASS
ok      terraform-provider-junos-device/patch    0.008s
```

| Category | Tests | Status |
|----------|-------|--------|
| Keyed list operations | 12 | ✅ |
| Scalar leaf operations | 8 | ✅ |
| Compound key (community-name) | 2 | ✅ |
| Leaf-list operations | 5 | ✅ |
| Schema-aware leaf map | 6 | ✅ |
| Edge cases (empty diff, no-op) | 11 | ✅ |
| **Total** | **44** | **ALL PASS** |

### Key Unit Tests Added
| Test Name | Validates |
|-----------|-----------|
| `TestCreateDiffPatch_CompoundKeyDeleteIncludesChoiceSibling` | Community delete emits `<add/>` alongside `<community-name>` |
| `TestLeafMapWithSchema_CompoundKeyEmitsAllChildren` | Leaf map produces paths for all compound key parts |
| `TestCreateDiffPatch_KeyedListRenameWithDescendantsUsesEntryOperations` | Two-pass algorithm: children of deleted parent have no operations |

### Generated Provider Tests
```
$ cd examples/providers/terraform-provider-junos-vmx-4-topo && go test ./...
ok      terraform-provider-junos-vmx-4-topo           0.004s
ok      terraform-provider-junos-vmx-4-topo/netconf   0.003s
ok      terraform-provider-junos-vmx-4-topo/patch     0.007s
```
All packages pass.

---

## Live Device Test Results — Full Terraform Lifecycle

### Test Configuration
- **Resource:** `terraform-provider-junos-vmx-4-topo.vpaa_config-base-config`
- **Group name:** `base-config`
- **Workspace:** `tests/terraform_fullstack_vmx4/`
- **Config scope:** Full device config (system, interfaces, protocols, routing, policy, firewall, SNMP, etc.)

### CREATE ✅

```
$ TF_CLI_CONFIG_FILE=./terraformrc terraform apply -auto-approve

terraform-provider-junos-vmx-4-topo.vpaa_config-base-config: Creating...
terraform-provider-junos-vmx-4-topo.vpaa_config-base-config: Creation complete after 21s

Apply complete! Resources: 1 added, 0 changed, 0 destroyed.
```

**Verification (NETCONF via ncclient):**
```python
# Confirmed group "base-config" exists with host_name = "tf-fullstack-test"
>>> config.xpath('//groups[name="base-config"]/system/host-name')
['tf-fullstack-test']
```

### UPDATE ✅

Changed `host_name` from `"tf-fullstack-test"` → `"tf-update-test-v2"` in `main.tf`.

```
$ TF_CLI_CONFIG_FILE=./terraformrc terraform apply -auto-approve

terraform-provider-junos-vmx-4-topo.vpaa_config-base-config: Modifying...
terraform-provider-junos-vmx-4-topo.vpaa_config-base-config: Modifications complete after 37s

Apply complete! Resources: 0 added, 1 changed, 0 destroyed.
```

**Verification (NETCONF via ncclient):**
```python
# Confirmed host_name updated within the base-config group
>>> config.xpath('//groups[name="base-config"]/system/host-name')
['tf-update-test-v2']
```

**Patch payload characteristics:**
- Only the changed leaf (`host_name`) received `nc:operation="replace"`
- No spurious deletes generated — provider reads only managed attributes
- Community entries correctly include `<add/>` sibling (no choice-ident errors)

### DESTROY ✅

```
$ TF_CLI_CONFIG_FILE=./terraformrc terraform destroy -auto-approve

terraform-provider-junos-vmx-4-topo.vpaa_config-base-config: Destroying...
terraform-provider-junos-vmx-4-topo.vpaa_config-base-config: Destruction complete after 16s

Destroy complete! Resources: 1 destroyed.
```

**Verification (NETCONF via ncclient):**
```python
# Confirmed group "base-config" no longer exists on device
>>> config.xpath('//groups[name="base-config"]')
[]
# "base-config group successfully removed!"
```

### Lifecycle Summary

| Operation | Duration | Result | Verification |
|-----------|----------|--------|--------------|
| **CREATE** | 21s | ✅ PASS | Group created, hostname set |
| **UPDATE** | 37s | ✅ PASS | Hostname changed in-place |
| **DESTROY** | 16s | ✅ PASS | Group removed from device |

---

## Patch Engine Fixes Applied

### Fix 1: `default-operation=none` → `merge`

| | Before | After |
|-|--------|-------|
| **Behavior** | NETCONF rejected partial edits | Partial config merges work |
| **File** | `netconf/client.go` | `netconf/client.go` |
| **Impact** | CREATE worked, UPDATE/DELETE failed | All operations work |

### Fix 2: Scalar leaf delete text

| | Before | After |
|-|--------|-------|
| **Output** | `<leaf nc:operation="delete">old-value</leaf>` | `<leaf nc:operation="delete"/>` |
| **Impact** | Junos rejected delete with stale text | Clean deletes accepted |

### Fix 3: Compound key support (community-name)

| | Before | After |
|-|--------|-------|
| **Output** | `<community nc:operation="delete"><community-name>X</community-name></community>` | `<community nc:operation="delete"><community-name>X</community-name><add/></community>` |
| **Error** | `expecting <choice-ident>` | No errors |
| **Root cause** | Leaf map didn't emit choice-ident siblings for compound keys | `leafmap.go` handles space-separated keys; `patch.go` two-pass preserves sibling elements |

### Fix 4: Provider Read scope (resolved by regeneration)

| | Old provider (`vmx-template-test`) | New provider (`vmx-4-topo`) |
|-|-----------------------------------|-----------------------------|
| **Input XML** | Minimal template (few attributes) | Full device config (1473 lines) |
| **Read behavior** | Read entire device, diff against small plan → spurious deletes | Read matches plan scope — no spurious deletes |
| **UPDATE/DESTROY** | ⚠️ BLOCKED | ✅ PASS |

---

## Known Limitations

| Issue | Impact | Severity |
|-------|--------|----------|
| UTF-8 em-dash drift (`—` → double-encoded on read-back) | Cosmetic diff on description fields | Low |
| `extension_service/request_response/grpc` not in YANG | Shows as drift on read-back | Low |
| `apply-groups` ordering not managed | Groups applied in device order, not Terraform order | Low |
| `<groups>` container not in YANG schema | Must be removed from input XML before generation | N/A (build-time) |

---

## Key Files

| File | Purpose |
|------|---------|
| `terraform_provider/patch/patch.go` | Two-pass CreateDiffPatch, scalar delete fix |
| `terraform_provider/patch/leafmap.go` | Compound key support |
| `terraform_provider/patch/patch_test.go` | All unit tests (44 tests) |
| `terraform_provider/netconf/client.go` | default-operation=merge fix |
| `examples/providers/terraform-provider-junos-vmx-4-topo/` | Generated provider (built, tested) |
| `examples/patch_test/vpaa_config.xml` | Device config input (groups removed) |
| `examples/terraform_files/vpaa_config.tf` | Generated .tf file |
| `tests/terraform_fullstack_vmx4/` | Live test workspace |

---

## How to Reproduce

```bash
# 1. Build the provider
cd examples/providers/terraform-provider-junos-vmx-4-topo
go build -o terraform-provider-junos-vmx-4-topo

# 2. Run unit tests
cd terraform_provider && go test ./patch/ -v

# 3. Run generated provider tests
cd examples/providers/terraform-provider-junos-vmx-4-topo && go test ./...

# 4. Run Terraform full-stack test
cd tests/terraform_fullstack_vmx4
TF_CLI_CONFIG_FILE=./terraformrc terraform apply -auto-approve

# 5. Update a value in main.tf, then re-apply
TF_CLI_CONFIG_FILE=./terraformrc terraform apply -auto-approve

# 6. Destroy
TF_CLI_CONFIG_FILE=./terraformrc terraform destroy -auto-approve

# 7. Debug mode (writes patch XML to /tmp/)
TF_CLI_CONFIG_FILE=./terraformrc JUNOS_TF_DEBUG_PATCH=1 terraform apply -auto-approve
# Check: /tmp/debug_vmx4_patch.xml, /tmp/debug_diffmap.txt
```

## Credentials (for live tests)
- vpaa: 10.54.21.174:830, root/pass
- Python venv: `/Users/patelv/Desktop/merge/venv`
