package main

import (
	"context"
	"encoding/xml"
	"fmt"
	"os"
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
		XMLName   xml.Name                 `xml:"groups"`
		Name      string                   `xml:"name"`
		Interface []xmlInterfacesInterface `xml:"interfaces>interface"`
	} `xml:"groups"`
}
type xmlInterfacesInterface struct {
	XMLName      xml.Name           `xml:"interface"`
	Name         *string            `xml:"name,omitempty"`
	Description  *string            `xml:"description,omitempty"`
	Vlan_tagging *string            `xml:"vlan-tagging,omitempty"`
	Mtu          *int64             `xml:"mtu,omitempty"`
	Unit         []xmlInterfaceUnit `xml:"unit,omitempty"`
}
type xmlInterfaceUnit struct {
	XMLName     xml.Name        `xml:"unit"`
	Name        *string         `xml:"name,omitempty"`
	Description *string         `xml:"description,omitempty"`
	Vlan_id     *string         `xml:"vlan-id,omitempty"`
	Family      []xmlUnitFamily `xml:"family,omitempty"`
}
type xmlUnitFamily struct {
	XMLName xml.Name         `xml:"family"`
	Inet    []xmlFamilyInet  `xml:"inet,omitempty"`
	Inet6   []xmlFamilyInet6 `xml:"inet6,omitempty"`
}
type xmlFamilyInet struct {
	XMLName xml.Name         `xml:"inet"`
	Address []xmlInetAddress `xml:"address,omitempty"`
}
type xmlInetAddress struct {
	XMLName xml.Name `xml:"address"`
	Name    *string  `xml:"name,omitempty"`
}
type xmlFamilyInet6 struct {
	XMLName xml.Name          `xml:"inet6"`
	Address []xmlInet6Address `xml:"address,omitempty"`
}
type xmlInet6Address struct {
	XMLName xml.Name `xml:"address"`
	Name    *string  `xml:"name,omitempty"`
}

// Collecting objects from the .tf file
type InterfacesModel struct {
	ResourceName types.String `tfsdk:"resource_name"`
	Interface    types.List   `tfsdk:"interface"`
}

func (o InterfacesModel) AttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"interface": types.ListType{ElemType: types.ObjectType{AttrTypes: InterfaceModel{}.AttrTypes()}},
	}
}
func (o InterfacesModel) Attributes() map[string]schema.Attribute {
	return map[string]schema.Attribute{
		"interface": schema.ListNestedAttribute{
			Optional: true,
			NestedObject: schema.NestedAttributeObject{
				Attributes: InterfaceModel{}.Attributes(),
			},
		},
	}
}

type InterfaceModel struct {
	Name         types.String `tfsdk:"name"`
	Description  types.String `tfsdk:"description"`
	Vlan_tagging types.Bool   `tfsdk:"vlan_tagging"`
	Mtu          types.Int64  `tfsdk:"mtu"`
	Unit         types.List   `tfsdk:"unit"`
}

