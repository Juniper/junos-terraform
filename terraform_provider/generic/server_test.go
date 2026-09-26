package generic

import (
	"bytes"
	"compress/gzip"
	"context"
	"encoding/xml"
	"strings"
	"testing"

	"terraform_provider/patch"

	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
)

// fakeDevice is a NETCONF client holding a device's configuration in memory.
type fakeDevice struct {
	config  string // content of <configuration>
	reads   int
	loads   []string
	patches []string
	commits int
	// onPatch, if set, gives the configuration after a patch; by default a
	// patch leaves it unchanged.
	onPatch func(patch string) string
	closed  bool
}

func (f *fakeDevice) Close() error                                    { f.closed = true; return nil }
func (f *fakeDevice) DeleteConfig(string, bool) (string, error)       { return "", nil }
func (f *fakeDevice) MarshalGroup(string, interface{}) error          { return nil }
func (f *fakeDevice) SendTransaction(string, interface{}, bool) error { return nil }
func (f *fakeDevice) SendCommit() error                               { f.commits++; return nil }

func (f *fakeDevice) MarshalConfig(obj interface{}) error {
	f.reads++
	return xml.Unmarshal([]byte(`<configuration junos:changed-seconds="1" xmlns:junos="http://xml.juniper.net/junos">`+f.config+`</configuration>`), obj)
}

// SendDirectTransaction loads a configuration; the fake replaces its own.
func (f *fakeDevice) SendDirectTransaction(obj interface{}, _ bool) error {
	b, err := xml.Marshal(obj)
	if err != nil {
		return err
	}
	f.loads = append(f.loads, string(b))
	f.config = strings.TrimSuffix(strings.TrimPrefix(string(b), "<configuration>"), "</configuration>")
	return nil
}

func (f *fakeDevice) SendUpdate(_ string, diff string, _ bool) error {
	f.patches = append(f.patches, diff)
	if f.onPatch != nil {
		f.config = f.onPatch(diff)
	}
	return nil
}

// testServer returns a server for convertSchema using client.
func testServer(t *testing.T, client *fakeDevice) (*Server, tftypes.Object) {
	t.Helper()
	ResetSchema()
	t.Cleanup(ResetSchema)
	s := NewServer("junos-test", []byte(convertSchema))
	l, err := s.load()
	if err != nil {
		t.Fatal(err)
	}
	s.client = client
	return s, l.typ
}

// resourceValue is the resource's value for a configuration's content.
func resourceValue(t *testing.T, typ tftypes.Object, name, config string) tftypes.Value {
	t.Helper()
	nodes, _ := convertFixture(t)
	tree, err := patch.BuildTree([]byte("<configuration>" + config + "</configuration>"))
	if err != nil {
		t.Fatal(err)
	}
	v, err := ConfigToValue(tree, nodes, typ)
	if err != nil {
		t.Fatal(err)
	}
	var attrs map[string]tftypes.Value
	if err := v.As(&attrs); err != nil {
		t.Fatal(err)
	}
	attrs[resourceNameAttribute] = tftypes.NewValue(tftypes.String, name)
	return tftypes.NewValue(typ, attrs)
}

func dv(t *testing.T, typ tftypes.Type, v tftypes.Value) *tfprotov6.DynamicValue {
	t.Helper()
	d, err := dynamicValue(typ, v)
	if err != nil {
		t.Fatal(err)
	}
	return d
}

func decode(t *testing.T, typ tftypes.Type, d *tfprotov6.DynamicValue) tftypes.Value {
	t.Helper()
	v, err := d.Unmarshal(typ)
	if err != nil {
		t.Fatal(err)
	}
	return v
}

func noDiags(t *testing.T, diags []*tfprotov6.Diagnostic) {
	t.Helper()
	for _, d := range diags {
		t.Errorf("%s: %s", d.Summary, d.Detail)
	}
	if len(diags) > 0 {
		t.FailNow()
	}
}

const (
	hostA = `<system><host-name>a</host-name></system>`
	hostB = `<system><host-name>b</host-name></system>`
)

func TestGetProviderSchema(t *testing.T) {
	s, _ := testServer(t, &fakeDevice{})
	resp, _ := s.GetProviderSchema(context.Background(), &tfprotov6.GetProviderSchemaRequest{})
	noDiags(t, resp.Diagnostics)
	rs := resp.ResourceSchemas["terraform-provider-junos-test"]
	if rs == nil || attribute(rs.Block.Attributes, "system") == nil || attribute(rs.Block.Attributes, "version") == nil {
		t.Fatalf("resource schema: %+v", rs)
	}
	if attribute(resp.Provider.Block.Attributes, "host") == nil {
		t.Fatal("provider schema has no host")
	}
}

