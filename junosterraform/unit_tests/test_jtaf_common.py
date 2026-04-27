import json
import os
import tempfile
import unittest
import xml.etree.ElementTree as ElementTree

from junosterraform import jtaf_common


class TestGetXpaths(unittest.TestCase):
    """Tests for get_xpaths function."""

    def test_get_xpaths_simple(self):
        """Test XPath extraction from simple XML tree."""
        xml_string = """<root>
            <interfaces>
                <interface>
                    <name>eth0</name>
                </interface>
            </interfaces>
        </root>"""
        root = ElementTree.fromstring(xml_string)
        result = jtaf_common.get_xpaths(root)

        self.assertIn('', result)  # Root path
        self.assertIn('interfaces', result)
        self.assertIn('interfaces/interface', result)
        self.assertIn('interfaces/interface/name', result)

    def test_get_xpaths_nested(self):
        """Test XPath extraction from nested XML tree."""
        xml_string = """<root>
            <system>
                <hostname>router1</hostname>
                <domain-name>example.com</domain-name>
            </system>
        </root>"""
        root = ElementTree.fromstring(xml_string)
        result = jtaf_common.get_xpaths(root)

        self.assertIn('system', result)
        self.assertIn('system/hostname', result)
        self.assertIn('system/domain-name', result)

    def test_get_xpaths_empty(self):
        """Test XPath extraction from minimal XML tree."""
        root = ElementTree.Element("root")
        result = jtaf_common.get_xpaths(root)

        self.assertIn('', result)


class TestUniqueXpaths(unittest.TestCase):
    """Tests for unique_xpaths function."""

    def test_unique_xpaths_preserves_groups_name(self):
        """Test that 'groups/name' is preserved."""
        paths = {
            'groups/name': True,
            'interfaces': True,
            'system': True
        }
        result = jtaf_common.unique_xpaths(paths)

        self.assertIn('groups/name', result)
        self.assertIn('interfaces', result)
        self.assertIn('system', result)

    def test_unique_xpaths_preserves_groups_prefix(self):
        """Test that 'groups/' prefix is preserved on paths."""
        paths = {
            'groups/interfaces': True,
            'groups/system': True,
            'vlans': True
        }
        result = jtaf_common.unique_xpaths(paths)

        self.assertIn('groups/interfaces', result)
        self.assertIn('groups/system', result)
        self.assertIn('vlans', result)

    def test_unique_xpaths_empty(self):
        """Test unique_xpaths with empty dict."""
        result = jtaf_common.unique_xpaths({})
        self.assertEqual(result, [])

    def test_unique_xpaths_duplicates(self):
        """Test unique_xpaths handles duplicates properly."""
        paths = {
            'interfaces': True,
            'interfaces': True,
            'system': True
        }
        result = jtaf_common.unique_xpaths(paths)

        # Should have unique entries
        self.assertEqual(len(result), 2)
        self.assertIn('interfaces', result)
        self.assertIn('system', result)


class TestGetPath(unittest.TestCase):
    """Tests for get_path function."""

    def test_get_path_simple(self):
        """Test path construction from parent hierarchy."""
        parent = [
            {"name": "root"},
            {"name": "level1"},
            {"name": "level2"},
            {"name": "level3"},
            {"name": "level4"}
        ]
        result = jtaf_common.get_path(parent)

        self.assertEqual(result, "level2/level3/level4")

    def test_get_path_minimum(self):
        """Test path construction with minimum hierarchy."""
        parent = [
            {"name": "root"},
            {"name": "level1"},
            {"name": "level2"}
        ]
        result = jtaf_common.get_path(parent)

        self.assertEqual(result, "level2")

    def test_get_path_no_name_field(self):
        """Test path construction with items missing 'name' field."""
        parent = [
            {"name": "root"},
            {"name": "level1"},
            {"name": "level2"},
            {"text": "no_name_here"},
            {"name": "level4"}
        ]
        result = jtaf_common.get_path(parent)

        # Items without 'name' are skipped
        self.assertIn("level4", result)

    def test_get_path_empty_list(self):
        """Test path construction with empty list."""
        result = jtaf_common.get_path([])
        self.assertEqual(result, "")


