# Patch Engine Specification

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
               AlignXMLOrderToReference() → ordered XML
                      ↓
               NETCONF RPC → Junos device
```

---

## LeafMapWithSchema — Flattening XML to Path→Value Map

**Location:** `terraform_provider/patch/leafmap.go`

### Behaviors

#### Empty Container/List Handling

- **Given** a node has no children AND `info.Kind` is Container or List, **When** `LeafMapWithSchema()` processes it, **Then** the node is skipped (no leaves emitted)

#### Empty Leaf (YANG `empty` Type)

- **Given** `info.Kind == Leaf` AND the element has no text content, **When** processed, **Then** emit `path = ""` (presence-only semantics)
- **Example:** `<vlan-tagging/>` → `"configuration/interfaces/interface/vlan-tagging" = ""`
- **Example:** `<disable/>` → `"configuration/protocols/ospf/area/interface/disable" = ""`

#### Ordered-by-User Leaf-Lists

- **Given** `info.OrderedByUser == true` AND `info.Kind == LeafList`, **When** processed, **Then**:
  - A per-tag position counter is tracked
  - Each entry emits: `path/tag[pos=N] = value`
  - This enables reordering detection (position change → Replace diff, not invisible)
- **Example:** VRRP virtual-address list: `[pos=0] = "10.0.0.1"`, `[pos=1] = "10.0.0.2"`

#### Keyed List Structural Shortcut

- **Given** `info.Kind == List` AND `info.ListKey` is a single field (not compound), **When** the list entry has NO material leaves except the key itself, **Then** emit only `path/listKey = keyValue`
- **Given** compound keys (space-separated) OR material non-key leaves exist, **When** processed, **Then** recurse normally into all children

#### Key Child Exclusion

- **Given** a child element's tag matches the list key name, **When** recursing into list children, **Then** the key child is NOT emitted separately (it's already expressed in the path predicate)

---

## ComputeDiff — Calculating CRUD Operations

**Location:** `terraform_provider/patch/diff.go`

### Behaviors

#### Diff Generation Rules

| Given | When | Then |
|-------|------|------|
| Path exists in state only | Computing diff | Emit `Delete` with `OldVal=stateVal, NewVal=""` |
| Path exists in both, values differ | Computing diff | Emit `Replace` with `OldVal=stateVal, NewVal=planVal` |
| Path exists in plan only | Computing diff | Emit `Create` with `OldVal="", NewVal=planVal` |
| Path exists in both, values identical | Computing diff | Path is omitted from diff (no-op) |

#### UTF-8 Double-Encoding Repair

- **Given** a string contains runes in range 0x80–0x9F (UTF-8 lead bytes misinterpreted as Latin-1), **When** `NormalizeLeafMapUTF8()` is called, **Then** attempt byte-level repair
- **Given** the repaired string is shorter than original AND valid UTF-8, **When** repair succeeds, **Then** use repaired version
- **Example:** em-dash U+2014 (UTF-8: `E2 80 94`) double-encoded as `"â\x80\x94"` → repaired to `"—"`

---

## CreateDiffPatchWithSchema — Building NETCONF XML

**Location:** `terraform_provider/patch/patch.go`

### Behaviors

#### Two-Pass Patch Building

**Pass 1: Keyed List Entry Operations**
- **Given** a leaf is a list key (e.g., `interface[name=ge-0/0/0]/name`), **When** building patch, **Then** apply the operation to the PARENT node (`interface[name=ge-0/0/0]` gets `nc:operation`)
- The key leaf itself is NOT emitted as a separate operation element

**Pass 2: Leaf-Level Operations**

| Given Operation | When Building XML | Then Output |
|---|---|---|
| Create | Leaf node | `<tag nc:operation="create">newVal</tag>` |
| Replace | Leaf node | `<tag nc:operation="replace">newVal</tag>` |
| Delete + leaf-list | Leaf-list entry | `<tag>oldVal</tag>` (Junos needs old value to match for deletion) |
| Delete + scalar | Scalar leaf | `<tag nc:operation="delete"/>` (no text; Junos rejects text on scalar delete) |
| Delete | Parent already has `nc:operation="delete"` | Child omits its own operation attribute |

#### Operation Ordering Rules

- **Given** a diff map with mixed operations, **When** `orderedChanges()` sorts them, **Then**:
  - Delete operations come first (precedence 0)
  - Replace operations second (precedence 1)
  - Create operations last (precedence 2)
- **Given** same operation type AND Delete, **When** ordering, **Then** deeper paths first (post-order: delete leaf before parent)
- **Given** same operation type AND Create, **When** ordering, **Then** shallower paths first (pre-order: create parent before leaf)
- **Given** all else equal, **When** ordering, **Then** sort lexicographically by path

#### Container Delete Coalescing

- **Given** ALL leaves under a container are Deletes (no Creates or Replaces) AND schema info is available, **When** `coalesceContainerDeletes()` runs, **Then**:
  - Emit single `<container nc:operation="delete"/>` instead of individual leaf deletes
  - Individual leaf delete entries are removed from the diff

---

## AlignXMLOrderToReference — Ordering Elements

**Location:** `terraform_provider/patch/order.go`

### Behaviors

#### Reference Order Priority

- **Given** an entry is present in the reference XML, **When** sorting siblings, **Then** sort by reference position
- **Given** an entry is absent from reference, **When** sorting, **Then** place deterministically after all reference entries

#### Instance-Aware Ordering

- **Given** a parent is a keyed list, **When** ordering children, **Then** each list entry gets independent child ordering
- **Example:** `interface[name=eth0]/unit` order is independent of `interface[name=eth1]/unit` order

#### Sort Key Construction

- **Given** a child matches a reference entry, **When** building sort key, **Then** `hasReference=true, referenceRank=N`
- **Given** a child is not in reference, **When** building sort key, **Then** `hasReference=false, referenceRank=2^30` (sorts last)
- **Given** same reference rank, **When** tie-breaking, **Then** compare by `tag < identity < text`

#### Node Identity Rules

| Given Node Kind | When Computing Identity | Then Identity String |
|---|---|---|
| List | Has key attribute | `"tag[key=value]"` (e.g., `"interface[name=ge-0/0/0]"`) |
| LeafList | Has text value | `"tag[value=text]"` (e.g., `"address[value=10.0.0.1]"`) |
| Other with text | Has text | `"tag=text"` |
| Other without text | No text | `"tag"` |

#### No Reference Behavior

- **Given** `referenceXML` is nil, **When** alignment runs, **Then** use deterministic fallback ordering (tag, identity, text)

---

## ProcessSchema — Schema Index Construction

**Location:** `terraform_provider/patch/process_schema.go`

### Behaviors

#### Path Construction

- **Given** node has `path=""` and is root, **When** processing, **Then** skip (walk children only)
- **Given** node has a name, **When** processing, **Then** full path = `parentPath + "/" + nodeName`

#### Type Resolution

| Given `node.Type` | When Processing | Then `info.Kind` |
|---|---|---|
| `"container"` | Type resolved | `KindContainer` |
| `"list"` | Type resolved | `KindList`; capture `node.Key` as `ListKey` |
| `"leaf"` | Type resolved | `KindLeaf`; parse leaf-type |
| `"leaf-list"` | Type resolved | `KindLeafList` |

- **Given** `node.OrderedBy == "user"`, **When** processing, **Then** `info.OrderedByUser = true`

#### Leaf Type Classification

| Given `leaf-type` | When Classifying | Then `Base` |
|---|---|---|
| `"string"` | Classifying | `LeafString` |
| `"union"` | Classifying | `LeafUnion`; index all union branches |
| `"enumeration"` | Classifying | `LeafEnum`; build enum value set |
| `"empty"` | Classifying | `LeafOther` (presence leaf) |

---

## Corner Cases

### Special Characters in Path Keys

- **Given** a list key contains slashes (e.g., `ge-0/0/0`), **When** `splitPathRespectingQuotes()` runs, **Then** slashes inside `[key=value]` brackets are NOT treated as path separators
- **Given** a list key contains nested brackets (e.g., `ALLOW[ALL]`), **When** parsing, **Then** nested brackets inside the key value are preserved

### Container Presence vs. Content

- **Given** an empty container in XML (e.g., `<ospf></ospf>`), **When** building leaf map, **Then** NO entries emitted (avoids spurious diffs between "container exists" and "container with no config")

### Compound List Keys

- **Given** a list has `key = "choice-ident choice-value"` (space-separated), **When** processing, **Then** both fields are used together to construct the path predicate

---

## Testing

### How to Run

```bash
# All patch engine tests
cd terraform_provider && go test ./patch/ -v

