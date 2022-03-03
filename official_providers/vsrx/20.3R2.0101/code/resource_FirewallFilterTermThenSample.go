
package main

import (
    "encoding/xml"
    "fmt"
    "github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)


// v_ is appended before every variable so it doesn't give any conflict
// with any keyword in golang. ex - interface is keyword in golang
type xmlFirewallFilterTermThenSample struct {
	XMLName xml.Name `xml:"configuration"`
	Groups  struct {
		XMLName	xml.Name	`xml:"groups"`
		Name	string	`xml:"name"`
		V_filter  struct {
			XMLName xml.Name `xml:"filter"`
			V_name  *string  `xml:"name,omitempty"`
			V_term  struct {
				XMLName xml.Name `xml:"term"`
				V_name__1  *string  `xml:"name,omitempty"`
				V_then  struct {
					XMLName xml.Name `xml:"then"`
					V_sample  *string  `xml:"sample,omitempty"`
				} `xml:"then"`
			} `xml:"term"`
		} `xml:"firewall>filter"`
	} `xml:"groups"`
	ApplyGroup string `xml:"apply-groups"`
}

// v_ is appended before every variable so it doesn't give any conflict
// with any keyword in golang. ex- interface is keyword in golang
func junosFirewallFilterTermThenSampleCreate(d *schema.ResourceData, m interface{}) error {

	var err error
	client := m.(*ProviderConfig)

    id := d.Get("resource_name").(string)
     	V_name := d.Get("name").(string)
	V_name__1 := d.Get("name__1").(string)
	V_sample := d.Get("sample").(string)


	config := xmlFirewallFilterTermThenSample{}
	config.ApplyGroup = id
	config.Groups.Name = id
	config.Groups.V_filter.V_name = &V_name
	config.Groups.V_filter.V_term.V_name__1 = &V_name__1
	config.Groups.V_filter.V_term.V_then.V_sample = &V_sample

    err = client.SendTransaction("", config, false)
    check(err)
    
    d.SetId(fmt.Sprintf("%s_%s", client.Host, id))
    
	return junosFirewallFilterTermThenSampleRead(d,m)
}

func junosFirewallFilterTermThenSampleRead(d *schema.ResourceData, m interface{}) error {

	var err error
	client := m.(*ProviderConfig)

    id := d.Get("resource_name").(string)
    
	config := &xmlFirewallFilterTermThenSample{}

	err = client.MarshalGroup(id, config)
	check(err)
 	d.Set("name", config.Groups.V_filter.V_name)
	d.Set("name__1", config.Groups.V_filter.V_term.V_name__1)
	d.Set("sample", config.Groups.V_filter.V_term.V_then.V_sample)

	return nil
}

func junosFirewallFilterTermThenSampleUpdate(d *schema.ResourceData, m interface{}) error {

	var err error
	client := m.(*ProviderConfig)

    id := d.Get("resource_name").(string)
     	V_name := d.Get("name").(string)
	V_name__1 := d.Get("name__1").(string)
	V_sample := d.Get("sample").(string)


	config := xmlFirewallFilterTermThenSample{}
	config.ApplyGroup = id
	config.Groups.Name = id
	config.Groups.V_filter.V_name = &V_name
	config.Groups.V_filter.V_term.V_name__1 = &V_name__1
	config.Groups.V_filter.V_term.V_then.V_sample = &V_sample

    err = client.SendTransaction(id, config, false)
    check(err)
    
	return junosFirewallFilterTermThenSampleRead(d,m)
}

func junosFirewallFilterTermThenSampleDelete(d *schema.ResourceData, m interface{}) error {

	var err error
	client := m.(*ProviderConfig)

    id := d.Get("resource_name").(string)
    
    _, err = client.DeleteConfigNoCommit(id)
    check(err)

    d.SetId("")
    
	return nil
}

func junosFirewallFilterTermThenSample() *schema.Resource {
	return &schema.Resource{
		Create: junosFirewallFilterTermThenSampleCreate,
		Read: junosFirewallFilterTermThenSampleRead,
		Update: junosFirewallFilterTermThenSampleUpdate,
		Delete: junosFirewallFilterTermThenSampleDelete,

        Schema: map[string]*schema.Schema{
            "resource_name": &schema.Schema{
                Type:    schema.TypeString,
                Required: true,
            },
			"name": &schema.Schema{
				Type:    schema.TypeString,
				Optional: true,
				Description:    "xpath is: config.Groups.V_filter",
			},
			"name__1": &schema.Schema{
				Type:    schema.TypeString,
				Optional: true,
				Description:    "xpath is: config.Groups.V_filter.V_term",
			},
			"sample": &schema.Schema{
				Type:    schema.TypeString,
				Optional: true,
				Description:    "xpath is: config.Groups.V_filter.V_term.V_then. Sample the packet",
			},
		},
	}
}