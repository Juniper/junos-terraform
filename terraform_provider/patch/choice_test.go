package patch

import (
	"strings"
	"testing"
)

// A choice and its cases have no element in the XML: their nodes are indexed
// directly under the choice's parent.
func TestCompileSchemaFlattensChoices(t *testing.T) {
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
		if _, ok := idx.Lookup(p); !ok {
			t.Errorf("%s not indexed", p)
		}
	}
	for _, p := range []string{"access/address-pool/address_choice", "access/address-pool/address_choice/case_1/address"} {
		if _, ok := idx.Lookup(p); ok {
			t.Errorf("choice or case indexed: %s", p)
		}
	}
	pool, ok := idx.find("access/address-pool")
	if !ok || idx.NumChildren(pool) != 3 {
		t.Fatalf("address-pool children: %v %v", pool, ok)
	}
	var names []string
	for i := 0; i < idx.NumChildren(pool); i++ {
		names = append(names, idx.Name(idx.Child(pool, i)))
	}
	if strings.Join(names, " ") != "name address address-range" {
		t.Fatalf("address-pool children %v, want schema order", names)
	}
}
