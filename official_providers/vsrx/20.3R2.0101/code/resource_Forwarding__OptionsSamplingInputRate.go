
package main

import (
    "encoding/xml"
    "fmt"
    "github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)


// v_ is appended before every variable so it doesn't give any conflict
// with any keyword in golang. ex - interface is keyword in golang
type xmlForwarding__OptionsSamplingInputRate struct {
	XMLName xml.Name `xml:"configuration"`
	Groups  struct {
		XMLName	xml.Name	`xml:"groups"`
		Name	string	`xml:"name"`
		V_input  struct {
			XMLName xml.Name `xml:"input"`
			V_rate  *string  `xml:"rate,omitempty"`
		} `xml:"forwarding-options>sampling>input"`
	} `xml:"groups"`
	ApplyGroup string `xml:"apply-groups"`
}

// v_ is appended before every variable so it doesn't give any conflict
// with any keyword in golang. ex- interface is keyword in golang
func junosForwarding__OptionsSamplingInputRateCreate(d *schema.ResourceData, m interface{}) error {

	var err error
	client := m.(*ProviderConfig)

    id := d.Get("resource_name").(string)
     	V_rate := d.Get("rate").(string)


	config := xmlForwarding__OptionsSamplingInputRate{}
	config.ApplyGroup = id
	config.Groups.Name = id
	config.Groups.V_input.V_rate = &V_rate

    err = client.SendTransaction("", config, false)
    check(err)
    
    d.SetId(fmt.Sprintf("%s_%s", client.Host, id))
    
	return junosForwarding__OptionsSamplingInputRateRead(d,m)
}

func junosForwarding__OptionsSamplingInputRateRead(d *schema.ResourceData, m interface{}) error {

	var err error
	client := m.(*ProviderConfig)

    id := d.Get("resource_name").(string)
    
	config := &xmlForwarding__OptionsSamplingInputRate{}

	err = client.MarshalGroup(id, config)
	check(err)
 	d.Set("rate", config.Groups.V_input.V_rate)

	return nil
}

func junosForwarding__OptionsSamplingInputRateUpdate(d *schema.ResourceData, m interface{}) error {

	var err error
	client := m.(*ProviderConfig)

    id := d.Get("resource_name").(string)
     	V_rate := d.Get("rate").(string)


	config := xmlForwarding__OptionsSamplingInputRate{}
	config.ApplyGroup = id
	config.Groups.Name = id
	config.Groups.V_input.V_rate = &V_rate

    err = client.SendTransaction(id, config, false)
    check(err)
    
	return junosForwarding__OptionsSamplingInputRateRead(d,m)
}

func junosForwarding__OptionsSamplingInputRateDelete(d *schema.ResourceData, m interface{}) error {

	var err error
	client := m.(*ProviderConfig)

    id := d.Get("resource_name").(string)
    
    _, err = client.DeleteConfigNoCommit(id)
    check(err)

    d.SetId("")
    
	return nil
}

func junosForwarding__OptionsSamplingInputRate() *schema.Resource {
	return &schema.Resource{
		Create: junosForwarding__OptionsSamplingInputRateCreate,
		Read: junosForwarding__OptionsSamplingInputRateRead,
		Update: junosForwarding__OptionsSamplingInputRateUpdate,
		Delete: junosForwarding__OptionsSamplingInputRateDelete,

        Schema: map[string]*schema.Schema{
            "resource_name": &schema.Schema{
                Type:    schema.TypeString,
                Required: true,
            },
			"rate": &schema.Schema{
				Type:    schema.TypeString,
				Optional: true,
				Description:    "xpath is: config.Groups.V_input. Ratio of packets to be sampled (1 out of N)",
			},
		},
	}
}