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
    python netconf_mock/netconf_mock_server.py \
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
import logging
import re
import signal
import traceback
import xml.etree.ElementTree as ET
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


logger = logging.getLogger("netconf-mock")


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
        logger.debug("auth attempt user=%s accepted=%s", username, username == self.username and password == self.password)
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
        logger.debug("session opened device=%s", self._state.name)

    def subsystem_requested(self, subsystem: str) -> bool:
        logger.debug("subsystem requested device=%s subsystem=%s", self._state.name, subsystem)
        return subsystem == "netconf"

    def session_started(self) -> None:
        logger.debug("session started device=%s", self._state.name)
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
        # Prefer XML parsing so attribute quoting/formatting differences do not break matching.
        try:
            root = ET.fromstring(xml_text)
            message_id = root.attrib.get("message-id")
            if message_id:
                return message_id
        except ET.ParseError:
            pass

        # Fallback regex supports single or double quotes and optional spaces.
        m = re.search(r"message-id\s*=\s*(['\"])(.*?)\1", xml_text)
        return m.group(2) if m else "0"

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
        logger.debug("device=%s op=%s detail=%s", self._state.name, op, detail)

    def _handle_load_configuration(self, xml_text: str, message_id: str) -> bool:
        if "<load-configuration" not in xml_text:
            return False

        group_name = self._extract_group_name(xml_text)
        cfg = self._extract_configuration(xml_text)
        if group_name and cfg:
            self._state.candidate_groups[group_name] = cfg
            self._state.submitted_xml_by_group[group_name] = cfg
            self._append_history("load-configuration", f"group={group_name}")
        self._send_frame(self._ok_reply(message_id))
        return True

    def _handle_edit_delete(self, xml_text: str, message_id: str) -> bool:
        if "<edit-config>" not in xml_text or 'operation="delete"' not in xml_text:
            return False

        group_name = self._extract_group_name(xml_text)
        if group_name:
            self._state.candidate_groups.pop(group_name, None)
            self._append_history("edit-config-delete", f"group={group_name}")
        self._send_frame(self._ok_reply(message_id))
        return True

    def _handle_discard_changes(self, xml_text: str, message_id: str) -> bool:
        if "<discard-changes" not in xml_text:
            return False

        self._state.candidate_groups = copy.deepcopy(self._state.running_groups)
        self._append_history("discard-changes", "candidate reset from running")
        self._send_frame(self._ok_reply(message_id))
        return True

    def _handle_commit(self, xml_text: str, message_id: str) -> bool:
        if "<commit" not in xml_text:
            return False

        self._state.running_groups = copy.deepcopy(self._state.candidate_groups)
        self._append_history("commit", f"groups={len(self._state.running_groups)}")
        self._send_frame(self._ok_reply(message_id))
        return True

    def _handle_get_configuration(self, xml_text: str, message_id: str) -> bool:
        if "<get-configuration>" not in xml_text:
            return False

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
        return True

    def _handle_lock_unlock(self, xml_text: str, message_id: str) -> bool:
        if "<lock>" not in xml_text and "<unlock>" not in xml_text:
            return False

        self._send_frame(self._ok_reply(message_id))
        return True

    def _handle_rpc(self, xml_text: str) -> None:
        message_id = self._extract_message_id(xml_text)
        self._state.rpc_log.append(xml_text)

        # Ignore client hello after server hello.
        if "<hello" in xml_text:
            return

        handlers = (
            self._handle_load_configuration,
            self._handle_edit_delete,
            self._handle_discard_changes,
            self._handle_commit,
            self._handle_get_configuration,
            self._handle_lock_unlock,
        )
        for handler in handlers:
            if handler(xml_text, message_id):
                return

        # Default success reply for unrecognized requests.
        self._append_history("unknown", "default-ok")
        self._send_frame(self._ok_reply(message_id))


async def run_server(args: argparse.Namespace) -> None:
    stop_event = asyncio.Event()
    _install_loop_exception_handler()
    _install_signal_handlers(stop_event)

    device_specs = _collect_device_specs(args)
    logger.info("starting NETCONF mock for %d device listeners", len(device_specs))

    host_key = asyncssh.generate_private_key("ssh-rsa")
    servers, device_states = await _start_device_listeners(args, host_key, device_specs)

    logger.info("all device listeners started")
    await stop_event.wait()

    await _close_servers(servers)
    logger.info("all device listeners closed")

    _dump_state_if_requested(args.state_dump, device_states)