class TestCheckPath(unittest.TestCase):
    """Tests for check_path function."""

    def test_check_path_found(self):
        """Test check_path when path is in allowed list."""
        paths = ["interfaces/eth0", "level2/interfaces/eth0"]
        node = [
            {"name": "root"},
            {"name": "level1"},
            {"name": "level2"},
            {"name": "interfaces"},
            {"name": "eth0"}
        ]
        result = jtaf_common.check_path(paths, node)

        self.assertTrue(result)

    def test_check_path_not_found(self):
        """Test check_path when path is not in allowed list."""
        paths = ["interfaces/eth0"]
        node = [
            {"name": "root"},
            {"name": "level1"},
            {"name": "level2"},
            {"name": "system"},
            {"name": "hostname"}
        ]
        result = jtaf_common.check_path(paths, node)

        self.assertFalse(result)

    def test_check_path_empty_path(self):
        """Test check_path returns True for empty computed path."""
        paths = ["something"]
        node = [
            {"name": "root"},
            {"name": "level1"},
            {"name": "level2"}
        ]
        result = jtaf_common.check_path(paths, node)

        # With only 3 elements, get_path returns "level2", which is not in paths,
        # so this should return False
        self.assertFalse(result)


class TestCheckForChoice(unittest.TestCase):
    """Tests for check_for_choice function."""

    def test_check_for_choice_no_choice(self):
        """Test check_for_choice with non-choice element."""
        elem = {
            "type": "leaf",
            "name": "hostname"
        }
        result = jtaf_common.check_for_choice(elem)

        self.assertEqual(result, [])

    def test_check_for_choice_with_children(self):
        """Test check_for_choice with choice element containing children."""
        elem = {
            "type": "choice",
            "name": "choice-elem",
            "children": [
                {
                    "name": "case1",
                    "children": [{"name": "option1"}]
                },
                {
                    "name": "case2",
                    "children": [{"name": "option2"}]
                }
            ]
        }
        result = jtaf_common.check_for_choice(elem)

        self.assertEqual(len(result), 2)
        self.assertIn([{"name": "option1"}], result)
        self.assertIn([{"name": "option2"}], result)

    def test_check_for_choice_no_children(self):
        """Test check_for_choice with choice element but no children."""
        elem = {
            "type": "choice",
            "name": "choice-elem"
        }
        result = jtaf_common.check_for_choice(elem)

        self.assertEqual(result, [])

    def test_check_for_choice_children_without_subchildren(self):
        """Test check_for_choice with children lacking 'children' field."""
        elem = {
            "type": "choice",
            "name": "choice-elem",
            "children": [
                {"name": "case1"},
                {"name": "case2"}
            ]
        }
        result = jtaf_common.check_for_choice(elem)

        self.assertEqual(result, [])


class TestCheckForEnums(unittest.TestCase):
    """Tests for check_for_enums function."""

    def test_check_for_enums_choice_ident_leaf(self):
        """Test check_for_enums with choice-ident leaf element."""
        elem = {
            "type": "leaf",
            "name": "choice-ident",
            "enums": [
                {"id": "option1"},
                {"id": "option2"}
            ]
        }
        node_parent = [
            {},
            {"key": ["existing_key"]},
            {}
        ]
        result = jtaf_common.check_for_enums(elem, node_parent)

        self.assertEqual(len(result), 2)
        self.assertEqual(result[0]["name"], "option1")
        self.assertEqual(result[0]["type"], "leaf")
        self.assertEqual(result[0]["leaf-type"], "string")

    def test_check_for_enums_not_choice_ident(self):
        """Test check_for_enums with non-choice-ident element."""
        elem = {
            "type": "leaf",
            "name": "regular-leaf",
            "enums": [{"id": "option1"}]
        }
        node_parent = [
            {},
            {"key": ["existing_key"]},
            {}
        ]
        result = jtaf_common.check_for_enums(elem, node_parent)

        self.assertEqual(result, [])

    def test_check_for_enums_filters_by_parent_key(self):
        """Test check_for_enums filters enums that exist in parent key."""
        elem = {
            "type": "leaf",
            "name": "choice-ident",
            "enums": [
                {"id": "in_key"},
                {"id": "not_in_key"}
            ]
        }
        node_parent = [
            {},
            {"key": ["in_key"]},
            {}
        ]
        result = jtaf_common.check_for_enums(elem, node_parent)

        # only "not_in_key" should be included
        self.assertEqual(len(result), 1)
        self.assertEqual(result[0]["name"], "not_in_key")


