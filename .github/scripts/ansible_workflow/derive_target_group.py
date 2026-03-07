#!/usr/bin/env python3
"""Derive the first apply-group name from rendered Junos XML config."""

from __future__ import annotations

import argparse
import re
from pathlib import Path


def parse_args() -> argparse.Namespace:
    parser = argparse.ArgumentParser(
        description="Extract the first <groups><name> value from rendered host XML."
    )
    parser.add_argument("--config-xml", required=True, help="Path to rendered host XML file.")
    return parser.parse_args()


def main() -> int:
    args = parse_args()
    xml_path = Path(args.config_xml)
    text = xml_path.read_text(encoding="utf-8")
    match = re.search(r"<groups>\s*<name>([^<]+)</name>", text)
    if not match:
        raise RuntimeError(f"failed to derive target group from rendered XML: {xml_path}")
    print(match.group(1).strip())
    return 0


if __name__ == "__main__":
    raise SystemExit(main())
