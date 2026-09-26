package netconf

import (
	"errors"
	"strings"
	"testing"

	netconf "nemith.io/netconf"
)

const replyNS = `xmlns:nc="urn:ietf:params:xml:ns:netconf:base:1.0"`

func TestParseReplyData(t *testing.T) {
	data, err := parseReply([]byte(`<nc:rpc-reply ` + replyNS + `><configuration><system/></configuration></nc:rpc-reply>`))
	if err != nil || !strings.Contains(data, "<configuration><system/></configuration>") {
		t.Fatalf("data %q, err %v", data, err)
	}
}

// As a Junos device answers an RPC it rejects.
func TestParseReplyError(t *testing.T) {
	_, err := parseReply([]byte(`<nc:rpc-reply ` + replyNS + `>
<nc:rpc-error>
<nc:error-type>protocol</nc:error-type>
<nc:error-tag>operation-failed</nc:error-tag>
<nc:error-severity>error</nc:error-severity>
<nc:error-message>syntax error</nc:error-message>
<nc:error-info><nc:bad-element>get-no-such-information</nc:bad-element></nc:error-info>
</nc:rpc-error>
</nc:rpc-reply>`))
	var errs netconf.RPCErrors
	if !errors.As(err, &errs) || len(errs) != 1 || errs[0].Message != "syntax error" {
		t.Fatalf("err %v", err)
	}
}

// A warning is not an error: the reply's content is returned.
func TestParseReplyWarning(t *testing.T) {
	data, err := parseReply([]byte(`<nc:rpc-reply ` + replyNS + `>
<nc:rpc-error><nc:error-severity>warning</nc:error-severity><nc:error-message>statement has no contents; ignored</nc:error-message></nc:rpc-error>
<nc:ok/>
</nc:rpc-reply>`))
	if err != nil || !strings.Contains(data, "ok") {
		t.Fatalf("data %q, err %v", data, err)
	}
}

// A missing statement on delete is reported so updateRawConfig and
// DeleteConfig can ignore it.
func TestParseReplyMissingDelete(t *testing.T) {
	_, err := parseReply([]byte(`<nc:rpc-reply ` + replyNS + `>
<nc:rpc-error><nc:error-type>application</nc:error-type><nc:error-tag>data-missing</nc:error-tag><nc:error-severity>error</nc:error-severity><nc:error-message>statement not found</nc:error-message></nc:rpc-error>
</nc:rpc-reply>`))
	if !isMissingDeleteError(err) {
		t.Fatalf("err %v is not a missing-delete error", err)
	}
}
