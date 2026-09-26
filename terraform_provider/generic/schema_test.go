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

func TestValidateNames(t *testing.T) {
	ok := []patch.SchemaNode{{Name: "firewall", Type: "container", Children: []patch.SchemaNode{
		{Name: "AH_header", Type: "leaf"}, {Name: "ESP_header", Type: "leaf"},
	}}}
	if err := ValidateNames(ok, "configuration"); err != nil {
		t.Fatalf("valid names rejected: %v", err)
	}
	collide := []patch.SchemaNode{{Name: "a-b", Type: "leaf"}, {Name: "a_b", Type: "leaf"}}
	if err := ValidateNames(collide, "configuration"); err == nil {
		t.Fatal("expected an error for a-b and a_b")
	}
	invalid := []patch.SchemaNode{{Name: "system", Type: "container", Children: []patch.SchemaNode{{Name: "802.1x", Type: "leaf"}}}}
	if err := ValidateNames(invalid, "configuration"); err == nil {
		t.Fatal("expected an error for a name starting with a digit")
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
	s, typ := BuildSchema([]patch.SchemaNode{
		{Name: "groups", Type: "container"},
		{Name: "apply-groups", Type: "leaf-list"},
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
	})
	attrs := s.Block.Attributes

	rn := attribute(attrs, "resource_name")
	if rn == nil || !rn.Required || !rn.Type.Is(tftypes.String) {
		t.Fatalf("resource_name: %+v", rn)
	}
	for _, skipped := range []string{"groups", "apply_groups"} {
		if attribute(attrs, skipped) != nil {
			t.Errorf("%s should be left out", skipped)
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
