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
    running_config: str = ""
    candidate_config: str = ""
    submitted_xml: str = ""
    rpc_log: list[str] = field(default_factory=list)
    history: list[dict[str, str]] = field(default_factory=list)

    def __post_init__(self) -> None:
        """Initialize candidate configuration from running state."""
        # Candidate starts as a copy of running, similar to Junos candidate model.
        self.candidate_config = self.running_config

    def snapshot(self) -> dict:
        """Return a serializable point-in-time view of device state."""
        return {
            "name": self.name,
            "running_config": self.running_config,
            "candidate_config": self.candidate_config,
            "submitted_xml": self.submitted_xml,
            "rpc_log": self.rpc_log,
            "history": self.history,
        }


class DeviceSSHServer(asyncssh.SSHServer):
    def __init__(self, username: str, password: str, state: DeviceState, disable_auth: bool):
        """Store auth settings and bound device state for one SSH server."""
        self.username = username
        self.password = password
        self.state = state
        self.disable_auth = disable_auth

    def begin_auth(self, username: str) -> bool:
        """Tell asyncssh whether password authentication should proceed."""
        logger.info("device=%s begin_auth username=%s", self.state.name, username)
        if self.disable_auth:
            logger.info("device=%s auth disabled; accepting without credentials", self.state.name)
            return False
        return True

    def password_auth_supported(self) -> bool:
        """Advertise password-auth support to the SSH client."""
        logger.info("device=%s password_auth_supported=true", self.state.name)
        return True

    def validate_password(self, username: str, password: str) -> bool:
        """Validate credentials against configured static username/password."""
        accepted = username == self.username and password == self.password
        logger.info("device=%s validate_password username=%s accepted=%s", self.state.name, username, accepted)
        return accepted

    def session_requested(self) -> asyncssh.SSHServerSession:
        """Create a per-connection NETCONF session bound to this device."""
        logger.info("device=%s session_requested", self.state.name)
        return DeviceSession(self.state)


