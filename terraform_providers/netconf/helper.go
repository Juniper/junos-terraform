package netconf

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"strings"
	"sync"

	driver "github.com/davedotdev/go-netconf/drivers/driver"
	sshdriver "github.com/davedotdev/go-netconf/drivers/ssh"

	"golang.org/x/crypto/ssh"
)

const groupStrXML = `<load-configuration action="merge" format="xml">
%s
</load-configuration>
`

const deleteStr = `<edit-config>
	<target>
		<candidate/>
	</target>
	<default-operation>none</default-operation> 
	<config>
		<configuration>
			<groups operation="delete">
				<name>%s</name>
			</groups>
			<apply-groups operation="delete">%s</apply-groups>
		</configuration>
	</config>
</edit-config>`

const commitStr = `<commit/>`

const getGroupXMLStr = `<get-configuration>
  <configuration>
  <groups><name>%s</name></groups>
  </configuration>
</get-configuration>
`

// GoNCClient type for storing data and wrapping functions
type GoNCClient struct {
	Driver driver.Driver
	Lock   sync.RWMutex
}

// Close is a functional thing to close the Driver
func (g *GoNCClient) Close() error {
	g.Driver = nil
	return nil
}

// updateRawConfig deletes group data and replaces it (for Update in TF)
func (g *GoNCClient) updateRawConfig(applyGroup string, netconfCall string, commit bool) (string, error) {

	g.Lock.Lock()
	defer g.Lock.Unlock()

	if err := g.Driver.Dial(); err != nil {
		return "", err
	}

	deleteString := fmt.Sprintf(deleteStr, applyGroup, applyGroup)

	if _, err := g.Driver.SendRaw(deleteString); err != nil {
		errInternal := g.Driver.Close()
		return "", fmt.Errorf("driver error: %+v, driver close error: %s", err, errInternal)
	}

	groupString := fmt.Sprintf(groupStrXML, netconfCall)

	reply, err := g.Driver.SendRaw(groupString)
	if err != nil {
		errInternal := g.Driver.Close()
		return "", fmt.Errorf("driver error: %+v, driver close error: %s", err, errInternal)
	}

	if commit {
		if _, err = g.Driver.SendRaw(commitStr); err != nil {
			errInternal := g.Driver.Close()
			return "", fmt.Errorf("driver error: %+v, driver close error: %s", err, errInternal)
		}
	}

	if err := g.Driver.Close(); err != nil {
		return "", fmt.Errorf("driver close error: %s", err)
	}
	return reply.Data, nil
}

// DeleteConfig is a wrapper for driver.SendRaw()
func (g *GoNCClient) DeleteConfig(applyGroup string, commit bool) (string, error) {

	g.Lock.Lock()
	defer g.Lock.Unlock()

	if err := g.Driver.Dial(); err != nil {
		return "", err
	}

	deleteString := fmt.Sprintf(deleteStr, applyGroup, applyGroup)

	reply, err := g.Driver.SendRaw(deleteString)
	if err != nil {
		errInternal := g.Driver.Close()
		return "", fmt.Errorf("driver error: %+v, driver close error: %s", err, errInternal)
	}
	if commit {
		if _, err = g.Driver.SendRaw(commitStr); err != nil {
			errInternal := g.Driver.Close()
			return "", fmt.Errorf("driver error: %+v, driver close error: %s", err, errInternal)
		}

	}

	output := strings.Replace(reply.Data, "\n", "", -1)

	if err := g.Driver.Close(); err != nil {
		return "", err
	}
	return output, nil
}

// SendCommit is a wrapper for driver.SendRaw()
func (g *GoNCClient) SendCommit() error {
	g.Lock.Lock()
	defer g.Lock.Unlock()

	if err := g.Driver.Dial(); err != nil {
		return err
	}
	if _, err := g.Driver.SendRaw(commitStr); err != nil {
		return err
	}
	return nil
}

// MarshalGroup accepts a struct of type X and then marshals data onto it
func (g *GoNCClient) MarshalGroup(id string, obj interface{}) error {

	reply, err := g.readRawGroup(id)
	if err != nil {
		return err
	}

	if err = xml.Unmarshal([]byte(reply), &obj); err != nil {
		return err
	}
	return nil
}

// SendTransaction is a method that unmarshal the XML, creates the transaction and passes in a commit
func (g *GoNCClient) SendTransaction(id string, obj interface{}, commit bool) error {
	cfg, err := xml.Marshal(obj)
	if err != nil {
		return err
	}
	// updateRawConfig deletes old group by, re-creates it then commits.
	// As far as Junos cares, it's an edit.
	if id != "" {
		if _, err = g.updateRawConfig(id, string(cfg), commit); err != nil {
			return err
		}
		return nil
	}
	if _, err = g.sendRawConfig(string(cfg), commit); err != nil {
		return err
	}
	return nil
}

// sendRawConfig is a wrapper for driver.SendRaw()
func (g *GoNCClient) sendRawConfig(netconfCall string, commit bool) (string, error) {
	g.Lock.Lock()
	defer g.Lock.Unlock()

	if err := g.Driver.Dial(); err != nil {
		return "", err
	}
	groupString := fmt.Sprintf(groupStrXML, netconfCall)

	reply, err := g.Driver.SendRaw(groupString)
	if err != nil {
		errInternal := g.Driver.Close()
		return "", fmt.Errorf("driver error: %+v, driver close error: %s", err, errInternal)
	}
	if commit {
		_, err = g.Driver.SendRaw(commitStr)
		if err != nil {
			errInternal := g.Driver.Close()
			return "", fmt.Errorf("driver error: %+v, driver close error: %s", err, errInternal)
		}
	}
	if err = g.Driver.Close(); err != nil {
		return "", err
	}
	return reply.Data, nil
}

// readRawGroup is a helper function
func (g *GoNCClient) readRawGroup(applyGroup string) (string, error) {
	g.Lock.Lock()
	defer g.Lock.Unlock()

	if err := g.Driver.Dial(); err != nil {
		return "", err
	}
	getGroupXMLString := fmt.Sprintf(getGroupXMLStr, applyGroup)

	reply, err := g.Driver.SendRaw(getGroupXMLString)
	if err != nil {
		errInternal := g.Driver.Close()
		return "", fmt.Errorf("driver error: %+v, driver close error: %s", err, errInternal)
	}

	if err = g.Driver.Close(); err != nil {
		return "", err
	}
	return reply.Data, nil
}

func publicKeyFile(file string) ssh.AuthMethod {
	buffer, err := ioutil.ReadFile(file)
	if err != nil {
		return nil
	}

	key, err := ssh.ParsePrivateKey(buffer)
	if err != nil {
		return nil
	}
	return ssh.PublicKeys(key)
}

// NewClient returns go-netconf new client driver
func NewClient(username string, password string, sshKey string, address string, port int) (Client, error) {

	// Dummy interface var ready for loading from inputs
	var nconf driver.Driver

	d := driver.New(sshdriver.New())

	nc := d.(*sshdriver.DriverSSH)

	nc.Host = address
	nc.Port = port

	// SSH keys takes priority over password based
	if sshKey != "" {
		nc.SSHConfig = &ssh.ClientConfig{
			User: username,
			Auth: []ssh.AuthMethod{
				publicKeyFile(sshKey),
			},
			HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		}
	} else {
		// Sort yourself out with SSH. Easiest to do that here.
		nc.SSHConfig = &ssh.ClientConfig{
			User:            username,
			Auth:            []ssh.AuthMethod{ssh.Password(password)},
			HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		}
	}

	nconf = nc

	return &GoNCClient{Driver: nconf}, nil
}
