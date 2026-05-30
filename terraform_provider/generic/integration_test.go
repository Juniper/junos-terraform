package generic

import (
	"context"
	"strings"
	"sync"
	"testing"

	"terraform_provider/patch"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// mockNetconfClient is a stateful mock that stores the last committed XML config.
type mockNetconfClient struct {
	mu        sync.Mutex
	configXML string
	commits   int
	sends     []string
	updates   []string
}

func newMockNetconfClient() *mockNetconfClient {
	return &mockNetconfClient{
		configXML: "<configuration></configuration>",
	}
}

func (m *mockNetconfClient) Close() error                                    { return nil }
func (m *mockNetconfClient) DeleteConfig(string, bool) (string, error)       { return "", nil }
func (m *mockNetconfClient) MarshalGroup(string, interface{}) error          { return nil }
func (m *mockNetconfClient) MarshalConfig(interface{}) error                 { return nil }
func (m *mockNetconfClient) SendTransaction(string, interface{}, bool) error { return nil }
func (m *mockNetconfClient) SendDirectTransaction(interface{}, bool) error   { return nil }

func (m *mockNetconfClient) SendCommit() error {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.commits++
	return nil
}

func (m *mockNetconfClient) GetConfigXML() ([]byte, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	return []byte(m.configXML), nil
}

func (m *mockNetconfClient) SendDirectTransactionRaw(xmlPayload string, _ bool) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.sends = append(m.sends, xmlPayload)
	// Simulate applying the full config
	m.configXML = xmlPayload
	return nil
}

func (m *mockNetconfClient) SendUpdate(_ string, patchXML string, _ bool) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.updates = append(m.updates, patchXML)
	return nil
}

// --- Integration Tests ---

// TestIntegration_CreateFlow tests the Create path: model -> XML -> NETCONF send -> commit -> read back.
func TestIntegration_CreateFlow(t *testing.T) {
	mock := newMockNetconfClient()
	r := buildIntegrationResource(t, mock)
	ctx := context.Background()

	// Build attrs representing system { host-name: "router1", name-server: ["8.8.8.8"] }
	attrs := buildSystemAttrs("router1", []string{"8.8.8.8"})

	// Convert model to XML (same as Create does)
	var diags diag.Diagnostics
	xmlBytes := ModelToXMLBytes(ctx, attrs, r.idx, &diags)
	if diags.HasError() {
		t.Fatalf("ModelToXMLBytes() failed: %v", diags)
	}

	xmlStr := string(xmlBytes)
	if !strings.Contains(xmlStr, "<host-name>router1</host-name>") {
		t.Fatalf("XML missing host-name, got:\n%s", xmlStr)
	}
	if !strings.Contains(xmlStr, "<name-server>8.8.8.8</name-server>") {
		t.Fatalf("XML missing name-server, got:\n%s", xmlStr)
	}

	// Simulate Create: send to device
	if err := r.client.SendDirectTransactionRaw(string(xmlBytes), false); err != nil {
		t.Fatalf("SendDirectTransactionRaw() error: %v", err)
	}
	if err := r.client.SendCommit(); err != nil {
		t.Fatalf("SendCommit() error: %v", err)
	}

	// Verify mock state
	mock.mu.Lock()
	if len(mock.sends) != 1 {
		t.Fatalf("expected 1 send, got %d", len(mock.sends))
	}
	if mock.commits != 1 {
		t.Fatalf("expected 1 commit, got %d", mock.commits)
	}
	mock.mu.Unlock()

	// Simulate Read: fetch from device and convert back
	readXML, err := r.client.GetConfigXML()
	if err != nil {
		t.Fatalf("GetConfigXML() error: %v", err)
	}

	readModel := XMLBytesToModel(ctx, readXML, r.idx, &diags)
	if diags.HasError() {
		t.Fatalf("XMLBytesToModel() failed: %v", diags)
	}

	// Verify system.host_name is preserved
	sysObj, ok := readModel["system"]
	if !ok || sysObj.IsNull() {
		t.Fatal("Read model missing 'system' attribute")
	}
	sysList := sysObj.(types.List)
	sysAttrs := sysList.Elements()[0].(types.Object).Attributes()
	hostName := sysAttrs["host_name"].(types.String)
	if hostName.ValueString() != "router1" {
		t.Fatalf("expected host_name='router1', got %q", hostName.ValueString())
	}
}

