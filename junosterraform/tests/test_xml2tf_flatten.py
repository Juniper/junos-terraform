import importlib.util
from importlib.machinery import SourceFileLoader
import os
import tempfile
import unittest


def _load_xml2tf():
    repo_root = os.path.abspath(os.path.join(os.path.dirname(__file__), "..", ".."))
    script_path = os.path.join(repo_root, "junosterraform", "jtaf-xml2tf")
    loader = SourceFileLoader("jtaf_xml2tf", script_path)
    spec = importlib.util.spec_from_loader("jtaf_xml2tf", loader)
    assert spec is not None
    module = importlib.util.module_from_spec(spec)
    assert spec.loader is not None
    spec.loader.exec_module(module)
    return module


class TestXMLToTerraformFlatten(unittest.TestCase):

    def test_applied_groups_are_flattened_into_base_config(self):
        repo_root = os.path.abspath(os.path.join(os.path.dirname(__file__), "..", ".."))
        script_path = os.path.join(repo_root, "junosterraform", "jtaf-xml2tf")

        loader = SourceFileLoader("jtaf_xml2tf", script_path)
        spec = importlib.util.spec_from_loader("jtaf_xml2tf", loader)
        assert spec is not None
        module = importlib.util.module_from_spec(spec)
        assert spec.loader is not None
        spec.loader.exec_module(module)

        type_lookup = {
            "system": {"type": "container"},
            "system/host_name": {"type": "leaf"},
            "system/services": {"type": "container"},
            "system/services/ssh": {"type": "container"},
        }

        xml_content = """<configuration>
  <apply-groups>overlay</apply-groups>
  <system>
    <host-name>router1</host-name>
  </system>
  <groups>
    <name>overlay</name>
    <system>
      <services>
        <ssh></ssh>
      </services>
    </system>
  </groups>
  <groups>
    <name>unused</name>
    <system>
      <services>
        <telnet></telnet>
      </services>
    </system>
  </groups>
</configuration>
"""

        with tempfile.NamedTemporaryFile("w", suffix=".xml", delete=False) as handle:
            handle.write(xml_content)
            xml_path = handle.name

        try:
            rendered = module.parse_xml_to_hcl(xml_path, "vmx", "router1", type_lookup)
        finally:
            os.remove(xml_path)

        self.assertIsNotNone(rendered)
        self.assertEqual(rendered.count('resource "terraform-provider-junos-vmx"'), 1)
        self.assertIn("router1-base-config", rendered)
        self.assertNotIn("apply-groups", rendered)
        self.assertIn("ssh = [", rendered)
        self.assertNotIn("unused", rendered)
        self.assertNotIn("telnet", rendered)

    def test_groups_are_kept_when_requested(self):
        """--groups converts a group-based configuration as it stands."""
        type_lookup = {
            "system": {"type": "container"},
            "system/host_name": {"type": "leaf"},
            "system/services": {"type": "container"},
            "system/services/ssh": {"type": "container"},
            "apply_groups": {"type": "leaf-list"},
            "groups": {"type": "list", "key": "name"},
            "groups/name": {"type": "leaf"},
            "groups/system": {"type": "container"},
            "groups/system/services": {"type": "container"},
            "groups/system/services/ssh": {"type": "container"},
        }

        xml_content = """<configuration>
  <apply-groups>overlay</apply-groups>
  <system>
    <host-name>router1</host-name>
  </system>
  <groups>
    <name>overlay</name>
    <system>
      <services>
        <ssh></ssh>
      </services>
    </system>
  </groups>
</configuration>
"""
        with tempfile.NamedTemporaryFile("w", suffix=".xml", delete=False) as handle:
            handle.write(xml_content)
            xml_path = handle.name

        try:
            # A fresh module each time: the converter remembers paths it has
            # warned about and dropped.
            kept = _load_xml2tf().parse_xml_to_hcl(
                xml_path, "vmx", "router1", type_lookup, True)
            flattened = _load_xml2tf().parse_xml_to_hcl(
                xml_path, "vmx", "router1", type_lookup)
        finally:
            os.remove(xml_path)

        self.assertIsNotNone(kept)
        # The group survives with its name and its body.
        self.assertIn("groups = [", kept)
        self.assertIn("overlay", kept)
        self.assertIn("apply_groups", kept)
        # and nothing from it is copied into the base hierarchy.
        base = kept.split("groups = [")[0]
        self.assertNotIn("ssh", base)

        # The default is unchanged: the group body is merged in and the group
        # itself is gone.
        self.assertNotIn("groups = [", flattened)
        self.assertNotIn("apply_groups", flattened)
        self.assertIn("ssh = [", flattened)

    def test_common_tf_created_with_locals_for_shared_config(self):
        repo_root = os.path.abspath(os.path.join(os.path.dirname(__file__), "..", ".."))
        script_path = os.path.join(repo_root, "junosterraform", "jtaf-xml2tf")

        loader = SourceFileLoader("jtaf_xml2tf", script_path)
        spec = importlib.util.spec_from_loader("jtaf_xml2tf", loader)
        assert spec is not None
        module = importlib.util.module_from_spec(spec)
        assert spec.loader is not None
        spec.loader.exec_module(module)

        # Synthetic parsed dicts — no XML or type_lookup needed.
        # "shared_key" is identical on both devices; "host_name" is device-specific.
        parsed_by_host = {
            "device1": {"host_name": "device1", "shared_key": [{"value": "same"}]},
            "device2": {"host_name": "device2", "shared_key": [{"value": "same"}]},
        }

        shared_locals, extraction_points, group_members = module.walk_and_extract(parsed_by_host)

        with tempfile.TemporaryDirectory() as output_dir:
            common_tf_path = os.path.join(output_dir, "common.tf")
            with open(common_tf_path, "w") as f:
                f.write(module.generate_locals_block(shared_locals, group_members))

            for hostname, parsed_data in parsed_by_host.items():
                hcl = module.generate_hcl_resources(
                    parsed_data, "vmx", hostname, extraction_points
                )
                with open(os.path.join(output_dir, f"{hostname}.tf"), "w") as f:
                    f.write(hcl)

            self.assertTrue(os.path.exists(common_tf_path))
            self.assertIn("locals {", open(common_tf_path).read())

            for hostname in parsed_by_host:
                device_tf = os.path.join(output_dir, f"{hostname}.tf")
                self.assertTrue(os.path.exists(device_tf))
                self.assertIn("local.", open(device_tf).read())
