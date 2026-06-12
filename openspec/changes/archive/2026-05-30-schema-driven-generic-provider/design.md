# Design: Schema-Driven Generic Provider

## Overview

Replace the Jinja2-generated per-node Go code with a fixed-size **generic runtime library** that interprets `trimmed_schema.json` at runtime. The generated provider output shrinks from ~7 code blocks per YANG node (potentially 100K+ lines) to a thin wiring layer (~200 lines) plus the JSON schema.

## Architecture

```
┌──────────────────────────────────────────────────────────────────────────┐
│  BEFORE (Current)                                                        │
│                                                                          │
│  pyang → JSON → Jinja2 template → MASSIVE resource_config_provider.go    │
│                                                                          │
│  ┌────────────────────────────────────┐                                  │
│  │ xml_A struct { ... }               │ ×N nodes                         │
│  │ xml_A_B struct { ... }             │                                  │
│  │ A_B_Model struct + AttrTypes()     │                                  │
│  │ A_B_Model.Attributes()             │                                  │
│  │ modelToConfig() unrolled loops     │                                  │
│  │ configToModel() unrolled loops     │                                  │
│  └────────────────────────────────────┘                                  │
├──────────────────────────────────────────────────────────────────────────┤
│  AFTER (Proposed)                                                        │
│                                                                          │
│  pyang → JSON → jtaf-provider copies fixed Go library + embeds JSON      │
│                                                                          │
│  ┌────────────────────────────────────┐  ┌─────────────────────────────┐ │
│  │ trimmed_schema.json (embedded)     │  │ generic/ package (FIXED)    │ │
│  │ provider.go (wiring, ~100 lines)   │  │  schema_loader.go          │ │
│  │ config.go (unchanged)              │  │  tf_schema_builder.go      │ │
│  │ main.go (unchanged)               │  │  xml_tree.go               │ │
│  └────────────────────────────────────┘  │  model_convert.go          │ │
│                                          │  resource.go               │ │
│                                          └─────────────────────────────┘ │
└──────────────────────────────────────────────────────────────────────────┘
```

## Key Components

### 1. Schema Loader (`generic/schema_loader.go`)

Reuses and extends the existing `patch.SchemaNode` / `patch.NodeInfo` types. Loads `trimmed_schema.json` into an indexed tree that drives all other components.

```go
// SchemaIndex holds the parsed schema for runtime use
type SchemaIndex struct {
    Root     *SchemaNode
    ByPath   map[string]*NodeInfo  // reuse from patch package
    TopLevel []*SchemaNode         // direct children of "configuration"
}

func LoadSchema(jsonBytes []byte) (*SchemaIndex, error)
```

The existing `patch.UnmarshalTrimmedSchemaIndex` already does 90% of this work. We extend it to also preserve the tree structure (not just the flat index).

### 2. Terraform Schema Builder (`generic/tf_schema_builder.go`)

Dynamically constructs `schema.Schema` from the schema tree. Each YANG node maps to a Terraform schema attribute:

| YANG Type | Terraform Schema Type |
|-----------|----------------------|
| `leaf` | `schema.StringAttribute` |
| `leaf-list` | `schema.ListAttribute{ElementType: types.StringType}` |
| `container` | `schema.SingleNestedAttribute` |
| `list` | `schema.ListNestedAttribute` |

```go
func BuildTerraformSchema(idx *SchemaIndex) schema.Schema {
    attrs := map[string]schema.Attribute{
        "resource_name": schema.StringAttribute{
            Required:      true,
            PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplace()},
        },
    }
    for _, top := range idx.TopLevel {
        attrs[normalize(top.Name)] = buildNestedAttribute(top)
    }
    return schema.Schema{Attributes: attrs}
}

func buildNestedAttribute(node *SchemaNode) schema.Attribute {
    // Recursively builds nested schema from YANG tree
}
```

### 3. Generic XML Tree (`generic/xml_tree.go`)

Reuses `patch.Node` and `patch.BuildTree` for XML parsing. Adds XML serialization from the generic tree (inverse of BuildTree):

```go
// SerializeTree converts a *patch.Node tree back to XML bytes
func SerializeTree(root *patch.Node) ([]byte, error)

// BuildTreeFromModel constructs a *patch.Node tree from Terraform model state
func BuildTreeFromModel(ctx context.Context, model map[string]any, idx *SchemaIndex) *patch.Node
```

### 4. Model Converter (`generic/model_convert.go`)

Bidirectional conversion between Terraform Framework types and generic XML trees. Replaces the generated `modelToConfig()` / `configToModel()` functions.

```go
// ModelToXMLBytes converts Terraform plan/state into XML config bytes
func ModelToXMLBytes(ctx context.Context, attrs map[string]attr.Value, idx *SchemaIndex) ([]byte, error)

// XMLBytesToModel converts device XML config into Terraform state values
func XMLBytesToModel(ctx context.Context, xmlBytes []byte, idx *SchemaIndex) (map[string]attr.Value, error)
```

**Key insight**: The Terraform Plugin Framework allows working with `attr.Value` maps directly via `req.Plan.GetAttribute()` and `resp.State.SetAttribute()`. We don't need typed Go structs for the model.

