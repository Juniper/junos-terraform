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
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func junosQFXCommitCreate(d *schema.ResourceData, m interface{}) error {

	id := d.Get("resource_name").(string)

	pcfg := m.(*ProviderConfig)
	client, err := pcfg.Client()

	if err != nil {
		return err
	}

	//client := m.(*ProviderConfig)

	err = client.SendCommit()

	if err != nil {
		return err
	}

	d.SetId(fmt.Sprintf("%s_%s", pcfg.Cfg.Host, id))

	err = client.Close()
	if err != nil {
		return err
	}

	return junosQFXCommitRead(d, m)
}

func junosQFXCommitRead(d *schema.ResourceData, m interface{}) error {

	return nil
}

func junosQFXCommitUpdate(d *schema.ResourceData, m interface{}) error {

	return nil
}

func junosQFXCommitDelete(d *schema.ResourceData, m interface{}) error {

	return nil
}

func junosQFXCommit() *schema.Resource {
	return &schema.Resource{
		Create: junosQFXCommitCreate,
		Read:   junosQFXCommitRead,
		Update: junosQFXCommitUpdate,
		Delete: junosQFXCommitDelete,

		Schema: map[string]*schema.Schema{
			"resource_name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
		},
	}
}
