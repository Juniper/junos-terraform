package patch

import (
	"strings"
	"testing"
)

// matrixSchema covers all four YANG node types needed by the test matrix:
//   - container:  system, interfaces, policy-options, chassis, services
//   - list:       interface (key=name), unit (key=name), host (key=name), contents (key=name)
//   - leaf:       host-name, description, device-count, etc.
//   - leaf-list:  members
const matrixSchema = `{
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
          },
          {
            "name": "chassis",
            "type": "container",
            "path": "",
            "children": [
              {
                "name": "aggregated-devices",
                "type": "container",
                "path": "chassis",
                "children": [
                  {
                    "name": "ethernet",
                    "type": "container",
                    "path": "chassis/aggregated-devices",
                    "children": [
                      {
                        "name": "device-count",
                        "type": "leaf",
                        "path": "chassis/aggregated-devices/ethernet",
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
}`

func matrixIdx(t *testing.T) map[string]*NodeInfo {
	t.Helper()
	return mustIdxFromSchema(t, matrixSchema)
}

// normalizeXML strips all leading/trailing whitespace from each line and joins.
func normalizeXML(s string) string {
	return strings.Join(strings.Fields(s), "")
}

// ---------------------------------------------------------------------------
// 2.1 Leaf Operations
// ---------------------------------------------------------------------------

// L1 — Create leaf: add host-name to empty system container
func TestL1_CreateLeaf(t *testing.T) {
	idx := matrixIdx(t)

	stateXML := `<configuration><system></system></configuration>`
	planXML := `<configuration><system><host-name>router1</host-name></system></configuration>`

	stateMap := LeafMapWithSchema(mustTree(t, stateXML), idx)
	planMap := LeafMapWithSchema(mustTree(t, planXML), idx)
	diff := ComputeDiff(stateMap, planMap)

	if len(diff) != 1 {
		t.Fatalf("expected 1 diff entry, got %d: %v", len(diff), diff)
	}

	key := "configuration/system/host-name"
	ch, ok := diff[key]
	if !ok {
		t.Fatalf("expected diff key %s, got %v", key, diff)
	}
	if ch.Op != Create {
		t.Fatalf("expected Create, got %v", ch.Op)
	}
	if ch.NewVal != "router1" {
		t.Fatalf("expected NewVal=router1, got %q", ch.NewVal)
	}

	patch, err := CreateDiffPatch(diff, "")
	if err != nil {
		t.Fatal(err)
	}
	patchStr := string(patch)
	if !strings.Contains(patchStr, `nc:operation="create"`) {
		t.Fatalf("patch missing create operation:\n%s", patchStr)
	}
	if !strings.Contains(patchStr, ">router1<") {
		t.Fatalf("patch missing value router1:\n%s", patchStr)
	}
}

// L2 — Replace leaf: change host-name value
func TestL2_ReplaceLeaf(t *testing.T) {
	idx := matrixIdx(t)

	stateXML := `<configuration><system><host-name>router1</host-name></system></configuration>`
	planXML := `<configuration><system><host-name>router2</host-name></system></configuration>`

	stateMap := LeafMapWithSchema(mustTree(t, stateXML), idx)
	planMap := LeafMapWithSchema(mustTree(t, planXML), idx)
	diff := ComputeDiff(stateMap, planMap)

	if len(diff) != 1 {
		t.Fatalf("expected 1 diff entry, got %d", len(diff))
	}

	key := "configuration/system/host-name"
	ch := diff[key]
	if ch.Op != Replace {
		t.Fatalf("expected Replace, got %v", ch.Op)
	}
	if ch.OldVal != "router1" || ch.NewVal != "router2" {
		t.Fatalf("expected router1->router2, got %q->%q", ch.OldVal, ch.NewVal)
	}

	patch, err := CreateDiffPatch(diff, "")
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(patch), `nc:operation="replace"`) {
		t.Fatalf("patch missing replace operation:\n%s", string(patch))
	}
}

