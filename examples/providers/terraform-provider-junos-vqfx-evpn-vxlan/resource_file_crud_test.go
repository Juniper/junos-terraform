package main

import (
	"context"
	"errors"
	"os"
	"path/filepath"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// fileResourceSchema builds the resource schema used to initialize typed plan/state values.
func fileResourceSchema(t *testing.T, rf *resourceFile) resource.SchemaResponse {
	t.Helper()
	resp := resource.SchemaResponse{}
	rf.Schema(context.Background(), resource.SchemaRequest{}, &resp)
	return resp
}

// TestResourceFileCreateReadUpdateDeleteEndToEnd verifies full CRUD behavior.
func TestResourceFileCreateReadUpdateDeleteEndToEnd(t *testing.T) {
	ctx := context.Background()
	dir := t.TempDir()
	rf := &resourceFile{dir: dir}
	schemaResp := fileResourceSchema(t, rf)

	createPlan := fileModel{
		Name:     types.StringValue("managed.txt"),
		Contents: types.StringValue("initial"),
	}
	createReqPlan := tfsdk.Plan{Schema: schemaResp.Schema}
	if diags := createReqPlan.Set(ctx, createPlan); diags.HasError() {
		t.Fatalf("failed to build create plan: %v", diags)
	}
	createReq := resource.CreateRequest{Plan: createReqPlan}
	createResp := &resource.CreateResponse{State: tfsdk.State{Schema: schemaResp.Schema}}
	rf.Create(ctx, createReq, createResp)
	if createResp.Diagnostics.HasError() {
		t.Fatalf("Create() diagnostics: %v", createResp.Diagnostics)
	}

	readState := tfsdk.State{Schema: schemaResp.Schema}
	if diags := readState.Set(ctx, fileModel{Name: types.StringValue("managed.txt")}); diags.HasError() {
		t.Fatalf("failed to build read state: %v", diags)
	}
	readReq := resource.ReadRequest{State: readState}
	readResp := &resource.ReadResponse{State: tfsdk.State{Schema: schemaResp.Schema}}
	rf.Read(ctx, readReq, readResp)
	if readResp.Diagnostics.HasError() {
		t.Fatalf("Read() diagnostics: %v", readResp.Diagnostics)
	}

	updatePlan := tfsdk.Plan{Schema: schemaResp.Schema}
	if diags := updatePlan.Set(ctx, fileModel{Name: types.StringValue("managed.txt"), Contents: types.StringValue("updated")}); diags.HasError() {
		t.Fatalf("failed to build update plan: %v", diags)
	}
	updateReq := resource.UpdateRequest{Plan: updatePlan}
	updateResp := &resource.UpdateResponse{State: tfsdk.State{Schema: schemaResp.Schema}}
	rf.Update(ctx, updateReq, updateResp)
	if updateResp.Diagnostics.HasError() {
		t.Fatalf("Update() diagnostics: %v", updateResp.Diagnostics)
	}

	deleteState := tfsdk.State{Schema: schemaResp.Schema}
	if diags := deleteState.Set(ctx, fileModel{Name: types.StringValue("managed.txt")}); diags.HasError() {
		t.Fatalf("failed to build delete state: %v", diags)
	}
	deleteReq := resource.DeleteRequest{State: deleteState}
	deleteResp := &resource.DeleteResponse{}
	rf.Delete(ctx, deleteReq, deleteResp)
	if deleteResp.Diagnostics.HasError() {
		t.Fatalf("Delete() diagnostics: %v", deleteResp.Diagnostics)
	}

	if _, err := os.Stat(filepath.Join(dir, "managed.txt")); !errors.Is(err, os.ErrNotExist) {
		t.Fatalf("expected managed file to be deleted, stat err: %v", err)
	}
}

// TestResourceFileCrudDiagPaths verifies unset request payloads panic in framework decoding.
func TestResourceFileCrudDiagPaths(t *testing.T) {
	ctx := context.Background()
	rf := &resourceFile{dir: t.TempDir()}

	tests := []struct {
		name string
		fn   func()
	}{
		{
			name: "create",
			fn: func() {
				rf.Create(ctx, resource.CreateRequest{}, &resource.CreateResponse{})
			},
		},
		{
			name: "read",
			fn: func() {
				rf.Read(ctx, resource.ReadRequest{}, &resource.ReadResponse{})
			},
		},
		{
			name: "update",
			fn: func() {
				rf.Update(ctx, resource.UpdateRequest{}, &resource.UpdateResponse{})
			},
		},
		{
			name: "delete",
			fn: func() {
				rf.Delete(ctx, resource.DeleteRequest{}, &resource.DeleteResponse{})
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			defer func() {
				if recover() == nil {
					t.Fatalf("expected panic when request payload is unset")
				}
			}()
			tc.fn()
		})
	}
}

// TestResourceFileCrudFilesystemErrors verifies diagnostics for filesystem failures.
func TestResourceFileCrudFilesystemErrors(t *testing.T) {
	ctx := context.Background()
	rf := &resourceFile{dir: filepath.Join(t.TempDir(), "does-not-exist")}
	schemaResp := fileResourceSchema(t, rf)

	plan := tfsdk.Plan{Schema: schemaResp.Schema}
	if diags := plan.Set(ctx, fileModel{Name: types.StringValue("f.txt"), Contents: types.StringValue("c")}); diags.HasError() {
		t.Fatalf("failed to build plan: %v", diags)
	}

	createResp := &resource.CreateResponse{State: tfsdk.State{Schema: schemaResp.Schema}}
	rf.Create(ctx, resource.CreateRequest{Plan: plan}, createResp)
	if !createResp.Diagnostics.HasError() {
		t.Fatalf("expected Create() filesystem diagnostic")
	}

	readState := tfsdk.State{Schema: schemaResp.Schema}
	if diags := readState.Set(ctx, fileModel{Name: types.StringValue("f.txt")}); diags.HasError() {
		t.Fatalf("failed to build state: %v", diags)
	}
	readResp := &resource.ReadResponse{State: tfsdk.State{Schema: schemaResp.Schema}}
	rf.Read(ctx, resource.ReadRequest{State: readState}, readResp)
	if !readResp.Diagnostics.HasError() {
		t.Fatalf("expected Read() filesystem diagnostic")
	}

	updateResp := &resource.UpdateResponse{State: tfsdk.State{Schema: schemaResp.Schema}}
	rf.Update(ctx, resource.UpdateRequest{Plan: plan}, updateResp)
	if !updateResp.Diagnostics.HasError() {
		t.Fatalf("expected Update() filesystem diagnostic")
	}

	deleteResp := &resource.DeleteResponse{}
	rf.Delete(ctx, resource.DeleteRequest{State: readState}, deleteResp)
	if deleteResp.Diagnostics.HasError() {
		t.Fatalf("expected Delete() missing file path to be treated as no-op")
	}
}
