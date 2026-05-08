import importlib.util
import json
import pathlib
import sys
from importlib.machinery import SourceFileLoader

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


def _basic_schema() -> dict:
    return {
        "root": {
            "name": "root",
            "children": [
                {
                    "name": "configuration",
                    "children": [
                        {
                            "name": "system",
                            "type": "container",
                            "children": [
                                {"name": "host-name", "type": "leaf"},
                                {"name": "product-name", "type": "leaf"},
                                {"name": "profile", "type": "leaf"},
                                {
                                    "name": "services",
                                    "type": "container",
                                    "children": [{"name": "ssh", "type": "leaf"}],
                                },
                            ],
                        },
                        {
                            "name": "routing-options",
                            "type": "container",
                            "children": [{"name": "router-id", "type": "leaf"}],
                        },
                        {
                            "name": "vlans",
                            "type": "container",
                            "children": [{"name": "enabled", "type": "leaf"}],
                        },
                        {
                            "name": "security",
                            "type": "container",
                            "children": [{"name": "policies", "type": "leaf"}],
                        },
                    ],
                }
            ],
        }
    }


def _write_provider(tmp_path: pathlib.Path, provider_dir_name: str, role_name: str) -> pathlib.Path:
    provider_dir = tmp_path / provider_dir_name
    (provider_dir / "roles" / role_name).mkdir(parents=True, exist_ok=True)
    schema_file = provider_dir / "trimmed_schema.json"
    schema_file.write_text(json.dumps(_basic_schema()))
    return schema_file


def _write_xml(
    xml_file: pathlib.Path,
    *,
    host_name: str,
    product_name: str,
    profile: str,
    router_id: str | None = None,
    vlan_enabled: str | None = None,
    security_policies: str | None = None,
) -> None:
    parts = [
        "<configuration>",
        "<system>",
        f"<host-name>{host_name}</host-name>",
        f"<product-name>{product_name}</product-name>",
        f"<profile>{profile}</profile>",
        "<services><ssh/></services>",
        "</system>",
    ]
    if router_id is not None:
        parts.extend([
            "<routing-options>",
            f"<router-id>{router_id}</router-id>",
            "</routing-options>",
        ])
    if vlan_enabled is not None:
        parts.extend([
            "<vlans>",
            f"<enabled>{vlan_enabled}</enabled>",
            "</vlans>",
        ])
    if security_policies is not None:
        parts.extend([
            "<security>",
            f"<policies>{security_policies}</policies>",
            "</security>",
        ])
    parts.append("</configuration>")
    xml_file.write_text("".join(parts))


def test_prepare_and_detect_device_type(xml2yaml_mod):
    assert xml2yaml_mod.prepare_tag("host-name.v4") == "host_name_v4"
    assert xml2yaml_mod.detect_device_type({"system": {"product_name": "QFX5110"}}) == "qfx"
    assert xml2yaml_mod.detect_device_type({"chassis": {"product_name": "SRX300"}}) == "srx"
    assert xml2yaml_mod.detect_device_type({"system": {"product_name": "Unknown"}}) is None


def test_keyed_list_intersection_and_subtraction(xml2yaml_mod):
    left = {
        "policy_options": {
            "policy_statement": [
                {"name": "A", "term": 1, "action": "accept"},
                {"name": "B", "term": 2, "action": "reject"},
            ]
        }
    }
    right = {
        "policy_options": {
            "policy_statement": [
                {"name": "A", "term": 1, "action": "accept", "comment": "preserved"},
                {"name": "C", "term": 3, "action": "accept"},
            ]
        }
    }

    shared = xml2yaml_mod.intersect_values(left, right)
    assert shared == {
        "policy_options": {
            "policy_statement": [{"name": "A", "term": 1, "action": "accept"}]
        }
    }

    left_delta = xml2yaml_mod.subtract_common(left, shared)
    assert left_delta == {
        "policy_options": {
            "policy_statement": [{"name": "B", "term": 2, "action": "reject"}]
        }
    }


def test_bool_field_is_not_used_as_list_identity_key(xml2yaml_mod):
    items = [
        [{"community_name": "dc1-leaf1", "add": True}],
        [{"community_name": "dc1-leaf2", "add": True}],
    ]

    assert xml2yaml_mod._infer_identity_key_for_dict_lists(items) == "community_name"

    left = {
        "policy_options": {
            "policy_statement": [{
                "name": "IPCLOS_BGP_EXP",
                "term": [{
                    "name": "loopback",
                    "then": {
                        "community": [{
                            "community_name": "dc1-leaf1",
                            "add": True,
                        }]
                    },
                }],
            }]
        }
    }
    right = {
        "policy_options": {
            "policy_statement": [{
                "name": "IPCLOS_BGP_EXP",
                "term": [{
                    "name": "loopback",
                    "then": {
                        "community": [{
                            "community_name": "dc1-leaf2",
                            "add": True,
                        }]
                    },
                }],
            }]
        }
    }

    shared = xml2yaml_mod.intersect_values(left, right)
    assert shared == {
        "policy_options": {
            "policy_statement": [{
                "name": "IPCLOS_BGP_EXP",
                "term": [{
                    "name": "loopback",
                    "then": {},
                }],
            }]
        }
    }

    left_delta = xml2yaml_mod.subtract_common(left, shared)
    assert left_delta == left


