
package main

import (
    "encoding/xml"
    "fmt"
    "github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)


// v_ is appended before every variable so it doesn't give any conflict
// with any keyword in golang. ex - interface is keyword in golang
type xmlSystemServicesWeb__ManagementHttpsInterface struct {
	XMLName xml.Name `xml:"configuration"`
	Groups  struct {
		XMLName	xml.Name	`xml:"groups"`
		Name	string	`xml:"name"`
		V_https  struct {
			XMLName xml.Name `xml:"https"`
			V_interface  *string  `xml:"interface,omitempty"`
		} `xml:"system>services>web-management>https"`
	} `xml:"groups"`
	ApplyGroup string `xml:"apply-groups"`
}

// v_ is appended before every variable so it doesn't give any conflict
// with any keyword in golang. ex- interface is keyword in golang
func junosSystemServicesWeb__ManagementHttpsInterfaceCreate(d *schema.ResourceData, m interface{}) error {

	var err error
	client := m.(*ProviderConfig)

    id := d.Get("resource_name").(string)
     	V_interface := d.Get("interface").(string)


	config := xmlSystemServicesWeb__ManagementHttpsInterface{}
	config.ApplyGroup = id
	config.Groups.Name = id
	config.Groups.V_https.V_interface = &V_interface

    err = client.SendTransaction("", config, false)
    check(err)
    
    d.SetId(fmt.Sprintf("%s_%s", client.Host, id))
    
	return junosSystemServicesWeb__ManagementHttpsInterfaceRead(d,m)
}

func junosSystemServicesWeb__ManagementHttpsInterfaceRead(d *schema.ResourceData, m interface{}) error {

	var err error
	client := m.(*ProviderConfig)

    id := d.Get("resource_name").(string)
    
	config := &xmlSystemServicesWeb__ManagementHttpsInterface{}

	err = client.MarshalGroup(id, config)
	check(err)
 	d.Set("interface", config.Groups.V_https.V_interface)

	return nil
}

func junosSystemServicesWeb__ManagementHttpsInterfaceUpdate(d *schema.ResourceData, m interface{}) error {

	var err error
	client := m.(*ProviderConfig)

    id := d.Get("resource_name").(string)
     	V_interface := d.Get("interface").(string)


	config := xmlSystemServicesWeb__ManagementHttpsInterface{}
	config.ApplyGroup = id
	config.Groups.Name = id
	config.Groups.V_https.V_interface = &V_interface

    err = client.SendTransaction(id, config, false)
    check(err)
    
	return junosSystemServicesWeb__ManagementHttpsInterfaceRead(d,m)
}

func junosSystemServicesWeb__ManagementHttpsInterfaceDelete(d *schema.ResourceData, m interface{}) error {

	var err error
	client := m.(*ProviderConfig)

    id := d.Get("resource_name").(string)
    
    _, err = client.DeleteConfigNoCommit(id)
    check(err)

    d.SetId("")
    
	return nil
}

func junosSystemServicesWeb__ManagementHttpsInterface() *schema.Resource {
	return &schema.Resource{
		Create: junosSystemServicesWeb__ManagementHttpsInterfaceCreate,
		Read: junosSystemServicesWeb__ManagementHttpsInterfaceRead,
		Update: junosSystemServicesWeb__ManagementHttpsInterfaceUpdate,
		Delete: junosSystemServicesWeb__ManagementHttpsInterfaceDelete,

        Schema: map[string]*schema.Schema{
            "resource_name": &schema.Schema{
                Type:    schema.TypeString,
                Required: true,
            },
			"interface": &schema.Schema{
				Type:    schema.TypeString,
				Optional: true,
				Description:    "xpath is: config.Groups.V_https. Interfaces that accept HTTPS access",
			},
		},
	}
}