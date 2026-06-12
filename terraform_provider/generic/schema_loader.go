package generic

import (
	"encoding/json"

	"terraform_provider/patch"
)

// SchemaIndex holds the parsed schema tree and flat index for runtime use.
type SchemaIndex struct {
	TopLevel []patch.SchemaNode
	ByPath   map[string]*patch.NodeInfo
}

// LoadSchema parses trimmed schema JSON into a SchemaIndex containing both
// the tree structure (TopLevel) and the flat path→NodeInfo index (ByPath).
func LoadSchema(jsonBytes []byte) (*SchemaIndex, error) {
	var w patch.TrimmedSchemaWrapper
	if err := json.Unmarshal(jsonBytes, &w); err != nil {
		return nil, err
	}

	// Build the flat index using the existing patch function.
	idx, err := patch.UnmarshalTrimmedSchemaIndex(string(jsonBytes))
	if err != nil {
		return nil, err
	}

	// Extract top-level children (direct children of configuration root).
	var topLevel []patch.SchemaNode
	for _, root := range w.Root.Children {
		if root.Name == "configuration" || root.Name == "root" {
			topLevel = root.Children
			break
		}
		topLevel = append(topLevel, root)
	}

	return &SchemaIndex{
		TopLevel: topLevel,
		ByPath:   idx,
	}, nil
}
