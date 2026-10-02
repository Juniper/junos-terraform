package generic

import (
	"testing"

	"terraform_provider/patch"

	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
)

func TestSanitizeName(t *testing.T) {
	tests := []struct{ in, want string }{
		{"apply-groups", "apply_groups"},
		{"802.1x", "802_1x"},
		{"simple", "simple"},
		{"a-b.c", "a_b_c"},
		{"AH_header", "ah_header"},
	}
	for _, tt := range tests {
		if got := SanitizeName(tt.in); got != tt.want {
			t.Errorf("SanitizeName(%q) = %q, want %q", tt.in, got, tt.want)
		}
	}
}

func attribute(attrs []*tfprotov6.SchemaAttribute, name string) *tfprotov6.SchemaAttribute {
	for _, a := range attrs {
		if a.Name == name {
			return a
		}
	}
	return nil
}

func TestBuildSchema(t *testing.T) {
	s, typ := BuildSchema(patch.CompileSchema([]patch.SchemaNode{
		{Name: "system", Type: "container", Children: []patch.SchemaNode{
			{Name: "host-name", Type: "leaf"},
			{Name: "name-server", Type: "leaf-list"},
		}},
		{Name: "interfaces", Type: "container", Children: []patch.SchemaNode{
			{Name: "interface", Type: "list", Key: "name", Children: []patch.SchemaNode{
				{Name: "name", Type: "leaf"},
				{Name: "unit", Type: "list", Key: "name", Children: []patch.SchemaNode{{Name: "name", Type: "leaf"}}},
			}},
		}},
	}))
	attrs := s.Block.Attributes

	rn := attribute(attrs, "resource_name")
	if rn == nil || !rn.Required || !rn.Type.Is(tftypes.String) {
		t.Fatalf("resource_name: %+v", rn)
	}
	// A schema without groups gives a resource without them, and no error.
	for _, absent := range []string{"groups", "apply_groups"} {
		if attribute(attrs, absent) != nil {
			t.Errorf("%s is not in the schema, so it should not be an attribute", absent)
		}
	}

	system := attribute(attrs, "system")
	if system == nil || system.NestedType == nil || system.NestedType.Nesting != tfprotov6.SchemaObjectNestingModeList || !system.Optional {
		t.Fatalf("system: %+v", system)
	}
	if hn := attribute(system.NestedType.Attributes, "host_name"); hn == nil || !hn.Type.Is(tftypes.String) || !hn.Optional {
		t.Fatalf("host_name: %+v", hn)
	}
	if ns := attribute(system.NestedType.Attributes, "name_server"); ns == nil || !ns.Type.Is(tftypes.List{ElementType: tftypes.String}) {
		t.Fatalf("name_server: %+v", ns)
	}

	iface := attribute(attribute(attrs, "interfaces").NestedType.Attributes, "interface")
	if iface == nil || attribute(iface.NestedType.Attributes, "unit") == nil {
		t.Fatalf("interface: %+v", iface)
	}

	// The value type matches the schema.
	want := tftypes.List{ElementType: tftypes.Object{AttributeTypes: map[string]tftypes.Type{
		"interface": tftypes.List{ElementType: tftypes.Object{AttributeTypes: map[string]tftypes.Type{
			"name": tftypes.String,
			"unit": tftypes.List{ElementType: tftypes.Object{AttributeTypes: map[string]tftypes.Type{"name": tftypes.String}}},
		}}},
	}}}
	if !typ.AttributeTypes["interfaces"].Is(want) {
		t.Fatalf("interfaces type %s, want %s", typ.AttributeTypes["interfaces"], want)
	}
	if _, ok := typ.AttributeTypes["groups"]; ok {
		t.Fatal("groups in the value type")
	}
}

// A provider generated with --groups carries them, and they are ordinary
// configuration: groups is a list keyed by name holding the same hierarchy as
// the base configuration, and apply-groups is a list of strings.
func TestBuildSchemaWithGroups(t *testing.T) {
	s, typ := BuildSchema(patch.CompileSchema([]patch.SchemaNode{
		{Name: "apply-groups", Type: "leaf-list", OrderedBy: "user"},
		{Name: "groups", Type: "list", Key: "name", Children: []patch.SchemaNode{
			{Name: "name", Type: "leaf"},
			{Name: "system", Type: "container", Children: []patch.SchemaNode{
				{Name: "host-name", Type: "leaf"},
			}},
		}},
		{Name: "system", Type: "container", Children: []patch.SchemaNode{
			{Name: "host-name", Type: "leaf"},
		}},
	}))
	attrs := s.Block.Attributes

	ag := attribute(attrs, "apply_groups")
	if ag == nil || !ag.Type.Is(tftypes.List{ElementType: tftypes.String}) {
		t.Fatalf("apply_groups: %+v", ag)
	}

	groups := attribute(attrs, "groups")
	if groups == nil || groups.NestedType == nil {
		t.Fatalf("groups: %+v", groups)
	}
	if attribute(groups.NestedType.Attributes, "name") == nil {
		t.Error("a group is identified by its name")
	}

	// The hierarchy inside a group has the same shape as the base hierarchy.
	inGroup := attribute(groups.NestedType.Attributes, "system")
	base := attribute(attrs, "system")
	if inGroup == nil || base == nil {
		t.Fatalf("system in group %+v, at the top %+v", inGroup, base)
	}
	if !typ.AttributeTypes["groups"].(tftypes.List).ElementType.(tftypes.Object).
		AttributeTypes["system"].Is(typ.AttributeTypes["system"]) {
		t.Error("system inside a group should have the same type as at the top")
	}
}
