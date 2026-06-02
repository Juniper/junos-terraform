package generic

import (
	"context"
	"strings"
	"testing"

	"terraform_provider/patch"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// testSchemaJSON is a small representative YANG-like schema for testing.
const testSchemaJSON = `{
  "root": {
    "children": [
      {
        "name": "interfaces",
        "type": "list",
        "key": "name",
        "path": "",
        "children": [
          {
            "name": "name",
            "type": "leaf",
            "leaf-type": "string",
            "path": "interfaces"
          },
          {
            "name": "description",
            "type": "leaf",
            "leaf-type": "string",
            "path": "interfaces"
          },
          {
            "name": "mtu",
            "type": "leaf",
            "leaf-type": "string",
            "path": "interfaces"
          },
          {
            "name": "unit",
            "type": "list",
            "key": "name",
            "path": "interfaces",
            "children": [
              {
                "name": "name",
                "type": "leaf",
                "leaf-type": "string",
                "path": "interfaces/unit"
              },
              {
                "name": "description",
                "type": "leaf",
                "leaf-type": "string",
                "path": "interfaces/unit"
              },
              {
                "name": "family",
                "type": "container",
                "path": "interfaces/unit",
                "children": [
                  {
                    "name": "inet",
                    "type": "container",
                    "path": "interfaces/unit/family",
                    "children": [
                      {
                        "name": "address",
                        "type": "list",
                        "key": "name",
                        "path": "interfaces/unit/family/inet",
                        "children": [
                          {
                            "name": "name",
                            "type": "leaf",
                            "leaf-type": "string",
                            "path": "interfaces/unit/family/inet/address"
                          }
                        ]
                      }
                    ]
                  }
                ]
              }
            ]
          }
        ]
      },
      {
        "name": "system",
        "type": "container",
        "path": "",
        "children": [
          {
            "name": "host-name",
            "type": "leaf",
            "leaf-type": "string",
            "path": "system"
          },
          {
            "name": "name-server",
            "type": "leaf-list",
            "leaf-type": "string",
            "path": "system"
          }
        ]
      },
      {
        "name": "protocols",
        "type": "container",
        "path": "",
        "children": [
          {
            "name": "lldp",
            "type": "container",
            "path": "protocols",
            "children": []
          }
        ]
      }
    ]
  }
}`

func TestLoadSchema(t *testing.T) {
	idx, err := LoadSchema([]byte(testSchemaJSON))
	if err != nil {
		t.Fatalf("LoadSchema() error: %v", err)
	}
	if idx == nil {
		t.Fatal("LoadSchema() returned nil")
	}
	if len(idx.TopLevel) != 3 {
		t.Fatalf("expected 3 top-level nodes, got %d", len(idx.TopLevel))
	}
	if idx.TopLevel[0].Name != "interfaces" {
		t.Fatalf("expected first top-level to be 'interfaces', got %q", idx.TopLevel[0].Name)
	}
	if idx.TopLevel[1].Name != "system" {
		t.Fatalf("expected second top-level to be 'system', got %q", idx.TopLevel[1].Name)
	}

	// Verify flat index has entries
	if len(idx.ByPath) == 0 {
		t.Fatal("ByPath index is empty")
	}
}

func TestBuildTerraformSchema(t *testing.T) {
	idx, err := LoadSchema([]byte(testSchemaJSON))
	if err != nil {
		t.Fatalf("LoadSchema() error: %v", err)
	}

	s := BuildTerraformSchema(idx)

	// Should have resource_name + 3 top-level attrs
	if len(s.Attributes) != 4 {
		t.Fatalf("expected 4 attributes, got %d", len(s.Attributes))
	}

	// resource_name should be required
	rn, ok := s.Attributes["resource_name"]
	if !ok {
		t.Fatal("missing resource_name attribute")
	}
	if !rn.(schema.StringAttribute).Required {
		t.Fatal("resource_name should be required")
	}

	// interfaces should be a ListNestedAttribute
	ifaces, ok := s.Attributes["interfaces"]
	if !ok {
		t.Fatal("missing interfaces attribute")
	}
	if _, ok := ifaces.(schema.ListNestedAttribute); !ok {
		t.Fatalf("interfaces should be ListNestedAttribute, got %T", ifaces)
	}

	// system should be a ListNestedAttribute (containers are lists with max 1 element)
	sys, ok := s.Attributes["system"]
	if !ok {
		t.Fatal("missing system attribute")
	}
	if _, ok := sys.(schema.ListNestedAttribute); !ok {
		t.Fatalf("system should be ListNestedAttribute, got %T", sys)
	}

	// protocols.lldp is an empty container — should be a ListNestedAttribute
	protocols, ok := s.Attributes["protocols"]
	if !ok {
		t.Fatal("missing protocols attribute")
	}
	protoNested, ok := protocols.(schema.ListNestedAttribute)
	if !ok {
		t.Fatalf("protocols should be ListNestedAttribute, got %T", protocols)
	}
	lldp, ok := protoNested.NestedObject.Attributes["lldp"]
	if !ok {
		t.Fatal("missing protocols.lldp attribute")
	}
	if _, ok := lldp.(schema.ListNestedAttribute); !ok {
		t.Fatalf("protocols.lldp (empty container) should be ListNestedAttribute, got %T", lldp)
	}
}

func TestModelToXMLBytes_SimpleLeaf(t *testing.T) {
	idx, err := LoadSchema([]byte(testSchemaJSON))
	if err != nil {
		t.Fatalf("LoadSchema() error: %v", err)
	}

	ctx := context.Background()
	var diags diag.Diagnostics

	// Build attrs for system container with host-name (containers are single-element lists)
	systemAttrs := map[string]attr.Value{
		"host_name":   types.StringValue("router1"),
		"name_server": types.ListNull(types.StringType),
	}
	systemAttrTypes := map[string]attr.Type{
		"host_name":   types.StringType,
		"name_server": types.ListType{ElemType: types.StringType},
	}
	systemObj, d := types.ObjectValue(systemAttrTypes, systemAttrs)
	diags.Append(d...)
	if diags.HasError() {
		t.Fatalf("failed to build system object: %v", diags)
	}
	systemList, d := types.ListValue(types.ObjectType{AttrTypes: systemAttrTypes}, []attr.Value{systemObj})
	diags.Append(d...)
	if diags.HasError() {
		t.Fatalf("failed to build system list: %v", diags)
	}

	protocolsObjType := types.ObjectType{AttrTypes: containerAttrTypes(idx.TopLevel[2])}
	attrs := map[string]attr.Value{
		"interfaces": types.ListNull(types.ObjectType{AttrTypes: containerAttrTypes(idx.TopLevel[0])}),
		"system":     systemList,
		"protocols":  types.ListNull(protocolsObjType),
	}

	xmlBytes := ModelToXMLBytes(ctx, attrs, idx, &diags)
	if diags.HasError() {
		t.Fatalf("ModelToXMLBytes() errors: %v", diags)
	}

	xmlStr := string(xmlBytes)
	if !strings.Contains(xmlStr, "<host-name>router1</host-name>") {
		t.Fatalf("expected XML to contain host-name element, got:\n%s", xmlStr)
	}
	if !strings.Contains(xmlStr, "<system>") {
		t.Fatalf("expected XML to contain <system> element, got:\n%s", xmlStr)
	}
}

func TestModelToXMLBytes_LeafList(t *testing.T) {
	idx, err := LoadSchema([]byte(testSchemaJSON))
	if err != nil {
		t.Fatalf("LoadSchema() error: %v", err)
	}

	ctx := context.Background()
	var diags diag.Diagnostics

	nameServers, _ := types.ListValue(types.StringType, []attr.Value{
		types.StringValue("8.8.8.8"),
		types.StringValue("8.8.4.4"),
	})

	systemAttrs := map[string]attr.Value{
		"host_name":   types.StringValue("router1"),
		"name_server": nameServers,
	}
	systemAttrTypes := map[string]attr.Type{
		"host_name":   types.StringType,
		"name_server": types.ListType{ElemType: types.StringType},
	}
	systemObj, _ := types.ObjectValue(systemAttrTypes, systemAttrs)
	systemList, _ := types.ListValue(types.ObjectType{AttrTypes: systemAttrTypes}, []attr.Value{systemObj})

	protocolsObjType := types.ObjectType{AttrTypes: containerAttrTypes(idx.TopLevel[2])}
	attrs := map[string]attr.Value{
		"interfaces": types.ListNull(types.ObjectType{AttrTypes: containerAttrTypes(idx.TopLevel[0])}),
		"system":     systemList,
		"protocols":  types.ListNull(protocolsObjType),
	}

	xmlBytes := ModelToXMLBytes(ctx, attrs, idx, &diags)
	if diags.HasError() {
		t.Fatalf("ModelToXMLBytes() errors: %v", diags)
	}

	xmlStr := string(xmlBytes)
	if !strings.Contains(xmlStr, "<name-server>8.8.8.8</name-server>") {
		t.Fatalf("expected XML to contain first name-server, got:\n%s", xmlStr)
	}
	if !strings.Contains(xmlStr, "<name-server>8.8.4.4</name-server>") {
		t.Fatalf("expected XML to contain second name-server, got:\n%s", xmlStr)
	}
}

func TestXMLBytesToModel_SimpleLeaf(t *testing.T) {
	idx, err := LoadSchema([]byte(testSchemaJSON))
	if err != nil {
		t.Fatalf("LoadSchema() error: %v", err)
	}

	ctx := context.Background()
	var diags diag.Diagnostics

	xmlInput := []byte(`<?xml version="1.0" encoding="UTF-8"?>
<configuration>
  <system>
    <host-name>router1</host-name>
    <name-server>8.8.8.8</name-server>
    <name-server>1.1.1.1</name-server>
  </system>
</configuration>`)

	result := XMLBytesToModel(ctx, xmlInput, idx, &diags)
	if diags.HasError() {
		t.Fatalf("XMLBytesToModel() errors: %v", diags)
	}

	systemVal, ok := result["system"]
	if !ok {
		t.Fatal("missing 'system' in result")
	}
	if systemVal.IsNull() {
		t.Fatal("system value is null")
	}

	sysList, ok := systemVal.(types.List)
	if !ok {
		t.Fatalf("system should be List, got %T", systemVal)
	}
	if len(sysList.Elements()) != 1 {
		t.Fatalf("system list should have 1 element, got %d", len(sysList.Elements()))
	}

	sysObj, ok := sysList.Elements()[0].(types.Object)
	if !ok {
		t.Fatalf("system[0] should be Object, got %T", sysList.Elements()[0])
	}

	hostName := sysObj.Attributes()["host_name"]
	if hostName.(types.String).ValueString() != "router1" {
		t.Fatalf("expected host_name='router1', got %q", hostName.(types.String).ValueString())
	}

	nameServerVal := sysObj.Attributes()["name_server"]
	if nameServerVal.IsNull() {
		t.Fatal("name_server should not be null")
	}
	nsList, ok := nameServerVal.(types.List)
	if !ok {
		t.Fatalf("name_server should be List, got %T", nameServerVal)
	}
	if len(nsList.Elements()) != 2 {
		t.Fatalf("expected 2 name-servers, got %d", len(nsList.Elements()))
	}
}

func TestXMLBytesToModel_List(t *testing.T) {
	idx, err := LoadSchema([]byte(testSchemaJSON))
	if err != nil {
		t.Fatalf("LoadSchema() error: %v", err)
	}

	ctx := context.Background()
	var diags diag.Diagnostics

	xmlInput := []byte(`<?xml version="1.0" encoding="UTF-8"?>
<configuration>
  <interfaces>
    <name>ge-0/0/0</name>
    <description>uplink</description>
    <mtu>9000</mtu>
  </interfaces>
  <interfaces>
    <name>ge-0/0/1</name>
    <description>downlink</description>
  </interfaces>
</configuration>`)

	result := XMLBytesToModel(ctx, xmlInput, idx, &diags)
	if diags.HasError() {
		t.Fatalf("XMLBytesToModel() errors: %v", diags)
	}

	ifacesVal, ok := result["interfaces"]
	if !ok {
		t.Fatal("missing 'interfaces' in result")
	}
	ifacesList, ok := ifacesVal.(types.List)
	if !ok {
		t.Fatalf("interfaces should be List, got %T", ifacesVal)
	}
	if len(ifacesList.Elements()) != 2 {
		t.Fatalf("expected 2 interfaces, got %d", len(ifacesList.Elements()))
	}

	// Check first interface
	first := ifacesList.Elements()[0].(types.Object)
	name := first.Attributes()["name"].(types.String)
	if name.ValueString() != "ge-0/0/0" {
		t.Fatalf("expected first interface name='ge-0/0/0', got %q", name.ValueString())
	}
	mtu := first.Attributes()["mtu"].(types.String)
	if mtu.ValueString() != "9000" {
		t.Fatalf("expected first interface mtu='9000', got %q", mtu.ValueString())
	}
}

func TestSerializeTree_RoundTrip(t *testing.T) {
	xmlInput := []byte(`<?xml version="1.0" encoding="UTF-8"?>
<configuration>
  <system>
    <host-name>router1</host-name>
  </system>
</configuration>`)

	tree, err := patch.BuildTree(xmlInput)
	if err != nil {
		t.Fatalf("BuildTree() error: %v", err)
	}

	serialized, err := SerializeTree(tree)
	if err != nil {
		t.Fatalf("SerializeTree() error: %v", err)
	}

	// Re-parse the serialized output
	tree2, err := patch.BuildTree(serialized)
	if err != nil {
		t.Fatalf("BuildTree(serialized) error: %v", err)
	}

	if tree2.Tag != "configuration" {
		t.Fatalf("expected root tag 'configuration', got %q", tree2.Tag)
	}
	if len(tree2.Children) != 1 {
		t.Fatalf("expected 1 child, got %d", len(tree2.Children))
	}
	if tree2.Children[0].Tag != "system" {
		t.Fatalf("expected child tag 'system', got %q", tree2.Children[0].Tag)
	}
	if tree2.Children[0].Children[0].Tag != "host-name" {
		t.Fatalf("expected grandchild tag 'host-name', got %q", tree2.Children[0].Children[0].Tag)
	}
	if tree2.Children[0].Children[0].Text != "router1" {
		t.Fatalf("expected text 'router1', got %q", tree2.Children[0].Children[0].Text)
	}
}

func TestModelToXML_RoundTrip(t *testing.T) {
	idx, err := LoadSchema([]byte(testSchemaJSON))
	if err != nil {
		t.Fatalf("LoadSchema() error: %v", err)
	}

	ctx := context.Background()
	var diags diag.Diagnostics

	// Start with XML, convert to model, then back to XML
	xmlInput := []byte(`<?xml version="1.0" encoding="UTF-8"?>
<configuration>
  <system>
    <host-name>router1</host-name>
    <name-server>8.8.8.8</name-server>
  </system>
</configuration>`)

	// XML → Model
	modelAttrs := XMLBytesToModel(ctx, xmlInput, idx, &diags)
	if diags.HasError() {
		t.Fatalf("XMLBytesToModel() errors: %v", diags)
	}

	// Model → XML
	xmlOutput := ModelToXMLBytes(ctx, modelAttrs, idx, &diags)
	if diags.HasError() {
		t.Fatalf("ModelToXMLBytes() errors: %v", diags)
	}

	// Verify the output contains the expected elements
	xmlStr := string(xmlOutput)
	if !strings.Contains(xmlStr, "<host-name>router1</host-name>") {
		t.Fatalf("round-trip lost host-name, got:\n%s", xmlStr)
	}
	if !strings.Contains(xmlStr, "<name-server>8.8.8.8</name-server>") {
		t.Fatalf("round-trip lost name-server, got:\n%s", xmlStr)
	}
}

func TestNormalizeUnknowns_Recursive(t *testing.T) {
	addressObjType := types.ObjectType{AttrTypes: map[string]attr.Type{
		"name": types.StringType,
	}}
	inetObjType := types.ObjectType{AttrTypes: map[string]attr.Type{
		"address": types.ListType{ElemType: addressObjType},
	}}
	familyObjType := types.ObjectType{AttrTypes: map[string]attr.Type{
		"inet": types.ListType{ElemType: inetObjType},
	}}
	unitObjType := types.ObjectType{AttrTypes: map[string]attr.Type{
		"name":        types.StringType,
		"description": types.StringType,
		"family":      types.ListType{ElemType: familyObjType},
	}}
	ifaceObjType := types.ObjectType{AttrTypes: map[string]attr.Type{
		"name":        types.StringType,
		"description": types.StringType,
		"unit":        types.ListType{ElemType: unitObjType},
	}}

	unitObj, diags := types.ObjectValue(unitObjType.AttrTypes, map[string]attr.Value{
		"name":        types.StringValue("0"),
		"description": types.StringUnknown(),
		"family":      types.ListUnknown(familyObjType),
	})
	if diags.HasError() {
		t.Fatalf("failed to build unit object: %v", diags)
	}

	unitList, diags := types.ListValue(unitObjType, []attr.Value{unitObj})
	if diags.HasError() {
		t.Fatalf("failed to build unit list: %v", diags)
	}

	ifaceObj, diags := types.ObjectValue(ifaceObjType.AttrTypes, map[string]attr.Value{
		"name":        types.StringValue("ge-0/0/0"),
		"description": types.StringUnknown(),
		"unit":        unitList,
	})
	if diags.HasError() {
		t.Fatalf("failed to build interface object: %v", diags)
	}

	ifaces, diags := types.ListValue(ifaceObjType, []attr.Value{ifaceObj})
	if diags.HasError() {
		t.Fatalf("failed to build interfaces list: %v", diags)
	}

	normalized := normalizeUnknowns(ifaces)
	listVal, ok := normalized.(types.List)
	if !ok {
		t.Fatalf("expected normalized list, got %T", normalized)
	}

	iface := listVal.Elements()[0].(types.Object)
	if !iface.Attributes()["description"].IsNull() || iface.Attributes()["description"].IsUnknown() {
		t.Fatal("top-level unknown leaf should be converted to null")
	}

	unit := iface.Attributes()["unit"].(types.List).Elements()[0].(types.Object)
	if !unit.Attributes()["description"].IsNull() || unit.Attributes()["description"].IsUnknown() {
		t.Fatal("nested unknown leaf should be converted to null")
	}

	family := unit.Attributes()["family"]
	if !family.IsNull() || family.IsUnknown() {
		t.Fatal("unknown nested list should be converted to typed null")
	}
}

func TestNormalizeName(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"host-name", "host_name"},
		{"name-server", "name_server"},
		{"ge-0/0/0", "ge_0/0/0"},
		{"vlan.id", "vlan_id"},
		{"simple", "simple"},
	}
	for _, tc := range tests {
		got := normalizeName(tc.input)
		if got != tc.expected {
			t.Errorf("normalizeName(%q) = %q, want %q", tc.input, got, tc.expected)
		}
	}
}

func TestNewConfigResource_ValidSchema(t *testing.T) {
	factory := NewConfigResource(testSchemaJSON)
	r := factory()
	if r == nil {
		t.Fatal("NewConfigResource() returned nil resource")
	}

	cr, ok := r.(*ConfigResource)
	if !ok {
		t.Fatalf("expected *ConfigResource, got %T", r)
	}
	if cr.idx == nil {
		t.Fatal("ConfigResource.idx is nil")
	}
	if len(cr.tfSchema.Attributes) == 0 {
		t.Fatal("ConfigResource.tfSchema has no attributes")
	}
}

func TestNewConfigResource_InvalidSchema(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Fatal("expected panic on invalid schema")
		}
	}()
	factory := NewConfigResource("{invalid json")
	factory()
}
