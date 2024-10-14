package main

import (
	"context"
	"encoding/xml"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// v_ is appended before every variable so it doesn't give any conflict
// with any keyword in golang. ex - interface is keyword in golang
type xmlInterfaces struct {
	XMLName xml.Name `xml:"configuration"`
	Groups  struct {
		XMLName     xml.Name `xml:"groups"`
		Name        string   `xml:"name"`
		V_interface struct {
			XMLName        xml.Name `xml:"interface"`
			V_name         *string  `xml:"name,omitempty"`
			V_description  *string  `xml:"description,omitempty"`
			V_mtu          *int64   `xml:"mtu,omitempty"`
			V_vlan_tagging *string  `xml:"vlan-tagging,omitempty"`
			V_units        struct {
				XMLName       xml.Name `xml:"unit"`
				V_name        *string  `xml:"name,omitempty"`
				V_description *string  `xml:"description,omitempty"`
				V_vlan_id     *int32   `xml:"vlan-id,omitempty"`
				V_family      struct {
					XMLName xml.Name `xml:"family,omitempty"`
					V_inet  struct {
						XMLName   xml.Name `xml:"inet,omitempty"`
						V_address struct {
							XMLName xml.Name `xml:"address,omitempty"`
							V_name  *string  `xml:"name,omitempty"`
						} `xml:"address"`
					} `xml:"inet"`
					V_inet6 struct {
						XMLName   xml.Name `xml:"inet6,omitempty"`
						V_address struct {
							XMLName xml.Name `xml:"address,omitempty"`
							V_name  *string  `xml:"name,omitempty"`
						} `xml:"address"`
					} `xml:"inet6"`
				} `xml:"family"`
			} `xml:"unit"`
		} `xml:"interfaces>interface"`
	} `xml:"groups"`
}

// Collects the objects from the .tf file
type InterfacesModel struct {
	ResourceName types.String `tfsdk:"resource_name"`
	Name         types.String `tfsdk:"name"`
	Description  types.String `tfsdk:"description"`
	Mtu          types.Int64  `tfsdk:"mtu"`
	Vlan_tagging types.Bool   `tfsdk:"vlan_tagging"`
	Units        []UnitsModel `tfsdk:"units"`
}

type UnitsModel struct {
	Name        types.String  `tfsdk:"name"`
	Description types.String  `tfsdk:"description"`
	Vlan_id     types.Int32   `tfsdk:"vlan_id"`
	Family      []FamilyModel `tfsdk:"family"`
}

type FamilyModel struct {
	Inet  types.List   `tfsdk:"inet"`
	Inet6 []Inet6Model `tfsdk:"inet6"`
}

type InetModel struct {
	Address types.String `tfsdk:"address"`
}

type Inet6Model struct {
	Address types.String `tfsdk:"address"`
}

func (o UnitsModel) Attributes() map[string]schema.Attribute {
	return map[string]schema.Attribute{
		"name": schema.StringAttribute{
			Required:            true,
			MarkdownDescription: "xpath is: `config.Groups.V_units.V_name`",
		},
		"description": schema.StringAttribute{
			Optional:            true,
			MarkdownDescription: "xpath is `config.Groups.V_units.V_description`",
		},
		"vlan_id": schema.Int32Attribute{
			Optional:            true,
			MarkdownDescription: "xpath is `config.Groups.V_units.V_vlan-id`",
		},
		"family": schema.ListNestedAttribute{
			Optional: true,
			NestedObject: schema.NestedAttributeObject{
				Attributes: FamilyModel{}.Attributes(),
			},
		},
	}
}

func (o UnitsModel) AttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"name":        types.StringType,
		"description": types.StringType,
		"vlan_id":     types.StringType,
		"family":      types.ListType{ElemType: types.ObjectType{AttrTypes: FamilyModel{}.AttrTypes()}},
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

func (o FamilyModel) AttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"inet":  types.ListType{ElemType: types.ObjectType{AttrTypes: InetModel{}.AttrTypes()}},
		"inet6": types.ListType{ElemType: types.ObjectType{AttrTypes: Inet6Model{}.AttrTypes()}},
	}
}

func (o InetModel) Attributes() map[string]schema.Attribute {
	return map[string]schema.Attribute{
		"address": schema.StringAttribute{
			Optional:            true,
			MarkdownDescription: "xpath is `config.Groups.V_units.V_family.V_inet.V_address`",
		},
	}
}

func (o InetModel) AttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"address": types.StringType,
	}
}

func (o Inet6Model) Attributes() map[string]schema.Attribute {
	return map[string]schema.Attribute{
		"address": schema.StringAttribute{
			Optional:            true,
			MarkdownDescription: "xpath is `config.Groups.V_units.V_family.V_inet.V_address`",
		},
	}
}

