#!/usr/bin/env python3
"""Mutate the first IPv4 CIDR found under YAML files."""

from __future__ import annotations

import argparse
import re
from pathlib import Path


CIDR_RE = re.compile(r"\b\d{1,3}(?:\.\d{1,3}){3}/\d{1,2}\b")


def parse_args() -> argparse.Namespace:
    parser = argparse.ArgumentParser()
    parser.add_argument("--directory", default=".", help="Directory containing YAML files")
    parser.add_argument("--replacement", required=True, help="CIDR replacement value")
    return parser.parse_args()


def main() -> int:
    args = parse_args()
    root = Path(args.directory)
    yaml_files = sorted(list(root.glob("*.yaml")) + list(root.glob("*.yml")))

    for yaml_file in yaml_files:
        text = yaml_file.read_text(encoding="utf-8")
        match = CIDR_RE.search(text)
        if not match:
            continue

        original = match.group(0)
        updated = text[: match.start()] + args.replacement + text[match.end() :]
        yaml_file.write_text(updated, encoding="utf-8")
        print(f"mutated {yaml_file.name}: {original} -> {args.replacement}")
        return 0

    raise RuntimeError(f"no IPv4 CIDR found in YAML files under {root}")


if __name__ == "__main__":
    raise SystemExit(main())
