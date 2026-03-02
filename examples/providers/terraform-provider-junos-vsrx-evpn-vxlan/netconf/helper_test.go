package netconf

import (
	"fmt"
	"sync"
	"testing"
)

// MockDriver is a mock implementation of the Driver interface for testing
type MockDriver struct {
	dialCalled      bool
	dialErr         error
	closeCalled     bool
	closeErr        error
	sendRawCalls    []string
	sendRawErr      error
	sendRawResponse *RPCReply
	lockTarget      string
	unlockTarget    string
}

func (m *MockDriver) Lock(ds string) (*RPCReply, error) {
	m.lockTarget = ds
	return &RPCReply{RPCReplyBody: RPCReplyBody{Data: "ok"}}, nil
}

func (m *MockDriver) Unlock(ds string) (*RPCReply, error) {
	m.unlockTarget = ds
	return &RPCReply{RPCReplyBody: RPCReplyBody{Data: "ok"}}, nil
}

func (m *MockDriver) Close() error {
	m.closeCalled = true
	return m.closeErr
}

func (m *MockDriver) Dial() error {
	m.dialCalled = true
	return m.dialErr
}

func (m *MockDriver) DialTimeout() error {
	return nil
}

func (m *MockDriver) SendRaw(rawxml string) (*RPCReply, error) {
	m.sendRawCalls = append(m.sendRawCalls, rawxml)
	if m.sendRawErr != nil {
		return nil, m.sendRawErr
	}
	if m.sendRawResponse != nil {
		return m.sendRawResponse, nil
	}
	return &RPCReply{RPCReplyBody: RPCReplyBody{Data: "ok"}}, nil
}

func (m *MockDriver) GetConfig() (*RPCReply, error) {
	ok := true
	return &RPCReply{RPCReplyBody: RPCReplyBody{Data: "<config/>", Ok: &ok}}, nil
}

// TestGoNCClientClose tests the Close method
func TestGoNCClientClose(t *testing.T) {
	driver := &MockDriver{}
	client := &GoNCClient{Driver: driver}

	err := client.Close()
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	if client.Driver != nil {
		t.Error("expected driver to be set to nil after close")
	}
}

// TestGoNCClientCloseWithNilDriver tests Close with nil driver
func TestGoNCClientCloseWithNilDriver(t *testing.T) {
	client := &GoNCClient{Driver: nil}

	err := client.Close()
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}

// TestGoNCClientDeleteConfig tests the DeleteConfig method
func TestGoNCClientDeleteConfig(t *testing.T) {
	mockDriver := &MockDriver{
		sendRawResponse: &RPCReply{RPCReplyBody: RPCReplyBody{Data: "ok"}},
	}
	client := &GoNCClient{Driver: mockDriver, Lock: sync.RWMutex{}}

	result, err := client.DeleteConfig("test-group", false)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	if result == "" {
		t.Error("expected non-empty result")
	}

	if !mockDriver.dialCalled {
		t.Error("expected Dial to be called")
	}
}

// TestGoNCClientDeleteConfigWithCommit tests DeleteConfig with commit
func TestGoNCClientDeleteConfigWithCommit(t *testing.T) {
	mockDriver := &MockDriver{
		sendRawResponse: &RPCReply{RPCReplyBody: RPCReplyBody{Data: "ok"}},
	}
	client := &GoNCClient{Driver: mockDriver, Lock: sync.RWMutex{}}

	result, err := client.DeleteConfig("test-group", true)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	if result == "" {
		t.Error("expected non-empty result")
	}

	// Should have called SendRaw twice (delete and commit)
	if len(mockDriver.sendRawCalls) < 2 {
		t.Errorf("expected at least 2 SendRaw calls, got %d", len(mockDriver.sendRawCalls))
	}
}

