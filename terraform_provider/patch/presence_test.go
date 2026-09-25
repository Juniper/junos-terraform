package patch

import (
	"strings"
	"testing"
)

// presenceSchema: bgp group multipath is a presence container, empty in the
// trimmed schema; bgp and group have children and are not.
const presenceSchema = `{
  "path": "",
  "root": {"children": [{"name": "configuration", "type": "container", "path": "", "children": [
    {"name": "protocols", "type": "container", "path": "", "children": [
      {"name": "bgp", "type": "container", "path": "protocols", "children": [
        {"name": "group", "type": "list", "key": "name", "path": "protocols/bgp", "children": [
          {"name": "name", "type": "leaf", "path": "protocols/bgp/group", "leaf-type": "string"},
          {"name": "type", "type": "leaf", "path": "protocols/bgp/group", "leaf-type": "string"},
          {"name": "multipath", "type": "container", "path": "protocols/bgp/group"}
        ]}
      ]}
    ]}
  ]}]}
}`

func presenceDiff(t *testing.T, stateXML, planXML string) (map[string]Change, string) {
	t.Helper()
	idx := mustIdxFromSchema(t, presenceSchema)
	diff := ComputeDiff(LeafMapWithSchema(mustTree(t, stateXML), idx), LeafMapWithSchema(mustTree(t, planXML), idx))
	// CreateDiffPatch, as the provider calls it.
	patch, err := CreateDiffPatch(diff, "")
	if err != nil {
		t.Fatalf("CreateDiffPatch: %v", err)
	}
	return diff, string(patch)
}

const (
	groupWithout = `<configuration><protocols><bgp><group><name>g</name><type>external</type></group></bgp></protocols></configuration>`
	groupWith    = `<configuration><protocols><bgp><group><name>g</name><type>external</type><multipath/></group></bgp></protocols></configuration>`
	// As Junos returns it: whitespace inside the empty element.
	groupWithDevice = "<configuration><protocols><bgp><group><name>g</name><type>external</type><multipath>\n</multipath></group></bgp></protocols></configuration>"
)

func TestPresenceContainerCreate(t *testing.T) {
	diff, patch := presenceDiff(t, groupWithout, groupWith)
	ch, ok := diff["configuration/protocols/bgp/group[name=g]/multipath"]
	if len(diff) != 1 || !ok || ch.Op != Create {
		t.Fatalf("expected one Create of multipath, got %v", diff)
	}
	if !strings.Contains(patch, `<multipath nc:operation="create"/>`) {
		t.Fatalf("patch does not create multipath:\n%s", patch)
	}
}

func TestPresenceContainerDelete(t *testing.T) {
	diff, patch := presenceDiff(t, groupWithDevice, groupWithout)
	ch, ok := diff["configuration/protocols/bgp/group[name=g]/multipath"]
	if len(diff) != 1 || !ok || ch.Op != Delete {
		t.Fatalf("expected one Delete of multipath, got %v", diff)
	}
	if !strings.Contains(patch, `<multipath nc:operation="delete"/>`) {
		t.Fatalf("patch does not delete multipath:\n%s", patch)
	}
}

func TestPresenceContainerUnchanged(t *testing.T) {
	if diff, _ := presenceDiff(t, groupWithDevice, groupWith); len(diff) != 0 {
		t.Fatalf("expected no diff, got %v", diff)
	}
}

// An enclosing container that is empty because its content is not modeled
// must not be taken for a presence container, or it would be deleted.
func TestEmptyStructuralContainerNotPresence(t *testing.T) {
	diff, patch := presenceDiff(t, `<configuration><protocols><bgp/></protocols></configuration>`, groupWithout)
	for path, ch := range diff {
		if ch.Op == Delete {
			t.Fatalf("unexpected Delete of %s: %v", path, diff)
		}
	}
	if strings.Contains(patch, `operation="delete"`) {
		t.Fatalf("patch deletes something:\n%s", patch)
	}
}
