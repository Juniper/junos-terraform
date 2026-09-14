package generic

import (
	"testing"

	"terraform_provider/patch"

	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func TestSanitizeName(t *testing.T) {
	tests := []struct{ in, want string }{
		{"apply-groups", "apply_groups"},
		{"802.1x", "802_1x"},
		{"simple", "simple"},
		{"a-b.c", "a_b_c"},
	}
	for _, tt := range tests {
		if got := SanitizeName(tt.in); got != tt.want {
			t.Errorf("SanitizeName(%q) = %q, want %q", tt.in, got, tt.want)
		}
	}
}

func TestBuildSchema_SingleLeaf(t *testing.T) {
	roots := []patch.SchemaNode{{
		Name: "configuration", Type: "container",
		Children: []patch.SchemaNode{
			{Name: "host-name", Type: "leaf", LeafType: "string"},
		},
	}}

	s := BuildSchema(roots)
	attr, ok := s.Attributes["host_name"]
	if !ok {
		t.Fatal("expected host_name attribute")
	}
	if _, ok := attr.(schema.StringAttribute); !ok {
		t.Fatalf("expected StringAttribute, got %T", attr)
	}
}

func TestBuildSchema_LeafList(t *testing.T) {
	roots := []patch.SchemaNode{{
		Name: "configuration", Type: "container",
		Children: []patch.SchemaNode{
			{Name: "name-server", Type: "leaf-list", LeafType: "string"},
		},
	}}

	s := BuildSchema(roots)
	attr, ok := s.Attributes["name_server"]
	if !ok {
		t.Fatal("expected name_server attribute")
	}
	la, ok := attr.(schema.ListAttribute)
	if !ok {
		t.Fatalf("expected ListAttribute, got %T", attr)
	}
	if la.ElementType != types.StringType {
		t.Fatalf("expected StringType element, got %v", la.ElementType)
	}
}

func TestBuildSchema_NestedContainer(t *testing.T) {
	roots := []patch.SchemaNode{{
		Name: "configuration", Type: "container",
		Children: []patch.SchemaNode{
			{Name: "interfaces", Type: "container", Children: []patch.SchemaNode{
				{Name: "interface", Type: "list", Children: []patch.SchemaNode{
					{Name: "name", Type: "leaf", LeafType: "string"},
					{Name: "unit", Type: "list", Children: []patch.SchemaNode{
						{Name: "name", Type: "leaf", LeafType: "string"},
					}},
				}},
			}},
		},
	}}

	s := BuildSchema(roots)
	ifaces, ok := s.Attributes["interfaces"]
	if !ok {
		t.Fatal("expected interfaces attribute")
	}
	lna, ok := ifaces.(schema.ListNestedAttribute)
	if !ok {
		t.Fatalf("expected ListNestedAttribute, got %T", ifaces)
	}
	ifaceAttr, ok := lna.NestedObject.Attributes["interface"]
	if !ok {
		t.Fatal("expected interface nested attribute")
	}
	ifaceLNA, ok := ifaceAttr.(schema.ListNestedAttribute)
	if !ok {
		t.Fatalf("expected ListNestedAttribute for interface, got %T", ifaceAttr)
	}
	if _, ok := ifaceLNA.NestedObject.Attributes["name"]; !ok {
		t.Fatal("expected name attribute inside interface")
	}
	if _, ok := ifaceLNA.NestedObject.Attributes["unit"]; !ok {
		t.Fatal("expected unit attribute inside interface")
	}
}

func TestBuildSchema_ResourceName(t *testing.T) {
	roots := []patch.SchemaNode{{
		Name: "configuration", Type: "container",
		Children: []patch.SchemaNode{},
	}}

	s := BuildSchema(roots)
	rn, ok := s.Attributes["resource_name"]
	if !ok {
		t.Fatal("expected resource_name attribute")
	}
	sa, ok := rn.(schema.StringAttribute)
	if !ok {
		t.Fatalf("expected StringAttribute, got %T", rn)
	}
	if !sa.Required {
		t.Fatal("resource_name should be Required")
	}
}

func TestBuildSchema_SkipsGroupsAndApplyGroups(t *testing.T) {
	roots := []patch.SchemaNode{{
		Name: "configuration", Type: "container",
		Children: []patch.SchemaNode{
			{Name: "groups", Type: "container"},
			{Name: "apply-groups", Type: "leaf-list"},
			{Name: "system", Type: "container"},
		},
	}}

	s := BuildSchema(roots)
	if _, ok := s.Attributes["groups"]; ok {
		t.Fatal("groups should be filtered out")
	}
	if _, ok := s.Attributes["apply_groups"]; ok {
		t.Fatal("apply-groups should be filtered out")
	}
	if _, ok := s.Attributes["system"]; !ok {
		t.Fatal("system should be present")
	}
}
