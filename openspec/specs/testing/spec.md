# Testing Specification

Comprehensive testing guide for JTAF — all test suites, how to run them, coverage strategies, and patterns.

---

## Test Suite Overview

| Suite | Language | Location | Scope |
|-------|----------|----------|-------|
| Patch Engine | Go | `terraform_provider/patch/*_test.go` | Unit: leaf maps, diffs, XML ordering |
| Provider Core | Go | `terraform_provider/*_test.go` | Unit: provider config, resources, CRUD |
| NETCONF Client | Go | `terraform_provider/netconf/*_test.go` | Unit: SSH transport, RPC framing |
| pyang Plugin | Python | `jtaf_pyang_plugin/tests/test_jtaf_json.py` | Unit: YANG → JSON conversion |
| Mock Server | Python | `netconf_mock/tests/test_netconf_mock_server.py` | Unit: mock state machine, RPC handling |
| Live Device | Python | `tests/live_*.py` | Integration: real device connectivity/config |
| Fullstack E2E | Terraform | `tests/terraform_fullstack_test/` | E2E: Terraform plan/apply against mock |

---

## How to Run Each Suite

### Go — Patch Engine Tests

```bash
cd terraform_provider && go test ./patch/ -v
```

**Key test files:**
- `patch/patch_test.go` — `CreateDiffPatchWithSchema`, XML generation
- `patch/matrix_test.go` — Table-driven diff scenarios (additions, deletions, modifications)
- `patch/corner_case_test.go` — Empty types, ordered leaf-lists, identityref, multi-key lists
- `patch/order_test.go` — `AlignXMLOrderToReference` reordering
- `patch/schema_phase1_test.go` — `ProcessSchema` JSON loading and type indexing

**Run specific test:**
```bash
cd terraform_provider && go test ./patch/ -run TestComputeDiff_Additions -v
```

**Pattern: table-driven tests**
```go
tests := []struct {
    name     string
    current  map[string]string
    desired  map[string]string
    expected []DiffEntry
}{...}
for _, tt := range tests {
    t.Run(tt.name, func(t *testing.T) {
        got := ComputeDiff(tt.current, tt.desired)
        // assertions
    })
}
```

---

### Go — Provider Core Tests

```bash
cd terraform_provider && go test . -v
```

**Key test files:**
- `provider_test.go` — Schema validation, metadata, type name
- `provider_build_test.go` — Provider assembly, resource registration
- `provider_runtime_test.go` — Runtime lifecycle with mock client injection
- `config_test.go` — Config struct, Client() factory
- `resource_config_provider_test.go` — Config resource CRUD
- `resource_file_test.go` — File resource schema
- `resource_file_crud_test.go` — File resource CRUD operations
- `resource_file_helpers_test.go` — File I/O helpers
- `main_test.go` — Provider entry point smoke test

**Run with race detection:**
```bash
cd terraform_provider && go test . -race -v
```

**Pattern: mock client injection**
```go
func TestConfigResource_Create(t *testing.T) {
    mockClient := &MockNetconfClient{...}
    // Override providerClientFactory
    oldFactory := providerClientFactory
    providerClientFactory = func(cfg Config) (NetconfClient, error) {
        return mockClient, nil
    }
    defer func() { providerClientFactory = oldFactory }()
    // Test CRUD operation
}
```

---

### Go — NETCONF Client Tests

```bash
cd terraform_provider && go test ./netconf/ -v
```

**Key test files:**
- `netconf/transport_test.go` — SSH connection, session setup
- `netconf/rpc_test.go` — RPC send/receive, message framing
- `netconf/low_level_ssh_test.go` — Low-level SSH operations

---

### Python — pyang Plugin Tests

```bash
pytest jtaf_pyang_plugin/tests/ -v
```

**Key test file:** `test_jtaf_json.py`
- Container/list/leaf/leaf-list recognition
- Type resolution (empty, identityref, union, enumeration)
- Choice/case flattening
- Ordered-by-user flag preservation

---

### Python — Mock Server Tests

```bash
pytest netconf_mock/tests/ -v
```

**Key test file:** `test_netconf_mock_server.py`
- DeviceState lifecycle (commit, discard)
- Load-configuration with merge/replace
- Edit-config patch mode (recursive merge, delete)
- Group delete handling
- Multi-device isolation

---

### Python — Live Device Tests

```bash
# Connectivity check (no config changes)
python tests/live_connectivity_check.py

# Config commit test (modifies device)
python tests/live_commit_test.py

# Matrix test (multiple configs across scenarios)
python tests/live_matrix_test.py
```

