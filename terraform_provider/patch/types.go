package patch

// ChangeType represents the CRUD operation for a single leaf diff.
type ChangeType int

const (
	Create  ChangeType = iota
	Replace            // value exists in both state and plan but differs
	Delete             // value exists in state but not in plan
)

// Node is a single element in the parsed XML tree.
type Node struct {
	Tag       string
	Attrs     map[string]string // standard XML attributes (not nc:operation)
	Operation string            // nc:operation value: "create", "replace", "delete", or ""
	Text      string            // character data between open/close tags
	Children  []*Node
	Parent    *Node
}

// Change holds a single leaf-level diff entry produced by ComputeDiff.
type Change struct {
	Op     ChangeType
	OldVal string // empty for Create
	NewVal string // empty for Delete
}