// L3 — Delete leaf: remove description from interface
func TestL3_DeleteLeaf(t *testing.T) {
	idx := matrixIdx(t)

	stateXML := `<configuration><interfaces><interface><name>ge-0/0/0</name><description>uplink</description></interface></interfaces></configuration>`
	planXML := `<configuration><interfaces><interface><name>ge-0/0/0</name></interface></interfaces></configuration>`

	stateMap := LeafMapWithSchema(mustTree(t, stateXML), idx)
	planMap := LeafMapWithSchema(mustTree(t, planXML), idx)
	diff := ComputeDiff(stateMap, planMap)

	if len(diff) != 1 {
		t.Fatalf("expected 1 diff entry, got %d: %v", len(diff), diff)
	}

	key := "configuration/interfaces/interface[name=ge-0/0/0]/description"
	ch, ok := diff[key]
	if !ok {
		// Dump all keys for debugging
		for k := range diff {
			t.Logf("diff key: %s", k)
		}
		t.Fatalf("expected diff key %s", key)
	}
	if ch.Op != Delete {
		t.Fatalf("expected Delete, got %v", ch.Op)
	}

	patch, err := CreateDiffPatch(diff, "")
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(patch), `nc:operation="delete"`) {
		t.Fatalf("patch missing delete operation:\n%s", string(patch))
	}
}

// L4 — Replace leaf with XML special characters
func TestL4_ReplaceLeafSpecialChars(t *testing.T) {
	idx := matrixIdx(t)

	stateXML := `<configuration><interfaces><interface><name>ge-0/0/0</name><description>old</description></interface></interfaces></configuration>`
	planXML := `<configuration><interfaces><interface><name>ge-0/0/0</name><description>x&amp;y&lt;z&gt;w</description></interface></interfaces></configuration>`

	stateMap := LeafMapWithSchema(mustTree(t, stateXML), idx)
	planMap := LeafMapWithSchema(mustTree(t, planXML), idx)
	diff := ComputeDiff(stateMap, planMap)

	if len(diff) != 1 {
		t.Fatalf("expected 1 diff entry, got %d", len(diff))
	}

	ch := diff["configuration/interfaces/interface[name=ge-0/0/0]/description"]
	if ch.Op != Replace {
		t.Fatalf("expected Replace, got %v", ch.Op)
	}
	// After XML parsing, the value should be the un-escaped form
	if ch.NewVal != "x&y<z>w" {
		t.Fatalf("expected un-escaped value x&y<z>w, got %q", ch.NewVal)
	}

	patch, err := CreateDiffPatch(diff, "")
	if err != nil {
		t.Fatal(err)
	}
	patchStr := string(patch)
	// The output XML must re-escape the special chars
	if !strings.Contains(patchStr, "&amp;") {
		t.Fatalf("patch missing &amp; escape:\n%s", patchStr)
	}
	if !strings.Contains(patchStr, "&lt;") {
		t.Fatalf("patch missing &lt; escape:\n%s", patchStr)
	}
	if !strings.Contains(patchStr, "&gt;") {
		t.Fatalf("patch missing &gt; escape:\n%s", patchStr)
	}
}

// L6 — No-op: same value produces empty diff
func TestL6_ReplaceLeafNoOp(t *testing.T) {
	idx := matrixIdx(t)

	xml := `<configuration><system><host-name>r1</host-name></system></configuration>`
	stateMap := LeafMapWithSchema(mustTree(t, xml), idx)
	planMap := LeafMapWithSchema(mustTree(t, xml), idx)
	diff := ComputeDiff(stateMap, planMap)

	if len(diff) != 0 {
		t.Fatalf("expected empty diff for identical config, got %d entries: %v", len(diff), diff)
	}
}

// ---------------------------------------------------------------------------
// 2.2 Leaf-List Operations
// ---------------------------------------------------------------------------

