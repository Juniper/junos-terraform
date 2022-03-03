
package main

import (
    "encoding/xml"
    "fmt"
    "github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)


// v_ is appended before every variable so it doesn't give any conflict
// with any keyword in golang. ex - interface is keyword in golang
type xmlRouting__InstancesInstanceRouting__OptionsStaticRouteNext__Hop struct {
	XMLName xml.Name `xml:"configuration"`
	Groups  struct {
		XMLName	xml.Name	`xml:"groups"`
		Name	string	`xml:"name"`
		V_instance  struct {
			XMLName xml.Name `xml:"instance"`
			V_name  *string  `xml:"name,omitempty"`
			V_route  struct {
				XMLName xml.Name `xml:"route"`
				V_name__1  *string  `xml:"name,omitempty"`
				V_next__hop  *string  `xml:"next-hop,omitempty"`
			} `xml:"routing-options>static>route"`
		} `xml:"routing-instances>instance"`
	} `xml:"groups"`
	ApplyGroup string `xml:"apply-groups"`
}

// v_ is appended before every variable so it doesn't give any conflict
// with any keyword in golang. ex- interface is keyword in golang
func junosRouting__InstancesInstanceRouting__OptionsStaticRouteNext__HopCreate(d *schema.ResourceData, m interface{}) error {

	var err error
	client := m.(*ProviderConfig)

    id := d.Get("resource_name").(string)
     	V_name := d.Get("name").(string)
	V_name__1 := d.Get("name__1").(string)
	V_next__hop := d.Get("next__hop").(string)


	config := xmlRouting__InstancesInstanceRouting__OptionsStaticRouteNext__Hop{}
	config.ApplyGroup = id
	config.Groups.Name = id
	config.Groups.V_instance.V_name = &V_name
	config.Groups.V_instance.V_route.V_name__1 = &V_name__1
	config.Groups.V_instance.V_route.V_next__hop = &V_next__hop

    err = client.SendTransaction("", config, false)
    check(err)
    
    d.SetId(fmt.Sprintf("%s_%s", client.Host, id))
    
	return junosRouting__InstancesInstanceRouting__OptionsStaticRouteNext__HopRead(d,m)
}

func junosRouting__InstancesInstanceRouting__OptionsStaticRouteNext__HopRead(d *schema.ResourceData, m interface{}) error {

	var err error
	client := m.(*ProviderConfig)

    id := d.Get("resource_name").(string)
    
	config := &xmlRouting__InstancesInstanceRouting__OptionsStaticRouteNext__Hop{}

	err = client.MarshalGroup(id, config)
	check(err)
 	d.Set("name", config.Groups.V_instance.V_name)
	d.Set("name__1", config.Groups.V_instance.V_route.V_name__1)
	d.Set("next__hop", config.Groups.V_instance.V_route.V_next__hop)

	return nil
}

func junosRouting__InstancesInstanceRouting__OptionsStaticRouteNext__HopUpdate(d *schema.ResourceData, m interface{}) error {

	var err error
	client := m.(*ProviderConfig)

    id := d.Get("resource_name").(string)
     	V_name := d.Get("name").(string)
	V_name__1 := d.Get("name__1").(string)
	V_next__hop := d.Get("next__hop").(string)


	config := xmlRouting__InstancesInstanceRouting__OptionsStaticRouteNext__Hop{}
	config.ApplyGroup = id
	config.Groups.Name = id
	config.Groups.V_instance.V_name = &V_name
	config.Groups.V_instance.V_route.V_name__1 = &V_name__1
	config.Groups.V_instance.V_route.V_next__hop = &V_next__hop

    err = client.SendTransaction(id, config, false)
    check(err)
    
	return junosRouting__InstancesInstanceRouting__OptionsStaticRouteNext__HopRead(d,m)
}

func junosRouting__InstancesInstanceRouting__OptionsStaticRouteNext__HopDelete(d *schema.ResourceData, m interface{}) error {

	var err error
	client := m.(*ProviderConfig)

    id := d.Get("resource_name").(string)
    
    _, err = client.DeleteConfigNoCommit(id)
    check(err)

    d.SetId("")
    
	return nil
}

func junosRouting__InstancesInstanceRouting__OptionsStaticRouteNext__Hop() *schema.Resource {
	return &schema.Resource{
		Create: junosRouting__InstancesInstanceRouting__OptionsStaticRouteNext__HopCreate,
		Read: junosRouting__InstancesInstanceRouting__OptionsStaticRouteNext__HopRead,
		Update: junosRouting__InstancesInstanceRouting__OptionsStaticRouteNext__HopUpdate,
		Delete: junosRouting__InstancesInstanceRouting__OptionsStaticRouteNext__HopDelete,

        Schema: map[string]*schema.Schema{
            "resource_name": &schema.Schema{
                Type:    schema.TypeString,
                Required: true,
            },
			"name": &schema.Schema{
				Type:    schema.TypeString,
				Optional: true,
				Description:    "xpath is: config.Groups.V_instance",
			},
			"name__1": &schema.Schema{
				Type:    schema.TypeString,
				Optional: true,
				Description:    "xpath is: config.Groups.V_instance.V_route",
			},
			"next__hop": &schema.Schema{
				Type:    schema.TypeString,
				Optional: true,
				Description:    "xpath is: config.Groups.V_instance.V_route. Next hop to destination",
			},
		},
	}
}