package patch

import (
	"strings"
	"sync"
)

// Schema is a compiled configuration schema: one fixed-size record per node,
// with each distinct name stored once. It answers the patch engine's lookups
// by schema path and gives the generic provider the node tree, in schema
// order, without keeping a path string or a struct with pointers per node.
//
// Nodes are laid out breadth first from the <configuration> node (ID 0), so a
// node's children are consecutive records, located by the first one's ID and
// their number. YANG choice and case nodes are flattened into their parent,
// and siblings with the same name are merged, as in the configuration XML.
type Schema struct {
	nodes []schemaRecord
	names []string // names[0] is ""

	nameIndexOnce sync.Once
	nameIndex     map[string]uint32
}

// SchemaNodeID identifies a node in a Schema. The <configuration> node is 0.
type SchemaNodeID uint32

type schemaRecord struct {
	name        uint32 // index in names
	key         uint32 // index in names of a list's key ("" if none)
	firstChild  uint32
	numChildren uint32
	kind        NodeKind
	flags       uint8
}

const (
	flagPresence uint8 = 1 << iota
	flagOrderedByUser
)

// NodeInfo is what the patch engine needs to know about a schema node.
type NodeInfo struct {
	Name string
	Kind NodeKind
	// ListKey is a list's key leaf, or its key leaves separated by spaces.
	ListKey string
	// OrderedByUser: the order of the list's or leaf-list's entries is
	// meaningful.
	OrderedByUser bool
	// Presence: the container's existence is configuration.
	Presence bool
}

// Len returns the number of nodes, <configuration> included.
func (s *Schema) Len() int { return len(s.nodes) }

// Name returns a node's YANG name.
func (s *Schema) Name(id SchemaNodeID) string { return s.names[s.nodes[id].name] }

// Kind returns a node's kind.
func (s *Schema) Kind(id SchemaNodeID) NodeKind { return s.nodes[id].kind }

// Key returns a list's key leaf, or its key leaves separated by spaces.
func (s *Schema) Key(id SchemaNodeID) string { return s.names[s.nodes[id].key] }

// NumChildren returns the number of a node's children.
func (s *Schema) NumChildren(id SchemaNodeID) int { return int(s.nodes[id].numChildren) }

// Child returns a node's i-th child, in schema order.
func (s *Schema) Child(id SchemaNodeID, i int) SchemaNodeID {
	return SchemaNodeID(s.nodes[id].firstChild + uint32(i))
}

// Info returns what the patch engine needs to know about a node.
func (s *Schema) Info(id SchemaNodeID) NodeInfo {
	r := s.nodes[id]
	return NodeInfo{
		Name:          s.names[r.name],
		Kind:          r.kind,
		ListKey:       s.names[r.key],
		OrderedByUser: r.flags&flagOrderedByUser != 0,
		Presence:      r.flags&flagPresence != 0,
	}
}

// Lookup returns the node at a schema path relative to <configuration>, such
// as "system/services/ssh", with no list predicates. A nil Schema has no
// nodes.
func (s *Schema) Lookup(path string) (NodeInfo, bool) {
	id, ok := s.find(path)
	if !ok {
		return NodeInfo{}, false
	}
	return s.Info(id), true
}

func (s *Schema) find(path string) (SchemaNodeID, bool) {
	if s == nil || path == "" {
		return 0, false
	}
	s.nameIndexOnce.Do(func() {
		s.nameIndex = make(map[string]uint32, len(s.names))
		for i, n := range s.names {
			s.nameIndex[n] = uint32(i)
		}
	})
	id := SchemaNodeID(0)
	for seg := range strings.SplitSeq(path, "/") {
		name, ok := s.nameIndex[seg]
		if !ok || seg == "" {
			return 0, false
		}
		r := s.nodes[id]
		found := false
		for c := r.firstChild; c < r.firstChild+r.numChildren; c++ {
			if s.nodes[c].name == name {
				id, found = SchemaNodeID(c), true
				break
			}
		}
		if !found {
			return 0, false
		}
	}
	return id, true
}

