## MODIFIED Requirements

### Requirement: Provider implements CRUD via NETCONF
The Terraform provider SHALL implement Create, Read, Update, and Delete operations using NETCONF RPCs, with all operations driven by the schema index rather than per-schema typed structs.

#### Scenario: Create sends full configuration
- **WHEN** Terraform calls Create with a plan
- **THEN** the provider SHALL convert the plan to XML via generic tree walk, send `load-configuration` via NETCONF, and commit

#### Scenario: Read returns current device state
- **WHEN** Terraform calls Read
- **THEN** the provider SHALL call `get-configuration` via NETCONF, convert the XML response to Terraform state via generic tree walk, and return the state

#### Scenario: Update computes minimal diff patch
- **WHEN** Terraform calls Update with a new plan
- **THEN** the provider SHALL convert both plan and current state to XML, use `ComputeDiff` to find changes, generate a NETCONF `edit-config` patch via `CreateDiffPatch`, apply it, and commit

#### Scenario: Delete removes managed configuration
- **WHEN** Terraform calls Delete
- **THEN** the provider SHALL compute a diff between current state and empty config, generate delete operations via `CreateDiffPatch`, apply via `edit-config`, and commit

### Requirement: Provider registers single config resource
The provider SHALL register exactly one resource type whose name is derived from the device type (e.g., `terraform-provider-junos-vqfx`).

#### Scenario: Resource metadata
- **WHEN** Terraform queries resource metadata
- **THEN** the resource type name SHALL be `terraform-provider-junos-<device-type>`

### Requirement: Provider schema built from embedded JSON
The provider SHALL build its resource schema dynamically from the embedded `trimmed_schema.json` at startup using the dynamic schema builder.

#### Scenario: Schema available on first RPC
- **WHEN** Terraform calls `GetProviderSchema`
- **THEN** the provider SHALL return the complete nested attribute schema built from the JSON