func (o InterfaceModel) AttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"name":         types.StringType,
		"description":  types.StringType,
		"vlan_tagging": types.BoolType,
		"mtu":          types.Int64Type,
		"unit":         types.ListType{ElemType: types.ObjectType{AttrTypes: UnitModel{}.AttrTypes()}},
	}
}
func (o InterfaceModel) Attributes() map[string]schema.Attribute {
	return map[string]schema.Attribute{
		"name": schema.StringAttribute{
			Optional:            true,
			MarkdownDescription: "xpath is `config.Groups.Interface.Name`",
		},
		"description": schema.StringAttribute{
			Optional:            true,
			MarkdownDescription: "xpath is `config.Groups.Interface.Description`",
		},
		"vlan_tagging": schema.BoolAttribute{
			Optional:            true,
			MarkdownDescription: "xpath is `config.Groups.Interface.Vlan_tagging`",
		},
		"mtu": schema.Int64Attribute{
			Optional:            true,
			MarkdownDescription: "xpath is `config.Groups.Interface.Mtu`",
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
	Vlan_id     types.String `tfsdk:"vlan_id"`
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
			MarkdownDescription: "xpath is `config.Groups.Unit.Vlan_id`",
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

	var var_interfaces_interface []InterfaceModel
	resp.Diagnostics.Append(plan.Interface.ElementsAs(ctx, &var_interfaces_interface, false)...)
	if resp.Diagnostics.HasError() {
		return
	}
	config.Groups.Interface = make([]xmlInterfacesInterface, len(var_interfaces_interface))
	for i_interface, v_interface := range var_interfaces_interface {
		config.Groups.Interface[i_interface].Name = v_interface.Name.ValueStringPointer()
		config.Groups.Interface[i_interface].Description = v_interface.Description.ValueStringPointer()
		if v_interface.Vlan_tagging.ValueBool() {
			empty := ""
			config.Groups.Interface[i_interface].Vlan_tagging = &empty
		}
		config.Groups.Interface[i_interface].Mtu = v_interface.Mtu.ValueInt64Pointer()

		var var_interface_unit []UnitModel
		resp.Diagnostics.Append(v_interface.Unit.ElementsAs(ctx, &var_interface_unit, false)...)
		if resp.Diagnostics.HasError() {
			return
		}
		config.Groups.Interface[i_interface].Unit = make([]xmlInterfaceUnit, len(var_interface_unit))
		for i_unit, v_unit := range var_interface_unit {
			config.Groups.Interface[i_interface].Unit[i_unit].Name = v_unit.Name.ValueStringPointer()
			config.Groups.Interface[i_interface].Unit[i_unit].Description = v_unit.Description.ValueStringPointer()
			config.Groups.Interface[i_interface].Unit[i_unit].Vlan_id = v_unit.Vlan_id.ValueStringPointer()

			var var_unit_family []FamilyModel
			resp.Diagnostics.Append(v_unit.Family.ElementsAs(ctx, &var_unit_family, false)...)
			if resp.Diagnostics.HasError() {
				return
			}
			config.Groups.Interface[i_interface].Unit[i_unit].Family = make([]xmlUnitFamily, len(var_unit_family))
			for i_family, v_family := range var_unit_family {

				var var_family_inet []InetModel
				resp.Diagnostics.Append(v_family.Inet.ElementsAs(ctx, &var_family_inet, false)...)
				if resp.Diagnostics.HasError() {
					return
				}
				config.Groups.Interface[i_interface].Unit[i_unit].Family[i_family].Inet = make([]xmlFamilyInet, len(var_family_inet))
				for i_inet, v_inet := range var_family_inet {

					var var_inet_address []AddressModel
					resp.Diagnostics.Append(v_inet.Address.ElementsAs(ctx, &var_inet_address, false)...)
					if resp.Diagnostics.HasError() {
						return
					}
					config.Groups.Interface[i_interface].Unit[i_unit].Family[i_family].Inet[i_inet].Address = make([]xmlInetAddress, len(var_inet_address))
					for i_address, v_address := range var_inet_address {
						config.Groups.Interface[i_interface].Unit[i_unit].Family[i_family].Inet[i_inet].Address[i_address].Name = v_address.Name.ValueStringPointer()
					}

				}

				var var_family_inet6 []Inet6Model
				resp.Diagnostics.Append(v_family.Inet6.ElementsAs(ctx, &var_family_inet6, false)...)
				if resp.Diagnostics.HasError() {
					return
				}
				config.Groups.Interface[i_interface].Unit[i_unit].Family[i_family].Inet6 = make([]xmlFamilyInet6, len(var_family_inet6))
				for i_inet6, v_inet6 := range var_family_inet6 {

					var var_inet6_address []AddressModel
					resp.Diagnostics.Append(v_inet6.Address.ElementsAs(ctx, &var_inet6_address, false)...)
					if resp.Diagnostics.HasError() {
						return
					}
					config.Groups.Interface[i_interface].Unit[i_unit].Family[i_family].Inet6[i_inet6].Address = make([]xmlInet6Address, len(var_inet6_address))
					for i_address, v_address := range var_inet6_address {
						config.Groups.Interface[i_interface].Unit[i_unit].Family[i_family].Inet6[i_inet6].Address[i_address].Name = v_address.Name.ValueStringPointer()
					}

				}

			}

		}

	}

	err := r.client.SendTransaction("", config, false)
	if err != nil {
		resp.Diagnostics.AddError("Failed while Sending", err.Error())
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

// ###################################################################################################################################################

// Read implements resource.Resource.
func (r *resourceInterfaces) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Get the data and set
	var state InterfacesModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	// Check for errors
	if resp.Diagnostics.HasError() {
		return
	}
	fmt.Fprintf(os.Stderr, "State: %#v", state)

	// Marshall group and check
	var config xmlInterfaces
	err := r.client.MarshalGroup(state.ResourceName.ValueString(), &config)
	if err != nil {
		resp.Diagnostics.AddError("Failed while Reading", err.Error())
		return
	}
	fmt.Fprintf(os.Stderr, "Config: %#v", config)

	// For now just see what in config and print to see what it look like

	// config has the data
	// I need to populate state with data from config similar to create (but reverse)

	// Begin reconstructing the model from Junos config
	state.Interface = types.ListNull(types.ObjectType{AttrTypes: InterfacesModel{}.AttrTypes()}) // Clear and rebuild

	interfaceList := make([]InterfaceModel, len(config.Groups.Interface))
	for i, iface := range config.Groups.Interface {
		var interfaceModel InterfaceModel
		interfaceModel.Name = types.StringPointerValue(iface.Name)
		interfaceModel.Description = types.StringPointerValue(iface.Description)
		interfaceModel.Vlan_tagging = types.BoolValue(iface.Vlan_tagging != nil)
		interfaceModel.Mtu = types.Int64PointerValue(iface.Mtu)

		// Units
		unitList := make([]UnitModel, len(iface.Unit))
		for j, unit := range iface.Unit {
			var unitModel UnitModel
			unitModel.Name = types.StringPointerValue(unit.Name)
			unitModel.Description = types.StringPointerValue(unit.Description)
			unitModel.Vlan_id = types.StringPointerValue(unit.Vlan_id)

			// Families
			familyList := make([]FamilyModel, len(unit.Family))
			for k, fam := range unit.Family {
				var familyModel FamilyModel

				// INET
				inetList := make([]InetModel, len(fam.Inet))
				for l, inet := range fam.Inet {
					var inetModel InetModel
					addressList := make([]AddressModel, len(inet.Address))
					for m, addr := range inet.Address {
						addressList[m].Name = types.StringPointerValue(addr.Name)
					}
					inetModel.Address, _ = types.ListValueFrom(ctx, types.ObjectType{AttrTypes: AddressModel{}.AttrTypes()}, addressList)
					inetList[l] = inetModel
				}
				familyModel.Inet, _ = types.ListValueFrom(ctx, types.ObjectType{AttrTypes: InetModel{}.AttrTypes()}, inetList)

				// INET6
				inet6List := make([]Inet6Model, len(fam.Inet6))
				for l, inet6 := range fam.Inet6 {
					var inet6Model Inet6Model
					addressList := make([]AddressModel, len(inet6.Address))
					for m, addr := range inet6.Address {
						addressList[m].Name = types.StringPointerValue(addr.Name)
					}
					inet6Model.Address, _ = types.ListValueFrom(ctx, types.ObjectType{AttrTypes: AddressModel{}.AttrTypes()}, addressList)
					inet6List[l] = inet6Model
				}
				familyModel.Inet6, _ = types.ListValueFrom(ctx, types.ObjectType{AttrTypes: Inet6Model{}.AttrTypes()}, inet6List)

				familyList[k] = familyModel
			}
			unitModel.Family, _ = types.ListValueFrom(ctx, types.ObjectType{AttrTypes: FamilyModel{}.AttrTypes()}, familyList)
			unitList[j] = unitModel
		}

		interfaceModel.Unit, _ = types.ListValueFrom(ctx, types.ObjectType{AttrTypes: UnitModel{}.AttrTypes()}, unitList)
		interfaceList[i] = interfaceModel
	}

	// Set interface list into state
	state.Interface, _ = types.ListValueFrom(ctx, types.ObjectType{AttrTypes: InterfaceModel{}.AttrTypes()}, interfaceList)

	// Write back to Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

// For now just see what in config and print to see what it look like

// config has the data
// I need to populate state with data from config similar to create (but reverse)

//state.Interface = config.Groups.Interface
//newState := InterfacesModel{}
//make([]xmlInterfaces)

// for i, m_interface := range config.Groups.Interface {
// 	interfacesModelInterface := types.List{}
// 	state.Interface[i].Name = types.StringPointerValue(m_interface.Name)
// 	state.Interface[i].Description = types.StringPointerValue(m_interface.Description)
// 	if m_interface.Vlan_tagging.ValueBool() {
// 		empty := ""
// 		state.Interface[i].Vlan_tagging = &empty
// 	}
// 	state.Interface[i].Mtu = types.Int64PointerValue(m_interface.Mtu)

// 	state.Interface[i].Unit := make([]UnitsModel, len(config.Groups.Interface[i].Unit))

// 	for u, xmlUnit := range config.Groups.Interface[i].Unit {
// 		state.Interface[i].Unit[u].Name = types.StringPointerValue(xmlUnit.Name)
// 		state.Interface[i].Unit[u].Description = types.StringPointerValue(xmlUnit.Description)
// 		state.Interface[i].Unit[u].Vlan_id = types.Int32PointerValue(xmlUnit.Vlan_id)

// 		state.Interface[i].Unit[u].Family := make([]FamilyModel, len(config.Groups.Interface[i].Unit[u].Family))
// 		for f, m_family := range config.Groups.Interface[i].Unit[u].Family {

// 			state.Interface[i].Unit[u].Family[f].Inet := make([]InetModel, len(config.Groups.Interface[i].Unit[u].Family[f].Inet))
// 			for in, m_inet := range config.Groups.Interface[i].Unit[u].Family[f].Inet {

// 				state.Interface[i].Unit[u].Family[f].Inet[in].Address := make([]AddressModel, len(config.Groups.Interface[i].Unit[u].Family[f].Inet[in]))
// 				for n, m_addrname := range config.Groups.Interface[i].Unit[u].Family[f].Inet[in].Address {
// 					state.Interface[i].Unit[u].Family[f].Inet[in].Address[n].Name = types.StringPointerValue(m_addrname.Name)
// 				}
// 			}

// 			state.Interface[i].Unit[u].Family[f].Inet6 := make([]Inet6Model, len(config.Groups.Interface[i].Unit[u].Family[f].Inet6))
// 			for in6, m_inet6 := range config.Groups.Interface[i].Unit[u].Family[f].Inet6 {

// 				state.Interface[i].Unit[u].Family[f].Inet6[in6].Address := make([]AddressModel, len(config.Groups.Interface[i].Unit[u].Family[f].Inet6[in6]))
// 				for n, m_addrname := range config.Groups.Interface[i].Unit[u].Family[f].Inet6[in6].Address {
// 					state.Interface[i].Unit[u].Family[f].Inet6[in6].Address[n].Name = types.StringPointerValue(m_addrname.Name)
// 				}
// 			}
// 		}
// 	}
// }

// 	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
// }

// ###################################################################################################################################################
// ###################################################################################################################################################
// ###################################################################################################################################################

// func (r *resourceInterfaces) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
// 	var state InterfacesModel
// 	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
// 	if resp.Diagnostics.HasError() {
// 		return
// 	}

// 	var config xmlInterfaces
// 	err := r.client.MarshalGroup(state.ResourceName.ValueString(), &config)
// 	if err != nil {
// 		resp.Diagnostics.AddError("Read Error", fmt.Sprintf("Unable to read config from device: %s", err))
// 		return
// 	}

// 	var interfaceModels []InterfaceModel
// 	for _, iface := range config.Groups.Interface {
// 		ifaceModel := InterfaceModel{
// 			Name:         types.StringPointerValue(iface.Name),
// 			Description:  types.StringNull(),
// 			Vlan_tagging: types.BoolValue(iface.Vlan_tagging != nil),
// 			Mtu:          types.Int64Null(),
// 		}
// 		if iface.Description != nil {
// 			ifaceModel.Description = types.StringValue(*iface.Description)
// 		}
// 		if iface.Mtu != nil {
// 			ifaceModel.Mtu = types.Int64Value(int64(*iface.Mtu))
// 		}

// 		var unitModels []UnitModel
// 		for _, unit := range iface.Unit {
// 			unitModel := UnitModel{
// 				Name:        types.StringPointerValue(unit.Name),
// 				Description: types.StringNull(),
// 				Vlan_id:     types.StringNull(),
// 			}
// 			if unit.Description != nil {
// 				unitModel.Description = types.StringValue(*unit.Description)
// 			}
// 			if unit.Vlan_id != nil {
// 				unitModel.Vlan_id = types.StringValue(*unit.Vlan_id)
// 			}

// 			// Families
// 			familyModel := FamilyModel{}

// 			if len(unit.Family.Inet.Address) > 0 {
// 				var inetAddrs []AddressModel
// 				for _, addr := range unit.Family.Inet.Address {
// 					if addr.Name != nil {
// 						inetAddrs = append(inetAddrs, AddressModel{
// 							Name: types.StringValue(*addr.Name),
// 						})
// 					}
// 				}
// 				listVal, diags := types.ListValueFrom(ctx, types.ObjectType{AttrTypes: AddressModel{}.AttrTypes()}, inetAddrs)
// 				resp.Diagnostics.Append(diags...)
// 				familyModel.Inet = listVal
// 			} else {
// 				familyModel.Inet = types.ListNull(types.ObjectType{AttrTypes: AddressModel{}.AttrTypes()})
// 			}

// 			if len(unit.Family.Inet6.Address) > 0 {
// 				var inet6Addrs []AddressModel
// 				for _, addr := range unit.Family.Inet6.Address {
// 					if addr.Name != nil {
// 						inet6Addrs = append(inet6Addrs, AddressModel{
// 							Name: types.StringValue(*addr.Name),
// 						})
// 					}
// 				}
// 				listVal, diags := types.ListValueFrom(ctx, types.ObjectType{AttrTypes: AddressModel{}.AttrTypes()}, inet6Addrs)
// 				resp.Diagnostics.Append(diags...)
// 				familyModel.Inet6 = listVal
// 			} else {
// 				familyModel.Inet6 = types.ListNull(types.ObjectType{AttrTypes: AddressModel{}.AttrTypes()})
// 			}

// 			// Wrap family in a list
// 			familyList, diags := types.ListValueFrom(ctx, types.ObjectType{AttrTypes: FamilyModel{}.AttrTypes()}, []FamilyModel{familyModel})
// 			resp.Diagnostics.Append(diags...)
// 			unitModel.Family = familyList

// 			unitModels = append(unitModels, unitModel)
// 		}

// 		unitList, diags := types.ListValueFrom(ctx, types.ObjectType{AttrTypes: UnitModel{}.AttrTypes()}, unitModels)
// 		resp.Diagnostics.Append(diags...)
// 		ifaceModel.Unit = unitList

// 		interfaceModels = append(interfaceModels, ifaceModel)
// 	}

// 	intfList, diags := types.ListValueFrom(ctx, types.ObjectType{AttrTypes: InterfaceModel{}.AttrTypes()}, interfaceModels)
// 	resp.Diagnostics.Append(diags...)
// 	state.Interface = intfList

// 	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
// }

// ###################################################################################################################################################
// ###################################################################################################################################################
// ###################################################################################################################################################
// ###################################################################################################################################################
// ###################################################################################################################################################
// ###################################################################################################################################################

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

	var var_interfaces_interface []InterfaceModel
	resp.Diagnostics.Append(plan.Interface.ElementsAs(ctx, &var_interfaces_interface, false)...)
	if resp.Diagnostics.HasError() {
		return
	}
	config.Groups.Interface = make([]xmlInterfacesInterface, len(var_interfaces_interface))
	for i_interface, v_interface := range var_interfaces_interface {
		config.Groups.Interface[i_interface].Name = v_interface.Name.ValueStringPointer()
		config.Groups.Interface[i_interface].Description = v_interface.Description.ValueStringPointer()
		if v_interface.Vlan_tagging.ValueBool() {
			empty := ""
			config.Groups.Interface[i_interface].Vlan_tagging = &empty
		}
		config.Groups.Interface[i_interface].Mtu = v_interface.Mtu.ValueInt64Pointer()

		var var_interface_unit []UnitModel
		resp.Diagnostics.Append(v_interface.Unit.ElementsAs(ctx, &var_interface_unit, false)...)
		if resp.Diagnostics.HasError() {
			return
		}
		config.Groups.Interface[i_interface].Unit = make([]xmlInterfaceUnit, len(var_interface_unit))
		for i_unit, v_unit := range var_interface_unit {
			config.Groups.Interface[i_interface].Unit[i_unit].Name = v_unit.Name.ValueStringPointer()
			config.Groups.Interface[i_interface].Unit[i_unit].Description = v_unit.Description.ValueStringPointer()
			config.Groups.Interface[i_interface].Unit[i_unit].Vlan_id = v_unit.Vlan_id.ValueStringPointer()

			var var_unit_family []FamilyModel
			resp.Diagnostics.Append(v_unit.Family.ElementsAs(ctx, &var_unit_family, false)...)
			if resp.Diagnostics.HasError() {
				return
			}
			config.Groups.Interface[i_interface].Unit[i_unit].Family = make([]xmlUnitFamily, len(var_unit_family))
			for i_family, v_family := range var_unit_family {

				var var_family_inet []InetModel
				resp.Diagnostics.Append(v_family.Inet.ElementsAs(ctx, &var_family_inet, false)...)
				if resp.Diagnostics.HasError() {
					return
				}
				config.Groups.Interface[i_interface].Unit[i_unit].Family[i_family].Inet = make([]xmlFamilyInet, len(var_family_inet))
				for i_inet, v_inet := range var_family_inet {

					var var_inet_address []AddressModel
					resp.Diagnostics.Append(v_inet.Address.ElementsAs(ctx, &var_inet_address, false)...)
					if resp.Diagnostics.HasError() {
						return
					}
					config.Groups.Interface[i_interface].Unit[i_unit].Family[i_family].Inet[i_inet].Address = make([]xmlInetAddress, len(var_inet_address))
					for i_address, v_address := range var_inet_address {
						config.Groups.Interface[i_interface].Unit[i_unit].Family[i_family].Inet[i_inet].Address[i_address].Name = v_address.Name.ValueStringPointer()
					}

				}

				var var_family_inet6 []Inet6Model
				resp.Diagnostics.Append(v_family.Inet6.ElementsAs(ctx, &var_family_inet6, false)...)
				if resp.Diagnostics.HasError() {
					return
				}
				config.Groups.Interface[i_interface].Unit[i_unit].Family[i_family].Inet6 = make([]xmlFamilyInet6, len(var_family_inet6))
				for i_inet6, v_inet6 := range var_family_inet6 {

					var var_inet6_address []AddressModel
					resp.Diagnostics.Append(v_inet6.Address.ElementsAs(ctx, &var_inet6_address, false)...)
					if resp.Diagnostics.HasError() {
						return
					}
					config.Groups.Interface[i_interface].Unit[i_unit].Family[i_family].Inet6[i_inet6].Address = make([]xmlInet6Address, len(var_inet6_address))
					for i_address, v_address := range var_inet6_address {
						config.Groups.Interface[i_interface].Unit[i_unit].Family[i_family].Inet6[i_inet6].Address[i_address].Name = v_address.Name.ValueStringPointer()
					}

				}

			}

		}

	}

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
