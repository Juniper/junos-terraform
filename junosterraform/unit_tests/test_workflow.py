import json
import os
from glob import glob
import subprocess
import shutil
import tempfile
import unittest


class TestWorkflow(unittest.TestCase):

    def test_workflow_basic(self):
        self.assertTrue(True)


# Note for any changes need to rerun "pip install ./junos-terraform"


def test_yang2go():
    repo_root = os.path.abspath(os.path.join(os.path.dirname(__file__), "..", ".."))
    yang_root = os.path.abspath(os.path.join(repo_root, "examples", "yang"))

    assert os.path.isdir(yang_root), f"YANG root does not exist: {yang_root}"

    exe = shutil.which("jtaf-yang2go")
    assert exe, "Could not find jtaf-yang2go on PATH"

    env = os.environ.copy()

    # Building full paths to YANG dirs / files
    common_dir = os.path.join(yang_root, "18.2", "18.2R3", "common")
    conf_glob = os.path.join(yang_root, "18.2", "18.2R3", "junos-qfx", "conf", "*.yang")
    conf_files = sorted(glob(conf_glob))

    assert os.path.isdir(common_dir), f"common_dir does not exist: {common_dir}"
    assert conf_files, f"No .yang files found under {conf_glob}"

    # Building list of all XMLs to pass in one command
    rel_xml_files = [
        "examples/evpn-vxlan-dc/dc1/dc1-borderleaf1.xml",
        "examples/evpn-vxlan-dc/dc1/dc1-borderleaf2.xml",
        "examples/evpn-vxlan-dc/dc1/dc1-leaf1.xml",
        "examples/evpn-vxlan-dc/dc1/dc1-leaf2.xml",
        "examples/evpn-vxlan-dc/dc1/dc1-leaf3.xml",
        "examples/evpn-vxlan-dc/dc1/dc1-spine1.xml",
        "examples/evpn-vxlan-dc/dc1/dc1-spine2.xml",
        "examples/evpn-vxlan-dc/dc2/dc2-spine1.xml",
        "examples/evpn-vxlan-dc/dc2/dc2-spine2.xml",
    ]
    xml_args = [os.path.join(repo_root, p) for p in rel_xml_files]

    for path in xml_args:
        assert os.path.exists(path), f"XML file does not exist: {path}"

    stdin_json = "{}"

    # yang2go command
    cmd = [
        exe,
        "-p",
        common_dir,
        *conf_files,
        "-x",
        *xml_args,
        "-t",
        "vqfx-evpn-vxlan",
    ]

    with tempfile.TemporaryDirectory(prefix="jtaf-yang2go-") as tmpdir:
        # Test generated provider with trimmed_schema.json in isolated temp workspace.
        generated_provider_dir = os.path.join(
            tmpdir, "terraform-provider-junos-vqfx-evpn-vxlan"
        )
        generated_trimmed_schema = os.path.join(
            generated_provider_dir, "trimmed_schema.json"
        )

        proc = subprocess.run(
            cmd,
            input=stdin_json,
            text=True,
            capture_output=True,
            check=False,
            cwd=tmpdir,
            env=env
        )

        assert proc.returncode == 0, (
            f"jtaf-yang2go failed:\nSTDOUT:\n{proc.stdout}\n\nSTDERR:\n{proc.stderr}"
        )

        assert os.path.isdir(generated_provider_dir), (
            f"Expected provider dir not created: {generated_provider_dir}"
        )
        assert os.path.exists(generated_trimmed_schema), (
            f"Expected trimmed_schema.json not found at {generated_trimmed_schema}"
        )

        with open(generated_trimmed_schema) as f:
            generated_json = json.load(f)

    # Validate generated schema shape without depending on committed generated fixtures.
    assert isinstance(generated_json, dict), "Generated trimmed_schema.json should be a JSON object"
    assert "root" in generated_json, "Generated schema missing 'root' key"
    assert isinstance(generated_json["root"], dict), "Generated schema 'root' should be an object"
    assert "children" in generated_json["root"], "Generated schema root missing 'children' key"
    assert generated_json["root"]["children"], "Generated schema root children should not be empty"


