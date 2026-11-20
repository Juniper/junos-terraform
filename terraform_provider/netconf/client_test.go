package netconf

import (
	"context"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"testing"
	"time"

	"golang.org/x/crypto/ssh"
)

// newMockClient creates a client with an injected RPC executor for deterministic tests.
func newMockClient(calls *[]string, ret string, err error) *GoNCClient {
	return &GoNCClient{
		Lock: sync.RWMutex{},
		exec: func(_ context.Context, op string) (string, error) {
			*calls = append(*calls, op)
			return ret, err
		},
	}
}

// TestSendUpdateBaseConfigPayload verifies patch updates do not require group name tags.
func TestSendUpdateBaseConfigPayload(t *testing.T) {
	calls := []string{}
	client := newMockClient(&calls, "<ok/>", nil)

	diff := `<configuration><system><host-name nc:operation="replace">leaf1</host-name></system></configuration>`
	if err := client.SendUpdate("base-config", diff, false); err != nil {
		t.Fatalf("SendUpdate returned error: %v", err)
	}

	if len(calls) != 1 {
		t.Fatalf("expected single edit-config operation, got %d", len(calls))
	}
	if !strings.Contains(calls[0], "<edit-config>") || !strings.Contains(calls[0], "<default-operation>merge</default-operation>") {
		t.Fatalf("expected patch edit-config envelope, got %q", calls[0])
	}
}

// TestMarshalRPCRequestUsesInnerXML verifies patch requests are wrapped in
// <rpc> while the operation body remains raw XML, not wrapper fields.
func TestMarshalRPCRequestUsesInnerXML(t *testing.T) {
	rpcXML, err := marshalRPCRequest("<edit-config><target><candidate/></target></edit-config>", "1")
	if err != nil {
		t.Fatalf("marshalRPCRequest() returned error: %v", err)
	}

	if !strings.Contains(rpcXML, `<rpc xmlns="urn:ietf:params:xml:ns:netconf:base:1.0" message-id="1">`) {
		t.Fatalf("expected rpc envelope, got %q", rpcXML)
	}
	if !strings.Contains(rpcXML, "<edit-config>") {
		t.Fatalf("expected edit-config payload in rpc, got %q", rpcXML)
	}
	if strings.Contains(rpcXML, "<Operation>") || strings.Contains(rpcXML, "<RawXML>") {
		t.Fatalf("expected raw payload without wrapper element, got %q", rpcXML)
	}
}

// TestSendCommitDiscardsOnCommitError verifies discard-changes is sent after commit failure.
func TestSendCommitDiscardsOnCommitError(t *testing.T) {
	calls := []string{}
	client := &GoNCClient{
		Lock: sync.RWMutex{},
		exec: func(_ context.Context, op string) (string, error) {
			calls = append(calls, op)
			if strings.TrimSpace(op) == commitStr {
				return "", errors.New("commit failed")
			}
			return "<ok/>", nil
		},
	}

	err := client.SendCommit()
	if err == nil {
		t.Fatal("expected commit error")
	}

	// Only the commit and the discard that follows it: apply-groups is
	// ordinary configuration and is not emitted on the resource's behalf.
	if len(calls) != 2 {
		t.Fatalf("expected commit then discard, got %#v", calls)
	}
	if strings.TrimSpace(calls[0]) != commitStr {
		t.Fatalf("expected commit first, got %q", calls[0])
	}
	if strings.TrimSpace(calls[len(calls)-1]) != discardChanges {
		t.Fatalf("expected discard-changes after commit failure, got %q", calls[len(calls)-1])
	}
}

// TestNewClientAllowsMissingSSHKeyPath verifies client creation tolerates unreadable key paths.
func TestNewClientAllowsMissingSSHKeyPath(t *testing.T) {
	client, err := NewClient("user", "", "/does/not/exist", "127.0.0.1", 830)
	if err != nil {
		t.Fatalf("expected no error for invalid ssh key path, got: %v", err)
	}
	if client == nil {
		t.Fatal("expected non-nil client")
	}
}

// TestGoNCClientCloseNoop verifies Close preserves no-op behavior.
func TestGoNCClientCloseNoop(t *testing.T) {
	client := &GoNCClient{}
	if err := client.Close(); err != nil {
		t.Fatalf("expected nil close error, got: %v", err)
	}
}

