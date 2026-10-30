"""Connection and inventory helpers for the YANG-choice encoding probe.

Credentials are supplied at runtime (inventory file, environment variable or
prompt) and are never written to disk or printed by this package.
"""
from __future__ import annotations

import json
import os
import re
import socket
from getpass import getpass

from ncclient import manager
from ncclient.operations import RPCError
from ncclient.xml_ import new_ele

DEFAULT_PORT = 830
ENV_DEVICES = "JTAF_PROBE_DEVICES"
ENV_PASSWORD = "JTAF_PROBE_PASSWORD"
ENV_USER = "JTAF_PROBE_USER"


class Target:
    """One device to probe."""

    def __init__(self, name, host, username, password, port=DEFAULT_PORT):
        self.name = name
        self.host = host
        self.username = username
        self.password = password
        self.port = int(port or DEFAULT_PORT)

    def __repr__(self):  # password deliberately omitted
        return f"<Target {self.name} {self.username}@{self.host}:{self.port}>"


def load_inventory(path=None):
    """Load a device inventory JSON file.

    Accepts a nested shape::

        {"device1": {"ip": "...", "port": 22, "username": "...",
                     "auth": {"type": "password", "password": "..."}}}

    and a flatter shape::

        {"device1": {"host": "...", "username": "...", "password": "..."}}
    """
    path = path or os.environ.get(ENV_DEVICES)
    if not path:
        return {}
    with open(path) as fh:
        raw = json.load(fh)

    devices = {}
    for name, entry in raw.items():
        if not isinstance(entry, dict):
            continue
        host = entry.get("ip") or entry.get("host") or entry.get("hostname")
        password = entry.get("password")
        if password is None:
            password = (entry.get("auth") or {}).get("password")
        devices[name] = {
            "host": host,
            "username": entry.get("username") or entry.get("user"),
            "password": password,
            # Inventories often record the SSH port; NETCONF still listens on 830.
            "port": entry.get("netconf_port") or DEFAULT_PORT,
        }
    return devices


def resolve_targets(names=None, devices_path=None, host=None, username=None,
                    port=None, password=None):
    """Build the list of targets from an inventory and/or explicit host flags."""
    if host:
        username = username or os.environ.get(ENV_USER)
        if not username:
            raise SystemExit("--user is required with --host "
                             f"(or set {ENV_USER})")
        password = password or os.environ.get(ENV_PASSWORD) or getpass(
            f"Password for {username}@{host}: ")
        return [Target(host, host, username, password, port or DEFAULT_PORT)]

    inventory = load_inventory(devices_path)
    if not inventory:
        raise SystemExit(
            "No devices. Pass --devices <file.json> (or set "
            f"{ENV_DEVICES}), or use --host/--user.")

    wanted = names or list(inventory)
    targets = []
    for name in wanted:
        entry = inventory.get(name)
        if entry is None:
            raise SystemExit(f"{name!r} is not in the inventory")
        pw = password or entry["password"] or os.environ.get(ENV_PASSWORD)
        if not pw:
            pw = getpass(f"Password for {entry['username']}@{name}: ")
        targets.append(Target(name, entry["host"], entry["username"], pw,
                              port or entry["port"]))
    return targets


def connect(target, timeout=45):
    return manager.connect(
        host=target.host,
        port=target.port,
        username=target.username,
        password=target.password,
        hostkey_verify=False,
        allow_agent=False,
        look_for_keys=False,
        timeout=timeout,
        device_params={"name": "junos"},
    )


def port_open(host, port, timeout=5):
    try:
        with socket.create_connection((host, port), timeout=timeout):
            return True
    except OSError:
        return False


def facts(conn):
    """Model, version and hostname, without requiring jnpr.junos."""
    xml = conn.dispatch(new_ele("get-software-information")).tostring.decode()

    def grab(tag):
        m = re.search(rf"<{tag}>([^<]+)</{tag}>", xml)
        return m.group(1) if m else "?"

    version = grab("junos-version")
    if version == "?":
        m = re.search(r"<package-information>.*?<comment>[^\[]*\[([^\]]+)\]", xml, re.S)
        version = m.group(1) if m else "?"
    return {"host_name": grab("host-name"), "model": grab("product-model"),
            "version": version}


def knob_committed(conn):
    """Is `system services netconf yang-compliant` in the committed config?"""
    reply = conn.get_config(source="running", filter=(
        "subtree",
        "<configuration><system><services><netconf/></services></system></configuration>"))
    return "yang-compliant" in reply.data_xml


def candidate_diff(conn):
    """Candidate-versus-committed diff as CLI text; '' when the candidate is clean."""
    rpc = new_ele("get-configuration")
    rpc.set("compare", "rollback")
    rpc.set("rollback", "0")
    rpc.set("format", "text")
    text = conn.dispatch(rpc).tostring.decode()
    start, end = text.find("<configuration-output>"), text.find("</configuration-output>")
    body = text[start + 22:end] if start != -1 and end != -1 else ""
    return "\n".join(line.rstrip() for line in body.strip().splitlines() if line.strip())


def try_load(conn, config_xml):
    """Load into the candidate and discard. Returns (accepted, detail)."""
    try:
        conn.edit_config(target="candidate", config=config_xml)
        return True, ""
    except RPCError as exc:
        return False, " ".join(str(exc).split())
    finally:
        try:
            conn.discard_changes()
        except RPCError:
            pass


def load_then_diff(conn, config_xml):
    """Load into the candidate, capture the resulting CLI diff, discard."""
    try:
        conn.edit_config(target="candidate", config=config_xml)
    except RPCError as exc:
        return None, " ".join(str(exc).split())
    try:
        return candidate_diff(conn), ""
    finally:
        try:
            conn.discard_changes()
        except RPCError:
            pass
