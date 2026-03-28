package netconf

import (
	"context"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"encoding/xml"
	"errors"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"testing"
	"time"

	"golang.org/x/crypto/ssh"
	netconf "nemith.io/netconf"
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

// TestDeleteConfigCallsExpectedOperations verifies delete then commit RPC sequencing.
func TestDeleteConfigCallsExpectedOperations(t *testing.T) {
	calls := []string{}
	client := newMockClient(&calls, "<ok/>", nil)

	_, err := client.DeleteConfig("base-config", true)
	if err != nil {
		t.Fatalf("DeleteConfig returned error: %v", err)
	}

	if len(calls) != 2 {
		t.Fatalf("expected 2 operations, got %d", len(calls))
	}
	if !strings.Contains(calls[0], "<edit-config>") {
		t.Fatalf("first call should be edit-config, got %q", calls[0])
	}
	if strings.TrimSpace(calls[1]) != commitStr {
		t.Fatalf("second call should be commit, got %q", calls[1])
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
	if !strings.Contains(calls[0], "<edit-config>") || !strings.Contains(calls[0], "<default-operation>none</default-operation>") {
		t.Fatalf("expected patch edit-config envelope, got %q", calls[0])
	}
}

// TestNetconfRPCSerializationUsesInnerXML verifies raw operations are embedded
// directly in <rpc> without wrapper elements like <Operation> or <RawXML>.
func TestNetconfRPCSerializationUsesInnerXML(t *testing.T) {
	rpcBytes, err := xml.Marshal(&netconf.RPC{
		MessageID: "1",
		Operation: []byte("<edit-config><target><candidate/></target></edit-config>"),
	})
	if err != nil {
		t.Fatalf("xml.Marshal() returned error: %v", err)
	}

	rpcXML := string(rpcBytes)
	if !strings.Contains(rpcXML, "<edit-config>") {
		t.Fatalf("expected edit-config payload in rpc, got %q", rpcXML)
	}
	if strings.Contains(rpcXML, "<Operation>") || strings.Contains(rpcXML, "<RawXML>") {
		t.Fatalf("expected raw payload without wrapper element, got %q", rpcXML)
	}
}

// TestSendTransactionWithIDReplacesGroup verifies ID-based transactions use update flow.
func TestSendTransactionWithIDReplacesGroup(t *testing.T) {
	calls := []string{}
	client := newMockClient(&calls, "<ok/>", nil)

	obj := struct {
		XMLName struct{} `xml:"configuration"`
		Groups  struct {
			Name string `xml:"name"`
		} `xml:"groups"`
	}{}
	obj.Groups.Name = "base-config"

	if err := client.SendTransaction("base-config", obj, false); err != nil {
		t.Fatalf("SendTransaction returned error: %v", err)
	}

	if len(calls) != 2 {
		t.Fatalf("expected 2 operations for update flow, got %d", len(calls))
	}
	if !strings.Contains(calls[0], "operation=\"delete\"") {
		t.Fatalf("expected delete operation first, got %q", calls[0])
	}
	if !strings.Contains(calls[1], "<load-configuration") {
		t.Fatalf("expected load-configuration second, got %q", calls[1])
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

	applyGroupsList = []string{"b", "a"}
	err := client.SendCommit()
	if err == nil {
		t.Fatal("expected commit error")
	}

	if len(calls) < 3 {
		t.Fatalf("expected apply-groups, commit, discard operations, got %d", len(calls))
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

// TestUpdateRawConfigMissingName verifies group name extraction failures are surfaced.
func TestUpdateRawConfigMissingName(t *testing.T) {
	client := newMockClient(&[]string{}, "<ok/>", nil)
	_, err := client.updateRawConfig("group", "<configuration></configuration>", false)
	if err == nil || !strings.Contains(err.Error(), "failed to extract") {
		t.Fatalf("expected extract error, got: %v", err)
	}
}

// TestSendRawConfigCommitFlow verifies raw config load followed by commit.
func TestSendRawConfigCommitFlow(t *testing.T) {
	calls := []string{}
	client := newMockClient(&calls, "<ok/>", nil)
	applyGroupsList = nil

	reply, err := client.sendRawConfig("<configuration><groups><name>z-group</name></groups></configuration>", true)
	if err != nil {
		t.Fatalf("sendRawConfig() returned error: %v", err)
	}
	if reply != "<ok/>" {
		t.Fatalf("unexpected reply: %q", reply)
	}
	if len(calls) != 2 {
		t.Fatalf("expected load and commit calls, got: %d", len(calls))
	}
	if strings.TrimSpace(calls[1]) != commitStr {
		t.Fatalf("expected commit as second call, got %q", calls[1])
	}
}

// TestSendRawConfigMissingName verifies missing group names return an error.
func TestSendRawConfigMissingName(t *testing.T) {
	client := newMockClient(&[]string{}, "<ok/>", nil)
	_, err := client.sendRawConfig("<configuration></configuration>", false)
	if err == nil || !strings.Contains(err.Error(), "failed to extract") {
		t.Fatalf("expected extract error, got: %v", err)
	}
}

// TestReadRawGroupUsesGetConfigRPC verifies group reads use get-configuration RPC.
func TestReadRawGroupUsesGetConfigRPC(t *testing.T) {
	calls := []string{}
	client := newMockClient(&calls, "<group/>", nil)

	reply, err := client.readRawGroup("base-config")
	if err != nil {
		t.Fatalf("readRawGroup() returned error: %v", err)
	}
	if reply != "<group/>" {
		t.Fatalf("unexpected reply: %q", reply)
	}
	if len(calls) != 1 || !strings.Contains(calls[0], "<get-configuration>") {
		t.Fatalf("expected get-configuration call, got %#v", calls)
	}
}

// TestMarshalGroupSuccessAndError verifies XML unmarshalling success and failure paths.
func TestMarshalGroupSuccessAndError(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		calls := []string{}
		reply := `<configuration><groups><name>g1</name></groups></configuration>`
		client := newMockClient(&calls, reply, nil)

		var out struct {
			XMLName xml.Name `xml:"configuration"`
			Groups  struct {
				Name string `xml:"name"`
			} `xml:"groups"`
		}

		if err := client.MarshalGroup("g1", &out); err != nil {
			t.Fatalf("MarshalGroup() returned error: %v", err)
		}
		if out.Groups.Name != "g1" {
			t.Fatalf("unexpected group name: %q", out.Groups.Name)
		}
	})

	t.Run("invalid xml", func(t *testing.T) {
		calls := []string{}
		client := newMockClient(&calls, `<bad`, nil)
		var out struct{}
		if err := client.MarshalGroup("g1", &out); err == nil {
			t.Fatalf("expected unmarshal error")
		}
	})
}

// TestSendTransactionEmptyIDPathAndMarshalError verifies empty-ID path and marshal failures.
func TestSendTransactionEmptyIDPathAndMarshalError(t *testing.T) {
	t.Run("empty id uses raw path", func(t *testing.T) {
		calls := []string{}
		client := newMockClient(&calls, "<ok/>", nil)
		applyGroupsList = nil

		obj := struct {
			XMLName xml.Name `xml:"configuration"`
			Groups  struct {
				Name string `xml:"name"`
			} `xml:"groups"`
		}{}
		obj.Groups.Name = "group-a"

		if err := client.SendTransaction("", obj, false); err != nil {
			t.Fatalf("SendTransaction() error: %v", err)
		}
		if len(calls) != 1 || !strings.Contains(calls[0], "<load-configuration") {
			t.Fatalf("expected single load-configuration call, got %#v", calls)
		}
	})

	t.Run("marshal error", func(t *testing.T) {
		client := newMockClient(&[]string{}, "", nil)
		obj := map[string]interface{}{"bad": make(chan int)}
		if err := client.SendTransaction("", obj, false); err == nil {
			t.Fatalf("expected marshal error")
		}
	})
}

// TestApplyGroupsHelpersSortAndFilter verifies helper list sanitization and sorting.
func TestApplyGroupsHelpersSortAndFilter(t *testing.T) {
	applyGroupsList = nil
	addToApplyGroupsList("b")
	addToApplyGroupsList("")
	addToApplyGroupsList("a")
	sortApplyGroupsList()

	if len(applyGroupsList) != 2 {
		t.Fatalf("expected empty value to be filtered out, got %#v", applyGroupsList)
	}
	if applyGroupsList[0] != "a" || applyGroupsList[1] != "b" {
		t.Fatalf("unexpected sorted list: %#v", applyGroupsList)
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

// TestUpdateRawConfigCommitAndErrorBranches verifies update commit and error branches.
func TestUpdateRawConfigCommitAndErrorBranches(t *testing.T) {
	t.Run("commit true success", func(t *testing.T) {
		calls := []string{}
		client := newMockClient(&calls, "<ok/>", nil)
		applyGroupsList = nil

		reply, err := client.updateRawConfig("grp", "<configuration><groups><name>grp</name></groups></configuration>", true)
		if err != nil {
			t.Fatalf("updateRawConfig() error: %v", err)
		}
		if reply != "<ok/>" {
			t.Fatalf("unexpected reply: %q", reply)
		}
		if len(calls) != 3 {
			t.Fatalf("expected 3 calls (delete, load, commit), got %d", len(calls))
		}
	})

	t.Run("delete error", func(t *testing.T) {
		client := &GoNCClient{
			Lock: sync.RWMutex{},
			exec: func(_ context.Context, op string) (string, error) {
				if strings.Contains(op, "operation=\"delete\"") {
					return "", errors.New("delete failed")
				}
				return "<ok/>", nil
			},
		}
		_, err := client.updateRawConfig("grp", "<configuration><groups><name>grp</name></groups></configuration>", false)
		if err == nil {
			t.Fatalf("expected delete error")
		}
	})

	t.Run("commit error", func(t *testing.T) {
		client := &GoNCClient{
			Lock: sync.RWMutex{},
			exec: func(_ context.Context, op string) (string, error) {
				if strings.TrimSpace(op) == commitStr {
					return "", errors.New("commit failed")
				}
				return "<ok/>", nil
			},
		}
		_, err := client.updateRawConfig("grp", "<configuration><groups><name>grp</name></groups></configuration>", true)
		if err == nil {
			t.Fatalf("expected commit error")
		}
	})
}

// TestDeleteConfigBranches verifies non-commit and commit-error delete behavior.
func TestDeleteConfigBranches(t *testing.T) {
	t.Run("without commit", func(t *testing.T) {
		calls := []string{}
		client := newMockClient(&calls, "<ok/>\n", nil)
		reply, err := client.DeleteConfig("grp", false)
		if err != nil {
			t.Fatalf("DeleteConfig() error: %v", err)
		}
		if reply != "<ok/>" {
			t.Fatalf("expected newline-stripped reply, got %q", reply)
		}
		if len(calls) != 1 {
			t.Fatalf("expected single delete call, got %d", len(calls))
		}
	})

	t.Run("commit error", func(t *testing.T) {
		client := &GoNCClient{
			Lock: sync.RWMutex{},
			exec: func(_ context.Context, op string) (string, error) {
				if strings.TrimSpace(op) == commitStr {
					return "", errors.New("commit failed")
				}
				return "<ok/>", nil
			},
		}
		_, err := client.DeleteConfig("grp", true)
		if err == nil {
			t.Fatalf("expected commit error")
		}
	})
}

// TestSendCommitSuccessAndApplyGroupError verifies commit success and apply-group RPC failure behavior.
func TestSendCommitSuccessAndApplyGroupError(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		calls := []string{}
		client := newMockClient(&calls, "<ok/>", nil)
		applyGroupsList = []string{"z", "", "a"}

		if err := client.SendCommit(); err != nil {
			t.Fatalf("SendCommit() error: %v", err)
		}
		if len(calls) < 2 {
			t.Fatalf("expected apply-group load then commit calls, got %d", len(calls))
		}
		if !strings.Contains(calls[0], "<apply-groups>a</apply-groups>") || !strings.Contains(calls[0], "<apply-groups>z</apply-groups>") {
			t.Fatalf("expected sorted apply groups in RPC, got %q", calls[0])
		}
	})

	t.Run("apply-group load error", func(t *testing.T) {
		client := &GoNCClient{
			Lock: sync.RWMutex{},
			exec: func(_ context.Context, op string) (string, error) {
				if strings.Contains(op, "<load-configuration") {
					return "", errors.New("load failed")
				}
				return "<ok/>", nil
			},
		}
		applyGroupsList = []string{"x"}
		if err := client.SendCommit(); err == nil {
			t.Fatalf("expected sendApplyGroupsLocked error")
		}
	})
}

// TestMarshalGroupReadError verifies read errors are propagated by MarshalGroup.
func TestMarshalGroupReadError(t *testing.T) {
	client := &GoNCClient{
		Lock: sync.RWMutex{},
		exec: func(_ context.Context, _ string) (string, error) {
			return "", errors.New("read failed")
		},
	}
	var out struct{}
	if err := client.MarshalGroup("x", &out); err == nil {
		t.Fatalf("expected read error")
	}
}

// TestSendRawConfigExecuteError verifies RPC execution errors are returned.
func TestSendRawConfigExecuteError(t *testing.T) {
	client := &GoNCClient{
		Lock: sync.RWMutex{},
		exec: func(_ context.Context, _ string) (string, error) {
			return "", errors.New("rpc failed")
		},
	}
	_, err := client.sendRawConfig("<configuration><groups><name>g</name></groups></configuration>", false)
	if err == nil {
		t.Fatalf("expected rpc error")
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
