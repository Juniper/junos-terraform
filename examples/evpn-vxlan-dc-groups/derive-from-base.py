#!/usr/bin/env python3
"""Derive the group-based example from the base-configuration example.

Each device keeps its identity in the base hierarchy — host name and the
apply-groups reference — and everything else moves into one group named after
the device's role, which is how a group-based network is usually organised.

Run from the repository root; the output is committed, so this only needs
running when examples/evpn-vxlan-dc changes.
"""
import pathlib
import sys
import xml.etree.ElementTree as ET

SOURCE = pathlib.Path("examples/evpn-vxlan-dc")
TARGET = pathlib.Path("examples/evpn-vxlan-dc-groups")

ROLES = ("borderleaf", "leaf", "spine", "firewall", "pe")


def role_of(name: str) -> str:
    for role in ROLES:
        if role in name:
            return role
    raise SystemExit(f"no role in {name}")


def indent(elem: ET.Element, level: int = 0) -> None:
    pad = "\n" + "    " * level
    if len(elem):
        if not (elem.text or "").strip():
            elem.text = pad + "    "
        for child in elem:
            indent(child, level + 1)
        if not (elem.tail or "").strip():
            elem.tail = pad
        if not (elem[-1].tail or "").strip():
            elem[-1].tail = pad
    elif level and not (elem.tail or "").strip():
        elem.tail = pad


def convert(path: pathlib.Path) -> str:
    device = path.stem
    role = role_of(device)

    source_root = ET.parse(path).getroot()
    config = source_root.find("configuration")
    if config is None:
        raise SystemExit(f"{path}: no configuration element")

    out_config = ET.Element("configuration")

    apply_groups = ET.SubElement(out_config, "apply-groups")
    apply_groups.text = role

    groups = ET.SubElement(out_config, "groups")
    ET.SubElement(groups, "name").text = role

    base_system = ET.SubElement(out_config, "system")

    for child in list(config):
        if child.tag == "version":
            continue
        if child.tag == "system":
            group_system = ET.SubElement(groups, "system")
            for leaf in list(child):
                # Host name identifies the device, so it stays where a reader
                # of the device's own configuration expects to find it.
                if leaf.tag == "host-name":
                    base_system.append(leaf)
                else:
                    group_system.append(leaf)
            continue
        groups.append(child)

    indent(out_config)
    body = ET.tostring(out_config, encoding="unicode")
    return f'<rpc-reply xmlns:junos="http://xml.juniper.net/junos/18.1R3/junos">\n{body}\n</rpc-reply>\n'


def main() -> int:
    if not SOURCE.is_dir():
        raise SystemExit(f"{SOURCE} not found; run from the repository root")

    written = 0
    for path in sorted(SOURCE.rglob("*.xml")):
        out_path = TARGET / path.relative_to(SOURCE)
        out_path.parent.mkdir(parents=True, exist_ok=True)
        out_path.write_text(convert(path), encoding="utf-8")
        written += 1
    print(f"wrote {written} files to {TARGET}")
    return 0


if __name__ == "__main__":
    sys.exit(main())
