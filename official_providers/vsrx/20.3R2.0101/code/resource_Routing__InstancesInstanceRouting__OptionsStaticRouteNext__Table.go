
package main

import (
    "encoding/xml"
    "fmt"
    "github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)


// v_ is appended before every variable so it doesn't give any conflict
// with any keyword in golang. ex - interface is keyword in golang
type xmlRouting__InstancesInstanceRouting__OptionsStaticRouteNext__Table struct {
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
				V_next__table  *string  `xml:"next-table,omitempty"`
			} `xml:"routing-options>static>route"`
		} `xml:"routing-instances>instance"`
	} `xml:"groups"`
	ApplyGroup string `xml:"apply-groups"`
}

// v_ is appended before every variable so it doesn't give any conflict
// with any keyword in golang. ex- interface is keyword in golang
func junosRouting__InstancesInstanceRouting__OptionsStaticRouteNext__TableCreate(d *schema.ResourceData, m interface{}) error {

	var err error
	client := m.(*ProviderConfig)

    id := d.Get("resource_name").(string)
     	V_name := d.Get("name").(string)
	V_name__1 := d.Get("name__1").(string)
	V_next__table := d.Get("next__table").(string)


	config := xmlRouting__InstancesInstanceRouting__OptionsStaticRouteNext__Table{}
	config.ApplyGroup = id
	config.Groups.Name = id
	config.Groups.V_instance.V_name = &V_name
	config.Groups.V_instance.V_route.V_name__1 = &V_name__1
	config.Groups.V_instance.V_route.V_next__table = &V_next__table

    err = client.SendTransaction("", config, false)
    check(err)
    
    d.SetId(fmt.Sprintf("%s_%s", client.Host, id))
    
	return junosRouting__InstancesInstanceRouting__OptionsStaticRouteNext__TableRead(d,m)
}

func junosRouting__InstancesInstanceRouting__OptionsStaticRouteNext__TableRead(d *schema.ResourceData, m interface{}) error {

	var err error
	client := m.(*ProviderConfig)

    id := d.Get("resource_name").(string)
    
	config := &xmlRouting__InstancesInstanceRouting__OptionsStaticRouteNext__Table{}

	err = client.MarshalGroup(id, config)
	check(err)
 	d.Set("name", config.Groups.V_instance.V_name)
	d.Set("name__1", config.Groups.V_instance.V_route.V_name__1)
	d.Set("next__table", config.Groups.V_instance.V_route.V_next__table)

	return nil
}

func junosRouting__InstancesInstanceRouting__OptionsStaticRouteNext__TableUpdate(d *schema.ResourceData, m interface{}) error {

	var err error
	client := m.(*ProviderConfig)

    id := d.Get("resource_name").(string)
     	V_name := d.Get("name").(string)
	V_name__1 := d.Get("name__1").(string)
	V_next__table := d.Get("next__table").(string)


	config := xmlRouting__InstancesInstanceRouting__OptionsStaticRouteNext__Table{}
	config.ApplyGroup = id
	config.Groups.Name = id
	config.Groups.V_instance.V_name = &V_name
	config.Groups.V_instance.V_route.V_name__1 = &V_name__1
	config.Groups.V_instance.V_route.V_next__table = &V_next__table

    err = client.SendTransaction(id, config, false)
    check(err)
    
	return junosRouting__InstancesInstanceRouting__OptionsStaticRouteNext__TableRead(d,m)
}

func junosRouting__InstancesInstanceRouting__OptionsStaticRouteNext__TableDelete(d *schema.ResourceData, m interface{}) error {

	var err error
	client := m.(*ProviderConfig)

    id := d.Get("resource_name").(string)
    
    _, err = client.DeleteConfigNoCommit(id)
    check(err)

    d.SetId("")
    
	return nil
}

func junosRouting__InstancesInstanceRouting__OptionsStaticRouteNext__Table() *schema.Resource {
	return &schema.Resource{
		Create: junosRouting__InstancesInstanceRouting__OptionsStaticRouteNext__TableCreate,
		Read: junosRouting__InstancesInstanceRouting__OptionsStaticRouteNext__TableRead,
		Update: junosRouting__InstancesInstanceRouting__OptionsStaticRouteNext__TableUpdate,
		Delete: junosRouting__InstancesInstanceRouting__OptionsStaticRouteNext__TableDelete,

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
			"next__table": &schema.Schema{
				Type:    schema.TypeString,
				Optional: true,
				Description:    "xpath is: config.Groups.V_instance.V_route. Next hop to another table",
			},
		},
	}
}