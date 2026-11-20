package patch

import (
	"strings"
	"testing"
)

// coalesceSchema: route-filters in the from container of keyed policy terms,
// and a system container with a nested container and plain leaves.
const coalesceSchema = `{
  "path": "",
  "root": {"children": [{"name": "configuration", "type": "container", "path": "", "children": [
    {"name": "policy-options", "type": "container", "path": "", "children": [
      {"name": "policy-statement", "type": "list", "key": "name", "path": "policy-options", "children": [
        {"name": "name", "type": "leaf", "path": "policy-options/policy-statement", "leaf-type": "string"},
        {"name": "term", "type": "list", "key": "name", "path": "policy-options/policy-statement", "children": [
          {"name": "name", "type": "leaf", "path": "policy-options/policy-statement/term", "leaf-type": "string"},
          {"name": "from", "type": "container", "path": "policy-options/policy-statement/term", "children": [
            {"name": "route-filter", "type": "list", "key": "address choice-ident choice-value", "ordered-by": "user", "path": "policy-options/policy-statement/term/from", "children": [
              {"name": "address", "type": "leaf", "path": "policy-options/policy-statement/term/from/route-filter", "leaf-type": "string"},
              {"name": "exact", "type": "leaf", "path": "policy-options/policy-statement/term/from/route-filter", "leaf-type": "string"}
            ]}
          ]}
        ]}
      ]}
    ]},
    {"name": "system", "type": "container", "path": "", "children": [
      {"name": "host-name", "type": "leaf", "path": "system", "leaf-type": "string"},
      {"name": "services", "type": "container", "path": "system", "children": [
        {"name": "ssh", "type": "container", "path": "system/services", "children": [
          {"name": "port", "type": "leaf", "path": "system/services/ssh", "leaf-type": "string"},
          {"name": "max-sessions-per-connection", "type": "leaf", "path": "system/services/ssh", "leaf-type": "string"}
        ]},
        {"name": "netconf-port", "type": "leaf", "path": "system/services", "leaf-type": "string"},
        {"name": "rest-port", "type": "leaf", "path": "system/services", "leaf-type": "string"},
        {"name": "web-port", "type": "leaf", "path": "system/services", "leaf-type": "string"}
      ]}
    ]}
  ]}]}
}`

func coalescePatch(t *testing.T, stateXML, planXML string) string {
	t.Helper()
	idx := mustIdxFromSchema(t, coalesceSchema)
	planMap := LeafMapWithSchema(mustTree(t, planXML), idx)
	diff := ComputeDiff(LeafMapWithSchema(mustTree(t, stateXML), idx), planMap)
	patch, err := CreateDiffPatchWithSchema(diff, planMap, "", idx)
	if err != nil {
		t.Fatalf("CreateDiffPatchWithSchema: %v", err)
	}
	return compactXML(string(patch))
}

func compactXML(s string) string {
	var b strings.Builder
	for _, line := range strings.Split(s, "\n") {
		b.WriteString(strings.TrimSpace(line))
	}
	return b.String()
}

func terms(ts ...string) string {
	return `<configuration><policy-options><policy-statement><name>p</name>` +
		strings.Join(ts, "") + `</policy-statement></policy-options></configuration>`
}

func term(name string, prefixes ...string) string {
	var from string
	if len(prefixes) > 0 {
		from = "<from>"
		for _, p := range prefixes {
			from += `<route-filter><address>` + p + `</address><exact/></route-filter>`
		}
		from += "</from>"
	}
	return `<term><name>` + name + `</name>` + from + `</term>`
}

// Removing one of two route-filters deletes that entry. Every change under
// from is a delete, but from keeps the other entry, so it is not deleted.
func TestCoalesceKeepsContainerWithRemainingEntries(t *testing.T) {
	got := coalescePatch(t,
		terms(term("t", "192.0.2.0/24", "198.51.100.0/24")),
		terms(term("t", "192.0.2.0/24")))
	if strings.Contains(got, `<from nc:operation="delete"`) {
		t.Fatalf("patch deletes from, which keeps 192.0.2.0/24:\n%s", got)
	}
	if !strings.Contains(got, `<route-filter nc:operation="delete"><address>198.51.100.0/24</address>`) {
		t.Fatalf("patch does not delete the 198.51.100.0/24 entry:\n%s", got)
	}
}

// A container removed entirely is deleted by its instance path: this term's
// from, with the policy and term names, not every term's.
func TestCoalesceKeepsKeysInContainerPath(t *testing.T) {
	got := coalescePatch(t,
		terms(term("a", "192.0.2.0/24", "198.51.100.0/24"), term("b", "203.0.113.0/24")),
		terms(term("a"), term("b", "203.0.113.0/24")))
	want := `<policy-statement><name>p</name><term><name>a</name><from nc:operation="delete"/></term></policy-statement>`
	if !strings.Contains(got, want) {
		t.Fatalf("patch missing %s:\n%s", want, got)
	}
	if strings.Contains(got, "<name>b</name>") || strings.Contains(got, "route-filter") {
		t.Fatalf("patch touches more than term a's from:\n%s", got)
	}
}

// The diff does not show unchanged leaves: netconf-port stays, so services is
// not deleted, although every change under it is a delete.
func TestCoalesceKeepsContainerWithUnchangedLeaf(t *testing.T) {
	got := coalescePatch(t,
		`<configuration><system><services><netconf-port>830</netconf-port><rest-port>8080</rest-port><web-port>80</web-port></services></system></configuration>`,
		`<configuration><system><services><netconf-port>830</netconf-port></services></system></configuration>`)
	want := `<configuration><system><services><rest-port nc:operation="delete"/><web-port nc:operation="delete"/></services></system></configuration>`
	if got != want {
		t.Fatalf("got\n%s\nwant\n%s", got, want)
	}
}

// A removed container inside a removed container: one delete, the outer one.
func TestCoalesceNestedEmitsOuterDeleteOnly(t *testing.T) {
	got := coalescePatch(t,
		`<configuration><system><host-name>r1</host-name><services><ssh><port>22</port><max-sessions-per-connection>4</max-sessions-per-connection></ssh><netconf-port>830</netconf-port><web-port>80</web-port></services></system></configuration>`,
		`<configuration><system><host-name>r1</host-name></system></configuration>`)
	want := `<configuration><system><services nc:operation="delete"/></system></configuration>`
	if got != want {
		t.Fatalf("got\n%s\nwant\n%s", got, want)
	}
}
