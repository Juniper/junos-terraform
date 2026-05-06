# Junos Configuration Test Matrix

> End-to-end test plan covering every structural way Junos config can be
> modified — from YANG model through NETCONF patch to Terraform provider to
> physical/mock devices.

---

## 1. The Four YANG Node Types

Every Junos configuration element maps to exactly one of these YANG types.
Tests must exercise each type independently and in combination.

| # | YANG Type | XML Shape | Terraform Type | Identity | Example |
|---|-----------|-----------|----------------|----------|---------|
| 1 | **container** | `<system>…</system>` | `types.List` (1-element) | Tag only | `system`, `chassis`, `protocols` |
| 2 | **list** | `<interface><name>ge-0/0/0</name>…</interface>` (repeating) | `types.List` (N-elements) | Tag + key child (`name`, `id`, `type`) | `interfaces/interface`, `syslog/host`, `bgp/group` |
| 3 | **leaf** | `<description>text</description>` | `types.String` / `types.Int64` | Tag + text | `host-name`, `device-count`, `vlan-id` |
| 4 | **leaf-list** | `<protocol>bgp</protocol><protocol>direct</protocol>` (repeating same tag) | `types.List` of scalars | Tag + text value | `routing-options/static/route/next-hop`, `policy-options/community/members` |

---

## 2. All Structural Mutation Types

These are every way a Junos configuration can change. Each row is a distinct
test scenario.

### 2.1 Leaf Operations

| ID | Operation | Before XML | After XML | NETCONF `nc:operation` | Diff Type |
|----|-----------|------------|-----------|------------------------|-----------|
| L1 | **Create leaf** | *(absent)* | `<description>new</description>` | `create` | Create |
| L2 | **Replace leaf** | `<description>old</description>` | `<description>new</description>` | `replace` | Replace |
| L3 | **Delete leaf** | `<description>old</description>` | *(absent)* | `delete` | Delete |
| L4 | **Replace leaf with special chars** | `<contact>a@b</contact>` | `<contact>x&amp;y</contact>` | `replace` | Replace |
| L5 | **Replace leaf with embedded newlines** | `<message>line1\nline2</message>` | `<message>line1\nline2\nline3</message>` | `replace` | Replace |
| L6 | **Replace leaf (same value, no-op)** | `<host-name>r1</host-name>` | `<host-name>r1</host-name>` | *(none)* | *(empty diff)* |

### 2.2 Leaf-List Operations

| ID | Operation | Before XML | After XML | NETCONF | Diff Type |
|----|-----------|------------|-----------|---------|-----------|
| LL1 | **Add entry to leaf-list** | `<members>a</members>` | `<members>a</members><members>b</members>` | `create` on `members[value=b]` | Create |
| LL2 | **Remove entry from leaf-list** | `<members>a</members><members>b</members>` | `<members>a</members>` | `delete` on `members[value=b]` | Delete |
| LL3 | **Replace entry in leaf-list** | `<members>a</members><members>b</members>` | `<members>a</members><members>c</members>` | `delete` b + `create` c | Delete + Create |
| LL4 | **Reorder leaf-list** | `<members>b</members><members>a</members>` | `<members>a</members><members>b</members>` | *(none — set semantics)* | *(empty diff)* |
| LL5 | **Delete all leaf-list entries** | `<members>a</members><members>b</members>` | *(absent)* | `delete` a + `delete` b | Delete × 2 |
| LL6 | **Create leaf-list from scratch** | *(absent)* | `<members>x</members><members>y</members>` | `create` x + `create` y | Create × 2 |

### 2.3 List Operations (Keyed by `<name>` or another key)

| ID | Operation | Before XML | After XML | NETCONF | Diff Type |
|----|-----------|------------|-----------|---------|-----------|
| K1 | **Add new list entry** | 1 `<interface>` | 2 `<interface>`s (new key) | `create` on `<interface>` entry | Create |
| K2 | **Delete list entry** | 2 `<interface>`s | 1 `<interface>` (removed key) | `delete` on `<interface>` entry | Delete |
| K3 | **Rename list key** | `<name>ge-0/0/0</name>` | `<name>ge-0/0/1</name>` | `delete` old entry + `create` new entry | Delete + Create |
| K4 | **Modify leaf inside list entry** | `<unit>…<description>old</description>` | `<unit>…<description>new</description>` | `replace` on nested leaf | Replace |
| K5 | **Add leaf inside existing list entry** | `<unit>…` (no description) | `<unit>…<description>new</description>` | `create` on nested leaf | Create |
| K6 | **Delete leaf inside list entry** | `<unit>…<description>old</description>` | `<unit>…` (no description) | `delete` on nested leaf | Delete |
| K7 | **Reorder list entries** | `[host=b, host=a]` | `[host=a, host=b]` | *(none — key identity, not position)* | *(alignment only)* |
| K8 | **Nested list: add entry in child list** | `host/contents` has 3 entries | `host/contents` has 4 entries | `create` on child list entry | Create |
| K9 | **Nested list: different order in different parent entries** | host[a] contents `[x,y]`, host[b] contents `[y,x]` | *(same values, different order)* | *(alignment per-instance)* | *(empty diff)* |
| K10 | **Key-only list entry (structural)** | *(absent)* | `<file><name>security</name></file>` | `create` on key leaf | Create |