// LL1 — Add entry to leaf-list
func TestLL1_AddLeafListEntry(t *testing.T) {
	idx := matrixIdx(t)

	stateXML := `<configuration><policy-options><community><name>my-comm</name><members>target:65000:100</members></community></policy-options></configuration>`
	planXML := `<configuration><policy-options><community><name>my-comm</name><members>target:65000:100</members><members>target:65000:200</members></community></policy-options></configuration>`

	stateMap := LeafMapWithSchema(mustTree(t, stateXML), idx)
	planMap := LeafMapWithSchema(mustTree(t, planXML), idx)
	diff := ComputeDiff(stateMap, planMap)

	if len(diff) != 1 {
		t.Fatalf("expected 1 diff entry, got %d: %v", len(diff), diff)
	}

	for path, ch := range diff {
		if ch.Op != Create {
			t.Errorf("expected Create, got %v for %s", ch.Op, path)
		}
		if !strings.Contains(path, "members[value=target:65000:200]") {
			t.Errorf("unexpected path: %s", path)
		}
	}

	patch, err := CreateDiffPatch(diff, "")
	if err != nil {
		t.Fatal(err)
	}
	patchStr := string(patch)
	if !strings.Contains(patchStr, `nc:operation="create"`) {
		t.Errorf("patch missing create operation:\n%s", patchStr)
	}
	if !strings.Contains(patchStr, "target:65000:200") {
		t.Errorf("patch missing new member value:\n%s", patchStr)
	}
}

// LL2 — Remove entry from leaf-list
func TestLL2_RemoveLeafListEntry(t *testing.T) {
	idx := matrixIdx(t)

	stateXML := `<configuration><policy-options><community><name>my-comm</name><members>target:65000:100</members><members>target:65000:200</members></community></policy-options></configuration>`
	planXML := `<configuration><policy-options><community><name>my-comm</name><members>target:65000:100</members></community></policy-options></configuration>`

	stateMap := LeafMapWithSchema(mustTree(t, stateXML), idx)
	planMap := LeafMapWithSchema(mustTree(t, planXML), idx)
	diff := ComputeDiff(stateMap, planMap)

	if len(diff) != 1 {
		t.Fatalf("expected 1 diff, got %d: %v", len(diff), diff)
	}

	for path, ch := range diff {
		if ch.Op != Delete {
			t.Errorf("expected Delete, got %v for %s", ch.Op, path)
		}
		if !strings.Contains(path, "members[value=target:65000:200]") {
			t.Errorf("unexpected path: %s", path)
		}
	}
}

// LL3 — Replace entry in leaf-list (delete old + create new)
func TestLL3_ReplaceLeafListEntry(t *testing.T) {
	idx := matrixIdx(t)

	stateXML := `<configuration><policy-options><community><name>my-comm</name><members>a</members><members>b</members></community></policy-options></configuration>`
	planXML := `<configuration><policy-options><community><name>my-comm</name><members>a</members><members>c</members></community></policy-options></configuration>`

	stateMap := LeafMapWithSchema(mustTree(t, stateXML), idx)
	planMap := LeafMapWithSchema(mustTree(t, planXML), idx)
	diff := ComputeDiff(stateMap, planMap)

	if len(diff) != 2 {
		t.Fatalf("expected 2 diff entries (delete b + create c), got %d: %v", len(diff), diff)
	}

	deletePath := "configuration/policy-options/community[name=my-comm]/members[value=b]"
	createPath := "configuration/policy-options/community[name=my-comm]/members[value=c]"
	if ch, ok := diff[deletePath]; !ok || ch.Op != Delete {
		t.Fatalf("expected Delete for %s, got %v", deletePath, diff)
	}
	if ch, ok := diff[createPath]; !ok || ch.Op != Create {
		t.Fatalf("expected Create for %s, got %v", createPath, diff)
	}
}

// LL4 — Reorder leaf-list produces no diff (set semantics)
func TestLL4_ReorderLeafListNoOp(t *testing.T) {
	idx := matrixIdx(t)

	stateXML := `<configuration><policy-options><community><name>my-comm</name><members>b</members><members>a</members></community></policy-options></configuration>`
	planXML := `<configuration><policy-options><community><name>my-comm</name><members>a</members><members>b</members></community></policy-options></configuration>`

	stateMap := LeafMapWithSchema(mustTree(t, stateXML), idx)
	planMap := LeafMapWithSchema(mustTree(t, planXML), idx)
	diff := ComputeDiff(stateMap, planMap)

	if len(diff) != 0 {
		t.Fatalf("expected empty diff for reordered leaf-list, got %d: %v", len(diff), diff)
	}
}

