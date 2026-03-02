package netconf

import (
	"bytes"
	"io"
	"testing"
)

// MockTransport is a mock implementation of Transport for testing
type MockTransport struct {
	ReceiveHelloCalled bool
	HelloResponse      *HelloMessage
	HelloErr           error

	SendHelloCalled bool
	SendHelloErr    error

	CloseCalled bool
	CloseErr    error

	SendCalled  bool
	SendErr     error
	SendData    []byte

	ReceiveCalled bool
	ReceiveErr    error
	ReceiveData   []byte
}

func (m *MockTransport) ReceiveHello() (*HelloMessage, error) {
	m.ReceiveHelloCalled = true
	if m.HelloErr != nil {
		return nil, m.HelloErr
	}
	if m.HelloResponse != nil {
		return m.HelloResponse, nil
	}
	return &HelloMessage{
		SessionID: 12345,
		Capabilities: []string{
			"urn:ietf:params:xml:ns:netconf:base:1.0",
		},
	}, nil
}

func (m *MockTransport) SendHello(h *HelloMessage) error {
	m.SendHelloCalled = true
	return m.SendHelloErr
}

func (m *MockTransport) Close() error {
	m.CloseCalled = true
	return m.CloseErr
}

func (m *MockTransport) Send(data []byte) error {
	m.SendCalled = true
	m.SendData = data
	return m.SendErr
}

func (m *MockTransport) Receive() ([]byte, error) {
	m.ReceiveCalled = true
	if m.ReceiveErr != nil {
		return nil, m.ReceiveErr
	}
	if m.ReceiveData != nil {
		return m.ReceiveData, nil
	}
	return []byte(`<rpc-reply xmlns="urn:ietf:params:xml:ns:netconf:base:1.0"><ok/></rpc-reply>`), nil
}

// TestSessionClose tests the Close method
func TestSessionClose(t *testing.T) {
	mockTransport := &MockTransport{}
	session := &Session{
		Transport: mockTransport,
		SessionID: 12345,
	}

	err := session.Close()
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	if !mockTransport.CloseCalled {
		t.Error("expected transport Close to be called")
	}
}

// TestSessionCloseWithError tests Close when transport returns error
func TestSessionCloseWithError(t *testing.T) {
	mockTransport := &MockTransport{
		CloseErr: io.ErrClosedPipe,
	}
	session := &Session{
		Transport: mockTransport,
	}

	err := session.Close()
	if err != io.ErrClosedPipe {
		t.Errorf("expected ErrClosedPipe, got %v", err)
	}
}

// TestNewSession tests creating a new session
func TestNewSession(t *testing.T) {
	mockTransport := &MockTransport{
		HelloResponse: &HelloMessage{
			SessionID: 54321,
			Capabilities: []string{
				"urn:ietf:params:xml:ns:netconf:base:1.0",
				"urn:ietf:params:xml:ns:netconf:capability:candidate:1.0",
			},
		},
	}

	session, err := NewSession(mockTransport)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if session == nil {
		t.Fatal("expected non-nil session")
	}

	if session.SessionID != 54321 {
		t.Errorf("session ID mismatch: expected 54321, got %d", session.SessionID)
	}

	if len(session.ServerCapabilities) != 2 {
		t.Errorf("expected 2 capabilities, got %d", len(session.ServerCapabilities))
	}

	if !mockTransport.ReceiveHelloCalled {
		t.Error("expected ReceiveHello to be called")
	}

	if !mockTransport.SendHelloCalled {
		t.Error("expected SendHello to be called")
	}
}

// TestNewSessionReceiveHelloError tests NewSession when ReceiveHello fails
func TestNewSessionReceiveHelloError(t *testing.T) {
	mockTransport := &MockTransport{
		HelloErr: io.EOF,
	}

	session, err := NewSession(mockTransport)
	if err == nil {
		t.Error("expected error")
	}

	if session != nil {
		t.Error("expected nil session")
	}

	if err != io.EOF {
		t.Errorf("error mismatch: expected io.EOF, got %v", err)
	}
}

// TestNewSessionSendHelloError tests NewSession when SendHello fails
func TestNewSessionSendHelloError(t *testing.T) {
	mockTransport := &MockTransport{
		HelloResponse: &HelloMessage{
			SessionID: 12345,
			Capabilities: []string{
				"urn:ietf:params:xml:ns:netconf:base:1.0",
			},
		},
		SendHelloErr: io.ErrClosedPipe,
	}

	session, err := NewSession(mockTransport)
	if err == nil {
		t.Error("expected error")
	}

	if session != nil {
		t.Error("expected nil session")
	}
}