**Prerequisites:**
- Reachable Junos device
- Credentials configured (typically env vars or hardcoded for lab)
- NETCONF enabled on device port 830

---

### E2E — Terraform Fullstack

```bash
cd tests/terraform_fullstack_test

# Point to local provider binary
export TF_CLI_CONFIG_FILE=./terraformrc

# Run against mock server
terraform init
terraform plan
terraform apply -auto-approve
terraform destroy -auto-approve
```

**Prerequisites:**
- Mock server running (`python netconf_mock/netconf_mock_server.py --device test:8301`)
- Provider binary built and installed
- `terraformrc` pointing to local provider path

---

## Running All Tests

```bash
# All Go tests with coverage
cd terraform_provider && go test ./... -coverprofile=coverage.out -v

# All Python tests
pytest jtaf_pyang_plugin/tests/ netconf_mock/tests/ -v

# Combined coverage report (Go)
cd terraform_provider && go tool cover -func=coverage.out

# Combined coverage report (Python)
pytest jtaf_pyang_plugin/tests/ netconf_mock/tests/ --cov --cov-report=term-missing
```

---

## Test Data Locations

| Data Type | Location | Used By |
|-----------|----------|---------|
| YANG models | `examples/yang/` | pyang plugin tests, code generation |
| XML configs | `examples/evpn-vxlan-dc/` | Provider generation, E2E |
| Terraform configs | `tests/terraform_fullstack_test/` | E2E apply/destroy |
| Patch test XML | `examples/patch_test/` | Patch engine unit tests |
| Schema JSON | `terraform_provider/trimmed_schema.json` | `ProcessSchema` tests |

---

## Test Patterns and Conventions

### Go Conventions

- **Table-driven tests**: All matrix/scenario tests use `[]struct{...}` with `t.Run()`
- **Parallel**: Independent subtests use `t.Parallel()`
- **Helpers**: Use `t.Helper()` for assertion functions
- **Temp files**: Use `t.TempDir()` for file resource tests
- **Naming**: `TestFunctionName_Scenario` (e.g., `TestComputeDiff_EmptyDesired`)

### Python Conventions

- **Fixtures**: Use `@pytest.fixture` for shared setup (DeviceState, mock sessions)
- **Parametrize**: Use `@pytest.mark.parametrize` for data-driven tests
- **Naming**: `test_behavior_description` (e.g., `test_commit_promotes_candidate`)
- **Assertions**: Plain `assert` statements (no unittest.TestCase)

---

## Coverage Targets

| Suite | Target | Current Focus |
|-------|--------|---------------|
| Patch Engine | >90% line coverage | Corner cases: empty, identityref, ordered leaf-lists |
| Provider Core | >80% | CRUD operations, error paths |
| NETCONF Client | >70% | Connection handling, RPC errors |
| Mock Server | >90% | All RPC types, state transitions |
| pyang Plugin | >80% | Type resolution, complex YANG constructs |

### Generating Coverage Reports

```bash
# Go HTML coverage report
cd terraform_provider && go test ./... -coverprofile=coverage.out && go tool cover -html=coverage.out -o coverage.html

# Python HTML coverage report
pytest netconf_mock/tests/ jtaf_pyang_plugin/tests/ --cov --cov-report=html

# Check a specific package coverage
cd terraform_provider && go test ./patch/ -coverprofile=patch_coverage.out && go tool cover -func=patch_coverage.out
```

---

## Adding New Tests

### When to Add Tests

- **Given** a new YANG type is supported, **When** adding it to the patch engine, **Then** add a matrix test case in `matrix_test.go`
- **Given** a new RPC type is handled, **When** adding it to mock server, **Then** add test in `test_netconf_mock_server.py`
- **Given** a bug is fixed, **When** writing the fix, **Then** add a regression test proving the fix
- **Given** a new CLI option is added, **When** implementing, **Then** add test for the option behavior

### Test File Naming

| Language | Convention | Example |
|----------|-----------|---------|
| Go | `{source}_test.go` | `patch.go` → `patch_test.go` |
| Go (scenarios) | `{topic}_test.go` | `corner_case_test.go` |
| Python | `test_{module}.py` | `test_netconf_mock_server.py` |

### Test Isolation

- **Given** Go tests, **When** they run in parallel, **Then** no shared mutable state between test functions
- **Given** Python mock tests, **When** testing state transitions, **Then** create fresh `DeviceState` per test
- **Given** file resource tests, **When** needing filesystem, **Then** use temp directories cleaned up automatically
