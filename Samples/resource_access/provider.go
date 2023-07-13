
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
	"log"

	gonetconf "github.com/vinpatel24/go-netconf/helpers/junos_helpers"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

// ProviderConfig is to hold client information
type ProviderConfig struct {
	Cfg Config
}

func check(e error) {
    if e != nil {
        panic(e)
    }
}

// Client Refreshes the client
func (PC *ProviderConfig) Client() (*gonetconf.GoNCClient, error) {

	client, err := PC.Cfg.Client()
	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	return client, nil
}

func returnProvider(d *schema.ResourceData) (interface{}, error) {
	config := Config{
		Host:     d.Get("host").(string),
		Port:     d.Get("port").(int),
		Username: d.Get("username").(string),
		Password: d.Get("password").(string),
		SSHKey:   d.Get("sshkey").(string),
	}

	return &ProviderConfig{config}, nil
}

// Provider returns a Terraform ResourceProvider.
func Provider() *schema.Provider {
	return &schema.Provider{

		Schema: map[string]*schema.Schema{
			"host": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},

			"port": &schema.Schema{
				Type:     schema.TypeInt,
				Required: true,
			},

			"username": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},

			"password": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"sshkey": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
		},

		ResourcesMap: map[string]*schema.Resource{
			"junos-device_AccessDomainMapApply__MacroData": junosAccessDomainMapApply__MacroData(),
			"junos-device_AccessGx__PlusPartition": junosAccessGx__PlusPartition(),
			"junos-device_AccessRadsecDestinationDynamic__Requests": junosAccessRadsecDestinationDynamic__Requests(),
			"junos-device_commit": junosCommit(),
		},
		ConfigureFunc: returnProvider,
	}
}
