# Patch Engine

Go package that computes and applies minimal NETCONF `edit-config` operations between desired and actual Junos configuration state. Lives at `terraform_provider/patch/`.

## Architecture

```
Current XML (device) ─┐
                      ├→ LeafMapWithSchema() → path→value maps
Desired XML (plan) ──┘
                      ↓
                ComputeDiff()
                      ↓
         Create | Replace | Delete operations
                      ↓
        CreateDiffPatchWithSchema() → NETCONF edit-config XML
                      ↓
               NETCONF RPC → Junos device
```

## Core Files

| File | Purpose |
|------|---------|
| `leafmap.go` | Flattens XML config into leaf path→value map using YANG schema |
| `diff.go` | Computes CRUD operations between state and plan leaf maps |
| `patch.go` | Builds NETCONF `<configuration>` XML with `nc:operation` attributes from diff |
| `tree.go` | XML tree parsing (`BuildTree`) and marshaling |
| `path.go` | Path string parsing, segment extraction, tree navigation (`ensurePath`) |
| `order.go` | Enforces element ordering for `ordered-by user` leaf-lists |
| `process_schema.go` | Parses trimmed JSON schema into `map[string]*NodeInfo` index |
| `types.go` | Core data structures: `Node`, `Change`, `ChangeType`, `NodeKind`, `NodeInfo` |

## Key Types

```go
type Node struct {
    Tag, Text, Operation string
    Attrs     map[string]string
    Children  []*Node
    Parent    *Node
}

type Change struct {
    Op     ChangeType  // Create | Replace | Delete
    OldVal string
    NewVal string
}

type NodeInfo struct {
    Path, Name    string
    Kind          NodeKind  // Container | List | Leaf | LeafList
    ListKey       string
    OrderedByUser bool
    Leaf          LeafInfo
}
```

## Key Functions

| Function | Signature | Purpose |
|----------|-----------|---------|
| `LeafMapWithSchema` | `(node *Node, idx map[string]*NodeInfo) map[string]string` | Flatten XML into path→value map, respecting list keys, leaf-lists, empty elements |
| `ComputeDiff` | `(stateMap, planMap map[string]string) map[string]Change` | Compute minimal CRUD operations between two leaf maps |
| `CreateDiffPatchWithSchema` | `(diffMap map[string]Change, groupName string, idx map[string]*NodeInfo) ([]byte, error)` | Build NETCONF edit-config XML from diff |
| `ProcessSchema` | `(rawJSON []byte) (map[string]*NodeInfo, error)` | Parse trimmed JSON schema into lookup index |
| `AlignXMLOrderToReference` | `(target, reference *Node)` | Reorder elements to match reference ordering |
| `BuildTree` | `(xmlBytes []byte) (*Node, error)` | Parse XML into Node tree |

## Corner Cases

- **Ordered leaf-lists** — VRRP virtual-address order matters for Junos; `order.go` aligns output to reference
- **YANG `empty` type** — Presence semantics (e.g., `<disable/>`, `<vlan-tagging/>`); represented as path→"" in leaf map
- **Container presence vs. content** — Empty containers should not generate spurious diffs
- **Multi-value list keys** — Interfaces keyed by `[name=X]`; compound keys like `key = "choice-ident choice-value"`
- **Special characters in keys** — Slashes, nested brackets in path segments handled by `splitPathRespectingQuotes()`
- **Container delete coalescing** — When all children deleted, coalesce to single container delete

## Tests

| Test File | Focus |
|-----------|-------|
| `patch_test.go` | Diff patch generation, operation attributes |
| `matrix_test.go` | Comprehensive CRUD matrix (create/replace/delete combinations) |
| `corner_case_test.go` | Special chars in keys, container delete coalescing |
| `schema_phase1_test.go` | Schema parsing, trimmed schema validation |
| `order_test.go` | Element ordering preservation |
