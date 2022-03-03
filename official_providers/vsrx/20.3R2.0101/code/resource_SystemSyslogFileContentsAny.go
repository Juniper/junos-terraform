
package main

import (
    "encoding/xml"
    "fmt"
    "github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)


// v_ is appended before every variable so it doesn't give any conflict
// with any keyword in golang. ex - interface is keyword in golang
type xmlSystemSyslogFileContentsAny struct {
	XMLName xml.Name `xml:"configuration"`
	Groups  struct {
		XMLName	xml.Name	`xml:"groups"`
		Name	string	`xml:"name"`
		V_file  struct {
			XMLName xml.Name `xml:"file"`
			V_name  *string  `xml:"name,omitempty"`
			V_contents  struct {
				XMLName xml.Name `xml:"contents"`
				V_name__1  *string  `xml:"name,omitempty"`
				V_any  *string  `xml:"any,omitempty"`
			} `xml:"contents"`
		} `xml:"system>syslog>file"`
	} `xml:"groups"`
	ApplyGroup string `xml:"apply-groups"`
}

// v_ is appended before every variable so it doesn't give any conflict
// with any keyword in golang. ex- interface is keyword in golang
func junosSystemSyslogFileContentsAnyCreate(d *schema.ResourceData, m interface{}) error {

	var err error
	client := m.(*ProviderConfig)

    id := d.Get("resource_name").(string)
     	V_name := d.Get("name").(string)
	V_name__1 := d.Get("name__1").(string)
	V_any := d.Get("any").(string)


	config := xmlSystemSyslogFileContentsAny{}
	config.ApplyGroup = id
	config.Groups.Name = id
	config.Groups.V_file.V_name = &V_name
	config.Groups.V_file.V_contents.V_name__1 = &V_name__1
	config.Groups.V_file.V_contents.V_any = &V_any

    err = client.SendTransaction("", config, false)
    check(err)
    
    d.SetId(fmt.Sprintf("%s_%s", client.Host, id))
    
	return junosSystemSyslogFileContentsAnyRead(d,m)
}

func junosSystemSyslogFileContentsAnyRead(d *schema.ResourceData, m interface{}) error {

	var err error
	client := m.(*ProviderConfig)

    id := d.Get("resource_name").(string)
    
	config := &xmlSystemSyslogFileContentsAny{}

	err = client.MarshalGroup(id, config)
	check(err)
 	d.Set("name", config.Groups.V_file.V_name)
	d.Set("name__1", config.Groups.V_file.V_contents.V_name__1)
	d.Set("any", config.Groups.V_file.V_contents.V_any)

	return nil
}

func junosSystemSyslogFileContentsAnyUpdate(d *schema.ResourceData, m interface{}) error {

	var err error
	client := m.(*ProviderConfig)

    id := d.Get("resource_name").(string)
     	V_name := d.Get("name").(string)
	V_name__1 := d.Get("name__1").(string)
	V_any := d.Get("any").(string)


	config := xmlSystemSyslogFileContentsAny{}
	config.ApplyGroup = id
	config.Groups.Name = id
	config.Groups.V_file.V_name = &V_name
	config.Groups.V_file.V_contents.V_name__1 = &V_name__1
	config.Groups.V_file.V_contents.V_any = &V_any

    err = client.SendTransaction(id, config, false)
    check(err)
    
	return junosSystemSyslogFileContentsAnyRead(d,m)
}

func junosSystemSyslogFileContentsAnyDelete(d *schema.ResourceData, m interface{}) error {

	var err error
	client := m.(*ProviderConfig)

    id := d.Get("resource_name").(string)
    
    _, err = client.DeleteConfigNoCommit(id)
    check(err)

    d.SetId("")
    
	return nil
}

func junosSystemSyslogFileContentsAny() *schema.Resource {
	return &schema.Resource{
		Create: junosSystemSyslogFileContentsAnyCreate,
		Read: junosSystemSyslogFileContentsAnyRead,
		Update: junosSystemSyslogFileContentsAnyUpdate,
		Delete: junosSystemSyslogFileContentsAnyDelete,

        Schema: map[string]*schema.Schema{
            "resource_name": &schema.Schema{
                Type:    schema.TypeString,
                Required: true,
            },
			"name": &schema.Schema{
				Type:    schema.TypeString,
				Optional: true,
				Description:    "xpath is: config.Groups.V_file",
			},
			"name__1": &schema.Schema{
				Type:    schema.TypeString,
				Optional: true,
				Description:    "xpath is: config.Groups.V_file.V_contents",
			},
			"any": &schema.Schema{
				Type:    schema.TypeString,
				Optional: true,
				Description:    "xpath is: config.Groups.V_file.V_contents. All levels",
			},
		},
	}
}