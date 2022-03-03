
package main

import (
    "encoding/xml"
    "fmt"
    "github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)


// v_ is appended before every variable so it doesn't give any conflict
// with any keyword in golang. ex - interface is keyword in golang
type xmlSecurityNatDestinationPoolAddressIpaddr struct {
	XMLName xml.Name `xml:"configuration"`
	Groups  struct {
		XMLName	xml.Name	`xml:"groups"`
		Name	string	`xml:"name"`
		V_pool  struct {
			XMLName xml.Name `xml:"pool"`
			V_name  *string  `xml:"name,omitempty"`
			V_address  struct {
				XMLName xml.Name `xml:"address"`
				V_ipaddr  *string  `xml:"ipaddr,omitempty"`
			} `xml:"address"`
		} `xml:"security>nat>destination>pool"`
	} `xml:"groups"`
	ApplyGroup string `xml:"apply-groups"`
}

// v_ is appended before every variable so it doesn't give any conflict
// with any keyword in golang. ex- interface is keyword in golang
func junosSecurityNatDestinationPoolAddressIpaddrCreate(d *schema.ResourceData, m interface{}) error {

	var err error
	client := m.(*ProviderConfig)

    id := d.Get("resource_name").(string)
     	V_name := d.Get("name").(string)
	V_ipaddr := d.Get("ipaddr").(string)


	config := xmlSecurityNatDestinationPoolAddressIpaddr{}
	config.ApplyGroup = id
	config.Groups.Name = id
	config.Groups.V_pool.V_name = &V_name
	config.Groups.V_pool.V_address.V_ipaddr = &V_ipaddr

    err = client.SendTransaction("", config, false)
    check(err)
    
    d.SetId(fmt.Sprintf("%s_%s", client.Host, id))
    
	return junosSecurityNatDestinationPoolAddressIpaddrRead(d,m)
}

func junosSecurityNatDestinationPoolAddressIpaddrRead(d *schema.ResourceData, m interface{}) error {

	var err error
	client := m.(*ProviderConfig)

    id := d.Get("resource_name").(string)
    
	config := &xmlSecurityNatDestinationPoolAddressIpaddr{}

	err = client.MarshalGroup(id, config)
	check(err)
 	d.Set("name", config.Groups.V_pool.V_name)
	d.Set("ipaddr", config.Groups.V_pool.V_address.V_ipaddr)

	return nil
}

func junosSecurityNatDestinationPoolAddressIpaddrUpdate(d *schema.ResourceData, m interface{}) error {

	var err error
	client := m.(*ProviderConfig)

    id := d.Get("resource_name").(string)
     	V_name := d.Get("name").(string)
	V_ipaddr := d.Get("ipaddr").(string)


	config := xmlSecurityNatDestinationPoolAddressIpaddr{}
	config.ApplyGroup = id
	config.Groups.Name = id
	config.Groups.V_pool.V_name = &V_name
	config.Groups.V_pool.V_address.V_ipaddr = &V_ipaddr

    err = client.SendTransaction(id, config, false)
    check(err)
    
	return junosSecurityNatDestinationPoolAddressIpaddrRead(d,m)
}

func junosSecurityNatDestinationPoolAddressIpaddrDelete(d *schema.ResourceData, m interface{}) error {

	var err error
	client := m.(*ProviderConfig)

    id := d.Get("resource_name").(string)
    
    _, err = client.DeleteConfigNoCommit(id)
    check(err)

    d.SetId("")
    
	return nil
}

func junosSecurityNatDestinationPoolAddressIpaddr() *schema.Resource {
	return &schema.Resource{
		Create: junosSecurityNatDestinationPoolAddressIpaddrCreate,
		Read: junosSecurityNatDestinationPoolAddressIpaddrRead,
		Update: junosSecurityNatDestinationPoolAddressIpaddrUpdate,
		Delete: junosSecurityNatDestinationPoolAddressIpaddrDelete,

        Schema: map[string]*schema.Schema{
            "resource_name": &schema.Schema{
                Type:    schema.TypeString,
                Required: true,
            },
			"name": &schema.Schema{
				Type:    schema.TypeString,
				Optional: true,
				Description:    "xpath is: config.Groups.V_pool",
			},
			"ipaddr": &schema.Schema{
				Type:    schema.TypeString,
				Optional: true,
				Description:    "xpath is: config.Groups.V_pool.V_address. IPv4 or IPv6 address or address range",
			},
		},
	}
}