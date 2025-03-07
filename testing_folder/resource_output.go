package main

import (
	"context"
	"encoding/xml"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"strings"
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
			Vlan_tagging *string  `xml:"vlan-tagging,omitempty"`
			Mtu          *int64   `xml:"mtu,omitempty"`
			Unit         struct {
				Name        *string `xml:"name,omitempty"`
				Description *string `xml:"description,omitempty"`
				Vlan_id     *string `xml:"vlan-id,omitempty"`
				Family      struct {
					Inet struct {
						Address struct {
							Name *string `xml:"name,omitempty"`
						}
					}
					Inet6 struct {
						Address struct {
							Name *string `xml:"name,omitempty"`
						}
					}
				}
			}
		} `xml:"interfaces>interface"`
	} `xml:"groups"`
}

// Collecting objects from the .tf file
type InterfacesModel struct {
	ResourceName types.String `tfsdk:"resource_name"`
	Name         types.String `tfsdk:"name"`
	Description  types.String `tfsdk:"description"`
	Vlan_tagging types.Bool   `tfsdk:"vlan_tagging"`
	Mtu          types.Int64  `tfsdk:"mtu"`
	Unit         types.List   `tfsdk:"unit"`
}

func (o InterfacesModel) AttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"name":         types.StringType,
		"description":  types.StringType,
		"vlan_tagging": types.BoolType,
		"mtu":          types.Int64Type,
		"unit":         types.ListType{ElemType: types.ObjectType{AttrTypes: UnitModel{}.AttrTypes()}},
	}
}
func (o InterfacesModel) Attributes() map[string]schema.Attribute {
	return map[string]schema.Attribute{
		"name": schema.StringAttribute{
			Optional:            true,
			MarkdownDescription: "xpth is `config.Groups.Interface.Name",
		},
		"description": schema.StringAttribute{
			Optional:            true,
			MarkdownDescription: "xpth is `config.Groups.Interface.Description",
		},
		"vlan_tagging": schema.BoolAttribute{
			Optional:            true,
			MarkdownDescription: "xpath is `config.Groups.Interface.Vlan_tagging",
		},
		"mtu": schema.Int64Attribute{
			Optional:            true,
			MarkdownDescription: "xpath is `config.Groups.Interface.Mtu",
		},
		"unit": schema.ListNestedAttribute{
			Optional: true,
			NestedObject: schema.NestedAttributeObject{
				Attributes: UnitModel{}.Attributes(),
			},
		},
	}
}

type UnitModel struct {
	Name        types.String `tfsdk:"name"`
	Description types.String `tfsdk:"description"`
	Vlan_id     types.String `tfsdk:"vlan-id"`
	Family      types.List   `tfsdk:"family"`
}

func (o UnitModel) AttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"name":        types.StringType,
		"description": types.StringType,
		"vlan_id":     types.StringType,
		"family":      types.ListType{ElemType: types.ObjectType{AttrTypes: FamilyModel{}.AttrTypes()}},
	}
}
func (o UnitModel) Attributes() map[string]schema.Attribute {
	return map[string]schema.Attribute{
		"name": schema.StringAttribute{
			Optional:            true,
			MarkdownDescription: "xpath is `config.Groups.Unit.Name`",
		},
		"description": schema.StringAttribute{
			Optional:            true,
			MarkdownDescription: "xpath is `config.Groups.Unit.Description`",
		},
		"vlan_id": schema.StringAttribute{
			Optional:            true,
			MarkdownDescription: "xpath is `config.Groups.Unit.Vlan-id`",
		},
		"family": schema.ListNestedAttribute{
			Optional: true,
			NestedObject: schema.NestedAttributeObject{
				Attributes: FamilyModel{}.Attributes(),
			},
		},
	}
}

type FamilyModel struct {
	Inet  types.List `tfsdk:"inet"`
	Inet6 types.List `tfsdk:"inet6"`
}

func (o FamilyModel) AttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"inet":  types.ListType{ElemType: types.ObjectType{AttrTypes: InetModel{}.AttrTypes()}},
		"inet6": types.ListType{ElemType: types.ObjectType{AttrTypes: Inet6Model{}.AttrTypes()}},
	}
}
func (o FamilyModel) Attributes() map[string]schema.Attribute {
	return map[string]schema.Attribute{
		"inet": schema.ListNestedAttribute{
			Optional: true,
			NestedObject: schema.NestedAttributeObject{
				Attributes: InetModel{}.Attributes(),
			},
		},
		"inet6": schema.ListNestedAttribute{
			Optional: true,
			NestedObject: schema.NestedAttributeObject{
				Attributes: Inet6Model{}.Attributes(),
			},
		},
	}
}

