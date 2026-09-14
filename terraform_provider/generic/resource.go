package generic

import (
	"context"
	"fmt"

	"terraform_provider/netconf"
	"terraform_provider/patch"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
)

// ConfigResource is the generic schema-driven resource.
type ConfigResource struct {
	client     netconf.Client
	host       string
	idx        map[string]*patch.NodeInfo
	nodes      []patch.SchemaNode
	tfSchema   schema.Schema
	schemaJSON string
}

// NewConfigResource creates a ConfigResource from pre-loaded schema data.
func NewConfigResource(idx map[string]*patch.NodeInfo, nodes []patch.SchemaNode, schemaJSON string) *ConfigResource {
	return &ConfigResource{
		idx:        idx,
		nodes:      nodes,
		tfSchema:   BuildSchema(nodes),
		schemaJSON: schemaJSON,
	}
}

// ProviderData is the interface the generic resource expects from provider configuration.
type ProviderData interface {
	GetClient() netconf.Client
	GetHost() string
}

func (r *ConfigResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	pd, ok := req.ProviderData.(ProviderData)
	if !ok {
		resp.Diagnostics.AddError("Unexpected provider data type",
			fmt.Sprintf("Expected generic.ProviderData, got %T", req.ProviderData))
		return
	}
	r.client = pd.GetClient()
	r.host = pd.GetHost()
}

func (r *ConfigResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = "terraform-provider-" + req.ProviderTypeName
}

func (r *ConfigResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = r.tfSchema
}

func (r *ConfigResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	resp.State = tfsdk.State{
		Schema: req.Plan.Schema,
		Raw:    req.Plan.Raw.Copy(),
	}
}

func (r *ConfigResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Preserve existing state — read-back from device not yet wired.
}

func (r *ConfigResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	resp.State = tfsdk.State{
		Schema: req.Plan.Schema,
		Raw:    req.Plan.Raw.Copy(),
	}
}

func (r *ConfigResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
}
