

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
				Description         *string  `xml:"description,omitempty"`
				Vlan_tagging         *string  `xml:"vlan-tagging,omitempty"`
				Mtu         *Int32  `xml:"mtu,omitempty"`
				Unit struct {
			 	Name         *string  `xml:"name,omitempty"`
			 	Description         *string  `xml:"description,omitempty"`
			 	Vlan_id         *string  `xml:"vlan-id,omitempty"`
				Family struct {
					Inet struct {
					Address struct {
			 	Name         *string  `xml:"name,omitempty"`
		}
		}
					Inet6 struct {
					Address struct {
			 	Name         *string  `xml:"name,omitempty"`
		}
		}
		}
		}
		} `xml:"interfaces>interface"`
	} `xml:"groups"`
}

// Collecting objects from the .tf file
type InterfaceModel struct {
	ResourceName	types.String `tfsdk:"resource_name"`
	Name	types.String `tfsdk:"name"`
	Description	types.String `tfsdk:"description"`
	Mtu	types.Int64 `tfsdk:"mtu"`
	Unit	[]UnitModel `tfsdk:"unit"`
	
}
func (o InterfaceModel) AttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"name": 	types.StringType,
		"description": 	types.StringType,
		"mtu": 	types.Int64Type,
	}
}
func (o InterfaceModel) Attributes() map[string]schema.Attribute {
	return map[string]schema.Attribute{
		"name": schema.StringAttribute{
			Optional: true,
			MaarkdownDescription: "xpth is `config.Groups.Interface.Name",
		},
		"description": schema.StringAttribute{
			Optional: true,
			MaarkdownDescription: "xpth is `config.Groups.Interface.Description",
		},
		"mtu": schema.Int64Attribute{
			Optional: true,
			MarkdownDescription: "xpath is `config.Groups.Interface.Mtu",
		},
		"unit": schema.NestedAttribute{
			Optional: true,
			NestedObject: schema.NestedAttributeObject{
				Attributes: UnitModel{}.Attributes(),
			},
		},
	}
}

type UnitModel struct {
	Name	types.String `tfsdk:"name"`
	Description	types.String `tfsdk:"description"`
	Vlan_id	types.String `tfsdk:"vlan-id"`
	Family	[]FamilyModel `tfsdk:"family"`
}
func (o UnitModel) AttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
	"name": 	types.StringType,
	"description": 	types.StringType,
	"vlan-id": 	types.StringType,
	}
}
func (o UnitModel) Attributes() map[string]schema.Attribute {
	return map[string]schema.Attribute{
	"name": schema.StringAttribute{
		Optional: true,
		MarkdownDescription: "xpath is `config.Groups.Unit.Name`",
	},
	"description": schema.StringAttribute{
		Optional: true,
		MarkdownDescription: "xpath is `config.Groups.Unit.Description`",
	},
	"vlan-id": schema.StringAttribute{
		Optional: true,
		MarkdownDescription: "xpath is `config.Groups.Unit.Vlan-id`",
	},
	"family": schema.NestedAttribute{
		Optional: true,
		NestedObject: schema.NestedAttributeObject{
			Attributes: FamilyModel{}.Attributes(),
		},
	},
	}
}

type FamilyModel struct {
	Inet	[]InetModel `tfsdk:"inet"`
	Inet6	[]Inet6Model `tfsdk:"inet6"`
}
func (o FamilyModel) AttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
	}
}
func (o FamilyModel) Attributes() map[string]schema.Attribute {
	return map[string]schema.Attribute{
	"inet": schema.NestedAttribute{
		Optional: true,
		NestedObject: schema.NestedAttributeObject{
			Attributes: InetModel{}.Attributes(),
		},
	},
	"inet6": schema.NestedAttribute{
		Optional: true,
		NestedObject: schema.NestedAttributeObject{
			Attributes: Inet6Model{}.Attributes(),
		},
	},
	}
}

type InetModel struct {
	Address	[]AddressModel `tfsdk:"address"`
}
func (o InetModel) AttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
	}
}
func (o InetModel) Attributes() map[string]schema.Attribute {
	return map[string]schema.Attribute{
	"address": schema.NestedAttribute{
		Optional: true,
		NestedObject: schema.NestedAttributeObject{
			Attributes: AddressModel{}.Attributes(),
		},
	},
	}
}

