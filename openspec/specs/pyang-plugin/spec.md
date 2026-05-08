# pyang Plugin

Custom pyang output plugin that converts YANG modules into a JSON schema tree used by JTAF code generation. Lives at `jtaf_pyang_plugin/`.

## Architecture

```
YANG modules (.yang files)
    ↓ pyang parser
YANG AST (pyang internal)
    ↓ jtaf_json.py plugin
JSON schema tree (stdout)
    ↓ piped to jtaf-provider
Filtered schema for code generation
```

## Files

| File | Purpose |
|------|---------|
| `jtaf_json.py` | pyang plugin — walks YANG AST, outputs JSON tree |
| `jtaf-pyang-plugindir` | Shell script that prints the plugin directory path |
| `tests/test_jtaf_json.py` | Integration tests for JSON output |

## Key Classes

| Class | Purpose |
|-------|---------|
| `FNode` | Represents a single YANG tree node (container/leaf/list/leaf-list) |
| `FNodeTree` | Tree structure with root `FNode` and identities list |

## YANG Type Mapping (`yang_coerce`)

The plugin maps YANG types to JSON-friendly representations:

| YANG Type | Handling |
|-----------|----------|
| `string`, `inet:*` | String value |
| `uint8/16/32/64`, `int8/16/32/64` | Numeric |
| `boolean` | Boolean |
| `empty` | Presence semantics (tag only, no value) |
| `identityref` | Namespace-qualified identity |
| `leafref` | Resolved to target type |
| `union` | First matching member type |
| `enumeration` | Enum values list |

## Usage

```bash
# Get plugin directory
jtaf-pyang-plugindir

# Run pyang with JTAF plugin
pyang --plugindir $(jtaf-pyang-plugindir) -f jtaf \
  -p <common-yang-dir> <yang-files> > schema.json

# Pipe directly to provider generation
pyang --plugindir $(jtaf-pyang-plugindir) -f jtaf \
  -p examples/yang/18.2/18.2R3/common \
  examples/yang/18.2/18.2R3/junos-qfx/conf/*.yang | \
  jtaf-provider -j - -x configs/*.xml -t vqfx
```

## JSON Output Format

The plugin produces a tree where each node contains:
- `name` — Element name
- `type` — YANG type or "container"/"list"
- `key` — List key field(s) (for YANG lists)
- `children` — Nested nodes
- `description` — YANG description text
- `ordered-by-user` — Boolean for ordered leaf-lists

## Supported YANG Versions

The plugin works with any version in the `yang/` repo (14.2 through 25.4). Common and device-specific modules are supported for: junos-qfx, junos-vmx, junos-ex, junos-srx, and others.

## Tests

- `jtaf_pyang_plugin/tests/test_jtaf_json.py` — Container handling, enum extraction, leaf-list output, identity-ref resolution