### 2.4 Container Operations

| ID | Operation | Before XML | After XML | NETCONF | Diff Type |
|----|-----------|------------|-----------|---------|-----------|
| C1 | **Create container** (with children) | *(absent)* | `<chassis><aggregated-devices>…` | `create` on contained leaves | Create |
| C2 | **Delete container** (all children) | `<chassis><aggregated-devices>…` | *(absent)* | `delete` on all contained leaves | Delete |
| C3 | **Empty container (presence)** | *(absent)* | `<ssh/>` | `create` empty element | Create |
| C4 | **Delete empty container** | `<ssh/>` | *(absent)* | `delete` empty element | Delete |
| C5 | **Modify children within container** | `<system><host-name>old</host-name>` | `<system><host-name>new</host-name>` | `replace` on leaf child | Replace |

### 2.5 Compound / Cross-Type Operations

| ID | Operation | Description | Components |
|----|-----------|-------------|------------|
| M1 | **Mixed: modify + add + delete in one apply** | Change description, add new interface, remove old one | Replace + Create + Delete |
| M2 | **Deep nesting: 4+ levels** | `system/syslog/host[name=log]/contents[name=any]/notice` | Nested list × 2 + leaf |
| M3 | **Multiple resources, parallel apply** | Two independent resources updated simultaneously | Tests provider parallelism |
| M4 | **Full resource replacement (overwrite)** | Entire config changed significantly | Fallback to `SendDirectTransaction` |
| M5 | **Out-of-band drift + reconcile** | External NETCONF change → Terraform detects and fixes | Read drift + patch |

---

## 3. The Pipeline Stages

Each test scenario exercises this full pipeline. A failure at any stage must
be diagnosable to that specific stage.

```
┌──────────────┐    ┌──────────────┐    ┌──────────────────┐    ┌──────────────┐    ┌────────────┐
│  Input XML   │───▶│  YANG Model  │───▶│ Terraform Provider│───▶│   .tf Files  │───▶│  Device /  │
│  (device     │    │  (schema)    │    │ (Go code gen)     │    │  (HCL)       │    │  Mock      │
│   config)    │    │              │    │                   │    │              │    │            │
└──────────────┘    └──────────────┘    └──────────────────┘    └──────────────┘    └────────────┘
      │                   │                      │                     │                  │
      ▼                   ▼                      ▼                     ▼                  ▼
   jtaf-xml2tf      trimmed_schema.json    resource_config_      terraform plan     NETCONF RPC
   (parse XML,      (node types, keys,     provider.go           terraform apply    load-config /
    generate HCL)    leaf constraints)     (marshal, diff,                          edit-config /
                                            patch, CRUD)                           get-config
```

### Stage-by-Stage Verification Points

| Stage | What to Verify | How |
|-------|----------------|-----|
| **XML → .tf** | HCL output matches expected structure, correct quoting, no trailing `\n` | Unit test `jtaf-xml2tf` with known XML input |
| **YANG → schema** | `trimmed_schema.json` has correct `type`, `key`, node hierarchy | Unit test `UnmarshalTrimmedSchemaIndex` |
| **Schema → LeafMap** | Leaf map paths use schema keys, leaf-list gets `[value=X]` | Unit test `LeafMapWithSchema` |
| **LeafMap → Diff** | `ComputeDiff` produces expected Create/Replace/Delete set | Unit test `ComputeDiff` |
| **Diff → Patch XML** | `CreateDiffPatch` emits correct `nc:operation` attributes, ordering | Unit test `CreateDiffPatch` |
| **Patch → Device** | Mock/device accepts edit-config, candidate state is correct | Integration test via mock |
| **Device → Read** | `MarshalGroup` / `MarshalConfig` reads back correct XML | Integration test |
| **Read → State** | `configToModel` + `alignConfigToReference` produces matching state | `terraform plan -detailed-exitcode` == 0 |

---

## 4. Test Fixture Design

Each test case needs three artifacts:

### 4.1 Input Artifacts

```
tests/
  cases/
    L1_create_leaf/
      before.xml          # Device config before the change (or empty)
      after.xml           # Device config after the change
      schema.json         # Trimmed schema covering the elements used
      expected_diff.json  # Expected ComputeDiff output
      expected_patch.xml  # Expected CreateDiffPatch output
    L2_replace_leaf/
      ...
    LL1_add_leaf_list/
      ...
```

