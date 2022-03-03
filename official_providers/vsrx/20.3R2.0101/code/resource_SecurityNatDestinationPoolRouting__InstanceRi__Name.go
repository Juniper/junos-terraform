
package main

import (
    "encoding/xml"
    "fmt"
    "github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)


// v_ is appended before every variable so it doesn't give any conflict
// with any keyword in golang. ex - interface is keyword in golang
type xmlSecurityNatDestinationPoolRouting__InstanceRi__Name struct {
	XMLName xml.Name `xml:"configuration"`
	Groups  struct {
		XMLName	xml.Name	`xml:"groups"`
		Name	string	`xml:"name"`
		V_pool  struct {
			XMLName xml.Name `xml:"pool"`
			V_name  *string  `xml:"name,omitempty"`
			V_routing__instance  struct {
				XMLName xml.Name `xml:"routing-instance"`
				V_ri__name  *string  `xml:"ri-name,omitempty"`
			} `xml:"routing-instance"`
		} `xml:"security>nat>destination>pool"`
	} `xml:"groups"`
	ApplyGroup string `xml:"apply-groups"`
}

// v_ is appended before every variable so it doesn't give any conflict
// with any keyword in golang. ex- interface is keyword in golang
func junosSecurityNatDestinationPoolRouting__InstanceRi__NameCreate(d *schema.ResourceData, m interface{}) error {

	var err error
	client := m.(*ProviderConfig)

    id := d.Get("resource_name").(string)
     	V_name := d.Get("name").(string)
	V_ri__name := d.Get("ri__name").(string)


	config := xmlSecurityNatDestinationPoolRouting__InstanceRi__Name{}
	config.ApplyGroup = id
	config.Groups.Name = id
	config.Groups.V_pool.V_name = &V_name
	config.Groups.V_pool.V_routing__instance.V_ri__name = &V_ri__name

    err = client.SendTransaction("", config, false)
    check(err)
    
    d.SetId(fmt.Sprintf("%s_%s", client.Host, id))
    
	return junosSecurityNatDestinationPoolRouting__InstanceRi__NameRead(d,m)
}

func junosSecurityNatDestinationPoolRouting__InstanceRi__NameRead(d *schema.ResourceData, m interface{}) error {

	var err error
	client := m.(*ProviderConfig)

    id := d.Get("resource_name").(string)
    
	config := &xmlSecurityNatDestinationPoolRouting__InstanceRi__Name{}

	err = client.MarshalGroup(id, config)
	check(err)
 	d.Set("name", config.Groups.V_pool.V_name)
	d.Set("ri__name", config.Groups.V_pool.V_routing__instance.V_ri__name)

	return nil
}

func junosSecurityNatDestinationPoolRouting__InstanceRi__NameUpdate(d *schema.ResourceData, m interface{}) error {

	var err error
	client := m.(*ProviderConfig)

    id := d.Get("resource_name").(string)
     	V_name := d.Get("name").(string)
	V_ri__name := d.Get("ri__name").(string)


	config := xmlSecurityNatDestinationPoolRouting__InstanceRi__Name{}
	config.ApplyGroup = id
	config.Groups.Name = id
	config.Groups.V_pool.V_name = &V_name
	config.Groups.V_pool.V_routing__instance.V_ri__name = &V_ri__name

    err = client.SendTransaction(id, config, false)
    check(err)
    
	return junosSecurityNatDestinationPoolRouting__InstanceRi__NameRead(d,m)
}

func junosSecurityNatDestinationPoolRouting__InstanceRi__NameDelete(d *schema.ResourceData, m interface{}) error {

	var err error
	client := m.(*ProviderConfig)

    id := d.Get("resource_name").(string)
    
    _, err = client.DeleteConfigNoCommit(id)
    check(err)

    d.SetId("")
    
	return nil
}

func junosSecurityNatDestinationPoolRouting__InstanceRi__Name() *schema.Resource {
	return &schema.Resource{
		Create: junosSecurityNatDestinationPoolRouting__InstanceRi__NameCreate,
		Read: junosSecurityNatDestinationPoolRouting__InstanceRi__NameRead,
		Update: junosSecurityNatDestinationPoolRouting__InstanceRi__NameUpdate,
		Delete: junosSecurityNatDestinationPoolRouting__InstanceRi__NameDelete,

        Schema: map[string]*schema.Schema{
            "resource_name": &schema.Schema{
                Type:    schema.TypeString,
                Required: true,
            },
			"name": &schema.Schema{
				Type:    schema.TypeString,
				Optional: true,
				Description:    "xpath is: config.Groups.V_pool",
			},
			"ri__name": &schema.Schema{
				Type:    schema.TypeString,
				Optional: true,
				Description:    "xpath is: config.Groups.V_pool.V_routing__instance. Routing-instance name",
			},
		},
	}
}