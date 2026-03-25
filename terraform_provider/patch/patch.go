package patch

import (
	"bytes"
	"fmt"
	"strings"
)

// CreateDiffPatch builds the Junos NETCONF <configuration> XML body from the
// diff map. The output is wrapped in the correct Junos configuration group so
// it targets the same group used by Create and Delete.
//
// The nc:operation attributes written here reference the xmlns:nc declaration
// that sendNetconfPatch places on the enclosing <config> element — no
// additional namespace declaration is required in this output.
//
// Example output for a single Replace:
//
//	<configuration>
//	  <groups>
//	    <name>MY_GROUP</name>
//	    <interfaces>
//	      <interface>
//	        <name>ge-0/0/0</name>
//	        <unit>
//	          <name>0</name>
//	          <description nc:operation="replace">new-desc</description>
//	        </unit>
//	      </interface>
//	    </interfaces>
//	  </groups>
//	</configuration>
func CreateDiffPatch(diffMap map[string]Change, groupName string) ([]byte, error) {
	// Root of the output tree
	root := &Node{Tag: "configuration"}

	// Junos configuration group wrapper
	groups := &Node{Tag: "groups", Parent: root}
	root.Children = append(root.Children, groups)

	groupNameNode := &Node{Tag: "name", Text: groupName, Parent: groups}
	groups.Children = append(groups.Children, groupNameNode)

	for path, change := range diffMap {
		segments := splitPathRespectingQuotes(path)
		if len(segments) == 0 {
			continue
		}

		// If the marshalConfig output includes "configuration" as the root
		// element in the path, strip it — our output tree already has that root.
		if segments[0] == "configuration" {
			segments = segments[1:]
		}
		if len(segments) == 0 {
			continue
		}

		// All segments except the last build the parent hierarchy.
		// The last segment is the leaf that carries the nc:operation.
		parentSegments := segments[:len(segments)-1]
		leafSegment := segments[len(segments)-1]

		// Walk/create the parent node tree under <groups>
		parent := ensurePath(groups, parentSegments)

		// Strip any key predicate from the leaf tag — keys are sibling elements,
		// not part of the tag name itself in XML.
		leafTag, _, _ := parseSegment(leafSegment)

		leaf := &Node{
			Tag:    leafTag,
			Parent: parent,
		}

		switch change.Op {
		case Create:
			leaf.Operation = "create"
			leaf.Text = change.NewVal
		case Replace:
			leaf.Operation = "replace"
			leaf.Text = change.NewVal
		case Delete:
			leaf.Operation = "delete"
			// No text content needed for delete — element will be self-closing
		}

		parent.Children = append(parent.Children, leaf)
	}

	return marshalNodeTree(root)
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
		buf.WriteString(fmt.Sprintf(` %s="%s"`, k, xmlEscape(v)))
	}

	// nc:operation attribute — references xmlns:nc on the <config> ancestor
	if n.Operation != "" {
		buf.WriteString(fmt.Sprintf(` nc:operation="%s"`, n.Operation))
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
