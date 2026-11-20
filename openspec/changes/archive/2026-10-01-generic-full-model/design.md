## Context

The generic provider embeds the pyang JSON schema and converts between the Terraform value and Junos configuration XML by walking it. With a full model the schema has about a million nodes (SRX 26.2: 1,035,369 paths), and every object in the configuration has all of its possible attributes, nearly all null.

## Decisions

### 1. Flatten `choice` and `case` when loading

YANG choices and cases have no XML element; their nodes appear under the choice's parent. `jtaf_common` flattens them when trimming a schema to example XML; a full model is not trimmed. `patch.CompileSchema` flattens them, merging same-named siblings, and both the patch engine's lookups and the Terraform schema use the result (decision 6). Flattening creates no name collisions in the SRX 26.2 model; `jtaf-provider --generic` checks the names when it generates the provider, and fails if a model ever does, so the provider does not check them each time it starts.

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
| generic, terraform-plugin-go, nothing excluded | 29.7s¹ | 14.1s | 1.95GB |

¹ Median of three plans. An earlier single plan took 51.6s: the device's refresh, about 20s of every plan, varies by several seconds between otherwise identical runs.

Profiles put the framework time in `Reify`/`NullifyCollectionBlocks`: for every node of the value they resolve its type from the schema root (`SchemaTypeAtTerraformPath`), rebuilding nested object types, so the cost grows with the configuration's size times the schema's. The generic resource uses no framework feature (no plan modifiers, defaults, validators or semantic equality), so it serves tfprotov6 directly: `PlanResourceChange` returns the proposed state (nothing is computed; a `resource_name` change requires replacement), `UpgradeResourceState` decodes stored JSON ignoring undefined attributes, and unsupported RPCs return errors.

### 6. A compact schema table, embedded compiled

With a full model, the provider held the schema twice, as the patch engine's index (a map from each node's full path to a struct with its leaf constraints) and as the tree of JSON nodes: 984MB after loading the SRX 26.2 model, 943k nodes, and 1s and 2.6GB of allocation to load, in each of the three processes OpenTofu starts for a plan. Precompiling those structures (gob) did not help: decoding built the same structs.

`patch.Schema` keeps one 20-byte record per node (name, kind, list key, flags, children), laid out breadth first so a node's children are consecutive, and each distinct name (8.7k) once. The patch engine looks a path up by walking it from `<configuration>`; the generic provider walks the nodes. It has no pointers, so its binary form is its records as held in memory, and reading it is a copy: `jtaf-provider --generic` compiles it with `cmd/compileschema` and the provider embeds it gzipped (3.0MB, against 2.8MB of JSON). Against the SRX300, medians of three plans:

| Provider | plan | outside the refresh | provider CPU | provider peak RSS | `GetProviderSchema` |
|---|---|---|---|---|---|
| path index, excludes | 27.8s | 7.5s | 2.7s | 0.54GB | 0.22s |
| path index, nothing excluded | 29.7s | 10.7s | 8.3s | 1.95GB | 1.03s |
| compiled table, excludes | 25.5s | 7.0s | 1.8s | 0.21GB | 0.06s |
| compiled table, nothing excluded | 26.6s | 8.0s | 2.8s | 0.38GB | 0.14s |

What a full model still costs is OpenTofu's handling of the larger schema. The index built a node's path by joining its parent's and stripping a leading `configuration`, which collapsed nodes named `configuration` into their parent (`snmp trap-group categories` took the kind of its leaf `configuration`); the table has them where they are, and otherwise answers as the index did for every path of the SRX 26.2 and 25.4 models.

## Risks

- A full model makes every object carry all its attributes in state; state files grow accordingly.
- Import is not supported.
- Generating a generic provider needs Go, to compile the schema.
