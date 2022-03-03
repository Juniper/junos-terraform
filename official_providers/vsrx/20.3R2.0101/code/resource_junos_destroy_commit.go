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
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func junosDestroyCommitCreate(d *schema.ResourceData, m interface{}) error {

	// var err error
	id := d.Get("resource_name").(string)

	client := m.(*ProviderConfig)

	d.SetId(fmt.Sprintf("%s_%s", client.Host, id))

	// err = client.Close()
	// if err != nil {
	// 	return err
	// }

	return nil
}

func junosDestroyCommitRead(d *schema.ResourceData, m interface{}) error {

	return nil
}

func junosDestroyCommitUpdate(d *schema.ResourceData, m interface{}) error {

	return nil
}

func junosDestroyCommitDelete(d *schema.ResourceData, m interface{}) error {

	var err error

	client := m.(*ProviderConfig)

	err = client.SendCommit()

	if err != nil {
		return err
	}

	err = client.Close()
	if err != nil {
		return err
	}

	return nil
}

func junosDestroyCommit() *schema.Resource {
	return &schema.Resource{
		Create: junosDestroyCommitCreate,
		Read:   junosDestroyCommitRead,
		Update: junosDestroyCommitUpdate,
		Delete: junosDestroyCommitDelete,

		Schema: map[string]*schema.Schema{
			"resource_name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
		},
	}
}
