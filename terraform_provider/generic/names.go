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

// attributeNodes returns the children of a schema node that are Terraform
// attributes, in schema order. The schema builder and the value converters
// both use it, so they agree on what the resource holds.
func attributeNodes(s *patch.Schema, parent patch.SchemaNodeID) []patch.SchemaNodeID {
	out := make([]patch.SchemaNodeID, 0, s.NumChildren(parent))
	for i := range s.NumChildren(parent) {
		c := s.Child(parent, i)
		if name := s.Name(c); name == "groups" || name == "apply-groups" {
			continue
		}
		out = append(out, c)
	}
	return out
}
