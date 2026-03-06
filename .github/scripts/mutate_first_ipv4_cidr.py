#!/usr/bin/env python3
"""Mutate the first IPv4 CIDR found in Terraform files."""

from __future__ import annotations

import argparse
import re
from pathlib import Path


CIDR_RE = re.compile(r"\b\d{1,3}(?:\.\d{1,3}){3}/\d{1,2}\b")


def parse_args() -> argparse.Namespace:
    parser = argparse.ArgumentParser()
    parser.add_argument(
        "--directory",
        default=".",
        help="Directory containing .tf files",
    )
    parser.add_argument(
        "--replacement",
        required=True,
        help="CIDR replacement value",
    )
    parser.add_argument(
        "--exclude",
        action="append",
        default=[],
        help="Terraform filename to skip; can be specified multiple times",
    )
    return parser.parse_args()


def main() -> int:
    args = parse_args()
    tf_dir = Path(args.directory)
    excluded = set(args.exclude)

    candidates = sorted(
        p for p in tf_dir.glob("*.tf") if p.name not in excluded
    )

    for tf_file in candidates:
        text = tf_file.read_text(encoding="utf-8")
        match = CIDR_RE.search(text)
        if not match:
            continue

        original = match.group(0)
        updated = text[:match.start()] + args.replacement + text[match.end():]
        tf_file.write_text(updated, encoding="utf-8")
        print(f"mutated {tf_file.name}: {original} -> {args.replacement}")
        return 0

    raise RuntimeError(f"no IPv4 CIDR found in terraform files under {tf_dir}")


if __name__ == "__main__":
    raise SystemExit(main())
