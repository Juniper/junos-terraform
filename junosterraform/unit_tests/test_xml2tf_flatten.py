import importlib.util
from importlib.machinery import SourceFileLoader
import os
import tempfile
import unittest


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