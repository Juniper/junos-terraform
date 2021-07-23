
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
	V_address__book  struct {
		XMLName xml.Name `xml:"address-book"`
		V_name  string  `xml:"name"`
		V_address__set  struct {
			XMLName xml.Name `xml:"address-set"`
			V_name__1  string  `xml:"name"`
			V_address  struct {
				XMLName xml.Name `xml:"address"`
				V_name__2  string  `xml:"name"`
			} `xml:"address"`
		} `xml:"address-set"`
	} `xml:"security>address-book"`
}

// v_ is appended before every variable so it doesn't give any conflict
// with any keyword in golang. ex- interface is keyword in golang
func junosSecurityAddress__BookAddress__SetAddressNameCreate(d *schema.ResourceData, m interface{}) error {

    pcfg := m.(*ProviderConfig)
    client, err := pcfg.Client()
    check(err)

    id := d.Get("resource_name").(string)
     	V_name := d.Get("name").(string)
	V_name__1 := d.Get("name__1").(string)
	V_name__2 := d.Get("name__2").(string)
	commit := true

	config := xmlSecurityAddress__BookAddress__SetAddressName{}
	config.V_address__book.V_name = V_name
	config.V_address__book.V_address__set.V_name__1 = V_name__1
	config.V_address__book.V_address__set.V_address.V_name__2 = V_name__2

    err = client.SendTransaction("", config, commit)
    check(err)
    
    d.SetId(fmt.Sprintf("%s_%s", pcfg.Cfg.Host, id))
    
    err = client.Close()
    check(err)
    
	return junosSecurityAddress__BookAddress__SetAddressNameRead(d,m)
}

func junosSecurityAddress__BookAddress__SetAddressNameRead(d *schema.ResourceData, m interface{}) error {

    pcfg := m.(*ProviderConfig)
    client, err := pcfg.Client()
    check(err)

    id := d.Get("resource_name").(string)
    
	config := &xmlSecurityAddress__BookAddress__SetAddressName{}

	err = client.MarshalGroup(id, config)
	check(err)
 	d.Set("name", config.V_address__book.V_name)
	d.Set("name__1", config.V_address__book.V_address__set.V_name__1)
	d.Set("name__2", config.V_address__book.V_address__set.V_address.V_name__2)

    err = client.Close()
    check(err)
    
	return nil
}

func junosSecurityAddress__BookAddress__SetAddressNameUpdate(d *schema.ResourceData, m interface{}) error {

    pcfg := m.(*ProviderConfig)
    client, err := pcfg.Client()
    check(err)

    id := d.Get("resource_name").(string)
     	V_name := d.Get("name").(string)
	V_name__1 := d.Get("name__1").(string)
	V_name__2 := d.Get("name__2").(string)
	commit := true

	config := xmlSecurityAddress__BookAddress__SetAddressName{}
	config.V_address__book.V_name = V_name
	config.V_address__book.V_address__set.V_name__1 = V_name__1
	config.V_address__book.V_address__set.V_address.V_name__2 = V_name__2

    err = client.SendTransaction(id, config, commit)
    check(err)
    
    err = client.Close()
    check(err)
    
	return junosSecurityAddress__BookAddress__SetAddressNameRead(d,m)
}

func junosSecurityAddress__BookAddress__SetAddressNameDelete(d *schema.ResourceData, m interface{}) error {

    pcfg := m.(*ProviderConfig)
    client, err := pcfg.Client()
    check(err)

    id := d.Get("resource_name").(string)
    
    _, err = client.DeleteConfig(id)
    check(err)

    d.SetId("")
    
    err = client.Close()
    check(err)
    
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
				Description:    "xpath is: config.V_address__book",
			},
			"name__1": &schema.Schema{
				Type:    schema.TypeString,
				Optional: true,
				Description:    "xpath is: config.V_address__book.V_address__set",
			},
			"name__2": &schema.Schema{
				Type:    schema.TypeString,
				Optional: true,
				Description:    "xpath is: config.V_address__book.V_address__set.V_address. Security address name",
			},
		},
	}
}