package netconf

import (
	"fmt"
	"net"
	"time"

	"golang.org/x/crypto/ssh"
)

// DriverSSH type is for creating an SSH based driver. Maintains state for session and connection. Implements Driver{}
type DriverSSH struct {
	Timeout   time.Duration     // Timeout for SSH timed sessions
	Port      int               // Target port
	Host      string            // Target hostname
	Target    string            // Target hostname:port
	Datastore string            // NETCONF datastore
	Conn      net.Conn          // Conn for session
	SSHConfig *ssh.ClientConfig // SSH Config
	Transport *TransportSSH     // Transport data
	Session   *Session          // Session data
}

// New creates a new instance of DriverSSH
func NewSSH() *DriverSSH {
	var t TransportSSH
	return &DriverSSH{
		Port:      LowLevelDefaultPort,
		Transport: &t,
	}
}

// SetDatastore sets the target datastore on the data structure
func (d *DriverSSH) SetDatastore(ds string) error {
	d.Datastore = ds
	return nil
}

// Dial function (call this after New())
func (d *DriverSSH) Dial() error {
	d.Target = fmt.Sprintf("%s:%d", d.Host, d.Port)

	err := d.Transport.DialSSH(d.Host, d.SSHConfig, d.Port)

	if err != nil {
		return err
	}

	d.Session, err = NewSession(d.Transport)

	if err != nil {
		return err
	}

	return nil
}

// DialTimeout function (call this after New())
func (d *DriverSSH) DialTimeout() error {
	d.Target = fmt.Sprintf("%s:%d", d.Host, d.Port)

	var err error

	d.Session, err = DialSSHTimeout(d.Target, d.SSHConfig, d.Timeout)

	if err != nil {
		return err
	}

	err = d.Transport.SetupSession()

	if err != nil {
		return err
	}

	return nil
}

// Close function closes the socket
func (d *DriverSSH) Close() error {

	// Close the SSH Session if we have one
	err := d.Session.Close()

	if err != nil {
		return err
	}

	return nil
}

// Lock the target datastore
func (d *DriverSSH) Lock(ds string) (*RPCReply, error) {
	reply, err := d.Session.Exec(MethodLock(ds))

	if err != nil {
		return reply, err
	}

	return reply, nil
}

// Unlock the target datastore
func (d *DriverSSH) Unlock(ds string) (*RPCReply, error) {
	reply, err := d.Session.Exec(MethodUnlock(ds))

	if err != nil {
		return reply, err
	}

	return reply, nil
}

// SendRaw sends a raw XML envelope
func (d *DriverSSH) SendRaw(rawxml string) (*RPCReply, error) {
	reply, err := d.Session.Exec(RawMethod(rawxml))

	if err != nil {
		return reply, err
	}

	return reply, nil
}

// GetConfig requests the contents of a datastore
func (d *DriverSSH) GetConfig() (*RPCReply, error) {
	reply, err := d.Session.Exec(MethodGetConfig(d.Datastore))

	if err != nil {
		return reply, err
	}

	return reply, nil
}