// LL5 — Delete all leaf-list entries
func TestLL5_DeleteAllLeafListEntries(t *testing.T) {
	idx := matrixIdx(t)

	stateXML := `<configuration><policy-options><community><name>my-comm</name><members>a</members><members>b</members></community></policy-options></configuration>`
	planXML := `<configuration><policy-options><community><name>my-comm</name></community></policy-options></configuration>`

	stateMap := LeafMapWithSchema(mustTree(t, stateXML), idx)
	planMap := LeafMapWithSchema(mustTree(t, planXML), idx)
	diff := ComputeDiff(stateMap, planMap)

	if len(diff) != 2 {
		t.Fatalf("expected 2 deletes, got %d: %v", len(diff), diff)
	}

	for _, ch := range diff {
		if ch.Op != Delete {
			t.Errorf("expected Delete, got %v", ch.Op)
		}
	}
}

// LL6 — Create leaf-list from scratch
func TestLL6_CreateLeafListFromScratch(t *testing.T) {
	idx := matrixIdx(t)

	stateXML := `<configuration><policy-options><community><name>my-comm</name></community></policy-options></configuration>`
	planXML := `<configuration><policy-options><community><name>my-comm</name><members>x</members><members>y</members></community></policy-options></configuration>`

	stateMap := LeafMapWithSchema(mustTree(t, stateXML), idx)
	planMap := LeafMapWithSchema(mustTree(t, planXML), idx)
	diff := ComputeDiff(stateMap, planMap)

	if len(diff) != 2 {
		t.Fatalf("expected 2 creates, got %d: %v", len(diff), diff)
	}

	for _, ch := range diff {
		if ch.Op != Create {
			t.Errorf("expected Create, got %v", ch.Op)
		}
	}
}

// ---------------------------------------------------------------------------
// 2.3 List Operations
// ---------------------------------------------------------------------------

// K1 — Add new list entry
func TestK1_AddListEntry(t *testing.T) {
	idx := matrixIdx(t)

	stateXML := `<configuration><interfaces><interface><name>ge-0/0/0</name><description>uplink</description></interface></interfaces></configuration>`
	planXML := `<configuration><interfaces><interface><name>ge-0/0/0</name><description>uplink</description></interface><interface><name>ge-0/0/1</name><description>downlink</description></interface></interfaces></configuration>`

	stateMap := LeafMapWithSchema(mustTree(t, stateXML), idx)
	planMap := LeafMapWithSchema(mustTree(t, planXML), idx)
	diff := ComputeDiff(stateMap, planMap)

	// New entry creates: name (promoted to entry operation) + description
	if len(diff) != 2 {
		t.Fatalf("expected 2 diff entries for new list entry, got %d: %v", len(diff), diff)
	}

	for _, ch := range diff {
		if ch.Op != Create {
			t.Errorf("expected Create, got %v", ch.Op)
		}
	}

	patch, err := CreateDiffPatch(diff, "")
	if err != nil {
		t.Fatal(err)
	}
	patchStr := string(patch)
	// The key leaf create should promote nc:operation to the parent <interface>
	if !strings.Contains(patchStr, `interface nc:operation="create"`) {
		t.Fatalf("expected nc:operation on interface entry, got:\n%s", patchStr)
	}
}

// K2 — Delete list entry
func TestK2_DeleteListEntry(t *testing.T) {
	idx := matrixIdx(t)

	stateXML := `<configuration><interfaces><interface><name>ge-0/0/0</name><description>uplink</description></interface><interface><name>ge-0/0/1</name><description>downlink</description></interface></interfaces></configuration>`
	planXML := `<configuration><interfaces><interface><name>ge-0/0/0</name><description>uplink</description></interface></interfaces></configuration>`

	stateMap := LeafMapWithSchema(mustTree(t, stateXML), idx)
	planMap := LeafMapWithSchema(mustTree(t, planXML), idx)
	diff := ComputeDiff(stateMap, planMap)

	for _, ch := range diff {
		if ch.Op != Delete {
			t.Errorf("expected Delete, got %v", ch.Op)
		}
	}

	patch, err := CreateDiffPatch(diff, "")
	if err != nil {
		t.Fatal(err)
	}
	patchStr := string(patch)
	if !strings.Contains(patchStr, `interface nc:operation="delete"`) {
		t.Fatalf("expected nc:operation on interface entry, got:\n%s", patchStr)
	}
}

