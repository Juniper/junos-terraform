#!/usr/bin/env python3
"""Assert expected NETCONF mock state invariants from JSON dump."""

from __future__ import annotations

import argparse
import json
from pathlib import Path


def parse_args() -> argparse.Namespace:
    parser = argparse.ArgumentParser()
    parser.add_argument("--state-dump", required=True, help="Path to netconf mock state JSON")
    parser.add_argument("--require-min-commits", type=int, default=1, help="Minimum total commit operations")
    parser.add_argument("--must-contain", action="append", default=[], help="String expected in running config")
    parser.add_argument("--must-not-contain", action="append", default=[], help="String forbidden in running config")
    parser.add_argument(
        "--only-device",
        default="",
        help="Optional device name to scope must-contain and must-not-contain checks",
    )
    parser.add_argument(
        "--only-group",
        default="",
        help="Optional group name to scope running config checks within each device",
    )
    return parser.parse_args()


def load_state_dump(path: str) -> dict:
    state_data = json.loads(Path(path).read_text(encoding="utf-8"))
    if not state_data:
        raise RuntimeError("state dump is empty")
    return state_data


def collect_state_metrics(state_data: dict) -> tuple[int, dict[str, str]]:
    commit_count = 0
    running_cfg_by_device: dict[str, str] = {}

    for device_name, device in state_data.items():
        history = device.get("history", [])
        commit_count += sum(1 for entry in history if entry.get("op") == "commit")

        running_groups = device.get("running_groups", {})
        running_cfg_by_device[device_name] = "\n".join(str(v) for v in running_groups.values())

        if not history:
            raise RuntimeError(f"device {device_name} has no RPC history")

    return commit_count, running_cfg_by_device


def scope_running_config_by_group(
    state_data: dict,
    devices_to_check: list[str],
    only_group: str,
) -> dict[str, str]:
    if not only_group:
        return {
            name: "\n".join(str(v) for v in state_data[name].get("running_groups", {}).values())
            for name in devices_to_check
        }

    scoped: dict[str, str] = {}
    for name in devices_to_check:
        groups = state_data[name].get("running_groups", {})
        scoped[name] = str(groups.get(only_group, ""))
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
    running_cfg_by_device: dict[str, str],
    only_device: str,
) -> None:
    for needle in must_contain:
        if any(needle in running_cfg_by_device[name] for name in devices_to_check):
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
    commit_count, running_cfg_by_device = collect_state_metrics(state_data)

    if commit_count < args.require_min_commits:
        raise RuntimeError(
            f"expected at least {args.require_min_commits} commit ops, got {commit_count}"
        )

    devices_to_check = resolve_devices_to_check(args.only_device, running_cfg_by_device)
    running_cfg_by_device = scope_running_config_by_group(
        state_data,
        devices_to_check,
        args.only_group,
    )
    assert_required_strings(args.must_contain, devices_to_check, running_cfg_by_device, args.only_device)
    assert_forbidden_strings(args.must_not_contain, devices_to_check, running_cfg_by_device)

    print(
        f"state assertions passed: devices={len(state_data)} commits={commit_count}"
    )
    return 0


if __name__ == "__main__":
    raise SystemExit(main())
