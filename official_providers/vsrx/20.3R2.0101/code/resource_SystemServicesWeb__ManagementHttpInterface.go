
package main

import (
    "encoding/xml"
    "fmt"
    "github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)


// v_ is appended before every variable so it doesn't give any conflict
// with any keyword in golang. ex - interface is keyword in golang
type xmlSystemServicesWeb__ManagementHttpInterface struct {
	XMLName xml.Name `xml:"configuration"`
	Groups  struct {
		XMLName	xml.Name	`xml:"groups"`
		Name	string	`xml:"name"`
		V_http  struct {
			XMLName xml.Name `xml:"http"`
			V_interface  *string  `xml:"interface,omitempty"`
		} `xml:"system>services>web-management>http"`
	} `xml:"groups"`
	ApplyGroup string `xml:"apply-groups"`
}

// v_ is appended before every variable so it doesn't give any conflict
// with any keyword in golang. ex- interface is keyword in golang
func junosSystemServicesWeb__ManagementHttpInterfaceCreate(d *schema.ResourceData, m interface{}) error {

	var err error
	client := m.(*ProviderConfig)

    id := d.Get("resource_name").(string)
     	V_interface := d.Get("interface").(string)


	config := xmlSystemServicesWeb__ManagementHttpInterface{}
	config.ApplyGroup = id
	config.Groups.Name = id
	config.Groups.V_http.V_interface = &V_interface

    err = client.SendTransaction("", config, false)
    check(err)
    
    d.SetId(fmt.Sprintf("%s_%s", client.Host, id))
    
	return junosSystemServicesWeb__ManagementHttpInterfaceRead(d,m)
}

func junosSystemServicesWeb__ManagementHttpInterfaceRead(d *schema.ResourceData, m interface{}) error {

	var err error
	client := m.(*ProviderConfig)

    id := d.Get("resource_name").(string)
    
	config := &xmlSystemServicesWeb__ManagementHttpInterface{}

	err = client.MarshalGroup(id, config)
	check(err)
 	d.Set("interface", config.Groups.V_http.V_interface)

	return nil
}

func junosSystemServicesWeb__ManagementHttpInterfaceUpdate(d *schema.ResourceData, m interface{}) error {

	var err error
	client := m.(*ProviderConfig)

    id := d.Get("resource_name").(string)
     	V_interface := d.Get("interface").(string)


	config := xmlSystemServicesWeb__ManagementHttpInterface{}
	config.ApplyGroup = id
	config.Groups.Name = id
	config.Groups.V_http.V_interface = &V_interface

    err = client.SendTransaction(id, config, false)
    check(err)
    
	return junosSystemServicesWeb__ManagementHttpInterfaceRead(d,m)
}

func junosSystemServicesWeb__ManagementHttpInterfaceDelete(d *schema.ResourceData, m interface{}) error {

	var err error
	client := m.(*ProviderConfig)

    id := d.Get("resource_name").(string)
    
    _, err = client.DeleteConfigNoCommit(id)
    check(err)

    d.SetId("")
    
	return nil
}

func junosSystemServicesWeb__ManagementHttpInterface() *schema.Resource {
	return &schema.Resource{
		Create: junosSystemServicesWeb__ManagementHttpInterfaceCreate,
		Read: junosSystemServicesWeb__ManagementHttpInterfaceRead,
		Update: junosSystemServicesWeb__ManagementHttpInterfaceUpdate,
		Delete: junosSystemServicesWeb__ManagementHttpInterfaceDelete,

        Schema: map[string]*schema.Schema{
            "resource_name": &schema.Schema{
                Type:    schema.TypeString,
                Required: true,
            },
			"interface": &schema.Schema{
				Type:    schema.TypeString,
				Optional: true,
				Description:    "xpath is: config.Groups.V_http. Interfaces that accept HTTP access",
			},
		},
	}
}