package patch

import (
	"strings"
	"testing"
)

// compoundKeySchema: route-filter as the YANG models it, a list keyed on
// "address choice-ident choice-value" — the form Junos returns with
// `system services netconf yang-compliant`.
const compoundKeySchema = `{
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
              {"name": "choice-ident", "type": "leaf", "path": "policy-options/policy-statement/term/from/route-filter", "leaf-type": "enumeration"},
              {"name": "choice-value", "type": "leaf", "path": "policy-options/policy-statement/term/from/route-filter", "leaf-type": "string"}
            ]}
          ]}
        ]}
      ]}
    ]}
  ]}]}
}`

func routeFilters(filters ...string) string {
	return `<configuration><policy-options><policy-statement><name>p</name><term><name>t</name><from>` +
		strings.Join(filters, "") + `</from></term></policy-statement></policy-options></configuration>`
}

func routeFilter(address, ident, value string) string {
	return `<route-filter><address>` + address + `</address><choice-ident>` + ident +
		`</choice-ident><choice-value>` + value + `</choice-value></route-filter>`
}

func compoundKeyDiff(t *testing.T, stateXML, planXML string) (map[string]Change, string) {
	t.Helper()
	idx := mustIdxFromSchema(t, compoundKeySchema)
	diff := ComputeDiff(LeafMapWithSchema(mustTree(t, stateXML), idx), LeafMapWithSchema(mustTree(t, planXML), idx))
	// CreateDiffPatch, as the provider calls it.
	patch, err := CreateDiffPatch(diff, "")
	if err != nil {
		t.Fatalf("CreateDiffPatch: %v", err)
	}
	return diff, string(patch)
}

// Each entry's patch element: the route-filter with its operation and keys.
func routeFilterOp(op, address, ident, value string) string {
	valueEl := "<choice-value/>"
	if value != "" {
		valueEl = "<choice-value>" + value + "</choice-value>"
	}
	return `<route-filter nc:operation="` + op + `">` +
		`<address>` + address + `</address>` +
		`<choice-ident>` + ident + `</choice-ident>` +
		valueEl + `</route-filter>`
}

func compact(s string) string {
	var b strings.Builder
	for _, line := range strings.Split(s, "\n") {
		b.WriteString(strings.TrimSpace(line))
	}
	return b.String()
}

func TestCompoundKeyUnchanged(t *testing.T) {
	// As Junos returns it: whitespace inside the empty choice-value.
	device := routeFilters(`<route-filter><address>::/0</address><choice-ident>exact</choice-ident><choice-value>
</choice-value></route-filter>`)
	if diff, _ := compoundKeyDiff(t, device, routeFilters(routeFilter("::/0", "exact", ""))); len(diff) != 0 {
		t.Fatalf("expected no diff, got %v", diff)
	}
}

// Every key identifies the entry: changing the match type is a different
// entry, so the old one is deleted and the new one created — not a replace of
// choice-ident, which would leave the old entry and add a second.
func TestCompoundKeyChangeIsDeleteAndCreate(t *testing.T) {
	_, patch := compoundKeyDiff(t,
		routeFilters(routeFilter("::/0", "exact", "")),
		routeFilters(routeFilter("::/0", "orlonger", "")))
	got := compact(patch)
	for _, want := range []string{
		routeFilterOp("delete", "::/0", "exact", ""),
		routeFilterOp("create", "::/0", "orlonger", ""),
	} {
		if !strings.Contains(got, want) {
			t.Fatalf("patch missing %s:\n%s", want, patch)
		}
	}
	if strings.Contains(got, `operation="replace"`) {
		t.Fatalf("patch replaces a key leaf:\n%s", patch)
	}
}

func TestCompoundKeyValueChangeIsDeleteAndCreate(t *testing.T) {
	_, patch := compoundKeyDiff(t,
		routeFilters(routeFilter("2001:db8::/48", "prefix-length-range", "/64-/64")),
		routeFilters(routeFilter("2001:db8::/48", "prefix-length-range", "/64-/128")))
	got := compact(patch)
	for _, want := range []string{
		routeFilterOp("delete", "2001:db8::/48", "prefix-length-range", "/64-/64"),
		routeFilterOp("create", "2001:db8::/48", "prefix-length-range", "/64-/128"),
	} {
		if !strings.Contains(got, want) {
			t.Fatalf("patch missing %s:\n%s", want, patch)
		}
	}
}

// Two entries on the same prefix are distinct: adding one leaves the other.
func TestCompoundKeySamePrefixAdd(t *testing.T) {
	_, patch := compoundKeyDiff(t,
		routeFilters(routeFilter("::/0", "exact", "")),
		routeFilters(routeFilter("::/0", "exact", ""), routeFilter("::/0", "orlonger", "")))
	got := compact(patch)
	if want := routeFilterOp("create", "::/0", "orlonger", ""); !strings.Contains(got, want) {
		t.Fatalf("patch missing %s:\n%s", want, patch)
	}
	if strings.Contains(got, "exact") || strings.Contains(got, "delete") {
		t.Fatalf("patch touches the unchanged entry:\n%s", patch)
	}
}

func TestCompoundKeySamePrefixRemove(t *testing.T) {
	_, patch := compoundKeyDiff(t,
		routeFilters(routeFilter("::/0", "exact", ""), routeFilter("::/0", "orlonger", "")),
		routeFilters(routeFilter("::/0", "exact", "")))
	got := compact(patch)
	if want := routeFilterOp("delete", "::/0", "orlonger", ""); !strings.Contains(got, want) {
		t.Fatalf("patch missing %s:\n%s", want, patch)
	}
	if strings.Contains(got, "exact") || strings.Contains(got, "create") {
		t.Fatalf("patch touches the unchanged entry:\n%s", patch)
	}
}

// Reordering keeps same-prefix entries apart: each follows its own position.
func TestCompoundKeyAlignOrder(t *testing.T) {
	idx := mustIdxFromSchema(t, compoundKeySchema)
	current := routeFilters(routeFilter("::/0", "orlonger", ""), routeFilter("::/0", "exact", ""))
	reference := routeFilters(routeFilter("::/0", "exact", ""), routeFilter("::/0", "orlonger", ""))
	aligned, err := AlignXMLOrderToReference([]byte(current), []byte(reference), idx)
	if err != nil {
		t.Fatalf("AlignXMLOrderToReference: %v", err)
	}
	got := string(aligned)
	if strings.Index(got, "exact") > strings.Index(got, "orlonger") {
		t.Fatalf("expected exact before orlonger:\n%s", got)
	}
}
