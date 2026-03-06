#!/usr/bin/env python3
"""Stateful multi-device NETCONF-over-SSH mock for Terraform CI.

This simulator models a simplified Junos config lifecycle per device:
- Each device listens on its own TCP socket.
- Each device has independent candidate/running group config state.
- <load-configuration> mutates candidate.
- <edit-config operation="delete"> mutates candidate.
- <commit/> copies candidate -> running.
- <discard-changes/> restores candidate from running.
- <get-configuration> returns running configuration for requested group.

Usage example:
  python .github/scripts/netconf_mock_server.py \
    --host 127.0.0.1 \
    --username ci-user \
    --password ci-password \
    --device dc1-leaf1:8301 --device dc1-leaf2:8302
"""

from __future__ import annotations

import argparse
import asyncio
import copy
import json
import re
import signal
from dataclasses import dataclass, field
from pathlib import Path

import asyncssh

MSG_SEP = "]]>]]>"

HELLO = (
    '<?xml version="1.0" encoding="UTF-8"?>'
    '<hello xmlns="urn:ietf:params:xml:ns:netconf:base:1.0">'
    '<capabilities>'
    '<capability>urn:ietf:params:netconf:base:1.0</capability>'
    '</capabilities>'
    '<session-id>100</session-id>'
    '</hello>'
)


@dataclass
class DeviceState:
    name: str
    running_groups: dict[str, str] = field(default_factory=dict)
    candidate_groups: dict[str, str] = field(default_factory=dict)
    submitted_xml_by_group: dict[str, str] = field(default_factory=dict)
    rpc_log: list[str] = field(default_factory=list)
    history: list[dict[str, str]] = field(default_factory=list)

    def __post_init__(self) -> None:
        # Candidate starts as a copy of running, similar to Junos candidate model.
        self.candidate_groups = copy.deepcopy(self.running_groups)

    def snapshot(self) -> dict:
        return {
            "name": self.name,
            "running_groups": self.running_groups,
            "candidate_groups": self.candidate_groups,
            "submitted_xml_by_group": self.submitted_xml_by_group,
            "rpc_log": self.rpc_log,
            "history": self.history,
        }


class DeviceSSHServer(asyncssh.SSHServer):
    def __init__(self, username: str, password: str, state: DeviceState):
        self.username = username
        self.password = password
        self.state = state

    def begin_auth(self, username: str) -> bool:
        return True

    def password_auth_supported(self) -> bool:
        return True

    def validate_password(self, username: str, password: str) -> bool:
        return username == self.username and password == self.password

    def session_requested(self) -> asyncssh.SSHServerSession:
        return DeviceSession(self.state)


