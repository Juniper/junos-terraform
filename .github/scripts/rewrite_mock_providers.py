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

    block_re = re.compile(
        r'provider\s+"junos-vqfx-evpn-vxlan"\s*\{.*?\}',
        re.DOTALL,
    )
    host_re = re.compile(r'host\s*=\s*"([^"]+)"')
    port_re = re.compile(r'port\s*=\s*\d+')

    hosts: list[str] = []
    for block in block_re.findall(text):
        hm = host_re.search(block)
        if hm:
            host = hm.group(1)
            if host not in hosts:
                hosts.append(host)

    host_port = {host: args.base_port + i for i, host in enumerate(hosts)}

    def rewrite_block(block: str) -> str:
        hm = host_re.search(block)
        if not hm:
            return block
        host = hm.group(1)
        port = host_port[host]
        block = host_re.sub(f'host     = "{args.bind_host}"', block)
        block = port_re.sub(f"port     = {port}", block)
        return block

    rewritten = block_re.sub(lambda m: rewrite_block(m.group(0)), text)
    providers_path.write_text(rewritten, encoding="utf-8")

    devices_path.parent.mkdir(parents=True, exist_ok=True)
    with devices_path.open("w", encoding="utf-8") as f:
        for host in hosts:
            f.write(f"{host}:{host_port[host]}\n")

    return 0


if __name__ == "__main__":
    raise SystemExit(main())