def test_derive_provider_key_from_role_directory(xml2yaml_mod, tmp_path):
    schema_file = _write_provider(
        tmp_path,
        "ansible-provider-junos-vqfx-ansible-role",
        "vqfx-ansible-role_role",
    )
    assert xml2yaml_mod.derive_provider_key(str(schema_file)) == "vqfx-ansible-role_role"


def test_parse_xml_to_payload(xml2yaml_mod, tmp_path):
    schema_file = _write_provider(
        tmp_path,
        "ansible-provider-junos-vqfx-ansible-role",
        "vqfx-ansible-role_role",
    )
    xml_file = tmp_path / "router1.xml"
    _write_xml(xml_file, host_name="router1", product_name="QFX5100", profile="qfx", router_id="1.1.1.1")

    schema = json.loads(schema_file.read_text())
    hostname, payload, dtype = xml2yaml_mod.parse_xml_to_payload(str(xml_file), schema)
    assert hostname == "router1"
    assert payload["system"]["host_name"] == "router1"
    assert payload["routing_options"]["router_id"] == "1.1.1.1"
    assert dtype == "qfx"


def test_generic_group_discovery_does_not_depend_on_hostnames(xml2yaml_mod):
    host_var_entries = [
        ("alpha", {"shared": {"ssh": True}, "feature": {"kind": "edge"}, "vlans": {"enabled": "true"}}),
        ("omega", {"shared": {"ssh": True}, "feature": {"kind": "edge"}, "vlans": {"enabled": "true"}}),
        ("one", {"shared": {"ssh": True}, "feature": {"kind": "core"}}),
        ("two", {"shared": {"ssh": True}, "feature": {"kind": "core"}}),
    ]

    groups = xml2yaml_mod.derive_common_host_groups(
        host_var_entries,
        min_hosts=2,
        max_groups=4,
        min_benefit_score=1,
        min_new_paths=0,
    )
    memberships = {frozenset(members) for members in groups.values()}
    assert frozenset({"alpha", "omega", "one", "two"}) in memberships
    assert frozenset({"alpha", "omega"}) in memberships
    assert frozenset({"one", "two"}) in memberships


def test_collapse_daisy_chains_removes_low_value_intermediate_node(xml2yaml_mod):
    common_groups = {
        "group:candidate1": ["h1", "h2", "h3", "h4"],
        "group:candidate2": ["h1", "h2", "h3"],
        "group:candidate3": ["h1", "h2"],
    }
    parent_map = {
        "group:candidate1": None,
        "group:candidate2": "group:candidate1",
        "group:candidate3": "group:candidate2",
    }
    payload_by_hostname = {
        "h1": {"common": {"x": 1}, "pair": {"y": 1}},
        "h2": {"common": {"x": 1}, "pair": {"y": 1}},
        "h3": {"common": {"x": 1}},
        "h4": {"common": {"x": 1}, "other": {"z": 1}},
    }

    reduced_groups, reduced_parent_map = xml2yaml_mod.collapse_daisy_chains(
        common_groups,
        parent_map,
        payload_by_hostname,
        max_single_child_chain=2,
        min_unique_paths=1,
    )

    assert "group:candidate2" not in reduced_groups
    assert reduced_parent_map["group:candidate3"] == "group:candidate1"