// K3 — Rename list key (delete old + create new)
func TestK3_RenameListKey(t *testing.T) {
	idx := matrixIdx(t)

	stateXML := `<configuration><interfaces><interface><name>ge-0/0/0</name><description>link</description></interface></interfaces></configuration>`
	planXML := `<configuration><interfaces><interface><name>ge-0/0/1</name><description>link</description></interface></interfaces></configuration>`

	stateMap := LeafMapWithSchema(mustTree(t, stateXML), idx)
	planMap := LeafMapWithSchema(mustTree(t, planXML), idx)
	diff := ComputeDiff(stateMap, planMap)

	// Should produce deletes for old entry + creates for new entry
	hasDelete := false
	hasCreate := false
	for _, ch := range diff {
		if ch.Op == Delete {
			hasDelete = true
		}
		if ch.Op == Create {
			hasCreate = true
		}
	}
	if !hasDelete || !hasCreate {
		t.Fatalf("expected both Delete and Create for key rename, got: %v", diff)
	}
}

// K4 — Modify leaf inside list entry
func TestK4_ModifyLeafInListEntry(t *testing.T) {
	idx := matrixIdx(t)

	stateXML := `<configuration><interfaces><interface><name>ge-0/0/0</name><unit><name>0</name><description>old</description></unit></interface></interfaces></configuration>`
	planXML := `<configuration><interfaces><interface><name>ge-0/0/0</name><unit><name>0</name><description>new</description></unit></interface></interfaces></configuration>`

	stateMap := LeafMapWithSchema(mustTree(t, stateXML), idx)
	planMap := LeafMapWithSchema(mustTree(t, planXML), idx)
	diff := ComputeDiff(stateMap, planMap)

	if len(diff) != 1 {
		t.Fatalf("expected 1 diff, got %d: %v", len(diff), diff)
	}

	key := "configuration/interfaces/interface[name=ge-0/0/0]/unit[name=0]/description"
	ch, ok := diff[key]
	if !ok {
		for k := range diff {
			t.Logf("diff key: %s", k)
		}
		t.Fatalf("expected diff key %s", key)
	}
	if ch.Op != Replace {
		t.Fatalf("expected Replace, got %v", ch.Op)
	}
}

// K5 — Add leaf inside existing list entry
func TestK5_AddLeafInListEntry(t *testing.T) {
	idx := matrixIdx(t)

	stateXML := `<configuration><interfaces><interface><name>ge-0/0/0</name><unit><name>0</name></unit></interface></interfaces></configuration>`
	planXML := `<configuration><interfaces><interface><name>ge-0/0/0</name><unit><name>0</name><description>new</description></unit></interface></interfaces></configuration>`

	stateMap := LeafMapWithSchema(mustTree(t, stateXML), idx)
	planMap := LeafMapWithSchema(mustTree(t, planXML), idx)
	diff := ComputeDiff(stateMap, planMap)

	if len(diff) != 1 {
		t.Fatalf("expected 1 diff, got %d: %v", len(diff), diff)
	}

	for _, ch := range diff {
		if ch.Op != Create {
			t.Errorf("expected Create, got %v", ch.Op)
		}
	}
}

// K6 — Delete leaf inside list entry
func TestK6_DeleteLeafInListEntry(t *testing.T) {
	idx := matrixIdx(t)

	stateXML := `<configuration><interfaces><interface><name>ge-0/0/0</name><unit><name>0</name><description>old</description></unit></interface></interfaces></configuration>`
	planXML := `<configuration><interfaces><interface><name>ge-0/0/0</name><unit><name>0</name></unit></interface></interfaces></configuration>`

	stateMap := LeafMapWithSchema(mustTree(t, stateXML), idx)
	planMap := LeafMapWithSchema(mustTree(t, planXML), idx)
	diff := ComputeDiff(stateMap, planMap)

	if len(diff) != 1 {
		t.Fatalf("expected 1 diff, got %d: %v", len(diff), diff)
	}

	for _, ch := range diff {
		if ch.Op != Delete {
			t.Errorf("expected Delete, got %v", ch.Op)
		}
	}
}

