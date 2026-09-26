## Context

The generic provider embeds the pyang JSON schema and converts between the Terraform value and Junos configuration XML by walking it. With a full model the schema has about a million nodes (SRX 26.2: 1,035,369 paths), and every object in the configuration has all of its possible attributes, nearly all null.

## Decisions

### 1. Flatten `choice` and `case` when loading

YANG choices and cases have no XML element; their nodes appear under the choice's parent. `jtaf_common` flattens them when trimming a schema to example XML; a full model is not trimmed. `patch.FlattenChoices` flattens them once, and both the index (`patch.BuildSchemaIndex`) and the Terraform schema are built from the flattened nodes, from one JSON parse. Flattening creates no name collisions in the SRX 26.2 model; `ValidateNames` fails the load if a model ever does.

### 2. Exclusions are command-line options

`--exclude PATH` removes a subtree (relative to `configuration`) before the schema is embedded; a path that does not exist is an error. Nothing is excluded by default, except the top-level `version`: it records the release the configuration was last committed with and cannot be set, so managing it would plan its deletion on every run. Trimming already drops it. Upstream's builder still leaves out `groups`/`apply-groups` attributes; excluding `groups` also keeps its ~500k nodes out of the embedded index.

### 3. Keyed converters

`ValueToConfig` emits each non-null attribute as an element in schema order, a list entry's keys first in key order (Junos requires it), an empty string as an empty element (YANG `empty`). `ConfigToValue` reads only schema elements, a missing element as null, repeated elements as list entries in document order. This matches the generated provider's typed structs, so state carries over.

### 4. CRUD

As the generated provider: Create loads the plan (merge) and commits; Read converts the device configuration aligned to the prior state's order; Update patches the diff between device and plan, commits, verifies, and falls back to loading the plan; Delete patches out the state. Update takes the new state from the configuration it read to verify (or, with no difference, the one it read first) instead of reading it again: each read is the whole configuration, several seconds on a small device. Unlike the generated provider, Read does not keep a planned top-level section the device lacks: that worked around structs that could not tell an absent section from an empty one, and hid the section being removed.

### 5. terraform-plugin-go instead of terraform-plugin-framework

Measured against an SRX300 (26.2R1.7) with a real configuration, full model minus groups, logical-systems, tenants, dynamic-profiles and version:

| Provider | plan | CPU | peak RSS |
|---|---|---|---|
| generated, trimmed schema | 29.8s | 2.2s | 0.08GB |
| generic, framework v1.18 | 295s | 653s | 1.2GB |
| generic, framework v1.19 | 298.5s | 657s | 1.2GB |
| generic, terraform-plugin-go | 30.1s | 5.8s | 0.5GB |
| generic, terraform-plugin-go, nothing excluded | 51.6s | 15.9s | 1.9GB |

Profiles put the framework time in `Reify`/`NullifyCollectionBlocks`: for every node of the value they resolve its type from the schema root (`SchemaTypeAtTerraformPath`), rebuilding nested object types, so the cost grows with the configuration's size times the schema's. The generic resource uses no framework feature (no plan modifiers, defaults, validators or semantic equality), so it serves tfprotov6 directly: `PlanResourceChange` returns the proposed state (nothing is computed; a `resource_name` change requires replacement), `UpgradeResourceState` decodes stored JSON ignoring undefined attributes, and unsupported RPCs return errors.

## Risks

- A full model makes every object carry all its attributes in state; state files grow accordingly.
- Import is not supported.
