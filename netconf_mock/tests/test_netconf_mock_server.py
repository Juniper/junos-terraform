import importlib.util
import json
import sys
from pathlib import Path

import pytest


pytest.importorskip("asyncssh")


MODULE_PATH = (
    Path(__file__).resolve().parents[1] / "netconf_mock_server.py"
)
SPEC = importlib.util.spec_from_file_location(
    "netconf_mock_server", MODULE_PATH
)
MODULE = importlib.util.module_from_spec(SPEC)
assert SPEC is not None and SPEC.loader is not None
sys.modules[SPEC.name] = MODULE
SPEC.loader.exec_module(MODULE)


class FakeChannel:
    def __init__(self):
        self.writes = []

    def write(self, payload: str) -> None:
        self.writes.append(payload)


@pytest.fixture
def state_and_session():
    state = MODULE.DeviceState(name="device-1")
    session = MODULE.DeviceSession(state)
    channel = FakeChannel()
    session.connection_made(channel)
    return state, session, channel


def base_group_xml(body: str) -> str:
    return f"<configuration><groups><name>base-config</name>{body}</groups></configuration>"


def test_parse_args_accepts_repeated_devices(monkeypatch):
    monkeypatch.setattr(
        sys,
        "argv",
        [
            "netconf_mock_server.py",
            "--device",
            "leaf1:8301",
            "--device",
            "leaf2:8302",
        ],
    )

    args = MODULE.parse_args()

    assert args.device == ["leaf1:8301", "leaf2:8302"]
    assert args.devices_file == ""


def test_parse_args_rejects_missing_device_source(monkeypatch):
    monkeypatch.setattr(sys, "argv", ["netconf_mock_server.py"])

    with pytest.raises(SystemExit):
        MODULE.parse_args()


def test_collect_device_specs_reads_escaped_newlines(tmp_path):
    devices_file = tmp_path / "devices.txt"
    devices_file.write_text(
        "leaf1:8301\\nleaf2:8302\n#comment\n",
        encoding="utf-8",
    )

    args = MODULE.argparse.Namespace(
        device=["spine1:8311"],
        devices_file=str(devices_file),
    )

    specs = MODULE._collect_device_specs(args)

    assert specs == ["spine1:8311", "leaf1:8301", "leaf2:8302"]


def test_device_ssh_server_validates_password_when_auth_enabled():
    state = MODULE.DeviceState(name="leaf1")
    server = MODULE.DeviceSSHServer(
        username="ci-user",
        password="ci-password",
        state=state,
        disable_auth=False,
    )

    assert server.begin_auth("ci-user") is True
    assert server.validate_password("ci-user", "ci-password") is True
    assert server.validate_password("ci-user", "wrong-password") is False


def test_load_configuration_updates_candidate_and_submitted(state_and_session):
    state, session, channel = state_and_session

    rpc = (
        '<rpc message-id="11">'
        "<load-configuration action=\"merge\" format=\"xml\">"
        "<configuration><groups><name>base-config</name>"
        "<system><host-name>leaf1</host-name></system>"
        "</groups></configuration>"
        "</load-configuration>"
        "</rpc>"
    )

    handled = session._handle_load_configuration(rpc, "11")

    assert handled is True
    assert "base-config" in state.candidate_groups
    assert state.submitted_xml_by_group["base-config"].startswith(
        "<configuration>"
    )
    assert 'message-id="11"' in channel.writes[-1]
    assert "<ok/>" in channel.writes[-1]


def test_edit_patch_replace_updates_existing_group(state_and_session):
    state, session, channel = state_and_session

    state.candidate_groups["base-config"] = (
        "<configuration><groups><name>base-config</name>"
        "<system><host-name>leaf1</host-name></system>"
        "</groups></configuration>"
    )

    rpc = (
        '<rpc message-id="11">'
        "<edit-config><target><candidate/></target>"
        '<default-operation>none</default-operation>'
        '<config xmlns:nc="urn:ietf:params:xml:ns:netconf:base:1.0">'
        '<configuration><system><host-name nc:operation="replace">leaf2</host-name>'
        "</system></configuration></config></edit-config>"
        "</rpc>"
    )

    handled = session._handle_edit_patch(rpc, "11")

    assert handled is True
    assert "leaf2" in state.candidate_groups["base-config"]
    assert "leaf1" not in state.candidate_groups["base-config"]
    assert state.history[-1]["op"] == "edit-config-patch"
    assert 'message-id="11"' in channel.writes[-1]


