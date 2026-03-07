package patch

import (
	"encoding/xml"
	"strings"
	"fmt"
    "io"
    "bytes"
)

// ------------------------- Type Definitions [START] -------------------------

// XML tree node
type Node struct {
	Name     string           
	Attrs    map[string]string 
	Text     string           
	Children []*Node           
    Parent   *Node 
}

// Diff struct
type Change struct {
    Op   string // "create" | "delete" | "replace"
    OldV string
    NewV string
}

// Minimal XML element with an operation
type PatchNode struct {
	XMLName xml.Name
	AttrOp  string  `xml:"operation,attr,omitempty"`
	Text    string  `xml:",chardata"`
	Children []*PatchNode `xml:",any"`
}

// ------------------------- Type Definitions [END] -------------------------

// ------------------------- Diff Functions [START] -------------------------

// Build XML tree
func BuildTree(xmlBytes []byte) (*Node, error) {
	dec := xml.NewDecoder(bytes.NewReader(xmlBytes))
	var stack []*Node
	var root *Node

	for {
		tok, err := dec.Token()
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}

		switch t := tok.(type) {
		case xml.StartElement:
			n := &Node{
				Name:  t.Name.Local,
				Attrs: map[string]string{},
			}
			for _, a := range t.Attr {
				n.Attrs[a.Name.Local] = a.Value
			}
			if len(stack) == 0 {
				root = n
			} else {
				parent := stack[len(stack)-1]
                n.Parent = parent
				parent.Children = append(parent.Children, n)
			}
			stack = append(stack, n)

		case xml.EndElement:
			// Pop
			stack = stack[:len(stack)-1]

		case xml.CharData:
			// Accumulate trimmed text on current node
			if len(stack) == 0 {
				continue
			}
			txt := strings.TrimSpace(string(t))
			if txt != "" {
				cur := stack[len(stack)-1]
				// Keep a single space between chunks if needed
				if cur.Text == "" {
					cur.Text = txt
				} else {
					cur.Text += " " + txt
				}
			}
		}
	}
	return root, nil
}

// LeafMap builds a flat path -> value map
func LeafMap(root *Node) map[string]string {
	out := make(map[string]string)

	var walk func(n *Node, stack []string)
	walk = func(n *Node, stack []string) {
		// Build a stable/unique path segment for this node
		seg := segmentWithSiblingIndex(n)
		stack = append(stack, seg)

		nameVal, hasName := childNameValue(n)
		otherKids := hasNonNameChildren(n)
		textVal := strings.TrimSpace(n.Text)

		switch {
		// Plain scalar leaf
		case !otherKids && n.Name != "name" && textVal != "":
			path := strings.Join(stack, "/")
                i := 0
                basePath := path

                for {
                    if _, ok := out[path]; !ok {
                        break
                    }
                    i++
                    path = fmt.Sprintf("%s#%d", basePath, i)
                }

                out[path] = textVal

		//  Node whose only real payload is <name>child</name>
		case !otherKids && hasName && n.Name != "name" && textVal == "":
			path := strings.Join(stack, "/")
			out[path] = nameVal

		default:
			// Keep walking into children, but skip <name> because it was already used as a predicate or as the leaf value.
			for _, ch := range n.Children {
				if ch.Name == "name" {
					continue
				}
				walk(ch, stack)
			}
		}
	}

	walk(root, nil)
	return out
}

// CreateDiffMap creates a map of differences between plan and state
func CreateDiffMap(planMap map[string]string, stateMap map[string]string, idx map[string]*NodeInfo) (changes map[string]Change){
	changes = make(map[string]Change)

	// Deletions & candidates for replace
    for k, lv := range stateMap {
        if rv, ok := planMap[k]; !ok {
            changes[k] = Change{Op: "delete", OldV: lv, NewV: ""}
        } else if rv != lv {
            changes[k] = Change{Op: "replace", OldV: lv, NewV: rv}
        }
    }

    // Creations
    for k, rv := range planMap {
        if _, ok := stateMap[k]; !ok {
            changes[k] = Change{Op: "create", OldV: "", NewV: rv}
        }
    }

	return changes
}

