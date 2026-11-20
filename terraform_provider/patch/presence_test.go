package patch

import (
	"strings"
	"testing"
)

// presenceSchema: bgp group multipath is a presence container, with a child
// as in the full model; bgp and group are not, and neither is traceoptions,
// which has no children here.
const presenceSchema = `{
  "path": "",
  "root": {"children": [{"name": "configuration", "type": "container", "path": "", "children": [
    {"name": "protocols", "type": "container", "path": "", "children": [
      {"name": "bgp", "type": "container", "path": "protocols", "children": [
        {"name": "group", "type": "list", "key": "name", "path": "protocols/bgp", "children": [
          {"name": "name", "type": "leaf", "path": "protocols/bgp/group", "leaf-type": "string"},
          {"name": "type", "type": "leaf", "path": "protocols/bgp/group", "leaf-type": "string"},
          {"name": "multipath", "type": "container", "presence": "enable multipath", "path": "protocols/bgp/group", "children": [
            {"name": "multiple-as", "type": "leaf", "path": "protocols/bgp/group/multipath", "leaf-type": "empty"}
          ]},
          {"name": "traceoptions", "type": "container", "path": "protocols/bgp/group"}
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

// A container without presence is not diffed when empty, even with no
// children in the schema.
func TestEmptyNonPresenceContainerNotDiffed(t *testing.T) {
	withTrace := `<configuration><protocols><bgp><group><name>g</name><type>external</type><traceoptions/></group></bgp></protocols></configuration>`
	if diff, _ := presenceDiff(t, groupWithout, withTrace); len(diff) != 0 {
		t.Fatalf("expected no diff, got %v", diff)
	}
}

// multipath with its child, as a group with multipath { multiple-as; }.
const groupWithMultipleAS = `<configuration><protocols><bgp><group><name>g</name><type>external</type><multipath><multiple-as/></multipath></group></bgp></protocols></configuration>`

// Removing a presence container with a child deletes the container: deleting
// only the child would leave it configured.
func TestPresenceContainerWithChildDelete(t *testing.T) {
	diff, patch := presenceDiff(t, groupWithMultipleAS, groupWithout)
	ch, ok := diff["configuration/protocols/bgp/group[name=g]/multipath"]
	if !ok || ch.Op != Delete {
		t.Fatalf("expected a Delete of multipath, got %v", diff)
	}
	if strings.Count(patch, "<multipath") != 1 || !strings.Contains(patch, `<multipath nc:operation="delete"/>`) ||
		strings.Contains(patch, "multiple-as") {
		t.Fatalf("patch should delete multipath alone:\n%s", patch)
	}
}

// Removing only the child keeps the container.
func TestPresenceContainerChildOnlyDelete(t *testing.T) {
	diff, patch := presenceDiff(t, groupWithMultipleAS, groupWith)
	if _, ok := diff["configuration/protocols/bgp/group[name=g]/multipath"]; ok || len(diff) != 1 {
		t.Fatalf("expected only the multiple-as Delete, got %v", diff)
	}
	if strings.Contains(patch, `<multipath nc:operation`) || !strings.Contains(patch, `<multiple-as nc:operation="delete"/>`) {
		t.Fatalf("patch should delete multiple-as only:\n%s", patch)
	}
}

// Adding a presence container with a child creates it once.
func TestPresenceContainerWithChildCreate(t *testing.T) {
	_, patch := presenceDiff(t, groupWithout, groupWithMultipleAS)
	if n := strings.Count(patch, "<multipath"); n != 1 || strings.Contains(patch, `<multipath nc:operation`) {
		t.Fatalf("patch should have one multipath, created with its child:\n%s", patch)
	}
	if !strings.Contains(patch, "<multiple-as") {
		t.Fatalf("patch does not add multiple-as:\n%s", patch)
	}
}

func TestPresenceContainerWithChildUnchanged(t *testing.T) {
	if diff, _ := presenceDiff(t, groupWithMultipleAS, groupWithMultipleAS); len(diff) != 0 {
		t.Fatalf("expected no diff, got %v", diff)
	}
}
