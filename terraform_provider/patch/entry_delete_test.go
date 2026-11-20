package patch

import (
	"strings"
	"testing"
)

// A deleted term is deleted by its key alone: a delete of its route-filter,
// deeper in its from, would follow the term's and fail on Junos with
// "statement not found: term a".
func TestDeletedEntryDropsNestedDeletes(t *testing.T) {
	idx := mustIdxFromSchema(t, coalesceSchema)
	diff := ComputeDiff(
		LeafMapWithSchema(mustTree(t, terms(term("a", "192.0.2.0/24"), term("b", "198.51.100.0/24"))), idx),
		LeafMapWithSchema(mustTree(t, terms(term("b", "198.51.100.0/24"))), idx))
	// CreateDiffPatch, as the provider calls it.
	patch, err := CreateDiffPatch(diff, "")
	if err != nil {
		t.Fatalf("CreateDiffPatch: %v", err)
	}
	got := compactXML(string(patch))
	want := `<configuration><policy-options><policy-statement><name>p</name><term nc:operation="delete"><name>a</name></term></policy-statement></policy-options></configuration>`
	if got != want {
		t.Fatalf("got\n%s\nwant\n%s", got, want)
	}
}

// Deleting a BGP neighbor or an address-book entry should send only the
// entry's name. Before this was fixed, the entry's other settings were sent
// as well, as empty tags (<local-address/>, <ip-prefix/>), and Junos rejected
// the whole change with "invalid ip address or hostname".
func TestDeletedEntryDropsValuedChildren(t *testing.T) {
	diff := map[string]Change{
		`configuration/protocols/bgp/group[name=g]/neighbor[name=192.0.2.1]/name`:          {Op: Delete, OldVal: "192.0.2.1"},
		`configuration/protocols/bgp/group[name=g]/neighbor[name=192.0.2.1]/local-address`: {Op: Delete, OldVal: "192.0.2.0"},
		`configuration/security/address-book[name=global]/address[name=a]/name`:            {Op: Delete, OldVal: "a"},
		`configuration/security/address-book[name=global]/address[name=a]/ip-prefix`:       {Op: Delete, OldVal: "198.51.100.0/24"},
	}
	patch, err := CreateDiffPatch(diff, "")
	if err != nil {
		t.Fatalf("CreateDiffPatch: %v", err)
	}
	got := compactXML(string(patch))
	for _, want := range []string{
		`<neighbor nc:operation="delete"><name>192.0.2.1</name></neighbor>`,
		`<address nc:operation="delete"><name>a</name></address>`,
	} {
		if !strings.Contains(got, want) {
			t.Fatalf("got\n%s\nwant it to contain\n%s", got, want)
		}
	}
}
