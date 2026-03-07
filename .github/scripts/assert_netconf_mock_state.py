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
    return parser.parse_args()


def main() -> int:
    args = parse_args()
    state_data = json.loads(Path(args.state_dump).read_text(encoding="utf-8"))
    if not state_data:
        raise RuntimeError("state dump is empty")

    commit_count = 0
    running_cfg_by_device: dict[str, str] = {}

    for device_name, device in state_data.items():
        history = device.get("history", [])
        commit_count += sum(1 for entry in history if entry.get("op") == "commit")

        running_groups = device.get("running_groups", {})
        running_cfg_by_device[device_name] = "\n".join(str(v) for v in running_groups.values())

        if not history:
            raise RuntimeError(f"device {device_name} has no RPC history")

    if commit_count < args.require_min_commits:
        raise RuntimeError(
            f"expected at least {args.require_min_commits} commit ops, got {commit_count}"
        )

    if args.only_device:
        if args.only_device not in running_cfg_by_device:
            raise RuntimeError(f"device '{args.only_device}' was not found in state dump")
        devices_to_check = [args.only_device]
    else:
        devices_to_check = sorted(running_cfg_by_device.keys())

    for needle in args.must_contain:
        if not any(needle in running_cfg_by_device[name] for name in devices_to_check):
            scope = args.only_device if args.only_device else "all devices"
            raise RuntimeError(f"expected to find '{needle}' in running config dump for {scope}")

    for needle in args.must_not_contain:
        offenders = [name for name in devices_to_check if needle in running_cfg_by_device[name]]
        if offenders:
            raise RuntimeError(
                f"did not expect to find '{needle}' in running config dump for devices: {', '.join(offenders)}"
            )

    print(
        f"state assertions passed: devices={len(state_data)} commits={commit_count}"
    )
    return 0


if __name__ == "__main__":
    raise SystemExit(main())