func gzipped(t *testing.T, s string) []byte {
	t.Helper()
	var b bytes.Buffer
	w := gzip.NewWriter(&b)
	if _, err := w.Write([]byte(s)); err != nil {
		t.Fatal(err)
	}
	if err := w.Close(); err != nil {
		t.Fatal(err)
	}
	return b.Bytes()
}

func TestGetProviderSchemaGzipped(t *testing.T) {
	ResetSchema()
	t.Cleanup(ResetSchema)
	resp, _ := NewServer("junos-test", gzipped(t, convertSchema)).GetProviderSchema(context.Background(), &tfprotov6.GetProviderSchemaRequest{})
	noDiags(t, resp.Diagnostics)
}

// A schema that does not load is reported, not a provider without resources.
func TestGetProviderSchemaReportsErrors(t *testing.T) {
	for name, schema := range map[string]string{
		"load error": `{"root": `,
		"collision": `{"root": {"children": [{"name": "configuration", "type": "container", "children": [
			{"name": "a-b", "type": "leaf"}, {"name": "a_b", "type": "leaf"}]}]}}`,
	} {
		ResetSchema()
		resp, _ := NewServer("junos-test", []byte(schema)).GetProviderSchema(context.Background(), &tfprotov6.GetProviderSchemaRequest{})
		if len(resp.Diagnostics) != 1 || resp.Diagnostics[0].Summary != "Invalid embedded schema" {
			t.Errorf("%s: diagnostics %+v", name, resp.Diagnostics)
		}
	}
	ResetSchema()
}

func TestPlanResourceChange(t *testing.T) {
	s, typ := testServer(t, &fakeDevice{})
	prior := resourceValue(t, typ, "r", hostA)
	plan := func(proposed tftypes.Value) *tfprotov6.PlanResourceChangeResponse {
		resp, _ := s.PlanResourceChange(context.Background(), &tfprotov6.PlanResourceChangeRequest{
			PriorState: dv(t, typ, prior), ProposedNewState: dv(t, typ, proposed),
		})
		noDiags(t, resp.Diagnostics)
		if !decode(t, typ, resp.PlannedState).Equal(proposed) {
			t.Fatal("planned state is not the proposed state")
		}
		return resp
	}
	if resp := plan(resourceValue(t, typ, "r", hostB)); len(resp.RequiresReplace) != 0 {
		t.Fatalf("a configuration change requires replacement: %v", resp.RequiresReplace)
	}
	if resp := plan(resourceValue(t, typ, "other", hostA)); len(resp.RequiresReplace) != 1 {
		t.Fatalf("a resource_name change does not require replacement")
	}
}

func apply(t *testing.T, s *Server, typ tftypes.Object, prior, planned tftypes.Value) tftypes.Value {
	t.Helper()
	resp, _ := s.ApplyResourceChange(context.Background(), &tfprotov6.ApplyResourceChangeRequest{
		PriorState: dv(t, typ, prior), PlannedState: dv(t, typ, planned),
	})
	noDiags(t, resp.Diagnostics)
	return decode(t, typ, resp.NewState)
}

func TestApplyCreate(t *testing.T) {
	f := &fakeDevice{}
	s, typ := testServer(t, f)
	planned := resourceValue(t, typ, "r", hostA)
	state := apply(t, s, typ, tftypes.NewValue(typ, nil), planned)
	if len(f.loads) != 1 || !strings.Contains(f.loads[0], "<host-name>a</host-name>") || f.commits != 1 {
		t.Fatalf("loads %v, commits %d", f.loads, f.commits)
	}
	if !state.Equal(planned) {
		t.Fatalf("state %s, want %s", state, planned)
	}
}

func TestApplyUpdate(t *testing.T) {
	f := &fakeDevice{config: hostA, onPatch: func(string) string { return hostB }}
	s, typ := testServer(t, f)
	planned := resourceValue(t, typ, "r", hostB)
	state := apply(t, s, typ, resourceValue(t, typ, "r", hostA), planned)
	if len(f.patches) != 1 || !strings.Contains(f.patches[0], `<host-name nc:operation="replace">b</host-name>`) {
		t.Fatalf("patches %v", f.patches)
	}
	if len(f.loads) != 0 || f.commits != 1 {
		t.Fatalf("fallback load %v, commits %d", f.loads, f.commits)
	}
	// The configuration before the patch and after it, which is the state.
	if f.reads != 2 {
		t.Fatalf("reads %d, want 2", f.reads)
	}
	if !state.Equal(planned) {
		t.Fatalf("state %s, want %s", state, planned)
	}
}

