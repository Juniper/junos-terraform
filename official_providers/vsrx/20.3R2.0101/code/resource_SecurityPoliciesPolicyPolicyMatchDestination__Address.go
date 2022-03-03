
package main

import (
    "encoding/xml"
    "fmt"
    "github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)


// v_ is appended before every variable so it doesn't give any conflict
// with any keyword in golang. ex - interface is keyword in golang
type xmlSecurityPoliciesPolicyPolicyMatchDestination__Address struct {
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
					V_destination__address  *string  `xml:"destination-address,omitempty"`
				} `xml:"match"`
			} `xml:"policy"`
		} `xml:"security>policies>policy"`
	} `xml:"groups"`
	ApplyGroup string `xml:"apply-groups"`
}

// v_ is appended before every variable so it doesn't give any conflict
// with any keyword in golang. ex- interface is keyword in golang
func junosSecurityPoliciesPolicyPolicyMatchDestination__AddressCreate(d *schema.ResourceData, m interface{}) error {

	var err error
	client := m.(*ProviderConfig)

    id := d.Get("resource_name").(string)
     	V_from__zone__name := d.Get("from__zone__name").(string)
	V_to__zone__name := d.Get("to__zone__name").(string)
	V_name := d.Get("name").(string)
	V_destination__address := d.Get("destination__address").(string)


	config := xmlSecurityPoliciesPolicyPolicyMatchDestination__Address{}
	config.ApplyGroup = id
	config.Groups.Name = id
	config.Groups.V_policy.V_from__zone__name = &V_from__zone__name
	config.Groups.V_policy.V_to__zone__name = &V_to__zone__name
	config.Groups.V_policy.V_policy__1.V_name = &V_name
	config.Groups.V_policy.V_policy__1.V_match.V_destination__address = &V_destination__address

    err = client.SendTransaction("", config, false)
    check(err)
    
    d.SetId(fmt.Sprintf("%s_%s", client.Host, id))
    
	return junosSecurityPoliciesPolicyPolicyMatchDestination__AddressRead(d,m)
}

func junosSecurityPoliciesPolicyPolicyMatchDestination__AddressRead(d *schema.ResourceData, m interface{}) error {

	var err error
	client := m.(*ProviderConfig)

    id := d.Get("resource_name").(string)
    
	config := &xmlSecurityPoliciesPolicyPolicyMatchDestination__Address{}

	err = client.MarshalGroup(id, config)
	check(err)
 	d.Set("from__zone__name", config.Groups.V_policy.V_from__zone__name)
	d.Set("to__zone__name", config.Groups.V_policy.V_to__zone__name)
	d.Set("name", config.Groups.V_policy.V_policy__1.V_name)
	d.Set("destination__address", config.Groups.V_policy.V_policy__1.V_match.V_destination__address)

	return nil
}

func junosSecurityPoliciesPolicyPolicyMatchDestination__AddressUpdate(d *schema.ResourceData, m interface{}) error {

	var err error
	client := m.(*ProviderConfig)

    id := d.Get("resource_name").(string)
     	V_from__zone__name := d.Get("from__zone__name").(string)
	V_to__zone__name := d.Get("to__zone__name").(string)
	V_name := d.Get("name").(string)
	V_destination__address := d.Get("destination__address").(string)


	config := xmlSecurityPoliciesPolicyPolicyMatchDestination__Address{}
	config.ApplyGroup = id
	config.Groups.Name = id
	config.Groups.V_policy.V_from__zone__name = &V_from__zone__name
	config.Groups.V_policy.V_to__zone__name = &V_to__zone__name
	config.Groups.V_policy.V_policy__1.V_name = &V_name
	config.Groups.V_policy.V_policy__1.V_match.V_destination__address = &V_destination__address

    err = client.SendTransaction(id, config, false)
    check(err)
    
	return junosSecurityPoliciesPolicyPolicyMatchDestination__AddressRead(d,m)
}

func junosSecurityPoliciesPolicyPolicyMatchDestination__AddressDelete(d *schema.ResourceData, m interface{}) error {

	var err error
	client := m.(*ProviderConfig)

    id := d.Get("resource_name").(string)
    
    _, err = client.DeleteConfigNoCommit(id)
    check(err)

    d.SetId("")
    
	return nil
}

func junosSecurityPoliciesPolicyPolicyMatchDestination__Address() *schema.Resource {
	return &schema.Resource{
		Create: junosSecurityPoliciesPolicyPolicyMatchDestination__AddressCreate,
		Read: junosSecurityPoliciesPolicyPolicyMatchDestination__AddressRead,
		Update: junosSecurityPoliciesPolicyPolicyMatchDestination__AddressUpdate,
		Delete: junosSecurityPoliciesPolicyPolicyMatchDestination__AddressDelete,

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
			"destination__address": &schema.Schema{
				Type:    schema.TypeString,
				Optional: true,
				Description:    "xpath is: config.Groups.V_policy.V_policy__1.V_match. ",
			},
		},
	}
}