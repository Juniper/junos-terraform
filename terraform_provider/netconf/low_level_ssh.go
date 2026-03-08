// Go NETCONF Client
//
// Copyright (c) 2013-2018, Juniper Networks, Inc. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package netconf

import (
	"context"
	"crypto/x509"
	"encoding/pem"
	"encoding/xml"
	"errors"
	"fmt"
	"io"
	"net"
	"os"
	"strings"
	"time"

	"golang.org/x/crypto/ssh"
	"golang.org/x/crypto/ssh/agent"
	netconfssh "nemith.io/netconf/transport/ssh"
)

const (
	// DefaultPort is the default SSH port used when communicating with
	// NETCONF
	LowLevelDefaultPort = 830
	// sshNetconfSubsystem sets the SSH subsystem to NETCONF
	sshNetconfSubsystem = "netconf"
)

// TransportSSH maintains the information necessary to communicate with the
// remote device over SSH
type TransportSSH struct {
	TransportBasicIO              // Embedded Transport basic IO base type
	SSHClient        *ssh.Client  // SSH Client
	SSHSession       *ssh.Session // SSH Client Session
	netconfTransport msgTransport
}

type msgTransport interface {
	MsgReader() (io.ReadCloser, error)
	MsgWriter() (io.WriteCloser, error)
	Close() error
}

// Close closes an existing SSH session and socket if they exist.
func (t *TransportSSH) Close() error {
	if t.netconfTransport != nil {
		if err := t.netconfTransport.Close(); !isBenignCloseError(err) {
			return err
		}
		return nil
	}

	if t.SSHSession != nil {
		if err := t.SSHSession.Close(); !isBenignCloseError(err) {
			return err
		}
	}

	err := error(nil)
	if t.SSHClient != nil {
		err = t.SSHClient.Close()
	}

	if !isBenignCloseError(err) {
		return err
	}

	return nil
}

// Send writes one NETCONF message using the underlying message-oriented transport.
func (t *TransportSSH) Send(data []byte) error {
	if t.netconfTransport == nil {
		return t.TransportBasicIO.Send(data)
	}

	w, err := t.netconfTransport.MsgWriter()
	if err != nil {
		return err
	}

	if _, err := w.Write(data); err != nil {
		_ = w.Close()
		return err
	}

	return w.Close()
}

// Receive reads one NETCONF message using the underlying message-oriented transport.
func (t *TransportSSH) Receive() ([]byte, error) {
	if t.netconfTransport == nil {
		return t.TransportBasicIO.Receive()
	}

	r, err := t.netconfTransport.MsgReader()
	if err != nil {
		return nil, err
	}
	defer func() {
		_ = r.Close()
	}()

	out, err := io.ReadAll(r)
	if err != nil {
		return nil, err
	}

	return out, nil
}

// SendHello sends a NETCONF hello message using this transport.
func (t *TransportSSH) SendHello(hello *HelloMessage) error {
	val, err := xml.Marshal(hello)
	if err != nil {
		return err
	}

	header := []byte(xml.Header)
	val = append(header, val...)
	return t.Send(val)
}

// ReceiveHello receives and decodes the remote NETCONF hello message.
func (t *TransportSSH) ReceiveHello() (*HelloMessage, error) {
	hello := new(HelloMessage)

	val, err := t.Receive()
	if err != nil {
		return hello, err
	}

	err = xml.Unmarshal(val, hello)
	return hello, err
}

// isBenignCloseError filters expected errors when the remote side closes first.
func isBenignCloseError(err error) bool {
	if err == nil {
		return true
	}

	if errors.Is(err, io.EOF) || strings.EqualFold(err.Error(), "EOF") {
		return true
	}

	errText := strings.ToLower(err.Error())
	return strings.Contains(errText, "closed pipe") ||
		strings.Contains(errText, "use of closed network connection")
}

// DialSSH connects and establishes SSH sessions
//
// target can be an IP address (e.g.) 172.16.1.1 which utlizes the default
// NETCONF over SSH port of 830.  Target can also specify a port with the
// following format <host>:<port (e.g 172.16.1.1:22)
//
// config takes a ssh.ClientConfig connection. See documentation for
// go.crypto/ssh for documentation.  There is a helper function SSHConfigPassword
// thar returns a ssh.ClientConfig for simple username/password authentication
func (t *TransportSSH) DialSSH(target string, config *ssh.ClientConfig, port int) error {
	if !strings.Contains(target, ":") {
		sshport := 0
		if port != 0 {
			sshport = port
		} else {
			sshport = LowLevelDefaultPort
		}
		target = fmt.Sprintf("%s:%d", target, sshport)
	}

	ctx := context.Background()
	if config != nil && config.Timeout > 0 {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, config.Timeout)
		defer cancel()
	}

	transport, err := netconfssh.Dial(ctx, "tcp", target, config)
	if err != nil {
		return err
	}

	t.netconfTransport = transport
	return nil
}