### 4.2 Minimal Schema for Test Cases

Each test case should use a self-contained schema snippet. Here is a
reference schema covering all four YANG types:

```json
{
  "path": "",
  "root": {
    "children": [
      {
        "name": "configuration",
        "type": "container",
        "path": "",
        "children": [
          {
            "name": "system",
            "type": "container",
            "path": "",
            "children": [
              {
                "name": "host-name",
                "type": "leaf",
                "path": "system",
                "leaf-type": "string"
              },
              {
                "name": "syslog",
                "type": "container",
                "path": "system",
                "children": [
                  {
                    "name": "host",
                    "type": "list",
                    "path": "system/syslog",
                    "key": "name",
                    "children": [
                      {
                        "name": "name",
                        "type": "leaf",
                        "path": "system/syslog/host",
                        "leaf-type": "string"
                      },
                      {
                        "name": "contents",
                        "type": "list",
                        "path": "system/syslog/host",
                        "key": "name",
                        "children": [
                          {
                            "name": "name",
                            "type": "leaf",
                            "path": "system/syslog/host/contents",
                            "leaf-type": "string"
                          },
                          {
                            "name": "any",
                            "type": "leaf",
                            "path": "system/syslog/host/contents",
                            "leaf-type": "empty"
                          },
                          {
                            "name": "notice",
                            "type": "leaf",
                            "path": "system/syslog/host/contents",
                            "leaf-type": "empty"
                          }
                        ]
                      }
                    ]
                  }
                ]
              },
              {
                "name": "services",
                "type": "container",
                "path": "system",
                "children": [
                  {
                    "name": "ssh",
                    "type": "container",
                    "path": "system/services",
                    "children": []
                  }
                ]
              }
            ]
          },
          {
            "name": "interfaces",
            "type": "container",
            "path": "",
            "children": [
              {
                "name": "interface",
                "type": "list",
                "path": "interfaces",
                "key": "name",
                "children": [
                  {
                    "name": "name",
                    "type": "leaf",
                    "path": "interfaces/interface",
                    "leaf-type": "string"
                  },
                  {
                    "name": "description",
                    "type": "leaf",
                    "path": "interfaces/interface",
                    "leaf-type": "string"
                  },
                  {
                    "name": "unit",
                    "type": "list",
                    "path": "interfaces/interface",
                    "key": "name",
                    "children": [
                      {
                        "name": "name",
                        "type": "leaf",
                        "path": "interfaces/interface/unit",
                        "leaf-type": "string"
                      },
                      {
                        "name": "description",
                        "type": "leaf",
                        "path": "interfaces/interface/unit",
                        "leaf-type": "string"
                      },
                      {
                        "name": "family",
                        "type": "container",
                        "path": "interfaces/interface/unit",
                        "children": [
                          {
                            "name": "inet",
                            "type": "container",
                            "path": "interfaces/interface/unit/family",
                            "children": [
                              {
                                "name": "address",
                                "type": "list",
                                "path": "interfaces/interface/unit/family/inet",
                                "key": "name",
                                "children": [
                                  {
                                    "name": "name",
                                    "type": "leaf",
                                    "path": "interfaces/interface/unit/family/inet/address",
                                    "leaf-type": "string"
                                  }
                                ]
                              }
                            ]
                          }
                        ]
                      }
                    ]
                  }
                ]
              }
            ]
          },
          {
            "name": "policy-options",
            "type": "container",
            "path": "",
            "children": [
              {
                "name": "community",
                "type": "list",
                "path": "policy-options",
                "key": "name",
                "children": [
                  {
                    "name": "name",
                    "type": "leaf",
                    "path": "policy-options/community",
                    "leaf-type": "string"
                  },
                  {
                    "name": "members",
                    "type": "leaf-list",
                    "path": "policy-options/community",
                    "leaf-type": "string"
                  }
                ]
              }
            ]
          }
        ]
      }
    ]
  }
}
```

---

## 5. Detailed Test Cases

### 5.1 Leaf Tests

#### L1 — Create Leaf

**Purpose:** Add a leaf that does not exist on the device.

**Before XML:**
```xml
<configuration>
  <system/>
</configuration>
```

**After XML:**
```xml
<configuration>
  <system>
    <host-name>router1</host-name>
  </system>
</configuration>
```

**Expected Diff:**
```
system/host-name → Create "router1"
```

**Expected Patch XML:**
```xml
<configuration>
  <system>
    <host-name nc:operation="create">router1</host-name>
  </system>
</configuration>
```

**Terraform verification:**
- `terraform plan` shows `+ host_name = "router1"`
- `terraform apply` succeeds
- Second `terraform plan` shows no changes (exit code 0)

