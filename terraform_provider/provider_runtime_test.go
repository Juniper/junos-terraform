package main

import (
	"context"
	"errors"
	"math/big"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
	"terraform_provider/netconf"
)

// TestProviderMetadata verifies provider type metadata is set correctly.
func TestProviderMetadata(t *testing.T) {
	p := &Provider{}
	resp := &provider.MetadataResponse{}
	p.Metadata(context.Background(), provider.MetadataRequest{}, resp)
	if resp.TypeName != "terraform_provider" {
		t.Fatalf("unexpected provider type name: %q", resp.TypeName)
	}
}

// TestProviderConfigureMissingConfigPanics documents current framework panic behavior.
func TestProviderConfigureMissingConfigPanics(t *testing.T) {
	p := &Provider{}
	defer func() {
		if recover() == nil {
			t.Fatalf("expected panic when ConfigureRequest.Config is unset")
		}
	}()

	resp := &provider.ConfigureResponse{}
	p.Configure(context.Background(), provider.ConfigureRequest{}, resp)
	if resp.Diagnostics.HasError() {
		t.Fatalf("unexpected diagnostics before panic")
	}
}

// providerConfigRaw creates a typed provider config payload for Configure tests.
func providerConfigRaw(t *testing.T, p *Provider) tfsdk.Config {
	t.Helper()

	ctx := context.Background()
	schemaResp := &provider.SchemaResponse{}
	p.Schema(ctx, provider.SchemaRequest{}, schemaResp)

	tfType := schemaResp.Schema.Type().TerraformType(ctx)
	raw := tftypes.NewValue(tfType, map[string]tftypes.Value{
		"host":     tftypes.NewValue(tftypes.String, "127.0.0.1"),
		"username": tftypes.NewValue(tftypes.String, "user"),
		"password": tftypes.NewValue(tftypes.String, "pass"),
		"port":     tftypes.NewValue(tftypes.Number, big.NewFloat(830)),
		"sshkey":   tftypes.NewValue(tftypes.String, ""),
	})

	return tfsdk.Config{Schema: schemaResp.Schema, Raw: raw}
}

// TestProviderConfigureSuccess verifies successful provider configuration.
func TestProviderConfigureSuccess(t *testing.T) {
	originalFactory := providerClientFactory
	t.Cleanup(func() { providerClientFactory = originalFactory })

	providerClientFactory = func(_ *Config) (netconf.Client, error) {
		return &fakeNetconfClient{}, nil
	}

	p := &Provider{}
	req := provider.ConfigureRequest{Config: providerConfigRaw(t, p)}
	resp := &provider.ConfigureResponse{}

	p.Configure(context.Background(), req, resp)
	if resp.Diagnostics.HasError() {
		t.Fatalf("unexpected diagnostics: %v", resp.Diagnostics)
	}
	if resp.ResourceData == nil {
		t.Fatalf("expected provider resource data")
	}
}

// TestProviderConfigureClientError verifies provider diagnostics on client creation failures.
func TestProviderConfigureClientError(t *testing.T) {
	originalFactory := providerClientFactory
	t.Cleanup(func() { providerClientFactory = originalFactory })

	providerClientFactory = func(_ *Config) (netconf.Client, error) {
		return nil, errors.New("client failed")
	}

	p := &Provider{}
	req := provider.ConfigureRequest{Config: providerConfigRaw(t, p)}
	resp := &provider.ConfigureResponse{}

	p.Configure(context.Background(), req, resp)
	if !resp.Diagnostics.HasError() {
		t.Fatalf("expected diagnostics on client factory error")
	}
}
