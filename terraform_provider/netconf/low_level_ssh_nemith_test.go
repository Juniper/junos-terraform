package netconf

import (
	"bytes"
	"io"
	"testing"
)

type mockMsgTransport struct {
	reader      io.ReadCloser
	writer      *bytes.Buffer
	msgReaderN  int
	msgWriterN  int
	closeCalled bool
}

func (m *mockMsgTransport) MsgReader() (io.ReadCloser, error) {
	m.msgReaderN++
	return m.reader, nil
}

func (m *mockMsgTransport) MsgWriter() (io.WriteCloser, error) {
	m.msgWriterN++
	return nopWriteCloser{m.writer}, nil
}

func (m *mockMsgTransport) Close() error {
	m.closeCalled = true
	return nil
}

type nopWriteCloser struct {
	io.Writer
}

func (n nopWriteCloser) Close() error { return nil }

func TestTransportSSHSendUsesMessageTransport(t *testing.T) {
	buf := bytes.NewBuffer(nil)
	m := &mockMsgTransport{writer: buf}
	transport := &TransportSSH{netconfTransport: m}

	payload := []byte("<rpc><get-config/></rpc>")
	if err := transport.Send(payload); err != nil {
		t.Fatalf("Send() failed: %v", err)
	}

	if m.msgWriterN != 1 {
		t.Fatalf("expected MsgWriter to be called once, got %d", m.msgWriterN)
	}

	if got := buf.String(); got != string(payload) {
		t.Fatalf("unexpected written payload: got %q want %q", got, string(payload))
	}
}

func TestTransportSSHReceiveUsesMessageTransport(t *testing.T) {
	want := []byte("<rpc-reply><ok/></rpc-reply>")
	m := &mockMsgTransport{reader: io.NopCloser(bytes.NewReader(want))}
	transport := &TransportSSH{netconfTransport: m}

	got, err := transport.Receive()
	if err != nil {
		t.Fatalf("Receive() failed: %v", err)
	}

	if m.msgReaderN != 1 {
		t.Fatalf("expected MsgReader to be called once, got %d", m.msgReaderN)
	}

	if !bytes.Equal(got, want) {
		t.Fatalf("unexpected read payload: got %q want %q", string(got), string(want))
	}
}

func TestTransportSSHCloseUsesMessageTransport(t *testing.T) {
	m := &mockMsgTransport{writer: bytes.NewBuffer(nil), reader: io.NopCloser(bytes.NewReader(nil))}
	transport := &TransportSSH{netconfTransport: m}

	if err := transport.Close(); err != nil {
		t.Fatalf("Close() failed: %v", err)
	}

	if !m.closeCalled {
		t.Fatal("expected underlying message transport Close to be called")
	}
}