def test_yang2ansible():

    repo_root = os.path.abspath(os.path.join(os.path.dirname(__file__), "..", ".."))
    yang_root = os.path.abspath(os.path.join(repo_root, "examples", "yang"))

    assert os.path.isdir(yang_root), f"YANG root does not exist: {yang_root}"

    exe = shutil.which("jtaf-yang2ansible")
    assert exe, "Could not find jtaf-yang2ansible on PATH"

    env = os.environ.copy()

    # Building full paths to YANG dirs / files
    common_dir = os.path.join(yang_root, "18.2", "18.2R3", "common")
    conf_glob = os.path.join(yang_root, "18.2", "18.2R3", "junos-qfx", "conf", "*.yang")
    conf_files = sorted(glob(conf_glob))

    assert os.path.isdir(common_dir), f"common_dir does not exist: {common_dir}"
    assert conf_files, f"No .yang files found under {conf_glob}"

    # Building list of all XMLs to pass in one command
    rel_xml_files = [
        "examples/evpn-vxlan-dc/dc1/dc1-borderleaf1.xml",
        "examples/evpn-vxlan-dc/dc1/dc1-borderleaf2.xml",
        "examples/evpn-vxlan-dc/dc1/dc1-leaf1.xml",
        "examples/evpn-vxlan-dc/dc1/dc1-leaf2.xml",
        "examples/evpn-vxlan-dc/dc1/dc1-leaf3.xml",
        "examples/evpn-vxlan-dc/dc1/dc1-spine1.xml",
        "examples/evpn-vxlan-dc/dc1/dc1-spine2.xml",
        "examples/evpn-vxlan-dc/dc2/dc2-spine1.xml",
        "examples/evpn-vxlan-dc/dc2/dc2-spine2.xml",
    ]
    xml_args = [os.path.join(repo_root, p) for p in rel_xml_files]

    for path in xml_args:
        assert os.path.exists(path), f"XML file does not exist: {path}"

    stdin_json = "{}"

    # yang2ansible command
    cmd = [
        exe,
        "-p",
        common_dir,
        *conf_files,
        "-x",
        *xml_args,
        "-t",
        "vqfx-ansible-role",
    ]

    with tempfile.TemporaryDirectory(prefix="jtaf-yang2ansible-") as ansible_dir:
        proc = subprocess.run(
            cmd,
            input=stdin_json,
            text=True,
            capture_output=True,
            check=False,
            cwd=ansible_dir,
            env=env
        )

        assert proc.returncode == 0, (
            f"jtaf-yang2ansible failed:\nSTDOUT:\n{proc.stdout}\n\nSTDERR:\n{proc.stderr}"
        )

        role_dir = os.path.join(ansible_dir, "ansible-provider-junos-vqfx-ansible-role")
        assert os.path.isdir(role_dir), (
            f"Expected ansible roles dir was not created: {role_dir}"
        )

        trimmed_schema_path = os.path.join(
            role_dir, "trimmed_schema.json"
        )

        # xml2yaml command
        cmd = [
            'jtaf-xml2yaml',
            "-j",
            trimmed_schema_path,
            "-x",
            *xml_args,
            "-d",
            "vqfx_ansible_files",
        ]  # noqa: E501
        proc = subprocess.run(
            cmd,
            input=stdin_json,
            text=True,
            capture_output=True,
            check=False,
            cwd=ansible_dir,
            env=env
        )

        assert proc.returncode == 0, (
            f"jtaf-xml2yaml failed:\nSTDOUT:\n{proc.stdout}\n\nSTDERR:\n{proc.stderr}"
        )

        ansible_files_dir = os.path.join(ansible_dir, "vqfx_ansible_files")
        assert os.path.isdir(ansible_files_dir), (
            f"Expected ansible files dir was not created: {ansible_files_dir}"
        )