class TestCalcElemPath(unittest.TestCase):
    """Tests for _calc_elem_path function."""

    def test_calc_elem_path_with_current_path(self):
        """Test element path calculation with existing current_path."""
        elem = {"name": "eth0"}
        result = jtaf_common._calc_elem_path("interfaces", elem)

        self.assertEqual(result, "interfaces/eth0")

    def test_calc_elem_path_empty_current(self):
        """Test element path calculation with empty current_path."""
        elem = {"name": "interfaces"}
        result = jtaf_common._calc_elem_path("", elem)

        self.assertEqual(result, "interfaces")

    def test_calc_elem_path_non_dict(self):
        """Test element path calculation with non-dict input."""
        result = jtaf_common._calc_elem_path("path", "string_not_dict")

        self.assertEqual(result, "")

    def test_calc_elem_path_missing_name(self):
        """Test element path calculation with missing 'name' key."""
        elem = {"type": "leaf"}
        result = jtaf_common._calc_elem_path("interfaces", elem)

        self.assertEqual(result, "interfaces/")


class TestMatchChoices(unittest.TestCase):
    """Tests for _match_choices function."""

    def test_match_choices_found(self):
        """Test _match_choices when matching choices exist."""
        paths = ["interfaces/eth0"]
        choices = [
            [{"name": "eth0"}],
            [{"name": "eth1"}]
        ]
        result = jtaf_common._match_choices(paths, "interfaces", choices)

        self.assertEqual(len(result), 1)
        self.assertIn([{"name": "eth0"}], result)

    def test_match_choices_not_found(self):
        """Test _match_choices when no matching choices exist."""
        paths = ["system/hostname"]
        choices = [
            [{"name": "eth0"}],
            [{"name": "eth1"}]
        ]
        result = jtaf_common._match_choices(paths, "interfaces", choices)

        self.assertEqual(result, [])

    def test_match_choices_multiple_matches(self):
        """Test _match_choices with multiple matching choices."""
        paths = ["interfaces/eth0", "interfaces/eth1"]
        choices = [
            [{"name": "eth0"}],
            [{"name": "eth1"}]
        ]
        result = jtaf_common._match_choices(paths, "interfaces", choices)

        self.assertEqual(len(result), 2)


class TestMatchEnums(unittest.TestCase):
    """Tests for _match_enums function."""

    def test_match_enums_found(self):
        """Test _match_enums when matching enums exist."""
        paths = ["protocol/tcp"]
        enums = [
            {"name": "tcp"},
            {"name": "udp"}
        ]
        result = jtaf_common._match_enums(paths, "protocol", enums)

        self.assertEqual(len(result), 1)
        self.assertIn({"name": "tcp"}, result)

    def test_match_enums_not_found(self):
        """Test _match_enums when no matching enums exist."""
        paths = ["protocol/icmp"]
        enums = [
            {"name": "tcp"},
            {"name": "udp"}
        ]
        result = jtaf_common._match_enums(paths, "protocol", enums)

        self.assertEqual(result, [])

    def test_match_enums_multiple_matches(self):
        """Test _match_enums with multiple matching enums."""
        paths = ["protocol/tcp", "protocol/udp"]
        enums = [
            {"name": "tcp"},
            {"name": "udp"},
            {"name": "icmp"}
        ]
        result = jtaf_common._match_enums(paths, "protocol", enums)

        self.assertEqual(len(result), 2)


