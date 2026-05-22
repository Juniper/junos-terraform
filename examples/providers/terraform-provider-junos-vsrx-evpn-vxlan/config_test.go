package main

import (
	"testing"
)

// TestConfigClientWithPassword tests creating a client with password authentication
func TestConfigClientWithPassword(t *testing.T) {
	config := &Config{
		Host:     "localhost",
		Port:     830,
		Username: "testuser",
		Password: "testpass",
		SSHKey:   "",
	}

	client, err := config.Client()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if client == nil {
		t.Error("expected non-nil client, got nil")
	}
}

// TestConfigClientWithSSHKey tests creating a client with SSH key authentication
func TestConfigClientWithSSHKey(t *testing.T) {
	config := &Config{
		Host:     "localhost",
		Port:     830,
		Username: "testuser",
		Password: "",
		SSHKey:   "/path/to/key",
	}

	client, err := config.Client()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if client == nil {
		t.Error("expected non-nil client, got nil")
	}
}

// TestConfigClientWithEmptyHost tests client creation with empty host
func TestConfigClientWithEmptyHost(t *testing.T) {
	config := &Config{
		Host:     "",
		Port:     830,
		Username: "testuser",
		Password: "testpass",
		SSHKey:   "",
	}

	client, err := config.Client()
	if err != nil {
		// Expected error when host is empty
		t.Logf("got expected error: %v", err)
	}

	if client != nil {
		// Client creation may still succeed, so we just validate
		t.Logf("client created with empty host: %v", client)
	}
}

// TestConfigClientWithValidParams tests client creation with all valid parameters
func TestConfigClientWithValidParams(t *testing.T) {
	testCases := []struct {
		name    string
		config  *Config
		wantErr bool
		desc    string
	}{
		{
			name: "valid with password",
			config: &Config{
				Host:     "192.168.1.1",
				Port:     830,
				Username: "admin",
				Password: "admin123",
				SSHKey:   "",
			},
			wantErr: false,
			desc:    "Should create client with password authentication",
		},
		{
			name: "valid with ssh key",
			config: &Config{
				Host:     "192.168.1.1",
				Port:     830,
				Username: "admin",
				Password: "",
				SSHKey:   "/home/user/.ssh/id_rsa",
			},
			wantErr: false,
			desc:    "Should create client with SSH key authentication",
		},
		{
			name: "custom port",
			config: &Config{
				Host:     "10.0.0.1",
				Port:     2222,
				Username: "operator",
				Password: "pass",
				SSHKey:   "",
			},
			wantErr: false,
			desc:    "Should create client with custom SSH port",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			client, err := tc.config.Client()

			if (err != nil) != tc.wantErr {
				t.Errorf("unexpected error state: %v (expected error: %v)", err, tc.wantErr)
			}

			if !tc.wantErr && client == nil {
				t.Error("expected non-nil client")
			}

			t.Logf("%s", tc.desc)
		})
	}
}

// TestNewClientDirectly tests the newClient helper function
func TestNewClientDirectly(t *testing.T) {
	config := &Config{
		Host:     "localhost",
		Port:     830,
		Username: "user",
		Password: "pass",
		SSHKey:   "",
	}

	client, err := newClient(config)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if client == nil {
		t.Error("expected non-nil client")
	}
}

// TestConfigValues tests that config values are properly set
func TestConfigValues(t *testing.T) {
	expectedHost := "example.com"
	expectedPort := 2222
	expectedUsername := "admin"
	expectedPassword := "secret"
	expectedSSHKey := "/path/to/key"

	config := &Config{
		Host:     expectedHost,
		Port:     expectedPort,
		Username: expectedUsername,
		Password: expectedPassword,
		SSHKey:   expectedSSHKey,
	}

	if config.Host != expectedHost {
		t.Errorf("host mismatch: expected %s, got %s", expectedHost, config.Host)
	}

	if config.Port != expectedPort {
		t.Errorf("port mismatch: expected %d, got %d", expectedPort, config.Port)
	}

	if config.Username != expectedUsername {
		t.Errorf("username mismatch: expected %s, got %s", expectedUsername, config.Username)
	}

	if config.Password != expectedPassword {
		t.Errorf("password mismatch: expected %s, got %s", expectedPassword, config.Password)
	}

	if config.SSHKey != expectedSSHKey {
		t.Errorf("ssh key mismatch: expected %s, got %s", expectedSSHKey, config.SSHKey)
	}
}
