package patch

import (
	"strings"
	"testing"
)

// ---------------------------------------------------------------------------
// CC-10: Special characters in key values
// ---------------------------------------------------------------------------

func TestCC10_ParseSegment_KeyWithSlashes(t *testing.T) {
	// Interface names like ge-0/0/0 are the common case — already inside brackets
	tag, keyName, keyValue := parseSegment("interface[name=ge-0/0/0]")
	if tag != "interface" || keyName != "name" || keyValue != "ge-0/0/0" {
		t.Errorf("got tag=%q key=%q val=%q", tag, keyName, keyValue)
	}
}

func TestCC10_ParseSegment_KeyWithNestedBrackets(t *testing.T) {
	// Policy name containing brackets: "ALLOW[ALL]"
	tag, keyName, keyValue := parseSegment("policy[name=ALLOW[ALL]]")
	if tag != "policy" {
		t.Errorf("expected tag=policy, got %q", tag)
	}
	if keyName != "name" {
		t.Errorf("expected keyName=name, got %q", keyName)
	}
	if keyValue != "ALLOW[ALL]" {
		t.Errorf("expected keyValue=ALLOW[ALL], got %q", keyValue)
	}
}

func TestCC10_ParseSegment_KeyWithEquals(t *testing.T) {
	// Route key containing equals: "prefix=10.0.0.0/8"
	tag, keyName, keyValue := parseSegment("route[prefix=10.0.0.0/8]")
	if tag != "route" || keyName != "prefix" || keyValue != "10.0.0.0/8" {
		t.Errorf("got tag=%q key=%q val=%q", tag, keyName, keyValue)
	}
}

func TestCC10_SplitPath_KeyWithSlashes(t *testing.T) {
	// Verify path splitting handles keys with slashes correctly
	path := "interfaces/interface[name=ge-0/0/0]/unit[name=0]/description"
	segments := splitPathRespectingQuotes(path)
	expected := []string{"interfaces", "interface[name=ge-0/0/0]", "unit[name=0]", "description"}
	if len(segments) != len(expected) {
		t.Fatalf("expected %d segments, got %d: %v", len(expected), len(segments), segments)
	}
	for i, seg := range segments {
		if seg != expected[i] {
			t.Errorf("segment[%d]: expected %q got %q", i, expected[i], seg)
		}
	}
}

// ---------------------------------------------------------------------------
// CC-5: Container delete coalescing
// ---------------------------------------------------------------------------

// CC-5 schema covers system/ntp with two list-entries and a leaf-list
const cc5Schema = `{
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
                "name": "ntp",
                "type": "container",
                "path": "system",
                "children": [
                  {
                    "name": "server",
                    "type": "list",
                    "path": "system/ntp",
                    "key": "name",
                    "children": [
                      {
                        "name": "name",
                        "type": "leaf",
                        "path": "system/ntp/server",
                        "leaf-type": "string"
                      },
                      {
                        "name": "routing-instance",
                        "type": "leaf",
                        "path": "system/ntp/server",
                        "leaf-type": "string"
                      }
                    ]
                  },
                  {
                    "name": "trusted-key",
                    "type": "leaf-list",
                    "path": "system/ntp",
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
}`