def test_xml2yaml_main_multi_provider_shared_all_yaml(xml2yaml_mod, tmp_path, monkeypatch):
    qfx_schema = _write_provider(
        tmp_path,
        "ansible-provider-junos-vqfx-ansible-role",
        "vqfx-ansible-role_role",
    )
    srx_schema = _write_provider(
        tmp_path,
        "ansible-provider-junos-srx-ansible-role",
        "srx-ansible-role_role",
    )

    qfx_xml1 = tmp_path / "qfx-a.xml"
    qfx_xml2 = tmp_path / "qfx-b.xml"
    qfx_xml3 = tmp_path / "qfx-c.xml"
    _write_xml(qfx_xml1, host_name="qfx-a", product_name="QFX5100", profile="qfx", router_id="1.1.1.1", vlan_enabled="true")
    _write_xml(qfx_xml2, host_name="qfx-b", product_name="QFX5100", profile="qfx", router_id="1.1.1.2", vlan_enabled="true")
    _write_xml(qfx_xml3, host_name="qfx-c", product_name="QFX5100", profile="qfx", router_id="1.1.1.3")

    srx_xml1 = tmp_path / "srx-a.xml"
    srx_xml2 = tmp_path / "srx-b.xml"
    _write_xml(srx_xml1, host_name="srx-a", product_name="SRX300", profile="srx", security_policies="strict")
    _write_xml(srx_xml2, host_name="srx-b", product_name="SRX300", profile="srx", security_policies="strict")

    out_dir = tmp_path / "out"

    monkeypatch.setattr(
        sys,
        "argv",
        [
            "jtaf-xml2yaml",
            "-j",
            str(qfx_schema),
            "-x",
            str(qfx_xml1),
            str(qfx_xml2),
            str(qfx_xml3),
            "-d",
            str(out_dir),
        ],
    )
    xml2yaml_mod.main()

    monkeypatch.setattr(
        sys,
        "argv",
        [
            "jtaf-xml2yaml",
            "-j",
            str(srx_schema),
            "-x",
            str(srx_xml1),
            str(srx_xml2),
            "-d",
            str(out_dir),
        ],
    )
    xml2yaml_mod.main()

    all_payload = yaml.safe_load((out_dir / "group_vars" / "all.yaml").read_text())
    registry = all_payload["meta"]["jtaf_registry"]
    qfx_entry = registry["providers"]["vqfx-ansible-role_role"]
    srx_entry = registry["providers"]["srx-ansible-role_role"]

    assert set(registry["providers"].keys()) == {"vqfx-ansible-role_role", "srx-ansible-role_role"}
    assert all_payload["system"]["services"]["ssh"] is True
    assert "product_name" not in all_payload["system"]

    qfx_range = qfx_entry["group_range"]
    srx_range = srx_entry["group_range"]
    assert qfx_range["end"] < srx_range["start"]

    qfx_root_path = out_dir / "group_vars" / qfx_entry["group_paths"][qfx_entry["root_group"]] / "all.yaml"
    srx_root_path = out_dir / "group_vars" / srx_entry["group_paths"][srx_entry["root_group"]] / "all.yaml"
    assert qfx_root_path.exists()
    assert srx_root_path.exists()

    qfx_root_payload = yaml.safe_load(qfx_root_path.read_text())
    srx_root_payload = yaml.safe_load(srx_root_path.read_text())
    assert qfx_root_payload["system"]["product_name"] == "QFX5100"
    assert srx_root_payload["system"]["product_name"] == "SRX300"

    hosts_text = (out_dir / "hosts").read_text()
    for hostname in ["qfx-a", "qfx-b", "qfx-c", "srx-a", "srx-b"]:
        assert hostname in hosts_text
    for group_name in qfx_entry["group_names"] + srx_entry["group_names"]:
        assert f"[{group_name}]" in hosts_text


def test_xml2yaml_main_removes_old_owned_groups_on_rerun(xml2yaml_mod, tmp_path, monkeypatch):
    schema_file = _write_provider(
        tmp_path,
        "ansible-provider-junos-vqfx-ansible-role",
        "vqfx-ansible-role_role",
    )
    xml1 = tmp_path / "first-a.xml"
    xml2 = tmp_path / "first-b.xml"
    _write_xml(xml1, host_name="a", product_name="QFX5100", profile="qfx", vlan_enabled="true")
    _write_xml(xml2, host_name="b", product_name="QFX5100", profile="qfx", vlan_enabled="true")

    out_dir = tmp_path / "out"
    monkeypatch.setattr(
        sys,
        "argv",
        ["jtaf-xml2yaml", "-j", str(schema_file), "-x", str(xml1), str(xml2), "-d", str(out_dir)],
    )
    xml2yaml_mod.main()

    first_all = yaml.safe_load((out_dir / "group_vars" / "all.yaml").read_text())
    first_entry = first_all["meta"]["jtaf_registry"]["providers"]["vqfx-ansible-role_role"]
    assert first_entry["group_names"]
    first_root_path = out_dir / "group_vars" / first_entry["group_paths"][first_entry["root_group"]]
    assert first_root_path.exists()

    xml3 = tmp_path / "second-a.xml"
    _write_xml(xml3, host_name="a", product_name="QFX5100", profile="qfx")
    monkeypatch.setattr(
        sys,
        "argv",
        ["jtaf-xml2yaml", "-j", str(schema_file), "-x", str(xml3), "-d", str(out_dir)],
    )
    xml2yaml_mod.main()

    second_all = yaml.safe_load((out_dir / "group_vars" / "all.yaml").read_text())
    second_entry = second_all["meta"]["jtaf_registry"]["providers"]["vqfx-ansible-role_role"]
    assert second_entry.get("group_names", []) == []
    assert not first_root_path.exists()

    hosts_text = (out_dir / "hosts").read_text()
    for old_group in first_entry["group_names"]:
        assert f"[{old_group}]" not in hosts_text


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


def test_element_node_returns_matching_schema_child(xml2yaml_mod):
    resources = {"name": "configuration", "children": [{"name": "system", "type": "container"}]}
    system_elem = xml2yaml_mod.ElementTree.fromstring("<system><host-name>r1</host-name></system>")
    node = xml2yaml_mod.element_node(system_elem, "system", resources)
    assert node["name"] == "system"
