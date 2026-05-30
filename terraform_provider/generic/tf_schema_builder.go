package generic

import (
	"strings"

	"terraform_provider/patch"

	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// BuildTerraformSchema constructs a Terraform resource schema dynamically from
// the SchemaIndex. Each YANG node maps to a Terraform attribute type.
func BuildTerraformSchema(idx *SchemaIndex) schema.Schema {
	attrs := map[string]schema.Attribute{
		"resource_name": schema.StringAttribute{
			Required:      true,
			PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplace()},
		},
	}

	for _, node := range idx.TopLevel {
		name := normalizeName(node.Name)
		attrs[name] = buildAttribute(node)
	}

	return schema.Schema{Attributes: attrs}
}

func buildAttribute(node patch.SchemaNode) schema.Attribute {
	switch node.Type {
	case "leaf":
		return schema.StringAttribute{
			Optional: true,
		}
	case "leaf-list":
		return schema.ListAttribute{
			ElementType: types.StringType,
			Optional:    true,
		}
	case "container":
		nested := buildNestedAttributes(node.Children)
		if len(nested) == 0 {
			// Empty container — model as optional string presence marker
			return schema.StringAttribute{
				Optional: true,
			}
		}
		return schema.ListNestedAttribute{
			Optional: true,
			NestedObject: schema.NestedAttributeObject{
				Attributes: nested,
			},
		}
	case "list":
		nested := buildNestedAttributes(node.Children)
		if len(nested) == 0 {
			return schema.ListAttribute{
				ElementType: types.StringType,
				Optional:    true,
			}
		}
		return schema.ListNestedAttribute{
			Optional: true,
			NestedObject: schema.NestedAttributeObject{
				Attributes: nested,
			},
		}
	default:
		return schema.StringAttribute{
			Optional: true,
		}
	}
}

func buildNestedAttributes(children []patch.SchemaNode) map[string]schema.Attribute {
	attrs := make(map[string]schema.Attribute, len(children))
	for _, child := range children {
		name := normalizeName(child.Name)
		attrs[name] = buildAttribute(child)
	}
	return attrs
}

// normalizeName converts a YANG name (with hyphens/dots) to a Terraform-safe attribute name.
func normalizeName(name string) string {
	name = strings.ReplaceAll(name, "-", "_")
	name = strings.ReplaceAll(name, ".", "_")
	return name
}