def test_edit_patch_create_initializes_default_group(state_and_session):
    state, session, _channel = state_and_session

    rpc = (
        '<rpc message-id="15">'
        "<edit-config><target><candidate/></target>"
        '<default-operation>none</default-operation>'
        '<config xmlns:nc="urn:ietf:params:xml:ns:netconf:base:1.0">'
        '<configuration><system><host-name nc:operation="create">leaf1</host-name>'
        "</system></configuration></config></edit-config>"
        "</rpc>"
    )

    handled = session._handle_edit_patch(rpc, "15")

    assert handled is True
    assert "base-config" in state.candidate_groups
    assert "<name>base-config</name>" in state.candidate_groups["base-config"]
    assert "leaf1" in state.candidate_groups["base-config"]


def test_handle_rpc_routes_patch_delete_before_group_delete(state_and_session):
    state, session, _channel = state_and_session

    state.candidate_groups["base-config"] = (
        "<configuration><groups><name>base-config</name><system>"
        "<services><ssh/></services></system></groups></configuration>"
    )

    rpc = (
        '<rpc message-id="16">'
        "<edit-config><target><candidate/></target>"
        '<default-operation>none</default-operation>'
        '<config xmlns:nc="urn:ietf:params:xml:ns:netconf:base:1.0">'
        '<configuration><system><services nc:operation="delete"/></system></configuration>'
        "</config></edit-config></rpc>"
    )

    session._handle_rpc(rpc)

    assert "<services>" not in state.candidate_groups["base-config"]
    assert state.history[-1]["op"] == "edit-config-patch"


def test_edit_patch_create_keyed_list_entry_appends_new_interface(state_and_session):
    state, session, _channel = state_and_session

    state.candidate_groups["base-config"] = base_group_xml(
        "<interfaces><interface><name>xe-0/0/0</name>"
        "<description>uplink-0</description></interface></interfaces>"
    )

    rpc = (
        '<rpc message-id="17">'
        "<edit-config><target><candidate/></target>"
        '<default-operation>none</default-operation>'
        '<config xmlns:nc="urn:ietf:params:xml:ns:netconf:base:1.0">'
        '<configuration><interfaces><interface nc:operation="create">'
        '<name>xe-0/0/1</name><description>uplink-1</description>'
        "</interface></interfaces></configuration></config></edit-config></rpc>"
    )

    assert session._handle_edit_patch(rpc, "17") is True
    updated = state.candidate_groups["base-config"]
    assert updated.count("<interface>") == 2
    assert "xe-0/0/1" in updated
    assert "uplink-1" in updated


def test_edit_patch_replace_leaf_in_keyed_nested_path(state_and_session):
    state, session, _channel = state_and_session

    state.candidate_groups["base-config"] = base_group_xml(
        "<interfaces><interface><name>lo0</name><unit><name>0</name>"
        "<family><inet><address><name>203.0.113.1/32</name>"
        "<description>old-desc</description></address></inet></family>"
        "</unit></interface></interfaces>"
    )

    rpc = (
        '<rpc message-id="18">'
        "<edit-config><target><candidate/></target>"
        '<default-operation>none</default-operation>'
        '<config xmlns:nc="urn:ietf:params:xml:ns:netconf:base:1.0">'
        '<configuration><interfaces><interface><name>lo0</name><unit><name>0</name>'
        '<family><inet><address><name>203.0.113.1/32</name>'
        '<description nc:operation="replace">new-desc</description>'
        "</address></inet></family></unit></interface></interfaces></configuration>"
        "</config></edit-config></rpc>"
    )

    assert session._handle_edit_patch(rpc, "18") is True
    updated = state.candidate_groups["base-config"]
    assert "new-desc" in updated
    assert "old-desc" not in updated


def test_edit_patch_delete_keyed_list_entry_removes_matching_interface(state_and_session):
    state, session, _channel = state_and_session

    state.candidate_groups["base-config"] = base_group_xml(
        "<interfaces>"
        "<interface><name>xe-0/0/0</name><description>keep</description></interface>"
        "<interface><name>xe-0/0/1</name><description>delete-me</description></interface>"
        "</interfaces>"
    )

    rpc = (
        '<rpc message-id="19">'
        "<edit-config><target><candidate/></target>"
        '<default-operation>none</default-operation>'
        '<config xmlns:nc="urn:ietf:params:xml:ns:netconf:base:1.0">'
        '<configuration><interfaces><interface nc:operation="delete">'
        "<name>xe-0/0/1</name></interface></interfaces></configuration>"
        "</config></edit-config></rpc>"
    )

    assert session._handle_edit_patch(rpc, "19") is True
    updated = state.candidate_groups["base-config"]
    assert "xe-0/0/1" not in updated
    assert "delete-me" not in updated
    assert "xe-0/0/0" in updated


