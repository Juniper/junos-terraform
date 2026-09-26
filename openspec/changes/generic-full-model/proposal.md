## Why

Issue #150 asks for a provider that interprets the pyang JSON schema at runtime and, as the next step, works with the full Junos model instead of a schema trimmed to the user's configuration. The generic provider from `go-generator-replacement` (#151) builds a Terraform schema from the JSON, but was not usable:

- A full model did not load: decimal64 range bounds are strings in the JSON, a few names have capitals Terraform rejects, and YANG `choice`/`case` nodes (5920 choices holding some 26000 nodes in the SRX 26.2 model) were dropped. `Resources()` swallowed the load error, so Terraform only reported an unsupported resource type.
- The resource did not talk to the device: `Read` kept the prior state, `Create`/`Update` copied the plan into state, `Delete` did nothing.
- The plan↔XML bridges worked on flat paths without list keys, collapsing every list to one entry.
- On terraform-plugin-framework, a plan with a full model against an SRX300 took about 5 minutes (650s of CPU), all in the framework's collection-block walks.

## What Changes

- Load a full model: numeric-string range bounds, lowercase attribute names with a name check, `choice`/`case` flattened, load errors reported as diagnostics.
- `jtaf-provider --generic --exclude PATH` (repeatable) leaves schema subtrees out, e.g. `groups`, `logical-systems`, `version`; the schema is embedded compact and gzipped.
- Keyed converters between the resource's value and configuration XML (`ValueToConfig`, `ConfigToValue`), replacing the flat-path bridges.
- Create/Read/Update/Delete against the device, as the generated provider does them.
- Serve the plugin protocol (tfprotov6) directly with terraform-plugin-go instead of terraform-plugin-framework. Configurations and state are unchanged: the provider settings, resource type and attribute shape stay the same. The same plan takes 30 seconds with 6s of CPU, the same as the generated provider with a trimmed schema.
- `jtaf-provider --generic` writes the main package (`main.go`, `embed_schema.go`).
- Requires YANG `presence` in the schema (#156) for presence containers, which have children in a full model.

## Capabilities

### New Capabilities
- `embedded-schema-loading`: full models, choice flattening, error reporting.
- `dynamic-schema-builder`: a tfprotov6 schema and its value type.
- `generic-plan-xml-bridge`: keyed converters between values and configuration XML.

### Modified Capabilities
- `terraform-provider`: generic mode serves tfprotov6 and implements CRUD.
- `code-generation`: `--exclude`, gzipped schema, generated `main.go`.

## Impact

- **Go source**: `terraform_provider/generic/` rewritten (no framework dependency); `patch/` gains `FlattenChoices`, `BuildSchemaIndex`, `MarshalTree` and numeric-string range bounds; the framework `provider.go` no longer registers a generic resource.
- **Python CLI**: `jtaf-provider --generic --exclude`; generated main package.
- **Users**: state from the generated provider loads into the generic provider unchanged (attributes the schema does not have are ignored). Import is not supported.
