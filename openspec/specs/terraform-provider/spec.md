# Terraform Provider Specification

Go implementation of a Terraform Plugin Framework provider for Junos configuration management via NETCONF. Lives at `terraform_provider/`.

## Architecture

```
Terraform CLI
    ↓ (Plugin Protocol)
Provider (main.go → provider.go)
    ↓
Config Resource (resource_config_provider.go)
    ↓ (uses patch engine for Update)
NETCONF Client (netconf/client.go)
    ↓ (SSH:830)
Junos Device
```

---

## Provider Configuration

**Location:** `terraform_provider/provider.go`

### Schema

```hcl
provider "junos-{device-type}" {
  host     = "10.54.21.174"     # Required: Device IP or hostname
  port     = 830                # Optional: NETCONF port (default 22)
  username = "root"             # Required: SSH username
  password = "secret"           # Required: SSH password (or use sshkey)
  sshkey   = "~/.ssh/id_rsa"   # Optional: SSH private key path
}
```

### Behaviors

#### Provider Initialization

- **Given** valid host, username, and password, **When** `Configure()` runs, **Then** a NETCONF client is created and stored for resource use
- **Given** client creation fails (SSH unreachable, auth failure), **When** `Configure()` runs, **Then** `resp.Diagnostics.AddError("failed to create client", err.Error())` is called and provider is unusable
- **Given** config parsing fails (missing required fields), **When** `Configure()` runs, **Then** diagnostics error appended and function returns early

#### Provider Metadata

- **Given** the provider is registered, **When** `Metadata()` is called, **Then** return `TypeName = "terraform_provider"` (overridden in generated code to `junos-{type}`)

#### Provider Resources

- **Given** the provider is initialized, **When** `Resources()` is called, **Then** return the `configResource` type (plus `fileResource` in dev mode)

#### Client Factory (Testability)

- **Given** `providerClientFactory` is set to default, **When** `Configure()` runs, **Then** `cfg.Client()` establishes real NETCONF connection
- **Given** `providerClientFactory` is overridden (tests), **When** `Configure()` runs, **Then** mock client is used instead

---

## NETCONF Client

**Location:** `terraform_provider/netconf/`

### Behaviors

#### Connection

- **Given** valid SSH credentials and reachable host, **When** `Client()` is called, **Then** establish SSH session and NETCONF subsystem
- **Given** SSH key is provided instead of password, **When** connecting, **Then** use key-based authentication
- **Given** host is unreachable or port closed, **When** connecting, **Then** return error with connection details

#### RPC Operations

| Given RPC | When Sent | Then Device Behavior |
|---|---|---|
| `<load-configuration>` with XML | Config needs staging | XML loaded into candidate datastore |
| `<edit-config>` with `nc:operation` attributes | Patch needs applying | Operations applied to candidate |
| `<commit/>` | Candidate is ready | Candidate promoted to running config |
| `<discard-changes/>` | Rollback needed | Candidate reset to match running |
| `<get-configuration>` | State refresh needed | Running config returned as XML |

#### Error Handling

- **Given** NETCONF RPC returns `<rpc-error>`, **When** response parsed, **Then** error propagated to caller with severity and message
- **Given** SSH session drops mid-RPC, **When** response waited, **Then** timeout error returned

---

## Config Resource — CRUD Lifecycle

**Location:** `terraform_provider/resource_config_provider.go` (Jinja2-generated)

### Behaviors

#### Create

- **Given** Terraform plan has new resource, **When** `Create()` executes, **Then**:
  1. Build XML from Terraform state attributes
  2. Send `<load-configuration>` RPC with the XML
  3. Send `<commit/>` RPC
  4. If commit succeeds, store state
  5. If commit fails, send `<discard-changes/>` and return error

#### Read

- **Given** resource exists in state, **When** `Read()` executes, **Then**:
  1. Send `<get-configuration>` for the resource's group
  2. Parse returned XML into Terraform state attributes
  3. If config not found on device, remove resource from state (triggers re-create)

#### Update

- **Given** Terraform plan differs from state, **When** `Update()` executes, **Then**:
  1. Build leaf maps for current state and desired plan (via `LeafMapWithSchema`)
  2. Compute diff (via `ComputeDiff`)
  3. If diff is empty, no-op
  4. Build patch XML (via `CreateDiffPatchWithSchema`)
  5. Align XML ordering to reference (via `AlignXMLOrderToReference`)
  6. Send `<edit-config>` RPC with patch
  7. Send `<commit/>`
  8. If commit fails, `<discard-changes/>` and return error

#### Delete

- **Given** resource marked for destruction, **When** `Delete()` executes, **Then**:
  1. Build delete XML with `nc:operation="delete"` on the resource group
  2. Send `<edit-config>` RPC
  3. Send `<commit/>`
  4. Clear resource from state

---

## File Resource (Development/Testing)

**Location:** `terraform_provider/resource_file.go`

### Behaviors

- **Given** `name` and `contents` provided, **When** `Create()` runs, **Then** write `contents` to `{managed_dir}/{name}` with mode 0644
- **Given** file exists on disk, **When** `Read()` runs, **Then** load file contents into state
- **Given** file not found, **When** `Read()` runs, **Then** add diagnostics error and clear state
- **Given** `name` changes, **When** plan computed, **Then** force resource replacement (new file)
- **Given** resource destroyed, **When** `Delete()` runs, **Then** remove file; if already missing, no error

---

## Testing

### How to Run

```bash
# All provider tests
cd terraform_provider && go test ./... -v

# Provider-only tests (no patch engine)
cd terraform_provider && go test . -v

# NETCONF client tests
cd terraform_provider && go test ./netconf/ -v

# Specific test
cd terraform_provider && go test . -run TestProvider_Configure -v

# With race detection
cd terraform_provider && go test ./... -race

# Coverage
cd terraform_provider && go test ./... -coverprofile=coverage.out && go tool cover -func=coverage.out
```

### Test Coverage Map

| Test File | What It Validates |
|-----------|-------------------|
| `provider_test.go` | Provider schema fields, metadata, type name |
| `provider_build_test.go` | Provider assembly, resource registration |
| `provider_runtime_test.go` | Runtime lifecycle, client injection, mock behaviors |
| `config_test.go` | Config struct creation, Client() factory |
| `resource_config_provider_test.go` | Config resource CRUD operations |
| `resource_file_test.go` | File resource schema, metadata |
| `resource_file_crud_test.go` | File create/read/update/delete operations |
| `resource_file_helpers_test.go` | File I/O helper functions, edge cases |
| `netconf/client_test.go` | NETCONF client connection, RPC sending |

### Key Test Behaviors

| Test Name | Given | When | Then |
|-----------|-------|------|------|
| `TestProvider_Configure` | Valid config model | Configure called | Client created without error |
| `TestProvider_Configure_MissingHost` | Empty host field | Configure called | Diagnostics contains error |
| `TestConfigResource_Create` | New resource plan | Create called | load-configuration + commit sent |
| `TestConfigResource_Update` | Plan differs from state | Update called | edit-config with patch + commit sent |
| `TestConfigResource_Delete` | Resource destruction | Delete called | edit-config delete + commit sent |
| `TestFileResource_ForceReplace` | Name field changed | Plan computed | Resource replaced (not updated in-place) |
