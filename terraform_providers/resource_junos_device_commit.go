// Copyright (c) 2017-2022, Juniper Networks Inc. All rights reserved.
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
	"context"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Collects the objects from the .tf file
type deviceCommitModel struct {
	ResourceName types.String `tfsdk:"resource_name"`
}

// Collects the data for the crud work
type resourceDeviceCommit struct {
	client ProviderConfig
}

var _ resource.ResourceWithConfigure = new(resourceDeviceCommit)

func (r *resourceDeviceCommit) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	r.client = req.ProviderData.(ProviderConfig)
}

// Metadata implements resource.Resource.
func (r *resourceDeviceCommit) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_JunosDeviceCommit"
}

// Schema implements resource.Resource.
func (r *resourceDeviceCommit) Schema(_ context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"resource_name": schema.StringAttribute{
				Required:      true,
				PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplace()},
			},
		},
	}
}

func (r *resourceDeviceCommit) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {

	var plan deviceCommitModel
	//id := plan.ResourceName.ValueString()

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)

	if err := r.client.SendCommit(); err != nil {

	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)

	if err := r.client.Close(); err != nil {

	}

}

func (r *resourceDeviceCommit) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {

}

func (r *resourceDeviceCommit) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {

}

func (r *resourceDeviceCommit) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {

}