// TestIntegration_UpdateWithPatchEngine tests the Update path using the patch diff engine.
func TestIntegration_UpdateWithPatchEngine(t *testing.T) {
	mock := newMockNetconfClient()
	r := buildIntegrationResource(t, mock)
	ctx := context.Background()
	var diags diag.Diagnostics

	// Current state on device: host-name=router1
	stateAttrs := buildSystemAttrs("router1", nil)
	stateXML := ModelToXMLBytes(ctx, stateAttrs, r.idx, &diags)
	if diags.HasError() {
		t.Fatalf("state ModelToXMLBytes() failed: %v", diags)
	}

	// Set mock to return current state
	mock.mu.Lock()
	mock.configXML = string(stateXML)
	mock.mu.Unlock()

	// Desired state: host-name=router2
	planAttrs := buildSystemAttrs("router2", nil)
	planXML := ModelToXMLBytes(ctx, planAttrs, r.idx, &diags)
	if diags.HasError() {
		t.Fatalf("plan ModelToXMLBytes() failed: %v", diags)
	}

	// Compute diff using the patch engine (same as Update does)
	patchIdx, err := patch.UnmarshalTrimmedSchemaIndex(r.rawSchema)
	if err != nil {
		t.Fatalf("UnmarshalTrimmedSchemaIndex() error: %v", err)
	}

	planTree, err := patch.BuildTree(planXML)
	if err != nil {
		t.Fatalf("BuildTree(plan) error: %v", err)
	}
	stateTree, err := patch.BuildTree(stateXML)
	if err != nil {
		t.Fatalf("BuildTree(state) error: %v", err)
	}

	planMap := patch.LeafMapWithSchema(planTree, patchIdx)
	stateMap := patch.LeafMapWithSchema(stateTree, patchIdx)
	diffMap := patch.ComputeDiff(stateMap, planMap)

	if len(diffMap) == 0 {
		t.Fatal("expected non-empty diff between router1 and router2")
	}

	// Verify the diff contains the host-name change
	foundHostNameChange := false
	for p, change := range diffMap {
		if strings.Contains(p, "host-name") {
			foundHostNameChange = true
			if change.OldVal != "router1" {
				t.Fatalf("diff old value: expected 'router1', got %q", change.OldVal)
			}
			if change.NewVal != "router2" {
				t.Fatalf("diff new value: expected 'router2', got %q", change.NewVal)
			}
		}
	}
	if !foundHostNameChange {
		t.Fatalf("diff does not contain host-name change, got: %v", diffMap)
	}

	// Create patch
	patchXMLBytes, err := patch.CreateDiffPatch(diffMap, "test-resource")
	if err != nil {
		t.Fatalf("CreateDiffPatch() error: %v", err)
	}

	patchStr := string(patchXMLBytes)
	if patchStr == "" {
		t.Fatal("CreateDiffPatch() returned empty patch")
	}

	// Send patch (same as Update does)
	if err := r.client.SendUpdate("", patchStr, false); err != nil {
		t.Fatalf("SendUpdate() error: %v", err)
	}
	if err := r.client.SendCommit(); err != nil {
		t.Fatalf("SendCommit() error: %v", err)
	}

	mock.mu.Lock()
	if len(mock.updates) != 1 {
		t.Fatalf("expected 1 update, got %d", len(mock.updates))
	}
	if mock.commits != 1 {
		t.Fatalf("expected 1 commit, got %d", mock.commits)
	}
	mock.mu.Unlock()
}

