package generic

import (
	"bytes"
	"encoding/xml"

	"terraform_provider/patch"
)

// SerializeTree converts a *patch.Node tree back to XML bytes (inverse of patch.BuildTree).
// Produces well-formed XML with a declaration header.
func SerializeTree(root *patch.Node) ([]byte, error) {
	var buf bytes.Buffer
	buf.WriteString(xml.Header)
	enc := xml.NewEncoder(&buf)
	enc.Indent("", "  ")
	if err := encodeNode(enc, root); err != nil {
		return nil, err
	}
	if err := enc.Flush(); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

// SerializeTreeNoHeader serializes without the XML declaration header.
func SerializeTreeNoHeader(root *patch.Node) ([]byte, error) {
	var buf bytes.Buffer
	enc := xml.NewEncoder(&buf)
	enc.Indent("", "  ")
	if err := encodeNode(enc, root); err != nil {
		return nil, err
	}
	if err := enc.Flush(); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func encodeNode(enc *xml.Encoder, node *patch.Node) error {
	start := xml.StartElement{
		Name: xml.Name{Local: node.Tag},
	}
	for k, v := range node.Attrs {
		start.Attr = append(start.Attr, xml.Attr{
			Name:  xml.Name{Local: k},
			Value: v,
		})
	}

	if len(node.Children) == 0 && node.Text == "" {
		// Self-closing empty element
		if err := enc.EncodeToken(start); err != nil {
			return err
		}
		return enc.EncodeToken(start.End())
	}

	if err := enc.EncodeToken(start); err != nil {
		return err
	}

	if node.Text != "" {
		if err := enc.EncodeToken(xml.CharData([]byte(node.Text))); err != nil {
			return err
		}
	}

	for _, child := range node.Children {
		if err := encodeNode(enc, child); err != nil {
			return err
		}
	}

	return enc.EncodeToken(start.End())
}
