package netconf

import (
	"context"
	"errors"
	"strings"
	"sync"
	"testing"
)

func newMockClient(calls *[]string, ret string, err error) *GoNCClient {
	return &GoNCClient{
		Lock: sync.RWMutex{},
		exec: func(_ context.Context, op string) (string, error) {
			*calls = append(*calls, op)
			return ret, err
		},
	}
}

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

func TestNewClientAllowsMissingSSHKeyPath(t *testing.T) {
	client, err := NewClient("user", "", "/does/not/exist", "127.0.0.1", 830)
	if err != nil {
		t.Fatalf("expected no error for invalid ssh key path, got: %v", err)
	}
	if client == nil {
		t.Fatal("expected non-nil client")
	}
}
