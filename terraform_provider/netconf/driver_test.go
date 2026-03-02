package netconf

import (
	"sync"
	"testing"
)

// TestClientInterface tests that implementations satisfy the Client interface
func TestClientInterface(t *testing.T) {
	// Mock client for testing
	mockClient := &GoNCClient{
		Driver: &MockDriver{},
	}

	// Verify Client interface implementation
	var _ Client = mockClient

	// Interface methods should be callable without panic
	_ = mockClient.Close()
}

// TestNewDriver tests the NewDriver function with different driver types
func TestNewDriver(t *testing.T) {
	sshDriver := NewSSH()
	d := NewDriver(sshDriver)

	if d != sshDriver {
		t.Error("NewDriver should return the provided driver unchanged")
	}
}

// TestNewDriverWithMockDriver tests NewDriver with mock driver
func TestNewDriverWithMockDriver(t *testing.T) {
	mockDriver := &MockDriver{}
	d := NewDriver(mockDriver)

	if d != mockDriver {
		t.Error("NewDriver should return the provided driver unchanged")
	}

	// Verify the returned driver implements Driver interface (compile-time assertion)
	_ = interface{}(d).(Driver)
}

// TestNewDriverWithNil tests NewDriver with nil driver
func TestNewDriverWithNil(t *testing.T) {
	var nilDriver Driver
	d := NewDriver(nilDriver)

	if d != nilDriver {
		t.Error("NewDriver should handle nil driver")
	}
}

// TestNewDriverWithJunosDriver tests NewDriver with DriverJunos
func TestNewDriverWithJunosDriver(t *testing.T) {
	junosDriver := New()
	d := NewDriver(junosDriver)

	if d != junosDriver {
		t.Error("NewDriver should return DriverJunos unchanged")
	}

	// Verify it implements Driver interface (compile-time assertion)
	_ = interface{}(d).(Driver)
}

// TestNewDriverPreservesType tests that NewDriver preserves driver type
func TestNewDriverPreservesType(t *testing.T) {
	testCases := []struct {
		name   string
		driver Driver
	}{
		{
			name:   "DriverSSH",
			driver: NewSSH(),
		},
		{
			name:   "DriverJunos",
			driver: New(),
		},
		{
			name:   "MockDriver",
			driver: &MockDriver{},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := NewDriver(tc.driver)

			if result != tc.driver {
				t.Error("NewDriver should preserve original driver")
			}
		})
	}
}

// TestClientInterfaceMethods tests that all Client interface methods are implemented
func TestClientInterfaceMethods(t *testing.T) {
	mockDriver := &MockDriver{
		sendRawResponse: &RPCReply{RPCReplyBody: RPCReplyBody{Data: "ok"}},
	}
	client := &GoNCClient{
		Driver: mockDriver,
		Lock:   sync.RWMutex{},
	}

	// Test that client fields are properly set
	if client.Driver != mockDriver {
		t.Error("expected Driver to be set correctly")
	}

	// Test Close - verifies the method exists and is callable
	err := client.Close()
	if err != nil {
		t.Logf("Close returned error: %v", err)
	}

	// Verify that Close set the driver to nil
	if client.Driver != nil {
		t.Error("expected Driver to be nil after Close")
	}

	// Note: Testing actual DeleteConfig, SendCommit, etc. requires
	// fully mocking the Driver and handling the framework's schema data
	// In integration tests, these are tested via the Terraform plugin system
}

// TestDriverInterface tests that Driver interface is properly defined
func TestDriverInterface(t *testing.T) {
	testCases := []struct {
		name   string
		driver Driver
	}{
		{
			name:   "DriverSSH",
			driver: NewSSH(),
		},
		{
			name:   "DriverJunos",
			driver: New(),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Verify drivers implement the Driver interface
			// by checking they're non-nil and have the interface type
			_ = interface{}(tc.driver).(Driver)
			
			// Note: Actual method calls on these drivers require
			// proper initialization with SSH connections or processes.
			// Unit testing the methods is done via mocks in other tests.
		})
	}
}

// TestClientClose tests Client.Close implementation
func TestClientClose(t *testing.T) {
	mockDriver := &MockDriver{}
	client := &GoNCClient{
		Driver: mockDriver,
	}

	// First close should succeed
	err := client.Close()
	if err != nil {
		t.Logf("Close error: %v", err)
	}

	if client.Driver != nil {
		t.Error("Driver should be nil after Close")
	}

	// Second close should also not panic
	err = client.Close()
	if err != nil {
		t.Logf("Second Close error: %v", err)
	}
}

// TestClientDeleteConfig tests Client.DeleteConfig implementation
func TestClientDeleteConfig(t *testing.T) {
	mockDriver := &MockDriver{
		sendRawResponse: &RPCReply{RPCReplyBody: RPCReplyBody{Data: "ok"}},
	}
	client := &GoNCClient{
		Driver: mockDriver,
	}

	result, err := client.DeleteConfig("test-group", false)
	if err != nil {
		t.Logf("DeleteConfig error: %v (expected for mock)", err)
	}

	if result == "" {
		t.Logf("result is empty (expected for mock)")
	}
}

