package patch

import (
	"strings"
	"testing"
)

// groupsTrimmedSchema is what --groups produces: groups repeats the
// configuration hierarchy keyed by name, apply-groups is ordered, and the same
// hierarchy also exists at the top for configuration that is not in a group.
const groupsTrimmedSchema = `{
  "path": "",
  "root": {
    "children": [
      {
        "name": "configuration",
        "type": "container",
        "path": "",
        "children": [
          {
            "name": "apply-groups",
            "type": "leaf-list",
            "path": "",
            "leaf-type": "string",
            "ordered-by": "user"
          },
          {
            "name": "groups",
            "type": "list",
            "path": "",
            "key": "name",
            "children": [
              {"name": "name", "type": "leaf", "path": "groups", "leaf-type": "string"},
              {
                "name": "system",
                "type": "container",
                "path": "groups",
                "children": [
                  {"name": "host-name", "type": "leaf", "path": "groups/system", "leaf-type": "string"},
                  {"name": "domain-name", "type": "leaf", "path": "groups/system", "leaf-type": "string"}
                ]
              }
            ]
          },
          {
            "name": "system",
            "type": "container",
            "path": "",
            "children": [
              {"name": "host-name", "type": "leaf", "path": "system", "leaf-type": "string"}
            ]
          }
        ]
      }
    ]
  }
}`

func groupsIdx(t *testing.T) *Schema {
	t.Helper()
	return mustIdxFromSchema(t, groupsTrimmedSchema)
}

// A change inside a group is emitted below its own groups element, not in the
// base hierarchy, which holds the device's own configuration.
func TestGroupLeafChangeStaysInsideItsGroup(t *testing.T) {
	idx := groupsIdx(t)
	state := `<configuration>
  <groups><name>base</name><system><host-name>old</host-name><domain-name>example.net</domain-name></system></groups>
</configuration>`
	plan := `<configuration>
  <groups><name>base</name><system><host-name>new</host-name><domain-name>example.net</domain-name></system></groups>
</configuration>`

	diff := ComputeDiff(LeafMapWithSchema(mustTree(t, state), idx),
		LeafMapWithSchema(mustTree(t, plan), idx))
	if len(diff) != 1 {
		t.Fatalf("expected only the changed leaf, got %v", diff)
	}
	if _, ok := diff[`configuration/groups[name=base]/system/host-name`]; !ok {
		t.Fatalf("change not tracked under its group: %v", diff)
	}

	out, err := CreateDiffPatch(diff, "")
	if err != nil {
		t.Fatal(err)
	}
	got := string(out)
	if !strings.Contains(got, "<groups>") || !strings.Contains(got, "<name>base</name>") {
		t.Fatalf("edit must target the group it belongs to:\n%s", got)
	}
	if strings.Contains(got, "domain-name") {
		t.Fatalf("unchanged leaf in the edit:\n%s", got)
	}
}

// The same path in two groups, and in the base hierarchy, are three different
// places; changing one leaves the others alone.
func TestSamePathInTwoGroupsIsTrackedSeparately(t *testing.T) {
	idx := groupsIdx(t)
	state := `<configuration>
  <groups><name>a</name><system><host-name>one</host-name></system></groups>
  <groups><name>b</name><system><host-name>two</host-name></system></groups>
  <system><host-name>device</host-name></system>
</configuration>`
	plan := `<configuration>
  <groups><name>a</name><system><host-name>changed</host-name></system></groups>
  <groups><name>b</name><system><host-name>two</host-name></system></groups>
  <system><host-name>device</host-name></system>
</configuration>`

	m := LeafMapWithSchema(mustTree(t, state), idx)
	for _, p := range []string{
		`configuration/groups[name=a]/system/host-name`,
		`configuration/groups[name=b]/system/host-name`,
		`configuration/system/host-name`,
	} {
		if _, ok := m[p]; !ok {
			t.Fatalf("%s missing from the leaf map: %v", p, m)
		}
	}

	diff := ComputeDiff(m, LeafMapWithSchema(mustTree(t, plan), idx))
	if len(diff) != 1 {
		t.Fatalf("only group a changed, got %v", diff)
	}
	if ch, ok := diff[`configuration/groups[name=a]/system/host-name`]; !ok || ch.NewVal != "changed" {
		t.Fatalf("expected group a to change, got %v", diff)
	}
}