def _install_loop_exception_handler() -> None:
    loop = asyncio.get_running_loop()

    def _loop_exception_handler(_loop: asyncio.AbstractEventLoop, context: dict) -> None:
        msg = context.get("message", "asyncio loop exception")
        exc = context.get("exception")
        logger.error("%s", msg)
        if exc is not None:
            logger.exception("event loop exception", exc_info=exc)
        else:
            logger.error("context: %s", context)

    loop.set_exception_handler(_loop_exception_handler)


def _install_signal_handlers(stop_event: asyncio.Event) -> None:
    loop = asyncio.get_running_loop()

    def _shutdown() -> None:
        logger.info("shutdown signal received")
        stop_event.set()

    for sig in (signal.SIGINT, signal.SIGTERM):
        loop.add_signal_handler(sig, _shutdown)


def _collect_device_specs(args: argparse.Namespace) -> list[str]:
    device_specs = list(args.device)
    if not args.devices_file:
        return device_specs

    logger.info("loading devices from file: %s", args.devices_file)
    file_data = Path(args.devices_file).read_text(encoding="utf-8")
    # Be tolerant of accidentally escaped newlines ("\\n") from shell-generated files.
    file_data = file_data.replace("\\n", "\n")
    for line in file_data.splitlines():
        entry = line.strip()
        if not entry or entry.startswith("#"):
            continue
        device_specs.append(entry)
    return device_specs


async def _start_device_listeners(
    args: argparse.Namespace,
    host_key: asyncssh.SSHKey,
    device_specs: list[str],
) -> tuple[list[asyncssh.SSHAcceptor], dict[str, DeviceState]]:
    servers: list[asyncssh.SSHAcceptor] = []
    device_states: dict[str, DeviceState] = {}

    for device_spec in device_specs:
        name, port_str = device_spec.split(":", 1)
        port = int(port_str)
        state = DeviceState(name=name)
        device_states[name] = state

        logger.info("binding device=%s host=%s port=%d", name, args.host, port)
        server = await asyncssh.create_server(
            lambda s=state: DeviceSSHServer(args.username, args.password, s),
            args.host,
            port,
            server_host_keys=[host_key],
            encoding="utf-8",
        )
        servers.append(server)

    return servers, device_states


async def _close_servers(servers: list[asyncssh.SSHAcceptor]) -> None:
    for server in servers:
        server.close()
        await server.wait_closed()


def _dump_state_if_requested(state_dump: str, device_states: dict[str, DeviceState]) -> None:
    if not state_dump:
        return

    dump_path = Path(state_dump)
    dump_path.parent.mkdir(parents=True, exist_ok=True)
    serialized = {name: state.snapshot() for name, state in device_states.items()}
    dump_path.write_text(json.dumps(serialized, indent=2), encoding="utf-8")
    logger.info("state dumped to %s", dump_path)


def parse_args() -> argparse.Namespace:
    parser = argparse.ArgumentParser()
    parser.add_argument("--host", default="127.0.0.1")
    parser.add_argument("--username", default="ci-user")
    parser.add_argument("--password", default="ci-password")
    parser.add_argument(
        "--device",
        action="append",
        default=[],
        help="Device listener in format <device-name>:<port>. May be repeated.",
    )
    parser.add_argument(
        "--devices-file",
        default="",
        help="Optional file with one <device-name>:<port> entry per line.",
    )
    parser.add_argument(
        "--state-dump",
        default="",
        help="Optional JSON path to dump per-device running/candidate/history on shutdown.",
    )
    parser.add_argument(
        "--log-level",
        default="INFO",
        help="Python logging level (DEBUG, INFO, WARNING, ERROR).",
    )
    args = parser.parse_args()
    if not args.device and not args.devices_file:
        parser.error("Provide at least one --device or a --devices-file")
    return args


def main() -> None:
    args = parse_args()
    logging.basicConfig(
        level=getattr(logging, args.log_level.upper(), logging.INFO),
        format="%(asctime)s %(levelname)s %(name)s: %(message)s",
    )
    try:
        logger.info("netconf mock booting")
        asyncio.run(run_server(args))
        logger.info("netconf mock exiting cleanly")
    except Exception as exc:  # pragma: no cover - fatal diagnostics path
        logger.error("fatal error in netconf mock: %s", exc)
        logger.error("python traceback:\n%s", traceback.format_exc())
        raise


if __name__ == "__main__":
    main()