class DeviceSession(asyncssh.SSHServerSession):
    def __init__(self, state: DeviceState):
        self._chan: asyncssh.SSHServerChannel | None = None
        self._buf = ""
        self._state = state

    def connection_made(self, chan: asyncssh.SSHServerChannel) -> None:
        self._chan = chan

    def subsystem_requested(self, subsystem: str) -> bool:
        return subsystem == "netconf"

    def session_started(self) -> None:
        self._send_frame(HELLO)

    def data_received(self, data: str, datatype: int | None = None) -> None:
        del datatype
        self._buf += data
        while MSG_SEP in self._buf:
            raw, self._buf = self._buf.split(MSG_SEP, 1)
            req = raw.strip()
            if req:
                self._handle_rpc(req)

    def _send_frame(self, payload: str) -> None:
        if self._chan is not None:
            self._chan.write(payload + MSG_SEP + "\n")

    @staticmethod
    def _extract_message_id(xml_text: str) -> str:
        m = re.search(r'message-id="([^"]+)"', xml_text)
        return m.group(1) if m else "0"

    @staticmethod
    def _extract_group_name(xml_text: str) -> str:
        m = re.search(r"<name>([^<]+)</name>", xml_text)
        return m.group(1) if m else ""

    @staticmethod
    def _extract_configuration(xml_text: str) -> str:
        m = re.search(r"(<configuration>.*?</configuration>)", xml_text, flags=re.DOTALL)
        return m.group(1) if m else ""

    def _ok_reply(self, message_id: str) -> str:
        return f'<rpc-reply message-id="{message_id}"><ok/></rpc-reply>'

    def _append_history(self, op: str, detail: str) -> None:
        self._state.history.append({"op": op, "detail": detail})

    def _handle_rpc(self, xml_text: str) -> None:
        message_id = self._extract_message_id(xml_text)
        self._state.rpc_log.append(xml_text)

        # Ignore client hello after server hello.
        if "<hello" in xml_text:
            return

        if "<load-configuration" in xml_text:
            group_name = self._extract_group_name(xml_text)
            cfg = self._extract_configuration(xml_text)
            if group_name and cfg:
                self._state.candidate_groups[group_name] = cfg
                self._state.submitted_xml_by_group[group_name] = cfg
                self._append_history("load-configuration", f"group={group_name}")
            self._send_frame(self._ok_reply(message_id))
            return

        if "<edit-config>" in xml_text and 'operation="delete"' in xml_text:
            group_name = self._extract_group_name(xml_text)
            if group_name:
                self._state.candidate_groups.pop(group_name, None)
                self._append_history("edit-config-delete", f"group={group_name}")
            self._send_frame(self._ok_reply(message_id))
            return

        if "<discard-changes" in xml_text:
            self._state.candidate_groups = copy.deepcopy(self._state.running_groups)
            self._append_history("discard-changes", "candidate reset from running")
            self._send_frame(self._ok_reply(message_id))
            return

        if "<commit" in xml_text:
            self._state.running_groups = copy.deepcopy(self._state.candidate_groups)
            self._append_history("commit", f"groups={len(self._state.running_groups)}")
            self._send_frame(self._ok_reply(message_id))
            return

        if "<get-configuration>" in xml_text:
            group_name = self._extract_group_name(xml_text)
            if group_name and group_name in self._state.running_groups:
                cfg = self._state.running_groups[group_name]
            else:
                cfg = (
                    "<configuration><groups>"
                    f"<name>{group_name}</name>"
                    "</groups></configuration>"
                )
            self._append_history("get-configuration", f"group={group_name}")
            reply = f'<rpc-reply message-id="{message_id}">{cfg}</rpc-reply>'
            self._send_frame(reply)
            return

        if "<lock>" in xml_text or "<unlock>" in xml_text:
            self._send_frame(self._ok_reply(message_id))
            return

        # Default success reply for unrecognized requests.
        self._append_history("unknown", "default-ok")
        self._send_frame(self._ok_reply(message_id))


async def run_server(args: argparse.Namespace) -> None:
    stop_event = asyncio.Event()
    host_key = asyncssh.generate_private_key("ssh-rsa")
    servers = []
    device_states: dict[str, DeviceState] = {}

    def _shutdown() -> None:
        stop_event.set()

    loop = asyncio.get_running_loop()
    for sig in (signal.SIGINT, signal.SIGTERM):
        loop.add_signal_handler(sig, _shutdown)

    for device_spec in args.device:
        name, port_str = device_spec.split(":", 1)
        port = int(port_str)
        state = DeviceState(name=name)
        device_states[name] = state

        server = await asyncssh.create_server(
            lambda s=state: DeviceSSHServer(args.username, args.password, s),
            args.host,
            port,
            server_host_keys=[host_key],
            encoding="utf-8",
        )
        servers.append(server)

    await stop_event.wait()

    for server in servers:
        server.close()
        await server.wait_closed()

    if args.state_dump:
        dump_path = Path(args.state_dump)
        dump_path.parent.mkdir(parents=True, exist_ok=True)
        serialized = {name: state.snapshot() for name, state in device_states.items()}
        dump_path.write_text(json.dumps(serialized, indent=2), encoding="utf-8")


def parse_args() -> argparse.Namespace:
    parser = argparse.ArgumentParser()
    parser.add_argument("--host", default="127.0.0.1")
    parser.add_argument("--username", default="ci-user")
    parser.add_argument("--password", default="ci-password")
    parser.add_argument(
        "--device",
        action="append",
        required=True,
        help="Device listener in format <device-name>:<port>. May be repeated.",
    )
    parser.add_argument(
        "--state-dump",
        default="",
        help="Optional JSON path to dump per-device running/candidate/history on shutdown.",
    )
    return parser.parse_args()


def main() -> None:
    args = parse_args()
    asyncio.run(run_server(args))


if __name__ == "__main__":
    main()
