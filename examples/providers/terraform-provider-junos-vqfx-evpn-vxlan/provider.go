package main

import (
	"context"

	"terraform-provider-junos-vqfx-evpn-vxlan/netconf"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var _ provider.Provider = new(Provider)

func newProvider() provider.Provider {
	return Provider{}
}

type Provider struct {
}

type providerModel struct {
	Host     types.String `tfsdk:"host"`
	Username types.String `tfsdk:"username"`
	Password types.String `tfsdk:"password"`
	Port     types.Int64  `tfsdk:"port"`
	SshKey   types.String `tfsdk:"sshkey"`
}

// ProviderConfig is to hold client information
type ProviderConfig struct {
	netconf.Client
	Host string
}

// Configure implements provider.Provider.
func (p Provider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	var config providerModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	clientConfig := Config{
		Host:     config.Host.ValueString(),
		Port:     int(config.Port.ValueInt64()),
		Username: config.Username.ValueString(),
		Password: config.Password.ValueString(),
		SSHKey:   config.SshKey.ValueString(),
	}

	client, err := clientConfig.Client()
	if err != nil {
		resp.Diagnostics.AddError("failed to create client", err.Error())
		return
	}

	var providerConfig ProviderConfig
	providerConfig.Client = client
	providerConfig.Host = config.Host.ValueString()

	resp.ResourceData = providerConfig
}

// DataSources implements provider.Provider.
func (p Provider) DataSources(_ context.Context) []func() datasource.DataSource {
	return nil
}

// Metadata implements provider.Provider.
func (p Provider) Metadata(_ context.Context, _ provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "junos-vqfx-evpn-vxlan"
}

// Resources implements provider.Provider.
func (p Provider) Resources(_ context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		func() resource.Resource { return new(resource_Apply_Groups) },
	}
}

// Schema implements provider.Provider.
func (p Provider) Schema(_ context.Context, _ provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"host": schema.StringAttribute{
				Required: true,
			},
			"username": schema.StringAttribute{
				Required: true,
			},
			"password": schema.StringAttribute{
				Optional: true,
			},
			"port": schema.Int64Attribute{
				Required: true,
				//Optional: true,
				//Computed: true,
				//Default: int64default.StaticInt64(22),
			},
			"sshkey": schema.StringAttribute{
				Optional: true,
				// Will need to add eventually
				//Validators: []validator.String{stringvalidator.AtLeastOneOf(path.MatchRoot("d")...)},
			},
		},
	}
}