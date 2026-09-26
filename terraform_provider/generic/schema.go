package generic

import (
	"terraform_provider/patch"

	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
)

// resourceNameAttribute is the resource's name in Terraform, not device
// configuration. Changing it replaces the resource.
const resourceNameAttribute = "resource_name"

// BuildSchema returns the resource's protocol schema and its value type, built
// from the schema nodes under <configuration>. A leaf is an optional string, a
// leaf-list an optional list of strings, and a container or list an optional
// list of nested objects (a container has at most one entry).
func BuildSchema(configNodes []patch.SchemaNode) (*tfprotov6.Schema, tftypes.Object) {
	attrs, types := buildAttributes(configNodes)
	attrs = append([]*tfprotov6.SchemaAttribute{{
		Name:     resourceNameAttribute,
		Type:     tftypes.String,
		Required: true,
	}}, attrs...)
	types[resourceNameAttribute] = tftypes.String
	return &tfprotov6.Schema{Block: &tfprotov6.SchemaBlock{Attributes: attrs}},
		tftypes.Object{AttributeTypes: types}
}

func buildAttributes(nodes []patch.SchemaNode) ([]*tfprotov6.SchemaAttribute, map[string]tftypes.Type) {
	nodes = attributeNodes(nodes)
	attrs := make([]*tfprotov6.SchemaAttribute, 0, len(nodes))
	types := make(map[string]tftypes.Type, len(nodes))
	for _, n := range nodes {
		name := SanitizeName(n.Name)
		switch n.Type {
		case "leaf":
			attrs = append(attrs, &tfprotov6.SchemaAttribute{Name: name, Type: tftypes.String, Optional: true})
			types[name] = tftypes.String
		case "leaf-list":
			t := tftypes.List{ElementType: tftypes.String}
			attrs = append(attrs, &tfprotov6.SchemaAttribute{Name: name, Type: t, Optional: true})
			types[name] = t
		case "container", "list":
			nested, nestedTypes := buildAttributes(n.Children)
			attrs = append(attrs, &tfprotov6.SchemaAttribute{
				Name:     name,
				Optional: true,
				NestedType: &tfprotov6.SchemaObject{
					Nesting:    tfprotov6.SchemaObjectNestingModeList,
					Attributes: nested,
				},
			})
			types[name] = tftypes.List{ElementType: tftypes.Object{AttributeTypes: nestedTypes}}
		}
	}
	return attrs, types
}