// CompileSchema compiles schema nodes, as the pyang plugin writes them under
// its root, into a Schema. The nodes are the <configuration> node, or the
// nodes under it.
func CompileSchema(roots []SchemaNode) *Schema {
	config := &SchemaNode{Name: "configuration", Type: "container", Children: roots}
	if len(roots) == 1 && roots[0].Name == "configuration" {
		config = &roots[0]
	}

	b := schemaBuilder{nameIDs: map[string]uint32{"": 0}, s: &Schema{names: []string{""}}}
	b.s.nodes = append(b.s.nodes, schemaRecord{name: b.name("configuration"), kind: KindContainer})
	// Breadth first: each queued node's children are appended together. A
	// queued node is one or more same-named siblings, merged.
	queue := [][]*SchemaNode{{config}}
	var children []*SchemaNode
	for i := 0; i < len(queue); i++ {
		children = children[:0]
		for _, n := range queue[i] {
			children = flattenChoices(children, n.Children)
		}
		groups := mergeSiblings(children)
		r := &b.s.nodes[i]
		r.firstChild = uint32(len(b.s.nodes))
		r.numChildren = uint32(len(groups))
		for _, g := range groups {
			b.s.nodes = append(b.s.nodes, b.record(g))
			queue = append(queue, g)
		}
	}
	return b.s
}

// flattenChoices appends nodes to out, with every YANG choice and case node
// replaced by its children.
func flattenChoices(out []*SchemaNode, nodes []SchemaNode) []*SchemaNode {
	for i := range nodes {
		n := &nodes[i]
		if n.Type == "choice" || n.Type == "case" {
			out = flattenChoices(out, n.Children)
			continue
		}
		if n.Name != "" {
			out = append(out, n)
		}
	}
	return out
}

// mergeSiblings groups siblings with the same name, which the configuration
// XML cannot tell apart (the same node in two cases of a choice), in the order
// of their first appearance.
func mergeSiblings(nodes []*SchemaNode) [][]*SchemaNode {
	groups := make([][]*SchemaNode, 0, len(nodes))
	var at map[string]int
	for _, n := range nodes {
		i := -1
		if len(nodes) <= 16 {
			for j, g := range groups {
				if g[0].Name == n.Name {
					i = j
					break
				}
			}
		} else {
			if at == nil {
				at = make(map[string]int, len(nodes))
			}
			if j, ok := at[n.Name]; ok {
				i = j
			} else {
				at[n.Name] = len(groups)
			}
		}
		if i < 0 {
			groups = append(groups, []*SchemaNode{n})
		} else {
			groups[i] = append(groups[i], n)
		}
	}
	return groups
}

type schemaBuilder struct {
	s       *Schema
	nameIDs map[string]uint32
}

func (b *schemaBuilder) name(n string) uint32 {
	id, ok := b.nameIDs[n]
	if !ok {
		id = uint32(len(b.s.names))
		b.s.names = append(b.s.names, n)
		b.nameIDs[n] = id
	}
	return id
}

// record returns the record for a node, or for same-named siblings merged: the
// last one's kind and key, and any one's flags.
func (b *schemaBuilder) record(group []*SchemaNode) schemaRecord {
	r := schemaRecord{name: b.name(group[0].Name)}
	for _, n := range group {
		switch n.Type {
		case "list":
			r.kind = KindList
			r.key = b.name(n.Key)
		case "leaf":
			r.kind = KindLeaf
		case "leaf-list":
			r.kind = KindLeafList
		default:
			// container, or a type the patch engine does not know: traversed
			// like a container.
			r.kind = KindContainer
		}
		if n.Presence != "" {
			r.flags |= flagPresence
		}
		if n.OrderedBy == "user" {
			r.flags |= flagOrderedByUser
		}
	}
	if r.kind != KindContainer {
		r.flags &^= flagPresence
	}
	return r
}
