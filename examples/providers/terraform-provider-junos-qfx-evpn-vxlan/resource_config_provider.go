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
		Forwarding_options []xml_Forwarding_options `xml:"forwarding-options,omitempty"`
		Interfaces []xml_Interfaces `xml:"interfaces,omitempty"`
		Policy_options []xml_Policy_options `xml:"policy-options,omitempty"`
		Protocols []xml_Protocols `xml:"protocols,omitempty"`
		Routing_instances []xml_Routing_instances `xml:"routing-instances,omitempty"`
		Routing_options []xml_Routing_options `xml:"routing-options,omitempty"`
		Snmp []xml_Snmp `xml:"snmp,omitempty"`
		Switch_options []xml_Switch_options `xml:"switch-options,omitempty"`
		System []xml_System `xml:"system,omitempty"`
		Vlans []xml_Vlans `xml:"vlans,omitempty"`
	}
}
type xml_Chassis struct {
	XMLName xml.Name `xml:"chassis"`
	Aggregated_devices []xml_Chassis_Aggregated_devices `xml:"aggregated-devices,omitempty"`
}

type xml_Forwarding_options struct {
	XMLName xml.Name `xml:"forwarding-options"`
	Storm_control_profiles []xml_Forwarding_options_Storm_control_profiles `xml:"storm-control-profiles,omitempty"`
}

type xml_Interfaces struct {
	XMLName xml.Name `xml:"interfaces"`
	Interface []xml_Interfaces_Interface `xml:"interface,omitempty"`
}

type xml_Policy_options struct {
	XMLName xml.Name `xml:"policy-options"`
	Policy_statement []xml_Policy_options_Policy_statement `xml:"policy-statement,omitempty"`
	Community []xml_Policy_options_Community `xml:"community,omitempty"`
}

type xml_Protocols struct {
	XMLName xml.Name `xml:"protocols"`
	Bgp []xml_Protocols_Bgp `xml:"bgp,omitempty"`
	Evpn []xml_Protocols_Evpn `xml:"evpn,omitempty"`
	Lldp []xml_Protocols_Lldp `xml:"lldp,omitempty"`
	Igmp_snooping []xml_Protocols_Igmp_snooping `xml:"igmp-snooping,omitempty"`
}

type xml_Routing_instances struct {
	XMLName xml.Name `xml:"routing-instances"`
	Instance []xml_Routing_instances_Instance `xml:"instance,omitempty"`
}

type xml_Routing_options struct {
	XMLName xml.Name `xml:"routing-options"`
	Static []xml_Routing_options_Static `xml:"static,omitempty"`
	Router_id         *string  `xml:"router-id,omitempty"`
	Forwarding_table []xml_Routing_options_Forwarding_table `xml:"forwarding-table,omitempty"`
}

type xml_Snmp struct {
	XMLName xml.Name `xml:"snmp"`
	Location         *string  `xml:"location,omitempty"`
	Contact         *string  `xml:"contact,omitempty"`
	Community []xml_Snmp_Community `xml:"community,omitempty"`
}

type xml_Switch_options struct {
	XMLName xml.Name `xml:"switch-options"`
	Vtep_source_interface []xml_Switch_options_Vtep_source_interface `xml:"vtep-source-interface,omitempty"`
	Route_distinguisher []xml_Switch_options_Route_distinguisher `xml:"route-distinguisher,omitempty"`
	Vrf_target []xml_Switch_options_Vrf_target `xml:"vrf-target,omitempty"`
}

type xml_System struct {
	XMLName xml.Name `xml:"system"`
	Login []xml_System_Login `xml:"login,omitempty"`
	Root_authentication []xml_System_Root_authentication `xml:"root-authentication,omitempty"`
	Host_name         *string  `xml:"host-name,omitempty"`
	Services []xml_System_Services `xml:"services,omitempty"`
	Syslog []xml_System_Syslog `xml:"syslog,omitempty"`
	Extensions []xml_System_Extensions `xml:"extensions,omitempty"`
}

type xml_Vlans struct {
	XMLName xml.Name `xml:"vlans"`
	Vlan []xml_Vlans_Vlan `xml:"vlan,omitempty"`
}


type xml_Chassis_Aggregated_devices struct {
	XMLName xml.Name `xml:"aggregated-devices"`
	Ethernet []xml_Chassis_Aggregated_devices_Ethernet `xml:"ethernet,omitempty"`
}
type xml_Forwarding_options_Storm_control_profiles struct {
	XMLName xml.Name `xml:"storm-control-profiles"`
	Name         *string  `xml:"name,omitempty"`
	All []xml_Forwarding_options_Storm_control_profiles_All `xml:"all,omitempty"`
}
type xml_Interfaces_Interface struct {
	XMLName xml.Name `xml:"interface"`
	Name         *string  `xml:"name,omitempty"`
	Description         *string  `xml:"description,omitempty"`
	Vlan_tagging         *string  `xml:"vlan-tagging,omitempty"`
	Esi []xml_Interfaces_Interface_Esi `xml:"esi,omitempty"`
	Ether_options []xml_Interfaces_Interface_Ether_options `xml:"ether-options,omitempty"`
	Aggregated_ether_options []xml_Interfaces_Interface_Aggregated_ether_options `xml:"aggregated-ether-options,omitempty"`
	Unit []xml_Interfaces_Interface_Unit `xml:"unit,omitempty"`
}
type xml_Policy_options_Policy_statement struct {
	XMLName xml.Name `xml:"policy-statement"`
	Name         *string  `xml:"name,omitempty"`
	Term []xml_Policy_options_Policy_statement_Term `xml:"term,omitempty"`
	Then []xml_Policy_options_Policy_statement_Then `xml:"then,omitempty"`
}
type xml_Policy_options_Community struct {
	XMLName xml.Name `xml:"community"`
	Name         *string  `xml:"name,omitempty"`
	Members         []*string  `xml:"members,omitempty"`
}
type xml_Protocols_Bgp struct {
	XMLName xml.Name `xml:"bgp"`
	Group []xml_Protocols_Bgp_Group `xml:"group,omitempty"`
}
type xml_Protocols_Evpn struct {
	XMLName xml.Name `xml:"evpn"`
	Encapsulation         *string  `xml:"encapsulation,omitempty"`
	Multicast_mode         *string  `xml:"multicast-mode,omitempty"`
	Default_gateway         *string  `xml:"default-gateway,omitempty"`
	Extended_vni_list         []*string  `xml:"extended-vni-list,omitempty"`
	No_core_isolation         *string  `xml:"no-core-isolation,omitempty"`
}
type xml_Protocols_Lldp struct {
	XMLName xml.Name `xml:"lldp"`
	Interface []xml_Protocols_Lldp_Interface `xml:"interface,omitempty"`
}
type xml_Protocols_Igmp_snooping struct {
	XMLName xml.Name `xml:"igmp-snooping"`
	Vlan []xml_Protocols_Igmp_snooping_Vlan `xml:"vlan,omitempty"`
}
type xml_Routing_instances_Instance struct {
	XMLName xml.Name `xml:"instance"`
	Name         *string  `xml:"name,omitempty"`
	Instance_type         *string  `xml:"instance-type,omitempty"`
	Interface []xml_Routing_instances_Instance_Interface `xml:"interface,omitempty"`
	Route_distinguisher []xml_Routing_instances_Instance_Route_distinguisher `xml:"route-distinguisher,omitempty"`
	Vrf_target []xml_Routing_instances_Instance_Vrf_target `xml:"vrf-target,omitempty"`
	Vrf_table_label []xml_Routing_instances_Instance_Vrf_table_label `xml:"vrf-table-label,omitempty"`
	Routing_options []xml_Routing_instances_Instance_Routing_options `xml:"routing-options,omitempty"`
	Protocols []xml_Routing_instances_Instance_Protocols `xml:"protocols,omitempty"`
}
type xml_Routing_options_Static struct {
	XMLName xml.Name `xml:"static"`
	Route []xml_Routing_options_Static_Route `xml:"route,omitempty"`
}
type xml_Routing_options_Forwarding_table struct {
	XMLName xml.Name `xml:"forwarding-table"`
	Export         []*string  `xml:"export,omitempty"`
	Ecmp_fast_reroute         *string  `xml:"ecmp-fast-reroute,omitempty"`
	Chained_composite_next_hop []xml_Routing_options_Forwarding_table_Chained_composite_next_hop `xml:"chained-composite-next-hop,omitempty"`
}
type xml_Snmp_Community struct {
	XMLName xml.Name `xml:"community"`
	Name         *string  `xml:"name,omitempty"`
	Authorization         *string  `xml:"authorization,omitempty"`
}
type xml_Switch_options_Vtep_source_interface struct {
	XMLName xml.Name `xml:"vtep-source-interface"`
	Interface_name         *string  `xml:"interface-name,omitempty"`
}
type xml_Switch_options_Route_distinguisher struct {
	XMLName xml.Name `xml:"route-distinguisher"`
	Rd_type         *string  `xml:"rd-type,omitempty"`
}
type xml_Switch_options_Vrf_target struct {
	XMLName xml.Name `xml:"vrf-target"`
	Community         *string  `xml:"community,omitempty"`
	Auto []xml_Switch_options_Vrf_target_Auto `xml:"auto,omitempty"`
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
type xml_System_Extensions struct {
	XMLName xml.Name `xml:"extensions"`
	Providers []xml_System_Extensions_Providers `xml:"providers,omitempty"`
}
type xml_Vlans_Vlan struct {
	XMLName xml.Name `xml:"vlan"`
	Name         *string  `xml:"name,omitempty"`
	Vlan_id         *string  `xml:"vlan-id,omitempty"`
	L3_interface         *string  `xml:"l3-interface,omitempty"`
	Vxlan []xml_Vlans_Vlan_Vxlan `xml:"vxlan,omitempty"`
}

type xml_Chassis_Aggregated_devices_Ethernet struct {
	XMLName xml.Name `xml:"ethernet"`
	Device_count         *string  `xml:"device-count,omitempty"`
}
type xml_Forwarding_options_Storm_control_profiles_All struct {
	XMLName xml.Name `xml:"all"`
}
type xml_Interfaces_Interface_Esi struct {
	XMLName xml.Name `xml:"esi"`
	Identifier         *string  `xml:"identifier,omitempty"`
	All_active         *string  `xml:"all-active,omitempty"`
}
type xml_Interfaces_Interface_Ether_options struct {
	XMLName xml.Name `xml:"ether-options"`
	Ieee_802_3ad []xml_Interfaces_Interface_Ether_options_Ieee_802_3ad `xml:"ieee-802.3ad,omitempty"`
}
type xml_Interfaces_Interface_Aggregated_ether_options struct {
	XMLName xml.Name `xml:"aggregated-ether-options"`
	Lacp []xml_Interfaces_Interface_Aggregated_ether_options_Lacp `xml:"lacp,omitempty"`
}
type xml_Interfaces_Interface_Unit struct {
	XMLName xml.Name `xml:"unit"`
	Name         *string  `xml:"name,omitempty"`
	Description         *string  `xml:"description,omitempty"`
	Vlan_id         *string  `xml:"vlan-id,omitempty"`
	Family []xml_Interfaces_Interface_Unit_Family `xml:"family,omitempty"`
	Mac         *string  `xml:"mac,omitempty"`
}
type xml_Policy_options_Policy_statement_Term struct {
	XMLName xml.Name `xml:"term"`
	Name         *string  `xml:"name,omitempty"`
	From []xml_Policy_options_Policy_statement_Term_From `xml:"from,omitempty"`
	Then []xml_Policy_options_Policy_statement_Term_Then `xml:"then,omitempty"`
}
type xml_Policy_options_Policy_statement_Then struct {
	XMLName xml.Name `xml:"then"`
	Load_balance []xml_Policy_options_Policy_statement_Then_Load_balance `xml:"load-balance,omitempty"`
}
type xml_Protocols_Bgp_Group struct {
	XMLName xml.Name `xml:"group"`
	Name         *string  `xml:"name,omitempty"`
	Type         *string  `xml:"type,omitempty"`
	Multihop []xml_Protocols_Bgp_Group_Multihop `xml:"multihop,omitempty"`
	Local_address         *string  `xml:"local-address,omitempty"`
	Mtu_discovery         *string  `xml:"mtu-discovery,omitempty"`
	Import         []*string  `xml:"import,omitempty"`
	Family []xml_Protocols_Bgp_Group_Family `xml:"family,omitempty"`
	Export         []*string  `xml:"export,omitempty"`
	Vpn_apply_export         *string  `xml:"vpn-apply-export,omitempty"`
	Cluster         *string  `xml:"cluster,omitempty"`
	Local_as []xml_Protocols_Bgp_Group_Local_as `xml:"local-as,omitempty"`
	Multipath []xml_Protocols_Bgp_Group_Multipath `xml:"multipath,omitempty"`
	Bfd_liveness_detection []xml_Protocols_Bgp_Group_Bfd_liveness_detection `xml:"bfd-liveness-detection,omitempty"`
	Allow         []*string  `xml:"allow,omitempty"`
	Neighbor []xml_Protocols_Bgp_Group_Neighbor `xml:"neighbor,omitempty"`
}
type xml_Protocols_Lldp_Interface struct {
	XMLName xml.Name `xml:"interface"`
	Name         *string  `xml:"name,omitempty"`
}
type xml_Protocols_Igmp_snooping_Vlan struct {
	XMLName xml.Name `xml:"vlan"`
	Name         *string  `xml:"name,omitempty"`
}
type xml_Routing_instances_Instance_Interface struct {
	XMLName xml.Name `xml:"interface"`
	Name         *string  `xml:"name,omitempty"`
}
type xml_Routing_instances_Instance_Route_distinguisher struct {
	XMLName xml.Name `xml:"route-distinguisher"`
	Rd_type         *string  `xml:"rd-type,omitempty"`
}
type xml_Routing_instances_Instance_Vrf_target struct {
	XMLName xml.Name `xml:"vrf-target"`
	Community         *string  `xml:"community,omitempty"`
}
type xml_Routing_instances_Instance_Vrf_table_label struct {
	XMLName xml.Name `xml:"vrf-table-label"`
}
type xml_Routing_instances_Instance_Routing_options struct {
	XMLName xml.Name `xml:"routing-options"`
	Auto_export []xml_Routing_instances_Instance_Routing_options_Auto_export `xml:"auto-export,omitempty"`
}
type xml_Routing_instances_Instance_Protocols struct {
	XMLName xml.Name `xml:"protocols"`
	Ospf []xml_Routing_instances_Instance_Protocols_Ospf `xml:"ospf,omitempty"`
	Evpn []xml_Routing_instances_Instance_Protocols_Evpn `xml:"evpn,omitempty"`
}
type xml_Routing_options_Static_Route struct {
	XMLName xml.Name `xml:"route"`
	Name         *string  `xml:"name,omitempty"`
	Next_hop         []*string  `xml:"next-hop,omitempty"`
}
type xml_Routing_options_Forwarding_table_Chained_composite_next_hop struct {
	XMLName xml.Name `xml:"chained-composite-next-hop"`
	Ingress []xml_Routing_options_Forwarding_table_Chained_composite_next_hop_Ingress `xml:"ingress,omitempty"`
}
type xml_Switch_options_Vrf_target_Auto struct {
	XMLName xml.Name `xml:"auto"`
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
type xml_System_Extensions_Providers struct {
	XMLName xml.Name `xml:"providers"`
	Name         *string  `xml:"name,omitempty"`
	License_type []xml_System_Extensions_Providers_License_type `xml:"license-type,omitempty"`
}
type xml_Vlans_Vlan_Vxlan struct {
	XMLName xml.Name `xml:"vxlan"`
	Vni         *string  `xml:"vni,omitempty"`
}

type xml_Interfaces_Interface_Ether_options_Ieee_802_3ad struct {
	XMLName xml.Name `xml:"ieee-802.3ad"`
	Bundle         *string  `xml:"bundle,omitempty"`
}
type xml_Interfaces_Interface_Aggregated_ether_options_Lacp struct {
	XMLName xml.Name `xml:"lacp"`
	Active         *string  `xml:"active,omitempty"`
	Periodic         *string  `xml:"periodic,omitempty"`
	System_id         *string  `xml:"system-id,omitempty"`
}
type xml_Interfaces_Interface_Unit_Family struct {
	XMLName xml.Name `xml:"family"`
	Inet []xml_Interfaces_Interface_Unit_Family_Inet `xml:"inet,omitempty"`
	Ethernet_switching []xml_Interfaces_Interface_Unit_Family_Ethernet_switching `xml:"ethernet-switching,omitempty"`
}
type xml_Policy_options_Policy_statement_Term_From struct {
	XMLName xml.Name `xml:"from"`
	Protocol         []*string  `xml:"protocol,omitempty"`
	Route_filter []xml_Policy_options_Policy_statement_Term_From_Route_filter `xml:"route-filter,omitempty"`
}
type xml_Policy_options_Policy_statement_Term_Then struct {
	XMLName xml.Name `xml:"then"`
	Community []xml_Policy_options_Policy_statement_Term_Then_Community `xml:"community,omitempty"`
	Accept         *string  `xml:"accept,omitempty"`
	Reject         *string  `xml:"reject,omitempty"`
}
type xml_Policy_options_Policy_statement_Then_Load_balance struct {
	XMLName xml.Name `xml:"load-balance"`
	Per_packet         *string  `xml:"per-packet,omitempty"`
}
type xml_Protocols_Bgp_Group_Multihop struct {
	XMLName xml.Name `xml:"multihop"`
	No_nexthop_change         *string  `xml:"no-nexthop-change,omitempty"`
}
type xml_Protocols_Bgp_Group_Family struct {
	XMLName xml.Name `xml:"family"`
	Evpn []xml_Protocols_Bgp_Group_Family_Evpn `xml:"evpn,omitempty"`
}
type xml_Protocols_Bgp_Group_Local_as struct {
	XMLName xml.Name `xml:"local-as"`
	As_number         *string  `xml:"as-number,omitempty"`
}
type xml_Protocols_Bgp_Group_Multipath struct {
	XMLName xml.Name `xml:"multipath"`
	Multiple_as         *string  `xml:"multiple-as,omitempty"`
}
type xml_Protocols_Bgp_Group_Bfd_liveness_detection struct {
	XMLName xml.Name `xml:"bfd-liveness-detection"`
	Minimum_interval         *string  `xml:"minimum-interval,omitempty"`
	Multiplier         *string  `xml:"multiplier,omitempty"`
}
type xml_Protocols_Bgp_Group_Neighbor struct {
	XMLName xml.Name `xml:"neighbor"`
	Name         *string  `xml:"name,omitempty"`
	Description         *string  `xml:"description,omitempty"`
	Peer_as         *string  `xml:"peer-as,omitempty"`
}
type xml_Routing_instances_Instance_Routing_options_Auto_export struct {
	XMLName xml.Name `xml:"auto-export"`
}
type xml_Routing_instances_Instance_Protocols_Ospf struct {
	XMLName xml.Name `xml:"ospf"`
	Export         []*string  `xml:"export,omitempty"`
	Area []xml_Routing_instances_Instance_Protocols_Ospf_Area `xml:"area,omitempty"`
}
type xml_Routing_instances_Instance_Protocols_Evpn struct {
	XMLName xml.Name `xml:"evpn"`
	Ip_prefix_routes []xml_Routing_instances_Instance_Protocols_Evpn_Ip_prefix_routes `xml:"ip-prefix-routes,omitempty"`
}
type xml_Routing_options_Forwarding_table_Chained_composite_next_hop_Ingress struct {
	XMLName xml.Name `xml:"ingress"`
	Evpn         *string  `xml:"evpn,omitempty"`
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
	Notice         *string  `xml:"notice,omitempty"`
	Info         *string  `xml:"info,omitempty"`
}
type xml_System_Extensions_Providers_License_type struct {
	XMLName xml.Name `xml:"license-type"`
	Name         *string  `xml:"name,omitempty"`
	Deployment_scope         []*string  `xml:"deployment-scope,omitempty"`
}

type xml_Interfaces_Interface_Unit_Family_Inet struct {
	XMLName xml.Name `xml:"inet"`
	Address []xml_Interfaces_Interface_Unit_Family_Inet_Address `xml:"address,omitempty"`
}
type xml_Interfaces_Interface_Unit_Family_Ethernet_switching struct {
	XMLName xml.Name `xml:"ethernet-switching"`
	Vlan []xml_Interfaces_Interface_Unit_Family_Ethernet_switching_Vlan `xml:"vlan,omitempty"`
}
type xml_Policy_options_Policy_statement_Term_From_Route_filter struct {
	XMLName xml.Name `xml:"route-filter"`
	Address         *string  `xml:"address,omitempty"`
	Exact         *string  `xml:"exact,omitempty"`
	Orlonger         *string  `xml:"orlonger,omitempty"`
	Prefix_length_range         *string  `xml:"prefix-length-range,omitempty"`
}
type xml_Policy_options_Policy_statement_Term_Then_Community struct {
	XMLName xml.Name `xml:"community"`
	Add         *string  `xml:"add,omitempty"`
	Community_name         *string  `xml:"community-name,omitempty"`
}
type xml_Protocols_Bgp_Group_Family_Evpn struct {
	XMLName xml.Name `xml:"evpn"`
	Signaling []xml_Protocols_Bgp_Group_Family_Evpn_Signaling `xml:"signaling,omitempty"`
}
type xml_Routing_instances_Instance_Protocols_Ospf_Area struct {
	XMLName xml.Name `xml:"area"`
	Name         *string  `xml:"name,omitempty"`
	Interface []xml_Routing_instances_Instance_Protocols_Ospf_Area_Interface `xml:"interface,omitempty"`
}
type xml_Routing_instances_Instance_Protocols_Evpn_Ip_prefix_routes struct {
	XMLName xml.Name `xml:"ip-prefix-routes"`
	Advertise         *string  `xml:"advertise,omitempty"`
	Encapsulation         *string  `xml:"encapsulation,omitempty"`
	Vni         *string  `xml:"vni,omitempty"`
	Export         []*string  `xml:"export,omitempty"`
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
type xml_Interfaces_Interface_Unit_Family_Ethernet_switching_Vlan struct {
	XMLName xml.Name `xml:"vlan"`
	Members         []*string  `xml:"members,omitempty"`
}
type xml_Protocols_Bgp_Group_Family_Evpn_Signaling struct {
	XMLName xml.Name `xml:"signaling"`
	Delay_route_advertisements []xml_Protocols_Bgp_Group_Family_Evpn_Signaling_Delay_route_advertisements `xml:"delay-route-advertisements,omitempty"`
}
type xml_Routing_instances_Instance_Protocols_Ospf_Area_Interface struct {
	XMLName xml.Name `xml:"interface"`
	Name         *string  `xml:"name,omitempty"`
	Metric         *string  `xml:"metric,omitempty"`
}

type xml_Protocols_Bgp_Group_Family_Evpn_Signaling_Delay_route_advertisements struct {
	XMLName xml.Name `xml:"delay-route-advertisements"`
	Minimum_delay []xml_Protocols_Bgp_Group_Family_Evpn_Signaling_Delay_route_advertisements_Minimum_delay `xml:"minimum-delay,omitempty"`
}

type xml_Protocols_Bgp_Group_Family_Evpn_Signaling_Delay_route_advertisements_Minimum_delay struct {
	XMLName xml.Name `xml:"minimum-delay"`
	Routing_uptime         *string  `xml:"routing-uptime,omitempty"`
}




// Collecting objects from the .tf file
type Groups_Model struct {
	ResourceName	types.String `tfsdk:"resource_name"`
	Chassis types.List `tfsdk:"chassis"`
	Forwarding_options types.List `tfsdk:"forwarding_options"`
	Interfaces types.List `tfsdk:"interfaces"`
	Policy_options types.List `tfsdk:"policy_options"`
	Protocols types.List `tfsdk:"protocols"`
	Routing_instances types.List `tfsdk:"routing_instances"`
	Routing_options types.List `tfsdk:"routing_options"`
	Snmp types.List `tfsdk:"snmp"`
	Switch_options types.List `tfsdk:"switch_options"`
	System types.List `tfsdk:"system"`
	Vlans types.List `tfsdk:"vlans"`
}
func (o Groups_Model) AttrTypes() map[string]attr.Type {
	return map[string]attr.Type {
		"chassis": 	types.ListType{ElemType: types.ObjectType{AttrTypes: Chassis_Model{}.AttrTypes()}},
		"forwarding_options": 	types.ListType{ElemType: types.ObjectType{AttrTypes: Forwarding_options_Model{}.AttrTypes()}},
		"interfaces": 	types.ListType{ElemType: types.ObjectType{AttrTypes: Interfaces_Model{}.AttrTypes()}},
		"policy_options": 	types.ListType{ElemType: types.ObjectType{AttrTypes: Policy_options_Model{}.AttrTypes()}},
		"protocols": 	types.ListType{ElemType: types.ObjectType{AttrTypes: Protocols_Model{}.AttrTypes()}},
		"routing_instances": 	types.ListType{ElemType: types.ObjectType{AttrTypes: Routing_instances_Model{}.AttrTypes()}},
		"routing_options": 	types.ListType{ElemType: types.ObjectType{AttrTypes: Routing_options_Model{}.AttrTypes()}},
		"snmp": 	types.ListType{ElemType: types.ObjectType{AttrTypes: Snmp_Model{}.AttrTypes()}},
		"switch_options": 	types.ListType{ElemType: types.ObjectType{AttrTypes: Switch_options_Model{}.AttrTypes()}},
		"system": 	types.ListType{ElemType: types.ObjectType{AttrTypes: System_Model{}.AttrTypes()}},
		"vlans": 	types.ListType{ElemType: types.ObjectType{AttrTypes: Vlans_Model{}.AttrTypes()}},
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
		"forwarding_options": schema.ListNestedAttribute{
			Optional: true,
			NestedObject: schema.NestedAttributeObject{
				Attributes: Forwarding_options_Model{}.Attributes(),
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
		"snmp": schema.ListNestedAttribute{
			Optional: true,
			NestedObject: schema.NestedAttributeObject{
				Attributes: Snmp_Model{}.Attributes(),
			},
		},
		"switch_options": schema.ListNestedAttribute{
			Optional: true,
			NestedObject: schema.NestedAttributeObject{
				Attributes: Switch_options_Model{}.Attributes(),
			},
		},
		"system": schema.ListNestedAttribute{
			Optional: true,
			NestedObject: schema.NestedAttributeObject{
				Attributes: System_Model{}.Attributes(),
			},
		},
		"vlans": schema.ListNestedAttribute{
			Optional: true,
			NestedObject: schema.NestedAttributeObject{
				Attributes: Vlans_Model{}.Attributes(),
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
type Forwarding_options_Model struct {
	Storm_control_profiles	types.List `tfsdk:"storm_control_profiles"`
}
func (o Forwarding_options_Model) AttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"storm_control_profiles": 	types.ListType{ElemType: types.ObjectType{AttrTypes: Forwarding_options_Storm_control_profiles_Model{}.AttrTypes()}},
	}
}
func (o Forwarding_options_Model) Attributes() map[string]schema.Attribute {
	return map[string]schema.Attribute{
		"storm_control_profiles": schema.ListNestedAttribute{
			Optional: true,
			NestedObject: schema.NestedAttributeObject{
				Attributes: Forwarding_options_Storm_control_profiles_Model{}.Attributes(),
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
	Community	types.List `tfsdk:"community"`
}
func (o Policy_options_Model) AttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"policy_statement": 	types.ListType{ElemType: types.ObjectType{AttrTypes: Policy_options_Policy_statement_Model{}.AttrTypes()}},
		"community": 	types.ListType{ElemType: types.ObjectType{AttrTypes: Policy_options_Community_Model{}.AttrTypes()}},
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
		"community": schema.ListNestedAttribute{
			Optional: true,
			NestedObject: schema.NestedAttributeObject{
				Attributes: Policy_options_Community_Model{}.Attributes(),
			},
		},
	}
}
type Protocols_Model struct {
	Bgp	types.List `tfsdk:"bgp"`
	Evpn	types.List `tfsdk:"evpn"`
	Lldp	types.List `tfsdk:"lldp"`
	Igmp_snooping	types.List `tfsdk:"igmp_snooping"`
}
func (o Protocols_Model) AttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"bgp": 	types.ListType{ElemType: types.ObjectType{AttrTypes: Protocols_Bgp_Model{}.AttrTypes()}},
		"evpn": 	types.ListType{ElemType: types.ObjectType{AttrTypes: Protocols_Evpn_Model{}.AttrTypes()}},
		"lldp": 	types.ListType{ElemType: types.ObjectType{AttrTypes: Protocols_Lldp_Model{}.AttrTypes()}},
		"igmp_snooping": 	types.ListType{ElemType: types.ObjectType{AttrTypes: Protocols_Igmp_snooping_Model{}.AttrTypes()}},
	}
}
func (o Protocols_Model) Attributes() map[string]schema.Attribute {
	return map[string]schema.Attribute{
		"bgp": schema.ListNestedAttribute{
			Optional: true,
			NestedObject: schema.NestedAttributeObject{
				Attributes: Protocols_Bgp_Model{}.Attributes(),
			},
		},
		"evpn": schema.ListNestedAttribute{
			Optional: true,
			NestedObject: schema.NestedAttributeObject{
				Attributes: Protocols_Evpn_Model{}.Attributes(),
			},
		},
		"lldp": schema.ListNestedAttribute{
			Optional: true,
			NestedObject: schema.NestedAttributeObject{
				Attributes: Protocols_Lldp_Model{}.Attributes(),
			},
		},
		"igmp_snooping": schema.ListNestedAttribute{
			Optional: true,
			NestedObject: schema.NestedAttributeObject{
				Attributes: Protocols_Igmp_snooping_Model{}.Attributes(),
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
	Router_id	types.String `tfsdk:"router_id"`
	Forwarding_table	types.List `tfsdk:"forwarding_table"`
}
func (o Routing_options_Model) AttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"static": 	types.ListType{ElemType: types.ObjectType{AttrTypes: Routing_options_Static_Model{}.AttrTypes()}},
		"router_id": 	types.StringType,
		"forwarding_table": 	types.ListType{ElemType: types.ObjectType{AttrTypes: Routing_options_Forwarding_table_Model{}.AttrTypes()}},
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
		"router_id": schema.StringAttribute{
			Optional: true,
			MarkdownDescription: "xpath is `config.Groups.Routing-options.Router-id`",
		},
		"forwarding_table": schema.ListNestedAttribute{
			Optional: true,
			NestedObject: schema.NestedAttributeObject{
				Attributes: Routing_options_Forwarding_table_Model{}.Attributes(),
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
type Switch_options_Model struct {
	Vtep_source_interface	types.List `tfsdk:"vtep_source_interface"`
	Route_distinguisher	types.List `tfsdk:"route_distinguisher"`
	Vrf_target	types.List `tfsdk:"vrf_target"`
}
func (o Switch_options_Model) AttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"vtep_source_interface": 	types.ListType{ElemType: types.ObjectType{AttrTypes: Switch_options_Vtep_source_interface_Model{}.AttrTypes()}},
		"route_distinguisher": 	types.ListType{ElemType: types.ObjectType{AttrTypes: Switch_options_Route_distinguisher_Model{}.AttrTypes()}},
		"vrf_target": 	types.ListType{ElemType: types.ObjectType{AttrTypes: Switch_options_Vrf_target_Model{}.AttrTypes()}},
	}
}
func (o Switch_options_Model) Attributes() map[string]schema.Attribute {
	return map[string]schema.Attribute{
		"vtep_source_interface": schema.ListNestedAttribute{
			Optional: true,
			NestedObject: schema.NestedAttributeObject{
				Attributes: Switch_options_Vtep_source_interface_Model{}.Attributes(),
			},
		},
		"route_distinguisher": schema.ListNestedAttribute{
			Optional: true,
			NestedObject: schema.NestedAttributeObject{
				Attributes: Switch_options_Route_distinguisher_Model{}.Attributes(),
			},
		},
		"vrf_target": schema.ListNestedAttribute{
			Optional: true,
			NestedObject: schema.NestedAttributeObject{
				Attributes: Switch_options_Vrf_target_Model{}.Attributes(),
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
	Extensions	types.List `tfsdk:"extensions"`
}
func (o System_Model) AttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"login": 	types.ListType{ElemType: types.ObjectType{AttrTypes: System_Login_Model{}.AttrTypes()}},
		"root_authentication": 	types.ListType{ElemType: types.ObjectType{AttrTypes: System_Root_authentication_Model{}.AttrTypes()}},
		"host_name": 	types.StringType,
		"services": 	types.ListType{ElemType: types.ObjectType{AttrTypes: System_Services_Model{}.AttrTypes()}},
		"syslog": 	types.ListType{ElemType: types.ObjectType{AttrTypes: System_Syslog_Model{}.AttrTypes()}},
		"extensions": 	types.ListType{ElemType: types.ObjectType{AttrTypes: System_Extensions_Model{}.AttrTypes()}},
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
		"extensions": schema.ListNestedAttribute{
			Optional: true,
			NestedObject: schema.NestedAttributeObject{
				Attributes: System_Extensions_Model{}.Attributes(),
			},
		},
	}
}
type Vlans_Model struct {
	Vlan	types.List `tfsdk:"vlan"`
}
func (o Vlans_Model) AttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"vlan": 	types.ListType{ElemType: types.ObjectType{AttrTypes: Vlans_Vlan_Model{}.AttrTypes()}},
	}
}
func (o Vlans_Model) Attributes() map[string]schema.Attribute {
	return map[string]schema.Attribute{
		"vlan": schema.ListNestedAttribute{
			Optional: true,
			NestedObject: schema.NestedAttributeObject{
				Attributes: Vlans_Vlan_Model{}.Attributes(),
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
type Forwarding_options_Storm_control_profiles_Model struct {
	Name	types.String `tfsdk:"name"`
	All	types.List `tfsdk:"all"`
}
func (o Forwarding_options_Storm_control_profiles_Model) AttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
	    "name": 	types.StringType,
	    "all": 	types.ListType{ElemType: types.ObjectType{AttrTypes: Forwarding_options_Storm_control_profiles_All_Model{}.AttrTypes()}},
	}
}
func (o Forwarding_options_Storm_control_profiles_Model) Attributes() map[string]schema.Attribute {
	return map[string]schema.Attribute{
	    "name": schema.StringAttribute{
		    Optional: true,
		    MarkdownDescription: "xpath is `config.Groups.Name.Storm_control_profiles`",
	    },
	    "all": schema.ListNestedAttribute{
		    Optional: true,
		    NestedObject: schema.NestedAttributeObject{
			    Attributes: Forwarding_options_Storm_control_profiles_All_Model{}.Attributes(),
	        },
        },
    }
}
type Interfaces_Interface_Model struct {
	Name	types.String `tfsdk:"name"`
	Description	types.String `tfsdk:"description"`
	Vlan_tagging	types.String `tfsdk:"vlan_tagging"`
	Esi	types.List `tfsdk:"esi"`
	Ether_options	types.List `tfsdk:"ether_options"`
	Aggregated_ether_options	types.List `tfsdk:"aggregated_ether_options"`
	Unit	types.List `tfsdk:"unit"`
}
func (o Interfaces_Interface_Model) AttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
	    "name": 	types.StringType,
	    "description": 	types.StringType,
	    "vlan_tagging": 	types.StringType,
	    "esi": 	types.ListType{ElemType: types.ObjectType{AttrTypes: Interfaces_Interface_Esi_Model{}.AttrTypes()}},
	    "ether_options": 	types.ListType{ElemType: types.ObjectType{AttrTypes: Interfaces_Interface_Ether_options_Model{}.AttrTypes()}},
	    "aggregated_ether_options": 	types.ListType{ElemType: types.ObjectType{AttrTypes: Interfaces_Interface_Aggregated_ether_options_Model{}.AttrTypes()}},
	    "unit": 	types.ListType{ElemType: types.ObjectType{AttrTypes: Interfaces_Interface_Unit_Model{}.AttrTypes()}},
	}
}
func (o Interfaces_Interface_Model) Attributes() map[string]schema.Attribute {
	return map[string]schema.Attribute{
	    "name": schema.StringAttribute{
		    Optional: true,
		    MarkdownDescription: "xpath is `config.Groups.Name.Interface`",
	    },
	    "description": schema.StringAttribute{
		    Optional: true,
		    MarkdownDescription: "xpath is `config.Groups.Description.Interface`",
	    },
	    "vlan_tagging": schema.StringAttribute{
		    Optional: true,
		    MarkdownDescription: "xpath is `config.Groups.Vlan-tagging.Interface`",
	    },
	    "esi": schema.ListNestedAttribute{
		    Optional: true,
		    NestedObject: schema.NestedAttributeObject{
			    Attributes: Interfaces_Interface_Esi_Model{}.Attributes(),
	        },
        },
	    "ether_options": schema.ListNestedAttribute{
		    Optional: true,
		    NestedObject: schema.NestedAttributeObject{
			    Attributes: Interfaces_Interface_Ether_options_Model{}.Attributes(),
	        },
        },
	    "aggregated_ether_options": schema.ListNestedAttribute{
		    Optional: true,
		    NestedObject: schema.NestedAttributeObject{
			    Attributes: Interfaces_Interface_Aggregated_ether_options_Model{}.Attributes(),
	        },
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
	Then	types.List `tfsdk:"then"`
}
func (o Policy_options_Policy_statement_Model) AttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
	    "name": 	types.StringType,
	    "term": 	types.ListType{ElemType: types.ObjectType{AttrTypes: Policy_options_Policy_statement_Term_Model{}.AttrTypes()}},
	    "then": 	types.ListType{ElemType: types.ObjectType{AttrTypes: Policy_options_Policy_statement_Then_Model{}.AttrTypes()}},
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
	    "then": schema.ListNestedAttribute{
		    Optional: true,
		    NestedObject: schema.NestedAttributeObject{
			    Attributes: Policy_options_Policy_statement_Then_Model{}.Attributes(),
	        },
        },
    }
}
type Policy_options_Community_Model struct {
	Name	types.String `tfsdk:"name"`
	Members	types.List `tfsdk:"members"`
}
func (o Policy_options_Community_Model) AttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
	    "name": 	types.StringType,
		"members": 	types.ListType{ElemType: types.StringType},
	}
}
func (o Policy_options_Community_Model) Attributes() map[string]schema.Attribute {
	return map[string]schema.Attribute{
	    "name": schema.StringAttribute{
		    Optional: true,
		    MarkdownDescription: "xpath is `config.Groups.Name.Community`",
	    },
		"members": schema.ListAttribute{
			ElementType:         types.StringType,
			Optional:            true,
			MarkdownDescription: "xpath is `config.Groups.Members.Community`",
		},
    }
}
type Protocols_Bgp_Model struct {
	Group	types.List `tfsdk:"group"`
}
func (o Protocols_Bgp_Model) AttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
	    "group": 	types.ListType{ElemType: types.ObjectType{AttrTypes: Protocols_Bgp_Group_Model{}.AttrTypes()}},
	}
}
func (o Protocols_Bgp_Model) Attributes() map[string]schema.Attribute {
	return map[string]schema.Attribute{
	    "group": schema.ListNestedAttribute{
		    Optional: true,
		    NestedObject: schema.NestedAttributeObject{
			    Attributes: Protocols_Bgp_Group_Model{}.Attributes(),
	        },
        },
    }
}
type Protocols_Evpn_Model struct {
	Encapsulation	types.String `tfsdk:"encapsulation"`
	Multicast_mode	types.String `tfsdk:"multicast_mode"`
	Default_gateway	types.String `tfsdk:"default_gateway"`
	Extended_vni_list	types.List `tfsdk:"extended_vni_list"`
	No_core_isolation	types.String `tfsdk:"no_core_isolation"`
}
func (o Protocols_Evpn_Model) AttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
	    "encapsulation": 	types.StringType,
	    "multicast_mode": 	types.StringType,
	    "default_gateway": 	types.StringType,
		"extended_vni_list": 	types.ListType{ElemType: types.StringType},
	    "no_core_isolation": 	types.StringType,
	}
}
func (o Protocols_Evpn_Model) Attributes() map[string]schema.Attribute {
	return map[string]schema.Attribute{
	    "encapsulation": schema.StringAttribute{
		    Optional: true,
		    MarkdownDescription: "xpath is `config.Groups.Encapsulation.Evpn`",
	    },
	    "multicast_mode": schema.StringAttribute{
		    Optional: true,
		    MarkdownDescription: "xpath is `config.Groups.Multicast-mode.Evpn`",
	    },
	    "default_gateway": schema.StringAttribute{
		    Optional: true,
		    MarkdownDescription: "xpath is `config.Groups.Default-gateway.Evpn`",
	    },
		"extended_vni_list": schema.ListAttribute{
			ElementType:         types.StringType,
			Optional:            true,
			MarkdownDescription: "xpath is `config.Groups.Extended-vni-list.Evpn`",
		},
	    "no_core_isolation": schema.StringAttribute{
		    Optional: true,
		    MarkdownDescription: "xpath is `config.Groups.No-core-isolation.Evpn`",
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
type Protocols_Igmp_snooping_Model struct {
	Vlan	types.List `tfsdk:"vlan"`
}
func (o Protocols_Igmp_snooping_Model) AttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
	    "vlan": 	types.ListType{ElemType: types.ObjectType{AttrTypes: Protocols_Igmp_snooping_Vlan_Model{}.AttrTypes()}},
	}
}
func (o Protocols_Igmp_snooping_Model) Attributes() map[string]schema.Attribute {
	return map[string]schema.Attribute{
	    "vlan": schema.ListNestedAttribute{
		    Optional: true,
		    NestedObject: schema.NestedAttributeObject{
			    Attributes: Protocols_Igmp_snooping_Vlan_Model{}.Attributes(),
	        },
        },
    }
}
type Routing_instances_Instance_Model struct {
	Name	types.String `tfsdk:"name"`
	Instance_type	types.String `tfsdk:"instance_type"`
	Interface	types.List `tfsdk:"interface"`
	Route_distinguisher	types.List `tfsdk:"route_distinguisher"`
	Vrf_target	types.List `tfsdk:"vrf_target"`
	Vrf_table_label	types.List `tfsdk:"vrf_table_label"`
	Routing_options	types.List `tfsdk:"routing_options"`
	Protocols	types.List `tfsdk:"protocols"`
}
func (o Routing_instances_Instance_Model) AttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
	    "name": 	types.StringType,
	    "instance_type": 	types.StringType,
	    "interface": 	types.ListType{ElemType: types.ObjectType{AttrTypes: Routing_instances_Instance_Interface_Model{}.AttrTypes()}},
	    "route_distinguisher": 	types.ListType{ElemType: types.ObjectType{AttrTypes: Routing_instances_Instance_Route_distinguisher_Model{}.AttrTypes()}},
	    "vrf_target": 	types.ListType{ElemType: types.ObjectType{AttrTypes: Routing_instances_Instance_Vrf_target_Model{}.AttrTypes()}},
	    "vrf_table_label": 	types.ListType{ElemType: types.ObjectType{AttrTypes: Routing_instances_Instance_Vrf_table_label_Model{}.AttrTypes()}},
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
	    "route_distinguisher": schema.ListNestedAttribute{
		    Optional: true,
		    NestedObject: schema.NestedAttributeObject{
			    Attributes: Routing_instances_Instance_Route_distinguisher_Model{}.Attributes(),
	        },
        },
	    "vrf_target": schema.ListNestedAttribute{
		    Optional: true,
		    NestedObject: schema.NestedAttributeObject{
			    Attributes: Routing_instances_Instance_Vrf_target_Model{}.Attributes(),
	        },
        },
	    "vrf_table_label": schema.ListNestedAttribute{
		    Optional: true,
		    NestedObject: schema.NestedAttributeObject{
			    Attributes: Routing_instances_Instance_Vrf_table_label_Model{}.Attributes(),
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
type Routing_options_Forwarding_table_Model struct {
	Export	types.List `tfsdk:"export"`
	Ecmp_fast_reroute	types.String `tfsdk:"ecmp_fast_reroute"`
	Chained_composite_next_hop	types.List `tfsdk:"chained_composite_next_hop"`
}
func (o Routing_options_Forwarding_table_Model) AttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"export": 	types.ListType{ElemType: types.StringType},
	    "ecmp_fast_reroute": 	types.StringType,
	    "chained_composite_next_hop": 	types.ListType{ElemType: types.ObjectType{AttrTypes: Routing_options_Forwarding_table_Chained_composite_next_hop_Model{}.AttrTypes()}},
	}
}
func (o Routing_options_Forwarding_table_Model) Attributes() map[string]schema.Attribute {
	return map[string]schema.Attribute{
		"export": schema.ListAttribute{
			ElementType:         types.StringType,
			Optional:            true,
			MarkdownDescription: "xpath is `config.Groups.Export.Forwarding_table`",
		},
	    "ecmp_fast_reroute": schema.StringAttribute{
		    Optional: true,
		    MarkdownDescription: "xpath is `config.Groups.Ecmp-fast-reroute.Forwarding_table`",
	    },
	    "chained_composite_next_hop": schema.ListNestedAttribute{
		    Optional: true,
		    NestedObject: schema.NestedAttributeObject{
			    Attributes: Routing_options_Forwarding_table_Chained_composite_next_hop_Model{}.Attributes(),
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
type Switch_options_Vtep_source_interface_Model struct {
	Interface_name	types.String `tfsdk:"interface_name"`
}
func (o Switch_options_Vtep_source_interface_Model) AttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
	    "interface_name": 	types.StringType,
	}
}
func (o Switch_options_Vtep_source_interface_Model) Attributes() map[string]schema.Attribute {
	return map[string]schema.Attribute{
	    "interface_name": schema.StringAttribute{
		    Optional: true,
		    MarkdownDescription: "xpath is `config.Groups.Interface-name.Vtep_source_interface`",
	    },
    }
}
type Switch_options_Route_distinguisher_Model struct {
	Rd_type	types.String `tfsdk:"rd_type"`
}
func (o Switch_options_Route_distinguisher_Model) AttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
	    "rd_type": 	types.StringType,
	}
}
func (o Switch_options_Route_distinguisher_Model) Attributes() map[string]schema.Attribute {
	return map[string]schema.Attribute{
	    "rd_type": schema.StringAttribute{
		    Optional: true,
		    MarkdownDescription: "xpath is `config.Groups.Rd-type.Route_distinguisher`",
	    },
    }
}
type Switch_options_Vrf_target_Model struct {
	Community	types.String `tfsdk:"community"`
	Auto	types.List `tfsdk:"auto"`
}
func (o Switch_options_Vrf_target_Model) AttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
	    "community": 	types.StringType,
	    "auto": 	types.ListType{ElemType: types.ObjectType{AttrTypes: Switch_options_Vrf_target_Auto_Model{}.AttrTypes()}},
	}
}
func (o Switch_options_Vrf_target_Model) Attributes() map[string]schema.Attribute {
	return map[string]schema.Attribute{
	    "community": schema.StringAttribute{
		    Optional: true,
		    MarkdownDescription: "xpath is `config.Groups.Community.Vrf_target`",
	    },
	    "auto": schema.ListNestedAttribute{
		    Optional: true,
		    NestedObject: schema.NestedAttributeObject{
			    Attributes: Switch_options_Vrf_target_Auto_Model{}.Attributes(),
	        },
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
type System_Extensions_Model struct {
	Providers	types.List `tfsdk:"providers"`
}
func (o System_Extensions_Model) AttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
	    "providers": 	types.ListType{ElemType: types.ObjectType{AttrTypes: System_Extensions_Providers_Model{}.AttrTypes()}},
	}
}
func (o System_Extensions_Model) Attributes() map[string]schema.Attribute {
	return map[string]schema.Attribute{
	    "providers": schema.ListNestedAttribute{
		    Optional: true,
		    NestedObject: schema.NestedAttributeObject{
			    Attributes: System_Extensions_Providers_Model{}.Attributes(),
	        },
        },
    }
}
type Vlans_Vlan_Model struct {
	Name	types.String `tfsdk:"name"`
	Vlan_id	types.String `tfsdk:"vlan_id"`
	L3_interface	types.String `tfsdk:"l3_interface"`
	Vxlan	types.List `tfsdk:"vxlan"`
}
func (o Vlans_Vlan_Model) AttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
	    "name": 	types.StringType,
	    "vlan_id": 	types.StringType,
	    "l3_interface": 	types.StringType,
	    "vxlan": 	types.ListType{ElemType: types.ObjectType{AttrTypes: Vlans_Vlan_Vxlan_Model{}.AttrTypes()}},
	}
}
func (o Vlans_Vlan_Model) Attributes() map[string]schema.Attribute {
	return map[string]schema.Attribute{
	    "name": schema.StringAttribute{
		    Optional: true,
		    MarkdownDescription: "xpath is `config.Groups.Name.Vlan`",
	    },
	    "vlan_id": schema.StringAttribute{
		    Optional: true,
		    MarkdownDescription: "xpath is `config.Groups.Vlan-id.Vlan`",
	    },
	    "l3_interface": schema.StringAttribute{
		    Optional: true,
		    MarkdownDescription: "xpath is `config.Groups.L3-interface.Vlan`",
	    },
	    "vxlan": schema.ListNestedAttribute{
		    Optional: true,
		    NestedObject: schema.NestedAttributeObject{
			    Attributes: Vlans_Vlan_Vxlan_Model{}.Attributes(),
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
type Forwarding_options_Storm_control_profiles_All_Model struct {
}
func (o Forwarding_options_Storm_control_profiles_All_Model) AttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
	}
}
func (o Forwarding_options_Storm_control_profiles_All_Model) Attributes() map[string]schema.Attribute {
	return map[string]schema.Attribute{
    }
}
type Interfaces_Interface_Esi_Model struct {
	Identifier	types.String `tfsdk:"identifier"`
	All_active	types.String `tfsdk:"all_active"`
}
func (o Interfaces_Interface_Esi_Model) AttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
	    "identifier": 	types.StringType,
	    "all_active": 	types.StringType,
	}
}
func (o Interfaces_Interface_Esi_Model) Attributes() map[string]schema.Attribute {
	return map[string]schema.Attribute{
	    "identifier": schema.StringAttribute{
		    Optional: true,
		    MarkdownDescription: "xpath is `config.Groups.Identifier.Esi`",
	    },
	    "all_active": schema.StringAttribute{
		    Optional: true,
		    MarkdownDescription: "xpath is `config.Groups.All-active.Esi`",
	    },
    }
}
type Interfaces_Interface_Ether_options_Model struct {
	Ieee_802_3ad	types.List `tfsdk:"ieee_802_3ad"`
}
func (o Interfaces_Interface_Ether_options_Model) AttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
	    "ieee_802_3ad": 	types.ListType{ElemType: types.ObjectType{AttrTypes: Interfaces_Interface_Ether_options_Ieee_802_3ad_Model{}.AttrTypes()}},
	}
}
func (o Interfaces_Interface_Ether_options_Model) Attributes() map[string]schema.Attribute {
	return map[string]schema.Attribute{
	    "ieee_802_3ad": schema.ListNestedAttribute{
		    Optional: true,
		    NestedObject: schema.NestedAttributeObject{
			    Attributes: Interfaces_Interface_Ether_options_Ieee_802_3ad_Model{}.Attributes(),
	        },
        },
    }
}
type Interfaces_Interface_Aggregated_ether_options_Model struct {
	Lacp	types.List `tfsdk:"lacp"`
}
func (o Interfaces_Interface_Aggregated_ether_options_Model) AttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
	    "lacp": 	types.ListType{ElemType: types.ObjectType{AttrTypes: Interfaces_Interface_Aggregated_ether_options_Lacp_Model{}.AttrTypes()}},
	}
}
func (o Interfaces_Interface_Aggregated_ether_options_Model) Attributes() map[string]schema.Attribute {
	return map[string]schema.Attribute{
	    "lacp": schema.ListNestedAttribute{
		    Optional: true,
		    NestedObject: schema.NestedAttributeObject{
			    Attributes: Interfaces_Interface_Aggregated_ether_options_Lacp_Model{}.Attributes(),
	        },
        },
    }
}
type Interfaces_Interface_Unit_Model struct {
	Name	types.String `tfsdk:"name"`
	Description	types.String `tfsdk:"description"`
	Vlan_id	types.String `tfsdk:"vlan_id"`
	Family	types.List `tfsdk:"family"`
	Mac	types.String `tfsdk:"mac"`
}
func (o Interfaces_Interface_Unit_Model) AttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
	    "name": 	types.StringType,
	    "description": 	types.StringType,
	    "vlan_id": 	types.StringType,
	    "family": 	types.ListType{ElemType: types.ObjectType{AttrTypes: Interfaces_Interface_Unit_Family_Model{}.AttrTypes()}},
	    "mac": 	types.StringType,
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
	    "mac": schema.StringAttribute{
		    Optional: true,
		    MarkdownDescription: "xpath is `config.Groups.Mac.Unit`",
	    },
    }
}
type Policy_options_Policy_statement_Term_Model struct {
	Name	types.String `tfsdk:"name"`
	From	types.List `tfsdk:"from"`
	Then	types.List `tfsdk:"then"`
}
func (o Policy_options_Policy_statement_Term_Model) AttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
	    "name": 	types.StringType,
	    "from": 	types.ListType{ElemType: types.ObjectType{AttrTypes: Policy_options_Policy_statement_Term_From_Model{}.AttrTypes()}},
	    "then": 	types.ListType{ElemType: types.ObjectType{AttrTypes: Policy_options_Policy_statement_Term_Then_Model{}.AttrTypes()}},
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
	    "then": schema.ListNestedAttribute{
		    Optional: true,
		    NestedObject: schema.NestedAttributeObject{
			    Attributes: Policy_options_Policy_statement_Term_Then_Model{}.Attributes(),
	        },
        },
    }
}
type Policy_options_Policy_statement_Then_Model struct {
	Load_balance	types.List `tfsdk:"load_balance"`
}
func (o Policy_options_Policy_statement_Then_Model) AttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
	    "load_balance": 	types.ListType{ElemType: types.ObjectType{AttrTypes: Policy_options_Policy_statement_Then_Load_balance_Model{}.AttrTypes()}},
	}
}
func (o Policy_options_Policy_statement_Then_Model) Attributes() map[string]schema.Attribute {
	return map[string]schema.Attribute{
	    "load_balance": schema.ListNestedAttribute{
		    Optional: true,
		    NestedObject: schema.NestedAttributeObject{
			    Attributes: Policy_options_Policy_statement_Then_Load_balance_Model{}.Attributes(),
	        },
        },
    }
}
type Protocols_Bgp_Group_Model struct {
	Name	types.String `tfsdk:"name"`
	Type	types.String `tfsdk:"type"`
	Multihop	types.List `tfsdk:"multihop"`
	Local_address	types.String `tfsdk:"local_address"`
	Mtu_discovery	types.String `tfsdk:"mtu_discovery"`
	Import	types.List `tfsdk:"import"`
	Family	types.List `tfsdk:"family"`
	Export	types.List `tfsdk:"export"`
	Vpn_apply_export	types.String `tfsdk:"vpn_apply_export"`
	Cluster	types.String `tfsdk:"cluster"`
	Local_as	types.List `tfsdk:"local_as"`
	Multipath	types.List `tfsdk:"multipath"`
	Bfd_liveness_detection	types.List `tfsdk:"bfd_liveness_detection"`
	Allow	types.List `tfsdk:"allow"`
	Neighbor	types.List `tfsdk:"neighbor"`
}
func (o Protocols_Bgp_Group_Model) AttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
	    "name": 	types.StringType,
	    "type": 	types.StringType,
	    "multihop": 	types.ListType{ElemType: types.ObjectType{AttrTypes: Protocols_Bgp_Group_Multihop_Model{}.AttrTypes()}},
	    "local_address": 	types.StringType,
	    "mtu_discovery": 	types.StringType,
		"import": 	types.ListType{ElemType: types.StringType},
	    "family": 	types.ListType{ElemType: types.ObjectType{AttrTypes: Protocols_Bgp_Group_Family_Model{}.AttrTypes()}},
		"export": 	types.ListType{ElemType: types.StringType},
	    "vpn_apply_export": 	types.StringType,
	    "cluster": 	types.StringType,
	    "local_as": 	types.ListType{ElemType: types.ObjectType{AttrTypes: Protocols_Bgp_Group_Local_as_Model{}.AttrTypes()}},
	    "multipath": 	types.ListType{ElemType: types.ObjectType{AttrTypes: Protocols_Bgp_Group_Multipath_Model{}.AttrTypes()}},
	    "bfd_liveness_detection": 	types.ListType{ElemType: types.ObjectType{AttrTypes: Protocols_Bgp_Group_Bfd_liveness_detection_Model{}.AttrTypes()}},
		"allow": 	types.ListType{ElemType: types.StringType},
	    "neighbor": 	types.ListType{ElemType: types.ObjectType{AttrTypes: Protocols_Bgp_Group_Neighbor_Model{}.AttrTypes()}},
	}
}
func (o Protocols_Bgp_Group_Model) Attributes() map[string]schema.Attribute {
	return map[string]schema.Attribute{
	    "name": schema.StringAttribute{
		    Optional: true,
		    MarkdownDescription: "xpath is `config.Groups.Name.Group`",
	    },
	    "type": schema.StringAttribute{
		    Optional: true,
		    MarkdownDescription: "xpath is `config.Groups.Type.Group`",
	    },
	    "multihop": schema.ListNestedAttribute{
		    Optional: true,
		    NestedObject: schema.NestedAttributeObject{
			    Attributes: Protocols_Bgp_Group_Multihop_Model{}.Attributes(),
	        },
        },
	    "local_address": schema.StringAttribute{
		    Optional: true,
		    MarkdownDescription: "xpath is `config.Groups.Local-address.Group`",
	    },
	    "mtu_discovery": schema.StringAttribute{
		    Optional: true,
		    MarkdownDescription: "xpath is `config.Groups.Mtu-discovery.Group`",
	    },
		"import": schema.ListAttribute{
			ElementType:         types.StringType,
			Optional:            true,
			MarkdownDescription: "xpath is `config.Groups.Import.Group`",
		},
	    "family": schema.ListNestedAttribute{
		    Optional: true,
		    NestedObject: schema.NestedAttributeObject{
			    Attributes: Protocols_Bgp_Group_Family_Model{}.Attributes(),
	        },
        },
		"export": schema.ListAttribute{
			ElementType:         types.StringType,
			Optional:            true,
			MarkdownDescription: "xpath is `config.Groups.Export.Group`",
		},
	    "vpn_apply_export": schema.StringAttribute{
		    Optional: true,
		    MarkdownDescription: "xpath is `config.Groups.Vpn-apply-export.Group`",
	    },
	    "cluster": schema.StringAttribute{
		    Optional: true,
		    MarkdownDescription: "xpath is `config.Groups.Cluster.Group`",
	    },
	    "local_as": schema.ListNestedAttribute{
		    Optional: true,
		    NestedObject: schema.NestedAttributeObject{
			    Attributes: Protocols_Bgp_Group_Local_as_Model{}.Attributes(),
	        },
        },
	    "multipath": schema.ListNestedAttribute{
		    Optional: true,
		    NestedObject: schema.NestedAttributeObject{
			    Attributes: Protocols_Bgp_Group_Multipath_Model{}.Attributes(),
	        },
        },
	    "bfd_liveness_detection": schema.ListNestedAttribute{
		    Optional: true,
		    NestedObject: schema.NestedAttributeObject{
			    Attributes: Protocols_Bgp_Group_Bfd_liveness_detection_Model{}.Attributes(),
	        },
        },
		"allow": schema.ListAttribute{
			ElementType:         types.StringType,
			Optional:            true,
			MarkdownDescription: "xpath is `config.Groups.Allow.Group`",
		},
	    "neighbor": schema.ListNestedAttribute{
		    Optional: true,
		    NestedObject: schema.NestedAttributeObject{
			    Attributes: Protocols_Bgp_Group_Neighbor_Model{}.Attributes(),
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
type Protocols_Igmp_snooping_Vlan_Model struct {
	Name	types.String `tfsdk:"name"`
}
func (o Protocols_Igmp_snooping_Vlan_Model) AttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
	    "name": 	types.StringType,
	}
}
func (o Protocols_Igmp_snooping_Vlan_Model) Attributes() map[string]schema.Attribute {
	return map[string]schema.Attribute{
	    "name": schema.StringAttribute{
		    Optional: true,
		    MarkdownDescription: "xpath is `config.Groups.Name.Vlan`",
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
type Routing_instances_Instance_Route_distinguisher_Model struct {
	Rd_type	types.String `tfsdk:"rd_type"`
}
func (o Routing_instances_Instance_Route_distinguisher_Model) AttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
	    "rd_type": 	types.StringType,
	}
}
func (o Routing_instances_Instance_Route_distinguisher_Model) Attributes() map[string]schema.Attribute {
	return map[string]schema.Attribute{
	    "rd_type": schema.StringAttribute{
		    Optional: true,
		    MarkdownDescription: "xpath is `config.Groups.Rd-type.Route_distinguisher`",
	    },
    }
}
type Routing_instances_Instance_Vrf_target_Model struct {
	Community	types.String `tfsdk:"community"`
}
func (o Routing_instances_Instance_Vrf_target_Model) AttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
	    "community": 	types.StringType,
	}
}
func (o Routing_instances_Instance_Vrf_target_Model) Attributes() map[string]schema.Attribute {
	return map[string]schema.Attribute{
	    "community": schema.StringAttribute{
		    Optional: true,
		    MarkdownDescription: "xpath is `config.Groups.Community.Vrf_target`",
	    },
    }
}
type Routing_instances_Instance_Vrf_table_label_Model struct {
}
func (o Routing_instances_Instance_Vrf_table_label_Model) AttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
	}
}
func (o Routing_instances_Instance_Vrf_table_label_Model) Attributes() map[string]schema.Attribute {
	return map[string]schema.Attribute{
    }
}
type Routing_instances_Instance_Routing_options_Model struct {
	Auto_export	types.List `tfsdk:"auto_export"`
}
func (o Routing_instances_Instance_Routing_options_Model) AttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
	    "auto_export": 	types.ListType{ElemType: types.ObjectType{AttrTypes: Routing_instances_Instance_Routing_options_Auto_export_Model{}.AttrTypes()}},
	}
}
func (o Routing_instances_Instance_Routing_options_Model) Attributes() map[string]schema.Attribute {
	return map[string]schema.Attribute{
	    "auto_export": schema.ListNestedAttribute{
		    Optional: true,
		    NestedObject: schema.NestedAttributeObject{
			    Attributes: Routing_instances_Instance_Routing_options_Auto_export_Model{}.Attributes(),
	        },
        },
    }
}
type Routing_instances_Instance_Protocols_Model struct {
	Ospf	types.List `tfsdk:"ospf"`
	Evpn	types.List `tfsdk:"evpn"`
}
func (o Routing_instances_Instance_Protocols_Model) AttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
	    "ospf": 	types.ListType{ElemType: types.ObjectType{AttrTypes: Routing_instances_Instance_Protocols_Ospf_Model{}.AttrTypes()}},
	    "evpn": 	types.ListType{ElemType: types.ObjectType{AttrTypes: Routing_instances_Instance_Protocols_Evpn_Model{}.AttrTypes()}},
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
	    "evpn": schema.ListNestedAttribute{
		    Optional: true,
		    NestedObject: schema.NestedAttributeObject{
			    Attributes: Routing_instances_Instance_Protocols_Evpn_Model{}.Attributes(),
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
type Routing_options_Forwarding_table_Chained_composite_next_hop_Model struct {
	Ingress	types.List `tfsdk:"ingress"`
}
func (o Routing_options_Forwarding_table_Chained_composite_next_hop_Model) AttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
	    "ingress": 	types.ListType{ElemType: types.ObjectType{AttrTypes: Routing_options_Forwarding_table_Chained_composite_next_hop_Ingress_Model{}.AttrTypes()}},
	}
}
func (o Routing_options_Forwarding_table_Chained_composite_next_hop_Model) Attributes() map[string]schema.Attribute {
	return map[string]schema.Attribute{
	    "ingress": schema.ListNestedAttribute{
		    Optional: true,
		    NestedObject: schema.NestedAttributeObject{
			    Attributes: Routing_options_Forwarding_table_Chained_composite_next_hop_Ingress_Model{}.Attributes(),
	        },
        },
    }
}
type Switch_options_Vrf_target_Auto_Model struct {
}
func (o Switch_options_Vrf_target_Auto_Model) AttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
	}
}
func (o Switch_options_Vrf_target_Auto_Model) Attributes() map[string]schema.Attribute {
	return map[string]schema.Attribute{
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
type System_Extensions_Providers_Model struct {
	Name	types.String `tfsdk:"name"`
	License_type	types.List `tfsdk:"license_type"`
}
func (o System_Extensions_Providers_Model) AttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
	    "name": 	types.StringType,
	    "license_type": 	types.ListType{ElemType: types.ObjectType{AttrTypes: System_Extensions_Providers_License_type_Model{}.AttrTypes()}},
	}
}
func (o System_Extensions_Providers_Model) Attributes() map[string]schema.Attribute {
	return map[string]schema.Attribute{
	    "name": schema.StringAttribute{
		    Optional: true,
		    MarkdownDescription: "xpath is `config.Groups.Name.Providers`",
	    },
	    "license_type": schema.ListNestedAttribute{
		    Optional: true,
		    NestedObject: schema.NestedAttributeObject{
			    Attributes: System_Extensions_Providers_License_type_Model{}.Attributes(),
	        },
        },
    }
}
type Vlans_Vlan_Vxlan_Model struct {
	Vni	types.String `tfsdk:"vni"`
}
func (o Vlans_Vlan_Vxlan_Model) AttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
	    "vni": 	types.StringType,
	}
}
func (o Vlans_Vlan_Vxlan_Model) Attributes() map[string]schema.Attribute {
	return map[string]schema.Attribute{
	    "vni": schema.StringAttribute{
		    Optional: true,
		    MarkdownDescription: "xpath is `config.Groups.Vni.Vxlan`",
	    },
    }
}