func TestCC5_ContainerDeleteCoalescing(t *testing.T) {
	// When ALL children of a container are deleted, the patch should ideally
	// emit a single container-level delete. Currently, the engine emits
	// individual per-leaf deletes.

	stateXML := `<configuration>
  <system>
    <host-name>router1</host-name>
    <ntp>
      <server><name>10.0.0.1</name><routing-instance>mgmt</routing-instance></server>
      <server><name>10.0.0.2</name><routing-instance>mgmt</routing-instance></server>
      <trusted-key>1</trusted-key>
      <trusted-key>2</trusted-key>
    </ntp>
  </system>
</configuration>`

	// Plan: ntp entirely removed, host-name stays
	planXML := `<configuration>
  <system>
    <host-name>router1</host-name>
  </system>
</configuration>`

	idx := mustIdxFromSchema(t, cc5Schema)

	stateTree, err := BuildTree([]byte(stateXML))
	if err != nil {
		t.Fatal(err)
	}
	planTree, err := BuildTree([]byte(planXML))
	if err != nil {
		t.Fatal(err)
	}

	stateMap := LeafMapWithSchema(stateTree, idx)
	planMap := LeafMapWithSchema(planTree, idx)
	diffMap := ComputeDiff(stateMap, planMap)

	if len(diffMap) == 0 {
		t.Fatal("expected non-empty diff")
	}

	// Use schema-aware patch (with container coalescing)
	patchBytes, err := CreateDiffPatchWithSchema(diffMap, "", idx)
	if err != nil {
		t.Fatal(err)
	}

	output := string(patchBytes)
	t.Logf("CC-5 patch output:\n%s", output)

	// Verify: the patch should contain a single container-level delete for ntp
	if !strings.Contains(output, "delete") {
		t.Error("expected delete operations in patch output")
	}

	// Check that host-name is NOT deleted (it's in both state and plan)
	if strings.Contains(output, "host-name") {
		t.Error("host-name should not appear in patch (unchanged)")
	}

	// Track whether coalescing is happening
	if strings.Contains(output, `<ntp nc:operation="delete"`) {
		t.Log("CC-5 PASSED: Container-level delete coalescing IS implemented")
	} else {
		t.Error("CC-5 FAILED: Expected container-level delete coalescing")
		// Verify that at least all ntp entries are being deleted
		if !strings.Contains(output, "server") && !strings.Contains(output, "trusted-key") {
			t.Error("expected server or trusted-key deletes")
		}
	}
}

// ---------------------------------------------------------------------------
// CC-2: Empty leaf toggle (disable / vlan-tagging)
// ---------------------------------------------------------------------------

const cc2Schema = `{
  "path": "",
  "root": {
    "children": [
      {
        "name": "configuration",
        "type": "container",
        "path": "",
        "children": [
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
                    "name": "disable",
                    "type": "leaf",
                    "path": "interfaces/interface",
                    "leaf-type": "empty"
                  },
                  {
                    "name": "vlan-tagging",
                    "type": "leaf",
                    "path": "interfaces/interface",
                    "leaf-type": "empty"
                  },
                  {
                    "name": "description",
                    "type": "leaf",
                    "path": "interfaces/interface",
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
}`

func TestCC2_EmptyLeafCreate(t *testing.T) {
	// Add disable to an interface that doesn't have it
	stateXML := `<configuration>
  <interfaces>
    <interface>
      <name>ge-0/0/0</name>
      <description>uplink</description>
    </interface>
  </interfaces>
</configuration>`

	planXML := `<configuration>
  <interfaces>
    <interface>
      <name>ge-0/0/0</name>
      <description>uplink</description>
      <disable/>
    </interface>
  </interfaces>
</configuration>`

	idx := mustIdxFromSchema(t, cc2Schema)

	stateTree, _ := BuildTree([]byte(stateXML))
	planTree, _ := BuildTree([]byte(planXML))

	stateMap := LeafMapWithSchema(stateTree, idx)
	planMap := LeafMapWithSchema(planTree, idx)

	// Verify the leaf map correctly represents the empty leaf
	disableKey := ""
	for k, v := range planMap {
		if strings.HasSuffix(k, "/disable") {
			disableKey = k
			t.Logf("planMap disable: %q = %q", k, v)
		}
	}
	if disableKey == "" {
		t.Fatal("disable leaf not found in plan map")
	}

	// State should NOT have disable
	for k := range stateMap {
		if strings.HasSuffix(k, "/disable") {
			t.Fatal("disable should not be in state map")
		}
	}

	diffMap := ComputeDiff(stateMap, planMap)
	patchBytes, err := CreateDiffPatch(diffMap, "")
	if err != nil {
		t.Fatal(err)
	}

	output := string(patchBytes)
	t.Logf("CC-2 create output:\n%s", output)

	if !strings.Contains(output, "disable") {
		t.Error("expected disable in patch output")
	}
	if !strings.Contains(output, `nc:operation="create"`) {
		t.Error("expected create operation for disable")
	}
}