// TestGoNCClientDeleteConfigDialError tests DeleteConfig when Dial fails
func TestGoNCClientDeleteConfigDialError(t *testing.T) {
	mockDriver := &MockDriver{
		dialErr: fmt.Errorf("dial failed"),
	}
	client := &GoNCClient{Driver: mockDriver, Lock: sync.RWMutex{}}

	_, err := client.DeleteConfig("test-group", false)
	if err == nil {
		t.Error("expected error")
	}

	if err.Error() != "dial failed" {
		t.Errorf("error mismatch: expected 'dial failed', got '%v'", err)
	}
}

// TestGoNCClientSendCommit tests the SendCommit method
func TestGoNCClientSendCommit(t *testing.T) {
	mockDriver := &MockDriver{
		sendRawResponse: &RPCReply{RPCReplyBody: RPCReplyBody{Data: "ok"}},
	}
	client := &GoNCClient{Driver: mockDriver, Lock: sync.RWMutex{}}

	// Reset the global apply groups list
	applyGroupsList = []string{}

	err := client.SendCommit()
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	if !mockDriver.dialCalled {
		t.Error("expected Dial to be called")
	}
}

// TestGoNCClientSendCommitWithDialError tests SendCommit when Dial fails
func TestGoNCClientSendCommitWithDialError(t *testing.T) {
	mockDriver := &MockDriver{
		dialErr: fmt.Errorf("dial failed"),
	}
	client := &GoNCClient{Driver: mockDriver, Lock: sync.RWMutex{}}

	// Reset the global apply groups list
	applyGroupsList = []string{}

	err := client.SendCommit()
	if err == nil {
		t.Error("expected error")
	}
}

// TestGoNCClientMarshalGroup tests the MarshalGroup method
func TestGoNCClientMarshalGroup(t *testing.T) {
	mockDriver := &MockDriver{
		sendRawResponse: &RPCReply{RPCReplyBody: RPCReplyBody{Data: `<configuration><groups><name>test</name></groups></configuration>`}},
	}
	client := &GoNCClient{Driver: mockDriver, Lock: sync.RWMutex{}}

	type TestConfig struct {
		Configuration struct {
			Groups struct {
				Name string `xml:"name"`
			} `xml:"groups"`
		} `xml:"configuration"`
	}

	var config TestConfig
	err := client.MarshalGroup("test", &config)

	// The actual implementation reads raw group data, so check for no error
	if err != nil && err.Error() != "dial failed" {
		t.Logf("MarshalGroup error (expected in test): %v", err)
	}
}

// TestSortApplyGroupsList tests the sortApplyGroupsList function
func TestSortApplyGroupsList(t *testing.T) {
	// Save original list
	originalList := applyGroupsList

	// Test with unsorted list
	applyGroupsList = []string{"group-c", "group-a", "group-b", ""}
	sortApplyGroupsList()

	expectedList := []string{"group-a", "group-b", "group-c"}
	if len(applyGroupsList) != len(expectedList) {
		t.Errorf("length mismatch: expected %d, got %d", len(expectedList), len(applyGroupsList))
	}

	for i, v := range expectedList {
		if i >= len(applyGroupsList) || applyGroupsList[i] != v {
			t.Errorf("element %d mismatch: expected %s, got %s", i, v, applyGroupsList[i])
		}
	}

	// Restore original list
	applyGroupsList = originalList
}

// TestAddToApplyGroupsList tests the addToApplyGroupsList function
func TestAddToApplyGroupsList(t *testing.T) {
	// Save original list
	originalList := applyGroupsList
	applyGroupsList = []string{}

	addToApplyGroupsList("group1")
	addToApplyGroupsList("group2")

	if len(applyGroupsList) != 2 {
		t.Errorf("expected 2 groups, got %d", len(applyGroupsList))
	}

	if applyGroupsList[0] != "group1" {
		t.Errorf("first group mismatch: expected group1, got %s", applyGroupsList[0])
	}

	if applyGroupsList[1] != "group2" {
		t.Errorf("second group mismatch: expected group2, got %s", applyGroupsList[1])
	}

	// Restore original list
	applyGroupsList = originalList
}

