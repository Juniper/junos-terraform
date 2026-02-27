import json
import os
from glob import glob
import subprocess
import shutil
import sys
import unittest
from subprocess import CalledProcessError

class TestWorkflow(unittest.TestCase):

    def test_workflow_basic(self):
        self.assertTrue(True)
        
# Note for any changes need to rerun "pip install ./junos-terraform"

def test_yang2go():
    repo_root = os.path.abspath(os.path.join(os.path.dirname(__file__), "..", ".."))
    yang_root = os.path.abspath(os.path.join(repo_root, "examples", "yang"))

    print("repo_root:", repo_root)
    print("yang_root:", yang_root)
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

    # Test generated provider with trimmed_schema.json
    generated_provider_dir = os.path.join(
        repo_root, "terraform-provider-junos-vqfx-evpn-vxlan"
    )
    generated_trimmed_schema = os.path.join(
        generated_provider_dir, "trimmed_schema.json"
    )

    # Remove any existing provider dir before running
    if os.path.exists(generated_provider_dir):
        shutil.rmtree(generated_provider_dir)

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

    print("CMD:", cmd)
    try:
        proc = subprocess.run(
            cmd,
            input=stdin_json,
            text=True,
            capture_output=True,
            check=True,
            cwd=repo_root,
            env=env
        )
    except CalledProcessError as e:
        # Debug output
        print("RETURNCODE:", e.returncode)
        print("STDOUT:\n", e.output)
        print("STDERR:\n", e.stderr)
        raise(e)

    # Debug output
    print("RETURNCODE:", proc.returncode)
    print("STDOUT:\n", proc.stdout)
    print("STDERR:\n", proc.stderr)

    assert proc.returncode == 0, (
        f"jtaf-yang2go failed:\nSTDOUT:\n{proc.stdout}\n\nSTDERR:\n{proc.stderr}"
    )

    assert os.path.isdir(generated_provider_dir), (
        f"Expected provider dir not created: {generated_provider_dir}"
    )
    assert os.path.exists(generated_trimmed_schema), (
        f"Expected trimmed_schema.json not found at {generated_trimmed_schema}"
    )

    with open(generated_trimmed_schema, "r") as f:
        generated_json = json.load(f)

    # Compare against expected trimmed_schema.json in examples/providers
    expected_trimmed_schema = os.path.join(
        repo_root,
        "examples",
        "providers",
        "terraform-provider-junos-vqfx-evpn-vxlan",
        "trimmed_schema.json",
    )
    assert os.path.exists(expected_trimmed_schema), (
        f"Expected example trimmed_schema.json not found at {expected_trimmed_schema}"
    )

    with open(expected_trimmed_schema, "r") as f:
        expected_json = json.load(f)

    assert generated_json == expected_json, "Generated trimmed_schema.json differs from expected example"

def test_yang2ansible():

    repo_root = os.path.abspath(os.path.join(os.path.dirname(__file__), "..", ".."))
    yang_root = os.path.abspath(os.path.join(repo_root, "examples", "yang"))

    print("repo_root:", repo_root)
    print("yang_root:", yang_root)
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

    ansible_dir = os.path.join(
        repo_root, "examples/ansible"
    )

    # Remove any existing ansible role dir before running
    if os.path.exists(ansible_dir+"/ansible-provider-junos-vqfx-ansible-role"):
        shutil.rmtree(ansible_dir+"/ansible-provider-junos-vqfx-ansible-role")

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

    print("CMD:", cmd)

    try:
        proc = subprocess.run(
            cmd,
            input=stdin_json,
            text=True,
            capture_output=True,
            check=True,
            cwd=ansible_dir,
            env=env
        )
    except CalledProcessError as e:
        # Debug output
        print("RETURNCODE:", e.returncode)
        print("STDOUT:\n", e.output)
        print("STDERR:\n", e.stderr)
        raise(e)

    # Debug output
    print("RETURNCODE:", proc.returncode)
    print("STDOUT:\n", proc.stdout)
    print("STDERR:\n", proc.stderr)

    print("Completed jtaf-yang2ansible execution")
    assert proc.returncode == 0, (
        f"jtaf-yang2ansible failed:\nSTDOUT:\n{proc.stdout}\n\nSTDERR:\n{proc.stderr}"
    )

    assert os.path.isdir(os.path.join(ansible_dir, "ansible-provider-junos-vqfx-ansible-role")), (
        f"Expected ansible roles dir was not created: {ansible_dir}+/ansible-provider-junos-vqfx-ansible-role"
    )
    
    trimmed_schema_path = os.path.join(
        ansible_dir, "ansible-provider-junos-vqfx-ansible-role", "trimmed_schema.json"
    )

    # Remove any existing ansible files dir before running
    if os.path.exists(ansible_dir+"/vqfx_ansible_files"):
        shutil.rmtree(ansible_dir+"/vqfx_ansible_files")

    stdin_json = "{}"

    # xml2yaml command
    cmd = [
        'jtaf-xml2yaml',
        "-j",
        trimmed_schema_path,
        "-x",
        *xml_args,
        "-d",
        "vqfx_ansible_files",
    ]
    print("CMD:", cmd)
    try:
        proc = subprocess.run(
            cmd,
            input=stdin_json,
            text=True,
            capture_output=True,
            check=True,
            cwd=ansible_dir,
            env=env
        )
    except CalledProcessError as e:
        # Debug output
        print("RETURNCODE:", e.returncode)
        print("STDOUT:\n", e.output)
        print("STDERR:\n", e.stderr)
        raise(e)

    # Debug output
    print("RETURNCODE:", proc.returncode)
    print("STDOUT:\n", proc.stdout)
    print("STDERR:\n", proc.stderr)

    assert proc.returncode == 0, (
        f"jtaf-xml2yaml failed:\nSTDOUT:\n{proc.stdout}\n\nSTDERR:\n{proc.stderr}"
    )

    assert os.path.isdir(os.path.join(ansible_dir, "vqfx_ansible_files")), (
        f"Expected ansible files dir was not created: {ansible_dir}+/vqfx_ansible_files"
    )