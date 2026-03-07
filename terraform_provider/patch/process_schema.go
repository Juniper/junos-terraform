package patch

import (
    "encoding/json"
	"strings"
)

// ------------------------- Type Definitions [START] -------------------------

type NodeKind uint8

const (
	KindContainer NodeKind = iota
	KindList
	KindLeaf
)

type LeafBase uint8

const (
	LeafString LeafBase = iota
	LeafBool
	LeafInt
	LeafUint
	LeafEnum
	LeafUnion
	LeafOther
)

// Raw JSON node
type SchemaNode struct {
	Name     string       `json:"name"`
	Type     string       `json:"type"`      // container | list | leaf
	Path     string       `json:"path"`      // often parent path
	Key      string       `json:"key"`       // list key
	LeafType string       `json:"leaf-type"` // leaf-only: string, union, etc.

    Children []SchemaNode `json:"children"` 
    // Union branches (when leaf-type == "union")
	Types []UnionType `json:"types"`
    // Constraints (sometimes on leaves, sometimes inside union branches)
	Lengths  []LenRange  `json:"lengths"`
	Ranges   []NumRange  `json:"ranges"`
	Patterns []string    `json:"patterns"`
	Enums    []EnumValue `json:"enums"` // if your trimmed schema ever includes enums
}

type UnionType struct {
	Type     string     `json:"type"`
	Path     string     `json:"path"`
	Patterns []string   `json:"patterns"`
	Ranges   []NumRange `json:"ranges"`
	Lengths  []LenRange `json:"lengths"`
	Enums    []EnumValue `json:"enums"`
}

type NumRange struct {
	Min  *float64 `json:"min"`
	Max  *float64 `json:"max"`
	Path string   `json:"path"`
}

type LenRange struct {
	Min  *int   `json:"min"`
	Max  *int   `json:"max"`
	Path string `json:"path"`
}

type EnumValue struct {
	Name  string `json:"name"`
	Value any `json:"value"`
}

type NodeInfo struct {
	Path   string
	Name   string
	Kind   NodeKind
	Parent string
	Children []string
	// Lists
	ListKey     string
	ListKeyPath string
	// Leaves
	Leaf LeafInfo
}

type LeafInfo struct {
	Base LeafBase
	Union []UnionBranch
	Patterns []string
	Ranges   []NumRange
	Lengths  []LenRange
	Enums map[string]struct{} // canonical set of enum names (if present)
}

type UnionBranch struct {
	Base     LeafBase
	Patterns []string
	Ranges   []NumRange
	Lengths  []LenRange
	Enums    map[string]struct{}
}

type TrimmedSchemaWrapper struct {
	Path string     `json:"path"`
	Root SchemaRoot `json:"root"`
}

type SchemaRoot struct {
	Children []SchemaNode `json:"children"`
}

// ------------------------- Type Definitions [END] -------------------------

// ------------------------- Process Trimmed Schema [START] -------------------------

// Go raw string literal -> compiled index
func UnmarshalTrimmedSchemaIndex(trimmedSchemaJSON string) (map[string]*NodeInfo, error) {

	var w TrimmedSchemaWrapper
	if err := json.Unmarshal([]byte(trimmedSchemaJSON), &w); err != nil {
		return nil, err
	}

	roots := w.Root.Children

	idx := make(map[string]*NodeInfo)

	// Walk and compile
	var walk func(n SchemaNode, parentFull string)

	walk = func(n SchemaNode, parentFull string) {
		full := canonicalFullPath(n, parentFull)
		if full == "" {
			// still walk children (some schemas have virtual root nodes)
			for _, c := range n.Children {
				walk(c, parentFull)
			}
			return
		}

		info := idx[full]
		if info == nil {
			info = &NodeInfo{
				Path:   full,
				Name:   n.Name,
				Parent: parentFull,
			}
			idx[full] = info
		} else {
			// if collisions happen, keep existing and merge below
			if info.Name == "" {
				info.Name = n.Name
			}
			if info.Parent == "" {
				info.Parent = parentFull
			}
		}

		// Kind
		switch n.Type {
		case "container":
			info.Kind = KindContainer
		case "list":
			info.Kind = KindList
		case "leaf":
			info.Kind = KindLeaf
		default:
			// unknown: treat like container-ish to keep traversal working
			info.Kind = KindContainer
		}

		// List metadata
		if info.Kind == KindList {
			info.ListKey = n.Key
			if n.Key != "" {
				info.ListKeyPath = joinPath(full, n.Key)
			}
		}

		// Leaf metadata
		if info.Kind == KindLeaf {
			compileLeafInfo(&info.Leaf, n)
		}

		// Children bookkeeping
		for _, c := range n.Children {
			childFull := canonicalFullPath(c, full)

			// record child path
			if childFull != "" {
				info.Children = appendUnique(info.Children, childFull)
			}

			walk(c, full)
		}
	}

	for _, r := range roots {
		// Some trimmed schemas include a top-level node with path:"" and children; still safe.
		walk(r, "")
	}

	return idx, nil
}

