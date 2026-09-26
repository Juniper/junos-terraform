package patch

import "testing"

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
