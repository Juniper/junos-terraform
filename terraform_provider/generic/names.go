package generic

import (
	"strings"

	"terraform_provider/patch"
)

// SanitizeName converts a YANG node name to a valid Terraform attribute name:
// lowercase letters, digits and underscores. A few Junos names have capitals
// (AH_header, ESP_header in firewall filters).
func SanitizeName(name string) string {
	return strings.ToLower(strings.NewReplacer("-", "_", ".", "_").Replace(name))
}

// attributeNodes returns the schema nodes that are Terraform attributes, in
// schema order. The schema builder and the value converters both use it, so
// they agree on what the resource holds.
func attributeNodes(nodes []patch.SchemaNode) []patch.SchemaNode {
	out := make([]patch.SchemaNode, 0, len(nodes))
	for _, n := range nodes {
		if n.Name == "" || n.Name == "groups" || n.Name == "apply-groups" {
			continue
		}
		switch n.Type {
		case "leaf", "leaf-list", "container", "list":
			out = append(out, n)
		}
	}
	return out
}
