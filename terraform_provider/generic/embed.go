package generic

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
	"fmt"
	"io"
	"sync"

	"terraform_provider/patch"
)

var (
	schemaOnce sync.Once
	schema     *patch.Schema
	schemaErr  error
)

// LoadSchema returns the embedded schema: a compiled schema (patch.Schema's
// binary form), or the pyang JSON (trimmed or a full model), either of them
// raw or gzipped. Safe to call from multiple goroutines; it is read once.
func LoadSchema(raw []byte) (*patch.Schema, error) {
	schemaOnce.Do(func() {
		data := raw
		// Detect gzip magic bytes and decompress if needed.
		if len(raw) >= 2 && raw[0] == 0x1f && raw[1] == 0x8b {
			r, err := gzip.NewReader(bytes.NewReader(raw))
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

		if patch.IsCompiledSchema(data) {
			schema, schemaErr = patch.UnmarshalSchema(data)
			return
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
		schema = patch.CompileSchema(w.Root.Children)
	})
	return schema, schemaErr
}

// ResetSchema allows tests to re-run LoadSchema with different data.
func ResetSchema() {
	schemaOnce = sync.Once{}
	schema = nil
	schemaErr = nil
}
