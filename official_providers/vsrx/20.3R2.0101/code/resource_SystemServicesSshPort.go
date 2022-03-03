
package main

import (
    "encoding/xml"
    "fmt"
    "github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)


// v_ is appended before every variable so it doesn't give any conflict
// with any keyword in golang. ex - interface is keyword in golang
type xmlSystemServicesSshPort struct {
	XMLName xml.Name `xml:"configuration"`
	Groups  struct {
		XMLName	xml.Name	`xml:"groups"`
		Name	string	`xml:"name"`
		V_ssh  struct {
			XMLName xml.Name `xml:"ssh"`
			V_port  *string  `xml:"port,omitempty"`
		} `xml:"system>services>ssh"`
	} `xml:"groups"`
	ApplyGroup string `xml:"apply-groups"`
}

// v_ is appended before every variable so it doesn't give any conflict
// with any keyword in golang. ex- interface is keyword in golang
func junosSystemServicesSshPortCreate(d *schema.ResourceData, m interface{}) error {

	var err error
	client := m.(*ProviderConfig)

    id := d.Get("resource_name").(string)
     	V_port := d.Get("port").(string)


	config := xmlSystemServicesSshPort{}
	config.ApplyGroup = id
	config.Groups.Name = id
	config.Groups.V_ssh.V_port = &V_port

    err = client.SendTransaction("", config, false)
    check(err)
    
    d.SetId(fmt.Sprintf("%s_%s", client.Host, id))
    
	return junosSystemServicesSshPortRead(d,m)
}

func junosSystemServicesSshPortRead(d *schema.ResourceData, m interface{}) error {

	var err error
	client := m.(*ProviderConfig)

    id := d.Get("resource_name").(string)
    
	config := &xmlSystemServicesSshPort{}

	err = client.MarshalGroup(id, config)
	check(err)
 	d.Set("port", config.Groups.V_ssh.V_port)

	return nil
}

func junosSystemServicesSshPortUpdate(d *schema.ResourceData, m interface{}) error {

	var err error
	client := m.(*ProviderConfig)

    id := d.Get("resource_name").(string)
     	V_port := d.Get("port").(string)


	config := xmlSystemServicesSshPort{}
	config.ApplyGroup = id
	config.Groups.Name = id
	config.Groups.V_ssh.V_port = &V_port

    err = client.SendTransaction(id, config, false)
    check(err)
    
	return junosSystemServicesSshPortRead(d,m)
}

func junosSystemServicesSshPortDelete(d *schema.ResourceData, m interface{}) error {

	var err error
	client := m.(*ProviderConfig)

    id := d.Get("resource_name").(string)
    
    _, err = client.DeleteConfigNoCommit(id)
    check(err)

    d.SetId("")
    
	return nil
}

func junosSystemServicesSshPort() *schema.Resource {
	return &schema.Resource{
		Create: junosSystemServicesSshPortCreate,
		Read: junosSystemServicesSshPortRead,
		Update: junosSystemServicesSshPortUpdate,
		Delete: junosSystemServicesSshPortDelete,

        Schema: map[string]*schema.Schema{
            "resource_name": &schema.Schema{
                Type:    schema.TypeString,
                Required: true,
            },
			"port": &schema.Schema{
				Type:    schema.TypeString,
				Optional: true,
				Description:    "xpath is: config.Groups.V_ssh. Port number to accept incoming connections",
			},
		},
	}
}