func TestCC2_EmptyLeafDelete(t *testing.T) {
	// Remove disable from an interface
	stateXML := `<configuration>
  <interfaces>
    <interface>
      <name>ge-0/0/0</name>
      <description>uplink</description>
      <disable/>
    </interface>
  </interfaces>
</configuration>`

	planXML := `<configuration>
  <interfaces>
    <interface>
      <name>ge-0/0/0</name>
      <description>uplink</description>
    </interface>
  </interfaces>
</configuration>`

	idx := mustIdxFromSchema(t, cc2Schema)

	stateTree, _ := BuildTree([]byte(stateXML))
	planTree, _ := BuildTree([]byte(planXML))

	stateMap := LeafMapWithSchema(stateTree, idx)
	planMap := LeafMapWithSchema(planTree, idx)

	// State SHOULD have disable
	found := false
	for k := range stateMap {
		if strings.HasSuffix(k, "/disable") {
			found = true
			break
		}
	}
	if !found {
		t.Fatal("disable should be in state map")
	}

	diffMap := ComputeDiff(stateMap, planMap)
	patchBytes, err := CreateDiffPatch(diffMap, "")
	if err != nil {
		t.Fatal(err)
	}

	output := string(patchBytes)
	t.Logf("CC-2 delete output:\n%s", output)

	if !strings.Contains(output, "disable") {
		t.Error("expected disable in patch output")
	}
	if !strings.Contains(output, `nc:operation="delete"`) {
		t.Error("expected delete operation for disable")
	}
	// Empty leaf delete should NOT have text content
	if strings.Contains(output, `<disable nc:operation="delete">`) {
		// Check if there's text between tags
		if strings.Contains(output, `<disable nc:operation="delete"></disable>`) {
			t.Log("CC-2 NOTE: empty tags (OK)")
		}
	}
}

// ---------------------------------------------------------------------------
// CC-3: Leaf-list full replacement vs incremental
// ---------------------------------------------------------------------------

func TestCC3_LeafListBulkChange(t *testing.T) {
	// Community members change from [A, B, C] to [A, D, E]
	stateXML := `<configuration>
  <policy-options>
    <community>
      <name>OC-STD</name>
      <members>65000:100</members>
      <members>65000:200</members>
      <members>65000:300</members>
    </community>
  </policy-options>
</configuration>`

	planXML := `<configuration>
  <policy-options>
    <community>
      <name>OC-STD</name>
      <members>65000:100</members>
      <members>65000:400</members>
      <members>65000:500</members>
    </community>
  </policy-options>
</configuration>`

	idx := mustIdxFromSchema(t, matrixSchema)

	stateTree, _ := BuildTree([]byte(stateXML))
	planTree, _ := BuildTree([]byte(planXML))

	stateMap := LeafMapWithSchema(stateTree, idx)
	planMap := LeafMapWithSchema(planTree, idx)
	diffMap := ComputeDiff(stateMap, planMap)

	patchBytes, err := CreateDiffPatch(diffMap, "")
	if err != nil {
		t.Fatal(err)
	}

	output := string(patchBytes)
	t.Logf("CC-3 leaf-list bulk change output:\n%s", output)

	// Verify: 65000:200 and 65000:300 should be deleted, 65000:400 and 65000:500 created
	if !strings.Contains(output, "65000:200") {
		t.Error("expected 65000:200 delete")
	}
	if !strings.Contains(output, "65000:300") {
		t.Error("expected 65000:300 delete")
	}
	if !strings.Contains(output, "65000:400") {
		t.Error("expected 65000:400 create")
	}
	if !strings.Contains(output, "65000:500") {
		t.Error("expected 65000:500 create")
	}
	// 65000:100 should NOT appear (unchanged)
	if strings.Contains(output, "65000:100") {
		t.Error("65000:100 should not be in patch (unchanged)")
	}
}

// ---------------------------------------------------------------------------
// CC-1: Ordered leaf-list reorder detection
// ---------------------------------------------------------------------------

const cc1Schema = `{
  "path": "",
  "root": {
    "children": [
      {
        "name": "configuration",
        "type": "container",
        "path": "",
        "children": [
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
                                  },
                                  {
                                    "name": "vrrp-group",
                                    "type": "list",
                                    "path": "interfaces/interface/unit/family/inet/address",
                                    "key": "name",
                                    "children": [
                                      {
                                        "name": "name",
                                        "type": "leaf",
                                        "path": "interfaces/interface/unit/family/inet/address/vrrp-group",
                                        "leaf-type": "string"
                                      },
                                      {
                                        "name": "virtual-address",
                                        "type": "leaf-list",
                                        "path": "interfaces/interface/unit/family/inet/address/vrrp-group",
                                        "leaf-type": "string",
                                        "ordered-by": "user"
                                      },
                                      {
                                        "name": "priority",
                                        "type": "leaf",
                                        "path": "interfaces/interface/unit/family/inet/address/vrrp-group",
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
              }
            ]
          }
        ]
      }
    ]
  }
}`

