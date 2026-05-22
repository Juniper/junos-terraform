package main

import (
	"errors"
	"os"
	"path"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/types"
)

// TestWriteManagedFileAndReadManagedFile verifies round-trip write/read helper behavior.
func TestWriteManagedFileAndReadManagedFile(t *testing.T) {
	dir := t.TempDir()
	model := fileModel{
		Name:     types.StringValue("sample.txt"),
		Contents: types.StringValue("hello"),
	}

	if err := writeManagedFile(dir, model); err != nil {
		t.Fatalf("writeManagedFile() error: %v", err)
	}

	state := fileModel{Name: types.StringValue("sample.txt")}
	if err := readManagedFile(dir, &state); err != nil {
		t.Fatalf("readManagedFile() error: %v", err)
	}
	if state.Contents.ValueString() != "hello" {
		t.Fatalf("unexpected contents: %q", state.Contents.ValueString())
	}
}

// TestReadManagedFileMissing verifies read helper errors for missing files.
func TestReadManagedFileMissing(t *testing.T) {
	dir := t.TempDir()
	state := fileModel{Name: types.StringValue("missing.txt")}

	err := readManagedFile(dir, &state)
	if err == nil {
		t.Fatalf("expected readManagedFile() error for missing file")
	}
}

// TestDeleteManagedFile verifies delete helper removes an existing file.
func TestDeleteManagedFile(t *testing.T) {
	dir := t.TempDir()
	filePath := path.Join(dir, "delete-me.txt")
	if err := os.WriteFile(filePath, []byte("x"), 0644); err != nil {
		t.Fatalf("setup write failed: %v", err)
	}

	state := fileModel{Name: types.StringValue("delete-me.txt")}
	if err := deleteManagedFile(dir, state); err != nil {
		t.Fatalf("deleteManagedFile() error: %v", err)
	}
	if _, err := os.Stat(filePath); !errors.Is(err, os.ErrNotExist) {
		t.Fatalf("expected file to be deleted, got stat err: %v", err)
	}
}

// TestDeleteManagedFileMissingIsNoop verifies delete helper ignores missing files.
func TestDeleteManagedFileMissingIsNoop(t *testing.T) {
	dir := t.TempDir()
	state := fileModel{Name: types.StringValue("missing.txt")}

	if err := deleteManagedFile(dir, state); err != nil {
		t.Fatalf("expected no error for missing file delete, got: %v", err)
	}
}
