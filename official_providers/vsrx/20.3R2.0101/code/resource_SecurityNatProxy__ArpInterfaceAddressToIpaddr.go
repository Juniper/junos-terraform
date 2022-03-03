
package main

import (
    "encoding/xml"
    "fmt"
    "github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)


// v_ is appended before every variable so it doesn't give any conflict
// with any keyword in golang. ex - interface is keyword in golang
type xmlSecurityNatProxy__ArpInterfaceAddressToIpaddr struct {
	XMLName xml.Name `xml:"configuration"`
	Groups  struct {
		XMLName	xml.Name	`xml:"groups"`
		Name	string	`xml:"name"`
		V_interface  struct {
			XMLName xml.Name `xml:"interface"`
			V_name  *string  `xml:"name,omitempty"`
			V_address  struct {
				XMLName xml.Name `xml:"address"`
				V_name__1  *string  `xml:"name,omitempty"`
				V_to  struct {
					XMLName xml.Name `xml:"to"`
					V_ipaddr  *string  `xml:"ipaddr,omitempty"`
				} `xml:"to"`
			} `xml:"address"`
		} `xml:"security>nat>proxy-arp>interface"`
	} `xml:"groups"`
	ApplyGroup string `xml:"apply-groups"`
}

// v_ is appended before every variable so it doesn't give any conflict
// with any keyword in golang. ex- interface is keyword in golang
func junosSecurityNatProxy__ArpInterfaceAddressToIpaddrCreate(d *schema.ResourceData, m interface{}) error {

	var err error
	client := m.(*ProviderConfig)

    id := d.Get("resource_name").(string)
     	V_name := d.Get("name").(string)
	V_name__1 := d.Get("name__1").(string)
	V_ipaddr := d.Get("ipaddr").(string)


	config := xmlSecurityNatProxy__ArpInterfaceAddressToIpaddr{}
	config.ApplyGroup = id
	config.Groups.Name = id
	config.Groups.V_interface.V_name = &V_name
	config.Groups.V_interface.V_address.V_name__1 = &V_name__1
	config.Groups.V_interface.V_address.V_to.V_ipaddr = &V_ipaddr

    err = client.SendTransaction("", config, false)
    check(err)
    
    d.SetId(fmt.Sprintf("%s_%s", client.Host, id))
    
	return junosSecurityNatProxy__ArpInterfaceAddressToIpaddrRead(d,m)
}

func junosSecurityNatProxy__ArpInterfaceAddressToIpaddrRead(d *schema.ResourceData, m interface{}) error {

	var err error
	client := m.(*ProviderConfig)

    id := d.Get("resource_name").(string)
    
	config := &xmlSecurityNatProxy__ArpInterfaceAddressToIpaddr{}

	err = client.MarshalGroup(id, config)
	check(err)
 	d.Set("name", config.Groups.V_interface.V_name)
	d.Set("name__1", config.Groups.V_interface.V_address.V_name__1)
	d.Set("ipaddr", config.Groups.V_interface.V_address.V_to.V_ipaddr)

	return nil
}

func junosSecurityNatProxy__ArpInterfaceAddressToIpaddrUpdate(d *schema.ResourceData, m interface{}) error {

	var err error
	client := m.(*ProviderConfig)

    id := d.Get("resource_name").(string)
     	V_name := d.Get("name").(string)
	V_name__1 := d.Get("name__1").(string)
	V_ipaddr := d.Get("ipaddr").(string)


	config := xmlSecurityNatProxy__ArpInterfaceAddressToIpaddr{}
	config.ApplyGroup = id
	config.Groups.Name = id
	config.Groups.V_interface.V_name = &V_name
	config.Groups.V_interface.V_address.V_name__1 = &V_name__1
	config.Groups.V_interface.V_address.V_to.V_ipaddr = &V_ipaddr

    err = client.SendTransaction(id, config, false)
    check(err)
    
	return junosSecurityNatProxy__ArpInterfaceAddressToIpaddrRead(d,m)
}

func junosSecurityNatProxy__ArpInterfaceAddressToIpaddrDelete(d *schema.ResourceData, m interface{}) error {

	var err error
	client := m.(*ProviderConfig)

    id := d.Get("resource_name").(string)
    
    _, err = client.DeleteConfigNoCommit(id)
    check(err)

    d.SetId("")
    
	return nil
}

func junosSecurityNatProxy__ArpInterfaceAddressToIpaddr() *schema.Resource {
	return &schema.Resource{
		Create: junosSecurityNatProxy__ArpInterfaceAddressToIpaddrCreate,
		Read: junosSecurityNatProxy__ArpInterfaceAddressToIpaddrRead,
		Update: junosSecurityNatProxy__ArpInterfaceAddressToIpaddrUpdate,
		Delete: junosSecurityNatProxy__ArpInterfaceAddressToIpaddrDelete,

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
				Description:    "xpath is: config.Groups.V_interface.V_address",
			},
			"ipaddr": &schema.Schema{
				Type:    schema.TypeString,
				Optional: true,
				Description:    "xpath is: config.Groups.V_interface.V_address.V_to. Upper limit of address range",
			},
		},
	}
}