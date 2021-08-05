
package main

import (
    "encoding/xml"
    "fmt"
    "github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)


// v_ is appended before every variable so it doesn't give any conflict
// with any keyword in golang. ex - interface is keyword in golang
type xmlSecurityNatSourceRule__SetRuleThenSource__NatPoolPool__Name struct {
	XMLName xml.Name `xml:"configuration"`
	V_rule__set  struct {
		XMLName xml.Name `xml:"rule-set"`
		V_name  string  `xml:"name"`
		V_rule  struct {
			XMLName xml.Name `xml:"rule"`
			V_name__1  string  `xml:"name"`
			V_pool  struct {
				XMLName xml.Name `xml:"pool"`
				V_pool__name  string  `xml:"pool-name"`
			} `xml:"then>source-nat>pool"`
		} `xml:"rule"`
	} `xml:"security>nat>source>rule-set"`
}

// v_ is appended before every variable so it doesn't give any conflict
// with any keyword in golang. ex- interface is keyword in golang
func junosSecurityNatSourceRule__SetRuleThenSource__NatPoolPool__NameCreate(d *schema.ResourceData, m interface{}) error {

    pcfg := m.(*ProviderConfig)
    client, err := pcfg.Client()
    check(err)

    id := d.Get("resource_name").(string)
     	V_name := d.Get("name").(string)
	V_name__1 := d.Get("name__1").(string)
	V_pool__name := d.Get("pool__name").(string)
	commit := true

	config := xmlSecurityNatSourceRule__SetRuleThenSource__NatPoolPool__Name{}
	config.V_rule__set.V_name = V_name
	config.V_rule__set.V_rule.V_name__1 = V_name__1
	config.V_rule__set.V_rule.V_pool.V_pool__name = V_pool__name

    err = client.SendTransaction("", config, commit)
    check(err)
    
    d.SetId(fmt.Sprintf("%s_%s", pcfg.Cfg.Host, id))
    
    err = client.Close()
    check(err)
    
	return junosSecurityNatSourceRule__SetRuleThenSource__NatPoolPool__NameRead(d,m)
}

func junosSecurityNatSourceRule__SetRuleThenSource__NatPoolPool__NameRead(d *schema.ResourceData, m interface{}) error {

    pcfg := m.(*ProviderConfig)
    client, err := pcfg.Client()
    check(err)

    id := d.Get("resource_name").(string)
    
	config := &xmlSecurityNatSourceRule__SetRuleThenSource__NatPoolPool__Name{}

	err = client.MarshalGroup(id, config)
	check(err)
 	d.Set("name", config.V_rule__set.V_name)
	d.Set("name__1", config.V_rule__set.V_rule.V_name__1)
	d.Set("pool__name", config.V_rule__set.V_rule.V_pool.V_pool__name)

    err = client.Close()
    check(err)
    
	return nil
}

func junosSecurityNatSourceRule__SetRuleThenSource__NatPoolPool__NameUpdate(d *schema.ResourceData, m interface{}) error {

    pcfg := m.(*ProviderConfig)
    client, err := pcfg.Client()
    check(err)

    id := d.Get("resource_name").(string)
     	V_name := d.Get("name").(string)
	V_name__1 := d.Get("name__1").(string)
	V_pool__name := d.Get("pool__name").(string)
	commit := true

	config := xmlSecurityNatSourceRule__SetRuleThenSource__NatPoolPool__Name{}
	config.V_rule__set.V_name = V_name
	config.V_rule__set.V_rule.V_name__1 = V_name__1
	config.V_rule__set.V_rule.V_pool.V_pool__name = V_pool__name

    err = client.SendTransaction(id, config, commit)
    check(err)
    
    err = client.Close()
    check(err)
    
	return junosSecurityNatSourceRule__SetRuleThenSource__NatPoolPool__NameRead(d,m)
}

func junosSecurityNatSourceRule__SetRuleThenSource__NatPoolPool__NameDelete(d *schema.ResourceData, m interface{}) error {

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

func junosSecurityNatSourceRule__SetRuleThenSource__NatPoolPool__Name() *schema.Resource {
	return &schema.Resource{
		Create: junosSecurityNatSourceRule__SetRuleThenSource__NatPoolPool__NameCreate,
		Read: junosSecurityNatSourceRule__SetRuleThenSource__NatPoolPool__NameRead,
		Update: junosSecurityNatSourceRule__SetRuleThenSource__NatPoolPool__NameUpdate,
		Delete: junosSecurityNatSourceRule__SetRuleThenSource__NatPoolPool__NameDelete,

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
			"name__1": &schema.Schema{
				Type:    schema.TypeString,
				Optional: true,
				Description:    "xpath is: config.V_rule__set.V_rule",
			},
			"pool__name": &schema.Schema{
				Type:    schema.TypeString,
				Optional: true,
				Description:    "xpath is: config.V_rule__set.V_rule.V_pool. Name of Source NAT pool",
			},
		},
	}
}