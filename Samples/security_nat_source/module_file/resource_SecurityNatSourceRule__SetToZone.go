
package main

import (
    "encoding/xml"
    "fmt"
    "github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)


// v_ is appended before every variable so it doesn't give any conflict
// with any keyword in golang. ex - interface is keyword in golang
type xmlSecurityNatSourceRule__SetToZone struct {
	XMLName xml.Name `xml:"configuration"`
	V_rule__set  struct {
		XMLName xml.Name `xml:"rule-set"`
		V_name  string  `xml:"name"`
		V_to  struct {
			XMLName xml.Name `xml:"to"`
			V_zone  string  `xml:"zone"`
		} `xml:"to"`
	} `xml:"security>nat>source>rule-set"`
}

// v_ is appended before every variable so it doesn't give any conflict
// with any keyword in golang. ex- interface is keyword in golang
func junosSecurityNatSourceRule__SetToZoneCreate(d *schema.ResourceData, m interface{}) error {

    pcfg := m.(*ProviderConfig)
    client, err := pcfg.Client()
    check(err)

    id := d.Get("resource_name").(string)
     	V_name := d.Get("name").(string)
	V_zone := d.Get("zone").(string)
	commit := true

	config := xmlSecurityNatSourceRule__SetToZone{}
	config.V_rule__set.V_name = V_name
	config.V_rule__set.V_to.V_zone = V_zone

    err = client.SendTransaction("", config, commit)
    check(err)
    
    d.SetId(fmt.Sprintf("%s_%s", pcfg.Cfg.Host, id))
    
    err = client.Close()
    check(err)
    
	return junosSecurityNatSourceRule__SetToZoneRead(d,m)
}

func junosSecurityNatSourceRule__SetToZoneRead(d *schema.ResourceData, m interface{}) error {

    pcfg := m.(*ProviderConfig)
    client, err := pcfg.Client()
    check(err)

    id := d.Get("resource_name").(string)
    
	config := &xmlSecurityNatSourceRule__SetToZone{}

	err = client.MarshalGroup(id, config)
	check(err)
 	d.Set("name", config.V_rule__set.V_name)
	d.Set("zone", config.V_rule__set.V_to.V_zone)

    err = client.Close()
    check(err)
    
	return nil
}

func junosSecurityNatSourceRule__SetToZoneUpdate(d *schema.ResourceData, m interface{}) error {

    pcfg := m.(*ProviderConfig)
    client, err := pcfg.Client()
    check(err)

    id := d.Get("resource_name").(string)
     	V_name := d.Get("name").(string)
	V_zone := d.Get("zone").(string)
	commit := true

	config := xmlSecurityNatSourceRule__SetToZone{}
	config.V_rule__set.V_name = V_name
	config.V_rule__set.V_to.V_zone = V_zone

    err = client.SendTransaction(id, config, commit)
    check(err)
    
    err = client.Close()
    check(err)
    
	return junosSecurityNatSourceRule__SetToZoneRead(d,m)
}

func junosSecurityNatSourceRule__SetToZoneDelete(d *schema.ResourceData, m interface{}) error {

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

func junosSecurityNatSourceRule__SetToZone() *schema.Resource {
	return &schema.Resource{
		Create: junosSecurityNatSourceRule__SetToZoneCreate,
		Read: junosSecurityNatSourceRule__SetToZoneRead,
		Update: junosSecurityNatSourceRule__SetToZoneUpdate,
		Delete: junosSecurityNatSourceRule__SetToZoneDelete,

        Schema: map[string]*schema.Schema{
            "resource_name": &schema.Schema{
                Type:    schema.TypeString,
                Required: true,
            },
			"name": &schema.Schema{
				Type:    schema.TypeString,
				Optional: true,
				Description:    "xpath is: config.V_rule__set",
			},
			"zone": &schema.Schema{
				Type:    schema.TypeString,
				Optional: true,
				Description:    "xpath is: config.V_rule__set.V_to. Destination zone list",
			},
		},
	}
}