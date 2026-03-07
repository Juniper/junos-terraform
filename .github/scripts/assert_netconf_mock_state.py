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
    return parser.parse_args()


def main() -> int:
    args = parse_args()
    state_data = json.loads(Path(args.state_dump).read_text(encoding="utf-8"))
    if not state_data:
        raise RuntimeError("state dump is empty")

    commit_count = 0
    running_cfg_chunks: list[str] = []

    for device_name, device in state_data.items():
        history = device.get("history", [])
        commit_count += sum(1 for entry in history if entry.get("op") == "commit")

        running_groups = device.get("running_groups", {})
        running_cfg_chunks.extend(str(v) for v in running_groups.values())

        if not history:
            raise RuntimeError(f"device {device_name} has no RPC history")

    if commit_count < args.require_min_commits:
        raise RuntimeError(
            f"expected at least {args.require_min_commits} commit ops, got {commit_count}"
        )

    running_cfg_text = "\n".join(running_cfg_chunks)

    for needle in args.must_contain:
        if needle not in running_cfg_text:
            raise RuntimeError(f"expected to find '{needle}' in running config dump")

    for needle in args.must_not_contain:
        if needle in running_cfg_text:
            raise RuntimeError(f"did not expect to find '{needle}' in running config dump")

    print(
        f"state assertions passed: devices={len(state_data)} commits={commit_count}"
    )
    return 0


if __name__ == "__main__":
    raise SystemExit(main())
