#!/usr/bin/env python3
"""Render CI NETCONF playbook from Jinja template using generated role name."""

from __future__ import annotations

import argparse
from pathlib import Path

from jinja2 import Template


def parse_args() -> argparse.Namespace:
    parser = argparse.ArgumentParser()
    parser.add_argument(
        "--roles-root",
        default="ansible-provider-junos-vqfx-evpn-vxlan/roles",
        help="Directory containing generated role subdirectories",
    )
    parser.add_argument(
        "--template-path",
        default=".github/ansible/ci-netconf-playbook.yml.j2",
        help="Path to Jinja2 playbook template",
    )
    parser.add_argument(
        "--output-path",
        default="ansible-provider-junos-vqfx-evpn-vxlan/ci-netconf-playbook.yml",
        help="Output path for rendered playbook",
    )
    return parser.parse_args()


def detect_generated_role_name(roles_root: Path) -> str:
    role_dirs = sorted(p.name for p in roles_root.iterdir() if p.is_dir())
    if not role_dirs:
        raise RuntimeError(f"failed to detect generated role under {roles_root}")
    return role_dirs[0]


def main() -> int:
    args = parse_args()
    roles_root = Path(args.roles_root)
    template_path = Path(args.template_path)
    output_path = Path(args.output_path)

    role_name = detect_generated_role_name(roles_root)
    template = Template(
        template_path.read_text(encoding="utf-8"),
        variable_start_string="[[",
        variable_end_string="]]",
    )
    rendered = template.render(generated_role_name=role_name)
    output_path.write_text(rendered, encoding="utf-8")
    print(f"rendered {output_path} with role={role_name}")
    return 0


if __name__ == "__main__":
    raise SystemExit(main())
