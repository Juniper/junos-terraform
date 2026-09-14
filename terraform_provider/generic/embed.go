package generic

import (
	"compress/gzip"
	"encoding/json"
	"fmt"
	"io"
	"strings"
	"sync"

	"terraform_provider/patch"
)

var (
	schemaOnce  sync.Once
	schemaIndex map[string]*patch.NodeInfo
	schemaNodes []patch.SchemaNode
	schemaErr   error
)

// LoadSchema parses trimmed-schema JSON (raw or gzipped) into the patch engine index.
// Safe to call from multiple goroutines; parsing happens exactly once.
func LoadSchema(raw []byte) (map[string]*patch.NodeInfo, []patch.SchemaNode, error) {
	schemaOnce.Do(func() {
		data := raw
		// Detect gzip magic bytes and decompress if needed.
		if len(raw) >= 2 && raw[0] == 0x1f && raw[1] == 0x8b {
			r, err := gzip.NewReader(strings.NewReader(string(raw)))
			if err != nil {
				schemaErr = fmt.Errorf("decompress schema: %w", err)
				return
			}
			defer func() { _ = r.Close() }()
			data, err = io.ReadAll(r)
			if err != nil {
				schemaErr = fmt.Errorf("read decompressed schema: %w", err)
				return
			}
		}

		var w patch.TrimmedSchemaWrapper
		if err := json.Unmarshal(data, &w); err != nil {
			schemaErr = fmt.Errorf("unmarshal schema JSON: %w", err)
			return
		}
		if len(w.Root.Children) == 0 {
			schemaErr = fmt.Errorf("schema JSON has no root children")
			return
		}

		schemaNodes = w.Root.Children
		schemaIndex, schemaErr = patch.UnmarshalTrimmedSchemaIndex(string(data))
	})
	return schemaIndex, schemaNodes, schemaErr
}

// ResetSchema allows tests to re-run LoadSchema with different data.
func ResetSchema() {
	schemaOnce = sync.Once{}
	schemaIndex = nil
	schemaNodes = nil
	schemaErr = nil
}
