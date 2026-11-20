#!/usr/bin/env python3
"""Fetch a mock device's committed configuration over NETCONF.

The mock writes its state dump on shutdown, so a check that runs while the
lifecycle is still going has to ask the device, which is what a real check
would do anyway.
"""

from __future__ import annotations

import argparse
import asyncio
import sys
from pathlib import Path

import asyncssh

MSG_SEP = "]]>]]>"


def parse_args() -> argparse.Namespace:
    parser = argparse.ArgumentParser(description=__doc__)
    parser.add_argument("--devices-file", required=True,
                        help="Path to host:port mapping file (<host>:<port> per line).")
    parser.add_argument("--target-host", required=True,
                        help="Device hostname as listed in devices-file.")
    parser.add_argument("--username", required=True)
    parser.add_argument("--password", required=True)
    parser.add_argument("--connect-host", default="127.0.0.1")
    parser.add_argument("--output", default="-",
                        help="Where to write the configuration XML, or - for stdout.")
    return parser.parse_args()


def read_host_port_map(devices_file: Path) -> dict[str, int]:
    mapping: dict[str, int] = {}
    for line in devices_file.read_text(encoding="utf-8").splitlines():
        line = line.strip()
        if not line:
            continue
        host, _, port = line.rpartition(":")
        mapping[host] = int(port)
    return mapping


async def recv_frame(reader) -> str:
    chunks: list[str] = []
    while True:
        chunk = await reader.read(4096)
        if not chunk:
            break
        chunks.append(chunk)
        if MSG_SEP in "".join(chunks):
            break
    return "".join(chunks).replace(MSG_SEP, "")


async def run(args: argparse.Namespace) -> str:
    ports = read_host_port_map(Path(args.devices_file))
    if args.target_host not in ports:
        raise SystemExit(f"{args.target_host} is not in {args.devices_file}")

    conn = await asyncssh.connect(
        args.connect_host,
        port=ports[args.target_host],
        username=args.username,
        password=args.password,
        known_hosts=None,
    )
    try:
        proc = await conn.create_process(subsystem="netconf")
        await recv_frame(proc.stdout)
        proc.stdin.write(
            '<hello xmlns="urn:ietf:params:xml:ns:netconf:base:1.0">'
            "<capabilities><capability>urn:ietf:params:netconf:base:1.0"
            "</capability></capabilities></hello>" + MSG_SEP + "\n"
        )
        proc.stdin.write(
            '<rpc xmlns="urn:ietf:params:xml:ns:netconf:base:1.0" message-id="read">'
            "<get-configuration><configuration></configuration></get-configuration>"
            "</rpc>" + MSG_SEP + "\n"
        )
        return await recv_frame(proc.stdout)
    finally:
        conn.close()
        await conn.wait_closed()


def main() -> int:
    args = parse_args()
    reply = asyncio.run(run(args))
    start = reply.find("<configuration")
    end = reply.rfind("</configuration>")
    config = reply[start:end + len("</configuration>")] if start != -1 and end != -1 else reply
    if args.output == "-":
        sys.stdout.write(config)
    else:
        Path(args.output).write_text(config, encoding="utf-8")
    return 0


if __name__ == "__main__":
    raise SystemExit(main())
