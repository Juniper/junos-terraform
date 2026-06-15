#!/usr/bin/env python3
"""Rewrite generated provider blocks to target local NETCONF mock listeners."""

from __future__ import annotations

import argparse
import re
from pathlib import Path


def parse_args() -> argparse.Namespace:
    parser = argparse.ArgumentParser(
        description="Rewrite generated Terraform providers to point at NETCONF mock listeners."
    )
    parser.add_argument(
        "--providers-file",
        required=True,
        help="Path to generated providers.tf.",
    )
    parser.add_argument(
        "--devices-file",
        required=True,
        help="Output file for host:port mapping used by netconf mock.",
    )
    parser.add_argument(
        "--bind-host",
        default="127.0.0.1",
        help="Host value written into provider blocks.",
    )
    parser.add_argument(
        "--base-port",
        type=int,
        default=8301,
        help="First mock listener port.",
    )
    return parser.parse_args()


def main() -> int:
    args = parse_args()

    providers_path = Path(args.providers_file)
    devices_path = Path(args.devices_file)

    text = providers_path.read_text(encoding="utf-8")

    # Match any variant of the junos-vqfx-evpn-vxlan provider (e.g. -trim suffix).
    block_re = re.compile(
        r'provider\s+"junos-vqfx-evpn-vxlan[^"]*"\s*\{.*?\}',
        re.DOTALL,
    )
    alias_re = re.compile(r'alias\s*=\s*"([^"]+)"')
    host_re = re.compile(r'(host\s*=\s*)"([^"]+)"')
    port_re = re.compile(r'(port\s*=\s*)\d+')

    # Collect aliases in order of appearance; each gets a unique port.
    aliases: list[str] = []
    for block in block_re.findall(text):
        am = alias_re.search(block)
        if am:
            alias = am.group(1)
            if alias not in aliases:
                aliases.append(alias)

    alias_port = {alias: args.base_port + i for i, alias in enumerate(aliases)}
    # devices file maps device-name (alias with underscores → hyphens) to port
    devices: list[tuple[str, int]] = [
        (alias.replace("_", "-"), port) for alias, port in alias_port.items()
    ]

    def rewrite_block(block: str) -> str:
        am = alias_re.search(block)
        if not am:
            return block
        alias = am.group(1)
        port = alias_port[alias]
        block = host_re.sub(lambda m: f'{m.group(1)}"{args.bind_host}"', block)
        block = port_re.sub(lambda m: f"{m.group(1)}{port}", block)
        return block

    rewritten = block_re.sub(lambda m: rewrite_block(m.group(0)), text)
    providers_path.write_text(rewritten, encoding="utf-8")

    devices_path.parent.mkdir(parents=True, exist_ok=True)
    with devices_path.open("w", encoding="utf-8") as f:
        for device_name, port in devices:
            f.write(f"{device_name}:{port}\n")

    return 0


if __name__ == "__main__":
    raise SystemExit(main())
