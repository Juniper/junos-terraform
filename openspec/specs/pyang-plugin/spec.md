# pyang Plugin Specification

Custom pyang output plugin that converts YANG modules into a JSON schema tree used by JTAF code generation. Lives at `jtaf_pyang_plugin/`.

## Architecture

```
YANG modules (.yang files)
    ↓ pyang parser (AST)
YANG abstract syntax tree
    ↓ jtaf_json.py plugin (emit phase)
JSON schema tree (stdout)
    ↓ piped to jtaf-provider
Filtered schema for code generation
```

---

## Files

| File | Purpose |
|------|---------|
| `jtaf_json.py` | pyang plugin — walks YANG AST, outputs JSON |
| `jtaf-pyang-plugindir` | Shell script that prints the plugin directory path |
| `tests/test_jtaf_json.py` | Integration tests for JSON output correctness |

---

## Behaviors

### Plugin Registration

- **Given** pyang is invoked with `--plugindir $(jtaf-pyang-plugindir) -f jtaf`, **When** pyang loads plugins, **Then** `jtaf_json.py` registers as output format `"jtaf"`
- **Given** `jtaf-pyang-plugindir` is called, **When** executed, **Then** print the absolute directory path containing `jtaf_json.py`

### YANG AST Walking

- **Given** YANG modules are parsed into AST, **When** the `emit()` function runs, **Then** the AST is walked depth-first and converted to `FNode` tree
- **Given** a YANG `container` statement, **When** processed, **Then** create `FNode` with `type="container"` and recurse into children
- **Given** a YANG `list` statement, **When** processed, **Then** create `FNode` with `type="list"` and capture `key` field
- **Given** a YANG `leaf` statement, **When** processed, **Then** create `FNode` with `type="leaf"` and resolve leaf type
- **Given** a YANG `leaf-list` statement, **When** processed, **Then** create `FNode` with `type="leaf-list"`

### Type Resolution

| Given YANG Type | When Resolved | Then JSON Representation |
|---|---|---|
| `string` | Type mapped | `"type": "string"` |
| `inet:ipv4-address` | Type mapped | `"type": "string"` (inet types → string) |
| `uint8`, `uint16`, `uint32`, `uint64` | Type mapped | `"type": "integer"` |
| `int8`, `int16`, `int32`, `int64` | Type mapped | `"type": "integer"` |
| `boolean` | Type mapped | `"type": "boolean"` |
| `empty` | Type mapped | `"type": "empty"` (presence semantics, no value) |
| `identityref` | Type mapped | `"type": "identityref"`, identity base recorded |
| `leafref` | Type mapped | Resolved to target leaf's type |
| `union` | Type mapped | `"type": "union"`, member types listed |
| `enumeration` | Type mapped | `"type": "enumeration"`, enum values listed |

### Special YANG Constructs

#### Choice/Case

- **Given** a YANG `choice` statement with `case` children, **When** processed, **Then** choice children are flattened into the parent container (no choice/case wrappers in JSON)
- **Given** overlapping case branches, **When** processed, **Then** all branches are included (runtime config determines which is active)

#### Ordered-by-User

- **Given** a `leaf-list` with `ordered-by user` statement, **When** processed, **Then** JSON node gets `"ordered-by": "user"` flag
- This flag is critical for the patch engine to track position-sensitive ordering

#### Augmentations

- **Given** a YANG `augment` statement targets a path, **When** the target module is loaded, **Then** augmented nodes appear as regular children of the target

#### Groupings and Uses

- **Given** a YANG `uses` statement references a `grouping`, **When** processed, **Then** grouping contents are expanded inline (no grouping references in JSON)

---

## JSON Output Format

### Node Structure

```json
{
  "name": "interfaces",
  "type": "container",
  "description": "Interface configuration",
  "children": [
    {
      "name": "interface",
      "type": "list",
      "key": "name",
      "children": [
        {"name": "name", "type": "leaf", "leaf-type": "string"},
        {"name": "description", "type": "leaf", "leaf-type": "string"},
        {"name": "disable", "type": "leaf", "leaf-type": "empty"},
        {"name": "unit", "type": "list", "key": "name", "children": [...]}
      ]
    }
  ]
}
```

### Root Structure

- **Given** YANG modules are processed, **When** JSON emitted, **Then** root object has:
  - `"name": "configuration"` — top-level Junos config container
  - `"children": [...]` — all top-level config sections
  - `"identities": [...]` — all YANG identity definitions (for identityref resolution)

---

## Usage

### Basic Invocation

```bash
# Print plugin directory
jtaf-pyang-plugindir
# Output: /path/to/jtaf_pyang_plugin/

# Generate JSON schema from YANG
pyang --plugindir $(jtaf-pyang-plugindir) -f jtaf \
  -p examples/yang/18.2/18.2R3/common \
  examples/yang/18.2/18.2R3/junos-qfx/conf/*.yang > schema.json

# Pipe directly to provider generation
pyang --plugindir $(jtaf-pyang-plugindir) -f jtaf \
  -p examples/yang/18.2/18.2R3/common \
  examples/yang/18.2/18.2R3/junos-qfx/conf/*.yang | \
  jtaf-provider -j - -x configs/*.xml -t vqfx
```

### Supported YANG Versions

The plugin works with any version in the `yang/` repo:
- 14.2, 16.1, 16.2, 17.1–17.4, 18.1–18.4
- 19.1–19.4, 20.1–20.4, 21.1–21.4
- 22.1–22.4, 23.1–23.4, 24.2–24.4, 25.2–25.4

Device families: junos-qfx, junos-vmx, junos-ex, junos-srx, junos-ptx, junos-acx

---

## Error Handling

- **Given** a YANG file has syntax errors, **When** pyang parses it, **Then** pyang reports errors to stderr (plugin never reached)
- **Given** a YANG type is unrecognized, **When** type mapping occurs, **Then** default to `"string"` type
- **Given** a `leafref` target path is unresolvable, **When** resolution attempted, **Then** fall back to `"string"` type

---

## Testing

### How to Run

```bash
# All pyang plugin tests
pytest jtaf_pyang_plugin/tests/ -v

# Specific test
pytest jtaf_pyang_plugin/tests/test_jtaf_json.py -v

# Single test function
pytest jtaf_pyang_plugin/tests/test_jtaf_json.py::test_container_output -v

# Integration test (requires YANG files)
pyang --plugindir $(jtaf-pyang-plugindir) -f jtaf \
  -p examples/yang/18.2/18.2R3/common \
  examples/yang/18.2/18.2R3/junos-qfx/conf/junos-qfx-conf-interfaces@2018-01-01.yang | \
  python -m json.tool > /dev/null && echo "Valid JSON"
```

### Test Coverage Map

| Test Case | Given | When | Then |
|-----------|-------|------|------|
| Container handling | YANG container statement | Plugin processes it | JSON node with `type: "container"` and children |
| List handling | YANG list with `key "name"` | Plugin processes it | JSON node with `type: "list"` and `key: "name"` |
| Leaf-list handling | YANG leaf-list statement | Plugin processes it | JSON node with `type: "leaf-list"` |
| Enum extraction | YANG `enumeration` type with values | Plugin processes it | `leaf-type: "enumeration"` with values list |
| Identity-ref | YANG `identityref` with base | Plugin processes it | `leaf-type: "identityref"` with base name |
| Empty type | YANG `empty` type leaf | Plugin processes it | `leaf-type: "empty"` |
| Ordered leaf-list | `ordered-by user` statement | Plugin processes it | `"ordered-by": "user"` flag in output |
