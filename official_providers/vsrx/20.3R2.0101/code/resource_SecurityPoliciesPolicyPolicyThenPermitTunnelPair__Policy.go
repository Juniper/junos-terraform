
package main

import (
    "encoding/xml"
    "fmt"
    "github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)


// v_ is appended before every variable so it doesn't give any conflict
// with any keyword in golang. ex - interface is keyword in golang
type xmlSecurityPoliciesPolicyPolicyThenPermitTunnelPair__Policy struct {
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
				V_tunnel  struct {
					XMLName xml.Name `xml:"tunnel"`
					V_pair__policy  *string  `xml:"pair-policy,omitempty"`
				} `xml:"then>permit>tunnel"`
			} `xml:"policy"`
		} `xml:"security>policies>policy"`
	} `xml:"groups"`
	ApplyGroup string `xml:"apply-groups"`
}

// v_ is appended before every variable so it doesn't give any conflict
// with any keyword in golang. ex- interface is keyword in golang
func junosSecurityPoliciesPolicyPolicyThenPermitTunnelPair__PolicyCreate(d *schema.ResourceData, m interface{}) error {

	var err error
	client := m.(*ProviderConfig)

    id := d.Get("resource_name").(string)
     	V_from__zone__name := d.Get("from__zone__name").(string)
	V_to__zone__name := d.Get("to__zone__name").(string)
	V_name := d.Get("name").(string)
	V_pair__policy := d.Get("pair__policy").(string)


	config := xmlSecurityPoliciesPolicyPolicyThenPermitTunnelPair__Policy{}
	config.ApplyGroup = id
	config.Groups.Name = id
	config.Groups.V_policy.V_from__zone__name = &V_from__zone__name
	config.Groups.V_policy.V_to__zone__name = &V_to__zone__name
	config.Groups.V_policy.V_policy__1.V_name = &V_name
	config.Groups.V_policy.V_policy__1.V_tunnel.V_pair__policy = &V_pair__policy

    err = client.SendTransaction("", config, false)
    check(err)
    
    d.SetId(fmt.Sprintf("%s_%s", client.Host, id))
    
	return junosSecurityPoliciesPolicyPolicyThenPermitTunnelPair__PolicyRead(d,m)
}

func junosSecurityPoliciesPolicyPolicyThenPermitTunnelPair__PolicyRead(d *schema.ResourceData, m interface{}) error {

	var err error
	client := m.(*ProviderConfig)

    id := d.Get("resource_name").(string)
    
	config := &xmlSecurityPoliciesPolicyPolicyThenPermitTunnelPair__Policy{}

	err = client.MarshalGroup(id, config)
	check(err)
 	d.Set("from__zone__name", config.Groups.V_policy.V_from__zone__name)
	d.Set("to__zone__name", config.Groups.V_policy.V_to__zone__name)
	d.Set("name", config.Groups.V_policy.V_policy__1.V_name)
	d.Set("pair__policy", config.Groups.V_policy.V_policy__1.V_tunnel.V_pair__policy)

	return nil
}

func junosSecurityPoliciesPolicyPolicyThenPermitTunnelPair__PolicyUpdate(d *schema.ResourceData, m interface{}) error {

	var err error
	client := m.(*ProviderConfig)

    id := d.Get("resource_name").(string)
     	V_from__zone__name := d.Get("from__zone__name").(string)
	V_to__zone__name := d.Get("to__zone__name").(string)
	V_name := d.Get("name").(string)
	V_pair__policy := d.Get("pair__policy").(string)


	config := xmlSecurityPoliciesPolicyPolicyThenPermitTunnelPair__Policy{}
	config.ApplyGroup = id
	config.Groups.Name = id
	config.Groups.V_policy.V_from__zone__name = &V_from__zone__name
	config.Groups.V_policy.V_to__zone__name = &V_to__zone__name
	config.Groups.V_policy.V_policy__1.V_name = &V_name
	config.Groups.V_policy.V_policy__1.V_tunnel.V_pair__policy = &V_pair__policy

    err = client.SendTransaction(id, config, false)
    check(err)
    
	return junosSecurityPoliciesPolicyPolicyThenPermitTunnelPair__PolicyRead(d,m)
}

func junosSecurityPoliciesPolicyPolicyThenPermitTunnelPair__PolicyDelete(d *schema.ResourceData, m interface{}) error {

	var err error
	client := m.(*ProviderConfig)

    id := d.Get("resource_name").(string)
    
    _, err = client.DeleteConfigNoCommit(id)
    check(err)

    d.SetId("")
    
	return nil
}

func junosSecurityPoliciesPolicyPolicyThenPermitTunnelPair__Policy() *schema.Resource {
	return &schema.Resource{
		Create: junosSecurityPoliciesPolicyPolicyThenPermitTunnelPair__PolicyCreate,
		Read: junosSecurityPoliciesPolicyPolicyThenPermitTunnelPair__PolicyRead,
		Update: junosSecurityPoliciesPolicyPolicyThenPermitTunnelPair__PolicyUpdate,
		Delete: junosSecurityPoliciesPolicyPolicyThenPermitTunnelPair__PolicyDelete,

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
			"pair__policy": &schema.Schema{
				Type:    schema.TypeString,
				Optional: true,
				Description:    "xpath is: config.Groups.V_policy.V_policy__1.V_tunnel. Policy in the reverse direction, to form a pair",
			},
		},
	}
}