# Specific test
cd terraform_provider && go test ./patch/ -run TestCC10_ParseSegment -v

# Corner case tests only
cd terraform_provider && go test ./patch/ -run TestCC -v

# Matrix tests (comprehensive CRUD combinations)
cd terraform_provider && go test ./patch/ -run TestMatrix -v

# Order alignment tests
cd terraform_provider && go test ./patch/ -run TestAlignOrder -v

# With race detection
cd terraform_provider && go test ./patch/ -race -v

# Coverage report
cd terraform_provider && go test ./patch/ -coverprofile=coverage.out && go tool cover -html=coverage.out
```

### Test Coverage Map

| Test File | What It Validates |
|-----------|-------------------|
| `patch_test.go` | Diff→patch XML generation, `nc:operation` attributes, group wrapping |
| `matrix_test.go` | All CRUD permutations (create/replace/delete in various combinations) |
| `corner_case_test.go` | Slashes in keys, nested brackets, container delete coalescing, UTF-8 repair |
| `schema_phase1_test.go` | JSON schema parsing, type resolution, list key extraction |
| `order_test.go` | Sibling reordering, per-instance ordering, extra entries sorting |

### Key Test Behaviors

| Test Name | Given | When | Then |
|-----------|-------|------|------|
| `TestCC10_ParseSegment_KeyWithSlashes` | Key value `ge-0/0/0` | Segment parsed | Slashes preserved in key, not split |
| `TestCC10_ParseSegment_KeyWithNestedBrackets` | Key value `ALLOW[ALL]` | Segment parsed | Nested brackets preserved |
| `TestCC5_ContainerDeleteCoalescing` | All children deleted | Coalescing runs | Single container delete emitted |
| `TestCC2_EmptyLeafToggle` | YANG empty type leaf | LeafMap built | Path emitted with value `""` |
| `TestCC6_UTF8DiffWithEncodingVariation` | Double-encoded UTF-8 | Normalize runs | Bytes repaired to valid UTF-8 |
| `TestAlignOrder_TopLevelListReorder` | List entries out of order | Alignment runs | Entries match reference order |
| `TestAlignOrder_ExtraEntriesInCurrent` | Entries not in reference | Alignment runs | Extras sorted after reference entries |
