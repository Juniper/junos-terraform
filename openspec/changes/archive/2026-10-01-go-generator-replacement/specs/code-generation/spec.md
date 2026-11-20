## MODIFIED Requirements

### Requirement: jtaf-provider generates a usable Terraform provider directory
The `jtaf-provider` CLI tool SHALL produce a complete, buildable Terraform provider directory. Instead of rendering Go source code via Jinja2, it SHALL copy the generic provider source and emit `trimmed_schema.json` into the output directory.

#### Scenario: Provider generation without -x
- **WHEN** `jtaf-provider -j schema.json -t vqfx` is run without `-x`
- **THEN** the output directory `terraform-provider-junos-vqfx/` SHALL contain the generic Go source files, `go.mod`, and `trimmed_schema.json` (full schema)

#### Scenario: Provider generation with -x
- **WHEN** `jtaf-provider -j schema.json -x config.xml -t vqfx` is run with `-x`
- **THEN** the output directory SHALL contain the generic Go source files and a filtered `trimmed_schema.json`

#### Scenario: Output directory is buildable
- **WHEN** `go build .` is run in the output directory
- **THEN** the build SHALL succeed producing a provider binary

#### Scenario: Jinja2 template no longer required
- **WHEN** `jtaf-provider` runs
- **THEN** it SHALL NOT read or require `resource_config_provider.go.j2`
