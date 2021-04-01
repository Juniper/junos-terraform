// Copyright (c) 2017-2021, Juniper Networks Inc. All rights reserved.
//
// License: Apache 2.0
//
// Redistribution and use in source and binary forms, with or without
// modification, are permitted provided that the following conditions are met:
//
// * Redistributions of source code must retain the above copyright
//   notice, this list of conditions and the following disclaimer.
//
// * Redistributions in binary form must reproduce the above copyright
//   notice, this list of conditions and the following disclaimer in the
//   documentation and/or other materials provided with the distribution.
//
// * Neither the name of the Juniper Networks nor the
//   names of its contributors may be used to endorse or promote products
//   derived from this software without specific prior written permission.
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
type xmlAccessRadsecDestinationDynamic__Requests struct {
	XMLName xml.Name `xml:"configuration"`
	V_destination  struct {
		XMLName xml.Name `xml:"destination"`
		V_name  string  `xml:"name"`
		V_dynamic-requests  struct {
			XMLName xml.Name `xml:"dynamic-requests"`
			V_apply-groups  string  `xml:"apply-groups"`
			V_apply-groups-except  string  `xml:"apply-groups-except"`
			V_apply-macro	struct {
				XMLName xml.Name `xml:"apply-macro"`
				V_name__1  string  `xml:"name"`
				V_data	struct {
					XMLName xml.Name `xml:"data"`
					V_name__2  string  `xml:"name"`
					V_value  string  `xml:"value"`
				} `xml:"data"`
			} `xml:"apply-macro"`
			V_source-address  string  `xml:"source-address"`
			V_source-port  string  `xml:"source-port"`
		} `xml:"dynamic-requests"`
	} `xml:"access>radsec>destination"`
}

// v_ is appended before every variable so it doesn't give any conflict
// with any keyword in golang. ex- interface is keyword in golang
func junosAccessRadsecDestinationDynamic__RequestsCreate(d *schema.ResourceData, m interface{}) error {

    pcfg := m.(*ProviderConfig)
    client, err := pcfg.Client()
    check(err)

    id := d.Get("resource_name").(string)
     	V_name := d.Get("name").(string)
	V_apply-groups := d.Get("apply-groups").(string)
	V_apply-groups-except := d.Get("apply-groups-except").(string)
	V_name__1 := d.Get("name__1").(string)
	V_name__2 := d.Get("name__2").(string)
	V_value := d.Get("value").(string)
	V_source-address := d.Get("source-address").(string)
	V_source-port := d.Get("source-port").(string)
	commit := true

	config := xmlAccessRadsecDestinationDynamic__Requests{}
	config.V_destination.V_name = V_name
	config.V_destination.V_dynamic-requests.V_apply-groups = V_apply-groups
	config.V_destination.V_dynamic-requests.V_apply-groups-except = V_apply-groups-except
	config.V_destination.V_dynamic-requests.V_apply-macro.V_name__1 = V_name__1
	config.V_destination.V_dynamic-requests.V_apply-macro.V_data.V_name__2 = V_name__2
	config.V_destination.V_dynamic-requests.V_apply-macro.V_data.V_value = V_value
	config.V_destination.V_dynamic-requests.V_source-address = V_source-address
	config.V_destination.V_dynamic-requests.V_source-port = V_source-port

    err = client.SendTransaction("", config, commit)
    check(err)
    
    d.SetId(fmt.Sprintf("%s_%s", pcfg.Cfg.Host, id))
    
    err = client.Close()
    check(err)
    
	return junosAccessRadsecDestinationDynamic__RequestsRead(d,m)
}

func junosAccessRadsecDestinationDynamic__RequestsRead(d *schema.ResourceData, m interface{}) error {

    pcfg := m.(*ProviderConfig)
    client, err := pcfg.Client()
    check(err)

    id := d.Get("resource_name").(string)
    
	config := &xmlAccessRadsecDestinationDynamic__Requests{}

	err = client.MarshalGroup(id, config)
	check(err)
 	d.Set("name", config.V_destination.V_name)
	d.Set("apply-groups", config.V_destination.V_dynamic-requests.V_apply-groups)
	d.Set("apply-groups-except", config.V_destination.V_dynamic-requests.V_apply-groups-except)
	d.Set("name__1", config.V_destination.V_dynamic-requests.V_apply-macro.V_name__1)
	d.Set("name__2", config.V_destination.V_dynamic-requests.V_apply-macro.V_data.V_name__2)
	d.Set("value", config.V_destination.V_dynamic-requests.V_apply-macro.V_data.V_value)
	d.Set("source-address", config.V_destination.V_dynamic-requests.V_source-address)
	d.Set("source-port", config.V_destination.V_dynamic-requests.V_source-port)

    err = client.Close()
    check(err)
    
	return nil
}

