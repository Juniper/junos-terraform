
package main

import (
    "encoding/xml"
    "fmt"
    "github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)


// v_ is appended before every variable so it doesn't give any conflict
// with any keyword in golang. ex - interface is keyword in golang
type xmlInterfacesInterfaceUnitFamilyInetAddressName struct {
	XMLName xml.Name `xml:"configuration"`
	Groups  struct {
		XMLName	xml.Name	`xml:"groups"`
		Name	string	`xml:"name"`
		V_interface  struct {
			XMLName xml.Name `xml:"interface"`
			V_name  *string  `xml:"name,omitempty"`
			V_unit  struct {
				XMLName xml.Name `xml:"unit"`
				V_name__1  *string  `xml:"name,omitempty"`
				V_address  struct {
					XMLName xml.Name `xml:"address"`
					V_name__2  *string  `xml:"name,omitempty"`
				} `xml:"family>inet>address"`
			} `xml:"unit"`
		} `xml:"interfaces>interface"`
	} `xml:"groups"`
	ApplyGroup string `xml:"apply-groups"`
}

// v_ is appended before every variable so it doesn't give any conflict
// with any keyword in golang. ex- interface is keyword in golang
func junosInterfacesInterfaceUnitFamilyInetAddressNameCreate(d *schema.ResourceData, m interface{}) error {

	var err error
	client := m.(*ProviderConfig)

    id := d.Get("resource_name").(string)
     	V_name := d.Get("name").(string)
	V_name__1 := d.Get("name__1").(string)
	V_name__2 := d.Get("name__2").(string)


	config := xmlInterfacesInterfaceUnitFamilyInetAddressName{}
	config.ApplyGroup = id
	config.Groups.Name = id
	config.Groups.V_interface.V_name = &V_name
	config.Groups.V_interface.V_unit.V_name__1 = &V_name__1
	config.Groups.V_interface.V_unit.V_address.V_name__2 = &V_name__2

    err = client.SendTransaction("", config, false)
    check(err)
    
    d.SetId(fmt.Sprintf("%s_%s", client.Host, id))
    
	return junosInterfacesInterfaceUnitFamilyInetAddressNameRead(d,m)
}

func junosInterfacesInterfaceUnitFamilyInetAddressNameRead(d *schema.ResourceData, m interface{}) error {

	var err error
	client := m.(*ProviderConfig)

    id := d.Get("resource_name").(string)
    
	config := &xmlInterfacesInterfaceUnitFamilyInetAddressName{}

	err = client.MarshalGroup(id, config)
	check(err)
 	d.Set("name", config.Groups.V_interface.V_name)
	d.Set("name__1", config.Groups.V_interface.V_unit.V_name__1)
	d.Set("name__2", config.Groups.V_interface.V_unit.V_address.V_name__2)

	return nil
}

func junosInterfacesInterfaceUnitFamilyInetAddressNameUpdate(d *schema.ResourceData, m interface{}) error {

	var err error
	client := m.(*ProviderConfig)

    id := d.Get("resource_name").(string)
     	V_name := d.Get("name").(string)
	V_name__1 := d.Get("name__1").(string)
	V_name__2 := d.Get("name__2").(string)


	config := xmlInterfacesInterfaceUnitFamilyInetAddressName{}
	config.ApplyGroup = id
	config.Groups.Name = id
	config.Groups.V_interface.V_name = &V_name
	config.Groups.V_interface.V_unit.V_name__1 = &V_name__1
	config.Groups.V_interface.V_unit.V_address.V_name__2 = &V_name__2

    err = client.SendTransaction(id, config, false)
    check(err)
    
	return junosInterfacesInterfaceUnitFamilyInetAddressNameRead(d,m)
}

func junosInterfacesInterfaceUnitFamilyInetAddressNameDelete(d *schema.ResourceData, m interface{}) error {

	var err error
	client := m.(*ProviderConfig)

    id := d.Get("resource_name").(string)
    
    _, err = client.DeleteConfigNoCommit(id)
    check(err)

    d.SetId("")
    
	return nil
}

func junosInterfacesInterfaceUnitFamilyInetAddressName() *schema.Resource {
	return &schema.Resource{
		Create: junosInterfacesInterfaceUnitFamilyInetAddressNameCreate,
		Read: junosInterfacesInterfaceUnitFamilyInetAddressNameRead,
		Update: junosInterfacesInterfaceUnitFamilyInetAddressNameUpdate,
		Delete: junosInterfacesInterfaceUnitFamilyInetAddressNameDelete,

        Schema: map[string]*schema.Schema{
            "resource_name": &schema.Schema{
                Type:    schema.TypeString,
                Required: true,
            },
			"name": &schema.Schema{
				Type:    schema.TypeString,
				Optional: true,
				Description:    "xpath is: config.Groups.V_interface",
			},
			"name__1": &schema.Schema{
				Type:    schema.TypeString,
				Optional: true,
				Description:    "xpath is: config.Groups.V_interface.V_unit",
			},
			"name__2": &schema.Schema{
				Type:    schema.TypeString,
				Optional: true,
				Description:    "xpath is: config.Groups.V_interface.V_unit.V_address. Interface address/destination prefix",
			},
		},
	}
}