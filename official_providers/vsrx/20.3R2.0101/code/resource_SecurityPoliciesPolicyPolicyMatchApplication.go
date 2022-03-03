
package main

import (
    "encoding/xml"
    "fmt"
    "github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)


// v_ is appended before every variable so it doesn't give any conflict
// with any keyword in golang. ex - interface is keyword in golang
type xmlSecurityPoliciesPolicyPolicyMatchApplication struct {
	XMLName xml.Name `xml:"configuration"`
	Groups  struct {
		XMLName	xml.Name	`xml:"groups"`
		Name	string	`xml:"name"`
		V_policy  struct {
			XMLName xml.Name `xml:"policy"`
			V_from__zone__name  *string  `xml:"from-zone-name,omitempty"`
			V_to__zone__name  *string  `xml:"to-zone-name,omitempty"`
			V_policy__1  struct {
				XMLName xml.Name `xml:"policy"`
				V_name  *string  `xml:"name,omitempty"`
				V_match  struct {
					XMLName xml.Name `xml:"match"`
					V_application  *string  `xml:"application,omitempty"`
				} `xml:"match"`
			} `xml:"policy"`
		} `xml:"security>policies>policy"`
	} `xml:"groups"`
	ApplyGroup string `xml:"apply-groups"`
}

// v_ is appended before every variable so it doesn't give any conflict
// with any keyword in golang. ex- interface is keyword in golang
func junosSecurityPoliciesPolicyPolicyMatchApplicationCreate(d *schema.ResourceData, m interface{}) error {

	var err error
	client := m.(*ProviderConfig)

    id := d.Get("resource_name").(string)
     	V_from__zone__name := d.Get("from__zone__name").(string)
	V_to__zone__name := d.Get("to__zone__name").(string)
	V_name := d.Get("name").(string)
	V_application := d.Get("application").(string)


	config := xmlSecurityPoliciesPolicyPolicyMatchApplication{}
	config.ApplyGroup = id
	config.Groups.Name = id
	config.Groups.V_policy.V_from__zone__name = &V_from__zone__name
	config.Groups.V_policy.V_to__zone__name = &V_to__zone__name
	config.Groups.V_policy.V_policy__1.V_name = &V_name
	config.Groups.V_policy.V_policy__1.V_match.V_application = &V_application

    err = client.SendTransaction("", config, false)
    check(err)
    
    d.SetId(fmt.Sprintf("%s_%s", client.Host, id))
    
	return junosSecurityPoliciesPolicyPolicyMatchApplicationRead(d,m)
}

func junosSecurityPoliciesPolicyPolicyMatchApplicationRead(d *schema.ResourceData, m interface{}) error {

	var err error
	client := m.(*ProviderConfig)

    id := d.Get("resource_name").(string)
    
	config := &xmlSecurityPoliciesPolicyPolicyMatchApplication{}

	err = client.MarshalGroup(id, config)
	check(err)
 	d.Set("from__zone__name", config.Groups.V_policy.V_from__zone__name)
	d.Set("to__zone__name", config.Groups.V_policy.V_to__zone__name)
	d.Set("name", config.Groups.V_policy.V_policy__1.V_name)
	d.Set("application", config.Groups.V_policy.V_policy__1.V_match.V_application)

	return nil
}

func junosSecurityPoliciesPolicyPolicyMatchApplicationUpdate(d *schema.ResourceData, m interface{}) error {

	var err error
	client := m.(*ProviderConfig)

    id := d.Get("resource_name").(string)
     	V_from__zone__name := d.Get("from__zone__name").(string)
	V_to__zone__name := d.Get("to__zone__name").(string)
	V_name := d.Get("name").(string)
	V_application := d.Get("application").(string)


	config := xmlSecurityPoliciesPolicyPolicyMatchApplication{}
	config.ApplyGroup = id
	config.Groups.Name = id
	config.Groups.V_policy.V_from__zone__name = &V_from__zone__name
	config.Groups.V_policy.V_to__zone__name = &V_to__zone__name
	config.Groups.V_policy.V_policy__1.V_name = &V_name
	config.Groups.V_policy.V_policy__1.V_match.V_application = &V_application

    err = client.SendTransaction(id, config, false)
    check(err)
    
	return junosSecurityPoliciesPolicyPolicyMatchApplicationRead(d,m)
}

func junosSecurityPoliciesPolicyPolicyMatchApplicationDelete(d *schema.ResourceData, m interface{}) error {

	var err error
	client := m.(*ProviderConfig)

    id := d.Get("resource_name").(string)
    
    _, err = client.DeleteConfigNoCommit(id)
    check(err)

    d.SetId("")
    
	return nil
}

func junosSecurityPoliciesPolicyPolicyMatchApplication() *schema.Resource {
	return &schema.Resource{
		Create: junosSecurityPoliciesPolicyPolicyMatchApplicationCreate,
		Read: junosSecurityPoliciesPolicyPolicyMatchApplicationRead,
		Update: junosSecurityPoliciesPolicyPolicyMatchApplicationUpdate,
		Delete: junosSecurityPoliciesPolicyPolicyMatchApplicationDelete,

        Schema: map[string]*schema.Schema{
            "resource_name": &schema.Schema{
                Type:    schema.TypeString,
                Required: true,
            },
			"from__zone__name": &schema.Schema{
				Type:    schema.TypeString,
				Optional: true,
				Description:    "xpath is: config.Groups.V_policy",
			},
			"to__zone__name": &schema.Schema{
				Type:    schema.TypeString,
				Optional: true,
				Description:    "xpath is: config.Groups.V_policy",
			},
			"name": &schema.Schema{
				Type:    schema.TypeString,
				Optional: true,
				Description:    "xpath is: config.Groups.V_policy.V_policy__1",
			},
			"application": &schema.Schema{
				Type:    schema.TypeString,
				Optional: true,
				Description:    "xpath is: config.Groups.V_policy.V_policy__1.V_match. ",
			},
		},
	}
}