type AddressModel struct {
	Name	types.String `tfsdk:"name"`
}
func (o AddressModel) AttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
	"name": 	types.StringType,
	}
}
func (o AddressModel) Attributes() map[string]schema.Attribute {
	return map[string]schema.Attribute{
	"name": schema.StringAttribute{
		Optional: true,
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
			"interface": schema.NestedAttribute{
				Optional: true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: InterfaceModel{}.Attributes(),
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
	config.Groups.Interface.Vlan-tagging = plan.Name.ValueStringPointer()
	config.Groups.Interface.Mtu = plan.Mtu.ValueInt64Pointer()
	
	var units []UnitModel
	resp.Diagnostics.Append(interface.unit.ElementsAs(ctx, &units, false)...)
	if resp.Diagnostics.HasError() {
		return
	}
	config.Groups.Interface.Unit.Name = unit.name.ValueStringPointer()
	config.Groups.Interface.Unit.Description = unit.description.ValueStringPointer()
	config.Groups.Interface.Unit.Vlan_id = unit.vlan-id.ValueStringPointer()
	
	var familys []FamilyModel
	resp.Diagnostics.Append(unit.family.ElementsAs(ctx, &familys, false)...)
	if resp.Diagnostics.HasError() {
		return
	}
	
	var inets []InetModel
	resp.Diagnostics.Append(family.inet.ElementsAs(ctx, &inets, false)...)
	if resp.Diagnostics.HasError() {
		return
	}
	
	var addresss []AddressModel
	resp.Diagnostics.Append(inet.address.ElementsAs(ctx, &addresss, false)...)
	if resp.Diagnostics.HasError() {
		return
	}
	config.Groups.Interface.Unit.Family.Inet.Address.Name = address.name.ValueStringPointer()
	
	var inet6s []Inet6Model
	resp.Diagnostics.Append(family.inet6.ElementsAs(ctx, &inet6s, false)...)
	if resp.Diagnostics.HasError() {
		return
	}
	
	var addresss []AddressModel
	resp.Diagnostics.Append(inet6.address.ElementsAs(ctx, &addresss, false)...)
	if resp.Diagnostics.HasError() {
		return
	}
	config.Groups.Interface.Unit.Family.Inet6.Address.Name = address.name.ValueStringPointer()
	
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
	id := config.Groups.Name
	err := r.client.MarshalGroup(id, config)
	if err != nil {
		resp.Diagnostics.AddError("Failed while Reading", err.Error())
		return
	}
	var newstate InterfacesModel
	newstate.ResourceName = types.StringValue(config.Groups.Name)
	newstate.Name = types.StringPointerValue(config.Groups.Interface.Name)
	newstate.Description = types.StringPointerValue(config.Groups.Interface.Description)
	newstate.Vlan-tagging = types.StringPointerValue(config.Groups.Interface.Vlan-tagging)
	newstate.Mtu = types.Int64PointerValue(config.Groups.Interface.Mtu)
	var units []UnitModel
	xmlunits = config.Groups.Unit[0]
	units = append(units, unit)
	var familys []FamilyModel
	xmlfamilys = config.Groups.Family[0]
	familys = append(familys, family)
	var inets []InetModel
	xmlinets = config.Groups.Inet[0]
	inets = append(inets, inet)
	var addresss []AddressModel
	xmladdresss = config.Groups.Address[0]
	addresss = append(addresss, address)
	newstate.name = types.StringPointerValue(config.Groups.Address.Name)
	var inet6s []Inet6Model
	xmlinet6s = config.Groups.Inet6[0]
	inet6s = append(inet6s, inet6)
	var addresss []AddressModel
	xmladdresss = config.Groups.Address[0]
	addresss = append(addresss, address)
	newstate.name = types.StringPointerValue(config.Groups.Address.Name)
	unit_list, diags := types.ListValueFrom(ctx, types.ObjectType{AttrTypes: unitsModel{}.AttrTypes()}, units)
	if resp.Diagnostics.HasError() {
		return
	}
	newstate.Units = units_list
	if err := resp.State.Set(ctx, config.Groups.interface.name); err != nil {
		return
	}
	if err := resp.State.Set(ctx, config.Groups.interface.description); err != nil {
		return
	}
	if err := resp.State.Set(ctx, config.Groups.interface.vlan-tagging); err != nil {
		return
	}
	if err := resp.State.Set(ctx, config.Groups.interface.mtu); err != nil {
		return
	}
	if err := resp.State.Set(ctx, config.Groups.interface.unit); err != nil {
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
	config.Groups.Interface.Name = plan.Name.ValueStringPointer()
	config.Groups.Interface.Description = plan.Name.ValueStringPointer()
	config.Groups.Interface.Vlan-tagging = ''
	config.Groups.Interface.Mtu = plan.Mtu.ValueInt64Pointer()
	
	var units []UnitModel
	resp.Diagnostics.Append(interface.unit.ElementsAs(ctx, &units, false)...)
	if resp.Diagnostics.HasError() {
		return
	}
	config.Groups.Interface.Unit.Name = unit.name.ValueStringPointer()
	config.Groups.Interface.Unit.Description = unit.description.ValueStringPointer()
	config.Groups.Interface.Unit.Vlan_id = unit.vlan-id.ValueStringPointer()
	
	var familys []FamilyModel
	resp.Diagnostics.Append(unit.family.ElementsAs(ctx, &familys, false)...)
	if resp.Diagnostics.HasError() {
		return
	}
	
	var inets []InetModel
	resp.Diagnostics.Append(family.inet.ElementsAs(ctx, &inets, false)...)
	if resp.Diagnostics.HasError() {
		return
	}
	
	var addresss []AddressModel
	resp.Diagnostics.Append(inet.address.ElementsAs(ctx, &addresss, false)...)
	if resp.Diagnostics.HasError() {
		return
	}
	config.Groups.Interface.Unit.Family.Inet.Address.Name = address.name.ValueStringPointer()
	
	var inet6s []Inet6Model
	resp.Diagnostics.Append(family.inet6.ElementsAs(ctx, &inet6s, false)...)
	if resp.Diagnostics.HasError() {
		return
	}
	
	var addresss []AddressModel
	resp.Diagnostics.Append(inet6.address.ElementsAs(ctx, &addresss, false)...)
	if resp.Diagnostics.HasError() {
		return
	}
	config.Groups.Interface.Unit.Family.Inet6.Address.Name = address.name.ValueStringPointer()
	
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
