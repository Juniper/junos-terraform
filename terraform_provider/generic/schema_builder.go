package generic

import (
	"strings"

	"terraform_provider/patch"

	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// SanitizeName converts a YANG node name to a valid Terraform attribute name.
func SanitizeName(name string) string {
	return strings.NewReplacer("-", "_", ".", "_").Replace(name)
}

// BuildSchema produces a schema.Schema from the pyang JSON SchemaNode tree.
func BuildSchema(roots []patch.SchemaNode) schema.Schema {
	// roots[0] is the "configuration" node; use its children.
	var configChildren []patch.SchemaNode
	if len(roots) > 0 {
		configChildren = roots[0].Children
	}

	attrs := buildAttributes(configChildren)
	attrs["resource_name"] = schema.StringAttribute{
		Required:      true,
		PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplace()},
	}

	return schema.Schema{Attributes: attrs}
}

func buildAttributes(nodes []patch.SchemaNode) map[string]schema.Attribute {
	attrs := make(map[string]schema.Attribute, len(nodes))
	for _, n := range nodes {
		if n.Name == "" || n.Name == "groups" || n.Name == "apply-groups" {
			continue
		}
		key := SanitizeName(n.Name)
		switch n.Type {
		case "leaf":
			attrs[key] = schema.StringAttribute{Optional: true}
		case "leaf-list":
			attrs[key] = schema.ListAttribute{
				ElementType: types.StringType,
				Optional:    true,
			}
		case "container", "list":
			nested := buildAttributes(n.Children)
			attrs[key] = schema.ListNestedAttribute{
				Optional: true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: nested,
				},
			}
		}
	}
	return attrs
}