class TestCheckChildren(unittest.TestCase):
    """Tests for check_children function."""

    def test_check_children_invalid_parent(self):
        """Test check_children when parent is not a dict."""
        paths = ["interfaces/eth0"]
        elem = {"name": "eth0"}
        # node_parent needs at least 2 elements, with [-2] being non-dict
        node_parent = [{"name": "root"}, "not_dict"]

        result = jtaf_common.check_children(paths, elem, node_parent, "interfaces")

        self.assertTrue(result)

    def test_check_children_direct_match(self):
        """Test check_children when element path matches directly."""
        paths = ["configuration/interfaces"]
        elem = {"name": "interfaces"}
        node_parent = [{}, {}]

        result = jtaf_common.check_children(paths, elem, node_parent, "configuration")

        self.assertTrue(result)

    def test_check_children_configuration_match(self):
        """Test check_children always returns True for configuration element."""
        paths = ["something/else"]
        elem = {"name": "configuration"}
        node_parent = [{}, {}]

        result = jtaf_common.check_children(paths, elem, node_parent, "")

        self.assertTrue(result)

    def test_check_children_no_children_in_parent(self):
        """Test check_children when parent has no children."""
        paths = ["interfaces"]
        elem = {"name": "eth0"}
        node_parent = [{"name": "interfaces"}, {}]

        result = jtaf_common.check_children(paths, elem, node_parent, "interfaces")

        self.assertTrue(result)


class TestRemoveTagsByName(unittest.TestCase):
    """Tests for remove_tags_by_name function."""

    def test_remove_tags_by_name_single_tag(self):
        """Test removing a single tag type."""
        xml_string = """<root>
            <system>
                <hostname>router1</hostname>
                <version>1.0</version>
            </system>
        </root>"""
        root = ElementTree.fromstring(xml_string)

        jtaf_common.remove_tags_by_name(root, ["version"])

        version = root.find(".//version")
        self.assertIsNone(version)
        hostname = root.find(".//hostname")
        self.assertIsNotNone(hostname)

    def test_remove_tags_by_name_multiple_tags(self):
        """Test removing multiple tag types."""
        xml_string = """<root>
            <system>
                <hostname>router1</hostname>
                <version>1.0</version>
                <model>vsrx</model>
            </system>
        </root>"""
        root = ElementTree.fromstring(xml_string)

        jtaf_common.remove_tags_by_name(root, ["version", "model"])

        version = root.find(".//version")
        model = root.find(".//model")
        hostname = root.find(".//hostname")

        self.assertIsNone(version)
        self.assertIsNone(model)
        self.assertIsNotNone(hostname)

    def test_remove_tags_by_name_nested(self):
        """Test removing tags from nested elements."""
        xml_string = """<root>
            <system>
                <hostname>router1</hostname>
                <versions>
                    <version>1.0</version>
                </versions>
            </system>
        </root>"""
        root = ElementTree.fromstring(xml_string)

        jtaf_common.remove_tags_by_name(root, ["version"])

        version = root.find(".//version")
        self.assertIsNone(version)

    def test_remove_tags_by_name_nonexistent_tag(self):
        """Test removing a tag that doesn't exist."""
        xml_string = """<root>
            <system>
                <hostname>router1</hostname>
            </system>
        </root>"""
        root = ElementTree.fromstring(xml_string)

        # Should not raise an error
        jtaf_common.remove_tags_by_name(root, ["nonexistent"])

        hostname = root.find(".//hostname")
        self.assertIsNotNone(hostname)


