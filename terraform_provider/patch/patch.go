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
	return CreateDiffPatchWithSchema(diffMap, nil, groupName, nil)
}

// CreateDiffPatchWithSchema is like CreateDiffPatch but accepts the plan's
// leaf map and a schema index for container delete coalescing. When both are
// non-nil and a schema container is being removed entirely — the plan keeps
// nothing under it — its leaf deletes are coalesced into a single
// container-level nc:operation="delete".
func CreateDiffPatchWithSchema(diffMap map[string]Change, planMap map[string]string, groupName string, idx map[string]*NodeInfo) ([]byte, error) {
	_ = groupName

	// Pre-pass: coalesce container deletes when schema and plan are available.
	if idx != nil && planMap != nil {
		diffMap = coalesceContainerDeletes(diffMap, planMap, idx)
	}

	diffMap = mergeContainerChanges(diffMap)

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
		// Positional leaf-list entries (path ends with [pos=N]) represent
		// ordered-by-user leaf-lists. A Replace means the value at that
		// position changed — emit delete of old + create of new.
		if p.keyName == "pos" {
			switch p.change.Op {
			case Create:
				leaf := &Node{Tag: p.tag, Parent: p.parent, Operation: "create", Text: p.change.NewVal}
				p.parent.Children = append(p.parent.Children, leaf)
			case Delete:
				leaf := &Node{Tag: p.tag, Parent: p.parent, Operation: "delete", Text: p.change.OldVal}
				p.parent.Children = append(p.parent.Children, leaf)
			case Replace:
				// Reorder: delete old value, create new value
				del := &Node{Tag: p.tag, Parent: p.parent, Operation: "delete", Text: p.change.OldVal}
				p.parent.Children = append(p.parent.Children, del)
				cre := &Node{Tag: p.tag, Parent: p.parent, Operation: "create", Text: p.change.NewVal}
				p.parent.Children = append(p.parent.Children, cre)
			}
			continue
		}

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

	_, parentKeys := parseSegmentKeys(parentSegments[len(parentSegments)-1])
	leafTag, _, _ := parseSegment(leafSegment)

	keyValue := change.NewVal
	if change.Op == Delete {
		keyValue = change.OldVal
	}

	// The leaf is one of the parent entry's keys, with the entry's value for
	// it. A single key must have a value; in a compound key one may be empty
	// (choice-value of route-filter ... exact), and is still the entry's key.
	isKey := false
	for _, k := range parentKeys {
		if k.name == leafTag && k.value == keyValue && (keyValue != "" || len(parentKeys) > 1) {
			isKey = true
			break
		}
	}
	if !isKey {
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

// coalesceContainerDeletes replaces the leaf deletes under a schema container
// that is being removed entirely with a single synthetic Delete of the
// container itself. This produces a compact <container nc:operation="delete"/>
// instead of N individual leaf deletes, which is both more efficient and avoids
// ordering issues on Junos.
//
// A container instance is removed entirely when the plan keeps no leaf under
// it; every change under it is then a Delete. The diff alone cannot tell: an
// unchanged sibling is not in it, and a container holding one deleted list
// entry and one kept one would otherwise be deleted with both. Containers are
// identified by their instance path, keys included, so a delete of
// policy-statement[name=a]/term[name=t]/from does not become one of every
// term's from. Each leaf counts toward its deepest container ancestor; when a
// container and one inside it are both removed, only the outer delete is kept.
func coalesceContainerDeletes(diffMap map[string]Change, planMap map[string]string, idx map[string]*NodeInfo) map[string]Change {
	// Leaf paths grouped by their deepest container ancestor instance.
	// A "container" here means a schema node of KindContainer (not KindList).
	containers := make(map[string][]string)

	for path := range diffMap {
		segments := splitPathRespectingQuotes(path)

		// Find the deepest container ancestor in the schema. segments[0] is
		// the configuration root and never a candidate.
		for depth := len(segments) - 1; depth >= 2; depth-- {
			instance := strings.Join(segments[:depth], "/")
			info, ok := idx[outputPathToSchemaPath(instance)]
			if !ok || info.Kind != KindContainer {
				continue
			}
			containers[instance] = append(containers[instance], path)
			break // only use the deepest container
		}
	}

	removed := make(map[string]bool)
	for instance, paths := range containers {
		if len(paths) >= 2 && !planHasPathUnder(planMap, instance) {
			removed[instance] = true
		}
	}

	if len(removed) == 0 {
		return diffMap
	}

	// Build new diffMap: remove the coalesced leaves and add a delete for each
	// removed container not inside another removed container.
	result := make(map[string]Change, len(diffMap))
	for path, change := range diffMap {
		if !underAny(path, removed) {
			result[path] = change
		}
	}
	for instance := range removed {
		if !underAny(instance, removed) {
			result[instance] = Change{Op: Delete, OldVal: "", NewVal: ""}
		}
	}

	return result
}

// mergeContainerChanges folds a changed presence container and the changes
// under it into one element. Only a presence container is in the leaf map
// with leaves under it, so a changed path with changes under it is one. A
// deleted container's own delete removes its children, so their deletes are
// dropped; a created one is created by creating its children, so its own
// create is dropped. Kept, each would add a second element beside the one
// holding the children.
func mergeContainerChanges(diffMap map[string]Change) map[string]Change {
	var deleted, created []string
	for path, change := range diffMap {
		if !hasChangeUnder(diffMap, path) {
			continue
		}
		switch change.Op {
		case Delete:
			deleted = append(deleted, path)
		case Create:
			created = append(created, path)
		}
	}
	if len(deleted)+len(created) == 0 {
		return diffMap
	}

	result := make(map[string]Change, len(diffMap))
	for path, change := range diffMap {
		result[path] = change
	}
	for _, path := range created {
		delete(result, path)
	}
	for _, container := range deleted {
		for path := range diffMap {
			if strings.HasPrefix(path, container+"/") {
				delete(result, path)
			}
		}
	}
	return result
}

// hasChangeUnder reports whether the diff changes anything below path.
func hasChangeUnder(diffMap map[string]Change, path string) bool {
	prefix := path + "/"
	for p := range diffMap {
		if strings.HasPrefix(p, prefix) {
			return true
		}
	}
	return false
}

// planHasPathUnder reports whether the plan keeps any leaf at or below path.
func planHasPathUnder(planMap map[string]string, path string) bool {
	prefix := path + "/"
	for p := range planMap {
		if p == path || strings.HasPrefix(p, prefix) {
			return true
		}
	}
	return false
}

// underAny reports whether path lies strictly below one of the containers.
func underAny(path string, containers map[string]bool) bool {
	for c := range containers {
		if strings.HasPrefix(path, c+"/") {
			return true
		}
	}
	return false
}
