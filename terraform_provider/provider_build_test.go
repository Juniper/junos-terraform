package main

import (
	"errors"
	"testing"

	"terraform_provider/netconf"

	"github.com/hashicorp/terraform-plugin-framework/types"
)

type fakeNetconfClient struct{}

// Close implements netconf.Client for unit tests.
func (f *fakeNetconfClient) Close() error { return nil }

// DeleteConfig implements netconf.Client for unit tests.
func (f *fakeNetconfClient) DeleteConfig(string, bool) (string, error) { return "", nil }

// SendCommit implements netconf.Client for unit tests.
func (f *fakeNetconfClient) SendCommit() error { return nil }

// MarshalGroup implements netconf.Client for unit tests.
func (f *fakeNetconfClient) MarshalGroup(string, interface{}) error { return nil }

// MarshalConfig implements netconf.Client for unit tests.
func (f *fakeNetconfClient) MarshalConfig(interface{}) error { return nil }

// SendTransaction implements netconf.Client for unit tests.
func (f *fakeNetconfClient) SendTransaction(string, interface{}, bool) error { return nil }

// SendDirectTransaction implements netconf.Client for unit tests.
func (f *fakeNetconfClient) SendDirectTransaction(interface{}, bool) error { return nil }

// SendUpdate implements netconf.Client for unit tests.
func (f *fakeNetconfClient) SendUpdate(string, string, bool) error { return nil }

// TestBuildProviderConfigSuccess verifies successful provider config construction.
func TestBuildProviderConfigSuccess(t *testing.T) {
	originalFactory := providerClientFactory
	t.Cleanup(func() { providerClientFactory = originalFactory })

	providerClientFactory = func(cfg *Config) (netconf.Client, error) {
		if cfg.Host != "127.0.0.1" || cfg.Port != 830 || cfg.Username != "user" {
			t.Fatalf("unexpected config passed to factory: %+v", cfg)
		}
		return &fakeNetconfClient{}, nil
	}

	model := providerModel{
		Host:     types.StringValue("127.0.0.1"),
		Username: types.StringValue("user"),
		Password: types.StringValue("pass"),
		Port:     types.Int64Value(830),
		SshKey:   types.StringValue(""),
	}

	cfg, err := buildProviderConfig(model)
	if err != nil {
		t.Fatalf("buildProviderConfig() returned error: %v", err)
	}
	if cfg.Host != "127.0.0.1" {
		t.Fatalf("unexpected host: %q", cfg.Host)
	}
	if cfg.Client == nil {
		t.Fatalf("expected non-nil netconf client")
	}
}

// TestBuildProviderConfigFactoryError verifies client factory errors are returned.
func TestBuildProviderConfigFactoryError(t *testing.T) {
	originalFactory := providerClientFactory
	t.Cleanup(func() { providerClientFactory = originalFactory })

	expectedErr := errors.New("boom")
	providerClientFactory = func(_ *Config) (netconf.Client, error) {
		return nil, expectedErr
	}

	_, err := buildProviderConfig(providerModel{})
	if !errors.Is(err, expectedErr) {
		t.Fatalf("expected %v, got %v", expectedErr, err)
	}
}
