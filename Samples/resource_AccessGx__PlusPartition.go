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
type xmlAccessGx__PlusPartition struct {
	XMLName xml.Name `xml:"configuration"`
	Groups  struct {
		XMLName	xml.Name	`xml:"groups"`
		Name	string	`xml:"name"`
		V_partition  struct {
			XMLName xml.Name `xml:"partition"`
			V_name  string  `xml:"name"`
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
			V_diameter-instance  string  `xml:"diameter-instance"`
			V_destination-realm  string  `xml:"destination-realm"`
			V_destination-host  string  `xml:"destination-host"`
		} `xml:"access>gx-plus>partition"`
	} `xml:"groups"`
	ApplyGroup string `xml:"apply-groups"`
}

// v_ is appended before every variable so it doesn't give any conflict
// with any keyword in golang. ex- interface is keyword in golang
func junosAccessGx__PlusPartitionCreate(d *schema.ResourceData, m interface{}) error {

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
	V_diameter-instance := d.Get("diameter-instance").(string)
	V_destination-realm := d.Get("destination-realm").(string)
	V_destination-host := d.Get("destination-host").(string)
	commit := true

	config := xmlAccessGx__PlusPartition{}
	config.ApplyGroup = id
	config.Groups.Name = id
	config.Groups.V_partition.V_name = V_name
	config.Groups.V_partition.V_apply-groups = V_apply-groups
	config.Groups.V_partition.V_apply-groups-except = V_apply-groups-except
	config.Groups.V_partition.V_apply-macro.V_name__1 = V_name__1
	config.Groups.V_partition.V_apply-macro.V_data.V_name__2 = V_name__2
	config.Groups.V_partition.V_apply-macro.V_data.V_value = V_value
	config.Groups.V_partition.V_diameter-instance = V_diameter-instance
	config.Groups.V_partition.V_destination-realm = V_destination-realm
	config.Groups.V_partition.V_destination-host = V_destination-host

    err = client.SendTransaction("", config, commit)
    check(err)
    
    d.SetId(fmt.Sprintf("%s_%s", pcfg.Cfg.Host, id))
    
    err = client.Close()
    check(err)
    
	return junosAccessGx__PlusPartitionRead(d,m)
}

func junosAccessGx__PlusPartitionRead(d *schema.ResourceData, m interface{}) error {

    pcfg := m.(*ProviderConfig)
    client, err := pcfg.Client()
    check(err)

    id := d.Get("resource_name").(string)
    
	config := &xmlAccessGx__PlusPartition{}

	err = client.MarshalGroup(id, config)
	check(err)
 	d.Set("name", config.Groups.V_partition.V_name)
	d.Set("apply-groups", config.Groups.V_partition.V_apply-groups)
	d.Set("apply-groups-except", config.Groups.V_partition.V_apply-groups-except)
	d.Set("name__1", config.Groups.V_partition.V_apply-macro.V_name__1)
	d.Set("name__2", config.Groups.V_partition.V_apply-macro.V_data.V_name__2)
	d.Set("value", config.Groups.V_partition.V_apply-macro.V_data.V_value)
	d.Set("diameter-instance", config.Groups.V_partition.V_diameter-instance)
	d.Set("destination-realm", config.Groups.V_partition.V_destination-realm)
	d.Set("destination-host", config.Groups.V_partition.V_destination-host)

    err = client.Close()
    check(err)
    
	return nil
}

func junosAccessGx__PlusPartitionUpdate(d *schema.ResourceData, m interface{}) error {

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
	V_diameter-instance := d.Get("diameter-instance").(string)
	V_destination-realm := d.Get("destination-realm").(string)
	V_destination-host := d.Get("destination-host").(string)
	commit := true

	config := xmlAccessGx__PlusPartition{}
	config.ApplyGroup = id
	config.Groups.Name = id
	config.Groups.V_partition.V_name = V_name
	config.Groups.V_partition.V_apply-groups = V_apply-groups
	config.Groups.V_partition.V_apply-groups-except = V_apply-groups-except
	config.Groups.V_partition.V_apply-macro.V_name__1 = V_name__1
	config.Groups.V_partition.V_apply-macro.V_data.V_name__2 = V_name__2
	config.Groups.V_partition.V_apply-macro.V_data.V_value = V_value
	config.Groups.V_partition.V_diameter-instance = V_diameter-instance
	config.Groups.V_partition.V_destination-realm = V_destination-realm
	config.Groups.V_partition.V_destination-host = V_destination-host

    err = client.SendTransaction(id, config, commit)
    check(err)
    
    err = client.Close()
    check(err)
    
	return junosAccessGx__PlusPartitionRead(d,m)
}

func junosAccessGx__PlusPartitionDelete(d *schema.ResourceData, m interface{}) error {

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

func junosAccessGx__PlusPartition() *schema.Resource {
	return &schema.Resource{
		Create: junosAccessGx__PlusPartitionCreate,
		Read: junosAccessGx__PlusPartitionRead,
		Update: junosAccessGx__PlusPartitionUpdate,
		Delete: junosAccessGx__PlusPartitionDelete,

        Schema: map[string]*schema.Schema{
            "resource_name": &schema.Schema{
                Type:    schema.TypeString,
                Required: true,
            },
			"name": &schema.Schema{
				Type:    schema.TypeString,
				Optional: true,
				Description:    "xpath is: config.Groups.V_partition. GX-PLUS partition name",
			},
			"apply-groups": &schema.Schema{
				Type:    schema.TypeString,
				Optional: true,
				Description:    "xpath is: config.Groups.V_partition. Groups from which to inherit configuration data",
			},
			"apply-groups-except": &schema.Schema{
				Type:    schema.TypeString,
				Optional: true,
				Description:    "xpath is: config.Groups.V_partition. Don't inherit configuration data from these groups",
			},
			"name__1": &schema.Schema{
				Type:    schema.TypeString,
				Optional: true,
				Description:    "xpath is: config.Groups.V_partition.V_apply-macro. Name of the macro to be expanded",
			},
			"name__2": &schema.Schema{
				Type:    schema.TypeString,
				Optional: true,
				Description:    "xpath is: config.Groups.V_partition.V_apply-macro.V_data. Keyword part of the keyword-value pair",
			},
			"value": &schema.Schema{
				Type:    schema.TypeString,
				Optional: true,
				Description:    "xpath is: config.Groups.V_partition.V_apply-macro.V_data. Value part of the keyword-value pair",
			},
			"diameter-instance": &schema.Schema{
				Type:    schema.TypeString,
				Optional: true,
				Description:    "xpath is: config.Groups.V_partition. GX-PLUS diameter instance",
			},
			"destination-realm": &schema.Schema{
				Type:    schema.TypeString,
				Optional: true,
				Description:    "xpath is: config.Groups.V_partition. GX-PLUS destination realm",
			},
			"destination-host": &schema.Schema{
				Type:    schema.TypeString,
				Optional: true,
				Description:    "xpath is: config.Groups.V_partition. GX-PLUS destination host",
			},
		},
	}
}