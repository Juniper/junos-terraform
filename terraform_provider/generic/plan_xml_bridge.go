package generic

import (
	"encoding/xml"
	"strings"

	"terraform_provider/patch"
)

// PlanToXML converts a flat path→value map (from Terraform plan) into Junos XML.
func PlanToXML(values map[string]string, index map[string]*patch.NodeInfo) ([]byte, error) {
	type element struct {
		XMLName  xml.Name
		Children []interface{}
		Value    string `xml:",chardata"`
	}

	root := &element{XMLName: xml.Name{Local: "configuration"}}
	for path, val := range values {
		parts := strings.Split(path, "/")
		cur := root
		for i, part := range parts {
			isLast := i == len(parts)-1
			var found *element
			for _, c := range cur.Children {
				if e, ok := c.(*element); ok && e.XMLName.Local == part {
					found = e
					break
				}
			}
			if found == nil {
				found = &element{XMLName: xml.Name{Local: part}}
				cur.Children = append(cur.Children, found)
			}
			if isLast {
				found.Value = val
			}
			cur = found
		}
	}

	return xml.MarshalIndent(root, "", "  ")
}

// StateToPathMap converts a nested Terraform attribute map into a flat path→value map
// using the schema index for YANG name lookups.
func StateToPathMap(attrs map[string]interface{}, nodes []patch.SchemaNode) map[string]string {
	result := make(map[string]string)
	walkAttrs(attrs, nodes, "", result)
	return result
}

func walkAttrs(attrs map[string]interface{}, nodes []patch.SchemaNode, prefix string, out map[string]string) {
	nodeByTFName := make(map[string]patch.SchemaNode, len(nodes))
	for _, n := range nodes {
		nodeByTFName[SanitizeName(n.Name)] = n
	}

	for key, val := range attrs {
		node, ok := nodeByTFName[key]
		if !ok {
			continue
		}
		yangName := node.Name
		path := yangName
		if prefix != "" {
			path = prefix + "/" + yangName
		}

		switch node.Type {
		case "leaf":
			if s, ok := val.(string); ok && s != "" {
				out[path] = s
			}
		case "leaf-list":
			if list, ok := val.([]interface{}); ok {
				for _, item := range list {
					if s, ok := item.(string); ok {
						out[path] = s // last wins for flat map; patch engine handles multiples via XML
					}
				}
			}
		case "container", "list":
			if list, ok := val.([]interface{}); ok {
				for _, item := range list {
					if m, ok := item.(map[string]interface{}); ok {
						walkAttrs(m, node.Children, path, out)
					}
				}
			}
		}
	}
}
