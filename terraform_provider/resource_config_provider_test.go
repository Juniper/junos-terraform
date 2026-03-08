package main

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/resource"
)

// TestResourceApplyGroupsStubMethods verifies stub CRUD methods remain no-op.
func TestResourceApplyGroupsStubMethods(t *testing.T) {
	r := &resource_Apply_Groups{}
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

// TestResourceApplyGroupsMetadataAndSchema verifies current stub metadata and schema.
func TestResourceApplyGroupsMetadataAndSchema(t *testing.T) {
	r := &resource_Apply_Groups{}
	ctx := context.Background()

	metadataResp := &resource.MetadataResponse{}
	r.Metadata(ctx, resource.MetadataRequest{}, metadataResp)
	if metadataResp.TypeName != "terraform_provider" {
		t.Fatalf("unexpected metadata type name: %q", metadataResp.TypeName)
	}

	schemaResp := &resource.SchemaResponse{}
	r.Schema(ctx, resource.SchemaRequest{}, schemaResp)
	if len(schemaResp.Schema.Attributes) != 0 {
		t.Fatalf("expected empty schema attributes for stub resource")
	}
}