// TestSessionExec tests the Exec method
func TestSessionExec(t *testing.T) {
	mockTransport := &MockTransport{
		ReceiveData: []byte(`<rpc-reply xmlns="urn:ietf:params:xml:ns:netconf:base:1.0"><ok/></rpc-reply>`),
	}

	session := &Session{
		Transport:      mockTransport,
		SessionID:      12345,
		ErrOnWarning:   false,
		ServerCapabilities: []string{
			"urn:ietf:params:xml:ns:netconf:base:1.0",
		},
	}

	methods := []RPCMethod{RawMethod("<get-config><source><running/></source></get-config>")}
	reply, err := session.Exec(methods...)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if reply == nil {
		t.Fatal("expected non-nil reply")
	}

	if !mockTransport.SendCalled {
		t.Error("expected Send to be called")
	}

	if !mockTransport.ReceiveCalled {
		t.Error("expected Receive to be called")
	}
}

// TestSessionExecMultipleMethods tests Exec with multiple methods
func TestSessionExecMultipleMethods(t *testing.T) {
	mockTransport := &MockTransport{
		ReceiveData: []byte(`<rpc-reply xmlns="urn:ietf:params:xml:ns:netconf:base:1.0"><ok/></rpc-reply>`),
	}

	session := &Session{
		Transport:      mockTransport,
		SessionID:      12345,
		ErrOnWarning:   false,
		ServerCapabilities: []string{
			"urn:ietf:params:xml:ns:netconf:base:1.0",
		},
	}

	methods := []RPCMethod{
		RawMethod("<operation1/>"),
		RawMethod("<operation2/>"),
		RawMethod("<operation3/>"),
	}

	reply, err := session.Exec(methods...)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if reply == nil {
		t.Fatal("expected non-nil reply")
	}
}

// TestSessionExecSendError tests Exec when Send fails
func TestSessionExecSendError(t *testing.T) {
	mockTransport := &MockTransport{
		SendErr: io.ErrClosedPipe,
	}

	session := &Session{
		Transport:      mockTransport,
		SessionID:      12345,
		ErrOnWarning:   false,
		ServerCapabilities: []string{
			"urn:ietf:params:xml:ns:netconf:base:1.0",
		},
	}

	methods := []RPCMethod{RawMethod("<get-config/>")}
	_, err := session.Exec(methods...)

	if err == nil {
		t.Error("expected error")
	}

	if err != io.ErrClosedPipe {
		t.Errorf("error mismatch: expected ErrClosedPipe, got %v", err)
	}
}

// TestSessionExecReceiveError tests Exec when Receive fails
func TestSessionExecReceiveError(t *testing.T) {
	mockTransport := &MockTransport{
		ReceiveErr: io.EOF,
	}

	session := &Session{
		Transport:      mockTransport,
		SessionID:      12345,
		ErrOnWarning:   false,
		ServerCapabilities: []string{
			"urn:ietf:params:xml:ns:netconf:base:1.0",
		},
	}

	methods := []RPCMethod{RawMethod("<get-config/>")}
	_, err := session.Exec(methods...)

	if err == nil {
		t.Error("expected error")
	}

	if err != io.EOF {
		t.Errorf("error mismatch: expected io.EOF, got %v", err)
	}
}

// TestSessionStructure tests the Session structure
func TestSessionStructure(t *testing.T) {
	mockTransport := &MockTransport{}
	session := &Session{
		Transport:          mockTransport,
		SessionID:          12345,
		ServerCapabilities: []string{"cap1", "cap2"},
		ErrOnWarning:       true,
	}

	if session.Transport != mockTransport {
		t.Error("transport mismatch")
	}

	if session.SessionID != 12345 {
		t.Errorf("session ID mismatch: expected 12345, got %d", session.SessionID)
	}

	if len(session.ServerCapabilities) != 2 {
		t.Errorf("expected 2 capabilities, got %d", len(session.ServerCapabilities))
	}

	if !session.ErrOnWarning {
		t.Error("ErrOnWarning should be true")
	}
}

// TestSessionWithDifferentCapabilities tests sessions with different server capabilities
func TestSessionWithDifferentCapabilities(t *testing.T) {
	testCases := []struct {
		name         string
		capabilities []string
	}{
		{
			name:         "base-1.0",
			capabilities: []string{"urn:ietf:params:xml:ns:netconf:base:1.0"},
		},
		{
			name: "base-1.0-with-candidate",
			capabilities: []string{
				"urn:ietf:params:xml:ns:netconf:base:1.0",
				"urn:ietf:params:xml:ns:netconf:capability:candidate:1.0",
			},
		},
		{
			name: "base-1.1-with-multiple",
			capabilities: []string{
				"urn:ietf:params:xml:ns:netconf:base:1.1",
				"urn:ietf:params:xml:ns:netconf:capability:candidate:1.0",
				"urn:ietf:params:xml:ns:netconf:capability:commit:1.0",
				"urn:ietf:params:xml:ns:netconf:capability:notification:1.0",
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mockTransport := &MockTransport{
				HelloResponse: &HelloMessage{
					SessionID:    12345,
					Capabilities: tc.capabilities,
				},
			}

			session, err := NewSession(mockTransport)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if len(session.ServerCapabilities) != len(tc.capabilities) {
				t.Errorf("capability count mismatch: expected %d, got %d",
					len(tc.capabilities), len(session.ServerCapabilities))
			}

			for i, cap := range tc.capabilities {
				if i >= len(session.ServerCapabilities) || session.ServerCapabilities[i] != cap {
					t.Errorf("capability %d mismatch: expected %s, got %s",
						i, cap, session.ServerCapabilities[i])
				}
			}
		})
	}
}

