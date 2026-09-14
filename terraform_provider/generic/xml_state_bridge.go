package generic

import (
	"encoding/xml"
	"strings"

	"terraform_provider/patch"
)

// xmlNode is a generic XML tree node for parsing device responses.
type xmlNode struct {
	XMLName  xml.Name
	Attrs    []xml.Attr `xml:",any,attr"`
	Content  string     `xml:",chardata"`
	Children []xmlNode  `xml:",any"`
}

// XMLToPathMap parses Junos XML and produces a flat path→value map
// using the schema index to resolve node types.
func XMLToPathMap(xmlBytes []byte, index map[string]*patch.NodeInfo) (map[string]string, error) {
	var root xmlNode
	if err := xml.Unmarshal(xmlBytes, &root); err != nil {
		return nil, err
	}

	result := make(map[string]string)
	walkXML(root.Children, "", result)
	return result, nil
}

func walkXML(nodes []xmlNode, prefix string, out map[string]string) {
	for _, n := range nodes {
		name := n.XMLName.Local
		path := name
		if prefix != "" {
			path = prefix + "/" + name
		}

		if len(n.Children) == 0 {
			val := strings.TrimSpace(n.Content)
			if val != "" {
				out[path] = val
			} else {
				// empty element = presence (e.g. <disable/>)
				out[path] = ""
			}
		} else {
			walkXML(n.Children, path, out)
		}
	}
}

// PathMapToStateAttrs converts a flat path→value map into a nested
// map[string]interface{} matching the Terraform attribute structure.
func PathMapToStateAttrs(pathMap map[string]string, nodes []patch.SchemaNode) map[string]interface{} {
	result := make(map[string]interface{})
	buildStateTree(pathMap, nodes, "", result)
	return result
}

func buildStateTree(pathMap map[string]string, nodes []patch.SchemaNode, prefix string, out map[string]interface{}) {
	for _, n := range nodes {
		yangName := n.Name
		tfName := SanitizeName(yangName)
		path := yangName
		if prefix != "" {
			path = prefix + "/" + yangName
		}

		switch n.Type {
		case "leaf":
			if val, ok := pathMap[path]; ok {
				out[tfName] = val
			}
		case "leaf-list":
			var vals []string
			// Collect all entries with this path prefix
			if val, ok := pathMap[path]; ok {
				vals = append(vals, val)
			}
			if len(vals) > 0 {
				ivals := make([]interface{}, len(vals))
				for i, v := range vals {
					ivals[i] = v
				}
				out[tfName] = ivals
			}
		case "container", "list":
			child := make(map[string]interface{})
			buildStateTree(pathMap, n.Children, path, child)
			if len(child) > 0 {
				out[tfName] = []interface{}{child}
			}
		}
	}
}