// CreateDiffPatch creates a string of the given changes map
func CreateDiffPatch(changes map[string]Change, group string) (string, error) {
	root := &PatchNode{XMLName: xml.Name{Local: "configuration"}}

	for path, change := range changes {
		segs := splitPathRespectQuotes(path)[1:] // drop leading "configuration"
		parent := ensurePath(root, segs[:len(segs)-1])
		leafName := segs[len(segs)-1]
        if i := strings.IndexByte(leafName, '#'); i >= 0 {
            leafName = leafName[:i]
        }
		switch change.Op {
		case "delete":
			node := &PatchNode{XMLName: xml.Name{Local: leafName}, AttrOp: "delete"}
			parent.Children = append(parent.Children, node)

		case "create":
			node := &PatchNode{XMLName: xml.Name{Local: leafName}, AttrOp: "create", Text: change.NewV}
			parent.Children = append(parent.Children, node)

		case "replace":
			delNode := &PatchNode{XMLName: xml.Name{Local: leafName}, AttrOp: "delete"}
			addNode := &PatchNode{XMLName: xml.Name{Local: leafName}, AttrOp: "create", Text: change.NewV}
			parent.Children = append(parent.Children, delNode, addNode)
		}
	}

	var buf bytes.Buffer
	buf.WriteString(xml.Header)
	enc := xml.NewEncoder(&buf)
	enc.Indent("", "  ")
	if err := enc.Encode(root); err != nil {
		return "", err
	}
	enc.Flush()

    diff := buf.String()

	return diff, nil
}

// ------------------------- Main Diff Functions [END] -------------------------

// ------------------------- Helpers [START] -------------------------

func segmentWithSiblingIndex(n *Node) string {
	// <name> child
	if v, ok := childNameValue(n); ok && v != "" {
		return fmt.Sprintf("%s[name=%q]", n.Name, v)
	}

	return n.Name
}

// childNameValue finds an immediate <name> childs text.
func childNameValue(n *Node) (string, bool) {
	for _, ch := range n.Children {
		if ch.Name == "name" {
			val := strings.TrimSpace(ch.Text)
			if val != "" {
				return val, true
			}
		}
	}
	return "", false
}

// hasNonNameChildren reports whether n has any element children other than <name>.
func hasNonNameChildren(n *Node) bool {
	for _, ch := range n.Children {
		if ch.Name != "" && ch.Name != "name" {
			return true
		}
	}
	return false
}

func splitPathRespectQuotes(path string) []string {
	var segs []string
	var sb strings.Builder
	inQuotes := false

	for i := 0; i < len(path); i++ {
		ch := path[i]

		switch ch {
		case '"':
			inQuotes = !inQuotes
			sb.WriteByte(ch)

		case '/':
			if inQuotes {
				sb.WriteByte(ch)
			} else {
				// new segment boundary
				if sb.Len() > 0 {
					segs = append(segs, sb.String())
					sb.Reset()
				}
			}

		default:
			sb.WriteByte(ch)
		}
	}

	if sb.Len() > 0 {
		segs = append(segs, sb.String())
	}

	// remove a leading empty segment if the path starts with '/'
	if len(segs) > 0 && segs[0] == "" {
		segs = segs[1:]
	}

	return segs
}

// ensurePath builds nested XML elements for the given path segments
func ensurePath(root *PatchNode, segs []string) *PatchNode {
	cur := root
	for _, s := range segs {
		// Parse name and optional key value
        name := s
        var key string
        if i := strings.IndexByte(s, '#'); i >= 0 {
            name = s[:i]
        }
        if i := strings.IndexByte(s, '['); i >= 0 {
            name = s[:i]
            // Extract value between quotes in [name="..."]
            start := strings.IndexByte(s[i:], '"')
            end := strings.LastIndexByte(s, '"')
            if start >= 0 && end > i+start {
                // adjust start to absolute index
                start = i + start
                key = s[start+1 : end]
            }
        }

        // Find or create the child element for this name
        var child *PatchNode
        for _, c := range cur.Children {
            if c.XMLName.Local == name {
                child = c
                break
            }
        }
        if child == nil {
            child = &PatchNode{XMLName: xml.Name{Local: name}}
            cur.Children = append(cur.Children, child)
        }

        // If we have a [name="..."], ensure the child has a <name> first-child
        if key != "" {
            hasName := false
            for _, c := range child.Children {
                if c.XMLName.Local == "name" {
                    hasName = true
                    break
                }
            }
            if !hasName {
                nameNode := &PatchNode{
                    XMLName: xml.Name{Local: "name"},
                    Text:    key,
                }
                child.Children = append([]*PatchNode{nameNode}, child.Children...)
            }
        }
		cur = child
	}
	return cur
}

// ------------------------- Helpers [END] -------------------------