
package main

import (
    "encoding/xml"
    "fmt"
    "github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)


// v_ is appended before every variable so it doesn't give any conflict
// with any keyword in golang. ex - interface is keyword in golang
type xmlSecurityNatDestinationRule__SetRuleDest__Nat__Rule__MatchDestination__PortName struct {
	XMLName xml.Name `xml:"configuration"`
	Groups  struct {
		XMLName	xml.Name	`xml:"groups"`
		Name	string	`xml:"name"`
		V_rule__set  struct {
			XMLName xml.Name `xml:"rule-set"`
			V_name  *string  `xml:"name,omitempty"`
			V_rule  struct {
				XMLName xml.Name `xml:"rule"`
				V_name__1  *string  `xml:"name,omitempty"`
				V_destination__port  struct {
					XMLName xml.Name `xml:"destination-port"`
					V_name__2  *string  `xml:"name,omitempty"`
				} `xml:"dest-nat-rule-match>destination-port"`
			} `xml:"rule"`
		} `xml:"security>nat>destination>rule-set"`
	} `xml:"groups"`
	ApplyGroup string `xml:"apply-groups"`
}

// v_ is appended before every variable so it doesn't give any conflict
// with any keyword in golang. ex- interface is keyword in golang
func junosSecurityNatDestinationRule__SetRuleDest__Nat__Rule__MatchDestination__PortNameCreate(d *schema.ResourceData, m interface{}) error {

	var err error
	client := m.(*ProviderConfig)

    id := d.Get("resource_name").(string)
     	V_name := d.Get("name").(string)
	V_name__1 := d.Get("name__1").(string)
	V_name__2 := d.Get("name__2").(string)


	config := xmlSecurityNatDestinationRule__SetRuleDest__Nat__Rule__MatchDestination__PortName{}
	config.ApplyGroup = id
	config.Groups.Name = id
	config.Groups.V_rule__set.V_name = &V_name
	config.Groups.V_rule__set.V_rule.V_name__1 = &V_name__1
	config.Groups.V_rule__set.V_rule.V_destination__port.V_name__2 = &V_name__2

    err = client.SendTransaction("", config, false)
    check(err)
    
    d.SetId(fmt.Sprintf("%s_%s", client.Host, id))
    
	return junosSecurityNatDestinationRule__SetRuleDest__Nat__Rule__MatchDestination__PortNameRead(d,m)
}

func junosSecurityNatDestinationRule__SetRuleDest__Nat__Rule__MatchDestination__PortNameRead(d *schema.ResourceData, m interface{}) error {

	var err error
	client := m.(*ProviderConfig)

    id := d.Get("resource_name").(string)
    
	config := &xmlSecurityNatDestinationRule__SetRuleDest__Nat__Rule__MatchDestination__PortName{}

	err = client.MarshalGroup(id, config)
	check(err)
 	d.Set("name", config.Groups.V_rule__set.V_name)
	d.Set("name__1", config.Groups.V_rule__set.V_rule.V_name__1)
	d.Set("name__2", config.Groups.V_rule__set.V_rule.V_destination__port.V_name__2)

	return nil
}

func junosSecurityNatDestinationRule__SetRuleDest__Nat__Rule__MatchDestination__PortNameUpdate(d *schema.ResourceData, m interface{}) error {

	var err error
	client := m.(*ProviderConfig)

    id := d.Get("resource_name").(string)
     	V_name := d.Get("name").(string)
	V_name__1 := d.Get("name__1").(string)
	V_name__2 := d.Get("name__2").(string)


	config := xmlSecurityNatDestinationRule__SetRuleDest__Nat__Rule__MatchDestination__PortName{}
	config.ApplyGroup = id
	config.Groups.Name = id
	config.Groups.V_rule__set.V_name = &V_name
	config.Groups.V_rule__set.V_rule.V_name__1 = &V_name__1
	config.Groups.V_rule__set.V_rule.V_destination__port.V_name__2 = &V_name__2

    err = client.SendTransaction(id, config, false)
    check(err)
    
	return junosSecurityNatDestinationRule__SetRuleDest__Nat__Rule__MatchDestination__PortNameRead(d,m)
}

func junosSecurityNatDestinationRule__SetRuleDest__Nat__Rule__MatchDestination__PortNameDelete(d *schema.ResourceData, m interface{}) error {

	var err error
	client := m.(*ProviderConfig)

    id := d.Get("resource_name").(string)
    
    _, err = client.DeleteConfigNoCommit(id)
    check(err)

    d.SetId("")
    
	return nil
}

func junosSecurityNatDestinationRule__SetRuleDest__Nat__Rule__MatchDestination__PortName() *schema.Resource {
	return &schema.Resource{
		Create: junosSecurityNatDestinationRule__SetRuleDest__Nat__Rule__MatchDestination__PortNameCreate,
		Read: junosSecurityNatDestinationRule__SetRuleDest__Nat__Rule__MatchDestination__PortNameRead,
		Update: junosSecurityNatDestinationRule__SetRuleDest__Nat__Rule__MatchDestination__PortNameUpdate,
		Delete: junosSecurityNatDestinationRule__SetRuleDest__Nat__Rule__MatchDestination__PortNameDelete,

        Schema: map[string]*schema.Schema{
            "resource_name": &schema.Schema{
                Type:    schema.TypeString,
                Required: true,
            },
			"name": &schema.Schema{
				Type:    schema.TypeString,
				Optional: true,
				Description:    "xpath is: config.Groups.V_rule__set",
			},
			"name__1": &schema.Schema{
				Type:    schema.TypeString,
				Optional: true,
				Description:    "xpath is: config.Groups.V_rule__set.V_rule",
			},
			"name__2": &schema.Schema{
				Type:    schema.TypeString,
				Optional: true,
				Description:    "xpath is: config.Groups.V_rule__set.V_rule.V_destination__port. Port or lower limit of port range",
			},
		},
	}
}