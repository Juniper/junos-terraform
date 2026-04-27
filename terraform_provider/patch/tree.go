package patch

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"io"
)

// BuildTree parses XML bytes into a *Node tree using a streaming token decoder.
// It handles the XML declaration header produced by xml.Header gracefully.
func BuildTree(xmlBytes []byte) (*Node, error) {
	decoder := xml.NewDecoder(bytes.NewReader(xmlBytes))

	var stack []*Node
	var root *Node

	for {
		tok, err := decoder.Token()
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, fmt.Errorf("failed parsing XML token stream: %w", err)
		}

		switch t := tok.(type) {

		case xml.StartElement:
			node := &Node{
				Tag:   t.Name.Local,
				Attrs: make(map[string]string),
			}
			for _, attr := range t.Attr {
				// Skip xmlns declarations — not needed in our tree
				if attr.Name.Space == "xmlns" || attr.Name.Local == "xmlns" {
					continue
				}
				node.Attrs[attr.Name.Local] = attr.Value
			}
			if len(stack) > 0 {
				parent := stack[len(stack)-1]
				node.Parent = parent
				parent.Children = append(parent.Children, node)
			}
			stack = append(stack, node)

		case xml.EndElement:
			if len(stack) == 0 {
				return nil, fmt.Errorf("unexpected end element </%s>", t.Name.Local)
			}
			popped := stack[len(stack)-1]
			stack = stack[:len(stack)-1]
			// When we pop the last element off the stack it is the root
			if len(stack) == 0 {
				root = popped
			}

		case xml.CharData:
			if len(stack) > 0 {
				text := string(bytes.TrimSpace([]byte(t)))
				if text != "" {
					top := stack[len(stack)-1]
					top.Text += text
				}
			}
		}
	}

	if root == nil {
		return nil, fmt.Errorf("no root element found in XML")
	}
	return root, nil
}