---

#### L2 — Replace Leaf

**Purpose:** Change the value of an existing leaf.

**Before XML:**
```xml
<configuration>
  <system>
    <host-name>router1</host-name>
  </system>
</configuration>
```

**After XML:**
```xml
<configuration>
  <system>
    <host-name>router2</host-name>
  </system>
</configuration>
```

**Expected Diff:**
```
system/host-name → Replace "router1" → "router2"
```

**Expected Patch XML:**
```xml
<configuration>
  <system>
    <host-name nc:operation="replace">router2</host-name>
  </system>
</configuration>
```

---

#### L3 — Delete Leaf

**Purpose:** Remove a leaf from the device.

**Before XML:**
```xml
<configuration>
  <interfaces>
    <interface>
      <name>ge-0/0/0</name>
      <description>uplink</description>
    </interface>
  </interfaces>
</configuration>
```

**After XML:**
```xml
<configuration>
  <interfaces>
    <interface>
      <name>ge-0/0/0</name>
    </interface>
  </interfaces>
</configuration>
```

**Expected Diff:**
```
interfaces/interface[name=ge-0/0/0]/description → Delete "uplink"
```

**Expected Patch XML:**
```xml
<configuration>
  <interfaces>
    <interface>
      <name>ge-0/0/0</name>
      <description nc:operation="delete">uplink</description>
    </interface>
  </interfaces>
</configuration>
```

---

### 5.2 Leaf-List Tests

#### LL1 — Add Entry to Leaf-List

**Purpose:** Append a new value to a leaf-list (set semantics).

**Before XML:**
```xml
<configuration>
  <policy-options>
    <community>
      <name>my-comm</name>
      <members>target:65000:100</members>
    </community>
  </policy-options>
</configuration>
```

**After XML:**
```xml
<configuration>
  <policy-options>
    <community>
      <name>my-comm</name>
      <members>target:65000:100</members>
      <members>target:65000:200</members>
    </community>
  </policy-options>
</configuration>
```

**Expected Diff:**
```
policy-options/community[name=my-comm]/members[value=target:65000:200] → Create "target:65000:200"
```

**Expected Patch XML:**
```xml
<configuration>
  <policy-options>
    <community>
      <name>my-comm</name>
      <members nc:operation="create">target:65000:200</members>
    </community>
  </policy-options>
</configuration>
```

---

#### LL2 — Remove Entry from Leaf-List

**Before XML:**
```xml
<configuration>
  <policy-options>
    <community>
      <name>my-comm</name>
      <members>target:65000:100</members>
      <members>target:65000:200</members>
    </community>
  </policy-options>
</configuration>
```

**After XML:**
```xml
<configuration>
  <policy-options>
    <community>
      <name>my-comm</name>
      <members>target:65000:100</members>
    </community>
  </policy-options>
</configuration>
```

**Expected Diff:**
```
policy-options/community[name=my-comm]/members[value=target:65000:200] → Delete "target:65000:200"
```

---

#### LL3 — Replace Entry in Leaf-List (Delete + Create)

**Before:** members = `[a, b]`
**After:** members = `[a, c]`

**Expected Diff:**
```
.../members[value=b] → Delete "b"
.../members[value=c] → Create "c"
```

---

#### LL4 — Reorder Leaf-List (No-Op)

**Purpose:** Verify leaf-lists use set semantics — order does not matter.

**Before:** `<members>b</members><members>a</members>`
**After:** `<members>a</members><members>b</members>`

**Expected:** Empty diff. `terraform plan` shows no changes.

---

### 5.3 List Tests

#### K1 — Add New List Entry

**Before XML:**
```xml
<configuration>
  <interfaces>
    <interface>
      <name>ge-0/0/0</name>
      <description>uplink</description>
    </interface>
  </interfaces>
</configuration>
```

**After XML:**
```xml
<configuration>
  <interfaces>
    <interface>
      <name>ge-0/0/0</name>
      <description>uplink</description>
    </interface>
    <interface>
      <name>ge-0/0/1</name>
      <description>downlink</description>
    </interface>
  </interfaces>
</configuration>
```

**Expected Diff:**
```
interfaces/interface[name=ge-0/0/1]/name        → Create "ge-0/0/1"
interfaces/interface[name=ge-0/0/1]/description  → Create "downlink"
```

**Expected Patch XML:**
```xml
<configuration>
  <interfaces>
    <interface nc:operation="create">
      <name>ge-0/0/1</name>
      <description>downlink</description>
    </interface>
  </interfaces>
</configuration>
```

Note: When the leaf being created IS the list key, `nc:operation="create"` is
promoted to the parent `<interface>` element.

---

#### K2 — Delete List Entry

