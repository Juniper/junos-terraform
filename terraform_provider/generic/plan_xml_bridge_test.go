package generic

import (
	"strings"
	"testing"

	"terraform_provider/patch"
)

func TestPlanToXML_SimpleLeaf(t *testing.T) {
	idx := map[string]*patch.NodeInfo{
		"interfaces/interface/name": {Path: "interfaces/interface/name", Kind: patch.KindLeaf},
	}
	values := map[string]string{
		"interfaces/interface/name": "ge-0/0/0",
	}
	got, err := PlanToXML(values, idx)
	if err != nil {
		t.Fatal(err)
	}
	xml := string(got)
	if !strings.Contains(xml, "<name>ge-0/0/0</name>") {
		t.Errorf("expected <name>ge-0/0/0</name>, got:\n%s", xml)
	}
	if !strings.Contains(xml, "<interface>") {
		t.Errorf("expected <interface> wrapper, got:\n%s", xml)
	}
}

func TestPlanToXML_EmptyValues(t *testing.T) {
	idx := map[string]*patch.NodeInfo{}
	values := map[string]string{}
	got, err := PlanToXML(values, idx)
	if err != nil {
		t.Fatal(err)
	}
	xml := string(got)
	if !strings.Contains(xml, "<configuration>") {
		t.Errorf("expected <configuration> root, got:\n%s", xml)
	}
}

func TestStateToPathMap_BasicLeaf(t *testing.T) {
	nodes := []patch.SchemaNode{
		{Name: "interfaces", Type: "container", Children: []patch.SchemaNode{
			{Name: "interface", Type: "list", Children: []patch.SchemaNode{
				{Name: "name", Type: "leaf"},
			}},
		}},
	}
	attrs := map[string]interface{}{
		"interfaces": []interface{}{
			map[string]interface{}{
				"interface": []interface{}{
					map[string]interface{}{
						"name": "ge-0/0/0",
					},
				},
			},
		},
	}
	result := StateToPathMap(attrs, nodes)
	if v, ok := result["interfaces/interface/name"]; !ok || v != "ge-0/0/0" {
		t.Errorf("expected interfaces/interface/name=ge-0/0/0, got %v", result)
	}
}

func TestStateToPathMap_NullLeafOmitted(t *testing.T) {
	nodes := []patch.SchemaNode{
		{Name: "host-name", Type: "leaf"},
	}
	attrs := map[string]interface{}{
		"host_name": "",
	}
	result := StateToPathMap(attrs, nodes)
	if _, ok := result["host-name"]; ok {
		t.Error("empty string should be omitted from path map")
	}
}