class TestFindParent(unittest.TestCase):
    """Tests for find_parent function."""

    def test_find_parent_direct_child(self):
        """Test finding parent of direct child."""
        xml_string = """<root>
            <system>
                <hostname>router1</hostname>
            </system>
        </root>"""
        root = ElementTree.fromstring(xml_string)
        system = root.find("system")

        parent = jtaf_common.find_parent(root, system)

        self.assertEqual(parent, root)

    def test_find_parent_nested_child(self):
        """Test finding parent of nested child."""
        xml_string = """<root>
            <system>
                <hostname>router1</hostname>
            </system>
        </root>"""
        root = ElementTree.fromstring(xml_string)
        hostname = root.find(".//hostname")
        system = root.find("system")

        parent = jtaf_common.find_parent(root, hostname)

        self.assertEqual(parent, system)

    def test_find_parent_not_found(self):
        """Test find_parent when child is not in tree."""
        xml_string = "<root><system><hostname>router1</hostname></system></root>"
        root = ElementTree.fromstring(xml_string)
        other_elem = ElementTree.Element("other")

        parent = jtaf_common.find_parent(root, other_elem)

        self.assertIsNone(parent)

    def test_find_parent_root_element(self):
        """Test find_parent when looking for root's parent."""
        root = ElementTree.Element("root")

        parent = jtaf_common.find_parent(root, root)

        self.assertIsNone(parent)


class TestWalkSchema(unittest.TestCase):
    """Tests for walk_schema function."""

    def test_walk_schema_dict(self):
        """Test walk_schema with dict node."""
        paths = ["interfaces"]
        schema = {
            "name": "configuration",
            "children": {
                "interfaces": {"name": "interfaces"}
            }
        }

        result = jtaf_common.walk_schema(paths, schema)

        self.assertIsInstance(result, dict)

    def test_walk_schema_list(self):
        """Test walk_schema with list node."""
        paths = ["interfaces"]
        schema = [
            {"name": "interfaces"},
            {"name": "system"}
        ]
        # Provide a proper parent structure with at least 2 elements
        parent = [{"name": "root"}, {"name": "level1"}]

        result = jtaf_common.walk_schema(paths, schema, parent)

        self.assertIsInstance(result, list)

    def test_walk_schema_scalar(self):
        """Test walk_schema with scalar value."""
        paths = []
        schema = "string_value"

        result = jtaf_common.walk_schema(paths, schema)

        self.assertEqual(result, "string_value")


class TestFilterJsonUsingXml(unittest.TestCase):
    """Tests for filter_json_using_xml function."""

    def test_filter_json_using_xml_with_element(self):
        """Test filter_json_using_xml with ElementTree element."""
        schema_dict = {
            "name": "configuration",
            "children": {}
        }

        xml_string = """<configuration>
            <system>
                <hostname>router1</hostname>
            </system>
        </configuration>"""

        # Create XML element at configuration level
        root = ElementTree.fromstring(xml_string)

        with tempfile.NamedTemporaryFile(mode='w', suffix='.json', delete=False) as f:
            json.dump(schema_dict, f)
            schema_file = f.name

        try:
            result = jtaf_common.filter_json_using_xml(schema_file, root)

            self.assertIsInstance(result, dict)
            self.assertEqual(result["name"], "configuration")
        finally:
            os.remove(schema_file)

    def test_filter_json_using_xml_removes_version_tags(self):
        """Test that filter_json_using_xml removes version tags."""
        schema_dict = {
            "name": "configuration",
            "children": {}
        }

        xml_string = """<root>
            <configuration>
                <version>18.2</version>
                <system>
                    <hostname>router1</hostname>
                </system>
            </configuration>
        </root>"""

        root = ElementTree.fromstring(xml_string)

        with tempfile.NamedTemporaryFile(mode='w', suffix='.json', delete=False) as f:
            json.dump(schema_dict, f)
            schema_file = f.name

        try:
            result = jtaf_common.filter_json_using_xml(schema_file, root)

            # The filtering should have removed version tag
            self.assertIsInstance(result, dict)
        finally:
            os.remove(schema_file)


