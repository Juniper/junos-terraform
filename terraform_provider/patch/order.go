package patch

import "sort"

// AlignXMLOrderToReference reorders current XML siblings to follow the order in
// the reference XML where possible, while keeping a deterministic fallback
// ordering for entries absent from the reference.
func AlignXMLOrderToReference(currentXML []byte, referenceXML []byte, idx map[string]*NodeInfo) ([]byte, error) {
	currentTree, err := BuildTree(currentXML)
	if err != nil {
		return nil, err
	}

	var referenceTree *Node
	if len(referenceXML) > 0 {
		referenceTree, err = BuildTree(referenceXML)
		if err != nil {
			return nil, err
		}
	}

	alignNodeOrder(currentTree, referenceSiblingOrders(referenceTree, idx), idx, currentTree.Tag)
	return marshalNodeTree(currentTree)
}

func alignNodeOrder(node *Node, ref map[string]map[string]int, idx map[string]*NodeInfo, instancePath string) {
	if len(node.Children) == 0 {
		return
	}

	referenceOrder := ref[instancePath]

	sort.SliceStable(node.Children, func(i, j int) bool {
		left := childSortKey(node.Children[i], instancePath, referenceOrder, idx)
		right := childSortKey(node.Children[j], instancePath, referenceOrder, idx)

		if left.hasReference != right.hasReference {
			return left.hasReference
		}
		if left.referenceRank != right.referenceRank {
			return left.referenceRank < right.referenceRank
		}
		if left.tag != right.tag {
			return left.tag < right.tag
		}
		if left.identity != right.identity {
			return left.identity < right.identity
		}
		return left.text < right.text
	})

	for _, child := range node.Children {
		childInstance := instanceIdentity(child, instancePath, idx)
		alignNodeOrder(child, ref, idx, childInstance)
	}
}

type sortKey struct {
	hasReference  bool
	referenceRank int
	tag           string
	identity      string
	text          string
}

func childSortKey(child *Node, parentInstancePath string, referenceOrder map[string]int, idx map[string]*NodeInfo) sortKey {
	identity := nodeIdentity(child, parentInstancePath, idx)
	rank, ok := referenceOrder[identity]
	if !ok {
		rank = 1 << 30
	}

	return sortKey{
		hasReference:  ok,
		referenceRank: rank,
		tag:           child.Tag,
		identity:      identity,
		text:          child.Text,
	}
}

func referenceSiblingOrders(root *Node, idx map[string]*NodeInfo) map[string]map[string]int {
	orders := make(map[string]map[string]int)
	if root == nil {
		return orders
	}

	var walk func(node *Node, instancePath string)
	walk = func(node *Node, instancePath string) {
		if _, ok := orders[instancePath]; !ok {
			orders[instancePath] = make(map[string]int)
		}
		for i, child := range node.Children {
			identity := nodeIdentity(child, instancePath, idx)
			if _, seen := orders[instancePath][identity]; !seen {
				orders[instancePath][identity] = i
			}
			childInstance := instanceIdentity(child, instancePath, idx)
			walk(child, childInstance)
		}
	}

	walk(root, root.Tag)
	return orders
}

// nodeIdentity returns a short identity string for sorting siblings under the
// same parent.  It does NOT include the parent path.
func nodeIdentity(node *Node, parentInstancePath string, idx map[string]*NodeInfo) string {
	schPath := schemaPathFromInstance(parentInstancePath, node.Tag)
	info := idx[schPath]
	if info == nil {
		if node.Text != "" {
			return node.Tag + "=" + node.Text
		}
		return node.Tag
	}

	switch info.Kind {
	case KindList:
		keyVal := findKeyChild(node, info.ListKey)
		if keyVal != "" {
			return node.Tag + "[" + info.ListKey + "=" + keyVal + "]"
		}
	case KindLeafList:
		return node.Tag + "[value=" + node.Text + "]"
	}

	if node.Text != "" {
		return node.Tag + "=" + node.Text
	}
	return node.Tag
}

// instanceIdentity returns a full instance-aware path for a child node,
// including keyed-list identity so that different list entries get distinct
// ordering buckets.
func instanceIdentity(child *Node, parentInstancePath string, idx map[string]*NodeInfo) string {
	base := parentInstancePath + "/" + child.Tag
	schPath := schemaPathFromInstance(parentInstancePath, child.Tag)
	info := idx[schPath]
	if info != nil && info.Kind == KindList && info.ListKey != "" {
		keyVal := findKeyChild(child, info.ListKey)
		if keyVal != "" {
			return base + "[" + info.ListKey + "=" + keyVal + "]"
		}
	}
	return base
}

// schemaPathFromInstance extracts the schema path for a child tag given an
// instance-aware parent path.  It strips keyed-list predicates like
// "host[name=log]" back to "host" so the result matches schema index keys.
func schemaPathFromInstance(parentInstancePath string, childTag string) string {
	return schemaPath(joinXMLPath(stripInstancePredicates(parentInstancePath), childTag))
}

// stripInstancePredicates removes [key=value] predicates from every segment
// of an instance path, returning the plain XML element path.
func stripInstancePredicates(path string) string {
	var out []byte
	inBracket := false
	for i := 0; i < len(path); i++ {
		if path[i] == '[' {
			inBracket = true
			continue
		}
		if path[i] == ']' {
			inBracket = false
			continue
		}
		if !inBracket {
			out = append(out, path[i])
		}
	}
	return string(out)
}

func joinXMLPath(parentPath string, tag string) string {
	if parentPath == "" {
		return tag
	}
	return parentPath + "/" + tag
}

func schemaPath(path string) string {
	path = normalizePath(path)
	if path == "configuration" {
		return ""
	}
	return normalizePath(path)
}
