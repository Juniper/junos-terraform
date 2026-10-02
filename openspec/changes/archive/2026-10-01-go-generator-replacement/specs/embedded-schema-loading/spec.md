## ADDED Requirements

### Requirement: Embed schema JSON in provider binary
The provider binary SHALL embed `trimmed_schema.json` (optionally gzip-compressed) using Go's `//go:embed` directive so the binary is self-contained with no external file dependencies.

#### Scenario: Compressed embedding
- **WHEN** the provider is built with a `trimmed_schema.json.gz` file present
- **THEN** the binary SHALL contain the compressed schema and decompress it at startup

#### Scenario: Raw embedding fallback
- **WHEN** the provider is built with only `trimmed_schema.json` (no .gz)
- **THEN** the binary SHALL embed the raw JSON directly

### Requirement: Parse schema into index at startup
The provider SHALL parse the embedded JSON into the existing `map[path]NodeInfo` index (from `patch/process_schema.go`) during provider initialization, before any Terraform RPC is served.

#### Scenario: Successful startup
- **WHEN** the provider process starts
- **THEN** it SHALL parse the schema and have the index ready before responding to `GetProviderSchema` RPC

#### Scenario: Corrupt or missing schema
- **WHEN** the embedded schema JSON is malformed or empty
- **THEN** the provider SHALL fail with a clear error message during startup (not silently serve empty schema)

### Requirement: Schema index shared across operations
The parsed schema index SHALL be built once and shared across all CRUD operations without re-parsing.

#### Scenario: Multiple terraform apply calls
- **WHEN** Terraform calls Create, then Read, then Update on the same provider process
- **THEN** all operations SHALL use the same pre-built schema index without re-parsing JSON
