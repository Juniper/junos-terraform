package patch

import (
	"fmt"
	"strings"
)

// junosListKeys contains the YANG list key element names common in Junos.
// "name" covers ~95% of cases (interfaces, units, BGP groups, policies, etc.).
// "id" and "type" handle a small number of edge-case list definitions.
var junosListKeys = map[string]bool{
	"name": true,
	"id":   true,
	"type": true,
}

// LeafMapWithSchema flattens an XML tree using schema-derived list keys and
// node kinds from trimmed_schema metadata.
//
// Behavior:
//   - list identity is derived from schema list key (not hardcoded key names)
//   - leaf-list entries are represented as distinct set elements by appending
//     [value=<text>] to the path segment, enabling add/remove diff semantics
//   - key leaf children are excluded from emitted leaves
func LeafMapWithSchema(root *Node, idx map[string]*NodeInfo) map[string]string {
	result := make(map[string]string)
	leafMapRecurseWithSchema(root, "", result, idx)
	return result
}

func leafMapRecurseWithSchema(node *Node, parentPath string, result map[string]string, idx map[string]*NodeInfo) {
	schemaPath := outputPathToSchemaPath(parentPath)
	segment := buildSegmentWithSchema(node, schemaPath, idx)

	currentPath := segment
	if parentPath != "" {
		currentPath = parentPath + "/" + segment
	}

	if len(node.Children) == 0 {
		// Skip empty containers/lists — they have no leaf content to diff.
		// Only emit actual leaves (YANG "empty" type like <any/>, <notice/>
		// or regular text leaves).
		leafSchemaPath := outputPathToSchemaPath(currentPath)
		if info, ok := idx[leafSchemaPath]; ok {
			if info.Kind == KindContainer || info.Kind == KindList {
				return
			}
			if info.Kind == KindLeafList {
				currentPath = currentPath + fmt.Sprintf("[value=%s]", node.Text)
			}
		} else if node.Text == "" {
			// Element not in schema and has no text — likely an empty
			// container or unrecognised element; skip it.
			return
		}

		result[currentPath] = node.Text
		return
	}

	if keyPath, keyValue, ok := structuralKeyedListLeaf(node, currentPath, idx); ok {
		result[keyPath] = keyValue
		return
	}

	// Process ALL children including key children.  Key leaves must appear
	// in the leaf map so ComputeDiff can detect new/removed list entries
	// (where the key leaf is the diff signal for entry-level operations).
	for _, child := range node.Children {
		leafMapRecurseWithSchema(child, currentPath, result, idx)
	}
}

func structuralKeyedListLeaf(node *Node, currentPath string, idx map[string]*NodeInfo) (string, string, bool) {
	schemaPath := outputPathToSchemaPath(currentPath)
	info, ok := idx[schemaPath]
	if !ok || info.Kind != KindList || info.ListKey == "" {
		return "", "", false
	}

	keyValue := keyedListValue(node, info.ListKey)
	if keyValue == "" || subtreeHasMaterialLeaves(node, currentPath, idx) {
		return "", "", false
	}

	return currentPath + "/" + info.ListKey, keyValue, true
}

func keyedListValue(node *Node, keyName string) string {
	for _, child := range node.Children {
		if child.Tag == keyName && child.Text != "" {
			return child.Text
		}
	}

	return ""
}

func subtreeHasMaterialLeaves(node *Node, currentPath string, idx map[string]*NodeInfo) bool {
	for _, child := range node.Children {
		if isKeyChildWithSchema(child, node, currentPath, idx) {
			continue
		}

		schemaPath := outputPathToSchemaPath(currentPath)
		segment := buildSegmentWithSchema(child, schemaPath, idx)
		childPath := segment
		if currentPath != "" {
			childPath = currentPath + "/" + segment
		}

		// A non-key child that is a list entry is material even if its only
		// descendant is its own key — nested list entries are real content.
		childSchemaPath := outputPathToSchemaPath(childPath)
		if info, ok := idx[childSchemaPath]; ok && info.Kind == KindList {
			return true
		}

		if len(child.Children) == 0 {
			// Even empty elements (YANG type "empty") represent material
			// config knobs; their presence prevents key-only early return.
			return true
		}

		if subtreeHasMaterialLeaves(child, childPath, idx) {
			return true
		}
	}

	return false
}

func buildSegmentWithSchema(node *Node, parentSchemaPath string, idx map[string]*NodeInfo) string {
	currentSchemaPath := joinPath(parentSchemaPath, node.Tag)
	if info, ok := idx[currentSchemaPath]; ok && info.Kind == KindList && info.ListKey != "" {
		for _, child := range node.Children {
			if child.Tag == info.ListKey && child.Text != "" {
				return fmt.Sprintf("%s[%s=%s]", node.Tag, info.ListKey, child.Text)
			}
		}
	}

	return buildSegment(node)
}

func isKeyChildWithSchema(child, parent *Node, parentOutputPath string, idx map[string]*NodeInfo) bool {
	parentSchemaPath := outputPathToSchemaPath(parentOutputPath)
	if info, ok := idx[parentSchemaPath]; ok && info.Kind == KindList && info.ListKey != "" {
		return child.Tag == info.ListKey && child.Text != ""
	}

	return false
}

func outputPathToSchemaPath(p string) string {
	if p == "" {
		return ""
	}

	segs := splitPathRespectingQuotes(p)
	out := make([]string, 0, len(segs))
	for i, seg := range segs {
		tag, _, _ := parseSegment(seg)

		if i == 0 && tag == "configuration" {
			continue
		}
		if len(out) == 0 && tag == "groups" {
			continue
		}
		if len(out) == 0 && tag == "name" {
			continue
		}

		if tag != "" {
			out = append(out, tag)
		}
	}

	return strings.Join(out, "/")
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
