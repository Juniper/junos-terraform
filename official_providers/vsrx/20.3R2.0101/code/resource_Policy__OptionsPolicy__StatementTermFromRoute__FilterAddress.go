
package main

import (
    "encoding/xml"
    "fmt"
    "github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)


// v_ is appended before every variable so it doesn't give any conflict
// with any keyword in golang. ex - interface is keyword in golang
type xmlPolicy__OptionsPolicy__StatementTermFromRoute__FilterAddress struct {
	XMLName xml.Name `xml:"configuration"`
	Groups  struct {
		XMLName	xml.Name	`xml:"groups"`
		Name	string	`xml:"name"`
		V_policy__statement  struct {
			XMLName xml.Name `xml:"policy-statement"`
			V_name  *string  `xml:"name,omitempty"`
			V_term  struct {
				XMLName xml.Name `xml:"term"`
				V_name__1  *string  `xml:"name,omitempty"`
				V_route__filter  struct {
					XMLName xml.Name `xml:"route-filter"`
					V_damping  *string  `xml:"damping,omitempty"`
					V_no__entropy__label__capability  *string  `xml:"no-entropy-label-capability,omitempty"`
					V_metric3  *string  `xml:"metric3,omitempty"`
					V_tag2  *string  `xml:"tag2,omitempty"`
					V_add__path  *string  `xml:"add-path,omitempty"`
					V_tunnel__end__point__address  *string  `xml:"tunnel-end-point-address,omitempty"`
					V_default__action  *string  `xml:"default-action,omitempty"`
					V_trace  *string  `xml:"trace,omitempty"`
					V_install__to__fib  *string  `xml:"install-to-fib,omitempty"`
					V_ssm__source  *string  `xml:"ssm-source,omitempty"`
					V_metric2  *string  `xml:"metric2,omitempty"`
					V_tag  *string  `xml:"tag,omitempty"`
					V_longer  *string  `xml:"longer,omitempty"`
					V_orlonger  *string  `xml:"orlonger,omitempty"`
					V_address__mask  *string  `xml:"address-mask,omitempty"`
					V_community  *string  `xml:"community,omitempty"`
					V_tunnel__attribute  *string  `xml:"tunnel-attribute,omitempty"`
					V_as__path__prepend  *string  `xml:"as-path-prepend,omitempty"`
					V_map__to__interface  *string  `xml:"map-to-interface,omitempty"`
					V_accept_reject  *string  `xml:"accept_reject,omitempty"`
					V_apply__advanced  *string  `xml:"apply-advanced,omitempty"`
					V_aigp__adjust  *string  `xml:"aigp-adjust,omitempty"`
					V_as__path__expand  *string  `xml:"as-path-expand,omitempty"`
					V_p2mp__lsp__root  *string  `xml:"p2mp-lsp-root,omitempty"`
					V_metric4  *string  `xml:"metric4,omitempty"`
					V_analyze  *string  `xml:"analyze,omitempty"`
					V_next  *string  `xml:"next,omitempty"`
					V_bgp__output__queue__priority  *string  `xml:"bgp-output-queue-priority,omitempty"`
					V_prefix__length__range  *string  `xml:"prefix-length-range,omitempty"`
					V_label  *string  `xml:"label,omitempty"`
					V_get__route__range  *string  `xml:"get-route-range,omitempty"`
					V_destination__class  *string  `xml:"destination-class,omitempty"`
					V_label__allocation__fallback__reject  *string  `xml:"label-allocation-fallback-reject,omitempty"`
					V_preference2  *string  `xml:"preference2,omitempty"`
					V_validation__state  *string  `xml:"validation-state,omitempty"`
					V_aigp__originate  *string  `xml:"aigp-originate,omitempty"`
					V_class  *string  `xml:"class,omitempty"`
					V_source__class  *string  `xml:"source-class,omitempty"`
					V_selected__mldp__egress  *string  `xml:"selected-mldp-egress,omitempty"`
					V_preference  *string  `xml:"preference,omitempty"`
					V_next__hop  *string  `xml:"next-hop,omitempty"`
					V_color2  *string  `xml:"color2,omitempty"`
					V_cos__next__hop__map  *string  `xml:"cos-next-hop-map,omitempty"`
					V_sr__te__template  *string  `xml:"sr-te-template,omitempty"`
					V_metric  *string  `xml:"metric,omitempty"`
					V_local__preference  *string  `xml:"local-preference,omitempty"`
					V_load__balance  *string  `xml:"load-balance,omitempty"`
					V_external  *string  `xml:"external,omitempty"`
					V_dynamic__tunnel__attributes  *string  `xml:"dynamic-tunnel-attributes,omitempty"`
					V_no__backup  *string  `xml:"no-backup,omitempty"`
					V_through  *string  `xml:"through,omitempty"`
					V_limit__bandwidth  *string  `xml:"limit-bandwidth,omitempty"`
					V_install__nexthop  *string  `xml:"install-nexthop,omitempty"`
					V_aggregate__bandwidth  *string  `xml:"aggregate-bandwidth,omitempty"`
					V_no__route__localize  *string  `xml:"no-route-localize,omitempty"`
					V_multipath__resolve  *string  `xml:"multipath-resolve,omitempty"`
					V_prefix__segment  *string  `xml:"prefix-segment,omitempty"`
					V_label__allocation  *string  `xml:"label-allocation,omitempty"`
					V_origin  *string  `xml:"origin,omitempty"`
					V_mhop__bfd__port  *string  `xml:"mhop-bfd-port,omitempty"`
					V_upto  *string  `xml:"upto,omitempty"`
					V_color  *string  `xml:"color,omitempty"`
					V_priority  *string  `xml:"priority,omitempty"`
					V_exact  *string  `xml:"exact,omitempty"`
					V_forwarding__class  *string  `xml:"forwarding-class,omitempty"`
					V_resolution__map  *string  `xml:"resolution-map,omitempty"`
					V_address  *string  `xml:"address,omitempty"`
				} `xml:"from>route-filter"`
			} `xml:"term"`
		} `xml:"policy-options>policy-statement"`
	} `xml:"groups"`
	ApplyGroup string `xml:"apply-groups"`
}

// v_ is appended before every variable so it doesn't give any conflict
// with any keyword in golang. ex- interface is keyword in golang
func junosPolicy__OptionsPolicy__StatementTermFromRoute__FilterAddressCreate(d *schema.ResourceData, m interface{}) error {

	var err error
	client := m.(*ProviderConfig)

    id := d.Get("resource_name").(string)
     	V_name := d.Get("name").(string)
	V_name__1 := d.Get("name__1").(string)
	V_damping := d.Get("damping").(string)
	V_no__entropy__label__capability := d.Get("no__entropy__label__capability").(string)
	V_metric3 := d.Get("metric3").(string)
	V_tag2 := d.Get("tag2").(string)
	V_add__path := d.Get("add__path").(string)
	V_tunnel__end__point__address := d.Get("tunnel__end__point__address").(string)
	V_default__action := d.Get("default__action").(string)
	V_trace := d.Get("trace").(string)
	V_install__to__fib := d.Get("install__to__fib").(string)
	V_ssm__source := d.Get("ssm__source").(string)
	V_metric2 := d.Get("metric2").(string)
	V_tag := d.Get("tag").(string)
	V_longer := d.Get("longer").(string)
	V_orlonger := d.Get("orlonger").(string)
	V_address__mask := d.Get("address__mask").(string)
	V_community := d.Get("community").(string)
	V_tunnel__attribute := d.Get("tunnel__attribute").(string)
	V_as__path__prepend := d.Get("as__path__prepend").(string)
	V_map__to__interface := d.Get("map__to__interface").(string)
	V_accept_reject := d.Get("accept_reject").(string)
	V_apply__advanced := d.Get("apply__advanced").(string)
	V_aigp__adjust := d.Get("aigp__adjust").(string)
	V_as__path__expand := d.Get("as__path__expand").(string)
	V_p2mp__lsp__root := d.Get("p2mp__lsp__root").(string)
	V_metric4 := d.Get("metric4").(string)
	V_analyze := d.Get("analyze").(string)
	V_next := d.Get("next").(string)
	V_bgp__output__queue__priority := d.Get("bgp__output__queue__priority").(string)
	V_prefix__length__range := d.Get("prefix__length__range").(string)
	V_label := d.Get("label").(string)
	V_get__route__range := d.Get("get__route__range").(string)
	V_destination__class := d.Get("destination__class").(string)
	V_label__allocation__fallback__reject := d.Get("label__allocation__fallback__reject").(string)
	V_preference2 := d.Get("preference2").(string)
	V_validation__state := d.Get("validation__state").(string)
	V_aigp__originate := d.Get("aigp__originate").(string)
	V_class := d.Get("class").(string)
	V_source__class := d.Get("source__class").(string)
	V_selected__mldp__egress := d.Get("selected__mldp__egress").(string)
	V_preference := d.Get("preference").(string)
	V_next__hop := d.Get("next__hop").(string)
	V_color2 := d.Get("color2").(string)
	V_cos__next__hop__map := d.Get("cos__next__hop__map").(string)
	V_sr__te__template := d.Get("sr__te__template").(string)
	V_metric := d.Get("metric").(string)
	V_local__preference := d.Get("local__preference").(string)
	V_load__balance := d.Get("load__balance").(string)
	V_external := d.Get("external").(string)
	V_dynamic__tunnel__attributes := d.Get("dynamic__tunnel__attributes").(string)
	V_no__backup := d.Get("no__backup").(string)
	V_through := d.Get("through").(string)
	V_limit__bandwidth := d.Get("limit__bandwidth").(string)
	V_install__nexthop := d.Get("install__nexthop").(string)
	V_aggregate__bandwidth := d.Get("aggregate__bandwidth").(string)
	V_no__route__localize := d.Get("no__route__localize").(string)
	V_multipath__resolve := d.Get("multipath__resolve").(string)
	V_prefix__segment := d.Get("prefix__segment").(string)
	V_label__allocation := d.Get("label__allocation").(string)
	V_origin := d.Get("origin").(string)
	V_mhop__bfd__port := d.Get("mhop__bfd__port").(string)
	V_upto := d.Get("upto").(string)
	V_color := d.Get("color").(string)
	V_priority := d.Get("priority").(string)
	V_exact := d.Get("exact").(string)
	V_forwarding__class := d.Get("forwarding__class").(string)
	V_resolution__map := d.Get("resolution__map").(string)
	V_address := d.Get("address").(string)


	config := xmlPolicy__OptionsPolicy__StatementTermFromRoute__FilterAddress{}
	config.ApplyGroup = id
	config.Groups.Name = id
	config.Groups.V_policy__statement.V_name = &V_name
	config.Groups.V_policy__statement.V_term.V_name__1 = &V_name__1
	if V_damping != "" { config.Groups.V_policy__statement.V_term.V_route__filter.V_damping = &V_damping }
	if V_no__entropy__label__capability != "" { config.Groups.V_policy__statement.V_term.V_route__filter.V_no__entropy__label__capability = &V_no__entropy__label__capability }
	if V_metric3 != "" { config.Groups.V_policy__statement.V_term.V_route__filter.V_metric3 = &V_metric3 }
	if V_tag2 != "" { config.Groups.V_policy__statement.V_term.V_route__filter.V_tag2 = &V_tag2 }
	if V_add__path != "" { config.Groups.V_policy__statement.V_term.V_route__filter.V_add__path = &V_add__path }
	if V_tunnel__end__point__address != "" { config.Groups.V_policy__statement.V_term.V_route__filter.V_tunnel__end__point__address = &V_tunnel__end__point__address }
	if V_default__action != "" { config.Groups.V_policy__statement.V_term.V_route__filter.V_default__action = &V_default__action }
	if V_trace != "" { config.Groups.V_policy__statement.V_term.V_route__filter.V_trace = &V_trace }
	if V_install__to__fib != "" { config.Groups.V_policy__statement.V_term.V_route__filter.V_install__to__fib = &V_install__to__fib }
	if V_ssm__source != "" { config.Groups.V_policy__statement.V_term.V_route__filter.V_ssm__source = &V_ssm__source }
	if V_metric2 != "" { config.Groups.V_policy__statement.V_term.V_route__filter.V_metric2 = &V_metric2 }
	if V_tag != "" { config.Groups.V_policy__statement.V_term.V_route__filter.V_tag = &V_tag }
	if V_longer != "" { config.Groups.V_policy__statement.V_term.V_route__filter.V_longer = &V_longer }
	if V_orlonger != "" { config.Groups.V_policy__statement.V_term.V_route__filter.V_orlonger = &V_orlonger }
	if V_address__mask != "" { config.Groups.V_policy__statement.V_term.V_route__filter.V_address__mask = &V_address__mask }
	if V_community != "" { config.Groups.V_policy__statement.V_term.V_route__filter.V_community = &V_community }
	if V_tunnel__attribute != "" { config.Groups.V_policy__statement.V_term.V_route__filter.V_tunnel__attribute = &V_tunnel__attribute }
	if V_as__path__prepend != "" { config.Groups.V_policy__statement.V_term.V_route__filter.V_as__path__prepend = &V_as__path__prepend }
	if V_map__to__interface != "" { config.Groups.V_policy__statement.V_term.V_route__filter.V_map__to__interface = &V_map__to__interface }
	if V_accept_reject != "" { config.Groups.V_policy__statement.V_term.V_route__filter.V_accept_reject = &V_accept_reject }
	if V_apply__advanced != "" { config.Groups.V_policy__statement.V_term.V_route__filter.V_apply__advanced = &V_apply__advanced }
	if V_aigp__adjust != "" { config.Groups.V_policy__statement.V_term.V_route__filter.V_aigp__adjust = &V_aigp__adjust }
	if V_as__path__expand != "" { config.Groups.V_policy__statement.V_term.V_route__filter.V_as__path__expand = &V_as__path__expand }
	if V_p2mp__lsp__root != "" { config.Groups.V_policy__statement.V_term.V_route__filter.V_p2mp__lsp__root = &V_p2mp__lsp__root }
	if V_metric4 != "" { config.Groups.V_policy__statement.V_term.V_route__filter.V_metric4 = &V_metric4 }
	if V_analyze != "" { config.Groups.V_policy__statement.V_term.V_route__filter.V_analyze = &V_analyze }
	if V_next != "" { config.Groups.V_policy__statement.V_term.V_route__filter.V_next = &V_next }
	if V_bgp__output__queue__priority != "" { config.Groups.V_policy__statement.V_term.V_route__filter.V_bgp__output__queue__priority = &V_bgp__output__queue__priority }
	if V_prefix__length__range != "" { config.Groups.V_policy__statement.V_term.V_route__filter.V_prefix__length__range = &V_prefix__length__range }
	if V_label != "" { config.Groups.V_policy__statement.V_term.V_route__filter.V_label = &V_label }
	if V_get__route__range != "" { config.Groups.V_policy__statement.V_term.V_route__filter.V_get__route__range = &V_get__route__range }
	if V_destination__class != "" { config.Groups.V_policy__statement.V_term.V_route__filter.V_destination__class = &V_destination__class }
	if V_label__allocation__fallback__reject != "" { config.Groups.V_policy__statement.V_term.V_route__filter.V_label__allocation__fallback__reject = &V_label__allocation__fallback__reject }
	if V_preference2 != "" { config.Groups.V_policy__statement.V_term.V_route__filter.V_preference2 = &V_preference2 }
	if V_validation__state != "" { config.Groups.V_policy__statement.V_term.V_route__filter.V_validation__state = &V_validation__state }
	if V_aigp__originate != "" { config.Groups.V_policy__statement.V_term.V_route__filter.V_aigp__originate = &V_aigp__originate }
	if V_class != "" { config.Groups.V_policy__statement.V_term.V_route__filter.V_class = &V_class }
	if V_source__class != "" { config.Groups.V_policy__statement.V_term.V_route__filter.V_source__class = &V_source__class }
	if V_selected__mldp__egress != "" { config.Groups.V_policy__statement.V_term.V_route__filter.V_selected__mldp__egress = &V_selected__mldp__egress }
	if V_preference != "" { config.Groups.V_policy__statement.V_term.V_route__filter.V_preference = &V_preference }
	if V_next__hop != "" { config.Groups.V_policy__statement.V_term.V_route__filter.V_next__hop = &V_next__hop }
	if V_color2 != "" { config.Groups.V_policy__statement.V_term.V_route__filter.V_color2 = &V_color2 }
	if V_cos__next__hop__map != "" { config.Groups.V_policy__statement.V_term.V_route__filter.V_cos__next__hop__map = &V_cos__next__hop__map }
	if V_sr__te__template != "" { config.Groups.V_policy__statement.V_term.V_route__filter.V_sr__te__template = &V_sr__te__template }
	if V_metric != "" { config.Groups.V_policy__statement.V_term.V_route__filter.V_metric = &V_metric }
	if V_local__preference != "" { config.Groups.V_policy__statement.V_term.V_route__filter.V_local__preference = &V_local__preference }
	if V_load__balance != "" { config.Groups.V_policy__statement.V_term.V_route__filter.V_load__balance = &V_load__balance }
	if V_external != "" { config.Groups.V_policy__statement.V_term.V_route__filter.V_external = &V_external }
	if V_dynamic__tunnel__attributes != "" { config.Groups.V_policy__statement.V_term.V_route__filter.V_dynamic__tunnel__attributes = &V_dynamic__tunnel__attributes }
	if V_no__backup != "" { config.Groups.V_policy__statement.V_term.V_route__filter.V_no__backup = &V_no__backup }
	if V_through != "" { config.Groups.V_policy__statement.V_term.V_route__filter.V_through = &V_through }
	if V_limit__bandwidth != "" { config.Groups.V_policy__statement.V_term.V_route__filter.V_limit__bandwidth = &V_limit__bandwidth }
	if V_install__nexthop != "" { config.Groups.V_policy__statement.V_term.V_route__filter.V_install__nexthop = &V_install__nexthop }
	if V_aggregate__bandwidth != "" { config.Groups.V_policy__statement.V_term.V_route__filter.V_aggregate__bandwidth = &V_aggregate__bandwidth }
	if V_no__route__localize != "" { config.Groups.V_policy__statement.V_term.V_route__filter.V_no__route__localize = &V_no__route__localize }
	if V_multipath__resolve != "" { config.Groups.V_policy__statement.V_term.V_route__filter.V_multipath__resolve = &V_multipath__resolve }
	if V_prefix__segment != "" { config.Groups.V_policy__statement.V_term.V_route__filter.V_prefix__segment = &V_prefix__segment }
	if V_label__allocation != "" { config.Groups.V_policy__statement.V_term.V_route__filter.V_label__allocation = &V_label__allocation }
	if V_origin != "" { config.Groups.V_policy__statement.V_term.V_route__filter.V_origin = &V_origin }
	if V_mhop__bfd__port != "" { config.Groups.V_policy__statement.V_term.V_route__filter.V_mhop__bfd__port = &V_mhop__bfd__port }
	if V_upto != "" { config.Groups.V_policy__statement.V_term.V_route__filter.V_upto = &V_upto }
	if V_color != "" { config.Groups.V_policy__statement.V_term.V_route__filter.V_color = &V_color }
	if V_priority != "" { config.Groups.V_policy__statement.V_term.V_route__filter.V_priority = &V_priority }
	if V_exact != "" { config.Groups.V_policy__statement.V_term.V_route__filter.V_exact = &V_exact }
	if V_forwarding__class != "" { config.Groups.V_policy__statement.V_term.V_route__filter.V_forwarding__class = &V_forwarding__class }
	if V_resolution__map != "" { config.Groups.V_policy__statement.V_term.V_route__filter.V_resolution__map = &V_resolution__map }
	config.Groups.V_policy__statement.V_term.V_route__filter.V_address = &V_address

    err = client.SendTransaction("", config, false)
    check(err)
    
    d.SetId(fmt.Sprintf("%s_%s", client.Host, id))
    
	return junosPolicy__OptionsPolicy__StatementTermFromRoute__FilterAddressRead(d,m)
}

func junosPolicy__OptionsPolicy__StatementTermFromRoute__FilterAddressRead(d *schema.ResourceData, m interface{}) error {

	var err error
	client := m.(*ProviderConfig)

    id := d.Get("resource_name").(string)
    
	config := &xmlPolicy__OptionsPolicy__StatementTermFromRoute__FilterAddress{}

	err = client.MarshalGroup(id, config)
	check(err)
 	d.Set("name", config.Groups.V_policy__statement.V_name)
	d.Set("name__1", config.Groups.V_policy__statement.V_term.V_name__1)
	if config.Groups.V_policy__statement.V_term.V_route__filter.V_damping != nil { 
		d.Set("damping", " ") } else { d.Set("damping", "")}

	if config.Groups.V_policy__statement.V_term.V_route__filter.V_no__entropy__label__capability != nil { 
		d.Set("no__entropy__label__capability", " ") } else { d.Set("no__entropy__label__capability", "")}

	if config.Groups.V_policy__statement.V_term.V_route__filter.V_metric3 != nil { 
		d.Set("metric3", " ") } else { d.Set("metric3", "")}

	if config.Groups.V_policy__statement.V_term.V_route__filter.V_tag2 != nil { 
		d.Set("tag2", " ") } else { d.Set("tag2", "")}

	if config.Groups.V_policy__statement.V_term.V_route__filter.V_add__path != nil { 
		d.Set("add__path", " ") } else { d.Set("add__path", "")}

	if config.Groups.V_policy__statement.V_term.V_route__filter.V_tunnel__end__point__address != nil { 
		d.Set("tunnel__end__point__address", " ") } else { d.Set("tunnel__end__point__address", "")}

	if config.Groups.V_policy__statement.V_term.V_route__filter.V_default__action != nil { 
		d.Set("default__action", " ") } else { d.Set("default__action", "")}

	if config.Groups.V_policy__statement.V_term.V_route__filter.V_trace != nil { 
		d.Set("trace", " ") } else { d.Set("trace", "")}

	if config.Groups.V_policy__statement.V_term.V_route__filter.V_install__to__fib != nil { 
		d.Set("install__to__fib", " ") } else { d.Set("install__to__fib", "")}

	if config.Groups.V_policy__statement.V_term.V_route__filter.V_ssm__source != nil { 
		d.Set("ssm__source", " ") } else { d.Set("ssm__source", "")}

	if config.Groups.V_policy__statement.V_term.V_route__filter.V_metric2 != nil { 
		d.Set("metric2", " ") } else { d.Set("metric2", "")}

	if config.Groups.V_policy__statement.V_term.V_route__filter.V_tag != nil { 
		d.Set("tag", " ") } else { d.Set("tag", "")}

	if config.Groups.V_policy__statement.V_term.V_route__filter.V_longer != nil { 
		d.Set("longer", " ") } else { d.Set("longer", "")}

	if config.Groups.V_policy__statement.V_term.V_route__filter.V_orlonger != nil { 
		d.Set("orlonger", " ") } else { d.Set("orlonger", "")}

	if config.Groups.V_policy__statement.V_term.V_route__filter.V_address__mask != nil { 
		d.Set("address__mask", " ") } else { d.Set("address__mask", "")}

	if config.Groups.V_policy__statement.V_term.V_route__filter.V_community != nil { 
		d.Set("community", " ") } else { d.Set("community", "")}

	if config.Groups.V_policy__statement.V_term.V_route__filter.V_tunnel__attribute != nil { 
		d.Set("tunnel__attribute", " ") } else { d.Set("tunnel__attribute", "")}

	if config.Groups.V_policy__statement.V_term.V_route__filter.V_as__path__prepend != nil { 
		d.Set("as__path__prepend", " ") } else { d.Set("as__path__prepend", "")}

	if config.Groups.V_policy__statement.V_term.V_route__filter.V_map__to__interface != nil { 
		d.Set("map__to__interface", " ") } else { d.Set("map__to__interface", "")}

	if config.Groups.V_policy__statement.V_term.V_route__filter.V_accept_reject != nil { 
		d.Set("accept_reject", " ") } else { d.Set("accept_reject", "")}

	if config.Groups.V_policy__statement.V_term.V_route__filter.V_apply__advanced != nil { 
		d.Set("apply__advanced", " ") } else { d.Set("apply__advanced", "")}

	if config.Groups.V_policy__statement.V_term.V_route__filter.V_aigp__adjust != nil { 
		d.Set("aigp__adjust", " ") } else { d.Set("aigp__adjust", "")}

	if config.Groups.V_policy__statement.V_term.V_route__filter.V_as__path__expand != nil { 
		d.Set("as__path__expand", " ") } else { d.Set("as__path__expand", "")}

	if config.Groups.V_policy__statement.V_term.V_route__filter.V_p2mp__lsp__root != nil { 
		d.Set("p2mp__lsp__root", " ") } else { d.Set("p2mp__lsp__root", "")}

	if config.Groups.V_policy__statement.V_term.V_route__filter.V_metric4 != nil { 
		d.Set("metric4", " ") } else { d.Set("metric4", "")}

	if config.Groups.V_policy__statement.V_term.V_route__filter.V_analyze != nil { 
		d.Set("analyze", " ") } else { d.Set("analyze", "")}

	if config.Groups.V_policy__statement.V_term.V_route__filter.V_next != nil { 
		d.Set("next", " ") } else { d.Set("next", "")}

	if config.Groups.V_policy__statement.V_term.V_route__filter.V_bgp__output__queue__priority != nil { 
		d.Set("bgp__output__queue__priority", " ") } else { d.Set("bgp__output__queue__priority", "")}

	if config.Groups.V_policy__statement.V_term.V_route__filter.V_prefix__length__range != nil { 
		d.Set("prefix__length__range", " ") } else { d.Set("prefix__length__range", "")}

	if config.Groups.V_policy__statement.V_term.V_route__filter.V_label != nil { 
		d.Set("label", " ") } else { d.Set("label", "")}

	if config.Groups.V_policy__statement.V_term.V_route__filter.V_get__route__range != nil { 
		d.Set("get__route__range", " ") } else { d.Set("get__route__range", "")}

	if config.Groups.V_policy__statement.V_term.V_route__filter.V_destination__class != nil { 
		d.Set("destination__class", " ") } else { d.Set("destination__class", "")}

	if config.Groups.V_policy__statement.V_term.V_route__filter.V_label__allocation__fallback__reject != nil { 
		d.Set("label__allocation__fallback__reject", " ") } else { d.Set("label__allocation__fallback__reject", "")}

	if config.Groups.V_policy__statement.V_term.V_route__filter.V_preference2 != nil { 
		d.Set("preference2", " ") } else { d.Set("preference2", "")}

	if config.Groups.V_policy__statement.V_term.V_route__filter.V_validation__state != nil { 
		d.Set("validation__state", " ") } else { d.Set("validation__state", "")}

	if config.Groups.V_policy__statement.V_term.V_route__filter.V_aigp__originate != nil { 
		d.Set("aigp__originate", " ") } else { d.Set("aigp__originate", "")}

	if config.Groups.V_policy__statement.V_term.V_route__filter.V_class != nil { 
		d.Set("class", " ") } else { d.Set("class", "")}

	if config.Groups.V_policy__statement.V_term.V_route__filter.V_source__class != nil { 
		d.Set("source__class", " ") } else { d.Set("source__class", "")}

	if config.Groups.V_policy__statement.V_term.V_route__filter.V_selected__mldp__egress != nil { 
		d.Set("selected__mldp__egress", " ") } else { d.Set("selected__mldp__egress", "")}

	if config.Groups.V_policy__statement.V_term.V_route__filter.V_preference != nil { 
		d.Set("preference", " ") } else { d.Set("preference", "")}

	if config.Groups.V_policy__statement.V_term.V_route__filter.V_next__hop != nil { 
		d.Set("next__hop", " ") } else { d.Set("next__hop", "")}

	if config.Groups.V_policy__statement.V_term.V_route__filter.V_color2 != nil { 
		d.Set("color2", " ") } else { d.Set("color2", "")}

	if config.Groups.V_policy__statement.V_term.V_route__filter.V_cos__next__hop__map != nil { 
		d.Set("cos__next__hop__map", " ") } else { d.Set("cos__next__hop__map", "")}

	if config.Groups.V_policy__statement.V_term.V_route__filter.V_sr__te__template != nil { 
		d.Set("sr__te__template", " ") } else { d.Set("sr__te__template", "")}

	if config.Groups.V_policy__statement.V_term.V_route__filter.V_metric != nil { 
		d.Set("metric", " ") } else { d.Set("metric", "")}

	if config.Groups.V_policy__statement.V_term.V_route__filter.V_local__preference != nil { 
		d.Set("local__preference", " ") } else { d.Set("local__preference", "")}

	if config.Groups.V_policy__statement.V_term.V_route__filter.V_load__balance != nil { 
		d.Set("load__balance", " ") } else { d.Set("load__balance", "")}

	if config.Groups.V_policy__statement.V_term.V_route__filter.V_external != nil { 
		d.Set("external", " ") } else { d.Set("external", "")}

	if config.Groups.V_policy__statement.V_term.V_route__filter.V_dynamic__tunnel__attributes != nil { 
		d.Set("dynamic__tunnel__attributes", " ") } else { d.Set("dynamic__tunnel__attributes", "")}

	if config.Groups.V_policy__statement.V_term.V_route__filter.V_no__backup != nil { 
		d.Set("no__backup", " ") } else { d.Set("no__backup", "")}

	if config.Groups.V_policy__statement.V_term.V_route__filter.V_through != nil { 
		d.Set("through", " ") } else { d.Set("through", "")}

	if config.Groups.V_policy__statement.V_term.V_route__filter.V_limit__bandwidth != nil { 
		d.Set("limit__bandwidth", " ") } else { d.Set("limit__bandwidth", "")}

	if config.Groups.V_policy__statement.V_term.V_route__filter.V_install__nexthop != nil { 
		d.Set("install__nexthop", " ") } else { d.Set("install__nexthop", "")}

	if config.Groups.V_policy__statement.V_term.V_route__filter.V_aggregate__bandwidth != nil { 
		d.Set("aggregate__bandwidth", " ") } else { d.Set("aggregate__bandwidth", "")}

	if config.Groups.V_policy__statement.V_term.V_route__filter.V_no__route__localize != nil { 
		d.Set("no__route__localize", " ") } else { d.Set("no__route__localize", "")}

	if config.Groups.V_policy__statement.V_term.V_route__filter.V_multipath__resolve != nil { 
		d.Set("multipath__resolve", " ") } else { d.Set("multipath__resolve", "")}

	if config.Groups.V_policy__statement.V_term.V_route__filter.V_prefix__segment != nil { 
		d.Set("prefix__segment", " ") } else { d.Set("prefix__segment", "")}

	if config.Groups.V_policy__statement.V_term.V_route__filter.V_label__allocation != nil { 
		d.Set("label__allocation", " ") } else { d.Set("label__allocation", "")}

	if config.Groups.V_policy__statement.V_term.V_route__filter.V_origin != nil { 
		d.Set("origin", " ") } else { d.Set("origin", "")}

	if config.Groups.V_policy__statement.V_term.V_route__filter.V_mhop__bfd__port != nil { 
		d.Set("mhop__bfd__port", " ") } else { d.Set("mhop__bfd__port", "")}

	if config.Groups.V_policy__statement.V_term.V_route__filter.V_upto != nil { 
		d.Set("upto", " ") } else { d.Set("upto", "")}

	if config.Groups.V_policy__statement.V_term.V_route__filter.V_color != nil { 
		d.Set("color", " ") } else { d.Set("color", "")}

	if config.Groups.V_policy__statement.V_term.V_route__filter.V_priority != nil { 
		d.Set("priority", " ") } else { d.Set("priority", "")}

	if config.Groups.V_policy__statement.V_term.V_route__filter.V_exact != nil { 
		d.Set("exact", " ") } else { d.Set("exact", "")}

	if config.Groups.V_policy__statement.V_term.V_route__filter.V_forwarding__class != nil { 
		d.Set("forwarding__class", " ") } else { d.Set("forwarding__class", "")}

	if config.Groups.V_policy__statement.V_term.V_route__filter.V_resolution__map != nil { 
		d.Set("resolution__map", " ") } else { d.Set("resolution__map", "")}

	d.Set("address", config.Groups.V_policy__statement.V_term.V_route__filter.V_address)

	return nil
}

func junosPolicy__OptionsPolicy__StatementTermFromRoute__FilterAddressUpdate(d *schema.ResourceData, m interface{}) error {

	var err error
	client := m.(*ProviderConfig)

    id := d.Get("resource_name").(string)
     	V_name := d.Get("name").(string)
	V_name__1 := d.Get("name__1").(string)
	V_damping := d.Get("damping").(string)
	V_no__entropy__label__capability := d.Get("no__entropy__label__capability").(string)
	V_metric3 := d.Get("metric3").(string)
	V_tag2 := d.Get("tag2").(string)
	V_add__path := d.Get("add__path").(string)
	V_tunnel__end__point__address := d.Get("tunnel__end__point__address").(string)
	V_default__action := d.Get("default__action").(string)
	V_trace := d.Get("trace").(string)
	V_install__to__fib := d.Get("install__to__fib").(string)
	V_ssm__source := d.Get("ssm__source").(string)
	V_metric2 := d.Get("metric2").(string)
	V_tag := d.Get("tag").(string)
	V_longer := d.Get("longer").(string)
	V_orlonger := d.Get("orlonger").(string)
	V_address__mask := d.Get("address__mask").(string)
	V_community := d.Get("community").(string)
	V_tunnel__attribute := d.Get("tunnel__attribute").(string)
	V_as__path__prepend := d.Get("as__path__prepend").(string)
	V_map__to__interface := d.Get("map__to__interface").(string)
	V_accept_reject := d.Get("accept_reject").(string)
	V_apply__advanced := d.Get("apply__advanced").(string)
	V_aigp__adjust := d.Get("aigp__adjust").(string)
	V_as__path__expand := d.Get("as__path__expand").(string)
	V_p2mp__lsp__root := d.Get("p2mp__lsp__root").(string)
	V_metric4 := d.Get("metric4").(string)
	V_analyze := d.Get("analyze").(string)
	V_next := d.Get("next").(string)
	V_bgp__output__queue__priority := d.Get("bgp__output__queue__priority").(string)
	V_prefix__length__range := d.Get("prefix__length__range").(string)
	V_label := d.Get("label").(string)
	V_get__route__range := d.Get("get__route__range").(string)
	V_destination__class := d.Get("destination__class").(string)
	V_label__allocation__fallback__reject := d.Get("label__allocation__fallback__reject").(string)
	V_preference2 := d.Get("preference2").(string)
	V_validation__state := d.Get("validation__state").(string)
	V_aigp__originate := d.Get("aigp__originate").(string)
	V_class := d.Get("class").(string)
	V_source__class := d.Get("source__class").(string)
	V_selected__mldp__egress := d.Get("selected__mldp__egress").(string)
	V_preference := d.Get("preference").(string)
	V_next__hop := d.Get("next__hop").(string)
	V_color2 := d.Get("color2").(string)
	V_cos__next__hop__map := d.Get("cos__next__hop__map").(string)
	V_sr__te__template := d.Get("sr__te__template").(string)
	V_metric := d.Get("metric").(string)
	V_local__preference := d.Get("local__preference").(string)
	V_load__balance := d.Get("load__balance").(string)
	V_external := d.Get("external").(string)
	V_dynamic__tunnel__attributes := d.Get("dynamic__tunnel__attributes").(string)
	V_no__backup := d.Get("no__backup").(string)
	V_through := d.Get("through").(string)
	V_limit__bandwidth := d.Get("limit__bandwidth").(string)
	V_install__nexthop := d.Get("install__nexthop").(string)
	V_aggregate__bandwidth := d.Get("aggregate__bandwidth").(string)
	V_no__route__localize := d.Get("no__route__localize").(string)
	V_multipath__resolve := d.Get("multipath__resolve").(string)
	V_prefix__segment := d.Get("prefix__segment").(string)
	V_label__allocation := d.Get("label__allocation").(string)
	V_origin := d.Get("origin").(string)
	V_mhop__bfd__port := d.Get("mhop__bfd__port").(string)
	V_upto := d.Get("upto").(string)
	V_color := d.Get("color").(string)
	V_priority := d.Get("priority").(string)
	V_exact := d.Get("exact").(string)
	V_forwarding__class := d.Get("forwarding__class").(string)
	V_resolution__map := d.Get("resolution__map").(string)
	V_address := d.Get("address").(string)


	config := xmlPolicy__OptionsPolicy__StatementTermFromRoute__FilterAddress{}
	config.ApplyGroup = id
	config.Groups.Name = id
	config.Groups.V_policy__statement.V_name = &V_name
	config.Groups.V_policy__statement.V_term.V_name__1 = &V_name__1
	if V_damping != "" { config.Groups.V_policy__statement.V_term.V_route__filter.V_damping = &V_damping }
	if V_no__entropy__label__capability != "" { config.Groups.V_policy__statement.V_term.V_route__filter.V_no__entropy__label__capability = &V_no__entropy__label__capability }
	if V_metric3 != "" { config.Groups.V_policy__statement.V_term.V_route__filter.V_metric3 = &V_metric3 }
	if V_tag2 != "" { config.Groups.V_policy__statement.V_term.V_route__filter.V_tag2 = &V_tag2 }
	if V_add__path != "" { config.Groups.V_policy__statement.V_term.V_route__filter.V_add__path = &V_add__path }
	if V_tunnel__end__point__address != "" { config.Groups.V_policy__statement.V_term.V_route__filter.V_tunnel__end__point__address = &V_tunnel__end__point__address }
	if V_default__action != "" { config.Groups.V_policy__statement.V_term.V_route__filter.V_default__action = &V_default__action }
	if V_trace != "" { config.Groups.V_policy__statement.V_term.V_route__filter.V_trace = &V_trace }
	if V_install__to__fib != "" { config.Groups.V_policy__statement.V_term.V_route__filter.V_install__to__fib = &V_install__to__fib }
	if V_ssm__source != "" { config.Groups.V_policy__statement.V_term.V_route__filter.V_ssm__source = &V_ssm__source }
	if V_metric2 != "" { config.Groups.V_policy__statement.V_term.V_route__filter.V_metric2 = &V_metric2 }
	if V_tag != "" { config.Groups.V_policy__statement.V_term.V_route__filter.V_tag = &V_tag }
	if V_longer != "" { config.Groups.V_policy__statement.V_term.V_route__filter.V_longer = &V_longer }
	if V_orlonger != "" { config.Groups.V_policy__statement.V_term.V_route__filter.V_orlonger = &V_orlonger }
	if V_address__mask != "" { config.Groups.V_policy__statement.V_term.V_route__filter.V_address__mask = &V_address__mask }
	if V_community != "" { config.Groups.V_policy__statement.V_term.V_route__filter.V_community = &V_community }
	if V_tunnel__attribute != "" { config.Groups.V_policy__statement.V_term.V_route__filter.V_tunnel__attribute = &V_tunnel__attribute }
	if V_as__path__prepend != "" { config.Groups.V_policy__statement.V_term.V_route__filter.V_as__path__prepend = &V_as__path__prepend }
	if V_map__to__interface != "" { config.Groups.V_policy__statement.V_term.V_route__filter.V_map__to__interface = &V_map__to__interface }
	if V_accept_reject != "" { config.Groups.V_policy__statement.V_term.V_route__filter.V_accept_reject = &V_accept_reject }
	if V_apply__advanced != "" { config.Groups.V_policy__statement.V_term.V_route__filter.V_apply__advanced = &V_apply__advanced }
	if V_aigp__adjust != "" { config.Groups.V_policy__statement.V_term.V_route__filter.V_aigp__adjust = &V_aigp__adjust }
	if V_as__path__expand != "" { config.Groups.V_policy__statement.V_term.V_route__filter.V_as__path__expand = &V_as__path__expand }
	if V_p2mp__lsp__root != "" { config.Groups.V_policy__statement.V_term.V_route__filter.V_p2mp__lsp__root = &V_p2mp__lsp__root }
	if V_metric4 != "" { config.Groups.V_policy__statement.V_term.V_route__filter.V_metric4 = &V_metric4 }
	if V_analyze != "" { config.Groups.V_policy__statement.V_term.V_route__filter.V_analyze = &V_analyze }
	if V_next != "" { config.Groups.V_policy__statement.V_term.V_route__filter.V_next = &V_next }
	if V_bgp__output__queue__priority != "" { config.Groups.V_policy__statement.V_term.V_route__filter.V_bgp__output__queue__priority = &V_bgp__output__queue__priority }
	if V_prefix__length__range != "" { config.Groups.V_policy__statement.V_term.V_route__filter.V_prefix__length__range = &V_prefix__length__range }
	if V_label != "" { config.Groups.V_policy__statement.V_term.V_route__filter.V_label = &V_label }
	if V_get__route__range != "" { config.Groups.V_policy__statement.V_term.V_route__filter.V_get__route__range = &V_get__route__range }
	if V_destination__class != "" { config.Groups.V_policy__statement.V_term.V_route__filter.V_destination__class = &V_destination__class }
	if V_label__allocation__fallback__reject != "" { config.Groups.V_policy__statement.V_term.V_route__filter.V_label__allocation__fallback__reject = &V_label__allocation__fallback__reject }
	if V_preference2 != "" { config.Groups.V_policy__statement.V_term.V_route__filter.V_preference2 = &V_preference2 }
	if V_validation__state != "" { config.Groups.V_policy__statement.V_term.V_route__filter.V_validation__state = &V_validation__state }
	if V_aigp__originate != "" { config.Groups.V_policy__statement.V_term.V_route__filter.V_aigp__originate = &V_aigp__originate }
	if V_class != "" { config.Groups.V_policy__statement.V_term.V_route__filter.V_class = &V_class }
	if V_source__class != "" { config.Groups.V_policy__statement.V_term.V_route__filter.V_source__class = &V_source__class }
	if V_selected__mldp__egress != "" { config.Groups.V_policy__statement.V_term.V_route__filter.V_selected__mldp__egress = &V_selected__mldp__egress }
	if V_preference != "" { config.Groups.V_policy__statement.V_term.V_route__filter.V_preference = &V_preference }
	if V_next__hop != "" { config.Groups.V_policy__statement.V_term.V_route__filter.V_next__hop = &V_next__hop }
	if V_color2 != "" { config.Groups.V_policy__statement.V_term.V_route__filter.V_color2 = &V_color2 }
	if V_cos__next__hop__map != "" { config.Groups.V_policy__statement.V_term.V_route__filter.V_cos__next__hop__map = &V_cos__next__hop__map }
	if V_sr__te__template != "" { config.Groups.V_policy__statement.V_term.V_route__filter.V_sr__te__template = &V_sr__te__template }
	if V_metric != "" { config.Groups.V_policy__statement.V_term.V_route__filter.V_metric = &V_metric }
	if V_local__preference != "" { config.Groups.V_policy__statement.V_term.V_route__filter.V_local__preference = &V_local__preference }
	if V_load__balance != "" { config.Groups.V_policy__statement.V_term.V_route__filter.V_load__balance = &V_load__balance }
	if V_external != "" { config.Groups.V_policy__statement.V_term.V_route__filter.V_external = &V_external }
	if V_dynamic__tunnel__attributes != "" { config.Groups.V_policy__statement.V_term.V_route__filter.V_dynamic__tunnel__attributes = &V_dynamic__tunnel__attributes }
	if V_no__backup != "" { config.Groups.V_policy__statement.V_term.V_route__filter.V_no__backup = &V_no__backup }
	if V_through != "" { config.Groups.V_policy__statement.V_term.V_route__filter.V_through = &V_through }
	if V_limit__bandwidth != "" { config.Groups.V_policy__statement.V_term.V_route__filter.V_limit__bandwidth = &V_limit__bandwidth }
	if V_install__nexthop != "" { config.Groups.V_policy__statement.V_term.V_route__filter.V_install__nexthop = &V_install__nexthop }
	if V_aggregate__bandwidth != "" { config.Groups.V_policy__statement.V_term.V_route__filter.V_aggregate__bandwidth = &V_aggregate__bandwidth }
	if V_no__route__localize != "" { config.Groups.V_policy__statement.V_term.V_route__filter.V_no__route__localize = &V_no__route__localize }
	if V_multipath__resolve != "" { config.Groups.V_policy__statement.V_term.V_route__filter.V_multipath__resolve = &V_multipath__resolve }
	if V_prefix__segment != "" { config.Groups.V_policy__statement.V_term.V_route__filter.V_prefix__segment = &V_prefix__segment }
	if V_label__allocation != "" { config.Groups.V_policy__statement.V_term.V_route__filter.V_label__allocation = &V_label__allocation }
	if V_origin != "" { config.Groups.V_policy__statement.V_term.V_route__filter.V_origin = &V_origin }
	if V_mhop__bfd__port != "" { config.Groups.V_policy__statement.V_term.V_route__filter.V_mhop__bfd__port = &V_mhop__bfd__port }
	if V_upto != "" { config.Groups.V_policy__statement.V_term.V_route__filter.V_upto = &V_upto }
	if V_color != "" { config.Groups.V_policy__statement.V_term.V_route__filter.V_color = &V_color }
	if V_priority != "" { config.Groups.V_policy__statement.V_term.V_route__filter.V_priority = &V_priority }
	if V_exact != "" { config.Groups.V_policy__statement.V_term.V_route__filter.V_exact = &V_exact }
	if V_forwarding__class != "" { config.Groups.V_policy__statement.V_term.V_route__filter.V_forwarding__class = &V_forwarding__class }
	if V_resolution__map != "" { config.Groups.V_policy__statement.V_term.V_route__filter.V_resolution__map = &V_resolution__map }
	config.Groups.V_policy__statement.V_term.V_route__filter.V_address = &V_address

    err = client.SendTransaction(id, config, false)
    check(err)
    
	return junosPolicy__OptionsPolicy__StatementTermFromRoute__FilterAddressRead(d,m)
}

func junosPolicy__OptionsPolicy__StatementTermFromRoute__FilterAddressDelete(d *schema.ResourceData, m interface{}) error {

	var err error
	client := m.(*ProviderConfig)

    id := d.Get("resource_name").(string)
    
    _, err = client.DeleteConfigNoCommit(id)
    check(err)

    d.SetId("")
    
	return nil
}

func junosPolicy__OptionsPolicy__StatementTermFromRoute__FilterAddress() *schema.Resource {
	return &schema.Resource{
		Create: junosPolicy__OptionsPolicy__StatementTermFromRoute__FilterAddressCreate,
		Read: junosPolicy__OptionsPolicy__StatementTermFromRoute__FilterAddressRead,
		Update: junosPolicy__OptionsPolicy__StatementTermFromRoute__FilterAddressUpdate,
		Delete: junosPolicy__OptionsPolicy__StatementTermFromRoute__FilterAddressDelete,

        Schema: map[string]*schema.Schema{
            "resource_name": &schema.Schema{
                Type:    schema.TypeString,
                Required: true,
            },
			"name": &schema.Schema{
				Type:    schema.TypeString,
				Optional: true,
				Description:    "xpath is: config.Groups.V_policy__statement",
			},
			"name__1": &schema.Schema{
				Type:    schema.TypeString,
				Optional: true,
				Description:    "xpath is: config.Groups.V_policy__statement.V_term",
			},
			"damping": &schema.Schema{
				Type:    schema.TypeString,
				Optional: true,
				Description:    "xpath is: config.Groups.V_policy__statement.V_term.V_route__filter",
			},
			"no__entropy__label__capability": &schema.Schema{
				Type:    schema.TypeString,
				Optional: true,
				Description:    "xpath is: config.Groups.V_policy__statement.V_term.V_route__filter",
			},
			"metric3": &schema.Schema{
				Type:    schema.TypeString,
				Optional: true,
				Description:    "xpath is: config.Groups.V_policy__statement.V_term.V_route__filter",
			},
			"tag2": &schema.Schema{
				Type:    schema.TypeString,
				Optional: true,
				Description:    "xpath is: config.Groups.V_policy__statement.V_term.V_route__filter",
			},
			"add__path": &schema.Schema{
				Type:    schema.TypeString,
				Optional: true,
				Description:    "xpath is: config.Groups.V_policy__statement.V_term.V_route__filter",
			},
			"tunnel__end__point__address": &schema.Schema{
				Type:    schema.TypeString,
				Optional: true,
				Description:    "xpath is: config.Groups.V_policy__statement.V_term.V_route__filter",
			},
			"default__action": &schema.Schema{
				Type:    schema.TypeString,
				Optional: true,
				Description:    "xpath is: config.Groups.V_policy__statement.V_term.V_route__filter",
			},
			"trace": &schema.Schema{
				Type:    schema.TypeString,
				Optional: true,
				Description:    "xpath is: config.Groups.V_policy__statement.V_term.V_route__filter",
			},
			"install__to__fib": &schema.Schema{
				Type:    schema.TypeString,
				Optional: true,
				Description:    "xpath is: config.Groups.V_policy__statement.V_term.V_route__filter",
			},
			"ssm__source": &schema.Schema{
				Type:    schema.TypeString,
				Optional: true,
				Description:    "xpath is: config.Groups.V_policy__statement.V_term.V_route__filter",
			},
			"metric2": &schema.Schema{
				Type:    schema.TypeString,
				Optional: true,
				Description:    "xpath is: config.Groups.V_policy__statement.V_term.V_route__filter",
			},
			"tag": &schema.Schema{
				Type:    schema.TypeString,
				Optional: true,
				Description:    "xpath is: config.Groups.V_policy__statement.V_term.V_route__filter",
			},
			"longer": &schema.Schema{
				Type:    schema.TypeString,
				Optional: true,
				Description:    "xpath is: config.Groups.V_policy__statement.V_term.V_route__filter",
			},
			"orlonger": &schema.Schema{
				Type:    schema.TypeString,
				Optional: true,
				Description:    "xpath is: config.Groups.V_policy__statement.V_term.V_route__filter",
			},
			"address__mask": &schema.Schema{
				Type:    schema.TypeString,
				Optional: true,
				Description:    "xpath is: config.Groups.V_policy__statement.V_term.V_route__filter",
			},
			"community": &schema.Schema{
				Type:    schema.TypeString,
				Optional: true,
				Description:    "xpath is: config.Groups.V_policy__statement.V_term.V_route__filter",
			},
			"tunnel__attribute": &schema.Schema{
				Type:    schema.TypeString,
				Optional: true,
				Description:    "xpath is: config.Groups.V_policy__statement.V_term.V_route__filter",
			},
			"as__path__prepend": &schema.Schema{
				Type:    schema.TypeString,
				Optional: true,
				Description:    "xpath is: config.Groups.V_policy__statement.V_term.V_route__filter",
			},
			"map__to__interface": &schema.Schema{
				Type:    schema.TypeString,
				Optional: true,
				Description:    "xpath is: config.Groups.V_policy__statement.V_term.V_route__filter",
			},
			"accept_reject": &schema.Schema{
				Type:    schema.TypeString,
				Optional: true,
				Description:    "xpath is: config.Groups.V_policy__statement.V_term.V_route__filter",
			},
			"apply__advanced": &schema.Schema{
				Type:    schema.TypeString,
				Optional: true,
				Description:    "xpath is: config.Groups.V_policy__statement.V_term.V_route__filter",
			},
			"aigp__adjust": &schema.Schema{
				Type:    schema.TypeString,
				Optional: true,
				Description:    "xpath is: config.Groups.V_policy__statement.V_term.V_route__filter",
			},
			"as__path__expand": &schema.Schema{
				Type:    schema.TypeString,
				Optional: true,
				Description:    "xpath is: config.Groups.V_policy__statement.V_term.V_route__filter",
			},
			"p2mp__lsp__root": &schema.Schema{
				Type:    schema.TypeString,
				Optional: true,
				Description:    "xpath is: config.Groups.V_policy__statement.V_term.V_route__filter",
			},
			"metric4": &schema.Schema{
				Type:    schema.TypeString,
				Optional: true,
				Description:    "xpath is: config.Groups.V_policy__statement.V_term.V_route__filter",
			},
			"analyze": &schema.Schema{
				Type:    schema.TypeString,
				Optional: true,
				Description:    "xpath is: config.Groups.V_policy__statement.V_term.V_route__filter",
			},
			"next": &schema.Schema{
				Type:    schema.TypeString,
				Optional: true,
				Description:    "xpath is: config.Groups.V_policy__statement.V_term.V_route__filter",
			},
			"bgp__output__queue__priority": &schema.Schema{
				Type:    schema.TypeString,
				Optional: true,
				Description:    "xpath is: config.Groups.V_policy__statement.V_term.V_route__filter",
			},
			"prefix__length__range": &schema.Schema{
				Type:    schema.TypeString,
				Optional: true,
				Description:    "xpath is: config.Groups.V_policy__statement.V_term.V_route__filter",
			},
			"label": &schema.Schema{
				Type:    schema.TypeString,
				Optional: true,
				Description:    "xpath is: config.Groups.V_policy__statement.V_term.V_route__filter",
			},
			"get__route__range": &schema.Schema{
				Type:    schema.TypeString,
				Optional: true,
				Description:    "xpath is: config.Groups.V_policy__statement.V_term.V_route__filter",
			},
			"destination__class": &schema.Schema{
				Type:    schema.TypeString,
				Optional: true,
				Description:    "xpath is: config.Groups.V_policy__statement.V_term.V_route__filter",
			},
			"label__allocation__fallback__reject": &schema.Schema{
				Type:    schema.TypeString,
				Optional: true,
				Description:    "xpath is: config.Groups.V_policy__statement.V_term.V_route__filter",
			},
			"preference2": &schema.Schema{
				Type:    schema.TypeString,
				Optional: true,
				Description:    "xpath is: config.Groups.V_policy__statement.V_term.V_route__filter",
			},
			"validation__state": &schema.Schema{
				Type:    schema.TypeString,
				Optional: true,
				Description:    "xpath is: config.Groups.V_policy__statement.V_term.V_route__filter",
			},
			"aigp__originate": &schema.Schema{
				Type:    schema.TypeString,
				Optional: true,
				Description:    "xpath is: config.Groups.V_policy__statement.V_term.V_route__filter",
			},
			"class": &schema.Schema{
				Type:    schema.TypeString,
				Optional: true,
				Description:    "xpath is: config.Groups.V_policy__statement.V_term.V_route__filter",
			},
			"source__class": &schema.Schema{
				Type:    schema.TypeString,
				Optional: true,
				Description:    "xpath is: config.Groups.V_policy__statement.V_term.V_route__filter",
			},
			"selected__mldp__egress": &schema.Schema{
				Type:    schema.TypeString,
				Optional: true,
				Description:    "xpath is: config.Groups.V_policy__statement.V_term.V_route__filter",
			},
			"preference": &schema.Schema{
				Type:    schema.TypeString,
				Optional: true,
				Description:    "xpath is: config.Groups.V_policy__statement.V_term.V_route__filter",
			},
			"next__hop": &schema.Schema{
				Type:    schema.TypeString,
				Optional: true,
				Description:    "xpath is: config.Groups.V_policy__statement.V_term.V_route__filter",
			},
			"color2": &schema.Schema{
				Type:    schema.TypeString,
				Optional: true,
				Description:    "xpath is: config.Groups.V_policy__statement.V_term.V_route__filter",
			},
			"cos__next__hop__map": &schema.Schema{
				Type:    schema.TypeString,
				Optional: true,
				Description:    "xpath is: config.Groups.V_policy__statement.V_term.V_route__filter",
			},
			"sr__te__template": &schema.Schema{
				Type:    schema.TypeString,
				Optional: true,
				Description:    "xpath is: config.Groups.V_policy__statement.V_term.V_route__filter",
			},
			"metric": &schema.Schema{
				Type:    schema.TypeString,
				Optional: true,
				Description:    "xpath is: config.Groups.V_policy__statement.V_term.V_route__filter",
			},
			"local__preference": &schema.Schema{
				Type:    schema.TypeString,
				Optional: true,
				Description:    "xpath is: config.Groups.V_policy__statement.V_term.V_route__filter",
			},
			"load__balance": &schema.Schema{
				Type:    schema.TypeString,
				Optional: true,
				Description:    "xpath is: config.Groups.V_policy__statement.V_term.V_route__filter",
			},
			"external": &schema.Schema{
				Type:    schema.TypeString,
				Optional: true,
				Description:    "xpath is: config.Groups.V_policy__statement.V_term.V_route__filter",
			},
			"dynamic__tunnel__attributes": &schema.Schema{
				Type:    schema.TypeString,
				Optional: true,
				Description:    "xpath is: config.Groups.V_policy__statement.V_term.V_route__filter",
			},
			"no__backup": &schema.Schema{
				Type:    schema.TypeString,
				Optional: true,
				Description:    "xpath is: config.Groups.V_policy__statement.V_term.V_route__filter",
			},
			"through": &schema.Schema{
				Type:    schema.TypeString,
				Optional: true,
				Description:    "xpath is: config.Groups.V_policy__statement.V_term.V_route__filter",
			},
			"limit__bandwidth": &schema.Schema{
				Type:    schema.TypeString,
				Optional: true,
				Description:    "xpath is: config.Groups.V_policy__statement.V_term.V_route__filter",
			},
			"install__nexthop": &schema.Schema{
				Type:    schema.TypeString,
				Optional: true,
				Description:    "xpath is: config.Groups.V_policy__statement.V_term.V_route__filter",
			},
			"aggregate__bandwidth": &schema.Schema{
				Type:    schema.TypeString,
				Optional: true,
				Description:    "xpath is: config.Groups.V_policy__statement.V_term.V_route__filter",
			},
			"no__route__localize": &schema.Schema{
				Type:    schema.TypeString,
				Optional: true,
				Description:    "xpath is: config.Groups.V_policy__statement.V_term.V_route__filter",
			},
			"multipath__resolve": &schema.Schema{
				Type:    schema.TypeString,
				Optional: true,
				Description:    "xpath is: config.Groups.V_policy__statement.V_term.V_route__filter",
			},
			"prefix__segment": &schema.Schema{
				Type:    schema.TypeString,
				Optional: true,
				Description:    "xpath is: config.Groups.V_policy__statement.V_term.V_route__filter",
			},
			"label__allocation": &schema.Schema{
				Type:    schema.TypeString,
				Optional: true,
				Description:    "xpath is: config.Groups.V_policy__statement.V_term.V_route__filter",
			},
			"origin": &schema.Schema{
				Type:    schema.TypeString,
				Optional: true,
				Description:    "xpath is: config.Groups.V_policy__statement.V_term.V_route__filter",
			},
			"mhop__bfd__port": &schema.Schema{
				Type:    schema.TypeString,
				Optional: true,
				Description:    "xpath is: config.Groups.V_policy__statement.V_term.V_route__filter",
			},
			"upto": &schema.Schema{
				Type:    schema.TypeString,
				Optional: true,
				Description:    "xpath is: config.Groups.V_policy__statement.V_term.V_route__filter",
			},
			"color": &schema.Schema{
				Type:    schema.TypeString,
				Optional: true,
				Description:    "xpath is: config.Groups.V_policy__statement.V_term.V_route__filter",
			},
			"priority": &schema.Schema{
				Type:    schema.TypeString,
				Optional: true,
				Description:    "xpath is: config.Groups.V_policy__statement.V_term.V_route__filter",
			},
			"exact": &schema.Schema{
				Type:    schema.TypeString,
				Optional: true,
				Description:    "xpath is: config.Groups.V_policy__statement.V_term.V_route__filter",
			},
			"forwarding__class": &schema.Schema{
				Type:    schema.TypeString,
				Optional: true,
				Description:    "xpath is: config.Groups.V_policy__statement.V_term.V_route__filter",
			},
			"resolution__map": &schema.Schema{
				Type:    schema.TypeString,
				Optional: true,
				Description:    "xpath is: config.Groups.V_policy__statement.V_term.V_route__filter",
			},
			"address": &schema.Schema{
				Type:    schema.TypeString,
				Optional: true,
				Description:    "xpath is: config.Groups.V_policy__statement.V_term.V_route__filter. IP address or hostname",
			},
		},
	}
}