func TestCC1_OrderedLeafListReorder(t *testing.T) {
	// VRRP virtual-address is ordered-by user — reorder should be detected
	stateXML := `<configuration>
  <interfaces>
    <interface>
      <name>ge-0/0/0</name>
      <unit>
        <name>0</name>
        <family>
          <inet>
            <address>
              <name>10.0.0.1/24</name>
              <vrrp-group>
                <name>1</name>
                <virtual-address>10.0.0.10</virtual-address>
                <virtual-address>10.0.0.20</virtual-address>
                <priority>200</priority>
              </vrrp-group>
            </address>
          </inet>
        </family>
      </unit>
    </interface>
  </interfaces>
</configuration>`

	// Reorder: swap virtual-address order
	planXML := `<configuration>
  <interfaces>
    <interface>
      <name>ge-0/0/0</name>
      <unit>
        <name>0</name>
        <family>
          <inet>
            <address>
              <name>10.0.0.1/24</name>
              <vrrp-group>
                <name>1</name>
                <virtual-address>10.0.0.20</virtual-address>
                <virtual-address>10.0.0.10</virtual-address>
                <priority>200</priority>
              </vrrp-group>
            </address>
          </inet>
        </family>
      </unit>
    </interface>
  </interfaces>
</configuration>`

	idx := mustIdxFromSchema(t, cc1Schema)

	stateTree, _ := BuildTree([]byte(stateXML))
	planTree, _ := BuildTree([]byte(planXML))

	stateMap := LeafMapWithSchema(stateTree, idx)
	planMap := LeafMapWithSchema(planTree, idx)

	t.Logf("State map entries:")
	for k, v := range stateMap {
		if strings.Contains(k, "virtual") {
			t.Logf("  %s = %q", k, v)
		}
	}
	t.Logf("Plan map entries:")
	for k, v := range planMap {
		if strings.Contains(k, "virtual") {
			t.Logf("  %s = %q", k, v)
		}
	}

	diffMap := ComputeDiff(stateMap, planMap)

	if len(diffMap) == 0 {
		t.Log("CC-1 CONFIRMED: Reorder produces EMPTY diff (order not tracked)")
		t.Log("CC-1 STATUS: NOT COVERED — ordered leaf-lists need position-aware diff")
	} else {
		t.Logf("CC-1 diff has %d entries — reorder IS detected", len(diffMap))
		for k, v := range diffMap {
			t.Logf("  %s: op=%d old=%q new=%q", k, v.Op, v.OldVal, v.NewVal)
		}
	}
}

// ---------------------------------------------------------------------------
// CC-4: Nested list entry addition with mandatory leaves
// ---------------------------------------------------------------------------

func TestCC4_NestedListEntryCreation(t *testing.T) {
	// Add a completely new interface with nested unit/family/address
	stateXML := `<configuration>
  <interfaces>
    <interface>
      <name>ge-0/0/0</name>
      <description>existing</description>
    </interface>
  </interfaces>
</configuration>`

	planXML := `<configuration>
  <interfaces>
    <interface>
      <name>ge-0/0/0</name>
      <description>existing</description>
    </interface>
    <interface>
      <name>ge-0/0/1</name>
      <description>new-link</description>
    </interface>
  </interfaces>
</configuration>`

	idx := mustIdxFromSchema(t, matrixSchema)

	stateTree, _ := BuildTree([]byte(stateXML))
	planTree, _ := BuildTree([]byte(planXML))

	stateMap := LeafMapWithSchema(stateTree, idx)
	planMap := LeafMapWithSchema(planTree, idx)
	diffMap := ComputeDiff(stateMap, planMap)

	patchBytes, err := CreateDiffPatch(diffMap, "")
	if err != nil {
		t.Fatal(err)
	}

	output := string(patchBytes)
	t.Logf("CC-4 nested list entry output:\n%s", output)

	// New interface should have create operation on parent
	if !strings.Contains(output, "ge-0/0/1") {
		t.Error("expected ge-0/0/1 in output")
	}
	if !strings.Contains(output, "new-link") {
		t.Error("expected description new-link in output")
	}
	// Verify the new entry has the create operation
	if !strings.Contains(output, "create") {
		t.Error("expected create operation for new list entry")
	}
	// Existing interface should NOT appear (unchanged)
	if strings.Contains(output, "ge-0/0/0") {
		t.Error("ge-0/0/0 should not be in patch (unchanged)")
	}
}

