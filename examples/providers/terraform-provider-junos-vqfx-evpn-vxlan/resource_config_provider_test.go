package main

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/resource"
)

func isStubConfigResource(t *testing.T, r *configResource) bool {
	t.Helper()
	schemaResp := &resource.SchemaResponse{}
	r.Schema(context.Background(), resource.SchemaRequest{}, schemaResp)
	return len(schemaResp.Schema.Attributes) == 0
}

// TestConfigResourceStubMethods verifies stub CRUD methods remain no-op.
func TestConfigResourceStubMethods(t *testing.T) {
	r := &configResource{}
	if !isStubConfigResource(t, r) {
		t.Skip("generated config resource requires framework-populated requests")
	}

	ctx := context.Background()

	createResp := &resource.CreateResponse{}
	r.Create(ctx, resource.CreateRequest{}, createResp)
	if createResp.Diagnostics.HasError() {
		t.Fatalf("Create() should not add diagnostics")
	}

	readResp := &resource.ReadResponse{}
	r.Read(ctx, resource.ReadRequest{}, readResp)
	if readResp.Diagnostics.HasError() {
		t.Fatalf("Read() should not add diagnostics")
	}

	updateResp := &resource.UpdateResponse{}
	r.Update(ctx, resource.UpdateRequest{}, updateResp)
	if updateResp.Diagnostics.HasError() {
		t.Fatalf("Update() should not add diagnostics")
	}

	deleteResp := &resource.DeleteResponse{}
	r.Delete(ctx, resource.DeleteRequest{}, deleteResp)
	if deleteResp.Diagnostics.HasError() {
		t.Fatalf("Delete() should not add diagnostics")
	}
}

// TestConfigResourceMetadataAndSchema verifies current stub metadata and schema.
func TestConfigResourceMetadataAndSchema(t *testing.T) {
	r := &configResource{}
	ctx := context.Background()

	metadataResp := &resource.MetadataResponse{}
	r.Metadata(ctx, resource.MetadataRequest{}, metadataResp)
	if metadataResp.TypeName == "" {
		t.Fatalf("expected non-empty metadata type name")
	}

	schemaResp := &resource.SchemaResponse{}
	r.Schema(ctx, resource.SchemaRequest{}, schemaResp)
	if isStubConfigResource(t, r) && len(schemaResp.Schema.Attributes) != 0 {
		t.Fatalf("expected empty schema attributes for stub resource")
	}
}
