package generic

import (
	"testing"

	"terraform_provider/patch"
)

func TestXMLToPathMap_SimpleLeaf(t *testing.T) {
	xml := `<configuration><interfaces><interface><name>ge-0/0/0</name></interface></interfaces></configuration>`
	idx := map[string]*patch.NodeInfo{}
	result, err := XMLToPathMap([]byte(xml), idx)
	if err != nil {
		t.Fatal(err)
	}
	if v, ok := result["interfaces/interface/name"]; !ok || v != "ge-0/0/0" {
		t.Errorf("expected interfaces/interface/name=ge-0/0/0, got %v", result)
	}
}

func TestXMLToPathMap_EmptyElement(t *testing.T) {
	xml := `<configuration><system><syslog><disable/></syslog></system></configuration>`
	idx := map[string]*patch.NodeInfo{}
	result, err := XMLToPathMap([]byte(xml), idx)
	if err != nil {
		t.Fatal(err)
	}
	if v, ok := result["system/syslog/disable"]; !ok || v != "" {
		t.Errorf("expected system/syslog/disable='', got %v", result)
	}
}

func TestXMLToPathMap_MissingElement(t *testing.T) {
	xml := `<configuration><system></system></configuration>`
	idx := map[string]*patch.NodeInfo{}
	result, err := XMLToPathMap([]byte(xml), idx)
	if err != nil {
		t.Fatal(err)
	}
	if _, ok := result["system/host-name"]; ok {
		t.Error("missing element should not appear in path map")
	}
}

func TestPathMapToStateAttrs_BasicLeaf(t *testing.T) {
	nodes := []patch.SchemaNode{
		{Name: "interfaces", Type: "container", Children: []patch.SchemaNode{
			{Name: "interface", Type: "list", Children: []patch.SchemaNode{
				{Name: "name", Type: "leaf"},
			}},
		}},
	}
	pathMap := map[string]string{
		"interfaces/interface/name": "ge-0/0/0",
	}
	result := PathMapToStateAttrs(pathMap, nodes)
	ifaces, ok := result["interfaces"]
	if !ok {
		t.Fatal("expected interfaces in result")
	}
	list, ok := ifaces.([]interface{})
	if !ok || len(list) == 0 {
		t.Fatal("expected non-empty interfaces list")
	}
	m, ok := list[0].(map[string]interface{})
	if !ok {
		t.Fatal("expected map inside interfaces")
	}
	ifaceList, ok := m["interface"].([]interface{})
	if !ok || len(ifaceList) == 0 {
		t.Fatal("expected interface list")
	}
	ifaceMap, ok := ifaceList[0].(map[string]interface{})
	if !ok {
		t.Fatal("expected map inside interface")
	}
	if ifaceMap["name"] != "ge-0/0/0" {
		t.Errorf("expected name=ge-0/0/0, got %v", ifaceMap["name"])
	}
}

func TestPathMapToStateAttrs_MissingLeaf(t *testing.T) {
	nodes := []patch.SchemaNode{
		{Name: "host-name", Type: "leaf"},
	}
	pathMap := map[string]string{}
	result := PathMapToStateAttrs(pathMap, nodes)
	if _, ok := result["host_name"]; ok {
		t.Error("missing leaf should not appear in state")
	}
}