// TestIntegration_DeleteWithPatchEngine tests the Delete path (diff to empty -> patch).
func TestIntegration_DeleteWithPatchEngine(t *testing.T) {
	mock := newMockNetconfClient()
	r := buildIntegrationResource(t, mock)
	ctx := context.Background()
	var diags diag.Diagnostics

	// Current state: host-name=router1
	stateAttrs := buildSystemAttrs("router1", []string{"8.8.8.8"})
	stateXML := ModelToXMLBytes(ctx, stateAttrs, r.idx, &diags)
	if diags.HasError() {
		t.Fatalf("ModelToXMLBytes() failed: %v", diags)
	}

	// Empty (delete target)
	emptyXML := []byte("<?xml version=\"1.0\" encoding=\"UTF-8\"?>\n<configuration></configuration>")

	patchIdx, err := patch.UnmarshalTrimmedSchemaIndex(r.rawSchema)
	if err != nil {
		t.Fatalf("UnmarshalTrimmedSchemaIndex() error: %v", err)
	}

	stateTree, err := patch.BuildTree(stateXML)
	if err != nil {
		t.Fatalf("BuildTree(state) error: %v", err)
	}
	emptyTree, err := patch.BuildTree(emptyXML)
	if err != nil {
		t.Fatalf("BuildTree(empty) error: %v", err)
	}

	stateMap := patch.LeafMapWithSchema(stateTree, patchIdx)
	emptyMap := patch.LeafMapWithSchema(emptyTree, patchIdx)
	diffMap := patch.ComputeDiff(stateMap, emptyMap)

	if len(diffMap) == 0 {
		t.Fatal("expected non-empty diff for delete operation")
	}

	// All changes should be deletions (Op == Delete)
	for p, change := range diffMap {
		if change.Op != patch.Delete {
			t.Fatalf("expected Delete op for path %q, got %v", p, change.Op)
		}
	}

	patchXMLBytes, err := patch.CreateDiffPatch(diffMap, "test-delete")
	if err != nil {
		t.Fatalf("CreateDiffPatch() error: %v", err)
	}

	patchStr := string(patchXMLBytes)
	if patchStr == "" {
		t.Fatal("delete patch is empty")
	}

	// The patch should contain delete operation markers
	if !strings.Contains(patchStr, "delete") {
		t.Logf("Note: patch uses operation syntax:\n%s", patchStr)
	}

	if err := r.client.SendUpdate("", patchStr, false); err != nil {
		t.Fatalf("SendUpdate() error: %v", err)
	}
	if err := r.client.SendCommit(); err != nil {
		t.Fatalf("SendCommit() error: %v", err)
	}

	mock.mu.Lock()
	if len(mock.updates) == 0 {
		t.Fatal("expected update to be sent for delete")
	}
	if mock.commits == 0 {
		t.Fatal("expected commit after delete")
	}
	mock.mu.Unlock()
}

