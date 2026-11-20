package patch

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
)

// ------------------------- Type Definitions [START] -------------------------

type NodeKind uint8

const (
	KindContainer NodeKind = iota
	KindList
	KindLeaf
	KindLeafList
)

// Raw JSON node
type SchemaNode struct {
	Name      string `json:"name"`
	Type      string `json:"type"`       // container | list | leaf
	Path      string `json:"path"`       // often parent path
	Key       string `json:"key"`        // list key
	LeafType  string `json:"leaf-type"`  // leaf-only: string, union, etc.
	OrderedBy string `json:"ordered-by"` // "user" for ordered leaf-lists/lists
	Presence  string `json:"presence"`   // container: YANG presence statement, if any

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
	Type     string      `json:"type"`
	Path     string      `json:"path"`
	Patterns []string    `json:"patterns"`
	Ranges   []NumRange  `json:"ranges"`
	Lengths  []LenRange  `json:"lengths"`
	Enums    []EnumValue `json:"enums"`
}

type NumRange struct {
	Min  *float64 `json:"min"`
	Max  *float64 `json:"max"`
	Path string   `json:"path"`
}

// UnmarshalJSON accepts a bound as a JSON number or a numeric string: the
// pyang plugin writes decimal64 bounds as strings to keep them exact
// ("9223372036.854775807").
func (r *NumRange) UnmarshalJSON(data []byte) error {
	var raw struct {
		Min  json.RawMessage `json:"min"`
		Max  json.RawMessage `json:"max"`
		Path string          `json:"path"`
	}
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}
	var err error
	if r.Min, err = parseBound(raw.Min); err != nil {
		return fmt.Errorf("range min: %w", err)
	}
	if r.Max, err = parseBound(raw.Max); err != nil {
		return fmt.Errorf("range max: %w", err)
	}
	r.Path = raw.Path
	return nil
}

func parseBound(raw json.RawMessage) (*float64, error) {
	if len(raw) == 0 || string(raw) == "null" {
		return nil, nil
	}
	text := string(raw)
	var s string
	if json.Unmarshal(raw, &s) == nil {
		text = s
	}
	v, err := strconv.ParseFloat(text, 64)
	if err != nil {
		return nil, err
	}
	return &v, nil
}

type LenRange struct {
	Min  *int   `json:"min"`
	Max  *int   `json:"max"`
	Path string `json:"path"`
}

type EnumValue struct {
	Name  string `json:"name"`
	Value any    `json:"value"`
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

// UnmarshalTrimmedSchemaIndex compiles a schema, as JSON from the pyang
// plugin, into a Schema.
func UnmarshalTrimmedSchemaIndex(trimmedSchemaJSON string) (*Schema, error) {
	var w TrimmedSchemaWrapper
	if err := json.Unmarshal([]byte(trimmedSchemaJSON), &w); err != nil {
		return nil, err
	}
	return CompileSchema(w.Root.Children), nil
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

// ------------------------- Helpers [END] -------------------------
