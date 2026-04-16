package main

// place holder module which is overwritten by jinja2

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/resource"
)

// Collects the data for the crud work
type configResource struct {
	client ProviderConfig //nolint:unused
}

// Create is a generated stub until template-backed logic is rendered.
func (r *configResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// stub method
}

// Read is a generated stub until template-backed logic is rendered.
func (r *configResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// stub method
}

// Update implements resource.Resource.
func (r *configResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// stub method
}

// Delete implements resource.Resource.
func (r *configResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// stub method
}

// Metadata implements resource.Resource.
func (r *configResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	// stub method
	resp.TypeName = "terraform_provider"
}

// Schema implements resource.Resource.
func (r *configResource) Schema(_ context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	// stub method
}