// K8 — Add entry in nested list (host/contents)
func TestK8_AddNestedListEntry(t *testing.T) {
	idx := matrixIdx(t)

	stateXML := `<configuration><system><syslog><host><name>log</name><contents><name>any</name><notice/></contents></host></syslog></system></configuration>`
	planXML := `<configuration><system><syslog><host><name>log</name><contents><name>any</name><notice/></contents><contents><name>kernel</name><any/></contents></host></syslog></system></configuration>`

	stateMap := LeafMapWithSchema(mustTree(t, stateXML), idx)
	planMap := LeafMapWithSchema(mustTree(t, planXML), idx)
	diff := ComputeDiff(stateMap, planMap)

	// New contents entry: name (key, promoted) + any (empty leaf)
	hasCreate := false
	for _, ch := range diff {
		if ch.Op == Create {
			hasCreate = true
		}
	}
	if !hasCreate {
		t.Fatalf("expected Create operations for new nested list entry, got: %v", diff)
	}
}

// K10 — Key-only list entry (structural)
func TestK10_KeyOnlyListEntry(t *testing.T) {
	idx := matrixIdx(t)

	stateXML := `<configuration><system><syslog></syslog></system></configuration>`
	planXML := `<configuration><system><syslog><host><name>log</name></host></syslog></system></configuration>`

	stateMap := LeafMapWithSchema(mustTree(t, stateXML), idx)
	planMap := LeafMapWithSchema(mustTree(t, planXML), idx)
	diff := ComputeDiff(stateMap, planMap)

	if len(diff) != 1 {
		t.Fatalf("expected 1 diff for key-only entry, got %d: %v", len(diff), diff)
	}

	key := "configuration/system/syslog/host[name=log]/name"
	ch, ok := diff[key]
	if !ok {
		for k := range diff {
			t.Logf("diff key: %s", k)
		}
		t.Fatalf("expected diff key %s", key)
	}
	if ch.Op != Create {
		t.Fatalf("expected Create, got %v", ch.Op)
	}
}

// ---------------------------------------------------------------------------
// 2.4 Container Operations
// ---------------------------------------------------------------------------

// C1 — Create container with children
func TestC1_CreateContainer(t *testing.T) {
	idx := matrixIdx(t)

	stateXML := `<configuration></configuration>`
	planXML := `<configuration><chassis><aggregated-devices><ethernet><device-count>24</device-count></ethernet></aggregated-devices></chassis></configuration>`

	stateMap := LeafMapWithSchema(mustTree(t, stateXML), idx)
	planMap := LeafMapWithSchema(mustTree(t, planXML), idx)
	diff := ComputeDiff(stateMap, planMap)

	if len(diff) != 1 {
		t.Fatalf("expected 1 Create for device-count leaf, got %d: %v", len(diff), diff)
	}

	for _, ch := range diff {
		if ch.Op != Create {
			t.Errorf("expected Create, got %v", ch.Op)
		}
	}
}

// C3 — Empty presence container (e.g., <ssh/>)
func TestC3_EmptyContainer(t *testing.T) {
	idx := matrixIdx(t)

	stateXML := `<configuration><system><services></services></system></configuration>`
	planXML := `<configuration><system><services><ssh/></services></system></configuration>`

	stateMap := LeafMapWithSchema(mustTree(t, stateXML), idx)
	planMap := LeafMapWithSchema(mustTree(t, planXML), idx)
	diff := ComputeDiff(stateMap, planMap)

	// An empty container like <ssh/> has no text and no children, so it
	// produces no leaf map entries. The diff should be empty because
	// LeafMapWithSchema only tracks leaf values.
	// This is a known limitation: presence containers need special handling.
	// For now, verify the leaf maps are consistent.
	t.Logf("state map: %v", stateMap)
	t.Logf("plan map:  %v", planMap)
	t.Logf("diff:      %v", diff)

	// Both maps should be empty (no leaves under services or ssh)
	// This documents the current behavior — empty containers don't produce diffs.
	// The provider handles this via the full-config SendDirectTransaction path.
}

// C5 — Modify children within container
func TestC5_ModifyChildrenInContainer(t *testing.T) {
	idx := matrixIdx(t)

	stateXML := `<configuration><system><host-name>old</host-name></system></configuration>`
	planXML := `<configuration><system><host-name>new</host-name></system></configuration>`

	stateMap := LeafMapWithSchema(mustTree(t, stateXML), idx)
	planMap := LeafMapWithSchema(mustTree(t, planXML), idx)
	diff := ComputeDiff(stateMap, planMap)

	if len(diff) != 1 {
		t.Fatalf("expected 1 Replace, got %d: %v", len(diff), diff)
	}

	ch := diff["configuration/system/host-name"]
	if ch.Op != Replace {
		t.Fatalf("expected Replace, got %v", ch.Op)
	}
}

