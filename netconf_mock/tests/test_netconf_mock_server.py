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
    assert channel.writes[-1].startswith(
        '<rpc-reply message-id="11"><ok/></rpc-reply>'
    )


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
    assert channel.writes[-1].startswith(
        '<rpc-reply message-id="30"><ok/></rpc-reply>'
    )
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


def test_dump_state_if_requested_writes_json(tmp_path):
    out_file = tmp_path / "state.json"
    state = MODULE.DeviceState(name="leaf1")
    state.running_groups["base-config"] = "<configuration/>"

    MODULE._dump_state_if_requested(str(out_file), {"leaf1": state})

    data = json.loads(out_file.read_text(encoding="utf-8"))
    assert "leaf1" in data
    assert data["leaf1"]["name"] == "leaf1"
    assert "base-config" in data["leaf1"]["running_groups"]