type InetModel struct {
	Address types.List `tfsdk:"address"`
}

func (o InetModel) AttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"address": types.ListType{ElemType: types.ObjectType{AttrTypes: AddressModel{}.AttrTypes()}},
	}
}
func (o InetModel) Attributes() map[string]schema.Attribute {
	return map[string]schema.Attribute{
		"address": schema.ListNestedAttribute{
			Optional: true,
			NestedObject: schema.NestedAttributeObject{
				Attributes: AddressModel{}.Attributes(),
			},
		},
	}
}

type Inet6Model struct {
	Address types.List `tfsdk:"address"`
}

func (o Inet6Model) AttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"address": types.ListType{ElemType: types.ObjectType{AttrTypes: AddressModel{}.AttrTypes()}},
	}
}
func (o Inet6Model) Attributes() map[string]schema.Attribute {
	return map[string]schema.Attribute{
		"address": schema.ListNestedAttribute{
			Optional: true,
			NestedObject: schema.NestedAttributeObject{
				Attributes: AddressModel{}.Attributes(),
			},
		},
	}
}

type AddressModel struct {
	Name types.String `tfsdk:"name"`
}

func (o AddressModel) AttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"name": types.StringType,
	}
}
func (o AddressModel) Attributes() map[string]schema.Attribute {
	return map[string]schema.Attribute{
		"name": schema.StringAttribute{
			Optional:            true,
			MarkdownDescription: "xpath is `config.Groups.Address.Name`",
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
			"interface": schema.ListNestedAttribute{
				Optional: true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: InterfacesModel{}.Attributes(),
				},
			},
		},
	}
}

// Create implements resource.Resource.
func (r *resourceInterfaces) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Get the Interfaces Model data and set
	var plan InterfacesModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	// Check for errors
	if resp.Diagnostics.HasError() {
		return
	}

	var config xmlInterfaces
	config.Groups.Name = plan.ResourceName.ValueString()
	config.Groups.Interface.Name = plan.Name.ValueStringPointer()
	config.Groups.Interface.Description = plan.Name.ValueStringPointer()
	if plan.Vlan_tagging.ValueBool() {
		empty := ""
		config.Groups.Interface.Vlan_tagging = &empty
	}
	config.Groups.Interface.Mtu = plan.Mtu.ValueInt64Pointer()

	var var_unit []UnitModel
	resp.Diagnostics.Append(plan.Unit.ElementsAs(ctx, &var_unit, false)...)
	if resp.Diagnostics.HasError() {
		return
	}
	config.Groups.Interface.Unit.Name = var_unit[0].Name.ValueStringPointer()
	config.Groups.Interface.Unit.Description = var_unit[0].Description.ValueStringPointer()
	config.Groups.Interface.Unit.Vlan_id = var_unit[0].Vlan_id.ValueStringPointer()

	var var_family []FamilyModel
	resp.Diagnostics.Append(var_unit[0].Family.ElementsAs(ctx, &var_family, false)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var var_inet []InetModel
	resp.Diagnostics.Append(var_family[0].Inet.ElementsAs(ctx, &var_inet, false)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var var_address []AddressModel
	resp.Diagnostics.Append(var_inet[0].Address.ElementsAs(ctx, &var_address, false)...)
	if resp.Diagnostics.HasError() {
		return
	}
	config.Groups.Interface.Unit.Family.Inet.Address.Name = var_address[0].Name.ValueStringPointer()

	var var_inet6 []Inet6Model
	resp.Diagnostics.Append(var_family[0].Inet6.ElementsAs(ctx, &var_inet6, false)...)
	if resp.Diagnostics.HasError() {
		return
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

}

// Update implements resource.Resource.
func (r *resourceInterfaces) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Get the Interfaces Model data and set
	var plan InterfacesModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	// Check for errors
	if resp.Diagnostics.HasError() {
		return
	}

	var config xmlInterfaces
	config.Groups.Name = plan.ResourceName.ValueString()
	config.Groups.Interface.Name = plan.Name.ValueStringPointer()
	config.Groups.Interface.Description = plan.Name.ValueStringPointer()
	if plan.Vlan_tagging.ValueBool() {
		empty := ""
		config.Groups.Interface.Vlan_tagging = &empty
	}
	config.Groups.Interface.Mtu = plan.Mtu.ValueInt64Pointer()

	err := r.client.SendTransaction("", config, false)
	if err != nil {
		resp.Diagnostics.AddError("Failed while Sending", err.Error())
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
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
