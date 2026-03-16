#!/usr/bin/env python3
"""Apply an out-of-band NETCONF drift change to the mock device."""

from __future__ import annotations

import argparse
import asyncio
from pathlib import Path

import asyncssh

MSG_SEP = "]]>]]>"


def parse_args() -> argparse.Namespace:
    parser = argparse.ArgumentParser(
        description="Apply an out-of-band NETCONF drift change to one mock device/group."
    )
    parser.add_argument(
        "--devices-file",
        required=True,
        help="Path to host:port mapping file (<host>:<port> per line).",
    )
    parser.add_argument(
        "--target-host",
        required=True,
        help="Target device hostname as listed in devices-file.",
    )
    parser.add_argument(
        "--target-group",
        required=True,
        help="Apply-group name to mutate.",
    )
    parser.add_argument(
        "--drift-ip",
        required=True,
        help="IPv4 CIDR value to inject as drift.",
    )
    parser.add_argument("--username", required=True, help="NETCONF username for mock auth.")
    parser.add_argument("--password", required=True, help="NETCONF password for mock auth.")
    parser.add_argument(
        "--connect-host",
        default="127.0.0.1",
        help="Address to connect to mock listeners.",
    )
    return parser.parse_args()


async def recv_frame(reader) -> str:
    buf = ""
    while MSG_SEP not in buf:
        chunk = await reader.read(4096)
        if not chunk:
            raise RuntimeError(
                "connection closed while waiting for NETCONF frame"
            )
        buf += chunk
    frame, _ = buf.split(MSG_SEP, 1)
    return frame.strip()


async def send_rpc(writer, reader, message_id: str, payload: str) -> str:
    rpc = "".join(
        [
            '<?xml version="1.0" encoding="UTF-8"?>',
            (
                f'<rpc message-id="{message_id}" '
                'xmlns="urn:ietf:params:xml:ns:netconf:base:1.0">'
            ),
            payload,
            "</rpc>",
        ]
    )
    writer.write(rpc + MSG_SEP + "\n")
    return await recv_frame(reader)


def read_host_port_map(devices_file: Path) -> dict[str, int]:
    host_port: dict[str, int] = {}
    for raw_line in devices_file.read_text(encoding="utf-8").splitlines():
        line = raw_line.strip()
        if not line:
            continue
        name, port = line.split(":", 1)
        host_port[name] = int(port)
    return host_port


async def run(args: argparse.Namespace) -> None:
    host_port = read_host_port_map(Path(args.devices_file))

    if args.target_host not in host_port:
        raise RuntimeError(
            (
                f"host {args.target_host} from state not found "
                "in mock-devices mapping"
            )
        )

    conn = await asyncssh.connect(
        args.connect_host,
        port=host_port[args.target_host],
        username=args.username,
        password=args.password,
        known_hosts=None,
    )
    proc = None
    try:
        proc = await conn.create_process(subsystem="netconf")
        reader = proc.stdout
        writer = proc.stdin

        # Consume server hello and send client hello.
        await recv_frame(reader)
        hello = "".join(
            [
                '<?xml version="1.0" encoding="UTF-8"?>',
                '<hello xmlns="urn:ietf:params:xml:ns:netconf:base:1.0">',
                "<capabilities><capability>",
                "urn:ietf:params:netconf:base:1.0",
                "</capability></capabilities>",
                "</hello>",
            ]
        )
        writer.write(hello + MSG_SEP + "\n")

        # Out-of-band change: alter interface IP in the managed apply-group.
        drift_cfg = (
            '<load-configuration action="merge" format="xml">'
            "<configuration><groups>"
            f"<name>{args.target_group}</name>"
            "<interfaces><interface><name>lo0</name><unit><name>0</name>"
            f"<family><inet><address><name>{args.drift_ip}</name>"
            "</address></inet></family>"
            "</unit></interface></interfaces>"
            "</groups></configuration>"
            "</load-configuration>"
        )

        for mid, payload in [
            ("1001", "<lock><target><candidate/></target></lock>"),
            ("1002", drift_cfg),
            ("1003", "<commit/>"),
            ("1004", "<unlock><target><candidate/></target></unlock>"),
        ]:
            reply = await send_rpc(writer, reader, mid, payload)
            if "<ok/>" not in reply:
                raise RuntimeError(
                    f"unexpected reply for message-id {mid}: {reply}"
                )
    finally:
        if proc is not None:
            try:
                proc.terminate()
            except Exception:
                pass
        conn.close()
        await conn.wait_closed()


def main() -> int:
    args = parse_args()
    asyncio.run(run(args))
    return 0


if __name__ == "__main__":
    raise SystemExit(main())
