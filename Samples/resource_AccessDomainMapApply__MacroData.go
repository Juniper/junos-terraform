
// Copyright (c) 2017-2021, Juniper Networks Inc. All rights reserved.
//
// License: Apache 2.0
//
// THIS SOFTWARE IS PROVIDED BY Juniper Networks, Inc. ''AS IS'' AND ANY
// EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE IMPLIED
// WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE ARE
// DISCLAIMED. IN NO EVENT SHALL Juniper Networks, Inc. BE LIABLE FOR ANY
// DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES
// (INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES;
// LOSS OF USE, DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND
// ON ANY THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT
// (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE OF THIS
// SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.
//

package main

import (
    "encoding/xml"
    "fmt"
    "github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)


// v_ is appended before every variable so it doesn't give any conflict
// with any keyword in golang. ex - interface is keyword in golang
type xmlAccessDomainMapApply__MacroData struct {
	XMLName xml.Name `xml:"configuration"`
	Groups  struct {
		XMLName	xml.Name	`xml:"groups"`
		Name	string	`xml:"name"`
		V_map  struct {
			XMLName xml.Name `xml:"map"`
			V_name  string  `xml:"name"`
			V_apply-macro  struct {
				XMLName xml.Name `xml:"apply-macro"`
				V_name__1  string  `xml:"name"`
				V_data  struct {
					XMLName xml.Name `xml:"data"`
					V_name__2  string  `xml:"name"`
					V_value  string  `xml:"value"`
				} `xml:"data"`
			} `xml:"apply-macro"`
		} `xml:"access>domain>map"`
	} `xml:"groups"`
	ApplyGroup string `xml:"apply-groups"`
}

// v_ is appended before every variable so it doesn't give any conflict
// with any keyword in golang. ex- interface is keyword in golang
func junosAccessDomainMapApply__MacroDataCreate(d *schema.ResourceData, m interface{}) error {

    pcfg := m.(*ProviderConfig)
    client, err := pcfg.Client()
    check(err)

    id := d.Get("resource_name").(string)
     	V_name := d.Get("name").(string)
	V_name__1 := d.Get("name__1").(string)
	V_name__2 := d.Get("name__2").(string)
	V_value := d.Get("value").(string)
	commit := true

	config := xmlAccessDomainMapApply__MacroData{}
	config.ApplyGroup = id
	config.Groups.Name = id
	config.Groups.V_map.V_name = V_name
	config.Groups.V_map.V_apply-macro.V_name__1 = V_name__1
	config.Groups.V_map.V_apply-macro.V_data.V_name__2 = V_name__2
	config.Groups.V_map.V_apply-macro.V_data.V_value = V_value

    err = client.SendTransaction("", config, commit)
    check(err)
    
    d.SetId(fmt.Sprintf("%s_%s", pcfg.Cfg.Host, id))
    
    err = client.Close()
    check(err)
    
	return junosAccessDomainMapApply__MacroDataRead(d,m)
}

func junosAccessDomainMapApply__MacroDataRead(d *schema.ResourceData, m interface{}) error {

    pcfg := m.(*ProviderConfig)
    client, err := pcfg.Client()
    check(err)

    id := d.Get("resource_name").(string)
    
	config := &xmlAccessDomainMapApply__MacroData{}

	err = client.MarshalGroup(id, config)
	check(err)
 	d.Set("name", config.Groups.V_map.V_name)
	d.Set("name__1", config.Groups.V_map.V_apply-macro.V_name__1)
	d.Set("name__2", config.Groups.V_map.V_apply-macro.V_data.V_name__2)
	d.Set("value", config.Groups.V_map.V_apply-macro.V_data.V_value)

    err = client.Close()
    check(err)
    
	return nil
}

func junosAccessDomainMapApply__MacroDataUpdate(d *schema.ResourceData, m interface{}) error {

    pcfg := m.(*ProviderConfig)
    client, err := pcfg.Client()
    check(err)

    id := d.Get("resource_name").(string)
     	V_name := d.Get("name").(string)
	V_name__1 := d.Get("name__1").(string)
	V_name__2 := d.Get("name__2").(string)
	V_value := d.Get("value").(string)
	commit := true

	config := xmlAccessDomainMapApply__MacroData{}
	config.ApplyGroup = id
	config.Groups.Name = id
	config.Groups.V_map.V_name = V_name
	config.Groups.V_map.V_apply-macro.V_name__1 = V_name__1
	config.Groups.V_map.V_apply-macro.V_data.V_name__2 = V_name__2
	config.Groups.V_map.V_apply-macro.V_data.V_value = V_value

    err = client.SendTransaction(id, config, commit)
    check(err)
    
    err = client.Close()
    check(err)
    
	return junosAccessDomainMapApply__MacroDataRead(d,m)
}

func junosAccessDomainMapApply__MacroDataDelete(d *schema.ResourceData, m interface{}) error {

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

func junosAccessDomainMapApply__MacroData() *schema.Resource {
	return &schema.Resource{
		Create: junosAccessDomainMapApply__MacroDataCreate,
		Read: junosAccessDomainMapApply__MacroDataRead,
		Update: junosAccessDomainMapApply__MacroDataUpdate,
		Delete: junosAccessDomainMapApply__MacroDataDelete,

        Schema: map[string]*schema.Schema{
            "resource_name": &schema.Schema{
                Type:    schema.TypeString,
                Required: true,
            },
			"name": &schema.Schema{
				Type:    schema.TypeString,
				Optional: true,
				Description:    "xpath is: config.Groups.V_map",
			},
			"name__1": &schema.Schema{
				Type:    schema.TypeString,
				Optional: true,
				Description:    "xpath is: config.Groups.V_map.V_apply-macro",
			},
			"name__2": &schema.Schema{
				Type:    schema.TypeString,
				Optional: true,
				Description:    "xpath is: config.Groups.V_map.V_apply-macro.V_data. Keyword part of the keyword-value pair",
			},
			"value": &schema.Schema{
				Type:    schema.TypeString,
				Optional: true,
				Description:    "xpath is: config.Groups.V_map.V_apply-macro.V_data. Value part of the keyword-value pair",
			},
		},
	}
}