def test_edit_patch_create_leaf_list_appends_new_value(state_and_session):
    state, session, _channel = state_and_session

    state.candidate_groups["base-config"] = base_group_xml(
        "<system><domain-search>example.com</domain-search></system>"
    )

    rpc = (
        '<rpc message-id="20">'
        "<edit-config><target><candidate/></target>"
        '<default-operation>none</default-operation>'
        '<config xmlns:nc="urn:ietf:params:xml:ns:netconf:base:1.0">'
        '<configuration><system><domain-search nc:operation="create">lab.example</domain-search>'
        "</system></configuration></config></edit-config></rpc>"
    )

    assert session._handle_edit_patch(rpc, "20") is True
    updated = state.candidate_groups["base-config"]
    assert updated.count("<domain-search>") == 2
    assert "example.com" in updated
    assert "lab.example" in updated


def test_edit_patch_delete_leaf_list_value_removes_only_matching_value(state_and_session):
    state, session, _channel = state_and_session

    state.candidate_groups["base-config"] = base_group_xml(
        "<system><domain-search>example.com</domain-search>"
        "<domain-search>lab.example</domain-search></system>"
    )

    rpc = (
        '<rpc message-id="21">'
        "<edit-config><target><candidate/></target>"
        '<default-operation>none</default-operation>'
        '<config xmlns:nc="urn:ietf:params:xml:ns:netconf:base:1.0">'
        '<configuration><system><domain-search nc:operation="delete">lab.example</domain-search>'
        "</system></configuration></config></edit-config></rpc>"
    )

    assert session._handle_edit_patch(rpc, "21") is True
    updated = state.candidate_groups["base-config"]
    assert "lab.example" not in updated
    assert "example.com" in updated
    assert updated.count("<domain-search>") == 1


def test_edit_patch_mixed_operations_updates_configuration_consistently(state_and_session):
    state, session, _channel = state_and_session

    state.candidate_groups["base-config"] = base_group_xml(
        "<system><host-name>leaf1</host-name><services><ssh/></services></system>"
    )

    rpc = (
        '<rpc message-id="22">'
        "<edit-config><target><candidate/></target>"
        '<default-operation>none</default-operation>'
        '<config xmlns:nc="urn:ietf:params:xml:ns:netconf:base:1.0">'
        '<configuration><system>'
        '<host-name nc:operation="replace">leaf2</host-name>'
        '<services nc:operation="delete"/>'
        '<domain-name nc:operation="create">example.net</domain-name>'
        "</system></configuration></config></edit-config></rpc>"
    )

    assert session._handle_edit_patch(rpc, "22") is True
    updated = state.candidate_groups["base-config"]
    assert "leaf2" in updated
    assert "leaf1" not in updated
    assert "<services>" not in updated
    assert "example.net" in updated


def test_commit_discard_and_delete_flow(state_and_session):
    state, session, _channel = state_and_session

    state.candidate_groups["base-config"] = (
        "<configuration><groups><name>base-config</name>"
        "</groups></configuration>"
    )
    assert session._handle_commit("<commit/>", "12") is True
    assert "base-config" in state.running_groups

    state.candidate_groups["temp"] = (
        "<configuration><groups><name>temp</name></groups></configuration>"
    )
    assert session._handle_discard_changes("<discard-changes/>", "13") is True
    assert "temp" not in state.candidate_groups

    delete_rpc = (
        '<rpc message-id="14">'
        "<edit-config><target><candidate/></target>"
        "<config><configuration><groups operation=\"delete\">"
        "<name>base-config</name>"
        "</groups></configuration></config></edit-config>"
        "</rpc>"
    )
    assert session._handle_edit_delete(delete_rpc, "14") is True
    assert "base-config" not in state.candidate_groups


def test_get_configuration_returns_group_or_fallback(state_and_session):
    state, session, channel = state_and_session

    state.running_groups["base-config"] = (
        "<configuration><groups><name>base-config</name>"
        "<system><services/></system></groups></configuration>"
    )

    rpc_existing = (
        '<rpc message-id="21"><get-configuration><configuration><groups><name>'
        "base-config"
        "</name></groups></configuration></get-configuration></rpc>"
    )
    assert session._handle_get_configuration(rpc_existing, "21") is True
    assert "services" in channel.writes[-1]

    rpc_missing = (
        '<rpc message-id="22"><get-configuration><configuration><groups><name>'
        "does-not-exist"
        "</name></groups></configuration></get-configuration></rpc>"
    )
    assert session._handle_get_configuration(rpc_missing, "22") is True
    assert "<name>does-not-exist</name>" in channel.writes[-1]


