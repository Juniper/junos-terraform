# Testing Framework

Test infrastructure across Go and Python for unit, integration, and end-to-end validation.

## Test Suites

### Go Unit Tests — `terraform_provider/`

```bash
cd terraform_provider && go test ./...
```

| Area | Test Files | Coverage |
|------|-----------|----------|
| **Provider** | `provider_test.go`, `provider_build_test.go`, `provider_runtime_test.go` | Schema, metadata, assembly, lifecycle |
| **Config** | `config_test.go` | Config creation, NETCONF client init |
| **Config Resource** | `resource_config_provider_test.go` | CRUD operations |
| **File Resource** | `resource_file_test.go`, `resource_file_crud_test.go`, `resource_file_helpers_test.go` | Schema, CRUD, file I/O helpers |
| **NETCONF Client** | `netconf/client_test.go` | Client behavior |

### Go Patch Engine Tests — `terraform_provider/patch/`

| Test File | Focus |
|-----------|-------|
| `patch_test.go` | Diff patch generation, NETCONF operation attributes |
| `matrix_test.go` | Comprehensive CRUD matrix (all create/replace/delete combos) |
| `corner_case_test.go` | Special chars in keys, container delete coalescing |
| `schema_phase1_test.go` | JSON schema parsing, trimmed schema validation |
| `order_test.go` | Element ordering preservation (ordered-by-user) |

### Python Unit Tests

```bash
pytest junosterraform/unit_tests/
pytest jtaf_pyang_plugin/tests/
pytest netconf_mock/tests/
```

| Test File | Coverage |
|-----------|----------|
| `junosterraform/unit_tests/test_jtaf_common.py` | XPath extraction, nesting, schema filtering |
| `junosterraform/unit_tests/test_workflow.py` | End-to-end jtaf-yang2go workflow |
| `junosterraform/unit_tests/test_xml2tf_flatten.py` | XML→HCL conversion, apply-groups flattening |
| `jtaf_pyang_plugin/tests/test_jtaf_json.py` | pyang plugin JSON output validation |
| `netconf_mock/tests/test_netconf_mock_server.py` | Mock server RPC handling, multi-device state |

### Fullstack / E2E Tests — `tests/`

| Directory | Purpose |
|-----------|---------|
| `tests/terraform_fullstack_vmx4/` | Terraform plan/apply/destroy against live or mock devices |

E2E tests use the generated provider against a target device (vpaa at 10.54.21.174:830) or the NETCONF mock server.

## Test Patterns

### Go Tests
- Use `testing` package with `testify` assertions
- Table-driven tests for matrix coverage
- XML string literals for input/expected comparisons
- `ProcessSchema()` with inline JSON for schema-aware tests

### Python Tests
- Use `pytest` with fixtures
- XML string comparisons via `lxml`
- Subprocess calls for CLI tool integration tests
- `conftest.py` for shared fixtures

## Running All Tests

```bash
# Go tests (provider + patch engine)
cd terraform_provider && go test ./... -v

# Python tests (all suites)
pytest junosterraform/unit_tests/ jtaf_pyang_plugin/tests/ netconf_mock/tests/ -v

# Single patch engine test
cd terraform_provider && go test ./patch/ -run TestCornerCase -v
```

## Test Data

| Location | Content |
|----------|---------|
| `examples/patch_test/` | XML configs for patch engine testing |
| `examples/yang/18.2/` | YANG models used in workflow tests |
| `examples/evpn-vxlan-dc/` | Multi-device XML configs for provider tests |
| `examples/terraform_files/` | Generated .tf files for validation |
