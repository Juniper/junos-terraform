package netconf

import (
	"context"
	"encoding/xml"
	"fmt"
	"os"
	"sort"
	"strings"
	"sync"

	"golang.org/x/crypto/ssh"
	netconf "nemith.io/netconf"
	netconfssh "nemith.io/netconf/transport/ssh"
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

const applyGroupXML = `<load-configuration action="merge" format="xml">
	%s
</load-configuration>
`

const discardChanges = `<discard-changes/>`

// defaultPort is the NETCONF-over-SSH default.
const defaultPort = 830

type configuration struct {
	ApplyGroup []string `xml:"apply-groups"`
}

type operationExecutor func(ctx context.Context, operation string) (string, error)

// GoNCClient implements the provider-facing NETCONF client API on top of nemith/netconf.
type GoNCClient struct {
	host      string
	port      int
	sshConfig *ssh.ClientConfig

	Lock sync.RWMutex
	exec operationExecutor
}

// Close keeps existing behavior contract for provider lifecycle hooks.
func (g *GoNCClient) Close() error {
	return nil
}

func (g *GoNCClient) execute(ctx context.Context, operation string) (string, error) {
	if g.exec != nil {
		return g.exec(ctx, operation)
	}

	target := fmt.Sprintf("%s:%d", g.host, g.port)
	transport, err := netconfssh.Dial(ctx, "tcp", target, g.sshConfig)
	if err != nil {
		return "", err
	}

	session, err := netconf.NewSession(transport)
	if err != nil {
		_ = transport.Close()
		return "", err
	}
	defer func() {
		_ = session.Close(context.Background())
	}()

	reply := struct {
		XMLName xml.Name `xml:"rpc-reply"`
		Data    string   `xml:",innerxml"`
	}{}

	if err := session.Exec(ctx, netconf.RawXML(operation), &reply); err != nil {
		return "", err
	}

	return reply.Data, nil
}

func (g *GoNCClient) updateRawConfig(applyGroup string, netconfCall string, commit bool) (string, error) {
	g.Lock.Lock()
	defer g.Lock.Unlock()

	ctx := context.Background()
	deleteString := fmt.Sprintf(deleteStr, applyGroup, applyGroup)
	if _, err := g.execute(ctx, deleteString); err != nil {
		return "", err
	}

	nameStart := strings.Index(netconfCall, "<name>")
	nameEnd := strings.Index(netconfCall, "</name>")
	if nameStart == -1 || nameEnd == -1 {
		return "", fmt.Errorf("failed to extract the group name from the netconfcall")
	}
	groupName := netconfCall[nameStart+6 : nameEnd]
	addToApplyGroupsList(groupName)

	groupString := fmt.Sprintf(groupStrXML, netconfCall)
	reply, err := g.execute(ctx, groupString)
	if err != nil {
		return "", err
	}

	if commit {
		if _, err := g.execute(ctx, commitStr); err != nil {
			return "", err
		}
	}

	return reply, nil
}

// DeleteConfig deletes the target apply-group and optionally commits.
func (g *GoNCClient) DeleteConfig(applyGroup string, commit bool) (string, error) {
	g.Lock.Lock()
	defer g.Lock.Unlock()

	ctx := context.Background()
	deleteString := fmt.Sprintf(deleteStr, applyGroup, applyGroup)
	reply, err := g.execute(ctx, deleteString)
	if err != nil {
		return "", err
	}

	if commit {
		if _, err := g.execute(ctx, commitStr); err != nil {
			return "", err
		}
	}

	return strings.ReplaceAll(reply, "\n", ""), nil
}

// SendCommit emits apply-groups in deterministic order and commits candidate config.
func (g *GoNCClient) SendCommit() error {
	g.Lock.Lock()
	defer g.Lock.Unlock()

	sortApplyGroupsList()
	if err := g.sendApplyGroupsLocked(context.Background()); err != nil {
		return err
	}

	if _, err := g.execute(context.Background(), commitStr); err != nil {
		_, _ = g.execute(context.Background(), discardChanges)
		return err
	}

	return nil
}

func (g *GoNCClient) sendApplyGroupsLocked(ctx context.Context) error {
	applyGroupsMutex.Lock()
	applyGroupsCopy := make([]string, len(applyGroupsList))
	copy(applyGroupsCopy, applyGroupsList)
	applyGroupsMutex.Unlock()

	var applyG configuration
	applyG.ApplyGroup = applyGroupsCopy

	cfg, err := xml.Marshal(applyG)
	if err != nil {
		return err
	}

	_, err = g.execute(ctx, fmt.Sprintf(applyGroupXML, string(cfg)))
	return err
}

// MarshalGroup fetches a group and unmarshals XML into obj.
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

var applyGroupsList []string
var applyGroupsMutex sync.Mutex

// SendTransaction updates or creates a config payload and optionally commits it.
func (g *GoNCClient) SendTransaction(id string, obj interface{}, commit bool) error {
	cfg, err := xml.Marshal(obj)
	if err != nil {
		return err
	}

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

func addToApplyGroupsList(id string) {
	applyGroupsMutex.Lock()
	defer applyGroupsMutex.Unlock()
	applyGroupsList = append(applyGroupsList, id)
}

func sortApplyGroupsList() {
	applyGroupsMutex.Lock()
	defer applyGroupsMutex.Unlock()

	filteredGroups := make([]string, 0, len(applyGroupsList))
	for _, group := range applyGroupsList {
		if group != "" {
			filteredGroups = append(filteredGroups, group)
		}
	}
	sort.Strings(filteredGroups)
	applyGroupsList = filteredGroups
}

func (g *GoNCClient) sendRawConfig(netconfCall string, commit bool) (string, error) {
	g.Lock.Lock()
	defer g.Lock.Unlock()

	nameStart := strings.Index(netconfCall, "<name>")
	nameEnd := strings.Index(netconfCall, "</name>")
	if nameStart == -1 || nameEnd == -1 {
		return "", fmt.Errorf("failed to extract the group name from the netconfCall")
	}
	groupName := netconfCall[nameStart+6 : nameEnd]
	addToApplyGroupsList(groupName)

	reply, err := g.execute(context.Background(), fmt.Sprintf(groupStrXML, netconfCall))
	if err != nil {
		return "", err
	}

	if commit {
		if _, err = g.execute(context.Background(), commitStr); err != nil {
			return "", err
		}
	}

	return reply, nil
}

func (g *GoNCClient) readRawGroup(applyGroup string) (string, error) {
	g.Lock.Lock()
	defer g.Lock.Unlock()

	return g.execute(context.Background(), fmt.Sprintf(getGroupXMLStr, applyGroup))
}

func publicKeyFile(file string) ssh.AuthMethod {
	buffer, err := os.ReadFile(file)
	if err != nil {
		return nil
	}

	key, err := ssh.ParsePrivateKey(buffer)
	if err != nil {
		return nil
	}
	return ssh.PublicKeys(key)
}

// NewClient returns a NETCONF client backed by nemith/netconf.
func NewClient(username, password, sshKey, address string, port int) (Client, error) {
	if port == 0 {
		port = defaultPort
	}

	cfg := &ssh.ClientConfig{
		User:            username,
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}

	if sshKey != "" {
		authMethod := publicKeyFile(sshKey)
		cfg.Auth = []ssh.AuthMethod{authMethod}
	} else {
		cfg.Auth = []ssh.AuthMethod{ssh.Password(password)}
	}

	return &GoNCClient{
		host:      address,
		port:      port,
		sshConfig: cfg,
	}, nil
}