**Before:** interfaces `[ge-0/0/0, ge-0/0/1]`
**After:** interfaces `[ge-0/0/0]`

**Expected Patch XML:**
```xml
<configuration>
  <interfaces>
    <interface nc:operation="delete">
      <name>ge-0/0/1</name>
    </interface>
  </interfaces>
</configuration>
```

---

#### K3 — Rename List Key (Delete Old + Create New)

**Before:** interface `ge-0/0/0`
**After:** interface `ge-0/0/1` (with same children)

**Expected:** Delete entire old entry, create entire new entry.

---

#### K7 — Reorder List Entries (Alignment Only)

**Purpose:** Verify `AlignXMLOrderToReference` reorders by key identity.

**Before (device order):** `[host=test-patch, host=log]`
**Reference (.tf order):** `[host=log, host=test-patch]`

**Expected:** After alignment, device XML matches reference order.
Diff is empty. `terraform plan` shows no changes.

---

#### K9 — Nested Lists with Different Per-Instance Ordering

**Purpose:** Verify instance-aware ordering — each parent list entry gets its
own ordering bucket (the bug we fixed in `order.go`).

**Reference:**
- `host[name=log]` → contents order: `[kernel, any, pfe]`
- `host[name=test-patch]` → contents order: `[any, kernel, pfe]`

**Device returns:** both hosts with contents order `[any, kernel, pfe]`

**Expected after alignment:**
- `host[name=log]` → reordered to `[kernel, any, pfe]` (matches reference)
- `host[name=test-patch]` → kept as `[any, kernel, pfe]` (already matches)

---

### 5.4 Container Tests

#### C1 — Create Container with Children

**Before:** No `<chassis>` element.
**After:** `<chassis><aggregated-devices><ethernet><device-count>24</device-count>…`

**Expected:** All contained leaves appear as Create operations.

---

#### C3 — Empty Container (Presence Container)

**Purpose:** Some YANG containers are "presence" containers — they exist or
they don't. In Junos, these appear as self-closing XML tags.

**Before:** No `<ssh/>` under `system/services`.
**After:** `<system><services><ssh/></services></system>`

**Expected:** The empty element is created. The provider must preserve the
empty tag through the Create → Read cycle.

---

### 5.5 Compound Tests

#### M1 — Mixed Operations in Single Apply

**Purpose:** Verify that a single `terraform apply` correctly handles
create + replace + delete in one pass.

**Changes:**
1. Replace `system/host-name` from "r1" to "r2" (leaf replace)
2. Add new interface `ge-0/0/2` (list create)
3. Remove community member `target:65000:300` (leaf-list delete)
4. Delete interface `ge-0/0/3` description (leaf delete within list)

**Expected Patch XML ordering:**
1. Deletes first (deepest first)
2. Then replacements
3. Then creates (shallowest first)

---

#### M4 — Full Overwrite (Fallback Path)

**Purpose:** When the diff is so large that targeted patching fails
verification, the provider falls back to `SendDirectTransaction`
(full `load-configuration action="merge"`).

**Test:** After the patch-based update, inject a discrepancy so the
verification diff is non-empty. Verify the provider retries with full
replace and succeeds.

---

#### M5 — Out-of-Band Drift Detection and Reconciliation

**Purpose:** Simulate an external NETCONF change (not via Terraform) and
verify that `terraform plan` detects the drift and `terraform apply` restores
the desired state.

**Steps:**
1. Initial `terraform apply` succeeds (no-op plan = exit 0)
2. Directly send NETCONF `load-configuration merge` to change an IP
3. `terraform plan -detailed-exitcode` = exit 2 (drift detected)
4. `terraform apply` sends targeted patch to restore desired state
5. No-op plan = exit 0

---

## 6. Testing Layers

### 6.1 Go Unit Tests (fastest, most isolated)

**Location:** `terraform_provider/patch/*_test.go`

These test the core diff and patch engine without any NETCONF or Terraform
involvement. Use the minimal schema JSON and raw XML strings.

```go
func TestL1_CreateLeaf(t *testing.T) {
    idx := mustIdxFromSchema(testSchema)

    stateXML := `<configuration><system/></configuration>`
    planXML  := `<configuration><system><host-name>router1</host-name></system></configuration>`

    stateMap := LeafMapWithSchema(mustTree(stateXML), idx)
    planMap  := LeafMapWithSchema(mustTree(planXML), idx)

    diff := ComputeDiff(stateMap, planMap)
    assert(len(diff) == 1)
    assert(diff["system/host-name"].Op == Create)
    assert(diff["system/host-name"].NewVal == "router1")

    patchXML, err := CreateDiffPatch(diff, "")
    assert(err == nil)
    assert(strings.Contains(string(patchXML), `nc:operation="create"`))
    assert(strings.Contains(string(patchXML), `>router1<`))
}
```