type Interfaces_Interface_Ether_options_Ieee_802_3ad_Model struct {
	Bundle	types.String `tfsdk:"bundle"`
}
func (o Interfaces_Interface_Ether_options_Ieee_802_3ad_Model) AttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
	    "bundle": 	types.StringType,
	}
}
func (o Interfaces_Interface_Ether_options_Ieee_802_3ad_Model) Attributes() map[string]schema.Attribute {
	return map[string]schema.Attribute{
	    "bundle": schema.StringAttribute{
		    Optional: true,
		    MarkdownDescription: "xpath is `config.Groups.Bundle.Ieee_802_3ad`",
	    },
    }
}
type Interfaces_Interface_Aggregated_ether_options_Lacp_Model struct {
	Active	types.String `tfsdk:"active"`
	Periodic	types.String `tfsdk:"periodic"`
	System_id	types.String `tfsdk:"system_id"`
}
func (o Interfaces_Interface_Aggregated_ether_options_Lacp_Model) AttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
	    "active": 	types.StringType,
	    "periodic": 	types.StringType,
	    "system_id": 	types.StringType,
	}
}
func (o Interfaces_Interface_Aggregated_ether_options_Lacp_Model) Attributes() map[string]schema.Attribute {
	return map[string]schema.Attribute{
	    "active": schema.StringAttribute{
		    Optional: true,
		    MarkdownDescription: "xpath is `config.Groups.Active.Lacp`",
	    },
	    "periodic": schema.StringAttribute{
		    Optional: true,
		    MarkdownDescription: "xpath is `config.Groups.Periodic.Lacp`",
	    },
	    "system_id": schema.StringAttribute{
		    Optional: true,
		    MarkdownDescription: "xpath is `config.Groups.System-id.Lacp`",
	    },
    }
}
type Interfaces_Interface_Unit_Family_Model struct {
	Inet	types.List `tfsdk:"inet"`
	Ethernet_switching	types.List `tfsdk:"ethernet_switching"`
}
func (o Interfaces_Interface_Unit_Family_Model) AttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
	    "inet": 	types.ListType{ElemType: types.ObjectType{AttrTypes: Interfaces_Interface_Unit_Family_Inet_Model{}.AttrTypes()}},
	    "ethernet_switching": 	types.ListType{ElemType: types.ObjectType{AttrTypes: Interfaces_Interface_Unit_Family_Ethernet_switching_Model{}.AttrTypes()}},
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
	    "ethernet_switching": schema.ListNestedAttribute{
		    Optional: true,
		    NestedObject: schema.NestedAttributeObject{
			    Attributes: Interfaces_Interface_Unit_Family_Ethernet_switching_Model{}.Attributes(),
	        },
        },
    }
}
type Policy_options_Policy_statement_Term_From_Model struct {
	Protocol	types.List `tfsdk:"protocol"`
	Route_filter	types.List `tfsdk:"route_filter"`
}
func (o Policy_options_Policy_statement_Term_From_Model) AttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"protocol": 	types.ListType{ElemType: types.StringType},
	    "route_filter": 	types.ListType{ElemType: types.ObjectType{AttrTypes: Policy_options_Policy_statement_Term_From_Route_filter_Model{}.AttrTypes()}},
	}
}
func (o Policy_options_Policy_statement_Term_From_Model) Attributes() map[string]schema.Attribute {
	return map[string]schema.Attribute{
		"protocol": schema.ListAttribute{
			ElementType:         types.StringType,
			Optional:            true,
			MarkdownDescription: "xpath is `config.Groups.Protocol.From`",
		},
	    "route_filter": schema.ListNestedAttribute{
		    Optional: true,
		    NestedObject: schema.NestedAttributeObject{
			    Attributes: Policy_options_Policy_statement_Term_From_Route_filter_Model{}.Attributes(),
	        },
        },
    }
}
type Policy_options_Policy_statement_Term_Then_Model struct {
	Community	types.List `tfsdk:"community"`
	Accept	types.String `tfsdk:"accept"`
	Reject	types.String `tfsdk:"reject"`
}
func (o Policy_options_Policy_statement_Term_Then_Model) AttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
	    "community": 	types.ListType{ElemType: types.ObjectType{AttrTypes: Policy_options_Policy_statement_Term_Then_Community_Model{}.AttrTypes()}},
	    "accept": 	types.StringType,
	    "reject": 	types.StringType,
	}
}
func (o Policy_options_Policy_statement_Term_Then_Model) Attributes() map[string]schema.Attribute {
	return map[string]schema.Attribute{
	    "community": schema.ListNestedAttribute{
		    Optional: true,
		    NestedObject: schema.NestedAttributeObject{
			    Attributes: Policy_options_Policy_statement_Term_Then_Community_Model{}.Attributes(),
	        },
        },
	    "accept": schema.StringAttribute{
		    Optional: true,
		    MarkdownDescription: "xpath is `config.Groups.Accept.Then`",
	    },
	    "reject": schema.StringAttribute{
		    Optional: true,
		    MarkdownDescription: "xpath is `config.Groups.Reject.Then`",
	    },
    }
}
type Policy_options_Policy_statement_Then_Load_balance_Model struct {
	Per_packet	types.String `tfsdk:"per_packet"`
}
func (o Policy_options_Policy_statement_Then_Load_balance_Model) AttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
	    "per_packet": 	types.StringType,
	}
}
func (o Policy_options_Policy_statement_Then_Load_balance_Model) Attributes() map[string]schema.Attribute {
	return map[string]schema.Attribute{
	    "per_packet": schema.StringAttribute{
		    Optional: true,
		    MarkdownDescription: "xpath is `config.Groups.Per-packet.Load_balance`",
	    },
    }
}
type Protocols_Bgp_Group_Multihop_Model struct {
	No_nexthop_change	types.String `tfsdk:"no_nexthop_change"`
}
func (o Protocols_Bgp_Group_Multihop_Model) AttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
	    "no_nexthop_change": 	types.StringType,
	}
}
func (o Protocols_Bgp_Group_Multihop_Model) Attributes() map[string]schema.Attribute {
	return map[string]schema.Attribute{
	    "no_nexthop_change": schema.StringAttribute{
		    Optional: true,
		    MarkdownDescription: "xpath is `config.Groups.No-nexthop-change.Multihop`",
	    },
    }
}
type Protocols_Bgp_Group_Family_Model struct {
	Evpn	types.List `tfsdk:"evpn"`
}
func (o Protocols_Bgp_Group_Family_Model) AttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
	    "evpn": 	types.ListType{ElemType: types.ObjectType{AttrTypes: Protocols_Bgp_Group_Family_Evpn_Model{}.AttrTypes()}},
	}
}
func (o Protocols_Bgp_Group_Family_Model) Attributes() map[string]schema.Attribute {
	return map[string]schema.Attribute{
	    "evpn": schema.ListNestedAttribute{
		    Optional: true,
		    NestedObject: schema.NestedAttributeObject{
			    Attributes: Protocols_Bgp_Group_Family_Evpn_Model{}.Attributes(),
	        },
        },
    }
}
type Protocols_Bgp_Group_Local_as_Model struct {
	As_number	types.String `tfsdk:"as_number"`
}
func (o Protocols_Bgp_Group_Local_as_Model) AttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
	    "as_number": 	types.StringType,
	}
}
func (o Protocols_Bgp_Group_Local_as_Model) Attributes() map[string]schema.Attribute {
	return map[string]schema.Attribute{
	    "as_number": schema.StringAttribute{
		    Optional: true,
		    MarkdownDescription: "xpath is `config.Groups.As-number.Local_as`",
	    },
    }
}
type Protocols_Bgp_Group_Multipath_Model struct {
	Multiple_as	types.String `tfsdk:"multiple_as"`
}
func (o Protocols_Bgp_Group_Multipath_Model) AttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
	    "multiple_as": 	types.StringType,
	}
}
func (o Protocols_Bgp_Group_Multipath_Model) Attributes() map[string]schema.Attribute {
	return map[string]schema.Attribute{
	    "multiple_as": schema.StringAttribute{
		    Optional: true,
		    MarkdownDescription: "xpath is `config.Groups.Multiple-as.Multipath`",
	    },
    }
}
type Protocols_Bgp_Group_Bfd_liveness_detection_Model struct {
	Minimum_interval	types.String `tfsdk:"minimum_interval"`
	Multiplier	types.String `tfsdk:"multiplier"`
}
func (o Protocols_Bgp_Group_Bfd_liveness_detection_Model) AttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
	    "minimum_interval": 	types.StringType,
	    "multiplier": 	types.StringType,
	}
}
func (o Protocols_Bgp_Group_Bfd_liveness_detection_Model) Attributes() map[string]schema.Attribute {
	return map[string]schema.Attribute{
	    "minimum_interval": schema.StringAttribute{
		    Optional: true,
		    MarkdownDescription: "xpath is `config.Groups.Minimum-interval.Bfd_liveness_detection`",
	    },
	    "multiplier": schema.StringAttribute{
		    Optional: true,
		    MarkdownDescription: "xpath is `config.Groups.Multiplier.Bfd_liveness_detection`",
	    },
    }
}
type Protocols_Bgp_Group_Neighbor_Model struct {
	Name	types.String `tfsdk:"name"`
	Description	types.String `tfsdk:"description"`
	Peer_as	types.String `tfsdk:"peer_as"`
}
func (o Protocols_Bgp_Group_Neighbor_Model) AttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
	    "name": 	types.StringType,
	    "description": 	types.StringType,
	    "peer_as": 	types.StringType,
	}
}
func (o Protocols_Bgp_Group_Neighbor_Model) Attributes() map[string]schema.Attribute {
	return map[string]schema.Attribute{
	    "name": schema.StringAttribute{
		    Optional: true,
		    MarkdownDescription: "xpath is `config.Groups.Name.Neighbor`",
	    },
	    "description": schema.StringAttribute{
		    Optional: true,
		    MarkdownDescription: "xpath is `config.Groups.Description.Neighbor`",
	    },
	    "peer_as": schema.StringAttribute{
		    Optional: true,
		    MarkdownDescription: "xpath is `config.Groups.Peer-as.Neighbor`",
	    },
    }
}
type Routing_instances_Instance_Routing_options_Auto_export_Model struct {
}
func (o Routing_instances_Instance_Routing_options_Auto_export_Model) AttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
	}
}
func (o Routing_instances_Instance_Routing_options_Auto_export_Model) Attributes() map[string]schema.Attribute {
	return map[string]schema.Attribute{
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
type Routing_instances_Instance_Protocols_Evpn_Model struct {
	Ip_prefix_routes	types.List `tfsdk:"ip_prefix_routes"`
}
func (o Routing_instances_Instance_Protocols_Evpn_Model) AttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
	    "ip_prefix_routes": 	types.ListType{ElemType: types.ObjectType{AttrTypes: Routing_instances_Instance_Protocols_Evpn_Ip_prefix_routes_Model{}.AttrTypes()}},
	}
}
func (o Routing_instances_Instance_Protocols_Evpn_Model) Attributes() map[string]schema.Attribute {
	return map[string]schema.Attribute{
	    "ip_prefix_routes": schema.ListNestedAttribute{
		    Optional: true,
		    NestedObject: schema.NestedAttributeObject{
			    Attributes: Routing_instances_Instance_Protocols_Evpn_Ip_prefix_routes_Model{}.Attributes(),
	        },
        },
    }
}
type Routing_options_Forwarding_table_Chained_composite_next_hop_Ingress_Model struct {
	Evpn	types.String `tfsdk:"evpn"`
}
func (o Routing_options_Forwarding_table_Chained_composite_next_hop_Ingress_Model) AttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
	    "evpn": 	types.StringType,
	}
}
func (o Routing_options_Forwarding_table_Chained_composite_next_hop_Ingress_Model) Attributes() map[string]schema.Attribute {
	return map[string]schema.Attribute{
	    "evpn": schema.StringAttribute{
		    Optional: true,
		    MarkdownDescription: "xpath is `config.Groups.Evpn.Ingress`",
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
	Notice	types.String `tfsdk:"notice"`
	Info	types.String `tfsdk:"info"`
}
func (o System_Syslog_File_Contents_Model) AttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
	    "name": 	types.StringType,
	    "any": 	types.StringType,
	    "notice": 	types.StringType,
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
	    "notice": schema.StringAttribute{
		    Optional: true,
		    MarkdownDescription: "xpath is `config.Groups.Notice.Contents`",
	    },
	    "info": schema.StringAttribute{
		    Optional: true,
		    MarkdownDescription: "xpath is `config.Groups.Info.Contents`",
	    },
    }
}
type System_Extensions_Providers_License_type_Model struct {
	Name	types.String `tfsdk:"name"`
	Deployment_scope	types.List `tfsdk:"deployment_scope"`
}
func (o System_Extensions_Providers_License_type_Model) AttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
	    "name": 	types.StringType,
		"deployment_scope": 	types.ListType{ElemType: types.StringType},
	}
}
func (o System_Extensions_Providers_License_type_Model) Attributes() map[string]schema.Attribute {
	return map[string]schema.Attribute{
	    "name": schema.StringAttribute{
		    Optional: true,
		    MarkdownDescription: "xpath is `config.Groups.Name.License_type`",
	    },
		"deployment_scope": schema.ListAttribute{
			ElementType:         types.StringType,
			Optional:            true,
			MarkdownDescription: "xpath is `config.Groups.Deployment-scope.License_type`",
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
type Interfaces_Interface_Unit_Family_Ethernet_switching_Model struct {
	Vlan	types.List `tfsdk:"vlan"`
}
func (o Interfaces_Interface_Unit_Family_Ethernet_switching_Model) AttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
	    "vlan": 	types.ListType{ElemType: types.ObjectType{AttrTypes: Interfaces_Interface_Unit_Family_Ethernet_switching_Vlan_Model{}.AttrTypes()}},
	}
}
func (o Interfaces_Interface_Unit_Family_Ethernet_switching_Model) Attributes() map[string]schema.Attribute {
	return map[string]schema.Attribute{
	    "vlan": schema.ListNestedAttribute{
		    Optional: true,
		    NestedObject: schema.NestedAttributeObject{
			    Attributes: Interfaces_Interface_Unit_Family_Ethernet_switching_Vlan_Model{}.Attributes(),
	        },
        },
    }
}
type Policy_options_Policy_statement_Term_From_Route_filter_Model struct {
	Address	types.String `tfsdk:"address"`
	Exact	types.String `tfsdk:"exact"`
	Orlonger	types.String `tfsdk:"orlonger"`
	Prefix_length_range	types.String `tfsdk:"prefix_length_range"`
}
func (o Policy_options_Policy_statement_Term_From_Route_filter_Model) AttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
	    "address": 	types.StringType,
	    "exact": 	types.StringType,
	    "orlonger": 	types.StringType,
	    "prefix_length_range": 	types.StringType,
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
	    "orlonger": schema.StringAttribute{
		    Optional: true,
		    MarkdownDescription: "xpath is `config.Groups.Orlonger.Route_filter`",
	    },
	    "prefix_length_range": schema.StringAttribute{
		    Optional: true,
		    MarkdownDescription: "xpath is `config.Groups.Prefix-length-range.Route_filter`",
	    },
    }
}
type Policy_options_Policy_statement_Term_Then_Community_Model struct {
	Add	types.String `tfsdk:"add"`
	Community_name	types.String `tfsdk:"community_name"`
}
func (o Policy_options_Policy_statement_Term_Then_Community_Model) AttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
	    "add": 	types.StringType,
	    "community_name": 	types.StringType,
	}
}
func (o Policy_options_Policy_statement_Term_Then_Community_Model) Attributes() map[string]schema.Attribute {
	return map[string]schema.Attribute{
	    "add": schema.StringAttribute{
		    Optional: true,
		    MarkdownDescription: "xpath is `config.Groups.Add.Community`",
	    },
	    "community_name": schema.StringAttribute{
		    Optional: true,
		    MarkdownDescription: "xpath is `config.Groups.Community-name.Community`",
	    },
    }
}
type Protocols_Bgp_Group_Family_Evpn_Model struct {
	Signaling	types.List `tfsdk:"signaling"`
}
func (o Protocols_Bgp_Group_Family_Evpn_Model) AttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
	    "signaling": 	types.ListType{ElemType: types.ObjectType{AttrTypes: Protocols_Bgp_Group_Family_Evpn_Signaling_Model{}.AttrTypes()}},
	}
}
func (o Protocols_Bgp_Group_Family_Evpn_Model) Attributes() map[string]schema.Attribute {
	return map[string]schema.Attribute{
	    "signaling": schema.ListNestedAttribute{
		    Optional: true,
		    NestedObject: schema.NestedAttributeObject{
			    Attributes: Protocols_Bgp_Group_Family_Evpn_Signaling_Model{}.Attributes(),
	        },
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
type Routing_instances_Instance_Protocols_Evpn_Ip_prefix_routes_Model struct {
	Advertise	types.String `tfsdk:"advertise"`
	Encapsulation	types.String `tfsdk:"encapsulation"`
	Vni	types.String `tfsdk:"vni"`
	Export	types.List `tfsdk:"export"`
}
func (o Routing_instances_Instance_Protocols_Evpn_Ip_prefix_routes_Model) AttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
	    "advertise": 	types.StringType,
	    "encapsulation": 	types.StringType,
	    "vni": 	types.StringType,
		"export": 	types.ListType{ElemType: types.StringType},
	}
}
func (o Routing_instances_Instance_Protocols_Evpn_Ip_prefix_routes_Model) Attributes() map[string]schema.Attribute {
	return map[string]schema.Attribute{
	    "advertise": schema.StringAttribute{
		    Optional: true,
		    MarkdownDescription: "xpath is `config.Groups.Advertise.Ip_prefix_routes`",
	    },
	    "encapsulation": schema.StringAttribute{
		    Optional: true,
		    MarkdownDescription: "xpath is `config.Groups.Encapsulation.Ip_prefix_routes`",
	    },
	    "vni": schema.StringAttribute{
		    Optional: true,
		    MarkdownDescription: "xpath is `config.Groups.Vni.Ip_prefix_routes`",
	    },
		"export": schema.ListAttribute{
			ElementType:         types.StringType,
			Optional:            true,
			MarkdownDescription: "xpath is `config.Groups.Export.Ip_prefix_routes`",
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
type Interfaces_Interface_Unit_Family_Ethernet_switching_Vlan_Model struct {
	Members	types.List `tfsdk:"members"`
}
func (o Interfaces_Interface_Unit_Family_Ethernet_switching_Vlan_Model) AttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"members": 	types.ListType{ElemType: types.StringType},
	}
}
func (o Interfaces_Interface_Unit_Family_Ethernet_switching_Vlan_Model) Attributes() map[string]schema.Attribute {
	return map[string]schema.Attribute{
		"members": schema.ListAttribute{
			ElementType:         types.StringType,
			Optional:            true,
			MarkdownDescription: "xpath is `config.Groups.Members.Vlan`",
		},
    }
}
type Protocols_Bgp_Group_Family_Evpn_Signaling_Model struct {
	Delay_route_advertisements	types.List `tfsdk:"delay_route_advertisements"`
}
func (o Protocols_Bgp_Group_Family_Evpn_Signaling_Model) AttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
	    "delay_route_advertisements": 	types.ListType{ElemType: types.ObjectType{AttrTypes: Protocols_Bgp_Group_Family_Evpn_Signaling_Delay_route_advertisements_Model{}.AttrTypes()}},
	}
}
func (o Protocols_Bgp_Group_Family_Evpn_Signaling_Model) Attributes() map[string]schema.Attribute {
	return map[string]schema.Attribute{
	    "delay_route_advertisements": schema.ListNestedAttribute{
		    Optional: true,
		    NestedObject: schema.NestedAttributeObject{
			    Attributes: Protocols_Bgp_Group_Family_Evpn_Signaling_Delay_route_advertisements_Model{}.Attributes(),
	        },
        },
    }
}
type Routing_instances_Instance_Protocols_Ospf_Area_Interface_Model struct {
	Name	types.String `tfsdk:"name"`
	Metric	types.String `tfsdk:"metric"`
}
func (o Routing_instances_Instance_Protocols_Ospf_Area_Interface_Model) AttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
	    "name": 	types.StringType,
	    "metric": 	types.StringType,
	}
}
func (o Routing_instances_Instance_Protocols_Ospf_Area_Interface_Model) Attributes() map[string]schema.Attribute {
	return map[string]schema.Attribute{
	    "name": schema.StringAttribute{
		    Optional: true,
		    MarkdownDescription: "xpath is `config.Groups.Name.Interface`",
	    },
	    "metric": schema.StringAttribute{
		    Optional: true,
		    MarkdownDescription: "xpath is `config.Groups.Metric.Interface`",
	    },
    }
}

