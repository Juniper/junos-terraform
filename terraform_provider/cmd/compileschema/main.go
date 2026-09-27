// Command compileschema compiles a schema, as JSON from the pyang plugin
// (raw or gzipped), into a patch.Schema's binary form (gzipped), which the
// generic provider embeds so that it does not parse the JSON each time it
// starts.
//
//	go run ./cmd/compileschema trimmed_schema.json.gz schema.bin.gz
package main

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
	"fmt"
	"io"
	"os"

	"terraform_provider/patch"
)

func main() {
	if len(os.Args) != 3 {
		fmt.Fprintln(os.Stderr, "usage: compileschema SCHEMA.json[.gz] OUT.bin.gz")
		os.Exit(2)
	}
	if err := run(os.Args[1], os.Args[2]); err != nil {
		fmt.Fprintln(os.Stderr, "compileschema:", err)
		os.Exit(1)
	}
}

func run(in, out string) error {
	data, err := os.ReadFile(in)
	if err != nil {
		return err
	}
	if len(data) >= 2 && data[0] == 0x1f && data[1] == 0x8b {
		r, err := gzip.NewReader(bytes.NewReader(data))
		if err != nil {
			return err
		}
		if data, err = io.ReadAll(r); err != nil {
			return err
		}
	}
	var w patch.TrimmedSchemaWrapper
	if err := json.Unmarshal(data, &w); err != nil {
		return fmt.Errorf("%s: %w", in, err)
	}
	if len(w.Root.Children) == 0 {
		return fmt.Errorf("%s: no root children", in)
	}
	compiled, err := patch.CompileSchema(w.Root.Children).MarshalBinary()
	if err != nil {
		return err
	}
	var buf bytes.Buffer
	zw, err := gzip.NewWriterLevel(&buf, gzip.BestCompression)
	if err != nil {
		return err
	}
	if _, err := zw.Write(compiled); err != nil {
		return err
	}
	if err := zw.Close(); err != nil {
		return err
	}
	return os.WriteFile(out, buf.Bytes(), 0o644)
}
