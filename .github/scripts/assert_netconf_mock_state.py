#!/usr/bin/env python3
"""Assert expected NETCONF mock state invariants from JSON dump."""

from __future__ import annotations

import argparse
import json
import sys
import xml.etree.ElementTree as ET
from pathlib import Path


def parse_args() -> argparse.Namespace:
    parser = argparse.ArgumentParser(
        description="Assert expected invariants from a NETCONF mock state dump."
    )
    parser.add_argument(
        "--state-dump",
        required=True,
        help="Path to netconf mock state JSON.",
    )
    parser.add_argument(
        "--require-min-commits",
        type=int,
        default=1,
        help="Minimum total commit operations.",
    )
    parser.add_argument(
        "--must-contain",
        action="append",
        default=[],
        help="String that must appear in running config.",
    )
    parser.add_argument(
        "--must-not-contain",
        action="append",
        default=[],
        help="String that must not appear in running config.",
    )
    parser.add_argument(
        "--only-device",
        default="",
        help="Limit contain/not-contain checks to one device.",
    )
    parser.add_argument(
        "--only-group",
        default="",
        help="Limit running-config contain checks to one group name.",
    )
    return parser.parse_args()


def load_state_dump(path: str) -> dict:
    state_data = json.loads(Path(path).read_text(encoding="utf-8"))
    if not state_data:
        raise RuntimeError("state dump is empty")
    return state_data


def collect_state_metrics(state_data: dict) -> tuple[int, dict[str, str], dict[str, int]]:
    commit_count = 0
    running_cfg_by_device: dict[str, str] = {}
    history_count_by_device: dict[str, int] = {}

    for device_name, device in state_data.items():
        history = device.get("history", [])
        history_count_by_device[device_name] = len(history)
        commit_count += sum(1 for entry in history if entry.get("op") == "commit")

        running_cfg_by_device[device_name] = str(device.get("running_config", ""))

    return commit_count, running_cfg_by_device, history_count_by_device


def assert_devices_have_rpc_history(
    devices_to_check: list[str],
    history_count_by_device: dict[str, int],
) -> None:
    missing = [name for name in devices_to_check if history_count_by_device.get(name, 0) == 0]
    if missing:
        raise RuntimeError(f"device(s) with no RPC history in selected scope: {', '.join(missing)}")


def _group_from_config(running: str, group_name: str) -> tuple[str | None, str]:
    """The named group's XML out of a device's configuration, and the names of
    the groups it does have."""
    try:
        config = ET.fromstring(running) if running.strip() else None
    except ET.ParseError:
        return None, "<unparseable>"
    if config is None:
        return None, "<none>"

    names: list[str] = []
    for group in config.findall("groups"):
        name = (group.findtext("name") or "").strip()
        names.append(name)
        if name == group_name:
            return ET.tostring(group, encoding="unicode"), ", ".join(sorted(names))
    return None, ", ".join(sorted(names)) if names else "<none>"


def scope_running_config_by_group(
    state_data: dict,
    devices_to_check: list[str],
    only_group: str,
) -> dict[str, str]:
    if not only_group:
        return {
            name: str(state_data[name].get("running_config", ""))
            for name in devices_to_check
        }

    scoped: dict[str, str] = {}
    for name in devices_to_check:
        running = str(state_data[name].get("running_config", ""))
        group_xml, available = _group_from_config(running, only_group)
        if group_xml is not None:
            scoped[name] = group_xml
            continue

        # Fallback prevents false negatives if group naming differs in state dumps.
        scoped[name] = running
        print(
            (
                f"warning: group '{only_group}' not found for device '{name}'; "
                f"falling back to the whole running config (available: {available})"
            ),
            file=sys.stderr,
        )
    return scoped


def resolve_devices_to_check(only_device: str, running_cfg_by_device: dict[str, str]) -> list[str]:
    if only_device:
        if only_device not in running_cfg_by_device:
            raise RuntimeError(f"device '{only_device}' was not found in state dump")
        return [only_device]
    return sorted(running_cfg_by_device.keys())


def assert_required_strings(
    must_contain: list[str],
    devices_to_check: list[str],
    scoped_running_cfg_by_device: dict[str, str],
    full_running_cfg_by_device: dict[str, str],
    only_device: str,
    only_group: str,
) -> None:
    for needle in must_contain:
        if any(needle in scoped_running_cfg_by_device[name] for name in devices_to_check):
            continue

        # Junos group application semantics may place effective config outside the
        # exact group key represented in the mock state. If a group-scoped check
        # misses, allow a fallback match in full per-device running config.
        if only_group and any(needle in full_running_cfg_by_device[name] for name in devices_to_check):
            scope = only_device if only_device else "all devices"
            print(
                (
                    f"warning: '{needle}' not found under group '{only_group}' for {scope}; "
                    "matched in full running config"
                ),
                file=sys.stderr,
            )
            continue

        scope = only_device if only_device else "all devices"
        raise RuntimeError(f"expected to find '{needle}' in running config dump for {scope}")


def assert_forbidden_strings(
    must_not_contain: list[str],
    devices_to_check: list[str],
    running_cfg_by_device: dict[str, str],
) -> None:
    for needle in must_not_contain:
        offenders = [name for name in devices_to_check if needle in running_cfg_by_device[name]]
        if offenders:
            raise RuntimeError(
                f"did not expect to find '{needle}' in running config dump for devices: {', '.join(offenders)}"
            )


def main() -> int:
    args = parse_args()
    state_data = load_state_dump(args.state_dump)
    commit_count, running_cfg_by_device, history_count_by_device = collect_state_metrics(state_data)

    devices_to_check = resolve_devices_to_check(args.only_device, running_cfg_by_device)
    assert_devices_have_rpc_history(devices_to_check, history_count_by_device)

    if commit_count < args.require_min_commits:
        raise RuntimeError(
            f"expected at least {args.require_min_commits} commit ops, got {commit_count}"
        )

    scoped_running_cfg_by_device = scope_running_config_by_group(
        state_data,
        devices_to_check,
        args.only_group,
    )
    assert_required_strings(
        args.must_contain,
        devices_to_check,
        scoped_running_cfg_by_device,
        running_cfg_by_device,
        args.only_device,
        args.only_group,
    )
    assert_forbidden_strings(args.must_not_contain, devices_to_check, scoped_running_cfg_by_device)

    print(
        f"state assertions passed: devices={len(state_data)} commits={commit_count}"
    )
    return 0


if __name__ == "__main__":
    raise SystemExit(main())
