package patch

import (
	"bytes"
	"fmt"
	"sort"
	"strings"
)

// CreateDiffPatch builds the Junos NETCONF <configuration> XML body from the
// diff map, targeting base configuration paths directly.
//
// The nc:operation attributes written here reference the xmlns:nc declaration
// that sendNetconfPatch places on the enclosing <config> element — no
// additional namespace declaration is required in this output.
//
// Example output for a single Replace:
//
//	<configuration>
//	  <interfaces>
//	    <interface>
//	      <name>ge-0/0/0</name>
//	      <unit>
//	        <name>0</name>
//	        <description nc:operation="replace">new-desc</description>
//	      </unit>
//	    </interface>
//	  </interfaces>
//	</configuration>
func CreateDiffPatch(diffMap map[string]Change, groupName string) ([]byte, error) {
	_ = groupName
	// Root of the output tree
	root := &Node{Tag: "configuration"}

	type pendingLeaf struct {
		parent  *Node
		tag     string
		keyName string
		change  Change
	}

	ordered := orderedChanges(diffMap)

	// Two-pass strategy:
	// Pass 1 — process key-entry operations and collect pending leaf ops.
	// This ensures parent nodes get nc:operation="delete" BEFORE we decide
	// whether a child leaf needs its own operation attribute.
	var pending []pendingLeaf
	for _, entry := range ordered {
		path := entry.path
		change := entry.change
		segments := splitPathRespectingQuotes(path)
		if len(segments) == 0 {
			continue
		}

		if segments[0] == "configuration" {
			segments = segments[1:]
		}
		if len(segments) > 0 {
			firstTag, _, _ := parseSegment(segments[0])
			if firstTag == "groups" {
				segments = segments[1:]
			}
		}
		if len(segments) == 0 {
			continue
		}

		parentSegments := segments[:len(segments)-1]
		leafSegment := segments[len(segments)-1]

		parent := ensurePath(root, parentSegments)

		if applyKeyedListEntryOperation(parent, parentSegments, leafSegment, change) {
			continue
		}

		leafTag, keyName, _ := parseSegment(leafSegment)
		pending = append(pending, pendingLeaf{
			parent: parent, tag: leafTag, keyName: keyName, change: change,
		})
	}

	// Pass 2 — create leaf nodes, inheriting context from pass-1 parent ops.
	for _, p := range pending {
		leaf := &Node{
			Tag:    p.tag,
			Parent: p.parent,
		}

		switch p.change.Op {
		case Create:
			leaf.Operation = "create"
			leaf.Text = p.change.NewVal
		case Replace:
			leaf.Operation = "replace"
			leaf.Text = p.change.NewVal
		case Delete:
			// Leaf-list entries (paths with [value=xxx]) need the old value
			// so Junos knows which instance to remove.
			// Scalar leaves must NOT include text — Junos rejects
			// <leaf nc:operation="delete">value</leaf> for scalar leaves.
			if p.keyName == "value" {
				leaf.Text = p.change.OldVal
			}
			// If the parent already has nc:operation="delete" (set by
			// applyKeyedListEntryOperation in pass 1), this leaf is just
			// a structural sibling — do NOT add an operation.  This is
			// critical for Junos compound-key lists where choice-ident
			// elements (e.g. <add/>) must appear WITHOUT an operation.
			if p.parent.Operation == "delete" {
				// structural child — no operation
			} else {
				leaf.Operation = "delete"
			}
		}

		p.parent.Children = append(p.parent.Children, leaf)
	}

	return marshalNodeTree(root)
}

func applyKeyedListEntryOperation(parent *Node, parentSegments []string, leafSegment string, change Change) bool {
	if len(parentSegments) == 0 {
		return false
	}

	_, parentKeyName, parentKeyValue := parseSegment(parentSegments[len(parentSegments)-1])
	leafTag, _, _ := parseSegment(leafSegment)
	if parentKeyName == "" || leafTag != parentKeyName {
		return false
	}

	keyValue := change.NewVal
	if change.Op == Delete {
		keyValue = change.OldVal
	}
	if keyValue == "" || keyValue != parentKeyValue {
		return false
	}

	switch change.Op {
	case Create:
		parent.Operation = "create"
	case Replace:
		parent.Operation = "replace"
	case Delete:
		parent.Operation = "delete"
	default:
		return false
	}

	return true
}

type orderedChange struct {
	path   string
	change Change
}

func orderedChanges(diffMap map[string]Change) []orderedChange {
	result := make([]orderedChange, 0, len(diffMap))
	for path, change := range diffMap {
		result = append(result, orderedChange{path: path, change: change})
	}

	sort.SliceStable(result, func(i, j int) bool {
		a := result[i]
		b := result[j]

		pa := opPriority(a.change.Op)
		pb := opPriority(b.change.Op)
		if pa != pb {
			return pa < pb
		}

		da := pathDepth(a.path)
		db := pathDepth(b.path)
		if a.change.Op == Delete {
			if da != db {
				return da > db
			}
		} else {
			if da != db {
				return da < db
			}
		}

		return a.path < b.path
	})

	return result
}

func opPriority(op ChangeType) int {
	switch op {
	case Delete:
		return 0
	case Replace:
		return 1
	case Create:
		return 2
	default:
		return 3
	}
}

func pathDepth(path string) int {
	if path == "" {
		return 0
	}
	return len(splitPathRespectingQuotes(path))
}

// marshalNodeTree serializes a *Node tree to indented XML bytes.
func marshalNodeTree(root *Node) ([]byte, error) {
	var buf bytes.Buffer
	if err := encodeNode(&buf, root, 0); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

// encodeNode recursively writes a node and all its descendants to buf.
func encodeNode(buf *bytes.Buffer, n *Node, depth int) error {
	indent := strings.Repeat("  ", depth)
	buf.WriteString(indent + "<" + n.Tag)

	// Standard XML attributes
	for k, v := range n.Attrs {
		if _, err := fmt.Fprintf(buf, ` %s="%s"`, k, xmlEscape(v)); err != nil {
			return err
		}
	}

	// nc:operation attribute — references xmlns:nc on the <config> ancestor
	if n.Operation != "" {
		if _, err := fmt.Fprintf(buf, ` nc:operation="%s"`, n.Operation); err != nil {
			return err
		}
	}

	// Self-closing for delete and empty nodes
	if len(n.Children) == 0 && n.Text == "" {
		buf.WriteString("/>\n")
		return nil
	}

	buf.WriteString(">")

	if len(n.Children) > 0 {
		buf.WriteString("\n")
		for _, child := range n.Children {
			if err := encodeNode(buf, child, depth+1); err != nil {
				return err
			}
		}
		buf.WriteString(indent + "</" + n.Tag + ">\n")
	} else {
		// Inline text with XML escaping
		buf.WriteString(xmlEscape(n.Text) + "</" + n.Tag + ">\n")
	}

	return nil
}

// xmlEscape escapes the five XML special characters in text content and
// attribute values.
func xmlEscape(s string) string {
	s = strings.ReplaceAll(s, "&", "&amp;") // must be first
	s = strings.ReplaceAll(s, "<", "&lt;")
	s = strings.ReplaceAll(s, ">", "&gt;")
	s = strings.ReplaceAll(s, "\"", "&quot;")
	s = strings.ReplaceAll(s, "'", "&apos;")
	return s
}
