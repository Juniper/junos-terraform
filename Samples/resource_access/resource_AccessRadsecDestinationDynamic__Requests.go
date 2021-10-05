
package main

import (
    "encoding/xml"
    "fmt"
    "github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)


// v_ is appended before every variable so it doesn't give any conflict
// with any keyword in golang. ex - interface is keyword in golang
type xmlAccessRadsecDestinationDynamic__Requests struct {
	XMLName xml.Name `xml:"configuration"`
	Groups  struct {
		XMLName	xml.Name	`xml:"groups"`
		Name	string	`xml:"name"`
		V_destination  struct {
			XMLName xml.Name `xml:"destination"`
			V_name  string  `xml:"name"`
			V_dynamic__requests  struct {
				XMLName xml.Name `xml:"dynamic-requests"`
				V_apply__groups  string  `xml:"apply-groups"`
				V_apply__groups__except  string  `xml:"apply-groups-except"`
				V_apply__macro	struct {
					XMLName xml.Name `xml:"apply-macro"`
					V_name__1  string  `xml:"name"`
					V_data	struct {
						XMLName xml.Name `xml:"data"`
						V_name__2  string  `xml:"name"`
						V_value  string  `xml:"value"`
					} `xml:"data"`
				} `xml:"apply-macro"`
				V_source__address  string  `xml:"source-address"`
				V_source__port  string  `xml:"source-port"`
			} `xml:"dynamic-requests"`
		} `xml:"access>radsec>destination"`
	} `xml:"groups"`
	ApplyGroup string `xml:"apply-groups"`
}

// v_ is appended before every variable so it doesn't give any conflict
// with any keyword in golang. ex- interface is keyword in golang
func junosAccessRadsecDestinationDynamic__RequestsCreate(d *schema.ResourceData, m interface{}) error {

    pcfg := m.(*ProviderConfig)
    client, err := pcfg.Client()
    check(err)

    id := d.Get("resource_name").(string)
     	V_name := d.Get("name").(string)
	V_apply__groups := d.Get("apply__groups").(string)
	V_apply__groups__except := d.Get("apply__groups__except").(string)
	V_name__1 := d.Get("name__1").(string)
	V_name__2 := d.Get("name__2").(string)
	V_value := d.Get("value").(string)
	V_source__address := d.Get("source__address").(string)
	V_source__port := d.Get("source__port").(string)
	commit := false

	config := xmlAccessRadsecDestinationDynamic__Requests{}
	config.ApplyGroup = id
	config.Groups.Name = id
	config.Groups.V_destination.V_name = V_name
	config.Groups.V_destination.V_dynamic__requests.V_apply__groups = V_apply__groups
	config.Groups.V_destination.V_dynamic__requests.V_apply__groups__except = V_apply__groups__except
	config.Groups.V_destination.V_dynamic__requests.V_apply__macro.V_name__1 = V_name__1
	config.Groups.V_destination.V_dynamic__requests.V_apply__macro.V_data.V_name__2 = V_name__2
	config.Groups.V_destination.V_dynamic__requests.V_apply__macro.V_data.V_value = V_value
	config.Groups.V_destination.V_dynamic__requests.V_source__address = V_source__address
	config.Groups.V_destination.V_dynamic__requests.V_source__port = V_source__port

    err = client.SendTransaction("", config, commit)
    check(err)
    
    d.SetId(fmt.Sprintf("%s_%s", pcfg.Cfg.Host, id))
    
    err = client.Close()
    check(err)
    
	return junosAccessRadsecDestinationDynamic__RequestsRead(d,m)
}

func junosAccessRadsecDestinationDynamic__RequestsRead(d *schema.ResourceData, m interface{}) error {

    pcfg := m.(*ProviderConfig)
    client, err := pcfg.Client()
    check(err)

    id := d.Get("resource_name").(string)
    
	config := &xmlAccessRadsecDestinationDynamic__Requests{}

	err = client.MarshalGroup(id, config)
	check(err)
 	d.Set("name", config.Groups.V_destination.V_name)
	d.Set("apply__groups", config.Groups.V_destination.V_dynamic__requests.V_apply__groups)
	d.Set("apply__groups__except", config.Groups.V_destination.V_dynamic__requests.V_apply__groups__except)
	d.Set("name__1", config.Groups.V_destination.V_dynamic__requests.V_apply__macro.V_name__1)
	d.Set("name__2", config.Groups.V_destination.V_dynamic__requests.V_apply__macro.V_data.V_name__2)
	d.Set("value", config.Groups.V_destination.V_dynamic__requests.V_apply__macro.V_data.V_value)
	d.Set("source__address", config.Groups.V_destination.V_dynamic__requests.V_source__address)
	d.Set("source__port", config.Groups.V_destination.V_dynamic__requests.V_source__port)

    err = client.Close()
    check(err)
    
	return nil
}

func junosAccessRadsecDestinationDynamic__RequestsUpdate(d *schema.ResourceData, m interface{}) error {

    pcfg := m.(*ProviderConfig)
    client, err := pcfg.Client()
    check(err)

    id := d.Get("resource_name").(string)
     	V_name := d.Get("name").(string)
	V_apply__groups := d.Get("apply__groups").(string)
	V_apply__groups__except := d.Get("apply__groups__except").(string)
	V_name__1 := d.Get("name__1").(string)
	V_name__2 := d.Get("name__2").(string)
	V_value := d.Get("value").(string)
	V_source__address := d.Get("source__address").(string)
	V_source__port := d.Get("source__port").(string)
	commit := false

	config := xmlAccessRadsecDestinationDynamic__Requests{}
	config.ApplyGroup = id
	config.Groups.Name = id
	config.Groups.V_destination.V_name = V_name
	config.Groups.V_destination.V_dynamic__requests.V_apply__groups = V_apply__groups
	config.Groups.V_destination.V_dynamic__requests.V_apply__groups__except = V_apply__groups__except
	config.Groups.V_destination.V_dynamic__requests.V_apply__macro.V_name__1 = V_name__1
	config.Groups.V_destination.V_dynamic__requests.V_apply__macro.V_data.V_name__2 = V_name__2
	config.Groups.V_destination.V_dynamic__requests.V_apply__macro.V_data.V_value = V_value
	config.Groups.V_destination.V_dynamic__requests.V_source__address = V_source__address
	config.Groups.V_destination.V_dynamic__requests.V_source__port = V_source__port

    err = client.SendTransaction(id, config, commit)
    check(err)
    
    err = client.Close()
    check(err)
    
	return junosAccessRadsecDestinationDynamic__RequestsRead(d,m)
}

func junosAccessRadsecDestinationDynamic__RequestsDelete(d *schema.ResourceData, m interface{}) error {

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

func junosAccessRadsecDestinationDynamic__Requests() *schema.Resource {
	return &schema.Resource{
		Create: junosAccessRadsecDestinationDynamic__RequestsCreate,
		Read: junosAccessRadsecDestinationDynamic__RequestsRead,
		Update: junosAccessRadsecDestinationDynamic__RequestsUpdate,
		Delete: junosAccessRadsecDestinationDynamic__RequestsDelete,

        Schema: map[string]*schema.Schema{
            "resource_name": &schema.Schema{
                Type:    schema.TypeString,
                Required: true,
            },
			"name": &schema.Schema{
				Type:    schema.TypeString,
				Optional: true,
				Description:    "xpath is: config.Groups.V_destination",
			},
			"apply__groups": &schema.Schema{
				Type:    schema.TypeString,
				Optional: true,
				Description:    "xpath is: config.Groups.V_destination.V_dynamic__requests. Groups from which to inherit configuration data",
			},
			"apply__groups__except": &schema.Schema{
				Type:    schema.TypeString,
				Optional: true,
				Description:    "xpath is: config.Groups.V_destination.V_dynamic__requests. Don't inherit configuration data from these groups",
			},
			"name__1": &schema.Schema{
				Type:    schema.TypeString,
				Optional: true,
				Description:    "xpath is: config.Groups.V_destination.V_dynamic__requests.V_apply__macro. Name of the macro to be expanded",
			},
			"name__2": &schema.Schema{
				Type:    schema.TypeString,
				Optional: true,
				Description:    "xpath is: config.Groups.V_destination.V_dynamic__requests.V_apply__macro.V_data. Keyword part of the keyword-value pair",
			},
			"value": &schema.Schema{
				Type:    schema.TypeString,
				Optional: true,
				Description:    "xpath is: config.Groups.V_destination.V_dynamic__requests.V_apply__macro.V_data. Value part of the keyword-value pair",
			},
			"source__address": &schema.Schema{
				Type:    schema.TypeString,
				Optional: true,
				Description:    "xpath is: config.Groups.V_destination.V_dynamic__requests. Source address for dynamic requests",
			},
			"source__port": &schema.Schema{
				Type:    schema.TypeString,
				Optional: true,
				Description:    "xpath is: config.Groups.V_destination.V_dynamic__requests. Source port for dynamic requests",
			},
		},
	}
}