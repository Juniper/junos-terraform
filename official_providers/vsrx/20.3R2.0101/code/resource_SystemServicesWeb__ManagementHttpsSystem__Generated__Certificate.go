
package main

import (
    "encoding/xml"
    "fmt"
    "github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)


// v_ is appended before every variable so it doesn't give any conflict
// with any keyword in golang. ex - interface is keyword in golang
type xmlSystemServicesWeb__ManagementHttpsSystem__Generated__Certificate struct {
	XMLName xml.Name `xml:"configuration"`
	Groups  struct {
		XMLName	xml.Name	`xml:"groups"`
		Name	string	`xml:"name"`
		V_https  struct {
			XMLName xml.Name `xml:"https"`
			V_system__generated__certificate  *string  `xml:"system-generated-certificate,omitempty"`
		} `xml:"system>services>web-management>https"`
	} `xml:"groups"`
	ApplyGroup string `xml:"apply-groups"`
}

// v_ is appended before every variable so it doesn't give any conflict
// with any keyword in golang. ex- interface is keyword in golang
func junosSystemServicesWeb__ManagementHttpsSystem__Generated__CertificateCreate(d *schema.ResourceData, m interface{}) error {

	var err error
	client := m.(*ProviderConfig)

    id := d.Get("resource_name").(string)
     	V_system__generated__certificate := d.Get("system__generated__certificate").(string)


	config := xmlSystemServicesWeb__ManagementHttpsSystem__Generated__Certificate{}
	config.ApplyGroup = id
	config.Groups.Name = id
	config.Groups.V_https.V_system__generated__certificate = &V_system__generated__certificate

    err = client.SendTransaction("", config, false)
    check(err)
    
    d.SetId(fmt.Sprintf("%s_%s", client.Host, id))
    
	return junosSystemServicesWeb__ManagementHttpsSystem__Generated__CertificateRead(d,m)
}

func junosSystemServicesWeb__ManagementHttpsSystem__Generated__CertificateRead(d *schema.ResourceData, m interface{}) error {

	var err error
	client := m.(*ProviderConfig)

    id := d.Get("resource_name").(string)
    
	config := &xmlSystemServicesWeb__ManagementHttpsSystem__Generated__Certificate{}

	err = client.MarshalGroup(id, config)
	check(err)
 	d.Set("system__generated__certificate", config.Groups.V_https.V_system__generated__certificate)

	return nil
}

func junosSystemServicesWeb__ManagementHttpsSystem__Generated__CertificateUpdate(d *schema.ResourceData, m interface{}) error {

	var err error
	client := m.(*ProviderConfig)

    id := d.Get("resource_name").(string)
     	V_system__generated__certificate := d.Get("system__generated__certificate").(string)


	config := xmlSystemServicesWeb__ManagementHttpsSystem__Generated__Certificate{}
	config.ApplyGroup = id
	config.Groups.Name = id
	config.Groups.V_https.V_system__generated__certificate = &V_system__generated__certificate

    err = client.SendTransaction(id, config, false)
    check(err)
    
	return junosSystemServicesWeb__ManagementHttpsSystem__Generated__CertificateRead(d,m)
}

func junosSystemServicesWeb__ManagementHttpsSystem__Generated__CertificateDelete(d *schema.ResourceData, m interface{}) error {

	var err error
	client := m.(*ProviderConfig)

    id := d.Get("resource_name").(string)
    
    _, err = client.DeleteConfigNoCommit(id)
    check(err)

    d.SetId("")
    
	return nil
}

func junosSystemServicesWeb__ManagementHttpsSystem__Generated__Certificate() *schema.Resource {
	return &schema.Resource{
		Create: junosSystemServicesWeb__ManagementHttpsSystem__Generated__CertificateCreate,
		Read: junosSystemServicesWeb__ManagementHttpsSystem__Generated__CertificateRead,
		Update: junosSystemServicesWeb__ManagementHttpsSystem__Generated__CertificateUpdate,
		Delete: junosSystemServicesWeb__ManagementHttpsSystem__Generated__CertificateDelete,

        Schema: map[string]*schema.Schema{
            "resource_name": &schema.Schema{
                Type:    schema.TypeString,
                Required: true,
            },
			"system__generated__certificate": &schema.Schema{
				Type:    schema.TypeString,
				Optional: true,
				Description:    "xpath is: config.Groups.V_https. X.509 certificate generated automatically by system",
			},
		},
	}
}