// SetupSession sorts out wiring
func (t *TransportSSH) SetupSession() error {
	if t.netconfTransport != nil {
		return nil
	}

	var err error
	t.SSHSession, err = t.SSHClient.NewSession()
	if err != nil {
		return err
	}

	writer, err := t.SSHSession.StdinPipe()
	if err != nil {
		return err
	}

	reader, err := t.SSHSession.StdoutPipe()
	if err != nil {
		return err
	}

	t.ReadWriteCloser = NewReadWriteCloser(reader, writer)
	return t.SSHSession.RequestSubsystem(sshNetconfSubsystem)
}

// NewSSHSession creates a new NETCONF session using an existing net.Conn.
func NewSSHSession(conn net.Conn, config *ssh.ClientConfig) (*Session, error) {
	t, err := connToTransport(conn, config)
	if err != nil {
		return nil, err
	}

	s, err := NewSession(t)
	if err != nil {
		return nil, err
	}

	return s, nil
}

// Dial creates a new NETCONF session using a SSH
// See TransportSSH.Dial for arguments.
func Dial(target string, config *ssh.ClientConfig, port int) (*Session, error) {
	t := TransportSSH{}
	err := t.DialSSH(target, config, port)

	if err != nil {
		return nil, err
	}

	s, err := NewSession(&t)

	if err != nil {
		return nil, err
	}

	return s, nil
}

// DialSSHTimeout creates a new NETCONF session using a SSH Transport with timeout.
// See TransportSSH.Dial for arguments.
// The timeout value is used for both connection establishment and Read/Write operations.
func DialSSHTimeout(target string, config *ssh.ClientConfig, timeout time.Duration) (*Session, error) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	transport, err := netconfssh.Dial(ctx, "tcp", target, config)
	if err != nil {
		return nil, err
	}

	t := &TransportSSH{netconfTransport: transport}
	s, err := NewSession(t)
	if err != nil {
		return nil, err
	}

	return s, nil
}

// SSHConfigPassword is a convenience function that takes a username and password
// and returns a new ssh.ClientConfig setup to pass that username and password.
// Convenience means that HostKey checks are disabled so it's probably less secure
func SSHConfigPassword(user string, pass string) *ssh.ClientConfig {
	return &ssh.ClientConfig{
		User: user,
		Auth: []ssh.AuthMethod{
			ssh.Password(pass),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}
}

// SSHConfigPubKeyFile is a convenience function that takes a username, private key
// and passphrase and returns a new ssh.ClientConfig setup to pass credentials
// to DialSSH
func SSHConfigPubKeyFile(user string, file string, passphrase string) (*ssh.ClientConfig, error) {
	buf, err := os.ReadFile(file)
	if err != nil {
		return nil, err
	}
	block, rest := pem.Decode(buf)
	if len(rest) > 0 {
		return nil, fmt.Errorf("pem: unable to decode file %s", file)
	}

	//nolint:staticcheck // deprecated encryption support preserved for compatibility
	if x509.IsEncryptedPEMBlock(block) {
		b, err := x509.DecryptPEMBlock(block, []byte(passphrase)) //nolint:staticcheck
		if err != nil {
			return nil, err
		}
		buf = pem.EncodeToMemory(&pem.Block{
			Type:  block.Type,
			Bytes: b,
		})
	}

	key, err := ssh.ParsePrivateKey(buf)
	if err != nil {
		return nil, err
	}
	return &ssh.ClientConfig{
		User: user,
		Auth: []ssh.AuthMethod{
			ssh.PublicKeys(key),
		},
	}, nil

}

// SSHConfigPubKeyAgent is a convience function that takes a username and
// returns a new ssh.Clientconfig setup to pass credentials received from
// an ssh agent
func SSHConfigPubKeyAgent(user string) (*ssh.ClientConfig, error) {
	c, err := net.Dial("unix", os.Getenv("SSH_AUTH_SOCK"))
	if err != nil {
		return nil, err
	}
	return &ssh.ClientConfig{
		User: user,
		Auth: []ssh.AuthMethod{
			ssh.PublicKeysCallback(agent.NewClient(c).Signers),
		},
	}, nil
}

func connToTransport(conn net.Conn, config *ssh.ClientConfig) (*TransportSSH, error) {
	c, chans, reqs, err := ssh.NewClientConn(conn, conn.RemoteAddr().String(), config)
	if err != nil {
		return nil, err
	}

	t := &TransportSSH{}
	t.SSHClient = ssh.NewClient(c, chans, reqs)

	err = t.SetupSession()
	if err != nil {
		return nil, err
	}

	return t, nil
}