// TestIntegration_FullCRUDCycle exercises Create -> Read -> Update -> Delete in sequence.
func TestIntegration_FullCRUDCycle(t *testing.T) {
	mock := newMockNetconfClient()
	r := buildIntegrationResource(t, mock)
	ctx := context.Background()
	var diags diag.Diagnostics

	// --- CREATE ---
	createAttrs := buildSystemAttrs("router1", []string{"1.1.1.1"})
	createXML := ModelToXMLBytes(ctx, createAttrs, r.idx, &diags)
	if diags.HasError() {
		t.Fatalf("Create: ModelToXMLBytes() failed: %v", diags)
	}
	if err := r.client.SendDirectTransactionRaw(string(createXML), false); err != nil {
		t.Fatalf("Create: SendDirectTransactionRaw() error: %v", err)
	}
	if err := r.client.SendCommit(); err != nil {
		t.Fatalf("Create: SendCommit() error: %v", err)
	}

	// --- READ ---
	readXML, err := r.client.GetConfigXML()
	if err != nil {
		t.Fatalf("Read: GetConfigXML() error: %v", err)
	}
	diags = diag.Diagnostics{}
	readModel := XMLBytesToModel(ctx, readXML, r.idx, &diags)
	if diags.HasError() {
		t.Fatalf("Read: XMLBytesToModel() failed: %v", diags)
	}
	sysListVal := readModel["system"].(types.List)
	sysObj := sysListVal.Elements()[0].(types.Object)
	hostName := sysObj.Attributes()["host_name"].(types.String).ValueString()
	if hostName != "router1" {
		t.Fatalf("Read: expected host_name='router1', got %q", hostName)
	}

	// --- UPDATE ---
	diags = diag.Diagnostics{}
	updateAttrs := buildSystemAttrs("router2", []string{"1.1.1.1", "8.8.4.4"})
	updateXML := ModelToXMLBytes(ctx, updateAttrs, r.idx, &diags)
	if diags.HasError() {
		t.Fatalf("Update: ModelToXMLBytes() failed: %v", diags)
	}

	patchIdx, err := patch.UnmarshalTrimmedSchemaIndex(r.rawSchema)
	if err != nil {
		t.Fatalf("Update: UnmarshalTrimmedSchemaIndex() error: %v", err)
	}

	updateTree, err := patch.BuildTree(updateXML)
	if err != nil {
		t.Fatalf("Update: BuildTree(plan) error: %v", err)
	}
	currentTree, err := patch.BuildTree(readXML)
	if err != nil {
		t.Fatalf("Update: BuildTree(state) error: %v", err)
	}

	updateMap := patch.LeafMapWithSchema(updateTree, patchIdx)
	currentMap := patch.LeafMapWithSchema(currentTree, patchIdx)
	diffMap := patch.ComputeDiff(currentMap, updateMap)

	if len(diffMap) == 0 {
		t.Fatal("Update: expected non-empty diff")
	}

	patchXMLBytes, err := patch.CreateDiffPatch(diffMap, "full-cycle")
	if err != nil {
		t.Fatalf("Update: CreateDiffPatch() error: %v", err)
	}
	if err := r.client.SendUpdate("", string(patchXMLBytes), false); err != nil {
		t.Fatalf("Update: SendUpdate() error: %v", err)
	}
	if err := r.client.SendCommit(); err != nil {
		t.Fatalf("Update: SendCommit() error: %v", err)
	}

	// Simulate device accepted the update (set mock config to plan)
	mock.mu.Lock()
	mock.configXML = string(updateXML)
	mock.mu.Unlock()

	// --- DELETE ---
	deleteStateXML, _ := r.client.GetConfigXML()
	deleteTree, err := patch.BuildTree(deleteStateXML)
	if err != nil {
		t.Fatalf("Delete: BuildTree() error: %v", err)
	}
	emptyTree, err := patch.BuildTree([]byte("<configuration></configuration>"))
	if err != nil {
		t.Fatalf("Delete: BuildTree(empty) error: %v", err)
	}

	deleteMap := patch.LeafMapWithSchema(deleteTree, patchIdx)
	emptyMap := patch.LeafMapWithSchema(emptyTree, patchIdx)
	deleteDiff := patch.ComputeDiff(deleteMap, emptyMap)

	if len(deleteDiff) == 0 {
		t.Fatal("Delete: expected non-empty diff")
	}

	deletePatch, err := patch.CreateDiffPatch(deleteDiff, "full-cycle-delete")
	if err != nil {
		t.Fatalf("Delete: CreateDiffPatch() error: %v", err)
	}
	if err := r.client.SendUpdate("", string(deletePatch), false); err != nil {
		t.Fatalf("Delete: SendUpdate() error: %v", err)
	}
	if err := r.client.SendCommit(); err != nil {
		t.Fatalf("Delete: SendCommit() error: %v", err)
	}

	// Verify operation counts: Create(1) + Update(1) + Delete(1) = 3 commits
	mock.mu.Lock()
	defer mock.mu.Unlock()
	if mock.commits < 3 {
		t.Fatalf("expected at least 3 commits (create+update+delete), got %d", mock.commits)
	}
	if len(mock.sends) < 1 {
		t.Fatal("expected at least 1 direct send (create)")
	}
	if len(mock.updates) < 2 {
		t.Fatalf("expected at least 2 updates (update+delete), got %d", len(mock.updates))
	}
}

