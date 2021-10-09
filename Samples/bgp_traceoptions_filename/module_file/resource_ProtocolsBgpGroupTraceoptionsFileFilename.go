
package main

import (
    "encoding/xml"
    "fmt"
    "github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)


// v_ is appended before every variable so it doesn't give any conflict
// with any keyword in golang. ex - interface is keyword in golang
type xmlProtocolsBgpGroupTraceoptionsFileFilename struct {
	XMLName xml.Name `xml:"configuration"`
	Groups  struct {
		XMLName	xml.Name	`xml:"groups"`
		Name	string	`xml:"name"`
		V_group  struct {
			XMLName xml.Name `xml:"group"`
			V_name  string  `xml:"name"`
			V_file  struct {
				XMLName xml.Name `xml:"file"`
				V_filename  string  `xml:"filename"`
			} `xml:"traceoptions>file"`
		} `xml:"protocols>bgp>group"`
	} `xml:"groups"`
	ApplyGroup string `xml:"apply-groups"`
}

// v_ is appended before every variable so it doesn't give any conflict
// with any keyword in golang. ex- interface is keyword in golang
func junosProtocolsBgpGroupTraceoptionsFileFilenameCreate(d *schema.ResourceData, m interface{}) error {

    pcfg := m.(*ProviderConfig)
    client, err := pcfg.Client()
    check(err)

    id := d.Get("resource_name").(string)
     	V_name := d.Get("name").(string)
	V_filename := d.Get("filename").(string)
	commit := false

	config := xmlProtocolsBgpGroupTraceoptionsFileFilename{}
	config.ApplyGroup = id
	config.Groups.Name = id
	config.Groups.V_group.V_name = V_name
	config.Groups.V_group.V_file.V_filename = V_filename

    err = client.SendTransaction("", config, commit)
    check(err)
    
    d.SetId(fmt.Sprintf("%s_%s", pcfg.Cfg.Host, id))
    
    err = client.Close()
    check(err)
    
	return junosProtocolsBgpGroupTraceoptionsFileFilenameRead(d,m)
}

func junosProtocolsBgpGroupTraceoptionsFileFilenameRead(d *schema.ResourceData, m interface{}) error {

    pcfg := m.(*ProviderConfig)
    client, err := pcfg.Client()
    check(err)

    id := d.Get("resource_name").(string)
    
	config := &xmlProtocolsBgpGroupTraceoptionsFileFilename{}

	err = client.MarshalGroup(id, config)
	check(err)
 	d.Set("name", config.Groups.V_group.V_name)
	d.Set("filename", config.Groups.V_group.V_file.V_filename)

    err = client.Close()
    check(err)
    
	return nil
}

func junosProtocolsBgpGroupTraceoptionsFileFilenameUpdate(d *schema.ResourceData, m interface{}) error {

    pcfg := m.(*ProviderConfig)
    client, err := pcfg.Client()
    check(err)

    id := d.Get("resource_name").(string)
     	V_name := d.Get("name").(string)
	V_filename := d.Get("filename").(string)
	commit := false

	config := xmlProtocolsBgpGroupTraceoptionsFileFilename{}
	config.ApplyGroup = id
	config.Groups.Name = id
	config.Groups.V_group.V_name = V_name
	config.Groups.V_group.V_file.V_filename = V_filename

    err = client.SendTransaction(id, config, commit)
    check(err)
    
    err = client.Close()
    check(err)
    
	return junosProtocolsBgpGroupTraceoptionsFileFilenameRead(d,m)
}

func junosProtocolsBgpGroupTraceoptionsFileFilenameDelete(d *schema.ResourceData, m interface{}) error {

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

func junosProtocolsBgpGroupTraceoptionsFileFilename() *schema.Resource {
	return &schema.Resource{
		Create: junosProtocolsBgpGroupTraceoptionsFileFilenameCreate,
		Read: junosProtocolsBgpGroupTraceoptionsFileFilenameRead,
		Update: junosProtocolsBgpGroupTraceoptionsFileFilenameUpdate,
		Delete: junosProtocolsBgpGroupTraceoptionsFileFilenameDelete,

        Schema: map[string]*schema.Schema{
            "resource_name": &schema.Schema{
                Type:    schema.TypeString,
                Required: true,
            },
			"name": &schema.Schema{
				Type:    schema.TypeString,
				Optional: true,
				Description:    "xpath is: config.Groups.V_group",
			},
			"filename": &schema.Schema{
				Type:    schema.TypeString,
				Optional: true,
				Description:    "xpath is: config.Groups.V_group.V_file. Name of file in which to write trace information",
			},
		},
	}
}