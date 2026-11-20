package netconf

import (
	"context"
	"errors"
	"testing"
)

// fakeSession answers each RPC with reply, or fails with err.
type fakeSession struct {
	reply  string
	err    error
	sent   []string
	closed bool
}

func (f *fakeSession) do(_ context.Context, op string) ([]byte, error) {
	f.sent = append(f.sent, op)
	if f.err != nil {
		return nil, f.err
	}
	return []byte(f.reply), nil
}

func (f *fakeSession) close() error {
	f.closed = true
	return nil
}

func sessionClient(sessions ...*fakeSession) (*GoNCClient, *int) {
	dials := 0
	return &GoNCClient{dial: func(context.Context) (rpcSession, error) {
		if dials >= len(sessions) {
			return nil, errors.New("no more sessions")
		}
		dials++
		return sessions[dials-1], nil
	}}, &dials
}

const okReply = `<rpc-reply xmlns="urn:ietf:params:xml:ns:netconf:base:1.0"><ok/></rpc-reply>`

// RPCs share one session.
func TestExecuteReusesSession(t *testing.T) {
	s := &fakeSession{reply: okReply}
	c, dials := sessionClient(s)
	for _, op := range []string{"<a/>", "<b/>", "<c/>"} {
		if _, err := c.execute(context.Background(), op); err != nil {
			t.Fatal(err)
		}
	}
	if *dials != 1 || len(s.sent) != 3 || s.closed {
		t.Fatalf("dials %d, sent %v, closed %v", *dials, s.sent, s.closed)
	}
}

// A transport error closes the session and is returned, without sending the
// RPC again; the next RPC opens a new session.
func TestExecuteTransportErrorReopens(t *testing.T) {
	broken := &fakeSession{err: errors.New("closed connection")}
	fresh := &fakeSession{reply: okReply}
	c, dials := sessionClient(broken, fresh)
	if _, err := c.execute(context.Background(), "<edit-config/>"); err == nil {
		t.Fatal("expected the transport error")
	}
	if len(broken.sent) != 1 || !broken.closed {
		t.Fatalf("broken session: sent %v, closed %v", broken.sent, broken.closed)
	}
	if _, err := c.execute(context.Background(), "<commit/>"); err != nil {
		t.Fatal(err)
	}
	if *dials != 2 || len(fresh.sent) != 1 {
		t.Fatalf("dials %d, fresh sent %v", *dials, fresh.sent)
	}
}

// An rpc-error is the device's answer, not a broken session: it is kept.
func TestExecuteRPCErrorKeepsSession(t *testing.T) {
	s := &fakeSession{reply: `<rpc-reply xmlns="urn:ietf:params:xml:ns:netconf:base:1.0"><rpc-error><error-severity>error</error-severity><error-message>syntax error</error-message></rpc-error></rpc-reply>`}
	c, dials := sessionClient(s)
	for i := 0; i < 2; i++ {
		if _, err := c.execute(context.Background(), "<bad/>"); err == nil {
			t.Fatal("expected the rpc-error")
		}
	}
	if *dials != 1 || s.closed {
		t.Fatalf("dials %d, closed %v", *dials, s.closed)
	}
}

func TestExecuteDialError(t *testing.T) {
	c, _ := sessionClient()
	if _, err := c.execute(context.Background(), "<a/>"); err == nil {
		t.Fatal("expected the dial error")
	}
	if c.session != nil {
		t.Fatal("session kept after a failed dial")
	}
}

func TestCloseClosesSession(t *testing.T) {
	s := &fakeSession{reply: okReply}
	c, _ := sessionClient(s)
	if _, err := c.execute(context.Background(), "<a/>"); err != nil {
		t.Fatal(err)
	}
	if err := c.Close(); err != nil || !s.closed || c.session != nil {
		t.Fatalf("close: err %v, closed %v", err, s.closed)
	}
	if err := c.Close(); err != nil {
		t.Fatalf("second close: %v", err)
	}
}