func (o Inet6Model) AttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"address": types.StringType,
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
				MarkdownDescription: "xpath is: `config.Groups.V_interface`",
			},
			"description": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "xpath is: `config.Groups.V_description`",
			},
			"mtu": schema.Int64Attribute{
				Optional:            true,
				MarkdownDescription: "xpath is: `config.Groups.V_mtu`",
			},
			"vlan_tagging": schema.BoolAttribute{
				Optional:            true,
				MarkdownDescription: "xpath is: `config.Groups.V_vlan_tagging`",
			},
			"units": schema.ListNestedAttribute{
				Optional: true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: UnitsModel{}.Attributes(),
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
	config.Groups.V_interface.V_name = plan.Name.ValueStringPointer()
	config.Groups.V_interface.V_description = plan.Description.ValueStringPointer()
	config.Groups.V_interface.V_mtu = plan.Mtu.ValueInt64Pointer()
	if plan.Vlan_tagging.ValueBool() {
		empty := ""
		config.Groups.V_interface.V_vlan_tagging = &empty
	}
	for _, unit := range plan.Units {
		config.Groups.V_interface.V_units.V_name = unit.Name.ValueStringPointer()
		config.Groups.V_interface.V_units.V_description = unit.Description.ValueStringPointer()
		config.Groups.V_interface.V_units.V_vlan_id = unit.Vlan_id.ValueInt32Pointer()
		var inets []InetModel
		resp.Diagnostics.Append(unit.Family[0].Inet.ElementsAs(ctx, &inets, false)...)
		if resp.Diagnostics.HasError() {
			return
		}
		for _, inet := range inets {
			config.Groups.V_interface.V_units.V_family.V_inet.V_address.V_name = inet.Address.ValueStringPointer()
		}
		for _, addrv6 := range unit.Family[0].Inet6 {
			config.Groups.V_interface.V_units.V_family.V_inet6.V_address.V_name = addrv6.Address.ValueStringPointer()
		}

		err := r.client.SendTransaction("", config, false)
		if err != nil {
			resp.Diagnostics.AddError("Failed while Sending", err.Error())
			return
		}
		resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
	}
	// resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)

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
	id := config.Groups.Name
	err := r.client.MarshalGroup(id, config)
	if err != nil {
		resp.Diagnostics.AddError("Failed while Reading", err.Error())
		return
	}

	var newstate InterfacesModel
	newstate.ResourceName = types.StringValue(config.Groups.Name)
	newstate.Name = types.StringPointerValue(config.Groups.V_interface.V_name)
	newstate.Description = types.StringPointerValue(config.Groups.V_interface.V_description)
	var units []UnitsModel
	for i, xmlunit := range config.Groups.V_interface.V_units {
		var unit UnitsModel
		for m, xmlfamily := range xmlunit.V_family {

		}
		units = append(units, unit)
	}
	units_list, diags := types.ListValueFrom(ctx, types.ObjectType{AttrTypes: UnitsModel{}.AttrTypes()}, units)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	newstate.Units = units_list
	config.Groups.V_interface.V_mtu = plan.Mtu.ValueInt64Pointer()

	// Check values
	if err := resp.State.Set(ctx, config.Groups.V_interface.V_name); err != nil {
		return
	}
	if err := resp.State.Set(ctx, config.Groups.V_interface.V_description); err != nil {
		return
	}
	if err := resp.State.Set(ctx, config.Groups.V_interface.V_mtu); err != nil {
		return
	}
	if err := resp.State.Set(ctx, config.Groups.V_interface.V_units.V_name); err != nil {
		return
	}
	if err := resp.State.Set(ctx, config.Groups.V_interface.V_units.V_description); err != nil {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

// Update implements resource.Resource.
func (r *resourceInterfaces) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Get the data and set
	var plan InterfacesModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	// Check for errors
	if resp.Diagnostics.HasError() {
		return
	}

	var config xmlInterfaces
	config.Groups.Name = plan.ResourceName.ValueString()
	config.Groups.V_interface.V_name = plan.Name.ValueStringPointer()
	config.Groups.V_interface.V_description = plan.Description.ValueStringPointer()
	config.Groups.V_interface.V_mtu = plan.Mtu.ValueInt64Pointer()
	if plan.Vlan_tagging.ValueBool() {
		empty := ""
		config.Groups.V_interface.V_vlan_tagging = &empty
	}
	for _, unit := range plan.Units {
		config.Groups.V_interface.V_units.V_name = unit.Name.ValueStringPointer()
		config.Groups.V_interface.V_units.V_description = unit.Description.ValueStringPointer()
		config.Groups.V_interface.V_units.V_vlan_id = unit.Vlan_id.ValueInt32Pointer()
		for _, addr := range unit.Family[0].Inet {
			config.Groups.V_interface.V_units.V_family.V_inet.V_address.V_name = addr.Address.ValueStringPointer()
		}
		for _, addrv6 := range unit.Family[0].Inet6 {
			config.Groups.V_interface.V_units.V_family.V_inet6.V_address.V_name = addrv6.Address.ValueStringPointer()
		}

		err := r.client.SendTransaction("", config, false)
		if err != nil {
			resp.Diagnostics.AddError("Failed while Sending", err.Error())
			return
		}
		resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
	}
	// resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
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