// TestSessionErrorOnWarning tests the ErrOnWarning flag
func TestSessionErrorOnWarning(t *testing.T) {
	mockTransport := &MockTransport{
		ReceiveData: []byte(`<rpc-reply xmlns="urn:ietf:params:xml:ns:netconf:base:1.0"><ok/></rpc-reply>`),
	}

	session := &Session{
		Transport:      mockTransport,
		SessionID:      12345,
		ErrOnWarning:   true, // Should treat warnings as errors
		ServerCapabilities: []string{
			"urn:ietf:params:xml:ns:netconf:base:1.0",
		},
	}

	methods := []RPCMethod{RawMethod("<test/>")}
	reply, err := session.Exec(methods...)

	if err != nil {
		t.Logf("error with ErrOnWarning=true: %v", err)
	}

	if reply == nil {
		t.Error("expected non-nil reply")
	}
}

// BenchmarkSessionExec benchmarks the Exec method
func BenchmarkSessionExec(b *testing.B) {
	mockTransport := &MockTransport{
		ReceiveData: []byte(`<rpc-reply xmlns="urn:ietf:params:xml:ns:netconf:base:1.0"><ok/></rpc-reply>`),
	}

	session := &Session{
		Transport:      mockTransport,
		SessionID:      12345,
		ErrOnWarning:   false,
		ServerCapabilities: []string{
			"urn:ietf:params:xml:ns:netconf:base:1.0",
		},
	}

	methods := []RPCMethod{RawMethod("<get-config><source><running/></source></get-config>")}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = session.Exec(methods...)
	}
}

// TestSessionMarshalXML tests XML marshaling of RPC messages
func TestSessionMarshalXML(t *testing.T) {
	mockTransport := &MockTransport{
		ReceiveData: []byte(`<rpc-reply xmlns="urn:ietf:params:xml:ns:netconf:base:1.0"><ok/></rpc-reply>`),
	}

	session := &Session{
		Transport:      mockTransport,
		SessionID:      12345,
		ErrOnWarning:   false,
		ServerCapabilities: []string{
			"urn:ietf:params:xml:ns:netconf:base:1.0",
		},
	}

	methods := []RPCMethod{RawMethod("<test-method/>")}
	reply, err := session.Exec(methods...)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if reply == nil {
		t.Fatal("expected non-nil reply")
	}

	// Verify the XML was sent
	if !mockTransport.SendCalled {
		t.Error("expected Send to be called")
	}

	// Check that XML header and rpc element were sent
	if len(mockTransport.SendData) == 0 {
		t.Error("expected data to be sent")
	}

	// Verify XML contains expected elements
	_ = string(mockTransport.SendData)
	if bytes.Index(mockTransport.SendData, []byte("<?xml")) == -1 {
		t.Error("expected XML header")
	}

	if bytes.Index(mockTransport.SendData, []byte("<rpc")) == -1 {
		t.Error("expected rpc element")
	}
}

// TestDefaultCapabilities tests DefaultCapabilities constant
func TestDefaultCapabilities(t *testing.T) {
	if len(DefaultCapabilities) == 0 {
		t.Error("DefaultCapabilities should not be empty")
	}

	// Check for base capability
	hasBase := false
	for _, cap := range DefaultCapabilities {
		if bytes.Contains([]byte(cap), []byte("base")) {
			hasBase = true
			break
		}
	}

	if !hasBase {
		t.Error("DefaultCapabilities should contain base capability")
	}
}

// TestSessionZeroValues tests zero values of Session
func TestSessionZeroValues(t *testing.T) {
	s := &Session{}

	if s.Transport != nil {
		t.Error("zero Transport should be nil")
	}

	if s.SessionID != 0 {
		t.Error("zero SessionID should be 0")
	}

	if s.ServerCapabilities != nil && len(s.ServerCapabilities) > 0 {
		t.Error("zero ServerCapabilities should be empty")
	}

	if s.ErrOnWarning {
		t.Error("zero ErrOnWarning should be false")
	}
}