**Full test list by ID:**

| Test Function | IDs Covered |
|---------------|-------------|
| `TestL1_CreateLeaf` | L1 |
| `TestL2_ReplaceLeaf` | L2 |
| `TestL3_DeleteLeaf` | L3 |
| `TestL4_ReplaceLeafSpecialChars` | L4 |
| `TestL6_ReplaceLeafNoOp` | L6 |
| `TestLL1_AddLeafListEntry` | LL1 |
| `TestLL2_RemoveLeafListEntry` | LL2 |
| `TestLL3_ReplaceLeafListEntry` | LL3 |
| `TestLL4_ReorderLeafListNoOp` | LL4 |
| `TestLL5_DeleteAllLeafListEntries` | LL5 |
| `TestLL6_CreateLeafListFromScratch` | LL6 |
| `TestK1_AddListEntry` | K1 |
| `TestK2_DeleteListEntry` | K2 |
| `TestK3_RenameListKey` | K3 |
| `TestK4_ModifyLeafInListEntry` | K4 |
| `TestK5_AddLeafInListEntry` | K5 |
| `TestK6_DeleteLeafInListEntry` | K6 |
| `TestK7_ReorderListEntries` | K7 (order.go) |
| `TestK8_AddNestedListEntry` | K8 |
| `TestK9_NestedListDifferentOrderPerInstance` | K9 (order.go) |
| `TestK10_KeyOnlyListEntry` | K10 |
| `TestC1_CreateContainer` | C1 |
| `TestC3_EmptyContainer` | C3 |
| `TestC5_ModifyChildrenInContainer` | C5 |
| `TestM1_MixedOperations` | M1 |

### 6.2 Mock Server Unit Tests (Python)

**Location:** `netconf_mock/tests/test_netconf_mock_server.py`

These test the mock server's ability to correctly apply each `nc:operation`.
For every patch XML generated by the Go layer, there should be a
corresponding Python test verifying the mock handles it correctly.

**New tests to add:**

| Test Function | IDs | What It Validates |
|---------------|-----|-------------------|
| `test_patch_create_leaf` | L1 | Mock adds `<host-name>` to empty `<system>` |
| `test_patch_replace_leaf` | L2 | Mock updates existing leaf text |
| `test_patch_delete_leaf` | L3 | Mock removes leaf element |
| `test_patch_create_leaf_list_entry` | LL1 | Mock appends new same-tag sibling |
| `test_patch_delete_leaf_list_entry` | LL2 | Mock removes specific leaf-list value |
| `test_patch_create_keyed_list_entry` | K1 | Mock appends new list entry |
| `test_patch_delete_keyed_list_entry` | K2 | Mock removes entry matching key |
| `test_patch_mixed_ops` | M1 | Mock applies mixed create+replace+delete |
| `test_readback_preserves_order` | K7 | get-configuration returns consistent order |

### 6.3 Integration Tests (Terraform + Mock)

**Location:** CI workflow (`go-terraform-provider.yml`) or local test harness

These test the full pipeline: `.tf` files → `terraform apply` →
provider NETCONF RPCs → mock server → readback → state comparison.

**Test flow for each scenario:**

```
1. Generate provider from YANG + scenario XMLs
2. Generate .tf files from scenario "after" XML
3. Start mock server seeded with scenario "before" XML
4. terraform apply (initial create)
5. terraform plan -detailed-exitcode (expect exit 0 = no-op)
6. Modify .tf to scenario "after2" XML (if testing update)
7. terraform apply (update via NETCONF patch)
8. terraform plan -detailed-exitcode (expect exit 0)
9. terraform destroy
```

---

## 7. Order-Specific Tests

Ordering is a critical concern because Junos XML element order is significant
for `<list>` and `<leaf-list>` but irrelevant for `<container>` children.

### 7.1 Test Matrix: Order Sensitivity

| Scenario | Element Type | Order Matters? | Provider Handling |
|----------|-------------|----------------|-------------------|
| `<interface>` siblings | list | Yes (key-based identity) | `AlignXMLOrderToReference` |
| `<contents>` inside host | nested list | Yes (per-instance) | Instance-aware alignment |
| `<members>` siblings | leaf-list | No (set semantics) | `[value=X]` path suffix |
| Container children (`<host-name>`, `<services>`) | mixed | No | Tag-based sorting |

### 7.2 Order Alignment Tests