class DeviceSession(asyncssh.SSHServerSession):
    @staticmethod
    def _local_name(tag: str) -> str:
        """Return XML local tag name without namespace prefix."""
        return tag.rsplit("}", 1)[-1] if "}" in tag else tag

    def __init__(self, state: DeviceState):
        """Initialize session buffers and attach target device state."""
        self._chan: asyncssh.SSHServerChannel | None = None
        self._buf = ""
        self._state = state

    def connection_made(self, chan: asyncssh.SSHServerChannel) -> None:
        """Capture the opened SSH channel for outbound NETCONF frames."""
        self._chan = chan
        logger.debug("session opened device=%s", self._state.name)

    def subsystem_requested(self, subsystem: str) -> bool:
        """Allow only the NETCONF subsystem."""
        logger.info("device=%s subsystem_requested subsystem=%s", self._state.name, subsystem)
        return subsystem == "netconf"

    def session_started(self) -> None:
        """Send initial NETCONF <hello> once subsystem starts."""
        logger.debug("session started device=%s", self._state.name)
        self._send_frame(HELLO)

    def data_received(self, data: str, datatype: int | None = None) -> None:
        """Buffer framed NETCONF payloads and dispatch each complete RPC."""
        del datatype
        self._buf += data
        while MSG_SEP in self._buf:
            raw, self._buf = self._buf.split(MSG_SEP, 1)
            req = raw.strip()
            if req:
                self._handle_rpc(req)

    def _send_frame(self, payload: str) -> None:
        """Write one NETCONF message followed by end-of-message separator."""
        if self._chan is not None:
            logger.debug("device=%s tx frame=%s", self._state.name, payload[:300])
            self._chan.write(payload + MSG_SEP + "\n")

    @staticmethod
    def _extract_message_id(xml_text: str) -> str:
        """Extract NETCONF message-id via XML parse, then regex fallback."""
        # Prefer XML parsing so attribute quoting/formatting differences do not break matching.
        try:
            root = ET.fromstring(xml_text)
            for attr_name, attr_value in root.attrib.items():
                if attr_name == "message-id" or attr_name.endswith("}message-id"):
                    if attr_value:
                        return attr_value
        except ET.ParseError:
            pass

        # Fallback regex supports optional namespace prefixes and quote styles.
        m = re.search(
            r"(?:[A-Za-z_][\w.\-]*:)?message-id\s*=\s*(['\"])(.*?)\1",
            xml_text,
        )
        return m.group(2) if m else ""

    @staticmethod
    def _extract_group_name(xml_text: str) -> str:
        """Extract configuration group name from XML payload."""
        # Prefer XML parsing and only consider <configuration><groups><name>.
        try:
            root = ET.fromstring(xml_text)
            for elem in root.iter():
                if DeviceSession._local_name(elem.tag) != "groups":
                    continue
                for child in list(elem):
                    if DeviceSession._local_name(child.tag) == "name" and child.text:
                        return child.text.strip()
        except ET.ParseError:
            pass

        # Regex fallback restricted to groups/name to avoid matching unrelated <name> tags.
        m = re.search(r"<groups>\s*<name>([^<]+)</name>", xml_text, flags=re.DOTALL)
        return m.group(1).strip() if m else ""

    @staticmethod
    def _parse_xml(xml_text: str) -> ET.Element | None:
        """Parse XML text to ElementTree root, returning None on parse errors."""
        try:
            return ET.fromstring(xml_text)
        except ET.ParseError:
            return None

    @staticmethod
    def _find_first_configuration(root: ET.Element) -> ET.Element | None:
        """Find first element named configuration in a parsed RPC tree."""
        for elem in root.iter():
            if DeviceSession._local_name(elem.tag) == "configuration":
                return elem
        return None

    @staticmethod
    def _extract_operation(elem: ET.Element) -> str:
        for attr_name, attr_value in elem.attrib.items():
            if DeviceSession._local_name(attr_name) == "operation":
                return attr_value
        return ""

    @staticmethod
    def _find_child(parent: ET.Element, local_name: str) -> ET.Element | None:
        for child in list(parent):
            if DeviceSession._local_name(child.tag) == local_name:
                return child
        return None

    @staticmethod
    def _extract_match_keys(elem: ET.Element) -> dict[str, str]:
        keys: dict[str, str] = {}
        for child in list(elem):
            if DeviceSession._extract_operation(child):
                continue
            if list(child):
                continue
            if child.text is None:
                continue
            local_name = DeviceSession._local_name(child.tag)
            if local_name == "name":
                keys[local_name] = child.text.strip()
        return keys

    @staticmethod
    def _matching_children(parent: ET.Element, patch_elem: ET.Element) -> list[ET.Element]:
        local_name = DeviceSession._local_name(patch_elem.tag)
        return [
            child
            for child in list(parent)
            if DeviceSession._local_name(child.tag) == local_name
        ]

    @staticmethod
    def _find_candidate_by_keys(
        candidates: list[ET.Element],
        match_keys: dict[str, str],
    ) -> ET.Element | None:
        for candidate in candidates:
            matches = True
            for key_name, key_text in match_keys.items():
                key_elem = DeviceSession._find_child(candidate, key_name)
                if key_elem is None or (key_elem.text or "").strip() != key_text:
                    matches = False
                    break
            if matches:
                return candidate
        return None

    @staticmethod
    def _find_leaf_value_candidate(
        candidates: list[ET.Element],
        patch_elem: ET.Element,
        operation: str,
        sibling_count: int,
    ) -> ET.Element | None:
        patch_text = (patch_elem.text or "").strip()
        for candidate in candidates:
            if (candidate.text or "").strip() == patch_text:
                return candidate

        # Leaf-list create operations should append a new sibling when no
        # existing leaf matches the requested value.
        if operation == "create":
            return None

        # Merge payloads can contain repeated simple siblings such as
        # <protocol>bgp</protocol><protocol>direct</protocol>. Treat these as
        # leaf-list entries and append distinct values instead of overwriting
        # the first matching tag.
        if sibling_count > 1:
            return None

        if len(candidates) == 1:
            return candidates[0]
        return None

    @staticmethod
    def _find_matching_patch_target(
        parent: ET.Element,
        patch_elem: ET.Element,
        operation: str,
        sibling_count: int = 1,
    ) -> ET.Element | None:
        candidates = DeviceSession._matching_children(parent, patch_elem)
        if not candidates:
            return None

        match_keys = DeviceSession._extract_match_keys(patch_elem)
        if match_keys:
            return DeviceSession._find_candidate_by_keys(candidates, match_keys)

        if not list(patch_elem):
            return DeviceSession._find_leaf_value_candidate(
                candidates,
                patch_elem,
                operation,
                sibling_count,
            )

        if len(candidates) == 1:
            return candidates[0]

        return None

    @staticmethod
    def _deletes_matched_keyed_entry(patch_elem: ET.Element) -> bool:
        match_keys = DeviceSession._extract_match_keys(patch_elem)
        if not match_keys:
            return False

        for child in list(patch_elem):
            if DeviceSession._extract_operation(child) != "delete":
                continue
            local_name = DeviceSession._local_name(child.tag)
            if local_name not in match_keys:
                continue
            if (child.text or "").strip() == match_keys[local_name]:
                return True

        return False

    @staticmethod
    def _extract_load_configuration_action(xml_text: str) -> str:
        root = DeviceSession._parse_xml(xml_text)
        if root is not None:
            for elem in root.iter():
                if DeviceSession._local_name(elem.tag) != "load-configuration":
                    continue
                return (elem.attrib.get("action") or "").strip().lower()

        m = re.search(r'<load-configuration[^>]*\baction="([^"]+)"', xml_text)
        return m.group(1).strip().lower() if m else ""

    def _candidate_root(self) -> ET.Element:
        """The candidate configuration as a tree, empty if nothing is staged."""
        parsed = self._parse_xml(self._state.candidate_config)
        if parsed is not None:
            return parsed
        return ET.Element("configuration")

    def _incoming_configuration(self, xml_text: str) -> ET.Element | None:
        """The <configuration> element of an RPC. Some client stacks send
        unbound namespace prefixes, which make the whole document unparseable,
        so fall back to the configuration element on its own."""
        root = self._parse_xml(xml_text)
        if root is not None:
            found = self._find_first_configuration(root)
            if found is not None:
                return found

        m = re.search(r"(<configuration>.*</configuration>)", xml_text, flags=re.DOTALL)
        return self._parse_xml(m.group(1)) if m else None

    def _apply_children(self, parent: ET.Element, children: list[ET.Element]) -> None:
        """Apply top-level patch children, counting repeated names so that a
        leaf-list such as apply-groups keeps every entry rather than collapsing
        into the first one."""
        name_counts: dict[str, int] = {}
        for child in children:
            local_name = self._local_name(child.tag)
            name_counts[local_name] = name_counts.get(local_name, 0) + 1
        for child in children:
            self._apply_patch_element(parent, child, name_counts)

    def _group_configuration(self, group_name: str) -> str:
        """The committed configuration filtered to one group, as Junos returns
        it for <get-configuration><configuration><groups><name>."""
        parsed = self._parse_xml(self._state.running_config)
        configuration = ET.Element("configuration")
        if parsed is None:
            return ET.tostring(configuration, encoding="unicode")

        for child in list(parsed):
            if self._local_name(child.tag) != "groups":
                continue
            name_elem = self._find_child(child, "name")
            if name_elem is not None and (name_elem.text or "").strip() == group_name:
                configuration.append(copy.deepcopy(child))

        return ET.tostring(configuration, encoding="unicode")

    def _iter_patch_children(self, patch_configuration: ET.Element) -> list[ET.Element]:
        # groups is configuration like any other node, so a patch below it is
        # applied where it is sent rather than lifted to the top.
        return list(patch_configuration)

    def _apply_patch_element(
        self,
        parent: ET.Element,
        patch_elem: ET.Element,
        sibling_name_counts: dict[str, int] | None = None,
    ) -> None:
        operation = self._extract_operation(patch_elem)
        local_name = self._local_name(patch_elem.tag)
        sibling_count = 1 if sibling_name_counts is None else sibling_name_counts.get(local_name, 1)
        target = self._find_matching_patch_target(parent, patch_elem, operation, sibling_count)

        if operation == "delete":
            if target is not None:
                parent.remove(target)
            return

        if target is None:
            target = ET.SubElement(parent, local_name)

        children = list(patch_elem)
        if not children:
            target.text = patch_elem.text
            return

        if target is not None and self._deletes_matched_keyed_entry(patch_elem):
            parent.remove(target)
            return

        child_name_counts: dict[str, int] = {}
        for child in children:
            child_local_name = self._local_name(child.tag)
            child_name_counts[child_local_name] = child_name_counts.get(child_local_name, 0) + 1

        for child in children:
            self._apply_patch_element(target, child, child_name_counts)

    def _extract_patch_configuration(self, xml_text: str) -> ET.Element | None:
        if "<edit-config>" not in xml_text or "<load-configuration" in xml_text:
            return None

        root = self._parse_xml(xml_text)
        if root is None:
            return None

        config_elem = next(
            (elem for elem in root.iter() if self._local_name(elem.tag) == "config"),
            None,
        )
        if config_elem is None:
            return None

        patch_configuration = self._find_first_configuration(config_elem)
        if patch_configuration is None:
            return None

        if not any(self._extract_operation(elem) for elem in patch_configuration.iter()):
            return None

        return patch_configuration

    def _apply_patch_configuration(self, patch_configuration: ET.Element) -> str | None:
        patch_children = self._iter_patch_children(patch_configuration)
        if not patch_children:
            return None

        configuration = self._candidate_root()
        self._apply_children(configuration, patch_children)

        updated_config = ET.tostring(configuration, encoding="unicode")
        self._state.candidate_config = updated_config
        self._state.submitted_xml = updated_config
        return updated_config

    def _handle_edit_patch(self, xml_text: str, message_id: str) -> bool:
        patch_configuration = self._extract_patch_configuration(xml_text)
        if patch_configuration is None:
            return False

        if self._apply_patch_configuration(patch_configuration) is None:
            return False

        self._append_history("edit-config-patch", "candidate updated")
        self._send_frame(self._ok_reply(message_id))
        return True

    def _ok_reply(self, message_id: str) -> str:
        """Build a minimal NETCONF <ok/> rpc-reply for a message id."""
        return (
            '<rpc-reply xmlns="urn:ietf:params:xml:ns:netconf:base:1.0" '
            f'message-id="{message_id}"><ok/></rpc-reply>'
        )

    def _append_history(self, op: str, detail: str) -> None:
        """Append an operation record to the per-device history list."""
        self._state.history.append({"op": op, "detail": detail})
        logger.debug("device=%s op=%s detail=%s", self._state.name, op, detail)

    def _handle_load_configuration(self, xml_text: str, message_id: str) -> bool:
        """Apply load-configuration payload to candidate state when present."""
        if "<load-configuration" not in xml_text:
            return False

        action = self._extract_load_configuration_action(xml_text)
        incoming = self._incoming_configuration(xml_text)
        if incoming is not None and list(incoming):
            if action == "merge":
                configuration = self._candidate_root()
                self._apply_children(configuration, list(incoming))
            else:
                configuration = copy.deepcopy(incoming)
            cfg = ET.tostring(configuration, encoding="unicode")
            self._state.candidate_config = cfg
            self._state.submitted_xml = cfg
            self._append_history(
                "load-configuration", f"action={action or 'replace'}"
            )
        self._send_frame(self._ok_reply(message_id))
        return True

    def _handle_edit_delete(self, xml_text: str, message_id: str) -> bool:
        """Handle an edit-config that deletes whole elements of the candidate."""
        if "<edit-config>" not in xml_text or 'operation="delete"' not in xml_text:
            return False

        incoming = self._incoming_configuration(xml_text)
        if incoming is not None:
            configuration = self._candidate_root()
            self._apply_children(configuration, list(incoming))
            cfg = ET.tostring(configuration, encoding="unicode")
            self._state.candidate_config = cfg
            self._state.submitted_xml = cfg
            self._append_history("edit-config-delete", "candidate updated")
        self._send_frame(self._ok_reply(message_id))
        return True

    def _handle_discard_changes(self, xml_text: str, message_id: str) -> bool:
        """Reset candidate configuration back to current running state."""
        if "<discard-changes" not in xml_text:
            return False

        self._state.candidate_config = self._state.running_config
        self._append_history("discard-changes", "candidate reset from running")
        self._send_frame(self._ok_reply(message_id))
        return True

    def _handle_commit(self, xml_text: str, message_id: str) -> bool:
        """Promote candidate configuration to running configuration."""
        if "<commit" not in xml_text:
            return False

        self._state.running_config = self._state.candidate_config
        self._append_history("commit", f"bytes={len(self._state.running_config)}")
        self._send_frame(self._ok_reply(message_id))
        return True

    def _handle_get_configuration(self, xml_text: str, message_id: str) -> bool:
        """Return the committed configuration, or one group of it if asked."""
        if "<get-configuration>" not in xml_text:
            return False

        group_name = self._extract_group_name(xml_text)
        if group_name:
            cfg = self._group_configuration(group_name)
        else:
            cfg = self._state.running_config or "<configuration/>"
        self._append_history("get-configuration", f"group={group_name}")
        reply = (
            '<rpc-reply xmlns="urn:ietf:params:xml:ns:netconf:base:1.0" '
            f'message-id="{message_id}">{cfg}</rpc-reply>'
        )
        self._send_frame(reply)
        return True

    def _handle_lock_unlock(self, xml_text: str, message_id: str) -> bool:
        """Acknowledge lock/unlock RPCs with an OK reply."""
        if "<lock>" not in xml_text and "<unlock>" not in xml_text:
            return False

        self._send_frame(self._ok_reply(message_id))
        return True

    def _handle_rpc(self, xml_text: str) -> None:
        """Route a NETCONF RPC to known handlers or default OK behavior."""
        message_id = self._extract_message_id(xml_text)
        self._state.rpc_log.append(xml_text)
        logger.debug(
            "device=%s rx message_id=%s rpc=%s",
            self._state.name,
            message_id if message_id else "<missing>",
            xml_text[:300],
        )

        # Ignore client hello after server hello.
        if "<hello" in xml_text:
            return

        if not message_id:
            self._append_history("invalid-rpc", "missing-message-id")
            logger.warning(
                "device=%s received rpc without parseable message-id; dropping request",
                self._state.name,
            )
            return

        handlers = (
            self._handle_load_configuration,
            self._handle_edit_patch,
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
    """Run multi-device NETCONF listeners until shutdown signal is received."""
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
    """Install verbose asyncio loop exception logging."""
    loop = asyncio.get_running_loop()

    def _loop_exception_handler(_loop: asyncio.AbstractEventLoop, context: dict) -> None:
        """Log uncaught event-loop exceptions with traceback context."""
        msg = context.get("message", "asyncio loop exception")
        exc = context.get("exception")
        logger.error("%s", msg)
        if exc is not None:
            logger.exception("event loop exception", exc_info=exc)
        else:
            logger.error("context: %s", context)

    loop.set_exception_handler(_loop_exception_handler)


def _install_signal_handlers(stop_event: asyncio.Event) -> None:
    """Install SIGINT/SIGTERM handlers that trigger graceful shutdown."""
    loop = asyncio.get_running_loop()

    def _shutdown() -> None:
        """Signal coroutine shutdown by setting the shared stop event."""
        logger.info("shutdown signal received")
        stop_event.set()

    for sig in (signal.SIGINT, signal.SIGTERM):
        loop.add_signal_handler(sig, _shutdown)


def _collect_device_specs(args: argparse.Namespace) -> list[str]:
    """Collect device specs from repeated args and optional file input."""
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
    """Create and bind one asyncssh server per device specification."""
    servers: list[asyncssh.SSHAcceptor] = []
    device_states: dict[str, DeviceState] = {}

    for device_spec in device_specs:
        name, port_str = device_spec.split(":", 1)
        port = int(port_str)
        state = DeviceState(name=name)
        device_states[name] = state

        logger.info("binding device=%s host=%s port=%d", name, args.host, port)
        server = await asyncssh.create_server(
            lambda s=state: DeviceSSHServer(args.username, args.password, s, args.disable_auth),
            args.host,
            port,
            server_host_keys=[host_key],
            encoding="utf-8",
        )
        servers.append(server)

    return servers, device_states


async def _close_servers(servers: list[asyncssh.SSHAcceptor]) -> None:
    """Close all asyncssh server listeners and await completion."""
    for server in servers:
        server.close()
        await server.wait_closed()


def _dump_state_if_requested(state_dump: str, device_states: dict[str, DeviceState]) -> None:
    """Write JSON state snapshot if a dump path is configured."""
    if not state_dump:
        return

    dump_path = Path(state_dump)
    dump_path.parent.mkdir(parents=True, exist_ok=True)
    serialized = {name: state.snapshot() for name, state in device_states.items()}
    dump_path.write_text(json.dumps(serialized, indent=2), encoding="utf-8")
    logger.info("state dumped to %s", dump_path)


def parse_args() -> argparse.Namespace:
    """Parse and validate command-line arguments for the mock server."""
    parser = argparse.ArgumentParser(
        description="Run a stateful NETCONF-over-SSH mock server for integration tests."
    )
    parser.add_argument("--host", default="127.0.0.1", help="Bind address for all device listeners.")
    parser.add_argument("--username", default="ci-user", help="Accepted NETCONF username.")
    parser.add_argument("--password", default="ci-password", help="Accepted NETCONF password.")
    parser.add_argument(
        "--disable-auth",
        action="store_true",
        help="Disable SSH authentication checks for mock-only compatibility testing.",
    )
    parser.add_argument(
        "--device",
        action="append",
        default=[],
        help="Device listener in format <device-name>:<port>. May be repeated.",
    )
    parser.add_argument(
        "--devices-file",
        default="",
        help="File with one <device-name>:<port> entry per line.",
    )
    parser.add_argument(
        "--state-dump",
        default="",
        help="JSON path to write per-device running/candidate/history at shutdown.",
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
    """CLI entry point for running the stateful NETCONF mock service."""
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