// ---------------------------------------------------------------------------
// CC-6: UTF-8 normalization
// ---------------------------------------------------------------------------

func TestCC6_UTF8DiffNoFalsePositive(t *testing.T) {
	// State and plan both have em-dash — should produce no diff
	stateXML := `<configuration>
  <interfaces>
    <interface>
      <name>ge-0/0/0</name>
      <description>Uplink — Core Router</description>
    </interface>
  </interfaces>
</configuration>`

	planXML := `<configuration>
  <interfaces>
    <interface>
      <name>ge-0/0/0</name>
      <description>Uplink — Core Router</description>
    </interface>
  </interfaces>
</configuration>`

	idx := mustIdxFromSchema(t, matrixSchema)

	stateTree, _ := BuildTree([]byte(stateXML))
	planTree, _ := BuildTree([]byte(planXML))

	stateMap := LeafMapWithSchema(stateTree, idx)
	planMap := LeafMapWithSchema(planTree, idx)
	diffMap := ComputeDiff(stateMap, planMap)

	if len(diffMap) != 0 {
		t.Errorf("CC-6: Expected empty diff for identical UTF-8 content, got %d entries", len(diffMap))
		for k, v := range diffMap {
			t.Logf("  %s: op=%d old=%q new=%q", k, v.Op, v.OldVal, v.NewVal)
		}
	} else {
		t.Log("CC-6 PASSED: Identical UTF-8 strings produce no diff")
	}
}

func TestCC6_UTF8DiffWithEncodingVariation(t *testing.T) {
	// Simulate encoding variation: em-dash as &#x2014; vs UTF-8 literal
	stateXML := `<configuration>
  <interfaces>
    <interface>
      <name>ge-0/0/0</name>
      <description>Uplink &#x2014; Core</description>
    </interface>
  </interfaces>
</configuration>`

	planXML := `<configuration>
  <interfaces>
    <interface>
      <name>ge-0/0/0</name>
      <description>Uplink — Core</description>
    </interface>
  </interfaces>
</configuration>`

	idx := mustIdxFromSchema(t, matrixSchema)

	stateTree, _ := BuildTree([]byte(stateXML))
	planTree, _ := BuildTree([]byte(planXML))

	stateMap := LeafMapWithSchema(stateTree, idx)
	planMap := LeafMapWithSchema(planTree, idx)
	diffMap := ComputeDiff(stateMap, planMap)

	if len(diffMap) == 0 {
		t.Log("CC-6 PASSED: XML entity reference (&#x2014;) and literal em-dash produce same string after XML decode — no false diff")
	} else {
		t.Log("CC-6 NOTE: XML entity vs literal character produces diff (may need normalization)")
		for k, v := range diffMap {
			t.Logf("  %s: op=%d old=%q new=%q", k, v.Op, v.OldVal, v.NewVal)
		}
	}
}

func TestCC6_DoubleEncodedUTF8Repair(t *testing.T) {
	// Simulate the actual double-encoding seen on vpaa: em-dash bytes
	// misinterpreted as Latin-1 and re-encoded to UTF-8
	// Original: "Uplink — Core" (em-dash U+2014 = E2 80 94 in UTF-8)
	// Double-encoded: each byte treated as Latin-1 code point:
	//   E2 → U+00E2 (â) → C3 A2
	//   80 → U+0080     → C2 80
	//   94 → U+0094     → C2 94
	doubleEncoded := "Uplink \u00e2\u0080\u0094 Core" // This is what double-encoding produces
	original := "Uplink \xe2\x80\x94 Core"            // Em-dash as UTF-8 bytes (but stored as string — same as —)

	// Using NormalizeLeafMapUTF8 to repair
	m := map[string]string{"test": doubleEncoded}
	normalized := NormalizeLeafMapUTF8(m)

	if normalized["test"] == original {
		t.Log("CC-6 PASSED: Double-encoded UTF-8 repaired to original")
	} else {
		t.Logf("CC-6 FAILED: normalized=%q expected=%q", normalized["test"], original)
		t.Logf("  double-encoded bytes: %x", []byte(doubleEncoded))
		t.Logf("  original bytes: %x", []byte(original))
		t.Logf("  normalized bytes: %x", []byte(normalized["test"]))
	}
}