| Test | Scenario | Input Order | Reference Order | Expected |
|------|----------|-------------|-----------------|----------|
| O1 | Top-level list reorder | `[B, A, C]` | `[A, B, C]` | Aligned to `[A, B, C]` |
| O2 | Nested list reorder | host A: `[y, x]`, host B: `[x, y]` | host A: `[x, y]`, host B: `[y, x]` | Each host aligned independently |
| O3 | Extra entry in current (not in reference) | `[A, B, C, D]` | `[B, A]` | `[B, A, C, D]` — ref entries first, extras after |
| O4 | Missing entry in current | `[A, C]` | `[A, B, C]` | `[A, C]` — B absent, order of A,C preserved |
| O5 | Leaf-list reorder | `<p>b</p><p>a</p>` | `<p>a</p><p>b</p>` | No diff (set semantics) |
| O6 | Empty reference | `[C, A, B]` | *(empty)* | Deterministic fallback: `[A, B, C]` |

---

## 8. Overwrite vs. Patch Semantics

The provider has two code paths for updates:

| Method | NETCONF RPC | When Used | Semantics |
|--------|-------------|-----------|-----------|
| **Patch** (`SendUpdate`) | `<edit-config>` with `nc:operation` on each leaf | Default for all updates | Surgical: only changed leaves are touched |
| **Overwrite** (`SendDirectTransaction`) | `<load-configuration action="merge">` with full XML | Fallback when patch verification fails | Full replace: entire resource config is sent |

### Test: Patch Verification Fallback (M4)

1. Apply initial config via `SendDirectTransaction` (Create)
2. Modify one leaf → `ComputeDiff` → `CreateDiffPatch` → `SendUpdate`
3. Read back → verify → diff should be empty
4. *Inject mock bug:* make readback return stale data after patch
5. Provider detects `remainingDiff > 0` → falls back to `SendDirectTransaction`
6. Read back → verify → diff empty → success

---

## 9. Checklist: Mapping to Existing Tests

| ID | Go Unit Test | Mock Test | CI Integration |
|----|-------------|-----------|----------------|
| L1 | ❌ Add | ❌ Add | ✅ (initial apply) |
| L2 | ✅ (TestCreateDiffPatch_ReplaceHostName) | ✅ (test_edit_patch_replace_leaf) | ✅ (mutate IP) |
| L3 | ❌ Add | ✅ (test_edit_patch_delete_keyed_nested_list_entry) | ❌ Add |
| L4 | ❌ Add | ❌ Add | ❌ Add |
| L5 | ❌ Add | ❌ Add | ❌ Add |
| L6 | ❌ Add (no-op) | ❌ N/A | ✅ (no-op plan) |
| LL1 | ✅ (TestLeafMapWithSchema_LeafListSetDiff) | ✅ (test_edit_patch_create_leaf_list_appends) | ❌ Add |
| LL2 | ✅ (same test) | ✅ (test_edit_patch_delete_leaf_list_value) | ❌ Add |
| LL3 | ❌ Add | ❌ Add | ❌ Add |
| LL4 | ❌ Add | ❌ N/A | ❌ Add |
| LL5 | ❌ Add | ❌ Add | ❌ Add |
| LL6 | ❌ Add | ❌ Add | ❌ Add |
| K1 | ❌ Add | ✅ (test_edit_patch_create_keyed_list_entry) | ❌ Add |
| K2 | ❌ Add | ✅ (test_edit_patch_delete_keyed_list_entry) | ❌ Add |
| K3 | ✅ (TestCreateDiffPatch_KeyedListRename) | ❌ Add | ❌ Add |
| K4 | ❌ Add | ✅ (test_edit_patch_replace_leaf_in_keyed_nested_path) | ✅ (mutate IP) |
| K5 | ❌ Add | ❌ Add | ❌ Add |
| K6 | ❌ Add | ❌ Add | ❌ Add |
| K7 | ✅ (TestAlignXMLOrderToReference_Reorders) | ❌ N/A | ❌ Add |
| K8 | ❌ Add | ❌ Add | ❌ Add |
| K9 | ❌ Add | ❌ N/A | ❌ Add |
| K10 | ✅ (TestLeafMapWithSchema_StructuralKeyed) | ❌ Add | ❌ Add |
| C1 | ❌ Add | ❌ Add | ✅ (initial apply) |
| C3 | ❌ Add | ❌ Add | ❌ Add |
| C5 | ❌ Add | ❌ Add | ✅ (mutate) |
| M1 | ❌ Add | ✅ (test_edit_patch_mixed_operations) | ❌ Add |
| M4 | ❌ Add | ❌ Add | ❌ Add |
| M5 | ❌ Add | ❌ Add | ✅ (drift inject + reconcile) |

### Coverage Summary

- **Existing coverage:** 10 of 31 scenarios have Go unit tests
- **Existing mock tests:** 10 of 31 have Python mock tests
- **CI integration:** 6 of 31 are exercised in the CI workflow
- **Gaps to fill:** 21 Go unit tests, 15 mock tests, ~18 CI scenarios

---

## 10. Quick-Start: Writing a New Test

### Go Unit Test (patch layer)

