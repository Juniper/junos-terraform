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
func ValueToConfig(v tftypes.Value, nodes []patch.SchemaNode) (*patch.Node, error) {
	root := &patch.Node{Tag: "configuration"}
	if v.IsNull() {
		return root, nil
	}
	var attrs map[string]tftypes.Value
	if err := v.As(&attrs); err != nil {
		return nil, fmt.Errorf("configuration: %w", err)
	}
	if err := fillElement(root, attrs, nodes, ""); err != nil {
		return nil, err
	}
	return root, nil
}

func fillElement(el *patch.Node, attrs map[string]tftypes.Value, nodes []patch.SchemaNode, key string) error {
	for _, n := range keysFirst(attributeNodes(nodes), key) {
		v, ok := attrs[SanitizeName(n.Name)]
		if !ok || v.IsNull() {
			continue
		}
		if !v.IsKnown() {
			return fmt.Errorf("%s: value is not known", n.Name)
		}
		switch n.Type {
		case "leaf":
			var s string
			if err := v.As(&s); err != nil {
				return fmt.Errorf("%s: %w", n.Name, err)
			}
			appendChild(el, &patch.Node{Tag: n.Name, Text: s})
		case "leaf-list":
			var items []tftypes.Value
			if err := v.As(&items); err != nil {
				return fmt.Errorf("%s: %w", n.Name, err)
			}
			for _, item := range items {
				var s string
				if err := item.As(&s); err != nil {
					return fmt.Errorf("%s: %w", n.Name, err)
				}
				appendChild(el, &patch.Node{Tag: n.Name, Text: s})
			}
		case "container", "list":
			var items []tftypes.Value
			if err := v.As(&items); err != nil {
				return fmt.Errorf("%s: %w", n.Name, err)
			}
			for _, item := range items {
				var itemAttrs map[string]tftypes.Value
				if err := item.As(&itemAttrs); err != nil {
					return fmt.Errorf("%s: %w", n.Name, err)
				}
				child := &patch.Node{Tag: n.Name}
				if err := fillElement(child, itemAttrs, n.Children, n.Key); err != nil {
					return fmt.Errorf("%s/%w", n.Name, err)
				}
				appendChild(el, child)
			}
		}
	}
	return nil
}

// keysFirst returns nodes with a list's key leaves first, in key order.
func keysFirst(nodes []patch.SchemaNode, key string) []patch.SchemaNode {
	keys := strings.Fields(key)
	if len(keys) == 0 {
		return nodes
	}
	out := make([]patch.SchemaNode, 0, len(nodes))
	for _, k := range keys {
		for _, n := range nodes {
			if n.Name == k {
				out = append(out, n)
			}
		}
	}
	for _, n := range nodes {
		if !slices.Contains(keys, n.Name) {
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
func ConfigToValue(root *patch.Node, nodes []patch.SchemaNode, typ tftypes.Object) (tftypes.Value, error) {
	return elementToObject(root, nodes, typ)
}

func elementToObject(el *patch.Node, nodes []patch.SchemaNode, typ tftypes.Object) (tftypes.Value, error) {
	byTag := make(map[string][]*patch.Node, len(el.Children))
	for _, c := range el.Children {
		byTag[c.Tag] = append(byTag[c.Tag], c)
	}

	vals := make(map[string]tftypes.Value, len(typ.AttributeTypes))
	for _, n := range attributeNodes(nodes) {
		name := SanitizeName(n.Name)
		at, ok := typ.AttributeTypes[name]
		if !ok {
			return tftypes.Value{}, fmt.Errorf("%s: no attribute %s in the resource type", n.Name, name)
		}
		els := byTag[n.Name]
		if len(els) == 0 {
			vals[name] = tftypes.NewValue(at, nil)
			continue
		}
		switch n.Type {
		case "leaf":
			vals[name] = tftypes.NewValue(tftypes.String, els[0].Text)
		case "leaf-list":
			items := make([]tftypes.Value, len(els))
			for i, e := range els {
				items[i] = tftypes.NewValue(tftypes.String, e.Text)
			}
			vals[name] = tftypes.NewValue(at, items)
		case "container", "list":
			lt, ok := at.(tftypes.List)
			if !ok {
				return tftypes.Value{}, fmt.Errorf("%s: attribute type %s is not a list", n.Name, at)
			}
			ot, ok := lt.ElementType.(tftypes.Object)
			if !ok {
				return tftypes.Value{}, fmt.Errorf("%s: list element type %s is not an object", n.Name, lt.ElementType)
			}
			items := make([]tftypes.Value, len(els))
			for i, e := range els {
				item, err := elementToObject(e, n.Children, ot)
				if err != nil {
					return tftypes.Value{}, fmt.Errorf("%s/%w", n.Name, err)
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
