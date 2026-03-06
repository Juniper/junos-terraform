#!/usr/bin/env python3
"""Minimal NETCONF-over-SSH mock for CI Terraform integration tests.

Implements enough NETCONF behavior for provider CRUD cycles:
- SSH password auth
- netconf subsystem
- hello exchange
- <load-configuration>, <edit-config>, <get-configuration>, <commit/>, <discard-changes/>
"""

from __future__ import annotations

import argparse
import asyncio
import re
import signal
from dataclasses import dataclass, field

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
class NetconfState:
    groups_by_name: dict[str, str] = field(default_factory=dict)


class MockSSHServer(asyncssh.SSHServer):
    def __init__(self, username: str, password: str, state: NetconfState):
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
        return NetconfSession(self.state)


class NetconfSession(asyncssh.SSHServerSession):
    def __init__(self, state: NetconfState):
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

    def _ok_reply(self, message_id: str) -> str:
        return f'<rpc-reply message-id="{message_id}"><ok/></rpc-reply>'

    def _handle_rpc(self, xml_text: str) -> None:
        message_id = self._extract_message_id(xml_text)

        # Ignore client hello after server hello.
        if "<hello" in xml_text:
            return

        if "<load-configuration" in xml_text:
            group_name = self._extract_group_name(xml_text)
            config_match = re.search(
                r"(<configuration>.*</configuration>)",
                xml_text,
                flags=re.DOTALL,
            )
            if group_name and config_match:
                self._state.groups_by_name[group_name] = config_match.group(1)
            self._send_frame(self._ok_reply(message_id))
            return

        if "<edit-config>" in xml_text and "operation=\"delete\"" in xml_text:
            group_name = self._extract_group_name(xml_text)
            if group_name:
                self._state.groups_by_name.pop(group_name, None)
            self._send_frame(self._ok_reply(message_id))
            return

        if "<get-configuration>" in xml_text:
            group_name = self._extract_group_name(xml_text)
            if group_name and group_name in self._state.groups_by_name:
                cfg = self._state.groups_by_name[group_name]
            else:
                cfg = (
                    "<configuration><groups>"
                    f"<name>{group_name}</name>"
                    "</groups></configuration>"
                )
            reply = f'<rpc-reply message-id="{message_id}">{cfg}</rpc-reply>'
            self._send_frame(reply)
            return

        if (
            "<commit" in xml_text
            or "<discard-changes" in xml_text
            or "<lock>" in xml_text
            or "<unlock>" in xml_text
        ):
            self._send_frame(self._ok_reply(message_id))
            return

        # Default success reply for unrecognized requests.
        self._send_frame(self._ok_reply(message_id))


async def run_server(host: str, port: int, username: str, password: str) -> None:
    state = NetconfState()
    stop_event = asyncio.Event()
    host_key = asyncssh.generate_private_key("ssh-rsa")

    def _shutdown() -> None:
        stop_event.set()

    loop = asyncio.get_running_loop()
    for sig in (signal.SIGINT, signal.SIGTERM):
        loop.add_signal_handler(sig, _shutdown)

    await asyncssh.create_server(
        lambda: MockSSHServer(username, password, state),
        host,
        port,
        server_host_keys=[host_key],
        encoding="utf-8",
    )

    await stop_event.wait()


def main() -> None:
    parser = argparse.ArgumentParser()
    parser.add_argument("--host", default="127.0.0.1")
    parser.add_argument("--port", type=int, default=8300)
    parser.add_argument("--username", default="ci-user")
    parser.add_argument("--password", default="ci-password")
    args = parser.parse_args()

    asyncio.run(run_server(args.host, args.port, args.username, args.password))


if __name__ == "__main__":
    main()