// A patch that leaves the device different from the plan is followed by a
// load of the whole plan.
func TestApplyUpdateFallback(t *testing.T) {
	f := &fakeDevice{config: hostA}
	s, typ := testServer(t, f)
	planned := resourceValue(t, typ, "r", hostB)
	state := apply(t, s, typ, resourceValue(t, typ, "r", hostA), planned)
	if len(f.patches) != 1 || len(f.loads) != 1 || f.commits != 2 || f.reads != 3 {
		t.Fatalf("patches %v, loads %v, commits %d, reads %d", f.patches, f.loads, f.commits, f.reads)
	}
	if !state.Equal(planned) {
		t.Fatalf("state %s, want %s", state, planned)
	}
}

// No difference: nothing is sent or committed.
func TestApplyUpdateUnchanged(t *testing.T) {
	f := &fakeDevice{config: hostA}
	s, typ := testServer(t, f)
	apply(t, s, typ, resourceValue(t, typ, "r", hostA), resourceValue(t, typ, "r", hostA))
	if len(f.patches)+len(f.loads)+f.commits != 0 || f.reads != 1 {
		t.Fatalf("patches %v, loads %v, commits %d, reads %d", f.patches, f.loads, f.commits, f.reads)
	}
}

func TestApplyDelete(t *testing.T) {
	f := &fakeDevice{config: hostA}
	s, typ := testServer(t, f)
	state := apply(t, s, typ, resourceValue(t, typ, "r", hostA), tftypes.NewValue(typ, nil))
	if len(f.patches) != 1 || !strings.Contains(f.patches[0], `<host-name nc:operation="delete"/>`) || f.commits != 1 {
		t.Fatalf("patches %v, commits %d", f.patches, f.commits)
	}
	if !state.IsNull() {
		t.Fatalf("state %s, want null", state)
	}
}

// Read returns the device's configuration, with entries ordered as in the
// prior state, the prior resource_name, and nothing the schema does not model.
func TestReadResource(t *testing.T) {
	rf := func(addr string) string {
		return `<route-filter><address>` + addr + `</address><choice-ident>exact</choice-ident><choice-value/></route-filter>`
	}
	policy := func(filters ...string) string {
		return `<policy-options><policy-statement><name>p</name><term><name>t</name><from>` +
			strings.Join(filters, "") + `</from></term></policy-statement></policy-options>`
	}
	f := &fakeDevice{config: `<chassis><alarm/></chassis>` + hostA + policy(rf("192.0.2.0/24"), rf("198.51.100.0/24"))}
	s, typ := testServer(t, f)
	prior := resourceValue(t, typ, "r", policy(rf("198.51.100.0/24"), rf("192.0.2.0/24")))

	resp, _ := s.ReadResource(context.Background(), &tfprotov6.ReadResourceRequest{CurrentState: dv(t, typ, prior)})
	noDiags(t, resp.Diagnostics)
	want := resourceValue(t, typ, "r", hostA+policy(rf("198.51.100.0/24"), rf("192.0.2.0/24")))
	if got := decode(t, typ, resp.NewState); !got.Equal(want) {
		t.Fatalf("state %s, want %s", got, want)
	}
}

// Stored state from an earlier schema: attributes it no longer has are
// dropped, attributes it did not have are null.
func TestUpgradeResourceState(t *testing.T) {
	s, typ := testServer(t, &fakeDevice{})
	resp, _ := s.UpgradeResourceState(context.Background(), &tfprotov6.UpgradeResourceStateRequest{
		RawState: &tfprotov6.RawState{JSON: []byte(`{"resource_name": "r", "gone": "x", "system": [{"host_name": "a", "old_leaf": "y"}]}`)},
	})
	noDiags(t, resp.Diagnostics)
	want := resourceValue(t, typ, "r", hostA)
	if got := decode(t, typ, resp.UpgradedState); !got.Equal(want) {
		t.Fatalf("state %s, want %s", got, want)
	}
}

// StopProvider closes the NETCONF session.
func TestStopProviderClosesClient(t *testing.T) {
	f := &fakeDevice{}
	s, _ := testServer(t, f)
	if _, err := s.StopProvider(context.Background(), &tfprotov6.StopProviderRequest{}); err != nil {
		t.Fatal(err)
	}
	if !f.closed {
		t.Fatal("client not closed")
	}
}