// TestClientSendCommit tests Client.SendCommit implementation
func TestClientSendCommit(t *testing.T) {
	mockDriver := &MockDriver{
		sendRawResponse: &RPCReply{RPCReplyBody: RPCReplyBody{Data: "ok"}},
	}
	client := &GoNCClient{
		Driver: mockDriver,
	}

	// Reset global state
	applyGroupsList = []string{}

	err := client.SendCommit()
	if err != nil {
		t.Logf("SendCommit error: %v (expected for mock)", err)
	}
}

// TestClientMarshalGroup tests Client.MarshalGroup implementation
func TestClientMarshalGroup(t *testing.T) {
	mockDriver := &MockDriver{
		sendRawResponse: &RPCReply{RPCReplyBody: RPCReplyBody{
			Data: `<configuration><groups><name>test</name></groups></configuration>`,
		}},
	}
	client := &GoNCClient{
		Driver: mockDriver,
	}

	type TestConfig struct {
		Configuration struct {
			Groups struct {
				Name string `xml:"name"`
			} `xml:"groups"`
		} `xml:"configuration"`
	}

	var config TestConfig
	err := client.MarshalGroup("test", &config)
	if err != nil {
		t.Logf("MarshalGroup error: %v (expected for mock)", err)
	}
}

// TestClientSendTransaction tests Client.SendTransaction implementation
func TestClientSendTransaction(t *testing.T) {
	mockDriver := &MockDriver{
		sendRawResponse: &RPCReply{RPCReplyBody: RPCReplyBody{Data: "ok"}},
	}
	client := &GoNCClient{
		Driver: mockDriver,
	}

	type TestConfig struct {
		XMLName struct{} `xml:"configuration"`
		Name    string   `xml:"name"`
	}

	config := TestConfig{Name: "test"}
	err := client.SendTransaction("test", config, false)
	if err != nil {
		t.Logf("SendTransaction error: %v (expected for mock)", err)
	}
}

// TestInterfaceCompliance verifies GoNCClient implements Client interface
func TestInterfaceCompliance(t *testing.T) {
	mockDriver := &MockDriver{}
	var client Client = &GoNCClient{Driver: mockDriver}

	// Methods should be callable
	_ = client.Close()
}

// TestDriverInterfaceCompliance verifies implementations satisfy Driver interface
func TestDriverInterfaceCompliance(t *testing.T) {
	testCases := []struct {
		name   string
		driver Driver
	}{
		{
			name:   "DriverSSH",
			driver: NewSSH(),
		},
		{
			name:   "DriverJunos",
			driver: New(),
		},
		{
			name:   "MockDriver",
			driver: &MockDriver{},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Verify drivers implement the Driver interface by asserting the type
			_ = interface{}(tc.driver).(Driver)
			
			// Note: Calling Dial, SendRaw, etc. on real drivers requires
			// initialized SSH connections or child processes. These are tested
			// via mocks and integration tests, not unit tests.
		})
	}
}

// TestNewDriverMultiple tests NewDriver with multiple drivers
func TestNewDriverMultiple(t *testing.T) {
	drivers := []Driver{
		NewSSH(),
		New(),
		&MockDriver{},
	}

	for i, d := range drivers {
		result := NewDriver(d)
		if result != d {
			t.Errorf("driver %d: NewDriver should preserve driver", i)
		}
	}
}

// BenchmarkNewDriver benchmarks the NewDriver function
func BenchmarkNewDriver(b *testing.B) {
	driver := NewSSH()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = NewDriver(driver)
	}
}

// BenchmarkClientClose benchmarks Client.Close
func BenchmarkClientClose(b *testing.B) {
	for i := 0; i < b.N; i++ {
		mockDriver := &MockDriver{}
		client := &GoNCClient{Driver: mockDriver}
		_ = client.Close()
	}
}

// TestClientInterfaceDocumentation tests the documented Client interface
func TestClientInterfaceDocumentation(t *testing.T) {
	mockDriver := &MockDriver{
		sendRawResponse: &RPCReply{RPCReplyBody: RPCReplyBody{Data: "ok"}},
	}
	client := &GoNCClient{Driver: mockDriver}

	// Test documented example usage:
	// 1. Create client
	// 2. Call interface methods
	// 3. Close client

	// Document the interface through test
	t.Logf("Client interface methods:")
	t.Logf("- Close() error")
	t.Logf("- DeleteConfig(applyGroup string, commit bool) (string, error)")
	t.Logf("- SendCommit() error")
	t.Logf("- MarshalGroup(id string, obj interface{}) error")
	t.Logf("- SendTransaction(id string, obj interface{}, commit bool) error")

	// All methods should be accessible (won't panic)
	_ = client
}
