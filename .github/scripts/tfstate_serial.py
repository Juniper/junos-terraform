#!/usr/bin/env python3
"""Read Terraform state JSON from stdin and print its serial value."""

from __future__ import annotations

import json
import sys


def main() -> int:
    state = json.load(sys.stdin)
    serial = state.get("serial")
    if serial is None:
        raise RuntimeError("terraform state JSON does not include 'serial'")
    print(serial)
    return 0


if __name__ == "__main__":
    raise SystemExit(main())
