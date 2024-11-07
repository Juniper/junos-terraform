
package main

import (
	"context"
	"encoding/xml"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
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

type xmlConfiguration struct {
	XMLName xml.Name `xml:"configuration"`
	Groups xmlGroups `xml:"groups"`
}

type xmlGroups struct {
	XMLName xml.Name `xml:"groups"`
	Name string   `xml:"name"`
	 xml `xml:""`
}
type xmlInterfaces struct {
    XMLName        xml.Name           `xml:"interfaces"`
    Interface []xmlInterface `xml:"interface"`
}
type xmlInterface struct {
    XMLName        xml.Name           `xml:"interface"`
    Name         *string            `xml:"name,omitempty"`
    Description         *string            `xml:"description,omitempty"`
    Mtu         *int64            `xml:"mtu,omitempty"`
    Vlan_tagging         *bool            `xml:"vlan-tagging,omitempty"`
    Unit []xmlUnit `xml:"unit"`
}
type xmlUnit struct {
    XMLName        xml.Name           `xml:"unit"`
    Name         *string            `xml:"name,omitempty"`
    Description         *string            `xml:"description,omitempty"`
    Vlan_id         *int32            `xml:"vlan-id,omitempty"`
}
// Collecting objects from the .tf file
type InterfacesModel struct {
	ResourceName types.String `tfsdk:"resource_name"`
	Interface        types.List `tfsdk:"interface"`
}
func (o InterfacesModel) Attributes() map[string]schema.Attributes {
    return map[string]schema.Attribute {
        "interface": schema.ListNestedAttribute {
            Optional : true,
            NestedObject: schema.NestedAttributeObject {
                Attributes: InterfaceModel{}.Attributes(),
            },
        },
    }
}
func (o InterfacesModel) AttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"interface":    types.ListType{ElemType: types.ObjectType{AttrTypes: Interface.AttrTypes()}}
	}
}
type InterfaceModel struct {
	Name        types.String `tfsdk:"name"`
	Description        types.String `tfsdk:"description"`
	Mtu        types.Int64 `tfsdk:"mtu"`
	Vlan_tagging        types.Bool `tfsdk:"vlan_tagging"`
	Unit        types.List `tfsdk:"unit"`
}
func (o InterfaceModel) Attributes() map[string]schema.Attributes {
    return map[string]schema.Attribute {
        "name": schema.StringAttribute {
            Required : true,
            MarkdownDescription: "xpath is `config.Groups.Interface.Name",
        }
        "description": schema.StringAttribute {
            Optional: true,
            MarkdownDescription: "xpath is `config.Groups.Interface.Description",
        }
        "mtu": schema.Int64Attribute {
            Optional: true,
            MarkdownDescription: "xpath is `config.Groups.Interface.Mtu",
        }
        "vlan_tagging": schema.BoolAttribute {
            Optional: true,
            MarkdownDescription: "xpath is `config.Groups.Interface.Vlan_tagging",
        }
        "unit": schema.ListNestedAttribute {
            Optional : true,
            NestedObject: schema.NestedAttributeObject {
                Attributes: UnitModel{}.Attributes(),
            },
        },
    }
}
func (o InterfaceModel) AttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"name":        types.StringType,
		"description":        types.StringType,
		"mtu":        types.Int64Type,
		"vlan_tagging":        types.BoolType,
		"unit":    types.ListType{ElemType: types.ObjectType{AttrTypes: Unit.AttrTypes()}}
	}
}
type UnitModel struct {
	Name        types.String `tfsdk:"name"`
	Description        types.String `tfsdk:"description"`
	Vlan_id        types.Int32 `tfsdk:"vlan_id"`
}
func (o UnitModel) Attributes() map[string]schema.Attributes {
    return map[string]schema.Attribute {
        "name": schema.StringAttribute {
            Required : true,
            MarkdownDescription: "xpath is `config.Groups.Unit.Name",
        }
        "description": schema.StringAttribute {
            Optional: true,
            MarkdownDescription: "xpath is `config.Groups.Unit.Description",
        }
        "vlan_id": schema.Int32Attribute {
            Optional: true,
            MarkdownDescription: "xpath is `config.Groups.Unit.Vlan_id",
        }
    }
}
func (o UnitModel) AttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"name":        types.StringType,
		"description":        types.StringType,
		"vlan_id":        types.Int32Type,
	}
}// Collects the data for the crud work
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
func (r *resource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Get the  Model data and set
	var plan Model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	// Check for errors
	if resp.Diagnostics.HasError() {
		return
	}
	var config xml
	config.Groups.Name = plan.ResourceName.ValueString()

    
    var interface xmlInterface
        config.Groups..Interface.Name = plan.Interface.Name.ValueStringPointer()
        config.Groups..Interface.Description = plan.Interface.Description.ValueStringPointer()
        config.Groups..Interface.Mtu = plan.Interface.Mtu.ValueInt64Pointer()
	if plan.Interface.Vlan_tagging.ValueBool(){
                empty := ""
                config.Groups.Interface.vlan_tagging = $empty
	}
	
            var unit []xmlUnit
            config.Groups..Unit = make([]xmlUnit, len(&unit))
            for m, n := range unit {
                config.Groups.Unit[m].Name = n.Name.ValueStringPointer()
                config.Groups.Unit[m].Description = n.Description.ValueStringPointer()
                config.Groups.Unit[m].Vlan_id = n.Vlan_id.ValueInt32Pointer()
            }


	err := r.client.SendTransaction("", config, false)
	if err != nil {
		resp.Diagnostics.AddError("Failed while Sending", err.Error())
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

// Read implements resource.Resource.
func (r *resource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var d diag.Diagnostics

	// Get the data and set
	var state Model
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	// Check for errors
	if resp.Diagnostics.HasError() {
		return
	}
        // Marshall group and check
	var config xml
	err := r.client.MarshalGroup(state.ResourceName.ValueString(), &config)
	if err != nil {
		resp.Diagnostics.AddError("Failed while Reading", err.Error())
		return
	}

	
    var interface xmlInterface
    resp.Diagnostics.Append(Interface.ElementAs(ctx, &interface)..)
    if resp.Diagnostics.HasError(){
        return
    }
    state.Interface.Name = types.StringPointerValue(config.Groups.Interface.Name)
    state.Interface.Description = types.StringPointerValue(config.Groups.Interface.Description)
    state.Interface.Mtu = types.Int64PointerValue(config.Groups.Interface.Mtu)
	var vlantagging *bool
    if config.Groups.Interface.Vlan_tagging != nil {
        b, err := strconv.ParseBool(*config.Groups.Interface.Vlan_tagging)
        if err == nil {
            vlantagging = &b
        }
    }
    state.Interface.Vlan_tagging = types.BoolPointerValue(vlantagging)
    var unit = make([]InterfaceModel, len(config.Groups.Interface.Unit
    
    for i, xmlUnit := range config.Groups.Interface.Unit {
        unit[i] = UnitModel {
            Name:      types.StringPointerValue(xmlUnit.Name),
            Description:      types.StringPointerValue(xmlUnit.Description),
            Vlan_id:      types.Int32PointerValue(xmlUnit.Vlan_id),
        }
    }

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
