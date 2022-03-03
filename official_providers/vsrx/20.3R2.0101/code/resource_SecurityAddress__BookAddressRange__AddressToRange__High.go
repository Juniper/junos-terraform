
package main

import (
    "encoding/xml"
    "fmt"
    "github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)


// v_ is appended before every variable so it doesn't give any conflict
// with any keyword in golang. ex - interface is keyword in golang
type xmlSecurityAddress__BookAddressRange__AddressToRange__High struct {
	XMLName xml.Name `xml:"configuration"`
	Groups  struct {
		XMLName	xml.Name	`xml:"groups"`
		Name	string	`xml:"name"`
		V_address__book  struct {
			XMLName xml.Name `xml:"address-book"`
			V_name  *string  `xml:"name,omitempty"`
			V_address  struct {
				XMLName xml.Name `xml:"address"`
				V_name__1  *string  `xml:"name,omitempty"`
				V_range__address  struct {
					XMLName xml.Name `xml:"range-address"`
					V_name__2  *string  `xml:"name,omitempty"`
					V_to  struct {
						XMLName xml.Name `xml:"to"`
						V_range__high  *string  `xml:"range-high,omitempty"`
					} `xml:"to"`
				} `xml:"range-address"`
			} `xml:"address"`
		} `xml:"security>address-book"`
	} `xml:"groups"`
	ApplyGroup string `xml:"apply-groups"`
}

// v_ is appended before every variable so it doesn't give any conflict
// with any keyword in golang. ex- interface is keyword in golang
func junosSecurityAddress__BookAddressRange__AddressToRange__HighCreate(d *schema.ResourceData, m interface{}) error {

	var err error
	client := m.(*ProviderConfig)

    id := d.Get("resource_name").(string)
     	V_name := d.Get("name").(string)
	V_name__1 := d.Get("name__1").(string)
	V_name__2 := d.Get("name__2").(string)
	V_range__high := d.Get("range__high").(string)


	config := xmlSecurityAddress__BookAddressRange__AddressToRange__High{}
	config.ApplyGroup = id
	config.Groups.Name = id
	config.Groups.V_address__book.V_name = &V_name
	config.Groups.V_address__book.V_address.V_name__1 = &V_name__1
	config.Groups.V_address__book.V_address.V_range__address.V_name__2 = &V_name__2
	config.Groups.V_address__book.V_address.V_range__address.V_to.V_range__high = &V_range__high

    err = client.SendTransaction("", config, false)
    check(err)
    
    d.SetId(fmt.Sprintf("%s_%s", client.Host, id))
    
	return junosSecurityAddress__BookAddressRange__AddressToRange__HighRead(d,m)
}

func junosSecurityAddress__BookAddressRange__AddressToRange__HighRead(d *schema.ResourceData, m interface{}) error {

	var err error
	client := m.(*ProviderConfig)

    id := d.Get("resource_name").(string)
    
	config := &xmlSecurityAddress__BookAddressRange__AddressToRange__High{}

	err = client.MarshalGroup(id, config)
	check(err)
 	d.Set("name", config.Groups.V_address__book.V_name)
	d.Set("name__1", config.Groups.V_address__book.V_address.V_name__1)
	d.Set("name__2", config.Groups.V_address__book.V_address.V_range__address.V_name__2)
	d.Set("range__high", config.Groups.V_address__book.V_address.V_range__address.V_to.V_range__high)

	return nil
}

func junosSecurityAddress__BookAddressRange__AddressToRange__HighUpdate(d *schema.ResourceData, m interface{}) error {

	var err error
	client := m.(*ProviderConfig)

    id := d.Get("resource_name").(string)
     	V_name := d.Get("name").(string)
	V_name__1 := d.Get("name__1").(string)
	V_name__2 := d.Get("name__2").(string)
	V_range__high := d.Get("range__high").(string)


	config := xmlSecurityAddress__BookAddressRange__AddressToRange__High{}
	config.ApplyGroup = id
	config.Groups.Name = id
	config.Groups.V_address__book.V_name = &V_name
	config.Groups.V_address__book.V_address.V_name__1 = &V_name__1
	config.Groups.V_address__book.V_address.V_range__address.V_name__2 = &V_name__2
	config.Groups.V_address__book.V_address.V_range__address.V_to.V_range__high = &V_range__high

    err = client.SendTransaction(id, config, false)
    check(err)
    
	return junosSecurityAddress__BookAddressRange__AddressToRange__HighRead(d,m)
}

func junosSecurityAddress__BookAddressRange__AddressToRange__HighDelete(d *schema.ResourceData, m interface{}) error {

	var err error
	client := m.(*ProviderConfig)

    id := d.Get("resource_name").(string)
    
    _, err = client.DeleteConfigNoCommit(id)
    check(err)

    d.SetId("")
    
	return nil
}

func junosSecurityAddress__BookAddressRange__AddressToRange__High() *schema.Resource {
	return &schema.Resource{
		Create: junosSecurityAddress__BookAddressRange__AddressToRange__HighCreate,
		Read: junosSecurityAddress__BookAddressRange__AddressToRange__HighRead,
		Update: junosSecurityAddress__BookAddressRange__AddressToRange__HighUpdate,
		Delete: junosSecurityAddress__BookAddressRange__AddressToRange__HighDelete,

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
				Description:    "xpath is: config.Groups.V_address__book.V_address",
			},
			"name__2": &schema.Schema{
				Type:    schema.TypeString,
				Optional: true,
				Description:    "xpath is: config.Groups.V_address__book.V_address.V_range__address",
			},
			"range__high": &schema.Schema{
				Type:    schema.TypeString,
				Optional: true,
				Description:    "xpath is: config.Groups.V_address__book.V_address.V_range__address.V_to. Upper limit of address range",
			},
		},
	}
}