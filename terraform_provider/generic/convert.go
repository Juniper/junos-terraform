package generic

import (
	"fmt"
	"slices"
	"strings"

	"terraform_provider/patch"

	"github.com/hashicorp/terraform-plugin-go/tftypes"
)

// ValueToConfig converts the resource's Terraform value into a <configuration>
// tree: each non-null attribute becomes an element named by its YANG node, in
// schema order. A leaf is one element (an empty string gives an empty element,
// as for a YANG empty leaf), a leaf-list one element per value, and a
// container or list one element per entry. A list entry's key leaves come
// first, in key order, as Junos requires.
func ValueToConfig(v tftypes.Value, s *patch.Schema) (*patch.Node, error) {
	root := &patch.Node{Tag: "configuration"}
	if v.IsNull() {
		return root, nil
	}
	var attrs map[string]tftypes.Value
	if err := v.As(&attrs); err != nil {
		return nil, fmt.Errorf("configuration: %w", err)
	}
	if err := fillElement(root, attrs, s, 0); err != nil {
		return nil, err
	}
	return root, nil
}

func fillElement(el *patch.Node, attrs map[string]tftypes.Value, s *patch.Schema, parent patch.SchemaNodeID) error {
	for _, n := range keysFirst(s, attributeNodes(s, parent), s.Key(parent)) {
		name := s.Name(n)
		v, ok := attrs[SanitizeName(name)]
		if !ok || v.IsNull() {
			continue
		}
		if !v.IsKnown() {
			return fmt.Errorf("%s: value is not known", name)
		}
		switch s.Kind(n) {
		case patch.KindLeaf:
			var text string
			if err := v.As(&text); err != nil {
				return fmt.Errorf("%s: %w", name, err)
			}
			appendChild(el, &patch.Node{Tag: name, Text: text})
		case patch.KindLeafList:
			var items []tftypes.Value
			if err := v.As(&items); err != nil {
				return fmt.Errorf("%s: %w", name, err)
			}
			for _, item := range items {
				var text string
				if err := item.As(&text); err != nil {
					return fmt.Errorf("%s: %w", name, err)
				}
				appendChild(el, &patch.Node{Tag: name, Text: text})
			}
		case patch.KindContainer, patch.KindList:
			var items []tftypes.Value
			if err := v.As(&items); err != nil {
				return fmt.Errorf("%s: %w", name, err)
			}
			for _, item := range items {
				var itemAttrs map[string]tftypes.Value
				if err := item.As(&itemAttrs); err != nil {
					return fmt.Errorf("%s: %w", name, err)
				}
				child := &patch.Node{Tag: name}
				if err := fillElement(child, itemAttrs, s, n); err != nil {
					return fmt.Errorf("%s/%w", name, err)
				}
				appendChild(el, child)
			}
		}
	}
	return nil
}

// keysFirst returns nodes with a list's key leaves first, in key order.
func keysFirst(s *patch.Schema, nodes []patch.SchemaNodeID, key string) []patch.SchemaNodeID {
	keys := strings.Fields(key)
	if len(keys) == 0 {
		return nodes
	}
	out := make([]patch.SchemaNodeID, 0, len(nodes))
	for _, k := range keys {
		for _, n := range nodes {
			if s.Name(n) == k {
				out = append(out, n)
			}
		}
	}
	for _, n := range nodes {
		if !slices.Contains(keys, s.Name(n)) {
			out = append(out, n)
		}
	}
	return out
}

// ConfigToValue converts a <configuration> tree into the resource's Terraform
// value of type typ. Only elements in the schema are read, so anything the
// provider does not model is ignored. A missing element is null; an element's
// text is a leaf's value (an empty element is ""); repeated elements are the
// entries of a leaf-list or list, in document order; a container is a list of
// one entry. Attributes that are not configuration (resource_name) are null.
func ConfigToValue(root *patch.Node, s *patch.Schema, typ tftypes.Object) (tftypes.Value, error) {
	return elementToObject(root, s, 0, typ)
}

func elementToObject(el *patch.Node, s *patch.Schema, parent patch.SchemaNodeID, typ tftypes.Object) (tftypes.Value, error) {
	byTag := make(map[string][]*patch.Node, len(el.Children))
	for _, c := range el.Children {
		byTag[c.Tag] = append(byTag[c.Tag], c)
	}

	vals := make(map[string]tftypes.Value, len(typ.AttributeTypes))
	for _, n := range attributeNodes(s, parent) {
		tag := s.Name(n)
		name := SanitizeName(tag)
		at, ok := typ.AttributeTypes[name]
		if !ok {
			return tftypes.Value{}, fmt.Errorf("%s: no attribute %s in the resource type", tag, name)
		}
		els := byTag[tag]
		if len(els) == 0 {
			vals[name] = tftypes.NewValue(at, nil)
			continue
		}
		switch s.Kind(n) {
		case patch.KindLeaf:
			vals[name] = tftypes.NewValue(tftypes.String, els[0].Text)
		case patch.KindLeafList:
			items := make([]tftypes.Value, len(els))
			for i, e := range els {
				items[i] = tftypes.NewValue(tftypes.String, e.Text)
			}
			vals[name] = tftypes.NewValue(at, items)
		case patch.KindContainer, patch.KindList:
			lt, ok := at.(tftypes.List)
			if !ok {
				return tftypes.Value{}, fmt.Errorf("%s: attribute type %s is not a list", tag, at)
			}
			ot, ok := lt.ElementType.(tftypes.Object)
			if !ok {
				return tftypes.Value{}, fmt.Errorf("%s: list element type %s is not an object", tag, lt.ElementType)
			}
			items := make([]tftypes.Value, len(els))
			for i, e := range els {
				item, err := elementToObject(e, s, n, ot)
				if err != nil {
					return tftypes.Value{}, fmt.Errorf("%s/%w", tag, err)
				}
				items[i] = item
			}
			vals[name] = tftypes.NewValue(at, items)
		}
	}

	for name, at := range typ.AttributeTypes {
		if _, ok := vals[name]; !ok {
			vals[name] = tftypes.NewValue(at, nil)
		}
	}
	return tftypes.NewValue(typ, vals), nil
}

func appendChild(parent, child *patch.Node) {
	child.Parent = parent
	parent.Children = append(parent.Children, child)
}
