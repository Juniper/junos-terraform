
package main

import (
    "encoding/xml"
    "fmt"
    "github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)


// v_ is appended before every variable so it doesn't give any conflict
// with any keyword in golang. ex - interface is keyword in golang
type xmlRouting__InstancesInstanceRouting__OptionsInstance__Import struct {
	XMLName xml.Name `xml:"configuration"`
	Groups  struct {
		XMLName	xml.Name	`xml:"groups"`
		Name	string	`xml:"name"`
		V_instance  struct {
			XMLName xml.Name `xml:"instance"`
			V_name  *string  `xml:"name,omitempty"`
			V_routing__options  struct {
				XMLName xml.Name `xml:"routing-options"`
				V_instance__import  *string  `xml:"instance-import,omitempty"`
			} `xml:"routing-options"`
		} `xml:"routing-instances>instance"`
	} `xml:"groups"`
	ApplyGroup string `xml:"apply-groups"`
}

// v_ is appended before every variable so it doesn't give any conflict
// with any keyword in golang. ex- interface is keyword in golang
func junosRouting__InstancesInstanceRouting__OptionsInstance__ImportCreate(d *schema.ResourceData, m interface{}) error {

	var err error
	client := m.(*ProviderConfig)

    id := d.Get("resource_name").(string)
     	V_name := d.Get("name").(string)
	V_instance__import := d.Get("instance__import").(string)


	config := xmlRouting__InstancesInstanceRouting__OptionsInstance__Import{}
	config.ApplyGroup = id
	config.Groups.Name = id
	config.Groups.V_instance.V_name = &V_name
	config.Groups.V_instance.V_routing__options.V_instance__import = &V_instance__import

    err = client.SendTransaction("", config, false)
    check(err)
    
    d.SetId(fmt.Sprintf("%s_%s", client.Host, id))
    
	return junosRouting__InstancesInstanceRouting__OptionsInstance__ImportRead(d,m)
}

func junosRouting__InstancesInstanceRouting__OptionsInstance__ImportRead(d *schema.ResourceData, m interface{}) error {

	var err error
	client := m.(*ProviderConfig)

    id := d.Get("resource_name").(string)
    
	config := &xmlRouting__InstancesInstanceRouting__OptionsInstance__Import{}

	err = client.MarshalGroup(id, config)
	check(err)
 	d.Set("name", config.Groups.V_instance.V_name)
	d.Set("instance__import", config.Groups.V_instance.V_routing__options.V_instance__import)

	return nil
}

func junosRouting__InstancesInstanceRouting__OptionsInstance__ImportUpdate(d *schema.ResourceData, m interface{}) error {

	var err error
	client := m.(*ProviderConfig)

    id := d.Get("resource_name").(string)
     	V_name := d.Get("name").(string)
	V_instance__import := d.Get("instance__import").(string)


	config := xmlRouting__InstancesInstanceRouting__OptionsInstance__Import{}
	config.ApplyGroup = id
	config.Groups.Name = id
	config.Groups.V_instance.V_name = &V_name
	config.Groups.V_instance.V_routing__options.V_instance__import = &V_instance__import

    err = client.SendTransaction(id, config, false)
    check(err)
    
	return junosRouting__InstancesInstanceRouting__OptionsInstance__ImportRead(d,m)
}

func junosRouting__InstancesInstanceRouting__OptionsInstance__ImportDelete(d *schema.ResourceData, m interface{}) error {

	var err error
	client := m.(*ProviderConfig)

    id := d.Get("resource_name").(string)
    
    _, err = client.DeleteConfigNoCommit(id)
    check(err)

    d.SetId("")
    
	return nil
}

func junosRouting__InstancesInstanceRouting__OptionsInstance__Import() *schema.Resource {
	return &schema.Resource{
		Create: junosRouting__InstancesInstanceRouting__OptionsInstance__ImportCreate,
		Read: junosRouting__InstancesInstanceRouting__OptionsInstance__ImportRead,
		Update: junosRouting__InstancesInstanceRouting__OptionsInstance__ImportUpdate,
		Delete: junosRouting__InstancesInstanceRouting__OptionsInstance__ImportDelete,

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
			"instance__import": &schema.Schema{
				Type:    schema.TypeString,
				Optional: true,
				Description:    "xpath is: config.Groups.V_instance.V_routing__options. Import policy for instance RIBs",
			},
		},
	}
}