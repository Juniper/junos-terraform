
package main

import (
    "encoding/xml"
    "fmt"
    "github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)


// v_ is appended before every variable so it doesn't give any conflict
// with any keyword in golang. ex - interface is keyword in golang
type xmlApplicationsApplicationProtocol struct {
	XMLName xml.Name `xml:"configuration"`
	Groups  struct {
		XMLName	xml.Name	`xml:"groups"`
		Name	string	`xml:"name"`
		V_application  struct {
			XMLName xml.Name `xml:"application"`
			V_name  *string  `xml:"name,omitempty"`
			V_protocol  *string  `xml:"protocol,omitempty"`
		} `xml:"applications>application"`
	} `xml:"groups"`
	ApplyGroup string `xml:"apply-groups"`
}

// v_ is appended before every variable so it doesn't give any conflict
// with any keyword in golang. ex- interface is keyword in golang
func junosApplicationsApplicationProtocolCreate(d *schema.ResourceData, m interface{}) error {

	var err error
	client := m.(*ProviderConfig)

    id := d.Get("resource_name").(string)
     	V_name := d.Get("name").(string)
	V_protocol := d.Get("protocol").(string)


	config := xmlApplicationsApplicationProtocol{}
	config.ApplyGroup = id
	config.Groups.Name = id
	config.Groups.V_application.V_name = &V_name
	config.Groups.V_application.V_protocol = &V_protocol

    err = client.SendTransaction("", config, false)
    check(err)
    
    d.SetId(fmt.Sprintf("%s_%s", client.Host, id))
    
	return junosApplicationsApplicationProtocolRead(d,m)
}

func junosApplicationsApplicationProtocolRead(d *schema.ResourceData, m interface{}) error {

	var err error
	client := m.(*ProviderConfig)

    id := d.Get("resource_name").(string)
    
	config := &xmlApplicationsApplicationProtocol{}

	err = client.MarshalGroup(id, config)
	check(err)
 	d.Set("name", config.Groups.V_application.V_name)
	d.Set("protocol", config.Groups.V_application.V_protocol)

	return nil
}

func junosApplicationsApplicationProtocolUpdate(d *schema.ResourceData, m interface{}) error {

	var err error
	client := m.(*ProviderConfig)

    id := d.Get("resource_name").(string)
     	V_name := d.Get("name").(string)
	V_protocol := d.Get("protocol").(string)


	config := xmlApplicationsApplicationProtocol{}
	config.ApplyGroup = id
	config.Groups.Name = id
	config.Groups.V_application.V_name = &V_name
	config.Groups.V_application.V_protocol = &V_protocol

    err = client.SendTransaction(id, config, false)
    check(err)
    
	return junosApplicationsApplicationProtocolRead(d,m)
}

func junosApplicationsApplicationProtocolDelete(d *schema.ResourceData, m interface{}) error {

	var err error
	client := m.(*ProviderConfig)

    id := d.Get("resource_name").(string)
    
    _, err = client.DeleteConfigNoCommit(id)
    check(err)

    d.SetId("")
    
	return nil
}

func junosApplicationsApplicationProtocol() *schema.Resource {
	return &schema.Resource{
		Create: junosApplicationsApplicationProtocolCreate,
		Read: junosApplicationsApplicationProtocolRead,
		Update: junosApplicationsApplicationProtocolUpdate,
		Delete: junosApplicationsApplicationProtocolDelete,

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
			"protocol": &schema.Schema{
				Type:    schema.TypeString,
				Optional: true,
				Description:    "xpath is: config.Groups.V_application. Match IP protocol type",
			},
		},
	}
}