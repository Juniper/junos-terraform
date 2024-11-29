package main

import (
	"context"
	"encoding/xml"
	"strconv"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Junos XML Hierarchy

type xmlInterfaces struct {
	XMLName xml.Name `xml:"configuration"`
	Groups  struct {
		XMLName   xml.Name `xml:"groups"`
		Name      string   `xml:"name"`
		Interface struct {
			XMLName      xml.Name `xml:"interface"`
			Name         *string  `xml:"name,omitempty"`
			Description  *string  `xml:"description,omitempty"`
			Mtu          *int64   `xml:"mtu,omitempty"`
			Vlan_tagging *string  `xml:"vlan-tagging,omitempty"`
		} `xml:"interfaces>interface"`
	} `xml:"groups"`
}

// Collecting objects from the .tf file
type InterfacesModel struct {
	ResourceName types.String `tfsdk:"resource_name"`
	Name         types.String `tfsdk:"name"`
	Description  types.String `tfsdk:"description"`
	Mtu          types.Int64  `tfsdk:"mtu"`
	Vlan_tagging types.Bool   `tfsdk:"vlan_tagging"`
}

func (o InterfacesModel) AttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"name":         types.StringType,
		"description":  types.StringType,
		"mtu":          types.Int64Type,
		"vlan_tagging": types.BoolType,
	}
}
func (o InterfacesModel) Attributes() map[string]schema.Attribute {
	return map[string]schema.Attribute{
		"name": schema.StringAttribute{
			Required:            true,
			MarkdownDescription: "xpath is `config.Groups.Interface.Name",
		},
		"description": schema.StringAttribute{
			Optional:            true,
			MarkdownDescription: "xpath is `config.Groups.Interface.Description",
		},
		"mtu": schema.Int64Attribute{
			Optional:            true,
			MarkdownDescription: "xpath is `config.Groups.Interface.Mtu",
		},
		"vlan_tagging": schema.BoolAttribute{
			Optional:            true,
			MarkdownDescription: "xpath is `config.Groups.Interface.Vlan_tagging",
		},
	}
}

// Collects the data for the crud work
type resourceInterfaces struct {
	client ProviderConfig
}

func (r *resourceInterfaces) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	r.client = req.ProviderData.(ProviderConfig)
}

// Metadata implements resource.Resource.
func (r *resourceInterfaces) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_Interfaces"
}

// Schema implements resource.Resource.
func (r *resourceInterfaces) Schema(_ context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"resource_name": schema.StringAttribute{
				Required:      true,
				PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplace()},
			},
			"name": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "xpath is: `config.Groups.Interface`",
			},
			"description": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "xpath is: `config.Groups.Description`",
			},
			"mtu": schema.Int64Attribute{
				Optional:            true,
				MarkdownDescription: "xpath is: `config.Groups.Mtu`",
			},
			"vlan_tagging": schema.BoolAttribute{
				Optional:            true,
				MarkdownDescription: "xpath is: `config.Groups.Vlan_tagging`",
			},
		},
	}
}

// Create implements resource.Resource.
func (r *resourceInterfaces) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Get the  Model data and set
	var plan InterfacesModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	// Check for errors
	if resp.Diagnostics.HasError() {
		return
	}

	var config xmlInterfaces
	config.Groups.Name = plan.ResourceName.ValueString()
	config.Groups.Interface.Name = plan.Name.ValueStringPointer()
	config.Groups.Interface.Description = plan.Description.ValueStringPointer()
	config.Groups.Interface.Mtu = plan.Mtu.ValueInt64Pointer()
	if plan.Vlan_tagging.ValueBool() {
		empty := ""
		config.Groups.Interface.Vlan_tagging = &empty
	}
	err := r.client.SendTransaction("", config, false)
	if err != nil {
		resp.Diagnostics.AddError("Failed while Sending", err.Error())
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

// Read implements resource.Resource.
func (r *resourceInterfaces) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {

	// Get the data and set
	var state InterfacesModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	// Check for errors
	if resp.Diagnostics.HasError() {
		return
	}
	// Marshall group and check
	var config xmlInterfaces
	err := r.client.MarshalGroup(state.ResourceName.ValueString(), &config)
	if err != nil {
		resp.Diagnostics.AddError("Failed while Reading", err.Error())
		return
	}

	state.Name = types.StringPointerValue(config.Groups.Interface.Name)
	state.Description = types.StringPointerValue(config.Groups.Interface.Description)
	state.Mtu = types.Int64PointerValue(config.Groups.Interface.Mtu)
	var vlanTagging *bool
	if config.Groups.Interface.Vlan_tagging != nil {
		b, err := strconv.ParseBool(*config.Groups.Interface.Vlan_tagging)
		if err == nil {
			vlanTagging = &b
		}
	}
	state.Vlan_tagging = types.BoolPointerValue(vlanTagging)
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

// Update implements resource.Resource.
func (r *resourceInterfaces) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {

}

// Delete implements resource.Resource.
func (r *resourceInterfaces) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state InterfacesModel
	d := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(d...)
	if resp.Diagnostics.HasError() {
		return
	}

	_, err := r.client.DeleteConfig(state.ResourceName.ValueString(), false)
	if err != nil {
		if strings.Contains(err.Error(), "ound") {
			return
		}
		resp.Diagnostics.AddError("Failed while deleting dile", err.Error())
		return
	}
}