// ------------------------- Process Trimmed Schema [END] -------------------------

// ------------------------- Helpers [START] -------------------------

func normalizePath(p string) string {
	p = strings.Trim(p, "/")

	// Strip "configuration" root in both forms
	if p == "configuration" {
		return ""
	}
	p = strings.TrimPrefix(p, "configuration/")
	return p
}

func canonicalFullPath(n SchemaNode, parentFull string) string {
	if n.Name == "" {
		return ""
	}

	parentFull = normalizePath(parentFull)

	// Your JSON "path" is usually the parent path (often includes "configuration/")
	p := normalizePath(n.Path)

	switch n.Type {
	case "leaf":
		// leaf full path should be parent-path + leaf name
		if p != "" {
			return joinPath(p, n.Name)
		}
		return joinPath(parentFull, n.Name)

	default: // container, list
		// container/list full path should be its parent + its name
		// Prefer n.Path if present (because your JSON is anchored under configuration)
		if p != "" {
			return joinPath(p, n.Name)
		}
		return joinPath(parentFull, n.Name)
	}
}

func joinPath(a, b string) string {
	a = normalizePath(a)
	b = normalizePath(b)
	if a == "" {
		return b
	}
	if b == "" {
		return a
	}
	return a + "/" + b
}

func appendUnique(xs []string, s string) []string {
	for _, x := range xs {
		if x == s {
			return xs
		}
	}
	return append(xs, s)
}

func compileLeafInfo(out *LeafInfo, n SchemaNode) {
	base := leafBaseFromLeafType(n.LeafType)
	out.Base = base

	// Direct constraints on the leaf node
	out.Patterns = append(out.Patterns, n.Patterns...)
	out.Ranges = append(out.Ranges, n.Ranges...)
	out.Lengths = append(out.Lengths, n.Lengths...)
	out.Enums = enumSetFromEnumValues(n.Enums, out.Enums)

	// Union branches
	if base == LeafUnion {
		for _, br := range n.Types {
			ub := UnionBranch{
				Base:     leafBaseFromLeafType(br.Type),
				Patterns: append([]string(nil), br.Patterns...),
				Ranges:   append([]NumRange(nil), br.Ranges...),
				Lengths:  append([]LenRange(nil), br.Lengths...),
				Enums:    enumSetFromEnumValues(br.Enums, nil),
			}
			out.Union = append(out.Union, ub)

			// Often useful to have a flattened view too:
			out.Patterns = append(out.Patterns, br.Patterns...)
			out.Ranges = append(out.Ranges, br.Ranges...)
			out.Lengths = append(out.Lengths, br.Lengths...)
			out.Enums = mergeEnumSets(out.Enums, ub.Enums)
		}
	}
}

func leafBaseFromLeafType(t string) LeafBase {
	switch strings.ToLower(strings.TrimSpace(t)) {
	case "string":
		return LeafString
	case "boolean", "bool":
		return LeafBool
	case "int8", "int16", "int32", "int64", "int":
		return LeafInt
	case "uint8", "uint16", "uint32", "uint64", "uint":
		return LeafUint
	case "enumeration", "enum":
		return LeafEnum
	case "union":
		return LeafUnion
	default:
		return LeafOther
	}
}

func enumSetFromEnumValues(vals []EnumValue, existing map[string]struct{}) map[string]struct{} {
	if len(vals) == 0 && existing != nil {
		return existing
	}
	if existing == nil {
		existing = make(map[string]struct{})
	}
	for _, v := range vals {
		if v.Name != "" {
			existing[v.Name] = struct{}{}
		}
	}
	return existing
}

func mergeEnumSets(a, b map[string]struct{}) map[string]struct{} {
	if a == nil {
		return b
	}
	for k := range b {
		a[k] = struct{}{}
	}
	return a
}

// ------------------------- Helpers [END] -------------------------