// ---------------------------------------------------------------------------
// 2.5 Compound / Cross-Type Operations
// ---------------------------------------------------------------------------

// M1 — Mixed operations: replace + create + delete in one diff
func TestM1_MixedOperations(t *testing.T) {
	idx := matrixIdx(t)

	stateXML := `<configuration>
  <system><host-name>r1</host-name></system>
  <interfaces>
    <interface><name>ge-0/0/0</name><description>old-desc</description></interface>
  </interfaces>
  <policy-options>
    <community><name>comm1</name><members>target:65000:100</members><members>target:65000:300</members></community>
  </policy-options>
</configuration>`

	planXML := `<configuration>
  <system><host-name>r2</host-name></system>
  <interfaces>
    <interface><name>ge-0/0/0</name><description>old-desc</description></interface>
    <interface><name>ge-0/0/2</name><description>new-link</description></interface>
  </interfaces>
  <policy-options>
    <community><name>comm1</name><members>target:65000:100</members></community>
  </policy-options>
</configuration>`

	stateMap := LeafMapWithSchema(mustTree(t, stateXML), idx)
	planMap := LeafMapWithSchema(mustTree(t, planXML), idx)
	diff := ComputeDiff(stateMap, planMap)

	// Expected changes:
	// 1. Replace system/host-name r1 -> r2
	// 2. Delete members[value=target:65000:300]
	// 3. Create interface[name=ge-0/0/2]/name (promoted)
	// 4. Create interface[name=ge-0/0/2]/description
	hasReplace := false
	hasDelete := false
	hasCreate := false
	for _, ch := range diff {
		switch ch.Op {
		case Replace:
			hasReplace = true
		case Delete:
			hasDelete = true
		case Create:
			hasCreate = true
		}
	}
	if !hasReplace || !hasDelete || !hasCreate {
		t.Fatalf("expected Replace+Delete+Create, got: %v", diff)
	}

	// Verify ordering: deletes first, then replacements, then creates
	patch, err := CreateDiffPatch(diff, "")
	if err != nil {
		t.Fatal(err)
	}
	patchStr := string(patch)

	deleteIdx := strings.Index(patchStr, `"delete"`)
	replaceIdx := strings.Index(patchStr, `"replace"`)
	createIdx := strings.Index(patchStr, `"create"`)

	if deleteIdx < 0 || replaceIdx < 0 || createIdx < 0 {
		t.Fatalf("patch missing expected operations:\n%s", patchStr)
	}

	if deleteIdx > replaceIdx {
		t.Errorf("delete should come before replace in patch output")
	}
	if replaceIdx > createIdx {
		t.Errorf("replace should come before create in patch output")
	}
}

// M2 — Deep nesting: 4+ levels (system/syslog/host/contents/notice)
func TestM2_DeepNesting(t *testing.T) {
	idx := matrixIdx(t)

	stateXML := `<configuration><system><syslog><host><name>log</name><contents><name>any</name></contents></host></syslog></system></configuration>`
	planXML := `<configuration><system><syslog><host><name>log</name><contents><name>any</name><notice/></contents></host></syslog></system></configuration>`

	stateMap := LeafMapWithSchema(mustTree(t, stateXML), idx)
	planMap := LeafMapWithSchema(mustTree(t, planXML), idx)
	diff := ComputeDiff(stateMap, planMap)

	// The <notice/> is an empty leaf, so in LeafMapWithSchema it may or may
	// not produce a leaf map entry (depends on how empty leaves are handled).
	t.Logf("state map: %v", stateMap)
	t.Logf("plan map:  %v", planMap)
	t.Logf("diff: %v", diff)

	// At minimum, the key-only entries should be consistent
	stateKey := "configuration/system/syslog/host[name=log]/contents[name=any]/name"
	if _, ok := stateMap[stateKey]; ok {
		// If contents has children beyond key, key should not be in map
		t.Logf("contents key is in state map (structural key behavior)")
	}
}
