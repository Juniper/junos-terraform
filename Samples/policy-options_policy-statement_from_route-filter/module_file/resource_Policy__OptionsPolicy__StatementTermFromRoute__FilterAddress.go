
package main

import (
    "encoding/xml"
    "fmt"
    "github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)


// v_ is appended before every variable so it doesn't give any conflict
// with any keyword in golang. ex - interface is keyword in golang
type xmlPolicy__OptionsPolicy__StatementTermFromRoute__FilterAddress struct {
	XMLName xml.Name `xml:"configuration"`
	V_policy__statement  struct {
		XMLName xml.Name `xml:"policy-statement"`
		V_name  string  `xml:"name"`
		V_term  struct {
			XMLName xml.Name `xml:"term"`
			V_name__1  string  `xml:"name"`
			V_route__filter  struct {
				XMLName xml.Name `xml:"route-filter"`
				V_choice__ident  string  `xml:"choice-ident"`
				V_choice__value  string  `xml:"choice-value"`
				V_address  string  `xml:"address"`
			} `xml:"from>route-filter"`
		} `xml:"term"`
	} `xml:"policy-options>policy-statement"`
}

// v_ is appended before every variable so it doesn't give any conflict
// with any keyword in golang. ex- interface is keyword in golang
func junosPolicy__OptionsPolicy__StatementTermFromRoute__FilterAddressCreate(d *schema.ResourceData, m interface{}) error {

    pcfg := m.(*ProviderConfig)
    client, err := pcfg.Client()
    check(err)

    id := d.Get("resource_name").(string)
     	V_name := d.Get("name").(string)
	V_name__1 := d.Get("name__1").(string)
	V_choice__ident := d.Get("choice__ident").(string)
	V_choice__value := d.Get("choice__value").(string)
	V_address := d.Get("address").(string)
	commit := true

	config := xmlPolicy__OptionsPolicy__StatementTermFromRoute__FilterAddress{}
	config.V_policy__statement.V_name = V_name
	config.V_policy__statement.V_term.V_name__1 = V_name__1
	config.V_policy__statement.V_term.V_route__filter.V_choice__ident = V_choice__ident
	config.V_policy__statement.V_term.V_route__filter.V_choice__value = V_choice__value
	config.V_policy__statement.V_term.V_route__filter.V_address = V_address

    err = client.SendTransaction("", config, commit)
    check(err)
    
    d.SetId(fmt.Sprintf("%s_%s", pcfg.Cfg.Host, id))
    
    err = client.Close()
    check(err)
    
	return junosPolicy__OptionsPolicy__StatementTermFromRoute__FilterAddressRead(d,m)
}

func junosPolicy__OptionsPolicy__StatementTermFromRoute__FilterAddressRead(d *schema.ResourceData, m interface{}) error {

    pcfg := m.(*ProviderConfig)
    client, err := pcfg.Client()
    check(err)

    id := d.Get("resource_name").(string)
    
	config := &xmlPolicy__OptionsPolicy__StatementTermFromRoute__FilterAddress{}

	err = client.MarshalGroup(id, config)
	check(err)
 	d.Set("name", config.V_policy__statement.V_name)
	d.Set("name__1", config.V_policy__statement.V_term.V_name__1)
	d.Set("choice__ident", config.V_policy__statement.V_term.V_route__filter.V_choice__ident)
	d.Set("choice__value", config.V_policy__statement.V_term.V_route__filter.V_choice__value)
	d.Set("address", config.V_policy__statement.V_term.V_route__filter.V_address)

    err = client.Close()
    check(err)
    
	return nil
}

func junosPolicy__OptionsPolicy__StatementTermFromRoute__FilterAddressUpdate(d *schema.ResourceData, m interface{}) error {

    pcfg := m.(*ProviderConfig)
    client, err := pcfg.Client()
    check(err)

    id := d.Get("resource_name").(string)
     	V_name := d.Get("name").(string)
	V_name__1 := d.Get("name__1").(string)
	V_choice__ident := d.Get("choice__ident").(string)
	V_choice__value := d.Get("choice__value").(string)
	V_address := d.Get("address").(string)
	commit := true

	config := xmlPolicy__OptionsPolicy__StatementTermFromRoute__FilterAddress{}
	config.V_policy__statement.V_name = V_name
	config.V_policy__statement.V_term.V_name__1 = V_name__1
	config.V_policy__statement.V_term.V_route__filter.V_choice__ident = V_choice__ident
	config.V_policy__statement.V_term.V_route__filter.V_choice__value = V_choice__value
	config.V_policy__statement.V_term.V_route__filter.V_address = V_address

    err = client.SendTransaction(id, config, commit)
    check(err)
    
    err = client.Close()
    check(err)
    
	return junosPolicy__OptionsPolicy__StatementTermFromRoute__FilterAddressRead(d,m)
}

func junosPolicy__OptionsPolicy__StatementTermFromRoute__FilterAddressDelete(d *schema.ResourceData, m interface{}) error {

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

func junosPolicy__OptionsPolicy__StatementTermFromRoute__FilterAddress() *schema.Resource {
	return &schema.Resource{
		Create: junosPolicy__OptionsPolicy__StatementTermFromRoute__FilterAddressCreate,
		Read: junosPolicy__OptionsPolicy__StatementTermFromRoute__FilterAddressRead,
		Update: junosPolicy__OptionsPolicy__StatementTermFromRoute__FilterAddressUpdate,
		Delete: junosPolicy__OptionsPolicy__StatementTermFromRoute__FilterAddressDelete,

        Schema: map[string]*schema.Schema{
            "resource_name": &schema.Schema{
                Type:    schema.TypeString,
                Required: true,
            },
			"name": &schema.Schema{
				Type:    schema.TypeString,
				Optional: true,
				Description:    "xpath is: config.V_policy__statement",
			},
			"name__1": &schema.Schema{
				Type:    schema.TypeString,
				Optional: true,
				Description:    "xpath is: config.V_policy__statement.V_term",
			},
			"choice__ident": &schema.Schema{
				Type:    schema.TypeString,
				Optional: true,
				Description:    "xpath is: config.V_policy__statement.V_term.V_route__filter",
			},
			"choice__value": &schema.Schema{
				Type:    schema.TypeString,
				Optional: true,
				Description:    "xpath is: config.V_policy__statement.V_term.V_route__filter",
			},
			"address": &schema.Schema{
				Type:    schema.TypeString,
				Optional: true,
				Description:    "xpath is: config.V_policy__statement.V_term.V_route__filter. IP address or hostname",
			},
		},
	}
}