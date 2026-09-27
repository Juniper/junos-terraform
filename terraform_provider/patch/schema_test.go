package patch

import (
	"reflect"
	"testing"
)

// schemaTestNodes has a choice with the same node in two cases, and nodes
// named "configuration" below the root, as Junos has (system archival
// configuration, snmp trap-group categories configuration).
var schemaTestNodes = []SchemaNode{{Name: "configuration", Type: "container", Children: []SchemaNode{
	{Name: "system", Type: "container", Children: []SchemaNode{
		{Name: "archival", Type: "container", Children: []SchemaNode{
			{Name: "configuration", Type: "container", Presence: "enable", Children: []SchemaNode{
				{Name: "transfer-on-commit", Type: "leaf"},
			}},
		}},
		{Name: "name-server", Type: "leaf-list", OrderedBy: "user"},
	}},
	{Name: "snmp", Type: "container", Children: []SchemaNode{
		{Name: "trap-group", Type: "list", Key: "group-name", Children: []SchemaNode{
			{Name: "group-name", Type: "leaf"},
			{Name: "categories", Type: "container", Children: []SchemaNode{{Name: "configuration", Type: "leaf"}}},
		}},
	}},
	{Name: "routing-options", Type: "container", Children: []SchemaNode{
		{Name: "c", Type: "choice", Children: []SchemaNode{
			{Name: "a", Type: "case", Children: []SchemaNode{{Name: "static", Type: "container", Children: []SchemaNode{{Name: "x", Type: "leaf"}}}}},
			{Name: "b", Type: "case", Children: []SchemaNode{{Name: "static", Type: "container", Children: []SchemaNode{{Name: "y", Type: "leaf"}}}}},
		}},
	}},
}}}

func TestSchemaLookup(t *testing.T) {
	s := CompileSchema(schemaTestNodes)
	for path, want := range map[string]NodeInfo{
		"system":                        {Name: "system", Kind: KindContainer},
		"system/archival/configuration": {Name: "configuration", Kind: KindContainer, Presence: true},
		"system/archival/configuration/transfer-on-commit": {Name: "transfer-on-commit", Kind: KindLeaf},
		"system/name-server":                       {Name: "name-server", Kind: KindLeafList, OrderedByUser: true},
		"snmp/trap-group":                          {Name: "trap-group", Kind: KindList, ListKey: "group-name"},
		"snmp/trap-group/categories":               {Name: "categories", Kind: KindContainer},
		"snmp/trap-group/categories/configuration": {Name: "configuration", Kind: KindLeaf},
		"routing-options/static/x":                 {Name: "x", Kind: KindLeaf},
		"routing-options/static/y":                 {Name: "y", Kind: KindLeaf},
	} {
		got, ok := s.Lookup(path)
		if !ok || got != want {
			t.Errorf("Lookup(%q) = %+v, %v; want %+v", path, got, ok, want)
		}
	}
	for _, path := range []string{"", "configuration", "system/archival/transfer-on-commit", "routing-options/c", "routing-options/c/a/static", "system//archival", "nope"} {
		if got, ok := s.Lookup(path); ok {
			t.Errorf("Lookup(%q) = %+v, want not found", path, got)
		}
	}
	if _, ok := (*Schema)(nil).Lookup("system"); ok {
		t.Error("nil Schema found a node")
	}

	// The two cases' static containers are one node, with both children.
	static, _ := s.find("routing-options/static")
	if s.NumChildren(static) != 2 {
		t.Errorf("merged static has %d children, want 2", s.NumChildren(static))
	}
}

func TestSchemaBinaryRoundTrip(t *testing.T) {
	s := CompileSchema(schemaTestNodes)
	data, err := s.MarshalBinary()
	if err != nil {
		t.Fatal(err)
	}
	if !IsCompiledSchema(data) {
		t.Fatal("not recognised as a compiled schema")
	}
	got, err := UnmarshalSchema(data)
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(got.nodes, s.nodes) || !reflect.DeepEqual(got.names, s.names) {
		t.Fatalf("round trip differs:\n got %+v %q\nwant %+v %q", got.nodes, got.names, s.nodes, s.names)
	}

	for _, bad := range [][]byte{
		nil,
		[]byte("JTAFSCH1"),
		data[:len(data)-1],
		append(append([]byte{}, data...), 0),
		[]byte(`{"root": {}}`),
	} {
		if _, err := UnmarshalSchema(bad); err == nil {
			t.Errorf("UnmarshalSchema(%d bytes) did not fail", len(bad))
		}
	}
	// A child out of range is rejected, not read later.
	corrupt := append([]byte{}, data...)
	corrupt[len(corrupt)-20+8] = 0xff
	if _, err := UnmarshalSchema(corrupt); err == nil {
		t.Error("out-of-range child accepted")
	}
}
