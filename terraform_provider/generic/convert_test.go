package generic

import (
	"strings"
	"testing"

	"terraform_provider/patch"

	"github.com/hashicorp/terraform-plugin-go/tftypes"
)

// convertSchema is shaped like an SRX configuration: presence containers and
// empty leaves, keyed lists nested in lists, a compound key, a leaf-list.
const convertSchema = `{"root": {"children": [{"name": "configuration", "type": "container", "children": [
  {"name": "version", "type": "leaf"},
  {"name": "system", "type": "container", "children": [
    {"name": "host-name", "type": "leaf"},
    {"name": "services", "type": "container", "children": [
      {"name": "netconf", "type": "container", "children": [
        {"name": "ssh", "type": "container", "presence": "enable ssh", "children": [{"name": "port", "type": "leaf"}]},
        {"name": "rfc-compliant", "type": "leaf", "leaf-type": "empty"}
      ]}
    ]}
  ]},
  {"name": "policy-options", "type": "container", "children": [
    {"name": "policy-statement", "type": "list", "key": "name", "children": [
      {"name": "name", "type": "leaf"},
      {"name": "term", "type": "list", "key": "name", "children": [
        {"name": "name", "type": "leaf"},
        {"name": "from", "type": "container", "children": [
          {"name": "family", "type": "leaf"},
          {"name": "route-filter", "type": "list", "key": "address choice-ident choice-value", "children": [
            {"name": "address", "type": "leaf"},
            {"name": "choice-ident", "type": "leaf"},
            {"name": "choice-value", "type": "leaf"}
          ]}
        ]},
        {"name": "then", "type": "container", "children": [{"name": "accept", "type": "leaf", "leaf-type": "empty"}]}
      ]}
    ]}
  ]},
  {"name": "protocols", "type": "container", "children": [
    {"name": "bgp", "type": "container", "children": [
      {"name": "group", "type": "list", "key": "name", "children": [
        {"name": "name", "type": "leaf"},
        {"name": "import", "type": "leaf-list"}
      ]}
    ]}
  ]}
]}]}}`

func convertFixture(t *testing.T) ([]patch.SchemaNode, tftypes.Object) {
	t.Helper()
	ResetSchema()
	t.Cleanup(ResetSchema)
	_, nodes, err := LoadSchema([]byte(convertSchema))
	if err != nil {
		t.Fatal(err)
	}
	_, typ := BuildSchema(nodes[0].Children)
	return nodes[0].Children, typ
}

func canonical(t *testing.T, xmlText string) (tftypes.Value, string) {
	t.Helper()
	nodes, typ := convertFixture(t)
	tree, err := patch.BuildTree([]byte(xmlText))
	if err != nil {
		t.Fatal(err)
	}
	v, err := ConfigToValue(tree, nodes, typ)
	if err != nil {
		t.Fatal(err)
	}
	out, err := ValueToConfig(v, nodes)
	if err != nil {
		t.Fatal(err)
	}
	b, err := patch.MarshalTree(out)
	if err != nil {
		t.Fatal(err)
	}
	return v, compactXML(string(b))
}

// Router output: Junos attributes, an element the schema does not model,
// whitespace in empty elements, keys after other leaves, several entries.
const deviceXML = `<configuration junos:changed-seconds="1" xmlns:junos="http://xml.juniper.net/junos/26.2R1.7/junos">
  <system>
    <services><netconf><ssh>
    </ssh><rfc-compliant/></netconf></services>
    <host-name>router</host-name>
    <not-modeled>x</not-modeled>
  </system>
  <policy-options>
    <policy-statement>
      <name>import</name>
      <term>
        <name>v6</name>
        <from>
          <family>inet6</family>
          <route-filter><address>2001:db8:1::/48</address><choice-ident>prefix-length-range</choice-ident><choice-value>/64-/64</choice-value></route-filter>
          <route-filter><address>2001:db8:2::/48</address><choice-ident>prefix-length-range</choice-ident><choice-value>/128-/128</choice-value></route-filter>
          <route-filter><address>::/0</address><choice-ident>exact</choice-ident><choice-value></choice-value></route-filter>
        </from>
        <then><accept/></then>
      </term>
      <term><then><accept/></then><name>v4</name></term>
    </policy-statement>
  </policy-options>
  <protocols><bgp><group><import>a</import><name>g</name><import>b</import></group></bgp></protocols>
</configuration>`

func TestConfigToValueAndBack(t *testing.T) {
	_, got := canonical(t, deviceXML)
	want := compactXML(`<configuration>
  <system>
    <host-name>router</host-name>
    <services><netconf><ssh/><rfc-compliant/></netconf></services>
  </system>
  <policy-options>
    <policy-statement>
      <name>import</name>
      <term>
        <name>v6</name>
        <from>
          <family>inet6</family>
          <route-filter><address>2001:db8:1::/48</address><choice-ident>prefix-length-range</choice-ident><choice-value>/64-/64</choice-value></route-filter>
          <route-filter><address>2001:db8:2::/48</address><choice-ident>prefix-length-range</choice-ident><choice-value>/128-/128</choice-value></route-filter>
          <route-filter><address>::/0</address><choice-ident>exact</choice-ident><choice-value/></route-filter>
        </from>
        <then><accept/></then>
      </term>
      <term><name>v4</name><then><accept/></then></term>
    </policy-statement>
  </policy-options>
  <protocols><bgp><group><name>g</name><import>a</import><import>b</import></group></bgp></protocols>
</configuration>`)
	if got != want {
		t.Fatalf("got\n%s\nwant\n%s", got, want)
	}
}

func TestConfigToValueRoundTrip(t *testing.T) {
	v, x := canonical(t, deviceXML)
	v2, x2 := canonical(t, x)
	if x != x2 {
		t.Fatalf("XML changed on a second round trip:\n%s\n%s", x, x2)
	}
	if !v.Equal(v2) {
		t.Fatalf("value changed on a second round trip:\n%s\n%s", v, v2)
	}
}

func TestConfigToValueNullAndEmpty(t *testing.T) {
	v, _ := canonical(t, `<configuration><system><services><netconf><rfc-compliant/></netconf></services></system></configuration>`)
	path := func(steps ...interface{}) tftypes.Value {
		p := tftypes.NewAttributePath()
		for _, s := range steps {
			switch s := s.(type) {
			case string:
				p = p.WithAttributeName(s)
			case int:
				p = p.WithElementKeyInt(s)
			}
		}
		got, _, err := tftypes.WalkAttributePath(v, p)
		if err != nil {
			t.Fatalf("%v: %v", p, err)
		}
		return got.(tftypes.Value)
	}
	netconf := func(attr string) tftypes.Value {
		return path("system", 0, "services", 0, "netconf", 0, attr)
	}
	var s string
	if err := netconf("rfc_compliant").As(&s); err != nil || s != "" {
		t.Fatalf("empty leaf: %q %v", s, err)
	}
	if !netconf("ssh").IsNull() {
		t.Fatal("missing container is not null")
	}
	if !path("system", 0, "host_name").IsNull() || !path("policy_options").IsNull() || !path("resource_name").IsNull() {
		t.Fatal("missing leaf, section or resource_name is not null")
	}
	if !path("version").IsNull() {
		t.Fatal("version is not null")
	}
}

func compactXML(s string) string {
	var b strings.Builder
	for _, line := range strings.Split(s, "\n") {
		b.WriteString(strings.TrimSpace(line))
	}
	return b.String()
}