type Protocols_Bgp_Group_Family_Evpn_Signaling_Delay_route_advertisements_Model struct {
	Minimum_delay	types.List `tfsdk:"minimum_delay"`
}
func (o Protocols_Bgp_Group_Family_Evpn_Signaling_Delay_route_advertisements_Model) AttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
	    "minimum_delay": 	types.ListType{ElemType: types.ObjectType{AttrTypes: Protocols_Bgp_Group_Family_Evpn_Signaling_Delay_route_advertisements_Minimum_delay_Model{}.AttrTypes()}},
	}
}
func (o Protocols_Bgp_Group_Family_Evpn_Signaling_Delay_route_advertisements_Model) Attributes() map[string]schema.Attribute {
	return map[string]schema.Attribute{
	    "minimum_delay": schema.ListNestedAttribute{
		    Optional: true,
		    NestedObject: schema.NestedAttributeObject{
			    Attributes: Protocols_Bgp_Group_Family_Evpn_Signaling_Delay_route_advertisements_Minimum_delay_Model{}.Attributes(),
	        },
        },
    }
}

type Protocols_Bgp_Group_Family_Evpn_Signaling_Delay_route_advertisements_Minimum_delay_Model struct {
	Routing_uptime	types.String `tfsdk:"routing_uptime"`
}
func (o Protocols_Bgp_Group_Family_Evpn_Signaling_Delay_route_advertisements_Minimum_delay_Model) AttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
	    "routing_uptime": 	types.StringType,
	}
}
func (o Protocols_Bgp_Group_Family_Evpn_Signaling_Delay_route_advertisements_Minimum_delay_Model) Attributes() map[string]schema.Attribute {
	return map[string]schema.Attribute{
	    "routing_uptime": schema.StringAttribute{
		    Optional: true,
		    MarkdownDescription: "xpath is `config.Groups.Routing-uptime.Minimum_delay`",
	    },
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
	resp.TypeName = req.ProviderTypeName + "_Apply_Groups"
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
			"forwarding_options": schema.ListNestedAttribute{
				Optional: true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: Forwarding_options_Model{}.Attributes(),
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
			"snmp": schema.ListNestedAttribute{
				Optional: true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: Snmp_Model{}.Attributes(),
				},
			},
			"switch_options": schema.ListNestedAttribute{
				Optional: true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: Switch_options_Model{}.Attributes(),
				},
			},
			"system": schema.ListNestedAttribute{
				Optional: true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: System_Model{}.Attributes(),
				},
			},
			"vlans": schema.ListNestedAttribute{
				Optional: true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: Vlans_Model{}.Attributes(),
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
	
    var var_forwarding_options []Forwarding_options_Model
    if plan.Forwarding_options.IsNull() {
        var_forwarding_options = []Forwarding_options_Model{}
    }else {
        resp.Diagnostics.Append(plan.Forwarding_options.ElementsAs(ctx, &var_forwarding_options, false)...)
        if resp.Diagnostics.HasError() {
            return
        }
    }
    config.Groups.Forwarding_options = make([]xml_Forwarding_options, len(var_forwarding_options))
   
    for i_forwarding_options, v_forwarding_options := range var_forwarding_options {
        var var_forwarding_options_storm_control_profiles []Forwarding_options_Storm_control_profiles_Model
        resp.Diagnostics.Append(v_forwarding_options.Storm_control_profiles.ElementsAs(ctx, &var_forwarding_options_storm_control_profiles, false)...)
        if resp.Diagnostics.HasError() {
            return
        }
	    config.Groups.Forwarding_options[i_forwarding_options].Storm_control_profiles = make([]xml_Forwarding_options_Storm_control_profiles, len(var_forwarding_options_storm_control_profiles))
        
		for i_forwarding_options_storm_control_profiles, v_forwarding_options_storm_control_profiles := range var_forwarding_options_storm_control_profiles {
            config.Groups.Forwarding_options[i_forwarding_options].Storm_control_profiles[i_forwarding_options_storm_control_profiles].Name = v_forwarding_options_storm_control_profiles.Name.ValueStringPointer()
            var var_forwarding_options_storm_control_profiles_all []Forwarding_options_Storm_control_profiles_All_Model
            resp.Diagnostics.Append(v_forwarding_options_storm_control_profiles.All.ElementsAs(ctx, &var_forwarding_options_storm_control_profiles_all, false)...)
            if resp.Diagnostics.HasError() {
                return
            }
	    config.Groups.Forwarding_options[i_forwarding_options].Storm_control_profiles[i_forwarding_options_storm_control_profiles].All = make([]xml_Forwarding_options_Storm_control_profiles_All, len(var_forwarding_options_storm_control_profiles_all))
        
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
            config.Groups.Interfaces[i_interfaces].Interface[i_interfaces_interface].Description = v_interfaces_interface.Description.ValueStringPointer()
            config.Groups.Interfaces[i_interfaces].Interface[i_interfaces_interface].Vlan_tagging = v_interfaces_interface.Vlan_tagging.ValueStringPointer()
            var var_interfaces_interface_esi []Interfaces_Interface_Esi_Model
            resp.Diagnostics.Append(v_interfaces_interface.Esi.ElementsAs(ctx, &var_interfaces_interface_esi, false)...)
            if resp.Diagnostics.HasError() {
                return
            }
	    config.Groups.Interfaces[i_interfaces].Interface[i_interfaces_interface].Esi = make([]xml_Interfaces_Interface_Esi, len(var_interfaces_interface_esi))
        
		for i_interfaces_interface_esi, v_interfaces_interface_esi := range var_interfaces_interface_esi {
            config.Groups.Interfaces[i_interfaces].Interface[i_interfaces_interface].Esi[i_interfaces_interface_esi].Identifier = v_interfaces_interface_esi.Identifier.ValueStringPointer()
            config.Groups.Interfaces[i_interfaces].Interface[i_interfaces_interface].Esi[i_interfaces_interface_esi].All_active = v_interfaces_interface_esi.All_active.ValueStringPointer()
        }
            var var_interfaces_interface_ether_options []Interfaces_Interface_Ether_options_Model
            resp.Diagnostics.Append(v_interfaces_interface.Ether_options.ElementsAs(ctx, &var_interfaces_interface_ether_options, false)...)
            if resp.Diagnostics.HasError() {
                return
            }
	    config.Groups.Interfaces[i_interfaces].Interface[i_interfaces_interface].Ether_options = make([]xml_Interfaces_Interface_Ether_options, len(var_interfaces_interface_ether_options))
        
		for i_interfaces_interface_ether_options, v_interfaces_interface_ether_options := range var_interfaces_interface_ether_options {
            var var_interfaces_interface_ether_options_ieee_802_3ad []Interfaces_Interface_Ether_options_Ieee_802_3ad_Model
            resp.Diagnostics.Append(v_interfaces_interface_ether_options.Ieee_802_3ad.ElementsAs(ctx, &var_interfaces_interface_ether_options_ieee_802_3ad, false)...)
            if resp.Diagnostics.HasError() {
                return
            }
	    config.Groups.Interfaces[i_interfaces].Interface[i_interfaces_interface].Ether_options[i_interfaces_interface_ether_options].Ieee_802_3ad = make([]xml_Interfaces_Interface_Ether_options_Ieee_802_3ad, len(var_interfaces_interface_ether_options_ieee_802_3ad))
        
		for i_interfaces_interface_ether_options_ieee_802_3ad, v_interfaces_interface_ether_options_ieee_802_3ad := range var_interfaces_interface_ether_options_ieee_802_3ad {
            config.Groups.Interfaces[i_interfaces].Interface[i_interfaces_interface].Ether_options[i_interfaces_interface_ether_options].Ieee_802_3ad[i_interfaces_interface_ether_options_ieee_802_3ad].Bundle = v_interfaces_interface_ether_options_ieee_802_3ad.Bundle.ValueStringPointer()
        }
        }
            var var_interfaces_interface_aggregated_ether_options []Interfaces_Interface_Aggregated_ether_options_Model
            resp.Diagnostics.Append(v_interfaces_interface.Aggregated_ether_options.ElementsAs(ctx, &var_interfaces_interface_aggregated_ether_options, false)...)
            if resp.Diagnostics.HasError() {
                return
            }
	    config.Groups.Interfaces[i_interfaces].Interface[i_interfaces_interface].Aggregated_ether_options = make([]xml_Interfaces_Interface_Aggregated_ether_options, len(var_interfaces_interface_aggregated_ether_options))
        
		for i_interfaces_interface_aggregated_ether_options, v_interfaces_interface_aggregated_ether_options := range var_interfaces_interface_aggregated_ether_options {
            var var_interfaces_interface_aggregated_ether_options_lacp []Interfaces_Interface_Aggregated_ether_options_Lacp_Model
            resp.Diagnostics.Append(v_interfaces_interface_aggregated_ether_options.Lacp.ElementsAs(ctx, &var_interfaces_interface_aggregated_ether_options_lacp, false)...)
            if resp.Diagnostics.HasError() {
                return
            }
	    config.Groups.Interfaces[i_interfaces].Interface[i_interfaces_interface].Aggregated_ether_options[i_interfaces_interface_aggregated_ether_options].Lacp = make([]xml_Interfaces_Interface_Aggregated_ether_options_Lacp, len(var_interfaces_interface_aggregated_ether_options_lacp))
        
		for i_interfaces_interface_aggregated_ether_options_lacp, v_interfaces_interface_aggregated_ether_options_lacp := range var_interfaces_interface_aggregated_ether_options_lacp {
            config.Groups.Interfaces[i_interfaces].Interface[i_interfaces_interface].Aggregated_ether_options[i_interfaces_interface_aggregated_ether_options].Lacp[i_interfaces_interface_aggregated_ether_options_lacp].Active = v_interfaces_interface_aggregated_ether_options_lacp.Active.ValueStringPointer()
            config.Groups.Interfaces[i_interfaces].Interface[i_interfaces_interface].Aggregated_ether_options[i_interfaces_interface_aggregated_ether_options].Lacp[i_interfaces_interface_aggregated_ether_options_lacp].Periodic = v_interfaces_interface_aggregated_ether_options_lacp.Periodic.ValueStringPointer()
            config.Groups.Interfaces[i_interfaces].Interface[i_interfaces_interface].Aggregated_ether_options[i_interfaces_interface_aggregated_ether_options].Lacp[i_interfaces_interface_aggregated_ether_options_lacp].System_id = v_interfaces_interface_aggregated_ether_options_lacp.System_id.ValueStringPointer()
        }
        }
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
            var var_interfaces_interface_unit_family_ethernet_switching []Interfaces_Interface_Unit_Family_Ethernet_switching_Model
            resp.Diagnostics.Append(v_interfaces_interface_unit_family.Ethernet_switching.ElementsAs(ctx, &var_interfaces_interface_unit_family_ethernet_switching, false)...)
            if resp.Diagnostics.HasError() {
                return
            }
	    config.Groups.Interfaces[i_interfaces].Interface[i_interfaces_interface].Unit[i_interfaces_interface_unit].Family[i_interfaces_interface_unit_family].Ethernet_switching = make([]xml_Interfaces_Interface_Unit_Family_Ethernet_switching, len(var_interfaces_interface_unit_family_ethernet_switching))
        
		for i_interfaces_interface_unit_family_ethernet_switching, v_interfaces_interface_unit_family_ethernet_switching := range var_interfaces_interface_unit_family_ethernet_switching {
            var var_interfaces_interface_unit_family_ethernet_switching_vlan []Interfaces_Interface_Unit_Family_Ethernet_switching_Vlan_Model
            resp.Diagnostics.Append(v_interfaces_interface_unit_family_ethernet_switching.Vlan.ElementsAs(ctx, &var_interfaces_interface_unit_family_ethernet_switching_vlan, false)...)
            if resp.Diagnostics.HasError() {
                return
            }
	    config.Groups.Interfaces[i_interfaces].Interface[i_interfaces_interface].Unit[i_interfaces_interface_unit].Family[i_interfaces_interface_unit_family].Ethernet_switching[i_interfaces_interface_unit_family_ethernet_switching].Vlan = make([]xml_Interfaces_Interface_Unit_Family_Ethernet_switching_Vlan, len(var_interfaces_interface_unit_family_ethernet_switching_vlan))
        
		for i_interfaces_interface_unit_family_ethernet_switching_vlan, v_interfaces_interface_unit_family_ethernet_switching_vlan := range var_interfaces_interface_unit_family_ethernet_switching_vlan {
			var var_interfaces_interface_unit_family_ethernet_switching_vlan_members []string
			resp.Diagnostics.Append(v_interfaces_interface_unit_family_ethernet_switching_vlan.Members.ElementsAs(ctx, &var_interfaces_interface_unit_family_ethernet_switching_vlan_members, false)...)
			for _, v_interfaces_interface_unit_family_ethernet_switching_vlan_members := range var_interfaces_interface_unit_family_ethernet_switching_vlan_members {
				config.Groups.Interfaces[i_interfaces].Interface[i_interfaces_interface].Unit[i_interfaces_interface_unit].Family[i_interfaces_interface_unit_family].Ethernet_switching[i_interfaces_interface_unit_family_ethernet_switching].Vlan[i_interfaces_interface_unit_family_ethernet_switching_vlan].Members = append(config.Groups.Interfaces[i_interfaces].Interface[i_interfaces_interface].Unit[i_interfaces_interface_unit].Family[i_interfaces_interface_unit_family].Ethernet_switching[i_interfaces_interface_unit_family_ethernet_switching].Vlan[i_interfaces_interface_unit_family_ethernet_switching_vlan].Members, &v_interfaces_interface_unit_family_ethernet_switching_vlan_members)
			}
        }
        }
        }
            config.Groups.Interfaces[i_interfaces].Interface[i_interfaces_interface].Unit[i_interfaces_interface_unit].Mac = v_interfaces_interface_unit.Mac.ValueStringPointer()
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
			var var_policy_options_policy_statement_term_from_protocol []string
			resp.Diagnostics.Append(v_policy_options_policy_statement_term_from.Protocol.ElementsAs(ctx, &var_policy_options_policy_statement_term_from_protocol, false)...)
			for _, v_policy_options_policy_statement_term_from_protocol := range var_policy_options_policy_statement_term_from_protocol {
				config.Groups.Policy_options[i_policy_options].Policy_statement[i_policy_options_policy_statement].Term[i_policy_options_policy_statement_term].From[i_policy_options_policy_statement_term_from].Protocol = append(config.Groups.Policy_options[i_policy_options].Policy_statement[i_policy_options_policy_statement].Term[i_policy_options_policy_statement_term].From[i_policy_options_policy_statement_term_from].Protocol, &v_policy_options_policy_statement_term_from_protocol)
			}
            var var_policy_options_policy_statement_term_from_route_filter []Policy_options_Policy_statement_Term_From_Route_filter_Model
            resp.Diagnostics.Append(v_policy_options_policy_statement_term_from.Route_filter.ElementsAs(ctx, &var_policy_options_policy_statement_term_from_route_filter, false)...)
            if resp.Diagnostics.HasError() {
                return
            }
	    config.Groups.Policy_options[i_policy_options].Policy_statement[i_policy_options_policy_statement].Term[i_policy_options_policy_statement_term].From[i_policy_options_policy_statement_term_from].Route_filter = make([]xml_Policy_options_Policy_statement_Term_From_Route_filter, len(var_policy_options_policy_statement_term_from_route_filter))
        
		for i_policy_options_policy_statement_term_from_route_filter, v_policy_options_policy_statement_term_from_route_filter := range var_policy_options_policy_statement_term_from_route_filter {
            config.Groups.Policy_options[i_policy_options].Policy_statement[i_policy_options_policy_statement].Term[i_policy_options_policy_statement_term].From[i_policy_options_policy_statement_term_from].Route_filter[i_policy_options_policy_statement_term_from_route_filter].Address = v_policy_options_policy_statement_term_from_route_filter.Address.ValueStringPointer()
            config.Groups.Policy_options[i_policy_options].Policy_statement[i_policy_options_policy_statement].Term[i_policy_options_policy_statement_term].From[i_policy_options_policy_statement_term_from].Route_filter[i_policy_options_policy_statement_term_from_route_filter].Exact = v_policy_options_policy_statement_term_from_route_filter.Exact.ValueStringPointer()
            config.Groups.Policy_options[i_policy_options].Policy_statement[i_policy_options_policy_statement].Term[i_policy_options_policy_statement_term].From[i_policy_options_policy_statement_term_from].Route_filter[i_policy_options_policy_statement_term_from_route_filter].Orlonger = v_policy_options_policy_statement_term_from_route_filter.Orlonger.ValueStringPointer()
            config.Groups.Policy_options[i_policy_options].Policy_statement[i_policy_options_policy_statement].Term[i_policy_options_policy_statement_term].From[i_policy_options_policy_statement_term_from].Route_filter[i_policy_options_policy_statement_term_from_route_filter].Prefix_length_range = v_policy_options_policy_statement_term_from_route_filter.Prefix_length_range.ValueStringPointer()
        }
        }
            var var_policy_options_policy_statement_term_then []Policy_options_Policy_statement_Term_Then_Model
            resp.Diagnostics.Append(v_policy_options_policy_statement_term.Then.ElementsAs(ctx, &var_policy_options_policy_statement_term_then, false)...)
            if resp.Diagnostics.HasError() {
                return
            }
	    config.Groups.Policy_options[i_policy_options].Policy_statement[i_policy_options_policy_statement].Term[i_policy_options_policy_statement_term].Then = make([]xml_Policy_options_Policy_statement_Term_Then, len(var_policy_options_policy_statement_term_then))
        
		for i_policy_options_policy_statement_term_then, v_policy_options_policy_statement_term_then := range var_policy_options_policy_statement_term_then {
            var var_policy_options_policy_statement_term_then_community []Policy_options_Policy_statement_Term_Then_Community_Model
            resp.Diagnostics.Append(v_policy_options_policy_statement_term_then.Community.ElementsAs(ctx, &var_policy_options_policy_statement_term_then_community, false)...)
            if resp.Diagnostics.HasError() {
                return
            }
	    config.Groups.Policy_options[i_policy_options].Policy_statement[i_policy_options_policy_statement].Term[i_policy_options_policy_statement_term].Then[i_policy_options_policy_statement_term_then].Community = make([]xml_Policy_options_Policy_statement_Term_Then_Community, len(var_policy_options_policy_statement_term_then_community))
        
		for i_policy_options_policy_statement_term_then_community, v_policy_options_policy_statement_term_then_community := range var_policy_options_policy_statement_term_then_community {
            config.Groups.Policy_options[i_policy_options].Policy_statement[i_policy_options_policy_statement].Term[i_policy_options_policy_statement_term].Then[i_policy_options_policy_statement_term_then].Community[i_policy_options_policy_statement_term_then_community].Add = v_policy_options_policy_statement_term_then_community.Add.ValueStringPointer()
            config.Groups.Policy_options[i_policy_options].Policy_statement[i_policy_options_policy_statement].Term[i_policy_options_policy_statement_term].Then[i_policy_options_policy_statement_term_then].Community[i_policy_options_policy_statement_term_then_community].Community_name = v_policy_options_policy_statement_term_then_community.Community_name.ValueStringPointer()
        }
            config.Groups.Policy_options[i_policy_options].Policy_statement[i_policy_options_policy_statement].Term[i_policy_options_policy_statement_term].Then[i_policy_options_policy_statement_term_then].Accept = v_policy_options_policy_statement_term_then.Accept.ValueStringPointer()
            config.Groups.Policy_options[i_policy_options].Policy_statement[i_policy_options_policy_statement].Term[i_policy_options_policy_statement_term].Then[i_policy_options_policy_statement_term_then].Reject = v_policy_options_policy_statement_term_then.Reject.ValueStringPointer()
        }
        }
            var var_policy_options_policy_statement_then []Policy_options_Policy_statement_Then_Model
            resp.Diagnostics.Append(v_policy_options_policy_statement.Then.ElementsAs(ctx, &var_policy_options_policy_statement_then, false)...)
            if resp.Diagnostics.HasError() {
                return
            }
	    config.Groups.Policy_options[i_policy_options].Policy_statement[i_policy_options_policy_statement].Then = make([]xml_Policy_options_Policy_statement_Then, len(var_policy_options_policy_statement_then))
        
		for i_policy_options_policy_statement_then, v_policy_options_policy_statement_then := range var_policy_options_policy_statement_then {
            var var_policy_options_policy_statement_then_load_balance []Policy_options_Policy_statement_Then_Load_balance_Model
            resp.Diagnostics.Append(v_policy_options_policy_statement_then.Load_balance.ElementsAs(ctx, &var_policy_options_policy_statement_then_load_balance, false)...)
            if resp.Diagnostics.HasError() {
                return
            }
	    config.Groups.Policy_options[i_policy_options].Policy_statement[i_policy_options_policy_statement].Then[i_policy_options_policy_statement_then].Load_balance = make([]xml_Policy_options_Policy_statement_Then_Load_balance, len(var_policy_options_policy_statement_then_load_balance))
        
		for i_policy_options_policy_statement_then_load_balance, v_policy_options_policy_statement_then_load_balance := range var_policy_options_policy_statement_then_load_balance {
            config.Groups.Policy_options[i_policy_options].Policy_statement[i_policy_options_policy_statement].Then[i_policy_options_policy_statement_then].Load_balance[i_policy_options_policy_statement_then_load_balance].Per_packet = v_policy_options_policy_statement_then_load_balance.Per_packet.ValueStringPointer()
        }
        }
        }
        var var_policy_options_community []Policy_options_Community_Model
        resp.Diagnostics.Append(v_policy_options.Community.ElementsAs(ctx, &var_policy_options_community, false)...)
        if resp.Diagnostics.HasError() {
            return
        }
	    config.Groups.Policy_options[i_policy_options].Community = make([]xml_Policy_options_Community, len(var_policy_options_community))
        
		for i_policy_options_community, v_policy_options_community := range var_policy_options_community {
            config.Groups.Policy_options[i_policy_options].Community[i_policy_options_community].Name = v_policy_options_community.Name.ValueStringPointer()
			var var_policy_options_community_members []string
			resp.Diagnostics.Append(v_policy_options_community.Members.ElementsAs(ctx, &var_policy_options_community_members, false)...)
			for _, v_policy_options_community_members := range var_policy_options_community_members {
				config.Groups.Policy_options[i_policy_options].Community[i_policy_options_community].Members = append(config.Groups.Policy_options[i_policy_options].Community[i_policy_options_community].Members, &v_policy_options_community_members)
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
        var var_protocols_bgp []Protocols_Bgp_Model
        resp.Diagnostics.Append(v_protocols.Bgp.ElementsAs(ctx, &var_protocols_bgp, false)...)
        if resp.Diagnostics.HasError() {
            return
        }
	    config.Groups.Protocols[i_protocols].Bgp = make([]xml_Protocols_Bgp, len(var_protocols_bgp))
        
		for i_protocols_bgp, v_protocols_bgp := range var_protocols_bgp {
            var var_protocols_bgp_group []Protocols_Bgp_Group_Model
            resp.Diagnostics.Append(v_protocols_bgp.Group.ElementsAs(ctx, &var_protocols_bgp_group, false)...)
            if resp.Diagnostics.HasError() {
                return
            }
	    config.Groups.Protocols[i_protocols].Bgp[i_protocols_bgp].Group = make([]xml_Protocols_Bgp_Group, len(var_protocols_bgp_group))
        
		for i_protocols_bgp_group, v_protocols_bgp_group := range var_protocols_bgp_group {
            config.Groups.Protocols[i_protocols].Bgp[i_protocols_bgp].Group[i_protocols_bgp_group].Name = v_protocols_bgp_group.Name.ValueStringPointer()
            config.Groups.Protocols[i_protocols].Bgp[i_protocols_bgp].Group[i_protocols_bgp_group].Type = v_protocols_bgp_group.Type.ValueStringPointer()
            var var_protocols_bgp_group_multihop []Protocols_Bgp_Group_Multihop_Model
            resp.Diagnostics.Append(v_protocols_bgp_group.Multihop.ElementsAs(ctx, &var_protocols_bgp_group_multihop, false)...)
            if resp.Diagnostics.HasError() {
                return
            }
	    config.Groups.Protocols[i_protocols].Bgp[i_protocols_bgp].Group[i_protocols_bgp_group].Multihop = make([]xml_Protocols_Bgp_Group_Multihop, len(var_protocols_bgp_group_multihop))
        
		for i_protocols_bgp_group_multihop, v_protocols_bgp_group_multihop := range var_protocols_bgp_group_multihop {
            config.Groups.Protocols[i_protocols].Bgp[i_protocols_bgp].Group[i_protocols_bgp_group].Multihop[i_protocols_bgp_group_multihop].No_nexthop_change = v_protocols_bgp_group_multihop.No_nexthop_change.ValueStringPointer()
        }
            config.Groups.Protocols[i_protocols].Bgp[i_protocols_bgp].Group[i_protocols_bgp_group].Local_address = v_protocols_bgp_group.Local_address.ValueStringPointer()
            config.Groups.Protocols[i_protocols].Bgp[i_protocols_bgp].Group[i_protocols_bgp_group].Mtu_discovery = v_protocols_bgp_group.Mtu_discovery.ValueStringPointer()
			var var_protocols_bgp_group_import []string
			resp.Diagnostics.Append(v_protocols_bgp_group.Import.ElementsAs(ctx, &var_protocols_bgp_group_import, false)...)
			for _, v_protocols_bgp_group_import := range var_protocols_bgp_group_import {
				config.Groups.Protocols[i_protocols].Bgp[i_protocols_bgp].Group[i_protocols_bgp_group].Import = append(config.Groups.Protocols[i_protocols].Bgp[i_protocols_bgp].Group[i_protocols_bgp_group].Import, &v_protocols_bgp_group_import)
			}
            var var_protocols_bgp_group_family []Protocols_Bgp_Group_Family_Model
            resp.Diagnostics.Append(v_protocols_bgp_group.Family.ElementsAs(ctx, &var_protocols_bgp_group_family, false)...)
            if resp.Diagnostics.HasError() {
                return
            }
	    config.Groups.Protocols[i_protocols].Bgp[i_protocols_bgp].Group[i_protocols_bgp_group].Family = make([]xml_Protocols_Bgp_Group_Family, len(var_protocols_bgp_group_family))
        
		for i_protocols_bgp_group_family, v_protocols_bgp_group_family := range var_protocols_bgp_group_family {
            var var_protocols_bgp_group_family_evpn []Protocols_Bgp_Group_Family_Evpn_Model
            resp.Diagnostics.Append(v_protocols_bgp_group_family.Evpn.ElementsAs(ctx, &var_protocols_bgp_group_family_evpn, false)...)
            if resp.Diagnostics.HasError() {
                return
            }
	    config.Groups.Protocols[i_protocols].Bgp[i_protocols_bgp].Group[i_protocols_bgp_group].Family[i_protocols_bgp_group_family].Evpn = make([]xml_Protocols_Bgp_Group_Family_Evpn, len(var_protocols_bgp_group_family_evpn))
        
		for i_protocols_bgp_group_family_evpn, v_protocols_bgp_group_family_evpn := range var_protocols_bgp_group_family_evpn {
            var var_protocols_bgp_group_family_evpn_signaling []Protocols_Bgp_Group_Family_Evpn_Signaling_Model
            resp.Diagnostics.Append(v_protocols_bgp_group_family_evpn.Signaling.ElementsAs(ctx, &var_protocols_bgp_group_family_evpn_signaling, false)...)
            if resp.Diagnostics.HasError() {
                return
            }
	    config.Groups.Protocols[i_protocols].Bgp[i_protocols_bgp].Group[i_protocols_bgp_group].Family[i_protocols_bgp_group_family].Evpn[i_protocols_bgp_group_family_evpn].Signaling = make([]xml_Protocols_Bgp_Group_Family_Evpn_Signaling, len(var_protocols_bgp_group_family_evpn_signaling))
        
		for i_protocols_bgp_group_family_evpn_signaling, v_protocols_bgp_group_family_evpn_signaling := range var_protocols_bgp_group_family_evpn_signaling {
            var var_protocols_bgp_group_family_evpn_signaling_delay_route_advertisements []Protocols_Bgp_Group_Family_Evpn_Signaling_Delay_route_advertisements_Model
            resp.Diagnostics.Append(v_protocols_bgp_group_family_evpn_signaling.Delay_route_advertisements.ElementsAs(ctx, &var_protocols_bgp_group_family_evpn_signaling_delay_route_advertisements, false)...)
            if resp.Diagnostics.HasError() {
                return
            }
	    config.Groups.Protocols[i_protocols].Bgp[i_protocols_bgp].Group[i_protocols_bgp_group].Family[i_protocols_bgp_group_family].Evpn[i_protocols_bgp_group_family_evpn].Signaling[i_protocols_bgp_group_family_evpn_signaling].Delay_route_advertisements = make([]xml_Protocols_Bgp_Group_Family_Evpn_Signaling_Delay_route_advertisements, len(var_protocols_bgp_group_family_evpn_signaling_delay_route_advertisements))
        
		for i_protocols_bgp_group_family_evpn_signaling_delay_route_advertisements, v_protocols_bgp_group_family_evpn_signaling_delay_route_advertisements := range var_protocols_bgp_group_family_evpn_signaling_delay_route_advertisements {
            var var_protocols_bgp_group_family_evpn_signaling_delay_route_advertisements_minimum_delay []Protocols_Bgp_Group_Family_Evpn_Signaling_Delay_route_advertisements_Minimum_delay_Model
            resp.Diagnostics.Append(v_protocols_bgp_group_family_evpn_signaling_delay_route_advertisements.Minimum_delay.ElementsAs(ctx, &var_protocols_bgp_group_family_evpn_signaling_delay_route_advertisements_minimum_delay, false)...)
            if resp.Diagnostics.HasError() {
                return
            }
	    config.Groups.Protocols[i_protocols].Bgp[i_protocols_bgp].Group[i_protocols_bgp_group].Family[i_protocols_bgp_group_family].Evpn[i_protocols_bgp_group_family_evpn].Signaling[i_protocols_bgp_group_family_evpn_signaling].Delay_route_advertisements[i_protocols_bgp_group_family_evpn_signaling_delay_route_advertisements].Minimum_delay = make([]xml_Protocols_Bgp_Group_Family_Evpn_Signaling_Delay_route_advertisements_Minimum_delay, len(var_protocols_bgp_group_family_evpn_signaling_delay_route_advertisements_minimum_delay))
        
		for i_protocols_bgp_group_family_evpn_signaling_delay_route_advertisements_minimum_delay, v_protocols_bgp_group_family_evpn_signaling_delay_route_advertisements_minimum_delay := range var_protocols_bgp_group_family_evpn_signaling_delay_route_advertisements_minimum_delay {
            config.Groups.Protocols[i_protocols].Bgp[i_protocols_bgp].Group[i_protocols_bgp_group].Family[i_protocols_bgp_group_family].Evpn[i_protocols_bgp_group_family_evpn].Signaling[i_protocols_bgp_group_family_evpn_signaling].Delay_route_advertisements[i_protocols_bgp_group_family_evpn_signaling_delay_route_advertisements].Minimum_delay[i_protocols_bgp_group_family_evpn_signaling_delay_route_advertisements_minimum_delay].Routing_uptime = v_protocols_bgp_group_family_evpn_signaling_delay_route_advertisements_minimum_delay.Routing_uptime.ValueStringPointer()
        }
        }
        }
        }
        }
			var var_protocols_bgp_group_export []string
			resp.Diagnostics.Append(v_protocols_bgp_group.Export.ElementsAs(ctx, &var_protocols_bgp_group_export, false)...)
			for _, v_protocols_bgp_group_export := range var_protocols_bgp_group_export {
				config.Groups.Protocols[i_protocols].Bgp[i_protocols_bgp].Group[i_protocols_bgp_group].Export = append(config.Groups.Protocols[i_protocols].Bgp[i_protocols_bgp].Group[i_protocols_bgp_group].Export, &v_protocols_bgp_group_export)
			}
            config.Groups.Protocols[i_protocols].Bgp[i_protocols_bgp].Group[i_protocols_bgp_group].Vpn_apply_export = v_protocols_bgp_group.Vpn_apply_export.ValueStringPointer()
            config.Groups.Protocols[i_protocols].Bgp[i_protocols_bgp].Group[i_protocols_bgp_group].Cluster = v_protocols_bgp_group.Cluster.ValueStringPointer()
            var var_protocols_bgp_group_local_as []Protocols_Bgp_Group_Local_as_Model
            resp.Diagnostics.Append(v_protocols_bgp_group.Local_as.ElementsAs(ctx, &var_protocols_bgp_group_local_as, false)...)
            if resp.Diagnostics.HasError() {
                return
            }
	    config.Groups.Protocols[i_protocols].Bgp[i_protocols_bgp].Group[i_protocols_bgp_group].Local_as = make([]xml_Protocols_Bgp_Group_Local_as, len(var_protocols_bgp_group_local_as))
        
		for i_protocols_bgp_group_local_as, v_protocols_bgp_group_local_as := range var_protocols_bgp_group_local_as {
            config.Groups.Protocols[i_protocols].Bgp[i_protocols_bgp].Group[i_protocols_bgp_group].Local_as[i_protocols_bgp_group_local_as].As_number = v_protocols_bgp_group_local_as.As_number.ValueStringPointer()
        }
            var var_protocols_bgp_group_multipath []Protocols_Bgp_Group_Multipath_Model
            resp.Diagnostics.Append(v_protocols_bgp_group.Multipath.ElementsAs(ctx, &var_protocols_bgp_group_multipath, false)...)
            if resp.Diagnostics.HasError() {
                return
            }
	    config.Groups.Protocols[i_protocols].Bgp[i_protocols_bgp].Group[i_protocols_bgp_group].Multipath = make([]xml_Protocols_Bgp_Group_Multipath, len(var_protocols_bgp_group_multipath))
        
		for i_protocols_bgp_group_multipath, v_protocols_bgp_group_multipath := range var_protocols_bgp_group_multipath {
            config.Groups.Protocols[i_protocols].Bgp[i_protocols_bgp].Group[i_protocols_bgp_group].Multipath[i_protocols_bgp_group_multipath].Multiple_as = v_protocols_bgp_group_multipath.Multiple_as.ValueStringPointer()
        }
            var var_protocols_bgp_group_bfd_liveness_detection []Protocols_Bgp_Group_Bfd_liveness_detection_Model
            resp.Diagnostics.Append(v_protocols_bgp_group.Bfd_liveness_detection.ElementsAs(ctx, &var_protocols_bgp_group_bfd_liveness_detection, false)...)
            if resp.Diagnostics.HasError() {
                return
            }
	    config.Groups.Protocols[i_protocols].Bgp[i_protocols_bgp].Group[i_protocols_bgp_group].Bfd_liveness_detection = make([]xml_Protocols_Bgp_Group_Bfd_liveness_detection, len(var_protocols_bgp_group_bfd_liveness_detection))
        
		for i_protocols_bgp_group_bfd_liveness_detection, v_protocols_bgp_group_bfd_liveness_detection := range var_protocols_bgp_group_bfd_liveness_detection {
            config.Groups.Protocols[i_protocols].Bgp[i_protocols_bgp].Group[i_protocols_bgp_group].Bfd_liveness_detection[i_protocols_bgp_group_bfd_liveness_detection].Minimum_interval = v_protocols_bgp_group_bfd_liveness_detection.Minimum_interval.ValueStringPointer()
            config.Groups.Protocols[i_protocols].Bgp[i_protocols_bgp].Group[i_protocols_bgp_group].Bfd_liveness_detection[i_protocols_bgp_group_bfd_liveness_detection].Multiplier = v_protocols_bgp_group_bfd_liveness_detection.Multiplier.ValueStringPointer()
        }
			var var_protocols_bgp_group_allow []string
			resp.Diagnostics.Append(v_protocols_bgp_group.Allow.ElementsAs(ctx, &var_protocols_bgp_group_allow, false)...)
			for _, v_protocols_bgp_group_allow := range var_protocols_bgp_group_allow {
				config.Groups.Protocols[i_protocols].Bgp[i_protocols_bgp].Group[i_protocols_bgp_group].Allow = append(config.Groups.Protocols[i_protocols].Bgp[i_protocols_bgp].Group[i_protocols_bgp_group].Allow, &v_protocols_bgp_group_allow)
			}
            var var_protocols_bgp_group_neighbor []Protocols_Bgp_Group_Neighbor_Model
            resp.Diagnostics.Append(v_protocols_bgp_group.Neighbor.ElementsAs(ctx, &var_protocols_bgp_group_neighbor, false)...)
            if resp.Diagnostics.HasError() {
                return
            }
	    config.Groups.Protocols[i_protocols].Bgp[i_protocols_bgp].Group[i_protocols_bgp_group].Neighbor = make([]xml_Protocols_Bgp_Group_Neighbor, len(var_protocols_bgp_group_neighbor))
        
		for i_protocols_bgp_group_neighbor, v_protocols_bgp_group_neighbor := range var_protocols_bgp_group_neighbor {
            config.Groups.Protocols[i_protocols].Bgp[i_protocols_bgp].Group[i_protocols_bgp_group].Neighbor[i_protocols_bgp_group_neighbor].Name = v_protocols_bgp_group_neighbor.Name.ValueStringPointer()
            config.Groups.Protocols[i_protocols].Bgp[i_protocols_bgp].Group[i_protocols_bgp_group].Neighbor[i_protocols_bgp_group_neighbor].Description = v_protocols_bgp_group_neighbor.Description.ValueStringPointer()
            config.Groups.Protocols[i_protocols].Bgp[i_protocols_bgp].Group[i_protocols_bgp_group].Neighbor[i_protocols_bgp_group_neighbor].Peer_as = v_protocols_bgp_group_neighbor.Peer_as.ValueStringPointer()
        }
        }
        }
        var var_protocols_evpn []Protocols_Evpn_Model
        resp.Diagnostics.Append(v_protocols.Evpn.ElementsAs(ctx, &var_protocols_evpn, false)...)
        if resp.Diagnostics.HasError() {
            return
        }
	    config.Groups.Protocols[i_protocols].Evpn = make([]xml_Protocols_Evpn, len(var_protocols_evpn))
        
		for i_protocols_evpn, v_protocols_evpn := range var_protocols_evpn {
            config.Groups.Protocols[i_protocols].Evpn[i_protocols_evpn].Encapsulation = v_protocols_evpn.Encapsulation.ValueStringPointer()
            config.Groups.Protocols[i_protocols].Evpn[i_protocols_evpn].Multicast_mode = v_protocols_evpn.Multicast_mode.ValueStringPointer()
            config.Groups.Protocols[i_protocols].Evpn[i_protocols_evpn].Default_gateway = v_protocols_evpn.Default_gateway.ValueStringPointer()
			var var_protocols_evpn_extended_vni_list []string
			resp.Diagnostics.Append(v_protocols_evpn.Extended_vni_list.ElementsAs(ctx, &var_protocols_evpn_extended_vni_list, false)...)
			for _, v_protocols_evpn_extended_vni_list := range var_protocols_evpn_extended_vni_list {
				config.Groups.Protocols[i_protocols].Evpn[i_protocols_evpn].Extended_vni_list = append(config.Groups.Protocols[i_protocols].Evpn[i_protocols_evpn].Extended_vni_list, &v_protocols_evpn_extended_vni_list)
			}
            config.Groups.Protocols[i_protocols].Evpn[i_protocols_evpn].No_core_isolation = v_protocols_evpn.No_core_isolation.ValueStringPointer()
        }
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
        var var_protocols_igmp_snooping []Protocols_Igmp_snooping_Model
        resp.Diagnostics.Append(v_protocols.Igmp_snooping.ElementsAs(ctx, &var_protocols_igmp_snooping, false)...)
        if resp.Diagnostics.HasError() {
            return
        }
	    config.Groups.Protocols[i_protocols].Igmp_snooping = make([]xml_Protocols_Igmp_snooping, len(var_protocols_igmp_snooping))
        
		for i_protocols_igmp_snooping, v_protocols_igmp_snooping := range var_protocols_igmp_snooping {
            var var_protocols_igmp_snooping_vlan []Protocols_Igmp_snooping_Vlan_Model
            resp.Diagnostics.Append(v_protocols_igmp_snooping.Vlan.ElementsAs(ctx, &var_protocols_igmp_snooping_vlan, false)...)
            if resp.Diagnostics.HasError() {
                return
            }
	    config.Groups.Protocols[i_protocols].Igmp_snooping[i_protocols_igmp_snooping].Vlan = make([]xml_Protocols_Igmp_snooping_Vlan, len(var_protocols_igmp_snooping_vlan))
        
		for i_protocols_igmp_snooping_vlan, v_protocols_igmp_snooping_vlan := range var_protocols_igmp_snooping_vlan {
            config.Groups.Protocols[i_protocols].Igmp_snooping[i_protocols_igmp_snooping].Vlan[i_protocols_igmp_snooping_vlan].Name = v_protocols_igmp_snooping_vlan.Name.ValueStringPointer()
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
            var var_routing_instances_instance_route_distinguisher []Routing_instances_Instance_Route_distinguisher_Model
            resp.Diagnostics.Append(v_routing_instances_instance.Route_distinguisher.ElementsAs(ctx, &var_routing_instances_instance_route_distinguisher, false)...)
            if resp.Diagnostics.HasError() {
                return
            }
	    config.Groups.Routing_instances[i_routing_instances].Instance[i_routing_instances_instance].Route_distinguisher = make([]xml_Routing_instances_Instance_Route_distinguisher, len(var_routing_instances_instance_route_distinguisher))
        
		for i_routing_instances_instance_route_distinguisher, v_routing_instances_instance_route_distinguisher := range var_routing_instances_instance_route_distinguisher {
            config.Groups.Routing_instances[i_routing_instances].Instance[i_routing_instances_instance].Route_distinguisher[i_routing_instances_instance_route_distinguisher].Rd_type = v_routing_instances_instance_route_distinguisher.Rd_type.ValueStringPointer()
        }
            var var_routing_instances_instance_vrf_target []Routing_instances_Instance_Vrf_target_Model
            resp.Diagnostics.Append(v_routing_instances_instance.Vrf_target.ElementsAs(ctx, &var_routing_instances_instance_vrf_target, false)...)
            if resp.Diagnostics.HasError() {
                return
            }
	    config.Groups.Routing_instances[i_routing_instances].Instance[i_routing_instances_instance].Vrf_target = make([]xml_Routing_instances_Instance_Vrf_target, len(var_routing_instances_instance_vrf_target))
        
		for i_routing_instances_instance_vrf_target, v_routing_instances_instance_vrf_target := range var_routing_instances_instance_vrf_target {
            config.Groups.Routing_instances[i_routing_instances].Instance[i_routing_instances_instance].Vrf_target[i_routing_instances_instance_vrf_target].Community = v_routing_instances_instance_vrf_target.Community.ValueStringPointer()
        }
            var var_routing_instances_instance_vrf_table_label []Routing_instances_Instance_Vrf_table_label_Model
            resp.Diagnostics.Append(v_routing_instances_instance.Vrf_table_label.ElementsAs(ctx, &var_routing_instances_instance_vrf_table_label, false)...)
            if resp.Diagnostics.HasError() {
                return
            }
	    config.Groups.Routing_instances[i_routing_instances].Instance[i_routing_instances_instance].Vrf_table_label = make([]xml_Routing_instances_Instance_Vrf_table_label, len(var_routing_instances_instance_vrf_table_label))
        
            var var_routing_instances_instance_routing_options []Routing_instances_Instance_Routing_options_Model
            resp.Diagnostics.Append(v_routing_instances_instance.Routing_options.ElementsAs(ctx, &var_routing_instances_instance_routing_options, false)...)
            if resp.Diagnostics.HasError() {
                return
            }
	    config.Groups.Routing_instances[i_routing_instances].Instance[i_routing_instances_instance].Routing_options = make([]xml_Routing_instances_Instance_Routing_options, len(var_routing_instances_instance_routing_options))
        
		for i_routing_instances_instance_routing_options, v_routing_instances_instance_routing_options := range var_routing_instances_instance_routing_options {
            var var_routing_instances_instance_routing_options_auto_export []Routing_instances_Instance_Routing_options_Auto_export_Model
            resp.Diagnostics.Append(v_routing_instances_instance_routing_options.Auto_export.ElementsAs(ctx, &var_routing_instances_instance_routing_options_auto_export, false)...)
            if resp.Diagnostics.HasError() {
                return
            }
	    config.Groups.Routing_instances[i_routing_instances].Instance[i_routing_instances_instance].Routing_options[i_routing_instances_instance_routing_options].Auto_export = make([]xml_Routing_instances_Instance_Routing_options_Auto_export, len(var_routing_instances_instance_routing_options_auto_export))
        
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
            config.Groups.Routing_instances[i_routing_instances].Instance[i_routing_instances_instance].Protocols[i_routing_instances_instance_protocols].Ospf[i_routing_instances_instance_protocols_ospf].Area[i_routing_instances_instance_protocols_ospf_area].Interface[i_routing_instances_instance_protocols_ospf_area_interface].Metric = v_routing_instances_instance_protocols_ospf_area_interface.Metric.ValueStringPointer()
        }
        }
        }
            var var_routing_instances_instance_protocols_evpn []Routing_instances_Instance_Protocols_Evpn_Model
            resp.Diagnostics.Append(v_routing_instances_instance_protocols.Evpn.ElementsAs(ctx, &var_routing_instances_instance_protocols_evpn, false)...)
            if resp.Diagnostics.HasError() {
                return
            }
	    config.Groups.Routing_instances[i_routing_instances].Instance[i_routing_instances_instance].Protocols[i_routing_instances_instance_protocols].Evpn = make([]xml_Routing_instances_Instance_Protocols_Evpn, len(var_routing_instances_instance_protocols_evpn))
        
		for i_routing_instances_instance_protocols_evpn, v_routing_instances_instance_protocols_evpn := range var_routing_instances_instance_protocols_evpn {
            var var_routing_instances_instance_protocols_evpn_ip_prefix_routes []Routing_instances_Instance_Protocols_Evpn_Ip_prefix_routes_Model
            resp.Diagnostics.Append(v_routing_instances_instance_protocols_evpn.Ip_prefix_routes.ElementsAs(ctx, &var_routing_instances_instance_protocols_evpn_ip_prefix_routes, false)...)
            if resp.Diagnostics.HasError() {
                return
            }
	    config.Groups.Routing_instances[i_routing_instances].Instance[i_routing_instances_instance].Protocols[i_routing_instances_instance_protocols].Evpn[i_routing_instances_instance_protocols_evpn].Ip_prefix_routes = make([]xml_Routing_instances_Instance_Protocols_Evpn_Ip_prefix_routes, len(var_routing_instances_instance_protocols_evpn_ip_prefix_routes))
        
		for i_routing_instances_instance_protocols_evpn_ip_prefix_routes, v_routing_instances_instance_protocols_evpn_ip_prefix_routes := range var_routing_instances_instance_protocols_evpn_ip_prefix_routes {
            config.Groups.Routing_instances[i_routing_instances].Instance[i_routing_instances_instance].Protocols[i_routing_instances_instance_protocols].Evpn[i_routing_instances_instance_protocols_evpn].Ip_prefix_routes[i_routing_instances_instance_protocols_evpn_ip_prefix_routes].Advertise = v_routing_instances_instance_protocols_evpn_ip_prefix_routes.Advertise.ValueStringPointer()
            config.Groups.Routing_instances[i_routing_instances].Instance[i_routing_instances_instance].Protocols[i_routing_instances_instance_protocols].Evpn[i_routing_instances_instance_protocols_evpn].Ip_prefix_routes[i_routing_instances_instance_protocols_evpn_ip_prefix_routes].Encapsulation = v_routing_instances_instance_protocols_evpn_ip_prefix_routes.Encapsulation.ValueStringPointer()
            config.Groups.Routing_instances[i_routing_instances].Instance[i_routing_instances_instance].Protocols[i_routing_instances_instance_protocols].Evpn[i_routing_instances_instance_protocols_evpn].Ip_prefix_routes[i_routing_instances_instance_protocols_evpn_ip_prefix_routes].Vni = v_routing_instances_instance_protocols_evpn_ip_prefix_routes.Vni.ValueStringPointer()
			var var_routing_instances_instance_protocols_evpn_ip_prefix_routes_export []string
			resp.Diagnostics.Append(v_routing_instances_instance_protocols_evpn_ip_prefix_routes.Export.ElementsAs(ctx, &var_routing_instances_instance_protocols_evpn_ip_prefix_routes_export, false)...)
			for _, v_routing_instances_instance_protocols_evpn_ip_prefix_routes_export := range var_routing_instances_instance_protocols_evpn_ip_prefix_routes_export {
				config.Groups.Routing_instances[i_routing_instances].Instance[i_routing_instances_instance].Protocols[i_routing_instances_instance_protocols].Evpn[i_routing_instances_instance_protocols_evpn].Ip_prefix_routes[i_routing_instances_instance_protocols_evpn_ip_prefix_routes].Export = append(config.Groups.Routing_instances[i_routing_instances].Instance[i_routing_instances_instance].Protocols[i_routing_instances_instance_protocols].Evpn[i_routing_instances_instance_protocols_evpn].Ip_prefix_routes[i_routing_instances_instance_protocols_evpn_ip_prefix_routes].Export, &v_routing_instances_instance_protocols_evpn_ip_prefix_routes_export)
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
        config.Groups.Routing_options[i_routing_options].Router_id = v_routing_options.Router_id.ValueStringPointer()
        var var_routing_options_forwarding_table []Routing_options_Forwarding_table_Model
        resp.Diagnostics.Append(v_routing_options.Forwarding_table.ElementsAs(ctx, &var_routing_options_forwarding_table, false)...)
        if resp.Diagnostics.HasError() {
            return
        }
	    config.Groups.Routing_options[i_routing_options].Forwarding_table = make([]xml_Routing_options_Forwarding_table, len(var_routing_options_forwarding_table))
        
		for i_routing_options_forwarding_table, v_routing_options_forwarding_table := range var_routing_options_forwarding_table {
			var var_routing_options_forwarding_table_export []string
			resp.Diagnostics.Append(v_routing_options_forwarding_table.Export.ElementsAs(ctx, &var_routing_options_forwarding_table_export, false)...)
			for _, v_routing_options_forwarding_table_export := range var_routing_options_forwarding_table_export {
				config.Groups.Routing_options[i_routing_options].Forwarding_table[i_routing_options_forwarding_table].Export = append(config.Groups.Routing_options[i_routing_options].Forwarding_table[i_routing_options_forwarding_table].Export, &v_routing_options_forwarding_table_export)
			}
            config.Groups.Routing_options[i_routing_options].Forwarding_table[i_routing_options_forwarding_table].Ecmp_fast_reroute = v_routing_options_forwarding_table.Ecmp_fast_reroute.ValueStringPointer()
            var var_routing_options_forwarding_table_chained_composite_next_hop []Routing_options_Forwarding_table_Chained_composite_next_hop_Model
            resp.Diagnostics.Append(v_routing_options_forwarding_table.Chained_composite_next_hop.ElementsAs(ctx, &var_routing_options_forwarding_table_chained_composite_next_hop, false)...)
            if resp.Diagnostics.HasError() {
                return
            }
	    config.Groups.Routing_options[i_routing_options].Forwarding_table[i_routing_options_forwarding_table].Chained_composite_next_hop = make([]xml_Routing_options_Forwarding_table_Chained_composite_next_hop, len(var_routing_options_forwarding_table_chained_composite_next_hop))
        
		for i_routing_options_forwarding_table_chained_composite_next_hop, v_routing_options_forwarding_table_chained_composite_next_hop := range var_routing_options_forwarding_table_chained_composite_next_hop {
            var var_routing_options_forwarding_table_chained_composite_next_hop_ingress []Routing_options_Forwarding_table_Chained_composite_next_hop_Ingress_Model
            resp.Diagnostics.Append(v_routing_options_forwarding_table_chained_composite_next_hop.Ingress.ElementsAs(ctx, &var_routing_options_forwarding_table_chained_composite_next_hop_ingress, false)...)
            if resp.Diagnostics.HasError() {
                return
            }
	    config.Groups.Routing_options[i_routing_options].Forwarding_table[i_routing_options_forwarding_table].Chained_composite_next_hop[i_routing_options_forwarding_table_chained_composite_next_hop].Ingress = make([]xml_Routing_options_Forwarding_table_Chained_composite_next_hop_Ingress, len(var_routing_options_forwarding_table_chained_composite_next_hop_ingress))
        
		for i_routing_options_forwarding_table_chained_composite_next_hop_ingress, v_routing_options_forwarding_table_chained_composite_next_hop_ingress := range var_routing_options_forwarding_table_chained_composite_next_hop_ingress {
            config.Groups.Routing_options[i_routing_options].Forwarding_table[i_routing_options_forwarding_table].Chained_composite_next_hop[i_routing_options_forwarding_table_chained_composite_next_hop].Ingress[i_routing_options_forwarding_table_chained_composite_next_hop_ingress].Evpn = v_routing_options_forwarding_table_chained_composite_next_hop_ingress.Evpn.ValueStringPointer()
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
	
    var var_switch_options []Switch_options_Model
    if plan.Switch_options.IsNull() {
        var_switch_options = []Switch_options_Model{}
    }else {
        resp.Diagnostics.Append(plan.Switch_options.ElementsAs(ctx, &var_switch_options, false)...)
        if resp.Diagnostics.HasError() {
            return
        }
    }
    config.Groups.Switch_options = make([]xml_Switch_options, len(var_switch_options))
   
    for i_switch_options, v_switch_options := range var_switch_options {
        var var_switch_options_vtep_source_interface []Switch_options_Vtep_source_interface_Model
        resp.Diagnostics.Append(v_switch_options.Vtep_source_interface.ElementsAs(ctx, &var_switch_options_vtep_source_interface, false)...)
        if resp.Diagnostics.HasError() {
            return
        }
	    config.Groups.Switch_options[i_switch_options].Vtep_source_interface = make([]xml_Switch_options_Vtep_source_interface, len(var_switch_options_vtep_source_interface))
        
		for i_switch_options_vtep_source_interface, v_switch_options_vtep_source_interface := range var_switch_options_vtep_source_interface {
            config.Groups.Switch_options[i_switch_options].Vtep_source_interface[i_switch_options_vtep_source_interface].Interface_name = v_switch_options_vtep_source_interface.Interface_name.ValueStringPointer()
        }
        var var_switch_options_route_distinguisher []Switch_options_Route_distinguisher_Model
        resp.Diagnostics.Append(v_switch_options.Route_distinguisher.ElementsAs(ctx, &var_switch_options_route_distinguisher, false)...)
        if resp.Diagnostics.HasError() {
            return
        }
	    config.Groups.Switch_options[i_switch_options].Route_distinguisher = make([]xml_Switch_options_Route_distinguisher, len(var_switch_options_route_distinguisher))
        
		for i_switch_options_route_distinguisher, v_switch_options_route_distinguisher := range var_switch_options_route_distinguisher {
            config.Groups.Switch_options[i_switch_options].Route_distinguisher[i_switch_options_route_distinguisher].Rd_type = v_switch_options_route_distinguisher.Rd_type.ValueStringPointer()
        }
        var var_switch_options_vrf_target []Switch_options_Vrf_target_Model
        resp.Diagnostics.Append(v_switch_options.Vrf_target.ElementsAs(ctx, &var_switch_options_vrf_target, false)...)
        if resp.Diagnostics.HasError() {
            return
        }
	    config.Groups.Switch_options[i_switch_options].Vrf_target = make([]xml_Switch_options_Vrf_target, len(var_switch_options_vrf_target))
        
		for i_switch_options_vrf_target, v_switch_options_vrf_target := range var_switch_options_vrf_target {
            config.Groups.Switch_options[i_switch_options].Vrf_target[i_switch_options_vrf_target].Community = v_switch_options_vrf_target.Community.ValueStringPointer()
            var var_switch_options_vrf_target_auto []Switch_options_Vrf_target_Auto_Model
            resp.Diagnostics.Append(v_switch_options_vrf_target.Auto.ElementsAs(ctx, &var_switch_options_vrf_target_auto, false)...)
            if resp.Diagnostics.HasError() {
                return
            }
	    config.Groups.Switch_options[i_switch_options].Vrf_target[i_switch_options_vrf_target].Auto = make([]xml_Switch_options_Vrf_target_Auto, len(var_switch_options_vrf_target_auto))
        
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
            config.Groups.System[i_system].Syslog[i_system_syslog].File[i_system_syslog_file].Contents[i_system_syslog_file_contents].Notice = v_system_syslog_file_contents.Notice.ValueStringPointer()
            config.Groups.System[i_system].Syslog[i_system_syslog].File[i_system_syslog_file].Contents[i_system_syslog_file_contents].Info = v_system_syslog_file_contents.Info.ValueStringPointer()
        }
        }
        }
        var var_system_extensions []System_Extensions_Model
        resp.Diagnostics.Append(v_system.Extensions.ElementsAs(ctx, &var_system_extensions, false)...)
        if resp.Diagnostics.HasError() {
            return
        }
	    config.Groups.System[i_system].Extensions = make([]xml_System_Extensions, len(var_system_extensions))
        
		for i_system_extensions, v_system_extensions := range var_system_extensions {
            var var_system_extensions_providers []System_Extensions_Providers_Model
            resp.Diagnostics.Append(v_system_extensions.Providers.ElementsAs(ctx, &var_system_extensions_providers, false)...)
            if resp.Diagnostics.HasError() {
                return
            }
	    config.Groups.System[i_system].Extensions[i_system_extensions].Providers = make([]xml_System_Extensions_Providers, len(var_system_extensions_providers))
        
		for i_system_extensions_providers, v_system_extensions_providers := range var_system_extensions_providers {
            config.Groups.System[i_system].Extensions[i_system_extensions].Providers[i_system_extensions_providers].Name = v_system_extensions_providers.Name.ValueStringPointer()
            var var_system_extensions_providers_license_type []System_Extensions_Providers_License_type_Model
            resp.Diagnostics.Append(v_system_extensions_providers.License_type.ElementsAs(ctx, &var_system_extensions_providers_license_type, false)...)
            if resp.Diagnostics.HasError() {
                return
            }
	    config.Groups.System[i_system].Extensions[i_system_extensions].Providers[i_system_extensions_providers].License_type = make([]xml_System_Extensions_Providers_License_type, len(var_system_extensions_providers_license_type))
        
		for i_system_extensions_providers_license_type, v_system_extensions_providers_license_type := range var_system_extensions_providers_license_type {
            config.Groups.System[i_system].Extensions[i_system_extensions].Providers[i_system_extensions_providers].License_type[i_system_extensions_providers_license_type].Name = v_system_extensions_providers_license_type.Name.ValueStringPointer()
			var var_system_extensions_providers_license_type_deployment_scope []string
			resp.Diagnostics.Append(v_system_extensions_providers_license_type.Deployment_scope.ElementsAs(ctx, &var_system_extensions_providers_license_type_deployment_scope, false)...)
			for _, v_system_extensions_providers_license_type_deployment_scope := range var_system_extensions_providers_license_type_deployment_scope {
				config.Groups.System[i_system].Extensions[i_system_extensions].Providers[i_system_extensions_providers].License_type[i_system_extensions_providers_license_type].Deployment_scope = append(config.Groups.System[i_system].Extensions[i_system_extensions].Providers[i_system_extensions_providers].License_type[i_system_extensions_providers_license_type].Deployment_scope, &v_system_extensions_providers_license_type_deployment_scope)
			}
        }
        }
        }
    }
	
    var var_vlans []Vlans_Model
    if plan.Vlans.IsNull() {
        var_vlans = []Vlans_Model{}
    }else {
        resp.Diagnostics.Append(plan.Vlans.ElementsAs(ctx, &var_vlans, false)...)
        if resp.Diagnostics.HasError() {
            return
        }
    }
    config.Groups.Vlans = make([]xml_Vlans, len(var_vlans))
   
    for i_vlans, v_vlans := range var_vlans {
        var var_vlans_vlan []Vlans_Vlan_Model
        resp.Diagnostics.Append(v_vlans.Vlan.ElementsAs(ctx, &var_vlans_vlan, false)...)
        if resp.Diagnostics.HasError() {
            return
        }
	    config.Groups.Vlans[i_vlans].Vlan = make([]xml_Vlans_Vlan, len(var_vlans_vlan))
        
		for i_vlans_vlan, v_vlans_vlan := range var_vlans_vlan {
            config.Groups.Vlans[i_vlans].Vlan[i_vlans_vlan].Name = v_vlans_vlan.Name.ValueStringPointer()
            config.Groups.Vlans[i_vlans].Vlan[i_vlans_vlan].Vlan_id = v_vlans_vlan.Vlan_id.ValueStringPointer()
            config.Groups.Vlans[i_vlans].Vlan[i_vlans_vlan].L3_interface = v_vlans_vlan.L3_interface.ValueStringPointer()
            var var_vlans_vlan_vxlan []Vlans_Vlan_Vxlan_Model
            resp.Diagnostics.Append(v_vlans_vlan.Vxlan.ElementsAs(ctx, &var_vlans_vlan_vxlan, false)...)
            if resp.Diagnostics.HasError() {
                return
            }
	    config.Groups.Vlans[i_vlans].Vlan[i_vlans_vlan].Vxlan = make([]xml_Vlans_Vlan_Vxlan, len(var_vlans_vlan_vxlan))
        
		for i_vlans_vlan_vxlan, v_vlans_vlan_vxlan := range var_vlans_vlan_vxlan {
            config.Groups.Vlans[i_vlans].Vlan[i_vlans_vlan].Vxlan[i_vlans_vlan_vxlan].Vni = v_vlans_vlan_vxlan.Vni.ValueStringPointer()
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
    state.Chassis = types.ListNull(types.ObjectType{AttrTypes: Groups_Model{}.AttrTypes()})
    chassis_List := make([]Chassis_Model, len(config.Groups.Chassis))
    for i_chassis, v_chassis := range config.Groups.Chassis {
        var chassis_model Chassis_Model
        chassis_aggregated_devices_List := make([]Chassis_Aggregated_devices_Model, len(v_chassis.Aggregated_devices))
        
		for i_chassis_aggregated_devices, v_chassis_aggregated_devices := range v_chassis.Aggregated_devices {
            var chassis_aggregated_devices_model Chassis_Aggregated_devices_Model
			chassis_aggregated_devices_List[i_chassis_aggregated_devices] = chassis_aggregated_devices_model
                
        chassis_aggregated_devices_ethernet_List := make([]Chassis_Aggregated_devices_Ethernet_Model, len(v_chassis_aggregated_devices.Ethernet))
        
		for i_chassis_aggregated_devices_ethernet, v_chassis_aggregated_devices_ethernet := range v_chassis_aggregated_devices.Ethernet {
            var chassis_aggregated_devices_ethernet_model Chassis_Aggregated_devices_Ethernet_Model
            chassis_aggregated_devices_ethernet_model.Device_count = types.StringPointerValue(v_chassis_aggregated_devices_ethernet.Device_count)
			chassis_aggregated_devices_ethernet_List[i_chassis_aggregated_devices_ethernet] = chassis_aggregated_devices_ethernet_model
        }
        chassis_aggregated_devices_model.Ethernet, _ = types.ListValueFrom(ctx, types.ObjectType{AttrTypes: Chassis_Aggregated_devices_Ethernet_Model{}.AttrTypes()}, chassis_aggregated_devices_ethernet_List)
        chassis_aggregated_devices_List[i_chassis_aggregated_devices] = chassis_aggregated_devices_model
        }
        chassis_model.Aggregated_devices, _ = types.ListValueFrom(ctx, types.ObjectType{AttrTypes: Chassis_Aggregated_devices_Model{}.AttrTypes()}, chassis_aggregated_devices_List)
        chassis_List[i_chassis] = chassis_model
    }
    state.Chassis, _ = types.ListValueFrom(ctx, types.ObjectType{AttrTypes: Chassis_Model{}.AttrTypes()}, chassis_List)
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
    state.Forwarding_options = types.ListNull(types.ObjectType{AttrTypes: Groups_Model{}.AttrTypes()})
    forwarding_options_List := make([]Forwarding_options_Model, len(config.Groups.Forwarding_options))
    for i_forwarding_options, v_forwarding_options := range config.Groups.Forwarding_options {
        var forwarding_options_model Forwarding_options_Model
        forwarding_options_storm_control_profiles_List := make([]Forwarding_options_Storm_control_profiles_Model, len(v_forwarding_options.Storm_control_profiles))
        
		for i_forwarding_options_storm_control_profiles, v_forwarding_options_storm_control_profiles := range v_forwarding_options.Storm_control_profiles {
            var forwarding_options_storm_control_profiles_model Forwarding_options_Storm_control_profiles_Model
            forwarding_options_storm_control_profiles_model.Name = types.StringPointerValue(v_forwarding_options_storm_control_profiles.Name)
			forwarding_options_storm_control_profiles_List[i_forwarding_options_storm_control_profiles] = forwarding_options_storm_control_profiles_model
			forwarding_options_storm_control_profiles_List[i_forwarding_options_storm_control_profiles] = forwarding_options_storm_control_profiles_model
                
        forwarding_options_storm_control_profiles_all_List := make([]Forwarding_options_Storm_control_profiles_All_Model, len(v_forwarding_options_storm_control_profiles.All))
        
        forwarding_options_storm_control_profiles_model.All, _ = types.ListValueFrom(ctx, types.ObjectType{AttrTypes: Forwarding_options_Storm_control_profiles_All_Model{}.AttrTypes()}, forwarding_options_storm_control_profiles_all_List)
        forwarding_options_storm_control_profiles_List[i_forwarding_options_storm_control_profiles] = forwarding_options_storm_control_profiles_model
        }
        forwarding_options_model.Storm_control_profiles, _ = types.ListValueFrom(ctx, types.ObjectType{AttrTypes: Forwarding_options_Storm_control_profiles_Model{}.AttrTypes()}, forwarding_options_storm_control_profiles_List)
        forwarding_options_List[i_forwarding_options] = forwarding_options_model
    }
    state.Forwarding_options, _ = types.ListValueFrom(ctx, types.ObjectType{AttrTypes: Forwarding_options_Model{}.AttrTypes()}, forwarding_options_List)
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
    state.Interfaces = types.ListNull(types.ObjectType{AttrTypes: Groups_Model{}.AttrTypes()})
    interfaces_List := make([]Interfaces_Model, len(config.Groups.Interfaces))
    for i_interfaces, v_interfaces := range config.Groups.Interfaces {
        var interfaces_model Interfaces_Model
        interfaces_interface_List := make([]Interfaces_Interface_Model, len(v_interfaces.Interface))
        
		for i_interfaces_interface, v_interfaces_interface := range v_interfaces.Interface {
            var interfaces_interface_model Interfaces_Interface_Model
            interfaces_interface_model.Name = types.StringPointerValue(v_interfaces_interface.Name)
			interfaces_interface_List[i_interfaces_interface] = interfaces_interface_model
            interfaces_interface_model.Description = types.StringPointerValue(v_interfaces_interface.Description)
			interfaces_interface_List[i_interfaces_interface] = interfaces_interface_model
            interfaces_interface_model.Vlan_tagging = types.StringPointerValue(v_interfaces_interface.Vlan_tagging)
			interfaces_interface_List[i_interfaces_interface] = interfaces_interface_model
			interfaces_interface_List[i_interfaces_interface] = interfaces_interface_model
                
        interfaces_interface_esi_List := make([]Interfaces_Interface_Esi_Model, len(v_interfaces_interface.Esi))
        
		for i_interfaces_interface_esi, v_interfaces_interface_esi := range v_interfaces_interface.Esi {
            var interfaces_interface_esi_model Interfaces_Interface_Esi_Model
            interfaces_interface_esi_model.Identifier = types.StringPointerValue(v_interfaces_interface_esi.Identifier)
			interfaces_interface_esi_List[i_interfaces_interface_esi] = interfaces_interface_esi_model
            interfaces_interface_esi_model.All_active = types.StringPointerValue(v_interfaces_interface_esi.All_active)
			interfaces_interface_esi_List[i_interfaces_interface_esi] = interfaces_interface_esi_model
        }
        interfaces_interface_model.Esi, _ = types.ListValueFrom(ctx, types.ObjectType{AttrTypes: Interfaces_Interface_Esi_Model{}.AttrTypes()}, interfaces_interface_esi_List)
        interfaces_interface_List[i_interfaces_interface] = interfaces_interface_model
			interfaces_interface_List[i_interfaces_interface] = interfaces_interface_model
                
        interfaces_interface_ether_options_List := make([]Interfaces_Interface_Ether_options_Model, len(v_interfaces_interface.Ether_options))
        
		for i_interfaces_interface_ether_options, v_interfaces_interface_ether_options := range v_interfaces_interface.Ether_options {
            var interfaces_interface_ether_options_model Interfaces_Interface_Ether_options_Model
			interfaces_interface_ether_options_List[i_interfaces_interface_ether_options] = interfaces_interface_ether_options_model
                
        interfaces_interface_ether_options_ieee_802_3ad_List := make([]Interfaces_Interface_Ether_options_Ieee_802_3ad_Model, len(v_interfaces_interface_ether_options.Ieee_802_3ad))
        
		for i_interfaces_interface_ether_options_ieee_802_3ad, v_interfaces_interface_ether_options_ieee_802_3ad := range v_interfaces_interface_ether_options.Ieee_802_3ad {
            var interfaces_interface_ether_options_ieee_802_3ad_model Interfaces_Interface_Ether_options_Ieee_802_3ad_Model
            interfaces_interface_ether_options_ieee_802_3ad_model.Bundle = types.StringPointerValue(v_interfaces_interface_ether_options_ieee_802_3ad.Bundle)
			interfaces_interface_ether_options_ieee_802_3ad_List[i_interfaces_interface_ether_options_ieee_802_3ad] = interfaces_interface_ether_options_ieee_802_3ad_model
        }
        interfaces_interface_ether_options_model.Ieee_802_3ad, _ = types.ListValueFrom(ctx, types.ObjectType{AttrTypes: Interfaces_Interface_Ether_options_Ieee_802_3ad_Model{}.AttrTypes()}, interfaces_interface_ether_options_ieee_802_3ad_List)
        interfaces_interface_ether_options_List[i_interfaces_interface_ether_options] = interfaces_interface_ether_options_model
        }
        interfaces_interface_model.Ether_options, _ = types.ListValueFrom(ctx, types.ObjectType{AttrTypes: Interfaces_Interface_Ether_options_Model{}.AttrTypes()}, interfaces_interface_ether_options_List)
        interfaces_interface_List[i_interfaces_interface] = interfaces_interface_model
			interfaces_interface_List[i_interfaces_interface] = interfaces_interface_model
                
        interfaces_interface_aggregated_ether_options_List := make([]Interfaces_Interface_Aggregated_ether_options_Model, len(v_interfaces_interface.Aggregated_ether_options))
        
		for i_interfaces_interface_aggregated_ether_options, v_interfaces_interface_aggregated_ether_options := range v_interfaces_interface.Aggregated_ether_options {
            var interfaces_interface_aggregated_ether_options_model Interfaces_Interface_Aggregated_ether_options_Model
			interfaces_interface_aggregated_ether_options_List[i_interfaces_interface_aggregated_ether_options] = interfaces_interface_aggregated_ether_options_model
                
        interfaces_interface_aggregated_ether_options_lacp_List := make([]Interfaces_Interface_Aggregated_ether_options_Lacp_Model, len(v_interfaces_interface_aggregated_ether_options.Lacp))
        
		for i_interfaces_interface_aggregated_ether_options_lacp, v_interfaces_interface_aggregated_ether_options_lacp := range v_interfaces_interface_aggregated_ether_options.Lacp {
            var interfaces_interface_aggregated_ether_options_lacp_model Interfaces_Interface_Aggregated_ether_options_Lacp_Model
            interfaces_interface_aggregated_ether_options_lacp_model.Active = types.StringPointerValue(v_interfaces_interface_aggregated_ether_options_lacp.Active)
			interfaces_interface_aggregated_ether_options_lacp_List[i_interfaces_interface_aggregated_ether_options_lacp] = interfaces_interface_aggregated_ether_options_lacp_model
            interfaces_interface_aggregated_ether_options_lacp_model.Periodic = types.StringPointerValue(v_interfaces_interface_aggregated_ether_options_lacp.Periodic)
			interfaces_interface_aggregated_ether_options_lacp_List[i_interfaces_interface_aggregated_ether_options_lacp] = interfaces_interface_aggregated_ether_options_lacp_model
            interfaces_interface_aggregated_ether_options_lacp_model.System_id = types.StringPointerValue(v_interfaces_interface_aggregated_ether_options_lacp.System_id)
			interfaces_interface_aggregated_ether_options_lacp_List[i_interfaces_interface_aggregated_ether_options_lacp] = interfaces_interface_aggregated_ether_options_lacp_model
        }
        interfaces_interface_aggregated_ether_options_model.Lacp, _ = types.ListValueFrom(ctx, types.ObjectType{AttrTypes: Interfaces_Interface_Aggregated_ether_options_Lacp_Model{}.AttrTypes()}, interfaces_interface_aggregated_ether_options_lacp_List)
        interfaces_interface_aggregated_ether_options_List[i_interfaces_interface_aggregated_ether_options] = interfaces_interface_aggregated_ether_options_model
        }
        interfaces_interface_model.Aggregated_ether_options, _ = types.ListValueFrom(ctx, types.ObjectType{AttrTypes: Interfaces_Interface_Aggregated_ether_options_Model{}.AttrTypes()}, interfaces_interface_aggregated_ether_options_List)
        interfaces_interface_List[i_interfaces_interface] = interfaces_interface_model
			interfaces_interface_List[i_interfaces_interface] = interfaces_interface_model
                
        interfaces_interface_unit_List := make([]Interfaces_Interface_Unit_Model, len(v_interfaces_interface.Unit))
        
		for i_interfaces_interface_unit, v_interfaces_interface_unit := range v_interfaces_interface.Unit {
            var interfaces_interface_unit_model Interfaces_Interface_Unit_Model
            interfaces_interface_unit_model.Name = types.StringPointerValue(v_interfaces_interface_unit.Name)
			interfaces_interface_unit_List[i_interfaces_interface_unit] = interfaces_interface_unit_model
            interfaces_interface_unit_model.Description = types.StringPointerValue(v_interfaces_interface_unit.Description)
			interfaces_interface_unit_List[i_interfaces_interface_unit] = interfaces_interface_unit_model
            interfaces_interface_unit_model.Vlan_id = types.StringPointerValue(v_interfaces_interface_unit.Vlan_id)
			interfaces_interface_unit_List[i_interfaces_interface_unit] = interfaces_interface_unit_model
			interfaces_interface_unit_List[i_interfaces_interface_unit] = interfaces_interface_unit_model
                
        interfaces_interface_unit_family_List := make([]Interfaces_Interface_Unit_Family_Model, len(v_interfaces_interface_unit.Family))
        
		for i_interfaces_interface_unit_family, v_interfaces_interface_unit_family := range v_interfaces_interface_unit.Family {
            var interfaces_interface_unit_family_model Interfaces_Interface_Unit_Family_Model
			interfaces_interface_unit_family_List[i_interfaces_interface_unit_family] = interfaces_interface_unit_family_model
                
        interfaces_interface_unit_family_inet_List := make([]Interfaces_Interface_Unit_Family_Inet_Model, len(v_interfaces_interface_unit_family.Inet))
        
		for i_interfaces_interface_unit_family_inet, v_interfaces_interface_unit_family_inet := range v_interfaces_interface_unit_family.Inet {
            var interfaces_interface_unit_family_inet_model Interfaces_Interface_Unit_Family_Inet_Model
			interfaces_interface_unit_family_inet_List[i_interfaces_interface_unit_family_inet] = interfaces_interface_unit_family_inet_model
                
        interfaces_interface_unit_family_inet_address_List := make([]Interfaces_Interface_Unit_Family_Inet_Address_Model, len(v_interfaces_interface_unit_family_inet.Address))
        
		for i_interfaces_interface_unit_family_inet_address, v_interfaces_interface_unit_family_inet_address := range v_interfaces_interface_unit_family_inet.Address {
            var interfaces_interface_unit_family_inet_address_model Interfaces_Interface_Unit_Family_Inet_Address_Model
            interfaces_interface_unit_family_inet_address_model.Name = types.StringPointerValue(v_interfaces_interface_unit_family_inet_address.Name)
			interfaces_interface_unit_family_inet_address_List[i_interfaces_interface_unit_family_inet_address] = interfaces_interface_unit_family_inet_address_model
        }
        interfaces_interface_unit_family_inet_model.Address, _ = types.ListValueFrom(ctx, types.ObjectType{AttrTypes: Interfaces_Interface_Unit_Family_Inet_Address_Model{}.AttrTypes()}, interfaces_interface_unit_family_inet_address_List)
        interfaces_interface_unit_family_inet_List[i_interfaces_interface_unit_family_inet] = interfaces_interface_unit_family_inet_model
        }
        interfaces_interface_unit_family_model.Inet, _ = types.ListValueFrom(ctx, types.ObjectType{AttrTypes: Interfaces_Interface_Unit_Family_Inet_Model{}.AttrTypes()}, interfaces_interface_unit_family_inet_List)
        interfaces_interface_unit_family_List[i_interfaces_interface_unit_family] = interfaces_interface_unit_family_model
			interfaces_interface_unit_family_List[i_interfaces_interface_unit_family] = interfaces_interface_unit_family_model
                
        interfaces_interface_unit_family_ethernet_switching_List := make([]Interfaces_Interface_Unit_Family_Ethernet_switching_Model, len(v_interfaces_interface_unit_family.Ethernet_switching))
        
		for i_interfaces_interface_unit_family_ethernet_switching, v_interfaces_interface_unit_family_ethernet_switching := range v_interfaces_interface_unit_family.Ethernet_switching {
            var interfaces_interface_unit_family_ethernet_switching_model Interfaces_Interface_Unit_Family_Ethernet_switching_Model
			interfaces_interface_unit_family_ethernet_switching_List[i_interfaces_interface_unit_family_ethernet_switching] = interfaces_interface_unit_family_ethernet_switching_model
                
        interfaces_interface_unit_family_ethernet_switching_vlan_List := make([]Interfaces_Interface_Unit_Family_Ethernet_switching_Vlan_Model, len(v_interfaces_interface_unit_family_ethernet_switching.Vlan))
        
		for i_interfaces_interface_unit_family_ethernet_switching_vlan, v_interfaces_interface_unit_family_ethernet_switching_vlan := range v_interfaces_interface_unit_family_ethernet_switching.Vlan {
            var interfaces_interface_unit_family_ethernet_switching_vlan_model Interfaces_Interface_Unit_Family_Ethernet_switching_Vlan_Model
			var var_interfaces_interface_unit_family_ethernet_switching_members []*string
			if v_interfaces_interface_unit_family_ethernet_switching_vlan.Members != nil {
				var_interfaces_interface_unit_family_ethernet_switching_members = make([]*string, len(v_interfaces_interface_unit_family_ethernet_switching_vlan.Members))
				copy(var_interfaces_interface_unit_family_ethernet_switching_members, v_interfaces_interface_unit_family_ethernet_switching_vlan.Members)
			}
			interfaces_interface_unit_family_ethernet_switching_vlan_model.Members, _ = types.ListValueFrom(ctx, types.StringType, var_interfaces_interface_unit_family_ethernet_switching_members)
			interfaces_interface_unit_family_ethernet_switching_vlan_List[i_interfaces_interface_unit_family_ethernet_switching_vlan] = interfaces_interface_unit_family_ethernet_switching_vlan_model
        }
        interfaces_interface_unit_family_ethernet_switching_model.Vlan, _ = types.ListValueFrom(ctx, types.ObjectType{AttrTypes: Interfaces_Interface_Unit_Family_Ethernet_switching_Vlan_Model{}.AttrTypes()}, interfaces_interface_unit_family_ethernet_switching_vlan_List)
        interfaces_interface_unit_family_ethernet_switching_List[i_interfaces_interface_unit_family_ethernet_switching] = interfaces_interface_unit_family_ethernet_switching_model
        }
        interfaces_interface_unit_family_model.Ethernet_switching, _ = types.ListValueFrom(ctx, types.ObjectType{AttrTypes: Interfaces_Interface_Unit_Family_Ethernet_switching_Model{}.AttrTypes()}, interfaces_interface_unit_family_ethernet_switching_List)
        interfaces_interface_unit_family_List[i_interfaces_interface_unit_family] = interfaces_interface_unit_family_model
        }
        interfaces_interface_unit_model.Family, _ = types.ListValueFrom(ctx, types.ObjectType{AttrTypes: Interfaces_Interface_Unit_Family_Model{}.AttrTypes()}, interfaces_interface_unit_family_List)
        interfaces_interface_unit_List[i_interfaces_interface_unit] = interfaces_interface_unit_model
            interfaces_interface_unit_model.Mac = types.StringPointerValue(v_interfaces_interface_unit.Mac)
			interfaces_interface_unit_List[i_interfaces_interface_unit] = interfaces_interface_unit_model
        }
        interfaces_interface_model.Unit, _ = types.ListValueFrom(ctx, types.ObjectType{AttrTypes: Interfaces_Interface_Unit_Model{}.AttrTypes()}, interfaces_interface_unit_List)
        interfaces_interface_List[i_interfaces_interface] = interfaces_interface_model
        }
        interfaces_model.Interface, _ = types.ListValueFrom(ctx, types.ObjectType{AttrTypes: Interfaces_Interface_Model{}.AttrTypes()}, interfaces_interface_List)
        interfaces_List[i_interfaces] = interfaces_model
    }
    state.Interfaces, _ = types.ListValueFrom(ctx, types.ObjectType{AttrTypes: Interfaces_Model{}.AttrTypes()}, interfaces_List)
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
    state.Policy_options = types.ListNull(types.ObjectType{AttrTypes: Groups_Model{}.AttrTypes()})
    policy_options_List := make([]Policy_options_Model, len(config.Groups.Policy_options))
    for i_policy_options, v_policy_options := range config.Groups.Policy_options {
        var policy_options_model Policy_options_Model
        policy_options_policy_statement_List := make([]Policy_options_Policy_statement_Model, len(v_policy_options.Policy_statement))
        
		for i_policy_options_policy_statement, v_policy_options_policy_statement := range v_policy_options.Policy_statement {
            var policy_options_policy_statement_model Policy_options_Policy_statement_Model
            policy_options_policy_statement_model.Name = types.StringPointerValue(v_policy_options_policy_statement.Name)
			policy_options_policy_statement_List[i_policy_options_policy_statement] = policy_options_policy_statement_model
			policy_options_policy_statement_List[i_policy_options_policy_statement] = policy_options_policy_statement_model
                
        policy_options_policy_statement_term_List := make([]Policy_options_Policy_statement_Term_Model, len(v_policy_options_policy_statement.Term))
        
		for i_policy_options_policy_statement_term, v_policy_options_policy_statement_term := range v_policy_options_policy_statement.Term {
            var policy_options_policy_statement_term_model Policy_options_Policy_statement_Term_Model
            policy_options_policy_statement_term_model.Name = types.StringPointerValue(v_policy_options_policy_statement_term.Name)
			policy_options_policy_statement_term_List[i_policy_options_policy_statement_term] = policy_options_policy_statement_term_model
			policy_options_policy_statement_term_List[i_policy_options_policy_statement_term] = policy_options_policy_statement_term_model
                
        policy_options_policy_statement_term_from_List := make([]Policy_options_Policy_statement_Term_From_Model, len(v_policy_options_policy_statement_term.From))
        
		for i_policy_options_policy_statement_term_from, v_policy_options_policy_statement_term_from := range v_policy_options_policy_statement_term.From {
            var policy_options_policy_statement_term_from_model Policy_options_Policy_statement_Term_From_Model
			var var_policy_options_policy_statement_term_protocol []*string
			if v_policy_options_policy_statement_term_from.Protocol != nil {
				var_policy_options_policy_statement_term_protocol = make([]*string, len(v_policy_options_policy_statement_term_from.Protocol))
				copy(var_policy_options_policy_statement_term_protocol, v_policy_options_policy_statement_term_from.Protocol)
			}
			policy_options_policy_statement_term_from_model.Protocol, _ = types.ListValueFrom(ctx, types.StringType, var_policy_options_policy_statement_term_protocol)
			policy_options_policy_statement_term_from_List[i_policy_options_policy_statement_term_from] = policy_options_policy_statement_term_from_model
			policy_options_policy_statement_term_from_List[i_policy_options_policy_statement_term_from] = policy_options_policy_statement_term_from_model
                
        policy_options_policy_statement_term_from_route_filter_List := make([]Policy_options_Policy_statement_Term_From_Route_filter_Model, len(v_policy_options_policy_statement_term_from.Route_filter))
        
		for i_policy_options_policy_statement_term_from_route_filter, v_policy_options_policy_statement_term_from_route_filter := range v_policy_options_policy_statement_term_from.Route_filter {
            var policy_options_policy_statement_term_from_route_filter_model Policy_options_Policy_statement_Term_From_Route_filter_Model
            policy_options_policy_statement_term_from_route_filter_model.Address = types.StringPointerValue(v_policy_options_policy_statement_term_from_route_filter.Address)
			policy_options_policy_statement_term_from_route_filter_List[i_policy_options_policy_statement_term_from_route_filter] = policy_options_policy_statement_term_from_route_filter_model
            policy_options_policy_statement_term_from_route_filter_model.Exact = types.StringPointerValue(v_policy_options_policy_statement_term_from_route_filter.Exact)
			policy_options_policy_statement_term_from_route_filter_List[i_policy_options_policy_statement_term_from_route_filter] = policy_options_policy_statement_term_from_route_filter_model
            policy_options_policy_statement_term_from_route_filter_model.Orlonger = types.StringPointerValue(v_policy_options_policy_statement_term_from_route_filter.Orlonger)
			policy_options_policy_statement_term_from_route_filter_List[i_policy_options_policy_statement_term_from_route_filter] = policy_options_policy_statement_term_from_route_filter_model
            policy_options_policy_statement_term_from_route_filter_model.Prefix_length_range = types.StringPointerValue(v_policy_options_policy_statement_term_from_route_filter.Prefix_length_range)
			policy_options_policy_statement_term_from_route_filter_List[i_policy_options_policy_statement_term_from_route_filter] = policy_options_policy_statement_term_from_route_filter_model
        }
        policy_options_policy_statement_term_from_model.Route_filter, _ = types.ListValueFrom(ctx, types.ObjectType{AttrTypes: Policy_options_Policy_statement_Term_From_Route_filter_Model{}.AttrTypes()}, policy_options_policy_statement_term_from_route_filter_List)
        policy_options_policy_statement_term_from_List[i_policy_options_policy_statement_term_from] = policy_options_policy_statement_term_from_model
        }
        policy_options_policy_statement_term_model.From, _ = types.ListValueFrom(ctx, types.ObjectType{AttrTypes: Policy_options_Policy_statement_Term_From_Model{}.AttrTypes()}, policy_options_policy_statement_term_from_List)
        policy_options_policy_statement_term_List[i_policy_options_policy_statement_term] = policy_options_policy_statement_term_model
			policy_options_policy_statement_term_List[i_policy_options_policy_statement_term] = policy_options_policy_statement_term_model
                
        policy_options_policy_statement_term_then_List := make([]Policy_options_Policy_statement_Term_Then_Model, len(v_policy_options_policy_statement_term.Then))
        
		for i_policy_options_policy_statement_term_then, v_policy_options_policy_statement_term_then := range v_policy_options_policy_statement_term.Then {
            var policy_options_policy_statement_term_then_model Policy_options_Policy_statement_Term_Then_Model
			policy_options_policy_statement_term_then_List[i_policy_options_policy_statement_term_then] = policy_options_policy_statement_term_then_model
                
        policy_options_policy_statement_term_then_community_List := make([]Policy_options_Policy_statement_Term_Then_Community_Model, len(v_policy_options_policy_statement_term_then.Community))
        
		for i_policy_options_policy_statement_term_then_community, v_policy_options_policy_statement_term_then_community := range v_policy_options_policy_statement_term_then.Community {
            var policy_options_policy_statement_term_then_community_model Policy_options_Policy_statement_Term_Then_Community_Model
            policy_options_policy_statement_term_then_community_model.Add = types.StringPointerValue(v_policy_options_policy_statement_term_then_community.Add)
			policy_options_policy_statement_term_then_community_List[i_policy_options_policy_statement_term_then_community] = policy_options_policy_statement_term_then_community_model
            policy_options_policy_statement_term_then_community_model.Community_name = types.StringPointerValue(v_policy_options_policy_statement_term_then_community.Community_name)
			policy_options_policy_statement_term_then_community_List[i_policy_options_policy_statement_term_then_community] = policy_options_policy_statement_term_then_community_model
        }
        policy_options_policy_statement_term_then_model.Community, _ = types.ListValueFrom(ctx, types.ObjectType{AttrTypes: Policy_options_Policy_statement_Term_Then_Community_Model{}.AttrTypes()}, policy_options_policy_statement_term_then_community_List)
        policy_options_policy_statement_term_then_List[i_policy_options_policy_statement_term_then] = policy_options_policy_statement_term_then_model
            policy_options_policy_statement_term_then_model.Accept = types.StringPointerValue(v_policy_options_policy_statement_term_then.Accept)
			policy_options_policy_statement_term_then_List[i_policy_options_policy_statement_term_then] = policy_options_policy_statement_term_then_model
            policy_options_policy_statement_term_then_model.Reject = types.StringPointerValue(v_policy_options_policy_statement_term_then.Reject)
			policy_options_policy_statement_term_then_List[i_policy_options_policy_statement_term_then] = policy_options_policy_statement_term_then_model
        }
        policy_options_policy_statement_term_model.Then, _ = types.ListValueFrom(ctx, types.ObjectType{AttrTypes: Policy_options_Policy_statement_Term_Then_Model{}.AttrTypes()}, policy_options_policy_statement_term_then_List)
        policy_options_policy_statement_term_List[i_policy_options_policy_statement_term] = policy_options_policy_statement_term_model
        }
        policy_options_policy_statement_model.Term, _ = types.ListValueFrom(ctx, types.ObjectType{AttrTypes: Policy_options_Policy_statement_Term_Model{}.AttrTypes()}, policy_options_policy_statement_term_List)
        policy_options_policy_statement_List[i_policy_options_policy_statement] = policy_options_policy_statement_model
			policy_options_policy_statement_List[i_policy_options_policy_statement] = policy_options_policy_statement_model
                
        policy_options_policy_statement_then_List := make([]Policy_options_Policy_statement_Then_Model, len(v_policy_options_policy_statement.Then))
        
		for i_policy_options_policy_statement_then, v_policy_options_policy_statement_then := range v_policy_options_policy_statement.Then {
            var policy_options_policy_statement_then_model Policy_options_Policy_statement_Then_Model
			policy_options_policy_statement_then_List[i_policy_options_policy_statement_then] = policy_options_policy_statement_then_model
                
        policy_options_policy_statement_then_load_balance_List := make([]Policy_options_Policy_statement_Then_Load_balance_Model, len(v_policy_options_policy_statement_then.Load_balance))
        
		for i_policy_options_policy_statement_then_load_balance, v_policy_options_policy_statement_then_load_balance := range v_policy_options_policy_statement_then.Load_balance {
            var policy_options_policy_statement_then_load_balance_model Policy_options_Policy_statement_Then_Load_balance_Model
            policy_options_policy_statement_then_load_balance_model.Per_packet = types.StringPointerValue(v_policy_options_policy_statement_then_load_balance.Per_packet)
			policy_options_policy_statement_then_load_balance_List[i_policy_options_policy_statement_then_load_balance] = policy_options_policy_statement_then_load_balance_model
        }
        policy_options_policy_statement_then_model.Load_balance, _ = types.ListValueFrom(ctx, types.ObjectType{AttrTypes: Policy_options_Policy_statement_Then_Load_balance_Model{}.AttrTypes()}, policy_options_policy_statement_then_load_balance_List)
        policy_options_policy_statement_then_List[i_policy_options_policy_statement_then] = policy_options_policy_statement_then_model
        }
        policy_options_policy_statement_model.Then, _ = types.ListValueFrom(ctx, types.ObjectType{AttrTypes: Policy_options_Policy_statement_Then_Model{}.AttrTypes()}, policy_options_policy_statement_then_List)
        policy_options_policy_statement_List[i_policy_options_policy_statement] = policy_options_policy_statement_model
        }
        policy_options_model.Policy_statement, _ = types.ListValueFrom(ctx, types.ObjectType{AttrTypes: Policy_options_Policy_statement_Model{}.AttrTypes()}, policy_options_policy_statement_List)
        policy_options_List[i_policy_options] = policy_options_model
        policy_options_community_List := make([]Policy_options_Community_Model, len(v_policy_options.Community))
        
		for i_policy_options_community, v_policy_options_community := range v_policy_options.Community {
            var policy_options_community_model Policy_options_Community_Model
            policy_options_community_model.Name = types.StringPointerValue(v_policy_options_community.Name)
			policy_options_community_List[i_policy_options_community] = policy_options_community_model
			var var_policy_options_members []*string
			if v_policy_options_community.Members != nil {
				var_policy_options_members = make([]*string, len(v_policy_options_community.Members))
				copy(var_policy_options_members, v_policy_options_community.Members)
			}
			policy_options_community_model.Members, _ = types.ListValueFrom(ctx, types.StringType, var_policy_options_members)
			policy_options_community_List[i_policy_options_community] = policy_options_community_model
        }
        policy_options_model.Community, _ = types.ListValueFrom(ctx, types.ObjectType{AttrTypes: Policy_options_Community_Model{}.AttrTypes()}, policy_options_community_List)
        policy_options_List[i_policy_options] = policy_options_model
    }
    state.Policy_options, _ = types.ListValueFrom(ctx, types.ObjectType{AttrTypes: Policy_options_Model{}.AttrTypes()}, policy_options_List)
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
    state.Protocols = types.ListNull(types.ObjectType{AttrTypes: Groups_Model{}.AttrTypes()})
    protocols_List := make([]Protocols_Model, len(config.Groups.Protocols))
    for i_protocols, v_protocols := range config.Groups.Protocols {
        var protocols_model Protocols_Model
        protocols_bgp_List := make([]Protocols_Bgp_Model, len(v_protocols.Bgp))
        
		for i_protocols_bgp, v_protocols_bgp := range v_protocols.Bgp {
            var protocols_bgp_model Protocols_Bgp_Model
			protocols_bgp_List[i_protocols_bgp] = protocols_bgp_model
                
        protocols_bgp_group_List := make([]Protocols_Bgp_Group_Model, len(v_protocols_bgp.Group))
        
		for i_protocols_bgp_group, v_protocols_bgp_group := range v_protocols_bgp.Group {
            var protocols_bgp_group_model Protocols_Bgp_Group_Model
            protocols_bgp_group_model.Name = types.StringPointerValue(v_protocols_bgp_group.Name)
			protocols_bgp_group_List[i_protocols_bgp_group] = protocols_bgp_group_model
            protocols_bgp_group_model.Type = types.StringPointerValue(v_protocols_bgp_group.Type)
			protocols_bgp_group_List[i_protocols_bgp_group] = protocols_bgp_group_model
			protocols_bgp_group_List[i_protocols_bgp_group] = protocols_bgp_group_model
                
        protocols_bgp_group_multihop_List := make([]Protocols_Bgp_Group_Multihop_Model, len(v_protocols_bgp_group.Multihop))
        
		for i_protocols_bgp_group_multihop, v_protocols_bgp_group_multihop := range v_protocols_bgp_group.Multihop {
            var protocols_bgp_group_multihop_model Protocols_Bgp_Group_Multihop_Model
            protocols_bgp_group_multihop_model.No_nexthop_change = types.StringPointerValue(v_protocols_bgp_group_multihop.No_nexthop_change)
			protocols_bgp_group_multihop_List[i_protocols_bgp_group_multihop] = protocols_bgp_group_multihop_model
        }
        protocols_bgp_group_model.Multihop, _ = types.ListValueFrom(ctx, types.ObjectType{AttrTypes: Protocols_Bgp_Group_Multihop_Model{}.AttrTypes()}, protocols_bgp_group_multihop_List)
        protocols_bgp_group_List[i_protocols_bgp_group] = protocols_bgp_group_model
            protocols_bgp_group_model.Local_address = types.StringPointerValue(v_protocols_bgp_group.Local_address)
			protocols_bgp_group_List[i_protocols_bgp_group] = protocols_bgp_group_model
            protocols_bgp_group_model.Mtu_discovery = types.StringPointerValue(v_protocols_bgp_group.Mtu_discovery)
			protocols_bgp_group_List[i_protocols_bgp_group] = protocols_bgp_group_model
			var var_protocols_bgp_import []*string
			if v_protocols_bgp_group.Import != nil {
				var_protocols_bgp_import = make([]*string, len(v_protocols_bgp_group.Import))
				copy(var_protocols_bgp_import, v_protocols_bgp_group.Import)
			}
			protocols_bgp_group_model.Import, _ = types.ListValueFrom(ctx, types.StringType, var_protocols_bgp_import)
			protocols_bgp_group_List[i_protocols_bgp_group] = protocols_bgp_group_model
			protocols_bgp_group_List[i_protocols_bgp_group] = protocols_bgp_group_model
                
        protocols_bgp_group_family_List := make([]Protocols_Bgp_Group_Family_Model, len(v_protocols_bgp_group.Family))
        
		for i_protocols_bgp_group_family, v_protocols_bgp_group_family := range v_protocols_bgp_group.Family {
            var protocols_bgp_group_family_model Protocols_Bgp_Group_Family_Model
			protocols_bgp_group_family_List[i_protocols_bgp_group_family] = protocols_bgp_group_family_model
                
        protocols_bgp_group_family_evpn_List := make([]Protocols_Bgp_Group_Family_Evpn_Model, len(v_protocols_bgp_group_family.Evpn))
        
		for i_protocols_bgp_group_family_evpn, v_protocols_bgp_group_family_evpn := range v_protocols_bgp_group_family.Evpn {
            var protocols_bgp_group_family_evpn_model Protocols_Bgp_Group_Family_Evpn_Model
			protocols_bgp_group_family_evpn_List[i_protocols_bgp_group_family_evpn] = protocols_bgp_group_family_evpn_model
                
        protocols_bgp_group_family_evpn_signaling_List := make([]Protocols_Bgp_Group_Family_Evpn_Signaling_Model, len(v_protocols_bgp_group_family_evpn.Signaling))
        
		for i_protocols_bgp_group_family_evpn_signaling, v_protocols_bgp_group_family_evpn_signaling := range v_protocols_bgp_group_family_evpn.Signaling {
            var protocols_bgp_group_family_evpn_signaling_model Protocols_Bgp_Group_Family_Evpn_Signaling_Model
			protocols_bgp_group_family_evpn_signaling_List[i_protocols_bgp_group_family_evpn_signaling] = protocols_bgp_group_family_evpn_signaling_model
                
        protocols_bgp_group_family_evpn_signaling_delay_route_advertisements_List := make([]Protocols_Bgp_Group_Family_Evpn_Signaling_Delay_route_advertisements_Model, len(v_protocols_bgp_group_family_evpn_signaling.Delay_route_advertisements))
        
		for i_protocols_bgp_group_family_evpn_signaling_delay_route_advertisements, v_protocols_bgp_group_family_evpn_signaling_delay_route_advertisements := range v_protocols_bgp_group_family_evpn_signaling.Delay_route_advertisements {
            var protocols_bgp_group_family_evpn_signaling_delay_route_advertisements_model Protocols_Bgp_Group_Family_Evpn_Signaling_Delay_route_advertisements_Model
			protocols_bgp_group_family_evpn_signaling_delay_route_advertisements_List[i_protocols_bgp_group_family_evpn_signaling_delay_route_advertisements] = protocols_bgp_group_family_evpn_signaling_delay_route_advertisements_model
                
        protocols_bgp_group_family_evpn_signaling_delay_route_advertisements_minimum_delay_List := make([]Protocols_Bgp_Group_Family_Evpn_Signaling_Delay_route_advertisements_Minimum_delay_Model, len(v_protocols_bgp_group_family_evpn_signaling_delay_route_advertisements.Minimum_delay))
        
		for i_protocols_bgp_group_family_evpn_signaling_delay_route_advertisements_minimum_delay, v_protocols_bgp_group_family_evpn_signaling_delay_route_advertisements_minimum_delay := range v_protocols_bgp_group_family_evpn_signaling_delay_route_advertisements.Minimum_delay {
            var protocols_bgp_group_family_evpn_signaling_delay_route_advertisements_minimum_delay_model Protocols_Bgp_Group_Family_Evpn_Signaling_Delay_route_advertisements_Minimum_delay_Model
            protocols_bgp_group_family_evpn_signaling_delay_route_advertisements_minimum_delay_model.Routing_uptime = types.StringPointerValue(v_protocols_bgp_group_family_evpn_signaling_delay_route_advertisements_minimum_delay.Routing_uptime)
			protocols_bgp_group_family_evpn_signaling_delay_route_advertisements_minimum_delay_List[i_protocols_bgp_group_family_evpn_signaling_delay_route_advertisements_minimum_delay] = protocols_bgp_group_family_evpn_signaling_delay_route_advertisements_minimum_delay_model
        }
        protocols_bgp_group_family_evpn_signaling_delay_route_advertisements_model.Minimum_delay, _ = types.ListValueFrom(ctx, types.ObjectType{AttrTypes: Protocols_Bgp_Group_Family_Evpn_Signaling_Delay_route_advertisements_Minimum_delay_Model{}.AttrTypes()}, protocols_bgp_group_family_evpn_signaling_delay_route_advertisements_minimum_delay_List)
        protocols_bgp_group_family_evpn_signaling_delay_route_advertisements_List[i_protocols_bgp_group_family_evpn_signaling_delay_route_advertisements] = protocols_bgp_group_family_evpn_signaling_delay_route_advertisements_model
        }
        protocols_bgp_group_family_evpn_signaling_model.Delay_route_advertisements, _ = types.ListValueFrom(ctx, types.ObjectType{AttrTypes: Protocols_Bgp_Group_Family_Evpn_Signaling_Delay_route_advertisements_Model{}.AttrTypes()}, protocols_bgp_group_family_evpn_signaling_delay_route_advertisements_List)
        protocols_bgp_group_family_evpn_signaling_List[i_protocols_bgp_group_family_evpn_signaling] = protocols_bgp_group_family_evpn_signaling_model
        }
        protocols_bgp_group_family_evpn_model.Signaling, _ = types.ListValueFrom(ctx, types.ObjectType{AttrTypes: Protocols_Bgp_Group_Family_Evpn_Signaling_Model{}.AttrTypes()}, protocols_bgp_group_family_evpn_signaling_List)
        protocols_bgp_group_family_evpn_List[i_protocols_bgp_group_family_evpn] = protocols_bgp_group_family_evpn_model
        }
        protocols_bgp_group_family_model.Evpn, _ = types.ListValueFrom(ctx, types.ObjectType{AttrTypes: Protocols_Bgp_Group_Family_Evpn_Model{}.AttrTypes()}, protocols_bgp_group_family_evpn_List)
        protocols_bgp_group_family_List[i_protocols_bgp_group_family] = protocols_bgp_group_family_model
        }
        protocols_bgp_group_model.Family, _ = types.ListValueFrom(ctx, types.ObjectType{AttrTypes: Protocols_Bgp_Group_Family_Model{}.AttrTypes()}, protocols_bgp_group_family_List)
        protocols_bgp_group_List[i_protocols_bgp_group] = protocols_bgp_group_model
			var var_protocols_bgp_export []*string
			if v_protocols_bgp_group.Export != nil {
				var_protocols_bgp_export = make([]*string, len(v_protocols_bgp_group.Export))
				copy(var_protocols_bgp_export, v_protocols_bgp_group.Export)
			}
			protocols_bgp_group_model.Export, _ = types.ListValueFrom(ctx, types.StringType, var_protocols_bgp_export)
			protocols_bgp_group_List[i_protocols_bgp_group] = protocols_bgp_group_model
            protocols_bgp_group_model.Vpn_apply_export = types.StringPointerValue(v_protocols_bgp_group.Vpn_apply_export)
			protocols_bgp_group_List[i_protocols_bgp_group] = protocols_bgp_group_model
            protocols_bgp_group_model.Cluster = types.StringPointerValue(v_protocols_bgp_group.Cluster)
			protocols_bgp_group_List[i_protocols_bgp_group] = protocols_bgp_group_model
			protocols_bgp_group_List[i_protocols_bgp_group] = protocols_bgp_group_model
                
        protocols_bgp_group_local_as_List := make([]Protocols_Bgp_Group_Local_as_Model, len(v_protocols_bgp_group.Local_as))
        
		for i_protocols_bgp_group_local_as, v_protocols_bgp_group_local_as := range v_protocols_bgp_group.Local_as {
            var protocols_bgp_group_local_as_model Protocols_Bgp_Group_Local_as_Model
            protocols_bgp_group_local_as_model.As_number = types.StringPointerValue(v_protocols_bgp_group_local_as.As_number)
			protocols_bgp_group_local_as_List[i_protocols_bgp_group_local_as] = protocols_bgp_group_local_as_model
        }
        protocols_bgp_group_model.Local_as, _ = types.ListValueFrom(ctx, types.ObjectType{AttrTypes: Protocols_Bgp_Group_Local_as_Model{}.AttrTypes()}, protocols_bgp_group_local_as_List)
        protocols_bgp_group_List[i_protocols_bgp_group] = protocols_bgp_group_model
			protocols_bgp_group_List[i_protocols_bgp_group] = protocols_bgp_group_model
                
        protocols_bgp_group_multipath_List := make([]Protocols_Bgp_Group_Multipath_Model, len(v_protocols_bgp_group.Multipath))
        
		for i_protocols_bgp_group_multipath, v_protocols_bgp_group_multipath := range v_protocols_bgp_group.Multipath {
            var protocols_bgp_group_multipath_model Protocols_Bgp_Group_Multipath_Model
            protocols_bgp_group_multipath_model.Multiple_as = types.StringPointerValue(v_protocols_bgp_group_multipath.Multiple_as)
			protocols_bgp_group_multipath_List[i_protocols_bgp_group_multipath] = protocols_bgp_group_multipath_model
        }
        protocols_bgp_group_model.Multipath, _ = types.ListValueFrom(ctx, types.ObjectType{AttrTypes: Protocols_Bgp_Group_Multipath_Model{}.AttrTypes()}, protocols_bgp_group_multipath_List)
        protocols_bgp_group_List[i_protocols_bgp_group] = protocols_bgp_group_model
			protocols_bgp_group_List[i_protocols_bgp_group] = protocols_bgp_group_model
                
        protocols_bgp_group_bfd_liveness_detection_List := make([]Protocols_Bgp_Group_Bfd_liveness_detection_Model, len(v_protocols_bgp_group.Bfd_liveness_detection))
        
		for i_protocols_bgp_group_bfd_liveness_detection, v_protocols_bgp_group_bfd_liveness_detection := range v_protocols_bgp_group.Bfd_liveness_detection {
            var protocols_bgp_group_bfd_liveness_detection_model Protocols_Bgp_Group_Bfd_liveness_detection_Model
            protocols_bgp_group_bfd_liveness_detection_model.Minimum_interval = types.StringPointerValue(v_protocols_bgp_group_bfd_liveness_detection.Minimum_interval)
			protocols_bgp_group_bfd_liveness_detection_List[i_protocols_bgp_group_bfd_liveness_detection] = protocols_bgp_group_bfd_liveness_detection_model
            protocols_bgp_group_bfd_liveness_detection_model.Multiplier = types.StringPointerValue(v_protocols_bgp_group_bfd_liveness_detection.Multiplier)
			protocols_bgp_group_bfd_liveness_detection_List[i_protocols_bgp_group_bfd_liveness_detection] = protocols_bgp_group_bfd_liveness_detection_model
        }
        protocols_bgp_group_model.Bfd_liveness_detection, _ = types.ListValueFrom(ctx, types.ObjectType{AttrTypes: Protocols_Bgp_Group_Bfd_liveness_detection_Model{}.AttrTypes()}, protocols_bgp_group_bfd_liveness_detection_List)
        protocols_bgp_group_List[i_protocols_bgp_group] = protocols_bgp_group_model
			var var_protocols_bgp_allow []*string
			if v_protocols_bgp_group.Allow != nil {
				var_protocols_bgp_allow = make([]*string, len(v_protocols_bgp_group.Allow))
				copy(var_protocols_bgp_allow, v_protocols_bgp_group.Allow)
			}
			protocols_bgp_group_model.Allow, _ = types.ListValueFrom(ctx, types.StringType, var_protocols_bgp_allow)
			protocols_bgp_group_List[i_protocols_bgp_group] = protocols_bgp_group_model
			protocols_bgp_group_List[i_protocols_bgp_group] = protocols_bgp_group_model
                
        protocols_bgp_group_neighbor_List := make([]Protocols_Bgp_Group_Neighbor_Model, len(v_protocols_bgp_group.Neighbor))
        
		for i_protocols_bgp_group_neighbor, v_protocols_bgp_group_neighbor := range v_protocols_bgp_group.Neighbor {
            var protocols_bgp_group_neighbor_model Protocols_Bgp_Group_Neighbor_Model
            protocols_bgp_group_neighbor_model.Name = types.StringPointerValue(v_protocols_bgp_group_neighbor.Name)
			protocols_bgp_group_neighbor_List[i_protocols_bgp_group_neighbor] = protocols_bgp_group_neighbor_model
            protocols_bgp_group_neighbor_model.Description = types.StringPointerValue(v_protocols_bgp_group_neighbor.Description)
			protocols_bgp_group_neighbor_List[i_protocols_bgp_group_neighbor] = protocols_bgp_group_neighbor_model
            protocols_bgp_group_neighbor_model.Peer_as = types.StringPointerValue(v_protocols_bgp_group_neighbor.Peer_as)
			protocols_bgp_group_neighbor_List[i_protocols_bgp_group_neighbor] = protocols_bgp_group_neighbor_model
        }
        protocols_bgp_group_model.Neighbor, _ = types.ListValueFrom(ctx, types.ObjectType{AttrTypes: Protocols_Bgp_Group_Neighbor_Model{}.AttrTypes()}, protocols_bgp_group_neighbor_List)
        protocols_bgp_group_List[i_protocols_bgp_group] = protocols_bgp_group_model
        }
        protocols_bgp_model.Group, _ = types.ListValueFrom(ctx, types.ObjectType{AttrTypes: Protocols_Bgp_Group_Model{}.AttrTypes()}, protocols_bgp_group_List)
        protocols_bgp_List[i_protocols_bgp] = protocols_bgp_model
        }
        protocols_model.Bgp, _ = types.ListValueFrom(ctx, types.ObjectType{AttrTypes: Protocols_Bgp_Model{}.AttrTypes()}, protocols_bgp_List)
        protocols_List[i_protocols] = protocols_model
        protocols_evpn_List := make([]Protocols_Evpn_Model, len(v_protocols.Evpn))
        
		for i_protocols_evpn, v_protocols_evpn := range v_protocols.Evpn {
            var protocols_evpn_model Protocols_Evpn_Model
            protocols_evpn_model.Encapsulation = types.StringPointerValue(v_protocols_evpn.Encapsulation)
			protocols_evpn_List[i_protocols_evpn] = protocols_evpn_model
            protocols_evpn_model.Multicast_mode = types.StringPointerValue(v_protocols_evpn.Multicast_mode)
			protocols_evpn_List[i_protocols_evpn] = protocols_evpn_model
            protocols_evpn_model.Default_gateway = types.StringPointerValue(v_protocols_evpn.Default_gateway)
			protocols_evpn_List[i_protocols_evpn] = protocols_evpn_model
			var var_protocols_extended_vni_list []*string
			if v_protocols_evpn.Extended_vni_list != nil {
				var_protocols_extended_vni_list = make([]*string, len(v_protocols_evpn.Extended_vni_list))
				copy(var_protocols_extended_vni_list, v_protocols_evpn.Extended_vni_list)
			}
			protocols_evpn_model.Extended_vni_list, _ = types.ListValueFrom(ctx, types.StringType, var_protocols_extended_vni_list)
			protocols_evpn_List[i_protocols_evpn] = protocols_evpn_model
            protocols_evpn_model.No_core_isolation = types.StringPointerValue(v_protocols_evpn.No_core_isolation)
			protocols_evpn_List[i_protocols_evpn] = protocols_evpn_model
        }
        protocols_model.Evpn, _ = types.ListValueFrom(ctx, types.ObjectType{AttrTypes: Protocols_Evpn_Model{}.AttrTypes()}, protocols_evpn_List)
        protocols_List[i_protocols] = protocols_model
        protocols_lldp_List := make([]Protocols_Lldp_Model, len(v_protocols.Lldp))
        
		for i_protocols_lldp, v_protocols_lldp := range v_protocols.Lldp {
            var protocols_lldp_model Protocols_Lldp_Model
			protocols_lldp_List[i_protocols_lldp] = protocols_lldp_model
                
        protocols_lldp_interface_List := make([]Protocols_Lldp_Interface_Model, len(v_protocols_lldp.Interface))
        
		for i_protocols_lldp_interface, v_protocols_lldp_interface := range v_protocols_lldp.Interface {
            var protocols_lldp_interface_model Protocols_Lldp_Interface_Model
            protocols_lldp_interface_model.Name = types.StringPointerValue(v_protocols_lldp_interface.Name)
			protocols_lldp_interface_List[i_protocols_lldp_interface] = protocols_lldp_interface_model
        }
        protocols_lldp_model.Interface, _ = types.ListValueFrom(ctx, types.ObjectType{AttrTypes: Protocols_Lldp_Interface_Model{}.AttrTypes()}, protocols_lldp_interface_List)
        protocols_lldp_List[i_protocols_lldp] = protocols_lldp_model
        }
        protocols_model.Lldp, _ = types.ListValueFrom(ctx, types.ObjectType{AttrTypes: Protocols_Lldp_Model{}.AttrTypes()}, protocols_lldp_List)
        protocols_List[i_protocols] = protocols_model
        protocols_igmp_snooping_List := make([]Protocols_Igmp_snooping_Model, len(v_protocols.Igmp_snooping))
        
		for i_protocols_igmp_snooping, v_protocols_igmp_snooping := range v_protocols.Igmp_snooping {
            var protocols_igmp_snooping_model Protocols_Igmp_snooping_Model
			protocols_igmp_snooping_List[i_protocols_igmp_snooping] = protocols_igmp_snooping_model
                
        protocols_igmp_snooping_vlan_List := make([]Protocols_Igmp_snooping_Vlan_Model, len(v_protocols_igmp_snooping.Vlan))
        
		for i_protocols_igmp_snooping_vlan, v_protocols_igmp_snooping_vlan := range v_protocols_igmp_snooping.Vlan {
            var protocols_igmp_snooping_vlan_model Protocols_Igmp_snooping_Vlan_Model
            protocols_igmp_snooping_vlan_model.Name = types.StringPointerValue(v_protocols_igmp_snooping_vlan.Name)
			protocols_igmp_snooping_vlan_List[i_protocols_igmp_snooping_vlan] = protocols_igmp_snooping_vlan_model
        }
        protocols_igmp_snooping_model.Vlan, _ = types.ListValueFrom(ctx, types.ObjectType{AttrTypes: Protocols_Igmp_snooping_Vlan_Model{}.AttrTypes()}, protocols_igmp_snooping_vlan_List)
        protocols_igmp_snooping_List[i_protocols_igmp_snooping] = protocols_igmp_snooping_model
        }
        protocols_model.Igmp_snooping, _ = types.ListValueFrom(ctx, types.ObjectType{AttrTypes: Protocols_Igmp_snooping_Model{}.AttrTypes()}, protocols_igmp_snooping_List)
        protocols_List[i_protocols] = protocols_model
    }
    state.Protocols, _ = types.ListValueFrom(ctx, types.ObjectType{AttrTypes: Protocols_Model{}.AttrTypes()}, protocols_List)
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
    state.Routing_instances = types.ListNull(types.ObjectType{AttrTypes: Groups_Model{}.AttrTypes()})
    routing_instances_List := make([]Routing_instances_Model, len(config.Groups.Routing_instances))
    for i_routing_instances, v_routing_instances := range config.Groups.Routing_instances {
        var routing_instances_model Routing_instances_Model
        routing_instances_instance_List := make([]Routing_instances_Instance_Model, len(v_routing_instances.Instance))
        
		for i_routing_instances_instance, v_routing_instances_instance := range v_routing_instances.Instance {
            var routing_instances_instance_model Routing_instances_Instance_Model
            routing_instances_instance_model.Name = types.StringPointerValue(v_routing_instances_instance.Name)
			routing_instances_instance_List[i_routing_instances_instance] = routing_instances_instance_model
            routing_instances_instance_model.Instance_type = types.StringPointerValue(v_routing_instances_instance.Instance_type)
			routing_instances_instance_List[i_routing_instances_instance] = routing_instances_instance_model
			routing_instances_instance_List[i_routing_instances_instance] = routing_instances_instance_model
                
        routing_instances_instance_interface_List := make([]Routing_instances_Instance_Interface_Model, len(v_routing_instances_instance.Interface))
        
		for i_routing_instances_instance_interface, v_routing_instances_instance_interface := range v_routing_instances_instance.Interface {
            var routing_instances_instance_interface_model Routing_instances_Instance_Interface_Model
            routing_instances_instance_interface_model.Name = types.StringPointerValue(v_routing_instances_instance_interface.Name)
			routing_instances_instance_interface_List[i_routing_instances_instance_interface] = routing_instances_instance_interface_model
        }
        routing_instances_instance_model.Interface, _ = types.ListValueFrom(ctx, types.ObjectType{AttrTypes: Routing_instances_Instance_Interface_Model{}.AttrTypes()}, routing_instances_instance_interface_List)
        routing_instances_instance_List[i_routing_instances_instance] = routing_instances_instance_model
			routing_instances_instance_List[i_routing_instances_instance] = routing_instances_instance_model
                
        routing_instances_instance_route_distinguisher_List := make([]Routing_instances_Instance_Route_distinguisher_Model, len(v_routing_instances_instance.Route_distinguisher))
        
		for i_routing_instances_instance_route_distinguisher, v_routing_instances_instance_route_distinguisher := range v_routing_instances_instance.Route_distinguisher {
            var routing_instances_instance_route_distinguisher_model Routing_instances_Instance_Route_distinguisher_Model
            routing_instances_instance_route_distinguisher_model.Rd_type = types.StringPointerValue(v_routing_instances_instance_route_distinguisher.Rd_type)
			routing_instances_instance_route_distinguisher_List[i_routing_instances_instance_route_distinguisher] = routing_instances_instance_route_distinguisher_model
        }
        routing_instances_instance_model.Route_distinguisher, _ = types.ListValueFrom(ctx, types.ObjectType{AttrTypes: Routing_instances_Instance_Route_distinguisher_Model{}.AttrTypes()}, routing_instances_instance_route_distinguisher_List)
        routing_instances_instance_List[i_routing_instances_instance] = routing_instances_instance_model
			routing_instances_instance_List[i_routing_instances_instance] = routing_instances_instance_model
                
        routing_instances_instance_vrf_target_List := make([]Routing_instances_Instance_Vrf_target_Model, len(v_routing_instances_instance.Vrf_target))
        
		for i_routing_instances_instance_vrf_target, v_routing_instances_instance_vrf_target := range v_routing_instances_instance.Vrf_target {
            var routing_instances_instance_vrf_target_model Routing_instances_Instance_Vrf_target_Model
            routing_instances_instance_vrf_target_model.Community = types.StringPointerValue(v_routing_instances_instance_vrf_target.Community)
			routing_instances_instance_vrf_target_List[i_routing_instances_instance_vrf_target] = routing_instances_instance_vrf_target_model
        }
        routing_instances_instance_model.Vrf_target, _ = types.ListValueFrom(ctx, types.ObjectType{AttrTypes: Routing_instances_Instance_Vrf_target_Model{}.AttrTypes()}, routing_instances_instance_vrf_target_List)
        routing_instances_instance_List[i_routing_instances_instance] = routing_instances_instance_model
			routing_instances_instance_List[i_routing_instances_instance] = routing_instances_instance_model
                
        routing_instances_instance_vrf_table_label_List := make([]Routing_instances_Instance_Vrf_table_label_Model, len(v_routing_instances_instance.Vrf_table_label))
        
        routing_instances_instance_model.Vrf_table_label, _ = types.ListValueFrom(ctx, types.ObjectType{AttrTypes: Routing_instances_Instance_Vrf_table_label_Model{}.AttrTypes()}, routing_instances_instance_vrf_table_label_List)
        routing_instances_instance_List[i_routing_instances_instance] = routing_instances_instance_model
			routing_instances_instance_List[i_routing_instances_instance] = routing_instances_instance_model
                
        routing_instances_instance_routing_options_List := make([]Routing_instances_Instance_Routing_options_Model, len(v_routing_instances_instance.Routing_options))
        
		for i_routing_instances_instance_routing_options, v_routing_instances_instance_routing_options := range v_routing_instances_instance.Routing_options {
            var routing_instances_instance_routing_options_model Routing_instances_Instance_Routing_options_Model
			routing_instances_instance_routing_options_List[i_routing_instances_instance_routing_options] = routing_instances_instance_routing_options_model
                
        routing_instances_instance_routing_options_auto_export_List := make([]Routing_instances_Instance_Routing_options_Auto_export_Model, len(v_routing_instances_instance_routing_options.Auto_export))
        
        routing_instances_instance_routing_options_model.Auto_export, _ = types.ListValueFrom(ctx, types.ObjectType{AttrTypes: Routing_instances_Instance_Routing_options_Auto_export_Model{}.AttrTypes()}, routing_instances_instance_routing_options_auto_export_List)
        routing_instances_instance_routing_options_List[i_routing_instances_instance_routing_options] = routing_instances_instance_routing_options_model
        }
        routing_instances_instance_model.Routing_options, _ = types.ListValueFrom(ctx, types.ObjectType{AttrTypes: Routing_instances_Instance_Routing_options_Model{}.AttrTypes()}, routing_instances_instance_routing_options_List)
        routing_instances_instance_List[i_routing_instances_instance] = routing_instances_instance_model
			routing_instances_instance_List[i_routing_instances_instance] = routing_instances_instance_model
                
        routing_instances_instance_protocols_List := make([]Routing_instances_Instance_Protocols_Model, len(v_routing_instances_instance.Protocols))
        
		for i_routing_instances_instance_protocols, v_routing_instances_instance_protocols := range v_routing_instances_instance.Protocols {
            var routing_instances_instance_protocols_model Routing_instances_Instance_Protocols_Model
			routing_instances_instance_protocols_List[i_routing_instances_instance_protocols] = routing_instances_instance_protocols_model
                
        routing_instances_instance_protocols_ospf_List := make([]Routing_instances_Instance_Protocols_Ospf_Model, len(v_routing_instances_instance_protocols.Ospf))
        
		for i_routing_instances_instance_protocols_ospf, v_routing_instances_instance_protocols_ospf := range v_routing_instances_instance_protocols.Ospf {
            var routing_instances_instance_protocols_ospf_model Routing_instances_Instance_Protocols_Ospf_Model
			var var_routing_instances_instance_protocols_export []*string
			if v_routing_instances_instance_protocols_ospf.Export != nil {
				var_routing_instances_instance_protocols_export = make([]*string, len(v_routing_instances_instance_protocols_ospf.Export))
				copy(var_routing_instances_instance_protocols_export, v_routing_instances_instance_protocols_ospf.Export)
			}
			routing_instances_instance_protocols_ospf_model.Export, _ = types.ListValueFrom(ctx, types.StringType, var_routing_instances_instance_protocols_export)
			routing_instances_instance_protocols_ospf_List[i_routing_instances_instance_protocols_ospf] = routing_instances_instance_protocols_ospf_model
			routing_instances_instance_protocols_ospf_List[i_routing_instances_instance_protocols_ospf] = routing_instances_instance_protocols_ospf_model
                
        routing_instances_instance_protocols_ospf_area_List := make([]Routing_instances_Instance_Protocols_Ospf_Area_Model, len(v_routing_instances_instance_protocols_ospf.Area))
        
		for i_routing_instances_instance_protocols_ospf_area, v_routing_instances_instance_protocols_ospf_area := range v_routing_instances_instance_protocols_ospf.Area {
            var routing_instances_instance_protocols_ospf_area_model Routing_instances_Instance_Protocols_Ospf_Area_Model
            routing_instances_instance_protocols_ospf_area_model.Name = types.StringPointerValue(v_routing_instances_instance_protocols_ospf_area.Name)
			routing_instances_instance_protocols_ospf_area_List[i_routing_instances_instance_protocols_ospf_area] = routing_instances_instance_protocols_ospf_area_model
			routing_instances_instance_protocols_ospf_area_List[i_routing_instances_instance_protocols_ospf_area] = routing_instances_instance_protocols_ospf_area_model
                
        routing_instances_instance_protocols_ospf_area_interface_List := make([]Routing_instances_Instance_Protocols_Ospf_Area_Interface_Model, len(v_routing_instances_instance_protocols_ospf_area.Interface))
        
		for i_routing_instances_instance_protocols_ospf_area_interface, v_routing_instances_instance_protocols_ospf_area_interface := range v_routing_instances_instance_protocols_ospf_area.Interface {
            var routing_instances_instance_protocols_ospf_area_interface_model Routing_instances_Instance_Protocols_Ospf_Area_Interface_Model
            routing_instances_instance_protocols_ospf_area_interface_model.Name = types.StringPointerValue(v_routing_instances_instance_protocols_ospf_area_interface.Name)
			routing_instances_instance_protocols_ospf_area_interface_List[i_routing_instances_instance_protocols_ospf_area_interface] = routing_instances_instance_protocols_ospf_area_interface_model
            routing_instances_instance_protocols_ospf_area_interface_model.Metric = types.StringPointerValue(v_routing_instances_instance_protocols_ospf_area_interface.Metric)
			routing_instances_instance_protocols_ospf_area_interface_List[i_routing_instances_instance_protocols_ospf_area_interface] = routing_instances_instance_protocols_ospf_area_interface_model
        }
        routing_instances_instance_protocols_ospf_area_model.Interface, _ = types.ListValueFrom(ctx, types.ObjectType{AttrTypes: Routing_instances_Instance_Protocols_Ospf_Area_Interface_Model{}.AttrTypes()}, routing_instances_instance_protocols_ospf_area_interface_List)
        routing_instances_instance_protocols_ospf_area_List[i_routing_instances_instance_protocols_ospf_area] = routing_instances_instance_protocols_ospf_area_model
        }
        routing_instances_instance_protocols_ospf_model.Area, _ = types.ListValueFrom(ctx, types.ObjectType{AttrTypes: Routing_instances_Instance_Protocols_Ospf_Area_Model{}.AttrTypes()}, routing_instances_instance_protocols_ospf_area_List)
        routing_instances_instance_protocols_ospf_List[i_routing_instances_instance_protocols_ospf] = routing_instances_instance_protocols_ospf_model
        }
        routing_instances_instance_protocols_model.Ospf, _ = types.ListValueFrom(ctx, types.ObjectType{AttrTypes: Routing_instances_Instance_Protocols_Ospf_Model{}.AttrTypes()}, routing_instances_instance_protocols_ospf_List)
        routing_instances_instance_protocols_List[i_routing_instances_instance_protocols] = routing_instances_instance_protocols_model
			routing_instances_instance_protocols_List[i_routing_instances_instance_protocols] = routing_instances_instance_protocols_model
                
        routing_instances_instance_protocols_evpn_List := make([]Routing_instances_Instance_Protocols_Evpn_Model, len(v_routing_instances_instance_protocols.Evpn))
        
		for i_routing_instances_instance_protocols_evpn, v_routing_instances_instance_protocols_evpn := range v_routing_instances_instance_protocols.Evpn {
            var routing_instances_instance_protocols_evpn_model Routing_instances_Instance_Protocols_Evpn_Model
			routing_instances_instance_protocols_evpn_List[i_routing_instances_instance_protocols_evpn] = routing_instances_instance_protocols_evpn_model
                
        routing_instances_instance_protocols_evpn_ip_prefix_routes_List := make([]Routing_instances_Instance_Protocols_Evpn_Ip_prefix_routes_Model, len(v_routing_instances_instance_protocols_evpn.Ip_prefix_routes))
        
		for i_routing_instances_instance_protocols_evpn_ip_prefix_routes, v_routing_instances_instance_protocols_evpn_ip_prefix_routes := range v_routing_instances_instance_protocols_evpn.Ip_prefix_routes {
            var routing_instances_instance_protocols_evpn_ip_prefix_routes_model Routing_instances_Instance_Protocols_Evpn_Ip_prefix_routes_Model
            routing_instances_instance_protocols_evpn_ip_prefix_routes_model.Advertise = types.StringPointerValue(v_routing_instances_instance_protocols_evpn_ip_prefix_routes.Advertise)
			routing_instances_instance_protocols_evpn_ip_prefix_routes_List[i_routing_instances_instance_protocols_evpn_ip_prefix_routes] = routing_instances_instance_protocols_evpn_ip_prefix_routes_model
            routing_instances_instance_protocols_evpn_ip_prefix_routes_model.Encapsulation = types.StringPointerValue(v_routing_instances_instance_protocols_evpn_ip_prefix_routes.Encapsulation)
			routing_instances_instance_protocols_evpn_ip_prefix_routes_List[i_routing_instances_instance_protocols_evpn_ip_prefix_routes] = routing_instances_instance_protocols_evpn_ip_prefix_routes_model
            routing_instances_instance_protocols_evpn_ip_prefix_routes_model.Vni = types.StringPointerValue(v_routing_instances_instance_protocols_evpn_ip_prefix_routes.Vni)
			routing_instances_instance_protocols_evpn_ip_prefix_routes_List[i_routing_instances_instance_protocols_evpn_ip_prefix_routes] = routing_instances_instance_protocols_evpn_ip_prefix_routes_model
			var var_routing_instances_instance_protocols_evpn_export []*string
			if v_routing_instances_instance_protocols_evpn_ip_prefix_routes.Export != nil {
				var_routing_instances_instance_protocols_evpn_export = make([]*string, len(v_routing_instances_instance_protocols_evpn_ip_prefix_routes.Export))
				copy(var_routing_instances_instance_protocols_evpn_export, v_routing_instances_instance_protocols_evpn_ip_prefix_routes.Export)
			}
			routing_instances_instance_protocols_evpn_ip_prefix_routes_model.Export, _ = types.ListValueFrom(ctx, types.StringType, var_routing_instances_instance_protocols_evpn_export)
			routing_instances_instance_protocols_evpn_ip_prefix_routes_List[i_routing_instances_instance_protocols_evpn_ip_prefix_routes] = routing_instances_instance_protocols_evpn_ip_prefix_routes_model
        }
        routing_instances_instance_protocols_evpn_model.Ip_prefix_routes, _ = types.ListValueFrom(ctx, types.ObjectType{AttrTypes: Routing_instances_Instance_Protocols_Evpn_Ip_prefix_routes_Model{}.AttrTypes()}, routing_instances_instance_protocols_evpn_ip_prefix_routes_List)
        routing_instances_instance_protocols_evpn_List[i_routing_instances_instance_protocols_evpn] = routing_instances_instance_protocols_evpn_model
        }
        routing_instances_instance_protocols_model.Evpn, _ = types.ListValueFrom(ctx, types.ObjectType{AttrTypes: Routing_instances_Instance_Protocols_Evpn_Model{}.AttrTypes()}, routing_instances_instance_protocols_evpn_List)
        routing_instances_instance_protocols_List[i_routing_instances_instance_protocols] = routing_instances_instance_protocols_model
        }
        routing_instances_instance_model.Protocols, _ = types.ListValueFrom(ctx, types.ObjectType{AttrTypes: Routing_instances_Instance_Protocols_Model{}.AttrTypes()}, routing_instances_instance_protocols_List)
        routing_instances_instance_List[i_routing_instances_instance] = routing_instances_instance_model
        }
        routing_instances_model.Instance, _ = types.ListValueFrom(ctx, types.ObjectType{AttrTypes: Routing_instances_Instance_Model{}.AttrTypes()}, routing_instances_instance_List)
        routing_instances_List[i_routing_instances] = routing_instances_model
    }
    state.Routing_instances, _ = types.ListValueFrom(ctx, types.ObjectType{AttrTypes: Routing_instances_Model{}.AttrTypes()}, routing_instances_List)
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
    state.Routing_options = types.ListNull(types.ObjectType{AttrTypes: Groups_Model{}.AttrTypes()})
    routing_options_List := make([]Routing_options_Model, len(config.Groups.Routing_options))
    for i_routing_options, v_routing_options := range config.Groups.Routing_options {
        var routing_options_model Routing_options_Model
        routing_options_static_List := make([]Routing_options_Static_Model, len(v_routing_options.Static))
        
		for i_routing_options_static, v_routing_options_static := range v_routing_options.Static {
            var routing_options_static_model Routing_options_Static_Model
			routing_options_static_List[i_routing_options_static] = routing_options_static_model
                
        routing_options_static_route_List := make([]Routing_options_Static_Route_Model, len(v_routing_options_static.Route))
        
		for i_routing_options_static_route, v_routing_options_static_route := range v_routing_options_static.Route {
            var routing_options_static_route_model Routing_options_Static_Route_Model
            routing_options_static_route_model.Name = types.StringPointerValue(v_routing_options_static_route.Name)
			routing_options_static_route_List[i_routing_options_static_route] = routing_options_static_route_model
			var var_routing_options_static_next_hop []*string
			if v_routing_options_static_route.Next_hop != nil {
				var_routing_options_static_next_hop = make([]*string, len(v_routing_options_static_route.Next_hop))
				copy(var_routing_options_static_next_hop, v_routing_options_static_route.Next_hop)
			}
			routing_options_static_route_model.Next_hop, _ = types.ListValueFrom(ctx, types.StringType, var_routing_options_static_next_hop)
			routing_options_static_route_List[i_routing_options_static_route] = routing_options_static_route_model
        }
        routing_options_static_model.Route, _ = types.ListValueFrom(ctx, types.ObjectType{AttrTypes: Routing_options_Static_Route_Model{}.AttrTypes()}, routing_options_static_route_List)
        routing_options_static_List[i_routing_options_static] = routing_options_static_model
        }
        routing_options_model.Static, _ = types.ListValueFrom(ctx, types.ObjectType{AttrTypes: Routing_options_Static_Model{}.AttrTypes()}, routing_options_static_List)
        routing_options_List[i_routing_options] = routing_options_model
        routing_options_forwarding_table_List := make([]Routing_options_Forwarding_table_Model, len(v_routing_options.Forwarding_table))
        
		for i_routing_options_forwarding_table, v_routing_options_forwarding_table := range v_routing_options.Forwarding_table {
            var routing_options_forwarding_table_model Routing_options_Forwarding_table_Model
			var var_routing_options_export []*string
			if v_routing_options_forwarding_table.Export != nil {
				var_routing_options_export = make([]*string, len(v_routing_options_forwarding_table.Export))
				copy(var_routing_options_export, v_routing_options_forwarding_table.Export)
			}
			routing_options_forwarding_table_model.Export, _ = types.ListValueFrom(ctx, types.StringType, var_routing_options_export)
			routing_options_forwarding_table_List[i_routing_options_forwarding_table] = routing_options_forwarding_table_model
            routing_options_forwarding_table_model.Ecmp_fast_reroute = types.StringPointerValue(v_routing_options_forwarding_table.Ecmp_fast_reroute)
			routing_options_forwarding_table_List[i_routing_options_forwarding_table] = routing_options_forwarding_table_model
			routing_options_forwarding_table_List[i_routing_options_forwarding_table] = routing_options_forwarding_table_model
                
        routing_options_forwarding_table_chained_composite_next_hop_List := make([]Routing_options_Forwarding_table_Chained_composite_next_hop_Model, len(v_routing_options_forwarding_table.Chained_composite_next_hop))
        
		for i_routing_options_forwarding_table_chained_composite_next_hop, v_routing_options_forwarding_table_chained_composite_next_hop := range v_routing_options_forwarding_table.Chained_composite_next_hop {
            var routing_options_forwarding_table_chained_composite_next_hop_model Routing_options_Forwarding_table_Chained_composite_next_hop_Model
			routing_options_forwarding_table_chained_composite_next_hop_List[i_routing_options_forwarding_table_chained_composite_next_hop] = routing_options_forwarding_table_chained_composite_next_hop_model
                
        routing_options_forwarding_table_chained_composite_next_hop_ingress_List := make([]Routing_options_Forwarding_table_Chained_composite_next_hop_Ingress_Model, len(v_routing_options_forwarding_table_chained_composite_next_hop.Ingress))
        
		for i_routing_options_forwarding_table_chained_composite_next_hop_ingress, v_routing_options_forwarding_table_chained_composite_next_hop_ingress := range v_routing_options_forwarding_table_chained_composite_next_hop.Ingress {
            var routing_options_forwarding_table_chained_composite_next_hop_ingress_model Routing_options_Forwarding_table_Chained_composite_next_hop_Ingress_Model
            routing_options_forwarding_table_chained_composite_next_hop_ingress_model.Evpn = types.StringPointerValue(v_routing_options_forwarding_table_chained_composite_next_hop_ingress.Evpn)
			routing_options_forwarding_table_chained_composite_next_hop_ingress_List[i_routing_options_forwarding_table_chained_composite_next_hop_ingress] = routing_options_forwarding_table_chained_composite_next_hop_ingress_model
        }
        routing_options_forwarding_table_chained_composite_next_hop_model.Ingress, _ = types.ListValueFrom(ctx, types.ObjectType{AttrTypes: Routing_options_Forwarding_table_Chained_composite_next_hop_Ingress_Model{}.AttrTypes()}, routing_options_forwarding_table_chained_composite_next_hop_ingress_List)
        routing_options_forwarding_table_chained_composite_next_hop_List[i_routing_options_forwarding_table_chained_composite_next_hop] = routing_options_forwarding_table_chained_composite_next_hop_model
        }
        routing_options_forwarding_table_model.Chained_composite_next_hop, _ = types.ListValueFrom(ctx, types.ObjectType{AttrTypes: Routing_options_Forwarding_table_Chained_composite_next_hop_Model{}.AttrTypes()}, routing_options_forwarding_table_chained_composite_next_hop_List)
        routing_options_forwarding_table_List[i_routing_options_forwarding_table] = routing_options_forwarding_table_model
        }
        routing_options_model.Forwarding_table, _ = types.ListValueFrom(ctx, types.ObjectType{AttrTypes: Routing_options_Forwarding_table_Model{}.AttrTypes()}, routing_options_forwarding_table_List)
        routing_options_List[i_routing_options] = routing_options_model
    }
    state.Routing_options, _ = types.ListValueFrom(ctx, types.ObjectType{AttrTypes: Routing_options_Model{}.AttrTypes()}, routing_options_List)
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
    state.Snmp = types.ListNull(types.ObjectType{AttrTypes: Groups_Model{}.AttrTypes()})
    snmp_List := make([]Snmp_Model, len(config.Groups.Snmp))
    for i_snmp, v_snmp := range config.Groups.Snmp {
        var snmp_model Snmp_Model
        snmp_community_List := make([]Snmp_Community_Model, len(v_snmp.Community))
        
		for i_snmp_community, v_snmp_community := range v_snmp.Community {
            var snmp_community_model Snmp_Community_Model
            snmp_community_model.Name = types.StringPointerValue(v_snmp_community.Name)
			snmp_community_List[i_snmp_community] = snmp_community_model
            snmp_community_model.Authorization = types.StringPointerValue(v_snmp_community.Authorization)
			snmp_community_List[i_snmp_community] = snmp_community_model
        }
        snmp_model.Community, _ = types.ListValueFrom(ctx, types.ObjectType{AttrTypes: Snmp_Community_Model{}.AttrTypes()}, snmp_community_List)
        snmp_List[i_snmp] = snmp_model
    }
    state.Snmp, _ = types.ListValueFrom(ctx, types.ObjectType{AttrTypes: Snmp_Model{}.AttrTypes()}, snmp_List)
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
    state.Switch_options = types.ListNull(types.ObjectType{AttrTypes: Groups_Model{}.AttrTypes()})
    switch_options_List := make([]Switch_options_Model, len(config.Groups.Switch_options))
    for i_switch_options, v_switch_options := range config.Groups.Switch_options {
        var switch_options_model Switch_options_Model
        switch_options_vtep_source_interface_List := make([]Switch_options_Vtep_source_interface_Model, len(v_switch_options.Vtep_source_interface))
        
		for i_switch_options_vtep_source_interface, v_switch_options_vtep_source_interface := range v_switch_options.Vtep_source_interface {
            var switch_options_vtep_source_interface_model Switch_options_Vtep_source_interface_Model
            switch_options_vtep_source_interface_model.Interface_name = types.StringPointerValue(v_switch_options_vtep_source_interface.Interface_name)
			switch_options_vtep_source_interface_List[i_switch_options_vtep_source_interface] = switch_options_vtep_source_interface_model
        }
        switch_options_model.Vtep_source_interface, _ = types.ListValueFrom(ctx, types.ObjectType{AttrTypes: Switch_options_Vtep_source_interface_Model{}.AttrTypes()}, switch_options_vtep_source_interface_List)
        switch_options_List[i_switch_options] = switch_options_model
        switch_options_route_distinguisher_List := make([]Switch_options_Route_distinguisher_Model, len(v_switch_options.Route_distinguisher))
        
		for i_switch_options_route_distinguisher, v_switch_options_route_distinguisher := range v_switch_options.Route_distinguisher {
            var switch_options_route_distinguisher_model Switch_options_Route_distinguisher_Model
            switch_options_route_distinguisher_model.Rd_type = types.StringPointerValue(v_switch_options_route_distinguisher.Rd_type)
			switch_options_route_distinguisher_List[i_switch_options_route_distinguisher] = switch_options_route_distinguisher_model
        }
        switch_options_model.Route_distinguisher, _ = types.ListValueFrom(ctx, types.ObjectType{AttrTypes: Switch_options_Route_distinguisher_Model{}.AttrTypes()}, switch_options_route_distinguisher_List)
        switch_options_List[i_switch_options] = switch_options_model
        switch_options_vrf_target_List := make([]Switch_options_Vrf_target_Model, len(v_switch_options.Vrf_target))
        
		for i_switch_options_vrf_target, v_switch_options_vrf_target := range v_switch_options.Vrf_target {
            var switch_options_vrf_target_model Switch_options_Vrf_target_Model
            switch_options_vrf_target_model.Community = types.StringPointerValue(v_switch_options_vrf_target.Community)
			switch_options_vrf_target_List[i_switch_options_vrf_target] = switch_options_vrf_target_model
			switch_options_vrf_target_List[i_switch_options_vrf_target] = switch_options_vrf_target_model
                
        switch_options_vrf_target_auto_List := make([]Switch_options_Vrf_target_Auto_Model, len(v_switch_options_vrf_target.Auto))
        
        switch_options_vrf_target_model.Auto, _ = types.ListValueFrom(ctx, types.ObjectType{AttrTypes: Switch_options_Vrf_target_Auto_Model{}.AttrTypes()}, switch_options_vrf_target_auto_List)
        switch_options_vrf_target_List[i_switch_options_vrf_target] = switch_options_vrf_target_model
        }
        switch_options_model.Vrf_target, _ = types.ListValueFrom(ctx, types.ObjectType{AttrTypes: Switch_options_Vrf_target_Model{}.AttrTypes()}, switch_options_vrf_target_List)
        switch_options_List[i_switch_options] = switch_options_model
    }
    state.Switch_options, _ = types.ListValueFrom(ctx, types.ObjectType{AttrTypes: Switch_options_Model{}.AttrTypes()}, switch_options_List)
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
    state.System = types.ListNull(types.ObjectType{AttrTypes: Groups_Model{}.AttrTypes()})
    system_List := make([]System_Model, len(config.Groups.System))
    for i_system, v_system := range config.Groups.System {
        var system_model System_Model
        system_login_List := make([]System_Login_Model, len(v_system.Login))
        
		for i_system_login, v_system_login := range v_system.Login {
            var system_login_model System_Login_Model
			system_login_List[i_system_login] = system_login_model
                
        system_login_user_List := make([]System_Login_User_Model, len(v_system_login.User))
        
		for i_system_login_user, v_system_login_user := range v_system_login.User {
            var system_login_user_model System_Login_User_Model
            system_login_user_model.Name = types.StringPointerValue(v_system_login_user.Name)
			system_login_user_List[i_system_login_user] = system_login_user_model
            system_login_user_model.Uid = types.StringPointerValue(v_system_login_user.Uid)
			system_login_user_List[i_system_login_user] = system_login_user_model
            system_login_user_model.Class = types.StringPointerValue(v_system_login_user.Class)
			system_login_user_List[i_system_login_user] = system_login_user_model
			system_login_user_List[i_system_login_user] = system_login_user_model
                
        system_login_user_authentication_List := make([]System_Login_User_Authentication_Model, len(v_system_login_user.Authentication))
        
		for i_system_login_user_authentication, v_system_login_user_authentication := range v_system_login_user.Authentication {
            var system_login_user_authentication_model System_Login_User_Authentication_Model
            system_login_user_authentication_model.Encrypted_password = types.StringPointerValue(v_system_login_user_authentication.Encrypted_password)
			system_login_user_authentication_List[i_system_login_user_authentication] = system_login_user_authentication_model
        }
        system_login_user_model.Authentication, _ = types.ListValueFrom(ctx, types.ObjectType{AttrTypes: System_Login_User_Authentication_Model{}.AttrTypes()}, system_login_user_authentication_List)
        system_login_user_List[i_system_login_user] = system_login_user_model
        }
        system_login_model.User, _ = types.ListValueFrom(ctx, types.ObjectType{AttrTypes: System_Login_User_Model{}.AttrTypes()}, system_login_user_List)
        system_login_List[i_system_login] = system_login_model
            system_login_model.Message = types.StringPointerValue(v_system_login.Message)
			system_login_List[i_system_login] = system_login_model
        }
        system_model.Login, _ = types.ListValueFrom(ctx, types.ObjectType{AttrTypes: System_Login_Model{}.AttrTypes()}, system_login_List)
        system_List[i_system] = system_model
        system_root_authentication_List := make([]System_Root_authentication_Model, len(v_system.Root_authentication))
        
		for i_system_root_authentication, v_system_root_authentication := range v_system.Root_authentication {
            var system_root_authentication_model System_Root_authentication_Model
            system_root_authentication_model.Encrypted_password = types.StringPointerValue(v_system_root_authentication.Encrypted_password)
			system_root_authentication_List[i_system_root_authentication] = system_root_authentication_model
        }
        system_model.Root_authentication, _ = types.ListValueFrom(ctx, types.ObjectType{AttrTypes: System_Root_authentication_Model{}.AttrTypes()}, system_root_authentication_List)
        system_List[i_system] = system_model
        system_services_List := make([]System_Services_Model, len(v_system.Services))
        
		for i_system_services, v_system_services := range v_system.Services {
            var system_services_model System_Services_Model
			system_services_List[i_system_services] = system_services_model
                
        system_services_ssh_List := make([]System_Services_Ssh_Model, len(v_system_services.Ssh))
        
		for i_system_services_ssh, v_system_services_ssh := range v_system_services.Ssh {
            var system_services_ssh_model System_Services_Ssh_Model
            system_services_ssh_model.Root_login = types.StringPointerValue(v_system_services_ssh.Root_login)
			system_services_ssh_List[i_system_services_ssh] = system_services_ssh_model
        }
        system_services_model.Ssh, _ = types.ListValueFrom(ctx, types.ObjectType{AttrTypes: System_Services_Ssh_Model{}.AttrTypes()}, system_services_ssh_List)
        system_services_List[i_system_services] = system_services_model
			system_services_List[i_system_services] = system_services_model
                
        system_services_extension_service_List := make([]System_Services_Extension_service_Model, len(v_system_services.Extension_service))
        
		for i_system_services_extension_service, v_system_services_extension_service := range v_system_services.Extension_service {
            var system_services_extension_service_model System_Services_Extension_service_Model
			system_services_extension_service_List[i_system_services_extension_service] = system_services_extension_service_model
                
        system_services_extension_service_request_response_List := make([]System_Services_Extension_service_Request_response_Model, len(v_system_services_extension_service.Request_response))
        
		for i_system_services_extension_service_request_response, v_system_services_extension_service_request_response := range v_system_services_extension_service.Request_response {
            var system_services_extension_service_request_response_model System_Services_Extension_service_Request_response_Model
			system_services_extension_service_request_response_List[i_system_services_extension_service_request_response] = system_services_extension_service_request_response_model
                
        system_services_extension_service_request_response_grpc_List := make([]System_Services_Extension_service_Request_response_Grpc_Model, len(v_system_services_extension_service_request_response.Grpc))
        
		for i_system_services_extension_service_request_response_grpc, v_system_services_extension_service_request_response_grpc := range v_system_services_extension_service_request_response.Grpc {
            var system_services_extension_service_request_response_grpc_model System_Services_Extension_service_Request_response_Grpc_Model
            system_services_extension_service_request_response_grpc_model.Max_connections = types.StringPointerValue(v_system_services_extension_service_request_response_grpc.Max_connections)
			system_services_extension_service_request_response_grpc_List[i_system_services_extension_service_request_response_grpc] = system_services_extension_service_request_response_grpc_model
        }
        system_services_extension_service_request_response_model.Grpc, _ = types.ListValueFrom(ctx, types.ObjectType{AttrTypes: System_Services_Extension_service_Request_response_Grpc_Model{}.AttrTypes()}, system_services_extension_service_request_response_grpc_List)
        system_services_extension_service_request_response_List[i_system_services_extension_service_request_response] = system_services_extension_service_request_response_model
        }
        system_services_extension_service_model.Request_response, _ = types.ListValueFrom(ctx, types.ObjectType{AttrTypes: System_Services_Extension_service_Request_response_Model{}.AttrTypes()}, system_services_extension_service_request_response_List)
        system_services_extension_service_List[i_system_services_extension_service] = system_services_extension_service_model
			system_services_extension_service_List[i_system_services_extension_service] = system_services_extension_service_model
                
        system_services_extension_service_notification_List := make([]System_Services_Extension_service_Notification_Model, len(v_system_services_extension_service.Notification))
        
		for i_system_services_extension_service_notification, v_system_services_extension_service_notification := range v_system_services_extension_service.Notification {
            var system_services_extension_service_notification_model System_Services_Extension_service_Notification_Model
			system_services_extension_service_notification_List[i_system_services_extension_service_notification] = system_services_extension_service_notification_model
                
        system_services_extension_service_notification_allow_clients_List := make([]System_Services_Extension_service_Notification_Allow_clients_Model, len(v_system_services_extension_service_notification.Allow_clients))
        
		for i_system_services_extension_service_notification_allow_clients, v_system_services_extension_service_notification_allow_clients := range v_system_services_extension_service_notification.Allow_clients {
            var system_services_extension_service_notification_allow_clients_model System_Services_Extension_service_Notification_Allow_clients_Model
			var var_system_services_extension_service_notification_address []*string
			if v_system_services_extension_service_notification_allow_clients.Address != nil {
				var_system_services_extension_service_notification_address = make([]*string, len(v_system_services_extension_service_notification_allow_clients.Address))
				copy(var_system_services_extension_service_notification_address, v_system_services_extension_service_notification_allow_clients.Address)
			}
			system_services_extension_service_notification_allow_clients_model.Address, _ = types.ListValueFrom(ctx, types.StringType, var_system_services_extension_service_notification_address)
			system_services_extension_service_notification_allow_clients_List[i_system_services_extension_service_notification_allow_clients] = system_services_extension_service_notification_allow_clients_model
        }
        system_services_extension_service_notification_model.Allow_clients, _ = types.ListValueFrom(ctx, types.ObjectType{AttrTypes: System_Services_Extension_service_Notification_Allow_clients_Model{}.AttrTypes()}, system_services_extension_service_notification_allow_clients_List)
        system_services_extension_service_notification_List[i_system_services_extension_service_notification] = system_services_extension_service_notification_model
        }
        system_services_extension_service_model.Notification, _ = types.ListValueFrom(ctx, types.ObjectType{AttrTypes: System_Services_Extension_service_Notification_Model{}.AttrTypes()}, system_services_extension_service_notification_List)
        system_services_extension_service_List[i_system_services_extension_service] = system_services_extension_service_model
        }
        system_services_model.Extension_service, _ = types.ListValueFrom(ctx, types.ObjectType{AttrTypes: System_Services_Extension_service_Model{}.AttrTypes()}, system_services_extension_service_List)
        system_services_List[i_system_services] = system_services_model
			system_services_List[i_system_services] = system_services_model
                
        system_services_netconf_List := make([]System_Services_Netconf_Model, len(v_system_services.Netconf))
        
		for i_system_services_netconf, v_system_services_netconf := range v_system_services.Netconf {
            var system_services_netconf_model System_Services_Netconf_Model
			system_services_netconf_List[i_system_services_netconf] = system_services_netconf_model
                
        system_services_netconf_ssh_List := make([]System_Services_Netconf_Ssh_Model, len(v_system_services_netconf.Ssh))
        
        system_services_netconf_model.Ssh, _ = types.ListValueFrom(ctx, types.ObjectType{AttrTypes: System_Services_Netconf_Ssh_Model{}.AttrTypes()}, system_services_netconf_ssh_List)
        system_services_netconf_List[i_system_services_netconf] = system_services_netconf_model
        }
        system_services_model.Netconf, _ = types.ListValueFrom(ctx, types.ObjectType{AttrTypes: System_Services_Netconf_Model{}.AttrTypes()}, system_services_netconf_List)
        system_services_List[i_system_services] = system_services_model
			system_services_List[i_system_services] = system_services_model
                
        system_services_rest_List := make([]System_Services_Rest_Model, len(v_system_services.Rest))
        
		for i_system_services_rest, v_system_services_rest := range v_system_services.Rest {
            var system_services_rest_model System_Services_Rest_Model
			system_services_rest_List[i_system_services_rest] = system_services_rest_model
                
        system_services_rest_http_List := make([]System_Services_Rest_Http_Model, len(v_system_services_rest.Http))
        
		for i_system_services_rest_http, v_system_services_rest_http := range v_system_services_rest.Http {
            var system_services_rest_http_model System_Services_Rest_Http_Model
            system_services_rest_http_model.Port = types.StringPointerValue(v_system_services_rest_http.Port)
			system_services_rest_http_List[i_system_services_rest_http] = system_services_rest_http_model
        }
        system_services_rest_model.Http, _ = types.ListValueFrom(ctx, types.ObjectType{AttrTypes: System_Services_Rest_Http_Model{}.AttrTypes()}, system_services_rest_http_List)
        system_services_rest_List[i_system_services_rest] = system_services_rest_model
            system_services_rest_model.Enable_explorer = types.StringPointerValue(v_system_services_rest.Enable_explorer)
			system_services_rest_List[i_system_services_rest] = system_services_rest_model
        }
        system_services_model.Rest, _ = types.ListValueFrom(ctx, types.ObjectType{AttrTypes: System_Services_Rest_Model{}.AttrTypes()}, system_services_rest_List)
        system_services_List[i_system_services] = system_services_model
        }
        system_model.Services, _ = types.ListValueFrom(ctx, types.ObjectType{AttrTypes: System_Services_Model{}.AttrTypes()}, system_services_List)
        system_List[i_system] = system_model
        system_syslog_List := make([]System_Syslog_Model, len(v_system.Syslog))
        
		for i_system_syslog, v_system_syslog := range v_system.Syslog {
            var system_syslog_model System_Syslog_Model
			system_syslog_List[i_system_syslog] = system_syslog_model
                
        system_syslog_user_List := make([]System_Syslog_User_Model, len(v_system_syslog.User))
        
		for i_system_syslog_user, v_system_syslog_user := range v_system_syslog.User {
            var system_syslog_user_model System_Syslog_User_Model
            system_syslog_user_model.Name = types.StringPointerValue(v_system_syslog_user.Name)
			system_syslog_user_List[i_system_syslog_user] = system_syslog_user_model
			system_syslog_user_List[i_system_syslog_user] = system_syslog_user_model
                
        system_syslog_user_contents_List := make([]System_Syslog_User_Contents_Model, len(v_system_syslog_user.Contents))
        
		for i_system_syslog_user_contents, v_system_syslog_user_contents := range v_system_syslog_user.Contents {
            var system_syslog_user_contents_model System_Syslog_User_Contents_Model
            system_syslog_user_contents_model.Name = types.StringPointerValue(v_system_syslog_user_contents.Name)
			system_syslog_user_contents_List[i_system_syslog_user_contents] = system_syslog_user_contents_model
            system_syslog_user_contents_model.Emergency = types.StringPointerValue(v_system_syslog_user_contents.Emergency)
			system_syslog_user_contents_List[i_system_syslog_user_contents] = system_syslog_user_contents_model
        }
        system_syslog_user_model.Contents, _ = types.ListValueFrom(ctx, types.ObjectType{AttrTypes: System_Syslog_User_Contents_Model{}.AttrTypes()}, system_syslog_user_contents_List)
        system_syslog_user_List[i_system_syslog_user] = system_syslog_user_model
        }
        system_syslog_model.User, _ = types.ListValueFrom(ctx, types.ObjectType{AttrTypes: System_Syslog_User_Model{}.AttrTypes()}, system_syslog_user_List)
        system_syslog_List[i_system_syslog] = system_syslog_model
			system_syslog_List[i_system_syslog] = system_syslog_model
                
        system_syslog_file_List := make([]System_Syslog_File_Model, len(v_system_syslog.File))
        
		for i_system_syslog_file, v_system_syslog_file := range v_system_syslog.File {
            var system_syslog_file_model System_Syslog_File_Model
            system_syslog_file_model.Name = types.StringPointerValue(v_system_syslog_file.Name)
			system_syslog_file_List[i_system_syslog_file] = system_syslog_file_model
			system_syslog_file_List[i_system_syslog_file] = system_syslog_file_model
                
        system_syslog_file_contents_List := make([]System_Syslog_File_Contents_Model, len(v_system_syslog_file.Contents))
        
		for i_system_syslog_file_contents, v_system_syslog_file_contents := range v_system_syslog_file.Contents {
            var system_syslog_file_contents_model System_Syslog_File_Contents_Model
            system_syslog_file_contents_model.Name = types.StringPointerValue(v_system_syslog_file_contents.Name)
			system_syslog_file_contents_List[i_system_syslog_file_contents] = system_syslog_file_contents_model
            system_syslog_file_contents_model.Any = types.StringPointerValue(v_system_syslog_file_contents.Any)
			system_syslog_file_contents_List[i_system_syslog_file_contents] = system_syslog_file_contents_model
            system_syslog_file_contents_model.Notice = types.StringPointerValue(v_system_syslog_file_contents.Notice)
			system_syslog_file_contents_List[i_system_syslog_file_contents] = system_syslog_file_contents_model
            system_syslog_file_contents_model.Info = types.StringPointerValue(v_system_syslog_file_contents.Info)
			system_syslog_file_contents_List[i_system_syslog_file_contents] = system_syslog_file_contents_model
        }
        system_syslog_file_model.Contents, _ = types.ListValueFrom(ctx, types.ObjectType{AttrTypes: System_Syslog_File_Contents_Model{}.AttrTypes()}, system_syslog_file_contents_List)
        system_syslog_file_List[i_system_syslog_file] = system_syslog_file_model
        }
        system_syslog_model.File, _ = types.ListValueFrom(ctx, types.ObjectType{AttrTypes: System_Syslog_File_Model{}.AttrTypes()}, system_syslog_file_List)
        system_syslog_List[i_system_syslog] = system_syslog_model
        }
        system_model.Syslog, _ = types.ListValueFrom(ctx, types.ObjectType{AttrTypes: System_Syslog_Model{}.AttrTypes()}, system_syslog_List)
        system_List[i_system] = system_model
        system_extensions_List := make([]System_Extensions_Model, len(v_system.Extensions))
        
		for i_system_extensions, v_system_extensions := range v_system.Extensions {
            var system_extensions_model System_Extensions_Model
			system_extensions_List[i_system_extensions] = system_extensions_model
                
        system_extensions_providers_List := make([]System_Extensions_Providers_Model, len(v_system_extensions.Providers))
        
		for i_system_extensions_providers, v_system_extensions_providers := range v_system_extensions.Providers {
            var system_extensions_providers_model System_Extensions_Providers_Model
            system_extensions_providers_model.Name = types.StringPointerValue(v_system_extensions_providers.Name)
			system_extensions_providers_List[i_system_extensions_providers] = system_extensions_providers_model
			system_extensions_providers_List[i_system_extensions_providers] = system_extensions_providers_model
                
        system_extensions_providers_license_type_List := make([]System_Extensions_Providers_License_type_Model, len(v_system_extensions_providers.License_type))
        
		for i_system_extensions_providers_license_type, v_system_extensions_providers_license_type := range v_system_extensions_providers.License_type {
            var system_extensions_providers_license_type_model System_Extensions_Providers_License_type_Model
            system_extensions_providers_license_type_model.Name = types.StringPointerValue(v_system_extensions_providers_license_type.Name)
			system_extensions_providers_license_type_List[i_system_extensions_providers_license_type] = system_extensions_providers_license_type_model
			var var_system_extensions_providers_deployment_scope []*string
			if v_system_extensions_providers_license_type.Deployment_scope != nil {
				var_system_extensions_providers_deployment_scope = make([]*string, len(v_system_extensions_providers_license_type.Deployment_scope))
				copy(var_system_extensions_providers_deployment_scope, v_system_extensions_providers_license_type.Deployment_scope)
			}
			system_extensions_providers_license_type_model.Deployment_scope, _ = types.ListValueFrom(ctx, types.StringType, var_system_extensions_providers_deployment_scope)
			system_extensions_providers_license_type_List[i_system_extensions_providers_license_type] = system_extensions_providers_license_type_model
        }
        system_extensions_providers_model.License_type, _ = types.ListValueFrom(ctx, types.ObjectType{AttrTypes: System_Extensions_Providers_License_type_Model{}.AttrTypes()}, system_extensions_providers_license_type_List)
        system_extensions_providers_List[i_system_extensions_providers] = system_extensions_providers_model
        }
        system_extensions_model.Providers, _ = types.ListValueFrom(ctx, types.ObjectType{AttrTypes: System_Extensions_Providers_Model{}.AttrTypes()}, system_extensions_providers_List)
        system_extensions_List[i_system_extensions] = system_extensions_model
        }
        system_model.Extensions, _ = types.ListValueFrom(ctx, types.ObjectType{AttrTypes: System_Extensions_Model{}.AttrTypes()}, system_extensions_List)
        system_List[i_system] = system_model
    }
    state.System, _ = types.ListValueFrom(ctx, types.ObjectType{AttrTypes: System_Model{}.AttrTypes()}, system_List)
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
    state.Vlans = types.ListNull(types.ObjectType{AttrTypes: Groups_Model{}.AttrTypes()})
    vlans_List := make([]Vlans_Model, len(config.Groups.Vlans))
    for i_vlans, v_vlans := range config.Groups.Vlans {
        var vlans_model Vlans_Model
        vlans_vlan_List := make([]Vlans_Vlan_Model, len(v_vlans.Vlan))
        
		for i_vlans_vlan, v_vlans_vlan := range v_vlans.Vlan {
            var vlans_vlan_model Vlans_Vlan_Model
            vlans_vlan_model.Name = types.StringPointerValue(v_vlans_vlan.Name)
			vlans_vlan_List[i_vlans_vlan] = vlans_vlan_model
            vlans_vlan_model.Vlan_id = types.StringPointerValue(v_vlans_vlan.Vlan_id)
			vlans_vlan_List[i_vlans_vlan] = vlans_vlan_model
            vlans_vlan_model.L3_interface = types.StringPointerValue(v_vlans_vlan.L3_interface)
			vlans_vlan_List[i_vlans_vlan] = vlans_vlan_model
			vlans_vlan_List[i_vlans_vlan] = vlans_vlan_model
                
        vlans_vlan_vxlan_List := make([]Vlans_Vlan_Vxlan_Model, len(v_vlans_vlan.Vxlan))
        
		for i_vlans_vlan_vxlan, v_vlans_vlan_vxlan := range v_vlans_vlan.Vxlan {
            var vlans_vlan_vxlan_model Vlans_Vlan_Vxlan_Model
            vlans_vlan_vxlan_model.Vni = types.StringPointerValue(v_vlans_vlan_vxlan.Vni)
			vlans_vlan_vxlan_List[i_vlans_vlan_vxlan] = vlans_vlan_vxlan_model
        }
        vlans_vlan_model.Vxlan, _ = types.ListValueFrom(ctx, types.ObjectType{AttrTypes: Vlans_Vlan_Vxlan_Model{}.AttrTypes()}, vlans_vlan_vxlan_List)
        vlans_vlan_List[i_vlans_vlan] = vlans_vlan_model
        }
        vlans_model.Vlan, _ = types.ListValueFrom(ctx, types.ObjectType{AttrTypes: Vlans_Vlan_Model{}.AttrTypes()}, vlans_vlan_List)
        vlans_List[i_vlans] = vlans_model
    }
    state.Vlans, _ = types.ListValueFrom(ctx, types.ObjectType{AttrTypes: Vlans_Model{}.AttrTypes()}, vlans_List)
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
	
    var var_forwarding_options []Forwarding_options_Model
    if plan.Forwarding_options.IsNull() {
        var_forwarding_options = []Forwarding_options_Model{}
    }else {
        resp.Diagnostics.Append(plan.Forwarding_options.ElementsAs(ctx, &var_forwarding_options, false)...)
        if resp.Diagnostics.HasError() {
            return
        }
    }
    config.Groups.Forwarding_options = make([]xml_Forwarding_options, len(var_forwarding_options))
    
    for i_forwarding_options, v_forwarding_options := range var_forwarding_options {
        var var_forwarding_options_storm_control_profiles []Forwarding_options_Storm_control_profiles_Model
        resp.Diagnostics.Append(v_forwarding_options.Storm_control_profiles.ElementsAs(ctx, &var_forwarding_options_storm_control_profiles, false)...)
        if resp.Diagnostics.HasError() {
            return
        }
	    config.Groups.Forwarding_options[i_forwarding_options].Storm_control_profiles = make([]xml_Forwarding_options_Storm_control_profiles, len(var_forwarding_options_storm_control_profiles))
        
		for i_forwarding_options_storm_control_profiles, v_forwarding_options_storm_control_profiles := range var_forwarding_options_storm_control_profiles {
            config.Groups.Forwarding_options[i_forwarding_options].Storm_control_profiles[i_forwarding_options_storm_control_profiles].Name = v_forwarding_options_storm_control_profiles.Name.ValueStringPointer()
            var var_forwarding_options_storm_control_profiles_all []Forwarding_options_Storm_control_profiles_All_Model
            resp.Diagnostics.Append(v_forwarding_options_storm_control_profiles.All.ElementsAs(ctx, &var_forwarding_options_storm_control_profiles_all, false)...)
            if resp.Diagnostics.HasError() {
                return
            }
	    config.Groups.Forwarding_options[i_forwarding_options].Storm_control_profiles[i_forwarding_options_storm_control_profiles].All = make([]xml_Forwarding_options_Storm_control_profiles_All, len(var_forwarding_options_storm_control_profiles_all))
        
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
            config.Groups.Interfaces[i_interfaces].Interface[i_interfaces_interface].Description = v_interfaces_interface.Description.ValueStringPointer()
            config.Groups.Interfaces[i_interfaces].Interface[i_interfaces_interface].Vlan_tagging = v_interfaces_interface.Vlan_tagging.ValueStringPointer()
            var var_interfaces_interface_esi []Interfaces_Interface_Esi_Model
            resp.Diagnostics.Append(v_interfaces_interface.Esi.ElementsAs(ctx, &var_interfaces_interface_esi, false)...)
            if resp.Diagnostics.HasError() {
                return
            }
	    config.Groups.Interfaces[i_interfaces].Interface[i_interfaces_interface].Esi = make([]xml_Interfaces_Interface_Esi, len(var_interfaces_interface_esi))
        
		for i_interfaces_interface_esi, v_interfaces_interface_esi := range var_interfaces_interface_esi {
            config.Groups.Interfaces[i_interfaces].Interface[i_interfaces_interface].Esi[i_interfaces_interface_esi].Identifier = v_interfaces_interface_esi.Identifier.ValueStringPointer()
            config.Groups.Interfaces[i_interfaces].Interface[i_interfaces_interface].Esi[i_interfaces_interface_esi].All_active = v_interfaces_interface_esi.All_active.ValueStringPointer()
        }
            var var_interfaces_interface_ether_options []Interfaces_Interface_Ether_options_Model
            resp.Diagnostics.Append(v_interfaces_interface.Ether_options.ElementsAs(ctx, &var_interfaces_interface_ether_options, false)...)
            if resp.Diagnostics.HasError() {
                return
            }
	    config.Groups.Interfaces[i_interfaces].Interface[i_interfaces_interface].Ether_options = make([]xml_Interfaces_Interface_Ether_options, len(var_interfaces_interface_ether_options))
        
		for i_interfaces_interface_ether_options, v_interfaces_interface_ether_options := range var_interfaces_interface_ether_options {
            var var_interfaces_interface_ether_options_ieee_802_3ad []Interfaces_Interface_Ether_options_Ieee_802_3ad_Model
            resp.Diagnostics.Append(v_interfaces_interface_ether_options.Ieee_802_3ad.ElementsAs(ctx, &var_interfaces_interface_ether_options_ieee_802_3ad, false)...)
            if resp.Diagnostics.HasError() {
                return
            }
	    config.Groups.Interfaces[i_interfaces].Interface[i_interfaces_interface].Ether_options[i_interfaces_interface_ether_options].Ieee_802_3ad = make([]xml_Interfaces_Interface_Ether_options_Ieee_802_3ad, len(var_interfaces_interface_ether_options_ieee_802_3ad))
        
		for i_interfaces_interface_ether_options_ieee_802_3ad, v_interfaces_interface_ether_options_ieee_802_3ad := range var_interfaces_interface_ether_options_ieee_802_3ad {
            config.Groups.Interfaces[i_interfaces].Interface[i_interfaces_interface].Ether_options[i_interfaces_interface_ether_options].Ieee_802_3ad[i_interfaces_interface_ether_options_ieee_802_3ad].Bundle = v_interfaces_interface_ether_options_ieee_802_3ad.Bundle.ValueStringPointer()
        }
        }
            var var_interfaces_interface_aggregated_ether_options []Interfaces_Interface_Aggregated_ether_options_Model
            resp.Diagnostics.Append(v_interfaces_interface.Aggregated_ether_options.ElementsAs(ctx, &var_interfaces_interface_aggregated_ether_options, false)...)
            if resp.Diagnostics.HasError() {
                return
            }
	    config.Groups.Interfaces[i_interfaces].Interface[i_interfaces_interface].Aggregated_ether_options = make([]xml_Interfaces_Interface_Aggregated_ether_options, len(var_interfaces_interface_aggregated_ether_options))
        
		for i_interfaces_interface_aggregated_ether_options, v_interfaces_interface_aggregated_ether_options := range var_interfaces_interface_aggregated_ether_options {
            var var_interfaces_interface_aggregated_ether_options_lacp []Interfaces_Interface_Aggregated_ether_options_Lacp_Model
            resp.Diagnostics.Append(v_interfaces_interface_aggregated_ether_options.Lacp.ElementsAs(ctx, &var_interfaces_interface_aggregated_ether_options_lacp, false)...)
            if resp.Diagnostics.HasError() {
                return
            }
	    config.Groups.Interfaces[i_interfaces].Interface[i_interfaces_interface].Aggregated_ether_options[i_interfaces_interface_aggregated_ether_options].Lacp = make([]xml_Interfaces_Interface_Aggregated_ether_options_Lacp, len(var_interfaces_interface_aggregated_ether_options_lacp))
        
		for i_interfaces_interface_aggregated_ether_options_lacp, v_interfaces_interface_aggregated_ether_options_lacp := range var_interfaces_interface_aggregated_ether_options_lacp {
            config.Groups.Interfaces[i_interfaces].Interface[i_interfaces_interface].Aggregated_ether_options[i_interfaces_interface_aggregated_ether_options].Lacp[i_interfaces_interface_aggregated_ether_options_lacp].Active = v_interfaces_interface_aggregated_ether_options_lacp.Active.ValueStringPointer()
            config.Groups.Interfaces[i_interfaces].Interface[i_interfaces_interface].Aggregated_ether_options[i_interfaces_interface_aggregated_ether_options].Lacp[i_interfaces_interface_aggregated_ether_options_lacp].Periodic = v_interfaces_interface_aggregated_ether_options_lacp.Periodic.ValueStringPointer()
            config.Groups.Interfaces[i_interfaces].Interface[i_interfaces_interface].Aggregated_ether_options[i_interfaces_interface_aggregated_ether_options].Lacp[i_interfaces_interface_aggregated_ether_options_lacp].System_id = v_interfaces_interface_aggregated_ether_options_lacp.System_id.ValueStringPointer()
        }
        }
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
            var var_interfaces_interface_unit_family_ethernet_switching []Interfaces_Interface_Unit_Family_Ethernet_switching_Model
            resp.Diagnostics.Append(v_interfaces_interface_unit_family.Ethernet_switching.ElementsAs(ctx, &var_interfaces_interface_unit_family_ethernet_switching, false)...)
            if resp.Diagnostics.HasError() {
                return
            }
	    config.Groups.Interfaces[i_interfaces].Interface[i_interfaces_interface].Unit[i_interfaces_interface_unit].Family[i_interfaces_interface_unit_family].Ethernet_switching = make([]xml_Interfaces_Interface_Unit_Family_Ethernet_switching, len(var_interfaces_interface_unit_family_ethernet_switching))
        
		for i_interfaces_interface_unit_family_ethernet_switching, v_interfaces_interface_unit_family_ethernet_switching := range var_interfaces_interface_unit_family_ethernet_switching {
            var var_interfaces_interface_unit_family_ethernet_switching_vlan []Interfaces_Interface_Unit_Family_Ethernet_switching_Vlan_Model
            resp.Diagnostics.Append(v_interfaces_interface_unit_family_ethernet_switching.Vlan.ElementsAs(ctx, &var_interfaces_interface_unit_family_ethernet_switching_vlan, false)...)
            if resp.Diagnostics.HasError() {
                return
            }
	    config.Groups.Interfaces[i_interfaces].Interface[i_interfaces_interface].Unit[i_interfaces_interface_unit].Family[i_interfaces_interface_unit_family].Ethernet_switching[i_interfaces_interface_unit_family_ethernet_switching].Vlan = make([]xml_Interfaces_Interface_Unit_Family_Ethernet_switching_Vlan, len(var_interfaces_interface_unit_family_ethernet_switching_vlan))
        
		for i_interfaces_interface_unit_family_ethernet_switching_vlan, v_interfaces_interface_unit_family_ethernet_switching_vlan := range var_interfaces_interface_unit_family_ethernet_switching_vlan {
			var var_interfaces_interface_unit_family_ethernet_switching_vlan_members []string
			resp.Diagnostics.Append(v_interfaces_interface_unit_family_ethernet_switching_vlan.Members.ElementsAs(ctx, &var_interfaces_interface_unit_family_ethernet_switching_vlan_members, false)...)
			for _, v_interfaces_interface_unit_family_ethernet_switching_vlan_members := range var_interfaces_interface_unit_family_ethernet_switching_vlan_members {
				config.Groups.Interfaces[i_interfaces].Interface[i_interfaces_interface].Unit[i_interfaces_interface_unit].Family[i_interfaces_interface_unit_family].Ethernet_switching[i_interfaces_interface_unit_family_ethernet_switching].Vlan[i_interfaces_interface_unit_family_ethernet_switching_vlan].Members = append(config.Groups.Interfaces[i_interfaces].Interface[i_interfaces_interface].Unit[i_interfaces_interface_unit].Family[i_interfaces_interface_unit_family].Ethernet_switching[i_interfaces_interface_unit_family_ethernet_switching].Vlan[i_interfaces_interface_unit_family_ethernet_switching_vlan].Members, &v_interfaces_interface_unit_family_ethernet_switching_vlan_members)
			}
        }
        }
        }
            config.Groups.Interfaces[i_interfaces].Interface[i_interfaces_interface].Unit[i_interfaces_interface_unit].Mac = v_interfaces_interface_unit.Mac.ValueStringPointer()
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
			var var_policy_options_policy_statement_term_from_protocol []string
			resp.Diagnostics.Append(v_policy_options_policy_statement_term_from.Protocol.ElementsAs(ctx, &var_policy_options_policy_statement_term_from_protocol, false)...)
			for _, v_policy_options_policy_statement_term_from_protocol := range var_policy_options_policy_statement_term_from_protocol {
				config.Groups.Policy_options[i_policy_options].Policy_statement[i_policy_options_policy_statement].Term[i_policy_options_policy_statement_term].From[i_policy_options_policy_statement_term_from].Protocol = append(config.Groups.Policy_options[i_policy_options].Policy_statement[i_policy_options_policy_statement].Term[i_policy_options_policy_statement_term].From[i_policy_options_policy_statement_term_from].Protocol, &v_policy_options_policy_statement_term_from_protocol)
			}
            var var_policy_options_policy_statement_term_from_route_filter []Policy_options_Policy_statement_Term_From_Route_filter_Model
            resp.Diagnostics.Append(v_policy_options_policy_statement_term_from.Route_filter.ElementsAs(ctx, &var_policy_options_policy_statement_term_from_route_filter, false)...)
            if resp.Diagnostics.HasError() {
                return
            }
	    config.Groups.Policy_options[i_policy_options].Policy_statement[i_policy_options_policy_statement].Term[i_policy_options_policy_statement_term].From[i_policy_options_policy_statement_term_from].Route_filter = make([]xml_Policy_options_Policy_statement_Term_From_Route_filter, len(var_policy_options_policy_statement_term_from_route_filter))
        
		for i_policy_options_policy_statement_term_from_route_filter, v_policy_options_policy_statement_term_from_route_filter := range var_policy_options_policy_statement_term_from_route_filter {
            config.Groups.Policy_options[i_policy_options].Policy_statement[i_policy_options_policy_statement].Term[i_policy_options_policy_statement_term].From[i_policy_options_policy_statement_term_from].Route_filter[i_policy_options_policy_statement_term_from_route_filter].Address = v_policy_options_policy_statement_term_from_route_filter.Address.ValueStringPointer()
            config.Groups.Policy_options[i_policy_options].Policy_statement[i_policy_options_policy_statement].Term[i_policy_options_policy_statement_term].From[i_policy_options_policy_statement_term_from].Route_filter[i_policy_options_policy_statement_term_from_route_filter].Exact = v_policy_options_policy_statement_term_from_route_filter.Exact.ValueStringPointer()
            config.Groups.Policy_options[i_policy_options].Policy_statement[i_policy_options_policy_statement].Term[i_policy_options_policy_statement_term].From[i_policy_options_policy_statement_term_from].Route_filter[i_policy_options_policy_statement_term_from_route_filter].Orlonger = v_policy_options_policy_statement_term_from_route_filter.Orlonger.ValueStringPointer()
            config.Groups.Policy_options[i_policy_options].Policy_statement[i_policy_options_policy_statement].Term[i_policy_options_policy_statement_term].From[i_policy_options_policy_statement_term_from].Route_filter[i_policy_options_policy_statement_term_from_route_filter].Prefix_length_range = v_policy_options_policy_statement_term_from_route_filter.Prefix_length_range.ValueStringPointer()
        }
        }
            var var_policy_options_policy_statement_term_then []Policy_options_Policy_statement_Term_Then_Model
            resp.Diagnostics.Append(v_policy_options_policy_statement_term.Then.ElementsAs(ctx, &var_policy_options_policy_statement_term_then, false)...)
            if resp.Diagnostics.HasError() {
                return
            }
	    config.Groups.Policy_options[i_policy_options].Policy_statement[i_policy_options_policy_statement].Term[i_policy_options_policy_statement_term].Then = make([]xml_Policy_options_Policy_statement_Term_Then, len(var_policy_options_policy_statement_term_then))
        
		for i_policy_options_policy_statement_term_then, v_policy_options_policy_statement_term_then := range var_policy_options_policy_statement_term_then {
            var var_policy_options_policy_statement_term_then_community []Policy_options_Policy_statement_Term_Then_Community_Model
            resp.Diagnostics.Append(v_policy_options_policy_statement_term_then.Community.ElementsAs(ctx, &var_policy_options_policy_statement_term_then_community, false)...)
            if resp.Diagnostics.HasError() {
                return
            }
	    config.Groups.Policy_options[i_policy_options].Policy_statement[i_policy_options_policy_statement].Term[i_policy_options_policy_statement_term].Then[i_policy_options_policy_statement_term_then].Community = make([]xml_Policy_options_Policy_statement_Term_Then_Community, len(var_policy_options_policy_statement_term_then_community))
        
		for i_policy_options_policy_statement_term_then_community, v_policy_options_policy_statement_term_then_community := range var_policy_options_policy_statement_term_then_community {
            config.Groups.Policy_options[i_policy_options].Policy_statement[i_policy_options_policy_statement].Term[i_policy_options_policy_statement_term].Then[i_policy_options_policy_statement_term_then].Community[i_policy_options_policy_statement_term_then_community].Add = v_policy_options_policy_statement_term_then_community.Add.ValueStringPointer()
            config.Groups.Policy_options[i_policy_options].Policy_statement[i_policy_options_policy_statement].Term[i_policy_options_policy_statement_term].Then[i_policy_options_policy_statement_term_then].Community[i_policy_options_policy_statement_term_then_community].Community_name = v_policy_options_policy_statement_term_then_community.Community_name.ValueStringPointer()
        }
            config.Groups.Policy_options[i_policy_options].Policy_statement[i_policy_options_policy_statement].Term[i_policy_options_policy_statement_term].Then[i_policy_options_policy_statement_term_then].Accept = v_policy_options_policy_statement_term_then.Accept.ValueStringPointer()
            config.Groups.Policy_options[i_policy_options].Policy_statement[i_policy_options_policy_statement].Term[i_policy_options_policy_statement_term].Then[i_policy_options_policy_statement_term_then].Reject = v_policy_options_policy_statement_term_then.Reject.ValueStringPointer()
        }
        }
            var var_policy_options_policy_statement_then []Policy_options_Policy_statement_Then_Model
            resp.Diagnostics.Append(v_policy_options_policy_statement.Then.ElementsAs(ctx, &var_policy_options_policy_statement_then, false)...)
            if resp.Diagnostics.HasError() {
                return
            }
	    config.Groups.Policy_options[i_policy_options].Policy_statement[i_policy_options_policy_statement].Then = make([]xml_Policy_options_Policy_statement_Then, len(var_policy_options_policy_statement_then))
        
		for i_policy_options_policy_statement_then, v_policy_options_policy_statement_then := range var_policy_options_policy_statement_then {
            var var_policy_options_policy_statement_then_load_balance []Policy_options_Policy_statement_Then_Load_balance_Model
            resp.Diagnostics.Append(v_policy_options_policy_statement_then.Load_balance.ElementsAs(ctx, &var_policy_options_policy_statement_then_load_balance, false)...)
            if resp.Diagnostics.HasError() {
                return
            }
	    config.Groups.Policy_options[i_policy_options].Policy_statement[i_policy_options_policy_statement].Then[i_policy_options_policy_statement_then].Load_balance = make([]xml_Policy_options_Policy_statement_Then_Load_balance, len(var_policy_options_policy_statement_then_load_balance))
        
		for i_policy_options_policy_statement_then_load_balance, v_policy_options_policy_statement_then_load_balance := range var_policy_options_policy_statement_then_load_balance {
            config.Groups.Policy_options[i_policy_options].Policy_statement[i_policy_options_policy_statement].Then[i_policy_options_policy_statement_then].Load_balance[i_policy_options_policy_statement_then_load_balance].Per_packet = v_policy_options_policy_statement_then_load_balance.Per_packet.ValueStringPointer()
        }
        }
        }
        var var_policy_options_community []Policy_options_Community_Model
        resp.Diagnostics.Append(v_policy_options.Community.ElementsAs(ctx, &var_policy_options_community, false)...)
        if resp.Diagnostics.HasError() {
            return
        }
	    config.Groups.Policy_options[i_policy_options].Community = make([]xml_Policy_options_Community, len(var_policy_options_community))
        
		for i_policy_options_community, v_policy_options_community := range var_policy_options_community {
            config.Groups.Policy_options[i_policy_options].Community[i_policy_options_community].Name = v_policy_options_community.Name.ValueStringPointer()
			var var_policy_options_community_members []string
			resp.Diagnostics.Append(v_policy_options_community.Members.ElementsAs(ctx, &var_policy_options_community_members, false)...)
			for _, v_policy_options_community_members := range var_policy_options_community_members {
				config.Groups.Policy_options[i_policy_options].Community[i_policy_options_community].Members = append(config.Groups.Policy_options[i_policy_options].Community[i_policy_options_community].Members, &v_policy_options_community_members)
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
        var var_protocols_bgp []Protocols_Bgp_Model
        resp.Diagnostics.Append(v_protocols.Bgp.ElementsAs(ctx, &var_protocols_bgp, false)...)
        if resp.Diagnostics.HasError() {
            return
        }
	    config.Groups.Protocols[i_protocols].Bgp = make([]xml_Protocols_Bgp, len(var_protocols_bgp))
        
		for i_protocols_bgp, v_protocols_bgp := range var_protocols_bgp {
            var var_protocols_bgp_group []Protocols_Bgp_Group_Model
            resp.Diagnostics.Append(v_protocols_bgp.Group.ElementsAs(ctx, &var_protocols_bgp_group, false)...)
            if resp.Diagnostics.HasError() {
                return
            }
	    config.Groups.Protocols[i_protocols].Bgp[i_protocols_bgp].Group = make([]xml_Protocols_Bgp_Group, len(var_protocols_bgp_group))
        
		for i_protocols_bgp_group, v_protocols_bgp_group := range var_protocols_bgp_group {
            config.Groups.Protocols[i_protocols].Bgp[i_protocols_bgp].Group[i_protocols_bgp_group].Name = v_protocols_bgp_group.Name.ValueStringPointer()
            config.Groups.Protocols[i_protocols].Bgp[i_protocols_bgp].Group[i_protocols_bgp_group].Type = v_protocols_bgp_group.Type.ValueStringPointer()
            var var_protocols_bgp_group_multihop []Protocols_Bgp_Group_Multihop_Model
            resp.Diagnostics.Append(v_protocols_bgp_group.Multihop.ElementsAs(ctx, &var_protocols_bgp_group_multihop, false)...)
            if resp.Diagnostics.HasError() {
                return
            }
	    config.Groups.Protocols[i_protocols].Bgp[i_protocols_bgp].Group[i_protocols_bgp_group].Multihop = make([]xml_Protocols_Bgp_Group_Multihop, len(var_protocols_bgp_group_multihop))
        
		for i_protocols_bgp_group_multihop, v_protocols_bgp_group_multihop := range var_protocols_bgp_group_multihop {
            config.Groups.Protocols[i_protocols].Bgp[i_protocols_bgp].Group[i_protocols_bgp_group].Multihop[i_protocols_bgp_group_multihop].No_nexthop_change = v_protocols_bgp_group_multihop.No_nexthop_change.ValueStringPointer()
        }
            config.Groups.Protocols[i_protocols].Bgp[i_protocols_bgp].Group[i_protocols_bgp_group].Local_address = v_protocols_bgp_group.Local_address.ValueStringPointer()
            config.Groups.Protocols[i_protocols].Bgp[i_protocols_bgp].Group[i_protocols_bgp_group].Mtu_discovery = v_protocols_bgp_group.Mtu_discovery.ValueStringPointer()
			var var_protocols_bgp_group_import []string
			resp.Diagnostics.Append(v_protocols_bgp_group.Import.ElementsAs(ctx, &var_protocols_bgp_group_import, false)...)
			for _, v_protocols_bgp_group_import := range var_protocols_bgp_group_import {
				config.Groups.Protocols[i_protocols].Bgp[i_protocols_bgp].Group[i_protocols_bgp_group].Import = append(config.Groups.Protocols[i_protocols].Bgp[i_protocols_bgp].Group[i_protocols_bgp_group].Import, &v_protocols_bgp_group_import)
			}
            var var_protocols_bgp_group_family []Protocols_Bgp_Group_Family_Model
            resp.Diagnostics.Append(v_protocols_bgp_group.Family.ElementsAs(ctx, &var_protocols_bgp_group_family, false)...)
            if resp.Diagnostics.HasError() {
                return
            }
	    config.Groups.Protocols[i_protocols].Bgp[i_protocols_bgp].Group[i_protocols_bgp_group].Family = make([]xml_Protocols_Bgp_Group_Family, len(var_protocols_bgp_group_family))
        
		for i_protocols_bgp_group_family, v_protocols_bgp_group_family := range var_protocols_bgp_group_family {
            var var_protocols_bgp_group_family_evpn []Protocols_Bgp_Group_Family_Evpn_Model
            resp.Diagnostics.Append(v_protocols_bgp_group_family.Evpn.ElementsAs(ctx, &var_protocols_bgp_group_family_evpn, false)...)
            if resp.Diagnostics.HasError() {
                return
            }
	    config.Groups.Protocols[i_protocols].Bgp[i_protocols_bgp].Group[i_protocols_bgp_group].Family[i_protocols_bgp_group_family].Evpn = make([]xml_Protocols_Bgp_Group_Family_Evpn, len(var_protocols_bgp_group_family_evpn))
        
		for i_protocols_bgp_group_family_evpn, v_protocols_bgp_group_family_evpn := range var_protocols_bgp_group_family_evpn {
            var var_protocols_bgp_group_family_evpn_signaling []Protocols_Bgp_Group_Family_Evpn_Signaling_Model
            resp.Diagnostics.Append(v_protocols_bgp_group_family_evpn.Signaling.ElementsAs(ctx, &var_protocols_bgp_group_family_evpn_signaling, false)...)
            if resp.Diagnostics.HasError() {
                return
            }
	    config.Groups.Protocols[i_protocols].Bgp[i_protocols_bgp].Group[i_protocols_bgp_group].Family[i_protocols_bgp_group_family].Evpn[i_protocols_bgp_group_family_evpn].Signaling = make([]xml_Protocols_Bgp_Group_Family_Evpn_Signaling, len(var_protocols_bgp_group_family_evpn_signaling))
        
		for i_protocols_bgp_group_family_evpn_signaling, v_protocols_bgp_group_family_evpn_signaling := range var_protocols_bgp_group_family_evpn_signaling {
            var var_protocols_bgp_group_family_evpn_signaling_delay_route_advertisements []Protocols_Bgp_Group_Family_Evpn_Signaling_Delay_route_advertisements_Model
            resp.Diagnostics.Append(v_protocols_bgp_group_family_evpn_signaling.Delay_route_advertisements.ElementsAs(ctx, &var_protocols_bgp_group_family_evpn_signaling_delay_route_advertisements, false)...)
            if resp.Diagnostics.HasError() {
                return
            }
	    config.Groups.Protocols[i_protocols].Bgp[i_protocols_bgp].Group[i_protocols_bgp_group].Family[i_protocols_bgp_group_family].Evpn[i_protocols_bgp_group_family_evpn].Signaling[i_protocols_bgp_group_family_evpn_signaling].Delay_route_advertisements = make([]xml_Protocols_Bgp_Group_Family_Evpn_Signaling_Delay_route_advertisements, len(var_protocols_bgp_group_family_evpn_signaling_delay_route_advertisements))
        
		for i_protocols_bgp_group_family_evpn_signaling_delay_route_advertisements, v_protocols_bgp_group_family_evpn_signaling_delay_route_advertisements := range var_protocols_bgp_group_family_evpn_signaling_delay_route_advertisements {
            var var_protocols_bgp_group_family_evpn_signaling_delay_route_advertisements_minimum_delay []Protocols_Bgp_Group_Family_Evpn_Signaling_Delay_route_advertisements_Minimum_delay_Model
            resp.Diagnostics.Append(v_protocols_bgp_group_family_evpn_signaling_delay_route_advertisements.Minimum_delay.ElementsAs(ctx, &var_protocols_bgp_group_family_evpn_signaling_delay_route_advertisements_minimum_delay, false)...)
            if resp.Diagnostics.HasError() {
                return
            }
	    config.Groups.Protocols[i_protocols].Bgp[i_protocols_bgp].Group[i_protocols_bgp_group].Family[i_protocols_bgp_group_family].Evpn[i_protocols_bgp_group_family_evpn].Signaling[i_protocols_bgp_group_family_evpn_signaling].Delay_route_advertisements[i_protocols_bgp_group_family_evpn_signaling_delay_route_advertisements].Minimum_delay = make([]xml_Protocols_Bgp_Group_Family_Evpn_Signaling_Delay_route_advertisements_Minimum_delay, len(var_protocols_bgp_group_family_evpn_signaling_delay_route_advertisements_minimum_delay))
        
		for i_protocols_bgp_group_family_evpn_signaling_delay_route_advertisements_minimum_delay, v_protocols_bgp_group_family_evpn_signaling_delay_route_advertisements_minimum_delay := range var_protocols_bgp_group_family_evpn_signaling_delay_route_advertisements_minimum_delay {
            config.Groups.Protocols[i_protocols].Bgp[i_protocols_bgp].Group[i_protocols_bgp_group].Family[i_protocols_bgp_group_family].Evpn[i_protocols_bgp_group_family_evpn].Signaling[i_protocols_bgp_group_family_evpn_signaling].Delay_route_advertisements[i_protocols_bgp_group_family_evpn_signaling_delay_route_advertisements].Minimum_delay[i_protocols_bgp_group_family_evpn_signaling_delay_route_advertisements_minimum_delay].Routing_uptime = v_protocols_bgp_group_family_evpn_signaling_delay_route_advertisements_minimum_delay.Routing_uptime.ValueStringPointer()
        }
        }
        }
        }
        }
			var var_protocols_bgp_group_export []string
			resp.Diagnostics.Append(v_protocols_bgp_group.Export.ElementsAs(ctx, &var_protocols_bgp_group_export, false)...)
			for _, v_protocols_bgp_group_export := range var_protocols_bgp_group_export {
				config.Groups.Protocols[i_protocols].Bgp[i_protocols_bgp].Group[i_protocols_bgp_group].Export = append(config.Groups.Protocols[i_protocols].Bgp[i_protocols_bgp].Group[i_protocols_bgp_group].Export, &v_protocols_bgp_group_export)
			}
            config.Groups.Protocols[i_protocols].Bgp[i_protocols_bgp].Group[i_protocols_bgp_group].Vpn_apply_export = v_protocols_bgp_group.Vpn_apply_export.ValueStringPointer()
            config.Groups.Protocols[i_protocols].Bgp[i_protocols_bgp].Group[i_protocols_bgp_group].Cluster = v_protocols_bgp_group.Cluster.ValueStringPointer()
            var var_protocols_bgp_group_local_as []Protocols_Bgp_Group_Local_as_Model
            resp.Diagnostics.Append(v_protocols_bgp_group.Local_as.ElementsAs(ctx, &var_protocols_bgp_group_local_as, false)...)
            if resp.Diagnostics.HasError() {
                return
            }
	    config.Groups.Protocols[i_protocols].Bgp[i_protocols_bgp].Group[i_protocols_bgp_group].Local_as = make([]xml_Protocols_Bgp_Group_Local_as, len(var_protocols_bgp_group_local_as))
        
		for i_protocols_bgp_group_local_as, v_protocols_bgp_group_local_as := range var_protocols_bgp_group_local_as {
            config.Groups.Protocols[i_protocols].Bgp[i_protocols_bgp].Group[i_protocols_bgp_group].Local_as[i_protocols_bgp_group_local_as].As_number = v_protocols_bgp_group_local_as.As_number.ValueStringPointer()
        }
            var var_protocols_bgp_group_multipath []Protocols_Bgp_Group_Multipath_Model
            resp.Diagnostics.Append(v_protocols_bgp_group.Multipath.ElementsAs(ctx, &var_protocols_bgp_group_multipath, false)...)
            if resp.Diagnostics.HasError() {
                return
            }
	    config.Groups.Protocols[i_protocols].Bgp[i_protocols_bgp].Group[i_protocols_bgp_group].Multipath = make([]xml_Protocols_Bgp_Group_Multipath, len(var_protocols_bgp_group_multipath))
        
		for i_protocols_bgp_group_multipath, v_protocols_bgp_group_multipath := range var_protocols_bgp_group_multipath {
            config.Groups.Protocols[i_protocols].Bgp[i_protocols_bgp].Group[i_protocols_bgp_group].Multipath[i_protocols_bgp_group_multipath].Multiple_as = v_protocols_bgp_group_multipath.Multiple_as.ValueStringPointer()
        }
            var var_protocols_bgp_group_bfd_liveness_detection []Protocols_Bgp_Group_Bfd_liveness_detection_Model
            resp.Diagnostics.Append(v_protocols_bgp_group.Bfd_liveness_detection.ElementsAs(ctx, &var_protocols_bgp_group_bfd_liveness_detection, false)...)
            if resp.Diagnostics.HasError() {
                return
            }
	    config.Groups.Protocols[i_protocols].Bgp[i_protocols_bgp].Group[i_protocols_bgp_group].Bfd_liveness_detection = make([]xml_Protocols_Bgp_Group_Bfd_liveness_detection, len(var_protocols_bgp_group_bfd_liveness_detection))
        
		for i_protocols_bgp_group_bfd_liveness_detection, v_protocols_bgp_group_bfd_liveness_detection := range var_protocols_bgp_group_bfd_liveness_detection {
            config.Groups.Protocols[i_protocols].Bgp[i_protocols_bgp].Group[i_protocols_bgp_group].Bfd_liveness_detection[i_protocols_bgp_group_bfd_liveness_detection].Minimum_interval = v_protocols_bgp_group_bfd_liveness_detection.Minimum_interval.ValueStringPointer()
            config.Groups.Protocols[i_protocols].Bgp[i_protocols_bgp].Group[i_protocols_bgp_group].Bfd_liveness_detection[i_protocols_bgp_group_bfd_liveness_detection].Multiplier = v_protocols_bgp_group_bfd_liveness_detection.Multiplier.ValueStringPointer()
        }
			var var_protocols_bgp_group_allow []string
			resp.Diagnostics.Append(v_protocols_bgp_group.Allow.ElementsAs(ctx, &var_protocols_bgp_group_allow, false)...)
			for _, v_protocols_bgp_group_allow := range var_protocols_bgp_group_allow {
				config.Groups.Protocols[i_protocols].Bgp[i_protocols_bgp].Group[i_protocols_bgp_group].Allow = append(config.Groups.Protocols[i_protocols].Bgp[i_protocols_bgp].Group[i_protocols_bgp_group].Allow, &v_protocols_bgp_group_allow)
			}
            var var_protocols_bgp_group_neighbor []Protocols_Bgp_Group_Neighbor_Model
            resp.Diagnostics.Append(v_protocols_bgp_group.Neighbor.ElementsAs(ctx, &var_protocols_bgp_group_neighbor, false)...)
            if resp.Diagnostics.HasError() {
                return
            }
	    config.Groups.Protocols[i_protocols].Bgp[i_protocols_bgp].Group[i_protocols_bgp_group].Neighbor = make([]xml_Protocols_Bgp_Group_Neighbor, len(var_protocols_bgp_group_neighbor))
        
		for i_protocols_bgp_group_neighbor, v_protocols_bgp_group_neighbor := range var_protocols_bgp_group_neighbor {
            config.Groups.Protocols[i_protocols].Bgp[i_protocols_bgp].Group[i_protocols_bgp_group].Neighbor[i_protocols_bgp_group_neighbor].Name = v_protocols_bgp_group_neighbor.Name.ValueStringPointer()
            config.Groups.Protocols[i_protocols].Bgp[i_protocols_bgp].Group[i_protocols_bgp_group].Neighbor[i_protocols_bgp_group_neighbor].Description = v_protocols_bgp_group_neighbor.Description.ValueStringPointer()
            config.Groups.Protocols[i_protocols].Bgp[i_protocols_bgp].Group[i_protocols_bgp_group].Neighbor[i_protocols_bgp_group_neighbor].Peer_as = v_protocols_bgp_group_neighbor.Peer_as.ValueStringPointer()
        }
        }
        }
        var var_protocols_evpn []Protocols_Evpn_Model
        resp.Diagnostics.Append(v_protocols.Evpn.ElementsAs(ctx, &var_protocols_evpn, false)...)
        if resp.Diagnostics.HasError() {
            return
        }
	    config.Groups.Protocols[i_protocols].Evpn = make([]xml_Protocols_Evpn, len(var_protocols_evpn))
        
		for i_protocols_evpn, v_protocols_evpn := range var_protocols_evpn {
            config.Groups.Protocols[i_protocols].Evpn[i_protocols_evpn].Encapsulation = v_protocols_evpn.Encapsulation.ValueStringPointer()
            config.Groups.Protocols[i_protocols].Evpn[i_protocols_evpn].Multicast_mode = v_protocols_evpn.Multicast_mode.ValueStringPointer()
            config.Groups.Protocols[i_protocols].Evpn[i_protocols_evpn].Default_gateway = v_protocols_evpn.Default_gateway.ValueStringPointer()
			var var_protocols_evpn_extended_vni_list []string
			resp.Diagnostics.Append(v_protocols_evpn.Extended_vni_list.ElementsAs(ctx, &var_protocols_evpn_extended_vni_list, false)...)
			for _, v_protocols_evpn_extended_vni_list := range var_protocols_evpn_extended_vni_list {
				config.Groups.Protocols[i_protocols].Evpn[i_protocols_evpn].Extended_vni_list = append(config.Groups.Protocols[i_protocols].Evpn[i_protocols_evpn].Extended_vni_list, &v_protocols_evpn_extended_vni_list)
			}
            config.Groups.Protocols[i_protocols].Evpn[i_protocols_evpn].No_core_isolation = v_protocols_evpn.No_core_isolation.ValueStringPointer()
        }
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
        var var_protocols_igmp_snooping []Protocols_Igmp_snooping_Model
        resp.Diagnostics.Append(v_protocols.Igmp_snooping.ElementsAs(ctx, &var_protocols_igmp_snooping, false)...)
        if resp.Diagnostics.HasError() {
            return
        }
	    config.Groups.Protocols[i_protocols].Igmp_snooping = make([]xml_Protocols_Igmp_snooping, len(var_protocols_igmp_snooping))
        
		for i_protocols_igmp_snooping, v_protocols_igmp_snooping := range var_protocols_igmp_snooping {
            var var_protocols_igmp_snooping_vlan []Protocols_Igmp_snooping_Vlan_Model
            resp.Diagnostics.Append(v_protocols_igmp_snooping.Vlan.ElementsAs(ctx, &var_protocols_igmp_snooping_vlan, false)...)
            if resp.Diagnostics.HasError() {
                return
            }
	    config.Groups.Protocols[i_protocols].Igmp_snooping[i_protocols_igmp_snooping].Vlan = make([]xml_Protocols_Igmp_snooping_Vlan, len(var_protocols_igmp_snooping_vlan))
        
		for i_protocols_igmp_snooping_vlan, v_protocols_igmp_snooping_vlan := range var_protocols_igmp_snooping_vlan {
            config.Groups.Protocols[i_protocols].Igmp_snooping[i_protocols_igmp_snooping].Vlan[i_protocols_igmp_snooping_vlan].Name = v_protocols_igmp_snooping_vlan.Name.ValueStringPointer()
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
            var var_routing_instances_instance_route_distinguisher []Routing_instances_Instance_Route_distinguisher_Model
            resp.Diagnostics.Append(v_routing_instances_instance.Route_distinguisher.ElementsAs(ctx, &var_routing_instances_instance_route_distinguisher, false)...)
            if resp.Diagnostics.HasError() {
                return
            }
	    config.Groups.Routing_instances[i_routing_instances].Instance[i_routing_instances_instance].Route_distinguisher = make([]xml_Routing_instances_Instance_Route_distinguisher, len(var_routing_instances_instance_route_distinguisher))
        
		for i_routing_instances_instance_route_distinguisher, v_routing_instances_instance_route_distinguisher := range var_routing_instances_instance_route_distinguisher {
            config.Groups.Routing_instances[i_routing_instances].Instance[i_routing_instances_instance].Route_distinguisher[i_routing_instances_instance_route_distinguisher].Rd_type = v_routing_instances_instance_route_distinguisher.Rd_type.ValueStringPointer()
        }
            var var_routing_instances_instance_vrf_target []Routing_instances_Instance_Vrf_target_Model
            resp.Diagnostics.Append(v_routing_instances_instance.Vrf_target.ElementsAs(ctx, &var_routing_instances_instance_vrf_target, false)...)
            if resp.Diagnostics.HasError() {
                return
            }
	    config.Groups.Routing_instances[i_routing_instances].Instance[i_routing_instances_instance].Vrf_target = make([]xml_Routing_instances_Instance_Vrf_target, len(var_routing_instances_instance_vrf_target))
        
		for i_routing_instances_instance_vrf_target, v_routing_instances_instance_vrf_target := range var_routing_instances_instance_vrf_target {
            config.Groups.Routing_instances[i_routing_instances].Instance[i_routing_instances_instance].Vrf_target[i_routing_instances_instance_vrf_target].Community = v_routing_instances_instance_vrf_target.Community.ValueStringPointer()
        }
            var var_routing_instances_instance_vrf_table_label []Routing_instances_Instance_Vrf_table_label_Model
            resp.Diagnostics.Append(v_routing_instances_instance.Vrf_table_label.ElementsAs(ctx, &var_routing_instances_instance_vrf_table_label, false)...)
            if resp.Diagnostics.HasError() {
                return
            }
	    config.Groups.Routing_instances[i_routing_instances].Instance[i_routing_instances_instance].Vrf_table_label = make([]xml_Routing_instances_Instance_Vrf_table_label, len(var_routing_instances_instance_vrf_table_label))
        
            var var_routing_instances_instance_routing_options []Routing_instances_Instance_Routing_options_Model
            resp.Diagnostics.Append(v_routing_instances_instance.Routing_options.ElementsAs(ctx, &var_routing_instances_instance_routing_options, false)...)
            if resp.Diagnostics.HasError() {
                return
            }
	    config.Groups.Routing_instances[i_routing_instances].Instance[i_routing_instances_instance].Routing_options = make([]xml_Routing_instances_Instance_Routing_options, len(var_routing_instances_instance_routing_options))
        
		for i_routing_instances_instance_routing_options, v_routing_instances_instance_routing_options := range var_routing_instances_instance_routing_options {
            var var_routing_instances_instance_routing_options_auto_export []Routing_instances_Instance_Routing_options_Auto_export_Model
            resp.Diagnostics.Append(v_routing_instances_instance_routing_options.Auto_export.ElementsAs(ctx, &var_routing_instances_instance_routing_options_auto_export, false)...)
            if resp.Diagnostics.HasError() {
                return
            }
	    config.Groups.Routing_instances[i_routing_instances].Instance[i_routing_instances_instance].Routing_options[i_routing_instances_instance_routing_options].Auto_export = make([]xml_Routing_instances_Instance_Routing_options_Auto_export, len(var_routing_instances_instance_routing_options_auto_export))
        
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
            config.Groups.Routing_instances[i_routing_instances].Instance[i_routing_instances_instance].Protocols[i_routing_instances_instance_protocols].Ospf[i_routing_instances_instance_protocols_ospf].Area[i_routing_instances_instance_protocols_ospf_area].Interface[i_routing_instances_instance_protocols_ospf_area_interface].Metric = v_routing_instances_instance_protocols_ospf_area_interface.Metric.ValueStringPointer()
        }
        }
        }
            var var_routing_instances_instance_protocols_evpn []Routing_instances_Instance_Protocols_Evpn_Model
            resp.Diagnostics.Append(v_routing_instances_instance_protocols.Evpn.ElementsAs(ctx, &var_routing_instances_instance_protocols_evpn, false)...)
            if resp.Diagnostics.HasError() {
                return
            }
	    config.Groups.Routing_instances[i_routing_instances].Instance[i_routing_instances_instance].Protocols[i_routing_instances_instance_protocols].Evpn = make([]xml_Routing_instances_Instance_Protocols_Evpn, len(var_routing_instances_instance_protocols_evpn))
        
		for i_routing_instances_instance_protocols_evpn, v_routing_instances_instance_protocols_evpn := range var_routing_instances_instance_protocols_evpn {
            var var_routing_instances_instance_protocols_evpn_ip_prefix_routes []Routing_instances_Instance_Protocols_Evpn_Ip_prefix_routes_Model
            resp.Diagnostics.Append(v_routing_instances_instance_protocols_evpn.Ip_prefix_routes.ElementsAs(ctx, &var_routing_instances_instance_protocols_evpn_ip_prefix_routes, false)...)
            if resp.Diagnostics.HasError() {
                return
            }
	    config.Groups.Routing_instances[i_routing_instances].Instance[i_routing_instances_instance].Protocols[i_routing_instances_instance_protocols].Evpn[i_routing_instances_instance_protocols_evpn].Ip_prefix_routes = make([]xml_Routing_instances_Instance_Protocols_Evpn_Ip_prefix_routes, len(var_routing_instances_instance_protocols_evpn_ip_prefix_routes))
        
		for i_routing_instances_instance_protocols_evpn_ip_prefix_routes, v_routing_instances_instance_protocols_evpn_ip_prefix_routes := range var_routing_instances_instance_protocols_evpn_ip_prefix_routes {
            config.Groups.Routing_instances[i_routing_instances].Instance[i_routing_instances_instance].Protocols[i_routing_instances_instance_protocols].Evpn[i_routing_instances_instance_protocols_evpn].Ip_prefix_routes[i_routing_instances_instance_protocols_evpn_ip_prefix_routes].Advertise = v_routing_instances_instance_protocols_evpn_ip_prefix_routes.Advertise.ValueStringPointer()
            config.Groups.Routing_instances[i_routing_instances].Instance[i_routing_instances_instance].Protocols[i_routing_instances_instance_protocols].Evpn[i_routing_instances_instance_protocols_evpn].Ip_prefix_routes[i_routing_instances_instance_protocols_evpn_ip_prefix_routes].Encapsulation = v_routing_instances_instance_protocols_evpn_ip_prefix_routes.Encapsulation.ValueStringPointer()
            config.Groups.Routing_instances[i_routing_instances].Instance[i_routing_instances_instance].Protocols[i_routing_instances_instance_protocols].Evpn[i_routing_instances_instance_protocols_evpn].Ip_prefix_routes[i_routing_instances_instance_protocols_evpn_ip_prefix_routes].Vni = v_routing_instances_instance_protocols_evpn_ip_prefix_routes.Vni.ValueStringPointer()
			var var_routing_instances_instance_protocols_evpn_ip_prefix_routes_export []string
			resp.Diagnostics.Append(v_routing_instances_instance_protocols_evpn_ip_prefix_routes.Export.ElementsAs(ctx, &var_routing_instances_instance_protocols_evpn_ip_prefix_routes_export, false)...)
			for _, v_routing_instances_instance_protocols_evpn_ip_prefix_routes_export := range var_routing_instances_instance_protocols_evpn_ip_prefix_routes_export {
				config.Groups.Routing_instances[i_routing_instances].Instance[i_routing_instances_instance].Protocols[i_routing_instances_instance_protocols].Evpn[i_routing_instances_instance_protocols_evpn].Ip_prefix_routes[i_routing_instances_instance_protocols_evpn_ip_prefix_routes].Export = append(config.Groups.Routing_instances[i_routing_instances].Instance[i_routing_instances_instance].Protocols[i_routing_instances_instance_protocols].Evpn[i_routing_instances_instance_protocols_evpn].Ip_prefix_routes[i_routing_instances_instance_protocols_evpn_ip_prefix_routes].Export, &v_routing_instances_instance_protocols_evpn_ip_prefix_routes_export)
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
        config.Groups.Routing_options[i_routing_options].Router_id = v_routing_options.Router_id.ValueStringPointer()
        var var_routing_options_forwarding_table []Routing_options_Forwarding_table_Model
        resp.Diagnostics.Append(v_routing_options.Forwarding_table.ElementsAs(ctx, &var_routing_options_forwarding_table, false)...)
        if resp.Diagnostics.HasError() {
            return
        }
	    config.Groups.Routing_options[i_routing_options].Forwarding_table = make([]xml_Routing_options_Forwarding_table, len(var_routing_options_forwarding_table))
        
		for i_routing_options_forwarding_table, v_routing_options_forwarding_table := range var_routing_options_forwarding_table {
			var var_routing_options_forwarding_table_export []string
			resp.Diagnostics.Append(v_routing_options_forwarding_table.Export.ElementsAs(ctx, &var_routing_options_forwarding_table_export, false)...)
			for _, v_routing_options_forwarding_table_export := range var_routing_options_forwarding_table_export {
				config.Groups.Routing_options[i_routing_options].Forwarding_table[i_routing_options_forwarding_table].Export = append(config.Groups.Routing_options[i_routing_options].Forwarding_table[i_routing_options_forwarding_table].Export, &v_routing_options_forwarding_table_export)
			}
            config.Groups.Routing_options[i_routing_options].Forwarding_table[i_routing_options_forwarding_table].Ecmp_fast_reroute = v_routing_options_forwarding_table.Ecmp_fast_reroute.ValueStringPointer()
            var var_routing_options_forwarding_table_chained_composite_next_hop []Routing_options_Forwarding_table_Chained_composite_next_hop_Model
            resp.Diagnostics.Append(v_routing_options_forwarding_table.Chained_composite_next_hop.ElementsAs(ctx, &var_routing_options_forwarding_table_chained_composite_next_hop, false)...)
            if resp.Diagnostics.HasError() {
                return
            }
	    config.Groups.Routing_options[i_routing_options].Forwarding_table[i_routing_options_forwarding_table].Chained_composite_next_hop = make([]xml_Routing_options_Forwarding_table_Chained_composite_next_hop, len(var_routing_options_forwarding_table_chained_composite_next_hop))
        
		for i_routing_options_forwarding_table_chained_composite_next_hop, v_routing_options_forwarding_table_chained_composite_next_hop := range var_routing_options_forwarding_table_chained_composite_next_hop {
            var var_routing_options_forwarding_table_chained_composite_next_hop_ingress []Routing_options_Forwarding_table_Chained_composite_next_hop_Ingress_Model
            resp.Diagnostics.Append(v_routing_options_forwarding_table_chained_composite_next_hop.Ingress.ElementsAs(ctx, &var_routing_options_forwarding_table_chained_composite_next_hop_ingress, false)...)
            if resp.Diagnostics.HasError() {
                return
            }
	    config.Groups.Routing_options[i_routing_options].Forwarding_table[i_routing_options_forwarding_table].Chained_composite_next_hop[i_routing_options_forwarding_table_chained_composite_next_hop].Ingress = make([]xml_Routing_options_Forwarding_table_Chained_composite_next_hop_Ingress, len(var_routing_options_forwarding_table_chained_composite_next_hop_ingress))
        
		for i_routing_options_forwarding_table_chained_composite_next_hop_ingress, v_routing_options_forwarding_table_chained_composite_next_hop_ingress := range var_routing_options_forwarding_table_chained_composite_next_hop_ingress {
            config.Groups.Routing_options[i_routing_options].Forwarding_table[i_routing_options_forwarding_table].Chained_composite_next_hop[i_routing_options_forwarding_table_chained_composite_next_hop].Ingress[i_routing_options_forwarding_table_chained_composite_next_hop_ingress].Evpn = v_routing_options_forwarding_table_chained_composite_next_hop_ingress.Evpn.ValueStringPointer()
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
	
    var var_switch_options []Switch_options_Model
    if plan.Switch_options.IsNull() {
        var_switch_options = []Switch_options_Model{}
    }else {
        resp.Diagnostics.Append(plan.Switch_options.ElementsAs(ctx, &var_switch_options, false)...)
        if resp.Diagnostics.HasError() {
            return
        }
    }
    config.Groups.Switch_options = make([]xml_Switch_options, len(var_switch_options))
    
    for i_switch_options, v_switch_options := range var_switch_options {
        var var_switch_options_vtep_source_interface []Switch_options_Vtep_source_interface_Model
        resp.Diagnostics.Append(v_switch_options.Vtep_source_interface.ElementsAs(ctx, &var_switch_options_vtep_source_interface, false)...)
        if resp.Diagnostics.HasError() {
            return
        }
	    config.Groups.Switch_options[i_switch_options].Vtep_source_interface = make([]xml_Switch_options_Vtep_source_interface, len(var_switch_options_vtep_source_interface))
        
		for i_switch_options_vtep_source_interface, v_switch_options_vtep_source_interface := range var_switch_options_vtep_source_interface {
            config.Groups.Switch_options[i_switch_options].Vtep_source_interface[i_switch_options_vtep_source_interface].Interface_name = v_switch_options_vtep_source_interface.Interface_name.ValueStringPointer()
        }
        var var_switch_options_route_distinguisher []Switch_options_Route_distinguisher_Model
        resp.Diagnostics.Append(v_switch_options.Route_distinguisher.ElementsAs(ctx, &var_switch_options_route_distinguisher, false)...)
        if resp.Diagnostics.HasError() {
            return
        }
	    config.Groups.Switch_options[i_switch_options].Route_distinguisher = make([]xml_Switch_options_Route_distinguisher, len(var_switch_options_route_distinguisher))
        
		for i_switch_options_route_distinguisher, v_switch_options_route_distinguisher := range var_switch_options_route_distinguisher {
            config.Groups.Switch_options[i_switch_options].Route_distinguisher[i_switch_options_route_distinguisher].Rd_type = v_switch_options_route_distinguisher.Rd_type.ValueStringPointer()
        }
        var var_switch_options_vrf_target []Switch_options_Vrf_target_Model
        resp.Diagnostics.Append(v_switch_options.Vrf_target.ElementsAs(ctx, &var_switch_options_vrf_target, false)...)
        if resp.Diagnostics.HasError() {
            return
        }
	    config.Groups.Switch_options[i_switch_options].Vrf_target = make([]xml_Switch_options_Vrf_target, len(var_switch_options_vrf_target))
        
		for i_switch_options_vrf_target, v_switch_options_vrf_target := range var_switch_options_vrf_target {
            config.Groups.Switch_options[i_switch_options].Vrf_target[i_switch_options_vrf_target].Community = v_switch_options_vrf_target.Community.ValueStringPointer()
            var var_switch_options_vrf_target_auto []Switch_options_Vrf_target_Auto_Model
            resp.Diagnostics.Append(v_switch_options_vrf_target.Auto.ElementsAs(ctx, &var_switch_options_vrf_target_auto, false)...)
            if resp.Diagnostics.HasError() {
                return
            }
	    config.Groups.Switch_options[i_switch_options].Vrf_target[i_switch_options_vrf_target].Auto = make([]xml_Switch_options_Vrf_target_Auto, len(var_switch_options_vrf_target_auto))
        
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
            config.Groups.System[i_system].Syslog[i_system_syslog].File[i_system_syslog_file].Contents[i_system_syslog_file_contents].Notice = v_system_syslog_file_contents.Notice.ValueStringPointer()
            config.Groups.System[i_system].Syslog[i_system_syslog].File[i_system_syslog_file].Contents[i_system_syslog_file_contents].Info = v_system_syslog_file_contents.Info.ValueStringPointer()
        }
        }
        }
        var var_system_extensions []System_Extensions_Model
        resp.Diagnostics.Append(v_system.Extensions.ElementsAs(ctx, &var_system_extensions, false)...)
        if resp.Diagnostics.HasError() {
            return
        }
	    config.Groups.System[i_system].Extensions = make([]xml_System_Extensions, len(var_system_extensions))
        
		for i_system_extensions, v_system_extensions := range var_system_extensions {
            var var_system_extensions_providers []System_Extensions_Providers_Model
            resp.Diagnostics.Append(v_system_extensions.Providers.ElementsAs(ctx, &var_system_extensions_providers, false)...)
            if resp.Diagnostics.HasError() {
                return
            }
	    config.Groups.System[i_system].Extensions[i_system_extensions].Providers = make([]xml_System_Extensions_Providers, len(var_system_extensions_providers))
        
		for i_system_extensions_providers, v_system_extensions_providers := range var_system_extensions_providers {
            config.Groups.System[i_system].Extensions[i_system_extensions].Providers[i_system_extensions_providers].Name = v_system_extensions_providers.Name.ValueStringPointer()
            var var_system_extensions_providers_license_type []System_Extensions_Providers_License_type_Model
            resp.Diagnostics.Append(v_system_extensions_providers.License_type.ElementsAs(ctx, &var_system_extensions_providers_license_type, false)...)
            if resp.Diagnostics.HasError() {
                return
            }
	    config.Groups.System[i_system].Extensions[i_system_extensions].Providers[i_system_extensions_providers].License_type = make([]xml_System_Extensions_Providers_License_type, len(var_system_extensions_providers_license_type))
        
		for i_system_extensions_providers_license_type, v_system_extensions_providers_license_type := range var_system_extensions_providers_license_type {
            config.Groups.System[i_system].Extensions[i_system_extensions].Providers[i_system_extensions_providers].License_type[i_system_extensions_providers_license_type].Name = v_system_extensions_providers_license_type.Name.ValueStringPointer()
			var var_system_extensions_providers_license_type_deployment_scope []string
			resp.Diagnostics.Append(v_system_extensions_providers_license_type.Deployment_scope.ElementsAs(ctx, &var_system_extensions_providers_license_type_deployment_scope, false)...)
			for _, v_system_extensions_providers_license_type_deployment_scope := range var_system_extensions_providers_license_type_deployment_scope {
				config.Groups.System[i_system].Extensions[i_system_extensions].Providers[i_system_extensions_providers].License_type[i_system_extensions_providers_license_type].Deployment_scope = append(config.Groups.System[i_system].Extensions[i_system_extensions].Providers[i_system_extensions_providers].License_type[i_system_extensions_providers_license_type].Deployment_scope, &v_system_extensions_providers_license_type_deployment_scope)
			}
        }
        }
        }
    }
	
    var var_vlans []Vlans_Model
    if plan.Vlans.IsNull() {
        var_vlans = []Vlans_Model{}
    }else {
        resp.Diagnostics.Append(plan.Vlans.ElementsAs(ctx, &var_vlans, false)...)
        if resp.Diagnostics.HasError() {
            return
        }
    }
    config.Groups.Vlans = make([]xml_Vlans, len(var_vlans))
    
    for i_vlans, v_vlans := range var_vlans {
        var var_vlans_vlan []Vlans_Vlan_Model
        resp.Diagnostics.Append(v_vlans.Vlan.ElementsAs(ctx, &var_vlans_vlan, false)...)
        if resp.Diagnostics.HasError() {
            return
        }
	    config.Groups.Vlans[i_vlans].Vlan = make([]xml_Vlans_Vlan, len(var_vlans_vlan))
        
		for i_vlans_vlan, v_vlans_vlan := range var_vlans_vlan {
            config.Groups.Vlans[i_vlans].Vlan[i_vlans_vlan].Name = v_vlans_vlan.Name.ValueStringPointer()
            config.Groups.Vlans[i_vlans].Vlan[i_vlans_vlan].Vlan_id = v_vlans_vlan.Vlan_id.ValueStringPointer()
            config.Groups.Vlans[i_vlans].Vlan[i_vlans_vlan].L3_interface = v_vlans_vlan.L3_interface.ValueStringPointer()
            var var_vlans_vlan_vxlan []Vlans_Vlan_Vxlan_Model
            resp.Diagnostics.Append(v_vlans_vlan.Vxlan.ElementsAs(ctx, &var_vlans_vlan_vxlan, false)...)
            if resp.Diagnostics.HasError() {
                return
            }
	    config.Groups.Vlans[i_vlans].Vlan[i_vlans_vlan].Vxlan = make([]xml_Vlans_Vlan_Vxlan, len(var_vlans_vlan_vxlan))
        
		for i_vlans_vlan_vxlan, v_vlans_vlan_vxlan := range var_vlans_vlan_vxlan {
            config.Groups.Vlans[i_vlans].Vlan[i_vlans_vlan].Vxlan[i_vlans_vlan_vxlan].Vni = v_vlans_vlan_vxlan.Vni.ValueStringPointer()
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
