
package main

import (
    "encoding/xml"
    "fmt"
    "github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)


// v_ is appended before every variable so it doesn't give any conflict
// with any keyword in golang. ex - interface is keyword in golang
type xmlApplicationsApplicationDestination__Port struct {
	XMLName xml.Name `xml:"configuration"`
	Groups  struct {
		XMLName	xml.Name	`xml:"groups"`
		Name	string	`xml:"name"`
		V_application  struct {
			XMLName xml.Name `xml:"application"`
			V_name  *string  `xml:"name,omitempty"`
			V_destination__port  *string  `xml:"destination-port,omitempty"`
		} `xml:"applications>application"`
	} `xml:"groups"`
	ApplyGroup string `xml:"apply-groups"`
}

// v_ is appended before every variable so it doesn't give any conflict
// with any keyword in golang. ex- interface is keyword in golang
func junosApplicationsApplicationDestination__PortCreate(d *schema.ResourceData, m interface{}) error {

	var err error
	client := m.(*ProviderConfig)

    id := d.Get("resource_name").(string)
     	V_name := d.Get("name").(string)
	V_destination__port := d.Get("destination__port").(string)


	config := xmlApplicationsApplicationDestination__Port{}
	config.ApplyGroup = id
	config.Groups.Name = id
	config.Groups.V_application.V_name = &V_name
	config.Groups.V_application.V_destination__port = &V_destination__port

    err = client.SendTransaction("", config, false)
    check(err)
    
    d.SetId(fmt.Sprintf("%s_%s", client.Host, id))
    
	return junosApplicationsApplicationDestination__PortRead(d,m)
}

func junosApplicationsApplicationDestination__PortRead(d *schema.ResourceData, m interface{}) error {

	var err error
	client := m.(*ProviderConfig)

    id := d.Get("resource_name").(string)
    
	config := &xmlApplicationsApplicationDestination__Port{}

	err = client.MarshalGroup(id, config)
	check(err)
 	d.Set("name", config.Groups.V_application.V_name)
	d.Set("destination__port", config.Groups.V_application.V_destination__port)

	return nil
}

func junosApplicationsApplicationDestination__PortUpdate(d *schema.ResourceData, m interface{}) error {

	var err error
	client := m.(*ProviderConfig)

    id := d.Get("resource_name").(string)
     	V_name := d.Get("name").(string)
	V_destination__port := d.Get("destination__port").(string)


	config := xmlApplicationsApplicationDestination__Port{}
	config.ApplyGroup = id
	config.Groups.Name = id
	config.Groups.V_application.V_name = &V_name
	config.Groups.V_application.V_destination__port = &V_destination__port

    err = client.SendTransaction(id, config, false)
    check(err)
    
	return junosApplicationsApplicationDestination__PortRead(d,m)
}

func junosApplicationsApplicationDestination__PortDelete(d *schema.ResourceData, m interface{}) error {

	var err error
	client := m.(*ProviderConfig)

    id := d.Get("resource_name").(string)
    
    _, err = client.DeleteConfigNoCommit(id)
    check(err)

    d.SetId("")
    
	return nil
}

func junosApplicationsApplicationDestination__Port() *schema.Resource {
	return &schema.Resource{
		Create: junosApplicationsApplicationDestination__PortCreate,
		Read: junosApplicationsApplicationDestination__PortRead,
		Update: junosApplicationsApplicationDestination__PortUpdate,
		Delete: junosApplicationsApplicationDestination__PortDelete,

        Schema: map[string]*schema.Schema{
            "resource_name": &schema.Schema{
                Type:    schema.TypeString,
                Required: true,
            },
			"name": &schema.Schema{
				Type:    schema.TypeString,
				Optional: true,
				Description:    "xpath is: config.Groups.V_application",
			},
			"destination__port": &schema.Schema{
				Type:    schema.TypeString,
				Optional: true,
				Description:    "xpath is: config.Groups.V_application. Match TCP/UDP destination port",
			},
		},
	}
}