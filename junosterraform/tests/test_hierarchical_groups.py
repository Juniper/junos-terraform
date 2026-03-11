import importlib.util
from importlib.machinery import SourceFileLoader
import json
import pathlib
import sys

import pytest
import yaml


REPO_ROOT = pathlib.Path(__file__).resolve().parents[2]


def _load_script_module(script_name: str, module_name: str):
    script_path = REPO_ROOT / "junosterraform" / script_name
    loader = SourceFileLoader(module_name, str(script_path))
    spec = importlib.util.spec_from_loader(module_name, loader)
    module = importlib.util.module_from_spec(spec)
    assert spec is not None and spec.loader is not None
    spec.loader.exec_module(module)
    return module


@pytest.fixture(scope="module")
def xml2yaml_mod():
    return _load_script_module("jtaf-xml2yaml", "jtaf_xml2yaml_mod")


@pytest.fixture(scope="module")
def ansible_mod():
    return _load_script_module("jtaf-ansible", "jtaf_ansible_mod")


def test_detect_device_type(xml2yaml_mod):
    assert xml2yaml_mod.detect_device_type({"system": {"product_name": "QFX5110"}}) == "qfx"
    assert xml2yaml_mod.detect_device_type({"chassis": {"product_name": "SRX300"}}) == "srx"
    assert xml2yaml_mod.detect_device_type({"system": {"product_name": "Unknown"}}) is None


def test_prepare_and_intersection_helpers(xml2yaml_mod):
    assert xml2yaml_mod.prepare_tag("host-name.v4") == "host_name_v4"

    shared = xml2yaml_mod.build_shared_payload([
        {"system": {"services": {"ssh": True}, "domain": "a"}},
        {"system": {"services": {"ssh": True}, "domain": "b"}},
    ])
    assert shared == {"system": {"services": {"ssh": True}}}

    delta = xml2yaml_mod.subtract_common(
        {"system": {"services": {"ssh": True}, "domain": "a"}},
        shared,
    )
    assert delta == {"system": {"domain": "a"}}


def test_extract_hierarchy_groups(xml2yaml_mod):
    groups = xml2yaml_mod.extract_hierarchy_groups([
        ("h1", {"a": 1}, "qfx"),
        ("h2", {"a": 2}, "qfx"),
        ("h3", {"a": 3}, None),
    ])
    assert sorted(groups.keys()) == ["all", "device_type:qfx"]
    assert len(groups["all"]) == 3
    assert len(groups["device_type:qfx"]) == 2


def test_yaml_writers(xml2yaml_mod, tmp_path):
    xml2yaml_mod.write_group_vars_flat(str(tmp_path), "all", {"system": {"x": 1}})
    xml2yaml_mod.write_group_vars_flat(str(tmp_path), "device_type:qfx", {"chassis": {"y": 2}})
    xml2yaml_mod.write_host_vars_flat(
        str(tmp_path),
        "leaf1",
        {"system": {"host_name": "leaf1", "services": {"ssh": True}}},
        {"system": {"services": {"ssh": True}}},
    )

    with open(tmp_path / "group_vars" / "all.yml") as f:
        assert yaml.safe_load(f) == {"system": {"x": 1}}
    with open(tmp_path / "group_vars" / "qfx" / "all.yml") as f:
        assert yaml.safe_load(f) == {"chassis": {"y": 2}}
    with open(tmp_path / "host_vars" / "leaf1.yaml") as f:
        assert yaml.safe_load(f) == {"system": {"host_name": "leaf1"}}


def test_parse_xml_to_payload(xml2yaml_mod, tmp_path):
    schema = {
        "root": {
            "name": "root",
            "children": [
                {
                    "name": "configuration",
                    "children": [
                        {"name": "system", "type": "container", "children": [
                            {"name": "host-name", "type": "leaf"},
                            {"name": "product-name", "type": "leaf"},
                        ]}
                    ],
                }
            ]
        }
    }
    xml_file = tmp_path / "router1.xml"
    xml_file.write_text(
        "<configuration><system><host-name>r1</host-name>"
        "<product-name>QFX5100</product-name></system></configuration>"
    )

    hostname, payload, dtype = xml2yaml_mod.parse_xml_to_payload(str(xml_file), schema)
    assert hostname == "r1"
    assert payload["system"]["host_name"] == "r1"
    assert dtype == "qfx"