```go
// In terraform_provider/patch/patch_test.go or a new file

func TestLL1_AddLeafListEntry(t *testing.T) {
    idx := mustIdx(t)  // Uses TrimmedSchemaJSON from patch_test.go

    stateXML := `<configuration>
      <policy-options>
        <community><name>my-comm</name>
          <members>target:65000:100</members>
        </community>
      </policy-options>
    </configuration>`

    planXML := `<configuration>
      <policy-options>
        <community><name>my-comm</name>
          <members>target:65000:100</members>
          <members>target:65000:200</members>
        </community>
      </policy-options>
    </configuration>`

    stateTree := mustTree(stateXML)
    planTree := mustTree(planXML)

    stateMap := LeafMapWithSchema(stateTree, idx)
    planMap := LeafMapWithSchema(planTree, idx)

    diff := ComputeDiff(stateMap, planMap)

    // Expect exactly one Create for the new member
    if len(diff) != 1 {
        t.Fatalf("expected 1 diff entry, got %d", len(diff))
    }

    for path, change := range diff {
        if change.Op != Create {
            t.Errorf("expected Create, got %v for %s", change.Op, path)
        }
        if !strings.Contains(path, "members[value=target:65000:200]") {
            t.Errorf("unexpected path: %s", path)
        }
    }

    patchXML, err := CreateDiffPatch(diff, "")
    if err != nil {
        t.Fatal(err)
    }

    patch := string(patchXML)
    if !strings.Contains(patch, `nc:operation="create"`) {
        t.Error("patch missing create operation")
    }
    if !strings.Contains(patch, `target:65000:200`) {
        t.Error("patch missing new member value")
    }
}
```

### Python Mock Test

```python
# In netconf_mock/tests/test_netconf_mock_server.py

def test_patch_create_leaf(state_and_session):
    state, session = state_and_session
    # Seed with existing config
    state.running_groups["base-config"] = (
        "<configuration><groups><name>base-config</name>"
        "<system/>"
        "</groups></configuration>"
    )
    state.candidate_groups = copy.deepcopy(state.running_groups)

    # Send edit-config with nc:operation="create" on host-name
    rpc = (
        '<rpc message-id="1">'
        '<edit-config>'
        '<target><candidate/></target>'
        '<default-operation>none</default-operation>'
        '<config xmlns:nc="urn:ietf:params:xml:ns:netconf:base:1.0">'
        '<configuration><groups><name>base-config</name>'
        '<system>'
        '<host-name nc:operation="create">router1</host-name>'
        '</system>'
        '</groups></configuration>'
        '</config></edit-config></rpc>'
    )
    session._handle_rpc(rpc)

    # Verify candidate has the new leaf
    candidate = state.candidate_groups["base-config"]
    assert "router1" in candidate
    assert "<host-name" in candidate  # or use ET.parse for structured check
```

### Terraform Integration Test (local, no CI)

```bash
#!/bin/bash
# Run from: examples/providers/terraform-provider-junos-vmx-4-topo/

# 1. Build provider
go install .

# 2. Start mock with "before" config
python3 ../../../netconf_mock/netconf_mock_server.py \
  --host 127.0.0.1 \
  --username test --password test \
  --device test-device:8301 &
MOCK_PID=$!

# 3. Write .terraformrc
cat > ~/.terraformrc <<EOF
provider_installation {
  dev_overrides {
    "hashicorp/junos-vmx-4-topo" = "$HOME/go/bin"
  }
  direct {}
}
EOF

# 4. Apply
cd /path/to/tf/files
terraform plan -refresh=false -no-color -out=tfplan
terraform apply -no-color -auto-approve tfplan

# 5. Verify no-op
terraform plan -refresh=false -detailed-exitcode
EXIT_CODE=$?
if [ "$EXIT_CODE" -ne 0 ]; then
  echo "FAIL: expected no-op plan after apply"
  exit 1
fi

# 6. Cleanup
kill $MOCK_PID
```

---

## 11. Physical Device Testing Notes

When testing against a real Junos device instead of the mock:

1. **Backup config first:** `show configuration | display xml | save /var/tmp/backup.xml`
2. **Use a lab/test device only** — never production
3. **Check Junos version compatibility:** The provider's YANG models must match the device's Junos version
4. **NETCONF must be enabled:** `set system services netconf ssh`
5. **The `providers.tf` must point to the real device IP/port**
6. **After testing, restore:** `load override /var/tmp/backup.xml` + `commit`

### Real-device test differences from mock:
- XML namespace prefixes may vary
- Element ordering follows Junos internal ordering (not alphabetical)
- `encoding/xml` in Go strips some whitespace that Junos preserves
- Login banners with literal `\n` need trailing-newline stripping (fixed in `jtaf-xml2tf`)
- `<groups>` inheritance and `apply-groups` ordering matter
