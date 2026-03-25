package patch

import "testing"

const testTrimmedSchema = `{
  "path": "",
  "root": {
    "children": [
      {
        "name": "configuration",
        "type": "container",
        "path": "",
        "children": [
          {
            "name": "foo",
            "type": "container",
            "path": "",
            "children": [
              {
                "name": "item",
                "type": "list",
                "path": "foo",
                "key": "address",
                "children": [
                  {
                    "name": "address",
                    "type": "leaf",
                    "path": "foo/item",
                    "leaf-type": "string"
                  },
                  {
                    "name": "value",
                    "type": "leaf",
                    "path": "foo/item",
                    "leaf-type": "string"
                  }
                ]
              },
              {
                "name": "members",
                "type": "leaf-list",
                "path": "foo",
                "leaf-type": "string"
              }
            ]
          }
        ]
      }
    ]
  }
}`

func mustTree(t *testing.T, xmlStr string) *Node {
	t.Helper()
	tree, err := BuildTree([]byte(xmlStr))
	if err != nil {
		t.Fatalf("BuildTree error: %v", err)
	}
	return tree
}

func mustIdx(t *testing.T) map[string]*NodeInfo {
	t.Helper()
	idx, err := UnmarshalTrimmedSchemaIndex(testTrimmedSchema)
	if err != nil {
		t.Fatalf("UnmarshalTrimmedSchemaIndex error: %v", err)
	}
	return idx
}

func TestLeafMapWithSchema_UsesSchemaListKey(t *testing.T) {
	idx := mustIdx(t)
	xmlStr := `<configuration>
  <groups>
    <name>g1</name>
    <foo>
      <item>
        <address>10.0.0.1</address>
        <value>alpha</value>
      </item>
    </foo>
  </groups>
</configuration>`

	m := LeafMapWithSchema(mustTree(t, xmlStr), idx)
	path := `configuration/groups[name=g1]/foo/item[address=10.0.0.1]/value`
	if got := m[path]; got != "alpha" {
		t.Fatalf("expected %s => alpha, got %q", path, got)
	}

	keyLeafPath := `configuration/groups[name=g1]/foo/item[address=10.0.0.1]/address`
	if _, exists := m[keyLeafPath]; exists {
		t.Fatalf("did not expect key leaf entry %s in leaf map", keyLeafPath)
	}
}

func TestLeafMapWithSchema_LeafListSetDiff(t *testing.T) {
	idx := mustIdx(t)
	stateXML := `<configuration>
  <groups>
    <name>g1</name>
    <foo>
      <members>a</members>
      <members>c</members>
    </foo>
  </groups>
</configuration>`
	planXML := `<configuration>
  <groups>
    <name>g1</name>
    <foo>
      <members>a</members>
      <members>b</members>
    </foo>
  </groups>
</configuration>`

	stateMap := LeafMapWithSchema(mustTree(t, stateXML), idx)
	planMap := LeafMapWithSchema(mustTree(t, planXML), idx)
	diff := ComputeDiff(stateMap, planMap)

	createPath := `configuration/groups[name=g1]/foo/members[value=b]`
	deletePath := `configuration/groups[name=g1]/foo/members[value=c]`
	commonPath := `configuration/groups[name=g1]/foo/members[value=a]`

	if ch, ok := diff[createPath]; !ok || ch.Op != Create {
		t.Fatalf("expected create for %s", createPath)
	}
	if ch, ok := diff[deletePath]; !ok || ch.Op != Delete {
		t.Fatalf("expected delete for %s", deletePath)
	}
	if _, ok := diff[commonPath]; ok {
		t.Fatalf("did not expect diff for unchanged leaf-list value %s", commonPath)
	}
}

func TestOrderedChanges_DeleteBeforeCreate_AndDepthRules(t *testing.T) {
	diffMap := map[string]Change{
		`configuration/groups[name=g1]/foo/bar/baz`: {Op: Delete, OldVal: "x"},
		`configuration/groups[name=g1]/foo/bar`:     {Op: Delete, OldVal: "x"},
		`configuration/groups[name=g1]/foo/alpha`:   {Op: Replace, OldVal: "a", NewVal: "b"},
		`configuration/groups[name=g1]/foo/new`:     {Op: Create, NewVal: "n"},
	}

	ordered := orderedChanges(diffMap)
	if len(ordered) != 4 {
		t.Fatalf("expected 4 ordered changes, got %d", len(ordered))
	}

	if ordered[0].change.Op != Delete || ordered[1].change.Op != Delete {
		t.Fatalf("expected first two operations to be deletes")
	}
	if pathDepth(ordered[0].path) < pathDepth(ordered[1].path) {
		t.Fatalf("expected deeper delete path first: %s then %s", ordered[0].path, ordered[1].path)
	}
	if ordered[2].change.Op != Replace || ordered[3].change.Op != Create {
		t.Fatalf("expected replace before create in trailing operations")
	}
}

func TestComputeDiff_ListKeyRename_ShowsDeleteAndCreate(t *testing.T) {
	stateMap := map[string]string{
		`configuration/groups[name=g1]/foo/item[address=10.0.0.1]/value`: "alpha",
	}
	planMap := map[string]string{
		`configuration/groups[name=g1]/foo/item[address=10.0.0.2]/value`: "alpha",
	}

	diff := ComputeDiff(stateMap, planMap)
	if len(diff) != 2 {
		t.Fatalf("expected 2 changes for key rename, got %d", len(diff))
	}
	if diff[`configuration/groups[name=g1]/foo/item[address=10.0.0.1]/value`].Op != Delete {
		t.Fatalf("expected delete for old key path")
	}
	if diff[`configuration/groups[name=g1]/foo/item[address=10.0.0.2]/value`].Op != Create {
		t.Fatalf("expected create for new key path")
	}
}