class TestLoadAndMergeXmls(unittest.TestCase):
    """Tests for load_and_merge_xmls function."""

    def test_load_and_merge_xmls_single_file(self):
        """Test merging a single XML file."""
        xml_content = """<root>
            <configuration>
                <system>
                    <hostname>router1</hostname>
                </system>
            </configuration>
        </root>"""

        with tempfile.NamedTemporaryFile(mode='w', suffix='.xml', delete=False) as f:
            f.write(xml_content)
            xml_file = f.name

        try:
            result = jtaf_common.load_and_merge_xmls([xml_file])

            self.assertEqual(result.tag, "configuration")
            system = result.find("system")
            self.assertIsNotNone(system)
        finally:
            os.remove(xml_file)

    def test_load_and_merge_xmls_multiple_files(self):
        """Test merging multiple XML files."""
        xml_content1 = """<root>
            <configuration>
                <system>
                    <hostname>router1</hostname>
                </system>
            </configuration>
        </root>"""

        xml_content2 = """<root>
            <configuration>
                <interfaces>
                    <interface>eth0</interface>
                </interfaces>
            </configuration>
        </root>"""

        with tempfile.NamedTemporaryFile(mode='w', suffix='.xml', delete=False) as f1:
            f1.write(xml_content1)
            xml_file1 = f1.name

        with tempfile.NamedTemporaryFile(mode='w', suffix='.xml', delete=False) as f2:
            f2.write(xml_content2)
            xml_file2 = f2.name

        try:
            result = jtaf_common.load_and_merge_xmls([xml_file1, xml_file2])

            self.assertEqual(result.tag, "configuration")
            system = result.find("system")
            interfaces = result.find("interfaces")
            self.assertIsNotNone(system)
            self.assertIsNotNone(interfaces)
        finally:
            os.remove(xml_file1)
            os.remove(xml_file2)

    def test_load_and_merge_xmls_missing_configuration(self):
        """Test load_and_merge_xmls with file missing configuration element."""
        xml_content = """<root>
            <system>
                <hostname>router1</hostname>
            </system>
        </root>"""

        with tempfile.NamedTemporaryFile(mode='w', suffix='.xml', delete=False) as f:
            f.write(xml_content)
            xml_file = f.name

        try:
            with self.assertRaises(ValueError):
                jtaf_common.load_and_merge_xmls([xml_file])
        finally:
            os.remove(xml_file)

    def test_load_and_merge_xmls_nonexistent_file(self):
        """Test load_and_merge_xmls with nonexistent file."""
        with self.assertRaises(FileNotFoundError):
            jtaf_common.load_and_merge_xmls(["/nonexistent/path/file.xml"])


class TestIntegration(unittest.TestCase):
    """Integration tests combining multiple functions."""

    def test_xpath_extraction_and_filtering(self):
        """Test XPath extraction followed by unique filtering."""
        xml_string = """<root>
            <interfaces>
                <interface>
                    <name>eth0</name>
                </interface>
            </interfaces>
            <system>
                <hostname>router1</hostname>
            </system>
        </root>"""

        root = ElementTree.fromstring(xml_string)
        xpaths = jtaf_common.get_xpaths(root)
        unique = jtaf_common.unique_xpaths(xpaths)

        self.assertIn('interfaces', unique)
        self.assertIn('interfaces/interface', unique)
        self.assertIn('system', unique)

    def test_path_validation_workflow(self):
        """Test complete path validation workflow."""
        paths = ["level2/interfaces/eth0", "level2/system/hostname"]

        # Valid path - matches "level2/interfaces/eth0"
        node1 = [
            {"name": "root"},
            {"name": "level1"},
            {"name": "level2"},
            {"name": "interfaces"},
            {"name": "eth0"}
        ]

        # Invalid path - doesn't match any allowed paths
        node2 = [
            {"name": "root"},
            {"name": "level1"},
            {"name": "level2"},
            {"name": "interfaces"},
            {"name": "eth1"}
        ]

        self.assertTrue(jtaf_common.check_path(paths, node1))
        self.assertFalse(jtaf_common.check_path(paths, node2))


if __name__ == "__main__":
    unittest.main()