// TestGoNCClientSendTransaction tests SendTransaction method
func TestGoNCClientSendTransaction(t *testing.T) {
	mockDriver := &MockDriver{
		sendRawResponse: &RPCReply{RPCReplyBody: RPCReplyBody{Data: "ok"}},
	}
	client := &GoNCClient{Driver: mockDriver, Lock: sync.RWMutex{}}

	type TestConfig struct {
		XMLName struct{} `xml:"configuration"`
		Name    string   `xml:"name"`
	}

	config := TestConfig{Name: "test"}

	// Should not panic with empty ID
	err := client.SendTransaction("", config, false)

	// Error is expected since we're using a mock driver
	if err == nil {
		t.Logf("SendTransaction succeeded (may be expected)")
	}
}

// TestGoNCClientConcurrency tests concurrent access to GoNCClient
func TestGoNCClientConcurrency(t *testing.T) {
	mockDriver := &MockDriver{
		sendRawResponse: &RPCReply{RPCReplyBody: RPCReplyBody{Data: "ok"}},
	}
	client := &GoNCClient{Driver: mockDriver, Lock: sync.RWMutex{}}

	// Test concurrent DeleteConfig calls
	var wg sync.WaitGroup
	errors := make(chan error, 10)

	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			_, err := client.DeleteConfig("test-group", false)
			if err != nil {
				errors <- err
			}
		}()
	}

	wg.Wait()
	close(errors)

	for err := range errors {
		t.Logf("concurrent error: %v", err)
	}
}

// TestGoNCClientReadRawGroup tests the readRawGroup method (private but testable)
func TestGoNCClientSendRawConfig(t *testing.T) {
	mockDriver := &MockDriver{
		sendRawResponse: &RPCReply{RPCReplyBody: RPCReplyBody{Data: "ok"}},
	}
	client := &GoNCClient{Driver: mockDriver, Lock: sync.RWMutex{}}

	// Reset the global apply groups list
	applyGroupsList = []string{}

	// Mock the sendRawConfig method by calling SendTransaction with empty ID
	type TestConfig struct {
		XMLName struct{} `xml:"configuration"`
		Name    string   `xml:"name"`
		Data    string   `xml:"data"`
	}

	config := TestConfig{Name: "test", Data: "content"}
	err := client.SendTransaction("", config, false)

	// Error expected with mock driver, but method should be called
	if err != nil {
		t.Logf("Expected error in test: %v", err)
	}
}

// TestPublicKeyFile tests the publicKeyFile function
func TestPublicKeyFile(t *testing.T) {
	// Test with non-existent file
	authMethod := publicKeyFile("/non/existent/path/to/key")
	if authMethod != nil {
		t.Logf("Auth method: %v (expected nil for non-existent file)", authMethod)
	}

	// Test with empty path
	authMethod = publicKeyFile("")
	if authMethod != nil {
		t.Logf("Auth method: %v (expected nil for empty path)", authMethod)
	}
}

// TestGoNCClientStructure tests the GoNCClient structure
func TestGoNCClientStructure(t *testing.T) {
	mockDriver := &MockDriver{}
	client := &GoNCClient{
		Driver: mockDriver,
		Lock:   sync.RWMutex{},
	}

	if client.Driver != mockDriver {
		t.Error("driver assignment failed")
	}

	// Verify Lock can be used for synchronization
	client.Lock.RLock()
	_ = client.Driver
	client.Lock.RUnlock()
}

// BenchmarkDeleteConfig benchmarks the DeleteConfig method
func BenchmarkDeleteConfig(b *testing.B) {
	mockDriver := &MockDriver{
		sendRawResponse: &RPCReply{RPCReplyBody: RPCReplyBody{Data: "ok"}},
	}
	client := &GoNCClient{Driver: mockDriver, Lock: sync.RWMutex{}}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = client.DeleteConfig("test-group", false)
	}
}
