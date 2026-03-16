#!/usr/bin/env python3
"""Read Terraform state JSON from stdin and print its serial value."""

from __future__ import annotations

import argparse
import json
import sys


def parse_args() -> argparse.Namespace:
    parser = argparse.ArgumentParser(
        description="Read Terraform state JSON from stdin and print the serial value."
    )
    return parser.parse_args()


def main() -> int:
    parse_args()
    state = json.load(sys.stdin)
    serial = state.get("serial")
    if serial is None:
        raise RuntimeError("terraform state JSON does not include 'serial'")
    print(serial)
    return 0


if __name__ == "__main__":
    raise SystemExit(main())