### 5. Generic Resource (`generic/resource.go`)

A single `configResource` implementation that handles CRUD for any schema:

```go
type configResource struct {
    client ProviderConfig
    idx    *SchemaIndex
}

func (r *configResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
    resp.Schema = BuildTerraformSchema(r.idx)
}

func (r *configResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
    // 1. Extract plan as generic attr.Value map
    // 2. Convert to XML via ModelToXMLBytes
    // 3. Send via NETCONF (same as today)
    // 4. Read back, convert to model, set state
}

// Update uses existing patch engine - no change needed
func (r *configResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
    // 1. Get plan XML and state XML (via ModelToXMLBytes)
    // 2. Use patch.LeafMapWithSchema + patch.ComputeDiff (unchanged)
    // 3. Apply patch via NETCONF (unchanged)
}
```

## Data Flow

```
                    Terraform Plan
                         │
                         ▼
              ┌─────────────────────┐
              │  Generic Resource   │
              │  (Create/Update)    │
              └─────────┬───────────┘
                        │
              ┌─────────▼───────────┐
              │  ModelToXMLBytes()  │  Walk schema tree + TF attr.Values → XML
              └─────────┬───────────┘
                        │
              ┌─────────▼───────────┐
              │  patch engine       │  LeafMapWithSchema → ComputeDiff → CreateDiffPatch
              │  (UNCHANGED)        │
              └─────────┬───────────┘
                        │
              ┌─────────▼───────────┐
              │  NETCONF Client     │  SendUpdate / SendDirectTransaction
              │  (UNCHANGED)        │
              └─────────┬───────────┘
                        │
              ┌─────────▼───────────┐
              │  Read back config   │  MarshalConfig → raw XML bytes
              └─────────┬───────────┘
                        │
              ┌─────────▼───────────┐
              │  XMLBytesToModel()  │  Parse XML → walk schema → build TF values
              └─────────┬───────────┘
                        │
                        ▼
                 Terraform State
```

## Integration With Existing Code

### What stays unchanged
- `patch/` package — all of it (already generic)
- `netconf/` package — all of it
- `main.go` — provider bootstrap
- `config.go` — NETCONF client setup
- pyang plugin (`jtaf_json.py`) — schema generation

### What changes
- `resource_config_provider.go.j2` → replaced by a minimal template that just embeds JSON and wires up the generic resource
- `jtaf-provider` script → copies `generic/` package into output, no longer renders massive template
- `provider.go.j2` → minor update: resource factory references generic resource

### NETCONF Client Interface

The existing NETCONF client interface needs one small adaptation. Currently `MarshalConfig` unmarshals into a typed struct:

```go
// Current: requires typed xml_Configuration struct
func (c *Client) MarshalConfig(v interface{}) error

// For generic approach: return raw XML bytes instead
func (c *Client) GetConfigXML() ([]byte, error)
```

We can add the new method alongside the existing one for backward compatibility.

## Schema → TF Attribute Mapping Details

```
YANG leaf (type string/int/etc.)  →  schema.StringAttribute (all values as strings)
YANG leaf (type empty)            →  schema.StringAttribute (presence = "")  
YANG leaf-list                    →  schema.ListAttribute{ElementType: StringType}
YANG container                    →  schema.SingleNestedAttribute{Attributes: ...}
YANG list                         →  schema.ListNestedAttribute{NestedObject: ...}
```

All leaf values are represented as `types.String` in Terraform (matching current behavior). This simplifies conversion and maintains compatibility with existing `.tf` files.

## Generated Output Structure

After running `jtaf-provider -j schema.json -t full`:

```
terraform-provider-junos-full/
├── main.go                    # unchanged bootstrap
├── provider.go                # wiring (references generic resource)
├── config.go                  # NETCONF client config
├── go.mod                     # module declaration
├── trimmed_schema.json        # the YANG schema (data, not code)
├── generic/                   # fixed-size runtime library
│   ├── schema_loader.go       # ~150 lines
│   ├── tf_schema_builder.go   # ~200 lines
│   ├── xml_tree.go            # ~150 lines
│   ├── model_convert.go       # ~300 lines
│   └── resource.go            # ~250 lines
├── patch/                     # unchanged
│   ├── diff.go
│   ├── leafmap.go
│   ├── ...
└── netconf/                   # unchanged
    ├── client.go
    └── ...
```

**Total fixed Go code: ~1050 lines** (constant regardless of YANG model size)
**Variable data: trimmed_schema.json** (grows with model but is just data, not code)

## Testing Strategy

1. **Unit tests for generic library** — test schema loading, TF schema building, model↔XML conversion with known schema fragments
2. **Existing patch tests** — unchanged, continue passing
3. **Integration tests** — run Terraform plan/apply against mock NETCONF server using a small schema subset
4. **Regression** — generate provider with both old and new approach for same XML-filtered schema; verify Terraform produces same plan

## Migration Path

1. Build `generic/` package with tests
2. Add new template mode to `jtaf-provider` (`--mode generic` flag, default initially stays as `legacy`)
3. Validate with existing test schemas
4. Flip default to `generic`
5. Remove old Jinja2 template code (optional, can keep for reference)