def test_handle_rpc_hello_and_unknown(state_and_session):
    state, session, channel = state_and_session

    session._handle_rpc("<hello><capabilities/></hello>")
    assert channel.writes == []

    session._handle_rpc('<rpc message-id="30"><foo/></rpc>')
    assert 'message-id="30"' in channel.writes[-1]
    assert "<ok/>" in channel.writes[-1]
    assert state.history[-1]["op"] == "unknown"


def test_extract_message_id_accepts_single_quoted_attributes():
    xml = "<rpc message-id='77'><lock><target><candidate/></target></lock></rpc>"

    assert MODULE.DeviceSession._extract_message_id(xml) == "77"


def test_extract_message_id_accepts_prefixed_attribute():
    xml = (
        '<rpc xmlns:nc="urn:ietf:params:xml:ns:netconf:base:1.0" '
        'nc:message-id="88"><lock><target><candidate/></target></lock></rpc>'
    )

    assert MODULE.DeviceSession._extract_message_id(xml) == "88"


def test_extract_message_id_returns_empty_when_missing():
    xml = "<rpc><lock><target><candidate/></target></lock></rpc>"

    assert MODULE.DeviceSession._extract_message_id(xml) == ""


def test_extract_group_name_prefers_groups_name_over_nested_name():
    xml = (
        '<rpc message-id="101">'
        '<load-configuration action="replace" format="xml">'
        '<configuration><groups><name>base-config</name>'
        '<interfaces><interface><name>lo0</name></interface></interfaces>'
        '</groups></configuration>'
        '</load-configuration></rpc>'
    )

    assert MODULE.DeviceSession._extract_group_name(xml) == "base-config"


def test_load_configuration_updates_all_groups_in_payload(state_and_session):
    state, session, _channel = state_and_session

    state.candidate_groups["base-config"] = (
        "<configuration><groups><name>base-config</name>"
        "<interfaces><interface><name>lo0</name><unit><name>0</name>"
        "<family><inet><address><name>203.0.113.250/32</name></address>"
        "</inet></family></unit></interface></interfaces>"
        "</groups></configuration>"
    )

    rpc = (
        '<rpc message-id="120">'
        '<load-configuration action="replace" format="xml">'
        '<configuration>'
        '<groups><name>overlay-config</name><protocols/></groups>'
        '<groups><name>base-config</name>'
        '<interfaces><interface><name>lo0</name><unit><name>0</name>'
        '<family><inet><address><name>203.0.113.10/32</name></address>'
        '</inet></family></unit></interface></interfaces>'
        '</groups>'
        '</configuration>'
        '</load-configuration>'
        '</rpc>'
    )

    handled = session._handle_load_configuration(rpc, "120")

    assert handled is True
    assert "overlay-config" in state.candidate_groups
    assert "base-config" in state.candidate_groups
    assert "203.0.113.10/32" in state.candidate_groups["base-config"]
    assert "203.0.113.250/32" not in state.candidate_groups["base-config"]


def test_load_configuration_regex_fallback_updates_all_groups(state_and_session):
    state, session, _channel = state_and_session

    state.candidate_groups["base-config"] = (
        "<configuration><groups><name>base-config</name>"
        "<interfaces><interface><name>lo0</name><unit><name>0</name>"
        "<family><inet><address><name>203.0.113.250/32</name></address>"
        "</inet></family></unit></interface></interfaces>"
        "</groups></configuration>"
    )

    # Deliberately malformed for XML parser (unbound prefixes) to force regex fallback.
    rpc = (
        '<nc:rpc message-id="121">'
        '<load-configuration action="replace" format="xml">'
        '<configuration>'
        '<groups><name>overlay-config</name><protocols/></groups>'
        '<groups><name>base-config</name>'
        '<interfaces><interface><name>lo0</name><unit><name>0</name>'
        '<family><inet><address><name>203.0.113.10/32</name></address>'
        '</inet></family></unit></interface></interfaces>'
        '</groups>'
        '</configuration>'
        '</load-configuration>'
        '</nc:rpc>'
    )

    handled = session._handle_load_configuration(rpc, "121")

    assert handled is True
    assert "overlay-config" in state.candidate_groups
    assert "base-config" in state.candidate_groups
    assert "203.0.113.10/32" in state.candidate_groups["base-config"]
    assert "203.0.113.250/32" not in state.candidate_groups["base-config"]


def test_dump_state_if_requested_writes_json(tmp_path):
    out_file = tmp_path / "state.json"
    state = MODULE.DeviceState(name="leaf1")
    state.running_groups["base-config"] = "<configuration/>"

    MODULE._dump_state_if_requested(str(out_file), {"leaf1": state})

    data = json.loads(out_file.read_text(encoding="utf-8"))
    assert "leaf1" in data
    assert data["leaf1"]["name"] == "leaf1"
    assert "base-config" in data["leaf1"]["running_groups"]
