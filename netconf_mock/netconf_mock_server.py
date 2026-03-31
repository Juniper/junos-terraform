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
    def __init__(self, username: str, password: str, state: DeviceState, disable_auth: bool):
        self.username = username
        self.password = password
        self.state = state
        self.disable_auth = disable_auth

    def begin_auth(self, username: str) -> bool:
        logger.info("device=%s begin_auth username=%s", self.state.name, username)
        if self.disable_auth:
            logger.info("device=%s auth disabled; accepting without credentials", self.state.name)
            return False
        return True

    def password_auth_supported(self) -> bool:
        logger.info("device=%s password_auth_supported=true", self.state.name)
        return True

    def validate_password(self, username: str, password: str) -> bool:
        accepted = username == self.username and password == self.password
        logger.info("device=%s validate_password username=%s accepted=%s", self.state.name, username, accepted)
        return accepted

    def session_requested(self) -> asyncssh.SSHServerSession:
        logger.info("device=%s session_requested", self.state.name)
        return DeviceSession(self.state)


class DeviceSession(asyncssh.SSHServerSession):
    @staticmethod
    def _local_name(tag: str) -> str:
        return tag.rsplit("}", 1)[-1] if "}" in tag else tag

    def __init__(self, state: DeviceState):
        self._chan: asyncssh.SSHServerChannel | None = None
        self._buf = ""
        self._state = state

    def connection_made(self, chan: asyncssh.SSHServerChannel) -> None:
        self._chan = chan
        logger.debug("session opened device=%s", self._state.name)

    def subsystem_requested(self, subsystem: str) -> bool:
        logger.info("device=%s subsystem_requested subsystem=%s", self._state.name, subsystem)
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
            logger.debug("device=%s tx frame=%s", self._state.name, payload[:300])
            self._chan.write(payload + MSG_SEP + "\n")

    @staticmethod
    def _extract_message_id(xml_text: str) -> str:
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
    def _extract_configuration(xml_text: str) -> str:
        m = re.search(r"(<configuration>.*?</configuration>)", xml_text, flags=re.DOTALL)
        return m.group(1) if m else ""

    @staticmethod
    def _parse_xml(xml_text: str) -> ET.Element | None:
        try:
            return ET.fromstring(xml_text)
        except ET.ParseError:
            return None

    @staticmethod
    def _find_first_configuration(root: ET.Element) -> ET.Element | None:
        for elem in root.iter():
            if DeviceSession._local_name(elem.tag) == "configuration":
                return elem
        return None

    @staticmethod
    def _extract_group_name_from_groups_elem(groups_elem: ET.Element) -> str:
        for child in list(groups_elem):
            if DeviceSession._local_name(child.tag) == "name" and child.text:
                return child.text.strip()
        return ""

    @staticmethod
    def _extract_groups_configurations_regex(xml_text: str) -> dict[str, str]:
        groups_by_name: dict[str, str] = {}
        for group_block in re.findall(r"(<groups>.*?</groups>)", xml_text, flags=re.DOTALL):
            name_match = re.search(r"<name>\s*([^<]+?)\s*</name>", group_block)
            if not name_match:
                continue
            group_name = name_match.group(1).strip()
            if not group_name:
                continue
            groups_by_name[group_name] = f"<configuration>{group_block}</configuration>"
        return groups_by_name

    @staticmethod
    def _extract_groups_from_configuration_set(xml_text: str) -> dict[str, str]:
        """Extract per-group payload from set-style load-configuration RPCs.

        Some client stacks send `<load-configuration format="text">` with
        `<configuration-set>` commands instead of XML `<configuration><groups>`.
        Capture `set groups <name> ...` lines so downstream assertions can still
        validate expected rendered content.
        """

        m = re.search(
            r"<configuration-set>\s*(.*?)\s*</configuration-set>",
            xml_text,
            flags=re.DOTALL,
        )
        if not m:
            return {}

        lines = [line.strip() for line in m.group(1).splitlines() if line.strip()]
        groups_lines: dict[str, list[str]] = {}
        for line in lines:
            lm = re.match(r"^set\s+groups\s+(\S+)\s+(.+)$", line)
            if not lm:
                continue
            group_name = lm.group(1)
            groups_lines.setdefault(group_name, []).append(line)

        return {name: "\n".join(group_lines) for name, group_lines in groups_lines.items()}

    @staticmethod
    def _extract_groups_configurations(xml_text: str) -> dict[str, str]:
        """Extract per-group configuration blobs from a load-configuration RPC."""

        groups_by_name: dict[str, str] = {}
        root = DeviceSession._parse_xml(xml_text)
        if root is None:
            return DeviceSession._extract_groups_configurations_regex(xml_text)

        config_elem = DeviceSession._find_first_configuration(root)
        if config_elem is None:
            return DeviceSession._extract_groups_configurations_regex(xml_text)

        for child in list(config_elem):
            if DeviceSession._local_name(child.tag) != "groups":
                continue
            group_name = DeviceSession._extract_group_name_from_groups_elem(child)
            if not group_name:
                continue
            groups_xml = ET.tostring(child, encoding="unicode")
            groups_by_name[group_name] = f"<configuration>{groups_xml}</configuration>"

        if groups_by_name:
            return groups_by_name
        return DeviceSession._extract_groups_configurations_regex(xml_text)

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
    ) -> ET.Element | None:
        patch_text = (patch_elem.text or "").strip()
        for candidate in candidates:
            if (candidate.text or "").strip() == patch_text:
                return candidate

        # Leaf-list create operations should append a new sibling when no
        # existing leaf matches the requested value.
        if operation == "create":
            return None

        if len(candidates) == 1:
            return candidates[0]
        return None

    @staticmethod
    def _find_matching_patch_target(
        parent: ET.Element,
        patch_elem: ET.Element,
        operation: str,
    ) -> ET.Element | None:
        candidates = DeviceSession._matching_children(parent, patch_elem)
        if not candidates:
            return None

        match_keys = DeviceSession._extract_match_keys(patch_elem)
        if match_keys:
            return DeviceSession._find_candidate_by_keys(candidates, match_keys)

        if not list(patch_elem):
            return DeviceSession._find_leaf_value_candidate(candidates, patch_elem, operation)

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

    def _merge_group_configuration(self, group_name: str, incoming_config: str) -> str:
        existing_root = self._load_group_configuration_root(group_name)
        incoming_root = self._parse_xml(incoming_config)
        if incoming_root is None:
            return incoming_config

        existing_groups = self._ensure_groups_elem(existing_root, group_name)
        incoming_groups = self._find_child(incoming_root, "groups")
        if incoming_groups is None:
            return incoming_config

        for child in list(incoming_groups):
            if self._local_name(child.tag) == "name":
                continue
            self._apply_patch_element(existing_groups, child)

        return ET.tostring(existing_root, encoding="unicode")

    def _resolve_patch_group_name(self) -> str:
        known_groups = set(self._state.candidate_groups)
        known_groups.update(self._state.running_groups)
        known_groups.update(self._state.submitted_xml_by_group)

        if "base-config" in known_groups:
            return "base-config"
        if len(known_groups) == 1:
            return next(iter(known_groups))
        if known_groups:
            return sorted(known_groups)[0]
        return "base-config"

    def _load_group_configuration_root(self, group_name: str) -> ET.Element:
        raw_config = (
            self._state.candidate_groups.get(group_name)
            or self._state.running_groups.get(group_name)
            or self._state.submitted_xml_by_group.get(group_name)
            or ""
        )
        if raw_config:
            parsed = self._parse_xml(raw_config)
            if parsed is not None:
                return parsed

        configuration = ET.Element("configuration")
        groups_elem = ET.SubElement(configuration, "groups")
        name_elem = ET.SubElement(groups_elem, "name")
        name_elem.text = group_name
        return configuration

    def _ensure_groups_elem(self, configuration: ET.Element, group_name: str) -> ET.Element:
        groups_elem = self._find_child(configuration, "groups")
        if groups_elem is None:
            groups_elem = ET.SubElement(configuration, "groups")

        name_elem = self._find_child(groups_elem, "name")
        if name_elem is None:
            name_elem = ET.Element("name")
            groups_elem.insert(0, name_elem)
        if not (name_elem.text or "").strip():
            name_elem.text = group_name
        return groups_elem

    def _apply_patch_element(self, parent: ET.Element, patch_elem: ET.Element) -> None:
        operation = self._extract_operation(patch_elem)
        target = self._find_matching_patch_target(parent, patch_elem, operation)

        if operation == "delete":
            if target is not None:
                parent.remove(target)
            return

        if target is None:
            target = ET.SubElement(parent, self._local_name(patch_elem.tag))

        children = list(patch_elem)
        if not children:
            target.text = patch_elem.text
            return

        if target is not None and self._deletes_matched_keyed_entry(patch_elem):
            parent.remove(target)
            return

        for child in children:
            self._apply_patch_element(target, child)

    def _handle_edit_patch(self, xml_text: str, message_id: str) -> bool:
        if "<edit-config>" not in xml_text or "<load-configuration" in xml_text:
            return False

        root = self._parse_xml(xml_text)
        if root is None:
            return False

        config_elem = None
        for elem in root.iter():
            if self._local_name(elem.tag) == "config":
                config_elem = elem
                break
        if config_elem is None:
            return False

        patch_configuration = self._find_first_configuration(config_elem)
        if patch_configuration is None:
            return False

        if not any(self._extract_operation(elem) for elem in patch_configuration.iter()):
            return False

        group_name = self._resolve_patch_group_name()
        configuration = self._load_group_configuration_root(group_name)
        groups_elem = self._ensure_groups_elem(configuration, group_name)

        for child in list(patch_configuration):
            self._apply_patch_element(groups_elem, child)

        updated_config = ET.tostring(configuration, encoding="unicode")
        self._state.candidate_groups[group_name] = updated_config
        self._state.submitted_xml_by_group[group_name] = updated_config
        self._append_history("edit-config-patch", f"group={group_name}")
        self._send_frame(self._ok_reply(message_id))
        return True

    def _ok_reply(self, message_id: str) -> str:
        return (
            '<rpc-reply xmlns="urn:ietf:params:xml:ns:netconf:base:1.0" '
            f'message-id="{message_id}"><ok/></rpc-reply>'
        )

    def _append_history(self, op: str, detail: str) -> None:
        self._state.history.append({"op": op, "detail": detail})
        logger.debug("device=%s op=%s detail=%s", self._state.name, op, detail)

    def _handle_load_configuration(self, xml_text: str, message_id: str) -> bool:
        if "<load-configuration" not in xml_text:
            return False

        action = self._extract_load_configuration_action(xml_text)
        groups_cfg = self._extract_groups_configurations(xml_text)
        if not groups_cfg:
            groups_cfg = self._extract_groups_from_configuration_set(xml_text)
        if groups_cfg:
            for group_name, cfg in groups_cfg.items():
                if action == "merge":
                    cfg = self._merge_group_configuration(group_name, cfg)
                self._state.candidate_groups[group_name] = cfg
                self._state.submitted_xml_by_group[group_name] = cfg
            self._append_history(
                "load-configuration",
                f"action={action or 'replace'} groups={','.join(sorted(groups_cfg.keys()))}",
            )
        else:
            # Fallback for malformed/minimal payloads.
            group_name = self._extract_group_name(xml_text)
            cfg = self._extract_configuration(xml_text)
            if group_name and cfg:
                if action == "merge":
                    cfg = self._merge_group_configuration(group_name, cfg)
                self._state.candidate_groups[group_name] = cfg
                self._state.submitted_xml_by_group[group_name] = cfg
                self._append_history(
                    "load-configuration",
                    f"action={action or 'replace'} group={group_name}",
                )
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
        reply = (
            '<rpc-reply xmlns="urn:ietf:params:xml:ns:netconf:base:1.0" '
            f'message-id="{message_id}">{cfg}</rpc-reply>'
        )
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
            lambda s=state: DeviceSSHServer(args.username, args.password, s, args.disable_auth),
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