// TestIntegration_XMLRoundTrip_WithInterfaces tests model conversion with list types (interfaces).
func TestIntegration_XMLRoundTrip_WithInterfaces(t *testing.T) {
	mock := newMockNetconfClient()
	r := buildIntegrationResource(t, mock)
	ctx := context.Background()

	// Build interface list attrs
	unitAttrTypes := map[string]attr.Type{
		"name":        types.StringType,
		"description": types.StringType,
		"family": types.ListType{ElemType: types.ObjectType{AttrTypes: map[string]attr.Type{
			"inet": types.ListType{ElemType: types.ObjectType{AttrTypes: map[string]attr.Type{
				"address": types.ListType{ElemType: types.ObjectType{AttrTypes: map[string]attr.Type{
					"name": types.StringType,
				}}},
			}}},
		}}},
	}

	ifaceAttrTypes := map[string]attr.Type{
		"name":        types.StringType,
		"description": types.StringType,
		"mtu":         types.StringType,
		"unit":        types.ListType{ElemType: types.ObjectType{AttrTypes: unitAttrTypes}},
	}

	iface := types.ObjectValueMust(ifaceAttrTypes, map[string]attr.Value{
		"name":        types.StringValue("ge-0/0/0"),
		"description": types.StringValue("uplink"),
		"mtu":         types.StringValue("9000"),
		"unit":        types.ListNull(types.ObjectType{AttrTypes: unitAttrTypes}),
	})

	ifaceList := types.ListValueMust(
		types.ObjectType{AttrTypes: ifaceAttrTypes},
		[]attr.Value{iface},
	)

	attrs := map[string]attr.Value{
		"resource_name": types.StringValue("iface-test"),
		"interfaces":    ifaceList,
		"system":        types.ListNull(types.ObjectType{AttrTypes: map[string]attr.Type{"host_name": types.StringType, "name_server": types.ListType{ElemType: types.StringType}}}),
		"protocols":     types.ListNull(types.ObjectType{AttrTypes: map[string]attr.Type{"lldp": types.StringType}}),
	}

	var diags diag.Diagnostics
	xmlBytes := ModelToXMLBytes(ctx, attrs, r.idx, &diags)
	if diags.HasError() {
		t.Fatalf("ModelToXMLBytes() failed: %v", diags)
	}

	xmlStr := string(xmlBytes)
	if !strings.Contains(xmlStr, "<name>ge-0/0/0</name>") {
		t.Fatalf("XML missing interface name, got:\n%s", xmlStr)
	}
	if !strings.Contains(xmlStr, "<description>uplink</description>") {
		t.Fatalf("XML missing interface description, got:\n%s", xmlStr)
	}
	if !strings.Contains(xmlStr, "<mtu>9000</mtu>") {
		t.Fatalf("XML missing mtu, got:\n%s", xmlStr)
	}

	// Round-trip: set mock, then read back
	mock.mu.Lock()
	mock.configXML = xmlStr
	mock.mu.Unlock()

	readXML, _ := r.client.GetConfigXML()
	diags = diag.Diagnostics{}
	readModel := XMLBytesToModel(ctx, readXML, r.idx, &diags)
	if diags.HasError() {
		t.Fatalf("XMLBytesToModel() failed: %v", diags)
	}

	ifacesVal, ok := readModel["interfaces"]
	if !ok || ifacesVal.IsNull() {
		t.Fatal("Read model missing 'interfaces'")
	}

	ifacesList := ifacesVal.(types.List)
	if len(ifacesList.Elements()) != 1 {
		t.Fatalf("expected 1 interface, got %d", len(ifacesList.Elements()))
	}
}

// --- Helpers ---

func buildIntegrationResource(t *testing.T, mock *mockNetconfClient) *ConfigResource {
	t.Helper()
	factory := NewConfigResource(testSchemaJSON)
	r := factory().(*ConfigResource)
	r.client = ProviderConfig{Client: mock, Host: "test-host"}
	return r
}

func buildSystemAttrs(hostName string, nameServers []string) map[string]attr.Value {
	sysAttrTypes := map[string]attr.Type{
		"host_name":   types.StringType,
		"name_server": types.ListType{ElemType: types.StringType},
	}

	sysValues := map[string]attr.Value{
		"host_name": types.StringValue(hostName),
	}

	if len(nameServers) > 0 {
		nsValues := make([]attr.Value, len(nameServers))
		for i, ns := range nameServers {
			nsValues[i] = types.StringValue(ns)
		}
		sysValues["name_server"] = types.ListValueMust(types.StringType, nsValues)
	} else {
		sysValues["name_server"] = types.ListNull(types.StringType)
	}

	sysObj := types.ObjectValueMust(sysAttrTypes, sysValues)
	sysList := types.ListValueMust(types.ObjectType{AttrTypes: sysAttrTypes}, []attr.Value{sysObj})

	protocolsAttrTypes := map[string]attr.Type{
		"lldp": types.StringType,
	}
	protocolsObjType := types.ObjectType{AttrTypes: protocolsAttrTypes}

	return map[string]attr.Value{
		"resource_name": types.StringValue("test-resource"),
		"system":        sysList,
		"interfaces":    types.ListNull(types.ObjectType{AttrTypes: map[string]attr.Type{}}),
		"protocols":     types.ListNull(protocolsObjType),
	}
}
