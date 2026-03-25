package patch

import "fmt"

// junosListKeys contains the YANG list key element names common in Junos.
// "name" covers ~95% of cases (interfaces, units, BGP groups, policies, etc.).
// "id" and "type" handle a small number of edge-case list definitions.
var junosListKeys = map[string]bool{
	"name": true,
	"id":   true,
	"type": true,
}

// LeafMap flattens a *Node tree into a map of XPath-style paths to leaf text
// values. Only nodes with no children (true leaves) and non-empty text are
// included. Keyed list entries encode the key in the path segment so siblings
// with different keys are kept distinct:
//
//	interfaces/interface[name=ge-0/0/0]/unit[name=0]/description → "uplink"
//
// Key elements themselves (e.g. <name>ge-0/0/0</name>) are NOT emitted as
// separate entries — they are encoded in the parent segment and cannot be
// independently patched (changing a key requires delete + create).
func LeafMap(root *Node) map[string]string {
	result := make(map[string]string)
	leafMapRecurse(root, "", result)
	return result
}

func leafMapRecurse(node *Node, parentPath string, result map[string]string) {
	segment := buildSegment(node)

	var currentPath string
	if parentPath == "" {
		currentPath = segment
	} else {
		currentPath = parentPath + "/" + segment
	}

	// Leaf node — record it and stop recursing
	if len(node.Children) == 0 {
		if node.Text != "" {
			result[currentPath] = node.Text
		}
		return
	}

	for _, child := range node.Children {
		// Skip the key child — it is already encoded in the current segment.
		// Emitting it separately would create spurious "delete key" operations
		// whenever an ancestor list entry is modified.
		if isKeyChild(child, node) {
			continue
		}
		leafMapRecurse(child, currentPath, result)
	}
}

// buildSegment returns "tag" for plain elements and "tag[keyName=keyValue]"
// for Junos YANG list entries whose first child is a recognised key element.
func buildSegment(node *Node) string {
	if len(node.Children) > 0 {
		first := node.Children[0]
		if junosListKeys[first.Tag] && first.Text != "" {
			return fmt.Sprintf("%s[%s=%s]", node.Tag, first.Tag, first.Text)
		}
	}
	return node.Tag
}

// isKeyChild returns true when child is the key element of a YANG list entry.
// The convention in Junos-generated XML is that the key is always the first
// child of the list element, so we check both position and tag name.
func isKeyChild(child, parent *Node) bool {
	if len(parent.Children) == 0 {
		return false
	}
	return junosListKeys[child.Tag] &&
		parent.Children[0] == child &&
		child.Text != ""
}