// TestReadRawConfigAsksForCommittedConfiguration checks that the read does not
// ask the device to resolve apply-groups. With inheritance resolved, a device
// reports a group's values in the base hierarchy too, and configuration that
// was never declared there would be read back as drift on every plan.
func TestReadRawConfigAsksForCommittedConfiguration(t *testing.T) {
	calls := []string{}
	client := newMockClient(&calls, "<configuration/>", nil)

	if _, err := client.readRawConfig(); err != nil {
		t.Fatalf("readRawConfig() returned error: %v", err)
	}
	if len(calls) != 1 {
		t.Fatalf("expected one call, got %#v", calls)
	}
	if !strings.Contains(calls[0], "<get-configuration>") {
		t.Fatalf("expected get-configuration, got %q", calls[0])
	}
	if strings.Contains(calls[0], "inherit") {
		t.Fatalf("read must not resolve inheritance: %q", calls[0])
	}
}

// TestPublicKeyFileErrorPaths verifies key loader failure paths.
func TestPublicKeyFileErrorPaths(t *testing.T) {
	if method := publicKeyFile("/does/not/exist"); method != nil {
		t.Fatalf("expected nil auth method for missing file")
	}

	dir := t.TempDir()
	keyPath := filepath.Join(dir, "invalid.key")
	if err := os.WriteFile(keyPath, []byte("not-a-key"), 0600); err != nil {
		t.Fatalf("failed to write key file: %v", err)
	}
	if method := publicKeyFile(keyPath); method != nil {
		t.Fatalf("expected nil auth method for invalid key")
	}
}

// TestNewClientDefaultPort verifies zero port maps to NETCONF default port.
func TestNewClientDefaultPort(t *testing.T) {
	client, err := NewClient("user", "pass", "", "127.0.0.1", 0)
	if err != nil {
		t.Fatalf("NewClient() error: %v", err)
	}
	gonc, ok := client.(*GoNCClient)
	if !ok {
		t.Fatalf("expected *GoNCClient")
	}
	if gonc.port != defaultPort {
		t.Fatalf("expected default port %d, got %d", defaultPort, gonc.port)
	}
}

// TestExecuteWithoutMockReturnsDialError verifies network execution returns dial errors.
func TestExecuteWithoutMockReturnsDialError(t *testing.T) {
	client := &GoNCClient{
		host:      "invalid-hostname-for-test",
		port:      830,
		sshConfig: &ssh.ClientConfig{HostKeyCallback: ssh.InsecureIgnoreHostKey()},
	}

	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel()

	_, err := client.execute(ctx, "<rpc/>")
	if err == nil {
		t.Fatalf("expected network execute to fail")
	}
}

// TestSendCommitOnlyCommits checks that committing does not write configuration
// of its own. apply-groups is part of the payload the resource sends, so a group
// reference the configuration does not declare must never reach the device.
func TestSendCommitOnlyCommits(t *testing.T) {
	calls := []string{}
	client := newMockClient(&calls, "<ok/>", nil)

	if err := client.SendCommit(); err != nil {
		t.Fatalf("SendCommit() error: %v", err)
	}
	if len(calls) != 1 {
		t.Fatalf("expected a single commit call, got %#v", calls)
	}
	if strings.TrimSpace(calls[0]) != commitStr {
		t.Fatalf("expected commit, got %q", calls[0])
	}
	if strings.Contains(calls[0], "apply-groups") {
		t.Fatalf("commit must not emit apply-groups: %q", calls[0])
	}
}

// TestNewClientWithValidSSHKey verifies SSH key auth path setup with a valid key.
func TestNewClientWithValidSSHKey(t *testing.T) {
	key, err := rsa.GenerateKey(rand.Reader, 1024)
	if err != nil {
		t.Fatalf("failed generating rsa key: %v", err)
	}
	keyDER := x509.MarshalPKCS1PrivateKey(key)
	keyPEM := pem.EncodeToMemory(&pem.Block{Type: "RSA PRIVATE KEY", Bytes: keyDER})

	keyPath := filepath.Join(t.TempDir(), "id_rsa")
	if err := os.WriteFile(keyPath, keyPEM, 0600); err != nil {
		t.Fatalf("failed writing key: %v", err)
	}

	client, err := NewClient("user", "", keyPath, "127.0.0.1", 830)
	if err != nil {
		t.Fatalf("NewClient() error: %v", err)
	}
	gonc := client.(*GoNCClient)
	if len(gonc.sshConfig.Auth) != 1 {
		t.Fatalf("expected one auth method")
	}
}
