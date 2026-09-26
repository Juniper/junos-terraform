package patch

import "testing"

// A choice and its cases have no element in the XML: their nodes are indexed
// directly under the choice's parent.
func TestFlattenChoices(t *testing.T) {
	idx, err := UnmarshalTrimmedSchemaIndex(`{"root": {"children": [{"name": "configuration", "type": "container", "children": [
		{"name": "access", "type": "container", "children": [
			{"name": "address-pool", "type": "list", "key": "name", "children": [
				{"name": "name", "type": "leaf"},
				{"name": "address_choice", "type": "choice", "children": [
					{"name": "case_1", "type": "case", "children": [{"name": "address", "type": "leaf"}]},
					{"name": "case_2", "type": "case", "children": [
						{"name": "address-range", "type": "container", "children": [{"name": "low", "type": "leaf"}]}]}
				]}
			]}
		]}
	]}]}}`)
	if err != nil {
		t.Fatal(err)
	}
	for _, p := range []string{"access/address-pool/address", "access/address-pool/address-range/low"} {
		if _, ok := idx[p]; !ok {
			t.Errorf("%s not indexed", p)
		}
	}
	for p := range idx {
		for _, bad := range []string{"address_choice", "case_1", "case_2"} {
			if containsSegment(p, bad) {
				t.Errorf("choice or case in indexed path %s", p)
			}
		}
	}
	pool := idx["access/address-pool"]
	if pool == nil || len(pool.Children) != 3 {
		t.Fatalf("address-pool children: %+v", pool)
	}
}

func containsSegment(path, seg string) bool {
	for _, s := range splitPathRespectingQuotes(path) {
		if s == seg {
			return true
		}
	}
	return false
}
