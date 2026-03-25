package patch

import "strings"

// splitPathRespectingQuotes splits a path string on '/' while treating
// bracket predicates as opaque — a '/' inside "[name=ge-0/0/0]" is not
// treated as a separator.
//
// Example:
//
//	"interfaces/interface[name=ge-0/0/0]/unit[name=0]/description"
//	→ ["interfaces", "interface[name=ge-0/0/0]", "unit[name=0]", "description"]
func splitPathRespectingQuotes(path string) []string {
	var segments []string
	var current strings.Builder
	inBracket := false

	for _, ch := range path {
		switch {
		case !inBracket && ch == '[':
			inBracket = true
			current.WriteRune(ch)
		case inBracket && ch == ']':
			inBracket = false
			current.WriteRune(ch)
		case !inBracket && ch == '/':
			if current.Len() > 0 {
				segments = append(segments, current.String())
				current.Reset()
			}
		default:
			current.WriteRune(ch)
		}
	}
	if current.Len() > 0 {
		segments = append(segments, current.String())
	}
	return segments
}

// parseSegment splits a path segment into its tag and optional key predicate.
//
//	"interface[name=ge-0/0/0]" → ("interface", "name", "ge-0/0/0")
//	"description"              → ("description", "", "")
func parseSegment(seg string) (tag, keyName, keyValue string) {
	idx := strings.Index(seg, "[")
	if idx == -1 {
		return seg, "", ""
	}
	tag = seg[:idx]
	predicate := seg[idx+1 : len(seg)-1] // strip outer [ and ]
	eqIdx := strings.Index(predicate, "=")
	if eqIdx == -1 {
		return tag, "", ""
	}
	keyName = predicate[:eqIdx]
	keyValue = strings.Trim(predicate[eqIdx+1:], "'\"") // strip optional quotes
	return
}

// ensurePath walks the node tree starting at current, creating intermediate
// nodes as needed for each segment, and returns the node at the end of the
// path.
//
// For keyed segments (e.g. "interface[name=ge-0/0/0]") it:
//  1. Looks for an existing child with matching tag AND key child value.
//  2. If not found, creates the element and injects a <name>ge-0/0/0</name>
//     child immediately so subsequent sibling leaf writes land in the right
//     list entry.
func ensurePath(current *Node, segments []string) *Node {
	for _, seg := range segments {
		tag, keyName, keyValue := parseSegment(seg)

		// Search for an existing child that matches this segment
		var found *Node
		for _, child := range current.Children {
			if child.Tag != tag {
				continue
			}
			// Plain element — first match wins
			if keyName == "" {
				found = child
				break
			}
			// Keyed element — must also match the key value
			if findKeyChild(child, keyName) == keyValue {
				found = child
				break
			}
		}

		if found == nil {
			found = &Node{Tag: tag, Parent: current}
			if keyName != "" {
				// Inject key child as the first child of this list entry
				keyNode := &Node{Tag: keyName, Text: keyValue, Parent: found}
				found.Children = append(found.Children, keyNode)
			}
			current.Children = append(current.Children, found)
		}

		current = found
	}
	return current
}

// findKeyChild returns the text of the first child whose tag matches keyName,
// used to disambiguate keyed list entries during ensurePath traversal.
func findKeyChild(node *Node, keyName string) string {
	for _, child := range node.Children {
		if child.Tag == keyName {
			return child.Text
		}
	}
	return ""
}
