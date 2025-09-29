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


// Junos XML Hierarchy

type xml_Configuration struct {
	XMLName xml.Name `xml:"configuration"`
	Groups struct {
		XMLName xml.Name `xml:"groups"`
		Name    *string   `xml:"name"`
		Chassis []xml_Chassis `xml:"chassis,omitempty"`
		Interfaces []xml_Interfaces `xml:"interfaces,omitempty"`
		Policy_options []xml_Policy_options `xml:"policy-options,omitempty"`
		Protocols []xml_Protocols `xml:"protocols,omitempty"`
		Routing_instances []xml_Routing_instances `xml:"routing-instances,omitempty"`
		Routing_options []xml_Routing_options `xml:"routing-options,omitempty"`
		Security []xml_Security `xml:"security,omitempty"`
		Snmp []xml_Snmp `xml:"snmp,omitempty"`
		System []xml_System `xml:"system,omitempty"`
	}
}
type xml_Chassis struct {
	XMLName xml.Name `xml:"chassis"`
	Aggregated_devices []xml_Chassis_Aggregated_devices `xml:"aggregated-devices,omitempty"`
}

type xml_Interfaces struct {
	XMLName xml.Name `xml:"interfaces"`
	Interface []xml_Interfaces_Interface `xml:"interface,omitempty"`
}

type xml_Policy_options struct {
	XMLName xml.Name `xml:"policy-options"`
	Policy_statement []xml_Policy_options_Policy_statement `xml:"policy-statement,omitempty"`
}

type xml_Protocols struct {
	XMLName xml.Name `xml:"protocols"`
	Lldp []xml_Protocols_Lldp `xml:"lldp,omitempty"`
}

type xml_Routing_instances struct {
	XMLName xml.Name `xml:"routing-instances"`
	Instance []xml_Routing_instances_Instance `xml:"instance,omitempty"`
}

type xml_Routing_options struct {
	XMLName xml.Name `xml:"routing-options"`
	Static []xml_Routing_options_Static `xml:"static,omitempty"`
}

type xml_Security struct {
	XMLName xml.Name `xml:"security"`
	Log []xml_Security_Log `xml:"log,omitempty"`
	Screen []xml_Security_Screen `xml:"screen,omitempty"`
	Policies []xml_Security_Policies `xml:"policies,omitempty"`
	Zones []xml_Security_Zones `xml:"zones,omitempty"`
}

type xml_Snmp struct {
	XMLName xml.Name `xml:"snmp"`
	Location         *string  `xml:"location,omitempty"`
	Contact         *string  `xml:"contact,omitempty"`
	Community []xml_Snmp_Community `xml:"community,omitempty"`
}

type xml_System struct {
	XMLName xml.Name `xml:"system"`
	Login []xml_System_Login `xml:"login,omitempty"`
	Root_authentication []xml_System_Root_authentication `xml:"root-authentication,omitempty"`
	Host_name         *string  `xml:"host-name,omitempty"`
	Services []xml_System_Services `xml:"services,omitempty"`
	Syslog []xml_System_Syslog `xml:"syslog,omitempty"`
	License []xml_System_License `xml:"license,omitempty"`
}


type xml_Chassis_Aggregated_devices struct {
	XMLName xml.Name `xml:"aggregated-devices"`
	Ethernet []xml_Chassis_Aggregated_devices_Ethernet `xml:"ethernet,omitempty"`
}
type xml_Interfaces_Interface struct {
	XMLName xml.Name `xml:"interface"`
	Name         *string  `xml:"name,omitempty"`
	Vlan_tagging         *string  `xml:"vlan-tagging,omitempty"`
	Unit []xml_Interfaces_Interface_Unit `xml:"unit,omitempty"`
}
type xml_Policy_options_Policy_statement struct {
	XMLName xml.Name `xml:"policy-statement"`
	Name         *string  `xml:"name,omitempty"`
	Term []xml_Policy_options_Policy_statement_Term `xml:"term,omitempty"`
}
type xml_Protocols_Lldp struct {
	XMLName xml.Name `xml:"lldp"`
	Interface []xml_Protocols_Lldp_Interface `xml:"interface,omitempty"`
}
type xml_Routing_instances_Instance struct {
	XMLName xml.Name `xml:"instance"`
	Name         *string  `xml:"name,omitempty"`
	Instance_type         *string  `xml:"instance-type,omitempty"`
	Interface []xml_Routing_instances_Instance_Interface `xml:"interface,omitempty"`
	Routing_options []xml_Routing_instances_Instance_Routing_options `xml:"routing-options,omitempty"`
	Protocols []xml_Routing_instances_Instance_Protocols `xml:"protocols,omitempty"`
}
type xml_Routing_options_Static struct {
	XMLName xml.Name `xml:"static"`
	Route []xml_Routing_options_Static_Route `xml:"route,omitempty"`
}
type xml_Security_Log struct {
	XMLName xml.Name `xml:"log"`
	Mode         *string  `xml:"mode,omitempty"`
	Report []xml_Security_Log_Report `xml:"report,omitempty"`
}
type xml_Security_Screen struct {
	XMLName xml.Name `xml:"screen"`
	Ids_option []xml_Security_Screen_Ids_option `xml:"ids-option,omitempty"`
}
type xml_Security_Policies struct {
	XMLName xml.Name `xml:"policies"`
	Policy []xml_Security_Policies_Policy `xml:"policy,omitempty"`
}
type xml_Security_Zones struct {
	XMLName xml.Name `xml:"zones"`
	Security_zone []xml_Security_Zones_Security_zone `xml:"security-zone,omitempty"`
}
type xml_Snmp_Community struct {
	XMLName xml.Name `xml:"community"`
	Name         *string  `xml:"name,omitempty"`
	Authorization         *string  `xml:"authorization,omitempty"`
}
type xml_System_Login struct {
	XMLName xml.Name `xml:"login"`
	User []xml_System_Login_User `xml:"user,omitempty"`
	Message         *string  `xml:"message,omitempty"`
}
type xml_System_Root_authentication struct {
	XMLName xml.Name `xml:"root-authentication"`
	Encrypted_password         *string  `xml:"encrypted-password,omitempty"`
}
type xml_System_Services struct {
	XMLName xml.Name `xml:"services"`
	Ssh []xml_System_Services_Ssh `xml:"ssh,omitempty"`
	Extension_service []xml_System_Services_Extension_service `xml:"extension-service,omitempty"`
	Netconf []xml_System_Services_Netconf `xml:"netconf,omitempty"`
	Rest []xml_System_Services_Rest `xml:"rest,omitempty"`
}
type xml_System_Syslog struct {
	XMLName xml.Name `xml:"syslog"`
	User []xml_System_Syslog_User `xml:"user,omitempty"`
	File []xml_System_Syslog_File `xml:"file,omitempty"`
}
type xml_System_License struct {
	XMLName xml.Name `xml:"license"`
	Autoupdate []xml_System_License_Autoupdate `xml:"autoupdate,omitempty"`
}

type xml_Chassis_Aggregated_devices_Ethernet struct {
	XMLName xml.Name `xml:"ethernet"`
	Device_count         *string  `xml:"device-count,omitempty"`
}
type xml_Interfaces_Interface_Unit struct {
	XMLName xml.Name `xml:"unit"`
	Name         *string  `xml:"name,omitempty"`
	Description         *string  `xml:"description,omitempty"`
	Vlan_id         *string  `xml:"vlan-id,omitempty"`
	Family []xml_Interfaces_Interface_Unit_Family `xml:"family,omitempty"`
}
type xml_Policy_options_Policy_statement_Term struct {
	XMLName xml.Name `xml:"term"`
	Name         *string  `xml:"name,omitempty"`
	From []xml_Policy_options_Policy_statement_Term_From `xml:"from,omitempty"`
}
type xml_Protocols_Lldp_Interface struct {
	XMLName xml.Name `xml:"interface"`
	Name         *string  `xml:"name,omitempty"`
}
type xml_Routing_instances_Instance_Interface struct {
	XMLName xml.Name `xml:"interface"`
	Name         *string  `xml:"name,omitempty"`
}
type xml_Routing_instances_Instance_Routing_options struct {
	XMLName xml.Name `xml:"routing-options"`
	Static []xml_Routing_instances_Instance_Routing_options_Static `xml:"static,omitempty"`
}
type xml_Routing_instances_Instance_Protocols struct {
	XMLName xml.Name `xml:"protocols"`
	Ospf []xml_Routing_instances_Instance_Protocols_Ospf `xml:"ospf,omitempty"`
}
type xml_Routing_options_Static_Route struct {
	XMLName xml.Name `xml:"route"`
	Name         *string  `xml:"name,omitempty"`
	Next_hop         []*string  `xml:"next-hop,omitempty"`
}
type xml_Security_Log_Report struct {
	XMLName xml.Name `xml:"report"`
}
type xml_Security_Screen_Ids_option struct {
	XMLName xml.Name `xml:"ids-option"`
	Name         *string  `xml:"name,omitempty"`
	Icmp []xml_Security_Screen_Ids_option_Icmp `xml:"icmp,omitempty"`
	Ip []xml_Security_Screen_Ids_option_Ip `xml:"ip,omitempty"`
	Tcp []xml_Security_Screen_Ids_option_Tcp `xml:"tcp,omitempty"`
}
type xml_Security_Policies_Policy struct {
	XMLName xml.Name `xml:"policy"`
	From_zone_name         *string  `xml:"from-zone-name,omitempty"`
	To_zone_name         *string  `xml:"to-zone-name,omitempty"`
	Policy []xml_Security_Policies_Policy_Policy `xml:"policy,omitempty"`
}
type xml_Security_Zones_Security_zone struct {
	XMLName xml.Name `xml:"security-zone"`
	Name         *string  `xml:"name,omitempty"`
	Tcp_rst         *string  `xml:"tcp-rst,omitempty"`
	Screen         *string  `xml:"screen,omitempty"`
	Host_inbound_traffic []xml_Security_Zones_Security_zone_Host_inbound_traffic `xml:"host-inbound-traffic,omitempty"`
	Interfaces []xml_Security_Zones_Security_zone_Interfaces `xml:"interfaces,omitempty"`
}
type xml_System_Login_User struct {
	XMLName xml.Name `xml:"user"`
	Name         *string  `xml:"name,omitempty"`
	Uid         *string  `xml:"uid,omitempty"`
	Class         *string  `xml:"class,omitempty"`
	Authentication []xml_System_Login_User_Authentication `xml:"authentication,omitempty"`
}
type xml_System_Services_Ssh struct {
	XMLName xml.Name `xml:"ssh"`
	Root_login         *string  `xml:"root-login,omitempty"`
}
type xml_System_Services_Extension_service struct {
	XMLName xml.Name `xml:"extension-service"`
	Request_response []xml_System_Services_Extension_service_Request_response `xml:"request-response,omitempty"`
	Notification []xml_System_Services_Extension_service_Notification `xml:"notification,omitempty"`
}
type xml_System_Services_Netconf struct {
	XMLName xml.Name `xml:"netconf"`
	Ssh []xml_System_Services_Netconf_Ssh `xml:"ssh,omitempty"`
}
type xml_System_Services_Rest struct {
	XMLName xml.Name `xml:"rest"`
	Http []xml_System_Services_Rest_Http `xml:"http,omitempty"`
	Enable_explorer         *string  `xml:"enable-explorer,omitempty"`
}
type xml_System_Syslog_User struct {
	XMLName xml.Name `xml:"user"`
	Name         *string  `xml:"name,omitempty"`
	Contents []xml_System_Syslog_User_Contents `xml:"contents,omitempty"`
}
type xml_System_Syslog_File struct {
	XMLName xml.Name `xml:"file"`
	Name         *string  `xml:"name,omitempty"`
	Contents []xml_System_Syslog_File_Contents `xml:"contents,omitempty"`
}
type xml_System_License_Autoupdate struct {
	XMLName xml.Name `xml:"autoupdate"`
	Url []xml_System_License_Autoupdate_Url `xml:"url,omitempty"`
}

type xml_Interfaces_Interface_Unit_Family struct {
	XMLName xml.Name `xml:"family"`
	Inet []xml_Interfaces_Interface_Unit_Family_Inet `xml:"inet,omitempty"`
}
type xml_Policy_options_Policy_statement_Term_From struct {
	XMLName xml.Name `xml:"from"`
	Route_filter []xml_Policy_options_Policy_statement_Term_From_Route_filter `xml:"route-filter,omitempty"`
}
type xml_Routing_instances_Instance_Routing_options_Static struct {
	XMLName xml.Name `xml:"static"`
	Route []xml_Routing_instances_Instance_Routing_options_Static_Route `xml:"route,omitempty"`
}
type xml_Routing_instances_Instance_Protocols_Ospf struct {
	XMLName xml.Name `xml:"ospf"`
	Export         []*string  `xml:"export,omitempty"`
	Area []xml_Routing_instances_Instance_Protocols_Ospf_Area `xml:"area,omitempty"`
}
type xml_Security_Screen_Ids_option_Icmp struct {
	XMLName xml.Name `xml:"icmp"`
	Ping_death         *string  `xml:"ping-death,omitempty"`
}
type xml_Security_Screen_Ids_option_Ip struct {
	XMLName xml.Name `xml:"ip"`
	Source_route_option         *string  `xml:"source-route-option,omitempty"`
	Tear_drop         *string  `xml:"tear-drop,omitempty"`
}
type xml_Security_Screen_Ids_option_Tcp struct {
	XMLName xml.Name `xml:"tcp"`
	Syn_flood []xml_Security_Screen_Ids_option_Tcp_Syn_flood `xml:"syn-flood,omitempty"`
	Land         *string  `xml:"land,omitempty"`
}
type xml_Security_Policies_Policy_Policy struct {
	XMLName xml.Name `xml:"policy"`
	Name         *string  `xml:"name,omitempty"`
	Match []xml_Security_Policies_Policy_Policy_Match `xml:"match,omitempty"`
	Then []xml_Security_Policies_Policy_Policy_Then `xml:"then,omitempty"`
}
type xml_Security_Zones_Security_zone_Host_inbound_traffic struct {
	XMLName xml.Name `xml:"host-inbound-traffic"`
	System_services []xml_Security_Zones_Security_zone_Host_inbound_traffic_System_services `xml:"system-services,omitempty"`
	Protocols []xml_Security_Zones_Security_zone_Host_inbound_traffic_Protocols `xml:"protocols,omitempty"`
}
type xml_Security_Zones_Security_zone_Interfaces struct {
	XMLName xml.Name `xml:"interfaces"`
	Name         *string  `xml:"name,omitempty"`
}
type xml_System_Login_User_Authentication struct {
	XMLName xml.Name `xml:"authentication"`
	Encrypted_password         *string  `xml:"encrypted-password,omitempty"`
}
type xml_System_Services_Extension_service_Request_response struct {
	XMLName xml.Name `xml:"request-response"`
	Grpc []xml_System_Services_Extension_service_Request_response_Grpc `xml:"grpc,omitempty"`
}
type xml_System_Services_Extension_service_Notification struct {
	XMLName xml.Name `xml:"notification"`
	Allow_clients []xml_System_Services_Extension_service_Notification_Allow_clients `xml:"allow-clients,omitempty"`
}
type xml_System_Services_Netconf_Ssh struct {
	XMLName xml.Name `xml:"ssh"`
}
type xml_System_Services_Rest_Http struct {
	XMLName xml.Name `xml:"http"`
	Port         *string  `xml:"port,omitempty"`
}
type xml_System_Syslog_User_Contents struct {
	XMLName xml.Name `xml:"contents"`
	Name         *string  `xml:"name,omitempty"`
	Emergency         *string  `xml:"emergency,omitempty"`
}
type xml_System_Syslog_File_Contents struct {
	XMLName xml.Name `xml:"contents"`
	Name         *string  `xml:"name,omitempty"`
	Any         *string  `xml:"any,omitempty"`
	Info         *string  `xml:"info,omitempty"`
}
type xml_System_License_Autoupdate_Url struct {
	XMLName xml.Name `xml:"url"`
	Name         *string  `xml:"name,omitempty"`
}

type xml_Interfaces_Interface_Unit_Family_Inet struct {
	XMLName xml.Name `xml:"inet"`
	Address []xml_Interfaces_Interface_Unit_Family_Inet_Address `xml:"address,omitempty"`
}
type xml_Policy_options_Policy_statement_Term_From_Route_filter struct {
	XMLName xml.Name `xml:"route-filter"`
	Address         *string  `xml:"address,omitempty"`
	Exact         *string  `xml:"exact,omitempty"`
	Accept         *string  `xml:"accept,omitempty"`
}
type xml_Routing_instances_Instance_Routing_options_Static_Route struct {
	XMLName xml.Name `xml:"route"`
	Name         *string  `xml:"name,omitempty"`
	Discard         *string  `xml:"discard,omitempty"`
}
type xml_Routing_instances_Instance_Protocols_Ospf_Area struct {
	XMLName xml.Name `xml:"area"`
	Name         *string  `xml:"name,omitempty"`
	Interface []xml_Routing_instances_Instance_Protocols_Ospf_Area_Interface `xml:"interface,omitempty"`
}
type xml_Security_Screen_Ids_option_Tcp_Syn_flood struct {
	XMLName xml.Name `xml:"syn-flood"`
	Alarm_threshold         *string  `xml:"alarm-threshold,omitempty"`
	Attack_threshold         *string  `xml:"attack-threshold,omitempty"`
	Source_threshold         *string  `xml:"source-threshold,omitempty"`
	Destination_threshold         *string  `xml:"destination-threshold,omitempty"`
	Timeout         *string  `xml:"timeout,omitempty"`
}
type xml_Security_Policies_Policy_Policy_Match struct {
	XMLName xml.Name `xml:"match"`
	Source_address         []*string  `xml:"source-address,omitempty"`
	Destination_address         []*string  `xml:"destination-address,omitempty"`
	Application         []*string  `xml:"application,omitempty"`
}
type xml_Security_Policies_Policy_Policy_Then struct {
	XMLName xml.Name `xml:"then"`
	Permit []xml_Security_Policies_Policy_Policy_Then_Permit `xml:"permit,omitempty"`
}
type xml_Security_Zones_Security_zone_Host_inbound_traffic_System_services struct {
	XMLName xml.Name `xml:"system-services"`
	Name         *string  `xml:"name,omitempty"`
}
type xml_Security_Zones_Security_zone_Host_inbound_traffic_Protocols struct {
	XMLName xml.Name `xml:"protocols"`
	Name         *string  `xml:"name,omitempty"`
}
type xml_System_Services_Extension_service_Request_response_Grpc struct {
	XMLName xml.Name `xml:"grpc"`
	Max_connections         *string  `xml:"max-connections,omitempty"`
}
type xml_System_Services_Extension_service_Notification_Allow_clients struct {
	XMLName xml.Name `xml:"allow-clients"`
	Address         []*string  `xml:"address,omitempty"`
}

type xml_Interfaces_Interface_Unit_Family_Inet_Address struct {
	XMLName xml.Name `xml:"address"`
	Name         *string  `xml:"name,omitempty"`
}
type xml_Routing_instances_Instance_Protocols_Ospf_Area_Interface struct {
	XMLName xml.Name `xml:"interface"`
	Name         *string  `xml:"name,omitempty"`
}
type xml_Security_Policies_Policy_Policy_Then_Permit struct {
	XMLName xml.Name `xml:"permit"`
}




// Collecting objects from the .tf file
type Groups_Model struct {
	ResourceName	types.String `tfsdk:"resource_name"`
	Chassis types.List `tfsdk:"chassis"`
	Interfaces types.List `tfsdk:"interfaces"`
	Policy_options types.List `tfsdk:"policy_options"`
	Protocols types.List `tfsdk:"protocols"`
	Routing_instances types.List `tfsdk:"routing_instances"`
	Routing_options types.List `tfsdk:"routing_options"`
	Security types.List `tfsdk:"security"`
	Snmp types.List `tfsdk:"snmp"`
	System types.List `tfsdk:"system"`
}
func (o Groups_Model) AttrTypes() map[string]attr.Type {
	return map[string]attr.Type {
		"chassis": 	types.ListType{ElemType: types.ObjectType{AttrTypes: Chassis_Model{}.AttrTypes()}},
		"interfaces": 	types.ListType{ElemType: types.ObjectType{AttrTypes: Interfaces_Model{}.AttrTypes()}},
		"policy_options": 	types.ListType{ElemType: types.ObjectType{AttrTypes: Policy_options_Model{}.AttrTypes()}},
		"protocols": 	types.ListType{ElemType: types.ObjectType{AttrTypes: Protocols_Model{}.AttrTypes()}},
		"routing_instances": 	types.ListType{ElemType: types.ObjectType{AttrTypes: Routing_instances_Model{}.AttrTypes()}},
		"routing_options": 	types.ListType{ElemType: types.ObjectType{AttrTypes: Routing_options_Model{}.AttrTypes()}},
		"security": 	types.ListType{ElemType: types.ObjectType{AttrTypes: Security_Model{}.AttrTypes()}},
		"snmp": 	types.ListType{ElemType: types.ObjectType{AttrTypes: Snmp_Model{}.AttrTypes()}},
		"system": 	types.ListType{ElemType: types.ObjectType{AttrTypes: System_Model{}.AttrTypes()}},
	}
}
func (o Groups_Model) Attributes() map[string]schema.Attribute {
	return map[string]schema.Attribute{
		"resource_name": schema.StringAttribute {
			Required: true,
			MarkdownDescription: "xpath is `config.Groups.resource_name`",
		},
		"chassis": schema.ListNestedAttribute{
			Optional: true,
			NestedObject: schema.NestedAttributeObject{
				Attributes: Chassis_Model{}.Attributes(),
			},
		},
		"interfaces": schema.ListNestedAttribute{
			Optional: true,
			NestedObject: schema.NestedAttributeObject{
				Attributes: Interfaces_Model{}.Attributes(),
			},
		},
		"policy_options": schema.ListNestedAttribute{
			Optional: true,
			NestedObject: schema.NestedAttributeObject{
				Attributes: Policy_options_Model{}.Attributes(),
			},
		},
		"protocols": schema.ListNestedAttribute{
			Optional: true,
			NestedObject: schema.NestedAttributeObject{
				Attributes: Protocols_Model{}.Attributes(),
			},
		},
		"routing_instances": schema.ListNestedAttribute{
			Optional: true,
			NestedObject: schema.NestedAttributeObject{
				Attributes: Routing_instances_Model{}.Attributes(),
			},
		},
		"routing_options": schema.ListNestedAttribute{
			Optional: true,
			NestedObject: schema.NestedAttributeObject{
				Attributes: Routing_options_Model{}.Attributes(),
			},
		},
		"security": schema.ListNestedAttribute{
			Optional: true,
			NestedObject: schema.NestedAttributeObject{
				Attributes: Security_Model{}.Attributes(),
			},
		},
		"snmp": schema.ListNestedAttribute{
			Optional: true,
			NestedObject: schema.NestedAttributeObject{
				Attributes: Snmp_Model{}.Attributes(),
			},
		},
		"system": schema.ListNestedAttribute{
			Optional: true,
			NestedObject: schema.NestedAttributeObject{
				Attributes: System_Model{}.Attributes(),
			},
		},
	}
}
type Chassis_Model struct {
	Aggregated_devices	types.List `tfsdk:"aggregated_devices"`
}
func (o Chassis_Model) AttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"aggregated_devices": 	types.ListType{ElemType: types.ObjectType{AttrTypes: Chassis_Aggregated_devices_Model{}.AttrTypes()}},
	}
}
func (o Chassis_Model) Attributes() map[string]schema.Attribute {
	return map[string]schema.Attribute{
		"aggregated_devices": schema.ListNestedAttribute{
			Optional: true,
			NestedObject: schema.NestedAttributeObject{
				Attributes: Chassis_Aggregated_devices_Model{}.Attributes(),
			},
		},
	}
}
type Interfaces_Model struct {
	Interface	types.List `tfsdk:"interface"`
}
func (o Interfaces_Model) AttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"interface": 	types.ListType{ElemType: types.ObjectType{AttrTypes: Interfaces_Interface_Model{}.AttrTypes()}},
	}
}
func (o Interfaces_Model) Attributes() map[string]schema.Attribute {
	return map[string]schema.Attribute{
		"interface": schema.ListNestedAttribute{
			Optional: true,
			NestedObject: schema.NestedAttributeObject{
				Attributes: Interfaces_Interface_Model{}.Attributes(),
			},
		},
	}
}
type Policy_options_Model struct {
	Policy_statement	types.List `tfsdk:"policy_statement"`
}
func (o Policy_options_Model) AttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"policy_statement": 	types.ListType{ElemType: types.ObjectType{AttrTypes: Policy_options_Policy_statement_Model{}.AttrTypes()}},
	}
}
func (o Policy_options_Model) Attributes() map[string]schema.Attribute {
	return map[string]schema.Attribute{
		"policy_statement": schema.ListNestedAttribute{
			Optional: true,
			NestedObject: schema.NestedAttributeObject{
				Attributes: Policy_options_Policy_statement_Model{}.Attributes(),
			},
		},
	}
}
type Protocols_Model struct {
	Lldp	types.List `tfsdk:"lldp"`
}
func (o Protocols_Model) AttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"lldp": 	types.ListType{ElemType: types.ObjectType{AttrTypes: Protocols_Lldp_Model{}.AttrTypes()}},
	}
}
func (o Protocols_Model) Attributes() map[string]schema.Attribute {
	return map[string]schema.Attribute{
		"lldp": schema.ListNestedAttribute{
			Optional: true,
			NestedObject: schema.NestedAttributeObject{
				Attributes: Protocols_Lldp_Model{}.Attributes(),
			},
		},
	}
}
type Routing_instances_Model struct {
	Instance	types.List `tfsdk:"instance"`
}
func (o Routing_instances_Model) AttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"instance": 	types.ListType{ElemType: types.ObjectType{AttrTypes: Routing_instances_Instance_Model{}.AttrTypes()}},
	}
}
func (o Routing_instances_Model) Attributes() map[string]schema.Attribute {
	return map[string]schema.Attribute{
		"instance": schema.ListNestedAttribute{
			Optional: true,
			NestedObject: schema.NestedAttributeObject{
				Attributes: Routing_instances_Instance_Model{}.Attributes(),
			},
		},
	}
}
type Routing_options_Model struct {
	Static	types.List `tfsdk:"static"`
}
func (o Routing_options_Model) AttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"static": 	types.ListType{ElemType: types.ObjectType{AttrTypes: Routing_options_Static_Model{}.AttrTypes()}},
	}
}
func (o Routing_options_Model) Attributes() map[string]schema.Attribute {
	return map[string]schema.Attribute{
		"static": schema.ListNestedAttribute{
			Optional: true,
			NestedObject: schema.NestedAttributeObject{
				Attributes: Routing_options_Static_Model{}.Attributes(),
			},
		},
	}
}
type Security_Model struct {
	Log	types.List `tfsdk:"log"`
	Screen	types.List `tfsdk:"screen"`
	Policies	types.List `tfsdk:"policies"`
	Zones	types.List `tfsdk:"zones"`
}
func (o Security_Model) AttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"log": 	types.ListType{ElemType: types.ObjectType{AttrTypes: Security_Log_Model{}.AttrTypes()}},
		"screen": 	types.ListType{ElemType: types.ObjectType{AttrTypes: Security_Screen_Model{}.AttrTypes()}},
		"policies": 	types.ListType{ElemType: types.ObjectType{AttrTypes: Security_Policies_Model{}.AttrTypes()}},
		"zones": 	types.ListType{ElemType: types.ObjectType{AttrTypes: Security_Zones_Model{}.AttrTypes()}},
	}
}
func (o Security_Model) Attributes() map[string]schema.Attribute {
	return map[string]schema.Attribute{
		"log": schema.ListNestedAttribute{
			Optional: true,
			NestedObject: schema.NestedAttributeObject{
				Attributes: Security_Log_Model{}.Attributes(),
			},
		},
		"screen": schema.ListNestedAttribute{
			Optional: true,
			NestedObject: schema.NestedAttributeObject{
				Attributes: Security_Screen_Model{}.Attributes(),
			},
		},
		"policies": schema.ListNestedAttribute{
			Optional: true,
			NestedObject: schema.NestedAttributeObject{
				Attributes: Security_Policies_Model{}.Attributes(),
			},
		},
		"zones": schema.ListNestedAttribute{
			Optional: true,
			NestedObject: schema.NestedAttributeObject{
				Attributes: Security_Zones_Model{}.Attributes(),
			},
		},
	}
}
type Snmp_Model struct {
	Location	types.String `tfsdk:"location"`
	Contact	types.String `tfsdk:"contact"`
	Community	types.List `tfsdk:"community"`
}
func (o Snmp_Model) AttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"location": 	types.StringType,
		"contact": 	types.StringType,
		"community": 	types.ListType{ElemType: types.ObjectType{AttrTypes: Snmp_Community_Model{}.AttrTypes()}},
	}
}
func (o Snmp_Model) Attributes() map[string]schema.Attribute {
	return map[string]schema.Attribute{
		"location": schema.StringAttribute{
			Optional: true,
			MarkdownDescription: "xpath is `config.Groups.Snmp.Location`",
		},
		"contact": schema.StringAttribute{
			Optional: true,
			MarkdownDescription: "xpath is `config.Groups.Snmp.Contact`",
		},
		"community": schema.ListNestedAttribute{
			Optional: true,
			NestedObject: schema.NestedAttributeObject{
				Attributes: Snmp_Community_Model{}.Attributes(),
			},
		},
	}
}
type System_Model struct {
	Login	types.List `tfsdk:"login"`
	Root_authentication	types.List `tfsdk:"root_authentication"`
	Host_name	types.String `tfsdk:"host_name"`
	Services	types.List `tfsdk:"services"`
	Syslog	types.List `tfsdk:"syslog"`
	License	types.List `tfsdk:"license"`
}
func (o System_Model) AttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"login": 	types.ListType{ElemType: types.ObjectType{AttrTypes: System_Login_Model{}.AttrTypes()}},
		"root_authentication": 	types.ListType{ElemType: types.ObjectType{AttrTypes: System_Root_authentication_Model{}.AttrTypes()}},
		"host_name": 	types.StringType,
		"services": 	types.ListType{ElemType: types.ObjectType{AttrTypes: System_Services_Model{}.AttrTypes()}},
		"syslog": 	types.ListType{ElemType: types.ObjectType{AttrTypes: System_Syslog_Model{}.AttrTypes()}},
		"license": 	types.ListType{ElemType: types.ObjectType{AttrTypes: System_License_Model{}.AttrTypes()}},
	}
}
func (o System_Model) Attributes() map[string]schema.Attribute {
	return map[string]schema.Attribute{
		"login": schema.ListNestedAttribute{
			Optional: true,
			NestedObject: schema.NestedAttributeObject{
				Attributes: System_Login_Model{}.Attributes(),
			},
		},
		"root_authentication": schema.ListNestedAttribute{
			Optional: true,
			NestedObject: schema.NestedAttributeObject{
				Attributes: System_Root_authentication_Model{}.Attributes(),
			},
		},
		"host_name": schema.StringAttribute{
			Optional: true,
			MarkdownDescription: "xpath is `config.Groups.System.Host-name`",
		},
		"services": schema.ListNestedAttribute{
			Optional: true,
			NestedObject: schema.NestedAttributeObject{
				Attributes: System_Services_Model{}.Attributes(),
			},
		},
		"syslog": schema.ListNestedAttribute{
			Optional: true,
			NestedObject: schema.NestedAttributeObject{
				Attributes: System_Syslog_Model{}.Attributes(),
			},
		},
		"license": schema.ListNestedAttribute{
			Optional: true,
			NestedObject: schema.NestedAttributeObject{
				Attributes: System_License_Model{}.Attributes(),
			},
		},
	}
}

type Chassis_Aggregated_devices_Model struct {
	Ethernet	types.List `tfsdk:"ethernet"`
}
func (o Chassis_Aggregated_devices_Model) AttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
	    "ethernet": 	types.ListType{ElemType: types.ObjectType{AttrTypes: Chassis_Aggregated_devices_Ethernet_Model{}.AttrTypes()}},
	}
}
func (o Chassis_Aggregated_devices_Model) Attributes() map[string]schema.Attribute {
	return map[string]schema.Attribute{
	    "ethernet": schema.ListNestedAttribute{
		    Optional: true,
		    NestedObject: schema.NestedAttributeObject{
			    Attributes: Chassis_Aggregated_devices_Ethernet_Model{}.Attributes(),
	        },
        },
    }
}
type Interfaces_Interface_Model struct {
	Name	types.String `tfsdk:"name"`
	Vlan_tagging	types.String `tfsdk:"vlan_tagging"`
	Unit	types.List `tfsdk:"unit"`
}
func (o Interfaces_Interface_Model) AttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
	    "name": 	types.StringType,
	    "vlan_tagging": 	types.StringType,
	    "unit": 	types.ListType{ElemType: types.ObjectType{AttrTypes: Interfaces_Interface_Unit_Model{}.AttrTypes()}},
	}
}
func (o Interfaces_Interface_Model) Attributes() map[string]schema.Attribute {
	return map[string]schema.Attribute{
	    "name": schema.StringAttribute{
		    Optional: true,
		    MarkdownDescription: "xpath is `config.Groups.Name.Interface`",
	    },
	    "vlan_tagging": schema.StringAttribute{
		    Optional: true,
		    MarkdownDescription: "xpath is `config.Groups.Vlan-tagging.Interface`",
	    },
	    "unit": schema.ListNestedAttribute{
		    Optional: true,
		    NestedObject: schema.NestedAttributeObject{
			    Attributes: Interfaces_Interface_Unit_Model{}.Attributes(),
	        },
        },
    }
}
type Policy_options_Policy_statement_Model struct {
	Name	types.String `tfsdk:"name"`
	Term	types.List `tfsdk:"term"`
}
func (o Policy_options_Policy_statement_Model) AttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
	    "name": 	types.StringType,
	    "term": 	types.ListType{ElemType: types.ObjectType{AttrTypes: Policy_options_Policy_statement_Term_Model{}.AttrTypes()}},
	}
}
func (o Policy_options_Policy_statement_Model) Attributes() map[string]schema.Attribute {
	return map[string]schema.Attribute{
	    "name": schema.StringAttribute{
		    Optional: true,
		    MarkdownDescription: "xpath is `config.Groups.Name.Policy_statement`",
	    },
	    "term": schema.ListNestedAttribute{
		    Optional: true,
		    NestedObject: schema.NestedAttributeObject{
			    Attributes: Policy_options_Policy_statement_Term_Model{}.Attributes(),
	        },
        },
    }
}
type Protocols_Lldp_Model struct {
	Interface	types.List `tfsdk:"interface"`
}
func (o Protocols_Lldp_Model) AttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
	    "interface": 	types.ListType{ElemType: types.ObjectType{AttrTypes: Protocols_Lldp_Interface_Model{}.AttrTypes()}},
	}
}
func (o Protocols_Lldp_Model) Attributes() map[string]schema.Attribute {
	return map[string]schema.Attribute{
	    "interface": schema.ListNestedAttribute{
		    Optional: true,
		    NestedObject: schema.NestedAttributeObject{
			    Attributes: Protocols_Lldp_Interface_Model{}.Attributes(),
	        },
        },
    }
}
type Routing_instances_Instance_Model struct {
	Name	types.String `tfsdk:"name"`
	Instance_type	types.String `tfsdk:"instance_type"`
	Interface	types.List `tfsdk:"interface"`
	Routing_options	types.List `tfsdk:"routing_options"`
	Protocols	types.List `tfsdk:"protocols"`
}
func (o Routing_instances_Instance_Model) AttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
	    "name": 	types.StringType,
	    "instance_type": 	types.StringType,
	    "interface": 	types.ListType{ElemType: types.ObjectType{AttrTypes: Routing_instances_Instance_Interface_Model{}.AttrTypes()}},
	    "routing_options": 	types.ListType{ElemType: types.ObjectType{AttrTypes: Routing_instances_Instance_Routing_options_Model{}.AttrTypes()}},
	    "protocols": 	types.ListType{ElemType: types.ObjectType{AttrTypes: Routing_instances_Instance_Protocols_Model{}.AttrTypes()}},
	}
}
func (o Routing_instances_Instance_Model) Attributes() map[string]schema.Attribute {
	return map[string]schema.Attribute{
	    "name": schema.StringAttribute{
		    Optional: true,
		    MarkdownDescription: "xpath is `config.Groups.Name.Instance`",
	    },
	    "instance_type": schema.StringAttribute{
		    Optional: true,
		    MarkdownDescription: "xpath is `config.Groups.Instance-type.Instance`",
	    },
	    "interface": schema.ListNestedAttribute{
		    Optional: true,
		    NestedObject: schema.NestedAttributeObject{
			    Attributes: Routing_instances_Instance_Interface_Model{}.Attributes(),
	        },
        },
	    "routing_options": schema.ListNestedAttribute{
		    Optional: true,
		    NestedObject: schema.NestedAttributeObject{
			    Attributes: Routing_instances_Instance_Routing_options_Model{}.Attributes(),
	        },
        },
	    "protocols": schema.ListNestedAttribute{
		    Optional: true,
		    NestedObject: schema.NestedAttributeObject{
			    Attributes: Routing_instances_Instance_Protocols_Model{}.Attributes(),
	        },
        },
    }
}
type Routing_options_Static_Model struct {
	Route	types.List `tfsdk:"route"`
}
func (o Routing_options_Static_Model) AttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
	    "route": 	types.ListType{ElemType: types.ObjectType{AttrTypes: Routing_options_Static_Route_Model{}.AttrTypes()}},
	}
}
func (o Routing_options_Static_Model) Attributes() map[string]schema.Attribute {
	return map[string]schema.Attribute{
	    "route": schema.ListNestedAttribute{
		    Optional: true,
		    NestedObject: schema.NestedAttributeObject{
			    Attributes: Routing_options_Static_Route_Model{}.Attributes(),
	        },
        },
    }
}
type Security_Log_Model struct {
	Mode	types.String `tfsdk:"mode"`
	Report	types.List `tfsdk:"report"`
}
func (o Security_Log_Model) AttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
	    "mode": 	types.StringType,
	    "report": 	types.ListType{ElemType: types.ObjectType{AttrTypes: Security_Log_Report_Model{}.AttrTypes()}},
	}
}
func (o Security_Log_Model) Attributes() map[string]schema.Attribute {
	return map[string]schema.Attribute{
	    "mode": schema.StringAttribute{
		    Optional: true,
		    MarkdownDescription: "xpath is `config.Groups.Mode.Log`",
	    },
	    "report": schema.ListNestedAttribute{
		    Optional: true,
		    NestedObject: schema.NestedAttributeObject{
			    Attributes: Security_Log_Report_Model{}.Attributes(),
	        },
        },
    }
}
type Security_Screen_Model struct {
	Ids_option	types.List `tfsdk:"ids_option"`
}
func (o Security_Screen_Model) AttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
	    "ids_option": 	types.ListType{ElemType: types.ObjectType{AttrTypes: Security_Screen_Ids_option_Model{}.AttrTypes()}},
	}
}
func (o Security_Screen_Model) Attributes() map[string]schema.Attribute {
	return map[string]schema.Attribute{
	    "ids_option": schema.ListNestedAttribute{
		    Optional: true,
		    NestedObject: schema.NestedAttributeObject{
			    Attributes: Security_Screen_Ids_option_Model{}.Attributes(),
	        },
        },
    }
}
type Security_Policies_Model struct {
	Policy	types.List `tfsdk:"policy"`
}
func (o Security_Policies_Model) AttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
	    "policy": 	types.ListType{ElemType: types.ObjectType{AttrTypes: Security_Policies_Policy_Model{}.AttrTypes()}},
	}
}
func (o Security_Policies_Model) Attributes() map[string]schema.Attribute {
	return map[string]schema.Attribute{
	    "policy": schema.ListNestedAttribute{
		    Optional: true,
		    NestedObject: schema.NestedAttributeObject{
			    Attributes: Security_Policies_Policy_Model{}.Attributes(),
	        },
        },
    }
}
type Security_Zones_Model struct {
	Security_zone	types.List `tfsdk:"security_zone"`
}
func (o Security_Zones_Model) AttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
	    "security_zone": 	types.ListType{ElemType: types.ObjectType{AttrTypes: Security_Zones_Security_zone_Model{}.AttrTypes()}},
	}
}
func (o Security_Zones_Model) Attributes() map[string]schema.Attribute {
	return map[string]schema.Attribute{
	    "security_zone": schema.ListNestedAttribute{
		    Optional: true,
		    NestedObject: schema.NestedAttributeObject{
			    Attributes: Security_Zones_Security_zone_Model{}.Attributes(),
	        },
        },
    }
}
type Snmp_Community_Model struct {
	Name	types.String `tfsdk:"name"`
	Authorization	types.String `tfsdk:"authorization"`
}
func (o Snmp_Community_Model) AttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
	    "name": 	types.StringType,
	    "authorization": 	types.StringType,
	}
}
func (o Snmp_Community_Model) Attributes() map[string]schema.Attribute {
	return map[string]schema.Attribute{
	    "name": schema.StringAttribute{
		    Optional: true,
		    MarkdownDescription: "xpath is `config.Groups.Name.Community`",
	    },
	    "authorization": schema.StringAttribute{
		    Optional: true,
		    MarkdownDescription: "xpath is `config.Groups.Authorization.Community`",
	    },
    }
}
type System_Login_Model struct {
	User	types.List `tfsdk:"user"`
	Message	types.String `tfsdk:"message"`
}
func (o System_Login_Model) AttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
	    "user": 	types.ListType{ElemType: types.ObjectType{AttrTypes: System_Login_User_Model{}.AttrTypes()}},
	    "message": 	types.StringType,
	}
}
func (o System_Login_Model) Attributes() map[string]schema.Attribute {
	return map[string]schema.Attribute{
	    "user": schema.ListNestedAttribute{
		    Optional: true,
		    NestedObject: schema.NestedAttributeObject{
			    Attributes: System_Login_User_Model{}.Attributes(),
	        },
        },
	    "message": schema.StringAttribute{
		    Optional: true,
		    MarkdownDescription: "xpath is `config.Groups.Message.Login`",
	    },
    }
}
type System_Root_authentication_Model struct {
	Encrypted_password	types.String `tfsdk:"encrypted_password"`
}
func (o System_Root_authentication_Model) AttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
	    "encrypted_password": 	types.StringType,
	}
}
func (o System_Root_authentication_Model) Attributes() map[string]schema.Attribute {
	return map[string]schema.Attribute{
	    "encrypted_password": schema.StringAttribute{
		    Optional: true,
		    MarkdownDescription: "xpath is `config.Groups.Encrypted-password.Root_authentication`",
	    },
    }
}
type System_Services_Model struct {
	Ssh	types.List `tfsdk:"ssh"`
	Extension_service	types.List `tfsdk:"extension_service"`
	Netconf	types.List `tfsdk:"netconf"`
	Rest	types.List `tfsdk:"rest"`
}
func (o System_Services_Model) AttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
	    "ssh": 	types.ListType{ElemType: types.ObjectType{AttrTypes: System_Services_Ssh_Model{}.AttrTypes()}},
	    "extension_service": 	types.ListType{ElemType: types.ObjectType{AttrTypes: System_Services_Extension_service_Model{}.AttrTypes()}},
	    "netconf": 	types.ListType{ElemType: types.ObjectType{AttrTypes: System_Services_Netconf_Model{}.AttrTypes()}},
	    "rest": 	types.ListType{ElemType: types.ObjectType{AttrTypes: System_Services_Rest_Model{}.AttrTypes()}},
	}
}
func (o System_Services_Model) Attributes() map[string]schema.Attribute {
	return map[string]schema.Attribute{
	    "ssh": schema.ListNestedAttribute{
		    Optional: true,
		    NestedObject: schema.NestedAttributeObject{
			    Attributes: System_Services_Ssh_Model{}.Attributes(),
	        },
        },
	    "extension_service": schema.ListNestedAttribute{
		    Optional: true,
		    NestedObject: schema.NestedAttributeObject{
			    Attributes: System_Services_Extension_service_Model{}.Attributes(),
	        },
        },
	    "netconf": schema.ListNestedAttribute{
		    Optional: true,
		    NestedObject: schema.NestedAttributeObject{
			    Attributes: System_Services_Netconf_Model{}.Attributes(),
	        },
        },
	    "rest": schema.ListNestedAttribute{
		    Optional: true,
		    NestedObject: schema.NestedAttributeObject{
			    Attributes: System_Services_Rest_Model{}.Attributes(),
	        },
        },
    }
}
type System_Syslog_Model struct {
	User	types.List `tfsdk:"user"`
	File	types.List `tfsdk:"file"`
}
func (o System_Syslog_Model) AttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
	    "user": 	types.ListType{ElemType: types.ObjectType{AttrTypes: System_Syslog_User_Model{}.AttrTypes()}},
	    "file": 	types.ListType{ElemType: types.ObjectType{AttrTypes: System_Syslog_File_Model{}.AttrTypes()}},
	}
}
func (o System_Syslog_Model) Attributes() map[string]schema.Attribute {
	return map[string]schema.Attribute{
	    "user": schema.ListNestedAttribute{
		    Optional: true,
		    NestedObject: schema.NestedAttributeObject{
			    Attributes: System_Syslog_User_Model{}.Attributes(),
	        },
        },
	    "file": schema.ListNestedAttribute{
		    Optional: true,
		    NestedObject: schema.NestedAttributeObject{
			    Attributes: System_Syslog_File_Model{}.Attributes(),
	        },
        },
    }
}
type System_License_Model struct {
	Autoupdate	types.List `tfsdk:"autoupdate"`
}
func (o System_License_Model) AttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
	    "autoupdate": 	types.ListType{ElemType: types.ObjectType{AttrTypes: System_License_Autoupdate_Model{}.AttrTypes()}},
	}
}
func (o System_License_Model) Attributes() map[string]schema.Attribute {
	return map[string]schema.Attribute{
	    "autoupdate": schema.ListNestedAttribute{
		    Optional: true,
		    NestedObject: schema.NestedAttributeObject{
			    Attributes: System_License_Autoupdate_Model{}.Attributes(),
	        },
        },
    }
}

type Chassis_Aggregated_devices_Ethernet_Model struct {
	Device_count	types.String `tfsdk:"device_count"`
}
func (o Chassis_Aggregated_devices_Ethernet_Model) AttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
	    "device_count": 	types.StringType,
	}
}
func (o Chassis_Aggregated_devices_Ethernet_Model) Attributes() map[string]schema.Attribute {
	return map[string]schema.Attribute{
	    "device_count": schema.StringAttribute{
		    Optional: true,
		    MarkdownDescription: "xpath is `config.Groups.Device-count.Ethernet`",
	    },
    }
}
type Interfaces_Interface_Unit_Model struct {
	Name	types.String `tfsdk:"name"`
	Description	types.String `tfsdk:"description"`
	Vlan_id	types.String `tfsdk:"vlan_id"`
	Family	types.List `tfsdk:"family"`
}
func (o Interfaces_Interface_Unit_Model) AttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
	    "name": 	types.StringType,
	    "description": 	types.StringType,
	    "vlan_id": 	types.StringType,
	    "family": 	types.ListType{ElemType: types.ObjectType{AttrTypes: Interfaces_Interface_Unit_Family_Model{}.AttrTypes()}},
	}
}
func (o Interfaces_Interface_Unit_Model) Attributes() map[string]schema.Attribute {
	return map[string]schema.Attribute{
	    "name": schema.StringAttribute{
		    Optional: true,
		    MarkdownDescription: "xpath is `config.Groups.Name.Unit`",
	    },
	    "description": schema.StringAttribute{
		    Optional: true,
		    MarkdownDescription: "xpath is `config.Groups.Description.Unit`",
	    },
	    "vlan_id": schema.StringAttribute{
		    Optional: true,
		    MarkdownDescription: "xpath is `config.Groups.Vlan-id.Unit`",
	    },
	    "family": schema.ListNestedAttribute{
		    Optional: true,
		    NestedObject: schema.NestedAttributeObject{
			    Attributes: Interfaces_Interface_Unit_Family_Model{}.Attributes(),
	        },
        },
    }
}
type Policy_options_Policy_statement_Term_Model struct {
	Name	types.String `tfsdk:"name"`
	From	types.List `tfsdk:"from"`
}
func (o Policy_options_Policy_statement_Term_Model) AttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
	    "name": 	types.StringType,
	    "from": 	types.ListType{ElemType: types.ObjectType{AttrTypes: Policy_options_Policy_statement_Term_From_Model{}.AttrTypes()}},
	}
}
func (o Policy_options_Policy_statement_Term_Model) Attributes() map[string]schema.Attribute {
	return map[string]schema.Attribute{
	    "name": schema.StringAttribute{
		    Optional: true,
		    MarkdownDescription: "xpath is `config.Groups.Name.Term`",
	    },
	    "from": schema.ListNestedAttribute{
		    Optional: true,
		    NestedObject: schema.NestedAttributeObject{
			    Attributes: Policy_options_Policy_statement_Term_From_Model{}.Attributes(),
	        },
        },
    }
}
type Protocols_Lldp_Interface_Model struct {
	Name	types.String `tfsdk:"name"`
}
func (o Protocols_Lldp_Interface_Model) AttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
	    "name": 	types.StringType,
	}
}
func (o Protocols_Lldp_Interface_Model) Attributes() map[string]schema.Attribute {
	return map[string]schema.Attribute{
	    "name": schema.StringAttribute{
		    Optional: true,
		    MarkdownDescription: "xpath is `config.Groups.Name.Interface`",
	    },
    }
}
type Routing_instances_Instance_Interface_Model struct {
	Name	types.String `tfsdk:"name"`
}
func (o Routing_instances_Instance_Interface_Model) AttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
	    "name": 	types.StringType,
	}
}
func (o Routing_instances_Instance_Interface_Model) Attributes() map[string]schema.Attribute {
	return map[string]schema.Attribute{
	    "name": schema.StringAttribute{
		    Optional: true,
		    MarkdownDescription: "xpath is `config.Groups.Name.Interface`",
	    },
    }
}
type Routing_instances_Instance_Routing_options_Model struct {
	Static	types.List `tfsdk:"static"`
}
func (o Routing_instances_Instance_Routing_options_Model) AttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
	    "static": 	types.ListType{ElemType: types.ObjectType{AttrTypes: Routing_instances_Instance_Routing_options_Static_Model{}.AttrTypes()}},
	}
}
func (o Routing_instances_Instance_Routing_options_Model) Attributes() map[string]schema.Attribute {
	return map[string]schema.Attribute{
	    "static": schema.ListNestedAttribute{
		    Optional: true,
		    NestedObject: schema.NestedAttributeObject{
			    Attributes: Routing_instances_Instance_Routing_options_Static_Model{}.Attributes(),
	        },
        },
    }
}
type Routing_instances_Instance_Protocols_Model struct {
	Ospf	types.List `tfsdk:"ospf"`
}
func (o Routing_instances_Instance_Protocols_Model) AttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
	    "ospf": 	types.ListType{ElemType: types.ObjectType{AttrTypes: Routing_instances_Instance_Protocols_Ospf_Model{}.AttrTypes()}},
	}
}
func (o Routing_instances_Instance_Protocols_Model) Attributes() map[string]schema.Attribute {
	return map[string]schema.Attribute{
	    "ospf": schema.ListNestedAttribute{
		    Optional: true,
		    NestedObject: schema.NestedAttributeObject{
			    Attributes: Routing_instances_Instance_Protocols_Ospf_Model{}.Attributes(),
	        },
        },
    }
}
type Routing_options_Static_Route_Model struct {
	Name	types.String `tfsdk:"name"`
	Next_hop	types.List `tfsdk:"next_hop"`
}
func (o Routing_options_Static_Route_Model) AttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
	    "name": 	types.StringType,
		"next_hop": 	types.ListType{ElemType: types.StringType},
	}
}
func (o Routing_options_Static_Route_Model) Attributes() map[string]schema.Attribute {
	return map[string]schema.Attribute{
	    "name": schema.StringAttribute{
		    Optional: true,
		    MarkdownDescription: "xpath is `config.Groups.Name.Route`",
	    },
		"next_hop": schema.ListAttribute{
			ElementType:         types.StringType,
			Optional:            true,
			MarkdownDescription: "xpath is `config.Groups.Next-hop.Route`",
		},
    }
}
type Security_Log_Report_Model struct {
}
func (o Security_Log_Report_Model) AttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
	}
}
func (o Security_Log_Report_Model) Attributes() map[string]schema.Attribute {
	return map[string]schema.Attribute{
    }
}
type Security_Screen_Ids_option_Model struct {
	Name	types.String `tfsdk:"name"`
	Icmp	types.List `tfsdk:"icmp"`
	Ip	types.List `tfsdk:"ip"`
	Tcp	types.List `tfsdk:"tcp"`
}
func (o Security_Screen_Ids_option_Model) AttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
	    "name": 	types.StringType,
	    "icmp": 	types.ListType{ElemType: types.ObjectType{AttrTypes: Security_Screen_Ids_option_Icmp_Model{}.AttrTypes()}},
	    "ip": 	types.ListType{ElemType: types.ObjectType{AttrTypes: Security_Screen_Ids_option_Ip_Model{}.AttrTypes()}},
	    "tcp": 	types.ListType{ElemType: types.ObjectType{AttrTypes: Security_Screen_Ids_option_Tcp_Model{}.AttrTypes()}},
	}
}
func (o Security_Screen_Ids_option_Model) Attributes() map[string]schema.Attribute {
	return map[string]schema.Attribute{
	    "name": schema.StringAttribute{
		    Optional: true,
		    MarkdownDescription: "xpath is `config.Groups.Name.Ids_option`",
	    },
	    "icmp": schema.ListNestedAttribute{
		    Optional: true,
		    NestedObject: schema.NestedAttributeObject{
			    Attributes: Security_Screen_Ids_option_Icmp_Model{}.Attributes(),
	        },
        },
	    "ip": schema.ListNestedAttribute{
		    Optional: true,
		    NestedObject: schema.NestedAttributeObject{
			    Attributes: Security_Screen_Ids_option_Ip_Model{}.Attributes(),
	        },
        },
	    "tcp": schema.ListNestedAttribute{
		    Optional: true,
		    NestedObject: schema.NestedAttributeObject{
			    Attributes: Security_Screen_Ids_option_Tcp_Model{}.Attributes(),
	        },
        },
    }
}
type Security_Policies_Policy_Model struct {
	From_zone_name	types.String `tfsdk:"from_zone_name"`
	To_zone_name	types.String `tfsdk:"to_zone_name"`
	Policy	types.List `tfsdk:"policy"`
}
func (o Security_Policies_Policy_Model) AttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
	    "from_zone_name": 	types.StringType,
	    "to_zone_name": 	types.StringType,
	    "policy": 	types.ListType{ElemType: types.ObjectType{AttrTypes: Security_Policies_Policy_Policy_Model{}.AttrTypes()}},
	}
}
func (o Security_Policies_Policy_Model) Attributes() map[string]schema.Attribute {
	return map[string]schema.Attribute{
	    "from_zone_name": schema.StringAttribute{
		    Optional: true,
		    MarkdownDescription: "xpath is `config.Groups.From-zone-name.Policy`",
	    },
	    "to_zone_name": schema.StringAttribute{
		    Optional: true,
		    MarkdownDescription: "xpath is `config.Groups.To-zone-name.Policy`",
	    },
	    "policy": schema.ListNestedAttribute{
		    Optional: true,
		    NestedObject: schema.NestedAttributeObject{
			    Attributes: Security_Policies_Policy_Policy_Model{}.Attributes(),
	        },
        },
    }
}
type Security_Zones_Security_zone_Model struct {
	Name	types.String `tfsdk:"name"`
	Tcp_rst	types.String `tfsdk:"tcp_rst"`
	Screen	types.String `tfsdk:"screen"`
	Host_inbound_traffic	types.List `tfsdk:"host_inbound_traffic"`
	Interfaces	types.List `tfsdk:"interfaces"`
}
func (o Security_Zones_Security_zone_Model) AttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
	    "name": 	types.StringType,
	    "tcp_rst": 	types.StringType,
	    "screen": 	types.StringType,
	    "host_inbound_traffic": 	types.ListType{ElemType: types.ObjectType{AttrTypes: Security_Zones_Security_zone_Host_inbound_traffic_Model{}.AttrTypes()}},
	    "interfaces": 	types.ListType{ElemType: types.ObjectType{AttrTypes: Security_Zones_Security_zone_Interfaces_Model{}.AttrTypes()}},
	}
}
func (o Security_Zones_Security_zone_Model) Attributes() map[string]schema.Attribute {
	return map[string]schema.Attribute{
	    "name": schema.StringAttribute{
		    Optional: true,
		    MarkdownDescription: "xpath is `config.Groups.Name.Security_zone`",
	    },
	    "tcp_rst": schema.StringAttribute{
		    Optional: true,
		    MarkdownDescription: "xpath is `config.Groups.Tcp-rst.Security_zone`",
	    },
	    "screen": schema.StringAttribute{
		    Optional: true,
		    MarkdownDescription: "xpath is `config.Groups.Screen.Security_zone`",
	    },
	    "host_inbound_traffic": schema.ListNestedAttribute{
		    Optional: true,
		    NestedObject: schema.NestedAttributeObject{
			    Attributes: Security_Zones_Security_zone_Host_inbound_traffic_Model{}.Attributes(),
	        },
        },
	    "interfaces": schema.ListNestedAttribute{
		    Optional: true,
		    NestedObject: schema.NestedAttributeObject{
			    Attributes: Security_Zones_Security_zone_Interfaces_Model{}.Attributes(),
	        },
        },
    }
}
type System_Login_User_Model struct {
	Name	types.String `tfsdk:"name"`
	Uid	types.String `tfsdk:"uid"`
	Class	types.String `tfsdk:"class"`
	Authentication	types.List `tfsdk:"authentication"`
}
func (o System_Login_User_Model) AttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
	    "name": 	types.StringType,
	    "uid": 	types.StringType,
	    "class": 	types.StringType,
	    "authentication": 	types.ListType{ElemType: types.ObjectType{AttrTypes: System_Login_User_Authentication_Model{}.AttrTypes()}},
	}
}
func (o System_Login_User_Model) Attributes() map[string]schema.Attribute {
	return map[string]schema.Attribute{
	    "name": schema.StringAttribute{
		    Optional: true,
		    MarkdownDescription: "xpath is `config.Groups.Name.User`",
	    },
	    "uid": schema.StringAttribute{
		    Optional: true,
		    MarkdownDescription: "xpath is `config.Groups.Uid.User`",
	    },
	    "class": schema.StringAttribute{
		    Optional: true,
		    MarkdownDescription: "xpath is `config.Groups.Class.User`",
	    },
	    "authentication": schema.ListNestedAttribute{
		    Optional: true,
		    NestedObject: schema.NestedAttributeObject{
			    Attributes: System_Login_User_Authentication_Model{}.Attributes(),
	        },
        },
    }
}
type System_Services_Ssh_Model struct {
	Root_login	types.String `tfsdk:"root_login"`
}
func (o System_Services_Ssh_Model) AttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
	    "root_login": 	types.StringType,
	}
}
func (o System_Services_Ssh_Model) Attributes() map[string]schema.Attribute {
	return map[string]schema.Attribute{
	    "root_login": schema.StringAttribute{
		    Optional: true,
		    MarkdownDescription: "xpath is `config.Groups.Root-login.Ssh`",
	    },
    }
}
type System_Services_Extension_service_Model struct {
	Request_response	types.List `tfsdk:"request_response"`
	Notification	types.List `tfsdk:"notification"`
}
func (o System_Services_Extension_service_Model) AttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
	    "request_response": 	types.ListType{ElemType: types.ObjectType{AttrTypes: System_Services_Extension_service_Request_response_Model{}.AttrTypes()}},
	    "notification": 	types.ListType{ElemType: types.ObjectType{AttrTypes: System_Services_Extension_service_Notification_Model{}.AttrTypes()}},
	}
}
func (o System_Services_Extension_service_Model) Attributes() map[string]schema.Attribute {
	return map[string]schema.Attribute{
	    "request_response": schema.ListNestedAttribute{
		    Optional: true,
		    NestedObject: schema.NestedAttributeObject{
			    Attributes: System_Services_Extension_service_Request_response_Model{}.Attributes(),
	        },
        },
	    "notification": schema.ListNestedAttribute{
		    Optional: true,
		    NestedObject: schema.NestedAttributeObject{
			    Attributes: System_Services_Extension_service_Notification_Model{}.Attributes(),
	        },
        },
    }
}
type System_Services_Netconf_Model struct {
	Ssh	types.List `tfsdk:"ssh"`
}
func (o System_Services_Netconf_Model) AttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
	    "ssh": 	types.ListType{ElemType: types.ObjectType{AttrTypes: System_Services_Netconf_Ssh_Model{}.AttrTypes()}},
	}
}
func (o System_Services_Netconf_Model) Attributes() map[string]schema.Attribute {
	return map[string]schema.Attribute{
	    "ssh": schema.ListNestedAttribute{
		    Optional: true,
		    NestedObject: schema.NestedAttributeObject{
			    Attributes: System_Services_Netconf_Ssh_Model{}.Attributes(),
	        },
        },
    }
}
type System_Services_Rest_Model struct {
	Http	types.List `tfsdk:"http"`
	Enable_explorer	types.String `tfsdk:"enable_explorer"`
}
func (o System_Services_Rest_Model) AttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
	    "http": 	types.ListType{ElemType: types.ObjectType{AttrTypes: System_Services_Rest_Http_Model{}.AttrTypes()}},
	    "enable_explorer": 	types.StringType,
	}
}
func (o System_Services_Rest_Model) Attributes() map[string]schema.Attribute {
	return map[string]schema.Attribute{
	    "http": schema.ListNestedAttribute{
		    Optional: true,
		    NestedObject: schema.NestedAttributeObject{
			    Attributes: System_Services_Rest_Http_Model{}.Attributes(),
	        },
        },
	    "enable_explorer": schema.StringAttribute{
		    Optional: true,
		    MarkdownDescription: "xpath is `config.Groups.Enable-explorer.Rest`",
	    },
    }
}
type System_Syslog_User_Model struct {
	Name	types.String `tfsdk:"name"`
	Contents	types.List `tfsdk:"contents"`
}
func (o System_Syslog_User_Model) AttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
	    "name": 	types.StringType,
	    "contents": 	types.ListType{ElemType: types.ObjectType{AttrTypes: System_Syslog_User_Contents_Model{}.AttrTypes()}},
	}
}
func (o System_Syslog_User_Model) Attributes() map[string]schema.Attribute {
	return map[string]schema.Attribute{
	    "name": schema.StringAttribute{
		    Optional: true,
		    MarkdownDescription: "xpath is `config.Groups.Name.User`",
	    },
	    "contents": schema.ListNestedAttribute{
		    Optional: true,
		    NestedObject: schema.NestedAttributeObject{
			    Attributes: System_Syslog_User_Contents_Model{}.Attributes(),
	        },
        },
    }
}
type System_Syslog_File_Model struct {
	Name	types.String `tfsdk:"name"`
	Contents	types.List `tfsdk:"contents"`
}
func (o System_Syslog_File_Model) AttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
	    "name": 	types.StringType,
	    "contents": 	types.ListType{ElemType: types.ObjectType{AttrTypes: System_Syslog_File_Contents_Model{}.AttrTypes()}},
	}
}
func (o System_Syslog_File_Model) Attributes() map[string]schema.Attribute {
	return map[string]schema.Attribute{
	    "name": schema.StringAttribute{
		    Optional: true,
		    MarkdownDescription: "xpath is `config.Groups.Name.File`",
	    },
	    "contents": schema.ListNestedAttribute{
		    Optional: true,
		    NestedObject: schema.NestedAttributeObject{
			    Attributes: System_Syslog_File_Contents_Model{}.Attributes(),
	        },
        },
    }
}
type System_License_Autoupdate_Model struct {
	Url	types.List `tfsdk:"url"`
}
func (o System_License_Autoupdate_Model) AttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
	    "url": 	types.ListType{ElemType: types.ObjectType{AttrTypes: System_License_Autoupdate_Url_Model{}.AttrTypes()}},
	}
}
func (o System_License_Autoupdate_Model) Attributes() map[string]schema.Attribute {
	return map[string]schema.Attribute{
	    "url": schema.ListNestedAttribute{
		    Optional: true,
		    NestedObject: schema.NestedAttributeObject{
			    Attributes: System_License_Autoupdate_Url_Model{}.Attributes(),
	        },
        },
    }
}

type Interfaces_Interface_Unit_Family_Model struct {
	Inet	types.List `tfsdk:"inet"`
}
func (o Interfaces_Interface_Unit_Family_Model) AttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
	    "inet": 	types.ListType{ElemType: types.ObjectType{AttrTypes: Interfaces_Interface_Unit_Family_Inet_Model{}.AttrTypes()}},
	}
}
func (o Interfaces_Interface_Unit_Family_Model) Attributes() map[string]schema.Attribute {
	return map[string]schema.Attribute{
	    "inet": schema.ListNestedAttribute{
		    Optional: true,
		    NestedObject: schema.NestedAttributeObject{
			    Attributes: Interfaces_Interface_Unit_Family_Inet_Model{}.Attributes(),
	        },
        },
    }
}
type Policy_options_Policy_statement_Term_From_Model struct {
	Route_filter	types.List `tfsdk:"route_filter"`
}
func (o Policy_options_Policy_statement_Term_From_Model) AttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
	    "route_filter": 	types.ListType{ElemType: types.ObjectType{AttrTypes: Policy_options_Policy_statement_Term_From_Route_filter_Model{}.AttrTypes()}},
	}
}
func (o Policy_options_Policy_statement_Term_From_Model) Attributes() map[string]schema.Attribute {
	return map[string]schema.Attribute{
	    "route_filter": schema.ListNestedAttribute{
		    Optional: true,
		    NestedObject: schema.NestedAttributeObject{
			    Attributes: Policy_options_Policy_statement_Term_From_Route_filter_Model{}.Attributes(),
	        },
        },
    }
}
type Routing_instances_Instance_Routing_options_Static_Model struct {
	Route	types.List `tfsdk:"route"`
}
func (o Routing_instances_Instance_Routing_options_Static_Model) AttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
	    "route": 	types.ListType{ElemType: types.ObjectType{AttrTypes: Routing_instances_Instance_Routing_options_Static_Route_Model{}.AttrTypes()}},
	}
}
func (o Routing_instances_Instance_Routing_options_Static_Model) Attributes() map[string]schema.Attribute {
	return map[string]schema.Attribute{
	    "route": schema.ListNestedAttribute{
		    Optional: true,
		    NestedObject: schema.NestedAttributeObject{
			    Attributes: Routing_instances_Instance_Routing_options_Static_Route_Model{}.Attributes(),
	        },
        },
    }
}
type Routing_instances_Instance_Protocols_Ospf_Model struct {
	Export	types.List `tfsdk:"export"`
	Area	types.List `tfsdk:"area"`
}
func (o Routing_instances_Instance_Protocols_Ospf_Model) AttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"export": 	types.ListType{ElemType: types.StringType},
	    "area": 	types.ListType{ElemType: types.ObjectType{AttrTypes: Routing_instances_Instance_Protocols_Ospf_Area_Model{}.AttrTypes()}},
	}
}
func (o Routing_instances_Instance_Protocols_Ospf_Model) Attributes() map[string]schema.Attribute {
	return map[string]schema.Attribute{
		"export": schema.ListAttribute{
			ElementType:         types.StringType,
			Optional:            true,
			MarkdownDescription: "xpath is `config.Groups.Export.Ospf`",
		},
	    "area": schema.ListNestedAttribute{
		    Optional: true,
		    NestedObject: schema.NestedAttributeObject{
			    Attributes: Routing_instances_Instance_Protocols_Ospf_Area_Model{}.Attributes(),
	        },
        },
    }
}
type Security_Screen_Ids_option_Icmp_Model struct {
	Ping_death	types.String `tfsdk:"ping_death"`
}
func (o Security_Screen_Ids_option_Icmp_Model) AttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
	    "ping_death": 	types.StringType,
	}
}
func (o Security_Screen_Ids_option_Icmp_Model) Attributes() map[string]schema.Attribute {
	return map[string]schema.Attribute{
	    "ping_death": schema.StringAttribute{
		    Optional: true,
		    MarkdownDescription: "xpath is `config.Groups.Ping-death.Icmp`",
	    },
    }
}
type Security_Screen_Ids_option_Ip_Model struct {
	Source_route_option	types.String `tfsdk:"source_route_option"`
	Tear_drop	types.String `tfsdk:"tear_drop"`
}
func (o Security_Screen_Ids_option_Ip_Model) AttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
	    "source_route_option": 	types.StringType,
	    "tear_drop": 	types.StringType,
	}
}
func (o Security_Screen_Ids_option_Ip_Model) Attributes() map[string]schema.Attribute {
	return map[string]schema.Attribute{
	    "source_route_option": schema.StringAttribute{
		    Optional: true,
		    MarkdownDescription: "xpath is `config.Groups.Source-route-option.Ip`",
	    },
	    "tear_drop": schema.StringAttribute{
		    Optional: true,
		    MarkdownDescription: "xpath is `config.Groups.Tear-drop.Ip`",
	    },
    }
}
type Security_Screen_Ids_option_Tcp_Model struct {
	Syn_flood	types.List `tfsdk:"syn_flood"`
	Land	types.String `tfsdk:"land"`
}
func (o Security_Screen_Ids_option_Tcp_Model) AttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
	    "syn_flood": 	types.ListType{ElemType: types.ObjectType{AttrTypes: Security_Screen_Ids_option_Tcp_Syn_flood_Model{}.AttrTypes()}},
	    "land": 	types.StringType,
	}
}
func (o Security_Screen_Ids_option_Tcp_Model) Attributes() map[string]schema.Attribute {
	return map[string]schema.Attribute{
	    "syn_flood": schema.ListNestedAttribute{
		    Optional: true,
		    NestedObject: schema.NestedAttributeObject{
			    Attributes: Security_Screen_Ids_option_Tcp_Syn_flood_Model{}.Attributes(),
	        },
        },
	    "land": schema.StringAttribute{
		    Optional: true,
		    MarkdownDescription: "xpath is `config.Groups.Land.Tcp`",
	    },
    }
}
type Security_Policies_Policy_Policy_Model struct {
	Name	types.String `tfsdk:"name"`
	Match	types.List `tfsdk:"match"`
	Then	types.List `tfsdk:"then"`
}
func (o Security_Policies_Policy_Policy_Model) AttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
	    "name": 	types.StringType,
	    "match": 	types.ListType{ElemType: types.ObjectType{AttrTypes: Security_Policies_Policy_Policy_Match_Model{}.AttrTypes()}},
	    "then": 	types.ListType{ElemType: types.ObjectType{AttrTypes: Security_Policies_Policy_Policy_Then_Model{}.AttrTypes()}},
	}
}
func (o Security_Policies_Policy_Policy_Model) Attributes() map[string]schema.Attribute {
	return map[string]schema.Attribute{
	    "name": schema.StringAttribute{
		    Optional: true,
		    MarkdownDescription: "xpath is `config.Groups.Name.Policy`",
	    },
	    "match": schema.ListNestedAttribute{
		    Optional: true,
		    NestedObject: schema.NestedAttributeObject{
			    Attributes: Security_Policies_Policy_Policy_Match_Model{}.Attributes(),
	        },
        },
	    "then": schema.ListNestedAttribute{
		    Optional: true,
		    NestedObject: schema.NestedAttributeObject{
			    Attributes: Security_Policies_Policy_Policy_Then_Model{}.Attributes(),
	        },
        },
    }
}
type Security_Zones_Security_zone_Host_inbound_traffic_Model struct {
	System_services	types.List `tfsdk:"system_services"`
	Protocols	types.List `tfsdk:"protocols"`
}
func (o Security_Zones_Security_zone_Host_inbound_traffic_Model) AttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
	    "system_services": 	types.ListType{ElemType: types.ObjectType{AttrTypes: Security_Zones_Security_zone_Host_inbound_traffic_System_services_Model{}.AttrTypes()}},
	    "protocols": 	types.ListType{ElemType: types.ObjectType{AttrTypes: Security_Zones_Security_zone_Host_inbound_traffic_Protocols_Model{}.AttrTypes()}},
	}
}
func (o Security_Zones_Security_zone_Host_inbound_traffic_Model) Attributes() map[string]schema.Attribute {
	return map[string]schema.Attribute{
	    "system_services": schema.ListNestedAttribute{
		    Optional: true,
		    NestedObject: schema.NestedAttributeObject{
			    Attributes: Security_Zones_Security_zone_Host_inbound_traffic_System_services_Model{}.Attributes(),
	        },
        },
	    "protocols": schema.ListNestedAttribute{
		    Optional: true,
		    NestedObject: schema.NestedAttributeObject{
			    Attributes: Security_Zones_Security_zone_Host_inbound_traffic_Protocols_Model{}.Attributes(),
	        },
        },
    }
}
type Security_Zones_Security_zone_Interfaces_Model struct {
	Name	types.String `tfsdk:"name"`
}
func (o Security_Zones_Security_zone_Interfaces_Model) AttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
	    "name": 	types.StringType,
	}
}
func (o Security_Zones_Security_zone_Interfaces_Model) Attributes() map[string]schema.Attribute {
	return map[string]schema.Attribute{
	    "name": schema.StringAttribute{
		    Optional: true,
		    MarkdownDescription: "xpath is `config.Groups.Name.Interfaces`",
	    },
    }
}
type System_Login_User_Authentication_Model struct {
	Encrypted_password	types.String `tfsdk:"encrypted_password"`
}
func (o System_Login_User_Authentication_Model) AttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
	    "encrypted_password": 	types.StringType,
	}
}
func (o System_Login_User_Authentication_Model) Attributes() map[string]schema.Attribute {
	return map[string]schema.Attribute{
	    "encrypted_password": schema.StringAttribute{
		    Optional: true,
		    MarkdownDescription: "xpath is `config.Groups.Encrypted-password.Authentication`",
	    },
    }
}
type System_Services_Extension_service_Request_response_Model struct {
	Grpc	types.List `tfsdk:"grpc"`
}
func (o System_Services_Extension_service_Request_response_Model) AttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
	    "grpc": 	types.ListType{ElemType: types.ObjectType{AttrTypes: System_Services_Extension_service_Request_response_Grpc_Model{}.AttrTypes()}},
	}
}
func (o System_Services_Extension_service_Request_response_Model) Attributes() map[string]schema.Attribute {
	return map[string]schema.Attribute{
	    "grpc": schema.ListNestedAttribute{
		    Optional: true,
		    NestedObject: schema.NestedAttributeObject{
			    Attributes: System_Services_Extension_service_Request_response_Grpc_Model{}.Attributes(),
	        },
        },
    }
}
type System_Services_Extension_service_Notification_Model struct {
	Allow_clients	types.List `tfsdk:"allow_clients"`
}
func (o System_Services_Extension_service_Notification_Model) AttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
	    "allow_clients": 	types.ListType{ElemType: types.ObjectType{AttrTypes: System_Services_Extension_service_Notification_Allow_clients_Model{}.AttrTypes()}},
	}
}
func (o System_Services_Extension_service_Notification_Model) Attributes() map[string]schema.Attribute {
	return map[string]schema.Attribute{
	    "allow_clients": schema.ListNestedAttribute{
		    Optional: true,
		    NestedObject: schema.NestedAttributeObject{
			    Attributes: System_Services_Extension_service_Notification_Allow_clients_Model{}.Attributes(),
	        },
        },
    }
}
type System_Services_Netconf_Ssh_Model struct {
}
func (o System_Services_Netconf_Ssh_Model) AttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
	}
}
func (o System_Services_Netconf_Ssh_Model) Attributes() map[string]schema.Attribute {
	return map[string]schema.Attribute{
    }
}
type System_Services_Rest_Http_Model struct {
	Port	types.String `tfsdk:"port"`
}
func (o System_Services_Rest_Http_Model) AttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
	    "port": 	types.StringType,
	}
}
func (o System_Services_Rest_Http_Model) Attributes() map[string]schema.Attribute {
	return map[string]schema.Attribute{
	    "port": schema.StringAttribute{
		    Optional: true,
		    MarkdownDescription: "xpath is `config.Groups.Port.Http`",
	    },
    }
}
type System_Syslog_User_Contents_Model struct {
	Name	types.String `tfsdk:"name"`
	Emergency	types.String `tfsdk:"emergency"`
}
func (o System_Syslog_User_Contents_Model) AttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
	    "name": 	types.StringType,
	    "emergency": 	types.StringType,
	}
}
func (o System_Syslog_User_Contents_Model) Attributes() map[string]schema.Attribute {
	return map[string]schema.Attribute{
	    "name": schema.StringAttribute{
		    Optional: true,
		    MarkdownDescription: "xpath is `config.Groups.Name.Contents`",
	    },
	    "emergency": schema.StringAttribute{
		    Optional: true,
		    MarkdownDescription: "xpath is `config.Groups.Emergency.Contents`",
	    },
    }
}
type System_Syslog_File_Contents_Model struct {
	Name	types.String `tfsdk:"name"`
	Any	types.String `tfsdk:"any"`
	Info	types.String `tfsdk:"info"`
}
func (o System_Syslog_File_Contents_Model) AttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
	    "name": 	types.StringType,
	    "any": 	types.StringType,
	    "info": 	types.StringType,
	}
}
func (o System_Syslog_File_Contents_Model) Attributes() map[string]schema.Attribute {
	return map[string]schema.Attribute{
	    "name": schema.StringAttribute{
		    Optional: true,
		    MarkdownDescription: "xpath is `config.Groups.Name.Contents`",
	    },
	    "any": schema.StringAttribute{
		    Optional: true,
		    MarkdownDescription: "xpath is `config.Groups.Any.Contents`",
	    },
	    "info": schema.StringAttribute{
		    Optional: true,
		    MarkdownDescription: "xpath is `config.Groups.Info.Contents`",
	    },
    }
}
type System_License_Autoupdate_Url_Model struct {
	Name	types.String `tfsdk:"name"`
}
func (o System_License_Autoupdate_Url_Model) AttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
	    "name": 	types.StringType,
	}
}
func (o System_License_Autoupdate_Url_Model) Attributes() map[string]schema.Attribute {
	return map[string]schema.Attribute{
	    "name": schema.StringAttribute{
		    Optional: true,
		    MarkdownDescription: "xpath is `config.Groups.Name.Url`",
	    },
    }
}

type Interfaces_Interface_Unit_Family_Inet_Model struct {
	Address	types.List `tfsdk:"address"`
}
func (o Interfaces_Interface_Unit_Family_Inet_Model) AttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
	    "address": 	types.ListType{ElemType: types.ObjectType{AttrTypes: Interfaces_Interface_Unit_Family_Inet_Address_Model{}.AttrTypes()}},
	}
}
func (o Interfaces_Interface_Unit_Family_Inet_Model) Attributes() map[string]schema.Attribute {
	return map[string]schema.Attribute{
	    "address": schema.ListNestedAttribute{
		    Optional: true,
		    NestedObject: schema.NestedAttributeObject{
			    Attributes: Interfaces_Interface_Unit_Family_Inet_Address_Model{}.Attributes(),
	        },
        },
    }
}
type Policy_options_Policy_statement_Term_From_Route_filter_Model struct {
	Address	types.String `tfsdk:"address"`
	Exact	types.String `tfsdk:"exact"`
	Accept	types.String `tfsdk:"accept"`
}
func (o Policy_options_Policy_statement_Term_From_Route_filter_Model) AttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
	    "address": 	types.StringType,
	    "exact": 	types.StringType,
	    "accept": 	types.StringType,
	}
}
func (o Policy_options_Policy_statement_Term_From_Route_filter_Model) Attributes() map[string]schema.Attribute {
	return map[string]schema.Attribute{
	    "address": schema.StringAttribute{
		    Optional: true,
		    MarkdownDescription: "xpath is `config.Groups.Address.Route_filter`",
	    },
	    "exact": schema.StringAttribute{
		    Optional: true,
		    MarkdownDescription: "xpath is `config.Groups.Exact.Route_filter`",
	    },
	    "accept": schema.StringAttribute{
		    Optional: true,
		    MarkdownDescription: "xpath is `config.Groups.Accept.Route_filter`",
	    },
    }
}
type Routing_instances_Instance_Routing_options_Static_Route_Model struct {
	Name	types.String `tfsdk:"name"`
	Discard	types.String `tfsdk:"discard"`
}
func (o Routing_instances_Instance_Routing_options_Static_Route_Model) AttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
	    "name": 	types.StringType,
	    "discard": 	types.StringType,
	}
}
func (o Routing_instances_Instance_Routing_options_Static_Route_Model) Attributes() map[string]schema.Attribute {
	return map[string]schema.Attribute{
	    "name": schema.StringAttribute{
		    Optional: true,
		    MarkdownDescription: "xpath is `config.Groups.Name.Route`",
	    },
	    "discard": schema.StringAttribute{
		    Optional: true,
		    MarkdownDescription: "xpath is `config.Groups.Discard.Route`",
	    },
    }
}
type Routing_instances_Instance_Protocols_Ospf_Area_Model struct {
	Name	types.String `tfsdk:"name"`
	Interface	types.List `tfsdk:"interface"`
}
func (o Routing_instances_Instance_Protocols_Ospf_Area_Model) AttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
	    "name": 	types.StringType,
	    "interface": 	types.ListType{ElemType: types.ObjectType{AttrTypes: Routing_instances_Instance_Protocols_Ospf_Area_Interface_Model{}.AttrTypes()}},
	}
}
func (o Routing_instances_Instance_Protocols_Ospf_Area_Model) Attributes() map[string]schema.Attribute {
	return map[string]schema.Attribute{
	    "name": schema.StringAttribute{
		    Optional: true,
		    MarkdownDescription: "xpath is `config.Groups.Name.Area`",
	    },
	    "interface": schema.ListNestedAttribute{
		    Optional: true,
		    NestedObject: schema.NestedAttributeObject{
			    Attributes: Routing_instances_Instance_Protocols_Ospf_Area_Interface_Model{}.Attributes(),
	        },
        },
    }
}
type Security_Screen_Ids_option_Tcp_Syn_flood_Model struct {
	Alarm_threshold	types.String `tfsdk:"alarm_threshold"`
	Attack_threshold	types.String `tfsdk:"attack_threshold"`
	Source_threshold	types.String `tfsdk:"source_threshold"`
	Destination_threshold	types.String `tfsdk:"destination_threshold"`
	Timeout	types.String `tfsdk:"timeout"`
}
func (o Security_Screen_Ids_option_Tcp_Syn_flood_Model) AttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
	    "alarm_threshold": 	types.StringType,
	    "attack_threshold": 	types.StringType,
	    "source_threshold": 	types.StringType,
	    "destination_threshold": 	types.StringType,
	    "timeout": 	types.StringType,
	}
}
func (o Security_Screen_Ids_option_Tcp_Syn_flood_Model) Attributes() map[string]schema.Attribute {
	return map[string]schema.Attribute{
	    "alarm_threshold": schema.StringAttribute{
		    Optional: true,
		    MarkdownDescription: "xpath is `config.Groups.Alarm-threshold.Syn_flood`",
	    },
	    "attack_threshold": schema.StringAttribute{
		    Optional: true,
		    MarkdownDescription: "xpath is `config.Groups.Attack-threshold.Syn_flood`",
	    },
	    "source_threshold": schema.StringAttribute{
		    Optional: true,
		    MarkdownDescription: "xpath is `config.Groups.Source-threshold.Syn_flood`",
	    },
	    "destination_threshold": schema.StringAttribute{
		    Optional: true,
		    MarkdownDescription: "xpath is `config.Groups.Destination-threshold.Syn_flood`",
	    },
	    "timeout": schema.StringAttribute{
		    Optional: true,
		    MarkdownDescription: "xpath is `config.Groups.Timeout.Syn_flood`",
	    },
    }
}
type Security_Policies_Policy_Policy_Match_Model struct {
	Source_address	types.List `tfsdk:"source_address"`
	Destination_address	types.List `tfsdk:"destination_address"`
	Application	types.List `tfsdk:"application"`
}
func (o Security_Policies_Policy_Policy_Match_Model) AttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"source_address": 	types.ListType{ElemType: types.StringType},
		"destination_address": 	types.ListType{ElemType: types.StringType},
		"application": 	types.ListType{ElemType: types.StringType},
	}
}
func (o Security_Policies_Policy_Policy_Match_Model) Attributes() map[string]schema.Attribute {
	return map[string]schema.Attribute{
		"source_address": schema.ListAttribute{
			ElementType:         types.StringType,
			Optional:            true,
			MarkdownDescription: "xpath is `config.Groups.Source-address.Match`",
		},
		"destination_address": schema.ListAttribute{
			ElementType:         types.StringType,
			Optional:            true,
			MarkdownDescription: "xpath is `config.Groups.Destination-address.Match`",
		},
		"application": schema.ListAttribute{
			ElementType:         types.StringType,
			Optional:            true,
			MarkdownDescription: "xpath is `config.Groups.Application.Match`",
		},
    }
}
type Security_Policies_Policy_Policy_Then_Model struct {
	Permit	types.List `tfsdk:"permit"`
}
func (o Security_Policies_Policy_Policy_Then_Model) AttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
	    "permit": 	types.ListType{ElemType: types.ObjectType{AttrTypes: Security_Policies_Policy_Policy_Then_Permit_Model{}.AttrTypes()}},
	}
}
func (o Security_Policies_Policy_Policy_Then_Model) Attributes() map[string]schema.Attribute {
	return map[string]schema.Attribute{
	    "permit": schema.ListNestedAttribute{
		    Optional: true,
		    NestedObject: schema.NestedAttributeObject{
			    Attributes: Security_Policies_Policy_Policy_Then_Permit_Model{}.Attributes(),
	        },
        },
    }
}
type Security_Zones_Security_zone_Host_inbound_traffic_System_services_Model struct {
	Name	types.String `tfsdk:"name"`
}
func (o Security_Zones_Security_zone_Host_inbound_traffic_System_services_Model) AttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
	    "name": 	types.StringType,
	}
}
func (o Security_Zones_Security_zone_Host_inbound_traffic_System_services_Model) Attributes() map[string]schema.Attribute {
	return map[string]schema.Attribute{
	    "name": schema.StringAttribute{
		    Optional: true,
		    MarkdownDescription: "xpath is `config.Groups.Name.System_services`",
	    },
    }
}
type Security_Zones_Security_zone_Host_inbound_traffic_Protocols_Model struct {
	Name	types.String `tfsdk:"name"`
}
func (o Security_Zones_Security_zone_Host_inbound_traffic_Protocols_Model) AttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
	    "name": 	types.StringType,
	}
}
func (o Security_Zones_Security_zone_Host_inbound_traffic_Protocols_Model) Attributes() map[string]schema.Attribute {
	return map[string]schema.Attribute{
	    "name": schema.StringAttribute{
		    Optional: true,
		    MarkdownDescription: "xpath is `config.Groups.Name.Protocols`",
	    },
    }
}
type System_Services_Extension_service_Request_response_Grpc_Model struct {
	Max_connections	types.String `tfsdk:"max_connections"`
}
func (o System_Services_Extension_service_Request_response_Grpc_Model) AttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
	    "max_connections": 	types.StringType,
	}
}
func (o System_Services_Extension_service_Request_response_Grpc_Model) Attributes() map[string]schema.Attribute {
	return map[string]schema.Attribute{
	    "max_connections": schema.StringAttribute{
		    Optional: true,
		    MarkdownDescription: "xpath is `config.Groups.Max-connections.Grpc`",
	    },
    }
}
type System_Services_Extension_service_Notification_Allow_clients_Model struct {
	Address	types.List `tfsdk:"address"`
}
func (o System_Services_Extension_service_Notification_Allow_clients_Model) AttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"address": 	types.ListType{ElemType: types.StringType},
	}
}
func (o System_Services_Extension_service_Notification_Allow_clients_Model) Attributes() map[string]schema.Attribute {
	return map[string]schema.Attribute{
		"address": schema.ListAttribute{
			ElementType:         types.StringType,
			Optional:            true,
			MarkdownDescription: "xpath is `config.Groups.Address.Allow_clients`",
		},
    }
}

type Interfaces_Interface_Unit_Family_Inet_Address_Model struct {
	Name	types.String `tfsdk:"name"`
}
func (o Interfaces_Interface_Unit_Family_Inet_Address_Model) AttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
	    "name": 	types.StringType,
	}
}
func (o Interfaces_Interface_Unit_Family_Inet_Address_Model) Attributes() map[string]schema.Attribute {
	return map[string]schema.Attribute{
	    "name": schema.StringAttribute{
		    Optional: true,
		    MarkdownDescription: "xpath is `config.Groups.Name.Address`",
	    },
    }
}
type Routing_instances_Instance_Protocols_Ospf_Area_Interface_Model struct {
	Name	types.String `tfsdk:"name"`
}
func (o Routing_instances_Instance_Protocols_Ospf_Area_Interface_Model) AttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
	    "name": 	types.StringType,
	}
}
func (o Routing_instances_Instance_Protocols_Ospf_Area_Interface_Model) Attributes() map[string]schema.Attribute {
	return map[string]schema.Attribute{
	    "name": schema.StringAttribute{
		    Optional: true,
		    MarkdownDescription: "xpath is `config.Groups.Name.Interface`",
	    },
    }
}
type Security_Policies_Policy_Policy_Then_Permit_Model struct {
}
func (o Security_Policies_Policy_Policy_Then_Permit_Model) AttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
	}
}
func (o Security_Policies_Policy_Policy_Then_Permit_Model) Attributes() map[string]schema.Attribute {
	return map[string]schema.Attribute{
    }
}



// Collects the data for the crud work
type resource_Apply_Groups struct {
	client ProviderConfig
}

func (r *resource_Apply_Groups) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	r.client = req.ProviderData.(ProviderConfig)
}
// Metadata implements resource.Resource.
func (r *resource_Apply_Groups) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = "terraform-provider-" + req.ProviderTypeName
}
// Schema implements resource.Resource.
func (r *resource_Apply_Groups) Schema(_ context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"resource_name": schema.StringAttribute{
				Required:      true,
				PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplace()},
			},
			"chassis": schema.ListNestedAttribute{
				Optional: true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: Chassis_Model{}.Attributes(),
				},
			},
			"interfaces": schema.ListNestedAttribute{
				Optional: true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: Interfaces_Model{}.Attributes(),
				},
			},
			"policy_options": schema.ListNestedAttribute{
				Optional: true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: Policy_options_Model{}.Attributes(),
				},
			},
			"protocols": schema.ListNestedAttribute{
				Optional: true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: Protocols_Model{}.Attributes(),
				},
			},
			"routing_instances": schema.ListNestedAttribute{
				Optional: true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: Routing_instances_Model{}.Attributes(),
				},
			},
			"routing_options": schema.ListNestedAttribute{
				Optional: true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: Routing_options_Model{}.Attributes(),
				},
			},
			"security": schema.ListNestedAttribute{
				Optional: true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: Security_Model{}.Attributes(),
				},
			},
			"snmp": schema.ListNestedAttribute{
				Optional: true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: Snmp_Model{}.Attributes(),
				},
			},
			"system": schema.ListNestedAttribute{
				Optional: true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: System_Model{}.Attributes(),
				},
			},
		},
	}
}




// Create implements resource.Resource.
func (r *resource_Apply_Groups) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	
	var plan Groups_Model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	// Check for errors
	if resp.Diagnostics.HasError() {
		return
	}
	var config xml_Configuration
	config.Groups.Name = plan.ResourceName.ValueStringPointer()
    
	
    var var_chassis []Chassis_Model
    if plan.Chassis.IsNull() {
        var_chassis = []Chassis_Model{}
    }else {
        resp.Diagnostics.Append(plan.Chassis.ElementsAs(ctx, &var_chassis, false)...)
        if resp.Diagnostics.HasError() {
            return
        }
    }
    config.Groups.Chassis = make([]xml_Chassis, len(var_chassis))
   
    for i_chassis, v_chassis := range var_chassis {
        var var_chassis_aggregated_devices []Chassis_Aggregated_devices_Model
        resp.Diagnostics.Append(v_chassis.Aggregated_devices.ElementsAs(ctx, &var_chassis_aggregated_devices, false)...)
        if resp.Diagnostics.HasError() {
            return
        }
	    config.Groups.Chassis[i_chassis].Aggregated_devices = make([]xml_Chassis_Aggregated_devices, len(var_chassis_aggregated_devices))
        
		for i_chassis_aggregated_devices, v_chassis_aggregated_devices := range var_chassis_aggregated_devices {
            var var_chassis_aggregated_devices_ethernet []Chassis_Aggregated_devices_Ethernet_Model
            resp.Diagnostics.Append(v_chassis_aggregated_devices.Ethernet.ElementsAs(ctx, &var_chassis_aggregated_devices_ethernet, false)...)
            if resp.Diagnostics.HasError() {
                return
            }
	    config.Groups.Chassis[i_chassis].Aggregated_devices[i_chassis_aggregated_devices].Ethernet = make([]xml_Chassis_Aggregated_devices_Ethernet, len(var_chassis_aggregated_devices_ethernet))
        
		for i_chassis_aggregated_devices_ethernet, v_chassis_aggregated_devices_ethernet := range var_chassis_aggregated_devices_ethernet {
            config.Groups.Chassis[i_chassis].Aggregated_devices[i_chassis_aggregated_devices].Ethernet[i_chassis_aggregated_devices_ethernet].Device_count = v_chassis_aggregated_devices_ethernet.Device_count.ValueStringPointer()
        }
        }
    }
	
    var var_interfaces []Interfaces_Model
    if plan.Interfaces.IsNull() {
        var_interfaces = []Interfaces_Model{}
    }else {
        resp.Diagnostics.Append(plan.Interfaces.ElementsAs(ctx, &var_interfaces, false)...)
        if resp.Diagnostics.HasError() {
            return
        }
    }
    config.Groups.Interfaces = make([]xml_Interfaces, len(var_interfaces))
   
    for i_interfaces, v_interfaces := range var_interfaces {
        var var_interfaces_interface []Interfaces_Interface_Model
        resp.Diagnostics.Append(v_interfaces.Interface.ElementsAs(ctx, &var_interfaces_interface, false)...)
        if resp.Diagnostics.HasError() {
            return
        }
	    config.Groups.Interfaces[i_interfaces].Interface = make([]xml_Interfaces_Interface, len(var_interfaces_interface))
        
		for i_interfaces_interface, v_interfaces_interface := range var_interfaces_interface {
            config.Groups.Interfaces[i_interfaces].Interface[i_interfaces_interface].Name = v_interfaces_interface.Name.ValueStringPointer()
            config.Groups.Interfaces[i_interfaces].Interface[i_interfaces_interface].Vlan_tagging = v_interfaces_interface.Vlan_tagging.ValueStringPointer()
            var var_interfaces_interface_unit []Interfaces_Interface_Unit_Model
            resp.Diagnostics.Append(v_interfaces_interface.Unit.ElementsAs(ctx, &var_interfaces_interface_unit, false)...)
            if resp.Diagnostics.HasError() {
                return
            }
	    config.Groups.Interfaces[i_interfaces].Interface[i_interfaces_interface].Unit = make([]xml_Interfaces_Interface_Unit, len(var_interfaces_interface_unit))
        
		for i_interfaces_interface_unit, v_interfaces_interface_unit := range var_interfaces_interface_unit {
            config.Groups.Interfaces[i_interfaces].Interface[i_interfaces_interface].Unit[i_interfaces_interface_unit].Name = v_interfaces_interface_unit.Name.ValueStringPointer()
            config.Groups.Interfaces[i_interfaces].Interface[i_interfaces_interface].Unit[i_interfaces_interface_unit].Description = v_interfaces_interface_unit.Description.ValueStringPointer()
            config.Groups.Interfaces[i_interfaces].Interface[i_interfaces_interface].Unit[i_interfaces_interface_unit].Vlan_id = v_interfaces_interface_unit.Vlan_id.ValueStringPointer()
            var var_interfaces_interface_unit_family []Interfaces_Interface_Unit_Family_Model
            resp.Diagnostics.Append(v_interfaces_interface_unit.Family.ElementsAs(ctx, &var_interfaces_interface_unit_family, false)...)
            if resp.Diagnostics.HasError() {
                return
            }
	    config.Groups.Interfaces[i_interfaces].Interface[i_interfaces_interface].Unit[i_interfaces_interface_unit].Family = make([]xml_Interfaces_Interface_Unit_Family, len(var_interfaces_interface_unit_family))
        
		for i_interfaces_interface_unit_family, v_interfaces_interface_unit_family := range var_interfaces_interface_unit_family {
            var var_interfaces_interface_unit_family_inet []Interfaces_Interface_Unit_Family_Inet_Model
            resp.Diagnostics.Append(v_interfaces_interface_unit_family.Inet.ElementsAs(ctx, &var_interfaces_interface_unit_family_inet, false)...)
            if resp.Diagnostics.HasError() {
                return
            }
	    config.Groups.Interfaces[i_interfaces].Interface[i_interfaces_interface].Unit[i_interfaces_interface_unit].Family[i_interfaces_interface_unit_family].Inet = make([]xml_Interfaces_Interface_Unit_Family_Inet, len(var_interfaces_interface_unit_family_inet))
        
		for i_interfaces_interface_unit_family_inet, v_interfaces_interface_unit_family_inet := range var_interfaces_interface_unit_family_inet {
            var var_interfaces_interface_unit_family_inet_address []Interfaces_Interface_Unit_Family_Inet_Address_Model
            resp.Diagnostics.Append(v_interfaces_interface_unit_family_inet.Address.ElementsAs(ctx, &var_interfaces_interface_unit_family_inet_address, false)...)
            if resp.Diagnostics.HasError() {
                return
            }
	    config.Groups.Interfaces[i_interfaces].Interface[i_interfaces_interface].Unit[i_interfaces_interface_unit].Family[i_interfaces_interface_unit_family].Inet[i_interfaces_interface_unit_family_inet].Address = make([]xml_Interfaces_Interface_Unit_Family_Inet_Address, len(var_interfaces_interface_unit_family_inet_address))
        
		for i_interfaces_interface_unit_family_inet_address, v_interfaces_interface_unit_family_inet_address := range var_interfaces_interface_unit_family_inet_address {
            config.Groups.Interfaces[i_interfaces].Interface[i_interfaces_interface].Unit[i_interfaces_interface_unit].Family[i_interfaces_interface_unit_family].Inet[i_interfaces_interface_unit_family_inet].Address[i_interfaces_interface_unit_family_inet_address].Name = v_interfaces_interface_unit_family_inet_address.Name.ValueStringPointer()
        }
        }
        }
        }
        }
    }
	
    var var_policy_options []Policy_options_Model
    if plan.Policy_options.IsNull() {
        var_policy_options = []Policy_options_Model{}
    }else {
        resp.Diagnostics.Append(plan.Policy_options.ElementsAs(ctx, &var_policy_options, false)...)
        if resp.Diagnostics.HasError() {
            return
        }
    }
    config.Groups.Policy_options = make([]xml_Policy_options, len(var_policy_options))
   
    for i_policy_options, v_policy_options := range var_policy_options {
        var var_policy_options_policy_statement []Policy_options_Policy_statement_Model
        resp.Diagnostics.Append(v_policy_options.Policy_statement.ElementsAs(ctx, &var_policy_options_policy_statement, false)...)
        if resp.Diagnostics.HasError() {
            return
        }
	    config.Groups.Policy_options[i_policy_options].Policy_statement = make([]xml_Policy_options_Policy_statement, len(var_policy_options_policy_statement))
        
		for i_policy_options_policy_statement, v_policy_options_policy_statement := range var_policy_options_policy_statement {
            config.Groups.Policy_options[i_policy_options].Policy_statement[i_policy_options_policy_statement].Name = v_policy_options_policy_statement.Name.ValueStringPointer()
            var var_policy_options_policy_statement_term []Policy_options_Policy_statement_Term_Model
            resp.Diagnostics.Append(v_policy_options_policy_statement.Term.ElementsAs(ctx, &var_policy_options_policy_statement_term, false)...)
            if resp.Diagnostics.HasError() {
                return
            }
	    config.Groups.Policy_options[i_policy_options].Policy_statement[i_policy_options_policy_statement].Term = make([]xml_Policy_options_Policy_statement_Term, len(var_policy_options_policy_statement_term))
        
		for i_policy_options_policy_statement_term, v_policy_options_policy_statement_term := range var_policy_options_policy_statement_term {
            config.Groups.Policy_options[i_policy_options].Policy_statement[i_policy_options_policy_statement].Term[i_policy_options_policy_statement_term].Name = v_policy_options_policy_statement_term.Name.ValueStringPointer()
            var var_policy_options_policy_statement_term_from []Policy_options_Policy_statement_Term_From_Model
            resp.Diagnostics.Append(v_policy_options_policy_statement_term.From.ElementsAs(ctx, &var_policy_options_policy_statement_term_from, false)...)
            if resp.Diagnostics.HasError() {
                return
            }
	    config.Groups.Policy_options[i_policy_options].Policy_statement[i_policy_options_policy_statement].Term[i_policy_options_policy_statement_term].From = make([]xml_Policy_options_Policy_statement_Term_From, len(var_policy_options_policy_statement_term_from))
        
		for i_policy_options_policy_statement_term_from, v_policy_options_policy_statement_term_from := range var_policy_options_policy_statement_term_from {
            var var_policy_options_policy_statement_term_from_route_filter []Policy_options_Policy_statement_Term_From_Route_filter_Model
            resp.Diagnostics.Append(v_policy_options_policy_statement_term_from.Route_filter.ElementsAs(ctx, &var_policy_options_policy_statement_term_from_route_filter, false)...)
            if resp.Diagnostics.HasError() {
                return
            }
	    config.Groups.Policy_options[i_policy_options].Policy_statement[i_policy_options_policy_statement].Term[i_policy_options_policy_statement_term].From[i_policy_options_policy_statement_term_from].Route_filter = make([]xml_Policy_options_Policy_statement_Term_From_Route_filter, len(var_policy_options_policy_statement_term_from_route_filter))
        
		for i_policy_options_policy_statement_term_from_route_filter, v_policy_options_policy_statement_term_from_route_filter := range var_policy_options_policy_statement_term_from_route_filter {
            config.Groups.Policy_options[i_policy_options].Policy_statement[i_policy_options_policy_statement].Term[i_policy_options_policy_statement_term].From[i_policy_options_policy_statement_term_from].Route_filter[i_policy_options_policy_statement_term_from_route_filter].Address = v_policy_options_policy_statement_term_from_route_filter.Address.ValueStringPointer()
            config.Groups.Policy_options[i_policy_options].Policy_statement[i_policy_options_policy_statement].Term[i_policy_options_policy_statement_term].From[i_policy_options_policy_statement_term_from].Route_filter[i_policy_options_policy_statement_term_from_route_filter].Exact = v_policy_options_policy_statement_term_from_route_filter.Exact.ValueStringPointer()
            config.Groups.Policy_options[i_policy_options].Policy_statement[i_policy_options_policy_statement].Term[i_policy_options_policy_statement_term].From[i_policy_options_policy_statement_term_from].Route_filter[i_policy_options_policy_statement_term_from_route_filter].Accept = v_policy_options_policy_statement_term_from_route_filter.Accept.ValueStringPointer()
        }
        }
        }
        }
    }
	
    var var_protocols []Protocols_Model
    if plan.Protocols.IsNull() {
        var_protocols = []Protocols_Model{}
    }else {
        resp.Diagnostics.Append(plan.Protocols.ElementsAs(ctx, &var_protocols, false)...)
        if resp.Diagnostics.HasError() {
            return
        }
    }
    config.Groups.Protocols = make([]xml_Protocols, len(var_protocols))
   
    for i_protocols, v_protocols := range var_protocols {
        var var_protocols_lldp []Protocols_Lldp_Model
        resp.Diagnostics.Append(v_protocols.Lldp.ElementsAs(ctx, &var_protocols_lldp, false)...)
        if resp.Diagnostics.HasError() {
            return
        }
	    config.Groups.Protocols[i_protocols].Lldp = make([]xml_Protocols_Lldp, len(var_protocols_lldp))
        
		for i_protocols_lldp, v_protocols_lldp := range var_protocols_lldp {
            var var_protocols_lldp_interface []Protocols_Lldp_Interface_Model
            resp.Diagnostics.Append(v_protocols_lldp.Interface.ElementsAs(ctx, &var_protocols_lldp_interface, false)...)
            if resp.Diagnostics.HasError() {
                return
            }
	    config.Groups.Protocols[i_protocols].Lldp[i_protocols_lldp].Interface = make([]xml_Protocols_Lldp_Interface, len(var_protocols_lldp_interface))
        
		for i_protocols_lldp_interface, v_protocols_lldp_interface := range var_protocols_lldp_interface {
            config.Groups.Protocols[i_protocols].Lldp[i_protocols_lldp].Interface[i_protocols_lldp_interface].Name = v_protocols_lldp_interface.Name.ValueStringPointer()
        }
        }
    }
	
    var var_routing_instances []Routing_instances_Model
    if plan.Routing_instances.IsNull() {
        var_routing_instances = []Routing_instances_Model{}
    }else {
        resp.Diagnostics.Append(plan.Routing_instances.ElementsAs(ctx, &var_routing_instances, false)...)
        if resp.Diagnostics.HasError() {
            return
        }
    }
    config.Groups.Routing_instances = make([]xml_Routing_instances, len(var_routing_instances))
   
    for i_routing_instances, v_routing_instances := range var_routing_instances {
        var var_routing_instances_instance []Routing_instances_Instance_Model
        resp.Diagnostics.Append(v_routing_instances.Instance.ElementsAs(ctx, &var_routing_instances_instance, false)...)
        if resp.Diagnostics.HasError() {
            return
        }
	    config.Groups.Routing_instances[i_routing_instances].Instance = make([]xml_Routing_instances_Instance, len(var_routing_instances_instance))
        
		for i_routing_instances_instance, v_routing_instances_instance := range var_routing_instances_instance {
            config.Groups.Routing_instances[i_routing_instances].Instance[i_routing_instances_instance].Name = v_routing_instances_instance.Name.ValueStringPointer()
            config.Groups.Routing_instances[i_routing_instances].Instance[i_routing_instances_instance].Instance_type = v_routing_instances_instance.Instance_type.ValueStringPointer()
            var var_routing_instances_instance_interface []Routing_instances_Instance_Interface_Model
            resp.Diagnostics.Append(v_routing_instances_instance.Interface.ElementsAs(ctx, &var_routing_instances_instance_interface, false)...)
            if resp.Diagnostics.HasError() {
                return
            }
	    config.Groups.Routing_instances[i_routing_instances].Instance[i_routing_instances_instance].Interface = make([]xml_Routing_instances_Instance_Interface, len(var_routing_instances_instance_interface))
        
		for i_routing_instances_instance_interface, v_routing_instances_instance_interface := range var_routing_instances_instance_interface {
            config.Groups.Routing_instances[i_routing_instances].Instance[i_routing_instances_instance].Interface[i_routing_instances_instance_interface].Name = v_routing_instances_instance_interface.Name.ValueStringPointer()
        }
            var var_routing_instances_instance_routing_options []Routing_instances_Instance_Routing_options_Model
            resp.Diagnostics.Append(v_routing_instances_instance.Routing_options.ElementsAs(ctx, &var_routing_instances_instance_routing_options, false)...)
            if resp.Diagnostics.HasError() {
                return
            }
	    config.Groups.Routing_instances[i_routing_instances].Instance[i_routing_instances_instance].Routing_options = make([]xml_Routing_instances_Instance_Routing_options, len(var_routing_instances_instance_routing_options))
        
		for i_routing_instances_instance_routing_options, v_routing_instances_instance_routing_options := range var_routing_instances_instance_routing_options {
            var var_routing_instances_instance_routing_options_static []Routing_instances_Instance_Routing_options_Static_Model
            resp.Diagnostics.Append(v_routing_instances_instance_routing_options.Static.ElementsAs(ctx, &var_routing_instances_instance_routing_options_static, false)...)
            if resp.Diagnostics.HasError() {
                return
            }
	    config.Groups.Routing_instances[i_routing_instances].Instance[i_routing_instances_instance].Routing_options[i_routing_instances_instance_routing_options].Static = make([]xml_Routing_instances_Instance_Routing_options_Static, len(var_routing_instances_instance_routing_options_static))
        
		for i_routing_instances_instance_routing_options_static, v_routing_instances_instance_routing_options_static := range var_routing_instances_instance_routing_options_static {
            var var_routing_instances_instance_routing_options_static_route []Routing_instances_Instance_Routing_options_Static_Route_Model
            resp.Diagnostics.Append(v_routing_instances_instance_routing_options_static.Route.ElementsAs(ctx, &var_routing_instances_instance_routing_options_static_route, false)...)
            if resp.Diagnostics.HasError() {
                return
            }
	    config.Groups.Routing_instances[i_routing_instances].Instance[i_routing_instances_instance].Routing_options[i_routing_instances_instance_routing_options].Static[i_routing_instances_instance_routing_options_static].Route = make([]xml_Routing_instances_Instance_Routing_options_Static_Route, len(var_routing_instances_instance_routing_options_static_route))
        
		for i_routing_instances_instance_routing_options_static_route, v_routing_instances_instance_routing_options_static_route := range var_routing_instances_instance_routing_options_static_route {
            config.Groups.Routing_instances[i_routing_instances].Instance[i_routing_instances_instance].Routing_options[i_routing_instances_instance_routing_options].Static[i_routing_instances_instance_routing_options_static].Route[i_routing_instances_instance_routing_options_static_route].Name = v_routing_instances_instance_routing_options_static_route.Name.ValueStringPointer()
            config.Groups.Routing_instances[i_routing_instances].Instance[i_routing_instances_instance].Routing_options[i_routing_instances_instance_routing_options].Static[i_routing_instances_instance_routing_options_static].Route[i_routing_instances_instance_routing_options_static_route].Discard = v_routing_instances_instance_routing_options_static_route.Discard.ValueStringPointer()
        }
        }
        }
            var var_routing_instances_instance_protocols []Routing_instances_Instance_Protocols_Model
            resp.Diagnostics.Append(v_routing_instances_instance.Protocols.ElementsAs(ctx, &var_routing_instances_instance_protocols, false)...)
            if resp.Diagnostics.HasError() {
                return
            }
	    config.Groups.Routing_instances[i_routing_instances].Instance[i_routing_instances_instance].Protocols = make([]xml_Routing_instances_Instance_Protocols, len(var_routing_instances_instance_protocols))
        
		for i_routing_instances_instance_protocols, v_routing_instances_instance_protocols := range var_routing_instances_instance_protocols {
            var var_routing_instances_instance_protocols_ospf []Routing_instances_Instance_Protocols_Ospf_Model
            resp.Diagnostics.Append(v_routing_instances_instance_protocols.Ospf.ElementsAs(ctx, &var_routing_instances_instance_protocols_ospf, false)...)
            if resp.Diagnostics.HasError() {
                return
            }
	    config.Groups.Routing_instances[i_routing_instances].Instance[i_routing_instances_instance].Protocols[i_routing_instances_instance_protocols].Ospf = make([]xml_Routing_instances_Instance_Protocols_Ospf, len(var_routing_instances_instance_protocols_ospf))
        
		for i_routing_instances_instance_protocols_ospf, v_routing_instances_instance_protocols_ospf := range var_routing_instances_instance_protocols_ospf {
			var var_routing_instances_instance_protocols_ospf_export []string
			resp.Diagnostics.Append(v_routing_instances_instance_protocols_ospf.Export.ElementsAs(ctx, &var_routing_instances_instance_protocols_ospf_export, false)...)
			for _, v_routing_instances_instance_protocols_ospf_export := range var_routing_instances_instance_protocols_ospf_export {
				config.Groups.Routing_instances[i_routing_instances].Instance[i_routing_instances_instance].Protocols[i_routing_instances_instance_protocols].Ospf[i_routing_instances_instance_protocols_ospf].Export = append(config.Groups.Routing_instances[i_routing_instances].Instance[i_routing_instances_instance].Protocols[i_routing_instances_instance_protocols].Ospf[i_routing_instances_instance_protocols_ospf].Export, &v_routing_instances_instance_protocols_ospf_export)
			}
            var var_routing_instances_instance_protocols_ospf_area []Routing_instances_Instance_Protocols_Ospf_Area_Model
            resp.Diagnostics.Append(v_routing_instances_instance_protocols_ospf.Area.ElementsAs(ctx, &var_routing_instances_instance_protocols_ospf_area, false)...)
            if resp.Diagnostics.HasError() {
                return
            }
	    config.Groups.Routing_instances[i_routing_instances].Instance[i_routing_instances_instance].Protocols[i_routing_instances_instance_protocols].Ospf[i_routing_instances_instance_protocols_ospf].Area = make([]xml_Routing_instances_Instance_Protocols_Ospf_Area, len(var_routing_instances_instance_protocols_ospf_area))
        
		for i_routing_instances_instance_protocols_ospf_area, v_routing_instances_instance_protocols_ospf_area := range var_routing_instances_instance_protocols_ospf_area {
            config.Groups.Routing_instances[i_routing_instances].Instance[i_routing_instances_instance].Protocols[i_routing_instances_instance_protocols].Ospf[i_routing_instances_instance_protocols_ospf].Area[i_routing_instances_instance_protocols_ospf_area].Name = v_routing_instances_instance_protocols_ospf_area.Name.ValueStringPointer()
            var var_routing_instances_instance_protocols_ospf_area_interface []Routing_instances_Instance_Protocols_Ospf_Area_Interface_Model
            resp.Diagnostics.Append(v_routing_instances_instance_protocols_ospf_area.Interface.ElementsAs(ctx, &var_routing_instances_instance_protocols_ospf_area_interface, false)...)
            if resp.Diagnostics.HasError() {
                return
            }
	    config.Groups.Routing_instances[i_routing_instances].Instance[i_routing_instances_instance].Protocols[i_routing_instances_instance_protocols].Ospf[i_routing_instances_instance_protocols_ospf].Area[i_routing_instances_instance_protocols_ospf_area].Interface = make([]xml_Routing_instances_Instance_Protocols_Ospf_Area_Interface, len(var_routing_instances_instance_protocols_ospf_area_interface))
        
		for i_routing_instances_instance_protocols_ospf_area_interface, v_routing_instances_instance_protocols_ospf_area_interface := range var_routing_instances_instance_protocols_ospf_area_interface {
            config.Groups.Routing_instances[i_routing_instances].Instance[i_routing_instances_instance].Protocols[i_routing_instances_instance_protocols].Ospf[i_routing_instances_instance_protocols_ospf].Area[i_routing_instances_instance_protocols_ospf_area].Interface[i_routing_instances_instance_protocols_ospf_area_interface].Name = v_routing_instances_instance_protocols_ospf_area_interface.Name.ValueStringPointer()
        }
        }
        }
        }
        }
    }
	
    var var_routing_options []Routing_options_Model
    if plan.Routing_options.IsNull() {
        var_routing_options = []Routing_options_Model{}
    }else {
        resp.Diagnostics.Append(plan.Routing_options.ElementsAs(ctx, &var_routing_options, false)...)
        if resp.Diagnostics.HasError() {
            return
        }
    }
    config.Groups.Routing_options = make([]xml_Routing_options, len(var_routing_options))
   
    for i_routing_options, v_routing_options := range var_routing_options {
        var var_routing_options_static []Routing_options_Static_Model
        resp.Diagnostics.Append(v_routing_options.Static.ElementsAs(ctx, &var_routing_options_static, false)...)
        if resp.Diagnostics.HasError() {
            return
        }
	    config.Groups.Routing_options[i_routing_options].Static = make([]xml_Routing_options_Static, len(var_routing_options_static))
        
		for i_routing_options_static, v_routing_options_static := range var_routing_options_static {
            var var_routing_options_static_route []Routing_options_Static_Route_Model
            resp.Diagnostics.Append(v_routing_options_static.Route.ElementsAs(ctx, &var_routing_options_static_route, false)...)
            if resp.Diagnostics.HasError() {
                return
            }
	    config.Groups.Routing_options[i_routing_options].Static[i_routing_options_static].Route = make([]xml_Routing_options_Static_Route, len(var_routing_options_static_route))
        
		for i_routing_options_static_route, v_routing_options_static_route := range var_routing_options_static_route {
            config.Groups.Routing_options[i_routing_options].Static[i_routing_options_static].Route[i_routing_options_static_route].Name = v_routing_options_static_route.Name.ValueStringPointer()
			var var_routing_options_static_route_next_hop []string
			resp.Diagnostics.Append(v_routing_options_static_route.Next_hop.ElementsAs(ctx, &var_routing_options_static_route_next_hop, false)...)
			for _, v_routing_options_static_route_next_hop := range var_routing_options_static_route_next_hop {
				config.Groups.Routing_options[i_routing_options].Static[i_routing_options_static].Route[i_routing_options_static_route].Next_hop = append(config.Groups.Routing_options[i_routing_options].Static[i_routing_options_static].Route[i_routing_options_static_route].Next_hop, &v_routing_options_static_route_next_hop)
			}
        }
        }
    }
	
    var var_security []Security_Model
    if plan.Security.IsNull() {
        var_security = []Security_Model{}
    }else {
        resp.Diagnostics.Append(plan.Security.ElementsAs(ctx, &var_security, false)...)
        if resp.Diagnostics.HasError() {
            return
        }
    }
    config.Groups.Security = make([]xml_Security, len(var_security))
   
    for i_security, v_security := range var_security {
        var var_security_log []Security_Log_Model
        resp.Diagnostics.Append(v_security.Log.ElementsAs(ctx, &var_security_log, false)...)
        if resp.Diagnostics.HasError() {
            return
        }
	    config.Groups.Security[i_security].Log = make([]xml_Security_Log, len(var_security_log))
        
		for i_security_log, v_security_log := range var_security_log {
            config.Groups.Security[i_security].Log[i_security_log].Mode = v_security_log.Mode.ValueStringPointer()
            var var_security_log_report []Security_Log_Report_Model
            resp.Diagnostics.Append(v_security_log.Report.ElementsAs(ctx, &var_security_log_report, false)...)
            if resp.Diagnostics.HasError() {
                return
            }
	    config.Groups.Security[i_security].Log[i_security_log].Report = make([]xml_Security_Log_Report, len(var_security_log_report))
        
        }
        var var_security_screen []Security_Screen_Model
        resp.Diagnostics.Append(v_security.Screen.ElementsAs(ctx, &var_security_screen, false)...)
        if resp.Diagnostics.HasError() {
            return
        }
	    config.Groups.Security[i_security].Screen = make([]xml_Security_Screen, len(var_security_screen))
        
		for i_security_screen, v_security_screen := range var_security_screen {
            var var_security_screen_ids_option []Security_Screen_Ids_option_Model
            resp.Diagnostics.Append(v_security_screen.Ids_option.ElementsAs(ctx, &var_security_screen_ids_option, false)...)
            if resp.Diagnostics.HasError() {
                return
            }
	    config.Groups.Security[i_security].Screen[i_security_screen].Ids_option = make([]xml_Security_Screen_Ids_option, len(var_security_screen_ids_option))
        
		for i_security_screen_ids_option, v_security_screen_ids_option := range var_security_screen_ids_option {
            config.Groups.Security[i_security].Screen[i_security_screen].Ids_option[i_security_screen_ids_option].Name = v_security_screen_ids_option.Name.ValueStringPointer()
            var var_security_screen_ids_option_icmp []Security_Screen_Ids_option_Icmp_Model
            resp.Diagnostics.Append(v_security_screen_ids_option.Icmp.ElementsAs(ctx, &var_security_screen_ids_option_icmp, false)...)
            if resp.Diagnostics.HasError() {
                return
            }
	    config.Groups.Security[i_security].Screen[i_security_screen].Ids_option[i_security_screen_ids_option].Icmp = make([]xml_Security_Screen_Ids_option_Icmp, len(var_security_screen_ids_option_icmp))
        
		for i_security_screen_ids_option_icmp, v_security_screen_ids_option_icmp := range var_security_screen_ids_option_icmp {
            config.Groups.Security[i_security].Screen[i_security_screen].Ids_option[i_security_screen_ids_option].Icmp[i_security_screen_ids_option_icmp].Ping_death = v_security_screen_ids_option_icmp.Ping_death.ValueStringPointer()
        }
            var var_security_screen_ids_option_ip []Security_Screen_Ids_option_Ip_Model
            resp.Diagnostics.Append(v_security_screen_ids_option.Ip.ElementsAs(ctx, &var_security_screen_ids_option_ip, false)...)
            if resp.Diagnostics.HasError() {
                return
            }
	    config.Groups.Security[i_security].Screen[i_security_screen].Ids_option[i_security_screen_ids_option].Ip = make([]xml_Security_Screen_Ids_option_Ip, len(var_security_screen_ids_option_ip))
        
		for i_security_screen_ids_option_ip, v_security_screen_ids_option_ip := range var_security_screen_ids_option_ip {
            config.Groups.Security[i_security].Screen[i_security_screen].Ids_option[i_security_screen_ids_option].Ip[i_security_screen_ids_option_ip].Source_route_option = v_security_screen_ids_option_ip.Source_route_option.ValueStringPointer()
            config.Groups.Security[i_security].Screen[i_security_screen].Ids_option[i_security_screen_ids_option].Ip[i_security_screen_ids_option_ip].Tear_drop = v_security_screen_ids_option_ip.Tear_drop.ValueStringPointer()
        }
            var var_security_screen_ids_option_tcp []Security_Screen_Ids_option_Tcp_Model
            resp.Diagnostics.Append(v_security_screen_ids_option.Tcp.ElementsAs(ctx, &var_security_screen_ids_option_tcp, false)...)
            if resp.Diagnostics.HasError() {
                return
            }
	    config.Groups.Security[i_security].Screen[i_security_screen].Ids_option[i_security_screen_ids_option].Tcp = make([]xml_Security_Screen_Ids_option_Tcp, len(var_security_screen_ids_option_tcp))
        
		for i_security_screen_ids_option_tcp, v_security_screen_ids_option_tcp := range var_security_screen_ids_option_tcp {
            var var_security_screen_ids_option_tcp_syn_flood []Security_Screen_Ids_option_Tcp_Syn_flood_Model
            resp.Diagnostics.Append(v_security_screen_ids_option_tcp.Syn_flood.ElementsAs(ctx, &var_security_screen_ids_option_tcp_syn_flood, false)...)
            if resp.Diagnostics.HasError() {
                return
            }
	    config.Groups.Security[i_security].Screen[i_security_screen].Ids_option[i_security_screen_ids_option].Tcp[i_security_screen_ids_option_tcp].Syn_flood = make([]xml_Security_Screen_Ids_option_Tcp_Syn_flood, len(var_security_screen_ids_option_tcp_syn_flood))
        
		for i_security_screen_ids_option_tcp_syn_flood, v_security_screen_ids_option_tcp_syn_flood := range var_security_screen_ids_option_tcp_syn_flood {
            config.Groups.Security[i_security].Screen[i_security_screen].Ids_option[i_security_screen_ids_option].Tcp[i_security_screen_ids_option_tcp].Syn_flood[i_security_screen_ids_option_tcp_syn_flood].Alarm_threshold = v_security_screen_ids_option_tcp_syn_flood.Alarm_threshold.ValueStringPointer()
            config.Groups.Security[i_security].Screen[i_security_screen].Ids_option[i_security_screen_ids_option].Tcp[i_security_screen_ids_option_tcp].Syn_flood[i_security_screen_ids_option_tcp_syn_flood].Attack_threshold = v_security_screen_ids_option_tcp_syn_flood.Attack_threshold.ValueStringPointer()
            config.Groups.Security[i_security].Screen[i_security_screen].Ids_option[i_security_screen_ids_option].Tcp[i_security_screen_ids_option_tcp].Syn_flood[i_security_screen_ids_option_tcp_syn_flood].Source_threshold = v_security_screen_ids_option_tcp_syn_flood.Source_threshold.ValueStringPointer()
            config.Groups.Security[i_security].Screen[i_security_screen].Ids_option[i_security_screen_ids_option].Tcp[i_security_screen_ids_option_tcp].Syn_flood[i_security_screen_ids_option_tcp_syn_flood].Destination_threshold = v_security_screen_ids_option_tcp_syn_flood.Destination_threshold.ValueStringPointer()
            config.Groups.Security[i_security].Screen[i_security_screen].Ids_option[i_security_screen_ids_option].Tcp[i_security_screen_ids_option_tcp].Syn_flood[i_security_screen_ids_option_tcp_syn_flood].Timeout = v_security_screen_ids_option_tcp_syn_flood.Timeout.ValueStringPointer()
        }
            config.Groups.Security[i_security].Screen[i_security_screen].Ids_option[i_security_screen_ids_option].Tcp[i_security_screen_ids_option_tcp].Land = v_security_screen_ids_option_tcp.Land.ValueStringPointer()
        }
        }
        }
        var var_security_policies []Security_Policies_Model
        resp.Diagnostics.Append(v_security.Policies.ElementsAs(ctx, &var_security_policies, false)...)
        if resp.Diagnostics.HasError() {
            return
        }
	    config.Groups.Security[i_security].Policies = make([]xml_Security_Policies, len(var_security_policies))
        
		for i_security_policies, v_security_policies := range var_security_policies {
            var var_security_policies_policy []Security_Policies_Policy_Model
            resp.Diagnostics.Append(v_security_policies.Policy.ElementsAs(ctx, &var_security_policies_policy, false)...)
            if resp.Diagnostics.HasError() {
                return
            }
	    config.Groups.Security[i_security].Policies[i_security_policies].Policy = make([]xml_Security_Policies_Policy, len(var_security_policies_policy))
        
		for i_security_policies_policy, v_security_policies_policy := range var_security_policies_policy {
            config.Groups.Security[i_security].Policies[i_security_policies].Policy[i_security_policies_policy].From_zone_name = v_security_policies_policy.From_zone_name.ValueStringPointer()
            config.Groups.Security[i_security].Policies[i_security_policies].Policy[i_security_policies_policy].To_zone_name = v_security_policies_policy.To_zone_name.ValueStringPointer()
            var var_security_policies_policy_policy []Security_Policies_Policy_Policy_Model
            resp.Diagnostics.Append(v_security_policies_policy.Policy.ElementsAs(ctx, &var_security_policies_policy_policy, false)...)
            if resp.Diagnostics.HasError() {
                return
            }
	    config.Groups.Security[i_security].Policies[i_security_policies].Policy[i_security_policies_policy].Policy = make([]xml_Security_Policies_Policy_Policy, len(var_security_policies_policy_policy))
        
		for i_security_policies_policy_policy, v_security_policies_policy_policy := range var_security_policies_policy_policy {
            config.Groups.Security[i_security].Policies[i_security_policies].Policy[i_security_policies_policy].Policy[i_security_policies_policy_policy].Name = v_security_policies_policy_policy.Name.ValueStringPointer()
            var var_security_policies_policy_policy_match []Security_Policies_Policy_Policy_Match_Model
            resp.Diagnostics.Append(v_security_policies_policy_policy.Match.ElementsAs(ctx, &var_security_policies_policy_policy_match, false)...)
            if resp.Diagnostics.HasError() {
                return
            }
	    config.Groups.Security[i_security].Policies[i_security_policies].Policy[i_security_policies_policy].Policy[i_security_policies_policy_policy].Match = make([]xml_Security_Policies_Policy_Policy_Match, len(var_security_policies_policy_policy_match))
        
		for i_security_policies_policy_policy_match, v_security_policies_policy_policy_match := range var_security_policies_policy_policy_match {
			var var_security_policies_policy_policy_match_source_address []string
			resp.Diagnostics.Append(v_security_policies_policy_policy_match.Source_address.ElementsAs(ctx, &var_security_policies_policy_policy_match_source_address, false)...)
			for _, v_security_policies_policy_policy_match_source_address := range var_security_policies_policy_policy_match_source_address {
				config.Groups.Security[i_security].Policies[i_security_policies].Policy[i_security_policies_policy].Policy[i_security_policies_policy_policy].Match[i_security_policies_policy_policy_match].Source_address = append(config.Groups.Security[i_security].Policies[i_security_policies].Policy[i_security_policies_policy].Policy[i_security_policies_policy_policy].Match[i_security_policies_policy_policy_match].Source_address, &v_security_policies_policy_policy_match_source_address)
			}
			var var_security_policies_policy_policy_match_destination_address []string
			resp.Diagnostics.Append(v_security_policies_policy_policy_match.Destination_address.ElementsAs(ctx, &var_security_policies_policy_policy_match_destination_address, false)...)
			for _, v_security_policies_policy_policy_match_destination_address := range var_security_policies_policy_policy_match_destination_address {
				config.Groups.Security[i_security].Policies[i_security_policies].Policy[i_security_policies_policy].Policy[i_security_policies_policy_policy].Match[i_security_policies_policy_policy_match].Destination_address = append(config.Groups.Security[i_security].Policies[i_security_policies].Policy[i_security_policies_policy].Policy[i_security_policies_policy_policy].Match[i_security_policies_policy_policy_match].Destination_address, &v_security_policies_policy_policy_match_destination_address)
			}
			var var_security_policies_policy_policy_match_application []string
			resp.Diagnostics.Append(v_security_policies_policy_policy_match.Application.ElementsAs(ctx, &var_security_policies_policy_policy_match_application, false)...)
			for _, v_security_policies_policy_policy_match_application := range var_security_policies_policy_policy_match_application {
				config.Groups.Security[i_security].Policies[i_security_policies].Policy[i_security_policies_policy].Policy[i_security_policies_policy_policy].Match[i_security_policies_policy_policy_match].Application = append(config.Groups.Security[i_security].Policies[i_security_policies].Policy[i_security_policies_policy].Policy[i_security_policies_policy_policy].Match[i_security_policies_policy_policy_match].Application, &v_security_policies_policy_policy_match_application)
			}
        }
            var var_security_policies_policy_policy_then []Security_Policies_Policy_Policy_Then_Model
            resp.Diagnostics.Append(v_security_policies_policy_policy.Then.ElementsAs(ctx, &var_security_policies_policy_policy_then, false)...)
            if resp.Diagnostics.HasError() {
                return
            }
	    config.Groups.Security[i_security].Policies[i_security_policies].Policy[i_security_policies_policy].Policy[i_security_policies_policy_policy].Then = make([]xml_Security_Policies_Policy_Policy_Then, len(var_security_policies_policy_policy_then))
        
		for i_security_policies_policy_policy_then, v_security_policies_policy_policy_then := range var_security_policies_policy_policy_then {
            var var_security_policies_policy_policy_then_permit []Security_Policies_Policy_Policy_Then_Permit_Model
            resp.Diagnostics.Append(v_security_policies_policy_policy_then.Permit.ElementsAs(ctx, &var_security_policies_policy_policy_then_permit, false)...)
            if resp.Diagnostics.HasError() {
                return
            }
	    config.Groups.Security[i_security].Policies[i_security_policies].Policy[i_security_policies_policy].Policy[i_security_policies_policy_policy].Then[i_security_policies_policy_policy_then].Permit = make([]xml_Security_Policies_Policy_Policy_Then_Permit, len(var_security_policies_policy_policy_then_permit))
        
        }
        }
        }
        }
        var var_security_zones []Security_Zones_Model
        resp.Diagnostics.Append(v_security.Zones.ElementsAs(ctx, &var_security_zones, false)...)
        if resp.Diagnostics.HasError() {
            return
        }
	    config.Groups.Security[i_security].Zones = make([]xml_Security_Zones, len(var_security_zones))
        
		for i_security_zones, v_security_zones := range var_security_zones {
            var var_security_zones_security_zone []Security_Zones_Security_zone_Model
            resp.Diagnostics.Append(v_security_zones.Security_zone.ElementsAs(ctx, &var_security_zones_security_zone, false)...)
            if resp.Diagnostics.HasError() {
                return
            }
	    config.Groups.Security[i_security].Zones[i_security_zones].Security_zone = make([]xml_Security_Zones_Security_zone, len(var_security_zones_security_zone))
        
		for i_security_zones_security_zone, v_security_zones_security_zone := range var_security_zones_security_zone {
            config.Groups.Security[i_security].Zones[i_security_zones].Security_zone[i_security_zones_security_zone].Name = v_security_zones_security_zone.Name.ValueStringPointer()
            config.Groups.Security[i_security].Zones[i_security_zones].Security_zone[i_security_zones_security_zone].Tcp_rst = v_security_zones_security_zone.Tcp_rst.ValueStringPointer()
            config.Groups.Security[i_security].Zones[i_security_zones].Security_zone[i_security_zones_security_zone].Screen = v_security_zones_security_zone.Screen.ValueStringPointer()
            var var_security_zones_security_zone_host_inbound_traffic []Security_Zones_Security_zone_Host_inbound_traffic_Model
            resp.Diagnostics.Append(v_security_zones_security_zone.Host_inbound_traffic.ElementsAs(ctx, &var_security_zones_security_zone_host_inbound_traffic, false)...)
            if resp.Diagnostics.HasError() {
                return
            }
	    config.Groups.Security[i_security].Zones[i_security_zones].Security_zone[i_security_zones_security_zone].Host_inbound_traffic = make([]xml_Security_Zones_Security_zone_Host_inbound_traffic, len(var_security_zones_security_zone_host_inbound_traffic))
        
		for i_security_zones_security_zone_host_inbound_traffic, v_security_zones_security_zone_host_inbound_traffic := range var_security_zones_security_zone_host_inbound_traffic {
            var var_security_zones_security_zone_host_inbound_traffic_system_services []Security_Zones_Security_zone_Host_inbound_traffic_System_services_Model
            resp.Diagnostics.Append(v_security_zones_security_zone_host_inbound_traffic.System_services.ElementsAs(ctx, &var_security_zones_security_zone_host_inbound_traffic_system_services, false)...)
            if resp.Diagnostics.HasError() {
                return
            }
	    config.Groups.Security[i_security].Zones[i_security_zones].Security_zone[i_security_zones_security_zone].Host_inbound_traffic[i_security_zones_security_zone_host_inbound_traffic].System_services = make([]xml_Security_Zones_Security_zone_Host_inbound_traffic_System_services, len(var_security_zones_security_zone_host_inbound_traffic_system_services))
        
		for i_security_zones_security_zone_host_inbound_traffic_system_services, v_security_zones_security_zone_host_inbound_traffic_system_services := range var_security_zones_security_zone_host_inbound_traffic_system_services {
            config.Groups.Security[i_security].Zones[i_security_zones].Security_zone[i_security_zones_security_zone].Host_inbound_traffic[i_security_zones_security_zone_host_inbound_traffic].System_services[i_security_zones_security_zone_host_inbound_traffic_system_services].Name = v_security_zones_security_zone_host_inbound_traffic_system_services.Name.ValueStringPointer()
        }
            var var_security_zones_security_zone_host_inbound_traffic_protocols []Security_Zones_Security_zone_Host_inbound_traffic_Protocols_Model
            resp.Diagnostics.Append(v_security_zones_security_zone_host_inbound_traffic.Protocols.ElementsAs(ctx, &var_security_zones_security_zone_host_inbound_traffic_protocols, false)...)
            if resp.Diagnostics.HasError() {
                return
            }
	    config.Groups.Security[i_security].Zones[i_security_zones].Security_zone[i_security_zones_security_zone].Host_inbound_traffic[i_security_zones_security_zone_host_inbound_traffic].Protocols = make([]xml_Security_Zones_Security_zone_Host_inbound_traffic_Protocols, len(var_security_zones_security_zone_host_inbound_traffic_protocols))
        
		for i_security_zones_security_zone_host_inbound_traffic_protocols, v_security_zones_security_zone_host_inbound_traffic_protocols := range var_security_zones_security_zone_host_inbound_traffic_protocols {
            config.Groups.Security[i_security].Zones[i_security_zones].Security_zone[i_security_zones_security_zone].Host_inbound_traffic[i_security_zones_security_zone_host_inbound_traffic].Protocols[i_security_zones_security_zone_host_inbound_traffic_protocols].Name = v_security_zones_security_zone_host_inbound_traffic_protocols.Name.ValueStringPointer()
        }
        }
            var var_security_zones_security_zone_interfaces []Security_Zones_Security_zone_Interfaces_Model
            resp.Diagnostics.Append(v_security_zones_security_zone.Interfaces.ElementsAs(ctx, &var_security_zones_security_zone_interfaces, false)...)
            if resp.Diagnostics.HasError() {
                return
            }
	    config.Groups.Security[i_security].Zones[i_security_zones].Security_zone[i_security_zones_security_zone].Interfaces = make([]xml_Security_Zones_Security_zone_Interfaces, len(var_security_zones_security_zone_interfaces))
        
		for i_security_zones_security_zone_interfaces, v_security_zones_security_zone_interfaces := range var_security_zones_security_zone_interfaces {
            config.Groups.Security[i_security].Zones[i_security_zones].Security_zone[i_security_zones_security_zone].Interfaces[i_security_zones_security_zone_interfaces].Name = v_security_zones_security_zone_interfaces.Name.ValueStringPointer()
        }
        }
        }
    }
	
    var var_snmp []Snmp_Model
    if plan.Snmp.IsNull() {
        var_snmp = []Snmp_Model{}
    }else {
        resp.Diagnostics.Append(plan.Snmp.ElementsAs(ctx, &var_snmp, false)...)
        if resp.Diagnostics.HasError() {
            return
        }
    }
    config.Groups.Snmp = make([]xml_Snmp, len(var_snmp))
   
    for i_snmp, v_snmp := range var_snmp {
        config.Groups.Snmp[i_snmp].Location = v_snmp.Location.ValueStringPointer()
        config.Groups.Snmp[i_snmp].Contact = v_snmp.Contact.ValueStringPointer()
        var var_snmp_community []Snmp_Community_Model
        resp.Diagnostics.Append(v_snmp.Community.ElementsAs(ctx, &var_snmp_community, false)...)
        if resp.Diagnostics.HasError() {
            return
        }
	    config.Groups.Snmp[i_snmp].Community = make([]xml_Snmp_Community, len(var_snmp_community))
        
		for i_snmp_community, v_snmp_community := range var_snmp_community {
            config.Groups.Snmp[i_snmp].Community[i_snmp_community].Name = v_snmp_community.Name.ValueStringPointer()
            config.Groups.Snmp[i_snmp].Community[i_snmp_community].Authorization = v_snmp_community.Authorization.ValueStringPointer()
        }
    }
	
    var var_system []System_Model
    if plan.System.IsNull() {
        var_system = []System_Model{}
    }else {
        resp.Diagnostics.Append(plan.System.ElementsAs(ctx, &var_system, false)...)
        if resp.Diagnostics.HasError() {
            return
        }
    }
    config.Groups.System = make([]xml_System, len(var_system))
   
    for i_system, v_system := range var_system {
        var var_system_login []System_Login_Model
        resp.Diagnostics.Append(v_system.Login.ElementsAs(ctx, &var_system_login, false)...)
        if resp.Diagnostics.HasError() {
            return
        }
	    config.Groups.System[i_system].Login = make([]xml_System_Login, len(var_system_login))
        
		for i_system_login, v_system_login := range var_system_login {
            var var_system_login_user []System_Login_User_Model
            resp.Diagnostics.Append(v_system_login.User.ElementsAs(ctx, &var_system_login_user, false)...)
            if resp.Diagnostics.HasError() {
                return
            }
	    config.Groups.System[i_system].Login[i_system_login].User = make([]xml_System_Login_User, len(var_system_login_user))
        
		for i_system_login_user, v_system_login_user := range var_system_login_user {
            config.Groups.System[i_system].Login[i_system_login].User[i_system_login_user].Name = v_system_login_user.Name.ValueStringPointer()
            config.Groups.System[i_system].Login[i_system_login].User[i_system_login_user].Uid = v_system_login_user.Uid.ValueStringPointer()
            config.Groups.System[i_system].Login[i_system_login].User[i_system_login_user].Class = v_system_login_user.Class.ValueStringPointer()
            var var_system_login_user_authentication []System_Login_User_Authentication_Model
            resp.Diagnostics.Append(v_system_login_user.Authentication.ElementsAs(ctx, &var_system_login_user_authentication, false)...)
            if resp.Diagnostics.HasError() {
                return
            }
	    config.Groups.System[i_system].Login[i_system_login].User[i_system_login_user].Authentication = make([]xml_System_Login_User_Authentication, len(var_system_login_user_authentication))
        
		for i_system_login_user_authentication, v_system_login_user_authentication := range var_system_login_user_authentication {
            config.Groups.System[i_system].Login[i_system_login].User[i_system_login_user].Authentication[i_system_login_user_authentication].Encrypted_password = v_system_login_user_authentication.Encrypted_password.ValueStringPointer()
        }
        }
            config.Groups.System[i_system].Login[i_system_login].Message = v_system_login.Message.ValueStringPointer()
        }
        var var_system_root_authentication []System_Root_authentication_Model
        resp.Diagnostics.Append(v_system.Root_authentication.ElementsAs(ctx, &var_system_root_authentication, false)...)
        if resp.Diagnostics.HasError() {
            return
        }
	    config.Groups.System[i_system].Root_authentication = make([]xml_System_Root_authentication, len(var_system_root_authentication))
        
		for i_system_root_authentication, v_system_root_authentication := range var_system_root_authentication {
            config.Groups.System[i_system].Root_authentication[i_system_root_authentication].Encrypted_password = v_system_root_authentication.Encrypted_password.ValueStringPointer()
        }
        config.Groups.System[i_system].Host_name = v_system.Host_name.ValueStringPointer()
        var var_system_services []System_Services_Model
        resp.Diagnostics.Append(v_system.Services.ElementsAs(ctx, &var_system_services, false)...)
        if resp.Diagnostics.HasError() {
            return
        }
	    config.Groups.System[i_system].Services = make([]xml_System_Services, len(var_system_services))
        
		for i_system_services, v_system_services := range var_system_services {
            var var_system_services_ssh []System_Services_Ssh_Model
            resp.Diagnostics.Append(v_system_services.Ssh.ElementsAs(ctx, &var_system_services_ssh, false)...)
            if resp.Diagnostics.HasError() {
                return
            }
	    config.Groups.System[i_system].Services[i_system_services].Ssh = make([]xml_System_Services_Ssh, len(var_system_services_ssh))
        
		for i_system_services_ssh, v_system_services_ssh := range var_system_services_ssh {
            config.Groups.System[i_system].Services[i_system_services].Ssh[i_system_services_ssh].Root_login = v_system_services_ssh.Root_login.ValueStringPointer()
        }
            var var_system_services_extension_service []System_Services_Extension_service_Model
            resp.Diagnostics.Append(v_system_services.Extension_service.ElementsAs(ctx, &var_system_services_extension_service, false)...)
            if resp.Diagnostics.HasError() {
                return
            }
	    config.Groups.System[i_system].Services[i_system_services].Extension_service = make([]xml_System_Services_Extension_service, len(var_system_services_extension_service))
        
		for i_system_services_extension_service, v_system_services_extension_service := range var_system_services_extension_service {
            var var_system_services_extension_service_request_response []System_Services_Extension_service_Request_response_Model
            resp.Diagnostics.Append(v_system_services_extension_service.Request_response.ElementsAs(ctx, &var_system_services_extension_service_request_response, false)...)
            if resp.Diagnostics.HasError() {
                return
            }
	    config.Groups.System[i_system].Services[i_system_services].Extension_service[i_system_services_extension_service].Request_response = make([]xml_System_Services_Extension_service_Request_response, len(var_system_services_extension_service_request_response))
        
		for i_system_services_extension_service_request_response, v_system_services_extension_service_request_response := range var_system_services_extension_service_request_response {
            var var_system_services_extension_service_request_response_grpc []System_Services_Extension_service_Request_response_Grpc_Model
            resp.Diagnostics.Append(v_system_services_extension_service_request_response.Grpc.ElementsAs(ctx, &var_system_services_extension_service_request_response_grpc, false)...)
            if resp.Diagnostics.HasError() {
                return
            }
	    config.Groups.System[i_system].Services[i_system_services].Extension_service[i_system_services_extension_service].Request_response[i_system_services_extension_service_request_response].Grpc = make([]xml_System_Services_Extension_service_Request_response_Grpc, len(var_system_services_extension_service_request_response_grpc))
        
		for i_system_services_extension_service_request_response_grpc, v_system_services_extension_service_request_response_grpc := range var_system_services_extension_service_request_response_grpc {
            config.Groups.System[i_system].Services[i_system_services].Extension_service[i_system_services_extension_service].Request_response[i_system_services_extension_service_request_response].Grpc[i_system_services_extension_service_request_response_grpc].Max_connections = v_system_services_extension_service_request_response_grpc.Max_connections.ValueStringPointer()
        }
        }
            var var_system_services_extension_service_notification []System_Services_Extension_service_Notification_Model
            resp.Diagnostics.Append(v_system_services_extension_service.Notification.ElementsAs(ctx, &var_system_services_extension_service_notification, false)...)
            if resp.Diagnostics.HasError() {
                return
            }
	    config.Groups.System[i_system].Services[i_system_services].Extension_service[i_system_services_extension_service].Notification = make([]xml_System_Services_Extension_service_Notification, len(var_system_services_extension_service_notification))
        
		for i_system_services_extension_service_notification, v_system_services_extension_service_notification := range var_system_services_extension_service_notification {
            var var_system_services_extension_service_notification_allow_clients []System_Services_Extension_service_Notification_Allow_clients_Model
            resp.Diagnostics.Append(v_system_services_extension_service_notification.Allow_clients.ElementsAs(ctx, &var_system_services_extension_service_notification_allow_clients, false)...)
            if resp.Diagnostics.HasError() {
                return
            }
	    config.Groups.System[i_system].Services[i_system_services].Extension_service[i_system_services_extension_service].Notification[i_system_services_extension_service_notification].Allow_clients = make([]xml_System_Services_Extension_service_Notification_Allow_clients, len(var_system_services_extension_service_notification_allow_clients))
        
		for i_system_services_extension_service_notification_allow_clients, v_system_services_extension_service_notification_allow_clients := range var_system_services_extension_service_notification_allow_clients {
			var var_system_services_extension_service_notification_allow_clients_address []string
			resp.Diagnostics.Append(v_system_services_extension_service_notification_allow_clients.Address.ElementsAs(ctx, &var_system_services_extension_service_notification_allow_clients_address, false)...)
			for _, v_system_services_extension_service_notification_allow_clients_address := range var_system_services_extension_service_notification_allow_clients_address {
				config.Groups.System[i_system].Services[i_system_services].Extension_service[i_system_services_extension_service].Notification[i_system_services_extension_service_notification].Allow_clients[i_system_services_extension_service_notification_allow_clients].Address = append(config.Groups.System[i_system].Services[i_system_services].Extension_service[i_system_services_extension_service].Notification[i_system_services_extension_service_notification].Allow_clients[i_system_services_extension_service_notification_allow_clients].Address, &v_system_services_extension_service_notification_allow_clients_address)
			}
        }
        }
        }
            var var_system_services_netconf []System_Services_Netconf_Model
            resp.Diagnostics.Append(v_system_services.Netconf.ElementsAs(ctx, &var_system_services_netconf, false)...)
            if resp.Diagnostics.HasError() {
                return
            }
	    config.Groups.System[i_system].Services[i_system_services].Netconf = make([]xml_System_Services_Netconf, len(var_system_services_netconf))
        
		for i_system_services_netconf, v_system_services_netconf := range var_system_services_netconf {
            var var_system_services_netconf_ssh []System_Services_Netconf_Ssh_Model
            resp.Diagnostics.Append(v_system_services_netconf.Ssh.ElementsAs(ctx, &var_system_services_netconf_ssh, false)...)
            if resp.Diagnostics.HasError() {
                return
            }
	    config.Groups.System[i_system].Services[i_system_services].Netconf[i_system_services_netconf].Ssh = make([]xml_System_Services_Netconf_Ssh, len(var_system_services_netconf_ssh))
        
        }
            var var_system_services_rest []System_Services_Rest_Model
            resp.Diagnostics.Append(v_system_services.Rest.ElementsAs(ctx, &var_system_services_rest, false)...)
            if resp.Diagnostics.HasError() {
                return
            }
	    config.Groups.System[i_system].Services[i_system_services].Rest = make([]xml_System_Services_Rest, len(var_system_services_rest))
        
		for i_system_services_rest, v_system_services_rest := range var_system_services_rest {
            var var_system_services_rest_http []System_Services_Rest_Http_Model
            resp.Diagnostics.Append(v_system_services_rest.Http.ElementsAs(ctx, &var_system_services_rest_http, false)...)
            if resp.Diagnostics.HasError() {
                return
            }
	    config.Groups.System[i_system].Services[i_system_services].Rest[i_system_services_rest].Http = make([]xml_System_Services_Rest_Http, len(var_system_services_rest_http))
        
		for i_system_services_rest_http, v_system_services_rest_http := range var_system_services_rest_http {
            config.Groups.System[i_system].Services[i_system_services].Rest[i_system_services_rest].Http[i_system_services_rest_http].Port = v_system_services_rest_http.Port.ValueStringPointer()
        }
            config.Groups.System[i_system].Services[i_system_services].Rest[i_system_services_rest].Enable_explorer = v_system_services_rest.Enable_explorer.ValueStringPointer()
        }
        }
        var var_system_syslog []System_Syslog_Model
        resp.Diagnostics.Append(v_system.Syslog.ElementsAs(ctx, &var_system_syslog, false)...)
        if resp.Diagnostics.HasError() {
            return
        }
	    config.Groups.System[i_system].Syslog = make([]xml_System_Syslog, len(var_system_syslog))
        
		for i_system_syslog, v_system_syslog := range var_system_syslog {
            var var_system_syslog_user []System_Syslog_User_Model
            resp.Diagnostics.Append(v_system_syslog.User.ElementsAs(ctx, &var_system_syslog_user, false)...)
            if resp.Diagnostics.HasError() {
                return
            }
	    config.Groups.System[i_system].Syslog[i_system_syslog].User = make([]xml_System_Syslog_User, len(var_system_syslog_user))
        
		for i_system_syslog_user, v_system_syslog_user := range var_system_syslog_user {
            config.Groups.System[i_system].Syslog[i_system_syslog].User[i_system_syslog_user].Name = v_system_syslog_user.Name.ValueStringPointer()
            var var_system_syslog_user_contents []System_Syslog_User_Contents_Model
            resp.Diagnostics.Append(v_system_syslog_user.Contents.ElementsAs(ctx, &var_system_syslog_user_contents, false)...)
            if resp.Diagnostics.HasError() {
                return
            }
	    config.Groups.System[i_system].Syslog[i_system_syslog].User[i_system_syslog_user].Contents = make([]xml_System_Syslog_User_Contents, len(var_system_syslog_user_contents))
        
		for i_system_syslog_user_contents, v_system_syslog_user_contents := range var_system_syslog_user_contents {
            config.Groups.System[i_system].Syslog[i_system_syslog].User[i_system_syslog_user].Contents[i_system_syslog_user_contents].Name = v_system_syslog_user_contents.Name.ValueStringPointer()
            config.Groups.System[i_system].Syslog[i_system_syslog].User[i_system_syslog_user].Contents[i_system_syslog_user_contents].Emergency = v_system_syslog_user_contents.Emergency.ValueStringPointer()
        }
        }
            var var_system_syslog_file []System_Syslog_File_Model
            resp.Diagnostics.Append(v_system_syslog.File.ElementsAs(ctx, &var_system_syslog_file, false)...)
            if resp.Diagnostics.HasError() {
                return
            }
	    config.Groups.System[i_system].Syslog[i_system_syslog].File = make([]xml_System_Syslog_File, len(var_system_syslog_file))
        
		for i_system_syslog_file, v_system_syslog_file := range var_system_syslog_file {
            config.Groups.System[i_system].Syslog[i_system_syslog].File[i_system_syslog_file].Name = v_system_syslog_file.Name.ValueStringPointer()
            var var_system_syslog_file_contents []System_Syslog_File_Contents_Model
            resp.Diagnostics.Append(v_system_syslog_file.Contents.ElementsAs(ctx, &var_system_syslog_file_contents, false)...)
            if resp.Diagnostics.HasError() {
                return
            }
	    config.Groups.System[i_system].Syslog[i_system_syslog].File[i_system_syslog_file].Contents = make([]xml_System_Syslog_File_Contents, len(var_system_syslog_file_contents))
        
		for i_system_syslog_file_contents, v_system_syslog_file_contents := range var_system_syslog_file_contents {
            config.Groups.System[i_system].Syslog[i_system_syslog].File[i_system_syslog_file].Contents[i_system_syslog_file_contents].Name = v_system_syslog_file_contents.Name.ValueStringPointer()
            config.Groups.System[i_system].Syslog[i_system_syslog].File[i_system_syslog_file].Contents[i_system_syslog_file_contents].Any = v_system_syslog_file_contents.Any.ValueStringPointer()
            config.Groups.System[i_system].Syslog[i_system_syslog].File[i_system_syslog_file].Contents[i_system_syslog_file_contents].Info = v_system_syslog_file_contents.Info.ValueStringPointer()
        }
        }
        }
        var var_system_license []System_License_Model
        resp.Diagnostics.Append(v_system.License.ElementsAs(ctx, &var_system_license, false)...)
        if resp.Diagnostics.HasError() {
            return
        }
	    config.Groups.System[i_system].License = make([]xml_System_License, len(var_system_license))
        
		for i_system_license, v_system_license := range var_system_license {
            var var_system_license_autoupdate []System_License_Autoupdate_Model
            resp.Diagnostics.Append(v_system_license.Autoupdate.ElementsAs(ctx, &var_system_license_autoupdate, false)...)
            if resp.Diagnostics.HasError() {
                return
            }
	    config.Groups.System[i_system].License[i_system_license].Autoupdate = make([]xml_System_License_Autoupdate, len(var_system_license_autoupdate))
        
		for i_system_license_autoupdate, v_system_license_autoupdate := range var_system_license_autoupdate {
            var var_system_license_autoupdate_url []System_License_Autoupdate_Url_Model
            resp.Diagnostics.Append(v_system_license_autoupdate.Url.ElementsAs(ctx, &var_system_license_autoupdate_url, false)...)
            if resp.Diagnostics.HasError() {
                return
            }
	    config.Groups.System[i_system].License[i_system_license].Autoupdate[i_system_license_autoupdate].Url = make([]xml_System_License_Autoupdate_Url, len(var_system_license_autoupdate_url))
        
		for i_system_license_autoupdate_url, v_system_license_autoupdate_url := range var_system_license_autoupdate_url {
            config.Groups.System[i_system].License[i_system_license].Autoupdate[i_system_license_autoupdate].Url[i_system_license_autoupdate_url].Name = v_system_license_autoupdate_url.Name.ValueStringPointer()
        }
        }
        }
    }
	
	err := r.client.SendTransaction(plan.ResourceName.ValueString(), config, false)
	if err != nil {
		resp.Diagnostics.AddError("Failed while adding group", err.Error())
		return
	}
	commit_err := r.client.SendCommit()
	if commit_err != nil {
		resp.Diagnostics.AddError("Failed while committing apply-group", commit_err.Error())
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}




func (r *resource_Apply_Groups) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
    
    var state Groups_Model
    resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
    if resp.Diagnostics.HasError() {
        return
    }

    var config xml_Configuration
    err := r.client.MarshalGroup(state.ResourceName.ValueString(), &config)
    if err != nil {
        resp.Diagnostics.AddError("Failed to read group", err.Error())
        return
    }
    // Initialize chassis as Null; only materialize when we have elements
    state.Chassis =
        types.ListNull(types.ObjectType{AttrTypes: Chassis_Model{}.AttrTypes()})

    // Build list from device
    chassis_List := make([]Chassis_Model,
        len(config.Groups.Chassis))

    for i_chassis, v_chassis := range config.Groups.Chassis {
        var chassis_model Chassis_Model
            
        // Build aggregated-devices list
        chassis_aggregated_devices_List := make([]Chassis_Aggregated_devices_Model, len(v_chassis.Aggregated_devices))

        
		for i_chassis_aggregated_devices, v_chassis_aggregated_devices := range v_chassis.Aggregated_devices {
            var chassis_aggregated_devices_model Chassis_Aggregated_devices_Model

            chassis_aggregated_devices_List[i_chassis_aggregated_devices] =
                chassis_aggregated_devices_model
                
        // Build ethernet list
        chassis_aggregated_devices_ethernet_List := make([]Chassis_Aggregated_devices_Ethernet_Model, len(v_chassis_aggregated_devices.Ethernet))

        
		for i_chassis_aggregated_devices_ethernet, v_chassis_aggregated_devices_ethernet := range v_chassis_aggregated_devices.Ethernet {
            var chassis_aggregated_devices_ethernet_model Chassis_Aggregated_devices_Ethernet_Model
            // leaf -> keep pointer semantics
            chassis_aggregated_devices_ethernet_model.Device_count =
                types.StringPointerValue(v_chassis_aggregated_devices_ethernet.Device_count)

            chassis_aggregated_devices_ethernet_List[i_chassis_aggregated_devices_ethernet] =
                chassis_aggregated_devices_ethernet_model
        }

        // Write ethernet field as Null when empty, else concrete list
        if len(chassis_aggregated_devices_ethernet_List) == 0 {
            chassis_aggregated_devices_model.Ethernet =
                types.ListNull(types.ObjectType{AttrTypes: Chassis_Aggregated_devices_Ethernet_Model{}.AttrTypes()})
        } else {
            chassis_aggregated_devices_model.Ethernet, _ =
                types.ListValueFrom(ctx,
                    types.ObjectType{AttrTypes: Chassis_Aggregated_devices_Ethernet_Model{}.AttrTypes()},
                    chassis_aggregated_devices_ethernet_List,
                )
        }
        chassis_aggregated_devices_List[i_chassis_aggregated_devices] = chassis_aggregated_devices_model
        }

        // Write aggregated-devices field as Null when empty, else concrete list
        if len(chassis_aggregated_devices_List) == 0 {
            chassis_model.Aggregated_devices =
                types.ListNull(types.ObjectType{AttrTypes: Chassis_Aggregated_devices_Model{}.AttrTypes()})
        } else {
            chassis_model.Aggregated_devices, _ =
                types.ListValueFrom(ctx,
                    types.ObjectType{AttrTypes: Chassis_Aggregated_devices_Model{}.AttrTypes()},
                    chassis_aggregated_devices_List,
                )
        }
        chassis_List[i_chassis] = chassis_model

        chassis_List[i_chassis] = chassis_model
    }

    // Write parent list as Null when empty
    if len(chassis_List) == 0 {
        state.Chassis =
            types.ListNull(types.ObjectType{AttrTypes: Chassis_Model{}.AttrTypes()})
    } else {
        state.Chassis, _ =
            types.ListValueFrom(ctx,
                types.ObjectType{AttrTypes: Chassis_Model{}.AttrTypes()},
                chassis_List,
            )
    }
    // Initialize interfaces as Null; only materialize when we have elements
    state.Interfaces =
        types.ListNull(types.ObjectType{AttrTypes: Interfaces_Model{}.AttrTypes()})

    // Build list from device
    interfaces_List := make([]Interfaces_Model,
        len(config.Groups.Interfaces))

    for i_interfaces, v_interfaces := range config.Groups.Interfaces {
        var interfaces_model Interfaces_Model
            
        // Build interface list
        interfaces_interface_List := make([]Interfaces_Interface_Model, len(v_interfaces.Interface))

        
		for i_interfaces_interface, v_interfaces_interface := range v_interfaces.Interface {
            var interfaces_interface_model Interfaces_Interface_Model
            // leaf -> keep pointer semantics
            interfaces_interface_model.Name =
                types.StringPointerValue(v_interfaces_interface.Name)

            interfaces_interface_List[i_interfaces_interface] =
                interfaces_interface_model
            // leaf -> keep pointer semantics
            interfaces_interface_model.Vlan_tagging =
                types.StringPointerValue(v_interfaces_interface.Vlan_tagging)

            interfaces_interface_List[i_interfaces_interface] =
                interfaces_interface_model

            interfaces_interface_List[i_interfaces_interface] =
                interfaces_interface_model
                
        // Build unit list
        interfaces_interface_unit_List := make([]Interfaces_Interface_Unit_Model, len(v_interfaces_interface.Unit))

        
		for i_interfaces_interface_unit, v_interfaces_interface_unit := range v_interfaces_interface.Unit {
            var interfaces_interface_unit_model Interfaces_Interface_Unit_Model
            // leaf -> keep pointer semantics
            interfaces_interface_unit_model.Name =
                types.StringPointerValue(v_interfaces_interface_unit.Name)

            interfaces_interface_unit_List[i_interfaces_interface_unit] =
                interfaces_interface_unit_model
            // leaf -> keep pointer semantics
            interfaces_interface_unit_model.Description =
                types.StringPointerValue(v_interfaces_interface_unit.Description)

            interfaces_interface_unit_List[i_interfaces_interface_unit] =
                interfaces_interface_unit_model
            // leaf -> keep pointer semantics
            interfaces_interface_unit_model.Vlan_id =
                types.StringPointerValue(v_interfaces_interface_unit.Vlan_id)

            interfaces_interface_unit_List[i_interfaces_interface_unit] =
                interfaces_interface_unit_model

            interfaces_interface_unit_List[i_interfaces_interface_unit] =
                interfaces_interface_unit_model
                
        // Build family list
        interfaces_interface_unit_family_List := make([]Interfaces_Interface_Unit_Family_Model, len(v_interfaces_interface_unit.Family))

        
		for i_interfaces_interface_unit_family, v_interfaces_interface_unit_family := range v_interfaces_interface_unit.Family {
            var interfaces_interface_unit_family_model Interfaces_Interface_Unit_Family_Model

            interfaces_interface_unit_family_List[i_interfaces_interface_unit_family] =
                interfaces_interface_unit_family_model
                
        // Build inet list
        interfaces_interface_unit_family_inet_List := make([]Interfaces_Interface_Unit_Family_Inet_Model, len(v_interfaces_interface_unit_family.Inet))

        
		for i_interfaces_interface_unit_family_inet, v_interfaces_interface_unit_family_inet := range v_interfaces_interface_unit_family.Inet {
            var interfaces_interface_unit_family_inet_model Interfaces_Interface_Unit_Family_Inet_Model

            interfaces_interface_unit_family_inet_List[i_interfaces_interface_unit_family_inet] =
                interfaces_interface_unit_family_inet_model
                
        // Build address list
        interfaces_interface_unit_family_inet_address_List := make([]Interfaces_Interface_Unit_Family_Inet_Address_Model, len(v_interfaces_interface_unit_family_inet.Address))

        
		for i_interfaces_interface_unit_family_inet_address, v_interfaces_interface_unit_family_inet_address := range v_interfaces_interface_unit_family_inet.Address {
            var interfaces_interface_unit_family_inet_address_model Interfaces_Interface_Unit_Family_Inet_Address_Model
            // leaf -> keep pointer semantics
            interfaces_interface_unit_family_inet_address_model.Name =
                types.StringPointerValue(v_interfaces_interface_unit_family_inet_address.Name)

            interfaces_interface_unit_family_inet_address_List[i_interfaces_interface_unit_family_inet_address] =
                interfaces_interface_unit_family_inet_address_model
        }

        // Write address field as Null when empty, else concrete list
        if len(interfaces_interface_unit_family_inet_address_List) == 0 {
            interfaces_interface_unit_family_inet_model.Address =
                types.ListNull(types.ObjectType{AttrTypes: Interfaces_Interface_Unit_Family_Inet_Address_Model{}.AttrTypes()})
        } else {
            interfaces_interface_unit_family_inet_model.Address, _ =
                types.ListValueFrom(ctx,
                    types.ObjectType{AttrTypes: Interfaces_Interface_Unit_Family_Inet_Address_Model{}.AttrTypes()},
                    interfaces_interface_unit_family_inet_address_List,
                )
        }
        interfaces_interface_unit_family_inet_List[i_interfaces_interface_unit_family_inet] = interfaces_interface_unit_family_inet_model
        }

        // Write inet field as Null when empty, else concrete list
        if len(interfaces_interface_unit_family_inet_List) == 0 {
            interfaces_interface_unit_family_model.Inet =
                types.ListNull(types.ObjectType{AttrTypes: Interfaces_Interface_Unit_Family_Inet_Model{}.AttrTypes()})
        } else {
            interfaces_interface_unit_family_model.Inet, _ =
                types.ListValueFrom(ctx,
                    types.ObjectType{AttrTypes: Interfaces_Interface_Unit_Family_Inet_Model{}.AttrTypes()},
                    interfaces_interface_unit_family_inet_List,
                )
        }
        interfaces_interface_unit_family_List[i_interfaces_interface_unit_family] = interfaces_interface_unit_family_model
        }

        // Write family field as Null when empty, else concrete list
        if len(interfaces_interface_unit_family_List) == 0 {
            interfaces_interface_unit_model.Family =
                types.ListNull(types.ObjectType{AttrTypes: Interfaces_Interface_Unit_Family_Model{}.AttrTypes()})
        } else {
            interfaces_interface_unit_model.Family, _ =
                types.ListValueFrom(ctx,
                    types.ObjectType{AttrTypes: Interfaces_Interface_Unit_Family_Model{}.AttrTypes()},
                    interfaces_interface_unit_family_List,
                )
        }
        interfaces_interface_unit_List[i_interfaces_interface_unit] = interfaces_interface_unit_model
        }

        // Write unit field as Null when empty, else concrete list
        if len(interfaces_interface_unit_List) == 0 {
            interfaces_interface_model.Unit =
                types.ListNull(types.ObjectType{AttrTypes: Interfaces_Interface_Unit_Model{}.AttrTypes()})
        } else {
            interfaces_interface_model.Unit, _ =
                types.ListValueFrom(ctx,
                    types.ObjectType{AttrTypes: Interfaces_Interface_Unit_Model{}.AttrTypes()},
                    interfaces_interface_unit_List,
                )
        }
        interfaces_interface_List[i_interfaces_interface] = interfaces_interface_model
        }

        // Write interface field as Null when empty, else concrete list
        if len(interfaces_interface_List) == 0 {
            interfaces_model.Interface =
                types.ListNull(types.ObjectType{AttrTypes: Interfaces_Interface_Model{}.AttrTypes()})
        } else {
            interfaces_model.Interface, _ =
                types.ListValueFrom(ctx,
                    types.ObjectType{AttrTypes: Interfaces_Interface_Model{}.AttrTypes()},
                    interfaces_interface_List,
                )
        }
        interfaces_List[i_interfaces] = interfaces_model

        interfaces_List[i_interfaces] = interfaces_model
    }

    // Write parent list as Null when empty
    if len(interfaces_List) == 0 {
        state.Interfaces =
            types.ListNull(types.ObjectType{AttrTypes: Interfaces_Model{}.AttrTypes()})
    } else {
        state.Interfaces, _ =
            types.ListValueFrom(ctx,
                types.ObjectType{AttrTypes: Interfaces_Model{}.AttrTypes()},
                interfaces_List,
            )
    }
    // Initialize policy-options as Null; only materialize when we have elements
    state.Policy_options =
        types.ListNull(types.ObjectType{AttrTypes: Policy_options_Model{}.AttrTypes()})

    // Build list from device
    policy_options_List := make([]Policy_options_Model,
        len(config.Groups.Policy_options))

    for i_policy_options, v_policy_options := range config.Groups.Policy_options {
        var policy_options_model Policy_options_Model
            
        // Build policy-statement list
        policy_options_policy_statement_List := make([]Policy_options_Policy_statement_Model, len(v_policy_options.Policy_statement))

        
		for i_policy_options_policy_statement, v_policy_options_policy_statement := range v_policy_options.Policy_statement {
            var policy_options_policy_statement_model Policy_options_Policy_statement_Model
            // leaf -> keep pointer semantics
            policy_options_policy_statement_model.Name =
                types.StringPointerValue(v_policy_options_policy_statement.Name)

            policy_options_policy_statement_List[i_policy_options_policy_statement] =
                policy_options_policy_statement_model

            policy_options_policy_statement_List[i_policy_options_policy_statement] =
                policy_options_policy_statement_model
                
        // Build term list
        policy_options_policy_statement_term_List := make([]Policy_options_Policy_statement_Term_Model, len(v_policy_options_policy_statement.Term))

        
		for i_policy_options_policy_statement_term, v_policy_options_policy_statement_term := range v_policy_options_policy_statement.Term {
            var policy_options_policy_statement_term_model Policy_options_Policy_statement_Term_Model
            // leaf -> keep pointer semantics
            policy_options_policy_statement_term_model.Name =
                types.StringPointerValue(v_policy_options_policy_statement_term.Name)

            policy_options_policy_statement_term_List[i_policy_options_policy_statement_term] =
                policy_options_policy_statement_term_model

            policy_options_policy_statement_term_List[i_policy_options_policy_statement_term] =
                policy_options_policy_statement_term_model
                
        // Build from list
        policy_options_policy_statement_term_from_List := make([]Policy_options_Policy_statement_Term_From_Model, len(v_policy_options_policy_statement_term.From))

        
		for i_policy_options_policy_statement_term_from, v_policy_options_policy_statement_term_from := range v_policy_options_policy_statement_term.From {
            var policy_options_policy_statement_term_from_model Policy_options_Policy_statement_Term_From_Model

            policy_options_policy_statement_term_from_List[i_policy_options_policy_statement_term_from] =
                policy_options_policy_statement_term_from_model
                
        // Build route-filter list
        policy_options_policy_statement_term_from_route_filter_List := make([]Policy_options_Policy_statement_Term_From_Route_filter_Model, len(v_policy_options_policy_statement_term_from.Route_filter))

        
		for i_policy_options_policy_statement_term_from_route_filter, v_policy_options_policy_statement_term_from_route_filter := range v_policy_options_policy_statement_term_from.Route_filter {
            var policy_options_policy_statement_term_from_route_filter_model Policy_options_Policy_statement_Term_From_Route_filter_Model
            // leaf -> keep pointer semantics
            policy_options_policy_statement_term_from_route_filter_model.Address =
                types.StringPointerValue(v_policy_options_policy_statement_term_from_route_filter.Address)

            policy_options_policy_statement_term_from_route_filter_List[i_policy_options_policy_statement_term_from_route_filter] =
                policy_options_policy_statement_term_from_route_filter_model
            // leaf -> keep pointer semantics
            policy_options_policy_statement_term_from_route_filter_model.Exact =
                types.StringPointerValue(v_policy_options_policy_statement_term_from_route_filter.Exact)

            policy_options_policy_statement_term_from_route_filter_List[i_policy_options_policy_statement_term_from_route_filter] =
                policy_options_policy_statement_term_from_route_filter_model
            // leaf -> keep pointer semantics
            policy_options_policy_statement_term_from_route_filter_model.Accept =
                types.StringPointerValue(v_policy_options_policy_statement_term_from_route_filter.Accept)

            policy_options_policy_statement_term_from_route_filter_List[i_policy_options_policy_statement_term_from_route_filter] =
                policy_options_policy_statement_term_from_route_filter_model
        }

        // Write route-filter field as Null when empty, else concrete list
        if len(policy_options_policy_statement_term_from_route_filter_List) == 0 {
            policy_options_policy_statement_term_from_model.Route_filter =
                types.ListNull(types.ObjectType{AttrTypes: Policy_options_Policy_statement_Term_From_Route_filter_Model{}.AttrTypes()})
        } else {
            policy_options_policy_statement_term_from_model.Route_filter, _ =
                types.ListValueFrom(ctx,
                    types.ObjectType{AttrTypes: Policy_options_Policy_statement_Term_From_Route_filter_Model{}.AttrTypes()},
                    policy_options_policy_statement_term_from_route_filter_List,
                )
        }
        policy_options_policy_statement_term_from_List[i_policy_options_policy_statement_term_from] = policy_options_policy_statement_term_from_model
        }

        // Write from field as Null when empty, else concrete list
        if len(policy_options_policy_statement_term_from_List) == 0 {
            policy_options_policy_statement_term_model.From =
                types.ListNull(types.ObjectType{AttrTypes: Policy_options_Policy_statement_Term_From_Model{}.AttrTypes()})
        } else {
            policy_options_policy_statement_term_model.From, _ =
                types.ListValueFrom(ctx,
                    types.ObjectType{AttrTypes: Policy_options_Policy_statement_Term_From_Model{}.AttrTypes()},
                    policy_options_policy_statement_term_from_List,
                )
        }
        policy_options_policy_statement_term_List[i_policy_options_policy_statement_term] = policy_options_policy_statement_term_model
        }

        // Write term field as Null when empty, else concrete list
        if len(policy_options_policy_statement_term_List) == 0 {
            policy_options_policy_statement_model.Term =
                types.ListNull(types.ObjectType{AttrTypes: Policy_options_Policy_statement_Term_Model{}.AttrTypes()})
        } else {
            policy_options_policy_statement_model.Term, _ =
                types.ListValueFrom(ctx,
                    types.ObjectType{AttrTypes: Policy_options_Policy_statement_Term_Model{}.AttrTypes()},
                    policy_options_policy_statement_term_List,
                )
        }
        policy_options_policy_statement_List[i_policy_options_policy_statement] = policy_options_policy_statement_model
        }

        // Write policy-statement field as Null when empty, else concrete list
        if len(policy_options_policy_statement_List) == 0 {
            policy_options_model.Policy_statement =
                types.ListNull(types.ObjectType{AttrTypes: Policy_options_Policy_statement_Model{}.AttrTypes()})
        } else {
            policy_options_model.Policy_statement, _ =
                types.ListValueFrom(ctx,
                    types.ObjectType{AttrTypes: Policy_options_Policy_statement_Model{}.AttrTypes()},
                    policy_options_policy_statement_List,
                )
        }
        policy_options_List[i_policy_options] = policy_options_model

        policy_options_List[i_policy_options] = policy_options_model
    }

    // Write parent list as Null when empty
    if len(policy_options_List) == 0 {
        state.Policy_options =
            types.ListNull(types.ObjectType{AttrTypes: Policy_options_Model{}.AttrTypes()})
    } else {
        state.Policy_options, _ =
            types.ListValueFrom(ctx,
                types.ObjectType{AttrTypes: Policy_options_Model{}.AttrTypes()},
                policy_options_List,
            )
    }
    // Initialize protocols as Null; only materialize when we have elements
    state.Protocols =
        types.ListNull(types.ObjectType{AttrTypes: Protocols_Model{}.AttrTypes()})

    // Build list from device
    protocols_List := make([]Protocols_Model,
        len(config.Groups.Protocols))

    for i_protocols, v_protocols := range config.Groups.Protocols {
        var protocols_model Protocols_Model
            
        // Build lldp list
        protocols_lldp_List := make([]Protocols_Lldp_Model, len(v_protocols.Lldp))

        
		for i_protocols_lldp, v_protocols_lldp := range v_protocols.Lldp {
            var protocols_lldp_model Protocols_Lldp_Model

            protocols_lldp_List[i_protocols_lldp] =
                protocols_lldp_model
                
        // Build interface list
        protocols_lldp_interface_List := make([]Protocols_Lldp_Interface_Model, len(v_protocols_lldp.Interface))

        
		for i_protocols_lldp_interface, v_protocols_lldp_interface := range v_protocols_lldp.Interface {
            var protocols_lldp_interface_model Protocols_Lldp_Interface_Model
            // leaf -> keep pointer semantics
            protocols_lldp_interface_model.Name =
                types.StringPointerValue(v_protocols_lldp_interface.Name)

            protocols_lldp_interface_List[i_protocols_lldp_interface] =
                protocols_lldp_interface_model
        }

        // Write interface field as Null when empty, else concrete list
        if len(protocols_lldp_interface_List) == 0 {
            protocols_lldp_model.Interface =
                types.ListNull(types.ObjectType{AttrTypes: Protocols_Lldp_Interface_Model{}.AttrTypes()})
        } else {
            protocols_lldp_model.Interface, _ =
                types.ListValueFrom(ctx,
                    types.ObjectType{AttrTypes: Protocols_Lldp_Interface_Model{}.AttrTypes()},
                    protocols_lldp_interface_List,
                )
        }
        protocols_lldp_List[i_protocols_lldp] = protocols_lldp_model
        }

        // Write lldp field as Null when empty, else concrete list
        if len(protocols_lldp_List) == 0 {
            protocols_model.Lldp =
                types.ListNull(types.ObjectType{AttrTypes: Protocols_Lldp_Model{}.AttrTypes()})
        } else {
            protocols_model.Lldp, _ =
                types.ListValueFrom(ctx,
                    types.ObjectType{AttrTypes: Protocols_Lldp_Model{}.AttrTypes()},
                    protocols_lldp_List,
                )
        }
        protocols_List[i_protocols] = protocols_model

        protocols_List[i_protocols] = protocols_model
    }

    // Write parent list as Null when empty
    if len(protocols_List) == 0 {
        state.Protocols =
            types.ListNull(types.ObjectType{AttrTypes: Protocols_Model{}.AttrTypes()})
    } else {
        state.Protocols, _ =
            types.ListValueFrom(ctx,
                types.ObjectType{AttrTypes: Protocols_Model{}.AttrTypes()},
                protocols_List,
            )
    }
    // Initialize routing-instances as Null; only materialize when we have elements
    state.Routing_instances =
        types.ListNull(types.ObjectType{AttrTypes: Routing_instances_Model{}.AttrTypes()})

    // Build list from device
    routing_instances_List := make([]Routing_instances_Model,
        len(config.Groups.Routing_instances))

    for i_routing_instances, v_routing_instances := range config.Groups.Routing_instances {
        var routing_instances_model Routing_instances_Model
            
        // Build instance list
        routing_instances_instance_List := make([]Routing_instances_Instance_Model, len(v_routing_instances.Instance))

        
		for i_routing_instances_instance, v_routing_instances_instance := range v_routing_instances.Instance {
            var routing_instances_instance_model Routing_instances_Instance_Model
            // leaf -> keep pointer semantics
            routing_instances_instance_model.Name =
                types.StringPointerValue(v_routing_instances_instance.Name)

            routing_instances_instance_List[i_routing_instances_instance] =
                routing_instances_instance_model
            // leaf -> keep pointer semantics
            routing_instances_instance_model.Instance_type =
                types.StringPointerValue(v_routing_instances_instance.Instance_type)

            routing_instances_instance_List[i_routing_instances_instance] =
                routing_instances_instance_model

            routing_instances_instance_List[i_routing_instances_instance] =
                routing_instances_instance_model
                
        // Build interface list
        routing_instances_instance_interface_List := make([]Routing_instances_Instance_Interface_Model, len(v_routing_instances_instance.Interface))

        
		for i_routing_instances_instance_interface, v_routing_instances_instance_interface := range v_routing_instances_instance.Interface {
            var routing_instances_instance_interface_model Routing_instances_Instance_Interface_Model
            // leaf -> keep pointer semantics
            routing_instances_instance_interface_model.Name =
                types.StringPointerValue(v_routing_instances_instance_interface.Name)

            routing_instances_instance_interface_List[i_routing_instances_instance_interface] =
                routing_instances_instance_interface_model
        }

        // Write interface field as Null when empty, else concrete list
        if len(routing_instances_instance_interface_List) == 0 {
            routing_instances_instance_model.Interface =
                types.ListNull(types.ObjectType{AttrTypes: Routing_instances_Instance_Interface_Model{}.AttrTypes()})
        } else {
            routing_instances_instance_model.Interface, _ =
                types.ListValueFrom(ctx,
                    types.ObjectType{AttrTypes: Routing_instances_Instance_Interface_Model{}.AttrTypes()},
                    routing_instances_instance_interface_List,
                )
        }
        routing_instances_instance_List[i_routing_instances_instance] = routing_instances_instance_model

            routing_instances_instance_List[i_routing_instances_instance] =
                routing_instances_instance_model
                
        // Build routing-options list
        routing_instances_instance_routing_options_List := make([]Routing_instances_Instance_Routing_options_Model, len(v_routing_instances_instance.Routing_options))

        
		for i_routing_instances_instance_routing_options, v_routing_instances_instance_routing_options := range v_routing_instances_instance.Routing_options {
            var routing_instances_instance_routing_options_model Routing_instances_Instance_Routing_options_Model

            routing_instances_instance_routing_options_List[i_routing_instances_instance_routing_options] =
                routing_instances_instance_routing_options_model
                
        // Build static list
        routing_instances_instance_routing_options_static_List := make([]Routing_instances_Instance_Routing_options_Static_Model, len(v_routing_instances_instance_routing_options.Static))

        
		for i_routing_instances_instance_routing_options_static, v_routing_instances_instance_routing_options_static := range v_routing_instances_instance_routing_options.Static {
            var routing_instances_instance_routing_options_static_model Routing_instances_Instance_Routing_options_Static_Model

            routing_instances_instance_routing_options_static_List[i_routing_instances_instance_routing_options_static] =
                routing_instances_instance_routing_options_static_model
                
        // Build route list
        routing_instances_instance_routing_options_static_route_List := make([]Routing_instances_Instance_Routing_options_Static_Route_Model, len(v_routing_instances_instance_routing_options_static.Route))

        
		for i_routing_instances_instance_routing_options_static_route, v_routing_instances_instance_routing_options_static_route := range v_routing_instances_instance_routing_options_static.Route {
            var routing_instances_instance_routing_options_static_route_model Routing_instances_Instance_Routing_options_Static_Route_Model
            // leaf -> keep pointer semantics
            routing_instances_instance_routing_options_static_route_model.Name =
                types.StringPointerValue(v_routing_instances_instance_routing_options_static_route.Name)

            routing_instances_instance_routing_options_static_route_List[i_routing_instances_instance_routing_options_static_route] =
                routing_instances_instance_routing_options_static_route_model
            // leaf -> keep pointer semantics
            routing_instances_instance_routing_options_static_route_model.Discard =
                types.StringPointerValue(v_routing_instances_instance_routing_options_static_route.Discard)

            routing_instances_instance_routing_options_static_route_List[i_routing_instances_instance_routing_options_static_route] =
                routing_instances_instance_routing_options_static_route_model
        }

        // Write route field as Null when empty, else concrete list
        if len(routing_instances_instance_routing_options_static_route_List) == 0 {
            routing_instances_instance_routing_options_static_model.Route =
                types.ListNull(types.ObjectType{AttrTypes: Routing_instances_Instance_Routing_options_Static_Route_Model{}.AttrTypes()})
        } else {
            routing_instances_instance_routing_options_static_model.Route, _ =
                types.ListValueFrom(ctx,
                    types.ObjectType{AttrTypes: Routing_instances_Instance_Routing_options_Static_Route_Model{}.AttrTypes()},
                    routing_instances_instance_routing_options_static_route_List,
                )
        }
        routing_instances_instance_routing_options_static_List[i_routing_instances_instance_routing_options_static] = routing_instances_instance_routing_options_static_model
        }

        // Write static field as Null when empty, else concrete list
        if len(routing_instances_instance_routing_options_static_List) == 0 {
            routing_instances_instance_routing_options_model.Static =
                types.ListNull(types.ObjectType{AttrTypes: Routing_instances_Instance_Routing_options_Static_Model{}.AttrTypes()})
        } else {
            routing_instances_instance_routing_options_model.Static, _ =
                types.ListValueFrom(ctx,
                    types.ObjectType{AttrTypes: Routing_instances_Instance_Routing_options_Static_Model{}.AttrTypes()},
                    routing_instances_instance_routing_options_static_List,
                )
        }
        routing_instances_instance_routing_options_List[i_routing_instances_instance_routing_options] = routing_instances_instance_routing_options_model
        }

        // Write routing-options field as Null when empty, else concrete list
        if len(routing_instances_instance_routing_options_List) == 0 {
            routing_instances_instance_model.Routing_options =
                types.ListNull(types.ObjectType{AttrTypes: Routing_instances_Instance_Routing_options_Model{}.AttrTypes()})
        } else {
            routing_instances_instance_model.Routing_options, _ =
                types.ListValueFrom(ctx,
                    types.ObjectType{AttrTypes: Routing_instances_Instance_Routing_options_Model{}.AttrTypes()},
                    routing_instances_instance_routing_options_List,
                )
        }
        routing_instances_instance_List[i_routing_instances_instance] = routing_instances_instance_model

            routing_instances_instance_List[i_routing_instances_instance] =
                routing_instances_instance_model
                
        // Build protocols list
        routing_instances_instance_protocols_List := make([]Routing_instances_Instance_Protocols_Model, len(v_routing_instances_instance.Protocols))

        
		for i_routing_instances_instance_protocols, v_routing_instances_instance_protocols := range v_routing_instances_instance.Protocols {
            var routing_instances_instance_protocols_model Routing_instances_Instance_Protocols_Model

            routing_instances_instance_protocols_List[i_routing_instances_instance_protocols] =
                routing_instances_instance_protocols_model
                
        // Build ospf list
        routing_instances_instance_protocols_ospf_List := make([]Routing_instances_Instance_Protocols_Ospf_Model, len(v_routing_instances_instance_protocols.Ospf))

        
		for i_routing_instances_instance_protocols_ospf, v_routing_instances_instance_protocols_ospf := range v_routing_instances_instance_protocols.Ospf {
            var routing_instances_instance_protocols_ospf_model Routing_instances_Instance_Protocols_Ospf_Model
            // leaf-list -> write Null when nil OR empty (avoid [] when absent)
            if v_routing_instances_instance_protocols_ospf.Export == nil ||
               len(v_routing_instances_instance_protocols_ospf.Export) == 0 {
                routing_instances_instance_protocols_ospf_model.Export =
                    types.ListNull(types.StringType)
            } else {
                src_routing_instances_instance_protocols_export :=
                    v_routing_instances_instance_protocols_ospf.Export
                vals_routing_instances_instance_protocols_export := make([]*string, len(src_routing_instances_instance_protocols_export))
                copy(vals_routing_instances_instance_protocols_export, src_routing_instances_instance_protocols_export)
                routing_instances_instance_protocols_ospf_model.Export, _ =
                    types.ListValueFrom(ctx, types.StringType, vals_routing_instances_instance_protocols_export)
            }

            routing_instances_instance_protocols_ospf_List[i_routing_instances_instance_protocols_ospf] =
                routing_instances_instance_protocols_ospf_model

            routing_instances_instance_protocols_ospf_List[i_routing_instances_instance_protocols_ospf] =
                routing_instances_instance_protocols_ospf_model
                
        // Build area list
        routing_instances_instance_protocols_ospf_area_List := make([]Routing_instances_Instance_Protocols_Ospf_Area_Model, len(v_routing_instances_instance_protocols_ospf.Area))

        
		for i_routing_instances_instance_protocols_ospf_area, v_routing_instances_instance_protocols_ospf_area := range v_routing_instances_instance_protocols_ospf.Area {
            var routing_instances_instance_protocols_ospf_area_model Routing_instances_Instance_Protocols_Ospf_Area_Model
            // leaf -> keep pointer semantics
            routing_instances_instance_protocols_ospf_area_model.Name =
                types.StringPointerValue(v_routing_instances_instance_protocols_ospf_area.Name)

            routing_instances_instance_protocols_ospf_area_List[i_routing_instances_instance_protocols_ospf_area] =
                routing_instances_instance_protocols_ospf_area_model

            routing_instances_instance_protocols_ospf_area_List[i_routing_instances_instance_protocols_ospf_area] =
                routing_instances_instance_protocols_ospf_area_model
                
        // Build interface list
        routing_instances_instance_protocols_ospf_area_interface_List := make([]Routing_instances_Instance_Protocols_Ospf_Area_Interface_Model, len(v_routing_instances_instance_protocols_ospf_area.Interface))

        
		for i_routing_instances_instance_protocols_ospf_area_interface, v_routing_instances_instance_protocols_ospf_area_interface := range v_routing_instances_instance_protocols_ospf_area.Interface {
            var routing_instances_instance_protocols_ospf_area_interface_model Routing_instances_Instance_Protocols_Ospf_Area_Interface_Model
            // leaf -> keep pointer semantics
            routing_instances_instance_protocols_ospf_area_interface_model.Name =
                types.StringPointerValue(v_routing_instances_instance_protocols_ospf_area_interface.Name)

            routing_instances_instance_protocols_ospf_area_interface_List[i_routing_instances_instance_protocols_ospf_area_interface] =
                routing_instances_instance_protocols_ospf_area_interface_model
        }

        // Write interface field as Null when empty, else concrete list
        if len(routing_instances_instance_protocols_ospf_area_interface_List) == 0 {
            routing_instances_instance_protocols_ospf_area_model.Interface =
                types.ListNull(types.ObjectType{AttrTypes: Routing_instances_Instance_Protocols_Ospf_Area_Interface_Model{}.AttrTypes()})
        } else {
            routing_instances_instance_protocols_ospf_area_model.Interface, _ =
                types.ListValueFrom(ctx,
                    types.ObjectType{AttrTypes: Routing_instances_Instance_Protocols_Ospf_Area_Interface_Model{}.AttrTypes()},
                    routing_instances_instance_protocols_ospf_area_interface_List,
                )
        }
        routing_instances_instance_protocols_ospf_area_List[i_routing_instances_instance_protocols_ospf_area] = routing_instances_instance_protocols_ospf_area_model
        }

        // Write area field as Null when empty, else concrete list
        if len(routing_instances_instance_protocols_ospf_area_List) == 0 {
            routing_instances_instance_protocols_ospf_model.Area =
                types.ListNull(types.ObjectType{AttrTypes: Routing_instances_Instance_Protocols_Ospf_Area_Model{}.AttrTypes()})
        } else {
            routing_instances_instance_protocols_ospf_model.Area, _ =
                types.ListValueFrom(ctx,
                    types.ObjectType{AttrTypes: Routing_instances_Instance_Protocols_Ospf_Area_Model{}.AttrTypes()},
                    routing_instances_instance_protocols_ospf_area_List,
                )
        }
        routing_instances_instance_protocols_ospf_List[i_routing_instances_instance_protocols_ospf] = routing_instances_instance_protocols_ospf_model
        }

        // Write ospf field as Null when empty, else concrete list
        if len(routing_instances_instance_protocols_ospf_List) == 0 {
            routing_instances_instance_protocols_model.Ospf =
                types.ListNull(types.ObjectType{AttrTypes: Routing_instances_Instance_Protocols_Ospf_Model{}.AttrTypes()})
        } else {
            routing_instances_instance_protocols_model.Ospf, _ =
                types.ListValueFrom(ctx,
                    types.ObjectType{AttrTypes: Routing_instances_Instance_Protocols_Ospf_Model{}.AttrTypes()},
                    routing_instances_instance_protocols_ospf_List,
                )
        }
        routing_instances_instance_protocols_List[i_routing_instances_instance_protocols] = routing_instances_instance_protocols_model
        }

        // Write protocols field as Null when empty, else concrete list
        if len(routing_instances_instance_protocols_List) == 0 {
            routing_instances_instance_model.Protocols =
                types.ListNull(types.ObjectType{AttrTypes: Routing_instances_Instance_Protocols_Model{}.AttrTypes()})
        } else {
            routing_instances_instance_model.Protocols, _ =
                types.ListValueFrom(ctx,
                    types.ObjectType{AttrTypes: Routing_instances_Instance_Protocols_Model{}.AttrTypes()},
                    routing_instances_instance_protocols_List,
                )
        }
        routing_instances_instance_List[i_routing_instances_instance] = routing_instances_instance_model
        }

        // Write instance field as Null when empty, else concrete list
        if len(routing_instances_instance_List) == 0 {
            routing_instances_model.Instance =
                types.ListNull(types.ObjectType{AttrTypes: Routing_instances_Instance_Model{}.AttrTypes()})
        } else {
            routing_instances_model.Instance, _ =
                types.ListValueFrom(ctx,
                    types.ObjectType{AttrTypes: Routing_instances_Instance_Model{}.AttrTypes()},
                    routing_instances_instance_List,
                )
        }
        routing_instances_List[i_routing_instances] = routing_instances_model

        routing_instances_List[i_routing_instances] = routing_instances_model
    }

    // Write parent list as Null when empty
    if len(routing_instances_List) == 0 {
        state.Routing_instances =
            types.ListNull(types.ObjectType{AttrTypes: Routing_instances_Model{}.AttrTypes()})
    } else {
        state.Routing_instances, _ =
            types.ListValueFrom(ctx,
                types.ObjectType{AttrTypes: Routing_instances_Model{}.AttrTypes()},
                routing_instances_List,
            )
    }
    // Initialize routing-options as Null; only materialize when we have elements
    state.Routing_options =
        types.ListNull(types.ObjectType{AttrTypes: Routing_options_Model{}.AttrTypes()})

    // Build list from device
    routing_options_List := make([]Routing_options_Model,
        len(config.Groups.Routing_options))

    for i_routing_options, v_routing_options := range config.Groups.Routing_options {
        var routing_options_model Routing_options_Model
            
        // Build static list
        routing_options_static_List := make([]Routing_options_Static_Model, len(v_routing_options.Static))

        
		for i_routing_options_static, v_routing_options_static := range v_routing_options.Static {
            var routing_options_static_model Routing_options_Static_Model

            routing_options_static_List[i_routing_options_static] =
                routing_options_static_model
                
        // Build route list
        routing_options_static_route_List := make([]Routing_options_Static_Route_Model, len(v_routing_options_static.Route))

        
		for i_routing_options_static_route, v_routing_options_static_route := range v_routing_options_static.Route {
            var routing_options_static_route_model Routing_options_Static_Route_Model
            // leaf -> keep pointer semantics
            routing_options_static_route_model.Name =
                types.StringPointerValue(v_routing_options_static_route.Name)

            routing_options_static_route_List[i_routing_options_static_route] =
                routing_options_static_route_model
            // leaf-list -> write Null when nil OR empty (avoid [] when absent)
            if v_routing_options_static_route.Next_hop == nil ||
               len(v_routing_options_static_route.Next_hop) == 0 {
                routing_options_static_route_model.Next_hop =
                    types.ListNull(types.StringType)
            } else {
                src_routing_options_static_next_hop :=
                    v_routing_options_static_route.Next_hop
                vals_routing_options_static_next_hop := make([]*string, len(src_routing_options_static_next_hop))
                copy(vals_routing_options_static_next_hop, src_routing_options_static_next_hop)
                routing_options_static_route_model.Next_hop, _ =
                    types.ListValueFrom(ctx, types.StringType, vals_routing_options_static_next_hop)
            }

            routing_options_static_route_List[i_routing_options_static_route] =
                routing_options_static_route_model
        }

        // Write route field as Null when empty, else concrete list
        if len(routing_options_static_route_List) == 0 {
            routing_options_static_model.Route =
                types.ListNull(types.ObjectType{AttrTypes: Routing_options_Static_Route_Model{}.AttrTypes()})
        } else {
            routing_options_static_model.Route, _ =
                types.ListValueFrom(ctx,
                    types.ObjectType{AttrTypes: Routing_options_Static_Route_Model{}.AttrTypes()},
                    routing_options_static_route_List,
                )
        }
        routing_options_static_List[i_routing_options_static] = routing_options_static_model
        }

        // Write static field as Null when empty, else concrete list
        if len(routing_options_static_List) == 0 {
            routing_options_model.Static =
                types.ListNull(types.ObjectType{AttrTypes: Routing_options_Static_Model{}.AttrTypes()})
        } else {
            routing_options_model.Static, _ =
                types.ListValueFrom(ctx,
                    types.ObjectType{AttrTypes: Routing_options_Static_Model{}.AttrTypes()},
                    routing_options_static_List,
                )
        }
        routing_options_List[i_routing_options] = routing_options_model

        routing_options_List[i_routing_options] = routing_options_model
    }

    // Write parent list as Null when empty
    if len(routing_options_List) == 0 {
        state.Routing_options =
            types.ListNull(types.ObjectType{AttrTypes: Routing_options_Model{}.AttrTypes()})
    } else {
        state.Routing_options, _ =
            types.ListValueFrom(ctx,
                types.ObjectType{AttrTypes: Routing_options_Model{}.AttrTypes()},
                routing_options_List,
            )
    }
    // Initialize security as Null; only materialize when we have elements
    state.Security =
        types.ListNull(types.ObjectType{AttrTypes: Security_Model{}.AttrTypes()})

    // Build list from device
    security_List := make([]Security_Model,
        len(config.Groups.Security))

    for i_security, v_security := range config.Groups.Security {
        var security_model Security_Model
            
        // Build log list
        security_log_List := make([]Security_Log_Model, len(v_security.Log))

        
		for i_security_log, v_security_log := range v_security.Log {
            var security_log_model Security_Log_Model
            // leaf -> keep pointer semantics
            security_log_model.Mode =
                types.StringPointerValue(v_security_log.Mode)

            security_log_List[i_security_log] =
                security_log_model

            security_log_List[i_security_log] =
                security_log_model
                
        // Build report list
        security_log_report_List := make([]Security_Log_Report_Model, len(v_security_log.Report))

        

        // Write report field as Null when empty, else concrete list
        if len(security_log_report_List) == 0 {
            security_log_model.Report =
                types.ListNull(types.ObjectType{AttrTypes: Security_Log_Report_Model{}.AttrTypes()})
        } else {
            security_log_model.Report, _ =
                types.ListValueFrom(ctx,
                    types.ObjectType{AttrTypes: Security_Log_Report_Model{}.AttrTypes()},
                    security_log_report_List,
                )
        }
        security_log_List[i_security_log] = security_log_model
        }

        // Write log field as Null when empty, else concrete list
        if len(security_log_List) == 0 {
            security_model.Log =
                types.ListNull(types.ObjectType{AttrTypes: Security_Log_Model{}.AttrTypes()})
        } else {
            security_model.Log, _ =
                types.ListValueFrom(ctx,
                    types.ObjectType{AttrTypes: Security_Log_Model{}.AttrTypes()},
                    security_log_List,
                )
        }
        security_List[i_security] = security_model
            
        // Build screen list
        security_screen_List := make([]Security_Screen_Model, len(v_security.Screen))

        
		for i_security_screen, v_security_screen := range v_security.Screen {
            var security_screen_model Security_Screen_Model

            security_screen_List[i_security_screen] =
                security_screen_model
                
        // Build ids-option list
        security_screen_ids_option_List := make([]Security_Screen_Ids_option_Model, len(v_security_screen.Ids_option))

        
		for i_security_screen_ids_option, v_security_screen_ids_option := range v_security_screen.Ids_option {
            var security_screen_ids_option_model Security_Screen_Ids_option_Model
            // leaf -> keep pointer semantics
            security_screen_ids_option_model.Name =
                types.StringPointerValue(v_security_screen_ids_option.Name)

            security_screen_ids_option_List[i_security_screen_ids_option] =
                security_screen_ids_option_model

            security_screen_ids_option_List[i_security_screen_ids_option] =
                security_screen_ids_option_model
                
        // Build icmp list
        security_screen_ids_option_icmp_List := make([]Security_Screen_Ids_option_Icmp_Model, len(v_security_screen_ids_option.Icmp))

        
		for i_security_screen_ids_option_icmp, v_security_screen_ids_option_icmp := range v_security_screen_ids_option.Icmp {
            var security_screen_ids_option_icmp_model Security_Screen_Ids_option_Icmp_Model
            // leaf -> keep pointer semantics
            security_screen_ids_option_icmp_model.Ping_death =
                types.StringPointerValue(v_security_screen_ids_option_icmp.Ping_death)

            security_screen_ids_option_icmp_List[i_security_screen_ids_option_icmp] =
                security_screen_ids_option_icmp_model
        }

        // Write icmp field as Null when empty, else concrete list
        if len(security_screen_ids_option_icmp_List) == 0 {
            security_screen_ids_option_model.Icmp =
                types.ListNull(types.ObjectType{AttrTypes: Security_Screen_Ids_option_Icmp_Model{}.AttrTypes()})
        } else {
            security_screen_ids_option_model.Icmp, _ =
                types.ListValueFrom(ctx,
                    types.ObjectType{AttrTypes: Security_Screen_Ids_option_Icmp_Model{}.AttrTypes()},
                    security_screen_ids_option_icmp_List,
                )
        }
        security_screen_ids_option_List[i_security_screen_ids_option] = security_screen_ids_option_model

            security_screen_ids_option_List[i_security_screen_ids_option] =
                security_screen_ids_option_model
                
        // Build ip list
        security_screen_ids_option_ip_List := make([]Security_Screen_Ids_option_Ip_Model, len(v_security_screen_ids_option.Ip))

        
		for i_security_screen_ids_option_ip, v_security_screen_ids_option_ip := range v_security_screen_ids_option.Ip {
            var security_screen_ids_option_ip_model Security_Screen_Ids_option_Ip_Model
            // leaf -> keep pointer semantics
            security_screen_ids_option_ip_model.Source_route_option =
                types.StringPointerValue(v_security_screen_ids_option_ip.Source_route_option)

            security_screen_ids_option_ip_List[i_security_screen_ids_option_ip] =
                security_screen_ids_option_ip_model
            // leaf -> keep pointer semantics
            security_screen_ids_option_ip_model.Tear_drop =
                types.StringPointerValue(v_security_screen_ids_option_ip.Tear_drop)

            security_screen_ids_option_ip_List[i_security_screen_ids_option_ip] =
                security_screen_ids_option_ip_model
        }

        // Write ip field as Null when empty, else concrete list
        if len(security_screen_ids_option_ip_List) == 0 {
            security_screen_ids_option_model.Ip =
                types.ListNull(types.ObjectType{AttrTypes: Security_Screen_Ids_option_Ip_Model{}.AttrTypes()})
        } else {
            security_screen_ids_option_model.Ip, _ =
                types.ListValueFrom(ctx,
                    types.ObjectType{AttrTypes: Security_Screen_Ids_option_Ip_Model{}.AttrTypes()},
                    security_screen_ids_option_ip_List,
                )
        }
        security_screen_ids_option_List[i_security_screen_ids_option] = security_screen_ids_option_model

            security_screen_ids_option_List[i_security_screen_ids_option] =
                security_screen_ids_option_model
                
        // Build tcp list
        security_screen_ids_option_tcp_List := make([]Security_Screen_Ids_option_Tcp_Model, len(v_security_screen_ids_option.Tcp))

        
		for i_security_screen_ids_option_tcp, v_security_screen_ids_option_tcp := range v_security_screen_ids_option.Tcp {
            var security_screen_ids_option_tcp_model Security_Screen_Ids_option_Tcp_Model

            security_screen_ids_option_tcp_List[i_security_screen_ids_option_tcp] =
                security_screen_ids_option_tcp_model
                
        // Build syn-flood list
        security_screen_ids_option_tcp_syn_flood_List := make([]Security_Screen_Ids_option_Tcp_Syn_flood_Model, len(v_security_screen_ids_option_tcp.Syn_flood))

        
		for i_security_screen_ids_option_tcp_syn_flood, v_security_screen_ids_option_tcp_syn_flood := range v_security_screen_ids_option_tcp.Syn_flood {
            var security_screen_ids_option_tcp_syn_flood_model Security_Screen_Ids_option_Tcp_Syn_flood_Model
            // leaf -> keep pointer semantics
            security_screen_ids_option_tcp_syn_flood_model.Alarm_threshold =
                types.StringPointerValue(v_security_screen_ids_option_tcp_syn_flood.Alarm_threshold)

            security_screen_ids_option_tcp_syn_flood_List[i_security_screen_ids_option_tcp_syn_flood] =
                security_screen_ids_option_tcp_syn_flood_model
            // leaf -> keep pointer semantics
            security_screen_ids_option_tcp_syn_flood_model.Attack_threshold =
                types.StringPointerValue(v_security_screen_ids_option_tcp_syn_flood.Attack_threshold)

            security_screen_ids_option_tcp_syn_flood_List[i_security_screen_ids_option_tcp_syn_flood] =
                security_screen_ids_option_tcp_syn_flood_model
            // leaf -> keep pointer semantics
            security_screen_ids_option_tcp_syn_flood_model.Source_threshold =
                types.StringPointerValue(v_security_screen_ids_option_tcp_syn_flood.Source_threshold)

            security_screen_ids_option_tcp_syn_flood_List[i_security_screen_ids_option_tcp_syn_flood] =
                security_screen_ids_option_tcp_syn_flood_model
            // leaf -> keep pointer semantics
            security_screen_ids_option_tcp_syn_flood_model.Destination_threshold =
                types.StringPointerValue(v_security_screen_ids_option_tcp_syn_flood.Destination_threshold)

            security_screen_ids_option_tcp_syn_flood_List[i_security_screen_ids_option_tcp_syn_flood] =
                security_screen_ids_option_tcp_syn_flood_model
            // leaf -> keep pointer semantics
            security_screen_ids_option_tcp_syn_flood_model.Timeout =
                types.StringPointerValue(v_security_screen_ids_option_tcp_syn_flood.Timeout)

            security_screen_ids_option_tcp_syn_flood_List[i_security_screen_ids_option_tcp_syn_flood] =
                security_screen_ids_option_tcp_syn_flood_model
        }

        // Write syn-flood field as Null when empty, else concrete list
        if len(security_screen_ids_option_tcp_syn_flood_List) == 0 {
            security_screen_ids_option_tcp_model.Syn_flood =
                types.ListNull(types.ObjectType{AttrTypes: Security_Screen_Ids_option_Tcp_Syn_flood_Model{}.AttrTypes()})
        } else {
            security_screen_ids_option_tcp_model.Syn_flood, _ =
                types.ListValueFrom(ctx,
                    types.ObjectType{AttrTypes: Security_Screen_Ids_option_Tcp_Syn_flood_Model{}.AttrTypes()},
                    security_screen_ids_option_tcp_syn_flood_List,
                )
        }
        security_screen_ids_option_tcp_List[i_security_screen_ids_option_tcp] = security_screen_ids_option_tcp_model
            // leaf -> keep pointer semantics
            security_screen_ids_option_tcp_model.Land =
                types.StringPointerValue(v_security_screen_ids_option_tcp.Land)

            security_screen_ids_option_tcp_List[i_security_screen_ids_option_tcp] =
                security_screen_ids_option_tcp_model
        }

        // Write tcp field as Null when empty, else concrete list
        if len(security_screen_ids_option_tcp_List) == 0 {
            security_screen_ids_option_model.Tcp =
                types.ListNull(types.ObjectType{AttrTypes: Security_Screen_Ids_option_Tcp_Model{}.AttrTypes()})
        } else {
            security_screen_ids_option_model.Tcp, _ =
                types.ListValueFrom(ctx,
                    types.ObjectType{AttrTypes: Security_Screen_Ids_option_Tcp_Model{}.AttrTypes()},
                    security_screen_ids_option_tcp_List,
                )
        }
        security_screen_ids_option_List[i_security_screen_ids_option] = security_screen_ids_option_model
        }

        // Write ids-option field as Null when empty, else concrete list
        if len(security_screen_ids_option_List) == 0 {
            security_screen_model.Ids_option =
                types.ListNull(types.ObjectType{AttrTypes: Security_Screen_Ids_option_Model{}.AttrTypes()})
        } else {
            security_screen_model.Ids_option, _ =
                types.ListValueFrom(ctx,
                    types.ObjectType{AttrTypes: Security_Screen_Ids_option_Model{}.AttrTypes()},
                    security_screen_ids_option_List,
                )
        }
        security_screen_List[i_security_screen] = security_screen_model
        }

        // Write screen field as Null when empty, else concrete list
        if len(security_screen_List) == 0 {
            security_model.Screen =
                types.ListNull(types.ObjectType{AttrTypes: Security_Screen_Model{}.AttrTypes()})
        } else {
            security_model.Screen, _ =
                types.ListValueFrom(ctx,
                    types.ObjectType{AttrTypes: Security_Screen_Model{}.AttrTypes()},
                    security_screen_List,
                )
        }
        security_List[i_security] = security_model
            
        // Build policies list
        security_policies_List := make([]Security_Policies_Model, len(v_security.Policies))

        
		for i_security_policies, v_security_policies := range v_security.Policies {
            var security_policies_model Security_Policies_Model

            security_policies_List[i_security_policies] =
                security_policies_model
                
        // Build policy list
        security_policies_policy_List := make([]Security_Policies_Policy_Model, len(v_security_policies.Policy))

        
		for i_security_policies_policy, v_security_policies_policy := range v_security_policies.Policy {
            var security_policies_policy_model Security_Policies_Policy_Model
            // leaf -> keep pointer semantics
            security_policies_policy_model.From_zone_name =
                types.StringPointerValue(v_security_policies_policy.From_zone_name)

            security_policies_policy_List[i_security_policies_policy] =
                security_policies_policy_model
            // leaf -> keep pointer semantics
            security_policies_policy_model.To_zone_name =
                types.StringPointerValue(v_security_policies_policy.To_zone_name)

            security_policies_policy_List[i_security_policies_policy] =
                security_policies_policy_model

            security_policies_policy_List[i_security_policies_policy] =
                security_policies_policy_model
                
        // Build policy list
        security_policies_policy_policy_List := make([]Security_Policies_Policy_Policy_Model, len(v_security_policies_policy.Policy))

        
		for i_security_policies_policy_policy, v_security_policies_policy_policy := range v_security_policies_policy.Policy {
            var security_policies_policy_policy_model Security_Policies_Policy_Policy_Model
            // leaf -> keep pointer semantics
            security_policies_policy_policy_model.Name =
                types.StringPointerValue(v_security_policies_policy_policy.Name)

            security_policies_policy_policy_List[i_security_policies_policy_policy] =
                security_policies_policy_policy_model

            security_policies_policy_policy_List[i_security_policies_policy_policy] =
                security_policies_policy_policy_model
                
        // Build match list
        security_policies_policy_policy_match_List := make([]Security_Policies_Policy_Policy_Match_Model, len(v_security_policies_policy_policy.Match))

        
		for i_security_policies_policy_policy_match, v_security_policies_policy_policy_match := range v_security_policies_policy_policy.Match {
            var security_policies_policy_policy_match_model Security_Policies_Policy_Policy_Match_Model
            // leaf-list -> write Null when nil OR empty (avoid [] when absent)
            if v_security_policies_policy_policy_match.Source_address == nil ||
               len(v_security_policies_policy_policy_match.Source_address) == 0 {
                security_policies_policy_policy_match_model.Source_address =
                    types.ListNull(types.StringType)
            } else {
                src_security_policies_policy_policy_source_address :=
                    v_security_policies_policy_policy_match.Source_address
                vals_security_policies_policy_policy_source_address := make([]*string, len(src_security_policies_policy_policy_source_address))
                copy(vals_security_policies_policy_policy_source_address, src_security_policies_policy_policy_source_address)
                security_policies_policy_policy_match_model.Source_address, _ =
                    types.ListValueFrom(ctx, types.StringType, vals_security_policies_policy_policy_source_address)
            }

            security_policies_policy_policy_match_List[i_security_policies_policy_policy_match] =
                security_policies_policy_policy_match_model
            // leaf-list -> write Null when nil OR empty (avoid [] when absent)
            if v_security_policies_policy_policy_match.Destination_address == nil ||
               len(v_security_policies_policy_policy_match.Destination_address) == 0 {
                security_policies_policy_policy_match_model.Destination_address =
                    types.ListNull(types.StringType)
            } else {
                src_security_policies_policy_policy_destination_address :=
                    v_security_policies_policy_policy_match.Destination_address
                vals_security_policies_policy_policy_destination_address := make([]*string, len(src_security_policies_policy_policy_destination_address))
                copy(vals_security_policies_policy_policy_destination_address, src_security_policies_policy_policy_destination_address)
                security_policies_policy_policy_match_model.Destination_address, _ =
                    types.ListValueFrom(ctx, types.StringType, vals_security_policies_policy_policy_destination_address)
            }

            security_policies_policy_policy_match_List[i_security_policies_policy_policy_match] =
                security_policies_policy_policy_match_model
            // leaf-list -> write Null when nil OR empty (avoid [] when absent)
            if v_security_policies_policy_policy_match.Application == nil ||
               len(v_security_policies_policy_policy_match.Application) == 0 {
                security_policies_policy_policy_match_model.Application =
                    types.ListNull(types.StringType)
            } else {
                src_security_policies_policy_policy_application :=
                    v_security_policies_policy_policy_match.Application
                vals_security_policies_policy_policy_application := make([]*string, len(src_security_policies_policy_policy_application))
                copy(vals_security_policies_policy_policy_application, src_security_policies_policy_policy_application)
                security_policies_policy_policy_match_model.Application, _ =
                    types.ListValueFrom(ctx, types.StringType, vals_security_policies_policy_policy_application)
            }

            security_policies_policy_policy_match_List[i_security_policies_policy_policy_match] =
                security_policies_policy_policy_match_model
        }

        // Write match field as Null when empty, else concrete list
        if len(security_policies_policy_policy_match_List) == 0 {
            security_policies_policy_policy_model.Match =
                types.ListNull(types.ObjectType{AttrTypes: Security_Policies_Policy_Policy_Match_Model{}.AttrTypes()})
        } else {
            security_policies_policy_policy_model.Match, _ =
                types.ListValueFrom(ctx,
                    types.ObjectType{AttrTypes: Security_Policies_Policy_Policy_Match_Model{}.AttrTypes()},
                    security_policies_policy_policy_match_List,
                )
        }
        security_policies_policy_policy_List[i_security_policies_policy_policy] = security_policies_policy_policy_model

            security_policies_policy_policy_List[i_security_policies_policy_policy] =
                security_policies_policy_policy_model
                
        // Build then list
        security_policies_policy_policy_then_List := make([]Security_Policies_Policy_Policy_Then_Model, len(v_security_policies_policy_policy.Then))

        
		for i_security_policies_policy_policy_then, v_security_policies_policy_policy_then := range v_security_policies_policy_policy.Then {
            var security_policies_policy_policy_then_model Security_Policies_Policy_Policy_Then_Model

            security_policies_policy_policy_then_List[i_security_policies_policy_policy_then] =
                security_policies_policy_policy_then_model
                
        // Build permit list
        security_policies_policy_policy_then_permit_List := make([]Security_Policies_Policy_Policy_Then_Permit_Model, len(v_security_policies_policy_policy_then.Permit))

        

        // Write permit field as Null when empty, else concrete list
        if len(security_policies_policy_policy_then_permit_List) == 0 {
            security_policies_policy_policy_then_model.Permit =
                types.ListNull(types.ObjectType{AttrTypes: Security_Policies_Policy_Policy_Then_Permit_Model{}.AttrTypes()})
        } else {
            security_policies_policy_policy_then_model.Permit, _ =
                types.ListValueFrom(ctx,
                    types.ObjectType{AttrTypes: Security_Policies_Policy_Policy_Then_Permit_Model{}.AttrTypes()},
                    security_policies_policy_policy_then_permit_List,
                )
        }
        security_policies_policy_policy_then_List[i_security_policies_policy_policy_then] = security_policies_policy_policy_then_model
        }

        // Write then field as Null when empty, else concrete list
        if len(security_policies_policy_policy_then_List) == 0 {
            security_policies_policy_policy_model.Then =
                types.ListNull(types.ObjectType{AttrTypes: Security_Policies_Policy_Policy_Then_Model{}.AttrTypes()})
        } else {
            security_policies_policy_policy_model.Then, _ =
                types.ListValueFrom(ctx,
                    types.ObjectType{AttrTypes: Security_Policies_Policy_Policy_Then_Model{}.AttrTypes()},
                    security_policies_policy_policy_then_List,
                )
        }
        security_policies_policy_policy_List[i_security_policies_policy_policy] = security_policies_policy_policy_model
        }

        // Write policy field as Null when empty, else concrete list
        if len(security_policies_policy_policy_List) == 0 {
            security_policies_policy_model.Policy =
                types.ListNull(types.ObjectType{AttrTypes: Security_Policies_Policy_Policy_Model{}.AttrTypes()})
        } else {
            security_policies_policy_model.Policy, _ =
                types.ListValueFrom(ctx,
                    types.ObjectType{AttrTypes: Security_Policies_Policy_Policy_Model{}.AttrTypes()},
                    security_policies_policy_policy_List,
                )
        }
        security_policies_policy_List[i_security_policies_policy] = security_policies_policy_model
        }

        // Write policy field as Null when empty, else concrete list
        if len(security_policies_policy_List) == 0 {
            security_policies_model.Policy =
                types.ListNull(types.ObjectType{AttrTypes: Security_Policies_Policy_Model{}.AttrTypes()})
        } else {
            security_policies_model.Policy, _ =
                types.ListValueFrom(ctx,
                    types.ObjectType{AttrTypes: Security_Policies_Policy_Model{}.AttrTypes()},
                    security_policies_policy_List,
                )
        }
        security_policies_List[i_security_policies] = security_policies_model
        }

        // Write policies field as Null when empty, else concrete list
        if len(security_policies_List) == 0 {
            security_model.Policies =
                types.ListNull(types.ObjectType{AttrTypes: Security_Policies_Model{}.AttrTypes()})
        } else {
            security_model.Policies, _ =
                types.ListValueFrom(ctx,
                    types.ObjectType{AttrTypes: Security_Policies_Model{}.AttrTypes()},
                    security_policies_List,
                )
        }
        security_List[i_security] = security_model
            
        // Build zones list
        security_zones_List := make([]Security_Zones_Model, len(v_security.Zones))

        
		for i_security_zones, v_security_zones := range v_security.Zones {
            var security_zones_model Security_Zones_Model

            security_zones_List[i_security_zones] =
                security_zones_model
                
        // Build security-zone list
        security_zones_security_zone_List := make([]Security_Zones_Security_zone_Model, len(v_security_zones.Security_zone))

        
		for i_security_zones_security_zone, v_security_zones_security_zone := range v_security_zones.Security_zone {
            var security_zones_security_zone_model Security_Zones_Security_zone_Model
            // leaf -> keep pointer semantics
            security_zones_security_zone_model.Name =
                types.StringPointerValue(v_security_zones_security_zone.Name)

            security_zones_security_zone_List[i_security_zones_security_zone] =
                security_zones_security_zone_model
            // leaf -> keep pointer semantics
            security_zones_security_zone_model.Tcp_rst =
                types.StringPointerValue(v_security_zones_security_zone.Tcp_rst)

            security_zones_security_zone_List[i_security_zones_security_zone] =
                security_zones_security_zone_model
            // leaf -> keep pointer semantics
            security_zones_security_zone_model.Screen =
                types.StringPointerValue(v_security_zones_security_zone.Screen)

            security_zones_security_zone_List[i_security_zones_security_zone] =
                security_zones_security_zone_model

            security_zones_security_zone_List[i_security_zones_security_zone] =
                security_zones_security_zone_model
                
        // Build host-inbound-traffic list
        security_zones_security_zone_host_inbound_traffic_List := make([]Security_Zones_Security_zone_Host_inbound_traffic_Model, len(v_security_zones_security_zone.Host_inbound_traffic))

        
		for i_security_zones_security_zone_host_inbound_traffic, v_security_zones_security_zone_host_inbound_traffic := range v_security_zones_security_zone.Host_inbound_traffic {
            var security_zones_security_zone_host_inbound_traffic_model Security_Zones_Security_zone_Host_inbound_traffic_Model

            security_zones_security_zone_host_inbound_traffic_List[i_security_zones_security_zone_host_inbound_traffic] =
                security_zones_security_zone_host_inbound_traffic_model
                
        // Build system-services list
        security_zones_security_zone_host_inbound_traffic_system_services_List := make([]Security_Zones_Security_zone_Host_inbound_traffic_System_services_Model, len(v_security_zones_security_zone_host_inbound_traffic.System_services))

        
		for i_security_zones_security_zone_host_inbound_traffic_system_services, v_security_zones_security_zone_host_inbound_traffic_system_services := range v_security_zones_security_zone_host_inbound_traffic.System_services {
            var security_zones_security_zone_host_inbound_traffic_system_services_model Security_Zones_Security_zone_Host_inbound_traffic_System_services_Model
            // leaf -> keep pointer semantics
            security_zones_security_zone_host_inbound_traffic_system_services_model.Name =
                types.StringPointerValue(v_security_zones_security_zone_host_inbound_traffic_system_services.Name)

            security_zones_security_zone_host_inbound_traffic_system_services_List[i_security_zones_security_zone_host_inbound_traffic_system_services] =
                security_zones_security_zone_host_inbound_traffic_system_services_model
        }

        // Write system-services field as Null when empty, else concrete list
        if len(security_zones_security_zone_host_inbound_traffic_system_services_List) == 0 {
            security_zones_security_zone_host_inbound_traffic_model.System_services =
                types.ListNull(types.ObjectType{AttrTypes: Security_Zones_Security_zone_Host_inbound_traffic_System_services_Model{}.AttrTypes()})
        } else {
            security_zones_security_zone_host_inbound_traffic_model.System_services, _ =
                types.ListValueFrom(ctx,
                    types.ObjectType{AttrTypes: Security_Zones_Security_zone_Host_inbound_traffic_System_services_Model{}.AttrTypes()},
                    security_zones_security_zone_host_inbound_traffic_system_services_List,
                )
        }
        security_zones_security_zone_host_inbound_traffic_List[i_security_zones_security_zone_host_inbound_traffic] = security_zones_security_zone_host_inbound_traffic_model

            security_zones_security_zone_host_inbound_traffic_List[i_security_zones_security_zone_host_inbound_traffic] =
                security_zones_security_zone_host_inbound_traffic_model
                
        // Build protocols list
        security_zones_security_zone_host_inbound_traffic_protocols_List := make([]Security_Zones_Security_zone_Host_inbound_traffic_Protocols_Model, len(v_security_zones_security_zone_host_inbound_traffic.Protocols))

        
		for i_security_zones_security_zone_host_inbound_traffic_protocols, v_security_zones_security_zone_host_inbound_traffic_protocols := range v_security_zones_security_zone_host_inbound_traffic.Protocols {
            var security_zones_security_zone_host_inbound_traffic_protocols_model Security_Zones_Security_zone_Host_inbound_traffic_Protocols_Model
            // leaf -> keep pointer semantics
            security_zones_security_zone_host_inbound_traffic_protocols_model.Name =
                types.StringPointerValue(v_security_zones_security_zone_host_inbound_traffic_protocols.Name)

            security_zones_security_zone_host_inbound_traffic_protocols_List[i_security_zones_security_zone_host_inbound_traffic_protocols] =
                security_zones_security_zone_host_inbound_traffic_protocols_model
        }

        // Write protocols field as Null when empty, else concrete list
        if len(security_zones_security_zone_host_inbound_traffic_protocols_List) == 0 {
            security_zones_security_zone_host_inbound_traffic_model.Protocols =
                types.ListNull(types.ObjectType{AttrTypes: Security_Zones_Security_zone_Host_inbound_traffic_Protocols_Model{}.AttrTypes()})
        } else {
            security_zones_security_zone_host_inbound_traffic_model.Protocols, _ =
                types.ListValueFrom(ctx,
                    types.ObjectType{AttrTypes: Security_Zones_Security_zone_Host_inbound_traffic_Protocols_Model{}.AttrTypes()},
                    security_zones_security_zone_host_inbound_traffic_protocols_List,
                )
        }
        security_zones_security_zone_host_inbound_traffic_List[i_security_zones_security_zone_host_inbound_traffic] = security_zones_security_zone_host_inbound_traffic_model
        }

        // Write host-inbound-traffic field as Null when empty, else concrete list
        if len(security_zones_security_zone_host_inbound_traffic_List) == 0 {
            security_zones_security_zone_model.Host_inbound_traffic =
                types.ListNull(types.ObjectType{AttrTypes: Security_Zones_Security_zone_Host_inbound_traffic_Model{}.AttrTypes()})
        } else {
            security_zones_security_zone_model.Host_inbound_traffic, _ =
                types.ListValueFrom(ctx,
                    types.ObjectType{AttrTypes: Security_Zones_Security_zone_Host_inbound_traffic_Model{}.AttrTypes()},
                    security_zones_security_zone_host_inbound_traffic_List,
                )
        }
        security_zones_security_zone_List[i_security_zones_security_zone] = security_zones_security_zone_model

            security_zones_security_zone_List[i_security_zones_security_zone] =
                security_zones_security_zone_model
                
        // Build interfaces list
        security_zones_security_zone_interfaces_List := make([]Security_Zones_Security_zone_Interfaces_Model, len(v_security_zones_security_zone.Interfaces))

        
		for i_security_zones_security_zone_interfaces, v_security_zones_security_zone_interfaces := range v_security_zones_security_zone.Interfaces {
            var security_zones_security_zone_interfaces_model Security_Zones_Security_zone_Interfaces_Model
            // leaf -> keep pointer semantics
            security_zones_security_zone_interfaces_model.Name =
                types.StringPointerValue(v_security_zones_security_zone_interfaces.Name)

            security_zones_security_zone_interfaces_List[i_security_zones_security_zone_interfaces] =
                security_zones_security_zone_interfaces_model
        }

        // Write interfaces field as Null when empty, else concrete list
        if len(security_zones_security_zone_interfaces_List) == 0 {
            security_zones_security_zone_model.Interfaces =
                types.ListNull(types.ObjectType{AttrTypes: Security_Zones_Security_zone_Interfaces_Model{}.AttrTypes()})
        } else {
            security_zones_security_zone_model.Interfaces, _ =
                types.ListValueFrom(ctx,
                    types.ObjectType{AttrTypes: Security_Zones_Security_zone_Interfaces_Model{}.AttrTypes()},
                    security_zones_security_zone_interfaces_List,
                )
        }
        security_zones_security_zone_List[i_security_zones_security_zone] = security_zones_security_zone_model
        }

        // Write security-zone field as Null when empty, else concrete list
        if len(security_zones_security_zone_List) == 0 {
            security_zones_model.Security_zone =
                types.ListNull(types.ObjectType{AttrTypes: Security_Zones_Security_zone_Model{}.AttrTypes()})
        } else {
            security_zones_model.Security_zone, _ =
                types.ListValueFrom(ctx,
                    types.ObjectType{AttrTypes: Security_Zones_Security_zone_Model{}.AttrTypes()},
                    security_zones_security_zone_List,
                )
        }
        security_zones_List[i_security_zones] = security_zones_model
        }

        // Write zones field as Null when empty, else concrete list
        if len(security_zones_List) == 0 {
            security_model.Zones =
                types.ListNull(types.ObjectType{AttrTypes: Security_Zones_Model{}.AttrTypes()})
        } else {
            security_model.Zones, _ =
                types.ListValueFrom(ctx,
                    types.ObjectType{AttrTypes: Security_Zones_Model{}.AttrTypes()},
                    security_zones_List,
                )
        }
        security_List[i_security] = security_model

        security_List[i_security] = security_model
    }

    // Write parent list as Null when empty
    if len(security_List) == 0 {
        state.Security =
            types.ListNull(types.ObjectType{AttrTypes: Security_Model{}.AttrTypes()})
    } else {
        state.Security, _ =
            types.ListValueFrom(ctx,
                types.ObjectType{AttrTypes: Security_Model{}.AttrTypes()},
                security_List,
            )
    }
    // Initialize snmp as Null; only materialize when we have elements
    state.Snmp =
        types.ListNull(types.ObjectType{AttrTypes: Snmp_Model{}.AttrTypes()})

    // Build list from device
    snmp_List := make([]Snmp_Model,
        len(config.Groups.Snmp))

    for i_snmp, v_snmp := range config.Groups.Snmp {
        var snmp_model Snmp_Model
            
        // Build community list
        snmp_community_List := make([]Snmp_Community_Model, len(v_snmp.Community))

        
		for i_snmp_community, v_snmp_community := range v_snmp.Community {
            var snmp_community_model Snmp_Community_Model
            // leaf -> keep pointer semantics
            snmp_community_model.Name =
                types.StringPointerValue(v_snmp_community.Name)

            snmp_community_List[i_snmp_community] =
                snmp_community_model
            // leaf -> keep pointer semantics
            snmp_community_model.Authorization =
                types.StringPointerValue(v_snmp_community.Authorization)

            snmp_community_List[i_snmp_community] =
                snmp_community_model
        }

        // Write community field as Null when empty, else concrete list
        if len(snmp_community_List) == 0 {
            snmp_model.Community =
                types.ListNull(types.ObjectType{AttrTypes: Snmp_Community_Model{}.AttrTypes()})
        } else {
            snmp_model.Community, _ =
                types.ListValueFrom(ctx,
                    types.ObjectType{AttrTypes: Snmp_Community_Model{}.AttrTypes()},
                    snmp_community_List,
                )
        }
        snmp_List[i_snmp] = snmp_model

        snmp_List[i_snmp] = snmp_model
    }

    // Write parent list as Null when empty
    if len(snmp_List) == 0 {
        state.Snmp =
            types.ListNull(types.ObjectType{AttrTypes: Snmp_Model{}.AttrTypes()})
    } else {
        state.Snmp, _ =
            types.ListValueFrom(ctx,
                types.ObjectType{AttrTypes: Snmp_Model{}.AttrTypes()},
                snmp_List,
            )
    }
    // Initialize system as Null; only materialize when we have elements
    state.System =
        types.ListNull(types.ObjectType{AttrTypes: System_Model{}.AttrTypes()})

    // Build list from device
    system_List := make([]System_Model,
        len(config.Groups.System))

    for i_system, v_system := range config.Groups.System {
        var system_model System_Model
            
        // Build login list
        system_login_List := make([]System_Login_Model, len(v_system.Login))

        
		for i_system_login, v_system_login := range v_system.Login {
            var system_login_model System_Login_Model

            system_login_List[i_system_login] =
                system_login_model
                
        // Build user list
        system_login_user_List := make([]System_Login_User_Model, len(v_system_login.User))

        
		for i_system_login_user, v_system_login_user := range v_system_login.User {
            var system_login_user_model System_Login_User_Model
            // leaf -> keep pointer semantics
            system_login_user_model.Name =
                types.StringPointerValue(v_system_login_user.Name)

            system_login_user_List[i_system_login_user] =
                system_login_user_model
            // leaf -> keep pointer semantics
            system_login_user_model.Uid =
                types.StringPointerValue(v_system_login_user.Uid)

            system_login_user_List[i_system_login_user] =
                system_login_user_model
            // leaf -> keep pointer semantics
            system_login_user_model.Class =
                types.StringPointerValue(v_system_login_user.Class)

            system_login_user_List[i_system_login_user] =
                system_login_user_model

            system_login_user_List[i_system_login_user] =
                system_login_user_model
                
        // Build authentication list
        system_login_user_authentication_List := make([]System_Login_User_Authentication_Model, len(v_system_login_user.Authentication))

        
		for i_system_login_user_authentication, v_system_login_user_authentication := range v_system_login_user.Authentication {
            var system_login_user_authentication_model System_Login_User_Authentication_Model
            // leaf -> keep pointer semantics
            system_login_user_authentication_model.Encrypted_password =
                types.StringPointerValue(v_system_login_user_authentication.Encrypted_password)

            system_login_user_authentication_List[i_system_login_user_authentication] =
                system_login_user_authentication_model
        }

        // Write authentication field as Null when empty, else concrete list
        if len(system_login_user_authentication_List) == 0 {
            system_login_user_model.Authentication =
                types.ListNull(types.ObjectType{AttrTypes: System_Login_User_Authentication_Model{}.AttrTypes()})
        } else {
            system_login_user_model.Authentication, _ =
                types.ListValueFrom(ctx,
                    types.ObjectType{AttrTypes: System_Login_User_Authentication_Model{}.AttrTypes()},
                    system_login_user_authentication_List,
                )
        }
        system_login_user_List[i_system_login_user] = system_login_user_model
        }

        // Write user field as Null when empty, else concrete list
        if len(system_login_user_List) == 0 {
            system_login_model.User =
                types.ListNull(types.ObjectType{AttrTypes: System_Login_User_Model{}.AttrTypes()})
        } else {
            system_login_model.User, _ =
                types.ListValueFrom(ctx,
                    types.ObjectType{AttrTypes: System_Login_User_Model{}.AttrTypes()},
                    system_login_user_List,
                )
        }
        system_login_List[i_system_login] = system_login_model
            // leaf -> keep pointer semantics
            system_login_model.Message =
                types.StringPointerValue(v_system_login.Message)

            system_login_List[i_system_login] =
                system_login_model
        }

        // Write login field as Null when empty, else concrete list
        if len(system_login_List) == 0 {
            system_model.Login =
                types.ListNull(types.ObjectType{AttrTypes: System_Login_Model{}.AttrTypes()})
        } else {
            system_model.Login, _ =
                types.ListValueFrom(ctx,
                    types.ObjectType{AttrTypes: System_Login_Model{}.AttrTypes()},
                    system_login_List,
                )
        }
        system_List[i_system] = system_model
            
        // Build root-authentication list
        system_root_authentication_List := make([]System_Root_authentication_Model, len(v_system.Root_authentication))

        
		for i_system_root_authentication, v_system_root_authentication := range v_system.Root_authentication {
            var system_root_authentication_model System_Root_authentication_Model
            // leaf -> keep pointer semantics
            system_root_authentication_model.Encrypted_password =
                types.StringPointerValue(v_system_root_authentication.Encrypted_password)

            system_root_authentication_List[i_system_root_authentication] =
                system_root_authentication_model
        }

        // Write root-authentication field as Null when empty, else concrete list
        if len(system_root_authentication_List) == 0 {
            system_model.Root_authentication =
                types.ListNull(types.ObjectType{AttrTypes: System_Root_authentication_Model{}.AttrTypes()})
        } else {
            system_model.Root_authentication, _ =
                types.ListValueFrom(ctx,
                    types.ObjectType{AttrTypes: System_Root_authentication_Model{}.AttrTypes()},
                    system_root_authentication_List,
                )
        }
        system_List[i_system] = system_model
            
        // Build services list
        system_services_List := make([]System_Services_Model, len(v_system.Services))

        
		for i_system_services, v_system_services := range v_system.Services {
            var system_services_model System_Services_Model

            system_services_List[i_system_services] =
                system_services_model
                
        // Build ssh list
        system_services_ssh_List := make([]System_Services_Ssh_Model, len(v_system_services.Ssh))

        
		for i_system_services_ssh, v_system_services_ssh := range v_system_services.Ssh {
            var system_services_ssh_model System_Services_Ssh_Model
            // leaf -> keep pointer semantics
            system_services_ssh_model.Root_login =
                types.StringPointerValue(v_system_services_ssh.Root_login)

            system_services_ssh_List[i_system_services_ssh] =
                system_services_ssh_model
        }

        // Write ssh field as Null when empty, else concrete list
        if len(system_services_ssh_List) == 0 {
            system_services_model.Ssh =
                types.ListNull(types.ObjectType{AttrTypes: System_Services_Ssh_Model{}.AttrTypes()})
        } else {
            system_services_model.Ssh, _ =
                types.ListValueFrom(ctx,
                    types.ObjectType{AttrTypes: System_Services_Ssh_Model{}.AttrTypes()},
                    system_services_ssh_List,
                )
        }
        system_services_List[i_system_services] = system_services_model

            system_services_List[i_system_services] =
                system_services_model
                
        // Build extension-service list
        system_services_extension_service_List := make([]System_Services_Extension_service_Model, len(v_system_services.Extension_service))

        
		for i_system_services_extension_service, v_system_services_extension_service := range v_system_services.Extension_service {
            var system_services_extension_service_model System_Services_Extension_service_Model

            system_services_extension_service_List[i_system_services_extension_service] =
                system_services_extension_service_model
                
        // Build request-response list
        system_services_extension_service_request_response_List := make([]System_Services_Extension_service_Request_response_Model, len(v_system_services_extension_service.Request_response))

        
		for i_system_services_extension_service_request_response, v_system_services_extension_service_request_response := range v_system_services_extension_service.Request_response {
            var system_services_extension_service_request_response_model System_Services_Extension_service_Request_response_Model

            system_services_extension_service_request_response_List[i_system_services_extension_service_request_response] =
                system_services_extension_service_request_response_model
                
        // Build grpc list
        system_services_extension_service_request_response_grpc_List := make([]System_Services_Extension_service_Request_response_Grpc_Model, len(v_system_services_extension_service_request_response.Grpc))

        
		for i_system_services_extension_service_request_response_grpc, v_system_services_extension_service_request_response_grpc := range v_system_services_extension_service_request_response.Grpc {
            var system_services_extension_service_request_response_grpc_model System_Services_Extension_service_Request_response_Grpc_Model
            // leaf -> keep pointer semantics
            system_services_extension_service_request_response_grpc_model.Max_connections =
                types.StringPointerValue(v_system_services_extension_service_request_response_grpc.Max_connections)

            system_services_extension_service_request_response_grpc_List[i_system_services_extension_service_request_response_grpc] =
                system_services_extension_service_request_response_grpc_model
        }

        // Write grpc field as Null when empty, else concrete list
        if len(system_services_extension_service_request_response_grpc_List) == 0 {
            system_services_extension_service_request_response_model.Grpc =
                types.ListNull(types.ObjectType{AttrTypes: System_Services_Extension_service_Request_response_Grpc_Model{}.AttrTypes()})
        } else {
            system_services_extension_service_request_response_model.Grpc, _ =
                types.ListValueFrom(ctx,
                    types.ObjectType{AttrTypes: System_Services_Extension_service_Request_response_Grpc_Model{}.AttrTypes()},
                    system_services_extension_service_request_response_grpc_List,
                )
        }
        system_services_extension_service_request_response_List[i_system_services_extension_service_request_response] = system_services_extension_service_request_response_model
        }

        // Write request-response field as Null when empty, else concrete list
        if len(system_services_extension_service_request_response_List) == 0 {
            system_services_extension_service_model.Request_response =
                types.ListNull(types.ObjectType{AttrTypes: System_Services_Extension_service_Request_response_Model{}.AttrTypes()})
        } else {
            system_services_extension_service_model.Request_response, _ =
                types.ListValueFrom(ctx,
                    types.ObjectType{AttrTypes: System_Services_Extension_service_Request_response_Model{}.AttrTypes()},
                    system_services_extension_service_request_response_List,
                )
        }
        system_services_extension_service_List[i_system_services_extension_service] = system_services_extension_service_model

            system_services_extension_service_List[i_system_services_extension_service] =
                system_services_extension_service_model
                
        // Build notification list
        system_services_extension_service_notification_List := make([]System_Services_Extension_service_Notification_Model, len(v_system_services_extension_service.Notification))

        
		for i_system_services_extension_service_notification, v_system_services_extension_service_notification := range v_system_services_extension_service.Notification {
            var system_services_extension_service_notification_model System_Services_Extension_service_Notification_Model

            system_services_extension_service_notification_List[i_system_services_extension_service_notification] =
                system_services_extension_service_notification_model
                
        // Build allow-clients list
        system_services_extension_service_notification_allow_clients_List := make([]System_Services_Extension_service_Notification_Allow_clients_Model, len(v_system_services_extension_service_notification.Allow_clients))

        
		for i_system_services_extension_service_notification_allow_clients, v_system_services_extension_service_notification_allow_clients := range v_system_services_extension_service_notification.Allow_clients {
            var system_services_extension_service_notification_allow_clients_model System_Services_Extension_service_Notification_Allow_clients_Model
            // leaf-list -> write Null when nil OR empty (avoid [] when absent)
            if v_system_services_extension_service_notification_allow_clients.Address == nil ||
               len(v_system_services_extension_service_notification_allow_clients.Address) == 0 {
                system_services_extension_service_notification_allow_clients_model.Address =
                    types.ListNull(types.StringType)
            } else {
                src_system_services_extension_service_notification_address :=
                    v_system_services_extension_service_notification_allow_clients.Address
                vals_system_services_extension_service_notification_address := make([]*string, len(src_system_services_extension_service_notification_address))
                copy(vals_system_services_extension_service_notification_address, src_system_services_extension_service_notification_address)
                system_services_extension_service_notification_allow_clients_model.Address, _ =
                    types.ListValueFrom(ctx, types.StringType, vals_system_services_extension_service_notification_address)
            }

            system_services_extension_service_notification_allow_clients_List[i_system_services_extension_service_notification_allow_clients] =
                system_services_extension_service_notification_allow_clients_model
        }

        // Write allow-clients field as Null when empty, else concrete list
        if len(system_services_extension_service_notification_allow_clients_List) == 0 {
            system_services_extension_service_notification_model.Allow_clients =
                types.ListNull(types.ObjectType{AttrTypes: System_Services_Extension_service_Notification_Allow_clients_Model{}.AttrTypes()})
        } else {
            system_services_extension_service_notification_model.Allow_clients, _ =
                types.ListValueFrom(ctx,
                    types.ObjectType{AttrTypes: System_Services_Extension_service_Notification_Allow_clients_Model{}.AttrTypes()},
                    system_services_extension_service_notification_allow_clients_List,
                )
        }
        system_services_extension_service_notification_List[i_system_services_extension_service_notification] = system_services_extension_service_notification_model
        }

        // Write notification field as Null when empty, else concrete list
        if len(system_services_extension_service_notification_List) == 0 {
            system_services_extension_service_model.Notification =
                types.ListNull(types.ObjectType{AttrTypes: System_Services_Extension_service_Notification_Model{}.AttrTypes()})
        } else {
            system_services_extension_service_model.Notification, _ =
                types.ListValueFrom(ctx,
                    types.ObjectType{AttrTypes: System_Services_Extension_service_Notification_Model{}.AttrTypes()},
                    system_services_extension_service_notification_List,
                )
        }
        system_services_extension_service_List[i_system_services_extension_service] = system_services_extension_service_model
        }

        // Write extension-service field as Null when empty, else concrete list
        if len(system_services_extension_service_List) == 0 {
            system_services_model.Extension_service =
                types.ListNull(types.ObjectType{AttrTypes: System_Services_Extension_service_Model{}.AttrTypes()})
        } else {
            system_services_model.Extension_service, _ =
                types.ListValueFrom(ctx,
                    types.ObjectType{AttrTypes: System_Services_Extension_service_Model{}.AttrTypes()},
                    system_services_extension_service_List,
                )
        }
        system_services_List[i_system_services] = system_services_model

            system_services_List[i_system_services] =
                system_services_model
                
        // Build netconf list
        system_services_netconf_List := make([]System_Services_Netconf_Model, len(v_system_services.Netconf))

        
		for i_system_services_netconf, v_system_services_netconf := range v_system_services.Netconf {
            var system_services_netconf_model System_Services_Netconf_Model

            system_services_netconf_List[i_system_services_netconf] =
                system_services_netconf_model
                
        // Build ssh list
        system_services_netconf_ssh_List := make([]System_Services_Netconf_Ssh_Model, len(v_system_services_netconf.Ssh))

        

        // Write ssh field as Null when empty, else concrete list
        if len(system_services_netconf_ssh_List) == 0 {
            system_services_netconf_model.Ssh =
                types.ListNull(types.ObjectType{AttrTypes: System_Services_Netconf_Ssh_Model{}.AttrTypes()})
        } else {
            system_services_netconf_model.Ssh, _ =
                types.ListValueFrom(ctx,
                    types.ObjectType{AttrTypes: System_Services_Netconf_Ssh_Model{}.AttrTypes()},
                    system_services_netconf_ssh_List,
                )
        }
        system_services_netconf_List[i_system_services_netconf] = system_services_netconf_model
        }

        // Write netconf field as Null when empty, else concrete list
        if len(system_services_netconf_List) == 0 {
            system_services_model.Netconf =
                types.ListNull(types.ObjectType{AttrTypes: System_Services_Netconf_Model{}.AttrTypes()})
        } else {
            system_services_model.Netconf, _ =
                types.ListValueFrom(ctx,
                    types.ObjectType{AttrTypes: System_Services_Netconf_Model{}.AttrTypes()},
                    system_services_netconf_List,
                )
        }
        system_services_List[i_system_services] = system_services_model

            system_services_List[i_system_services] =
                system_services_model
                
        // Build rest list
        system_services_rest_List := make([]System_Services_Rest_Model, len(v_system_services.Rest))

        
		for i_system_services_rest, v_system_services_rest := range v_system_services.Rest {
            var system_services_rest_model System_Services_Rest_Model

            system_services_rest_List[i_system_services_rest] =
                system_services_rest_model
                
        // Build http list
        system_services_rest_http_List := make([]System_Services_Rest_Http_Model, len(v_system_services_rest.Http))

        
		for i_system_services_rest_http, v_system_services_rest_http := range v_system_services_rest.Http {
            var system_services_rest_http_model System_Services_Rest_Http_Model
            // leaf -> keep pointer semantics
            system_services_rest_http_model.Port =
                types.StringPointerValue(v_system_services_rest_http.Port)

            system_services_rest_http_List[i_system_services_rest_http] =
                system_services_rest_http_model
        }

        // Write http field as Null when empty, else concrete list
        if len(system_services_rest_http_List) == 0 {
            system_services_rest_model.Http =
                types.ListNull(types.ObjectType{AttrTypes: System_Services_Rest_Http_Model{}.AttrTypes()})
        } else {
            system_services_rest_model.Http, _ =
                types.ListValueFrom(ctx,
                    types.ObjectType{AttrTypes: System_Services_Rest_Http_Model{}.AttrTypes()},
                    system_services_rest_http_List,
                )
        }
        system_services_rest_List[i_system_services_rest] = system_services_rest_model
            // leaf -> keep pointer semantics
            system_services_rest_model.Enable_explorer =
                types.StringPointerValue(v_system_services_rest.Enable_explorer)

            system_services_rest_List[i_system_services_rest] =
                system_services_rest_model
        }

        // Write rest field as Null when empty, else concrete list
        if len(system_services_rest_List) == 0 {
            system_services_model.Rest =
                types.ListNull(types.ObjectType{AttrTypes: System_Services_Rest_Model{}.AttrTypes()})
        } else {
            system_services_model.Rest, _ =
                types.ListValueFrom(ctx,
                    types.ObjectType{AttrTypes: System_Services_Rest_Model{}.AttrTypes()},
                    system_services_rest_List,
                )
        }
        system_services_List[i_system_services] = system_services_model
        }

        // Write services field as Null when empty, else concrete list
        if len(system_services_List) == 0 {
            system_model.Services =
                types.ListNull(types.ObjectType{AttrTypes: System_Services_Model{}.AttrTypes()})
        } else {
            system_model.Services, _ =
                types.ListValueFrom(ctx,
                    types.ObjectType{AttrTypes: System_Services_Model{}.AttrTypes()},
                    system_services_List,
                )
        }
        system_List[i_system] = system_model
            
        // Build syslog list
        system_syslog_List := make([]System_Syslog_Model, len(v_system.Syslog))

        
		for i_system_syslog, v_system_syslog := range v_system.Syslog {
            var system_syslog_model System_Syslog_Model

            system_syslog_List[i_system_syslog] =
                system_syslog_model
                
        // Build user list
        system_syslog_user_List := make([]System_Syslog_User_Model, len(v_system_syslog.User))

        
		for i_system_syslog_user, v_system_syslog_user := range v_system_syslog.User {
            var system_syslog_user_model System_Syslog_User_Model
            // leaf -> keep pointer semantics
            system_syslog_user_model.Name =
                types.StringPointerValue(v_system_syslog_user.Name)

            system_syslog_user_List[i_system_syslog_user] =
                system_syslog_user_model

            system_syslog_user_List[i_system_syslog_user] =
                system_syslog_user_model
                
        // Build contents list
        system_syslog_user_contents_List := make([]System_Syslog_User_Contents_Model, len(v_system_syslog_user.Contents))

        
		for i_system_syslog_user_contents, v_system_syslog_user_contents := range v_system_syslog_user.Contents {
            var system_syslog_user_contents_model System_Syslog_User_Contents_Model
            // leaf -> keep pointer semantics
            system_syslog_user_contents_model.Name =
                types.StringPointerValue(v_system_syslog_user_contents.Name)

            system_syslog_user_contents_List[i_system_syslog_user_contents] =
                system_syslog_user_contents_model
            // leaf -> keep pointer semantics
            system_syslog_user_contents_model.Emergency =
                types.StringPointerValue(v_system_syslog_user_contents.Emergency)

            system_syslog_user_contents_List[i_system_syslog_user_contents] =
                system_syslog_user_contents_model
        }

        // Write contents field as Null when empty, else concrete list
        if len(system_syslog_user_contents_List) == 0 {
            system_syslog_user_model.Contents =
                types.ListNull(types.ObjectType{AttrTypes: System_Syslog_User_Contents_Model{}.AttrTypes()})
        } else {
            system_syslog_user_model.Contents, _ =
                types.ListValueFrom(ctx,
                    types.ObjectType{AttrTypes: System_Syslog_User_Contents_Model{}.AttrTypes()},
                    system_syslog_user_contents_List,
                )
        }
        system_syslog_user_List[i_system_syslog_user] = system_syslog_user_model
        }

        // Write user field as Null when empty, else concrete list
        if len(system_syslog_user_List) == 0 {
            system_syslog_model.User =
                types.ListNull(types.ObjectType{AttrTypes: System_Syslog_User_Model{}.AttrTypes()})
        } else {
            system_syslog_model.User, _ =
                types.ListValueFrom(ctx,
                    types.ObjectType{AttrTypes: System_Syslog_User_Model{}.AttrTypes()},
                    system_syslog_user_List,
                )
        }
        system_syslog_List[i_system_syslog] = system_syslog_model

            system_syslog_List[i_system_syslog] =
                system_syslog_model
                
        // Build file list
        system_syslog_file_List := make([]System_Syslog_File_Model, len(v_system_syslog.File))

        
		for i_system_syslog_file, v_system_syslog_file := range v_system_syslog.File {
            var system_syslog_file_model System_Syslog_File_Model
            // leaf -> keep pointer semantics
            system_syslog_file_model.Name =
                types.StringPointerValue(v_system_syslog_file.Name)

            system_syslog_file_List[i_system_syslog_file] =
                system_syslog_file_model

            system_syslog_file_List[i_system_syslog_file] =
                system_syslog_file_model
                
        // Build contents list
        system_syslog_file_contents_List := make([]System_Syslog_File_Contents_Model, len(v_system_syslog_file.Contents))

        
		for i_system_syslog_file_contents, v_system_syslog_file_contents := range v_system_syslog_file.Contents {
            var system_syslog_file_contents_model System_Syslog_File_Contents_Model
            // leaf -> keep pointer semantics
            system_syslog_file_contents_model.Name =
                types.StringPointerValue(v_system_syslog_file_contents.Name)

            system_syslog_file_contents_List[i_system_syslog_file_contents] =
                system_syslog_file_contents_model
            // leaf -> keep pointer semantics
            system_syslog_file_contents_model.Any =
                types.StringPointerValue(v_system_syslog_file_contents.Any)

            system_syslog_file_contents_List[i_system_syslog_file_contents] =
                system_syslog_file_contents_model
            // leaf -> keep pointer semantics
            system_syslog_file_contents_model.Info =
                types.StringPointerValue(v_system_syslog_file_contents.Info)

            system_syslog_file_contents_List[i_system_syslog_file_contents] =
                system_syslog_file_contents_model
        }

        // Write contents field as Null when empty, else concrete list
        if len(system_syslog_file_contents_List) == 0 {
            system_syslog_file_model.Contents =
                types.ListNull(types.ObjectType{AttrTypes: System_Syslog_File_Contents_Model{}.AttrTypes()})
        } else {
            system_syslog_file_model.Contents, _ =
                types.ListValueFrom(ctx,
                    types.ObjectType{AttrTypes: System_Syslog_File_Contents_Model{}.AttrTypes()},
                    system_syslog_file_contents_List,
                )
        }
        system_syslog_file_List[i_system_syslog_file] = system_syslog_file_model
        }

        // Write file field as Null when empty, else concrete list
        if len(system_syslog_file_List) == 0 {
            system_syslog_model.File =
                types.ListNull(types.ObjectType{AttrTypes: System_Syslog_File_Model{}.AttrTypes()})
        } else {
            system_syslog_model.File, _ =
                types.ListValueFrom(ctx,
                    types.ObjectType{AttrTypes: System_Syslog_File_Model{}.AttrTypes()},
                    system_syslog_file_List,
                )
        }
        system_syslog_List[i_system_syslog] = system_syslog_model
        }

        // Write syslog field as Null when empty, else concrete list
        if len(system_syslog_List) == 0 {
            system_model.Syslog =
                types.ListNull(types.ObjectType{AttrTypes: System_Syslog_Model{}.AttrTypes()})
        } else {
            system_model.Syslog, _ =
                types.ListValueFrom(ctx,
                    types.ObjectType{AttrTypes: System_Syslog_Model{}.AttrTypes()},
                    system_syslog_List,
                )
        }
        system_List[i_system] = system_model
            
        // Build license list
        system_license_List := make([]System_License_Model, len(v_system.License))

        
		for i_system_license, v_system_license := range v_system.License {
            var system_license_model System_License_Model

            system_license_List[i_system_license] =
                system_license_model
                
        // Build autoupdate list
        system_license_autoupdate_List := make([]System_License_Autoupdate_Model, len(v_system_license.Autoupdate))

        
		for i_system_license_autoupdate, v_system_license_autoupdate := range v_system_license.Autoupdate {
            var system_license_autoupdate_model System_License_Autoupdate_Model

            system_license_autoupdate_List[i_system_license_autoupdate] =
                system_license_autoupdate_model
                
        // Build url list
        system_license_autoupdate_url_List := make([]System_License_Autoupdate_Url_Model, len(v_system_license_autoupdate.Url))

        
		for i_system_license_autoupdate_url, v_system_license_autoupdate_url := range v_system_license_autoupdate.Url {
            var system_license_autoupdate_url_model System_License_Autoupdate_Url_Model
            // leaf -> keep pointer semantics
            system_license_autoupdate_url_model.Name =
                types.StringPointerValue(v_system_license_autoupdate_url.Name)

            system_license_autoupdate_url_List[i_system_license_autoupdate_url] =
                system_license_autoupdate_url_model
        }

        // Write url field as Null when empty, else concrete list
        if len(system_license_autoupdate_url_List) == 0 {
            system_license_autoupdate_model.Url =
                types.ListNull(types.ObjectType{AttrTypes: System_License_Autoupdate_Url_Model{}.AttrTypes()})
        } else {
            system_license_autoupdate_model.Url, _ =
                types.ListValueFrom(ctx,
                    types.ObjectType{AttrTypes: System_License_Autoupdate_Url_Model{}.AttrTypes()},
                    system_license_autoupdate_url_List,
                )
        }
        system_license_autoupdate_List[i_system_license_autoupdate] = system_license_autoupdate_model
        }

        // Write autoupdate field as Null when empty, else concrete list
        if len(system_license_autoupdate_List) == 0 {
            system_license_model.Autoupdate =
                types.ListNull(types.ObjectType{AttrTypes: System_License_Autoupdate_Model{}.AttrTypes()})
        } else {
            system_license_model.Autoupdate, _ =
                types.ListValueFrom(ctx,
                    types.ObjectType{AttrTypes: System_License_Autoupdate_Model{}.AttrTypes()},
                    system_license_autoupdate_List,
                )
        }
        system_license_List[i_system_license] = system_license_model
        }

        // Write license field as Null when empty, else concrete list
        if len(system_license_List) == 0 {
            system_model.License =
                types.ListNull(types.ObjectType{AttrTypes: System_License_Model{}.AttrTypes()})
        } else {
            system_model.License, _ =
                types.ListValueFrom(ctx,
                    types.ObjectType{AttrTypes: System_License_Model{}.AttrTypes()},
                    system_license_List,
                )
        }
        system_List[i_system] = system_model

        system_List[i_system] = system_model
    }

    // Write parent list as Null when empty
    if len(system_List) == 0 {
        state.System =
            types.ListNull(types.ObjectType{AttrTypes: System_Model{}.AttrTypes()})
    } else {
        state.System, _ =
            types.ListValueFrom(ctx,
                types.ObjectType{AttrTypes: System_Model{}.AttrTypes()},
                system_List,
            )
    }

    // Persist final state once after all parents processed
    resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}





// Update implements resource.Resource.
func (r *resource_Apply_Groups) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	
	var plan Groups_Model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	// Check for errors
	if resp.Diagnostics.HasError() {
		return
	}
	var config xml_Configuration
	config.Groups.Name = plan.ResourceName.ValueStringPointer()
    
	
    var var_chassis []Chassis_Model
    if plan.Chassis.IsNull() {
        var_chassis = []Chassis_Model{}
    }else {
        resp.Diagnostics.Append(plan.Chassis.ElementsAs(ctx, &var_chassis, false)...)
        if resp.Diagnostics.HasError() {
            return
        }
    }
    config.Groups.Chassis = make([]xml_Chassis, len(var_chassis))
    
    for i_chassis, v_chassis := range var_chassis {
        var var_chassis_aggregated_devices []Chassis_Aggregated_devices_Model
        resp.Diagnostics.Append(v_chassis.Aggregated_devices.ElementsAs(ctx, &var_chassis_aggregated_devices, false)...)
        if resp.Diagnostics.HasError() {
            return
        }
	    config.Groups.Chassis[i_chassis].Aggregated_devices = make([]xml_Chassis_Aggregated_devices, len(var_chassis_aggregated_devices))
        
		for i_chassis_aggregated_devices, v_chassis_aggregated_devices := range var_chassis_aggregated_devices {
            var var_chassis_aggregated_devices_ethernet []Chassis_Aggregated_devices_Ethernet_Model
            resp.Diagnostics.Append(v_chassis_aggregated_devices.Ethernet.ElementsAs(ctx, &var_chassis_aggregated_devices_ethernet, false)...)
            if resp.Diagnostics.HasError() {
                return
            }
	    config.Groups.Chassis[i_chassis].Aggregated_devices[i_chassis_aggregated_devices].Ethernet = make([]xml_Chassis_Aggregated_devices_Ethernet, len(var_chassis_aggregated_devices_ethernet))
        
		for i_chassis_aggregated_devices_ethernet, v_chassis_aggregated_devices_ethernet := range var_chassis_aggregated_devices_ethernet {
            config.Groups.Chassis[i_chassis].Aggregated_devices[i_chassis_aggregated_devices].Ethernet[i_chassis_aggregated_devices_ethernet].Device_count = v_chassis_aggregated_devices_ethernet.Device_count.ValueStringPointer()
        }
        }
    }
	
    var var_interfaces []Interfaces_Model
    if plan.Interfaces.IsNull() {
        var_interfaces = []Interfaces_Model{}
    }else {
        resp.Diagnostics.Append(plan.Interfaces.ElementsAs(ctx, &var_interfaces, false)...)
        if resp.Diagnostics.HasError() {
            return
        }
    }
    config.Groups.Interfaces = make([]xml_Interfaces, len(var_interfaces))
    
    for i_interfaces, v_interfaces := range var_interfaces {
        var var_interfaces_interface []Interfaces_Interface_Model
        resp.Diagnostics.Append(v_interfaces.Interface.ElementsAs(ctx, &var_interfaces_interface, false)...)
        if resp.Diagnostics.HasError() {
            return
        }
	    config.Groups.Interfaces[i_interfaces].Interface = make([]xml_Interfaces_Interface, len(var_interfaces_interface))
        
		for i_interfaces_interface, v_interfaces_interface := range var_interfaces_interface {
            config.Groups.Interfaces[i_interfaces].Interface[i_interfaces_interface].Name = v_interfaces_interface.Name.ValueStringPointer()
            config.Groups.Interfaces[i_interfaces].Interface[i_interfaces_interface].Vlan_tagging = v_interfaces_interface.Vlan_tagging.ValueStringPointer()
            var var_interfaces_interface_unit []Interfaces_Interface_Unit_Model
            resp.Diagnostics.Append(v_interfaces_interface.Unit.ElementsAs(ctx, &var_interfaces_interface_unit, false)...)
            if resp.Diagnostics.HasError() {
                return
            }
	    config.Groups.Interfaces[i_interfaces].Interface[i_interfaces_interface].Unit = make([]xml_Interfaces_Interface_Unit, len(var_interfaces_interface_unit))
        
		for i_interfaces_interface_unit, v_interfaces_interface_unit := range var_interfaces_interface_unit {
            config.Groups.Interfaces[i_interfaces].Interface[i_interfaces_interface].Unit[i_interfaces_interface_unit].Name = v_interfaces_interface_unit.Name.ValueStringPointer()
            config.Groups.Interfaces[i_interfaces].Interface[i_interfaces_interface].Unit[i_interfaces_interface_unit].Description = v_interfaces_interface_unit.Description.ValueStringPointer()
            config.Groups.Interfaces[i_interfaces].Interface[i_interfaces_interface].Unit[i_interfaces_interface_unit].Vlan_id = v_interfaces_interface_unit.Vlan_id.ValueStringPointer()
            var var_interfaces_interface_unit_family []Interfaces_Interface_Unit_Family_Model
            resp.Diagnostics.Append(v_interfaces_interface_unit.Family.ElementsAs(ctx, &var_interfaces_interface_unit_family, false)...)
            if resp.Diagnostics.HasError() {
                return
            }
	    config.Groups.Interfaces[i_interfaces].Interface[i_interfaces_interface].Unit[i_interfaces_interface_unit].Family = make([]xml_Interfaces_Interface_Unit_Family, len(var_interfaces_interface_unit_family))
        
		for i_interfaces_interface_unit_family, v_interfaces_interface_unit_family := range var_interfaces_interface_unit_family {
            var var_interfaces_interface_unit_family_inet []Interfaces_Interface_Unit_Family_Inet_Model
            resp.Diagnostics.Append(v_interfaces_interface_unit_family.Inet.ElementsAs(ctx, &var_interfaces_interface_unit_family_inet, false)...)
            if resp.Diagnostics.HasError() {
                return
            }
	    config.Groups.Interfaces[i_interfaces].Interface[i_interfaces_interface].Unit[i_interfaces_interface_unit].Family[i_interfaces_interface_unit_family].Inet = make([]xml_Interfaces_Interface_Unit_Family_Inet, len(var_interfaces_interface_unit_family_inet))
        
		for i_interfaces_interface_unit_family_inet, v_interfaces_interface_unit_family_inet := range var_interfaces_interface_unit_family_inet {
            var var_interfaces_interface_unit_family_inet_address []Interfaces_Interface_Unit_Family_Inet_Address_Model
            resp.Diagnostics.Append(v_interfaces_interface_unit_family_inet.Address.ElementsAs(ctx, &var_interfaces_interface_unit_family_inet_address, false)...)
            if resp.Diagnostics.HasError() {
                return
            }
	    config.Groups.Interfaces[i_interfaces].Interface[i_interfaces_interface].Unit[i_interfaces_interface_unit].Family[i_interfaces_interface_unit_family].Inet[i_interfaces_interface_unit_family_inet].Address = make([]xml_Interfaces_Interface_Unit_Family_Inet_Address, len(var_interfaces_interface_unit_family_inet_address))
        
		for i_interfaces_interface_unit_family_inet_address, v_interfaces_interface_unit_family_inet_address := range var_interfaces_interface_unit_family_inet_address {
            config.Groups.Interfaces[i_interfaces].Interface[i_interfaces_interface].Unit[i_interfaces_interface_unit].Family[i_interfaces_interface_unit_family].Inet[i_interfaces_interface_unit_family_inet].Address[i_interfaces_interface_unit_family_inet_address].Name = v_interfaces_interface_unit_family_inet_address.Name.ValueStringPointer()
        }
        }
        }
        }
        }
    }
	
    var var_policy_options []Policy_options_Model
    if plan.Policy_options.IsNull() {
        var_policy_options = []Policy_options_Model{}
    }else {
        resp.Diagnostics.Append(plan.Policy_options.ElementsAs(ctx, &var_policy_options, false)...)
        if resp.Diagnostics.HasError() {
            return
        }
    }
    config.Groups.Policy_options = make([]xml_Policy_options, len(var_policy_options))
    
    for i_policy_options, v_policy_options := range var_policy_options {
        var var_policy_options_policy_statement []Policy_options_Policy_statement_Model
        resp.Diagnostics.Append(v_policy_options.Policy_statement.ElementsAs(ctx, &var_policy_options_policy_statement, false)...)
        if resp.Diagnostics.HasError() {
            return
        }
	    config.Groups.Policy_options[i_policy_options].Policy_statement = make([]xml_Policy_options_Policy_statement, len(var_policy_options_policy_statement))
        
		for i_policy_options_policy_statement, v_policy_options_policy_statement := range var_policy_options_policy_statement {
            config.Groups.Policy_options[i_policy_options].Policy_statement[i_policy_options_policy_statement].Name = v_policy_options_policy_statement.Name.ValueStringPointer()
            var var_policy_options_policy_statement_term []Policy_options_Policy_statement_Term_Model
            resp.Diagnostics.Append(v_policy_options_policy_statement.Term.ElementsAs(ctx, &var_policy_options_policy_statement_term, false)...)
            if resp.Diagnostics.HasError() {
                return
            }
	    config.Groups.Policy_options[i_policy_options].Policy_statement[i_policy_options_policy_statement].Term = make([]xml_Policy_options_Policy_statement_Term, len(var_policy_options_policy_statement_term))
        
		for i_policy_options_policy_statement_term, v_policy_options_policy_statement_term := range var_policy_options_policy_statement_term {
            config.Groups.Policy_options[i_policy_options].Policy_statement[i_policy_options_policy_statement].Term[i_policy_options_policy_statement_term].Name = v_policy_options_policy_statement_term.Name.ValueStringPointer()
            var var_policy_options_policy_statement_term_from []Policy_options_Policy_statement_Term_From_Model
            resp.Diagnostics.Append(v_policy_options_policy_statement_term.From.ElementsAs(ctx, &var_policy_options_policy_statement_term_from, false)...)
            if resp.Diagnostics.HasError() {
                return
            }
	    config.Groups.Policy_options[i_policy_options].Policy_statement[i_policy_options_policy_statement].Term[i_policy_options_policy_statement_term].From = make([]xml_Policy_options_Policy_statement_Term_From, len(var_policy_options_policy_statement_term_from))
        
		for i_policy_options_policy_statement_term_from, v_policy_options_policy_statement_term_from := range var_policy_options_policy_statement_term_from {
            var var_policy_options_policy_statement_term_from_route_filter []Policy_options_Policy_statement_Term_From_Route_filter_Model
            resp.Diagnostics.Append(v_policy_options_policy_statement_term_from.Route_filter.ElementsAs(ctx, &var_policy_options_policy_statement_term_from_route_filter, false)...)
            if resp.Diagnostics.HasError() {
                return
            }
	    config.Groups.Policy_options[i_policy_options].Policy_statement[i_policy_options_policy_statement].Term[i_policy_options_policy_statement_term].From[i_policy_options_policy_statement_term_from].Route_filter = make([]xml_Policy_options_Policy_statement_Term_From_Route_filter, len(var_policy_options_policy_statement_term_from_route_filter))
        
		for i_policy_options_policy_statement_term_from_route_filter, v_policy_options_policy_statement_term_from_route_filter := range var_policy_options_policy_statement_term_from_route_filter {
            config.Groups.Policy_options[i_policy_options].Policy_statement[i_policy_options_policy_statement].Term[i_policy_options_policy_statement_term].From[i_policy_options_policy_statement_term_from].Route_filter[i_policy_options_policy_statement_term_from_route_filter].Address = v_policy_options_policy_statement_term_from_route_filter.Address.ValueStringPointer()
            config.Groups.Policy_options[i_policy_options].Policy_statement[i_policy_options_policy_statement].Term[i_policy_options_policy_statement_term].From[i_policy_options_policy_statement_term_from].Route_filter[i_policy_options_policy_statement_term_from_route_filter].Exact = v_policy_options_policy_statement_term_from_route_filter.Exact.ValueStringPointer()
            config.Groups.Policy_options[i_policy_options].Policy_statement[i_policy_options_policy_statement].Term[i_policy_options_policy_statement_term].From[i_policy_options_policy_statement_term_from].Route_filter[i_policy_options_policy_statement_term_from_route_filter].Accept = v_policy_options_policy_statement_term_from_route_filter.Accept.ValueStringPointer()
        }
        }
        }
        }
    }
	
    var var_protocols []Protocols_Model
    if plan.Protocols.IsNull() {
        var_protocols = []Protocols_Model{}
    }else {
        resp.Diagnostics.Append(plan.Protocols.ElementsAs(ctx, &var_protocols, false)...)
        if resp.Diagnostics.HasError() {
            return
        }
    }
    config.Groups.Protocols = make([]xml_Protocols, len(var_protocols))
    
    for i_protocols, v_protocols := range var_protocols {
        var var_protocols_lldp []Protocols_Lldp_Model
        resp.Diagnostics.Append(v_protocols.Lldp.ElementsAs(ctx, &var_protocols_lldp, false)...)
        if resp.Diagnostics.HasError() {
            return
        }
	    config.Groups.Protocols[i_protocols].Lldp = make([]xml_Protocols_Lldp, len(var_protocols_lldp))
        
		for i_protocols_lldp, v_protocols_lldp := range var_protocols_lldp {
            var var_protocols_lldp_interface []Protocols_Lldp_Interface_Model
            resp.Diagnostics.Append(v_protocols_lldp.Interface.ElementsAs(ctx, &var_protocols_lldp_interface, false)...)
            if resp.Diagnostics.HasError() {
                return
            }
	    config.Groups.Protocols[i_protocols].Lldp[i_protocols_lldp].Interface = make([]xml_Protocols_Lldp_Interface, len(var_protocols_lldp_interface))
        
		for i_protocols_lldp_interface, v_protocols_lldp_interface := range var_protocols_lldp_interface {
            config.Groups.Protocols[i_protocols].Lldp[i_protocols_lldp].Interface[i_protocols_lldp_interface].Name = v_protocols_lldp_interface.Name.ValueStringPointer()
        }
        }
    }
	
    var var_routing_instances []Routing_instances_Model
    if plan.Routing_instances.IsNull() {
        var_routing_instances = []Routing_instances_Model{}
    }else {
        resp.Diagnostics.Append(plan.Routing_instances.ElementsAs(ctx, &var_routing_instances, false)...)
        if resp.Diagnostics.HasError() {
            return
        }
    }
    config.Groups.Routing_instances = make([]xml_Routing_instances, len(var_routing_instances))
    
    for i_routing_instances, v_routing_instances := range var_routing_instances {
        var var_routing_instances_instance []Routing_instances_Instance_Model
        resp.Diagnostics.Append(v_routing_instances.Instance.ElementsAs(ctx, &var_routing_instances_instance, false)...)
        if resp.Diagnostics.HasError() {
            return
        }
	    config.Groups.Routing_instances[i_routing_instances].Instance = make([]xml_Routing_instances_Instance, len(var_routing_instances_instance))
        
		for i_routing_instances_instance, v_routing_instances_instance := range var_routing_instances_instance {
            config.Groups.Routing_instances[i_routing_instances].Instance[i_routing_instances_instance].Name = v_routing_instances_instance.Name.ValueStringPointer()
            config.Groups.Routing_instances[i_routing_instances].Instance[i_routing_instances_instance].Instance_type = v_routing_instances_instance.Instance_type.ValueStringPointer()
            var var_routing_instances_instance_interface []Routing_instances_Instance_Interface_Model
            resp.Diagnostics.Append(v_routing_instances_instance.Interface.ElementsAs(ctx, &var_routing_instances_instance_interface, false)...)
            if resp.Diagnostics.HasError() {
                return
            }
	    config.Groups.Routing_instances[i_routing_instances].Instance[i_routing_instances_instance].Interface = make([]xml_Routing_instances_Instance_Interface, len(var_routing_instances_instance_interface))
        
		for i_routing_instances_instance_interface, v_routing_instances_instance_interface := range var_routing_instances_instance_interface {
            config.Groups.Routing_instances[i_routing_instances].Instance[i_routing_instances_instance].Interface[i_routing_instances_instance_interface].Name = v_routing_instances_instance_interface.Name.ValueStringPointer()
        }
            var var_routing_instances_instance_routing_options []Routing_instances_Instance_Routing_options_Model
            resp.Diagnostics.Append(v_routing_instances_instance.Routing_options.ElementsAs(ctx, &var_routing_instances_instance_routing_options, false)...)
            if resp.Diagnostics.HasError() {
                return
            }
	    config.Groups.Routing_instances[i_routing_instances].Instance[i_routing_instances_instance].Routing_options = make([]xml_Routing_instances_Instance_Routing_options, len(var_routing_instances_instance_routing_options))
        
		for i_routing_instances_instance_routing_options, v_routing_instances_instance_routing_options := range var_routing_instances_instance_routing_options {
            var var_routing_instances_instance_routing_options_static []Routing_instances_Instance_Routing_options_Static_Model
            resp.Diagnostics.Append(v_routing_instances_instance_routing_options.Static.ElementsAs(ctx, &var_routing_instances_instance_routing_options_static, false)...)
            if resp.Diagnostics.HasError() {
                return
            }
	    config.Groups.Routing_instances[i_routing_instances].Instance[i_routing_instances_instance].Routing_options[i_routing_instances_instance_routing_options].Static = make([]xml_Routing_instances_Instance_Routing_options_Static, len(var_routing_instances_instance_routing_options_static))
        
		for i_routing_instances_instance_routing_options_static, v_routing_instances_instance_routing_options_static := range var_routing_instances_instance_routing_options_static {
            var var_routing_instances_instance_routing_options_static_route []Routing_instances_Instance_Routing_options_Static_Route_Model
            resp.Diagnostics.Append(v_routing_instances_instance_routing_options_static.Route.ElementsAs(ctx, &var_routing_instances_instance_routing_options_static_route, false)...)
            if resp.Diagnostics.HasError() {
                return
            }
	    config.Groups.Routing_instances[i_routing_instances].Instance[i_routing_instances_instance].Routing_options[i_routing_instances_instance_routing_options].Static[i_routing_instances_instance_routing_options_static].Route = make([]xml_Routing_instances_Instance_Routing_options_Static_Route, len(var_routing_instances_instance_routing_options_static_route))
        
		for i_routing_instances_instance_routing_options_static_route, v_routing_instances_instance_routing_options_static_route := range var_routing_instances_instance_routing_options_static_route {
            config.Groups.Routing_instances[i_routing_instances].Instance[i_routing_instances_instance].Routing_options[i_routing_instances_instance_routing_options].Static[i_routing_instances_instance_routing_options_static].Route[i_routing_instances_instance_routing_options_static_route].Name = v_routing_instances_instance_routing_options_static_route.Name.ValueStringPointer()
            config.Groups.Routing_instances[i_routing_instances].Instance[i_routing_instances_instance].Routing_options[i_routing_instances_instance_routing_options].Static[i_routing_instances_instance_routing_options_static].Route[i_routing_instances_instance_routing_options_static_route].Discard = v_routing_instances_instance_routing_options_static_route.Discard.ValueStringPointer()
        }
        }
        }
            var var_routing_instances_instance_protocols []Routing_instances_Instance_Protocols_Model
            resp.Diagnostics.Append(v_routing_instances_instance.Protocols.ElementsAs(ctx, &var_routing_instances_instance_protocols, false)...)
            if resp.Diagnostics.HasError() {
                return
            }
	    config.Groups.Routing_instances[i_routing_instances].Instance[i_routing_instances_instance].Protocols = make([]xml_Routing_instances_Instance_Protocols, len(var_routing_instances_instance_protocols))
        
		for i_routing_instances_instance_protocols, v_routing_instances_instance_protocols := range var_routing_instances_instance_protocols {
            var var_routing_instances_instance_protocols_ospf []Routing_instances_Instance_Protocols_Ospf_Model
            resp.Diagnostics.Append(v_routing_instances_instance_protocols.Ospf.ElementsAs(ctx, &var_routing_instances_instance_protocols_ospf, false)...)
            if resp.Diagnostics.HasError() {
                return
            }
	    config.Groups.Routing_instances[i_routing_instances].Instance[i_routing_instances_instance].Protocols[i_routing_instances_instance_protocols].Ospf = make([]xml_Routing_instances_Instance_Protocols_Ospf, len(var_routing_instances_instance_protocols_ospf))
        
		for i_routing_instances_instance_protocols_ospf, v_routing_instances_instance_protocols_ospf := range var_routing_instances_instance_protocols_ospf {
			var var_routing_instances_instance_protocols_ospf_export []string
			resp.Diagnostics.Append(v_routing_instances_instance_protocols_ospf.Export.ElementsAs(ctx, &var_routing_instances_instance_protocols_ospf_export, false)...)
			for _, v_routing_instances_instance_protocols_ospf_export := range var_routing_instances_instance_protocols_ospf_export {
				config.Groups.Routing_instances[i_routing_instances].Instance[i_routing_instances_instance].Protocols[i_routing_instances_instance_protocols].Ospf[i_routing_instances_instance_protocols_ospf].Export = append(config.Groups.Routing_instances[i_routing_instances].Instance[i_routing_instances_instance].Protocols[i_routing_instances_instance_protocols].Ospf[i_routing_instances_instance_protocols_ospf].Export, &v_routing_instances_instance_protocols_ospf_export)
			}
            var var_routing_instances_instance_protocols_ospf_area []Routing_instances_Instance_Protocols_Ospf_Area_Model
            resp.Diagnostics.Append(v_routing_instances_instance_protocols_ospf.Area.ElementsAs(ctx, &var_routing_instances_instance_protocols_ospf_area, false)...)
            if resp.Diagnostics.HasError() {
                return
            }
	    config.Groups.Routing_instances[i_routing_instances].Instance[i_routing_instances_instance].Protocols[i_routing_instances_instance_protocols].Ospf[i_routing_instances_instance_protocols_ospf].Area = make([]xml_Routing_instances_Instance_Protocols_Ospf_Area, len(var_routing_instances_instance_protocols_ospf_area))
        
		for i_routing_instances_instance_protocols_ospf_area, v_routing_instances_instance_protocols_ospf_area := range var_routing_instances_instance_protocols_ospf_area {
            config.Groups.Routing_instances[i_routing_instances].Instance[i_routing_instances_instance].Protocols[i_routing_instances_instance_protocols].Ospf[i_routing_instances_instance_protocols_ospf].Area[i_routing_instances_instance_protocols_ospf_area].Name = v_routing_instances_instance_protocols_ospf_area.Name.ValueStringPointer()
            var var_routing_instances_instance_protocols_ospf_area_interface []Routing_instances_Instance_Protocols_Ospf_Area_Interface_Model
            resp.Diagnostics.Append(v_routing_instances_instance_protocols_ospf_area.Interface.ElementsAs(ctx, &var_routing_instances_instance_protocols_ospf_area_interface, false)...)
            if resp.Diagnostics.HasError() {
                return
            }
	    config.Groups.Routing_instances[i_routing_instances].Instance[i_routing_instances_instance].Protocols[i_routing_instances_instance_protocols].Ospf[i_routing_instances_instance_protocols_ospf].Area[i_routing_instances_instance_protocols_ospf_area].Interface = make([]xml_Routing_instances_Instance_Protocols_Ospf_Area_Interface, len(var_routing_instances_instance_protocols_ospf_area_interface))
        
		for i_routing_instances_instance_protocols_ospf_area_interface, v_routing_instances_instance_protocols_ospf_area_interface := range var_routing_instances_instance_protocols_ospf_area_interface {
            config.Groups.Routing_instances[i_routing_instances].Instance[i_routing_instances_instance].Protocols[i_routing_instances_instance_protocols].Ospf[i_routing_instances_instance_protocols_ospf].Area[i_routing_instances_instance_protocols_ospf_area].Interface[i_routing_instances_instance_protocols_ospf_area_interface].Name = v_routing_instances_instance_protocols_ospf_area_interface.Name.ValueStringPointer()
        }
        }
        }
        }
        }
    }
	
    var var_routing_options []Routing_options_Model
    if plan.Routing_options.IsNull() {
        var_routing_options = []Routing_options_Model{}
    }else {
        resp.Diagnostics.Append(plan.Routing_options.ElementsAs(ctx, &var_routing_options, false)...)
        if resp.Diagnostics.HasError() {
            return
        }
    }
    config.Groups.Routing_options = make([]xml_Routing_options, len(var_routing_options))
    
    for i_routing_options, v_routing_options := range var_routing_options {
        var var_routing_options_static []Routing_options_Static_Model
        resp.Diagnostics.Append(v_routing_options.Static.ElementsAs(ctx, &var_routing_options_static, false)...)
        if resp.Diagnostics.HasError() {
            return
        }
	    config.Groups.Routing_options[i_routing_options].Static = make([]xml_Routing_options_Static, len(var_routing_options_static))
        
		for i_routing_options_static, v_routing_options_static := range var_routing_options_static {
            var var_routing_options_static_route []Routing_options_Static_Route_Model
            resp.Diagnostics.Append(v_routing_options_static.Route.ElementsAs(ctx, &var_routing_options_static_route, false)...)
            if resp.Diagnostics.HasError() {
                return
            }
	    config.Groups.Routing_options[i_routing_options].Static[i_routing_options_static].Route = make([]xml_Routing_options_Static_Route, len(var_routing_options_static_route))
        
		for i_routing_options_static_route, v_routing_options_static_route := range var_routing_options_static_route {
            config.Groups.Routing_options[i_routing_options].Static[i_routing_options_static].Route[i_routing_options_static_route].Name = v_routing_options_static_route.Name.ValueStringPointer()
			var var_routing_options_static_route_next_hop []string
			resp.Diagnostics.Append(v_routing_options_static_route.Next_hop.ElementsAs(ctx, &var_routing_options_static_route_next_hop, false)...)
			for _, v_routing_options_static_route_next_hop := range var_routing_options_static_route_next_hop {
				config.Groups.Routing_options[i_routing_options].Static[i_routing_options_static].Route[i_routing_options_static_route].Next_hop = append(config.Groups.Routing_options[i_routing_options].Static[i_routing_options_static].Route[i_routing_options_static_route].Next_hop, &v_routing_options_static_route_next_hop)
			}
        }
        }
    }
	
    var var_security []Security_Model
    if plan.Security.IsNull() {
        var_security = []Security_Model{}
    }else {
        resp.Diagnostics.Append(plan.Security.ElementsAs(ctx, &var_security, false)...)
        if resp.Diagnostics.HasError() {
            return
        }
    }
    config.Groups.Security = make([]xml_Security, len(var_security))
    
    for i_security, v_security := range var_security {
        var var_security_log []Security_Log_Model
        resp.Diagnostics.Append(v_security.Log.ElementsAs(ctx, &var_security_log, false)...)
        if resp.Diagnostics.HasError() {
            return
        }
	    config.Groups.Security[i_security].Log = make([]xml_Security_Log, len(var_security_log))
        
		for i_security_log, v_security_log := range var_security_log {
            config.Groups.Security[i_security].Log[i_security_log].Mode = v_security_log.Mode.ValueStringPointer()
            var var_security_log_report []Security_Log_Report_Model
            resp.Diagnostics.Append(v_security_log.Report.ElementsAs(ctx, &var_security_log_report, false)...)
            if resp.Diagnostics.HasError() {
                return
            }
	    config.Groups.Security[i_security].Log[i_security_log].Report = make([]xml_Security_Log_Report, len(var_security_log_report))
        
        }
        var var_security_screen []Security_Screen_Model
        resp.Diagnostics.Append(v_security.Screen.ElementsAs(ctx, &var_security_screen, false)...)
        if resp.Diagnostics.HasError() {
            return
        }
	    config.Groups.Security[i_security].Screen = make([]xml_Security_Screen, len(var_security_screen))
        
		for i_security_screen, v_security_screen := range var_security_screen {
            var var_security_screen_ids_option []Security_Screen_Ids_option_Model
            resp.Diagnostics.Append(v_security_screen.Ids_option.ElementsAs(ctx, &var_security_screen_ids_option, false)...)
            if resp.Diagnostics.HasError() {
                return
            }
	    config.Groups.Security[i_security].Screen[i_security_screen].Ids_option = make([]xml_Security_Screen_Ids_option, len(var_security_screen_ids_option))
        
		for i_security_screen_ids_option, v_security_screen_ids_option := range var_security_screen_ids_option {
            config.Groups.Security[i_security].Screen[i_security_screen].Ids_option[i_security_screen_ids_option].Name = v_security_screen_ids_option.Name.ValueStringPointer()
            var var_security_screen_ids_option_icmp []Security_Screen_Ids_option_Icmp_Model
            resp.Diagnostics.Append(v_security_screen_ids_option.Icmp.ElementsAs(ctx, &var_security_screen_ids_option_icmp, false)...)
            if resp.Diagnostics.HasError() {
                return
            }
	    config.Groups.Security[i_security].Screen[i_security_screen].Ids_option[i_security_screen_ids_option].Icmp = make([]xml_Security_Screen_Ids_option_Icmp, len(var_security_screen_ids_option_icmp))
        
		for i_security_screen_ids_option_icmp, v_security_screen_ids_option_icmp := range var_security_screen_ids_option_icmp {
            config.Groups.Security[i_security].Screen[i_security_screen].Ids_option[i_security_screen_ids_option].Icmp[i_security_screen_ids_option_icmp].Ping_death = v_security_screen_ids_option_icmp.Ping_death.ValueStringPointer()
        }
            var var_security_screen_ids_option_ip []Security_Screen_Ids_option_Ip_Model
            resp.Diagnostics.Append(v_security_screen_ids_option.Ip.ElementsAs(ctx, &var_security_screen_ids_option_ip, false)...)
            if resp.Diagnostics.HasError() {
                return
            }
	    config.Groups.Security[i_security].Screen[i_security_screen].Ids_option[i_security_screen_ids_option].Ip = make([]xml_Security_Screen_Ids_option_Ip, len(var_security_screen_ids_option_ip))
        
		for i_security_screen_ids_option_ip, v_security_screen_ids_option_ip := range var_security_screen_ids_option_ip {
            config.Groups.Security[i_security].Screen[i_security_screen].Ids_option[i_security_screen_ids_option].Ip[i_security_screen_ids_option_ip].Source_route_option = v_security_screen_ids_option_ip.Source_route_option.ValueStringPointer()
            config.Groups.Security[i_security].Screen[i_security_screen].Ids_option[i_security_screen_ids_option].Ip[i_security_screen_ids_option_ip].Tear_drop = v_security_screen_ids_option_ip.Tear_drop.ValueStringPointer()
        }
            var var_security_screen_ids_option_tcp []Security_Screen_Ids_option_Tcp_Model
            resp.Diagnostics.Append(v_security_screen_ids_option.Tcp.ElementsAs(ctx, &var_security_screen_ids_option_tcp, false)...)
            if resp.Diagnostics.HasError() {
                return
            }
	    config.Groups.Security[i_security].Screen[i_security_screen].Ids_option[i_security_screen_ids_option].Tcp = make([]xml_Security_Screen_Ids_option_Tcp, len(var_security_screen_ids_option_tcp))
        
		for i_security_screen_ids_option_tcp, v_security_screen_ids_option_tcp := range var_security_screen_ids_option_tcp {
            var var_security_screen_ids_option_tcp_syn_flood []Security_Screen_Ids_option_Tcp_Syn_flood_Model
            resp.Diagnostics.Append(v_security_screen_ids_option_tcp.Syn_flood.ElementsAs(ctx, &var_security_screen_ids_option_tcp_syn_flood, false)...)
            if resp.Diagnostics.HasError() {
                return
            }
	    config.Groups.Security[i_security].Screen[i_security_screen].Ids_option[i_security_screen_ids_option].Tcp[i_security_screen_ids_option_tcp].Syn_flood = make([]xml_Security_Screen_Ids_option_Tcp_Syn_flood, len(var_security_screen_ids_option_tcp_syn_flood))
        
		for i_security_screen_ids_option_tcp_syn_flood, v_security_screen_ids_option_tcp_syn_flood := range var_security_screen_ids_option_tcp_syn_flood {
            config.Groups.Security[i_security].Screen[i_security_screen].Ids_option[i_security_screen_ids_option].Tcp[i_security_screen_ids_option_tcp].Syn_flood[i_security_screen_ids_option_tcp_syn_flood].Alarm_threshold = v_security_screen_ids_option_tcp_syn_flood.Alarm_threshold.ValueStringPointer()
            config.Groups.Security[i_security].Screen[i_security_screen].Ids_option[i_security_screen_ids_option].Tcp[i_security_screen_ids_option_tcp].Syn_flood[i_security_screen_ids_option_tcp_syn_flood].Attack_threshold = v_security_screen_ids_option_tcp_syn_flood.Attack_threshold.ValueStringPointer()
            config.Groups.Security[i_security].Screen[i_security_screen].Ids_option[i_security_screen_ids_option].Tcp[i_security_screen_ids_option_tcp].Syn_flood[i_security_screen_ids_option_tcp_syn_flood].Source_threshold = v_security_screen_ids_option_tcp_syn_flood.Source_threshold.ValueStringPointer()
            config.Groups.Security[i_security].Screen[i_security_screen].Ids_option[i_security_screen_ids_option].Tcp[i_security_screen_ids_option_tcp].Syn_flood[i_security_screen_ids_option_tcp_syn_flood].Destination_threshold = v_security_screen_ids_option_tcp_syn_flood.Destination_threshold.ValueStringPointer()
            config.Groups.Security[i_security].Screen[i_security_screen].Ids_option[i_security_screen_ids_option].Tcp[i_security_screen_ids_option_tcp].Syn_flood[i_security_screen_ids_option_tcp_syn_flood].Timeout = v_security_screen_ids_option_tcp_syn_flood.Timeout.ValueStringPointer()
        }
            config.Groups.Security[i_security].Screen[i_security_screen].Ids_option[i_security_screen_ids_option].Tcp[i_security_screen_ids_option_tcp].Land = v_security_screen_ids_option_tcp.Land.ValueStringPointer()
        }
        }
        }
        var var_security_policies []Security_Policies_Model
        resp.Diagnostics.Append(v_security.Policies.ElementsAs(ctx, &var_security_policies, false)...)
        if resp.Diagnostics.HasError() {
            return
        }
	    config.Groups.Security[i_security].Policies = make([]xml_Security_Policies, len(var_security_policies))
        
		for i_security_policies, v_security_policies := range var_security_policies {
            var var_security_policies_policy []Security_Policies_Policy_Model
            resp.Diagnostics.Append(v_security_policies.Policy.ElementsAs(ctx, &var_security_policies_policy, false)...)
            if resp.Diagnostics.HasError() {
                return
            }
	    config.Groups.Security[i_security].Policies[i_security_policies].Policy = make([]xml_Security_Policies_Policy, len(var_security_policies_policy))
        
		for i_security_policies_policy, v_security_policies_policy := range var_security_policies_policy {
            config.Groups.Security[i_security].Policies[i_security_policies].Policy[i_security_policies_policy].From_zone_name = v_security_policies_policy.From_zone_name.ValueStringPointer()
            config.Groups.Security[i_security].Policies[i_security_policies].Policy[i_security_policies_policy].To_zone_name = v_security_policies_policy.To_zone_name.ValueStringPointer()
            var var_security_policies_policy_policy []Security_Policies_Policy_Policy_Model
            resp.Diagnostics.Append(v_security_policies_policy.Policy.ElementsAs(ctx, &var_security_policies_policy_policy, false)...)
            if resp.Diagnostics.HasError() {
                return
            }
	    config.Groups.Security[i_security].Policies[i_security_policies].Policy[i_security_policies_policy].Policy = make([]xml_Security_Policies_Policy_Policy, len(var_security_policies_policy_policy))
        
		for i_security_policies_policy_policy, v_security_policies_policy_policy := range var_security_policies_policy_policy {
            config.Groups.Security[i_security].Policies[i_security_policies].Policy[i_security_policies_policy].Policy[i_security_policies_policy_policy].Name = v_security_policies_policy_policy.Name.ValueStringPointer()
            var var_security_policies_policy_policy_match []Security_Policies_Policy_Policy_Match_Model
            resp.Diagnostics.Append(v_security_policies_policy_policy.Match.ElementsAs(ctx, &var_security_policies_policy_policy_match, false)...)
            if resp.Diagnostics.HasError() {
                return
            }
	    config.Groups.Security[i_security].Policies[i_security_policies].Policy[i_security_policies_policy].Policy[i_security_policies_policy_policy].Match = make([]xml_Security_Policies_Policy_Policy_Match, len(var_security_policies_policy_policy_match))
        
		for i_security_policies_policy_policy_match, v_security_policies_policy_policy_match := range var_security_policies_policy_policy_match {
			var var_security_policies_policy_policy_match_source_address []string
			resp.Diagnostics.Append(v_security_policies_policy_policy_match.Source_address.ElementsAs(ctx, &var_security_policies_policy_policy_match_source_address, false)...)
			for _, v_security_policies_policy_policy_match_source_address := range var_security_policies_policy_policy_match_source_address {
				config.Groups.Security[i_security].Policies[i_security_policies].Policy[i_security_policies_policy].Policy[i_security_policies_policy_policy].Match[i_security_policies_policy_policy_match].Source_address = append(config.Groups.Security[i_security].Policies[i_security_policies].Policy[i_security_policies_policy].Policy[i_security_policies_policy_policy].Match[i_security_policies_policy_policy_match].Source_address, &v_security_policies_policy_policy_match_source_address)
			}
			var var_security_policies_policy_policy_match_destination_address []string
			resp.Diagnostics.Append(v_security_policies_policy_policy_match.Destination_address.ElementsAs(ctx, &var_security_policies_policy_policy_match_destination_address, false)...)
			for _, v_security_policies_policy_policy_match_destination_address := range var_security_policies_policy_policy_match_destination_address {
				config.Groups.Security[i_security].Policies[i_security_policies].Policy[i_security_policies_policy].Policy[i_security_policies_policy_policy].Match[i_security_policies_policy_policy_match].Destination_address = append(config.Groups.Security[i_security].Policies[i_security_policies].Policy[i_security_policies_policy].Policy[i_security_policies_policy_policy].Match[i_security_policies_policy_policy_match].Destination_address, &v_security_policies_policy_policy_match_destination_address)
			}
			var var_security_policies_policy_policy_match_application []string
			resp.Diagnostics.Append(v_security_policies_policy_policy_match.Application.ElementsAs(ctx, &var_security_policies_policy_policy_match_application, false)...)
			for _, v_security_policies_policy_policy_match_application := range var_security_policies_policy_policy_match_application {
				config.Groups.Security[i_security].Policies[i_security_policies].Policy[i_security_policies_policy].Policy[i_security_policies_policy_policy].Match[i_security_policies_policy_policy_match].Application = append(config.Groups.Security[i_security].Policies[i_security_policies].Policy[i_security_policies_policy].Policy[i_security_policies_policy_policy].Match[i_security_policies_policy_policy_match].Application, &v_security_policies_policy_policy_match_application)
			}
        }
            var var_security_policies_policy_policy_then []Security_Policies_Policy_Policy_Then_Model
            resp.Diagnostics.Append(v_security_policies_policy_policy.Then.ElementsAs(ctx, &var_security_policies_policy_policy_then, false)...)
            if resp.Diagnostics.HasError() {
                return
            }
	    config.Groups.Security[i_security].Policies[i_security_policies].Policy[i_security_policies_policy].Policy[i_security_policies_policy_policy].Then = make([]xml_Security_Policies_Policy_Policy_Then, len(var_security_policies_policy_policy_then))
        
		for i_security_policies_policy_policy_then, v_security_policies_policy_policy_then := range var_security_policies_policy_policy_then {
            var var_security_policies_policy_policy_then_permit []Security_Policies_Policy_Policy_Then_Permit_Model
            resp.Diagnostics.Append(v_security_policies_policy_policy_then.Permit.ElementsAs(ctx, &var_security_policies_policy_policy_then_permit, false)...)
            if resp.Diagnostics.HasError() {
                return
            }
	    config.Groups.Security[i_security].Policies[i_security_policies].Policy[i_security_policies_policy].Policy[i_security_policies_policy_policy].Then[i_security_policies_policy_policy_then].Permit = make([]xml_Security_Policies_Policy_Policy_Then_Permit, len(var_security_policies_policy_policy_then_permit))
        
        }
        }
        }
        }
        var var_security_zones []Security_Zones_Model
        resp.Diagnostics.Append(v_security.Zones.ElementsAs(ctx, &var_security_zones, false)...)
        if resp.Diagnostics.HasError() {
            return
        }
	    config.Groups.Security[i_security].Zones = make([]xml_Security_Zones, len(var_security_zones))
        
		for i_security_zones, v_security_zones := range var_security_zones {
            var var_security_zones_security_zone []Security_Zones_Security_zone_Model
            resp.Diagnostics.Append(v_security_zones.Security_zone.ElementsAs(ctx, &var_security_zones_security_zone, false)...)
            if resp.Diagnostics.HasError() {
                return
            }
	    config.Groups.Security[i_security].Zones[i_security_zones].Security_zone = make([]xml_Security_Zones_Security_zone, len(var_security_zones_security_zone))
        
		for i_security_zones_security_zone, v_security_zones_security_zone := range var_security_zones_security_zone {
            config.Groups.Security[i_security].Zones[i_security_zones].Security_zone[i_security_zones_security_zone].Name = v_security_zones_security_zone.Name.ValueStringPointer()
            config.Groups.Security[i_security].Zones[i_security_zones].Security_zone[i_security_zones_security_zone].Tcp_rst = v_security_zones_security_zone.Tcp_rst.ValueStringPointer()
            config.Groups.Security[i_security].Zones[i_security_zones].Security_zone[i_security_zones_security_zone].Screen = v_security_zones_security_zone.Screen.ValueStringPointer()
            var var_security_zones_security_zone_host_inbound_traffic []Security_Zones_Security_zone_Host_inbound_traffic_Model
            resp.Diagnostics.Append(v_security_zones_security_zone.Host_inbound_traffic.ElementsAs(ctx, &var_security_zones_security_zone_host_inbound_traffic, false)...)
            if resp.Diagnostics.HasError() {
                return
            }
	    config.Groups.Security[i_security].Zones[i_security_zones].Security_zone[i_security_zones_security_zone].Host_inbound_traffic = make([]xml_Security_Zones_Security_zone_Host_inbound_traffic, len(var_security_zones_security_zone_host_inbound_traffic))
        
		for i_security_zones_security_zone_host_inbound_traffic, v_security_zones_security_zone_host_inbound_traffic := range var_security_zones_security_zone_host_inbound_traffic {
            var var_security_zones_security_zone_host_inbound_traffic_system_services []Security_Zones_Security_zone_Host_inbound_traffic_System_services_Model
            resp.Diagnostics.Append(v_security_zones_security_zone_host_inbound_traffic.System_services.ElementsAs(ctx, &var_security_zones_security_zone_host_inbound_traffic_system_services, false)...)
            if resp.Diagnostics.HasError() {
                return
            }
	    config.Groups.Security[i_security].Zones[i_security_zones].Security_zone[i_security_zones_security_zone].Host_inbound_traffic[i_security_zones_security_zone_host_inbound_traffic].System_services = make([]xml_Security_Zones_Security_zone_Host_inbound_traffic_System_services, len(var_security_zones_security_zone_host_inbound_traffic_system_services))
        
		for i_security_zones_security_zone_host_inbound_traffic_system_services, v_security_zones_security_zone_host_inbound_traffic_system_services := range var_security_zones_security_zone_host_inbound_traffic_system_services {
            config.Groups.Security[i_security].Zones[i_security_zones].Security_zone[i_security_zones_security_zone].Host_inbound_traffic[i_security_zones_security_zone_host_inbound_traffic].System_services[i_security_zones_security_zone_host_inbound_traffic_system_services].Name = v_security_zones_security_zone_host_inbound_traffic_system_services.Name.ValueStringPointer()
        }
            var var_security_zones_security_zone_host_inbound_traffic_protocols []Security_Zones_Security_zone_Host_inbound_traffic_Protocols_Model
            resp.Diagnostics.Append(v_security_zones_security_zone_host_inbound_traffic.Protocols.ElementsAs(ctx, &var_security_zones_security_zone_host_inbound_traffic_protocols, false)...)
            if resp.Diagnostics.HasError() {
                return
            }
	    config.Groups.Security[i_security].Zones[i_security_zones].Security_zone[i_security_zones_security_zone].Host_inbound_traffic[i_security_zones_security_zone_host_inbound_traffic].Protocols = make([]xml_Security_Zones_Security_zone_Host_inbound_traffic_Protocols, len(var_security_zones_security_zone_host_inbound_traffic_protocols))
        
		for i_security_zones_security_zone_host_inbound_traffic_protocols, v_security_zones_security_zone_host_inbound_traffic_protocols := range var_security_zones_security_zone_host_inbound_traffic_protocols {
            config.Groups.Security[i_security].Zones[i_security_zones].Security_zone[i_security_zones_security_zone].Host_inbound_traffic[i_security_zones_security_zone_host_inbound_traffic].Protocols[i_security_zones_security_zone_host_inbound_traffic_protocols].Name = v_security_zones_security_zone_host_inbound_traffic_protocols.Name.ValueStringPointer()
        }
        }
            var var_security_zones_security_zone_interfaces []Security_Zones_Security_zone_Interfaces_Model
            resp.Diagnostics.Append(v_security_zones_security_zone.Interfaces.ElementsAs(ctx, &var_security_zones_security_zone_interfaces, false)...)
            if resp.Diagnostics.HasError() {
                return
            }
	    config.Groups.Security[i_security].Zones[i_security_zones].Security_zone[i_security_zones_security_zone].Interfaces = make([]xml_Security_Zones_Security_zone_Interfaces, len(var_security_zones_security_zone_interfaces))
        
		for i_security_zones_security_zone_interfaces, v_security_zones_security_zone_interfaces := range var_security_zones_security_zone_interfaces {
            config.Groups.Security[i_security].Zones[i_security_zones].Security_zone[i_security_zones_security_zone].Interfaces[i_security_zones_security_zone_interfaces].Name = v_security_zones_security_zone_interfaces.Name.ValueStringPointer()
        }
        }
        }
    }
	
    var var_snmp []Snmp_Model
    if plan.Snmp.IsNull() {
        var_snmp = []Snmp_Model{}
    }else {
        resp.Diagnostics.Append(plan.Snmp.ElementsAs(ctx, &var_snmp, false)...)
        if resp.Diagnostics.HasError() {
            return
        }
    }
    config.Groups.Snmp = make([]xml_Snmp, len(var_snmp))
    
    for i_snmp, v_snmp := range var_snmp {
        config.Groups.Snmp[i_snmp].Location = v_snmp.Location.ValueStringPointer()
        config.Groups.Snmp[i_snmp].Contact = v_snmp.Contact.ValueStringPointer()
        var var_snmp_community []Snmp_Community_Model
        resp.Diagnostics.Append(v_snmp.Community.ElementsAs(ctx, &var_snmp_community, false)...)
        if resp.Diagnostics.HasError() {
            return
        }
	    config.Groups.Snmp[i_snmp].Community = make([]xml_Snmp_Community, len(var_snmp_community))
        
		for i_snmp_community, v_snmp_community := range var_snmp_community {
            config.Groups.Snmp[i_snmp].Community[i_snmp_community].Name = v_snmp_community.Name.ValueStringPointer()
            config.Groups.Snmp[i_snmp].Community[i_snmp_community].Authorization = v_snmp_community.Authorization.ValueStringPointer()
        }
    }
	
    var var_system []System_Model
    if plan.System.IsNull() {
        var_system = []System_Model{}
    }else {
        resp.Diagnostics.Append(plan.System.ElementsAs(ctx, &var_system, false)...)
        if resp.Diagnostics.HasError() {
            return
        }
    }
    config.Groups.System = make([]xml_System, len(var_system))
    
    for i_system, v_system := range var_system {
        var var_system_login []System_Login_Model
        resp.Diagnostics.Append(v_system.Login.ElementsAs(ctx, &var_system_login, false)...)
        if resp.Diagnostics.HasError() {
            return
        }
	    config.Groups.System[i_system].Login = make([]xml_System_Login, len(var_system_login))
        
		for i_system_login, v_system_login := range var_system_login {
            var var_system_login_user []System_Login_User_Model
            resp.Diagnostics.Append(v_system_login.User.ElementsAs(ctx, &var_system_login_user, false)...)
            if resp.Diagnostics.HasError() {
                return
            }
	    config.Groups.System[i_system].Login[i_system_login].User = make([]xml_System_Login_User, len(var_system_login_user))
        
		for i_system_login_user, v_system_login_user := range var_system_login_user {
            config.Groups.System[i_system].Login[i_system_login].User[i_system_login_user].Name = v_system_login_user.Name.ValueStringPointer()
            config.Groups.System[i_system].Login[i_system_login].User[i_system_login_user].Uid = v_system_login_user.Uid.ValueStringPointer()
            config.Groups.System[i_system].Login[i_system_login].User[i_system_login_user].Class = v_system_login_user.Class.ValueStringPointer()
            var var_system_login_user_authentication []System_Login_User_Authentication_Model
            resp.Diagnostics.Append(v_system_login_user.Authentication.ElementsAs(ctx, &var_system_login_user_authentication, false)...)
            if resp.Diagnostics.HasError() {
                return
            }
	    config.Groups.System[i_system].Login[i_system_login].User[i_system_login_user].Authentication = make([]xml_System_Login_User_Authentication, len(var_system_login_user_authentication))
        
		for i_system_login_user_authentication, v_system_login_user_authentication := range var_system_login_user_authentication {
            config.Groups.System[i_system].Login[i_system_login].User[i_system_login_user].Authentication[i_system_login_user_authentication].Encrypted_password = v_system_login_user_authentication.Encrypted_password.ValueStringPointer()
        }
        }
            config.Groups.System[i_system].Login[i_system_login].Message = v_system_login.Message.ValueStringPointer()
        }
        var var_system_root_authentication []System_Root_authentication_Model
        resp.Diagnostics.Append(v_system.Root_authentication.ElementsAs(ctx, &var_system_root_authentication, false)...)
        if resp.Diagnostics.HasError() {
            return
        }
	    config.Groups.System[i_system].Root_authentication = make([]xml_System_Root_authentication, len(var_system_root_authentication))
        
		for i_system_root_authentication, v_system_root_authentication := range var_system_root_authentication {
            config.Groups.System[i_system].Root_authentication[i_system_root_authentication].Encrypted_password = v_system_root_authentication.Encrypted_password.ValueStringPointer()
        }
        config.Groups.System[i_system].Host_name = v_system.Host_name.ValueStringPointer()
        var var_system_services []System_Services_Model
        resp.Diagnostics.Append(v_system.Services.ElementsAs(ctx, &var_system_services, false)...)
        if resp.Diagnostics.HasError() {
            return
        }
	    config.Groups.System[i_system].Services = make([]xml_System_Services, len(var_system_services))
        
		for i_system_services, v_system_services := range var_system_services {
            var var_system_services_ssh []System_Services_Ssh_Model
            resp.Diagnostics.Append(v_system_services.Ssh.ElementsAs(ctx, &var_system_services_ssh, false)...)
            if resp.Diagnostics.HasError() {
                return
            }
	    config.Groups.System[i_system].Services[i_system_services].Ssh = make([]xml_System_Services_Ssh, len(var_system_services_ssh))
        
		for i_system_services_ssh, v_system_services_ssh := range var_system_services_ssh {
            config.Groups.System[i_system].Services[i_system_services].Ssh[i_system_services_ssh].Root_login = v_system_services_ssh.Root_login.ValueStringPointer()
        }
            var var_system_services_extension_service []System_Services_Extension_service_Model
            resp.Diagnostics.Append(v_system_services.Extension_service.ElementsAs(ctx, &var_system_services_extension_service, false)...)
            if resp.Diagnostics.HasError() {
                return
            }
	    config.Groups.System[i_system].Services[i_system_services].Extension_service = make([]xml_System_Services_Extension_service, len(var_system_services_extension_service))
        
		for i_system_services_extension_service, v_system_services_extension_service := range var_system_services_extension_service {
            var var_system_services_extension_service_request_response []System_Services_Extension_service_Request_response_Model
            resp.Diagnostics.Append(v_system_services_extension_service.Request_response.ElementsAs(ctx, &var_system_services_extension_service_request_response, false)...)
            if resp.Diagnostics.HasError() {
                return
            }
	    config.Groups.System[i_system].Services[i_system_services].Extension_service[i_system_services_extension_service].Request_response = make([]xml_System_Services_Extension_service_Request_response, len(var_system_services_extension_service_request_response))
        
		for i_system_services_extension_service_request_response, v_system_services_extension_service_request_response := range var_system_services_extension_service_request_response {
            var var_system_services_extension_service_request_response_grpc []System_Services_Extension_service_Request_response_Grpc_Model
            resp.Diagnostics.Append(v_system_services_extension_service_request_response.Grpc.ElementsAs(ctx, &var_system_services_extension_service_request_response_grpc, false)...)
            if resp.Diagnostics.HasError() {
                return
            }
	    config.Groups.System[i_system].Services[i_system_services].Extension_service[i_system_services_extension_service].Request_response[i_system_services_extension_service_request_response].Grpc = make([]xml_System_Services_Extension_service_Request_response_Grpc, len(var_system_services_extension_service_request_response_grpc))
        
		for i_system_services_extension_service_request_response_grpc, v_system_services_extension_service_request_response_grpc := range var_system_services_extension_service_request_response_grpc {
            config.Groups.System[i_system].Services[i_system_services].Extension_service[i_system_services_extension_service].Request_response[i_system_services_extension_service_request_response].Grpc[i_system_services_extension_service_request_response_grpc].Max_connections = v_system_services_extension_service_request_response_grpc.Max_connections.ValueStringPointer()
        }
        }
            var var_system_services_extension_service_notification []System_Services_Extension_service_Notification_Model
            resp.Diagnostics.Append(v_system_services_extension_service.Notification.ElementsAs(ctx, &var_system_services_extension_service_notification, false)...)
            if resp.Diagnostics.HasError() {
                return
            }
	    config.Groups.System[i_system].Services[i_system_services].Extension_service[i_system_services_extension_service].Notification = make([]xml_System_Services_Extension_service_Notification, len(var_system_services_extension_service_notification))
        
		for i_system_services_extension_service_notification, v_system_services_extension_service_notification := range var_system_services_extension_service_notification {
            var var_system_services_extension_service_notification_allow_clients []System_Services_Extension_service_Notification_Allow_clients_Model
            resp.Diagnostics.Append(v_system_services_extension_service_notification.Allow_clients.ElementsAs(ctx, &var_system_services_extension_service_notification_allow_clients, false)...)
            if resp.Diagnostics.HasError() {
                return
            }
	    config.Groups.System[i_system].Services[i_system_services].Extension_service[i_system_services_extension_service].Notification[i_system_services_extension_service_notification].Allow_clients = make([]xml_System_Services_Extension_service_Notification_Allow_clients, len(var_system_services_extension_service_notification_allow_clients))
        
		for i_system_services_extension_service_notification_allow_clients, v_system_services_extension_service_notification_allow_clients := range var_system_services_extension_service_notification_allow_clients {
			var var_system_services_extension_service_notification_allow_clients_address []string
			resp.Diagnostics.Append(v_system_services_extension_service_notification_allow_clients.Address.ElementsAs(ctx, &var_system_services_extension_service_notification_allow_clients_address, false)...)
			for _, v_system_services_extension_service_notification_allow_clients_address := range var_system_services_extension_service_notification_allow_clients_address {
				config.Groups.System[i_system].Services[i_system_services].Extension_service[i_system_services_extension_service].Notification[i_system_services_extension_service_notification].Allow_clients[i_system_services_extension_service_notification_allow_clients].Address = append(config.Groups.System[i_system].Services[i_system_services].Extension_service[i_system_services_extension_service].Notification[i_system_services_extension_service_notification].Allow_clients[i_system_services_extension_service_notification_allow_clients].Address, &v_system_services_extension_service_notification_allow_clients_address)
			}
        }
        }
        }
            var var_system_services_netconf []System_Services_Netconf_Model
            resp.Diagnostics.Append(v_system_services.Netconf.ElementsAs(ctx, &var_system_services_netconf, false)...)
            if resp.Diagnostics.HasError() {
                return
            }
	    config.Groups.System[i_system].Services[i_system_services].Netconf = make([]xml_System_Services_Netconf, len(var_system_services_netconf))
        
		for i_system_services_netconf, v_system_services_netconf := range var_system_services_netconf {
            var var_system_services_netconf_ssh []System_Services_Netconf_Ssh_Model
            resp.Diagnostics.Append(v_system_services_netconf.Ssh.ElementsAs(ctx, &var_system_services_netconf_ssh, false)...)
            if resp.Diagnostics.HasError() {
                return
            }
	    config.Groups.System[i_system].Services[i_system_services].Netconf[i_system_services_netconf].Ssh = make([]xml_System_Services_Netconf_Ssh, len(var_system_services_netconf_ssh))
        
        }
            var var_system_services_rest []System_Services_Rest_Model
            resp.Diagnostics.Append(v_system_services.Rest.ElementsAs(ctx, &var_system_services_rest, false)...)
            if resp.Diagnostics.HasError() {
                return
            }
	    config.Groups.System[i_system].Services[i_system_services].Rest = make([]xml_System_Services_Rest, len(var_system_services_rest))
        
		for i_system_services_rest, v_system_services_rest := range var_system_services_rest {
            var var_system_services_rest_http []System_Services_Rest_Http_Model
            resp.Diagnostics.Append(v_system_services_rest.Http.ElementsAs(ctx, &var_system_services_rest_http, false)...)
            if resp.Diagnostics.HasError() {
                return
            }
	    config.Groups.System[i_system].Services[i_system_services].Rest[i_system_services_rest].Http = make([]xml_System_Services_Rest_Http, len(var_system_services_rest_http))
        
		for i_system_services_rest_http, v_system_services_rest_http := range var_system_services_rest_http {
            config.Groups.System[i_system].Services[i_system_services].Rest[i_system_services_rest].Http[i_system_services_rest_http].Port = v_system_services_rest_http.Port.ValueStringPointer()
        }
            config.Groups.System[i_system].Services[i_system_services].Rest[i_system_services_rest].Enable_explorer = v_system_services_rest.Enable_explorer.ValueStringPointer()
        }
        }
        var var_system_syslog []System_Syslog_Model
        resp.Diagnostics.Append(v_system.Syslog.ElementsAs(ctx, &var_system_syslog, false)...)
        if resp.Diagnostics.HasError() {
            return
        }
	    config.Groups.System[i_system].Syslog = make([]xml_System_Syslog, len(var_system_syslog))
        
		for i_system_syslog, v_system_syslog := range var_system_syslog {
            var var_system_syslog_user []System_Syslog_User_Model
            resp.Diagnostics.Append(v_system_syslog.User.ElementsAs(ctx, &var_system_syslog_user, false)...)
            if resp.Diagnostics.HasError() {
                return
            }
	    config.Groups.System[i_system].Syslog[i_system_syslog].User = make([]xml_System_Syslog_User, len(var_system_syslog_user))
        
		for i_system_syslog_user, v_system_syslog_user := range var_system_syslog_user {
            config.Groups.System[i_system].Syslog[i_system_syslog].User[i_system_syslog_user].Name = v_system_syslog_user.Name.ValueStringPointer()
            var var_system_syslog_user_contents []System_Syslog_User_Contents_Model
            resp.Diagnostics.Append(v_system_syslog_user.Contents.ElementsAs(ctx, &var_system_syslog_user_contents, false)...)
            if resp.Diagnostics.HasError() {
                return
            }
	    config.Groups.System[i_system].Syslog[i_system_syslog].User[i_system_syslog_user].Contents = make([]xml_System_Syslog_User_Contents, len(var_system_syslog_user_contents))
        
		for i_system_syslog_user_contents, v_system_syslog_user_contents := range var_system_syslog_user_contents {
            config.Groups.System[i_system].Syslog[i_system_syslog].User[i_system_syslog_user].Contents[i_system_syslog_user_contents].Name = v_system_syslog_user_contents.Name.ValueStringPointer()
            config.Groups.System[i_system].Syslog[i_system_syslog].User[i_system_syslog_user].Contents[i_system_syslog_user_contents].Emergency = v_system_syslog_user_contents.Emergency.ValueStringPointer()
        }
        }
            var var_system_syslog_file []System_Syslog_File_Model
            resp.Diagnostics.Append(v_system_syslog.File.ElementsAs(ctx, &var_system_syslog_file, false)...)
            if resp.Diagnostics.HasError() {
                return
            }
	    config.Groups.System[i_system].Syslog[i_system_syslog].File = make([]xml_System_Syslog_File, len(var_system_syslog_file))
        
		for i_system_syslog_file, v_system_syslog_file := range var_system_syslog_file {
            config.Groups.System[i_system].Syslog[i_system_syslog].File[i_system_syslog_file].Name = v_system_syslog_file.Name.ValueStringPointer()
            var var_system_syslog_file_contents []System_Syslog_File_Contents_Model
            resp.Diagnostics.Append(v_system_syslog_file.Contents.ElementsAs(ctx, &var_system_syslog_file_contents, false)...)
            if resp.Diagnostics.HasError() {
                return
            }
	    config.Groups.System[i_system].Syslog[i_system_syslog].File[i_system_syslog_file].Contents = make([]xml_System_Syslog_File_Contents, len(var_system_syslog_file_contents))
        
		for i_system_syslog_file_contents, v_system_syslog_file_contents := range var_system_syslog_file_contents {
            config.Groups.System[i_system].Syslog[i_system_syslog].File[i_system_syslog_file].Contents[i_system_syslog_file_contents].Name = v_system_syslog_file_contents.Name.ValueStringPointer()
            config.Groups.System[i_system].Syslog[i_system_syslog].File[i_system_syslog_file].Contents[i_system_syslog_file_contents].Any = v_system_syslog_file_contents.Any.ValueStringPointer()
            config.Groups.System[i_system].Syslog[i_system_syslog].File[i_system_syslog_file].Contents[i_system_syslog_file_contents].Info = v_system_syslog_file_contents.Info.ValueStringPointer()
        }
        }
        }
        var var_system_license []System_License_Model
        resp.Diagnostics.Append(v_system.License.ElementsAs(ctx, &var_system_license, false)...)
        if resp.Diagnostics.HasError() {
            return
        }
	    config.Groups.System[i_system].License = make([]xml_System_License, len(var_system_license))
        
		for i_system_license, v_system_license := range var_system_license {
            var var_system_license_autoupdate []System_License_Autoupdate_Model
            resp.Diagnostics.Append(v_system_license.Autoupdate.ElementsAs(ctx, &var_system_license_autoupdate, false)...)
            if resp.Diagnostics.HasError() {
                return
            }
	    config.Groups.System[i_system].License[i_system_license].Autoupdate = make([]xml_System_License_Autoupdate, len(var_system_license_autoupdate))
        
		for i_system_license_autoupdate, v_system_license_autoupdate := range var_system_license_autoupdate {
            var var_system_license_autoupdate_url []System_License_Autoupdate_Url_Model
            resp.Diagnostics.Append(v_system_license_autoupdate.Url.ElementsAs(ctx, &var_system_license_autoupdate_url, false)...)
            if resp.Diagnostics.HasError() {
                return
            }
	    config.Groups.System[i_system].License[i_system_license].Autoupdate[i_system_license_autoupdate].Url = make([]xml_System_License_Autoupdate_Url, len(var_system_license_autoupdate_url))
        
		for i_system_license_autoupdate_url, v_system_license_autoupdate_url := range var_system_license_autoupdate_url {
            config.Groups.System[i_system].License[i_system_license].Autoupdate[i_system_license_autoupdate].Url[i_system_license_autoupdate_url].Name = v_system_license_autoupdate_url.Name.ValueStringPointer()
        }
        }
        }
    }
	
	err := r.client.SendTransaction(plan.ResourceName.ValueString(), config, false)
	if err != nil {
		resp.Diagnostics.AddError("Failed while Sending", err.Error())
		return
	}
	commit_err := r.client.SendCommit()
	if commit_err != nil {
		resp.Diagnostics.AddError("Failed while committing apply-group", commit_err.Error())
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}



// Delete implements resource.Resource.
func (r *resource_Apply_Groups) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state Groups_Model
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
		resp.Diagnostics.AddError("Failed while deleting configuration", err.Error())
		return
	}
    commit_err := r.client.SendCommit()
	if commit_err != nil {
		resp.Diagnostics.AddError("Failed while committing apply-group", commit_err.Error())
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}
