package generic

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
	"os"
	"runtime"
	"testing"
	"time"

	"terraform_provider/patch"

	"github.com/hashicorp/terraform-plugin-go/tftypes"
)

// These measurements answer whether a provider can serve a full, unfiltered
// Junos model from memory, and what the compiled schema buys over the pyang
// JSON. They need a real schema and are skipped otherwise:
//
//	JTAF_FULL_SCHEMA=/path/junos.json go test ./generic -run TestFullSchema -v -timeout 30m

func fullSchemaPath(t *testing.T) string {
	t.Helper()
	path := os.Getenv("JTAF_FULL_SCHEMA")
	if path == "" {
		t.Skip("set JTAF_FULL_SCHEMA to a pyang schema JSON path to run this measurement")
	}
	return path
}

func readSchema(t *testing.T, path string) []byte {
	t.Helper()
	raw, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read schema: %v", err)
	}
	return raw
}

func heapMB() float64 {
	// Two cycles: after churning a few hundred MB, one collection can leave
	// garbage behind and inflate the reading.
	runtime.GC()
	runtime.GC()
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	return float64(m.HeapAlloc) / (1 << 20)
}

func mb(n int) float64 { return float64(n) / (1 << 20) }

func gzipLen(t *testing.T, data []byte) int {
	t.Helper()
	var buf bytes.Buffer
	w := gzip.NewWriter(&buf)
	if _, err := w.Write(data); err != nil {
		t.Fatalf("gzip: %v", err)
	}
	if err := w.Close(); err != nil {
		t.Fatalf("gzip close: %v", err)
	}
	return buf.Len()
}

// TestFullSchemaPhases breaks provider startup down by phase, so it is clear
// which one dominates and how much of the pyang tree survives compilation.
func TestFullSchemaPhases(t *testing.T) {
	path := fullSchemaPath(t)

	// Measured before the file is read, so every later figure is the heap the
	// provider is actually holding.
	baseline := heapMB()
	raw := readSchema(t, path)
	t.Logf("schema file:       %.1f MB", mb(len(raw)))

	start := time.Now()
	var w patch.TrimmedSchemaWrapper
	if err := json.Unmarshal(raw, &w); err != nil {
		t.Fatalf("unmarshal schema: %v", err)
	}
	unmarshalDur := time.Since(start)
	raw = nil
	pyangHeap := heapMB()

	start = time.Now()
	schema := patch.CompileSchema(w.Root.Children)
	compileDur := time.Since(start)

	// Release the pyang tree: only the compiled schema stays reachable, which
	// is what the provider holds for the rest of its life.
	w = patch.TrimmedSchemaWrapper{}
	compiledHeap := heapMB()

	start = time.Now()
	_, objType := BuildSchema(schema)
	buildDur := time.Since(start)

	start = time.Now()
	typeJSON, err := objType.MarshalJSON()
	if err != nil {
		t.Fatalf("marshal value type: %v", err)
	}
	marshalDur := time.Since(start)

	t.Logf("1. json.Unmarshal:  %v", unmarshalDur)
	t.Logf("2. CompileSchema:   %v (%d nodes)", compileDur, schema.Len())
	t.Logf("3. BuildSchema:     %v", buildDur)
	t.Logf("4. marshal type:    %v", marshalDur)
	t.Logf("pyang tree heap:    %.1f MB", pyangHeap-baseline)
	t.Logf("STEADY heap:        %.1f MB", compiledHeap-baseline)
	t.Logf("released by compile: %.1f MB", pyangHeap-compiledHeap)
	t.Logf("schema wire size:   %.1f MB", mb(len(typeJSON)))

	runtime.KeepAlive(schema)
}

// TestFullSchemaCompiledVsJSON quantifies what cmd/compileschema buys: the
// embedded form's size, and how much of startup it removes.
func TestFullSchemaCompiledVsJSON(t *testing.T) {
	path := fullSchemaPath(t)
	raw := readSchema(t, path)
	jsonSize, jsonGz := len(raw), gzipLen(t, raw)

	start := time.Now()
	var w patch.TrimmedSchemaWrapper
	if err := json.Unmarshal(raw, &w); err != nil {
		t.Fatalf("unmarshal schema: %v", err)
	}
	fromJSON := patch.CompileSchema(w.Root.Children)
	jsonDur := time.Since(start)
	wantNodes := fromJSON.Len()

	compiled, err := fromJSON.MarshalBinary()
	if err != nil {
		t.Fatalf("marshal compiled schema: %v", err)
	}
	if !patch.IsCompiledSchema(compiled) {
		t.Fatal("MarshalBinary output is not recognised as a compiled schema")
	}
	compiledGz := gzipLen(t, compiled)

	// Drop everything but the compiled bytes so the load below is measured on
	// its own, the way a provider start sees it.
	raw, w, fromJSON = nil, patch.TrimmedSchemaWrapper{}, nil
	runtime.GC()

	start = time.Now()
	fromCompiled, err := patch.UnmarshalSchema(compiled)
	if err != nil {
		t.Fatalf("unmarshal compiled schema: %v", err)
	}
	compiledDur := time.Since(start)

	if fromCompiled.Len() != wantNodes {
		t.Errorf("compiled schema has %d nodes, JSON path has %d",
			fromCompiled.Len(), wantNodes)
	}

	// Steady heap is reported by TestFullSchemaPhases; both paths end at the
	// same *patch.Schema, so it is not repeated here.
	t.Logf("JSON:      %.1f MB (%.1f MB gzipped), %v to load",
		mb(jsonSize), mb(jsonGz), jsonDur)
	t.Logf("compiled:  %.1f MB (%.1f MB gzipped), %v to load",
		mb(len(compiled)), mb(compiledGz), compiledDur)
	t.Logf("gain:      %.0fx faster to load, %.0fx smaller raw, %.1fx smaller gzipped",
		float64(jsonDur)/float64(compiledDur),
		float64(jsonSize)/float64(len(compiled)),
		float64(jsonGz)/float64(compiledGz))

	runtime.KeepAlive(fromCompiled)
}

// TestFullSchemaServesSchema exercises the path Terraform calls first, so a
// model that cannot be served fails here rather than at terraform plan.
func TestFullSchemaServesSchema(t *testing.T) {
	raw := readSchema(t, fullSchemaPath(t))

	ResetSchema()
	t.Cleanup(ResetSchema)

	start := time.Now()
	schema, err := LoadSchema(raw)
	if err != nil {
		t.Fatalf("load schema: %v", err)
	}
	loadDur := time.Since(start)

	protoSchema, objType := BuildSchema(schema)
	if protoSchema.Block == nil {
		t.Fatal("provider served a schema with no block")
	}
	if len(objType.AttributeTypes) == 0 {
		t.Fatal("provider served a schema with no attributes")
	}
	if _, err := objType.MarshalJSON(); err != nil {
		t.Fatalf("value type is not serialisable: %v", err)
	}

	// An empty resource must convert through the served type, which is the
	// conversion every plan and apply starts from.
	root, err := ValueToConfig(tftypes.NewValue(objType, nil), schema)
	if err != nil {
		t.Fatalf("null value does not convert: %v", err)
	}
	if root == nil {
		t.Fatal("expected a configuration root")
	}

	t.Logf("LoadSchema:       %v", loadDur)
	t.Logf("top-level attrs:  %d", len(objType.AttributeTypes))
	t.Logf("nodes:            %d", schema.Len())
}
