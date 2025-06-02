package main

import (
	"context"
	"os"
	"path"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var _ resource.ResourceWithConfigure = new(resourceFile)

type resourceFile struct {
	dir string
}

// Configure implements resource.ResourceWithConfigure.
func (r *resourceFile) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	r.dir = req.ProviderData.(string)
}

type fileModel struct {
	Name     types.String `tfsdk:"name"`
	Contents types.String `tfsdk:"contents"`
}

// Metadata implements resource.Resource.
func (r *resourceFile) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_file"
}

// Schema implements resource.Resource.
func (r *resourceFile) Schema(_ context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"name": schema.StringAttribute{
				Required:      true,
				PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplace()},
			},
			"contents": schema.StringAttribute{
				Required: true,
			},
		},
	}
}

// Create implements resource.Resource.
func (r *resourceFile) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan fileModel
	d := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(d...)
	if resp.Diagnostics.HasError() {
		return
	}
	err := os.WriteFile(path.Join(r.dir, plan.Name.ValueString()), []byte(plan.Contents.ValueString()), 0644)
	if err != nil {
		resp.Diagnostics.AddError("failed writing file", err.Error())
		return
	}
	d = resp.State.Set(ctx, &plan)
	resp.Diagnostics.Append(d...)
}

// Read implements resource.Resource.
func (r *resourceFile) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state fileModel
	d := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(d...)
	if resp.Diagnostics.HasError() {
		return
	}
	data, err := os.ReadFile(path.Join(r.dir, state.Name.ValueString()))
	if err != nil {
		resp.Diagnostics.AddError("failed reading file", err.Error())
		return
	}
	state.Contents = types.StringValue(string(data))
	d = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(d...)
}

// Update implements resource.Resource.
func (r *resourceFile) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan fileModel
	d := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(d...)
	if resp.Diagnostics.HasError() {
		return
	}
	err := os.WriteFile(path.Join(r.dir, plan.Name.ValueString()), []byte(plan.Contents.ValueString()), 0644)
	if err != nil {
		resp.Diagnostics.AddError("failed writing file", err.Error())
		return
	}
	d = resp.State.Set(ctx, &plan)
	resp.Diagnostics.Append(d...)
}

// Delete implements resource.Resource.
func (r *resourceFile) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state fileModel
	d := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(d...)
	if resp.Diagnostics.HasError() {
		return
	}
	err := os.Remove(path.Join(r.dir, state.Name.ValueString()))
	if err != nil {
		if strings.Contains(err.Error(), "ound") {
			return
		}
		resp.Diagnostics.AddError("failed deleting file", err.Error())
		return
	}
}