func junosAccessRadsecDestinationDynamic__RequestsUpdate(d *schema.ResourceData, m interface{}) error {

    pcfg := m.(*ProviderConfig)
    client, err := pcfg.Client()
    check(err)

    id := d.Get("resource_name").(string)
     	V_name := d.Get("name").(string)
	V_apply-groups := d.Get("apply-groups").(string)
	V_apply-groups-except := d.Get("apply-groups-except").(string)
	V_name__1 := d.Get("name__1").(string)
	V_name__2 := d.Get("name__2").(string)
	V_value := d.Get("value").(string)
	V_source-address := d.Get("source-address").(string)
	V_source-port := d.Get("source-port").(string)
	commit := true

	config := xmlAccessRadsecDestinationDynamic__Requests{}
	config.V_destination.V_name = V_name
	config.V_destination.V_dynamic-requests.V_apply-groups = V_apply-groups
	config.V_destination.V_dynamic-requests.V_apply-groups-except = V_apply-groups-except
	config.V_destination.V_dynamic-requests.V_apply-macro.V_name__1 = V_name__1
	config.V_destination.V_dynamic-requests.V_apply-macro.V_data.V_name__2 = V_name__2
	config.V_destination.V_dynamic-requests.V_apply-macro.V_data.V_value = V_value
	config.V_destination.V_dynamic-requests.V_source-address = V_source-address
	config.V_destination.V_dynamic-requests.V_source-port = V_source-port

    err = client.SendTransaction(id, config, commit)
    check(err)
    
    err = client.Close()
    check(err)
    
	return junosAccessRadsecDestinationDynamic__RequestsRead(d,m)
}

func junosAccessRadsecDestinationDynamic__RequestsDelete(d *schema.ResourceData, m interface{}) error {

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

func junosAccessRadsecDestinationDynamic__Requests() *schema.Resource {
	return &schema.Resource{
		Create: junosAccessRadsecDestinationDynamic__RequestsCreate,
		Read: junosAccessRadsecDestinationDynamic__RequestsRead,
		Update: junosAccessRadsecDestinationDynamic__RequestsUpdate,
		Delete: junosAccessRadsecDestinationDynamic__RequestsDelete,

        Schema: map[string]*schema.Schema{
            "resource_name": &schema.Schema{
                Type:    schema.TypeString,
                Required: true,
            },
			"name": &schema.Schema{
				Type:    schema.TypeString,
				Optional: true,
				Description:    "xpath is: config.V_destination",
			},
			"apply-groups": &schema.Schema{
				Type:    schema.TypeString,
				Optional: true,
				Description:    "xpath is: config.V_destination.V_dynamic-requests. Groups from which to inherit configuration data",
			},
			"apply-groups-except": &schema.Schema{
				Type:    schema.TypeString,
				Optional: true,
				Description:    "xpath is: config.V_destination.V_dynamic-requests. Don't inherit configuration data from these groups",
			},
			"name__1": &schema.Schema{
				Type:    schema.TypeString,
				Optional: true,
				Description:    "xpath is: config.V_destination.V_dynamic-requests.V_apply-macro. Name of the macro to be expanded",
			},
			"name__2": &schema.Schema{
				Type:    schema.TypeString,
				Optional: true,
				Description:    "xpath is: config.V_destination.V_dynamic-requests.V_apply-macro.V_data. Keyword part of the keyword-value pair",
			},
			"value": &schema.Schema{
				Type:    schema.TypeString,
				Optional: true,
				Description:    "xpath is: config.V_destination.V_dynamic-requests.V_apply-macro.V_data. Value part of the keyword-value pair",
			},
			"source-address": &schema.Schema{
				Type:    schema.TypeString,
				Optional: true,
				Description:    "xpath is: config.V_destination.V_dynamic-requests. Source address for dynamic requests",
			},
			"source-port": &schema.Schema{
				Type:    schema.TypeString,
				Optional: true,
				Description:    "xpath is: config.V_destination.V_dynamic-requests. Source port for dynamic requests",
			},
		},
	}
}