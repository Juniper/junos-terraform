
package main

import (
    "encoding/xml"
    "fmt"
    "github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)


// v_ is appended before every variable so it doesn't give any conflict
// with any keyword in golang. ex - interface is keyword in golang
type xmlProtocolsBgpGroupTraceoptionsFile struct {
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
				V_replace  string  `xml:"replace"`
				V_size  string  `xml:"size"`
				V_files  string  `xml:"files"`
				V_no__stamp  string  `xml:"no-stamp"`
			} `xml:"traceoptions>file"`
		} `xml:"protocols>bgp>group"`
	} `xml:"groups"`
	ApplyGroup string `xml:"apply-groups"`
}

// v_ is appended before every variable so it doesn't give any conflict
// with any keyword in golang. ex- interface is keyword in golang
func junosProtocolsBgpGroupTraceoptionsFileCreate(d *schema.ResourceData, m interface{}) error {

    pcfg := m.(*ProviderConfig)
    client, err := pcfg.Client()
    check(err)

    id := d.Get("resource_name").(string)
     	V_name := d.Get("name").(string)
	V_filename := d.Get("filename").(string)
	V_replace := d.Get("replace").(string)
	V_size := d.Get("size").(string)
	V_files := d.Get("files").(string)
	V_no__stamp := d.Get("no__stamp").(string)
	commit := false

	config := xmlProtocolsBgpGroupTraceoptionsFile{}
	config.ApplyGroup = id
	config.Groups.Name = id
	config.Groups.V_group.V_name = V_name
	config.Groups.V_group.V_file.V_filename = V_filename
	config.Groups.V_group.V_file.V_replace = V_replace
	config.Groups.V_group.V_file.V_size = V_size
	config.Groups.V_group.V_file.V_files = V_files
	config.Groups.V_group.V_file.V_no__stamp = V_no__stamp

    err = client.SendTransaction("", config, commit)
    check(err)
    
    d.SetId(fmt.Sprintf("%s_%s", pcfg.Cfg.Host, id))
    
    err = client.Close()
    check(err)
    
	return junosProtocolsBgpGroupTraceoptionsFileRead(d,m)
}

func junosProtocolsBgpGroupTraceoptionsFileRead(d *schema.ResourceData, m interface{}) error {

    pcfg := m.(*ProviderConfig)
    client, err := pcfg.Client()
    check(err)

    id := d.Get("resource_name").(string)
    
	config := &xmlProtocolsBgpGroupTraceoptionsFile{}

	err = client.MarshalGroup(id, config)
	check(err)
 	d.Set("name", config.Groups.V_group.V_name)
	d.Set("filename", config.Groups.V_group.V_file.V_filename)
	d.Set("replace", config.Groups.V_group.V_file.V_replace)
	d.Set("size", config.Groups.V_group.V_file.V_size)
	d.Set("files", config.Groups.V_group.V_file.V_files)
	d.Set("no__stamp", config.Groups.V_group.V_file.V_no__stamp)

    err = client.Close()
    check(err)
    
	return nil
}

func junosProtocolsBgpGroupTraceoptionsFileUpdate(d *schema.ResourceData, m interface{}) error {

    pcfg := m.(*ProviderConfig)
    client, err := pcfg.Client()
    check(err)

    id := d.Get("resource_name").(string)
     	V_name := d.Get("name").(string)
	V_filename := d.Get("filename").(string)
	V_replace := d.Get("replace").(string)
	V_size := d.Get("size").(string)
	V_files := d.Get("files").(string)
	V_no__stamp := d.Get("no__stamp").(string)
	commit := false

	config := xmlProtocolsBgpGroupTraceoptionsFile{}
	config.ApplyGroup = id
	config.Groups.Name = id
	config.Groups.V_group.V_name = V_name
	config.Groups.V_group.V_file.V_filename = V_filename
	config.Groups.V_group.V_file.V_replace = V_replace
	config.Groups.V_group.V_file.V_size = V_size
	config.Groups.V_group.V_file.V_files = V_files
	config.Groups.V_group.V_file.V_no__stamp = V_no__stamp

    err = client.SendTransaction(id, config, commit)
    check(err)
    
    err = client.Close()
    check(err)
    
	return junosProtocolsBgpGroupTraceoptionsFileRead(d,m)
}

func junosProtocolsBgpGroupTraceoptionsFileDelete(d *schema.ResourceData, m interface{}) error {

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

func junosProtocolsBgpGroupTraceoptionsFile() *schema.Resource {
	return &schema.Resource{
		Create: junosProtocolsBgpGroupTraceoptionsFileCreate,
		Read: junosProtocolsBgpGroupTraceoptionsFileRead,
		Update: junosProtocolsBgpGroupTraceoptionsFileUpdate,
		Delete: junosProtocolsBgpGroupTraceoptionsFileDelete,

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
			"replace": &schema.Schema{
				Type:    schema.TypeString,
				Optional: true,
				Description:    "xpath is: config.Groups.V_group.V_file. Replace trace file rather than appending to it",
			},
			"size": &schema.Schema{
				Type:    schema.TypeString,
				Optional: true,
				Description:    "xpath is: config.Groups.V_group.V_file. Maximum trace file size",
			},
			"files": &schema.Schema{
				Type:    schema.TypeString,
				Optional: true,
				Description:    "xpath is: config.Groups.V_group.V_file. Maximum number of trace files",
			},
			"no__stamp": &schema.Schema{
				Type:    schema.TypeString,
				Optional: true,
				Description:    "xpath is: config.Groups.V_group.V_file. Do not timestamp trace file",
			},
		},
	}
}