// apply-groups is ordered: Junos applies groups in the order they are listed.
func TestApplyGroupsOrderIsSignificant(t *testing.T) {
	idx := groupsIdx(t)
	one := `<configuration><apply-groups>a</apply-groups><apply-groups>b</apply-groups></configuration>`
	other := `<configuration><apply-groups>b</apply-groups><apply-groups>a</apply-groups></configuration>`

	same := ComputeDiff(LeafMapWithSchema(mustTree(t, one), idx),
		LeafMapWithSchema(mustTree(t, one), idx))
	if len(same) != 0 {
		t.Fatalf("an unchanged list should produce no diff, got %v", same)
	}

	reordered := ComputeDiff(LeafMapWithSchema(mustTree(t, one), idx),
		LeafMapWithSchema(mustTree(t, other), idx))
	if len(reordered) == 0 {
		t.Fatal("reordering apply-groups must be a change")
	}
}

// A reference that is no longer declared is removed from the device.
func TestApplyGroupsEntryRemoved(t *testing.T) {
	idx := groupsIdx(t)
	state := `<configuration><apply-groups>a</apply-groups><apply-groups>b</apply-groups></configuration>`
	plan := `<configuration><apply-groups>a</apply-groups></configuration>`

	diff := ComputeDiff(LeafMapWithSchema(mustTree(t, state), idx),
		LeafMapWithSchema(mustTree(t, plan), idx))
	if len(diff) == 0 {
		t.Fatal("removing a reference must be a change")
	}
	found := false
	for _, ch := range diff {
		if ch.Op == Delete {
			found = true
		}
	}
	if !found {
		t.Fatalf("expected a delete, got %v", diff)
	}
}

// Removing a managed group is one operation on the group, not one per leaf.
func TestRemovingGroupCoalescesToOneDelete(t *testing.T) {
	idx := groupsIdx(t)
	state := `<configuration>
  <groups><name>gone</name><system><host-name>h</host-name><domain-name>d</domain-name></system></groups>
  <groups><name>kept</name><system><host-name>k</host-name></system></groups>
</configuration>`
	plan := `<configuration>
  <groups><name>kept</name><system><host-name>k</host-name></system></groups>
</configuration>`

	stateMap := LeafMapWithSchema(mustTree(t, state), idx)
	planMap := LeafMapWithSchema(mustTree(t, plan), idx)
	diff := ComputeDiff(stateMap, planMap)

	out, err := CreateDiffPatchWithSchema(diff, planMap, "", idx)
	if err != nil {
		t.Fatal(err)
	}
	got := string(out)
	if strings.Count(got, `nc:operation="delete"`) != 1 {
		t.Fatalf("expected a single delete for the whole group:\n%s", got)
	}
	if !strings.Contains(got, "<name>gone</name>") {
		t.Fatalf("the delete must name the group being removed:\n%s", got)
	}
	if strings.Contains(got, "<name>kept</name>") {
		t.Fatalf("the group that stays must not appear in the edit:\n%s", got)
	}
}

// Junos group names are often wildcards; the key must survive the leaf map and
// come back unchanged in the patch.
func TestWildcardGroupKeyRoundTrips(t *testing.T) {
	idx := groupsIdx(t)
	state := `<configuration>
  <groups><name>ge-*</name><system><host-name>old</host-name></system></groups>
</configuration>`
	plan := `<configuration>
  <groups><name>ge-*</name><system><host-name>new</host-name></system></groups>
</configuration>`

	m := LeafMapWithSchema(mustTree(t, state), idx)
	if _, ok := m[`configuration/groups[name=ge-*]/system/host-name`]; !ok {
		t.Fatalf("wildcard key did not survive the leaf map: %v", m)
	}

	diff := ComputeDiff(m, LeafMapWithSchema(mustTree(t, plan), idx))
	out, err := CreateDiffPatch(diff, "")
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(out), "<name>ge-*</name>") {
		t.Fatalf("wildcard key lost in the patch:\n%s", out)
	}
}
