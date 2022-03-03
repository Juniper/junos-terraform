
package main

import (
    "encoding/xml"
    "fmt"
    "github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)


// v_ is appended before every variable so it doesn't give any conflict
// with any keyword in golang. ex - interface is keyword in golang
type xmlSecurityAddress__BookAddress__SetAddressName struct {
	XMLName xml.Name `xml:"configuration"`
	Groups  struct {
		XMLName	xml.Name	`xml:"groups"`
		Name	string	`xml:"name"`
		V_address__book  struct {
			XMLName xml.Name `xml:"address-book"`
			V_name  *string  `xml:"name,omitempty"`
			V_address__set  struct {
				XMLName xml.Name `xml:"address-set"`
				V_name__1  *string  `xml:"name,omitempty"`
				V_address  struct {
					XMLName xml.Name `xml:"address"`
					V_name__2  *string  `xml:"name,omitempty"`
				} `xml:"address"`
			} `xml:"address-set"`
		} `xml:"security>address-book"`
	} `xml:"groups"`
	ApplyGroup string `xml:"apply-groups"`
}

// v_ is appended before every variable so it doesn't give any conflict
// with any keyword in golang. ex- interface is keyword in golang
func junosSecurityAddress__BookAddress__SetAddressNameCreate(d *schema.ResourceData, m interface{}) error {

	var err error
	client := m.(*ProviderConfig)

    id := d.Get("resource_name").(string)
     	V_name := d.Get("name").(string)
	V_name__1 := d.Get("name__1").(string)
	V_name__2 := d.Get("name__2").(string)


	config := xmlSecurityAddress__BookAddress__SetAddressName{}
	config.ApplyGroup = id
	config.Groups.Name = id
	config.Groups.V_address__book.V_name = &V_name
	config.Groups.V_address__book.V_address__set.V_name__1 = &V_name__1
	config.Groups.V_address__book.V_address__set.V_address.V_name__2 = &V_name__2

    err = client.SendTransaction("", config, false)
    check(err)
    
    d.SetId(fmt.Sprintf("%s_%s", client.Host, id))
    
	return junosSecurityAddress__BookAddress__SetAddressNameRead(d,m)
}

func junosSecurityAddress__BookAddress__SetAddressNameRead(d *schema.ResourceData, m interface{}) error {

	var err error
	client := m.(*ProviderConfig)

    id := d.Get("resource_name").(string)
    
	config := &xmlSecurityAddress__BookAddress__SetAddressName{}

	err = client.MarshalGroup(id, config)
	check(err)
 	d.Set("name", config.Groups.V_address__book.V_name)
	d.Set("name__1", config.Groups.V_address__book.V_address__set.V_name__1)
	d.Set("name__2", config.Groups.V_address__book.V_address__set.V_address.V_name__2)

	return nil
}

func junosSecurityAddress__BookAddress__SetAddressNameUpdate(d *schema.ResourceData, m interface{}) error {

	var err error
	client := m.(*ProviderConfig)

    id := d.Get("resource_name").(string)
     	V_name := d.Get("name").(string)
	V_name__1 := d.Get("name__1").(string)
	V_name__2 := d.Get("name__2").(string)


	config := xmlSecurityAddress__BookAddress__SetAddressName{}
	config.ApplyGroup = id
	config.Groups.Name = id
	config.Groups.V_address__book.V_name = &V_name
	config.Groups.V_address__book.V_address__set.V_name__1 = &V_name__1
	config.Groups.V_address__book.V_address__set.V_address.V_name__2 = &V_name__2

    err = client.SendTransaction(id, config, false)
    check(err)
    
	return junosSecurityAddress__BookAddress__SetAddressNameRead(d,m)
}

func junosSecurityAddress__BookAddress__SetAddressNameDelete(d *schema.ResourceData, m interface{}) error {

	var err error
	client := m.(*ProviderConfig)

    id := d.Get("resource_name").(string)
    
    _, err = client.DeleteConfigNoCommit(id)
    check(err)

    d.SetId("")
    
	return nil
}

func junosSecurityAddress__BookAddress__SetAddressName() *schema.Resource {
	return &schema.Resource{
		Create: junosSecurityAddress__BookAddress__SetAddressNameCreate,
		Read: junosSecurityAddress__BookAddress__SetAddressNameRead,
		Update: junosSecurityAddress__BookAddress__SetAddressNameUpdate,
		Delete: junosSecurityAddress__BookAddress__SetAddressNameDelete,

        Schema: map[string]*schema.Schema{
            "resource_name": &schema.Schema{
                Type:    schema.TypeString,
                Required: true,
            },
			"name": &schema.Schema{
				Type:    schema.TypeString,
				Optional: true,
				Description:    "xpath is: config.Groups.V_address__book",
			},
			"name__1": &schema.Schema{
				Type:    schema.TypeString,
				Optional: true,
				Description:    "xpath is: config.Groups.V_address__book.V_address__set",
			},
			"name__2": &schema.Schema{
				Type:    schema.TypeString,
				Optional: true,
				Description:    "xpath is: config.Groups.V_address__book.V_address__set.V_address. Security address name",
			},
		},
	}
}