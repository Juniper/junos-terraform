#!/usr/bin/env python3
"""Build Ansible mock inventory and host:port mapping from a hosts file."""

from __future__ import annotations

import argparse
from pathlib import Path


def parse_args() -> argparse.Namespace:
    parser = argparse.ArgumentParser(
        description="Generate mock Ansible inventory and host:port device map from a hosts file."
    )
    parser.add_argument("--hosts-file", required=True, help="Source hosts file (Ansible-style host list).")
    parser.add_argument("--inventory-file", required=True, help="Output inventory file path.")
    parser.add_argument("--devices-file", required=True, help="Output host:port device map path.")
    parser.add_argument("--bind-host", default="127.0.0.1", help="Mock bind address written to inventory.")
    parser.add_argument("--base-port", type=int, default=8301, help="Starting NETCONF port for host mapping.")
    return parser.parse_args()


def parse_hosts(hosts_file: Path) -> list[str]:
    hosts: list[str] = []
    for raw_line in hosts_file.read_text(encoding="utf-8").splitlines():
        line = raw_line.strip()
        if not line or line.startswith("#"):
            continue
        if line.startswith("[") and line.endswith("]"):
            continue
        host = line.split()[0]
        if host not in hosts:
            hosts.append(host)
    if not hosts:
        raise RuntimeError(f"no hosts discovered in {hosts_file}")
    return hosts


def write_inventory(inventory_file: Path, hosts: list[str], bind_host: str, base_port: int) -> None:
    lines = ["[all]"]
    for idx, host in enumerate(hosts):
        lines.append(f"{host} ansible_host={bind_host} ansible_port={base_port + idx}")
    lines.append("")
    lines.append("[all:vars]")
    lines.append("ansible_connection=local")
    lines.append("ansible_python_interpreter=auto_silent")
    lines.append("")
    inventory_file.parent.mkdir(parents=True, exist_ok=True)
    inventory_file.write_text("\n".join(lines), encoding="utf-8")


def write_devices(devices_file: Path, hosts: list[str], base_port: int) -> None:
    lines = [f"{host}:{base_port + idx}" for idx, host in enumerate(hosts)]
    devices_file.parent.mkdir(parents=True, exist_ok=True)
    devices_file.write_text("\n".join(lines) + "\n", encoding="utf-8")


def main() -> int:
    args = parse_args()
    hosts = parse_hosts(Path(args.hosts_file))
    write_inventory(Path(args.inventory_file), hosts, args.bind_host, args.base_port)
    write_devices(Path(args.devices_file), hosts, args.base_port)
    print(f"prepared inventory for {len(hosts)} hosts")
    return 0


if __name__ == "__main__":
    raise SystemExit(main())