def test_xml2yaml_main_end_to_end(xml2yaml_mod, tmp_path, monkeypatch):
    schema = {
        "root": {
            "name": "root",
            "children": [
                {
                    "name": "configuration",
                    "children": [
                        {"name": "system", "type": "container", "children": [
                            {"name": "host-name", "type": "leaf"},
                            {"name": "product-name", "type": "leaf"},
                            {"name": "services", "type": "container", "children": [
                                {"name": "ssh", "type": "leaf"},
                            ]},
                        ]}
                    ],
                }
            ]
        }
    }
    schema_file = tmp_path / "trimmed_schema.json"
    schema_file.write_text(json.dumps(schema))

    xml1 = tmp_path / "a.xml"
    xml1.write_text(
        "<configuration><system><host-name>a</host-name><product-name>QFX5100</product-name>"
        "<services><ssh/></services></system></configuration>"
    )
    xml2 = tmp_path / "b.xml"
    xml2.write_text(
        "<configuration><system><host-name>b</host-name><product-name>QFX5100</product-name>"
        "<services><ssh/></services></system></configuration>"
    )

    out_dir = tmp_path / "out"
    monkeypatch.setattr(
        sys,
        "argv",
        [
            "jtaf-xml2yaml",
            "-j",
            str(schema_file),
            "-x",
            str(xml1),
            str(xml2),
            "-d",
            str(out_dir),
        ],
    )
    xml2yaml_mod.main()

    assert (out_dir / "hosts").exists()
    assert (out_dir / "group_vars" / "all.yml").exists()
    assert (out_dir / "group_vars" / "qfx" / "all.yml").exists()
    assert (out_dir / "host_vars" / "a.yaml").exists()
    assert (out_dir / "host_vars" / "b.yaml").exists()
    assert "[device_qfx]" in (out_dir / "hosts").read_text()


def test_xml2yaml_main_missing_configuration_raises(xml2yaml_mod, tmp_path, monkeypatch):
    bad_schema = {"root": {"name": "root", "children": [{"name": "not-configuration"}]}}
    schema_file = tmp_path / "bad.json"
    schema_file.write_text(json.dumps(bad_schema))
    xml1 = tmp_path / "a.xml"
    xml1.write_text("<configuration><system/></configuration>")

    monkeypatch.setattr(
        sys,
        "argv",
        ["jtaf-xml2yaml", "-j", str(schema_file), "-x", str(xml1), "-d", str(tmp_path / "out")],
    )
    with pytest.raises(ValueError):
        xml2yaml_mod.main()


def test_jtaf_ansible_main_generates_role(ansible_mod, tmp_path, monkeypatch):
    schema_file = tmp_path / "schema.json"
    schema_file.write_text(json.dumps({"root": {"children": [{"name": "configuration", "children": []}]}}))
    xml_file = tmp_path / "cfg.xml"
    xml_file.write_text("<configuration><system><host-name>r1</host-name></system></configuration>")

    monkeypatch.chdir(tmp_path)
    monkeypatch.setattr(
        sys,
        "argv",
        [
            "jtaf-ansible",
            "-j",
            str(schema_file),
            "-x",
            str(xml_file),
            "-t",
            "qfx",
        ],
    )
    ansible_mod.main()

    out = tmp_path / "ansible-provider-junos-qfx"
    tasks = out / "roles" / "qfx_role" / "tasks" / "main.yml"
    filters_py = out / "filter_plugins" / "jtaf_filters.py"
    assert tasks.exists()
    assert filters_py.exists()
    task_text = tasks.read_text()
    assert "Merge variables from hierarchy" in task_text
    assert "jtaf_apply_merge_directives" in task_text
    assert (out / "group_vars" / "all.yml").exists()
    assert (out / "trimmed_schema.json").exists()


def test_elem_to_dict_list_leaflist_and_container(xml2yaml_mod):
    resources = {
        "name": "root",
        "children": [
            {
                "name": "configuration",
                "type": "container",
                "children": [
                    {
                        "name": "interfaces",
                        "type": "list",
                        "children": [
                            {"name": "name", "type": "leaf"},
                            {"name": "description", "type": "leaf"},
                        ],
                    },
                    {
                        "name": "vlans",
                        "type": "container",
                        "children": [
                            {"name": "vlan-id", "type": "leaf-list"},
                        ],
                    },
                ],
            }
        ],
    }
    xml = (
        "<configuration>"
        "<interfaces><name>et-0/0/0</name><description>uplink</description></interfaces>"
        "<interfaces><name>et-0/0/1</name><description>downlink</description></interfaces>"
        "<vlans><vlan-id>100</vlan-id></vlans>"
        "</configuration>"
    )
    root = xml2yaml_mod.ElementTree.fromstring(xml)
    payload = xml2yaml_mod.elem_to_dict(root, "", resources)

    assert payload["interfaces"][0]["name"] == "et-0/0/0"
    assert payload["interfaces"][1]["description"] == "downlink"
    assert payload["vlans"]["vlan_id"] == "100"


def test_element_node_and_write_group_no_payload(xml2yaml_mod, tmp_path):
    resources = {"name": "configuration", "children": [{"name": "system", "type": "container"}]}
    system_elem = xml2yaml_mod.ElementTree.fromstring("<system><host-name>r1</host-name></system>")
    node = xml2yaml_mod.element_node(system_elem, "system", resources)
    assert node["name"] == "system"
    xml2yaml_mod.write_group_vars_flat(str(tmp_path), "all", {})
    assert not (tmp_path / "group_vars" / "all.yml").exists()
